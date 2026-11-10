package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/chubuk/mcp-is/internal/ratelimit"
	"github.com/chubuk/mcp-is/internal/resilience"
	"github.com/chubuk/mcp-is/internal/service"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// Chain-тесты: MCP-клиент → реальные handler'ы → service → InSales-клиент →
// httptest-мок. Закрывают пробелы «нет контрактов create-тулов против моков
// InSales» и «e2e не вызывает живые handler'ы». Правила:
//   - MCP_TOOLS_MODE=all (create_* живут вне core-поверхности);
//   - не вызывать store_info/get_shop_info — fetchStoreTitle ходит в сеть мимо HTTPClient.

type recordedCall struct {
	Method string
	Path   string // r.URL.Path без query
	Body   map[string]any
}

type chainMock struct {
	mu      sync.Mutex
	calls   []recordedCall
	respond func(recordedCall) string
	server  *httptest.Server
}

func newChainMock(t *testing.T, respond func(recordedCall) string) *chainMock {
	t.Helper()
	m := &chainMock{respond: respond}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&body)
		}
		rec := recordedCall{Method: r.Method, Path: r.URL.Path, Body: body}
		m.mu.Lock()
		m.calls = append(m.calls, rec)
		out := m.respond(rec)
		m.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(out))
	}))
	t.Cleanup(m.server.Close)
	return m
}

func (m *chainMock) last() recordedCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.calls) == 0 {
		return recordedCall{}
	}
	return m.calls[len(m.calls)-1]
}

func (m *chainMock) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.calls)
}

func jsonResponder(payload string) func(recordedCall) string {
	return func(recordedCall) string { return payload }
}

// fakeClientProvider — шов service.ClientProvider: возвращает клиент,
// направленный на httptest-мок (паттерн newTestClient из insales-тестов).
type fakeClientProvider struct{ c *insales.InSalesClient }

func (f fakeClientProvider) GetClient(ctx context.Context) (*insales.InSalesClient, context.Context, error) {
	ctx = context.WithValue(ctx, core.ShopReadOnlyKey, false)
	return f.c, ctx, nil
}

func newTestInSalesClient(srv *httptest.Server) *insales.InSalesClient {
	return &insales.InSalesClient{
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
		CircuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
		RateTracker:    ratelimit.NewRateTracker(500, 5*time.Minute),
		BreakerPool:    ratelimit.NewPerKeyBreaker(5, 30*time.Second),
	}
}

type chainHarness struct {
	srv *MCPServer
	cl  *client.Client
}

func newChainHarness(t *testing.T, mock *chainMock) *chainHarness {
	t.Helper()
	t.Setenv("MCP_TOOLS_MODE", "all")
	provider := fakeClientProvider{c: newTestInSalesClient(mock.server)}
	srv := NewMCPServer(service.NewService(provider), false, false, "RU", nil, nil)
	cl, err := client.NewInProcessClient(srv.mcpServer)
	if err != nil {
		t.Fatalf("in-process client: %v", err)
	}
	e2eStart(t, cl)
	return &chainHarness{srv: srv, cl: cl}
}

func callTool(t *testing.T, h *chainHarness, ctx context.Context, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	res, err := h.cl.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{Name: name, Arguments: args},
	})
	if err != nil {
		t.Fatalf("%s call: %v", name, err)
	}
	if res.IsError {
		t.Fatalf("%s errored: %+v", name, res.Content)
	}
	return res
}

// fillRequiredArgs добирает обязательные аргументы из inputSchema тула:
// строки → "1" (или первый enum), числа → 1, boolean → true. Enum берётся
// первым допустым значением — иначе валидация схемы отклонит вызов.
func fillRequiredArgs(t *testing.T, h *chainHarness, name string) map[string]any {
	t.Helper()
	st := h.srv.mcpServer.GetTool(name)
	if st == nil {
		t.Fatalf("tool %q not registered", name)
	}
	args := map[string]any{}
	for _, req := range st.Tool.InputSchema.Required {
		if _, ok := args[req]; ok {
			continue
		}
		def, _ := st.Tool.InputSchema.Properties[req].(map[string]any)
		typ, _ := def["type"].(string)
		switch typ {
		case "number", "integer":
			args[req] = float64(1)
		case "boolean":
			args[req] = true
		case "array":
			args[req] = []any{}
		case "object":
			args[req] = map[string]any{}
		default:
			if enum, ok := def["enum"].([]string); ok && len(enum) > 0 {
				args[req] = enum[0]
			} else {
				args[req] = "1"
			}
		}
	}
	return args
}

type createCall struct {
	name     string
	args     map[string]any // поверх авто-fill обязательных
	wantPath string
	wantRoot string   // "" = тело без корня (плоское)
	wantKeys []string // ключи, которые обязаны быть в корне (или в теле при root="")
}

func runCreateCalls(t *testing.T, calls []createCall) {
	t.Helper()
	for _, call := range calls {
		t.Run(call.name, func(t *testing.T) {
			mock := newChainMock(t, jsonResponder(`{"id":42}`))
			h := newChainHarness(t, mock)
			args := fillRequiredArgs(t, h, call.name)
			for k, v := range call.args {
				args[k] = v
			}
			callTool(t, h, context.Background(), call.name, args)

			rec := mock.last()
			if rec.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", rec.Method)
			}
			if rec.Path != call.wantPath {
				t.Errorf("path = %s, want %s", rec.Path, call.wantPath)
			}
			if _, ok := rec.Body["idempotency_key"]; ok {
				t.Error("idempotency_key leaked into request body")
			}
			root := rec.Body
			if call.wantRoot != "" {
				r, ok := rec.Body[call.wantRoot].(map[string]any)
				if !ok {
					t.Fatalf("body root %q missing or not object: %v", call.wantRoot, rec.Body)
				}
				if _, ok := r["idempotency_key"]; ok {
					t.Error("idempotency_key leaked into request body root")
				}
				root = r
			}
			for _, k := range call.wantKeys {
				if _, ok := root[k]; !ok {
					t.Errorf("body key %q missing: %v", k, root)
				}
			}
		})
	}
}

// genericCreateCalls — create-тулы на genericCreateHandler (один путь кода:
// args → svc.CreateXxx → клиент). Путь и root-ключ тела — контракт с InSales.
var genericCreateCalls = []createCall{
	{"create_delivery_variant", nil, "/delivery_variants.json", "delivery_variant", []string{"title"}},
	{"create_payment_gateway", nil, "/payment_gateways.json", "payment_gateway", []string{"title"}},
	{"create_domain", nil, "/domains.json", "domain", []string{"domain"}},
	{"create_client_group", nil, "/client_groups.json", "client_group", []string{"title"}},
	{"create_custom_status", nil, "/custom_statuses.json", "custom_status", []string{"title"}},
	{"create_variant_field", nil, "/variant_fields.json", "variant_field", []string{"title"}},
	{"create_collection_field", nil, "/collection_fields.json", "collection_field", []string{"title"}},
	{"create_collection_filter", nil, "/collection_filters.json", "filter", []string{"title"}},
	{"create_js_tag", nil, "/js_tags.json", "js_tag", []string{"title"}},
	{"create_external_discount", map[string]any{"description": "d"}, "/external_discounts.json", "external_discount", []string{"description"}},
	{"create_discount_collections_group", nil, "/discount_collections_groups.json", "discount_collections_group", []string{"title"}},
	{"create_warehouse", nil, "/warehouses.json", "warehouse", []string{"title"}},
	{"create_sales_channel", nil, "/sales_channels.json", "sales_channel", []string{"title"}},
	{"create_stock_currency", nil, "/stock_currencies.json", "stock_currency", []string{"title"}},
	{"create_price_kind", nil, "/price_kinds.json", "price_kind", []string{"title"}},
	{"create_option_name", nil, "/option_names.json", "option_name", []string{"title"}},
	{"create_application_action", nil, "/application_actions.json", "application_action", []string{"title"}},
	{"create_application_charge", nil, "/application_charges.json", "application_charge", []string{"name"}},
	{"create_application_widget", nil, "/application_widgets.json", "application_widget", []string{"title"}},
	{"create_recurring_application_charge", nil, "/recurring_application_charge.json", "recurring_application_charge", []string{"name"}},
	{"create_oauth_app", nil, "/o_auth_apps.json", "o_auth_app", []string{"title"}},
	{"create_user_task", nil, "/user_tasks.json", "user_task", []string{"title"}},
	{"create_supplementary", nil, "/products/1/supplementaries.json", "", []string{"supplementary_ids"}},
	{"create_similar", nil, "/products/1/similars.json", "", []string{"similar_ids"}},
	{"create_outlet", nil, "/outlets.json", "outlet", []string{"title"}},
	{"create_marketplace", nil, "/marketplaces.json", "marketplace", []string{"name"}},
}

// customCreateCalls — create-тулы с кастомными handler'ами (позиционные
// аргументы, особые тела). Особые пути и переименования полей зафиксированы
// литералами — дрейф обязан падать тестом.
var customCreateCalls = []createCall{
	{"create_product", nil, "/products.json", "product", []string{"title", "variants_attributes"}},
	{"create_order", map[string]any{"client_id": float64(7)}, "/orders.json", "order", []string{"client_id"}},
	{"create_client", nil, "/clients.json", "client", []string{"name", "surname"}},
	{"create_category", nil, "/categories.json", "category", []string{"title"}},
	{"create_collection", nil, "/collections.json", "collection", []string{"title"}},
	{"create_webhook", nil, "/webhooks.json", "webhook", []string{"topic", "address", "format_type"}},
	{"create_page", nil, "/pages.json", "page", []string{"title", "content"}},
	{"create_accessory", nil, "/accessories.json", "accessory", []string{"name"}},
	{"create_property_value", nil, "/accessories/1/values.json", "accessory_value", []string{"name"}},
	{"create_menu", nil, "/menus.json", "menu", []string{"name", "handle"}},
	{"create_menu_item", map[string]any{"related_object_id": "Blog_1"}, "/menus/1/items.json", "menu_item", []string{"related_object_id"}},
	{"create_redirect", nil, "/redirects.json", "redirect", []string{"url", "target"}},
	{"create_blog", nil, "/blogs.json", "blog", []string{"title", "handle"}},
	{"create_blog_post", nil, "/blogs/1/articles.json", "article", []string{"title", "content", "published_at", "author"}},
	{"create_discount_code", nil, "/discount_codes.json", "discount_code", nil},
	// add_* (не create_*): имя тулза — add_product_to_collection, строка нужна
	// для контракта POST /collects.json; в инвариант «51 create_*» не входит.
	{"add_product_to_collection", nil, "/collects.json", "collect", []string{"collection_id", "product_id"}},
	{"create_property", nil, "/properties.json", "property", nil},
	{"create_characteristic", nil, "/properties/1/characteristics.json", "characteristic", nil},
	{"create_review", nil, "/reviews.json", "review", nil},
	{"create_file", nil, "/files.json", "file", nil},
	{"create_product_field", map[string]any{"handle": "fld1", "type": "ProductField::TextField"}, "/product_fields.json", "product_field", []string{"handle", "type"}},
	{"create_product_field_value", map[string]any{"product_field_id": "1"}, "/products/1/product_field_values.json", "product_field_value", []string{"product_field_id"}},
	{"create_variant_field_value", map[string]any{"variant_field_id": "1"}, "/variants/1/variant_field_values.json", "variant_field_value", []string{"variant_field_id"}},
	{"create_option_value", nil, "/option_names/1/option_values.json", "option_value", []string{"title"}},
	{"create_bonus_system_transaction", nil, "/clients/1/bonus_system_transactions.json", "bonus_system_transaction", []string{"bonus_points"}},
	{"create_provider_token", nil, "/dialogs/provider_tokens.json", "", []string{"messenger_type"}},
}

func TestChainCreateTools_Generic(t *testing.T) {
	runCreateCalls(t, genericCreateCalls)
}

func TestChainCreateTools_CustomHandlers(t *testing.T) {
	runCreateCalls(t, customCreateCalls)
}

// TestChainCreateTools_CoverAllRegistered — каждые create_* из реестра имеют
// ровно одну строку в таблицах; дрейф каталога ловится здесь, а не жёстким len.
func TestChainCreateTools_CoverAllRegistered(t *testing.T) {
	covered := map[string]bool{}
	coveredCreates := 0
	for _, call := range append(append([]createCall{}, genericCreateCalls...), customCreateCalls...) {
		if covered[call.name] {
			t.Errorf("duplicate chain row for %q", call.name)
		}
		covered[call.name] = true
		if strings.HasPrefix(call.name, "create_") {
			coveredCreates++
		}
	}
	missing := 0
	for name := range toolMeta {
		if !strings.HasPrefix(name, "create_") {
			continue
		}
		if !covered[name] {
			t.Errorf("create tool %q has no chain row", name)
			missing++
		}
	}
	if missing == 0 && coveredCreates != 51 {
		t.Errorf("create_* chain rows = %d, want 51", coveredCreates)
	}
}

// TestChainReadTools_ResponseMapping — маппинг ответов InSales в structuredContent.
func TestChainReadTools_ResponseMapping(t *testing.T) {
	mock := newChainMock(t, func(rec recordedCall) string {
		if rec.Method == http.MethodGet && strings.HasPrefix(rec.Path, "/products/") {
			return `{"id":1,"title":"t"}`
		}
		return `[{"id":1,"title":"t"}]`
	})
	h := newChainHarness(t, mock)

	res := callTool(t, h, context.Background(), "get_product", map[string]any{"product_id": "1"})
	sc, _ := res.StructuredContent.(map[string]any)
	if sc["title"] != "t" {
		t.Errorf("get_product structured = %v, want title=t", sc)
	}
	if rec := mock.last(); rec.Path != "/products/1.json" || rec.Method != http.MethodGet {
		t.Errorf("get_product request = %s %s", rec.Method, rec.Path)
	}

	res = callTool(t, h, context.Background(), "list_products", map[string]any{})
	sc, _ = res.StructuredContent.(map[string]any)
	items, _ := sc["items"].([]any)
	if len(items) != 1 {
		t.Errorf("list_products items = %v, want 1", sc["items"])
	}
	if total, _ := sc["total"].(int); total != 1 {
		if totalF, _ := sc["total"].(float64); totalF != 1 {
			t.Errorf("list_products total = %v, want 1", sc["total"])
		}
	}
	if rec := mock.last(); rec.Path != "/products.json" {
		t.Errorf("list_products path = %s", rec.Path)
	}
}

// TestChainUpdateProduct_DiffPath — update идёт с pre-fetch GET и возвращает
// changed_fields {old,new} (цепочка updateWithDiff → diffFromArgs).
func TestChainUpdateProduct_DiffPath(t *testing.T) {
	mock := newChainMock(t, func(rec recordedCall) string {
		if rec.Method == http.MethodGet {
			return `{"id":1,"title":"old"}`
		}
		return `{"id":1,"title":"new"}`
	})
	h := newChainHarness(t, mock)

	res := callTool(t, h, context.Background(), "update_product", map[string]any{
		"product_id": "1",
		"title":      "new",
	})
	if mock.count() != 2 {
		t.Fatalf("want pre-fetch GET + PUT, got %d requests: %+v", mock.count(), mock.calls)
	}
	get, put := mock.calls[0], mock.calls[1]
	if get.Method != http.MethodGet || get.Path != "/products/1.json" {
		t.Errorf("pre-fetch = %s %s", get.Method, get.Path)
	}
	if put.Method != http.MethodPut || put.Path != "/products/1.json" {
		t.Errorf("update = %s %s", put.Method, put.Path)
	}
	sc, _ := res.StructuredContent.(map[string]any)
	cf, _ := sc["changed_fields"].(map[string]any)
	ch, _ := cf["title"].(map[string]any)
	if ch == nil {
		t.Fatalf("changed_fields.title missing: %v", sc)
	}
	if ch["old"] != "old" || ch["new"] != "new" {
		t.Errorf("changed_fields.title = %v, want old=new?", ch)
	}
}

// TestChainUpdateBlogPost_PayloadMapping — body_html уезжает как content,
// image_url → image_attributes{src}, роутинг-id (blog_id/post_id) не утекают
// в article, published_at проходит как есть (регрессия B1/B2/B6).
func TestChainUpdateBlogPost_PayloadMapping(t *testing.T) {
	mock := newChainMock(t, func(rec recordedCall) string {
		if rec.Method == http.MethodGet {
			return `{"id":7,"title":"old","content":"<p>old</p>"}`
		}
		return `{"id":7,"title":"old","content":"<p>new</p>"}`
	})
	h := newChainHarness(t, mock)

	callTool(t, h, context.Background(), "update_blog_post", map[string]any{
		"blog_id":      "1",
		"post_id":      "7",
		"body_html":    "<p>new</p>",
		"published_at": "2026-12-01 12:00:00",
		"image_url":    "https://shop.example/img.jpg",
	})
	put := mock.calls[len(mock.calls)-1]
	if put.Method != http.MethodPut || put.Path != "/blogs/1/articles/7.json" {
		t.Fatalf("update = %s %s", put.Method, put.Path)
	}
	article, ok := put.Body["article"].(map[string]any)
	if !ok {
		t.Fatalf("article root missing: %v", put.Body)
	}
	if article["content"] != "<p>new</p>" {
		t.Errorf("content = %v, want <p>new</p>", article["content"])
	}
	if _, ok := article["body_html"]; ok {
		t.Errorf("body_html leaked into article: %v", article)
	}
	if _, ok := article["blog_id"]; ok {
		t.Errorf("blog_id leaked into article: %v", article)
	}
	if _, ok := article["post_id"]; ok {
		t.Errorf("post_id leaked into article: %v", article)
	}
	if article["published_at"] != "2026-12-01 12:00:00" {
		t.Errorf("published_at = %v", article["published_at"])
	}
	img, _ := article["image_attributes"].(map[string]any)
	if img["src"] != "https://shop.example/img.jpg" {
		t.Errorf("image_attributes = %v, want src", article["image_attributes"])
	}
}

// TestChainIdempotentReplay — полная цепочка с idempotency_key: второй вызов
// не трогает InSales и возвращает сохранённый ответ с idempotent_replay.
func TestChainIdempotentReplay(t *testing.T) {
	mock := newChainMock(t, jsonResponder(`{"id":42}`))
	h := newChainHarness(t, mock)

	store := newFakeStore()
	idemRepo = store
	t.Cleanup(func() { idemRepo = nil })

	ctx := context.WithValue(context.Background(), auth.KeyIDKey, "ai-key-1")
	args := fillRequiredArgs(t, h, "create_warehouse")
	args["idempotency_key"] = "k1"
	callTool(t, h, ctx, "create_warehouse", args)

	args2 := fillRequiredArgs(t, h, "create_warehouse")
	args2["idempotency_key"] = "k1"
	res := callTool(t, h, ctx, "create_warehouse", args2)

	if n := mock.count(); n != 1 {
		t.Fatalf("InSales hits = %d, want 1 (replay must not re-execute)", n)
	}
	sc, _ := res.StructuredContent.(map[string]any)
	if sc["idempotent_replay"] != true {
		t.Errorf("idempotent_replay flag missing: %v", sc)
	}
}

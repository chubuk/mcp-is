package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/pdn"
	"github.com/chubuk/mcp-is/internal/service"
	"github.com/chubuk/mcp-is/internal/storage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	mcpServer     *server.MCPServer
	httpServer    *server.StreamableHTTPServer
	readOnly      bool
	autoPDN       bool
	serverCountry string
}

type mcpContextKey string

const (
	mcpShopNameKey   mcpContextKey = "shop_name"
	mcpKeyNameKey    mcpContextKey = "key_name"
	mcpReadOnlyKey   mcpContextKey = "read_only"
	mcpShopPDNKey    mcpContextKey = "shop_pdn_enabled"
	mcpAIKeyPDNKey   mcpContextKey = "ai_key_pdn_enabled"
	mcpNestedCallKey mcpContextKey = "nested_call" // вызов через call_insales_tool (guard видимости пропускает)
)

// Версия каталога инструментов (W7): бампать при любом изменении схем/состава
// тулов — по ней интеграции понимают, что надо перепроверить свои скрипты.
const (
	catalogVersion = "2026-11-10.1"
	catalogDate    = "2026-11-10"
)

// catalogChanges — рукописный список последних изменений интерфейса (store_info).
var catalogChanges = []string{
	"update_blog_post/create_blog_post: добавлены published_at/author/image_url и SEO-поля; body_html теперь реально сохраняется (→article[content]); published_at при создании больше не «2025-01-01»",
	"field_values ({handle: значение}) у create/update_product, update_variant, create/update_collection, create/update_order, create/update_client — доп. поля магазина; handle'ы — в описании параметра",
	"field_values у create/update_category УБРАН: категории InSales (дерево склада) не хранят доп. поля — используйте коллекции (create/update_collection)",
	"create_order: доп. поля пишутся отдельным PUT после создания (InSales игнорирует их при POST /orders.json) — значения сохраняются",
	"core расширен: товары, заказы и блог видны в tools/list напрямую (49 из N); остальное — search_insales_tools → call_insales_tool",
	"переименования полей приняты через алиасы (title→name, event→topic, url→address, amount→bonus_points, quantity→quantity_at_warehouse0 и др.) — старые имена работают",
	"лишние аргументы больше не проскакивают в API — смотрите ignored_keys; при непустом ignored_keys сверьте схему тула",
	"удалены мёртвые тулы (list/create/update/delete_meta_field, list_currencies, get_supplementary) — используйте list_product_fields / stock_currencies",
}

// filterGlobalReadOnly — флаг read-only всего сервера; заполняется в toolFilter.
var filterGlobalReadOnly bool

// catalogCounts — фактические размеры каталога для этого ctx: видимых и всего.
func catalogCounts(ctx context.Context) (visible, total int) {
	total = len(registeredToolMeta)
	for _, m := range registeredToolMeta {
		if toolHiddenReason(ctx, m.Name, m.IsWrite, filterGlobalReadOnly) == "" {
			visible++
		}
	}
	return visible, total
}

const baseInstructions = `MCP-сервер для магазина InSales (mcp-is).

## Discovery
Default tools/list is a core surface (~49 tools): products, orders, clients, categories, collections, blog, product fields, pipelines, export.
For anything else use search_insales_tools(query) then call_insales_tool(name, args). All InSales tools stay callable.
Set MCP_TOOLS_MODE=all to list every tool at once (more tokens, lower selection accuracy).

## Data conventions
- List results: {items, total, page, per_page}. total = items in this response, NOT the catalog size (catalog size: get_products_count / export_products).
- IDs are strings. Field values include handle/title when available.
- Argument values must be language-invariant: pass API ids/handles/enum codes as-is (do not translate).

## Read-only
In read-only sessions write tools are hidden and call_insales_tool rejects them.

## Rate limits
300 requests/min per API key. Near the limit responses slow down instead of failing; on HTTP 429 wait for Retry-After seconds and retry, then pace requests.

## Pipelines
search_product_details, order_full_details, client_full_profile, product_catalog_search, collection_products — multi-step reads in one call.`

// effectiveAlias implements the read-side default chain (S2.3):
// stored connection_name → slug(shop name) → "mcp-is".
func effectiveAlias(connName, shopName string) string {
	if connName != "" {
		return connName
	}
	if shopName != "" {
		return core.Slugify(shopName)
	}
	return "mcp-is"
}

// buildInstructions prepends the per-connection identity block (first ~512 chars
// are self-contained) to the static server guide.
func buildInstructions(ctx context.Context, base string) string {
	connName, _ := ctx.Value(auth.ConnectionNameKey).(string)
	shopName, _ := ctx.Value(auth.ShopNameKey).(string)
	shopURL, _ := ctx.Value(auth.ShopURLKey).(string)
	shopID, _ := ctx.Value(auth.ShopIDKey).(string)

	var header string
	// Баннер совместимости + честные счётчики (W6/W7): первые строки — то, что
	// агент должен прочитать до работы с данными. Плоско и коротко.
	visible, total := catalogCounts(ctx)
	banner := fmt.Sprintf(
		"Каталог схем %s (%s): настроились раньше — сверьте store_info (catalog_changes) с датой своих скриптов, потом работа с данными.\n"+
			"Показано %d из %d тулов; полный каталог — search_insales_tools → call_insales_tool.\n"+
			"«Невозможно» — только после ошибки API или ignored_keys; иначе порядок: store_info → search → call_insales_tool.\n\n",
		catalogVersion, catalogDate, visible, total)
	if shopID == "" {
		header = "ВНИМАНИЕ: AI-ключ не привязан к магазину — запросы к данным отклоняются. " +
			"Попросите владельца привязать ключ к магазину в личном кабинете (раздел «AI-ключи»).\n\n"
	} else {
		connName = effectiveAlias(connName, shopName)
		header = fmt.Sprintf("Connection: %s\nStore: %q (%s)\nScope: read + write for this store only.\n"+
			"All tools operate exclusively on this store; no store parameter exists.\n"+
			"When the user names this store or %q, use this connection.\n"+
			"If the user names a different store, say so — do not guess and do not write.\n"+
			"Call store_info to confirm the target store before any write.\n\n",
			connName, shopName, shopURL, connName)

		if v, ok := ctx.Value(core.ShopDeleteDeniedKey).(bool); ok && v {
			header += "Deletion is denied on this store: delete_/remove_/decline_* tools are hidden.\n\n"
		}

		// Права ключа (из снимка probe) + прямой REST недоступен.
		if v, ok := ctx.Value(auth.CapabilitiesJSONKey).(string); ok && v != "" {
			snap := core.ParseSnapshot(v)
			if snap != nil {
				summary := core.BuildSummary(snap, nil)
				if len(summary.DeniedDomains) > 0 {
					header += fmt.Sprintf("НЕТ ПРАВ на ресурсы: %s — эти инструменты скрыты. "+
						"Включите доступ в настройках приложения InSales.\n", strings.Join(summary.DeniedDomains, ", "))
				}
				if !summary.CredentialsOK {
					header += "ВНИМАНИЕ: учётные данные InSales не прошли проверку (401).\n"
				}
				header += fmt.Sprintf("Права проверены: %s. Актуальное состояние — через store_info.\n\n", summary.CheckedAt)
			}
		}
		header += "Прямой REST к InSales через этот ключ недоступен — используйте инструменты или call_insales_tool.\n" +
			"При update передавайте только изменяемые поля: пустое значение затирает поле. Для явной очистки — clear_fields.\n\n"
	}
	return banner + header + base
}

func NewMCPServer(svc *service.Service, readOnly, autoPDN bool, serverCountry string, auditRepo *storage.AuditRepository, idemRepo *storage.IdempotencyRepository) *MCPServer {
	hooks := &server.Hooks{}
	hooks.AddAfterInitialize(func(ctx context.Context, _ any, _ *mcp.InitializeRequest, result *mcp.InitializeResult) {
		result.Instructions = buildInstructions(ctx, baseInstructions)
	})

	s := server.NewMCPServer(
		"mcp-is",
		"0.12.0",
		server.WithToolCapabilities(true),
		server.WithToolFilter(toolFilter(readOnly)),
		server.WithRecovery(),
		server.WithHooks(hooks),
		server.WithOutputSchemaValidation(),
		server.WithInstructions(baseInstructions),
	)

	SetIdempotencyRepo(idemRepo)
	RegisterTools(s, svc)
	RegisterDiscoveryTools(s, svc, auditRepo)
	if auditRepo != nil {
		s.Use(auditServerMiddleware(s, auditRepo))
	}

	httpServer := server.NewStreamableHTTPServer(s)

	return &MCPServer{
		mcpServer:     s,
		httpServer:    httpServer,
		readOnly:      readOnly,
		autoPDN:       autoPDN,
		serverCountry: serverCountry,
	}
}

func (s *MCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Limit body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	ctx := r.Context()

	shopName, _ := ctx.Value(core.ShopNameKey).(string)
	keyName, _ := ctx.Value(auth.KeyNameKey).(string)
	aiKeyPDN, _ := ctx.Value(auth.KeyPDNKey).(bool)
	readOnly := s.readOnly
	if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
		readOnly = readOnly || v
	}

	// Auto-PDN logic:
	// 1. If auto-PDN disabled → no protection
	// 2. If server NOT in Russia AND shop IS Russian → protect
	// 3. If server IS in Russia AND shop IS Russian AND LLM is foreign → protect
	shopJurisdiction := "INTL"
	llmJurisdiction := "UNKNOWN"
	if shopName != "" {
		shopJurisdiction = pdn.DetectJurisdiction(shopName)
	}
	if keyName != "" {
		llmJurisdiction = pdn.DetectLLMProvider(keyName)
	}
	shopPDN := false
	if v, ok := ctx.Value(auth.ShopPDNEnabledKey).(bool); ok {
		shopPDN = v // shops.pdn_enabled — принудительная защита независимо от auto-PDN
	}
	if s.autoPDN {
		serverInRU := s.serverCountry == "RU"
		shopInRU := shopJurisdiction == "RU"
		llmForeign := llmJurisdiction != "RU"

		if (!serverInRU && shopInRU) || (serverInRU && shopInRU && llmForeign) {
			shopPDN = true
		}
	}

	// Manual PDN override takes precedence
	if shopPDN || aiKeyPDN {
		slog.Debug("PDN protection active",
			"shop_jurisdiction", shopJurisdiction,
			"llm_provider", llmJurisdiction,
			"shop_pdn_flag", shopPDN,
			"manual_ai_key", aiKeyPDN,
		)
	}

	if shopName != "" || keyName != "" {
		ctx = context.WithValue(ctx, mcpShopNameKey, shopName)
		ctx = context.WithValue(ctx, mcpKeyNameKey, keyName)
		ctx = context.WithValue(ctx, mcpReadOnlyKey, readOnly)
		ctx = context.WithValue(ctx, mcpShopPDNKey, shopPDN)
		ctx = context.WithValue(ctx, mcpAIKeyPDNKey, aiKeyPDN)
	}

	s.httpServer.ServeHTTP(w, r.WithContext(ctx))
}

func (s *MCPServer) NotifyToolsChanged() {
	s.mcpServer.SendNotificationToAllClients(
		"notifications/tools/list_changed",
		nil,
	)
	slog.Info("sent tools/list_changed notification")
}

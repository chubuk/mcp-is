package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/chubuk/mcp-is/internal/pdn"
	"github.com/chubuk/mcp-is/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func formatResult(v any) string {
	if v == nil {
		return "null"
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

// wrapEnvelope гарантирует, что payload сериализуется в JSON-объект:
// structuredContent по MCP-спеке (CallToolResult, 2025-06-18) — только объект,
// никогда массив/скаляр. total = len(items) — семантика домашнего конверта
// (ProductsPage.Total = len(result)).
func wrapEnvelope(v any) any {
	if v == nil {
		return v // omitempty: поле отсутствует — спека это разрешает
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() != reflect.Uint8 {
		return map[string]any{"items": v, "total": rv.Len()}
	}
	if b, err := json.Marshal(v); err == nil && len(b) > 0 && b[0] != '{' {
		return map[string]any{"value": v} // скаляры, *[]T, json.RawMessage с массивом и пр.
	}
	return v
}

// toolResult returns the same payload as pretty JSON text and as structuredContent
// (MCP dual form; clients may use either).
func toolResult(v any) *mcp.CallToolResult {
	v = wrapEnvelope(v)
	return mcp.NewToolResultStructured(v, formatResult(v))
}

func isHTMLResponse(s string) bool {
	trimmed := strings.TrimSpace(s)
	return strings.HasPrefix(trimmed, "<!DOCTYPE") ||
		strings.HasPrefix(trimmed, "<html") ||
		strings.HasPrefix(trimmed, "<HTML") ||
		strings.Contains(s, "<!DOCTYPE html") ||
		strings.Contains(s, "<!DOCTYPE HTML")
}

func formatResultOrError(v any) (string, bool) {
	s := formatResult(v)
	if isHTMLResponse(s) {
		return "InSales API returned an error page. This endpoint may not be available for your API key.", true
	}
	return s, false
}

func cleanError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if isHTMLResponse(msg) {
		// Extract status code if present (e.g. "API error 404: <html...")
		if idx := strings.Index(msg, "<!"); idx > 0 {
			return fmt.Errorf("InSales API error: %s", strings.TrimSpace(msg[:idx]))
		}
		return fmt.Errorf("InSales API returned an error page (endpoint may not be available for this API key)")
	}
	return err
}

func toolError(err error) *mcp.CallToolResult {
	if errors.Is(err, core.ErrBillingExpired) || errors.Is(err, core.ErrBillingUnpaid) {
		return mcp.NewToolResultError("Подписка магазина на приложение не оплачена или период закончился. Данные недоступны до оплаты — продлите приложение в админке InSales.")
	}
	if errors.Is(err, core.ErrKeyNotBound) {
		return mcp.NewToolResultError("AI-ключ не привязан к магазину. Попросите владельца привязать ключ к магазину в личном кабинете (раздел «AI-ключи») и повторите запрос.")
	}
	if errors.Is(err, insales.ErrForbidden) {
		return mcp.NewToolResultError("Нет прав для данной операции — проверьте настройки прав приложения в InSales (см. store_info).")
	}
	if errors.Is(err, insales.ErrUnauthorized) {
		return mcp.NewToolResultError("Неверные учётные данные InSales — обновите ключ магазина в личном кабинете.")
	}
	return mcp.NewToolResultError(cleanError(err).Error())
}

// MutationResult — подробный конверт результата мутации (план 3c #4).
// LLM нужна максимальная подробность: что именно изменено/удалено.
type MutationResult struct {
	OK            bool                   `json:"ok"`
	Action        string                 `json:"action"` // "update" | "delete"
	Resource      string                 `json:"resource"`
	ID            string                 `json:"id"`
	Title         string                 `json:"title,omitempty"`
	ChangedFields map[string]FieldChange `json:"changed_fields,omitempty"`
	IgnoredKeys   []string               `json:"ignored_keys,omitempty"`
	Data          any                    `json:"data,omitempty"`    // полное итоговое состояние (update)
	Deleted       any                    `json:"deleted,omitempty"` // снимок до удаления (delete)
	Message       string                 `json:"message,omitempty"`
}

// FieldChange — old/new значение поля.
type FieldChange struct {
	Old any `json:"old"`
	New any `json:"new"`
}

// deleteResult формирует конверт удаления с pre-fetch снимка.
// Если pre-fetch вернул 404 — deleted:null + пояснение.
func deleteResult(resource, id string, preFetch func() (any, error), doDelete func() error) *mcp.CallToolResult {
	var snapshot any
	if preFetch != nil {
		obj, err := preFetch()
		if err == nil {
			snapshot = obj
		} else if !errors.Is(err, insales.ErrNotFound) {
			slog.Warn("delete pre-fetch failed", "resource", resource, "id", id, "error", err)
		}
	}
	if err := doDelete(); err != nil {
		return toolError(err)
	}
	msg := fmt.Sprintf("%s удалён", resource)
	if snapshot == nil {
		msg = fmt.Sprintf("%s удалён (объект не найден на момент удаления)", resource)
	}
	return toolResult(MutationResult{
		OK:       true,
		Action:   "delete",
		Resource: resource,
		ID:       id,
		Deleted:  snapshot,
		Message:  msg,
	})
}

// updateResult формирует конверт обновления с diff (pre-fetch → post-fetch).
func updateResult(resource, id, title string, preFetch func() (any, error), doUpdate func() (any, error), changedFields map[string]FieldChange, ignored []string) *mcp.CallToolResult {
	// Pre-fetch для old-значений (если ещё не переданы).
	if changedFields == nil && preFetch != nil {
		changedFields = map[string]FieldChange{}
	}
	var updated any
	if doUpdate != nil {
		obj, err := doUpdate()
		if err != nil {
			return toolError(err)
		}
		updated = obj
	}
	return toolResult(MutationResult{
		OK:            true,
		Action:        "update",
		Resource:      resource,
		ID:            id,
		Title:         title,
		ChangedFields: changedFields,
		IgnoredKeys:   ignored,
		Data:          updated,
	})
}

// asMap нормализует объект (map или struct из маппера) в map[string]any.
func asMap(obj any) map[string]any {
	switch t := obj.(type) {
	case map[string]any:
		return t
	case nil:
		return nil
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

// fieldOf возвращает значение поля объекта по ключу аргумента.
// Особый случай: цена товара живёт в variants[0].price (product.price = null).
func fieldOf(obj any, key string) any {
	m := asMap(obj)
	if m == nil {
		return nil
	}
	if key == "price" {
		if v, ok := m["variants"].([]any); ok && len(v) > 0 {
			if vm, ok := v[0].(map[string]any); ok {
				return vm["price"]
			}
		}
	}
	return m[key]
}

// diffFromArgs строит changed_fields {old,new} для присланных аргументов:
// old — из pre-fetch снимка, new — значение из аргументов (то, что ушло в InSales).
func diffFromArgs(oldObj any, args map[string]any, ignored []string) map[string]FieldChange {
	skip := map[string]bool{}
	for _, k := range ignored {
		skip[k] = true
	}
	diff := map[string]FieldChange{}
	for k, v := range args {
		if skip[k] || k == "id" || strings.HasSuffix(k, "_id") {
			continue
		}
		diff[k] = FieldChange{Old: fieldOf(oldObj, k), New: v}
	}
	if len(diff) == 0 {
		return nil
	}
	return diff
}

// updateWithDiff — единый конверт update-мутации: pre-fetch снимка → diff → обновление.
// +1 GET на мутацию (осознанно, план 3c): LLM видит old/new каждого изменённого поля.
// mask (опционально) — ПДн-маскировка: применяется к Data и old-значениям diff'а.
func updateWithDiff(resource, id string, args map[string]any, ignored []string,
	preFetch func() (any, error), doUpdate func() (any, error), mask func(any) any) *mcp.CallToolResult {
	var old any
	if preFetch != nil {
		if obj, err := preFetch(); err == nil {
			old = obj
		} else if !errors.Is(err, insales.ErrNotFound) {
			slog.Warn("update pre-fetch failed", "resource", resource, "id", id, "error", err)
		}
	}
	var updated any
	if doUpdate != nil {
		obj, err := doUpdate()
		if err != nil {
			return toolError(err)
		}
		updated = obj
	}
	changed := diffFromArgs(old, args, ignored)
	if mask != nil {
		updated = mask(updated)
		for k, fc := range changed {
			// Скаляр маскируем через map-обёртку — MaskAny маскирует по имени PII-ключа.
			if masked, ok := mask(map[string]any{k: fc.Old}).(map[string]any); ok {
				fc.Old = masked[k]
			}
			changed[k] = fc
		}
	}
	return toolResult(MutationResult{
		OK:            true,
		Action:        "update",
		Resource:      resource,
		ID:            id,
		Title:         fmt.Sprint(fieldOf(updated, "title")),
		ChangedFields: changed,
		IgnoredKeys:   ignored,
		Data:          updated,
	})
}

// getFloatArg parses a numeric tool argument (float64, int, json.Number or
// numeric string). Absent key → def; present but unparseable → error.
func getFloatArg(args map[string]any, key string, def float64) (float64, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return def, nil
	}
	switch t := v.(type) {
	case float64:
		return t, nil
	case float32:
		return float64(t), nil
	case int:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case json.Number:
		f, err := t.Float64()
		if err != nil {
			return 0, fmt.Errorf("%s must be a number, got %q", key, t)
		}
		return f, nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a number, got %q", key, t)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("%s must be a number", key)
	}
}

func getFloat(args map[string]any, key string, def float64) float64 {
	f, _ := getFloatArg(args, key, def)
	return f
}

// pageArgs parses page/per_page for list tools with clamping: page ≥ 1,
// 1 ≤ per_page ≤ 250; per_page < 1 means "unset" → default 25. Garbage values
// are a tool error, not a silent default.
func pageArgs(args map[string]any) (page, perPage int, err error) {
	pf, err := getFloatArg(args, "page", 1)
	if err != nil {
		return 0, 0, err
	}
	ppf, err := getFloatArg(args, "per_page", 25)
	if err != nil {
		return 0, 0, err
	}
	page = int(pf)
	if page < 1 {
		page = 1
	}
	perPage = int(ppf)
	if perPage < 1 {
		perPage = 25
	}
	if perPage > 250 {
		perPage = 250
	}
	return page, perPage, nil
}

// applyPDNToResult маскирует ПДн в результате инструмента при включённом режиме.
// Значение целиком JSON round-trip-ится (typed-структуры → generic-дерево → маска →
// обратно в исходный Go-тип), чтобы вызывающие хендлеры и типизированные конверты
// продолжали работать (masked.([]insales.Order) и т.п.).
func applyPDNToResult(ctx context.Context, result any, toolName string) any {
	shopPDN, _ := ctx.Value(mcpShopPDNKey).(bool)
	aiKeyPDN, _ := ctx.Value(mcpAIKeyPDNKey).(bool)

	if (!shopPDN && !aiKeyPDN) || result == nil {
		return result
	}

	b, err := json.Marshal(result)
	if err != nil {
		return result
	}
	var generic any
	if err := json.Unmarshal(b, &generic); err != nil {
		return result
	}
	masked := pdn.MaskAny(generic)
	mb, err := json.Marshal(masked)
	if err != nil {
		return masked
	}
	out := reflect.New(reflect.TypeOf(result))
	if err := json.Unmarshal(mb, out.Interface()); err != nil {
		return masked
	}
	return out.Elem().Interface()
}

// === SEO ===

var seoReadArgs = map[string]bool{"resource_type": true, "resource_id": true, "blog_id": true}

var seoWriteArgs = map[string]bool{
	"resource_type": true, "resource_id": true, "blog_id": true, "clear_fields": true,
	"html_title": true, "meta_keywords": true, "meta_description": true, "seo_description": true,
}

// stringArg reads a string tool argument, coercing numbers (IDs often arrive as numbers).
func stringArg(args map[string]any, key string) string {
	switch v := args[key].(type) {
	case string:
		return v
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.Itoa(v)
	case json.Number:
		return v.String()
	}
	return ""
}

func seoResultMap(resourceType, resourceID string, fields *insales.SEOFields) map[string]any {
	out := map[string]any{
		"resource_type":    resourceType,
		"resource_id":      resourceID,
		"html_title":       fields.HTMLTitle,
		"meta_keywords":    fields.MetaKeywords,
		"meta_description": fields.MetaDescription,
	}
	if resourceType == "collection" {
		out["seo_description"] = fields.SEODescription
	}
	return out
}

func getSeoHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		resourceType := strings.TrimSpace(stringArg(args, "resource_type"))
		resourceID := strings.TrimSpace(stringArg(args, "resource_id"))
		blogID := strings.TrimSpace(stringArg(args, "blog_id"))
		if resourceType == "" || resourceID == "" {
			return toolError(errors.New("resource_type and resource_id are required")), nil
		}
		ignored := filterKnown(args, seoReadArgs)
		fields, err := svc.GetSEO(ctx, resourceType, resourceID, blogID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(withIgnoredKeys(seoResultMap(resourceType, resourceID, fields), ignored)), nil
	}
}

func updateSeoHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		resourceType := strings.TrimSpace(stringArg(args, "resource_type"))
		resourceID := strings.TrimSpace(stringArg(args, "resource_id"))
		blogID := strings.TrimSpace(stringArg(args, "blog_id"))
		if resourceType == "" || resourceID == "" {
			return toolError(errors.New("resource_type and resource_id are required")), nil
		}
		if err := resolveKeywordsAlias(args); err != nil {
			return toolError(err), nil
		}
		ignored := filterKnown(args, seoWriteArgs)
		delete(args, "resource_type")
		return updateWithDiff("seo", resourceID, args, ignored,
			func() (any, error) {
				f, err := svc.GetSEO(ctx, resourceType, resourceID, blogID)
				if err != nil {
					return nil, err
				}
				return f, nil
			},
			func() (any, error) {
				f, err := svc.UpdateSEO(ctx, resourceType, resourceID, blogID, args)
				if err != nil {
					return nil, err
				}
				return seoResultMap(resourceType, resourceID, f), nil
			},
			nil,
		), nil
	}
}

// === ТОВАРЫ ===

func listProductsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetProducts(ctx, page, perPage, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(ProductsPage{Items: result, Total: len(result), Page: page, PerPage: perPage}), nil
	}
}

func getProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetProduct(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		price, err := req.RequireString("price")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateProduct(ctx, title, price, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// resolveKeywordsAlias maps the legacy keywords key to meta_keywords; passing
// both is a conflict.
func resolveKeywordsAlias(args map[string]any) error {
	if _, hasAlias := args["keywords"]; !hasAlias {
		return nil
	}
	if _, hasCanonical := args["meta_keywords"]; hasCanonical {
		return errors.New("передайте только meta_keywords: keywords — устаревший алиас")
	}
	args["meta_keywords"] = args["keywords"]
	delete(args, "keywords")
	return nil
}

// filterKnown drops non-whitelisted keys and returns them sorted (for ignored_keys).
func filterKnown(args map[string]any, known map[string]bool) []string {
	ignored := []string{}
	for k := range args {
		if !known[k] {
			ignored = append(ignored, k)
			delete(args, k)
		}
	}
	sort.Strings(ignored)
	return ignored
}

// filterArgsByTool фильтрует args по whitelist из inputSchema инструмента
// (toolArgWhitelist, заполняется в addTool). Нет записи — не трогаем (fail-open).
func filterArgsByTool(args map[string]any, toolName string) []string {
	wl, ok := toolArgWhitelist[toolName]
	if !ok {
		return nil
	}
	return filterKnown(args, wl)
}

// wrapArgsAudit оборачивает handler мутационного инструмента: отфильтровывает
// неизвестные ключи аргументов и добавляет их в результат как ignored_keys.
// Применяется в addTool для create/update/delete — единообразно, без дрейфа.
func wrapArgsAudit(handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		ignored := filterArgsByTool(args, req.Params.Name)
		res, err := handler(ctx, req)
		if err != nil || res == nil || len(ignored) == 0 {
			return res, err
		}
		if res.IsError {
			return res, err
		}
		payload := res.StructuredContent
		if payload == nil {
			return res, err
		}
		updated := withIgnoredKeys(payload, ignored)
		out := *res
		out.StructuredContent = updated
		if len(out.Content) > 0 {
			if tc, ok := out.Content[0].(mcp.TextContent); ok {
				tc.Text = formatResult(updated)
				out.Content[0] = tc
			}
		}
		return &out, err
	}
}

// withIgnoredKeys adds ignored_keys to an object result (additive key).
func withIgnoredKeys(result any, ignored []string) any {
	if len(ignored) == 0 {
		return result
	}
	b, err := json.Marshal(result)
	if err != nil {
		return result
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return result
	}
	m["ignored_keys"] = ignored
	return m
}

// productArgFields is the update_product argument whitelist (handler UX layer;
// the InSales client enforces the same list as a safety net).
var productArgFields = map[string]bool{
	"product_id": true, "clear_fields": true,
	"title": true, "description": true, "short_description": true,
	"vendor": true, "sku": true,
	"html_title": true, "meta_keywords": true, "meta_description": true,
	"category_id": true,
}

// orderArgFields is the update_order argument whitelist.
var orderArgFields = map[string]bool{
	"order_id": true, "status": true, "financial_status": true,
	"custom_status_permalink": true, "fulfillment_status": true,
	"comment": true, "manager_comment": true,
	"delivery_variant_id": true, "payment_gateway_id": true, "warehouse_id": true,
	"coupon":                 true,
	"order_lines_attributes": true, "discounts_attributes": true,
	"shipping_address_attributes": true,
}

func updateProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		if err := resolveKeywordsAlias(args); err != nil {
			return toolError(err), nil
		}
		ignored := filterKnown(args, productArgFields)
		return updateWithDiff("product", productID, args, ignored,
			func() (any, error) { return svc.GetProduct(ctx, productID) },
			func() (any, error) { return svc.UpdateProduct(ctx, productID, args) },
			nil,
		), nil
	}
}

func deleteProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("product", productID,
			func() (any, error) { return svc.GetProduct(ctx, productID) },
			func() error { return svc.DeleteProduct(ctx, productID) },
		), nil
	}
}

// === ЗАКАЗЫ ===

func listOrdersHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		filters := make(map[string]any)
		if v, ok := args["fulfillment_status"].(string); ok && v != "" {
			filters["fulfillment_status"] = v
		}
		if v, ok := args["updated_since"].(string); ok && v != "" {
			filters["updated_since"] = v
		}
		if v, ok := args["from_id"].(float64); ok && v > 0 {
			filters["from_id"] = int(v)
		}
		result, err := svc.GetOrders(ctx, page, perPage, filters)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "list_orders")
		orders, _ := masked.([]insales.Order)
		return toolResult(OrdersPage{Items: orders, Total: len(orders), Page: page, PerPage: perPage}), nil
	}
}

func getOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetOrderDetails(ctx, orderID)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "get_order")
		return toolResult(masked), nil
	}
}

func updateOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "order_id")
		ignored := filterKnown(args, orderArgFields)
		return updateWithDiff("order", orderID, args, ignored,
			func() (any, error) { return svc.GetOrderDetails(ctx, orderID) },
			func() (any, error) { return svc.UpdateOrder(ctx, orderID, args) },
			func(v any) any { return applyPDNToResult(ctx, v, "update_order") },
		), nil
	}
}

func createOrderHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		clientID := int(getFloat(args, "client_id", 0))
		// client_id optional when client / client_attributes present
		if clientID == 0 {
			if _, hasClient := args["client"]; !hasClient {
				if _, hasAttrs := args["client_attributes"]; !hasAttrs {
					if _, hasEmail := args["client_email"]; !hasEmail {
						return mcp.NewToolResultError("client_id or client (name/email/phone) is required"), nil
					}
				}
			}
		}
		result, err := svc.CreateOrder(ctx, clientID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(applyPDNToResult(ctx, result, "create_order")), nil
	}
}

// === КЛИЕНТЫ ===

func listClientsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetClients(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "list_clients")
		clients, _ := masked.([]any)
		return toolResult(ClientsPage{Items: clients, Total: len(clients), Page: page, PerPage: perPage}), nil
	}
}

func getClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetClient(ctx, clientID)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "get_client")
		return toolResult(masked), nil
	}
}

func createClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateClient(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		masked := applyPDNToResult(ctx, result, "create_client")
		return toolResult(masked), nil
	}
}

func updateClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("client", clientID, args, nil,
			func() (any, error) { return svc.GetClient(ctx, clientID) },
			func() (any, error) { return svc.UpdateClient(ctx, clientID, args) },
			func(v any) any { return applyPDNToResult(ctx, v, "update_client") },
		), nil
	}
}

func deleteClientHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("client", clientID,
			func() (any, error) { return svc.GetClient(ctx, clientID) },
			func() error { return svc.DeleteClient(ctx, clientID) },
		), nil
	}
}

// === КАТЕГОРИИ ===

func listCategoriesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCategories(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCategory(ctx, categoryID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateCategory(ctx, title, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("category", categoryID, args, nil,
			func() (any, error) { return svc.GetCategory(ctx, categoryID) },
			func() (any, error) { return svc.UpdateCategory(ctx, categoryID, args) },
			nil,
		), nil
	}
}

func deleteCategoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID, err := req.RequireString("category_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("category", categoryID,
			func() (any, error) { return svc.GetCategory(ctx, categoryID) },
			func() error { return svc.DeleteCategory(ctx, categoryID) },
		), nil
	}
}

// === КОЛЛЕКЦИИ ===

func listCollectionsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCollections(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCollection(ctx, collectionID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateCollection(ctx, title, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		title, _ := req.RequireString("title")
		return updateWithDiff("collection", collectionID, map[string]any{"title": title}, nil,
			func() (any, error) { return svc.GetCollection(ctx, collectionID) },
			func() (any, error) { return svc.UpdateCollection(ctx, collectionID, title) },
			nil,
		), nil
	}
}

func deleteCollectionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("collection", collectionID,
			func() (any, error) { return svc.GetCollection(ctx, collectionID) },
			func() error { return svc.DeleteCollection(ctx, collectionID) },
		), nil
	}
}

// === ИЗОБРАЖЕНИЯ ===

func listImagesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListImages(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func addImageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		url, err := req.RequireString("url")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.AddImage(ctx, productID, url, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteImageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		imageID, err := req.RequireString("image_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("image", imageID,
			nil,
			func() error { return svc.DeleteImage(ctx, productID, imageID) },
		), nil
	}
}

// === МЕТАПОЛЯ ===

func listMetaFieldsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.ListMetaFields(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateMetaField(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		metaFieldID, err := req.RequireString("meta_field_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("meta_field", metaFieldID, args, nil,
			nil,
			func() (any, error) { return svc.UpdateMetaField(ctx, metaFieldID, args) },
			nil,
		), nil
	}
}

func deleteMetaFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		metaFieldID, err := req.RequireString("meta_field_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("meta_field", metaFieldID,
			nil,
			func() error { return svc.DeleteMetaField(ctx, metaFieldID) },
		), nil
	}
}

// === ВЕБХУКИ ===

func listWebhooksHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListWebhooks(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		event, err := req.RequireString("event")
		if err != nil {
			return toolError(err), nil
		}
		url, err := req.RequireString("url")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateWebhook(ctx, event, url, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		webhookID, err := req.RequireString("webhook_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("webhook", webhookID, args, nil,
			nil,
			func() (any, error) { return svc.UpdateWebhook(ctx, webhookID, args) },
			nil,
		), nil
	}
}

func deleteWebhookHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		webhookID, err := req.RequireString("webhook_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("webhook", webhookID,
			nil,
			func() error { return svc.DeleteWebhook(ctx, webhookID) },
		), nil
	}
}

// === СИСТЕМНЫЕ ===

func getShopInfoHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.GetStoreDetails(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listCurrenciesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListCurrencies(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listDeliveryVariantsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListDeliveryVariants(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listPaymentGatewaysHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListPaymentGateways(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === СТРАНИЦЫ ===

func listPagesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListPages(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getPageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetPage(ctx, pageID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createPageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		content, err := req.RequireString("content")
		if err != nil {
			return toolError(err), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreatePage(ctx, title, content, handle)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updatePageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("page", pageID, args, nil,
			func() (any, error) { return svc.GetPage(ctx, pageID) },
			func() (any, error) { return svc.UpdatePage(ctx, pageID, args) },
			nil,
		), nil
	}
}

func deletePageHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pageID, err := req.RequireString("page_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("page", pageID,
			func() (any, error) { return svc.GetPage(ctx, pageID) },
			func() error { return svc.DeletePage(ctx, pageID) },
		), nil
	}
}

// === АКСЕССУАРЫ ===

func listAccessoriesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListAccessories(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetAccessory(ctx, accessoryID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateAccessory(ctx, title, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("accessory", accessoryID, args, nil,
			func() (any, error) { return svc.GetAccessory(ctx, accessoryID) },
			func() (any, error) { return svc.UpdateAccessory(ctx, accessoryID, args) },
			nil,
		), nil
	}
}

func deleteAccessoryHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("accessory", accessoryID,
			func() (any, error) { return svc.GetAccessory(ctx, accessoryID) },
			func() error { return svc.DeleteAccessory(ctx, accessoryID) },
		), nil
	}
}

// === ВАРИАНТЫ ===

func listVariantsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListVariants(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetVariant(ctx, productID, variantID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("variant", variantID, args, nil,
			func() (any, error) { return svc.GetVariant(ctx, productID, variantID) },
			func() (any, error) { return svc.UpdateVariant(ctx, productID, variantID, args) },
			nil,
		), nil
	}
}

func deleteVariantHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("variant", variantID,
			func() (any, error) { return svc.GetVariant(ctx, productID, variantID) },
			func() error { return svc.DeleteVariant(ctx, productID, variantID) },
		), nil
	}
}

// === ЗНАЧЕНИЯ СВОЙСТВ ===

func listPropertyValuesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListPropertyValues(ctx, accessoryID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createPropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreatePropertyValue(ctx, accessoryID, title)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updatePropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		propertyValueID, err := req.RequireString("property_value_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("property_value", propertyValueID, args, nil,
			nil, // GetPropertyValue в API нет
			func() (any, error) { return svc.UpdatePropertyValue(ctx, accessoryID, propertyValueID, args) },
			nil,
		), nil
	}
}

func deletePropertyValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessoryID, err := req.RequireString("accessory_id")
		if err != nil {
			return toolError(err), nil
		}
		propertyValueID, err := req.RequireString("property_value_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("property_value", propertyValueID,
			nil,
			func() error { return svc.DeletePropertyValue(ctx, accessoryID, propertyValueID) },
		), nil
	}
}

// === МЕНЮ ===

func listMenusHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListMenus(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetMenu(ctx, menuID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreateMenu(ctx, title, handle)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("menu", menuID, args, nil,
			func() (any, error) { return svc.GetMenu(ctx, menuID) },
			func() (any, error) { return svc.UpdateMenu(ctx, menuID, args) },
			nil,
		), nil
	}
}

func deleteMenuHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("menu", menuID,
			nil,
			func() error { return svc.DeleteMenu(ctx, menuID) },
		), nil
	}
}

// === РЕДИРЕКТЫ ===

func listRedirectsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListRedirects(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createRedirectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, err := req.RequireString("url")
		if err != nil {
			return toolError(err), nil
		}
		target, err := req.RequireString("target")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreateRedirect(ctx, url, target)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteRedirectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		redirectID, err := req.RequireString("redirect_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("redirect", redirectID,
			nil,
			func() error { return svc.DeleteRedirect(ctx, redirectID) },
		), nil
	}
}

// === БЛОГИ ===

func listBlogsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListBlogs(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetBlog(ctx, blogID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		handle, _ := req.RequireString("handle")
		result, err := svc.CreateBlog(ctx, title, handle)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("blog", blogID,
			nil,
			func() error { return svc.DeleteBlog(ctx, blogID) },
		), nil
	}
}

func updateBlogHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("blog", blogID, args, nil,
			func() (any, error) { return svc.GetBlog(ctx, blogID) },
			func() (any, error) { return svc.UpdateBlog(ctx, blogID, args) },
			nil,
		), nil
	}
}

// === ЗАПИСИ БЛОГА ===

func listBlogPostsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListBlogPosts(ctx, blogID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetBlogPost(ctx, blogID, postID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		title, err := req.RequireString("title")
		if err != nil {
			return toolError(err), nil
		}
		bodyHTML, err := req.RequireString("body_html")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreateBlogPost(ctx, blogID, title, bodyHTML)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("blog_post", postID, args, nil,
			func() (any, error) { return svc.GetBlogPost(ctx, blogID, postID) },
			func() (any, error) { return svc.UpdateBlogPost(ctx, blogID, postID, args) },
			nil,
		), nil
	}
}

func deleteBlogPostHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blogID, err := req.RequireString("blog_id")
		if err != nil {
			return toolError(err), nil
		}
		postID, err := req.RequireString("post_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("blog_post", postID,
			nil,
			func() error { return svc.DeleteBlogPost(ctx, blogID, postID) },
		), nil
	}
}

// === DISCOUNT CODES ===

func listDiscountCodesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListDiscountCodes(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codeID, err := req.RequireString("discount_code_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetDiscountCode(ctx, codeID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateDiscountCode(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codeID, err := req.RequireString("discount_code_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("discount_code", codeID, args, nil,
			func() (any, error) { return svc.GetDiscountCode(ctx, codeID) },
			func() (any, error) { return svc.UpdateDiscountCode(ctx, codeID, args) },
			nil,
		), nil
	}
}

func deleteDiscountCodeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codeID, err := req.RequireString("discount_code_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("discount_code", codeID,
			nil,
			func() error { return svc.DeleteDiscountCode(ctx, codeID) },
		), nil
	}
}

// === COLLECTS ===

func listCollectsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListCollects(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCollectsByProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCollectsByProduct(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCollectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		collectionID := int(getFloat(args, "collection_id", 0))
		productID := int(getFloat(args, "product_id", 0))
		if collectionID == 0 || productID == 0 {
			return toolError(fmt.Errorf("collection_id and product_id are required")), nil
		}
		result, err := svc.CreateCollect(ctx, collectionID, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateCollectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectID, err := req.RequireString("collect_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "collect_id")
		payload := map[string]any{}
		if v, ok := args["collection_id"]; ok && v != nil && v != "" {
			payload["collection_id"] = v
		}
		if v, ok := args["position"]; ok && v != nil {
			payload["position"] = v
		}
		if len(payload) == 0 {
			return mcp.NewToolResultError("collection_id or position is required"), nil
		}
		result, err := svc.UpdateCollect(ctx, collectID, payload)
		if err != nil {
			return toolError(err), nil
		}
		// GetCollect (одиночного) в API нет — diff без old.
		return updateWithDiff("collect", collectID, payload, nil, nil,
			func() (any, error) { return result, nil },
			nil,
		), nil
	}
}

func deleteCollectHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectID, err := req.RequireString("collect_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("collect", collectID,
			nil,
			func() error { return svc.DeleteCollect(ctx, collectID) },
		), nil
	}
}

// === PROPERTIES ===

func updateRecurringApplicationChargeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		monthly, ok := req.GetArguments()["monthly"]
		if !ok || monthly == nil || monthly == "" {
			return mcp.NewToolResultError("monthly is required"), nil
		}
		payload := map[string]any{"monthly": monthly}
		return updateWithDiff("recurring_application_charge", "", payload, nil, nil,
			func() (any, error) { return svc.UpdateRecurringApplicationCharge(ctx, payload) },
			nil,
		), nil
	}
}

func listPropertiesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListProperties(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getPropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetProperty(ctx, propertyID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createPropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateProperty(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updatePropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("property", propertyID, args, nil,
			func() (any, error) { return svc.GetProperty(ctx, propertyID) },
			func() (any, error) { return svc.UpdateProperty(ctx, propertyID, args) },
			nil,
		), nil
	}
}

func deletePropertyHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("property", propertyID,
			nil,
			func() error { return svc.DeleteProperty(ctx, propertyID) },
		), nil
	}
}

// === CHARACTERISTICS ===

func listCharacteristicsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListCharacteristics(ctx, propertyID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		charID, err := req.RequireString("characteristic_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetCharacteristic(ctx, propertyID, charID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateCharacteristic(ctx, propertyID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		charID, err := req.RequireString("characteristic_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("characteristic", charID, args, nil,
			func() (any, error) { return svc.GetCharacteristic(ctx, propertyID, charID) },
			func() (any, error) { return svc.UpdateCharacteristic(ctx, propertyID, charID, args) },
			nil,
		), nil
	}
}

func deleteCharacteristicHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		propertyID, err := req.RequireString("property_id")
		if err != nil {
			return toolError(err), nil
		}
		charID, err := req.RequireString("characteristic_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("characteristic", charID,
			nil,
			func() error { return svc.DeleteCharacteristic(ctx, propertyID, charID) },
		), nil
	}
}

// === REVIEWS ===

func listReviewsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		page, perPage, err := pageArgs(args)
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListReviews(ctx, page, perPage)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reviewID, err := req.RequireString("review_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetReview(ctx, reviewID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateReview(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reviewID, err := req.RequireString("review_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("review", reviewID, args, nil,
			func() (any, error) { return svc.GetReview(ctx, reviewID) },
			func() (any, error) { return svc.UpdateReview(ctx, reviewID, args) },
			nil,
		), nil
	}
}

func deleteReviewHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reviewID, err := req.RequireString("review_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("review", reviewID,
			nil,
			func() error { return svc.DeleteReview(ctx, reviewID) },
		), nil
	}
}

// === FILES ===

func listFilesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := svc.ListFiles(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func getFileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileID, err := req.RequireString("file_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetFile(ctx, fileID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createFileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		result, err := svc.CreateFile(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func deleteFileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileID, err := req.RequireString("file_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("file", fileID,
			nil,
			func() error { return svc.DeleteFile(ctx, fileID) },
		), nil
	}
}

// Generic CRUD handlers for round 3+ resources

func genericListHandler(svc *service.Service, listFn func(ctx context.Context) ([]any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := listFn(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericGetHandler(svc *service.Service, param string, getFn func(ctx context.Context, id string) (any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString(param)
		if err != nil {
			return toolError(err), nil
		}
		result, err := getFn(ctx, id)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericCreateHandler(svc *service.Service, createFn func(ctx context.Context, args map[string]any) (any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		ignored := filterArgsByTool(args, req.Params.Name)
		result, err := createFn(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(withIgnoredKeys(result, ignored)), nil
	}
}

func genericUpdateHandler(svc *service.Service, param string, updateFn func(ctx context.Context, id string, args map[string]any) (any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return genericUpdateGetHandler(param, nil, updateFn)
}

// genericUpdateGetHandler — generic update с pre-fetch снимка для changed_fields.
// getFn может быть nil (diff без old-значений).
func genericUpdateGetHandler(param string, getFn func(ctx context.Context, id string) (any, error), updateFn func(ctx context.Context, id string, args map[string]any) (any, error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString(param)
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		ignored := filterArgsByTool(args, req.Params.Name)
		var preFetch func() (any, error)
		if getFn != nil {
			preFetch = func() (any, error) { return getFn(ctx, id) }
		}
		return updateWithDiff("", id, args, ignored, preFetch,
			func() (any, error) { return updateFn(ctx, id, args) },
			nil,
		), nil
	}
}

func genericDeleteHandler(deleteFn func(ctx context.Context, id string) error, paramName, deletedMsg string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString(paramName)
		if err != nil {
			return toolError(err), nil
		}
		if err := deleteFn(ctx, id); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(deletedMsg), nil
	}
}

// nested helpers for parent-scoped resources (product_id + value_id, etc.)

func nestedListHandler(listFn func(ctx context.Context, parentID string) ([]any, error), parentParam string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		result, err := listFn(ctx, parentID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(map[string]any{
			"items": result,
			"total": len(result),
		}), nil
	}
}

func nestedGetHandler(getFn func(ctx context.Context, parentID, key string) (any, error), parentParam, keyParam string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		key, err := req.RequireString(keyParam)
		if err != nil {
			return toolError(err), nil
		}
		result, err := getFn(ctx, parentID, key)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func nestedUpdateHandler(updateFn func(ctx context.Context, parentID, key string, args map[string]any) (any, error), parentParam, keyParam string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		key, err := req.RequireString(keyParam)
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, parentParam)
		delete(args, keyParam)
		// body is only the value payload
		payload := map[string]any{}
		if v, ok := args["value"]; ok {
			payload["value"] = v
		}
		for k, v := range args {
			if k != "value" && k != "name" && k != "args" {
				payload[k] = v
			}
		}
		result, err := updateFn(ctx, parentID, key, payload)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func nestedDeleteHandler(deleteFn func(ctx context.Context, parentID, key string) error, parentParam, keyParam, deletedMsg string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parentID, err := req.RequireString(parentParam)
		if err != nil {
			return toolError(err), nil
		}
		key, err := req.RequireString(keyParam)
		if err != nil {
			return toolError(err), nil
		}
		if err := deleteFn(ctx, parentID, key); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(deletedMsg), nil
	}
}

func createProductFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		if _, ok := args["type"]; !ok {
			if ft, ok := args["field_type"].(string); ok && ft != "" {
				args["type"] = ft
			}
		}
		delete(args, "field_type")
		if _, ok := args["handle"].(string); !ok {
			return mcp.NewToolResultError("handle is required (language-invariant slug)"), nil
		}
		if _, ok := args["type"].(string); !ok {
			return mcp.NewToolResultError("type is required: ProductField::TextField | ProductField::TextArea"), nil
		}
		result, err := svc.CreateProductField(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func updateProductFieldHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fieldID, err := req.RequireString("field_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "field_id")
		delete(args, "field_type")
		if ft, ok := req.GetArguments()["field_type"].(string); ok && ft != "" {
			if _, has := args["type"]; !has {
				args["type"] = ft
			}
		}
		return updateWithDiff("product_field", fieldID, args, nil,
			func() (any, error) { return svc.GetProductField(ctx, fieldID) },
			func() (any, error) { return svc.UpdateProductField(ctx, fieldID, args) },
			nil,
		), nil
	}
}

func createProductFieldValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "product_id")
		if _, ok := args["product_field_id"]; !ok {
			if _, ok := args["product_field_handle"]; !ok {
				return mcp.NewToolResultError("product_field_id or product_field_handle is required"), nil
			}
		}
		result, err := svc.CreateProductFieldValue(ctx, productID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func createVariantFieldValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		variantID, err := req.RequireString("variant_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "variant_id")
		if _, ok := args["variant_field_id"]; !ok {
			if _, ok := args["variant_field_handle"]; !ok {
				return mcp.NewToolResultError("variant_field_id or variant_field_handle is required"), nil
			}
		}
		result, err := svc.CreateVariantFieldValue(ctx, variantID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === ORDER SUB-OPERATIONS (map into order[...] via UpdateOrder) ===

func orderSubHandler(svc *service.Service, build func(req mcp.CallToolRequest) (orderID string, body map[string]any, err error)) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, body, err := build(req)
		if err != nil {
			return toolError(err), nil
		}
		return updateWithDiff("order", orderID, body, nil,
			func() (any, error) { return svc.GetOrderDetails(ctx, orderID) },
			func() (any, error) { return svc.UpdateOrder(ctx, orderID, body) },
			func(v any) any { return applyPDNToResult(ctx, v, "update_order") },
		), nil
	}
}

func strArg(req mcp.CallToolRequest, key string) string {
	return req.GetString(key, "")
}

func floatArg(req mcp.CallToolRequest, key string) (float64, bool) {
	args := req.GetArguments()
	v, ok := args[key].(float64)
	return v, ok
}

func addOrderLineHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		productID := strArg(req, "product_id")
		variantID := strArg(req, "variant_id")
		if productID == "" && variantID == "" {
			return "", nil, errString("product_id or variant_id is required")
		}
		if productID != "" && variantID != "" {
			return "", nil, errString("pass only one of product_id or variant_id")
		}
		line := map[string]any{}
		if productID != "" {
			line["product_id"] = productID
		}
		if variantID != "" {
			line["variant_id"] = variantID
		}
		if q, ok := floatArg(req, "quantity"); ok {
			line["quantity"] = q
		} else {
			return "", nil, errString("quantity is required")
		}
		if v, ok := floatArg(req, "vat"); ok {
			line["vat"] = v
		}
		if p, ok := floatArg(req, "sale_price"); ok {
			line["sale_price"] = p
		}
		return orderID, map[string]any{"order_lines_attributes": []any{line}}, nil
	})
}

func updateOrderLineHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		lineID, err := req.RequireString("line_id")
		if err != nil {
			return "", nil, err
		}
		line := map[string]any{"id": lineID}
		if q, ok := floatArg(req, "quantity"); ok {
			line["quantity"] = q
		}
		if v, ok := floatArg(req, "vat"); ok {
			line["vat"] = v
		}
		if p, ok := floatArg(req, "sale_price"); ok {
			line["sale_price"] = p
		}
		if c := strArg(req, "comment"); c != "" {
			line["comment"] = c
		}
		return orderID, map[string]any{"order_lines_attributes": []any{line}}, nil
	})
}

func removeOrderLineHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		lineID, err := req.RequireString("line_id")
		if err != nil {
			return "", nil, err
		}
		return orderID, map[string]any{
			"order_lines_attributes": []any{
				map[string]any{"id": lineID, "_destroy": true},
			},
		}, nil
	})
}

func addOrderDiscountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		d, ok := floatArg(req, "discount")
		if !ok {
			return "", nil, errString("discount is required")
		}
		item := map[string]any{"discount": d}
		if t, ok := floatArg(req, "type_id"); ok {
			item["type_id"] = t
		}
		if desc := strArg(req, "description"); desc != "" {
			item["description"] = desc
		}
		if ids, ok := req.GetArguments()["line_ids"].([]any); ok && len(ids) > 0 {
			item["discount_order_lines_ids"] = ids
		}
		return orderID, map[string]any{"discounts_attributes": []any{item}}, nil
	})
}

func updateOrderDiscountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		discountID, err := req.RequireString("discount_id")
		if err != nil {
			return "", nil, err
		}
		item := map[string]any{"id": discountID}
		if d, ok := floatArg(req, "discount"); ok {
			item["discount"] = d
		}
		if t, ok := floatArg(req, "type_id"); ok {
			item["type_id"] = t
		}
		if desc := strArg(req, "description"); desc != "" {
			item["description"] = desc
		}
		return orderID, map[string]any{"discounts_attributes": []any{item}}, nil
	})
}

func removeOrderDiscountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		discountID, err := req.RequireString("discount_id")
		if err != nil {
			return "", nil, err
		}
		return orderID, map[string]any{
			"discounts_attributes": []any{
				map[string]any{"id": discountID, "_destroy": true},
			},
		}, nil
	})
}

func updateOrderShippingAddressHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		addr := map[string]any{}
		args := req.GetArguments()
		for _, k := range []string{"address", "full_locality_name", "country", "region", "city", "zip", "phone", "full_name", "comment"} {
			if v, ok := args[k]; ok && v != nil && v != "" {
				addr[k] = v
			}
		}
		if len(addr) == 0 {
			return "", nil, errString("at least one address field is required")
		}
		return orderID, map[string]any{"shipping_address_attributes": addr}, nil
	})
}

func updateOrderCustomStatusHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return orderSubHandler(svc, func(req mcp.CallToolRequest) (string, map[string]any, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return "", nil, err
		}
		permalink, err := req.RequireString("custom_status_permalink")
		if err != nil {
			return "", nil, err
		}
		body := map[string]any{"custom_status_permalink": permalink}
		if fs := strArg(req, "fulfillment_status"); fs != "" {
			body["fulfillment_status"] = fs
		}
		return orderID, body, nil
	})
}

type errString string

func (e errString) Error() string { return string(e) }

func createMenuItemHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		menuID, err := req.RequireString("menu_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		delete(args, "menu_id")
		rel, _ := args["related_object_id"].(string)
		if rel == "" {
			return mcp.NewToolResultError("related_object_id is required (Blog_1, Page_2, Collection_3, or URL)"), nil
		}
		if !looksLikeRelatedObjectID(rel) {
			return mcp.NewToolResultError(`related_object_id must be composite like "Blog_1", "Page_2", "Collection_3", or an absolute URL`), nil
		}
		result, err := svc.CreateMenuItem(ctx, menuID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func looksLikeRelatedObjectID(rel string) bool {
	if strings.HasPrefix(rel, "http://") || strings.HasPrefix(rel, "https://") {
		return true
	}
	// Type_id composite: Blog_1, Page_2, Collection_3, ...
	i := strings.IndexByte(rel, '_')
	if i <= 0 || i == len(rel)-1 {
		return false
	}
	for _, r := range rel[i+1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// === OPTION VALUES handlers (nested resource) ===

func genericListOptionValuesHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListOptionValues(ctx, optionNameID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericGetOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		optionValueID, err := req.RequireString("option_value_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.GetOptionValue(ctx, optionNameID, optionValueID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericCreateOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		result, err := svc.CreateOptionValue(ctx, optionNameID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func genericUpdateOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		optionValueID, err := req.RequireString("option_value_id")
		if err != nil {
			return toolError(err), nil
		}
		args := req.GetArguments()
		return updateWithDiff("option_value", optionValueID, args, nil,
			func() (any, error) { return svc.GetOptionValue(ctx, optionNameID, optionValueID) },
			func() (any, error) { return svc.UpdateOptionValue(ctx, optionNameID, optionValueID, args) },
			nil,
		), nil
	}
}

func genericDeleteOptionValueHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		optionNameID, err := req.RequireString("option_name_id")
		if err != nil {
			return toolError(err), nil
		}
		optionValueID, err := req.RequireString("option_value_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("option_value", optionValueID,
			nil,
			func() error { return svc.DeleteOptionValue(ctx, optionNameID, optionValueID) },
		), nil
	}
}

// === ROUND 4: Special handlers ===

func declineApplicationChargeHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("charge_id")
		if err != nil {
			return toolError(err), nil
		}
		return deleteResult("application_charge", id,
			nil,
			func() error { return svc.DeclineApplicationCharge(ctx, id) },
		), nil
	}
}

func createBonusSystemTransactionHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		clientID := int(getFloat(args, "client_id", 0))
		if clientID == 0 {
			return toolError(fmt.Errorf("client_id is required")), nil
		}
		result, err := svc.CreateBonusSystemTransaction(ctx, clientID, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listClientBonusTransactionsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListClientBonusSystemTransactions(ctx, clientID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func listSimilarsByProductHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		productID, err := req.RequireString("product_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ListSimilars(ctx, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === PROVIDER TOKEN handler ===

func createProviderTokenHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		messengerType, err := req.RequireString("messenger_type")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CreateProviderToken(ctx, messengerType)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === ACCOUNT handler ===

func updateAccountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		// GetAccount в API нет — diff без old.
		return updateWithDiff("account", "", args, nil, nil,
			func() (any, error) { return svc.UpdateAccount(ctx, args) },
			nil,
		), nil
	}
}

// === COUNT handlers ===

func getOrdersCountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		count, err := svc.GetOrdersCount(ctx, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(map[string]any{"count": count}), nil
	}
}

func getProductsCountHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		count, err := svc.GetProductsCount(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(map[string]any{"count": count}), nil
	}
}

// === PIPELINE HANDLERS ===

func searchProductDetailsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		query, _ := args["query"].(string)
		productID, _ := args["product_id"].(string)
		if query == "" && productID == "" {
			return toolError(fmt.Errorf("query или product_id обязательны")), nil
		}
		result, err := svc.SearchProductDetails(ctx, query, productID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func orderFullDetailsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, err := req.RequireString("order_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.OrderFullDetails(ctx, orderID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(applyPDNToResult(ctx, result, "order_full_details")), nil
	}
}

func clientFullProfileHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		clientID, err := req.RequireString("client_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.ClientFullProfile(ctx, clientID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(applyPDNToResult(ctx, result, "client_full_profile")), nil
	}
}

func productCatalogSearchHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		query, _ := args["query"].(string)
		result, err := svc.ProductCatalogSearch(ctx, query, args)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

func collectionProductsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collectionID, err := req.RequireString("collection_id")
		if err != nil {
			return toolError(err), nil
		}
		result, err := svc.CollectionProducts(ctx, collectionID)
		if err != nil {
			return toolError(err), nil
		}
		return toolResult(result), nil
	}
}

// === EXPORT handlers ===

// pdnActive — включён ли режим маскирования ПДн для этой сессии.
func pdnActive(ctx context.Context) bool {
	shopPDN, _ := ctx.Value(mcpShopPDNKey).(bool)
	aiKeyPDN, _ := ctx.Value(mcpAIKeyPDNKey).(bool)
	return shopPDN || aiKeyPDN
}

func exportProductsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		csv, err := svc.ExportProducts(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText(csv), nil
	}
}

func exportOrdersHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		csv, err := svc.ExportOrders(ctx)
		if err != nil {
			return toolError(err), nil
		}
		if pdnActive(ctx) {
			csv = pdn.MaskCSV(csv)
		}
		return mcp.NewToolResultText(csv), nil
	}
}

func exportClientsHandler(svc *service.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		csv, err := svc.ExportClients(ctx)
		if err != nil {
			return toolError(err), nil
		}
		if pdnActive(ctx) {
			csv = pdn.MaskCSV(csv)
		}
		return mcp.NewToolResultText(csv), nil
	}
}

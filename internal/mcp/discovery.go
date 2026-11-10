package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/service"
	"github.com/chubuk/mcp-is/internal/storage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// coreToolNames is the default visible tools/list surface (progressive discovery).
// Long-tail tools stay callable via call_insales_tool / search_insales_tools.
var coreToolNames = map[string]bool{
	// pipelines
	"search_product_details": true,
	"order_full_details":     true,
	"client_full_profile":    true,
	"product_catalog_search": true,
	"collection_products":    true,
	// catalog / core entities
	"list_products":          true,
	"get_product":            true,
	"get_seo":                true,
	"update_seo":             true,
	"list_orders":            true,
	"get_order":              true,
	"get_orders_count":       true,
	"list_clients":           true,
	"get_client":             true,
	"list_categories":        true,
	"get_category":           true,
	"list_collections":       true,
	"get_collection":         true,
	"get_shop_info":          true,
	"list_variants":          true,
	"get_variant":            true,
	"list_images":            true,
	"list_product_fields":    true,
	"get_product_field":      true,
	"list_delivery_variants": true,
	"list_payment_gateways":  true,
	// export
	"export_products": true,
	"export_orders":   true,
	"export_clients":  true,
	// товары: базовая запись (W5, бюджет core ≤50)
	"create_product": true,
	"update_product": true,
	"delete_product": true,
	"add_image":      true,
	"update_variant": true,
	// публикации (блог): полный цикл постов (W5; delete_blog остаётся long-tail)
	"list_blogs":       true,
	"get_blog":         true,
	"create_blog":      true,
	"update_blog":      true,
	"list_blog_posts":  true,
	"get_blog_post":    true,
	"create_blog_post": true,
	"update_blog_post": true,
	"delete_blog_post": true,
	// заказы: основные записи (W5)
	"create_order":   true,
	"update_order":   true,
	"add_order_line": true,
	// discovery meta-tools
	"search_insales_tools": true,
	"call_insales_tool":    true,
	"store_info":           true,
}

const (
	toolsModeCore = "core"
	toolsModeAll  = "all"
)

// retiredTools — удалённые тулы с человеческим объяснением взамен
// «unknown tool» (W7; каталог v2026-11-10: мёртвые эндпоинты InSales).
var retiredTools = map[string]string{
	"list_meta_fields":  "Тул удалён: /admin/meta_fields.json в API InSales не существует (возвращает HTML). Доп. поля товара — list_product_fields / list_product_field_values.",
	"get_meta_field":    "Тул удалён: /admin/meta_fields.json в API InSales не существует. Доп. поля товара — get_product_field.",
	"create_meta_field": "Тул удалён: эндпоинта нет в API InSales. Создавайте поля — create_product_field.",
	"update_meta_field": "Тул удалён: эндпоинта нет в API InSales. Поля — update_product_field.",
	"delete_meta_field": "Тул удалён: эндпоинта нет в API InSales. Поля — delete_product_field.",
	"list_currencies":   "Тул удалён: /admin/currencies.json не существует в API InSales. Валюты учёта — list_stock_currencies.",
	"get_supplementary": "Тул удалён: одиночного GET у supplementaries в API нет. Список — list_supplementaries.",
}

func toolsMode() string {
	if strings.EqualFold(os.Getenv("MCP_TOOLS_MODE"), toolsModeAll) {
		return toolsModeAll
	}
	return toolsModeCore
}

// isCoreTool reports whether the tool belongs in the default visible surface.
func isCoreTool(name string) bool {
	return coreToolNames[name]
}

// toolHiddenReason — причина, по которой тул скрыт для этого ctx ("" = видим).
// Единый источник для list-фильтра и guard'а вызовов: логика не должна расходиться.
// Цепочка приоритетов (план D4): cfg.ReadOnly → capabilities (denied) →
// is_read_only → delete-deny → политика ИИ-ключа → core/all.
func toolHiddenReason(ctx context.Context, name string, isWrite bool, globalReadOnly bool) string {
	readOnly := globalReadOnly
	if v, ok := ctx.Value(mcpReadOnlyKey).(bool); ok {
		readOnly = readOnly || v
	}
	if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
		readOnly = readOnly || v
	}
	if readOnly && isWrite {
		return "readonly"
	}
	deleteDenied := false
	if v, ok := ctx.Value(core.ShopDeleteDeniedKey).(bool); ok {
		deleteDenied = v
	}
	if deleteDenied && ToolMetaFor(name).Class == ClassDelete {
		return "delete_denied"
	}
	if v, ok := ctx.Value(auth.CapabilitiesJSONKey).(string); ok && v != "" {
		if snap := core.ParseSnapshot(v); snap != nil {
			domain := ToolDomain(name)
			if domain != "" && core.DomainStatus(snap, domain) == core.CapStatusDenied {
				return "no_scope"
			}
		}
	}
	if v, ok := ctx.Value(auth.CommandPolicyJSONKey).(string); ok && v != "" {
		if policy := parseCommandPolicy(v); policy != nil && !policy.IsCommandAllowed(name) {
			return "policy"
		}
	}
	if toolsMode() == toolsModeCore && !isCoreTool(name) {
		return "not_core"
	}
	return ""
}

// hiddenToolError — дружелюбная ошибка вместо голого «tool not found»: агент
// должен понимать, что тул существует, но скрыт, и как его вызвать (B5/hermes).
func hiddenToolError(name, reason string) *mcp.CallToolResult {
	var msg string
	switch reason {
	case "readonly":
		msg = fmt.Sprintf("Инструмент %q скрыт: сессия только чтение (is_read_only). Для записи снимите ограничение в настройках магазина.", name)
	case "delete_denied":
		msg = fmt.Sprintf("Инструмент %q скрыт: удаление запрещено настройками магазина (is_delete_denied).", name)
	case "no_scope":
		msg = fmt.Sprintf("Инструмент %q скрыт: нет прав ключа на ресурс %q — включите доступ в настройках приложения InSales (см. store_info).", name, ToolDomain(name))
	case "policy":
		msg = fmt.Sprintf("Инструмент %q отключён политикой ИИ-ключа — измените политику в настройках ключа.", name)
	case "not_core":
		visible, total := catalogCounts(context.Background())
		msg = fmt.Sprintf("Инструмент %q существует, но скрыт текущим списком (показано %d из %d). Вызовите через call_insales_tool(name=%q, args={...}) или найдите через search_insales_tools.", name, visible, total, name)
	default:
		msg = fmt.Sprintf("Инструмент %q скрыт.", name)
	}
	return mcp.NewToolResultError(msg)
}

// toolFilter hides write tools in read-only mode and applies progressive
// discovery (core surface) unless MCP_TOOLS_MODE=all.
func toolFilter(globalReadOnly bool) server.ToolFilterFunc {
	filterGlobalReadOnly = globalReadOnly
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		tools = augmentFieldValuesSchemas(ctx, tools)
		// mcp-go при tools/call прогоняет через фильтр только запрошенный тул
		// (passesToolFilters, len==1) — пропускаем его, чтобы handler вернул
		// понятную ошибку видимости вместо голого «tool not found».
		if len(tools) == 1 {
			return tools
		}
		filtered := make([]mcp.Tool, 0, len(tools))
		for _, t := range tools {
			isWrite := t.Annotations.ReadOnlyHint == nil || !*t.Annotations.ReadOnlyHint
			if toolHiddenReason(ctx, t.Name, isWrite, globalReadOnly) != "" {
				continue
			}
			filtered = append(filtered, t)
		}
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].Name < filtered[j].Name })
		return filtered
	}
}

// parseCommandPolicy разбирает JSON политики; кривой JSON → deny-all (fail-closed).
func parseCommandPolicy(raw string) *CommandPolicy {
	var p CommandPolicy
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		slog.Warn("corrupt command_policy, fail-closed", "error", err)
		return &CommandPolicy{Default: "deny"}
	}
	return &p
}

// RegisterDiscoveryTools adds search_insales_tools / call_insales_tool.
func RegisterDiscoveryTools(s *server.MCPServer, _ *service.Service, auditRepo *storage.AuditRepository) {
	addTool(s, mcp.NewTool("search_insales_tools",
		mcp.WithDescription("Search insales tools and return matching data in one call. Pipeline tool: combines lookup steps for common workflows. Returns JSON items relevant to the query."),
		mcp.WithToolTitle("Поиск инструментов InSales"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Keyword or tool name fragment, e.g. 'discount' or 'menu'")),
		mcp.WithString("resource", mcp.Description("Optional resource filter: product, order, client, collection, category, menu, delivery, payment, field")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), searchInSalesToolsHandler(s))

	addTool(s, mcp.NewTool("call_insales_tool",
		mcp.WithDescription("Manage call insales tool via InSales API. Execute any InSales MCP tool by name, including tools not listed by default. Use search_insales_tools first to discover the name. Arguments must use API language-invariant values (IDs, handles, enum codes)."),
		mcp.WithToolTitle("Вызов инструмента InSales"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Exact tool name from search_insales_tools")),
		mcp.WithObject("args", mcp.Description("Arguments object for the tool, same as calling it directly")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), callInSalesToolHandler(s, auditRepo))

	addTool(s, mcp.NewTool("store_info",
		mcp.WithDescription("Identify the store this connection is bound to. Returns connection alias, store name, domain and binding state — the ground truth before any write. All other tools operate on this store only."),
		mcp.WithToolTitle("Магазин этого соединения"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), storeInfoHandler())
}

// storeInfoHandler reports the connection identity from auth context — no InSales call.
func storeInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		shopID, _ := ctx.Value(auth.ShopIDKey).(string)
		shopName, _ := ctx.Value(auth.ShopNameKey).(string)
		connName, _ := ctx.Value(auth.ConnectionNameKey).(string)
		info := map[string]any{
			"connection_name": effectiveAlias(connName, shopName),
			"key_name":        ctx.Value(auth.KeyNameKey),
			"bound":           shopID != "",
			"store": map[string]any{
				"id":       shopID,
				"name":     shopName,
				"shop_url": ctx.Value(auth.ShopURLKey),
			},
			"read_only": false,
		}
		if v, ok := ctx.Value(mcpReadOnlyKey).(bool); ok {
			info["read_only"] = v
		}
		// Каталог инструментов (W6/W7): честные счётчики per-ctx + оглавление доменов.
		visible, total := catalogCounts(ctx)
		domainIndex := map[string]int{}
		for _, m := range registeredToolMeta {
			if d := ToolDomain(m.Name); d != "" {
				domainIndex[d]++
			}
		}
		info["tools"] = map[string]any{
			"total":   total,
			"visible": visible,
			"mode":    toolsMode(),
		}
		info["domain_index"] = domainIndex
		info["catalog_version"] = catalogVersion
		info["catalog_date"] = catalogDate
		info["catalog_changes"] = catalogChanges
		// Capabilities-блок: права ключа магазина (из снимка probe).
		if v, ok := ctx.Value(auth.CapabilitiesJSONKey).(string); ok && v != "" {
			snap := core.ParseSnapshot(v)
			if snap != nil {
				var checkedAt string
				summary := core.BuildSummary(snap, nil)
				checkedAt = summary.CheckedAt
				// Детали по доменам (включая since/source) — чтобы были понятны причины отказа.
				domainDetails := make([]map[string]any, 0, len(core.ProbeDomains))
				for _, dp := range core.ProbeDomains {
					st := core.DomainStatus(snap, dp.Domain)
					dc := snap.Domains[dp.Domain]
					domainDetails = append(domainDetails, map[string]any{
						"domain": dp.Domain,
						"status": st,
						"since":  dc.Since,
						"source": dc.Source,
					})
				}
				info["capabilities"] = map[string]any{
					"checked_at":     checkedAt,
					"stale":          summary.Stale,
					"credentials_ok": snap.CredentialsOK,
					"allowed":        summary.AllowedDomains,
					"denied":         summary.DeniedDomains,
					"unknown":        summary.UnknownDomains,
					"domains":        domainDetails,
					"error":          snap.Error,
				}
			}
		}
		if shopID == "" {
			info["message"] = "AI-ключ не привязан к магазину. Попросите владельца привязать ключ к магазину в личном кабинете (раздел «AI-ключи») и повторите запрос."
		}
		info["direct_rest"] = "Прямой REST к InSales через этот ключ недоступен — используйте инструменты или call_insales_tool."
		return toolResult(info), nil
	}
}

func searchInSalesToolsHandler(s *server.MCPServer) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := strings.ToLower(strings.TrimSpace(req.GetString("query", "")))
		resource := strings.ToLower(strings.TrimSpace(req.GetString("resource", "")))
		if query == "" && resource == "" {
			return mcp.NewToolResultError("query or resource is required"), nil
		}

		readOnly := false
		if v, ok := ctx.Value(mcpReadOnlyKey).(bool); ok {
			readOnly = v
		}
		if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
			readOnly = readOnly || v
		}
		deleteDenied := false
		if v, ok := ctx.Value(core.ShopDeleteDeniedKey).(bool); ok {
			deleteDenied = v
		}

		all := s.ListTools()
		matches := make([]map[string]any, 0, 20)
		for name, st := range all {
			if st == nil {
				continue
			}
			desc := st.Tool.Description
			title := st.Tool.Title
			blob := strings.ToLower(name + " " + desc + " " + title)
			if query != "" && !strings.Contains(blob, query) {
				continue
			}
			if resource != "" && !strings.Contains(blob, resource) {
				continue
			}
			isWrite := st.Tool.Annotations.ReadOnlyHint == nil || !*st.Tool.Annotations.ReadOnlyHint
			if readOnly && isWrite {
				continue
			}
			if deleteDenied && ToolMetaFor(name).Class == ClassDelete {
				continue
			}
			short := desc
			if len(short) > 160 {
				short = short[:157] + "..."
			}
			item := map[string]any{
				"name":        name,
				"description": short,
				"read_only":   !isWrite,
				"core":        isCoreTool(name),
			}
			if title != "" {
				item["title"] = title
			}
			matches = append(matches, item)
		}
		sort.Slice(matches, func(i, j int) bool {
			return matches[i]["name"].(string) < matches[j]["name"].(string)
		})
		_, total := catalogCounts(ctx)
		if len(matches) == 0 {
			return mcp.NewToolResultText(formatResult(map[string]any{
				"items":         []any{},
				"total":         0,
				"catalog_total": total,
				"hint":          fmt.Sprintf("No tools matched. Try a shorter query or omit resource. Каталог: %d инструментов; оглавление доменов — store_info.", total),
			})), nil
		}
		return mcp.NewToolResultText(formatResult(map[string]any{
			"items":         matches,
			"total":         len(matches),
			"catalog_total": total,
			"hint":          fmt.Sprintf("Найдено %d из %d. Call the tool via call_insales_tool with name and args, or call it directly if present in tools/list.", len(matches), total),
		})), nil
	}
}

func callInSalesToolHandler(s *server.MCPServer, auditRepo *storage.AuditRepository) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := strings.TrimSpace(req.GetString("name", ""))
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		if msg, gone := retiredTools[name]; gone {
			return mcp.NewToolResultError(msg), nil
		}
		st := s.GetTool(name)
		if st == nil || st.Handler == nil {
			return mcp.NewToolResultError(fmt.Sprintf("unknown tool %q; use search_insales_tools to find the name", name)), nil
		}

		// block write tools under read-only (same rule as toolFilter)
		readOnly := false
		if v, ok := ctx.Value(mcpReadOnlyKey).(bool); ok {
			readOnly = v
		}
		if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
			readOnly = readOnly || v
		}
		isWrite := st.Tool.Annotations.ReadOnlyHint == nil || !*st.Tool.Annotations.ReadOnlyHint
		if readOnly && isWrite {
			if auditRepo != nil {
				args := map[string]any{}
				if raw, ok := req.GetArguments()["args"].(map[string]any); ok {
					args = raw
				}
				recordDenied(auditRepo, ctx, name, args)
			}
			return mcp.NewToolResultError(fmt.Sprintf("tool %q is write; session is read-only", name)), nil
		}

		// block delete tools when the shop denies deletion (same rule as toolFilter)
		deleteDenied := false
		if v, ok := ctx.Value(core.ShopDeleteDeniedKey).(bool); ok {
			deleteDenied = v
		}
		if deleteDenied && ToolMetaFor(name).Class == ClassDelete {
			if auditRepo != nil {
				args := map[string]any{}
				if raw, ok := req.GetArguments()["args"].(map[string]any); ok {
					args = raw
				}
				recordDenied(auditRepo, ctx, name, args)
			}
			return mcp.NewToolResultError(fmt.Sprintf("tool %q deletes data; shop denies deletion", name)), nil
		}

		// Политика ИИ-ключа и capabilities: эффективное состояние = f(policy, capability);
		// deny побеждает, unavailable_unknown — fail-open (выполняем).
		policyOn := true
		if v, ok := ctx.Value(auth.CommandPolicyJSONKey).(string); ok && v != "" {
			policyOn = parseCommandPolicy(v).IsCommandAllowed(name)
		}
		capStatus := core.CapStatusUnknown
		var capSince, capSource string
		if v, ok := ctx.Value(auth.CapabilitiesJSONKey).(string); ok && v != "" {
			if snap := core.ParseSnapshot(v); snap != nil {
				if domain := ToolDomain(name); domain != "" {
					capStatus = core.DomainStatus(snap, domain)
					dc := snap.Domains[domain]
					capSince, capSource = dc.Since, dc.Source
				}
			}
		}
		switch core.EffectiveStatus(policyOn, capStatus) {
		case core.EffOffByPolicy:
			return mcp.NewToolResultError(fmt.Sprintf(
				"Инструмент %q отключён политикой ИИ-ключа — измените политику в настройках ключа", name)), nil
		case core.EffUnavailableNoScope:
			msg := fmt.Sprintf(
				"Нет прав ключа на ресурс %q — включите доступ в настройках приложения InSales; см. store_info", ToolDomain(name))
			if capSince != "" || capSource != "" {
				msg += fmt.Sprintf(" (с %s, источник: %s)", capSince, capSource)
			}
			return mcp.NewToolResultError(msg), nil
		}

		args := map[string]any{}
		if raw, ok := req.GetArguments()["args"]; ok {
			switch t := raw.(type) {
			case map[string]any:
				args = t
			case nil:
			default:
				return mcp.NewToolResultError("args must be an object"), nil
			}
		}

		sub := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      name,
				Arguments: args,
			},
		}
		ctx = context.WithValue(ctx, mcpNestedCallKey, true) // guard видимости пропускает nested-вызов
		if isWrite {
			// nested dispatch bypasses server middleware — audit write tools explicitly
			res, err := st.Handler(ctx, sub)
			action, ok := "call", true
			if err != nil || (res != nil && res.IsError) {
				action, ok = "error", false
			}
			if auditRepo != nil {
				if err := auditRepo.Insert(ctx, auditFromContext(ctx, name, action, args, ok)); err != nil {
					slog.Error("audit insert", "error", err, "tool", name)
				}
			}
			return res, err
		}
		return st.Handler(ctx, sub)
	}
}

package mcp

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/service"
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
	// discovery meta-tools
	"search_insales_tools": true,
	"call_insales_tool":    true,
}

const (
	toolsModeCore = "core"
	toolsModeAll  = "all"
)

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

// toolFilter hides write tools in read-only mode and applies progressive
// discovery (core surface) unless MCP_TOOLS_MODE=all.
func toolFilter(globalReadOnly bool) server.ToolFilterFunc {
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		readOnly := globalReadOnly
		if v, ok := ctx.Value(mcpReadOnlyKey).(bool); ok {
			readOnly = readOnly || v
		}
		if v, ok := ctx.Value(core.ShopReadOnlyKey).(bool); ok {
			readOnly = readOnly || v
		}

		mode := toolsMode()
		filtered := make([]mcp.Tool, 0, len(tools))
		for _, t := range tools {
			// nil ReadOnlyHint is treated as write (no leak into read-only list)
			isWrite := t.Annotations.ReadOnlyHint == nil || !*t.Annotations.ReadOnlyHint
			if readOnly && isWrite {
				continue
			}
			if mode == toolsModeCore && !isCoreTool(t.Name) {
				continue
			}
			filtered = append(filtered, t)
		}
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].Name < filtered[j].Name })
		return filtered
	}
}

// RegisterDiscoveryTools adds search_insales_tools / call_insales_tool.
func RegisterDiscoveryTools(s *server.MCPServer, _ *service.Service) {
	s.AddTool(mcp.NewTool("search_insales_tools",
		mcp.WithDescription("Find MCP tools by keyword when the tool you need is not in the default list. Returns matching tool names and short descriptions. Then call_insales_tool(name, args)."),
		mcp.WithToolTitle("Поиск инструментов InSales"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Keyword or tool name fragment, e.g. 'discount' or 'menu'")),
		mcp.WithString("resource", mcp.Description("Optional resource filter: product, order, client, collection, category, menu, delivery, payment, field")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), searchInSalesToolsHandler(s))

	s.AddTool(mcp.NewTool("call_insales_tool",
		mcp.WithDescription("Execute any InSales MCP tool by name, including tools not listed by default. Use search_insales_tools first to discover the name. Arguments must use API language-invariant values (IDs, handles, enum codes)."),
		mcp.WithToolTitle("Вызов инструмента InSales"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Exact tool name from search_insales_tools")),
		mcp.WithObject("args", mcp.Description("Arguments object for the tool, same as calling it directly")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(true),
	), callInSalesToolHandler(s))
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
		if len(matches) == 0 {
			return mcp.NewToolResultText(formatResult(map[string]any{
				"items": []any{},
				"total": 0,
				"hint":  "No tools matched. Try a shorter query or omit resource.",
			})), nil
		}
		return mcp.NewToolResultText(formatResult(map[string]any{
			"items": matches,
			"total": len(matches),
			"hint":  "Call the tool via call_insales_tool with name and args, or call it directly if present in tools/list.",
		})), nil
	}
}

func callInSalesToolHandler(s *server.MCPServer) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := strings.TrimSpace(req.GetString("name", ""))
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
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
			return mcp.NewToolResultError(fmt.Sprintf("tool %q is write; session is read-only", name)), nil
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
		return st.Handler(ctx, sub)
	}
}

package mcp

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func newDiscoveryTestServer() *server.MCPServer {
	s := server.NewMCPServer("test", "1.0")
	s.AddTool(mcp.NewTool("create_menu",
		mcp.WithDescription("Create a navigation menu item"),
		mcp.WithReadOnlyHintAnnotation(false),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		b, _ := json.Marshal(args)
		return mcp.NewToolResultText("created:" + string(b)), nil
	})
	s.AddTool(mcp.NewTool("list_menus",
		mcp.WithDescription("List navigation menus"),
		mcp.WithReadOnlyHintAnnotation(true),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(`{"items":[],"total":0}`), nil
	})
	RegisterDiscoveryTools(s, nil, nil)
	return s
}

func TestSearchInSalesTools(t *testing.T) {
	s := newDiscoveryTestServer()
	h := searchInSalesToolsHandler(s)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "search_insales_tools",
			Arguments: map[string]any{"query": "menu"},
		},
	}
	res, err := h(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil || len(res.Content) == 0 {
		t.Fatal("empty result")
	}
	text := res.Content[0].(mcp.TextContent).Text
	if !containsAll(text, "create_menu", "list_menus") {
		t.Errorf("search result missing tools: %s", text)
	}
}

func TestSearchInSalesToolsHidesWriteInReadOnly(t *testing.T) {
	s := newDiscoveryTestServer()
	h := searchInSalesToolsHandler(s)
	ctx := context.WithValue(context.Background(), mcpReadOnlyKey, true)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "search_insales_tools",
			Arguments: map[string]any{"query": "menu"},
		},
	}
	res, err := h(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	text := res.Content[0].(mcp.TextContent).Text
	if containsAll(text, "create_menu") {
		t.Errorf("write tool must be hidden in read-only: %s", text)
	}
	if !containsAll(text, "list_menus") {
		t.Errorf("read tool must stay visible: %s", text)
	}
}

func TestCallInSalesTool(t *testing.T) {
	s := newDiscoveryTestServer()
	h := callInSalesToolHandler(s, nil)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "call_insales_tool",
			Arguments: map[string]any{
				"name": "create_menu",
				"args": map[string]any{"title": "Main"},
			},
		},
	}
	res, err := h(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	text := res.Content[0].(mcp.TextContent).Text
	if !containsAll(text, "created:", "Main") {
		t.Errorf("call_insales_tool did not dispatch: %s", text)
	}
}

func TestCallInSalesToolBlocksWriteInReadOnly(t *testing.T) {
	s := newDiscoveryTestServer()
	h := callInSalesToolHandler(s, nil)
	ctx := context.WithValue(context.Background(), mcpReadOnlyKey, true)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "call_insales_tool",
			Arguments: map[string]any{
				"name": "create_menu",
				"args": map[string]any{"title": "Main"},
			},
		},
	}
	res, err := h(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Errorf("expected error for write in read-only, got %v", res)
	}
}

func TestCallInSalesToolUnknown(t *testing.T) {
	s := newDiscoveryTestServer()
	h := callInSalesToolHandler(s, nil)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "call_insales_tool",
			Arguments: map[string]any{"name": "nope_tool"},
		},
	}
	res, err := h(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Errorf("expected unknown-tool error, got %v", res)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !jsonContains(s, p) {
			return false
		}
	}
	return true
}

func jsonContains(s, sub string) bool {
	return len(s) >= len(sub) && (stringIndex(s, sub) >= 0)
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestToolResultDualPayload(t *testing.T) {
	payload := map[string]any{"items": []any{map[string]any{"id": "1"}}, "total": 1}
	res := toolResult(payload)
	if res.StructuredContent == nil {
		t.Fatal("StructuredContent missing")
	}
	if len(res.Content) == 0 {
		t.Fatal("text content missing")
	}
	text := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(text, "\"total\"") {
		t.Errorf("text mirror missing fields: %s", text)
	}
}

func TestLooksLikeRelatedObjectID(t *testing.T) {
	ok := []string{"Blog_1", "Page_2", "Collection_3", "https://example.com", "http://x.ru/a"}
	bad := []string{"", "Blog", "Blog_", "_1", "Page_xx", "just text"}
	for _, s := range ok {
		if !looksLikeRelatedObjectID(s) {
			t.Errorf("want valid: %q", s)
		}
	}
	for _, s := range bad {
		if looksLikeRelatedObjectID(s) {
			t.Errorf("want invalid: %q", s)
		}
	}
}

func TestCoreToolsHaveOutputSchema(t *testing.T) {
	want := []string{
		"list_products", "get_product", "list_orders", "get_order",
		"list_clients", "get_shop_info", "search_product_details",
		"order_full_details", "client_full_profile", "product_catalog_search",
		"collection_products",
	}
	s := server.NewMCPServer("test", "1.0")
	RegisterTools(s, nil)
	for _, name := range want {
		st := s.GetTool(name)
		if st == nil || st.Tool.Name == "" {
			t.Fatalf("tool %s not registered", name)
		}
		if st.Tool.OutputSchema.Type != "object" || len(st.Tool.OutputSchema.Properties) == 0 {
			t.Errorf("%s: outputSchema missing/empty: %+v", name, st.Tool.OutputSchema)
		}
	}
}

func TestListProductsEnvelope(t *testing.T) {
	got := ProductsPage{Items: nil, Total: 0, Page: 2, PerPage: 10}
	b, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"items", "total", "page", "per_page"} {
		if _, ok := m[k]; !ok {
			t.Errorf("missing key %q in %s", k, b)
		}
	}
}

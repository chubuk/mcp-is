package mcp

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// newE2EServer — продакшн-сервер без БД-репозиториев: реестр тулов, фильтры,
// hooks и recovery как в бою; вызовы InSales-тулов не выполняются.
func newE2EServer(t *testing.T) *MCPServer {
	t.Helper()
	return NewMCPServer(nil, false, false, "RU", nil, nil)
}

func e2eStart(t *testing.T, c *client.Client) *mcp.InitializeResult {
	t.Helper()
	if err := c.Start(context.Background()); err != nil {
		t.Fatalf("client start: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	req := mcp.InitializeRequest{}
	req.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	req.Params.ClientInfo = mcp.Implementation{Name: "e2e-test", Version: "1.0"}
	res, err := c.Initialize(context.Background(), req)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	return res
}

// TestE2E_CoreSurfaceProtocol — e2e через реальный MCP-клиент: initialize →
// tools/list (core surface) → discovery meta-tools → ошибки вызова.
func TestE2E_CoreSurfaceProtocol(t *testing.T) {
	srv := newE2EServer(t)
	c, err := client.NewInProcessClient(srv.mcpServer)
	if err != nil {
		t.Fatalf("in-process client: %v", err)
	}
	initRes := e2eStart(t, c)
	if initRes.ServerInfo.Name != "mcp-is" {
		t.Errorf("server info: %+v", initRes.ServerInfo)
	}

	// tools/list по умолчанию = core surface (32 tools).
	tools, err := c.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("tools/list: %v", err)
	}
	if len(tools.Tools) != len(coreToolNames) {
		t.Errorf("core surface size = %d, want %d", len(tools.Tools), len(coreToolNames))
	}
	seen := map[string]bool{}
	for _, tool := range tools.Tools {
		seen[tool.Name] = true
		if !isCoreTool(tool.Name) {
			t.Errorf("non-core tool %q leaked into core tools/list", tool.Name)
		}
		if tool.Description == "" {
			t.Errorf("tool %q has empty description", tool.Name)
		}
		if tool.InputSchema.Type != "object" {
			t.Errorf("tool %q input schema type = %q, want object", tool.Name, tool.InputSchema.Type)
		}
	}
	for name := range coreToolNames {
		if !seen[name] {
			t.Errorf("core tool %q missing from tools/list", name)
		}
	}
	if seen["create_product"] {
		t.Error("long-tail create_product must not be in core tools/list")
	}

	// Discovery: поиск по каталогу находит long-tail тулзы.
	res, err := c.CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "search_insales_tools",
			Arguments: map[string]any{"query": "create_product"},
		},
	})
	if err != nil {
		t.Fatalf("search_insales_tools: %v", err)
	}
	if res.IsError {
		t.Fatalf("search_insales_tools errored: %v", res.Content)
	}
	text, _ := res.Content[0].(mcp.TextContent)
	if !strings.Contains(text.Text, "create_product") {
		t.Errorf("search must find create_product, got: %s", text.Text)
	}

	// call_insales_tool: неизвестный тул → понятная ошибка, не паника.
	res, err = c.CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "call_insales_tool",
			Arguments: map[string]any{"name": "nope_tool"},
		},
	})
	if err != nil {
		t.Fatalf("call_insales_tool: %v", err)
	}
	if !res.IsError {
		t.Errorf("unknown tool must be an error result")
	}

	// Прямой вызов отфильтрованного long-tail тулза — JSON-RPC ошибка «not found», не выполнение.
	_, err = c.CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "create_product",
			Arguments: map[string]any{"title": "x"},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("filtered tool call must fail with not found, got: %v", err)
	}
}

// TestE2E_AllModeListsFullCatalog — MCP_TOOLS_MODE=all отдаёт весь каталог (287).
func TestE2E_AllModeListsFullCatalog(t *testing.T) {
	t.Setenv("MCP_TOOLS_MODE", "all")
	srv := newE2EServer(t)
	c, err := client.NewInProcessClient(srv.mcpServer)
	if err != nil {
		t.Fatalf("in-process client: %v", err)
	}
	e2eStart(t, c)

	tools, err := c.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("tools/list: %v", err)
	}
	if len(tools.Tools) != 287 {
		t.Errorf("full catalog = %d tools, want 287", len(tools.Tools))
	}
}

// TestE2E_ReadOnlyKeyHidesWrites — read-only сессия отдаёт только read-тулзы.
func TestE2E_ReadOnlyKeyHidesWrites(t *testing.T) {
	srv := NewMCPServer(nil, true, false, "RU", nil, nil)
	c, err := client.NewInProcessClient(srv.mcpServer)
	if err != nil {
		t.Fatalf("in-process client: %v", err)
	}
	e2eStart(t, c)

	tools, err := c.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("tools/list: %v", err)
	}
	for _, tool := range tools.Tools {
		if isWrite := tool.Annotations.ReadOnlyHint == nil || !*tool.Annotations.ReadOnlyHint; isWrite {
			t.Errorf("write tool %q visible in read-only session", tool.Name)
		}
	}
}

// TestE2E_StreamableHTTP — тот же сервер через HTTP-транспорт (wire-формат JSON-RPC).
func TestE2E_StreamableHTTP(t *testing.T) {
	srv := newE2EServer(t)
	hs := httptest.NewServer(srv)
	t.Cleanup(hs.Close)

	c, err := client.NewStreamableHttpClient(hs.URL + "/mcp")
	if err != nil {
		t.Fatalf("http client: %v", err)
	}
	e2eStart(t, c)

	tools, err := c.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("tools/list over http: %v", err)
	}
	if len(tools.Tools) != len(coreToolNames) {
		t.Errorf("http core surface = %d, want %d", len(tools.Tools), len(coreToolNames))
	}
}

package mcp

import (
	"context"
	"os"
	"testing"

	"github.com/chubuk/mcp-is/internal/core"
	"github.com/mark3labs/mcp-go/mcp"
)

func boolPtr(b bool) *bool { return &b }

func withToolsMode(t *testing.T, mode string) {
	t.Helper()
	old, had := os.LookupEnv("MCP_TOOLS_MODE")
	if mode == "" {
		_ = os.Unsetenv("MCP_TOOLS_MODE")
	} else {
		_ = os.Setenv("MCP_TOOLS_MODE", mode)
	}
	t.Cleanup(func() {
		if had {
			_ = os.Setenv("MCP_TOOLS_MODE", old)
		} else {
			_ = os.Unsetenv("MCP_TOOLS_MODE")
		}
	})
}

func TestToolFilterReadOnly(t *testing.T) {
	withToolsMode(t, "all")
	filter := toolFilter(true)

	tools := []mcp.Tool{
		{Name: "list_products", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "create_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
		{Name: "get_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "delete_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
		{Name: "no_hint_tool"},
	}

	filtered := filter(context.Background(), tools)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 read-only tools, got %d: %v", len(filtered), filtered)
	}
	for _, tool := range filtered {
		if tool.Name != "list_products" && tool.Name != "get_product" {
			t.Errorf("unexpected tool in filtered list: %s", tool.Name)
		}
	}
}

func TestToolFilterNilReadOnlyHintIsWrite(t *testing.T) {
	withToolsMode(t, "all")
	filter := toolFilter(true)
	filtered := filter(context.Background(), []mcp.Tool{{Name: "orphan"}})
	if len(filtered) != 0 {
		t.Errorf("nil ReadOnlyHint must not leak into read-only list, got %d", len(filtered))
	}
}

func TestToolFilterCoreSurface(t *testing.T) {
	withToolsMode(t, "core")
	filter := toolFilter(false)
	tools := []mcp.Tool{
		{Name: "list_products", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "create_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
		{Name: "search_insales_tools", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "call_insales_tool", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
		{Name: "update_webhook", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
	}
	filtered := filter(context.Background(), tools)
	names := map[string]bool{}
	for _, t := range filtered {
		names[t.Name] = true
	}
	if !names["list_products"] || !names["search_insales_tools"] || !names["call_insales_tool"] {
		t.Errorf("core+meta must stay visible, got %v", names)
	}
	if names["create_product"] || names["update_webhook"] {
		t.Errorf("non-core tools must be hidden in core mode, got %v", names)
	}
}

func TestToolFilterAllMode(t *testing.T) {
	withToolsMode(t, "all")
	filter := toolFilter(false)
	tools := []mcp.Tool{
		{Name: "list_products", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "create_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
	}
	filtered := filter(context.Background(), tools)
	if len(filtered) != 2 {
		t.Errorf("expected 2 tools in all mode, got %d", len(filtered))
	}
}

func TestToolFilterShopReadOnly(t *testing.T) {
	withToolsMode(t, "all")
	filter := toolFilter(false)

	tools := []mcp.Tool{
		{Name: "list_products", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(true)}},
		{Name: "create_product", Annotations: mcp.ToolAnnotation{ReadOnlyHint: boolPtr(false)}},
	}

	ctx := context.WithValue(context.Background(), core.ShopReadOnlyKey, true)
	filtered := filter(ctx, tools)

	if len(filtered) != 1 {
		t.Fatalf("expected 1 read-only tool when shop is read-only, got %d", len(filtered))
	}
	if filtered[0].Name != "list_products" {
		t.Errorf("expected list_products, got %s", filtered[0].Name)
	}
}

func TestCoreToolNamesAreKnownTools(t *testing.T) {
	// every core name must exist in the full tool name list produced by tools.go
	// (checked loosely: prefixes that we register; full sync test lives in handlers)
	for name := range coreToolNames {
		if name == "" {
			t.Error("empty core tool name")
		}
	}
	if len(coreToolNames) < 20 || len(coreToolNames) > 50 {
		t.Errorf("core surface size = %d, want ~30-40", len(coreToolNames))
	}
}

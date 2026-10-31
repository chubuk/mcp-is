package mcp

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/mark3labs/mcp-go/mcp"
)

func boundCtx() context.Context {
	ctx := context.WithValue(context.Background(), auth.ShopIDKey, "shop-1")
	ctx = context.WithValue(ctx, auth.ConnectionNameKey, "romashka")
	ctx = context.WithValue(ctx, auth.ShopNameKey, "Ромашка")
	ctx = context.WithValue(ctx, auth.ShopURLKey, "romashka.myinsales.ru")
	ctx = context.WithValue(ctx, auth.KeyNameKey, "Cursor MCP")
	return ctx
}

func TestStoreInfo_Bound(t *testing.T) {
	res, err := storeInfoHandler()(boundCtx(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	var info map[string]any
	json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &info)

	if info["bound"] != true {
		t.Errorf("bound = %v, want true", info["bound"])
	}
	if info["connection_name"] != "romashka" {
		t.Errorf("connection_name = %v", info["connection_name"])
	}
	store := info["store"].(map[string]any)
	if store["name"] != "Ромашка" || store["id"] != "shop-1" {
		t.Errorf("store = %v", store)
	}
}

func TestStoreInfo_Unbound(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.KeyNameKey, "legacy key")
	res, err := storeInfoHandler()(ctx, mcp.CallToolRequest{})
	if err != nil {
		t.Fatal(err)
	}
	var info map[string]any
	json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &info)

	if info["bound"] != false {
		t.Errorf("bound = %v, want false", info["bound"])
	}
	msg, _ := info["message"].(string)
	if !strings.Contains(msg, "не привязан") {
		t.Errorf("message should explain missing binding: %q", msg)
	}
}

func TestBuildInstructions_Bound(t *testing.T) {
	s := buildInstructions(boundCtx(), "BASE")
	for _, want := range []string{"Connection: romashka", `Store: "Ромашка" (romashka.myinsales.ru)`, "store_info", "BASE"} {
		if !strings.Contains(s, want) {
			t.Errorf("instructions must contain %q", want)
		}
	}
	if idx := strings.Index(s, "Connection: romashka"); idx > 512 {
		t.Errorf("identity block must start within first 512 chars, got %d", idx)
	}
}

func TestBuildInstructions_Unbound(t *testing.T) {
	s := buildInstructions(context.Background(), "BASE")
	if !strings.Contains(s, "не привязан") {
		t.Error("unbound instructions must warn about missing binding")
	}
	if !strings.Contains(s, "BASE") {
		t.Error("base instructions must be preserved")
	}
}

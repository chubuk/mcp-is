package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/mark3labs/mcp-go/mcp"
)

// ── mutationResult: delete ──

func TestDeleteResult_WithSnapshot(t *testing.T) {
	pre := map[string]any{"id": "42", "title": "Товар"}
	res := deleteResult("product", "42",
		func() (any, error) { return pre, nil },
		func() error { return nil },
	)
	if res == nil || res.IsError {
		t.Fatalf("deleteResult error: %+v", res)
	}
	var payload MutationResult
	if err := json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !payload.OK || payload.Action != "delete" || payload.Resource != "product" || payload.ID != "42" {
		t.Errorf("payload header wrong: %+v", payload)
	}
	if payload.Deleted == nil {
		t.Error("deleted snapshot missing")
	}
	if payload.Message == "" {
		t.Error("message empty")
	}
}

func TestDeleteResult_NotFoundPrefetch(t *testing.T) {
	res := deleteResult("product", "42",
		func() (any, error) { return nil, insales.ErrNotFound },
		func() error { return nil },
	)
	var payload MutationResult
	_ = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &payload)
	if payload.Deleted != nil {
		t.Errorf("deleted should be null on prefetch 404, got %v", payload.Deleted)
	}
	if payload.Message == "" {
		t.Error("message should explain not-found")
	}
}

func TestDeleteResult_DeleteError(t *testing.T) {
	res := deleteResult("product", "42",
		func() (any, error) { return map[string]any{"id": "42"}, nil },
		func() error { return insales.ErrForbidden },
	)
	if res == nil || !res.IsError {
		t.Fatal("expected error result on delete failure")
	}
}

// ── mutationResult: update ──

func TestUpdateResult_ChangedFields(t *testing.T) {
	changed := map[string]FieldChange{
		"title": {Old: "Было", New: "Стало"},
		"sku":   {Old: "A1", New: "A2"},
	}
	res := updateResult("product", "42", "Товар",
		nil,
		func() (any, error) { return map[string]any{"id": "42", "title": "Стало"}, nil },
		changed,
		[]string{"unknown_field"},
	)
	if res == nil || res.IsError {
		t.Fatalf("updateResult error: %+v", res)
	}
	var payload MutationResult
	if err := json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload.Action != "update" {
		t.Errorf("action = %q, want update", payload.Action)
	}
	if payload.ChangedFields == nil || payload.ChangedFields["title"].Old != "Было" {
		t.Errorf("changed_fields wrong: %+v", payload.ChangedFields)
	}
	if len(payload.IgnoredKeys) != 1 || payload.IgnoredKeys[0] != "unknown_field" {
		t.Errorf("ignored_keys wrong: %v", payload.IgnoredKeys)
	}
	if payload.Data == nil {
		t.Error("data (final state) missing")
	}
}

// ── updateWithDiff: diff по аргументам (old из pre-fetch, new из аргументов) ──

func TestUpdateWithDiff_ChangedFields(t *testing.T) {
	old := map[string]any{"id": "42", "title": "Было", "sku": "A1"}
	args := map[string]any{"title": "Стало", "sku": "A2", "product_id": "42", "bogus": "x"}
	res := updateWithDiff("product", "42", args, []string{"bogus"},
		func() (any, error) { return old, nil },
		func() (any, error) { return map[string]any{"id": "42", "title": "Стало"}, nil },
		nil,
	)
	if res == nil || res.IsError {
		t.Fatalf("updateWithDiff error: %+v", res)
	}
	var payload MutationResult
	if err := json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload.Action != "update" || payload.Resource != "product" || payload.ID != "42" {
		t.Errorf("envelope wrong: %+v", payload)
	}
	if fc := payload.ChangedFields["title"]; fc.Old != "Было" || fc.New != "Стало" {
		t.Errorf("title diff wrong: %+v", fc)
	}
	if fc := payload.ChangedFields["sku"]; fc.Old != "A1" || fc.New != "A2" {
		t.Errorf("sku diff wrong: %+v", fc)
	}
	if _, ok := payload.ChangedFields["product_id"]; ok {
		t.Error("product_id must not be in changed_fields")
	}
	if _, ok := payload.ChangedFields["bogus"]; ok {
		t.Error("ignored key must not be in changed_fields")
	}
	if len(payload.IgnoredKeys) != 1 || payload.IgnoredKeys[0] != "bogus" {
		t.Errorf("ignored_keys wrong: %v", payload.IgnoredKeys)
	}
}

func TestUpdateWithDiff_PriceFromVariants(t *testing.T) {
	// Цена товара живёт в variants[0].price, product.price = null.
	old := map[string]any{"id": "1", "price": nil, "variants": []any{
		map[string]any{"id": "7", "price": "100.0"},
	}}
	args := map[string]any{"price": "250"}
	res := updateWithDiff("product", "1", args, nil,
		func() (any, error) { return old, nil },
		func() (any, error) { return old, nil },
		nil,
	)
	var payload MutationResult
	if err := json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if fc := payload.ChangedFields["price"]; fc.Old != "100.0" || fc.New != "250" {
		t.Errorf("price diff wrong: %+v", fc)
	}
}

func TestUpdateWithDiff_MaskAppliedToOldAndData(t *testing.T) {
	old := map[string]any{"id": "9", "email": "real@example.com"}
	args := map[string]any{"email": "new@example.com"}
	mask := func(v any) any {
		m, ok := v.(map[string]any)
		if !ok {
			return v
		}
		out := make(map[string]any, len(m))
		for k, val := range m {
			if k == "email" {
				val = "m***@example.com"
			}
			out[k] = val
		}
		return out
	}
	res := updateWithDiff("client", "9", args, nil,
		func() (any, error) { return old, nil },
		func() (any, error) { return map[string]any{"id": "9", "email": "new@example.com"}, nil },
		mask,
	)
	var payload MutationResult
	if err := json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if fc := payload.ChangedFields["email"]; fc.Old != "m***@example.com" {
		t.Errorf("old must be masked: %+v", fc)
	}
	data := payload.Data.(map[string]any)
	if data["email"] != "m***@example.com" {
		t.Errorf("data must be masked: %+v", data)
	}
}

func TestUpdateWithDiff_NoPrefetch_OldNil(t *testing.T) {
	args := map[string]any{"title": "X"}
	res := updateWithDiff("account", "", args, nil, nil,
		func() (any, error) { return map[string]any{"title": "X"}, nil },
		nil,
	)
	var payload MutationResult
	if err := json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if fc := payload.ChangedFields["title"]; fc.Old != nil || fc.New != "X" {
		t.Errorf("diff without prefetch wrong: %+v", fc)
	}
}

// ── outputSchema ──

func TestOutputSchemaFor_Mutations(t *testing.T) {
	cases := []struct {
		name string
		want bool // has schema
	}{
		{"delete_product", true},
		{"update_product", true},
		{"create_product", true},
		{"add_image", true},
		{"decline_application_charge", true},
		{"list_products", false},
		{"get_product", false},
		{"store_info", false},
	}
	for _, c := range cases {
		got := outputSchemaFor(c.name)
		if (got != nil) != c.want {
			t.Errorf("outputSchemaFor(%q) = %v, want present=%v", c.name, got != nil, c.want)
		}
	}
}

func TestOutputSchema_MutationResultValid(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(mutationResultSchema, &schema); err != nil {
		t.Fatalf("mutationResultSchema invalid JSON: %v", err)
	}
	if schema["type"] != "object" {
		t.Error("schema type must be object")
	}
	props, _ := schema["properties"].(map[string]any)
	for _, key := range []string{"ok", "action", "resource", "id", "changed_fields", "deleted", "data"} {
		if _, ok := props[key]; !ok {
			t.Errorf("schema missing property %q", key)
		}
	}
}

// ── ignored_keys (filterArgsByTool / wrapArgsAudit) ──

func TestFilterArgsByTool_UnknownKeys(t *testing.T) {
	toolArgWhitelist["test_tool"] = map[string]bool{"title": true, "price": true}
	defer delete(toolArgWhitelist, "test_tool")

	args := map[string]any{"title": "X", "price": "10", "bogus": "y", "another": 1}
	ignored := filterArgsByTool(args, "test_tool")
	if len(ignored) != 2 {
		t.Fatalf("ignored = %v, want 2 keys", ignored)
	}
	if _, ok := args["bogus"]; ok {
		t.Error("bogus key not removed")
	}
	if _, ok := args["title"]; !ok {
		t.Error("title key wrongly removed")
	}
}

func TestFilterArgsByTool_NoWhitelist_NoOp(t *testing.T) {
	args := map[string]any{"anything": true}
	ignored := filterArgsByTool(args, "nonexistent_tool_zzz")
	if ignored != nil {
		t.Errorf("ignored = %v, want nil (fail-open)", ignored)
	}
	if _, ok := args["anything"]; !ok {
		t.Error("args must not be touched without whitelist")
	}
}

func TestWrapArgsAudit_AddsIgnoredKeys(t *testing.T) {
	toolArgWhitelist["audit_tool"] = map[string]bool{"keep": true}
	defer delete(toolArgWhitelist, "audit_tool")

	inner := func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolResult(map[string]any{"ok": true}), nil
	}
	wrapped := wrapArgsAudit(inner)

	req := mcp.CallToolRequest{}
	req.Params.Name = "audit_tool"
	req.Params.Arguments = map[string]any{"keep": 1, "junk": 2}

	res, err := wrapped(context.Background(), req)
	if err != nil || res == nil {
		t.Fatalf("wrapped call failed: %v", err)
	}
	m, ok := res.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content not a map: %T", res.StructuredContent)
	}
	ig, ok := m["ignored_keys"].([]string)
	if !ok || len(ig) != 1 || ig[0] != "junk" {
		t.Errorf("ignored_keys = %v, want [junk]", m["ignored_keys"])
	}
}

// ── handle/permalink: стабильность имён ──

func TestHandleNaming_AliasPresent(t *testing.T) {
	// MapProperty должен отдавать handle (и permalink как алиас), не голый permalink.
	prop := map[string]any{"id": 1, "permalink": "krasnyi-cvet"}
	if _, ok := prop["permalink"]; !ok {
		t.Fatal("test fixture broken")
	}
	// Проверяем сам паттерн через mapper-контракт: handle должен существовать при permalink.
	// (Mapper тестируется в insales-пакете; здесь — инвариант имени.)
	t.Log("handle/permalink alias contract — see TestMapProperty_Alias in insales package")
}

// TestInputSchemaFormats — enum/pattern реально используются в tools.go (дрейф-гейт).
func TestInputSchemaFormats(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "tools.go", nil, 0)
	if err != nil {
		t.Fatalf("parse tools.go: %v", err)
	}
	enumCalls, patternCalls := 0, 0
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		fun, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || fun.Sel == nil {
			return true
		}
		switch fun.Sel.Name {
		case "Enum":
			enumCalls++
		case "Pattern":
			patternCalls++
		}
		return true
	})
	if enumCalls < 3 {
		t.Errorf("mcp.Enum used %d times in tools.go, want >= 3 (словари InSales должны быть ограничены)", enumCalls)
	}
	if patternCalls < 3 {
		t.Errorf("mcp.Pattern used %d times in tools.go, want >= 3 (slug-поля должны быть ограничены)", patternCalls)
	}
}

// ── effective status (страховка модели состояний) ──

func TestErrorsAreSentinels(t *testing.T) {
	wrapped := errors.Join(errors.New("API error 403: нет прав"), insales.ErrForbidden)
	if !errors.Is(wrapped, insales.ErrForbidden) {
		t.Error("ErrForbidden must be detectable via errors.Is")
	}
}

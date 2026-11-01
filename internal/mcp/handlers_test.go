package mcp

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, "null"},
		{"string", "hello", `"hello"`},
		{"number", 42, "42"},
		{"bool", true, "true"},
		{"map", map[string]any{"key": "value"}, "{\n  \"key\": \"value\"\n}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatResult(tt.input)
			if result != tt.expected {
				t.Errorf("formatResult(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// structuredContent по MCP-спеке обязан сериализоваться в JSON-объект.
func TestWrapEnvelope(t *testing.T) {
	tests := []struct {
		name string
		in   any
	}{
		{"slice", []any{map[string]any{"id": "1"}}},
		{"dynamic slice in any", any([]any{map[string]any{"id": "1"}})},
		{"typed nil slice", []map[string]any(nil)},
		{"string", "hello"},
		{"int", 42},
		{"byte slice", []byte("raw")},
		{"raw message array", json.RawMessage(`[1,2]`)},
		{"raw message object", json.RawMessage(`{"a":1}`)},
		{"map", map[string]any{"a": 1}},
		{"ptr to slice", &[]int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := wrapEnvelope(tt.in)
			b, err := json.Marshal(out)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if b[0] != '{' {
				t.Errorf("marshals to %s, want JSON object", b)
			}
		})
	}
	// nil остаётся nil — поле omitempty, спека разрешает отсутствие.
	if got := wrapEnvelope(nil); got != nil {
		t.Errorf("wrapEnvelope(nil) = %v, want nil", got)
	}
	// слайс → домашний конверт {items,total}
	out := wrapEnvelope([]any{"a", "b"}).(map[string]any)
	if out["total"] != 2 {
		t.Errorf("total = %v, want 2", out["total"])
	}
	if items, ok := out["items"].([]any); !ok || len(items) != 2 {
		t.Errorf("items = %v, want 2-element slice", out["items"])
	}
	// объект не модифицируется
	obj := map[string]any{"a": 1}
	if wrapEnvelope(obj).(map[string]any)["a"] != 1 {
		t.Error("map payload modified")
	}
}

func TestToolResult_ObjectStructuredContent(t *testing.T) {
	res := toolResult([]any{map[string]any{"id": "1"}})
	if res.StructuredContent == nil {
		t.Fatal("structuredContent missing")
	}
	b, err := json.Marshal(res.StructuredContent)
	if err != nil || b[0] != '{' {
		t.Errorf("structuredContent = %s (%v), want JSON object", b, err)
	}
}

func TestGetFloat(t *testing.T) {
	args := map[string]any{
		"page":     float64(5),
		"per_page": float64(50),
	}

	if got := getFloat(args, "page", 1); got != 5 {
		t.Errorf("getFloat(page) = %v, want 5", got)
	}
	if got := getFloat(args, "missing", 10); got != 10 {
		t.Errorf("getFloat(missing) = %v, want 10", got)
	}
	if got := getFloat(nil, "key", 7); got != 7 {
		t.Errorf("getFloat(nil) = %v, want 7", got)
	}
}

func TestHashArgsStable(t *testing.T) {
	a := hashArgs(map[string]any{"b": 2, "a": 1})
	b := hashArgs(map[string]any{"a": 1, "b": 2})
	if a == "" || a != b {
		t.Errorf("hash not stable: %q vs %q", a, b)
	}
	if hashArgs(nil) != "" {
		t.Error("empty args should hash empty")
	}
}

func TestAuditFromContext(t *testing.T) {
	e := auditFromContext(context.Background(), "create_product", "call", map[string]any{"title": "x"}, true)
	if e.ToolName != "create_product" || e.Action != "call" || !e.ResultOK {
		t.Errorf("entry = %+v", e)
	}
	if e.ArgsHash == "" {
		t.Error("expected args hash")
	}
}

func TestServerJSONVersionMatches(t *testing.T) {
	raw, err := os.ReadFile("../../server.json")
	if err != nil {
		t.Skip("server.json not found")
	}
	var doc struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	// keep in sync with server.go NewMCPServer version literal
	src, err := os.ReadFile("server.go")
	if err != nil {
		t.Fatal(err)
	}
	want := regexp.MustCompile(`"(\d+\.\d+\.\d+)"`).FindStringSubmatch(string(src))
	if want == nil {
		t.Fatal("version literal not found in server.go")
	}
	if doc.Version != want[1] {
		t.Errorf("server.json version %q != server.go %q", doc.Version, want[1])
	}
}

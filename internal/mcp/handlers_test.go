package mcp

import (
	"context"
	"encoding/json"
	"os"
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
	if doc.Version != "0.6.0" {
		t.Errorf("server.json version %q != server.go 0.6.0", doc.Version)
	}
}

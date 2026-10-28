package mcp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chubuk/mcp-is/internal/core"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestIsHTMLResponse_Doctype(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"<!DOCTYPE html><html><body>Error</body></html>", true},
		{"<!DOCTYPE HTML><html><body>Error</body></HTML>", true},
		{"  <!DOCTYPE html><html>", true},
		{"API error 404: <!DOCTYPE html><html>", true},
		{"API error 500: <!DOCTYPE HTML><html>", true},
		{`{"error":"not found"}`, false},
		{"plain text error", false},
		{"", false},
	}
	for _, tt := range tests {
		got := isHTMLResponse(tt.input)
		if got != tt.want {
			t.Errorf("isHTMLResponse(%q) = %v, want %v", tt.input[:min(50, len(tt.input))], got, tt.want)
		}
	}
}

func TestCleanError_HTML(t *testing.T) {
	err := fmt.Errorf("API error 404: <!DOCTYPE html><html><body>Not Found</body></html>")
	cleaned := cleanError(err)
	if cleaned == nil {
		t.Fatal("expected non-nil error")
	}
	msg := cleaned.Error()
	if msg != "InSales API error: API error 404:" {
		t.Errorf("expected clean error message, got: %q", msg)
	}
}

func TestCleanError_NonHTML(t *testing.T) {
	err := fmt.Errorf("connection refused")
	cleaned := cleanError(err)
	if cleaned == nil {
		t.Fatal("expected non-nil error")
	}
	if cleaned.Error() != "connection refused" {
		t.Errorf("expected original message, got: %q", cleaned.Error())
	}
}

func TestCleanError_Nil(t *testing.T) {
	cleaned := cleanError(nil)
	if cleaned != nil {
		t.Errorf("expected nil, got: %v", cleaned)
	}
}

func TestToolError_HTML(t *testing.T) {
	err := fmt.Errorf("API error 404: <!DOCTYPE html><html><body>Error</body></html>")
	result := toolError(err)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	if textContent.Text != "InSales API error: API error 404:" {
		t.Errorf("expected clean error, got: %q", textContent.Text)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestToolError_BillingSentinels(t *testing.T) {
	wrapped := fmt.Errorf("shop s1: %w", core.ErrBillingUnpaid)
	res := toolError(wrapped)
	if !res.IsError {
		t.Fatal("expected IsError result")
	}
	msg := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(msg, "не оплачена") {
		t.Errorf("billing error text = %q, want explanation about unpaid period", msg)
	}
	res = toolError(fmt.Errorf("shop s1: %w", core.ErrBillingExpired))
	msg = res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(msg, "закончился") {
		t.Errorf("expired error text = %q", msg)
	}
}

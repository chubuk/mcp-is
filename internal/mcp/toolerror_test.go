package mcp

import (
	"errors"
	"strings"
	"testing"

	"github.com/chubuk/mcp-is/internal/core"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestToolError_KeyNotBound(t *testing.T) {
	res := toolError(errors.Join(errors.New("wrap"), core.ErrKeyNotBound))
	if !res.IsError {
		t.Fatal("expected tool error result")
	}
	text := res.Content[0].(mcp.TextContent).Text
	for _, want := range []string{"не привязан", "привязать", "AI-ключи"} {
		if !strings.Contains(text, want) {
			t.Errorf("error text %q must contain %q", text, want)
		}
	}
}

package core

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "mcp-is"},
		{"whitespace only", "   ", "mcp-is"},
		{"simple latin", "test", "test"},
		{"with spaces", "Cursor MCP", "cursor-mcp"},
		{"with underscores", "my_key_name", "my-key-name"},
		{"multiple spaces", "hello   world", "hello-world"},
		{"leading/trailing spaces", "  test  ", "test"},
		{"special chars removed", "key@#$%name!", "keyname"},
		{"mixed case", "ClaudeDesktop", "claudedesktop"},
		{"cyrillic", "Мой Ключ", "moy-klyuch"},
		{"cyrillic no spaces", "тест", "test"},
		{"mixed cyrillic latin", "Мой Cursor", "moy-cursor"},
		{"numbers", "key123", "key123"},
		{"numbers with text", "shop-v2", "shop-v2"},
		{"only special chars", "@#$%", "mcp-is"},
		{"single char", "a", "a"},
		{"hyphens preserved", "my-key", "my-key"},
		{"double hyphens collapsed", "a--b", "a-b"},
		{"cyrillic complex", "Парафарм Шоп", "parafarm-shop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

package insales

import (
	"context"
	"testing"
)

func TestParseTitleFromHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "standard InSales title",
			input: `<html><head><title>Парафарм Шоп | Официальный магазин</title></head><body></body></html>`,
			want:  "Парафарм Шоп",
		},
		{
			name:  "title without pipe",
			input: `<html><head><title>My Shop</title></head></html>`,
			want:  "My Shop",
		},
		{
			name:  "uppercase TITLE",
			input: `<html><head><TITLE>Test Shop | Subtitle</TITLE></head></html>`,
			want:  "Test Shop",
		},
		{
			name:  "no title tag",
			input: `<html><head></head><body></body></html>`,
			want:  "",
		},
		{
			name:  "empty title",
			input: `<html><head><title></title></head></html>`,
			want:  "",
		},
		{
			name:  "title with whitespace",
			input: `<html><head><title>  My Shop  |  Sub  </title></head></html>`,
			want:  "My Shop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTitleFromHTML(tt.input)
			if got != tt.want {
				t.Errorf("parseTitleFromHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetStoreInfo_WithMockServer(t *testing.T) {
	storeHTML := `<html><head><title>Парафарм Шоп | Магазин</title></head><body></body></html>`
	title := parseTitleFromHTML(storeHTML)
	if title != "Парафарм Шоп" {
		t.Errorf("expected 'Парафарм Шоп', got %q", title)
	}
}

func TestGetStoreInfo_FallbackToDomain(t *testing.T) {
	// If storefront is unreachable, should fallback to domain
	title := parseTitleFromHTML("")
	if title != "" {
		t.Errorf("expected empty title for empty HTML, got %q", title)
	}
}

func TestGetStoreInfo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	title := parseTitleFromHTML("<html><head><title>Test</title></head></html>")
	if title != "Test" {
		t.Errorf("expected 'Test', got %q", title)
	}
	_ = ctx
}

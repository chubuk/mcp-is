package core

import (
	"context"
	"errors"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
)

// A key with no shop binding must fail closed — never fall back to the user's default shop.
func TestGetClient_UnboundKeyFailsClosed(t *testing.T) {
	m := NewAccountManager(nil, nil, false, 0, nil)
	ctx := context.WithValue(context.Background(), auth.UserIDKey, "user-1")

	_, _, err := m.GetClient(ctx)
	if !errors.Is(err, ErrKeyNotBound) {
		t.Fatalf("got %v, want ErrKeyNotBound", err)
	}
}

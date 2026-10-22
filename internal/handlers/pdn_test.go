package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupPDNHandler(t *testing.T) *PDNHandler {
	t.Helper()
	dbURL := "postgres://admin:password@localhost:5432/mcp_is?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skip("PostgreSQL not available, skipping PDN handler tests")
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skip("PostgreSQL not reachable, skipping PDN handler tests")
	}
	t.Cleanup(func() { pool.Close() })
	return NewPDNHandler(pool, storage.NewSettingsRepository(pool))
}

// Regression: UpdateSettings used to only invalidate cache and never persist — checkbox bounced back.
func TestUpdateSettings_PDNPersists(t *testing.T) {
	h := setupPDNHandler(t)

	enabled, log, mask, auto := true, false, true, false
	body, _ := json.Marshal(PDNSettingsRequest{
		Enabled:     &enabled,
		Log:         &log,
		MaskPII:     &mask,
		AutoProtect: &auto,
	})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/pdn", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), auth.RoleKey, "super_admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	h.UpdateSettings(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var got PDNSettingsResponse
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	want := PDNSettingsResponse{Enabled: true, Log: false, MaskPII: true, AutoProtect: false}
	if got != want {
		t.Fatalf("after save got %+v, want %+v", got, want)
	}

	// fresh GetSettings must see the same values (DB round-trip)
	req2 := httptest.NewRequest(http.MethodGet, "/api/admin/pdn", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), auth.RoleKey, "super_admin"))
	w2 := httptest.NewRecorder()
	h.GetSettings(w2, req2)
	var again PDNSettingsResponse
	json.NewDecoder(w2.Body).Decode(&again)
	if again != want {
		t.Fatalf("after reload got %+v, want %+v", again, want)
	}
}

func TestUpdateSettings_PDNAdminDenied(t *testing.T) {
	h := setupPDNHandler(t)
	body, _ := json.Marshal(PDNSettingsRequest{})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/pdn", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), auth.RoleKey, "admin"))
	w := httptest.NewRecorder()
	h.UpdateSettings(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

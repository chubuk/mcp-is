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

func setupServerHandler(t *testing.T) (*ServerHandler, *pgxpool.Pool) {
	t.Helper()
	dbURL := "postgres://admin:password@localhost:5432/mcp_is?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skip("PostgreSQL not available, skipping server handler tests")
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skip("PostgreSQL not reachable, skipping server handler tests")
	}
	t.Cleanup(func() { pool.Close() })

	settingsRepo := storage.NewSettingsRepository(pool)
	return NewServerHandler(settingsRepo, 0), pool
}

func TestGetSettings_AccessDenied(t *testing.T) {
	h, _ := setupServerHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/server", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetSettings(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestGetSettings_Admin(t *testing.T) {
	h, _ := setupServerHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/server", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetSettings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp ServerSettingsResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Domain == "" {
		t.Error("expected domain in response")
	}
}

func TestGetSettings_SuperAdmin(t *testing.T) {
	h, _ := setupServerHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/server", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "super_admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetSettings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateSettings_AccessDenied(t *testing.T) {
	h, _ := setupServerHandler(t)

	body, _ := json.Marshal(ServerSettingsRequest{
		Domain: strPtr("new-domain.com"),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/server", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateSettings(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for admin, got %d", w.Code)
	}
}

func TestUpdateSettings_SuperAdmin(t *testing.T) {
	h, _ := setupServerHandler(t)

	body, _ := json.Marshal(ServerSettingsRequest{
		Domain: strPtr("test-domain.com"),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/server", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), auth.RoleKey, "super_admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateSettings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIssueCert_AccessDenied(t *testing.T) {
	h, _ := setupServerHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/server/cert-issue", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.IssueCert(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestRenewCert_AccessDenied(t *testing.T) {
	h, _ := setupServerHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/server/cert-renew", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.RenewCert(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func strPtr(s string) *string {
	return &s
}

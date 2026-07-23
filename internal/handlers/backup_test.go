package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupBackupHandler(t *testing.T) (*BackupHandler, *pgxpool.Pool) {
	t.Helper()
	dbURL := "postgres://admin:password@localhost:5432/mcp_is?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skip("PostgreSQL not available, skipping backup handler tests")
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skip("PostgreSQL not reachable, skipping backup handler tests")
	}
	t.Cleanup(func() { pool.Close() })

	return NewBackupHandler(pool), pool
}

func TestBackupGetSettings_AccessDenied(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetSettings(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestBackupGetSettings_Admin(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetSettings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp BackupSettingsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Method == "" {
		t.Error("expected method to be set")
	}
}

func TestBackupGetSettings_SuperAdmin(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "super_admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetSettings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBackupUpdateSettings_AdminDenied(t *testing.T) {
	h, _ := setupBackupHandler(t)

	body := `{"retention_days": 7}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/backup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateSettings(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for admin, got %d", w.Code)
	}
}

func TestBackupUpdateSettings_SuperAdmin(t *testing.T) {
	h, _ := setupBackupHandler(t)

	body := `{"retention_days": 7}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/backup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.RoleKey, "super_admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateSettings(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBackupGetLog_AccessDenied(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/log", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetLog(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestBackupGetLog_Admin(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/log", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetLog(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBackupGetIncidents_AccessDenied(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/incidents", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetIncidents(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestBackupGetIncidents_Admin(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/incidents", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetIncidents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBackupGetStatus_AccessDenied(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/status", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetStatus(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestBackupGetStatus_Admin(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/status", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp BackupStatusResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Load5m < 0 {
		t.Error("expected load_5m to be non-negative")
	}
}

func TestBackupRunBackup_UserDenied(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/backup/run", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.RunBackup(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestBackupRunBackup_AdminAllowed(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/backup/run", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// This will fail because backup.sh is not available in test environment
	// but it should return 500, not 403
	h.RunBackup(w, req)

	if w.Code == http.StatusForbidden {
		t.Error("admin should be allowed to run backup")
	}
}

func TestBackupGetSettings_InvalidJSON(t *testing.T) {
	h, _ := setupBackupHandler(t)

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/backup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.RoleKey, "super_admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateSettings(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", w.Code)
	}
}

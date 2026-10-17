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
	h, pool := setupBackupHandler(t)

	// seed one row so pagination fields are meaningful
	_, _ = pool.Exec(context.Background(),
		`INSERT INTO backup_log (status, size_bytes) VALUES ('ok', 1)`)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/log?page=1&per_page=2", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetLog(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Entries []BackupLogEntry `json:"entries"`
		Total   int              `json:"total"`
		Page    int              `json:"page"`
		PerPage int              `json:"per_page"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Page != 1 || resp.PerPage != 2 {
		t.Errorf("expected page=1 per_page=2, got %d/%d", resp.Page, resp.PerPage)
	}
	if resp.Total < 1 {
		t.Errorf("expected total >= 1, got %d", resp.Total)
	}
	if len(resp.Entries) > 2 {
		t.Errorf("expected at most 2 entries, got %d", len(resp.Entries))
	}
}

func TestBackupGetLog_InvalidPagingDefaults(t *testing.T) {
	h, _ := setupBackupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/log?page=0&per_page=9999", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetLog(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Page != 1 || resp.PerPage != 20 {
		t.Errorf("expected defaults 1/20, got %d/%d", resp.Page, resp.PerPage)
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
	h, pool := setupBackupHandler(t)

	_, _ = pool.Exec(context.Background(),
		`INSERT INTO backup_incidents (type, message, severity) VALUES ('test', 'msg', 'info')`)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup/incidents?page=1&per_page=2", nil)
	ctx := context.WithValue(req.Context(), auth.RoleKey, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetIncidents(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Incidents []BackupIncident `json:"incidents"`
		Total     int              `json:"total"`
		Page      int              `json:"page"`
		PerPage   int              `json:"per_page"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Page != 1 || resp.PerPage != 2 {
		t.Errorf("expected page=1 per_page=2, got %d/%d", resp.Page, resp.PerPage)
	}
	if resp.Total < 1 {
		t.Errorf("expected total >= 1, got %d", resp.Total)
	}
	if len(resp.Incidents) > 2 {
		t.Errorf("expected at most 2 incidents, got %d", len(resp.Incidents))
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

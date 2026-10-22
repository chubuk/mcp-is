package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/email"
)

func ctxWithRole(role string) context.Context {
	return context.WithValue(context.Background(), auth.RoleKey, role)
}

func TestValidateTemplate_NoPanic(t *testing.T) {
	err := email.ValidateTemplate("S", "{{.Name}", "", []string{"Name"}, map[string]any{"Name": "x"})
	if err == nil {
		t.Fatal("expected parse error without panic")
	}
}

func TestLoadSMTPConfigFromSettings_Empty(t *testing.T) {
	env := email.Config{Host: "env.smtp", Port: "587", From: "a@b.c"}
	cfg := LoadSMTPConfigFromSettings(nil, env, nil)
	if cfg.Host != "env.smtp" {
		t.Errorf("host = %q", cfg.Host)
	}
}

func TestGetSMTPSettings_AdminAllowed(t *testing.T) {
	aes, _ := auth.NewAESCipher("test-encryption-key-32-bytes!!!!")
	svc := email.NewService(email.Config{Host: "smtp.x", Port: "587", User: "u", Pass: "secret", From: "a@b.c", Domain: "example.com"})
	h := NewEmailAdminHandler(nil, nil, svc, aes)

	req := httptest.NewRequest("GET", "/api/admin/smtp/settings", nil).WithContext(ctxWithRole("admin"))
	rec := httptest.NewRecorder()
	h.GetSMTPSettings(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp smtpSettingsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.HasPass {
		t.Error("has_password should be true")
	}
	if bytes.Contains(rec.Body.Bytes(), []byte("secret")) {
		t.Error("password leaked in GET response")
	}
}

func TestGetSMTPSettings_UserDenied(t *testing.T) {
	svc := email.NewService(email.Config{Domain: "example.com"})
	h := NewEmailAdminHandler(nil, nil, svc, nil)
	req := httptest.NewRequest("GET", "/api/admin/smtp/settings", nil).WithContext(ctxWithRole("user"))
	rec := httptest.NewRecorder()
	h.GetSMTPSettings(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}

func TestUpdateSMTPSettings_KeepsPassword(t *testing.T) {
	aes, _ := auth.NewAESCipher("test-encryption-key-32-bytes!!!!")
	svc := email.NewService(email.Config{Host: "smtp.x", Port: "587", Pass: "oldpass", From: "a@b.c", Domain: "example.com"})
	h := NewEmailAdminHandler(nil, nil, svc, aes)

	body, _ := json.Marshal(map[string]any{"host": "smtp.y", "port": "465", "tls_mode": "ssl"})
	req := httptest.NewRequest("PUT", "/api/admin/smtp/settings", bytes.NewReader(body)).WithContext(ctxWithRole("super_admin"))
	rec := httptest.NewRecorder()
	h.UpdateSMTPSettings(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	cfg := svc.Config()
	if cfg.Host != "smtp.y" || cfg.Pass != "oldpass" {
		t.Errorf("cfg = %+v", cfg)
	}
	if cfg.TLSMode != email.TLSModeSSL {
		t.Errorf("tls = %q", cfg.TLSMode)
	}
}

func TestUpdateSMTPSettings_EmptyPassClears(t *testing.T) {
	aes, _ := auth.NewAESCipher("test-encryption-key-32-bytes!!!!")
	svc := email.NewService(email.Config{Host: "smtp.x", Port: "587", Pass: "oldpass", From: "a@b.c", Domain: "example.com"})
	h := NewEmailAdminHandler(nil, nil, svc, aes)

	body, _ := json.Marshal(map[string]any{"pass": ""})
	req := httptest.NewRequest("PUT", "/api/admin/smtp/settings", bytes.NewReader(body)).WithContext(ctxWithRole("super_admin"))
	rec := httptest.NewRecorder()
	h.UpdateSMTPSettings(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	if cfg := svc.Config(); cfg.Pass != "" {
		t.Errorf("Pass = %q, want empty after clear", cfg.Pass)
	}
}

func TestUpdateSMTPSettings_EnableWithoutHost(t *testing.T) {
	svc := email.NewService(email.Config{Domain: "example.com"})
	h := NewEmailAdminHandler(nil, nil, svc, nil)

	body, _ := json.Marshal(map[string]any{"enabled": true, "host": ""})
	req := httptest.NewRequest("PUT", "/api/admin/smtp/settings", bytes.NewReader(body)).WithContext(ctxWithRole("admin"))
	rec := httptest.NewRecorder()
	h.UpdateSMTPSettings(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 body=%s", rec.Code, rec.Body.String())
	}
	if !svc.IsDisabled() {
		t.Error("must stay disabled when enable without host")
	}
}

func TestUpdateSMTPSettings_EnablePersists(t *testing.T) {
	svc := email.NewService(email.Config{Host: "smtp.x", From: "a@b.c", Domain: "example.com", Disabled: true})
	h := NewEmailAdminHandler(nil, nil, svc, nil)

	body, _ := json.Marshal(map[string]any{"enabled": true, "host": "smtp.x"})
	req := httptest.NewRequest("PUT", "/api/admin/smtp/settings", bytes.NewReader(body)).WithContext(ctxWithRole("admin"))
	rec := httptest.NewRecorder()
	h.UpdateSMTPSettings(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	if svc.IsDisabled() {
		t.Fatal("expected enabled after PUT enabled=true")
	}
	var resp smtpSettingsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Enabled {
		t.Errorf("response enabled=false, want true: %s", rec.Body.String())
	}

	// disable sticks even with host
	body, _ = json.Marshal(map[string]any{"enabled": false, "host": "smtp.x"})
	req = httptest.NewRequest("PUT", "/api/admin/smtp/settings", bytes.NewReader(body)).WithContext(ctxWithRole("admin"))
	rec = httptest.NewRecorder()
	h.UpdateSMTPSettings(rec, req)
	if !svc.IsDisabled() {
		t.Error("Disabled must stick when enabled=false and host set")
	}
}

func TestTestTemplate_SMTPDisabled(t *testing.T) {
	svc := email.NewService(email.Config{Domain: "example.com"})
	h := NewEmailAdminHandler(nil, nil, svc, nil)

	body, _ := json.Marshal(map[string]any{"email": "user@example.com"})
	req := httptest.NewRequest("POST", "/api/admin/email/templates/test/test", bytes.NewReader(body)).
		WithContext(ctxWithRole("admin"))
	req.SetPathValue("key", "test")
	rec := httptest.NewRecorder()
	h.TestTemplate(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 body=%s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("SMTP")) {
		t.Errorf("expected SMTP error message, got %s", rec.Body.String())
	}
}

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

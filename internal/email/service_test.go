package email

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewService_DisabledWhenNoHost(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	if !svc.IsDisabled() {
		t.Error("expected service to be disabled when host is empty")
	}
}

func TestNewService_EnabledWhenHostSet(t *testing.T) {
	svc := NewService(Config{Host: "smtp.example.com", Port: "587", Domain: "example.com"})
	if svc.IsDisabled() {
		t.Error("expected service to be enabled when host is set")
	}
}

func TestIsDisabled(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	if !svc.IsDisabled() {
		t.Error("expected true")
	}

	svc2 := NewService(Config{Host: "smtp.example.com", Domain: "example.com"})
	if svc2.IsDisabled() {
		t.Error("expected false")
	}
}

func TestConfig(t *testing.T) {
	cfg := Config{Host: "smtp.example.com", Port: "587", User: "user", Pass: "pass", From: "from@example.com", Domain: "example.com"}
	svc := NewService(cfg)
	got := svc.Config()
	if got.Host != cfg.Host {
		t.Errorf("Host = %q, want %q", got.Host, cfg.Host)
	}
	if got.Port != cfg.Port {
		t.Errorf("Port = %q, want %q", got.Port, cfg.Port)
	}
}

func TestSendResetPassword_Disabled(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	err := svc.SendResetPassword("user@example.com", "TestUser", "abc123")
	if err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestSendVerifyEmail_Disabled(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	err := svc.SendVerifyEmail("user@example.com", "TestUser", "abc123")
	if err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestSendTest_Disabled(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	err := svc.SendTest("user@example.com")
	if err == nil {
		t.Error("expected error when SMTP disabled")
	}
}

func TestSendTest_Enabled(t *testing.T) {
	svc := NewService(Config{Host: "127.0.0.1", Port: "1", User: "u", Pass: "p", From: "f@e.com", Domain: "example.com"})
	err := svc.SendTest("user@example.com")
	if err == nil {
		return
	}
	t.Logf("expected SMTP error (no server): %v", err)
}

func TestRenderTemplate(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	body, err := svc.renderTemplate("reset_password", struct {
		Name     string
		ResetURL string
	}{Name: "Test", ResetURL: "https://example.com/reset?token=abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body == "" {
		t.Error("expected non-empty body")
	}
	if !containsStr(body, "https://example.com/reset?token=abc") {
		t.Error("expected body to contain reset URL")
	}
	if !containsStr(body, "Test") {
		t.Error("expected body to contain name")
	}
}

func TestRenderTemplate_VerifyEmail(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	body, err := svc.renderTemplate("verify_email", struct {
		Name      string
		VerifyURL string
	}{Name: "User", VerifyURL: "https://example.com/verify?token=xyz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(body, "https://example.com/verify?token=xyz") {
		t.Error("expected body to contain verify URL")
	}
}

func TestRenderTemplate_Nonexistent(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	// Empty template string from map lookup doesn't panic, just produces empty output
	body, err := svc.renderTemplate("nonexistent", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "" {
		t.Errorf("expected empty body for nonexistent template, got: %q", body)
	}
}

func TestSendResetPassword_Logging(t *testing.T) {
	handler := &captureHandler{}
	slog.SetDefault(slog.New(handler))

	svc := NewService(Config{Host: "", Domain: "mcp-is.ru"})
	err := svc.SendResetPassword("user@example.com", "TestUser", "token123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsStr(handler.output, "mcp-is.ru/reset-password?token=token123") {
		t.Errorf("expected log to contain reset URL, got: %s", handler.output)
	}
}

func TestSendVerifyEmail_Logging(t *testing.T) {
	handler := &captureHandler{}
	slog.SetDefault(slog.New(handler))

	svc := NewService(Config{Host: "", Domain: "mcp-is.ru"})
	err := svc.SendVerifyEmail("user@example.com", "TestUser", "verify456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsStr(handler.output, "mcp-is.ru/verify-email?token=verify456") {
		t.Errorf("expected log to contain verify URL, got: %s", handler.output)
	}
}

// helpers

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type captureHandler struct {
	output string
}

func (h *captureHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *captureHandler) Handle(_ context.Context, r slog.Record) error {
	h.output += r.Message + " "
	r.Attrs(func(a slog.Attr) bool {
		h.output += a.Key + "=" + a.Value.String() + " "
		return true
	})
	return nil
}
func (h *captureHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *captureHandler) WithGroup(_ string) slog.Handler      { return h }

package email

import (
	"context"
	"strings"
	"testing"
)

func TestNewService_DisabledWithoutHost(t *testing.T) {
	svc := NewService(Config{Domain: "example.com"})
	if !svc.IsDisabled() {
		t.Error("expected disabled when Host empty")
	}
}

func TestSendResetPassword_Disabled(t *testing.T) {
	svc := NewService(Config{Domain: "example.com"})
	if err := svc.SendResetPassword("user@example.com", "User", "tok"); err != nil {
		t.Errorf("disabled mode should not error, got %v", err)
	}
}

func TestSendTest_Disabled(t *testing.T) {
	svc := NewService(Config{Host: "", Domain: "example.com"})
	err := svc.SendTest("user@example.com")
	if err == nil {
		t.Error("expected error when SMTP disabled")
	}
}

func TestSendTemplate_RendersSubjectAndBody(t *testing.T) {
	svc := NewService(Config{Host: "smtp.example.com", Port: "587", From: "a@b.c", Domain: "example.com"})
	var gotTo, gotSubject, gotHTML, gotText string
	svc.sender = func(ctx context.Context, cfg Config, to, subject, htmlBody, textBody string) error {
		gotTo, gotSubject, gotHTML, gotText = to, subject, htmlBody, textBody
		return nil
	}

	err := svc.SendTemplate(context.Background(), "reset_password", "user@example.com", map[string]any{
		"Name": "Иван", "ResetURL": "https://example.com/reset?token=x",
	})
	if err != nil {
		t.Fatalf("SendTemplate: %v", err)
	}
	if gotTo != "user@example.com" {
		t.Errorf("to = %q", gotTo)
	}
	if !strings.Contains(gotSubject, "Сброс пароля") {
		t.Errorf("subject = %q", gotSubject)
	}
	if !strings.Contains(gotHTML, "Иван") || !strings.Contains(gotHTML, "example.com/reset") {
		t.Errorf("html missing data: %s", gotHTML[:min(120, len(gotHTML))])
	}
	if !strings.Contains(gotText, "Иван") {
		t.Errorf("text missing data: %q", gotText)
	}
}

func TestSendWithRetry_PermanentNoRetry(t *testing.T) {
	svc := NewService(Config{Host: "smtp.example.com", From: "a@b.c", Domain: "example.com"})
	attempts := 0
	svc.sender = func(ctx context.Context, cfg Config, to, subject, htmlBody, textBody string) error {
		attempts++
		return errPermanent
	}
	err := svc.sendWithRetry(context.Background(), "u@e.c", "s", "h", "t")
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Errorf("permanent error should not retry, attempts=%d", attempts)
	}
}

func TestSendWithRetry_SucceedsAfterTransient(t *testing.T) {
	svc := NewService(Config{Host: "smtp.example.com", From: "a@b.c", Domain: "example.com"})
	attempts := 0
	svc.sender = func(ctx context.Context, cfg Config, to, subject, htmlBody, textBody string) error {
		attempts++
		if attempts == 1 {
			return errTransient
		}
		return nil
	}
	// backoff 1s on first retry — acceptable in unit test
	if err := svc.sendWithRetry(context.Background(), "u@e.c", "s", "h", "t"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("attempts=%d, want 2", attempts)
	}
}

type permErr struct{}

func (permErr) Error() string { return "SMTP 550 permanent" }

type transErr struct{}

func (transErr) Error() string { return "connection timeout" }

var errPermanent = permErr{}
var errTransient = transErr{}

func TestValidateTemplate_OK(t *testing.T) {
	tpl, err := Default("reset_password")
	if err != nil {
		t.Fatal(err)
	}
	err = ValidateTemplate(tpl.Subject, tpl.HTMLBody, tpl.TextBody, tpl.RequiredVars,
		map[string]any{"Name": "T", "ResetURL": "https://x"})
	if err != nil {
		t.Errorf("ValidateTemplate: %v", err)
	}
}

func TestValidateTemplate_MissingVar(t *testing.T) {
	err := ValidateTemplate("Subj", "<p>no vars</p>", "", []string{"ResetURL"}, map[string]any{})
	if err == nil {
		t.Error("expected error for missing required var")
	}
}

func TestValidateTemplate_BadHTML(t *testing.T) {
	err := ValidateTemplate("S", "{{.Name}", "", nil, map[string]any{"Name": "x"})
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestReload(t *testing.T) {
	svc := NewService(Config{Domain: "example.com"})
	if !svc.IsDisabled() {
		t.Fatal("expected disabled")
	}
	svc.Reload(Config{Host: "smtp.x", Port: "465", TLSMode: TLSModeSSL, From: "a@b.c"})
	if svc.IsDisabled() {
		t.Error("expected enabled after reload")
	}
	if svc.Config().TLSMode != TLSModeSSL {
		t.Errorf("tls = %q", svc.Config().TLSMode)
	}
}

func TestReload_PreservesDisabled(t *testing.T) {
	svc := NewService(Config{Host: "smtp.x", From: "a@b.c", Domain: "example.com"})
	svc.Reload(Config{Host: "smtp.x", From: "a@b.c", Disabled: true})
	if !svc.IsDisabled() {
		t.Error("Disabled=true must be preserved when host is set")
	}
	svc.Reload(Config{Host: "smtp.x", From: "a@b.c", Disabled: false})
	if svc.IsDisabled() {
		t.Error("Disabled=false must be preserved when host is set")
	}
	svc.Reload(Config{Host: "", From: "a@b.c", Disabled: false})
	if !svc.IsDisabled() {
		t.Error("empty host must force disabled")
	}
}

func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore()
	tpl, err := s.Get("welcome")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(tpl.HTMLBody, "MCPURL") && !strings.Contains(tpl.HTMLBody, "{{.MCPURL}}") {
		t.Error("welcome template missing MCPURL")
	}
	list, _ := s.List()
	if len(list) < 4 {
		t.Errorf("expected >=4 templates, got %d", len(list))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

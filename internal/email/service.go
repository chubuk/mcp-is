package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/mail"
	"strings"
	"sync"
	"time"

	gomail "github.com/wneessen/go-mail"
)

// TLSMode selects the SMTP transport security.
type TLSMode string

const (
	TLSModeStartTLS TLSMode = "starttls"
	TLSModeSSL      TLSMode = "ssl"
	TLSModeNone     TLSMode = "none"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Pass     string
	From     string
	Domain   string
	TLSMode  TLSMode // default: starttls if empty
	Disabled bool
}

type Service struct {
	mu     sync.RWMutex
	cfg    Config
	store  TemplateStore
	sender func(ctx context.Context, cfg Config, to, subject, htmlBody, textBody string) error
}

// NewService builds a Service. store may be nil (inline defaults are used).
func NewService(cfg Config) *Service {
	if cfg.Host == "" {
		cfg.Disabled = true
		slog.Warn("SMTP not configured, emails will be logged only")
	}
	if cfg.TLSMode == "" {
		cfg.TLSMode = TLSModeStartTLS
	}
	s := &Service{cfg: cfg, store: NewMemoryStore()}
	s.sender = smtpSend
	return s
}

// SetTemplateStore replaces the template source (DB-backed store).
func (s *Service) SetTemplateStore(store TemplateStore) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if store != nil {
		s.store = store
	}
}

// Reload swaps runtime SMTP config without process restart.
func (s *Service) Reload(cfg Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cfg.Host == "" {
		cfg.Disabled = true
	} else {
		cfg.Disabled = false
	}
	if cfg.TLSMode == "" {
		cfg.TLSMode = TLSModeStartTLS
	}
	s.cfg = cfg
	slog.Info("email config reloaded", "host", cfg.Host, "port", cfg.Port, "tls", string(cfg.TLSMode), "disabled", cfg.Disabled)
}

func (s *Service) Config() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

func (s *Service) IsDisabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Disabled
}

func (s *Service) SendResetPassword(to, name, token string) error {
	cfg := s.Config()
	resetURL := fmt.Sprintf("https://%s/reset-password?token=%s", cfg.Domain, token)
	if cfg.Disabled {
		slog.Warn("SMTP disabled — reset URL (share manually)", "to", to, "url", resetURL)
		return nil
	}
	return s.SendTemplate(context.Background(), "reset_password", to, map[string]any{
		"Name": name, "ResetURL": resetURL,
	})
}

// SendTest sends a sample letter using the test template.
func (s *Service) SendTest(to string) error {
	if s.IsDisabled() {
		return fmt.Errorf("SMTP not configured")
	}
	return s.SendTemplate(context.Background(), "test", to, map[string]any{"Name": "Test"})
}

// SendTemplate renders key with data and sends an email with retries.
func (s *Service) SendTemplate(ctx context.Context, key, to string, data map[string]any) error {
	tpl, err := s.getTemplate(key)
	if err != nil {
		return err
	}

	subject, err := renderSubject(tpl.Subject, data)
	if err != nil {
		return fmt.Errorf("render subject: %w", err)
	}
	htmlBody, err := renderHTML(tpl.HTMLBody, data)
	if err != nil {
		return fmt.Errorf("render html: %w", err)
	}
	textBody := ""
	if tpl.TextBody != "" {
		textBody, err = renderText(tpl.TextBody, data)
		if err != nil {
			return fmt.Errorf("render text: %w", err)
		}
	} else {
		textBody = stripTags(htmlBody)
	}

	cfg := s.Config()
	if cfg.Disabled {
		slog.Info("email (SMTP disabled)", "to", to, "subject", subject, "template", key)
		return nil
	}

	return s.sendWithRetry(ctx, to, subject, htmlBody, textBody)
}

func (s *Service) getTemplate(key string) (Template, error) {
	s.mu.RLock()
	store := s.store
	s.mu.RUnlock()
	if store == nil {
		store = NewMemoryStore()
	}
	return store.Get(key)
}

// sendWithRetry: 3 attempts with backoff on transient errors (R2).
func (s *Service) sendWithRetry(ctx context.Context, to, subject, htmlBody, textBody string) error {
	cfg := s.Config()
	var lastErr error
	backoffs := []time.Duration{time.Second, 5 * time.Second, 30 * time.Second}

	for i, wait := range append([]time.Duration{0}, backoffs...) {
		if wait > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		}
		err := s.sender(ctx, cfg, to, subject, htmlBody, textBody)
		if err == nil {
			slog.Info("email sent", "to", to, "subject", subject)
			return nil
		}
		lastErr = err
		if !isTransient(err) {
			break
		}
		slog.Warn("email send attempt failed, will retry", "to", to, "attempt", i+1, "error", err)
	}
	slog.Error("failed to send email", "to", to, "error", lastErr)
	return fmt.Errorf("send email: %w", lastErr)
}

func isTransient(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// 4xx / connection issues are transient; 5xx permanent.
	if strings.Contains(msg, "5") && strings.Contains(msg, "SMTP") {
		return false
	}
	return strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "connection") ||
		strings.Contains(msg, "temporarily") ||
		strings.Contains(msg, "4.7") ||
		strings.Contains(msg, "4.4") ||
		strings.Contains(msg, "421") ||
		strings.Contains(msg, "450") ||
		strings.Contains(msg, "451") ||
		strings.Contains(msg, "452")
}

// smtpSend delivers one message via go-mail (RFC 2047 headers, multipart).
func smtpSend(ctx context.Context, cfg Config, to, subject, htmlBody, textBody string) error {
	port := 587
	switch cfg.TLSMode {
	case TLSModeSSL:
		port = 465
	case TLSModeNone:
		port = 25
	}
	if cfg.Port != "" {
		fmt.Sscanf(cfg.Port, "%d", &port)
	}

	opts := []gomail.Option{
		gomail.WithPort(port),
		gomail.WithTLSPolicy(tlsPolicy(cfg.TLSMode)),
		gomail.WithTimeout(15 * time.Second),
	}
	if cfg.User != "" {
		opts = append(opts, gomail.WithSMTPAuth(gomail.SMTPAuthPlain))
	}

	client, err := gomail.NewClient(cfg.Host, opts...)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	if cfg.User != "" {
		client.SetUsername(cfg.User)
		client.SetPassword(cfg.Pass)
	}

	msg := gomail.NewMsg()
	if err := msg.From(cfg.From); err != nil {
		return fmt.Errorf("from address: %w", err)
	}
	if err := msg.To(to); err != nil {
		return fmt.Errorf("to address: %w", err)
	}
	msg.SetGenHeader(gomail.HeaderSubject, subject)
	msg.SetDate()
	msg.SetMessageID()

	if textBody != "" {
		msg.SetBodyString(gomail.TypeTextPlain, textBody)
		msg.AddAlternativeString(gomail.TypeTextHTML, htmlBody)
	} else {
		msg.SetBodyString(gomail.TypeTextHTML, htmlBody)
	}

	return client.DialAndSendWithContext(ctx, msg)
}

func tlsPolicy(mode TLSMode) gomail.TLSPolicy {
	switch mode {
	case TLSModeSSL:
		return gomail.TLSMandatory
	case TLSModeNone:
		return gomail.NoTLS
	default:
		return gomail.TLSOpportunistic
	}
}

// --- rendering helpers (never template.Must on runtime data) ---

func renderSubject(subject string, data map[string]any) (string, error) {
	if subject == "" {
		return "", fmt.Errorf("empty subject")
	}
	tmpl, err := template.New("subject").Parse(subject)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderHTML(body string, data map[string]any) (string, error) {
	tmpl, err := template.New("html").Parse(body)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderText(body string, data map[string]any) (string, error) {
	tmpl, err := template.New("text").Parse(body)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ValidateTemplate checks parse + sample execution (used by admin PUT).
func ValidateTemplate(subject, htmlBody, textBody string, requiredVars []string, sample map[string]any) error {
	if _, err := renderSubject(subject, sample); err != nil {
		return fmt.Errorf("subject: %w", err)
	}
	if _, err := renderHTML(htmlBody, sample); err != nil {
		return fmt.Errorf("html: %w", err)
	}
	if textBody != "" {
		if _, err := renderText(textBody, sample); err != nil {
			return fmt.Errorf("text: %w", err)
		}
	}
	for _, v := range requiredVars {
		if !strings.Contains(htmlBody, "{{."+v) && !strings.Contains(htmlBody, "{{ ."+v) {
			return fmt.Errorf("required variable %q missing in html body", v)
		}
	}
	return nil
}

func stripTags(html string) string {
	var b strings.Builder
	inTag := false
	for _, r := range html {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

// ParseFrom validates an address like "MCP IS <noreply@mcp-is.ru>".
func ParseFrom(from string) error {
	_, err := mail.ParseAddress(from)
	return err
}

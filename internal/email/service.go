package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Pass     string
	From     string
	Domain   string
	Disabled bool
}

type Service struct {
	cfg       Config
	templates *template.Template
}

func NewService(cfg Config) *Service {
	s := &Service{cfg: cfg}
	if cfg.Host == "" {
		s.cfg.Disabled = true
		slog.Warn("SMTP not configured, emails will be logged only")
	}
	return s
}

func (s *Service) SendResetPassword(to, name, token string) error {
	resetURL := fmt.Sprintf("https://%s/reset-password?token=%s", s.cfg.Domain, token)

	data := struct {
		Name     string
		ResetURL string
	}{
		Name:     name,
		ResetURL: resetURL,
	}

	if s.cfg.Disabled {
		slog.Warn("SMTP disabled — reset URL (share manually)", "to", to, "url", resetURL)
		return nil
	}

	body, err := s.renderTemplate("reset_password", data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	return s.send(to, "Сброс пароля — MCP IS", body)
}

func (s *Service) SendVerifyEmail(to, name, token string) error {
	verifyURL := fmt.Sprintf("https://%s/verify-email?token=%s", s.cfg.Domain, token)

	data := struct {
		Name      string
		VerifyURL string
	}{
		Name:      name,
		VerifyURL: verifyURL,
	}

	if s.cfg.Disabled {
		slog.Warn("SMTP disabled — verify URL (share manually)", "to", to, "url", verifyURL)
		return nil
	}

	body, err := s.renderTemplate("verify_email", data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	return s.send(to, "Подтвердите email — MCP IS", body)
}

func (s *Service) SendTest(to string) error {
	if s.cfg.Disabled {
		return fmt.Errorf("SMTP not configured")
	}

	body, err := s.renderTemplate("reset_password", struct {
		Name     string
		ResetURL string
	}{Name: "Test", ResetURL: "https://example.com/test"})
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	return s.send(to, "SMTP Test — MCP IS", body)
}

func (s *Service) IsDisabled() bool {
	return s.cfg.Disabled
}

func (s *Service) Config() Config {
	return s.cfg
}

func (s *Service) renderTemplate(name string, data any) (string, error) {
	tmpl := template.Must(template.New(name).Parse(templates[name]))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Service) send(to, subject, body string) error {
	if s.cfg.Disabled {
		slog.Info("email (SMTP disabled)", "to", to, "subject", subject)
		return nil
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.cfg.From, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Pass, s.cfg.Host)

	err := smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg))
	if err != nil {
		slog.Error("failed to send email", "to", to, "error", err)
		return fmt.Errorf("send email: %w", err)
	}

	slog.Info("email sent", "to", to, "subject", subject)
	return nil
}

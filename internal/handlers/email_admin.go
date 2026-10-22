package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/email"
	"github.com/chubuk/mcp-is/internal/storage"
)

// EmailAdminHandler serves /api/admin/email/templates* and /api/admin/smtp/settings.
type EmailAdminHandler struct {
	tplRepo  *storage.EmailTemplateRepository
	settings *storage.SettingsRepository
	mail     *email.Service
	aes      *auth.AESCipher
}

func NewEmailAdminHandler(tplRepo *storage.EmailTemplateRepository, settings *storage.SettingsRepository, mail *email.Service, aes *auth.AESCipher) *EmailAdminHandler {
	return &EmailAdminHandler{tplRepo: tplRepo, settings: settings, mail: mail, aes: aes}
}

func (h *EmailAdminHandler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	role, _ := r.Context().Value(auth.RoleKey).(string)
	if role != "admin" && role != "super_admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return false
	}
	return true
}

// --- templates ---

type templateDTO struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	Subject      string   `json:"subject"`
	HTMLBody     string   `json:"html_body"`
	TextBody     string   `json:"text_body"`
	Description  string   `json:"description"`
	RequiredVars []string `json:"required_vars"`
	IsModified   bool     `json:"is_modified"`
	UpdatedAt    string   `json:"updated_at,omitempty"`
}

type templateSaveRequest struct {
	Name         string   `json:"name"`
	Subject      string   `json:"subject"`
	HTMLBody     string   `json:"html_body"`
	TextBody     string   `json:"text_body"`
	Description  string   `json:"description"`
	RequiredVars []string `json:"required_vars"`
}

type templatePreviewRequest struct {
	Subject  string         `json:"subject"`
	HTMLBody string         `json:"html_body"`
	TextBody string         `json:"text_body"`
	Sample   map[string]any `json:"sample"`
}

func (h *EmailAdminHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	items, err := h.tplRepo.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	// If DB empty — expose defaults so UI works before seed.
	if len(items) == 0 {
		for _, d := range email.DefaultList() {
			items = append(items, storage.EmailTemplate{
				Key: d.Key, Name: d.Name, Subject: d.Subject, Description: d.Description,
				RequiredVars: d.RequiredVars,
			})
		}
	}
	out := make([]templateDTO, 0, len(items))
	for _, t := range items {
		out = append(out, templateDTO{
			Key: t.Key, Name: t.Name, Subject: t.Subject, Description: t.Description,
			RequiredVars: t.RequiredVars, IsModified: t.IsModified, UpdatedAt: t.UpdatedAt,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

func (h *EmailAdminHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	key := r.PathValue("key")
	t, err := h.tplRepo.GetByKey(key)
	if err != nil {
		if d, derr := email.Default(key); derr == nil {
			writeJSON(w, http.StatusOK, templateDTO{
				Key: d.Key, Name: d.Name, Subject: d.Subject, HTMLBody: d.HTMLBody,
				TextBody: d.TextBody, Description: d.Description, RequiredVars: d.RequiredVars,
			})
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "шаблон не найден"})
		return
	}
	writeJSON(w, http.StatusOK, templateDTO{
		Key: t.Key, Name: t.Name, Subject: t.Subject, HTMLBody: t.HTMLBody,
		TextBody: t.TextBody, Description: t.Description, RequiredVars: t.RequiredVars,
		IsModified: t.IsModified, UpdatedAt: t.UpdatedAt,
	})
}

func (h *EmailAdminHandler) SaveTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	key := r.PathValue("key")
	var req templateSaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	base, err := email.Default(key)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "неизвестный шаблон"})
		return
	}
	if req.Name == "" {
		req.Name = base.Name
	}
	if req.RequiredVars == nil {
		req.RequiredVars = base.RequiredVars
	}
	if req.Subject == "" || req.HTMLBody == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "тема и тело обязательны"})
		return
	}

	sample := sampleFor(key, req.RequiredVars)
	if err := email.ValidateTemplate(req.Subject, req.HTMLBody, req.TextBody, req.RequiredVars, sample); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "шаблон некорректен: " + err.Error()})
		return
	}

	if err := h.tplRepo.Upsert(key, req.Name, req.Subject, req.HTMLBody, req.TextBody, req.Description, req.RequiredVars); err != nil {
		slog.Error("save email template", "key", key, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	slog.Info("email template saved", "key", key)
	writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

func (h *EmailAdminHandler) ResetTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	key := r.PathValue("key")
	d, err := email.Default(key)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "неизвестный шаблон"})
		return
	}
	if err := h.tplRepo.Upsert(key, d.Name, d.Subject, d.HTMLBody, d.TextBody, d.Description, d.RequiredVars); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	// mark as not modified: re-seed via delete+insert
	_ = h.tplRepo.Delete(key)
	_ = h.tplRepo.SeedIfMissing(d.Key, d.Name, d.Subject, d.HTMLBody, d.TextBody, d.Description, d.RequiredVars)
	writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
}

// PreviewTemplate renders HTML for a sandboxed iframe (never executes scripts server-side).
func (h *EmailAdminHandler) PreviewTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	var req templatePreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Sample == nil {
		req.Sample = map[string]any{"Name": "Пользователь", "ResetURL": "https://example.com", "VerifyURL": "https://example.com", "MCPURL": "https://example.com"}
	}
	// html/template escapes {{.Name}} etc.
	tmpl, perr := template.New("preview").Parse(req.HTMLBody)
	if perr != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ошибка в шаблоне: " + perr.Error()})
		return
	}
	var b strings.Builder
	if err := tmpl.Execute(&b, req.Sample); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ошибка рендера: " + err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"html": b.String()})
}

func (h *EmailAdminHandler) TestTemplate(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	if h.mail.IsDisabled() {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "SMTP выключен — письмо не отправлялось"})
		return
	}
	key := r.PathValue("key")
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email required"})
		return
	}
	// Ensure seeded
	if _, err := h.tplRepo.GetByKey(key); err != nil {
		if d, derr := email.Default(key); derr == nil {
			_ = h.tplRepo.SeedIfMissing(d.Key, d.Name, d.Subject, d.HTMLBody, d.TextBody, d.Description, d.RequiredVars)
		}
	}
	if err := h.mail.SendTemplate(r.Context(), key, req.Email, sampleFor(key, nil)); err != nil {
		slog.Error("template test send", "key", key, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "не удалось отправить: " + err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

func sampleFor(key string, required []string) map[string]any {
	s := map[string]any{
		"Name":      "Пользователь",
		"ResetURL":  "https://mcp-is.ru/reset-password?token=sample",
		"VerifyURL": "https://mcp-is.ru/verify-email?token=sample",
		"MCPURL":    "https://mcp-is.ru/",
	}
	_ = key
	_ = required
	return s
}

// --- SMTP settings ---

type smtpSettingsResponse struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	User    string `json:"user"`
	From    string `json:"from"`
	TLSMode string `json:"tls_mode"`
	HasPass bool   `json:"has_password"`
	Source  string `json:"source"` // "db" | "env"
	Error   string `json:"error,omitempty"`
}

type smtpSettingsRequest struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Host    *string `json:"host,omitempty"`
	Port    *string `json:"port,omitempty"`
	User    *string `json:"user,omitempty"`
	Pass    *string `json:"pass,omitempty"` // empty = keep existing
	From    *string `json:"from,omitempty"`
	TLSMode *string `json:"tls_mode,omitempty"`
}

func (h *EmailAdminHandler) GetSMTPSettings(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	cfg := h.mail.Config()
	writeJSON(w, http.StatusOK, smtpSettingsResponse{
		Enabled: !cfg.Disabled,
		Host:    cfg.Host,
		Port:    cfg.Port,
		User:    cfg.User,
		From:    cfg.From,
		TLSMode: string(cfg.TLSMode),
		HasPass: cfg.Pass != "",
		Source:  "db",
	})
}

func (h *EmailAdminHandler) setSetting(key, value string) error {
	if h.settings == nil {
		return nil
	}
	return h.settings.Set(key, value)
}

func (h *EmailAdminHandler) UpdateSMTPSettings(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}
	var req smtpSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	cfg := h.mail.Config()
	if req.Host != nil {
		cfg.Host = strings.TrimSpace(*req.Host)
	}
	if req.Port != nil {
		cfg.Port = strings.TrimSpace(*req.Port)
	}
	if req.User != nil {
		cfg.User = strings.TrimSpace(*req.User)
	}
	if req.From != nil {
		from := strings.TrimSpace(*req.From)
		if err := email.ParseFrom(from); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "некорректный адрес «От кого»"})
			return
		}
		cfg.From = from
	}
	if req.TLSMode != nil {
		mode := email.TLSMode(strings.ToLower(strings.TrimSpace(*req.TLSMode)))
		switch mode {
		case email.TLSModeSSL, email.TLSModeStartTLS, email.TLSModeNone:
			cfg.TLSMode = mode
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tls_mode: ssl | starttls | none"})
			return
		}
	}
	if req.Pass != nil && *req.Pass != "" {
		cfg.Pass = *req.Pass
		enc, err := h.aes.Encrypt(cfg.Pass)
		if err == nil {
			if err := h.setSetting("smtp_pass_encrypted", enc); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save settings"})
				return
			}
		}
	}

	// Final enabled state: enable requires a host.
	if cfg.Host == "" {
		if req.Enabled != nil && *req.Enabled {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Укажите хост SMTP, чтобы включить отправку"})
			return
		}
		cfg.Disabled = true
	} else if req.Enabled != nil {
		cfg.Disabled = !*req.Enabled
	}

	// Persist non-secret fields
	save := func(k, v string) bool {
		if err := h.setSetting(k, v); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save settings"})
			return false
		}
		return true
	}
	if req.Enabled != nil {
		val := "true"
		if cfg.Disabled {
			val = "false"
		}
		if !save("smtp_enabled", val) {
			return
		}
	}
	if req.Host != nil && !save("smtp_host", cfg.Host) {
		return
	}
	if req.Port != nil && !save("smtp_port", cfg.Port) {
		return
	}
	if req.User != nil && !save("smtp_user", cfg.User) {
		return
	}
	if req.From != nil && !save("smtp_from", cfg.From) {
		return
	}
	if req.TLSMode != nil && !save("smtp_tls_mode", string(cfg.TLSMode)) {
		return
	}

	h.mail.Reload(cfg)
	slog.Info("smtp settings updated", "host", cfg.Host, "port", cfg.Port, "tls", string(cfg.TLSMode), "from", cfg.From, "disabled", cfg.Disabled)
	h.GetSMTPSettings(w, r)
}

// SeedEmailTemplates inserts default templates into DB if missing (startup).
func SeedEmailTemplates(repo *storage.EmailTemplateRepository) {
	if repo == nil {
		return
	}
	for _, d := range email.DefaultList() {
		if err := repo.SeedIfMissing(d.Key, d.Name, d.Subject, d.HTMLBody, d.TextBody, d.Description, d.RequiredVars); err != nil {
			slog.Warn("seed email template", "key", d.Key, "error", err)
		}
	}
}

// LoadSMTPConfigFromSettings merges server_settings over env config.
func LoadSMTPConfigFromSettings(settings *storage.SettingsRepository, env email.Config, aes *auth.AESCipher) email.Config {
	cfg := env
	if settings == nil {
		return cfg
	}
	get := func(k string) string { return settings.GetDefault(k, "") }
	if v := get("smtp_host"); v != "" {
		cfg.Host = v
	}
	if v := get("smtp_port"); v != "" {
		cfg.Port = v
	}
	if v := get("smtp_user"); v != "" {
		cfg.User = v
	}
	if v := get("smtp_from"); v != "" {
		cfg.From = v
	}
	if v := get("smtp_tls_mode"); v != "" {
		cfg.TLSMode = email.TLSMode(v)
	}
	if v := get("smtp_enabled"); v == "false" {
		cfg.Disabled = true
	} else if cfg.Host != "" {
		cfg.Disabled = false
	}
	if enc := get("smtp_pass_encrypted"); enc != "" && aes != nil {
		if pass, err := aes.Decrypt(enc); err == nil {
			cfg.Pass = pass
		}
	}
	return cfg
}

package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
)

type ServerHandler struct {
	settingsRepo *storage.SettingsRepository
}

func NewServerHandler(settingsRepo *storage.SettingsRepository) *ServerHandler {
	return &ServerHandler{settingsRepo: settingsRepo}
}

type ServerSettingsResponse struct {
	Domain             string `json:"domain"`
	HTTPSEnabled       bool   `json:"https_enabled"`
	AutoRenew          bool   `json:"auto_renew"`
	RenewDaysBefore    int    `json:"renew_days_before"`
	CertExpiresAt      string `json:"cert_expires_at"`
	CertIssuer         string `json:"cert_issuer"`
	CertNeedsAttention bool   `json:"cert_needs_attention"`
	LastRenewError     string `json:"last_renew_error"`
	RetryCount         int    `json:"retry_count"`
	LastRetryAt        string `json:"last_retry_at"`
	LetsencryptEmail   string `json:"letsencrypt_email"`
	Error              string `json:"error,omitempty"`
}

type ServerSettingsRequest struct {
	Domain           *string `json:"domain,omitempty"`
	HTTPSEnabled     *bool   `json:"https_enabled,omitempty"`
	AutoRenew        *bool   `json:"auto_renew,omitempty"`
	RenewDaysBefore  *int    `json:"renew_days_before,omitempty"`
	LetsencryptEmail *string `json:"letsencrypt_email,omitempty"`
}

func (h *ServerHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, ServerSettingsResponse{Error: "access denied"})
		return
	}

	settings, err := h.settingsRepo.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ServerSettingsResponse{Error: "failed to load settings"})
		return
	}

	resp := ServerSettingsResponse{
		Domain:             h.settingsRepo.GetDefault("domain", "localhost"),
		HTTPSEnabled:       settings["https_enabled"] == "true",
		AutoRenew:          settings["auto_renew"] == "true",
		RenewDaysBefore:    30,
		CertExpiresAt:      settings["cert_expires_at"],
		CertIssuer:         settings["cert_issuer"],
		CertNeedsAttention: settings["cert_needs_attention"] == "true",
		LastRenewError:     settings["last_renew_error"],
		RetryCount:         0,
		LastRetryAt:        settings["last_retry_at"],
		LetsencryptEmail:   settings["letsencrypt_email"],
	}

	if v, ok := settings["renew_days_before"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &resp.RenewDaysBefore)
	}
	if v, ok := settings["retry_count"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &resp.RetryCount)
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ServerHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" {
		writeJSON(w, http.StatusForbidden, ServerSettingsResponse{Error: "only super_admin can change server settings"})
		return
	}

	var req ServerSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ServerSettingsResponse{Error: "invalid request body"})
		return
	}

	if req.Domain != nil {
		if err := h.settingsRepo.Set("domain", *req.Domain); err != nil {
			writeJSON(w, http.StatusInternalServerError, ServerSettingsResponse{Error: "failed to save domain"})
			return
		}
		slog.Info("domain updated", "domain", *req.Domain)
	}

	if req.HTTPSEnabled != nil {
		value := "false"
		if *req.HTTPSEnabled {
			value = "true"
		}
		_ = h.settingsRepo.Set("https_enabled", value)
	}

	if req.AutoRenew != nil {
		value := "false"
		if *req.AutoRenew {
			value = "true"
		}
		_ = h.settingsRepo.Set("auto_renew", value)
	}

	if req.RenewDaysBefore != nil {
		if *req.RenewDaysBefore < 14 {
			*req.RenewDaysBefore = 14
		}
		if *req.RenewDaysBefore > 60 {
			*req.RenewDaysBefore = 60
		}
		_ = h.settingsRepo.Set("renew_days_before", fmt.Sprintf("%d", *req.RenewDaysBefore))
	}

	if req.LetsencryptEmail != nil {
		_ = h.settingsRepo.Set("letsencrypt_email", *req.LetsencryptEmail)
	}

	h.GetSettings(w, r)
}

func (h *ServerHandler) IssueCert(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	domain := h.settingsRepo.GetDefault("domain", "localhost")
	if domain == "localhost" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "set a real domain before issuing a certificate"})
		return
	}

	email := h.settingsRepo.GetDefault("letsencrypt_email", "")
	if email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "set Let's Encrypt email first"})
		return
	}

	slog.Info("issuing certificate", "domain", domain, "email", email)
	h.settingsRepo.Set("cert_needs_attention", "false")
	h.settingsRepo.Set("last_renew_error", "")
	h.settingsRepo.Set("retry_count", "0")

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "certificate issuance started",
		"domain":  domain,
	})
}

func (h *ServerHandler) RenewCert(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(auth.RoleKey).(string)
	if role != "super_admin" && role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	domain := h.settingsRepo.GetDefault("domain", "localhost")
	slog.Info("renewing certificate", "domain", domain)

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "certificate renewal started",
		"domain":  domain,
	})
}

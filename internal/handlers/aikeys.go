package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
)

type AIKeyHandler struct {
	keyRepo storage.AIKeyRepositoryInterface
}

func NewAIKeyHandler(keyRepo storage.AIKeyRepositoryInterface) *AIKeyHandler {
	return &AIKeyHandler{keyRepo: keyRepo}
}

type AIKeyRequest struct {
	Name       string `json:"name"`
	PDNEnabled bool   `json:"pdn_enabled"`
}

type AIKeyResponse struct {
	LastUsedAt string `json:"last_used_at,omitempty"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key,omitempty"`
	KeyHash    string `json:"key_hash,omitempty"`
	MCPUrl     string `json:"mcp_url,omitempty"`
	PDNEnabled bool   `json:"pdn_enabled"`
	CreatedAt  string `json:"created_at,omitempty"`
	Error      string `json:"error,omitempty"`
}

func (h *AIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var req AIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: "invalid request body"})
		return
	}

	rawKey, keyHash, err := auth.GenerateAPIKey()
	if err != nil {
		slog.Error("generate key", "error", err)
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "internal error"})
		return
	}

	key, err := h.keyRepo.Create(userID, keyHash, req.Name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to create key"})
		return
	}

	if req.PDNEnabled {
		h.keyRepo.UpdatePDNEnabled(key.ID, true)
	}

	slog.Info("ai key created", "id", key.ID, "name", key.Name, "pdn", req.PDNEnabled)

	scheme := "https"
	if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" {
		scheme = "http"
	}
	host := r.Host
	if !strings.Contains(host, ":") {
		if scheme == "https" {
			host += ":443"
		} else {
			host += ":8080"
		}
	}
	mcpURL := scheme + "://" + host + "/mcp"

	writeJSON(w, http.StatusCreated, AIKeyResponse{
		ID: key.ID, Name: key.Name, Key: rawKey, KeyHash: keyHash, MCPUrl: mcpURL, PDNEnabled: req.PDNEnabled,
	})
}

func (h *AIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	keys, total, err := h.keyRepo.ListByUserIDFiltered(userID, page, perPage)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list keys"})
		return
	}

	result := make([]AIKeyResponse, len(keys))
	for i, k := range keys {
		result[i] = AIKeyResponse{
			ID: k.ID, Name: k.Name, KeyHash: k.KeyHash, PDNEnabled: k.PDNEnabled, CreatedAt: k.CreatedAt, LastUsedAt: k.LastUsedAt,
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":    result,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *AIKeyHandler) Update(w http.ResponseWriter, r *http.Request) {
	keyID := r.PathValue("id")

	var req AIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: "invalid request body"})
		return
	}

	if req.Name != "" {
		if err := h.keyRepo.UpdateName(keyID, req.Name); err != nil {
			writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to update name"})
			return
		}
	}

	if err := h.keyRepo.UpdatePDNEnabled(keyID, req.PDNEnabled); err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to update pdn"})
		return
	}

	writeJSON(w, http.StatusOK, AIKeyResponse{ID: keyID, Name: req.Name, PDNEnabled: req.PDNEnabled})
}

func (h *AIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	keyID := r.PathValue("id")

	if err := h.keyRepo.Delete(keyID); err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to delete key"})
		return
	}

	writeJSON(w, http.StatusOK, AIKeyResponse{ID: keyID})
}

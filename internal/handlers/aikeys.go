package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/storage"
)

type AIKeyHandler struct {
	keyRepo  storage.AIKeyRepositoryInterface
	shopRepo storage.ShopRepositoryInterface
	maxKeys  int // manual creation quota (0 = unlimited)
}

func NewAIKeyHandler(keyRepo storage.AIKeyRepositoryInterface, shopRepo storage.ShopRepositoryInterface, maxKeys int) *AIKeyHandler {
	return &AIKeyHandler{keyRepo: keyRepo, shopRepo: shopRepo, maxKeys: maxKeys}
}

// AIKeyRequest is a merge-patch: nil fields are left unchanged.
type AIKeyRequest struct {
	Name           *string `json:"name"`
	PDNEnabled     *bool   `json:"pdn_enabled"`
	ShopID         *string `json:"shop_id"`
	ConnectionName *string `json:"connection_name"`
}

type AIKeyResponse struct {
	LastUsedAt     string `json:"last_used_at,omitempty"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Key            string `json:"key,omitempty"`
	KeyHash        string `json:"key_hash,omitempty"`
	MCPUrl         string `json:"mcp_url,omitempty"`
	ShopID         string `json:"shop_id,omitempty"`
	ShopName       string `json:"shop_name,omitempty"`
	ConnectionName string `json:"connection_name,omitempty"`
	PDNEnabled     bool   `json:"pdn_enabled"`
	CreatedAt      string `json:"created_at,omitempty"`
	Error          string `json:"error,omitempty"`
}

var connectionNameRe = regexp.MustCompile(`^[a-z0-9-]{1,64}$`)

// uniqueKeyConnectionName returns base or base-2/-3... free of collisions among
// the user's stored connection names (excluding excludeID).
func uniqueKeyConnectionName(repo storage.AIKeyRepositoryInterface, userID, base, excludeID string) string {
	taken := map[string]bool{}
	if keys, err := repo.ListByUserID(userID); err == nil {
		for _, k := range keys {
			if k.ID != excludeID && k.ConnectionName != "" {
				taken[k.ConnectionName] = true
			}
		}
	}
	if !taken[base] {
		return base
	}
	for i := 2; ; i++ {
		cand := fmt.Sprintf("%s-%d", base, i)
		if !taken[cand] {
			return cand
		}
	}
}

func (h *AIKeyHandler) uniqueConnectionName(userID, base, excludeID string) string {
	return uniqueKeyConnectionName(h.keyRepo, userID, base, excludeID)
}

func (h *AIKeyHandler) resolveConnectionName(userID, requested, shopName, current, excludeID string) (string, error) {
	name := requested
	if name == "" {
		name = current
	}
	if name == "" {
		name = core.Slugify(shopName)
	}
	if !connectionNameRe.MatchString(name) {
		return "", fmt.Errorf("connection_name: допустимы только латинские буквы в нижнем регистре, цифры и дефис (1–64 символов)")
	}
	return h.uniqueConnectionName(userID, name, excludeID), nil
}

func (h *AIKeyHandler) ownedShop(userID, shopID string) (*storage.Shop, error) {
	shop, err := h.shopRepo.GetByID(shopID)
	if err != nil || shop.UserID != userID {
		return nil, fmt.Errorf("магазин не найден")
	}
	return shop, nil
}

func (h *AIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var req AIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: "invalid request body"})
		return
	}

	if h.maxKeys > 0 {
		if _, total, err := h.keyRepo.ListByUserIDFiltered(userID, 1, 1); err == nil && total >= h.maxKeys {
			writeJSON(w, http.StatusForbidden, AIKeyResponse{Error: fmt.Sprintf("Достигнут лимит ИИ-ключей (%d)", h.maxKeys)})
			return
		}
	}

	shopID := ""
	if req.ShopID != nil {
		shopID = strings.TrimSpace(*req.ShopID)
	}
	if shopID == "" {
		msg := "Укажите магазин (shop_id) — ИИ-ключ должен быть привязан к магазину."
		if shops, err := h.shopRepo.ListByUserID(userID); err == nil && len(shops) == 0 {
			msg = "Сначала добавьте магазин — ИИ-ключ привязывается к магазину."
		}
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: msg})
		return
	}

	shop, err := h.ownedShop(userID, shopID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: err.Error()})
		return
	}

	reqName := ""
	if req.Name != nil {
		reqName = *req.Name
	}
	connName, err := h.resolveConnectionName(userID, derefStr(req.ConnectionName), shop.Name, "", "")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: err.Error()})
		return
	}

	rawKey, keyHash, err := auth.GenerateAPIKey()
	if err != nil {
		slog.Error("generate key", "error", err)
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "internal error"})
		return
	}

	key, err := h.keyRepo.CreateForShop(userID, shop.ID, keyHash, reqName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to create key"})
		return
	}
	if err := h.keyRepo.UpdateConnectionName(key.ID, connName); err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to update key"})
		return
	}
	if req.PDNEnabled != nil && *req.PDNEnabled {
		if err := h.keyRepo.UpdatePDNEnabled(key.ID, true); err != nil {
			writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to update key"})
			return
		}
	}

	slog.Info("ai key created", "id", key.ID, "name", reqName, "shop_id", shop.ID, "connection_name", connName)

	mcpURL := buildMCPURL(r)

	writeJSON(w, http.StatusCreated, AIKeyResponse{
		ID: key.ID, Name: reqName, Key: rawKey, KeyHash: keyHash, MCPUrl: mcpURL,
		ShopID: shop.ID, ShopName: shop.Name, ConnectionName: connName,
		PDNEnabled: req.PDNEnabled != nil && *req.PDNEnabled,
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
			ID: k.ID, Name: k.Name, KeyHash: k.KeyHash,
			ShopID: k.ShopID, ShopName: k.ShopName, ConnectionName: k.ConnectionName,
			PDNEnabled: k.PDNEnabled, CreatedAt: k.CreatedAt, LastUsedAt: k.LastUsedAt,
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
	userID := r.Context().Value(auth.UserIDKey).(string)
	keyID := r.PathValue("id")

	key, err := h.keyRepo.GetByID(keyID)
	if err != nil || key.UserID != userID {
		writeJSON(w, http.StatusNotFound, AIKeyResponse{Error: "ключ не найден"})
		return
	}

	var req AIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: "invalid request body"})
		return
	}

	// Merge-patch: an unbound key cannot be saved/changed without a shop binding.
	newShopID := key.ShopID
	newShopName := key.ShopName
	if req.ShopID != nil && strings.TrimSpace(*req.ShopID) != "" {
		shop, err := h.ownedShop(userID, strings.TrimSpace(*req.ShopID))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: err.Error()})
			return
		}
		newShopID, newShopName = shop.ID, shop.Name
	} else if key.ShopID == "" {
		writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: "AI-ключ не привязан к магазину. Укажите shop_id при сохранении."})
		return
	}

	reqConn := derefStr(req.ConnectionName)
	if req.ShopID != nil || reqConn != "" {
		connName, err := h.resolveConnectionName(userID, reqConn, newShopName, key.ConnectionName, keyID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, AIKeyResponse{Error: err.Error()})
			return
		}
		if err := h.keyRepo.UpdateShopBinding(keyID, newShopID, connName); err != nil {
			writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to update key"})
			return
		}
	}

	if req.Name != nil && *req.Name != "" {
		if err := h.keyRepo.UpdateName(keyID, *req.Name); err != nil {
			writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to update name"})
			return
		}
	}
	if req.PDNEnabled != nil {
		if err := h.keyRepo.UpdatePDNEnabled(keyID, *req.PDNEnabled); err != nil {
			writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to update pdn"})
			return
		}
	}

	updated, err := h.keyRepo.GetByID(keyID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to read key"})
		return
	}
	writeJSON(w, http.StatusOK, AIKeyResponse{
		ID: updated.ID, Name: updated.Name,
		ShopID: updated.ShopID, ShopName: updated.ShopName, ConnectionName: updated.ConnectionName,
		PDNEnabled: updated.PDNEnabled,
	})
}

// Rotate issues a fresh secret for an existing key; the old secret stops
// authenticating immediately. Name, binding and connection_name are preserved.
func (h *AIKeyHandler) Rotate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	keyID := r.PathValue("id")

	key, err := h.keyRepo.GetByID(keyID)
	if err != nil || key.UserID != userID {
		writeJSON(w, http.StatusNotFound, AIKeyResponse{Error: "ключ не найден"})
		return
	}

	rawKey, keyHash, err := auth.GenerateAPIKey()
	if err != nil {
		slog.Error("rotate key", "error", err)
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "internal error"})
		return
	}
	if err := h.keyRepo.UpdateKeyHash(keyID, keyHash); err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to rotate key"})
		return
	}

	slog.Info("ai key rotated", "id", keyID)
	writeJSON(w, http.StatusOK, AIKeyResponse{
		ID: key.ID, Name: key.Name, Key: rawKey, KeyHash: keyHash, MCPUrl: buildMCPURL(r),
		ShopID: key.ShopID, ShopName: key.ShopName, ConnectionName: key.ConnectionName,
		PDNEnabled: key.PDNEnabled,
	})
}

func (h *AIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	keyID := r.PathValue("id")

	key, err := h.keyRepo.GetByID(keyID)
	if err != nil || key.UserID != userID {
		writeJSON(w, http.StatusNotFound, AIKeyResponse{Error: "ключ не найден"})
		return
	}

	if err := h.keyRepo.Delete(keyID); err != nil {
		writeJSON(w, http.StatusInternalServerError, AIKeyResponse{Error: "failed to delete key"})
		return
	}

	writeJSON(w, http.StatusOK, AIKeyResponse{ID: keyID})
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}

func buildMCPURL(r *http.Request) string {
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
	return scheme + "://" + host + "/mcp"
}

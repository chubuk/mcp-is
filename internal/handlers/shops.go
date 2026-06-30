package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/your-project/mcp-insales/internal/auth"
	"github.com/your-project/mcp-insales/internal/storage"
)

type ShopHandler struct {
	shopRepo storage.ShopRepositoryInterface
	aes      *auth.AESCipher
}

func NewShopHandler(shopRepo storage.ShopRepositoryInterface, aes *auth.AESCipher) *ShopHandler {
	return &ShopHandler{shopRepo: shopRepo, aes: aes}
}

type ShopRequest struct {
	ShopURL     string `json:"shop_url"`
	APIKeyID    string `json:"api_key_id"`
	APIPassword string `json:"api_password"`
	IsDefault   bool   `json:"is_default"`
	IsReadOnly  bool   `json:"is_read_only"`
}

type ShopResponse struct {
	ID         string `json:"id"`
	ShopURL    string `json:"shop_url"`
	APIKeyID   string `json:"api_key_id"`
	IsDefault  bool   `json:"is_default"`
	IsReadOnly bool   `json:"is_read_only"`
	CreatedAt  string `json:"created_at,omitempty"`
	Error      string `json:"error,omitempty"`
}

func (h *ShopHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var req ShopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ShopResponse{Error: "invalid request body"})
		return
	}

	if req.ShopURL == "" || req.APIKeyID == "" || req.APIPassword == "" {
		writeJSON(w, http.StatusBadRequest, ShopResponse{Error: "shop_url, api_key_id and api_password are required"})
		return
	}

	encrypted, err := h.aes.Encrypt(req.APIPassword)
	if err != nil {
		slog.Error("encrypt password", "error", err)
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "encryption error"})
		return
	}

	shop, err := h.shopRepo.Create(userID, req.ShopURL, req.APIKeyID, encrypted, req.IsDefault, req.IsReadOnly)
	if err != nil {
		slog.Warn("create shop failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to create shop"})
		return
	}

	slog.Info("shop created", "id", shop.ID, "url", shop.ShopURL)
	writeJSON(w, http.StatusCreated, ShopResponse{
		ID: shop.ID, ShopURL: shop.ShopURL, APIKeyID: shop.APIKeyID, IsDefault: shop.IsDefault, IsReadOnly: shop.IsReadOnly,
	})
}

func (h *ShopHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	shops, total, err := h.shopRepo.ListByUserIDFiltered(userID, page, perPage)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list shops"})
		return
	}

	result := make([]ShopResponse, len(shops))
	for i, s := range shops {
		result[i] = ShopResponse{
			ID: s.ID, ShopURL: s.ShopURL, APIKeyID: s.APIKeyID, IsDefault: s.IsDefault, IsReadOnly: s.IsReadOnly, CreatedAt: s.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items":    result,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *ShopHandler) Update(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("id")

	var req ShopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ShopResponse{Error: "invalid request body"})
		return
	}

	encrypted := ""
	if req.APIPassword != "" {
		var err error
		encrypted, err = h.aes.Encrypt(req.APIPassword)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "encryption error"})
			return
		}
	}

	// Get existing shop to preserve encrypted password if not provided
	existing, err := h.shopRepo.GetByID(shopID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ShopResponse{Error: "shop not found"})
		return
	}

	if encrypted == "" {
		encrypted = existing.APIPasswordEncrypted
	}

	if err := h.shopRepo.Update(shopID, req.ShopURL, req.APIKeyID, encrypted, req.IsDefault, req.IsReadOnly); err != nil {
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to update shop"})
		return
	}

	writeJSON(w, http.StatusOK, ShopResponse{ID: shopID, ShopURL: req.ShopURL, APIKeyID: req.APIKeyID, IsDefault: req.IsDefault, IsReadOnly: req.IsReadOnly})
}

func (h *ShopHandler) Delete(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("id")

	if err := h.shopRepo.Delete(shopID); err != nil {
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to delete shop"})
		return
	}

	writeJSON(w, http.StatusOK, ShopResponse{ID: shopID})
}

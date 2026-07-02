package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/your-project/mcp-insales/internal/auth"
	"github.com/your-project/mcp-insales/internal/integration/insales"
	"github.com/your-project/mcp-insales/internal/mcp"
	"github.com/your-project/mcp-insales/internal/storage"
)

type ShopHandler struct {
	shopRepo   storage.ShopRepositoryInterface
	aiKeyRepo  storage.AIKeyRepositoryInterface
	aes        *auth.AESCipher
	mcpServer  *mcp.Server
}

func NewShopHandler(shopRepo storage.ShopRepositoryInterface, aiKeyRepo storage.AIKeyRepositoryInterface, aes *auth.AESCipher, mcpServer *mcp.Server) *ShopHandler {
	return &ShopHandler{shopRepo: shopRepo, aiKeyRepo: aiKeyRepo, aes: aes, mcpServer: mcpServer}
}

type ShopRequest struct {
	Name        string `json:"name"`
	ShopURL     string `json:"shop_url"`
	APIKeyID    string `json:"api_key_id"`
	APIPassword string `json:"api_password"`
	IsDefault   bool   `json:"is_default"`
	IsReadOnly  bool   `json:"is_read_only"`
}

type ShopResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
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

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, ShopResponse{Error: "name is required"})
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

	shop, err := h.shopRepo.Create(userID, req.Name, req.ShopURL, req.APIKeyID, encrypted, req.IsDefault, req.IsReadOnly)
	if err != nil {
		slog.Warn("create shop failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to create shop"})
		return
	}

	slog.Info("shop created", "id", shop.ID, "name", shop.Name, "url", shop.ShopURL)
	writeJSON(w, http.StatusCreated, ShopResponse{
		ID: shop.ID, Name: shop.Name, ShopURL: shop.ShopURL, APIKeyID: shop.APIKeyID, IsDefault: shop.IsDefault, IsReadOnly: shop.IsReadOnly,
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
			ID: s.ID, Name: s.Name, ShopURL: s.ShopURL, APIKeyID: s.APIKeyID, IsDefault: s.IsDefault, IsReadOnly: s.IsReadOnly, CreatedAt: s.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
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

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, ShopResponse{Error: "name is required"})
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

	existing, err := h.shopRepo.GetByID(shopID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ShopResponse{Error: "shop not found"})
		return
	}

	if encrypted == "" {
		encrypted = existing.APIPasswordEncrypted
	}

	if req.ShopURL == "" {
		req.ShopURL = existing.ShopURL
	}
	if req.APIKeyID == "" {
		req.APIKeyID = existing.APIKeyID
	}

	if err := h.shopRepo.Update(shopID, req.Name, req.ShopURL, req.APIKeyID, encrypted, req.IsDefault, req.IsReadOnly); err != nil {
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to update shop"})
		return
	}

	writeJSON(w, http.StatusOK, ShopResponse{ID: shopID, Name: req.Name, ShopURL: req.ShopURL, APIKeyID: req.APIKeyID, IsDefault: req.IsDefault, IsReadOnly: req.IsReadOnly})
}

func (h *ShopHandler) Delete(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("id")

	if err := h.shopRepo.Delete(shopID); err != nil {
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to delete shop"})
		return
	}

	writeJSON(w, http.StatusOK, ShopResponse{ID: shopID})
}

type TestConnectionRequest struct {
	ShopURL     string `json:"shop_url"`
	APIKeyID    string `json:"api_key_id"`
	APIPassword string `json:"api_password"`
}

type TestConnectionResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Shop    *struct {
		Name     string `json:"name"`
		Domain   string `json:"domain"`
		Email    string `json:"email"`
		Currency string `json:"currency"`
	} `json:"shop,omitempty"`
}

func (h *ShopHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	var req TestConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, TestConnectionResponse{Error: "invalid request body"})
		return
	}

	if req.ShopURL == "" || req.APIKeyID == "" || req.APIPassword == "" {
		writeJSON(w, http.StatusBadRequest, TestConnectionResponse{Error: "shop_url, api_key_id and api_password are required"})
		return
	}

	client, err := insales.NewClient(req.ShopURL, req.APIKeyID, req.APIPassword)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TestConnectionResponse{Error: err.Error()})
		return
	}

	info, err := client.GetStoreInfo(r.Context())
	if err != nil {
		slog.Warn("test connection failed", "shop", req.ShopURL, "error", err)
		writeJSON(w, http.StatusOK, TestConnectionResponse{OK: false, Error: err.Error()})
		return
	}

	slog.Info("test connection success", "shop", req.ShopURL, "name", info.Name)
	writeJSON(w, http.StatusOK, TestConnectionResponse{
		OK:      true,
		Message: "Соединение успешно установлено",
		Shop: &struct {
			Name     string `json:"name"`
			Domain   string `json:"domain"`
			Email    string `json:"email"`
			Currency string `json:"currency"`
		}{
			Name:     info.Name,
			Domain:   info.Domain,
			Email:    info.Email,
			Currency: info.Currency,
		},
	})
}

type CheckCommandsResponse struct {
	Commands []string `json:"commands"`
	Count    int      `json:"count"`
	Error    string   `json:"error,omitempty"`
}

var readCommands = []string{
	"list_products", "get_product",
	"list_orders", "get_order",
	"list_clients", "get_client",
	"list_categories", "get_category",
	"list_collections", "get_collection",
	"list_images",
	"list_meta_fields",
	"list_webhooks",
	"list_pages", "get_page",
	"list_accessories", "get_accessory",
	"list_variants", "get_variant",
	"list_property_values",
	"list_menus", "get_menu",
	"list_redirects",
	"list_blogs", "get_blog",
	"list_blog_posts", "get_blog_post",
	"get_shop_info",
	"list_currencies", "list_delivery_variants", "list_payment_gateways",
}

func (h *ShopHandler) CheckCommands(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("id")

	shop, err := h.shopRepo.GetByID(shopID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, CheckCommandsResponse{Error: "shop not found"})
		return
	}

	var commands []string
	if shop.IsReadOnly {
		commands = readCommands
	} else {
		allTools := mcp.GetAllTools()
		commands = make([]string, len(allTools))
		for i, t := range allTools {
			commands[i] = t.Name
		}
	}

	commandsJSON, _ := json.Marshal(commands)
	_ = h.aiKeyRepo.UpdateCommandsByUserID(shop.UserID, string(commandsJSON))

	slog.Info("commands updated", "shop_id", shopID, "count", len(commands), "read_only", shop.IsReadOnly)
	writeJSON(w, http.StatusOK, CheckCommandsResponse{
		Commands: commands,
		Count:    len(commands),
	})
}

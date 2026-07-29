package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/chubuk/mcp-is/internal/mcp"
	"github.com/chubuk/mcp-is/internal/storage"
)

type ShopHandler struct {
	shopRepo   storage.ShopRepositoryInterface
	aiKeyRepo  storage.AIKeyRepositoryInterface
	aes        *auth.AESCipher
	mcpServer  *mcp.MCPServer
}

func NewShopHandler(shopRepo storage.ShopRepositoryInterface, aiKeyRepo storage.AIKeyRepositoryInterface, aes *auth.AESCipher, mcpServer *mcp.MCPServer) *ShopHandler {
	return &ShopHandler{shopRepo: shopRepo, aiKeyRepo: aiKeyRepo, aes: aes, mcpServer: mcpServer}
}

type ShopRequest struct {
	Name        string `json:"name"`
	ShopURL     string `json:"shop_url"`
	APIKeyID    string `json:"api_key_id"`
	APIPassword string `json:"api_password"`
	IsDefault   bool   `json:"is_default"`
	IsReadOnly  bool   `json:"is_read_only"`
	PDNEnabled  bool   `json:"pdn_enabled"`
}

type ShopResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ShopURL      string `json:"shop_url"`
	APIKeyID     string `json:"api_key_id"`
	IsDefault    bool   `json:"is_default"`
	IsReadOnly   bool   `json:"is_read_only"`
	PDNEnabled   bool   `json:"pdn_enabled"`
	Jurisdiction string `json:"jurisdiction,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Error        string `json:"error,omitempty"`
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
	role := r.Context().Value(auth.RoleKey).(string)
	isAdmin := role == "admin" || role == "super_admin"

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
		resp := ShopResponse{
			ID: s.ID, Name: s.Name, ShopURL: s.ShopURL, APIKeyID: s.APIKeyID, IsDefault: s.IsDefault, IsReadOnly: s.IsReadOnly, PDNEnabled: s.PDNEnabled, CreatedAt: s.CreatedAt,
		}
		if isAdmin {
			resp.Jurisdiction = s.Jurisdiction
		}
		result[i] = resp
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

	if err := h.shopRepo.Update(shopID, req.Name, req.ShopURL, req.APIKeyID, encrypted, req.IsDefault, req.IsReadOnly, req.PDNEnabled); err != nil {
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to update shop"})
		return
	}

	writeJSON(w, http.StatusOK, ShopResponse{ID: shopID, Name: req.Name, ShopURL: req.ShopURL, APIKeyID: req.APIKeyID, IsDefault: req.IsDefault, IsReadOnly: req.IsReadOnly, PDNEnabled: req.PDNEnabled})
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

type CheckCommandItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CheckCommandsResponse struct {
	Commands  []CheckCommandItem `json:"commands"`
	Count     int                `json:"count"`
	UpdatedAt string             `json:"updated_at,omitempty"`
	Error     string             `json:"error,omitempty"`
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
	"list_outlets", "get_outlet",
	"list_marketplaces", "get_marketplace",
	"get_orders_count", "get_products_count",
	"collection_products",
	"export_products", "export_orders", "export_clients",
	"search_product_details", "order_full_details", "client_full_profile", "product_catalog_search",
}

var allCommandNames = []string{
	"list_products", "get_product", "create_product", "update_product", "delete_product",
	"list_orders", "get_order", "update_order", "create_order",
	"list_clients", "get_client", "create_client", "update_client", "delete_client",
	"list_categories", "get_category", "create_category", "update_category", "delete_category",
	"list_collections", "get_collection", "create_collection", "update_collection", "delete_collection",
	"list_images", "add_image", "delete_image",
	"list_meta_fields", "create_meta_field", "update_meta_field", "delete_meta_field",
	"list_webhooks", "create_webhook", "update_webhook", "delete_webhook",
	"get_shop_info", "list_currencies", "list_delivery_variants", "list_payment_gateways",
	"list_pages", "get_page", "create_page", "update_page", "delete_page",
	"list_accessories", "get_accessory", "create_accessory", "update_accessory", "delete_accessory",
	"list_variants", "get_variant", "update_variant", "delete_variant",
	"list_property_values", "create_property_value", "update_property_value", "delete_property_value",
	"list_menus", "get_menu", "create_menu", "update_menu", "delete_menu",
	"list_redirects", "create_redirect", "delete_redirect",
	"list_blogs", "get_blog", "create_blog", "delete_blog",
	"list_blog_posts", "get_blog_post", "create_blog_post", "update_blog_post", "delete_blog_post",
	"list_outlets", "get_outlet", "create_outlet", "update_outlet", "delete_outlet",
	"list_marketplaces", "get_marketplace", "create_marketplace", "update_marketplace", "delete_marketplace",
	"create_provider_token", "update_account",
	"get_orders_count", "get_products_count",
	"collection_products",
	"export_products", "export_orders", "export_clients",
	"search_product_details", "order_full_details", "client_full_profile", "product_catalog_search",
}

var commandDescriptions = map[string]string{
	"list_products":         "Получить список товаров",
	"get_product":           "Получить товар по ID",
	"create_product":        "Создать товар",
	"update_product":        "Обновить товар",
	"delete_product":        "Удалить товар",
	"list_orders":           "Получить список заказов",
	"get_order":             "Получить заказ по ID",
	"update_order":          "Обновить заказ",
	"create_order":          "Создать заказ",
	"list_clients":          "Получить список клиентов",
	"get_client":            "Получить клиента по ID",
	"create_client":         "Создать клиента",
	"update_client":         "Обновить клиента",
	"delete_client":         "Удалить клиента",
	"list_categories":       "Получить список категорий",
	"get_category":          "Получить категорию по ID",
	"create_category":       "Создать категорию",
	"update_category":       "Обновить категорию",
	"delete_category":       "Удалить категорию",
	"list_collections":      "Получить список коллекций",
	"get_collection":        "Получить коллекцию по ID",
	"create_collection":     "Создать коллекцию",
	"update_collection":     "Обновить коллекцию",
	"delete_collection":     "Удалить коллекцию",
	"list_images":           "Получить изображения товара",
	"add_image":             "Добавить изображение",
	"delete_image":          "Удалить изображение",
	"list_meta_fields":      "Получить список метаполей",
	"create_meta_field":     "Создать метаполе",
	"update_meta_field":     "Обновить метаполе",
	"delete_meta_field":     "Удалить метаполе",
	"list_webhooks":         "Получить список вебхуков",
	"create_webhook":        "Создать вебхук",
	"update_webhook":        "Обновить вебхук",
	"delete_webhook":        "Удалить вебхук",
	"get_shop_info":         "Получить информацию о магазине",
	"list_currencies":       "Получить список валют",
	"list_delivery_variants": "Получить способы доставки",
	"list_payment_gateways": "Получить способы оплаты",
	"list_pages":            "Получить список страниц",
	"get_page":              "Получить страницу по ID",
	"create_page":           "Создать страницу",
	"update_page":           "Обновить страницу",
	"delete_page":           "Удалить страницу",
	"list_accessories":      "Получить список аксессуаров",
	"get_accessory":         "Получить аксессуар по ID",
	"create_accessory":      "Создать аксессуар",
	"update_accessory":      "Обновить аксессуар",
	"delete_accessory":      "Удалить аксессуар",
	"list_variants":         "Получить варианты товара",
	"get_variant":           "Получить вариант по ID",
	"update_variant":        "Обновить вариант",
	"delete_variant":        "Удалить вариант",
	"list_property_values":  "Получить значения свойств",
	"create_property_value": "Создать значение свойства",
	"update_property_value": "Обновить значение свойства",
	"delete_property_value": "Удалить значение свойства",
	"list_menus":            "Получить список меню",
	"get_menu":              "Получить меню по ID",
	"create_menu":           "Создать меню",
	"update_menu":           "Обновить меню",
	"delete_menu":           "Удалить меню",
	"list_redirects":        "Получить список редиректов",
	"create_redirect":       "Создать редирект",
	"delete_redirect":       "Удалить редирект",
	"list_blogs":            "Получить список блогов",
	"get_blog":              "Получить блог по ID",
	"create_blog":           "Создать блог",
	"delete_blog":           "Удалить блог",
	"list_blog_posts":       "Получить записи блога",
	"get_blog_post":         "Получить запись блога по ID",
	"create_blog_post":      "Создать запись блога",
	"update_blog_post":      "Обновить запись блога",
	"delete_blog_post":      "Удалить запись блога",
	"list_outlets":          "Получить список точек выдачи",
	"get_outlet":            "Получить точку выдачи по ID",
	"create_outlet":         "Создать точку выдачи",
	"update_outlet":         "Обновить точку выдачи",
	"delete_outlet":         "Удалить точку выдачи",
	"list_marketplaces":     "Получить список маркетплейсов",
	"get_marketplace":       "Получить маркетплейс по ID",
	"create_marketplace":    "Создать маркетплейс",
	"update_marketplace":    "Обновить маркетплейс",
	"delete_marketplace":    "Удалить маркетплейс",
	"create_provider_token": "Создать токен провайдера мессенджера",
	"update_account":        "Обновить настройки аккаунта",
	"get_orders_count":      "Получить количество заказов",
	"get_products_count":    "Получить количество товаров",
	"collection_products":   "Товары коллекции с полной информацией",
	"export_products":       "Экспорт товаров в CSV",
	"export_orders":         "Экспорт заказов в CSV",
	"export_clients":        "Экспорт клиентов в CSV",
	"search_product_details": "Поиск товара с полной информацией",
	"order_full_details":    "Полная информация о заказе",
	"client_full_profile":   "Полный профиль клиента",
	"product_catalog_search": "Поиск по каталогу с фильтрами",
}

func (h *ShopHandler) CheckCommands(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("id")

	shop, err := h.shopRepo.GetByID(shopID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, CheckCommandsResponse{Error: "shop not found"})
		return
	}

	var commandNames []string
	if shop.IsReadOnly {
		commandNames = readCommands
	} else {
		commandNames = allCommandNames
	}

	items := make([]CheckCommandItem, len(commandNames))
	for i, name := range commandNames {
		items[i] = CheckCommandItem{
			Name:        name,
			Description: commandDescriptions[name],
		}
	}

	commandsJSON, _ := json.Marshal(commandNames)
	_ = h.aiKeyRepo.UpdateCommandsByUserID(shop.UserID, string(commandsJSON))

	now := time.Now().Format("2006-01-02 15:04")
	slog.Info("commands updated", "shop_id", shopID, "count", len(commandNames), "read_only", shop.IsReadOnly)
	writeJSON(w, http.StatusOK, CheckCommandsResponse{
		Commands:  items,
		Count:     len(commandNames),
		UpdatedAt: now,
	})
}

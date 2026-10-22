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
	shopRepo  storage.ShopRepositoryInterface
	aiKeyRepo storage.AIKeyRepositoryInterface
	aes       *auth.AESCipher
	mcpServer *mcp.MCPServer
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

	shop, err := h.shopRepo.Create(userID, req.Name, req.ShopURL, req.APIKeyID, encrypted, req.IsDefault, req.IsReadOnly, req.PDNEnabled)
	if err != nil {
		slog.Warn("create shop failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, ShopResponse{Error: "failed to create shop"})
		return
	}

	slog.Info("shop created", "id", shop.ID, "name", shop.Name, "url", shop.ShopURL)
	writeJSON(w, http.StatusCreated, ShopResponse{
		ID: shop.ID, Name: shop.Name, ShopURL: shop.ShopURL, APIKeyID: shop.APIKeyID, IsDefault: shop.IsDefault, IsReadOnly: shop.IsReadOnly, PDNEnabled: shop.PDNEnabled,
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

// Code generated from tools.go — keep in sync (A2 drift test).
var readCommands = []string{
	"list_products",
	"get_product",
	"list_orders",
	"get_order",
	"list_clients",
	"get_client",
	"list_categories",
	"get_category",
	"list_collections",
	"get_collection",
	"list_images",
	"list_meta_fields",
	"list_webhooks",
	"get_shop_info",
	"list_currencies",
	"list_delivery_variants",
	"get_delivery_variant",
	"list_payment_gateways",
	"get_payment_gateway",
	"list_pages",
	"get_page",
	"list_accessories",
	"get_accessory",
	"list_variants",
	"get_variant",
	"list_property_values",
	"list_menus",
	"get_menu",
	"list_menu_items",
	"get_menu_item",
	"list_redirects",
	"list_blogs",
	"get_blog",
	"list_blog_posts",
	"get_blog_post",
	"list_discount_codes",
	"get_discount_code",
	"list_collects",
	"get_collects_by_product",
	"list_properties",
	"get_property",
	"list_characteristics",
	"get_characteristic",
	"list_reviews",
	"get_review",
	"list_files",
	"get_file",
	"list_domains",
	"get_domain",
	"list_client_groups",
	"get_client_group",
	"list_custom_statuses",
	"get_custom_status",
	"list_product_fields",
	"get_product_field",
	"list_product_field_values",
	"get_product_field_value",
	"list_product_accessory_links",
	"list_variant_field_values",
	"get_variant_field_value",
	"list_collection_field_values",
	"get_collection_field_value",
	"list_variant_fields",
	"get_variant_field",
	"list_collection_fields",
	"get_collection_field",
	"list_collection_filters",
	"get_collection_filter",
	"list_js_tags",
	"get_js_tag",
	"list_external_discounts",
	"get_external_discount",
	"list_discount_collections_groups",
	"get_discount_collections_group",
	"list_warehouses",
	"get_warehouse",
	"list_sales_channels",
	"get_sales_channel",
	"list_stock_currencies",
	"get_stock_currency",
	"list_price_kinds",
	"get_price_kind",
	"list_option_names",
	"get_option_name",
	"list_option_values",
	"get_option_value",
	"list_application_actions",
	"get_application_action",
	"list_application_charges",
	"get_application_charge",
	"list_application_widgets",
	"get_application_widget",
	"list_recurring_application_charges",
	"get_recurring_application_charge",
	"list_bonus_system_transactions",
	"get_bonus_system_transaction",
	"list_client_bonus_transactions",
	"list_oauth_apps",
	"get_oauth_app",
	"list_user_tasks",
	"get_user_task",
	"list_supplementaries",
	"get_supplementary",
	"list_similars",
	"list_units",
	"list_outlets",
	"get_outlet",
	"list_marketplaces",
	"get_marketplace",
	"get_orders_count",
	"get_products_count",
	"search_product_details",
	"order_full_details",
	"client_full_profile",
	"product_catalog_search",
	"collection_products",
	"export_products",
	"export_orders",
	"export_clients",
	"search_insales_tools",
}

var allCommandNames = []string{
	"list_products",
	"get_product",
	"create_product",
	"update_product",
	"delete_product",
	"list_orders",
	"get_order",
	"update_order",
	"create_order",
	"add_order_line",
	"update_order_line",
	"remove_order_line",
	"add_order_discount",
	"update_order_discount",
	"remove_order_discount",
	"update_order_shipping_address",
	"update_order_custom_status",
	"list_clients",
	"get_client",
	"create_client",
	"update_client",
	"delete_client",
	"list_categories",
	"get_category",
	"create_category",
	"update_category",
	"delete_category",
	"list_collections",
	"get_collection",
	"create_collection",
	"update_collection",
	"delete_collection",
	"list_images",
	"add_image",
	"delete_image",
	"list_meta_fields",
	"create_meta_field",
	"update_meta_field",
	"delete_meta_field",
	"list_webhooks",
	"create_webhook",
	"update_webhook",
	"delete_webhook",
	"get_shop_info",
	"list_currencies",
	"list_delivery_variants",
	"get_delivery_variant",
	"create_delivery_variant",
	"update_delivery_variant",
	"delete_delivery_variant",
	"list_payment_gateways",
	"get_payment_gateway",
	"create_payment_gateway",
	"update_payment_gateway",
	"delete_payment_gateway",
	"list_pages",
	"get_page",
	"create_page",
	"update_page",
	"delete_page",
	"list_accessories",
	"get_accessory",
	"create_accessory",
	"update_accessory",
	"delete_accessory",
	"list_variants",
	"get_variant",
	"update_variant",
	"delete_variant",
	"list_property_values",
	"create_property_value",
	"update_property_value",
	"delete_property_value",
	"list_menus",
	"get_menu",
	"create_menu",
	"update_menu",
	"delete_menu",
	"list_menu_items",
	"get_menu_item",
	"create_menu_item",
	"update_menu_item",
	"delete_menu_item",
	"list_redirects",
	"create_redirect",
	"delete_redirect",
	"list_blogs",
	"get_blog",
	"create_blog",
	"update_blog",
	"delete_blog",
	"list_blog_posts",
	"get_blog_post",
	"create_blog_post",
	"update_blog_post",
	"delete_blog_post",
	"list_discount_codes",
	"get_discount_code",
	"create_discount_code",
	"update_discount_code",
	"delete_discount_code",
	"list_collects",
	"get_collects_by_product",
	"add_product_to_collection",
	"remove_product_from_collection",
	"update_collect",
	"list_properties",
	"get_property",
	"create_property",
	"update_property",
	"delete_property",
	"list_characteristics",
	"get_characteristic",
	"create_characteristic",
	"update_characteristic",
	"delete_characteristic",
	"list_reviews",
	"get_review",
	"create_review",
	"update_review",
	"delete_review",
	"list_files",
	"get_file",
	"create_file",
	"delete_file",
	"list_domains",
	"get_domain",
	"create_domain",
	"update_domain",
	"delete_domain",
	"list_client_groups",
	"get_client_group",
	"create_client_group",
	"update_client_group",
	"delete_client_group",
	"list_custom_statuses",
	"get_custom_status",
	"create_custom_status",
	"update_custom_status",
	"delete_custom_status",
	"list_product_fields",
	"get_product_field",
	"create_product_field",
	"update_product_field",
	"delete_product_field",
	"list_product_field_values",
	"get_product_field_value",
	"create_product_field_value",
	"update_product_field_value",
	"delete_product_field_value",
	"list_product_accessory_links",
	"list_variant_field_values",
	"get_variant_field_value",
	"create_variant_field_value",
	"update_variant_field_value",
	"delete_variant_field_value",
	"list_collection_field_values",
	"get_collection_field_value",
	"update_collection_field_value",
	"list_variant_fields",
	"get_variant_field",
	"create_variant_field",
	"update_variant_field",
	"delete_variant_field",
	"list_collection_fields",
	"get_collection_field",
	"create_collection_field",
	"update_collection_field",
	"delete_collection_field",
	"list_collection_filters",
	"get_collection_filter",
	"create_collection_filter",
	"update_collection_filter",
	"delete_collection_filter",
	"list_js_tags",
	"get_js_tag",
	"create_js_tag",
	"update_js_tag",
	"delete_js_tag",
	"list_external_discounts",
	"get_external_discount",
	"create_external_discount",
	"update_external_discount",
	"delete_external_discount",
	"list_discount_collections_groups",
	"get_discount_collections_group",
	"create_discount_collections_group",
	"update_discount_collections_group",
	"delete_discount_collections_group",
	"list_warehouses",
	"get_warehouse",
	"create_warehouse",
	"update_warehouse",
	"delete_warehouse",
	"list_sales_channels",
	"get_sales_channel",
	"create_sales_channel",
	"update_sales_channel",
	"delete_sales_channel",
	"list_stock_currencies",
	"get_stock_currency",
	"create_stock_currency",
	"update_stock_currency",
	"delete_stock_currency",
	"list_price_kinds",
	"get_price_kind",
	"create_price_kind",
	"update_price_kind",
	"delete_price_kind",
	"list_option_names",
	"get_option_name",
	"create_option_name",
	"update_option_name",
	"delete_option_name",
	"list_option_values",
	"get_option_value",
	"create_option_value",
	"update_option_value",
	"delete_option_value",
	"list_application_actions",
	"get_application_action",
	"create_application_action",
	"update_application_action",
	"delete_application_action",
	"list_application_charges",
	"get_application_charge",
	"create_application_charge",
	"decline_application_charge",
	"list_application_widgets",
	"get_application_widget",
	"create_application_widget",
	"update_application_widget",
	"delete_application_widget",
	"list_recurring_application_charges",
	"get_recurring_application_charge",
	"create_recurring_application_charge",
	"update_recurring_application_charge",
	"delete_recurring_application_charge",
	"list_bonus_system_transactions",
	"get_bonus_system_transaction",
	"create_bonus_system_transaction",
	"list_client_bonus_transactions",
	"list_oauth_apps",
	"get_oauth_app",
	"create_oauth_app",
	"update_oauth_app",
	"delete_oauth_app",
	"list_user_tasks",
	"get_user_task",
	"create_user_task",
	"update_user_task",
	"delete_user_task",
	"list_supplementaries",
	"get_supplementary",
	"create_supplementary",
	"delete_supplementary",
	"list_similars",
	"create_similar",
	"delete_similar",
	"list_units",
	"list_outlets",
	"get_outlet",
	"create_outlet",
	"update_outlet",
	"delete_outlet",
	"list_marketplaces",
	"get_marketplace",
	"create_marketplace",
	"update_marketplace",
	"delete_marketplace",
	"create_provider_token",
	"update_account",
	"get_orders_count",
	"get_products_count",
	"search_product_details",
	"order_full_details",
	"client_full_profile",
	"product_catalog_search",
	"collection_products",
	"export_products",
	"export_orders",
	"export_clients",
	"search_insales_tools",
	"call_insales_tool",
}

var commandDescriptions = map[string]string{
	"list_products":                       "Список товаров",
	"get_product":                         "Товар по ID",
	"create_product":                      "Создать новый товар в магазине",
	"update_product":                      "Обновить данные товара",
	"delete_product":                      "Удалить товар из магазина",
	"list_orders":                         "Получить список заказов с фильтрацией",
	"get_order":                           "Получить детали заказа по ID",
	"update_order":                        "Обновить заказ: статус, комментарий, строки, скидки, адрес доставки. status — alias для financial_status",
	"create_order":                        "Создать заказ: клиент (client_id ИЛИ client name/email/phone) и опционально строки/адрес/оплата",
	"add_order_line":                      "Добавить строку заказа",
	"update_order_line":                   "Обновить строку заказа",
	"remove_order_line":                   "Удалить строку заказа",
	"add_order_discount":                  "Добавить скидку к заказу",
	"update_order_discount":               "Обновить скидку заказа",
	"remove_order_discount":               "Удалить скидку заказа",
	"update_order_shipping_address":       "Обновить адрес доставки",
	"update_order_custom_status":          "Обновить кастомный статус заказа",
	"list_clients":                        "Получить список клиентов с поиском",
	"get_client":                          "Получить информацию о клиенте по ID",
	"create_client":                       "Создать нового клиента",
	"update_client":                       "Обновить данные клиента",
	"delete_client":                       "Удалить клиента",
	"list_categories":                     "Получить список категорий товаров",
	"get_category":                        "Получить информацию о категории по ID",
	"create_category":                     "Создать новую категорию",
	"update_category":                     "Обновить категорию",
	"delete_category":                     "Удалить категорию",
	"list_collections":                    "Получить список коллекций (умных подборок)",
	"get_collection":                      "Получить информацию о коллекции по ID",
	"create_collection":                   "Создать новую коллекцию",
	"update_collection":                   "Обновить коллекцию",
	"delete_collection":                   "Удалить коллекцию",
	"list_images":                         "Получить список изображений товара",
	"add_image":                           "Добавить изображение к товару",
	"delete_image":                        "Удалить изображение товара",
	"list_meta_fields":                    "Получить список метаполей",
	"create_meta_field":                   "Создать метаполе для ресурса",
	"update_meta_field":                   "Обновить метаполе",
	"delete_meta_field":                   "Удалить метаполе",
	"list_webhooks":                       "Получить список подписанных вебхуков",
	"create_webhook":                      "Подписаться на событие InSales",
	"update_webhook":                      "Обновить URL вебхука",
	"delete_webhook":                      "Отписаться от вебхука",
	"get_shop_info":                       "Получить информацию о магазине (название, домен, валюта)",
	"list_currencies":                     "Получить список валют магазина",
	"list_delivery_variants":              "Способы доставки",
	"get_delivery_variant":                "Способ доставки",
	"create_delivery_variant":             "Создать способ доставки",
	"update_delivery_variant":             "Обновить способ доставки",
	"delete_delivery_variant":             "Удалить способ доставки",
	"list_payment_gateways":               "Способы оплаты",
	"get_payment_gateway":                 "Способ оплаты",
	"create_payment_gateway":              "Создать способ оплаты",
	"update_payment_gateway":              "Обновить способ оплаты",
	"delete_payment_gateway":              "Удалить способ оплаты",
	"list_pages":                          "Получить список страниц магазина",
	"get_page":                            "Получить страницу по ID",
	"create_page":                         "Создать новую страницу в магазине",
	"update_page":                         "Обновить данные страницы",
	"delete_page":                         "Удалить страницу",
	"list_accessories":                    "Получить список аксессуаров (характеристик товаров)",
	"get_accessory":                       "Получить аксессуар по ID",
	"create_accessory":                    "Создать новый аксессуар (характеристику)",
	"update_accessory":                    "Обновить аксессуар",
	"delete_accessory":                    "Удалить аксессуар",
	"list_variants":                       "Получить список вариантов товара (размеры, цвета)",
	"get_variant":                         "Получить информацию о варианте товара",
	"update_variant":                      "Обновить данные варианта товара",
	"delete_variant":                      "Удалить вариант товара",
	"list_property_values":                "Получить список значений аксессуара",
	"create_property_value":               "Создать новое значение для аксессуара",
	"update_property_value":               "Обновить значение аксессуара",
	"delete_property_value":               "Удалить значение аксессуара",
	"list_menus":                          "Получить список навигационных меню",
	"get_menu":                            "Получить меню по ID со списком пунктов",
	"create_menu":                         "Создать новое навигационное меню",
	"update_menu":                         "Обновить данные меню",
	"delete_menu":                         "Удалить меню",
	"list_menu_items":                     "Пункты меню",
	"get_menu_item":                       "Пункт меню",
	"create_menu_item":                    "Добавить пункт меню",
	"update_menu_item":                    "Обновить пункт меню",
	"delete_menu_item":                    "Удалить пункт меню",
	"list_redirects":                      "Получить список URL-редиректов",
	"create_redirect":                     "Создать URL-редирект",
	"delete_redirect":                     "Удалить URL-редирект",
	"list_blogs":                          "Получить список блогов магазина",
	"get_blog":                            "Получить информацию о блоге по ID",
	"create_blog":                         "Создать новый блог в магазине",
	"update_blog":                         "Обновить блог",
	"delete_blog":                         "Удалить блог",
	"list_blog_posts":                     "Получить список записей в блоге",
	"get_blog_post":                       "Получить запись блога по ID",
	"create_blog_post":                    "Создать новую запись в блоге",
	"update_blog_post":                    "Обновить запись блога",
	"delete_blog_post":                    "Удалить запись блога",
	"list_discount_codes":                 "Получить список промокодов",
	"get_discount_code":                   "Получить промокод по ID",
	"create_discount_code":                "Создать промокод",
	"update_discount_code":                "Обновить промокод",
	"delete_discount_code":                "Удалить промокод",
	"list_collects":                       "Получить список связей товар-коллекция",
	"get_collects_by_product":             "Получить коллекции товара",
	"add_product_to_collection":           "Добавить товар в коллекцию",
	"remove_product_from_collection":      "Удалить товар из коллекции",
	"update_collect":                      "Обновить связь товар-коллекция",
	"list_properties":                     "Получить список свойств товаров (Бренд, Цвет...)",
	"get_property":                        "Получить свойство по ID",
	"create_property":                     "Создать свойство товаров",
	"update_property":                     "Обновить свойство товаров",
	"delete_property":                     "Удалить свойство товаров",
	"list_characteristics":                "Получить список характеристик свойства (Nike, Adidas...)",
	"get_characteristic":                  "Получить характеристику по ID",
	"create_characteristic":               "Создать характеристику свойства",
	"update_characteristic":               "Обновить характеристику свойства",
	"delete_characteristic":               "Удалить характеристику свойства",
	"list_reviews":                        "Получить список отзывов",
	"get_review":                          "Получить отзыв по ID",
	"create_review":                       "Создать отзыв",
	"update_review":                       "Обновить отзыв",
	"delete_review":                       "Удалить отзыв",
	"list_files":                          "Получить список файлов магазина",
	"get_file":                            "Получить файл по ID",
	"create_file":                         "Загрузить файл (по URL)",
	"delete_file":                         "Удалить файл",
	"list_domains":                        "Список доменов магазина",
	"get_domain":                          "Получить домен по ID",
	"create_domain":                       "Добавить домен",
	"update_domain":                       "Обновить домен",
	"delete_domain":                       "Удалить домен",
	"list_client_groups":                  "Список групп клиентов",
	"get_client_group":                    "Получить группу клиентов",
	"create_client_group":                 "Создать группу клиентов",
	"update_client_group":                 "Обновить группу клиентов",
	"delete_client_group":                 "Удалить группу клиентов",
	"list_custom_statuses":                "Список кастомных статусов",
	"get_custom_status":                   "Получить кастомный статус",
	"create_custom_status":                "Создать кастомный статус",
	"update_custom_status":                "Обновить кастомный статус",
	"delete_custom_status":                "Удалить кастомный статус",
	"list_product_fields":                 "Список полей товаров",
	"get_product_field":                   "Поле товара",
	"create_product_field":                "Создать поле товара",
	"update_product_field":                "Обновить поле товара",
	"delete_product_field":                "Удалить поле товара",
	"list_product_field_values":           "Значения полей товара",
	"get_product_field_value":             "Значение поля товара",
	"create_product_field_value":          "Установить значение поля товара",
	"update_product_field_value":          "Обновить значение поля товара",
	"delete_product_field_value":          "Удалить значение поля товара",
	"list_product_accessory_links":        "Связи аксессуаров товара",
	"list_variant_field_values":           "Значения полей варианта",
	"get_variant_field_value":             "Значение поля варианта",
	"create_variant_field_value":          "Установить значение поля варианта",
	"update_variant_field_value":          "Обновить значение поля варианта",
	"delete_variant_field_value":          "Удалить значение поля варианта",
	"list_collection_field_values":        "Значения полей коллекции",
	"get_collection_field_value":          "Значение поля коллекции",
	"update_collection_field_value":       "Обновить значение поля коллекции",
	"list_variant_fields":                 "Список полей вариантов",
	"get_variant_field":                   "Получить поле варианта",
	"create_variant_field":                "Создать поле варианта",
	"update_variant_field":                "Обновить поле варианта",
	"delete_variant_field":                "Удалить поле варианта",
	"list_collection_fields":              "Список полей коллекций",
	"get_collection_field":                "Получить поле коллекции",
	"create_collection_field":             "Создать поле коллекции",
	"update_collection_field":             "Обновить поле коллекции",
	"delete_collection_field":             "Удалить поле коллекции",
	"list_collection_filters":             "Список фильтров коллекций",
	"get_collection_filter":               "Получить фильтр коллекции",
	"create_collection_filter":            "Создать фильтр коллекции",
	"update_collection_filter":            "Обновить фильтр коллекции",
	"delete_collection_filter":            "Удалить фильтр коллекции",
	"list_js_tags":                        "Список JS-тегов",
	"get_js_tag":                          "Получить JS-тег",
	"create_js_tag":                       "Создать JS-тег",
	"update_js_tag":                       "Обновить JS-тег",
	"delete_js_tag":                       "Удалить JS-тег",
	"list_external_discounts":             "Список внешних скидок",
	"get_external_discount":               "Получить внешнюю скидку",
	"create_external_discount":            "Создать внешнюю скидку",
	"update_external_discount":            "Обновить внешнюю скидку",
	"delete_external_discount":            "Удалить внешнюю скидку",
	"list_discount_collections_groups":    "Список групп коллекций скидок",
	"get_discount_collections_group":      "Получить группу коллекций скидок",
	"create_discount_collections_group":   "Создать группу коллекций скидок",
	"update_discount_collections_group":   "Обновить группу коллекций скидок",
	"delete_discount_collections_group":   "Удалить группу коллекций скидок",
	"list_warehouses":                     "Список складов",
	"get_warehouse":                       "Получить склад",
	"create_warehouse":                    "Создать склад",
	"update_warehouse":                    "Обновить склад",
	"delete_warehouse":                    "Удалить склад",
	"list_sales_channels":                 "Список каналов продаж",
	"get_sales_channel":                   "Получить канал продаж",
	"create_sales_channel":                "Создать канал продаж",
	"update_sales_channel":                "Обновить канал продаж",
	"delete_sales_channel":                "Удалить канал продаж",
	"list_stock_currencies":               "Список валют склада",
	"get_stock_currency":                  "Получить валюту склада",
	"create_stock_currency":               "Создать валюту склада",
	"update_stock_currency":               "Обновить валюту склада",
	"delete_stock_currency":               "Удалить валюту склада",
	"list_price_kinds":                    "Список типов цен",
	"get_price_kind":                      "Получить тип цены",
	"create_price_kind":                   "Создать тип цены",
	"update_price_kind":                   "Обновить тип цены",
	"delete_price_kind":                   "Удалить тип цены",
	"list_option_names":                   "Список имён опций (Размер, Цвет...)",
	"get_option_name":                     "Получить имя опции",
	"create_option_name":                  "Создать имя опции",
	"update_option_name":                  "Обновить имя опции",
	"delete_option_name":                  "Удалить имя опции",
	"list_option_values":                  "Список значений опции",
	"get_option_value":                    "Получить значение опции",
	"create_option_value":                 "Создать значение опции",
	"update_option_value":                 "Обновить значение опции",
	"delete_option_value":                 "Удалить значение опции",
	"list_application_actions":            "Список действий приложения",
	"get_application_action":              "Получить действие приложения",
	"create_application_action":           "Создать действие приложения",
	"update_application_action":           "Обновить действие приложения",
	"delete_application_action":           "Удалить действие приложения",
	"list_application_charges":            "Список счетов приложения",
	"get_application_charge":              "Получить счёт приложения",
	"create_application_charge":           "Создать счёт приложения",
	"decline_application_charge":          "Отклонить счёт приложения",
	"list_application_widgets":            "Список виджетов приложения",
	"get_application_widget":              "Получить виджет",
	"create_application_widget":           "Создать виджет",
	"update_application_widget":           "Обновить виджет",
	"delete_application_widget":           "Удалить виджет",
	"list_recurring_application_charges":  "Список рекуррентных счетов",
	"get_recurring_application_charge":    "Получить рекуррентный счёт",
	"create_recurring_application_charge": "Создать рекуррентный счёт",
	"update_recurring_application_charge": "Обновить рекуррентный счёт",
	"delete_recurring_application_charge": "Удалить рекуррентный счёт",
	"list_bonus_system_transactions":      "Список бонусных транзакций",
	"get_bonus_system_transaction":        "Получить бонусную транзакцию",
	"create_bonus_system_transaction":     "Создать бонусную транзакцию",
	"list_client_bonus_transactions":      "Список бонусных транзакций клиента",
	"list_oauth_apps":                     "Список OAuth приложений",
	"get_oauth_app":                       "Получить OAuth приложение",
	"create_oauth_app":                    "Создать OAuth приложение",
	"update_oauth_app":                    "Обновить OAuth приложение",
	"delete_oauth_app":                    "Удалить OAuth приложение",
	"list_user_tasks":                     "Список задач пользователя",
	"get_user_task":                       "Получить задачу",
	"create_user_task":                    "Создать задачу",
	"update_user_task":                    "Обновить задачу",
	"delete_user_task":                    "Удалить задачу",
	"list_supplementaries":                "Список дополнительных товаров",
	"get_supplementary":                   "Получить дополнительный товар",
	"create_supplementary":                "Создать дополнительный товар",
	"delete_supplementary":                "Удалить дополнительный товар",
	"list_similars":                       "Похожие товары",
	"create_similar":                      "Добавить похожий товар",
	"delete_similar":                      "Удалить похожий товар",
	"list_units":                          "Список единиц измерения",
	"list_outlets":                        "Список точек выдачи",
	"get_outlet":                          "Получить точку выдачи по ID",
	"create_outlet":                       "Создать точку выдачи",
	"update_outlet":                       "Обновить точку выдачи",
	"delete_outlet":                       "Удалить точку выдачи",
	"list_marketplaces":                   "Список маркетплейсов",
	"get_marketplace":                     "Получить маркетплейс по ID",
	"create_marketplace":                  "Создать маркетплейс. Типы: Marketplace::GoogleMerchant, Marketplace::LiquidYandexMarket, Marketplace::ModelYandexMarket",
	"update_marketplace":                  "Обновить маркетплейс",
	"delete_marketplace":                  "Удалить маркетплейс",
	"create_provider_token":               "Создать токен провайдера мессенджера",
	"update_account":                      "Обновить настройки аккаунта",
	"get_orders_count":                    "Количество заказов",
	"get_products_count":                  "Количество товаров",
	"search_product_details":              "Поиск товара с полной информацией: кастомные поля, метаполя, изображения. Возвращает полную карточку товара одним вызовом.",
	"order_full_details":                  "Полная информация о заказе: позиции, адрес доставки, клиент, кастомные поля. Всё одним вызовом.",
	"client_full_profile":                 "Полный профиль клиента: данные, последние заказы, бонусы, метаполя. Всё одним вызовом.",
	"product_catalog_search":              "Поиск по каталогу товаров с фильтрами, пагинацией и подсчётом.",
	"collection_products":                 "Товары коллекции с полной информацией: данные коллекции + все её товары.",
	"export_products":                     "Экспорт товаров в CSV. Возвращает CSV-файл с колонками: id, title, price, sku, vendor, quantity, created_at",
	"export_orders":                       "Экспорт заказов в CSV. Возвращает CSV-файл с колонками: id, number, status, fulfillment_status, total_price, client_name, created_at",
	"export_clients":                      "Экспорт клиентов в CSV. Возвращает CSV-файл с колонками: id, name, surname, email, phone, created_at",
	"search_insales_tools":                "Поиск инструментов InSales",
	"call_insales_tool":                   "Вызов инструмента InSales",
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

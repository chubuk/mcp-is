package mcp

// tool_domains.go — маппинг инструментов → {probe-домен, функциональная группа, класс операции}.
// Класс задан явно для каждой команды (в рантайме ничего не выводится);
// дефолт для незамапленного имени — самый строгий класс («удаление»).
// Дрейф-тест TestToolDomainCoverage сверяет полноту против tools.go.

// Классы операций.
const (
	ClassRead   = "read"   // list_/get_/search_*, *_full_details, *_count, *_profile, export_*
	ClassWrite  = "write"  // create_/add_/update_*
	ClassDelete = "delete" // delete_/remove_/decline_*
)

// ToolMeta — метаданные инструмента для фильтрации и UI.
type ToolMeta struct {
	Domain string // probe-домен InSales ("" = discovery, без домена)
	Group  string // функциональная группа (редактор политики)
	Class  string // read | write | delete
}

// toolMeta — полный маппинг инструментов.
var toolMeta = buildToolMeta()

func buildToolMeta() map[string]ToolMeta {
	m := make(map[string]ToolMeta, 300)

	// ── Discovery / служебные (без домена, всегда видны) ──
	register(m, "", "service", ClassRead,
		"search_insales_tools", "call_insales_tool", "store_info")

	// ── Каталог товаров (domain: products) ──
	register(m, "products", "catalog", ClassRead,
		"list_products", "get_product", "get_products_count",
		"list_variants", "get_variant",
		"list_images",
		"list_collects", "get_collects_by_product",
		"list_accessories", "get_accessory", "list_product_accessory_links",
		"list_supplementaries", "get_supplementary",
		"list_similars",
		"list_meta_fields",
		"list_product_fields", "get_product_field", "list_product_field_values", "get_product_field_value",
		"list_variant_fields", "get_variant_field", "list_variant_field_values", "get_variant_field_value",
		"get_seo",
		"search_product_details", "product_catalog_search",
		"export_products",
	)
	register(m, "products", "catalog", ClassWrite,
		"create_product", "update_product",
		"update_variant",
		"add_image",
		"update_collect",
		"create_accessory", "update_accessory",
		"create_supplementary",
		"create_similar",
		"create_meta_field", "update_meta_field",
		"create_product_field", "update_product_field",
		"create_product_field_value", "update_product_field_value",
		"create_variant_field", "update_variant_field",
		"create_variant_field_value", "update_variant_field_value",
		"update_seo",
	)
	register(m, "products", "catalog", ClassDelete,
		"delete_product",
		"delete_variant",
		"delete_image",
		"delete_accessory",
		"delete_supplementary",
		"delete_similar",
		"delete_meta_field",
		"delete_product_field",
		"delete_product_field_value",
		"delete_variant_field",
		"delete_variant_field_value",
	)

	// ── Свойства и характеристики (domain: properties) ──
	register(m, "properties", "properties", ClassRead,
		"list_properties", "get_property",
		"list_property_values",
		"list_characteristics", "get_characteristic",
		"list_option_names", "get_option_name",
		"list_option_values", "get_option_value",
		"list_units",
	)
	register(m, "properties", "properties", ClassWrite,
		"create_property", "update_property",
		"create_property_value", "update_property_value",
		"create_characteristic", "update_characteristic",
		"create_option_name", "update_option_name",
		"create_option_value", "update_option_value",
	)
	register(m, "properties", "properties", ClassDelete,
		"delete_property",
		"delete_property_value",
		"delete_characteristic",
		"delete_option_name",
		"delete_option_value",
	)

	// ── Категории и коллекции (domain: categories, collections) ──
	register(m, "categories", "categories_collections", ClassRead,
		"list_categories", "get_category",
	)
	register(m, "categories", "categories_collections", ClassWrite,
		"create_category", "update_category",
	)
	register(m, "categories", "categories_collections", ClassDelete,
		"delete_category",
	)
	register(m, "collections", "categories_collections", ClassRead,
		"list_collections", "get_collection",
		"list_collection_fields", "get_collection_field",
		"list_collection_field_values", "get_collection_field_value",
		"list_collection_filters", "get_collection_filter",
		"collection_products",
	)
	register(m, "collections", "categories_collections", ClassWrite,
		"create_collection", "update_collection",
		"create_collection_field", "update_collection_field",
		"update_collection_field_value",
		"create_collection_filter", "update_collection_filter",
		"add_product_to_collection",
	)
	register(m, "collections", "categories_collections", ClassDelete,
		"delete_collection",
		"delete_collection_field",
		"delete_collection_filter",
		"remove_product_from_collection",
	)

	// ── Заказы (domain: orders) ──
	register(m, "orders", "orders", ClassRead,
		"list_orders", "get_order", "get_orders_count",
		"list_custom_statuses", "get_custom_status",
		"order_full_details",
		"export_orders",
	)
	register(m, "orders", "orders", ClassWrite,
		"create_order", "update_order",
		"add_order_line", "update_order_line",
		"add_order_discount", "update_order_discount",
		"update_order_custom_status",
		"update_order_shipping_address",
		"create_custom_status", "update_custom_status",
	)
	register(m, "orders", "orders", ClassDelete,
		"remove_order_line",
		"remove_order_discount",
		"delete_custom_status",
	)

	// ── Клиенты (domain: clients) ──
	register(m, "clients", "clients", ClassRead,
		"list_clients", "get_client",
		"list_client_groups", "get_client_group",
		"list_bonus_system_transactions", "get_bonus_system_transaction",
		"list_client_bonus_transactions",
		"client_full_profile",
		"export_clients",
	)
	register(m, "clients", "clients", ClassWrite,
		"create_client", "update_client",
		"create_client_group", "update_client_group",
		"create_bonus_system_transaction",
	)
	register(m, "clients", "clients", ClassDelete,
		"delete_client",
		"delete_client_group",
	)

	// ── Блог (domain: blogs) ──
	register(m, "blogs", "blog", ClassRead,
		"list_blogs", "get_blog",
		"list_blog_posts", "get_blog_post",
	)
	register(m, "blogs", "blog", ClassWrite,
		"create_blog", "update_blog",
		"create_blog_post", "update_blog_post",
	)
	register(m, "blogs", "blog", ClassDelete,
		"delete_blog",
		"delete_blog_post",
	)

	// ── Страницы (domain: pages) ──
	register(m, "pages", "pages", ClassRead,
		"list_pages", "get_page",
	)
	register(m, "pages", "pages", ClassWrite,
		"create_page", "update_page",
	)
	register(m, "pages", "pages", ClassDelete,
		"delete_page",
	)

	// ── Меню и навигация (domain: menus, redirects) ──
	register(m, "menus", "menus_nav", ClassRead,
		"list_menus", "get_menu",
		"list_menu_items", "get_menu_item",
	)
	register(m, "menus", "menus_nav", ClassWrite,
		"create_menu", "update_menu",
		"create_menu_item", "update_menu_item",
	)
	register(m, "menus", "menus_nav", ClassDelete,
		"delete_menu",
		"delete_menu_item",
	)
	register(m, "redirects", "menus_nav", ClassRead,
		"list_redirects",
	)
	register(m, "redirects", "menus_nav", ClassWrite,
		"create_redirect",
	)
	register(m, "redirects", "menus_nav", ClassDelete,
		"delete_redirect",
	)

	// ── Отзывы (domain: reviews) ──
	register(m, "reviews", "reviews", ClassRead,
		"list_reviews", "get_review",
	)
	register(m, "reviews", "reviews", ClassWrite,
		"create_review", "update_review",
	)
	register(m, "reviews", "reviews", ClassDelete,
		"delete_review",
	)

	// ── Скидки и акции (domain: discounts) ──
	register(m, "discounts", "discounts", ClassRead,
		"list_discount_codes", "get_discount_code",
		"list_discount_collections_groups", "get_discount_collections_group",
		"list_external_discounts", "get_external_discount",
	)
	register(m, "discounts", "discounts", ClassWrite,
		"create_discount_code", "update_discount_code",
		"create_discount_collections_group", "update_discount_collections_group",
		"create_external_discount", "update_external_discount",
	)
	register(m, "discounts", "discounts", ClassDelete,
		"delete_discount_code",
		"delete_discount_collections_group",
		"delete_external_discount",
	)

	// ── Доставка и оплата (domain: delivery, payments) ──
	register(m, "delivery", "delivery_payments", ClassRead,
		"list_delivery_variants", "get_delivery_variant",
	)
	register(m, "delivery", "delivery_payments", ClassWrite,
		"create_delivery_variant", "update_delivery_variant",
	)
	register(m, "delivery", "delivery_payments", ClassDelete,
		"delete_delivery_variant",
	)
	register(m, "payments", "delivery_payments", ClassRead,
		"list_payment_gateways", "get_payment_gateway",
	)
	register(m, "payments", "delivery_payments", ClassWrite,
		"create_payment_gateway", "update_payment_gateway",
	)
	register(m, "payments", "delivery_payments", ClassDelete,
		"delete_payment_gateway",
	)

	// ── Склад и цены (domain: warehouse, outlets) ──
	register(m, "warehouse", "warehouse_prices", ClassRead,
		"list_warehouses", "get_warehouse",
		"list_stock_currencies", "get_stock_currency",
		"list_price_kinds", "get_price_kind",
		"list_sales_channels", "get_sales_channel",
	)
	register(m, "warehouse", "warehouse_prices", ClassWrite,
		"create_warehouse", "update_warehouse",
		"create_stock_currency", "update_stock_currency",
		"create_price_kind", "update_price_kind",
		"create_sales_channel", "update_sales_channel",
	)
	register(m, "warehouse", "warehouse_prices", ClassDelete,
		"delete_warehouse",
		"delete_stock_currency",
		"delete_price_kind",
		"delete_sales_channel",
	)
	register(m, "outlets", "warehouse_prices", ClassRead,
		"list_outlets", "get_outlet",
	)
	register(m, "outlets", "warehouse_prices", ClassWrite,
		"create_outlet", "update_outlet",
	)
	register(m, "outlets", "warehouse_prices", ClassDelete,
		"delete_outlet",
	)

	// ── Файлы (domain: files) ──
	register(m, "files", "files", ClassRead,
		"list_files", "get_file",
	)
	register(m, "files", "files", ClassWrite,
		"create_file",
	)
	register(m, "files", "files", ClassDelete,
		"delete_file",
	)

	// ── Домены (domain: domains) ──
	register(m, "domains", "domains", ClassRead,
		"list_domains", "get_domain",
	)
	register(m, "domains", "domains", ClassWrite,
		"create_domain", "update_domain",
	)
	register(m, "domains", "domains", ClassDelete,
		"delete_domain",
	)

	// ── Настройки магазина (domain: account) ──
	register(m, "account", "shop_settings", ClassRead,
		"get_shop_info", "list_currencies",
	)
	register(m, "account", "shop_settings", ClassWrite,
		"update_account",
	)

	// ── Вебхуки и интеграции (domain: webhooks, js_tags, marketplaces, apps) ──
	register(m, "webhooks", "webhooks_integrations", ClassRead,
		"list_webhooks",
	)
	register(m, "webhooks", "webhooks_integrations", ClassWrite,
		"create_webhook", "update_webhook",
	)
	register(m, "webhooks", "webhooks_integrations", ClassDelete,
		"delete_webhook",
	)
	register(m, "js_tags", "webhooks_integrations", ClassRead,
		"list_js_tags", "get_js_tag",
	)
	register(m, "js_tags", "webhooks_integrations", ClassWrite,
		"create_js_tag", "update_js_tag",
	)
	register(m, "js_tags", "webhooks_integrations", ClassDelete,
		"delete_js_tag",
	)
	register(m, "marketplaces", "webhooks_integrations", ClassRead,
		"list_marketplaces", "get_marketplace",
	)
	register(m, "marketplaces", "webhooks_integrations", ClassWrite,
		"create_marketplace", "update_marketplace",
	)
	register(m, "marketplaces", "webhooks_integrations", ClassDelete,
		"delete_marketplace",
	)
	register(m, "apps", "webhooks_integrations", ClassRead,
		"list_oauth_apps", "get_oauth_app",
		"list_application_actions", "get_application_action",
		"list_application_widgets", "get_application_widget",
	)
	register(m, "apps", "webhooks_integrations", ClassWrite,
		"create_oauth_app", "update_oauth_app",
		"create_application_action", "update_application_action",
		"create_application_widget", "update_application_widget",
		"create_provider_token",
	)
	register(m, "apps", "webhooks_integrations", ClassDelete,
		"delete_oauth_app",
	)

	// ── Приложения и биллинг (domain: apps) ──
	register(m, "apps", "apps_billing", ClassRead,
		"list_application_charges", "get_application_charge",
		"list_recurring_application_charges", "get_recurring_application_charge",
	)
	register(m, "apps", "apps_billing", ClassWrite,
		"create_application_charge",
		"create_recurring_application_charge", "update_recurring_application_charge",
	)
	register(m, "apps", "apps_billing", ClassDelete,
		"decline_application_charge",
		"delete_recurring_application_charge",
		"delete_application_action",
		"delete_application_widget",
	)

	// ── Задачи (domain: user_tasks) ──
	register(m, "user_tasks", "user_tasks", ClassRead,
		"list_user_tasks", "get_user_task",
	)
	register(m, "user_tasks", "user_tasks", ClassWrite,
		"create_user_task", "update_user_task",
	)
	register(m, "user_tasks", "user_tasks", ClassDelete,
		"delete_user_task",
	)

	return m
}

func register(m map[string]ToolMeta, domain, group, class string, names ...string) {
	for _, n := range names {
		m[n] = ToolMeta{Domain: domain, Group: group, Class: class}
	}
}

// ToolMetaFor возвращает метаданные инструмента.
// Неизвестное имя → самый строгий класс («удаление»), без домена.
func ToolMetaFor(name string) ToolMeta {
	if tm, ok := toolMeta[name]; ok {
		return tm
	}
	return ToolMeta{Domain: "", Group: "unknown", Class: ClassDelete}
}

// ToolDomain возвращает probe-домен инструмента ("" = без домена).
func ToolDomain(name string) string {
	return ToolMetaFor(name).Domain
}

// ToolClass возвращает класс операции инструмента.
func ToolClass(name string) string {
	return ToolMetaFor(name).Class
}

// ToolGroup возвращает функциональную группу инструмента.
func ToolGroup(name string) string {
	return ToolMetaFor(name).Group
}

// AllToolNames — все имена из маппинга (для drift-теста и UI).
func AllToolNames() []string {
	names := make([]string, 0, len(toolMeta))
	for n := range toolMeta {
		names = append(names, n)
	}
	return names
}

// GroupClassEnabled — политика: разрешён ли класс операций в группе.
type GroupClassEnabled struct {
	Read   bool `json:"read"`
	Write  bool `json:"write"`
	Delete bool `json:"delete"`
}

// CommandPolicy — политика ИИ-ключа (JSON в ai_keys.command_policy).
// Приоритет: override команды > класс в группе > default; deny побеждает.
type CommandPolicy struct {
	Default  string                       `json:"default"` // "allow" | "deny"
	Groups   map[string]GroupClassEnabled `json:"groups"`
	Commands map[string]bool              `json:"commands"` // override: true=allow, false=deny
}

// IsCommandAllowed вычисляет, разрешена ли команда политикой.
func (p *CommandPolicy) IsCommandAllowed(name string) bool {
	if p == nil {
		return true // '' / отсутствует = allow (обратная совместимость)
	}
	// 1. Override команды
	if allowed, ok := p.Commands[name]; ok {
		return allowed
	}
	// 2. Класс операции в группе
	meta := ToolMetaFor(name)
	if g, ok := p.Groups[meta.Group]; ok {
		var classOn bool
		switch meta.Class {
		case ClassRead:
			classOn = g.Read
		case ClassWrite:
			classOn = g.Write
		case ClassDelete:
			classOn = g.Delete
		default:
			classOn = false
		}
		return classOn
	}
	// 3. Default
	return p.Default != "deny"
}

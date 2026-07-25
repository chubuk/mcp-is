package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/chubuk/mcp-is/internal/service"
)

type ToolInfo struct {
	Name        string
	Description string
}

var toolRegistry = make(map[string]ToolInfo)

func GetToolNames() []string {
	names := make([]string, 0, len(toolRegistry))
	for name := range toolRegistry {
		names = append(names, name)
	}
	return names
}

func GetToolDescription(name string) string {
	if info, ok := toolRegistry[name]; ok {
		return info.Description
	}
	return ""
}

func GetAllToolInfos() []ToolInfo {
	infos := make([]ToolInfo, 0, len(toolRegistry))
	for _, info := range toolRegistry {
		infos = append(infos, info)
	}
	return infos
}

func RegisterTools(s *server.MCPServer, svc *service.Service) {
	registerTool := func(name, desc string, tool mcp.Tool, handler func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
		s.AddTool(tool, handler)
		toolRegistry[name] = ToolInfo{Name: name, Description: desc}
	}

	// === ТОВАРЫ ===
	registerTool("list_products", "Получить список товаров магазина с пагинацией и фильтрацией",
		mcp.NewTool("list_products",
		mcp.WithDescription("Получить список товаров магазина с пагинацией и фильтрацией"),
		mcp.WithNumber("page", mcp.Description("Номер страницы"), mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.Description("Товаров на странице (макс 250)"), mcp.DefaultNumber(25)),
		mcp.WithString("category_id", mcp.Description("Фильтр по категории")),
		mcp.WithString("collection_id", mcp.Description("Фильтр по коллекции")),
		mcp.WithString("q", mcp.Description("Поиск по названию/артикулу")),
		mcp.WithString("updated_since", mcp.Description("Дата обновления (ISO format)")),
		mcp.WithNumber("from_id", mcp.Description("ID начиная с которого получать данные")),
		mcp.WithBoolean("with_deleted", mcp.Description("Включить удалённые товары")),
		mcp.WithBoolean("deleted", mcp.Description("Только удалённые товары")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listProductsHandler(svc))

	s.AddTool(mcp.NewTool("get_product",
		mcp.WithDescription("Получить детальную информацию о товаре по ID"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getProductHandler(svc))

	s.AddTool(mcp.NewTool("create_product",
		mcp.WithDescription("Создать новый товар в магазине"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название товара")),
		mcp.WithString("price", mcp.Required(), mcp.Description("Цена (например 1000.00)")),
		mcp.WithString("description", mcp.Description("Описание товара")),
		mcp.WithString("vendor", mcp.Description("Производитель")),
		mcp.WithString("sku", mcp.Description("Артикул")),
		mcp.WithNumber("quantity", mcp.Description("Количество на складе")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
	), createProductHandler(svc))

	s.AddTool(mcp.NewTool("update_product",
		mcp.WithDescription("Обновить данные товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithString("price", mcp.Description("Цена")),
		mcp.WithNumber("quantity", mcp.Description("Количество")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateProductHandler(svc))

	s.AddTool(mcp.NewTool("delete_product",
		mcp.WithDescription("Удалить товар из магазина"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteProductHandler(svc))

	// === ЗАКАЗЫ ===
	s.AddTool(mcp.NewTool("list_orders",
		mcp.WithDescription("Получить список заказов с фильтрацией"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithString("status", mcp.Description("Статус: open, closed, canceled")),
		mcp.WithString("fulfillment_status", mcp.Description("Статус выполнения")),
		mcp.WithString("updated_since", mcp.Description("Дата обновления (ISO format)")),
		mcp.WithNumber("from_id", mcp.Description("ID начиная с которого получать данные")),
		mcp.WithBoolean("paid", mcp.Description("Оплачен ли")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listOrdersHandler(svc))

	s.AddTool(mcp.NewTool("get_order",
		mcp.WithDescription("Получить детали заказа по ID"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("ID заказа")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getOrderHandler(svc))

	s.AddTool(mcp.NewTool("update_order",
		mcp.WithDescription("Обновить статус или данные заказа"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("ID заказа")),
		mcp.WithString("status", mcp.Description("Статус")),
		mcp.WithString("fulfillment_status", mcp.Description("Статус выполнения")),
		mcp.WithString("comment", mcp.Description("Комментарий")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateOrderHandler(svc))

	s.AddTool(mcp.NewTool("create_order",
		mcp.WithDescription("Создать новый заказ вручную"),
		mcp.WithNumber("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithString("comment", mcp.Description("Комментарий")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createOrderHandler(svc))

	// === КЛИЕНТЫ ===
	s.AddTool(mcp.NewTool("list_clients",
		mcp.WithDescription("Получить список клиентов с поиском"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithString("q", mcp.Description("Поиск по имени/email/телефону")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listClientsHandler(svc))

	s.AddTool(mcp.NewTool("get_client",
		mcp.WithDescription("Получить информацию о клиенте по ID"),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getClientHandler(svc))

	s.AddTool(mcp.NewTool("create_client",
		mcp.WithDescription("Создать нового клиента"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Имя")),
		mcp.WithString("surname", mcp.Required(), mcp.Description("Фамилия")),
		mcp.WithString("middlename", mcp.Description("Отчество")),
		mcp.WithString("email", mcp.Description("Email")),
		mcp.WithString("phone", mcp.Description("Телефон")),
		mcp.WithString("password", mcp.Description("Пароль (минимум 6 символов)")),
		mcp.WithBoolean("registered", mcp.Description("Зарегистрирован (требуется email и пароль)")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createClientHandler(svc))

	s.AddTool(mcp.NewTool("update_client",
		mcp.WithDescription("Обновить данные клиента"),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithString("email", mcp.Description("Email")),
		mcp.WithString("phone", mcp.Description("Телефон")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateClientHandler(svc))

	s.AddTool(mcp.NewTool("delete_client",
		mcp.WithDescription("Удалить клиента"),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteClientHandler(svc))

	// === КАТЕГОРИИ ===
	s.AddTool(mcp.NewTool("list_categories",
		mcp.WithDescription("Получить список категорий товаров"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
	mcp.WithReadOnlyHintAnnotation(true),
	), listCategoriesHandler(svc))

	s.AddTool(mcp.NewTool("get_category",
		mcp.WithDescription("Получить информацию о категории по ID"),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("ID категории")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getCategoryHandler(svc))

	s.AddTool(mcp.NewTool("create_category",
		mcp.WithDescription("Создать новую категорию"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название категории")),
		mcp.WithNumber("parent_id", mcp.Description("ID родительской категории")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createCategoryHandler(svc))

	s.AddTool(mcp.NewTool("update_category",
		mcp.WithDescription("Обновить категорию"),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("ID категории")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateCategoryHandler(svc))

	s.AddTool(mcp.NewTool("delete_category",
		mcp.WithDescription("Удалить категорию"),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("ID категории")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteCategoryHandler(svc))

	// === КОЛЛЕКЦИИ ===
	s.AddTool(mcp.NewTool("list_collections",
		mcp.WithDescription("Получить список коллекций (умных подборок)"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
	mcp.WithReadOnlyHintAnnotation(true),
	), listCollectionsHandler(svc))

	s.AddTool(mcp.NewTool("get_collection",
		mcp.WithDescription("Получить информацию о коллекции по ID"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getCollectionHandler(svc))

	s.AddTool(mcp.NewTool("create_collection",
		mcp.WithDescription("Создать новую коллекцию"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название коллекции")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createCollectionHandler(svc))

	s.AddTool(mcp.NewTool("update_collection",
		mcp.WithDescription("Обновить коллекцию"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateCollectionHandler(svc))

	s.AddTool(mcp.NewTool("delete_collection",
		mcp.WithDescription("Удалить коллекцию"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteCollectionHandler(svc))

	// === ИЗОБРАЖЕНИЯ ===
	s.AddTool(mcp.NewTool("list_images",
		mcp.WithDescription("Получить список изображений товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listImagesHandler(svc))

	s.AddTool(mcp.NewTool("add_image",
		mcp.WithDescription("Добавить изображение к товару"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("url", mcp.Required(), mcp.Description("URL изображения")),
		mcp.WithString("filename", mcp.Description("Имя файла")),
		mcp.WithNumber("position", mcp.Description("Позиция")),
		mcp.WithReadOnlyHintAnnotation(false),
	), addImageHandler(svc))

	s.AddTool(mcp.NewTool("delete_image",
		mcp.WithDescription("Удалить изображение товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("image_id", mcp.Required(), mcp.Description("ID изображения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteImageHandler(svc))

	// === МЕТАПОЛЯ ===
	s.AddTool(mcp.NewTool("list_meta_fields",
		mcp.WithDescription("Получить список метаполей"),
		mcp.WithString("resource_type", mcp.Description("Тип ресурса: Product, Variant, Order, Client, Category, Page")),
		mcp.WithString("resource_id", mcp.Description("ID ресурса")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listMetaFieldsHandler(svc))

	s.AddTool(mcp.NewTool("create_meta_field",
		mcp.WithDescription("Создать метаполе для ресурса"),
		mcp.WithString("resource_type", mcp.Required(), mcp.Description("Тип ресурса")),
		mcp.WithString("resource_id", mcp.Required(), mcp.Description("ID ресурса")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Ключ")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Значение")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createMetaFieldHandler(svc))

	s.AddTool(mcp.NewTool("update_meta_field",
		mcp.WithDescription("Обновить метаполе"),
		mcp.WithString("meta_field_id", mcp.Required(), mcp.Description("ID метаполя")),
		mcp.WithString("value", mcp.Description("Значение")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateMetaFieldHandler(svc))

	s.AddTool(mcp.NewTool("delete_meta_field",
		mcp.WithDescription("Удалить метаполе"),
		mcp.WithString("meta_field_id", mcp.Required(), mcp.Description("ID метаполя")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteMetaFieldHandler(svc))

	// === ВЕБХУКИ ===
	s.AddTool(mcp.NewTool("list_webhooks",
		mcp.WithDescription("Получить список подписанных вебхуков"),
	mcp.WithReadOnlyHintAnnotation(true),
	), listWebhooksHandler(svc))

	s.AddTool(mcp.NewTool("create_webhook",
		mcp.WithDescription("Подписаться на событие InSales"),
		mcp.WithString("event", mcp.Required(), mcp.Description("Событие: order_create, product_update и т.д.")),
		mcp.WithString("url", mcp.Required(), mcp.Description("URL для получения событий")),
		mcp.WithString("format", mcp.Description("Формат: json, xml"), mcp.DefaultString("json")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createWebhookHandler(svc))

	s.AddTool(mcp.NewTool("update_webhook",
		mcp.WithDescription("Обновить URL вебхука"),
		mcp.WithString("webhook_id", mcp.Required(), mcp.Description("ID вебхука")),
		mcp.WithString("url", mcp.Description("Новый URL")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateWebhookHandler(svc))

	s.AddTool(mcp.NewTool("delete_webhook",
		mcp.WithDescription("Отписаться от вебхука"),
		mcp.WithString("webhook_id", mcp.Required(), mcp.Description("ID вебхука")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteWebhookHandler(svc))

	// === СИСТЕМНЫЕ ===
	s.AddTool(mcp.NewTool("get_shop_info",
		mcp.WithDescription("Получить информацию о магазине (название, домен, валюта)"),
	mcp.WithReadOnlyHintAnnotation(true),
	), getShopInfoHandler(svc))

	s.AddTool(mcp.NewTool("list_currencies",
		mcp.WithDescription("Получить список валют магазина"),
	mcp.WithReadOnlyHintAnnotation(true),
	), listCurrenciesHandler(svc))

	s.AddTool(mcp.NewTool("list_delivery_variants",
		mcp.WithDescription("Получить доступные способы доставки"),
	mcp.WithReadOnlyHintAnnotation(true),
	), listDeliveryVariantsHandler(svc))

	s.AddTool(mcp.NewTool("list_payment_gateways",
		mcp.WithDescription("Получить доступные способы оплаты"),
	mcp.WithReadOnlyHintAnnotation(true),
	), listPaymentGatewaysHandler(svc))

	// === СТРАНИЦЫ ===
	s.AddTool(mcp.NewTool("list_pages",
		mcp.WithDescription("Получить список страниц магазина"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
	mcp.WithReadOnlyHintAnnotation(true),
	), listPagesHandler(svc))

	s.AddTool(mcp.NewTool("get_page",
		mcp.WithDescription("Получить страницу по ID"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("ID страницы")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getPageHandler(svc))

	s.AddTool(mcp.NewTool("create_page",
		mcp.WithDescription("Создать новую страницу в магазине"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название страницы")),
		mcp.WithString("content", mcp.Required(), mcp.Description("HTML-содержимое")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createPageHandler(svc))

	s.AddTool(mcp.NewTool("update_page",
		mcp.WithDescription("Обновить данные страницы"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("ID страницы")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithString("content", mcp.Description("Содержимое")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updatePageHandler(svc))

	s.AddTool(mcp.NewTool("delete_page",
		mcp.WithDescription("Удалить страницу"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("ID страницы")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deletePageHandler(svc))

	// === АКСЕССУАРЫ ===
	s.AddTool(mcp.NewTool("list_accessories",
		mcp.WithDescription("Получить список аксессуаров (характеристик товаров)"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
	mcp.WithReadOnlyHintAnnotation(true),
	), listAccessoriesHandler(svc))

	s.AddTool(mcp.NewTool("get_accessory",
		mcp.WithDescription("Получить аксессуар по ID"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getAccessoryHandler(svc))

	s.AddTool(mcp.NewTool("create_accessory",
		mcp.WithDescription("Создать новый аксессуар (характеристику)"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название (например 'Цвет', 'Размер')")),
		mcp.WithNumber("max_count", mcp.Description("Максимальное количество опций в заказе")),
		mcp.WithNumber("min_count", mcp.Description("Минимальное количество опций в заказе")),
		mcp.WithString("permalink", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createAccessoryHandler(svc))

	s.AddTool(mcp.NewTool("update_accessory",
		mcp.WithDescription("Обновить аксессуар"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithNumber("position", mcp.Description("Позиция")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateAccessoryHandler(svc))

	s.AddTool(mcp.NewTool("delete_accessory",
		mcp.WithDescription("Удалить аксессуар"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteAccessoryHandler(svc))

	// === ВАРИАНТЫ ТОВАРОВ ===
	s.AddTool(mcp.NewTool("list_variants",
		mcp.WithDescription("Получить список вариантов товара (размеры, цвета)"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listVariantsHandler(svc))

	s.AddTool(mcp.NewTool("get_variant",
		mcp.WithDescription("Получить информацию о варианте товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("ID варианта")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getVariantHandler(svc))

	s.AddTool(mcp.NewTool("update_variant",
		mcp.WithDescription("Обновить данные варианта товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("ID варианта")),
		mcp.WithString("price", mcp.Description("Цена")),
		mcp.WithString("sku", mcp.Description("Артикул")),
		mcp.WithNumber("quantity", mcp.Description("Количество")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateVariantHandler(svc))

	s.AddTool(mcp.NewTool("delete_variant",
		mcp.WithDescription("Удалить вариант товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("ID варианта")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteVariantHandler(svc))

	// === ЗНАЧЕНИЯ СВОЙСТВ ===
	s.AddTool(mcp.NewTool("list_property_values",
		mcp.WithDescription("Получить список значений аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listPropertyValuesHandler(svc))

	s.AddTool(mcp.NewTool("create_property_value",
		mcp.WithDescription("Создать новое значение для аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название значения")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createPropertyValueHandler(svc))

	s.AddTool(mcp.NewTool("update_property_value",
		mcp.WithDescription("Обновить значение аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("property_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updatePropertyValueHandler(svc))

	s.AddTool(mcp.NewTool("delete_property_value",
		mcp.WithDescription("Удалить значение аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("property_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deletePropertyValueHandler(svc))

	// === МЕНЮ ===
	s.AddTool(mcp.NewTool("list_menus",
		mcp.WithDescription("Получить список навигационных меню"),
	mcp.WithReadOnlyHintAnnotation(true),
	), listMenusHandler(svc))

	s.AddTool(mcp.NewTool("get_menu",
		mcp.WithDescription("Получить меню по ID со списком пунктов"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("ID меню")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getMenuHandler(svc))

	s.AddTool(mcp.NewTool("create_menu",
		mcp.WithDescription("Создать новое навигационное меню"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название меню")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createMenuHandler(svc))

	s.AddTool(mcp.NewTool("update_menu",
		mcp.WithDescription("Обновить данные меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("ID меню")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateMenuHandler(svc))

	s.AddTool(mcp.NewTool("delete_menu",
		mcp.WithDescription("Удалить меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("ID меню")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteMenuHandler(svc))

	// === РЕДИРЕКТЫ ===
	s.AddTool(mcp.NewTool("list_redirects",
		mcp.WithDescription("Получить список URL-редиректов"),
	mcp.WithReadOnlyHintAnnotation(true),
	), listRedirectsHandler(svc))

	s.AddTool(mcp.NewTool("create_redirect",
		mcp.WithDescription("Создать URL-редирект"),
		mcp.WithString("url", mcp.Required(), mcp.Description("Старый URL")),
		mcp.WithString("target", mcp.Required(), mcp.Description("Новый URL")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createRedirectHandler(svc))

	s.AddTool(mcp.NewTool("delete_redirect",
		mcp.WithDescription("Удалить URL-редирект"),
		mcp.WithString("redirect_id", mcp.Required(), mcp.Description("ID редиректа")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteRedirectHandler(svc))

	// === БЛОГИ ===
	s.AddTool(mcp.NewTool("list_blogs",
		mcp.WithDescription("Получить список блогов магазина"),
	mcp.WithReadOnlyHintAnnotation(true),
	), listBlogsHandler(svc))

	s.AddTool(mcp.NewTool("get_blog",
		mcp.WithDescription("Получить информацию о блоге по ID"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getBlogHandler(svc))

	s.AddTool(mcp.NewTool("create_blog",
		mcp.WithDescription("Создать новый блог в магазине"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название блога")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createBlogHandler(svc))

	s.AddTool(mcp.NewTool("update_blog",
		mcp.WithDescription("Обновить блог"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("title", mcp.Description("Название блога")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithString("commentable", mcp.Description("Комментарии: no, yes, moderable")),
		mcp.WithBoolean("notify_about_comment", mcp.Description("Уведомлять о комментариях")),
		mcp.WithBoolean("captcha_enabled", mcp.Description("Капча при комментировании")),
		mcp.WithBoolean("disallow_links_in_comments", mcp.Description("Запретить ссылки в комментариях")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateBlogHandler(svc))

	s.AddTool(mcp.NewTool("delete_blog",
		mcp.WithDescription("Удалить блог"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteBlogHandler(svc))

	// === ЗАПИСИ БЛОГА ===
	s.AddTool(mcp.NewTool("list_blog_posts",
		mcp.WithDescription("Получить список записей в блоге"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
	mcp.WithReadOnlyHintAnnotation(true),
	), listBlogPostsHandler(svc))

	s.AddTool(mcp.NewTool("get_blog_post",
		mcp.WithDescription("Получить запись блога по ID"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("ID записи")),
	mcp.WithReadOnlyHintAnnotation(true),
	), getBlogPostHandler(svc))

	s.AddTool(mcp.NewTool("create_blog_post",
		mcp.WithDescription("Создать новую запись в блоге"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Заголовок")),
		mcp.WithString("body_html", mcp.Required(), mcp.Description("HTML-содержимое")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createBlogPostHandler(svc))

	s.AddTool(mcp.NewTool("update_blog_post",
		mcp.WithDescription("Обновить запись блога"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("ID записи")),
		mcp.WithString("title", mcp.Description("Заголовок")),
		mcp.WithString("body_html", mcp.Description("Содержимое")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateBlogPostHandler(svc))

	s.AddTool(mcp.NewTool("delete_blog_post",
		mcp.WithDescription("Удалить запись блога"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("ID записи")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteBlogPostHandler(svc))

	// === ПРОМОКОДЫ (Discount Codes) ===
	s.AddTool(mcp.NewTool("list_discount_codes",
		mcp.WithDescription("Получить список промокодов"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
	), listDiscountCodesHandler(svc))

	s.AddTool(mcp.NewTool("get_discount_code",
		mcp.WithDescription("Получить промокод по ID"),
		mcp.WithString("discount_code_id", mcp.Required(), mcp.Description("ID промокода")),
		mcp.WithReadOnlyHintAnnotation(true),
	), getDiscountCodeHandler(svc))

	s.AddTool(mcp.NewTool("create_discount_code",
		mcp.WithDescription("Создать промокод"),
		mcp.WithString("code", mcp.Required(), mcp.Description("Код промокода")),
		mcp.WithNumber("type_id", mcp.Required(), mcp.Description("Тип: 1=процент, 2=фиксированная сумма")),
		mcp.WithNumber("discount", mcp.Required(), mcp.Description("Скидка")),
		mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithString("expired_at", mcp.Description("Дата окончания (YYYY-MM-DD)")),
		mcp.WithNumber("min_price", mcp.Description("Минимальная сумма заказа")),
		mcp.WithBoolean("act_once", mcp.Description("Одноразовый")),
		mcp.WithBoolean("act_once_for_client", mcp.Description("Одноразовый для клиента")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createDiscountCodeHandler(svc))

	s.AddTool(mcp.NewTool("update_discount_code",
		mcp.WithDescription("Обновить промокод"),
		mcp.WithString("discount_code_id", mcp.Required(), mcp.Description("ID промокода")),
		mcp.WithString("code", mcp.Description("Код промокода")),
		mcp.WithNumber("discount", mcp.Description("Скидка")),
		mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithBoolean("disabled", mcp.Description("Отключить")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateDiscountCodeHandler(svc))

	s.AddTool(mcp.NewTool("delete_discount_code",
		mcp.WithDescription("Удалить промокод"),
		mcp.WithString("discount_code_id", mcp.Required(), mcp.Description("ID промокода")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteDiscountCodeHandler(svc))

	// === СВЯЗИ ТОВАР-КОЛЛЕКЦИЯ (Collects) ===
	s.AddTool(mcp.NewTool("list_collects",
		mcp.WithDescription("Получить список связей товар-коллекция"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
	), listCollectsHandler(svc))

	s.AddTool(mcp.NewTool("get_collects_by_product",
		mcp.WithDescription("Получить коллекции товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(true),
	), getCollectsByProductHandler(svc))

	s.AddTool(mcp.NewTool("add_product_to_collection",
		mcp.WithDescription("Добавить товар в коллекцию"),
		mcp.WithNumber("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithNumber("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createCollectHandler(svc))

	s.AddTool(mcp.NewTool("remove_product_from_collection",
		mcp.WithDescription("Удалить товар из коллекции"),
		mcp.WithString("collect_id", mcp.Required(), mcp.Description("ID связи")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteCollectHandler(svc))

	// === СВОЙСТВА ТОВАРОВ (Properties) ===
	s.AddTool(mcp.NewTool("list_properties",
		mcp.WithDescription("Получить список свойств товаров (Бренд, Цвет...)"),
		mcp.WithReadOnlyHintAnnotation(true),
	), listPropertiesHandler(svc))

	s.AddTool(mcp.NewTool("get_property",
		mcp.WithDescription("Получить свойство по ID"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithReadOnlyHintAnnotation(true),
	), getPropertyHandler(svc))

	s.AddTool(mcp.NewTool("create_property",
		mcp.WithDescription("Создать свойство товаров"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название свойства")),
		mcp.WithString("permalink", mcp.Description("URL-идентификатор")),
		mcp.WithBoolean("is_navigational", mcp.Description("Показывать в фильтре")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createPropertyHandler(svc))

	s.AddTool(mcp.NewTool("update_property",
		mcp.WithDescription("Обновить свойство товаров"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithBoolean("is_navigational", mcp.Description("Показывать в фильтре")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updatePropertyHandler(svc))

	s.AddTool(mcp.NewTool("delete_property",
		mcp.WithDescription("Удалить свойство товаров"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deletePropertyHandler(svc))

	// === ХАРАКТЕРИСТИКИ (Characteristics) ===
	s.AddTool(mcp.NewTool("list_characteristics",
		mcp.WithDescription("Получить список характеристик свойства (Nike, Adidas...)"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithReadOnlyHintAnnotation(true),
	), listCharacteristicsHandler(svc))

	s.AddTool(mcp.NewTool("get_characteristic",
		mcp.WithDescription("Получить характеристику по ID"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("characteristic_id", mcp.Required(), mcp.Description("ID характеристики")),
		mcp.WithReadOnlyHintAnnotation(true),
	), getCharacteristicHandler(svc))

	s.AddTool(mcp.NewTool("create_characteristic",
		mcp.WithDescription("Создать характеристику свойства"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название характеристики")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createCharacteristicHandler(svc))

	s.AddTool(mcp.NewTool("update_characteristic",
		mcp.WithDescription("Обновить характеристику свойства"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("characteristic_id", mcp.Required(), mcp.Description("ID характеристики")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateCharacteristicHandler(svc))

	s.AddTool(mcp.NewTool("delete_characteristic",
		mcp.WithDescription("Удалить характеристику свойства"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("characteristic_id", mcp.Required(), mcp.Description("ID характеристики")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteCharacteristicHandler(svc))

	// === ОТЗЫВЫ (Reviews) ===
	s.AddTool(mcp.NewTool("list_reviews",
		mcp.WithDescription("Получить список отзывов"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
	), listReviewsHandler(svc))

	s.AddTool(mcp.NewTool("get_review",
		mcp.WithDescription("Получить отзыв по ID"),
		mcp.WithString("review_id", mcp.Required(), mcp.Description("ID отзыва")),
		mcp.WithReadOnlyHintAnnotation(true),
	), getReviewHandler(svc))

	s.AddTool(mcp.NewTool("create_review",
		mcp.WithDescription("Создать отзыв"),
		mcp.WithNumber("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("author", mcp.Required(), mcp.Description("Автор")),
		mcp.WithString("content", mcp.Description("Текст отзыва")),
		mcp.WithString("email", mcp.Description("Email")),
		mcp.WithNumber("rating", mcp.Description("Оценка (1-10)")),
		mcp.WithBoolean("published", mcp.Description("Опубликован")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createReviewHandler(svc))

	s.AddTool(mcp.NewTool("update_review",
		mcp.WithDescription("Обновить отзыв"),
		mcp.WithString("review_id", mcp.Required(), mcp.Description("ID отзыва")),
		mcp.WithString("content", mcp.Description("Текст отзыва")),
		mcp.WithNumber("rating", mcp.Description("Оценка")),
		mcp.WithBoolean("published", mcp.Description("Опубликован")),
		mcp.WithString("manager_reply", mcp.Description("Ответ менеджера")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updateReviewHandler(svc))

	s.AddTool(mcp.NewTool("delete_review",
		mcp.WithDescription("Удалить отзыв"),
		mcp.WithString("review_id", mcp.Required(), mcp.Description("ID отзыва")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteReviewHandler(svc))

	// === ФАЙЛЫ (Files) ===
	s.AddTool(mcp.NewTool("list_files",
		mcp.WithDescription("Получить список файлов магазина"),
		mcp.WithReadOnlyHintAnnotation(true),
	), listFilesHandler(svc))

	s.AddTool(mcp.NewTool("get_file",
		mcp.WithDescription("Получить файл по ID"),
		mcp.WithString("file_id", mcp.Required(), mcp.Description("ID файла")),
		mcp.WithReadOnlyHintAnnotation(true),
	), getFileHandler(svc))

	s.AddTool(mcp.NewTool("create_file",
		mcp.WithDescription("Загрузить файл (по URL)"),
		mcp.WithString("src", mcp.Required(), mcp.Description("URL файла")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createFileHandler(svc))

	s.AddTool(mcp.NewTool("delete_file",
		mcp.WithDescription("Удалить файл"),
		mcp.WithString("file_id", mcp.Required(), mcp.Description("ID файла")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deleteFileHandler(svc))

	// === DOMAINS ===
	s.AddTool(mcp.NewTool("list_domains", mcp.WithDescription("Список доменов магазина"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListDomains))
	s.AddTool(mcp.NewTool("get_domain", mcp.WithDescription("Получить домен по ID"), mcp.WithString("domain_id", mcp.Required(), mcp.Description("ID домена")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "domain_id", svc.GetDomain))
	s.AddTool(mcp.NewTool("create_domain", mcp.WithDescription("Добавить домен"), mcp.WithString("domain", mcp.Required(), mcp.Description("Домен")), mcp.WithBoolean("main", mcp.Description("Основной домен")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateDomain))
	s.AddTool(mcp.NewTool("update_domain", mcp.WithDescription("Обновить домен"), mcp.WithString("domain_id", mcp.Required(), mcp.Description("ID домена")), mcp.WithBoolean("main", mcp.Description("Основной домен")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "domain_id", svc.UpdateDomain))
	s.AddTool(mcp.NewTool("delete_domain", mcp.WithDescription("Удалить домен"), mcp.WithString("domain_id", mcp.Required(), mcp.Description("ID домена")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteDomain, "domain_id", "Домен удалён"))

	// === CLIENT GROUPS ===
	s.AddTool(mcp.NewTool("list_client_groups", mcp.WithDescription("Список групп клиентов"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListClientGroups))
	s.AddTool(mcp.NewTool("get_client_group", mcp.WithDescription("Получить группу клиентов"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "group_id", svc.GetClientGroup))
	s.AddTool(mcp.NewTool("create_client_group", mcp.WithDescription("Создать группу клиентов"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithNumber("discount", mcp.Description("Скидка (%)")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateClientGroup))
	s.AddTool(mcp.NewTool("update_client_group", mcp.WithDescription("Обновить группу клиентов"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "group_id", svc.UpdateClientGroup))
	s.AddTool(mcp.NewTool("delete_client_group", mcp.WithDescription("Удалить группу клиентов"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteClientGroup, "group_id", "Группа удалена"))

	// === CUSTOM STATUSES ===
	s.AddTool(mcp.NewTool("list_custom_statuses", mcp.WithDescription("Список кастомных статусов"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListCustomStatuses))
	s.AddTool(mcp.NewTool("get_custom_status", mcp.WithDescription("Получить кастомный статус"), mcp.WithString("status_id", mcp.Required(), mcp.Description("ID статуса")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "status_id", svc.GetCustomStatus))
	s.AddTool(mcp.NewTool("create_custom_status", mcp.WithDescription("Создать кастомный статус"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("status_type", mcp.Required(), mcp.Description("Тип статуса")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateCustomStatus))
	s.AddTool(mcp.NewTool("update_custom_status", mcp.WithDescription("Обновить кастомный статус"), mcp.WithString("status_id", mcp.Required(), mcp.Description("ID статуса")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "status_id", svc.UpdateCustomStatus))
	s.AddTool(mcp.NewTool("delete_custom_status", mcp.WithDescription("Удалить кастомный статус"), mcp.WithString("status_id", mcp.Required(), mcp.Description("ID статуса")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteCustomStatus, "status_id", "Статус удалён"))

	// === PRODUCT FIELDS ===
	s.AddTool(mcp.NewTool("list_product_fields", mcp.WithDescription("Список полей товаров"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListProductFields))
	s.AddTool(mcp.NewTool("get_product_field", mcp.WithDescription("Получить поле товара"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "field_id", svc.GetProductField))
	s.AddTool(mcp.NewTool("create_product_field", mcp.WithDescription("Создать поле товара"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("field_type", mcp.Required(), mcp.Description("Тип поля")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateProductField))
	s.AddTool(mcp.NewTool("update_product_field", mcp.WithDescription("Обновить поле товара"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "field_id", svc.UpdateProductField))
	s.AddTool(mcp.NewTool("delete_product_field", mcp.WithDescription("Удалить поле товара"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteProductField, "field_id", "Поле удалено"))

	// === VARIANT FIELDS ===
	s.AddTool(mcp.NewTool("list_variant_fields", mcp.WithDescription("Список полей вариантов"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListVariantFields))
	s.AddTool(mcp.NewTool("get_variant_field", mcp.WithDescription("Получить поле варианта"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "field_id", svc.GetVariantField))
	s.AddTool(mcp.NewTool("create_variant_field", mcp.WithDescription("Создать поле варианта"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("field_type", mcp.Required(), mcp.Description("Тип поля")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateVariantField))
	s.AddTool(mcp.NewTool("update_variant_field", mcp.WithDescription("Обновить поле варианта"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "field_id", svc.UpdateVariantField))
	s.AddTool(mcp.NewTool("delete_variant_field", mcp.WithDescription("Удалить поле варианта"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteVariantField, "field_id", "Поле удалено"))

	// === COLLECTION FIELDS ===
	s.AddTool(mcp.NewTool("list_collection_fields", mcp.WithDescription("Список полей коллекций"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListCollectionFields))
	s.AddTool(mcp.NewTool("get_collection_field", mcp.WithDescription("Получить поле коллекции"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "field_id", svc.GetCollectionField))
	s.AddTool(mcp.NewTool("create_collection_field", mcp.WithDescription("Создать поле коллекции"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("field_type", mcp.Required(), mcp.Description("Тип поля")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateCollectionField))
	s.AddTool(mcp.NewTool("update_collection_field", mcp.WithDescription("Обновить поле коллекции"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "field_id", svc.UpdateCollectionField))
	s.AddTool(mcp.NewTool("delete_collection_field", mcp.WithDescription("Удалить поле коллекции"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteCollectionField, "field_id", "Поле удалено"))

	// === COLLECTION FILTERS ===
	s.AddTool(mcp.NewTool("list_collection_filters", mcp.WithDescription("Список фильтров коллекций"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListCollectionFilters))
	s.AddTool(mcp.NewTool("get_collection_filter", mcp.WithDescription("Получить фильтр коллекции"), mcp.WithString("filter_id", mcp.Required(), mcp.Description("ID фильтра")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "filter_id", svc.GetCollectionFilter))
	s.AddTool(mcp.NewTool("create_collection_filter", mcp.WithDescription("Создать фильтр коллекции"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("filter_type", mcp.Required(), mcp.Description("Тип фильтра")), mcp.WithNumber("collection_id", mcp.Required(), mcp.Description("ID коллекции")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateCollectionFilter))
	s.AddTool(mcp.NewTool("update_collection_filter", mcp.WithDescription("Обновить фильтр коллекции"), mcp.WithString("filter_id", mcp.Required(), mcp.Description("ID фильтра")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "filter_id", svc.UpdateCollectionFilter))
	s.AddTool(mcp.NewTool("delete_collection_filter", mcp.WithDescription("Удалить фильтр коллекции"), mcp.WithString("filter_id", mcp.Required(), mcp.Description("ID фильтра")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteCollectionFilter, "filter_id", "Фильтр удалён"))

	// === JS TAGS ===
	s.AddTool(mcp.NewTool("list_js_tags", mcp.WithDescription("Список JS-тегов"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListJsTags))
	s.AddTool(mcp.NewTool("get_js_tag", mcp.WithDescription("Получить JS-тег"), mcp.WithString("tag_id", mcp.Required(), mcp.Description("ID тега")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "tag_id", svc.GetJsTag))
	s.AddTool(mcp.NewTool("create_js_tag", mcp.WithDescription("Создать JS-тег"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("content", mcp.Required(), mcp.Description("Код JS")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateJsTag))
	s.AddTool(mcp.NewTool("update_js_tag", mcp.WithDescription("Обновить JS-тег"), mcp.WithString("tag_id", mcp.Required(), mcp.Description("ID тега")), mcp.WithString("title", mcp.Description("Название")), mcp.WithString("content", mcp.Description("Код JS")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "tag_id", svc.UpdateJsTag))
	s.AddTool(mcp.NewTool("delete_js_tag", mcp.WithDescription("Удалить JS-тег"), mcp.WithString("tag_id", mcp.Required(), mcp.Description("ID тега")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteJsTag, "tag_id", "Тег удалён"))

	// === EXTERNAL DISCOUNTS ===
	s.AddTool(mcp.NewTool("list_external_discounts", mcp.WithDescription("Список внешних скидок"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListExternalDiscounts))
	s.AddTool(mcp.NewTool("get_external_discount", mcp.WithDescription("Получить внешнюю скидку"), mcp.WithString("discount_id", mcp.Required(), mcp.Description("ID скидки")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "discount_id", svc.GetExternalDiscount))
	s.AddTool(mcp.NewTool("create_external_discount", mcp.WithDescription("Создать внешнюю скидку"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithNumber("percent", mcp.Description("Процент скидки")), mcp.WithNumber("amount", mcp.Description("Сумма скидки")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateExternalDiscount))
	s.AddTool(mcp.NewTool("update_external_discount", mcp.WithDescription("Обновить внешнюю скидку"), mcp.WithString("discount_id", mcp.Required(), mcp.Description("ID скидки")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "discount_id", svc.UpdateExternalDiscount))
	s.AddTool(mcp.NewTool("delete_external_discount", mcp.WithDescription("Удалить внешнюю скидку"), mcp.WithString("discount_id", mcp.Required(), mcp.Description("ID скидки")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteExternalDiscount, "discount_id", "Скидка удалена"))

	// === DISCOUNT COLLECTIONS GROUPS ===
	s.AddTool(mcp.NewTool("list_discount_collections_groups", mcp.WithDescription("Список групп коллекций скидок"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListDiscountCollectionsGroups))
	s.AddTool(mcp.NewTool("get_discount_collections_group", mcp.WithDescription("Получить группу коллекций скидок"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "group_id", svc.GetDiscountCollectionsGroup))
	s.AddTool(mcp.NewTool("create_discount_collections_group", mcp.WithDescription("Создать группу коллекций скидок"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateDiscountCollectionsGroup))
	s.AddTool(mcp.NewTool("update_discount_collections_group", mcp.WithDescription("Обновить группу коллекций скидок"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "group_id", svc.UpdateDiscountCollectionsGroup))
	s.AddTool(mcp.NewTool("delete_discount_collections_group", mcp.WithDescription("Удалить группу коллекций скидок"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteDiscountCollectionsGroup, "group_id", "Группа удалена"))

	// === WAREHOUSES ===
	s.AddTool(mcp.NewTool("list_warehouses", mcp.WithDescription("Список складов"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListWarehouses))
	s.AddTool(mcp.NewTool("get_warehouse", mcp.WithDescription("Получить склад"), mcp.WithString("warehouse_id", mcp.Required(), mcp.Description("ID склада")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "warehouse_id", svc.GetWarehouse))
	s.AddTool(mcp.NewTool("create_warehouse", mcp.WithDescription("Создать склад"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("address", mcp.Description("Адрес")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateWarehouse))
	s.AddTool(mcp.NewTool("update_warehouse", mcp.WithDescription("Обновить склад"), mcp.WithString("warehouse_id", mcp.Required(), mcp.Description("ID склада")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "warehouse_id", svc.UpdateWarehouse))
	s.AddTool(mcp.NewTool("delete_warehouse", mcp.WithDescription("Удалить склад"), mcp.WithString("warehouse_id", mcp.Required(), mcp.Description("ID склада")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteWarehouse, "warehouse_id", "Склад удалён"))

	// === SALES CHANNELS ===
	s.AddTool(mcp.NewTool("list_sales_channels", mcp.WithDescription("Список каналов продаж"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListSalesChannels))
	s.AddTool(mcp.NewTool("get_sales_channel", mcp.WithDescription("Получить канал продаж"), mcp.WithString("channel_id", mcp.Required(), mcp.Description("ID канала")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "channel_id", svc.GetSalesChannel))
	s.AddTool(mcp.NewTool("create_sales_channel", mcp.WithDescription("Создать канал продаж"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateSalesChannel))
	s.AddTool(mcp.NewTool("update_sales_channel", mcp.WithDescription("Обновить канал продаж"), mcp.WithString("channel_id", mcp.Required(), mcp.Description("ID канала")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "channel_id", svc.UpdateSalesChannel))
	s.AddTool(mcp.NewTool("delete_sales_channel", mcp.WithDescription("Удалить канал продаж"), mcp.WithString("channel_id", mcp.Required(), mcp.Description("ID канала")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteSalesChannel, "channel_id", "Канал удалён"))

	// === STOCK CURRENCIES ===
	s.AddTool(mcp.NewTool("list_stock_currencies", mcp.WithDescription("Список валют склада"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListStockCurrencies))
	s.AddTool(mcp.NewTool("get_stock_currency", mcp.WithDescription("Получить валюту склада"), mcp.WithString("currency_id", mcp.Required(), mcp.Description("ID валюты")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "currency_id", svc.GetStockCurrency))
	s.AddTool(mcp.NewTool("create_stock_currency", mcp.WithDescription("Создать валюту склада"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("code", mcp.Required(), mcp.Description("Код валюты")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateStockCurrency))
	s.AddTool(mcp.NewTool("update_stock_currency", mcp.WithDescription("Обновить валюту склада"), mcp.WithString("currency_id", mcp.Required(), mcp.Description("ID валюты")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "currency_id", svc.UpdateStockCurrency))
	s.AddTool(mcp.NewTool("delete_stock_currency", mcp.WithDescription("Удалить валюту склада"), mcp.WithString("currency_id", mcp.Required(), mcp.Description("ID валюты")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteStockCurrency, "currency_id", "Валюта удалена"))

	// === PRICE KINDS ===
	s.AddTool(mcp.NewTool("list_price_kinds", mcp.WithDescription("Список типов цен"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListPriceKinds))
	s.AddTool(mcp.NewTool("get_price_kind", mcp.WithDescription("Получить тип цены"), mcp.WithString("kind_id", mcp.Required(), mcp.Description("ID типа")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "kind_id", svc.GetPriceKind))
	s.AddTool(mcp.NewTool("create_price_kind", mcp.WithDescription("Создать тип цены"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreatePriceKind))
	s.AddTool(mcp.NewTool("update_price_kind", mcp.WithDescription("Обновить тип цены"), mcp.WithString("kind_id", mcp.Required(), mcp.Description("ID типа")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "kind_id", svc.UpdatePriceKind))
	s.AddTool(mcp.NewTool("delete_price_kind", mcp.WithDescription("Удалить тип цены"), mcp.WithString("kind_id", mcp.Required(), mcp.Description("ID типа")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeletePriceKind, "kind_id", "Тип цены удалён"))

	// === OPTION NAMES ===
	s.AddTool(mcp.NewTool("list_option_names", mcp.WithDescription("Список имён опций (Размер, Цвет...)"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListOptionNames))
	s.AddTool(mcp.NewTool("get_option_name", mcp.WithDescription("Получить имя опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "option_name_id", svc.GetOptionName))
	s.AddTool(mcp.NewTool("create_option_name", mcp.WithDescription("Создать имя опции"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithBoolean("navigational", mcp.Description("Показывать в фильтре")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateOptionName))
	s.AddTool(mcp.NewTool("update_option_name", mcp.WithDescription("Обновить имя опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "option_name_id", svc.UpdateOptionName))
	s.AddTool(mcp.NewTool("delete_option_name", mcp.WithDescription("Удалить имя опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteOptionName, "option_name_id", "Имя опции удалено"))

	// === OPTION VALUES ===
	s.AddTool(mcp.NewTool("list_option_values", mcp.WithDescription("Список значений опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithReadOnlyHintAnnotation(true)), genericListOptionValuesHandler(svc))
	s.AddTool(mcp.NewTool("get_option_value", mcp.WithDescription("Получить значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("option_value_id", mcp.Required(), mcp.Description("ID значения")), mcp.WithReadOnlyHintAnnotation(true)), genericGetOptionValueHandler(svc))
	s.AddTool(mcp.NewTool("create_option_value", mcp.WithDescription("Создать значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateOptionValueHandler(svc))
	s.AddTool(mcp.NewTool("update_option_value", mcp.WithDescription("Обновить значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("option_value_id", mcp.Required(), mcp.Description("ID значения")), mcp.WithString("title", mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateOptionValueHandler(svc))
	s.AddTool(mcp.NewTool("delete_option_value", mcp.WithDescription("Удалить значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("option_value_id", mcp.Required(), mcp.Description("ID значения")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteOptionValueHandler(svc))

	// === ROUND 4: APPLICATION ACTIONS ===
	s.AddTool(mcp.NewTool("list_application_actions", mcp.WithDescription("Список действий приложения"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListApplicationActions))
	s.AddTool(mcp.NewTool("get_application_action", mcp.WithDescription("Получить действие приложения"), mcp.WithString("action_id", mcp.Required(), mcp.Description("ID действия")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "action_id", svc.GetApplicationAction))
	s.AddTool(mcp.NewTool("create_application_action", mcp.WithDescription("Создать действие приложения"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("handle", mcp.Required(), mcp.Description("Идентификатор")), mcp.WithString("entity", mcp.Required(), mcp.Description("Сущность: order, product, client...")), mcp.WithString("url", mcp.Required(), mcp.Description("URL")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateApplicationAction))
	s.AddTool(mcp.NewTool("update_application_action", mcp.WithDescription("Обновить действие приложения"), mcp.WithString("action_id", mcp.Required(), mcp.Description("ID действия")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "action_id", svc.UpdateApplicationAction))
	s.AddTool(mcp.NewTool("delete_application_action", mcp.WithDescription("Удалить действие приложения"), mcp.WithString("action_id", mcp.Required(), mcp.Description("ID действия")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteApplicationAction, "action_id", "Действие удалено"))

	// === APPLICATION CHARGES ===
	s.AddTool(mcp.NewTool("list_application_charges", mcp.WithDescription("Список счетов приложения"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListApplicationCharges))
	s.AddTool(mcp.NewTool("get_application_charge", mcp.WithDescription("Получить счёт приложения"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "charge_id", svc.GetApplicationCharge))
	s.AddTool(mcp.NewTool("create_application_charge", mcp.WithDescription("Создать счёт приложения"), mcp.WithString("name", mcp.Required(), mcp.Description("Назначение")), mcp.WithNumber("price", mcp.Required(), mcp.Description("Сумма")), mcp.WithString("return_url", mcp.Description("URL уведомления")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateApplicationCharge))
	s.AddTool(mcp.NewTool("decline_application_charge", mcp.WithDescription("Отклонить счёт приложения"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")), mcp.WithReadOnlyHintAnnotation(false)), declineApplicationChargeHandler(svc))

	// === APPLICATION WIDGETS ===
	s.AddTool(mcp.NewTool("list_application_widgets", mcp.WithDescription("Список виджетов приложения"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListApplicationWidgets))
	s.AddTool(mcp.NewTool("get_application_widget", mcp.WithDescription("Получить виджет"), mcp.WithString("widget_id", mcp.Required(), mcp.Description("ID виджета")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "widget_id", svc.GetApplicationWidget))
	s.AddTool(mcp.NewTool("create_application_widget", mcp.WithDescription("Создать виджет"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateApplicationWidget))
	s.AddTool(mcp.NewTool("update_application_widget", mcp.WithDescription("Обновить виджет"), mcp.WithString("widget_id", mcp.Required(), mcp.Description("ID виджета")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "widget_id", svc.UpdateApplicationWidget))
	s.AddTool(mcp.NewTool("delete_application_widget", mcp.WithDescription("Удалить виджет"), mcp.WithString("widget_id", mcp.Required(), mcp.Description("ID виджета")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteApplicationWidget, "widget_id", "Виджет удалён"))

	// === RECURRING APPLICATION CHARGES ===
	s.AddTool(mcp.NewTool("list_recurring_application_charges", mcp.WithDescription("Список рекуррентных счетов"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListRecurringApplicationCharges))
	s.AddTool(mcp.NewTool("get_recurring_application_charge", mcp.WithDescription("Получить рекуррентный счёт"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "charge_id", svc.GetRecurringApplicationCharge))
	s.AddTool(mcp.NewTool("create_recurring_application_charge", mcp.WithDescription("Создать рекуррентный счёт"), mcp.WithString("name", mcp.Required(), mcp.Description("Назначение")), mcp.WithNumber("price", mcp.Required(), mcp.Description("Сумма")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateRecurringApplicationCharge))
	s.AddTool(mcp.NewTool("delete_recurring_application_charge", mcp.WithDescription("Удалить рекуррентный счёт"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteRecurringApplicationCharge, "charge_id", "Счёт удалён"))

	// === BONUS SYSTEM TRANSACTIONS ===
	s.AddTool(mcp.NewTool("list_bonus_system_transactions", mcp.WithDescription("Список бонусных транзакций"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListBonusSystemTransactions))
	s.AddTool(mcp.NewTool("get_bonus_system_transaction", mcp.WithDescription("Получить бонусную транзакцию"), mcp.WithString("transaction_id", mcp.Required(), mcp.Description("ID транзакции")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "transaction_id", svc.GetBonusSystemTransaction))
	s.AddTool(mcp.NewTool("create_bonus_system_transaction", mcp.WithDescription("Создать бонусную транзакцию"), mcp.WithNumber("client_id", mcp.Required(), mcp.Description("ID клиента")), mcp.WithNumber("amount", mcp.Required(), mcp.Description("Сумма баллов")), mcp.WithString("comment", mcp.Description("Комментарий")), mcp.WithReadOnlyHintAnnotation(false)), createBonusSystemTransactionHandler(svc))
	s.AddTool(mcp.NewTool("list_client_bonus_transactions", mcp.WithDescription("Список бонусных транзакций клиента"), mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")), mcp.WithReadOnlyHintAnnotation(true)), listClientBonusTransactionsHandler(svc))

	// === OAUTH APPS ===
	s.AddTool(mcp.NewTool("list_oauth_apps", mcp.WithDescription("Список OAuth приложений"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListOAuthApps))
	s.AddTool(mcp.NewTool("get_oauth_app", mcp.WithDescription("Получить OAuth приложение"), mcp.WithString("app_id", mcp.Required(), mcp.Description("ID приложения")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "app_id", svc.GetOAuthApp))
	s.AddTool(mcp.NewTool("create_oauth_app", mcp.WithDescription("Создать OAuth приложение"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateOAuthApp))
	s.AddTool(mcp.NewTool("update_oauth_app", mcp.WithDescription("Обновить OAuth приложение"), mcp.WithString("app_id", mcp.Required(), mcp.Description("ID приложения")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "app_id", svc.UpdateOAuthApp))
	s.AddTool(mcp.NewTool("delete_oauth_app", mcp.WithDescription("Удалить OAuth приложение"), mcp.WithString("app_id", mcp.Required(), mcp.Description("ID приложения")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteOAuthApp, "app_id", "Приложение удалено"))

	// === USER TASKS ===
	s.AddTool(mcp.NewTool("list_user_tasks", mcp.WithDescription("Список задач пользователя"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListUserTasks))
	s.AddTool(mcp.NewTool("get_user_task", mcp.WithDescription("Получить задачу"), mcp.WithString("task_id", mcp.Required(), mcp.Description("ID задачи")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "task_id", svc.GetUserTask))
	s.AddTool(mcp.NewTool("create_user_task", mcp.WithDescription("Создать задачу"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateUserTask))
	s.AddTool(mcp.NewTool("update_user_task", mcp.WithDescription("Обновить задачу"), mcp.WithString("task_id", mcp.Required(), mcp.Description("ID задачи")), mcp.WithReadOnlyHintAnnotation(false)), genericUpdateHandler(svc, "task_id", svc.UpdateUserTask))
	s.AddTool(mcp.NewTool("delete_user_task", mcp.WithDescription("Удалить задачу"), mcp.WithString("task_id", mcp.Required(), mcp.Description("ID задачи")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteUserTask, "task_id", "Задача удалена"))

	// === SUPPLEMENTARY ===
	s.AddTool(mcp.NewTool("list_supplementaries", mcp.WithDescription("Список дополнительных товаров"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListSupplementaries))
	s.AddTool(mcp.NewTool("get_supplementary", mcp.WithDescription("Получить дополнительный товар"), mcp.WithString("supplementary_id", mcp.Required(), mcp.Description("ID")), mcp.WithReadOnlyHintAnnotation(true)), genericGetHandler(svc, "supplementary_id", svc.GetSupplementary))
	s.AddTool(mcp.NewTool("create_supplementary", mcp.WithDescription("Создать дополнительный товар"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateSupplementary))
	s.AddTool(mcp.NewTool("delete_supplementary", mcp.WithDescription("Удалить дополнительный товар"), mcp.WithString("supplementary_id", mcp.Required(), mcp.Description("ID")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteSupplementary, "supplementary_id", "Удалено"))

	// === SIMILARS ===
	s.AddTool(mcp.NewTool("list_similars", mcp.WithDescription("Похожие товары"), mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")), mcp.WithReadOnlyHintAnnotation(true)), listSimilarsByProductHandler(svc))
	s.AddTool(mcp.NewTool("create_similar", mcp.WithDescription("Добавить похожий товар"), mcp.WithNumber("product_id", mcp.Required(), mcp.Description("ID товара")), mcp.WithNumber("similar_product_id", mcp.Required(), mcp.Description("ID похожего товара")), mcp.WithReadOnlyHintAnnotation(false)), genericCreateHandler(svc, svc.CreateSimilar))
	s.AddTool(mcp.NewTool("delete_similar", mcp.WithDescription("Удалить похожий товар"), mcp.WithString("similar_id", mcp.Required(), mcp.Description("ID связи")), mcp.WithReadOnlyHintAnnotation(false), mcp.WithDestructiveHintAnnotation(true)), genericDeleteHandler(svc.DeleteSimilar, "similar_id", "Удалено"))

	// === UNITS ===
	s.AddTool(mcp.NewTool("list_units", mcp.WithDescription("Список единиц измерения"), mcp.WithReadOnlyHintAnnotation(true)), genericListHandler(svc, svc.ListUnits))
}

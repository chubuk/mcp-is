package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/your-project/mcp-insales/internal/service"
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
	), listProductsHandler(svc))

	s.AddTool(mcp.NewTool("get_product",
		mcp.WithDescription("Получить детальную информацию о товаре по ID"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
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
		mcp.WithBoolean("paid", mcp.Description("Оплачен ли")),
	), listOrdersHandler(svc))

	s.AddTool(mcp.NewTool("get_order",
		mcp.WithDescription("Получить детали заказа по ID"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("ID заказа")),
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
	), listClientsHandler(svc))

	s.AddTool(mcp.NewTool("get_client",
		mcp.WithDescription("Получить информацию о клиенте по ID"),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
	), getClientHandler(svc))

	s.AddTool(mcp.NewTool("create_client",
		mcp.WithDescription("Создать нового клиента"),
		mcp.WithString("email", mcp.Required(), mcp.Description("Email")),
		mcp.WithString("first_name", mcp.Description("Имя")),
		mcp.WithString("last_name", mcp.Description("Фамилия")),
		mcp.WithString("phone", mcp.Description("Телефон")),
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
	), listCategoriesHandler(svc))

	s.AddTool(mcp.NewTool("get_category",
		mcp.WithDescription("Получить информацию о категории по ID"),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("ID категории")),
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
	), listCollectionsHandler(svc))

	s.AddTool(mcp.NewTool("get_collection",
		mcp.WithDescription("Получить информацию о коллекции по ID"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
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
	), getShopInfoHandler(svc))

	s.AddTool(mcp.NewTool("list_currencies",
		mcp.WithDescription("Получить список валют магазина"),
	), listCurrenciesHandler(svc))

	s.AddTool(mcp.NewTool("list_delivery_variants",
		mcp.WithDescription("Получить доступные способы доставки"),
	), listDeliveryVariantsHandler(svc))

	s.AddTool(mcp.NewTool("list_payment_gateways",
		mcp.WithDescription("Получить доступные способы оплаты"),
	), listPaymentGatewaysHandler(svc))

	// === СТРАНИЦЫ ===
	s.AddTool(mcp.NewTool("list_pages",
		mcp.WithDescription("Получить список страниц магазина"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
	), listPagesHandler(svc))

	s.AddTool(mcp.NewTool("get_page",
		mcp.WithDescription("Получить страницу по ID"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("ID страницы")),
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
	), listAccessoriesHandler(svc))

	s.AddTool(mcp.NewTool("get_accessory",
		mcp.WithDescription("Получить аксессуар по ID"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
	), getAccessoryHandler(svc))

	s.AddTool(mcp.NewTool("create_accessory",
		mcp.WithDescription("Создать новый аксессуар (характеристику)"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название (например 'Цвет', 'Размер')")),
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
	), listVariantsHandler(svc))

	s.AddTool(mcp.NewTool("get_variant",
		mcp.WithDescription("Получить информацию о варианте товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("ID варианта")),
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
	), listPropertyValuesHandler(svc))

	s.AddTool(mcp.NewTool("create_property_value",
		mcp.WithDescription("Создать новое значение для аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название значения")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createPropertyValueHandler(svc))

	s.AddTool(mcp.NewTool("update_property_value",
		mcp.WithDescription("Обновить значение аксессуара"),
		mcp.WithString("property_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
	), updatePropertyValueHandler(svc))

	s.AddTool(mcp.NewTool("delete_property_value",
		mcp.WithDescription("Удалить значение аксессуара"),
		mcp.WithString("property_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	), deletePropertyValueHandler(svc))

	// === МЕНЮ ===
	s.AddTool(mcp.NewTool("list_menus",
		mcp.WithDescription("Получить список навигационных меню"),
	), listMenusHandler(svc))

	s.AddTool(mcp.NewTool("get_menu",
		mcp.WithDescription("Получить меню по ID со списком пунктов"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("ID меню")),
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
	), listBlogsHandler(svc))

	s.AddTool(mcp.NewTool("get_blog",
		mcp.WithDescription("Получить информацию о блоге по ID"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
	), getBlogHandler(svc))

	s.AddTool(mcp.NewTool("create_blog",
		mcp.WithDescription("Создать новый блог в магазине"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название блога")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
	), createBlogHandler(svc))

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
	), listBlogPostsHandler(svc))

	s.AddTool(mcp.NewTool("get_blog_post",
		mcp.WithDescription("Получить запись блога по ID"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("ID записи")),
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
}

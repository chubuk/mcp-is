package mcp

import (
	"context"

	"github.com/chubuk/mcp-is/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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
			mcp.WithDescription("List product from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
			mcp.WithToolTitle("Список товаров"),
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
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
		), listProductsHandler(svc))

	s.AddTool(mcp.NewTool("get_product",
		mcp.WithDescription("Get a single product by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Товар по ID"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("Product id (string)")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getProductHandler(svc))

	s.AddTool(mcp.NewTool("create_product",
		mcp.WithDescription("Create or add a new product in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новый товар в магазине"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название товара")),
		mcp.WithString("price", mcp.Required(), mcp.Description("Цена (например 1000.00)")),
		mcp.WithString("description", mcp.Description("Описание товара")),
		mcp.WithString("vendor", mcp.Description("Производитель")),
		mcp.WithString("sku", mcp.Description("Артикул")),
		mcp.WithNumber("quantity", mcp.Description("Количество на складе")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createProductHandler(svc))

	s.AddTool(mcp.NewTool("update_product",
		mcp.WithDescription("Update an existing product by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить данные товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithString("price", mcp.Description("Цена")),
		mcp.WithNumber("quantity", mcp.Description("Количество")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateProductHandler(svc))

	s.AddTool(mcp.NewTool("delete_product",
		mcp.WithDescription("Delete or remove product on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить товар из магазина"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteProductHandler(svc))

	// === ЗАКАЗЫ ===
	s.AddTool(mcp.NewTool("list_orders",
		mcp.WithDescription("List order from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список заказов с фильтрацией"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithString("status", mcp.Description("Статус: open, closed, canceled")),
		mcp.WithString("fulfillment_status", mcp.Description("Статус выполнения")),
		mcp.WithString("updated_since", mcp.Description("Дата обновления (ISO format)")),
		mcp.WithNumber("from_id", mcp.Description("ID начиная с которого получать данные")),
		mcp.WithBoolean("paid", mcp.Description("Оплачен ли")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listOrdersHandler(svc))

	s.AddTool(mcp.NewTool("get_order",
		mcp.WithDescription("Get a single order by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить детали заказа по ID"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("ID заказа")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getOrderHandler(svc))

	s.AddTool(mcp.NewTool("update_order",
		mcp.WithDescription("Update an existing order by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить заказ: статус, комментарий, строки, скидки, адрес доставки. status — alias для financial_status"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("ID заказа")),
		mcp.WithString("status", mcp.Description("Alias для financial_status: pending | paid")),
		mcp.WithString("financial_status", mcp.Description("Финансовый статус: pending | paid")),
		mcp.WithString("custom_status_permalink", mcp.Description("Permalink кастомного статуса")),
		mcp.WithString("fulfillment_status", mcp.Description("Статус выполнения")),
		mcp.WithString("comment", mcp.Description("Комментарий клиента")),
		mcp.WithString("manager_comment", mcp.Description("Комментарий менеджера")),
		mcp.WithNumber("delivery_variant_id", mcp.Description("ID способа доставки")),
		mcp.WithNumber("payment_gateway_id", mcp.Description("ID платёжного шлюза")),
		mcp.WithNumber("warehouse_id", mcp.Description("ID склада")),
		mcp.WithString("coupon", mcp.Description("Промокод")),
		mcp.WithArray("order_lines_attributes", mcp.Description("Строки заказа: {id, quantity?, sale_price?, vat?, comment?, _destroy?}"),
			mcp.Items(map[string]any{
				"type":                 "object",
				"additionalProperties": true,
			})),
		mcp.WithArray("discounts_attributes", mcp.Description("Скидки: {discount, type_id? (1=percent,2=money), description?, id?, _destroy?, discount_order_lines_ids?}"),
			mcp.Items(map[string]any{
				"type":                 "object",
				"additionalProperties": true,
			})),
		mcp.WithObject("shipping_address_attributes", mcp.Description("Адрес доставки: {address?, full_locality_name?, ...}")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateOrderHandler(svc))

	s.AddTool(mcp.NewTool("create_order",
		mcp.WithDescription("Create or add a new order in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать заказ: клиент (client_id ИЛИ client name/email/phone) и опционально строки/адрес/оплата"),
		mcp.WithNumber("client_id", mcp.Description("ID существующего клиента")),
		mcp.WithObject("client", mcp.Description("Новый клиент: {name?, email?, phone?, consent_to_personal_data?}")),
		mcp.WithString("comment", mcp.Description("Комментарий")),
		mcp.WithString("manager_comment", mcp.Description("Комментарий менеджера")),
		mcp.WithString("status", mcp.Description("Alias для financial_status: pending | paid")),
		mcp.WithString("financial_status", mcp.Description("Финансовый статус: pending | paid")),
		mcp.WithNumber("delivery_variant_id", mcp.Description("ID способа доставки")),
		mcp.WithNumber("payment_gateway_id", mcp.Description("ID платёжного шлюза")),
		mcp.WithString("coupon", mcp.Description("Промокод")),
		mcp.WithArray("order_lines_attributes", mcp.Description("Строки: {product_id ИЛИ variant_id, quantity, vat?, sale_price?}"),
			mcp.Items(map[string]any{
				"type":                 "object",
				"additionalProperties": true,
			})),
		mcp.WithObject("shipping_address_attributes", mcp.Description("Адрес: {address?, full_locality_name?}")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createOrderHandler(svc))

	// === ORDER LINES / DISCOUNTS / SHIPPING (B1) ===
	s.AddTool(mcp.NewTool("add_order_line",
		mcp.WithDescription("Create or add a new order line in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Добавить строку заказа"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithString("product_id", mcp.Description("Product id (use either product_id or variant_id)")),
		mcp.WithString("variant_id", mcp.Description("Variant id")),
		mcp.WithNumber("quantity", mcp.Required(), mcp.Description("Quantity")),
		mcp.WithNumber("vat", mcp.Description("VAT, use force if needed")),
		mcp.WithNumber("sale_price", mcp.Description("Sale price override")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), addOrderLineHandler(svc))

	s.AddTool(mcp.NewTool("update_order_line",
		mcp.WithDescription("Update an existing order line by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить строку заказа"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithString("line_id", mcp.Required(), mcp.Description("Order line id")),
		mcp.WithNumber("quantity", mcp.Description("New quantity")),
		mcp.WithNumber("vat", mcp.Description("VAT")),
		mcp.WithNumber("sale_price", mcp.Description("Sale price")),
		mcp.WithString("comment", mcp.Description("Line comment")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateOrderLineHandler(svc))

	s.AddTool(mcp.NewTool("remove_order_line",
		mcp.WithDescription("Delete or remove order line on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить строку заказа"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithString("line_id", mcp.Required(), mcp.Description("Order line id")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), removeOrderLineHandler(svc))

	s.AddTool(mcp.NewTool("add_order_discount",
		mcp.WithDescription("Create or add a new order discount in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Добавить скидку к заказу"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithNumber("discount", mcp.Required(), mcp.Description("Discount amount")),
		mcp.WithNumber("type_id", mcp.Description("1=percent, 2=money")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithArray("line_ids", mcp.Description("Affected order line ids; all lines if omitted"), mcp.WithNumberItems()),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), addOrderDiscountHandler(svc))

	s.AddTool(mcp.NewTool("update_order_discount",
		mcp.WithDescription("Update an existing order discount by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить скидку заказа"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithString("discount_id", mcp.Required(), mcp.Description("Discount id")),
		mcp.WithNumber("discount", mcp.Description("Discount amount")),
		mcp.WithNumber("type_id", mcp.Description("1=percent, 2=money")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateOrderDiscountHandler(svc))

	s.AddTool(mcp.NewTool("remove_order_discount",
		mcp.WithDescription("Delete or remove order discount on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить скидку заказа"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithString("discount_id", mcp.Required(), mcp.Description("Discount id")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), removeOrderDiscountHandler(svc))

	s.AddTool(mcp.NewTool("update_order_shipping_address",
		mcp.WithDescription("Update an existing order shipping address by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить адрес доставки"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithString("address", mcp.Description("Street address")),
		mcp.WithString("full_locality_name", mcp.Description("Full locality name")),
		mcp.WithString("country", mcp.Description("Country")),
		mcp.WithString("region", mcp.Description("Region")),
		mcp.WithString("city", mcp.Description("City")),
		mcp.WithString("zip", mcp.Description("Zip")),
		mcp.WithString("phone", mcp.Description("Phone")),
		mcp.WithString("full_name", mcp.Description("Recipient full name")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateOrderShippingAddressHandler(svc))

	s.AddTool(mcp.NewTool("update_order_custom_status",
		mcp.WithDescription("Update an existing custom order status by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить кастомный статус заказа"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("Order id")),
		mcp.WithString("custom_status_permalink", mcp.Required(), mcp.Description("Custom status permalink, e.g. v-obrabotke")),
		mcp.WithString("fulfillment_status", mcp.Description("Optional fulfillment status")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateOrderCustomStatusHandler(svc))

	// === КЛИЕНТЫ ===
	s.AddTool(mcp.NewTool("list_clients",
		mcp.WithDescription("List client from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список клиентов с поиском"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithString("q", mcp.Description("Поиск по имени/email/телефону")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listClientsHandler(svc))

	s.AddTool(mcp.NewTool("get_client",
		mcp.WithDescription("Get a single client by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить информацию о клиенте по ID"),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getClientHandler(svc))

	s.AddTool(mcp.NewTool("create_client",
		mcp.WithDescription("Create or add a new client in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать нового клиента"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Имя")),
		mcp.WithString("surname", mcp.Required(), mcp.Description("Фамилия")),
		mcp.WithString("middlename", mcp.Description("Отчество")),
		mcp.WithString("email", mcp.Description("Email")),
		mcp.WithString("phone", mcp.Description("Телефон")),
		mcp.WithString("password", mcp.Description("Пароль (минимум 6 символов)")),
		mcp.WithBoolean("registered", mcp.Description("Зарегистрирован (требуется email и пароль)")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createClientHandler(svc))

	s.AddTool(mcp.NewTool("update_client",
		mcp.WithDescription("Update an existing client by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить данные клиента"),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithString("email", mcp.Description("Email")),
		mcp.WithString("phone", mcp.Description("Телефон")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateClientHandler(svc))

	s.AddTool(mcp.NewTool("delete_client",
		mcp.WithDescription("Delete or remove client on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить клиента"),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteClientHandler(svc))

	// === КАТЕГОРИИ ===
	s.AddTool(mcp.NewTool("list_categories",
		mcp.WithDescription("List category from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список категорий товаров"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listCategoriesHandler(svc))

	s.AddTool(mcp.NewTool("get_category",
		mcp.WithDescription("Get a single category by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить информацию о категории по ID"),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("ID категории")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getCategoryHandler(svc))

	s.AddTool(mcp.NewTool("create_category",
		mcp.WithDescription("Create or add a new category in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новую категорию"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название категории")),
		mcp.WithNumber("parent_id", mcp.Description("ID родительской категории")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createCategoryHandler(svc))

	s.AddTool(mcp.NewTool("update_category",
		mcp.WithDescription("Update an existing category by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить категорию"),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("ID категории")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateCategoryHandler(svc))

	s.AddTool(mcp.NewTool("delete_category",
		mcp.WithDescription("Delete or remove category on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить категорию"),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("ID категории")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteCategoryHandler(svc))

	// === КОЛЛЕКЦИИ ===
	s.AddTool(mcp.NewTool("list_collections",
		mcp.WithDescription("List collection from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список коллекций (умных подборок)"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listCollectionsHandler(svc))

	s.AddTool(mcp.NewTool("get_collection",
		mcp.WithDescription("Get a single collection by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить информацию о коллекции по ID"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getCollectionHandler(svc))

	s.AddTool(mcp.NewTool("create_collection",
		mcp.WithDescription("Create or add a new collection in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новую коллекцию"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название коллекции")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createCollectionHandler(svc))

	s.AddTool(mcp.NewTool("update_collection",
		mcp.WithDescription("Update an existing collection by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить коллекцию"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateCollectionHandler(svc))

	s.AddTool(mcp.NewTool("delete_collection",
		mcp.WithDescription("Delete or remove collection on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить коллекцию"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteCollectionHandler(svc))

	// === ИЗОБРАЖЕНИЯ ===
	s.AddTool(mcp.NewTool("list_images",
		mcp.WithDescription("List product image from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список изображений товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listImagesHandler(svc))

	s.AddTool(mcp.NewTool("add_image",
		mcp.WithDescription("Create or add a new product image in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Добавить изображение к товару"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("url", mcp.Required(), mcp.Description("URL изображения")),
		mcp.WithString("filename", mcp.Description("Имя файла")),
		mcp.WithNumber("position", mcp.Description("Позиция")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), addImageHandler(svc))

	s.AddTool(mcp.NewTool("delete_image",
		mcp.WithDescription("Delete or remove product image on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить изображение товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("image_id", mcp.Required(), mcp.Description("ID изображения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteImageHandler(svc))

	// === МЕТАПОЛЯ ===
	s.AddTool(mcp.NewTool("list_meta_fields",
		mcp.WithDescription("List meta field from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список метаполей"),
		mcp.WithString("resource_type", mcp.Description("Тип ресурса: Product, Variant, Order, Client, Category, Page")),
		mcp.WithString("resource_id", mcp.Description("ID ресурса")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listMetaFieldsHandler(svc))

	s.AddTool(mcp.NewTool("create_meta_field",
		mcp.WithDescription("Create or add a new meta field in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать метаполе для ресурса"),
		mcp.WithString("resource_type", mcp.Required(), mcp.Description("Тип ресурса")),
		mcp.WithString("resource_id", mcp.Required(), mcp.Description("ID ресурса")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Ключ")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Значение")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createMetaFieldHandler(svc))

	s.AddTool(mcp.NewTool("update_meta_field",
		mcp.WithDescription("Update an existing meta field by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить метаполе"),
		mcp.WithString("meta_field_id", mcp.Required(), mcp.Description("ID метаполя")),
		mcp.WithString("value", mcp.Description("Значение")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateMetaFieldHandler(svc))

	s.AddTool(mcp.NewTool("delete_meta_field",
		mcp.WithDescription("Delete or remove meta field on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить метаполе"),
		mcp.WithString("meta_field_id", mcp.Required(), mcp.Description("ID метаполя")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteMetaFieldHandler(svc))

	// === ВЕБХУКИ ===
	s.AddTool(mcp.NewTool("list_webhooks",
		mcp.WithDescription("List webhook from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список подписанных вебхуков"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listWebhooksHandler(svc))

	s.AddTool(mcp.NewTool("create_webhook",
		mcp.WithDescription("Create or add a new webhook in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Подписаться на событие InSales"),
		mcp.WithString("event", mcp.Required(), mcp.Description("Событие: order_create, product_update и т.д.")),
		mcp.WithString("url", mcp.Required(), mcp.Description("URL для получения событий")),
		mcp.WithString("format", mcp.Description("Формат: json, xml"), mcp.DefaultString("json")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createWebhookHandler(svc))

	s.AddTool(mcp.NewTool("update_webhook",
		mcp.WithDescription("Update an existing webhook by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить URL вебхука"),
		mcp.WithString("webhook_id", mcp.Required(), mcp.Description("ID вебхука")),
		mcp.WithString("url", mcp.Description("Новый URL")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateWebhookHandler(svc))

	s.AddTool(mcp.NewTool("delete_webhook",
		mcp.WithDescription("Delete or remove webhook on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Отписаться от вебхука"),
		mcp.WithString("webhook_id", mcp.Required(), mcp.Description("ID вебхука")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteWebhookHandler(svc))

	// === СИСТЕМНЫЕ ===
	s.AddTool(mcp.NewTool("get_shop_info",
		mcp.WithDescription("Get a single shop by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить информацию о магазине (название, домен, валюта)"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getShopInfoHandler(svc))

	s.AddTool(mcp.NewTool("list_currencies",
		mcp.WithDescription("List currency from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список валют магазина"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listCurrenciesHandler(svc))

	s.AddTool(mcp.NewTool("list_delivery_variants",
		mcp.WithDescription("List delivery variant from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Способы доставки"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listDeliveryVariantsHandler(svc))

	s.AddTool(mcp.NewTool("get_delivery_variant",
		mcp.WithDescription("Get a single delivery variant by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Способ доставки"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Delivery variant id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "variant_id", svc.GetDeliveryVariant))

	s.AddTool(mcp.NewTool("create_delivery_variant",
		mcp.WithDescription("Create or add a new delivery variant in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Создать способ доставки"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Title")),
		mcp.WithString("type", mcp.Required(), mcp.Description("DeliveryVariant::External|Fixed|LocationDepend|PickUp|PriceDepend")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithNumber("position", mcp.Description("Position")),
		mcp.WithBoolean("add_payment_gateways", mcp.Description("Add all payment gateways")),
		mcp.WithBoolean("customer_pickup", mcp.Description("Customer pickup")),
		mcp.WithBoolean("inverted", mcp.Description("Invert locations filter")),
		mcp.WithNumber("charge_up_to", mcp.Description("Free delivery from sum")),
		mcp.WithString("url", mcp.Description("External estimation URL")),
		mcp.WithString("api_version", mcp.Description("v1 or v2")),
		mcp.WithArray("delivery_locations_attributes", mcp.Description("Locations [{country?, region, area?, city?}]"),
			mcp.Items(map[string]any{"type": "object", "additionalProperties": true})),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateDeliveryVariant))

	s.AddTool(mcp.NewTool("update_delivery_variant",
		mcp.WithDescription("Update an existing delivery variant by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить способ доставки"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Delivery variant id")),
		mcp.WithString("title", mcp.Description("Title")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithNumber("position", mcp.Description("Position")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "variant_id", svc.UpdateDeliveryVariant))

	s.AddTool(mcp.NewTool("delete_delivery_variant",
		mcp.WithDescription("Delete or remove delivery variant on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить способ доставки"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Delivery variant id")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteDeliveryVariant, "variant_id", "Способ доставки удалён"))

	s.AddTool(mcp.NewTool("list_payment_gateways",
		mcp.WithDescription("List payment gateway from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Способы оплаты"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listPaymentGatewaysHandler(svc))

	s.AddTool(mcp.NewTool("get_payment_gateway",
		mcp.WithDescription("Get a single payment gateway by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Способ оплаты"),
		mcp.WithString("gateway_id", mcp.Required(), mcp.Description("Payment gateway id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "gateway_id", svc.GetPaymentGateway))

	s.AddTool(mcp.NewTool("create_payment_gateway",
		mcp.WithDescription("Create or add a new payment gateway in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Создать способ оплаты"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Title")),
		mcp.WithString("type", mcp.Required(), mcp.Description("PaymentGateway::Cod|Custom")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithNumber("margin", mcp.Description("Margin percent -100..1000")),
		mcp.WithNumber("position", mcp.Description("Position")),
		mcp.WithBoolean("add_delivery_variants", mcp.Description("Link all delivery variants")),
		mcp.WithArray("payment_delivery_variants_attributes", mcp.Description("[{delivery_variant_id}]"),
			mcp.Items(map[string]any{"type": "object", "additionalProperties": true})),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreatePaymentGateway))

	s.AddTool(mcp.NewTool("update_payment_gateway",
		mcp.WithDescription("Update an existing payment gateway by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить способ оплаты"),
		mcp.WithString("gateway_id", mcp.Required(), mcp.Description("Payment gateway id")),
		mcp.WithString("title", mcp.Description("Title")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithNumber("margin", mcp.Description("Margin percent")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "gateway_id", svc.UpdatePaymentGateway))

	s.AddTool(mcp.NewTool("delete_payment_gateway",
		mcp.WithDescription("Delete or remove payment gateway on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить способ оплаты"),
		mcp.WithString("gateway_id", mcp.Required(), mcp.Description("Payment gateway id")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeletePaymentGateway, "gateway_id", "Способ оплаты удалён"))

	// === СТРАНИЦЫ ===
	s.AddTool(mcp.NewTool("list_pages",
		mcp.WithDescription("List page from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список страниц магазина"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listPagesHandler(svc))

	s.AddTool(mcp.NewTool("get_page",
		mcp.WithDescription("Get a single page by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить страницу по ID"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("ID страницы")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getPageHandler(svc))

	s.AddTool(mcp.NewTool("create_page",
		mcp.WithDescription("Create or add a new page in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новую страницу в магазине"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название страницы")),
		mcp.WithString("content", mcp.Required(), mcp.Description("HTML-содержимое")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createPageHandler(svc))

	s.AddTool(mcp.NewTool("update_page",
		mcp.WithDescription("Update an existing page by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить данные страницы"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("ID страницы")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithString("content", mcp.Description("Содержимое")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updatePageHandler(svc))

	s.AddTool(mcp.NewTool("delete_page",
		mcp.WithDescription("Delete or remove page on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить страницу"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("ID страницы")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deletePageHandler(svc))

	// === АКСЕССУАРЫ ===
	s.AddTool(mcp.NewTool("list_accessories",
		mcp.WithDescription("List accessory from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список аксессуаров (характеристик товаров)"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listAccessoriesHandler(svc))

	s.AddTool(mcp.NewTool("get_accessory",
		mcp.WithDescription("Get a single accessory by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить аксессуар по ID"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getAccessoryHandler(svc))

	s.AddTool(mcp.NewTool("create_accessory",
		mcp.WithDescription("Create or add a new accessory in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новый аксессуар (характеристику)"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название (например 'Цвет', 'Размер')")),
		mcp.WithNumber("max_count", mcp.Description("Максимальное количество опций в заказе")),
		mcp.WithNumber("min_count", mcp.Description("Минимальное количество опций в заказе")),
		mcp.WithString("permalink", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createAccessoryHandler(svc))

	s.AddTool(mcp.NewTool("update_accessory",
		mcp.WithDescription("Update an existing accessory by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить аксессуар"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithNumber("position", mcp.Description("Позиция")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateAccessoryHandler(svc))

	s.AddTool(mcp.NewTool("delete_accessory",
		mcp.WithDescription("Delete or remove accessory on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить аксессуар"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteAccessoryHandler(svc))

	// === ВАРИАНТЫ ТОВАРОВ ===
	s.AddTool(mcp.NewTool("list_variants",
		mcp.WithDescription("List product variant from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список вариантов товара (размеры, цвета)"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listVariantsHandler(svc))

	s.AddTool(mcp.NewTool("get_variant",
		mcp.WithDescription("Get a single product variant by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить информацию о варианте товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("ID варианта")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getVariantHandler(svc))

	s.AddTool(mcp.NewTool("update_variant",
		mcp.WithDescription("Update an existing product variant by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить данные варианта товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("ID варианта")),
		mcp.WithString("price", mcp.Description("Цена")),
		mcp.WithString("sku", mcp.Description("Артикул")),
		mcp.WithNumber("quantity", mcp.Description("Количество")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateVariantHandler(svc))

	s.AddTool(mcp.NewTool("delete_variant",
		mcp.WithDescription("Delete or remove product variant on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить вариант товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("ID варианта")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteVariantHandler(svc))

	// === ЗНАЧЕНИЯ СВОЙСТВ ===
	s.AddTool(mcp.NewTool("list_property_values",
		mcp.WithDescription("List property value from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список значений аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listPropertyValuesHandler(svc))

	s.AddTool(mcp.NewTool("create_property_value",
		mcp.WithDescription("Create or add a new property value in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новое значение для аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название значения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createPropertyValueHandler(svc))

	s.AddTool(mcp.NewTool("update_property_value",
		mcp.WithDescription("Update an existing property value by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить значение аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("property_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updatePropertyValueHandler(svc))

	s.AddTool(mcp.NewTool("delete_property_value",
		mcp.WithDescription("Delete or remove property value on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить значение аксессуара"),
		mcp.WithString("accessory_id", mcp.Required(), mcp.Description("ID аксессуара")),
		mcp.WithString("property_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deletePropertyValueHandler(svc))

	// === МЕНЮ ===
	s.AddTool(mcp.NewTool("list_menus",
		mcp.WithDescription("List navigation menu from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список навигационных меню"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listMenusHandler(svc))

	s.AddTool(mcp.NewTool("get_menu",
		mcp.WithDescription("Get a single navigation menu by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить меню по ID со списком пунктов"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("ID меню")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getMenuHandler(svc))

	s.AddTool(mcp.NewTool("create_menu",
		mcp.WithDescription("Create or add a new navigation menu in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новое навигационное меню"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название меню")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createMenuHandler(svc))

	s.AddTool(mcp.NewTool("update_menu",
		mcp.WithDescription("Update an existing navigation menu by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить данные меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("ID меню")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateMenuHandler(svc))

	s.AddTool(mcp.NewTool("delete_menu",
		mcp.WithDescription("Delete or remove navigation menu on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("ID меню")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteMenuHandler(svc))

	// === MENU ITEMS (B2) ===
	s.AddTool(mcp.NewTool("list_menu_items",
		mcp.WithDescription("List navigation menu item from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Пункты меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("Menu id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedListHandler(func(ctx context.Context, parentID string) ([]any, error) {
		return svc.ListMenuItems(ctx, parentID)
	}, "menu_id"))

	s.AddTool(mcp.NewTool("get_menu_item",
		mcp.WithDescription("Get a single navigation menu item by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Пункт меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("Menu id")),
		mcp.WithString("item_id", mcp.Required(), mcp.Description("Menu item id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedGetHandler(func(ctx context.Context, parentID, key string) (any, error) {
		return svc.GetMenuItem(ctx, parentID, key)
	}, "menu_id", "item_id"))

	s.AddTool(mcp.NewTool("create_menu_item",
		mcp.WithDescription("Create or add a new navigation menu item in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Добавить пункт меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("Menu id")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Item title")),
		mcp.WithString("related_object_id", mcp.Required(), mcp.Description(`Composite ref: "Blog_1" | "Page_2" | "Collection_3" | https://...`)),
		mcp.WithNumber("position", mcp.Description("Sort position")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createMenuItemHandler(svc))

	s.AddTool(mcp.NewTool("update_menu_item",
		mcp.WithDescription("Update an existing navigation menu item by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить пункт меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("Menu id")),
		mcp.WithString("item_id", mcp.Required(), mcp.Description("Menu item id")),
		mcp.WithString("title", mcp.Description("Item title")),
		mcp.WithString("related_object_id", mcp.Description("Composite ref: Blog_1 | Page_2 | Collection_3 | URL")),
		mcp.WithNumber("position", mcp.Description("Sort position")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedUpdateHandler(func(ctx context.Context, parentID, key string, args map[string]any) (any, error) {
		return svc.UpdateMenuItem(ctx, parentID, key, args)
	}, "menu_id", "item_id"))

	s.AddTool(mcp.NewTool("delete_menu_item",
		mcp.WithDescription("Delete or remove navigation menu item on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить пункт меню"),
		mcp.WithString("menu_id", mcp.Required(), mcp.Description("Menu id")),
		mcp.WithString("item_id", mcp.Required(), mcp.Description("Menu item id")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedDeleteHandler(func(ctx context.Context, parentID, key string) error {
		return svc.DeleteMenuItem(ctx, parentID, key)
	}, "menu_id", "item_id", "Пункт меню удалён"))

	// === РЕДИРЕКТЫ ===
	s.AddTool(mcp.NewTool("list_redirects",
		mcp.WithDescription("List redirect from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список URL-редиректов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listRedirectsHandler(svc))

	s.AddTool(mcp.NewTool("create_redirect",
		mcp.WithDescription("Create or add a new redirect in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать URL-редирект"),
		mcp.WithString("url", mcp.Required(), mcp.Description("Старый URL")),
		mcp.WithString("target", mcp.Required(), mcp.Description("Новый URL")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createRedirectHandler(svc))

	s.AddTool(mcp.NewTool("delete_redirect",
		mcp.WithDescription("Delete or remove redirect on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить URL-редирект"),
		mcp.WithString("redirect_id", mcp.Required(), mcp.Description("ID редиректа")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteRedirectHandler(svc))

	// === БЛОГИ ===
	s.AddTool(mcp.NewTool("list_blogs",
		mcp.WithDescription("List blog from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список блогов магазина"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listBlogsHandler(svc))

	s.AddTool(mcp.NewTool("get_blog",
		mcp.WithDescription("Get a single blog by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить информацию о блоге по ID"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getBlogHandler(svc))

	s.AddTool(mcp.NewTool("create_blog",
		mcp.WithDescription("Create or add a new blog in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новый блог в магазине"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название блога")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createBlogHandler(svc))

	s.AddTool(mcp.NewTool("update_blog",
		mcp.WithDescription("Update an existing blog by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить блог"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("title", mcp.Description("Название блога")),
		mcp.WithString("handle", mcp.Description("URL-идентификатор")),
		mcp.WithString("commentable", mcp.Description("Комментарии: no, yes, moderable")),
		mcp.WithBoolean("notify_about_comment", mcp.Description("Уведомлять о комментариях")),
		mcp.WithBoolean("captcha_enabled", mcp.Description("Капча при комментировании")),
		mcp.WithBoolean("disallow_links_in_comments", mcp.Description("Запретить ссылки в комментариях")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateBlogHandler(svc))

	s.AddTool(mcp.NewTool("delete_blog",
		mcp.WithDescription("Delete or remove blog on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить блог"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteBlogHandler(svc))

	// === ЗАПИСИ БЛОГА ===
	s.AddTool(mcp.NewTool("list_blog_posts",
		mcp.WithDescription("List blog post from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список записей в блоге"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listBlogPostsHandler(svc))

	s.AddTool(mcp.NewTool("get_blog_post",
		mcp.WithDescription("Get a single blog post by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить запись блога по ID"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("ID записи")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getBlogPostHandler(svc))

	s.AddTool(mcp.NewTool("create_blog_post",
		mcp.WithDescription("Create or add a new blog post in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать новую запись в блоге"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Заголовок")),
		mcp.WithString("body_html", mcp.Required(), mcp.Description("HTML-содержимое")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createBlogPostHandler(svc))

	s.AddTool(mcp.NewTool("update_blog_post",
		mcp.WithDescription("Update an existing blog post by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить запись блога"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("ID записи")),
		mcp.WithString("title", mcp.Description("Заголовок")),
		mcp.WithString("body_html", mcp.Description("Содержимое")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateBlogPostHandler(svc))

	s.AddTool(mcp.NewTool("delete_blog_post",
		mcp.WithDescription("Delete or remove blog post on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить запись блога"),
		mcp.WithString("blog_id", mcp.Required(), mcp.Description("ID блога")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("ID записи")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteBlogPostHandler(svc))

	// === ПРОМОКОДЫ (Discount Codes) ===
	s.AddTool(mcp.NewTool("list_discount_codes",
		mcp.WithDescription("List discount code from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список промокодов"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listDiscountCodesHandler(svc))

	s.AddTool(mcp.NewTool("get_discount_code",
		mcp.WithDescription("Get a single discount code by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить промокод по ID"),
		mcp.WithString("discount_code_id", mcp.Required(), mcp.Description("ID промокода")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getDiscountCodeHandler(svc))

	s.AddTool(mcp.NewTool("create_discount_code",
		mcp.WithDescription("Create or add a new discount code in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать промокод"),
		mcp.WithString("code", mcp.Required(), mcp.Description("Код промокода")),
		mcp.WithNumber("type_id", mcp.Required(), mcp.Description("Тип: 1=процент, 2=фиксированная сумма")),
		mcp.WithNumber("discount", mcp.Required(), mcp.Description("Скидка")),
		mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithString("expired_at", mcp.Description("Дата окончания (YYYY-MM-DD)")),
		mcp.WithNumber("min_price", mcp.Description("Минимальная сумма заказа")),
		mcp.WithBoolean("act_once", mcp.Description("Одноразовый")),
		mcp.WithBoolean("act_once_for_client", mcp.Description("Одноразовый для клиента")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createDiscountCodeHandler(svc))

	s.AddTool(mcp.NewTool("update_discount_code",
		mcp.WithDescription("Update an existing discount code by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить промокод"),
		mcp.WithString("discount_code_id", mcp.Required(), mcp.Description("ID промокода")),
		mcp.WithString("code", mcp.Description("Код промокода")),
		mcp.WithNumber("discount", mcp.Description("Скидка")),
		mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithBoolean("disabled", mcp.Description("Отключить")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateDiscountCodeHandler(svc))

	s.AddTool(mcp.NewTool("delete_discount_code",
		mcp.WithDescription("Delete or remove discount code on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить промокод"),
		mcp.WithString("discount_code_id", mcp.Required(), mcp.Description("ID промокода")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteDiscountCodeHandler(svc))

	// === СВЯЗИ ТОВАР-КОЛЛЕКЦИЯ (Collects) ===
	s.AddTool(mcp.NewTool("list_collects",
		mcp.WithDescription("List product-collection link from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список связей товар-коллекция"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listCollectsHandler(svc))

	s.AddTool(mcp.NewTool("get_collects_by_product",
		mcp.WithDescription("Get a single product-collection link by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить коллекции товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getCollectsByProductHandler(svc))

	s.AddTool(mcp.NewTool("add_product_to_collection",
		mcp.WithDescription("Create or add a new collection in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Добавить товар в коллекцию"),
		mcp.WithNumber("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithNumber("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createCollectHandler(svc))

	s.AddTool(mcp.NewTool("remove_product_from_collection",
		mcp.WithDescription("Remove a product from a collection (collect link). Destructive on the link only."),
		mcp.WithToolTitle("Удалить товар из коллекции"),
		mcp.WithString("collect_id", mcp.Required(), mcp.Description("Collect link id")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteCollectHandler(svc))

	s.AddTool(mcp.NewTool("update_collect",
		mcp.WithDescription("Move a product to another collection or change position in the list. Uses collect link id."),
		mcp.WithToolTitle("Обновить связь товар-коллекция"),
		mcp.WithString("collect_id", mcp.Required(), mcp.Description("Collect link id")),
		mcp.WithNumber("collection_id", mcp.Description("Target collection id (move)")),
		mcp.WithNumber("position", mcp.Description("New position")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateCollectHandler(svc))

	// === СВОЙСТВА ТОВАРОВ (Properties) ===
	s.AddTool(mcp.NewTool("list_properties",
		mcp.WithDescription("List property from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список свойств товаров (Бренд, Цвет...)"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listPropertiesHandler(svc))

	s.AddTool(mcp.NewTool("get_property",
		mcp.WithDescription("Get a single property by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить свойство по ID"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getPropertyHandler(svc))

	s.AddTool(mcp.NewTool("create_property",
		mcp.WithDescription("Create or add a new property in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать свойство товаров"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название свойства")),
		mcp.WithString("permalink", mcp.Description("URL-идентификатор")),
		mcp.WithBoolean("is_navigational", mcp.Description("Показывать в фильтре")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createPropertyHandler(svc))

	s.AddTool(mcp.NewTool("update_property",
		mcp.WithDescription("Update an existing property by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить свойство товаров"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithBoolean("is_navigational", mcp.Description("Показывать в фильтре")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updatePropertyHandler(svc))

	s.AddTool(mcp.NewTool("delete_property",
		mcp.WithDescription("Delete or remove property on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить свойство товаров"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deletePropertyHandler(svc))

	// === ХАРАКТЕРИСТИКИ (Characteristics) ===
	s.AddTool(mcp.NewTool("list_characteristics",
		mcp.WithDescription("List characteristic from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список характеристик свойства (Nike, Adidas...)"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listCharacteristicsHandler(svc))

	s.AddTool(mcp.NewTool("get_characteristic",
		mcp.WithDescription("Get a single characteristic by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить характеристику по ID"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("characteristic_id", mcp.Required(), mcp.Description("ID характеристики")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getCharacteristicHandler(svc))

	s.AddTool(mcp.NewTool("create_characteristic",
		mcp.WithDescription("Create or add a new characteristic in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать характеристику свойства"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Название характеристики")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createCharacteristicHandler(svc))

	s.AddTool(mcp.NewTool("update_characteristic",
		mcp.WithDescription("Update an existing characteristic by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить характеристику свойства"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("characteristic_id", mcp.Required(), mcp.Description("ID характеристики")),
		mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateCharacteristicHandler(svc))

	s.AddTool(mcp.NewTool("delete_characteristic",
		mcp.WithDescription("Delete or remove characteristic on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить характеристику свойства"),
		mcp.WithString("property_id", mcp.Required(), mcp.Description("ID свойства")),
		mcp.WithString("characteristic_id", mcp.Required(), mcp.Description("ID характеристики")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteCharacteristicHandler(svc))

	// === ОТЗЫВЫ (Reviews) ===
	s.AddTool(mcp.NewTool("list_reviews",
		mcp.WithDescription("List review from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список отзывов"),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listReviewsHandler(svc))

	s.AddTool(mcp.NewTool("get_review",
		mcp.WithDescription("Get a single review by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить отзыв по ID"),
		mcp.WithString("review_id", mcp.Required(), mcp.Description("ID отзыва")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getReviewHandler(svc))

	s.AddTool(mcp.NewTool("create_review",
		mcp.WithDescription("Create or add a new review in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать отзыв"),
		mcp.WithNumber("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithString("author", mcp.Required(), mcp.Description("Автор")),
		mcp.WithString("content", mcp.Description("Текст отзыва")),
		mcp.WithString("email", mcp.Description("Email")),
		mcp.WithNumber("rating", mcp.Description("Оценка (1-10)")),
		mcp.WithBoolean("published", mcp.Description("Опубликован")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createReviewHandler(svc))

	s.AddTool(mcp.NewTool("update_review",
		mcp.WithDescription("Update an existing review by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить отзыв"),
		mcp.WithString("review_id", mcp.Required(), mcp.Description("ID отзыва")),
		mcp.WithString("content", mcp.Description("Текст отзыва")),
		mcp.WithNumber("rating", mcp.Description("Оценка")),
		mcp.WithBoolean("published", mcp.Description("Опубликован")),
		mcp.WithString("manager_reply", mcp.Description("Ответ менеджера")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateReviewHandler(svc))

	s.AddTool(mcp.NewTool("delete_review",
		mcp.WithDescription("Delete or remove review on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить отзыв"),
		mcp.WithString("review_id", mcp.Required(), mcp.Description("ID отзыва")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteReviewHandler(svc))

	// === ФАЙЛЫ (Files) ===
	s.AddTool(mcp.NewTool("list_files",
		mcp.WithDescription("List file from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Получить список файлов магазина"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listFilesHandler(svc))

	s.AddTool(mcp.NewTool("get_file",
		mcp.WithDescription("Get a single file by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить файл по ID"),
		mcp.WithString("file_id", mcp.Required(), mcp.Description("ID файла")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getFileHandler(svc))

	s.AddTool(mcp.NewTool("create_file",
		mcp.WithDescription("Create or add a new file in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Загрузить файл (по URL)"),
		mcp.WithString("src", mcp.Required(), mcp.Description("URL файла")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createFileHandler(svc))

	s.AddTool(mcp.NewTool("delete_file",
		mcp.WithDescription("Delete or remove file on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить файл"),
		mcp.WithString("file_id", mcp.Required(), mcp.Description("ID файла")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), deleteFileHandler(svc))

	// === DOMAINS ===
	s.AddTool(mcp.NewTool("list_domains", mcp.WithDescription("List domain from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список доменов магазина"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListDomains))
	s.AddTool(mcp.NewTool("get_domain", mcp.WithDescription("Get a single domain by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить домен по ID"), mcp.WithString("domain_id", mcp.Required(), mcp.Description("ID домена")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "domain_id", svc.GetDomain))
	s.AddTool(mcp.NewTool("create_domain", mcp.WithDescription("Create or add a new domain in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Добавить домен"), mcp.WithString("domain", mcp.Required(), mcp.Description("Домен")), mcp.WithBoolean("main", mcp.Description("Основной домен")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateDomain))
	s.AddTool(mcp.NewTool("update_domain", mcp.WithDescription("Update an existing domain by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить домен"), mcp.WithString("domain_id", mcp.Required(), mcp.Description("ID домена")), mcp.WithBoolean("main", mcp.Description("Основной домен")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "domain_id", svc.UpdateDomain))
	s.AddTool(mcp.NewTool("delete_domain", mcp.WithDescription("Delete or remove domain on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить домен"), mcp.WithString("domain_id", mcp.Required(), mcp.Description("ID домена")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteDomain, "domain_id", "Домен удалён"))

	// === CLIENT GROUPS ===
	s.AddTool(mcp.NewTool("list_client_groups", mcp.WithDescription("List client group from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список групп клиентов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListClientGroups))
	s.AddTool(mcp.NewTool("get_client_group", mcp.WithDescription("Get a single client group by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить группу клиентов"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "group_id", svc.GetClientGroup))
	s.AddTool(mcp.NewTool("create_client_group", mcp.WithDescription("Create or add a new client group in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать группу клиентов"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithNumber("discount", mcp.Description("Скидка (%)")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateClientGroup))
	s.AddTool(mcp.NewTool("update_client_group", mcp.WithDescription("Update an existing client group by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить группу клиентов"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "group_id", svc.UpdateClientGroup))
	s.AddTool(mcp.NewTool("delete_client_group", mcp.WithDescription("Delete or remove client group on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить группу клиентов"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteClientGroup, "group_id", "Группа удалена"))

	// === CUSTOM STATUSES ===
	s.AddTool(mcp.NewTool("list_custom_statuses", mcp.WithDescription("List custom order status from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список кастомных статусов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListCustomStatuses))
	s.AddTool(mcp.NewTool("get_custom_status", mcp.WithDescription("Get a single custom order status by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить кастомный статус"), mcp.WithString("status_id", mcp.Required(), mcp.Description("ID статуса")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "status_id", svc.GetCustomStatus))
	s.AddTool(mcp.NewTool("create_custom_status", mcp.WithDescription("Create or add a new custom order status in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать кастомный статус"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("status_type", mcp.Required(), mcp.Description("Тип статуса")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateCustomStatus))
	s.AddTool(mcp.NewTool("update_custom_status", mcp.WithDescription("Update an existing custom order status by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить кастомный статус"), mcp.WithString("status_id", mcp.Required(), mcp.Description("ID статуса")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "status_id", svc.UpdateCustomStatus))
	s.AddTool(mcp.NewTool("delete_custom_status", mcp.WithDescription("Delete or remove custom order status on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить кастомный статус"), mcp.WithString("status_id", mcp.Required(), mcp.Description("ID статуса")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteCustomStatus, "status_id", "Статус удалён"))

	// === PRODUCT FIELDS ===
	s.AddTool(mcp.NewTool("list_product_fields",
		mcp.WithDescription("List product custom field from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Список полей товаров"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListProductFields))
	s.AddTool(mcp.NewTool("get_product_field",
		mcp.WithDescription("Get a single product custom field by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Поле товара"),
		mcp.WithString("field_id", mcp.Required(), mcp.Description("Field id or handle")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "field_id", svc.GetProductField))
	s.AddTool(mcp.NewTool("create_product_field",
		mcp.WithDescription("Create or add a new product custom field in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Создать поле товара"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Display title")),
		mcp.WithString("handle", mcp.Required(), mcp.Description("API handle/slug, e.g. dosage — language-invariant")),
		mcp.WithString("type", mcp.Description("ProductField::TextField | ProductField::TextArea")),
		mcp.WithString("field_type", mcp.Description("Alias for type (legacy)")),
		mcp.WithBoolean("is_hidden", mcp.Description("Hide in admin UI")),
		mcp.WithBoolean("is_indexed", mcp.Description("Include in search index")),
		mcp.WithNumber("position", mcp.Description("Sort position")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createProductFieldHandler(svc))
	s.AddTool(mcp.NewTool("update_product_field",
		mcp.WithDescription("Update an existing product custom field by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить поле товара"),
		mcp.WithString("field_id", mcp.Required(), mcp.Description("Field id or handle")),
		mcp.WithString("title", mcp.Description("Display title")),
		mcp.WithString("handle", mcp.Description("API handle/slug")),
		mcp.WithString("type", mcp.Description("ProductField::TextField | ProductField::TextArea")),
		mcp.WithBoolean("is_hidden", mcp.Description("Hide in admin UI")),
		mcp.WithBoolean("is_indexed", mcp.Description("Include in search index")),
		mcp.WithNumber("position", mcp.Description("Sort position")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateProductFieldHandler(svc))
	s.AddTool(mcp.NewTool("delete_product_field",
		mcp.WithDescription("Delete or remove product custom field on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить поле товара"),
		mcp.WithString("field_id", mcp.Required(), mcp.Description("Field id or handle")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteProductField, "field_id", "Поле удалено"))

	// === PRODUCT FIELD VALUES (custom field values on a product) ===
	s.AddTool(mcp.NewTool("list_product_field_values",
		mcp.WithDescription("List product custom field value from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Значения полей товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("Product id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedListHandler(func(ctx context.Context, parentID string) ([]any, error) {
		return svc.ListProductFieldValues(ctx, parentID)
	}, "product_id"))
	s.AddTool(mcp.NewTool("get_product_field_value",
		mcp.WithDescription("Get a single product custom field value by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Значение поля товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("Product id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedGetHandler(func(ctx context.Context, parentID, key string) (any, error) {
		return svc.GetProductFieldValue(ctx, parentID, key)
	}, "product_id", "value_id"))
	s.AddTool(mcp.NewTool("create_product_field_value",
		mcp.WithDescription("Create or add a new product custom field value in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Установить значение поля товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("Product id")),
		mcp.WithString("product_field_id", mcp.Description("Field definition id")),
		mcp.WithString("product_field_handle", mcp.Description("Field handle, e.g. dosage")),
		mcp.WithString("value", mcp.Description("Value text")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createProductFieldValueHandler(svc))
	s.AddTool(mcp.NewTool("update_product_field_value",
		mcp.WithDescription("Update an existing product custom field value by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить значение поля товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("Product id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithString("value", mcp.Description("New value text")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedUpdateHandler(func(ctx context.Context, parentID, key string, args map[string]any) (any, error) {
		return svc.UpdateProductFieldValue(ctx, parentID, key, args)
	}, "product_id", "value_id"))
	s.AddTool(mcp.NewTool("delete_product_field_value",
		mcp.WithDescription("Delete or remove product custom field value on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить значение поля товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("Product id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedDeleteHandler(func(ctx context.Context, parentID, key string) error {
		return svc.DeleteProductFieldValue(ctx, parentID, key)
	}, "product_id", "value_id", "Значение удалено"))

	s.AddTool(mcp.NewTool("list_product_accessory_links",
		mcp.WithDescription("List accessory from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Связи аксессуаров товара"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("Product id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedListHandler(func(ctx context.Context, parentID string) ([]any, error) {
		return svc.ListProductAccessoryLinks(ctx, parentID)
	}, "product_id"))

	// === VARIANT FIELD VALUES ===
	s.AddTool(mcp.NewTool("list_variant_field_values",
		mcp.WithDescription("List variant custom field value from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Значения полей варианта"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Variant id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedListHandler(func(ctx context.Context, parentID string) ([]any, error) {
		return svc.ListVariantFieldValues(ctx, parentID)
	}, "variant_id"))
	s.AddTool(mcp.NewTool("get_variant_field_value",
		mcp.WithDescription("Get a single variant custom field value by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Значение поля варианта"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Variant id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedGetHandler(func(ctx context.Context, parentID, key string) (any, error) {
		return svc.GetVariantFieldValue(ctx, parentID, key)
	}, "variant_id", "value_id"))
	s.AddTool(mcp.NewTool("create_variant_field_value",
		mcp.WithDescription("Create or add a new variant custom field value in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."),
		mcp.WithToolTitle("Установить значение поля варианта"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Variant id")),
		mcp.WithString("variant_field_id", mcp.Description("Field definition id")),
		mcp.WithString("variant_field_handle", mcp.Description("Field handle")),
		mcp.WithString("value", mcp.Description("Value text")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createVariantFieldValueHandler(svc))
	s.AddTool(mcp.NewTool("update_variant_field_value",
		mcp.WithDescription("Update an existing variant custom field value by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить значение поля варианта"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Variant id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithString("value", mcp.Description("New value text")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedUpdateHandler(func(ctx context.Context, parentID, key string, args map[string]any) (any, error) {
		return svc.UpdateVariantFieldValue(ctx, parentID, key, args)
	}, "variant_id", "value_id"))
	s.AddTool(mcp.NewTool("delete_variant_field_value",
		mcp.WithDescription("Delete or remove variant custom field value on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."),
		mcp.WithToolTitle("Удалить значение поля варианта"),
		mcp.WithString("variant_id", mcp.Required(), mcp.Description("Variant id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedDeleteHandler(func(ctx context.Context, parentID, key string) error {
		return svc.DeleteVariantFieldValue(ctx, parentID, key)
	}, "variant_id", "value_id", "Значение удалено"))

	// === COLLECTION FIELD VALUES (GET/PUT only) ===
	s.AddTool(mcp.NewTool("list_collection_field_values",
		mcp.WithDescription("List collection custom field value from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."),
		mcp.WithToolTitle("Значения полей коллекции"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("Collection id")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedListHandler(func(ctx context.Context, parentID string) ([]any, error) {
		return svc.ListCollectionFieldValues(ctx, parentID)
	}, "collection_id"))
	s.AddTool(mcp.NewTool("get_collection_field_value",
		mcp.WithDescription("Get a single collection custom field value by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."),
		mcp.WithToolTitle("Значение поля коллекции"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("Collection id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedGetHandler(func(ctx context.Context, parentID, key string) (any, error) {
		return svc.GetCollectionFieldValue(ctx, parentID, key)
	}, "collection_id", "value_id"))
	s.AddTool(mcp.NewTool("update_collection_field_value",
		mcp.WithDescription("Update an existing collection custom field value by id. Only send fields that should change. Returns the updated object as JSON."),
		mcp.WithToolTitle("Обновить значение поля коллекции"),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("Collection id")),
		mcp.WithString("value_id", mcp.Required(), mcp.Description("Field value id OR field handle")),
		mcp.WithString("value", mcp.Description("New value text")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), nestedUpdateHandler(func(ctx context.Context, parentID, key string, args map[string]any) (any, error) {
		return svc.UpdateCollectionFieldValue(ctx, parentID, key, args)
	}, "collection_id", "value_id"))

	// === VARIANT FIELDS ===
	s.AddTool(mcp.NewTool("list_variant_fields", mcp.WithDescription("List variant custom field from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список полей вариантов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListVariantFields))
	s.AddTool(mcp.NewTool("get_variant_field", mcp.WithDescription("Get a single variant custom field by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить поле варианта"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "field_id", svc.GetVariantField))
	s.AddTool(mcp.NewTool("create_variant_field", mcp.WithDescription("Create or add a new variant custom field in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать поле варианта"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("field_type", mcp.Required(), mcp.Description("Тип поля")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateVariantField))
	s.AddTool(mcp.NewTool("update_variant_field", mcp.WithDescription("Update an existing variant custom field by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить поле варианта"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "field_id", svc.UpdateVariantField))
	s.AddTool(mcp.NewTool("delete_variant_field", mcp.WithDescription("Delete or remove variant custom field on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить поле варианта"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteVariantField, "field_id", "Поле удалено"))

	// === COLLECTION FIELDS ===
	s.AddTool(mcp.NewTool("list_collection_fields", mcp.WithDescription("List collection custom field from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список полей коллекций"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListCollectionFields))
	s.AddTool(mcp.NewTool("get_collection_field", mcp.WithDescription("Get a single collection custom field by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить поле коллекции"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "field_id", svc.GetCollectionField))
	s.AddTool(mcp.NewTool("create_collection_field", mcp.WithDescription("Create or add a new collection custom field in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать поле коллекции"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("field_type", mcp.Required(), mcp.Description("Тип поля")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateCollectionField))
	s.AddTool(mcp.NewTool("update_collection_field", mcp.WithDescription("Update an existing collection custom field by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить поле коллекции"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "field_id", svc.UpdateCollectionField))
	s.AddTool(mcp.NewTool("delete_collection_field", mcp.WithDescription("Delete or remove collection custom field on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить поле коллекции"), mcp.WithString("field_id", mcp.Required(), mcp.Description("ID поля")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteCollectionField, "field_id", "Поле удалено"))

	// === COLLECTION FILTERS ===
	s.AddTool(mcp.NewTool("list_collection_filters", mcp.WithDescription("List collection filter from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список фильтров коллекций"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListCollectionFilters))
	s.AddTool(mcp.NewTool("get_collection_filter", mcp.WithDescription("Get a single collection filter by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить фильтр коллекции"), mcp.WithString("filter_id", mcp.Required(), mcp.Description("ID фильтра")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "filter_id", svc.GetCollectionFilter))
	s.AddTool(mcp.NewTool("create_collection_filter", mcp.WithDescription("Create or add a new collection filter in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать фильтр коллекции"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("filter_type", mcp.Required(), mcp.Description("Тип фильтра")), mcp.WithNumber("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateCollectionFilter))
	s.AddTool(mcp.NewTool("update_collection_filter", mcp.WithDescription("Update an existing collection filter by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить фильтр коллекции"), mcp.WithString("filter_id", mcp.Required(), mcp.Description("ID фильтра")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "filter_id", svc.UpdateCollectionFilter))
	s.AddTool(mcp.NewTool("delete_collection_filter", mcp.WithDescription("Delete or remove collection filter on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить фильтр коллекции"), mcp.WithString("filter_id", mcp.Required(), mcp.Description("ID фильтра")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteCollectionFilter, "filter_id", "Фильтр удалён"))

	// === JS TAGS ===
	s.AddTool(mcp.NewTool("list_js_tags", mcp.WithDescription("List JS tag from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список JS-тегов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListJsTags))
	s.AddTool(mcp.NewTool("get_js_tag", mcp.WithDescription("Get a single JS tag by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить JS-тег"), mcp.WithString("tag_id", mcp.Required(), mcp.Description("ID тега")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "tag_id", svc.GetJsTag))
	s.AddTool(mcp.NewTool("create_js_tag", mcp.WithDescription("Create or add a new JS tag in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать JS-тег"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("content", mcp.Required(), mcp.Description("Код JS")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateJsTag))
	s.AddTool(mcp.NewTool("update_js_tag", mcp.WithDescription("Update an existing JS tag by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить JS-тег"), mcp.WithString("tag_id", mcp.Required(), mcp.Description("ID тега")), mcp.WithString("title", mcp.Description("Название")), mcp.WithString("content", mcp.Description("Код JS")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "tag_id", svc.UpdateJsTag))
	s.AddTool(mcp.NewTool("delete_js_tag", mcp.WithDescription("Delete or remove JS tag on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить JS-тег"), mcp.WithString("tag_id", mcp.Required(), mcp.Description("ID тега")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteJsTag, "tag_id", "Тег удалён"))

	// === EXTERNAL DISCOUNTS ===
	s.AddTool(mcp.NewTool("list_external_discounts", mcp.WithDescription("List external discount from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список внешних скидок"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListExternalDiscounts))
	s.AddTool(mcp.NewTool("get_external_discount", mcp.WithDescription("Get a single external discount by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить внешнюю скидку"), mcp.WithString("discount_id", mcp.Required(), mcp.Description("ID скидки")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "discount_id", svc.GetExternalDiscount))
	s.AddTool(mcp.NewTool("create_external_discount", mcp.WithDescription("Create or add a new external discount in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать внешнюю скидку"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithNumber("percent", mcp.Description("Процент скидки")), mcp.WithNumber("amount", mcp.Description("Сумма скидки")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateExternalDiscount))
	s.AddTool(mcp.NewTool("update_external_discount", mcp.WithDescription("Update an existing external discount by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить внешнюю скидку"), mcp.WithString("discount_id", mcp.Required(), mcp.Description("ID скидки")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "discount_id", svc.UpdateExternalDiscount))
	s.AddTool(mcp.NewTool("delete_external_discount", mcp.WithDescription("Delete or remove external discount on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить внешнюю скидку"), mcp.WithString("discount_id", mcp.Required(), mcp.Description("ID скидки")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteExternalDiscount, "discount_id", "Скидка удалена"))

	// === DISCOUNT COLLECTIONS GROUPS ===
	s.AddTool(mcp.NewTool("list_discount_collections_groups", mcp.WithDescription("List discount collections group from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список групп коллекций скидок"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListDiscountCollectionsGroups))
	s.AddTool(mcp.NewTool("get_discount_collections_group", mcp.WithDescription("Get a single discount collections group by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить группу коллекций скидок"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "group_id", svc.GetDiscountCollectionsGroup))
	s.AddTool(mcp.NewTool("create_discount_collections_group", mcp.WithDescription("Create or add a new discount collections group in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать группу коллекций скидок"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateDiscountCollectionsGroup))
	s.AddTool(mcp.NewTool("update_discount_collections_group", mcp.WithDescription("Update an existing discount collections group by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить группу коллекций скидок"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "group_id", svc.UpdateDiscountCollectionsGroup))
	s.AddTool(mcp.NewTool("delete_discount_collections_group", mcp.WithDescription("Delete or remove discount collections group on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить группу коллекций скидок"), mcp.WithString("group_id", mcp.Required(), mcp.Description("ID группы")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteDiscountCollectionsGroup, "group_id", "Группа удалена"))

	// === WAREHOUSES ===
	s.AddTool(mcp.NewTool("list_warehouses", mcp.WithDescription("List warehouse from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список складов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListWarehouses))
	s.AddTool(mcp.NewTool("get_warehouse", mcp.WithDescription("Get a single warehouse by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить склад"), mcp.WithString("warehouse_id", mcp.Required(), mcp.Description("ID склада")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "warehouse_id", svc.GetWarehouse))
	s.AddTool(mcp.NewTool("create_warehouse", mcp.WithDescription("Create or add a new warehouse in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать склад"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("address", mcp.Description("Адрес")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateWarehouse))
	s.AddTool(mcp.NewTool("update_warehouse", mcp.WithDescription("Update an existing warehouse by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить склад"), mcp.WithString("warehouse_id", mcp.Required(), mcp.Description("ID склада")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "warehouse_id", svc.UpdateWarehouse))
	s.AddTool(mcp.NewTool("delete_warehouse", mcp.WithDescription("Delete or remove warehouse on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить склад"), mcp.WithString("warehouse_id", mcp.Required(), mcp.Description("ID склада")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteWarehouse, "warehouse_id", "Склад удалён"))

	// === SALES CHANNELS ===
	s.AddTool(mcp.NewTool("list_sales_channels", mcp.WithDescription("List sales channel from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список каналов продаж"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListSalesChannels))
	s.AddTool(mcp.NewTool("get_sales_channel", mcp.WithDescription("Get a single sales channel by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить канал продаж"), mcp.WithString("channel_id", mcp.Required(), mcp.Description("ID канала")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "channel_id", svc.GetSalesChannel))
	s.AddTool(mcp.NewTool("create_sales_channel", mcp.WithDescription("Create or add a new sales channel in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать канал продаж"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateSalesChannel))
	s.AddTool(mcp.NewTool("update_sales_channel", mcp.WithDescription("Update an existing sales channel by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить канал продаж"), mcp.WithString("channel_id", mcp.Required(), mcp.Description("ID канала")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "channel_id", svc.UpdateSalesChannel))
	s.AddTool(mcp.NewTool("delete_sales_channel", mcp.WithDescription("Delete or remove sales channel on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить канал продаж"), mcp.WithString("channel_id", mcp.Required(), mcp.Description("ID канала")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteSalesChannel, "channel_id", "Канал удалён"))

	// === STOCK CURRENCIES ===
	s.AddTool(mcp.NewTool("list_stock_currencies", mcp.WithDescription("List stock currency from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список валют склада"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListStockCurrencies))
	s.AddTool(mcp.NewTool("get_stock_currency", mcp.WithDescription("Get a single stock currency by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить валюту склада"), mcp.WithString("currency_id", mcp.Required(), mcp.Description("ID валюты")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "currency_id", svc.GetStockCurrency))
	s.AddTool(mcp.NewTool("create_stock_currency", mcp.WithDescription("Create or add a new stock currency in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать валюту склада"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("code", mcp.Required(), mcp.Description("Код валюты")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateStockCurrency))
	s.AddTool(mcp.NewTool("update_stock_currency", mcp.WithDescription("Update an existing stock currency by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить валюту склада"), mcp.WithString("currency_id", mcp.Required(), mcp.Description("ID валюты")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "currency_id", svc.UpdateStockCurrency))
	s.AddTool(mcp.NewTool("delete_stock_currency", mcp.WithDescription("Delete or remove stock currency on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить валюту склада"), mcp.WithString("currency_id", mcp.Required(), mcp.Description("ID валюты")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteStockCurrency, "currency_id", "Валюта удалена"))

	// === PRICE KINDS ===
	s.AddTool(mcp.NewTool("list_price_kinds", mcp.WithDescription("List price kind from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список типов цен"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListPriceKinds))
	s.AddTool(mcp.NewTool("get_price_kind", mcp.WithDescription("Get a single price kind by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить тип цены"), mcp.WithString("kind_id", mcp.Required(), mcp.Description("ID типа")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "kind_id", svc.GetPriceKind))
	s.AddTool(mcp.NewTool("create_price_kind", mcp.WithDescription("Create or add a new price kind in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать тип цены"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreatePriceKind))
	s.AddTool(mcp.NewTool("update_price_kind", mcp.WithDescription("Update an existing price kind by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить тип цены"), mcp.WithString("kind_id", mcp.Required(), mcp.Description("ID типа")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "kind_id", svc.UpdatePriceKind))
	s.AddTool(mcp.NewTool("delete_price_kind", mcp.WithDescription("Delete or remove price kind on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить тип цены"), mcp.WithString("kind_id", mcp.Required(), mcp.Description("ID типа")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeletePriceKind, "kind_id", "Тип цены удалён"))

	// === OPTION NAMES ===
	s.AddTool(mcp.NewTool("list_option_names", mcp.WithDescription("List option name from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список имён опций (Размер, Цвет...)"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListOptionNames))
	s.AddTool(mcp.NewTool("get_option_name", mcp.WithDescription("Get a single option name by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить имя опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "option_name_id", svc.GetOptionName))
	s.AddTool(mcp.NewTool("create_option_name", mcp.WithDescription("Create or add a new option name in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать имя опции"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithBoolean("navigational", mcp.Description("Показывать в фильтре")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateOptionName))
	s.AddTool(mcp.NewTool("update_option_name", mcp.WithDescription("Update an existing option name by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить имя опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "option_name_id", svc.UpdateOptionName))
	s.AddTool(mcp.NewTool("delete_option_name", mcp.WithDescription("Delete or remove option name on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить имя опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteOptionName, "option_name_id", "Имя опции удалено"))

	// === OPTION VALUES ===
	s.AddTool(mcp.NewTool("list_option_values", mcp.WithDescription("List option value from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список значений опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListOptionValuesHandler(svc))
	s.AddTool(mcp.NewTool("get_option_value", mcp.WithDescription("Get a single option value by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("option_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetOptionValueHandler(svc))
	s.AddTool(mcp.NewTool("create_option_value", mcp.WithDescription("Create or add a new option value in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateOptionValueHandler(svc))
	s.AddTool(mcp.NewTool("update_option_value", mcp.WithDescription("Update an existing option value by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("option_value_id", mcp.Required(), mcp.Description("ID значения")), mcp.WithString("title", mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateOptionValueHandler(svc))
	s.AddTool(mcp.NewTool("delete_option_value", mcp.WithDescription("Delete or remove option value on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить значение опции"), mcp.WithString("option_name_id", mcp.Required(), mcp.Description("ID имени опции")), mcp.WithString("option_value_id", mcp.Required(), mcp.Description("ID значения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteOptionValueHandler(svc))

	// === ROUND 4: APPLICATION ACTIONS ===
	s.AddTool(mcp.NewTool("list_application_actions", mcp.WithDescription("List application action from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список действий приложения"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListApplicationActions))
	s.AddTool(mcp.NewTool("get_application_action", mcp.WithDescription("Get a single application action by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить действие приложения"), mcp.WithString("action_id", mcp.Required(), mcp.Description("ID действия")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "action_id", svc.GetApplicationAction))
	s.AddTool(mcp.NewTool("create_application_action", mcp.WithDescription("Create or add a new application action in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать действие приложения"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("handle", mcp.Required(), mcp.Description("Идентификатор")), mcp.WithString("entity", mcp.Required(), mcp.Description("Сущность: order, product, client...")), mcp.WithString("url", mcp.Required(), mcp.Description("URL")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateApplicationAction))
	s.AddTool(mcp.NewTool("update_application_action", mcp.WithDescription("Update an existing application action by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить действие приложения"), mcp.WithString("action_id", mcp.Required(), mcp.Description("ID действия")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "action_id", svc.UpdateApplicationAction))
	s.AddTool(mcp.NewTool("delete_application_action", mcp.WithDescription("Delete or remove application action on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить действие приложения"), mcp.WithString("action_id", mcp.Required(), mcp.Description("ID действия")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteApplicationAction, "action_id", "Действие удалено"))

	// === APPLICATION CHARGES ===
	s.AddTool(mcp.NewTool("list_application_charges", mcp.WithDescription("List application charge from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список счетов приложения"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListApplicationCharges))
	s.AddTool(mcp.NewTool("get_application_charge", mcp.WithDescription("Get a single application charge by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить счёт приложения"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "charge_id", svc.GetApplicationCharge))
	s.AddTool(mcp.NewTool("create_application_charge", mcp.WithDescription("Create or add a new application charge in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать счёт приложения"), mcp.WithString("name", mcp.Required(), mcp.Description("Назначение")), mcp.WithNumber("price", mcp.Required(), mcp.Description("Сумма")), mcp.WithString("return_url", mcp.Description("URL уведомления")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateApplicationCharge))
	s.AddTool(mcp.NewTool("decline_application_charge", mcp.WithDescription("Delete or remove application charge on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Отклонить счёт приложения"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), declineApplicationChargeHandler(svc))

	// === APPLICATION WIDGETS ===
	s.AddTool(mcp.NewTool("list_application_widgets", mcp.WithDescription("List application widget from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список виджетов приложения"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListApplicationWidgets))
	s.AddTool(mcp.NewTool("get_application_widget", mcp.WithDescription("Get a single application widget by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить виджет"), mcp.WithString("widget_id", mcp.Required(), mcp.Description("ID виджета")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "widget_id", svc.GetApplicationWidget))
	s.AddTool(mcp.NewTool("create_application_widget", mcp.WithDescription("Create or add a new application widget in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать виджет"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateApplicationWidget))
	s.AddTool(mcp.NewTool("update_application_widget", mcp.WithDescription("Update an existing application widget by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить виджет"), mcp.WithString("widget_id", mcp.Required(), mcp.Description("ID виджета")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "widget_id", svc.UpdateApplicationWidget))
	s.AddTool(mcp.NewTool("delete_application_widget", mcp.WithDescription("Delete or remove application widget on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить виджет"), mcp.WithString("widget_id", mcp.Required(), mcp.Description("ID виджета")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteApplicationWidget, "widget_id", "Виджет удалён"))

	// === RECURRING APPLICATION CHARGES ===
	s.AddTool(mcp.NewTool("list_recurring_application_charges", mcp.WithDescription("List recurring application charge from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список рекуррентных счетов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListRecurringApplicationCharges))
	s.AddTool(mcp.NewTool("get_recurring_application_charge", mcp.WithDescription("Get a single recurring application charge by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить рекуррентный счёт"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "charge_id", svc.GetRecurringApplicationCharge))
	s.AddTool(mcp.NewTool("create_recurring_application_charge", mcp.WithDescription("Create or add a new recurring application charge in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать рекуррентный счёт"), mcp.WithString("name", mcp.Required(), mcp.Description("Назначение")), mcp.WithNumber("price", mcp.Required(), mcp.Description("Сумма")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateRecurringApplicationCharge))
	s.AddTool(mcp.NewTool("update_recurring_application_charge", mcp.WithDescription("Update the recurring application charge monthly amount. InSales uses a singleton PUT path (no charge_id)."), mcp.WithToolTitle("Обновить рекуррентный счёт"), mcp.WithNumber("monthly", mcp.Required(), mcp.Description("Monthly payment amount")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateRecurringApplicationChargeHandler(svc))
	s.AddTool(mcp.NewTool("delete_recurring_application_charge", mcp.WithDescription("Delete or remove recurring application charge on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить рекуррентный счёт"), mcp.WithString("charge_id", mcp.Required(), mcp.Description("ID счёта")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteRecurringApplicationCharge, "charge_id", "Счёт удалён"))

	// === BONUS SYSTEM TRANSACTIONS ===
	s.AddTool(mcp.NewTool("list_bonus_system_transactions", mcp.WithDescription("List bonus transaction from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список бонусных транзакций"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListBonusSystemTransactions))
	s.AddTool(mcp.NewTool("get_bonus_system_transaction", mcp.WithDescription("Get a single bonus transaction by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить бонусную транзакцию"), mcp.WithString("transaction_id", mcp.Required(), mcp.Description("ID транзакции")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "transaction_id", svc.GetBonusSystemTransaction))
	s.AddTool(mcp.NewTool("create_bonus_system_transaction", mcp.WithDescription("Create or add a new bonus transaction in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать бонусную транзакцию"), mcp.WithNumber("client_id", mcp.Required(), mcp.Description("ID клиента")), mcp.WithNumber("amount", mcp.Required(), mcp.Description("Сумма баллов")), mcp.WithString("comment", mcp.Description("Комментарий")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createBonusSystemTransactionHandler(svc))
	s.AddTool(mcp.NewTool("list_client_bonus_transactions", mcp.WithDescription("List client from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список бонусных транзакций клиента"), mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listClientBonusTransactionsHandler(svc))

	// === OAUTH APPS ===
	s.AddTool(mcp.NewTool("list_oauth_apps", mcp.WithDescription("List OAuth app from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список OAuth приложений"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListOAuthApps))
	s.AddTool(mcp.NewTool("get_oauth_app", mcp.WithDescription("Get a single OAuth app by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить OAuth приложение"), mcp.WithString("app_id", mcp.Required(), mcp.Description("ID приложения")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "app_id", svc.GetOAuthApp))
	s.AddTool(mcp.NewTool("create_oauth_app", mcp.WithDescription("Create or add a new OAuth app in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать OAuth приложение"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateOAuthApp))
	s.AddTool(mcp.NewTool("update_oauth_app", mcp.WithDescription("Update an existing OAuth app by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить OAuth приложение"), mcp.WithString("app_id", mcp.Required(), mcp.Description("ID приложения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "app_id", svc.UpdateOAuthApp))
	s.AddTool(mcp.NewTool("delete_oauth_app", mcp.WithDescription("Delete or remove OAuth app on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить OAuth приложение"), mcp.WithString("app_id", mcp.Required(), mcp.Description("ID приложения")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteOAuthApp, "app_id", "Приложение удалено"))

	// === USER TASKS ===
	s.AddTool(mcp.NewTool("list_user_tasks", mcp.WithDescription("List user task from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список задач пользователя"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListUserTasks))
	s.AddTool(mcp.NewTool("get_user_task", mcp.WithDescription("Get a single user task by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить задачу"), mcp.WithString("task_id", mcp.Required(), mcp.Description("ID задачи")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "task_id", svc.GetUserTask))
	s.AddTool(mcp.NewTool("create_user_task", mcp.WithDescription("Create or add a new user task in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать задачу"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateUserTask))
	s.AddTool(mcp.NewTool("update_user_task", mcp.WithDescription("Update an existing user task by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить задачу"), mcp.WithString("task_id", mcp.Required(), mcp.Description("ID задачи")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "task_id", svc.UpdateUserTask))
	s.AddTool(mcp.NewTool("delete_user_task", mcp.WithDescription("Delete or remove user task on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить задачу"), mcp.WithString("task_id", mcp.Required(), mcp.Description("ID задачи")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteUserTask, "task_id", "Задача удалена"))

	// === SUPPLEMENTARY ===
	s.AddTool(mcp.NewTool("list_supplementaries", mcp.WithDescription("List supplementary-product link from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список дополнительных товаров"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListSupplementaries))
	s.AddTool(mcp.NewTool("get_supplementary", mcp.WithDescription("Get a single supplementary-product link by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить дополнительный товар"), mcp.WithString("supplementary_id", mcp.Required(), mcp.Description("ID")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "supplementary_id", svc.GetSupplementary))
	s.AddTool(mcp.NewTool("create_supplementary", mcp.WithDescription("Create or add a new supplementary-product link in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать дополнительный товар"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateSupplementary))
	s.AddTool(mcp.NewTool("delete_supplementary", mcp.WithDescription("Delete or remove supplementary-product link on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить дополнительный товар"), mcp.WithString("supplementary_id", mcp.Required(), mcp.Description("ID")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteSupplementary, "supplementary_id", "Удалено"))

	// === SIMILARS ===
	s.AddTool(mcp.NewTool("list_similars", mcp.WithDescription("List similar-product link from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Похожие товары"), mcp.WithString("product_id", mcp.Required(), mcp.Description("ID товара")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), listSimilarsByProductHandler(svc))
	s.AddTool(mcp.NewTool("create_similar", mcp.WithDescription("Create or add a new similar-product link in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Добавить похожий товар"), mcp.WithNumber("product_id", mcp.Required(), mcp.Description("ID товара")), mcp.WithNumber("similar_product_id", mcp.Required(), mcp.Description("ID похожего товара")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateSimilar))
	s.AddTool(mcp.NewTool("delete_similar", mcp.WithDescription("Delete or remove similar-product link on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить похожий товар"), mcp.WithString("similar_id", mcp.Required(), mcp.Description("ID связи")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteSimilar, "similar_id", "Удалено"))

	// === UNITS ===
	s.AddTool(mcp.NewTool("list_units", mcp.WithDescription("List unit from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список единиц измерения"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListUnits))

	// === OUTLETS ===
	s.AddTool(mcp.NewTool("list_outlets", mcp.WithDescription("List outlet from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список точек выдачи"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListOutlets))
	s.AddTool(mcp.NewTool("get_outlet", mcp.WithDescription("Get a single outlet by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить точку выдачи по ID"), mcp.WithString("outlet_id", mcp.Required(), mcp.Description("ID точки выдачи")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "outlet_id", svc.GetOutlet))
	s.AddTool(mcp.NewTool("create_outlet", mcp.WithDescription("Create or add a new outlet in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать точку выдачи"), mcp.WithString("title", mcp.Required(), mcp.Description("Название")), mcp.WithString("address", mcp.Required(), mcp.Description("Адрес")), mcp.WithString("price", mcp.Required(), mcp.Description("Стоимость доставки")), mcp.WithString("latitude", mcp.Required(), mcp.Description("Широта")), mcp.WithString("longitude", mcp.Required(), mcp.Description("Долгота")), mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateOutlet))
	s.AddTool(mcp.NewTool("update_outlet", mcp.WithDescription("Update an existing outlet by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить точку выдачи"), mcp.WithString("outlet_id", mcp.Required(), mcp.Description("ID точки выдачи")), mcp.WithString("title", mcp.Description("Название")), mcp.WithString("address", mcp.Description("Адрес")), mcp.WithString("price", mcp.Description("Стоимость доставки")), mcp.WithString("latitude", mcp.Description("Широта")), mcp.WithString("longitude", mcp.Description("Долгота")), mcp.WithString("description", mcp.Description("Описание")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "outlet_id", svc.UpdateOutlet))
	s.AddTool(mcp.NewTool("delete_outlet", mcp.WithDescription("Delete or remove outlet on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить точку выдачи"), mcp.WithString("outlet_id", mcp.Required(), mcp.Description("ID точки выдачи")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteOutlet, "outlet_id", "Точка выдачи удалена"))

	// === MARKETPLACES ===
	s.AddTool(mcp.NewTool("list_marketplaces", mcp.WithDescription("List marketplace from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Список маркетплейсов"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericListHandler(svc, svc.ListMarketplaces))
	s.AddTool(mcp.NewTool("get_marketplace", mcp.WithDescription("Get a single marketplace by id or handle. Returns the full object as JSON text and structuredContent. Prefer this over list_* when the id is already known."), mcp.WithToolTitle("Получить маркетплейс по ID"), mcp.WithString("marketplace_id", mcp.Required(), mcp.Description("ID маркетплейса")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericGetHandler(svc, "marketplace_id", svc.GetMarketplace))
	s.AddTool(mcp.NewTool("create_marketplace", mcp.WithDescription("Create or add a new marketplace in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать маркетплейс. Типы: Marketplace::GoogleMerchant, Marketplace::LiquidYandexMarket, Marketplace::ModelYandexMarket"), mcp.WithString("name", mcp.Required(), mcp.Description("Название")), mcp.WithString("type", mcp.Required(), mcp.Description("Тип: Marketplace::GoogleMerchant, Marketplace::LiquidYandexMarket, Marketplace::ModelYandexMarket")), mcp.WithString("page_encoding", mcp.Description("Кодировка (для Yandex): windows-1251, utf-8")), mcp.WithString("offer_liquid", mcp.Description("Liquid шаблон (для LiquidYandexMarket)")), mcp.WithString("vendor_id", mcp.Description("ID свойства (для ModelYandexMarket)")), mcp.WithString("image_style", mcp.Description("Стиль изображения (для ModelYandexMarket): thumb, compact, medium, large, original")), mcp.WithString("model_type", mcp.Description("Тип наименования (для ModelYandexMarket): name, sku")), mcp.WithString("shop_name", mcp.Description("Название магазина")), mcp.WithString("shop_company", mcp.Description("Название компании")), mcp.WithNumber("description_type", mcp.Description("Тип описания: 1=краткое, 2=полное")), mcp.WithString("local_delivery_cost", mcp.Description("Стоимость локальной доставки")), mcp.WithBoolean("adult", mcp.Description("Только для взрослых")), mcp.WithString("utm_label", mcp.Description("UTM метка")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericCreateHandler(svc, svc.CreateMarketplace))
	s.AddTool(mcp.NewTool("update_marketplace", mcp.WithDescription("Update an existing marketplace by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить маркетплейс"), mcp.WithString("marketplace_id", mcp.Required(), mcp.Description("ID маркетплейса")), mcp.WithString("name", mcp.Description("Название")), mcp.WithString("page_encoding", mcp.Description("Кодировка")), mcp.WithString("offer_liquid", mcp.Description("Liquid шаблон")), mcp.WithString("vendor_id", mcp.Description("ID свойства")), mcp.WithString("image_style", mcp.Description("Стиль изображения")), mcp.WithString("model_type", mcp.Description("Тип наименования")), mcp.WithString("shop_name", mcp.Description("Название магазина")), mcp.WithString("shop_company", mcp.Description("Название компании")), mcp.WithNumber("description_type", mcp.Description("Тип описания")), mcp.WithString("local_delivery_cost", mcp.Description("Стоимость локальной доставки")), mcp.WithBoolean("adult", mcp.Description("Только для взрослых")), mcp.WithString("utm_label", mcp.Description("UTM метка")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericUpdateHandler(svc, "marketplace_id", svc.UpdateMarketplace))
	s.AddTool(mcp.NewTool("delete_marketplace", mcp.WithDescription("Delete or remove marketplace on the InSales side. This is destructive and cannot be undone from this tool. Returns a short confirmation or the updated parent object."), mcp.WithToolTitle("Удалить маркетплейс"), mcp.WithString("marketplace_id", mcp.Required(), mcp.Description("ID маркетплейса")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), genericDeleteHandler(svc.DeleteMarketplace, "marketplace_id", "Маркетплейс удалён"))

	// === PROVIDER TOKEN ===
	s.AddTool(mcp.NewTool("create_provider_token", mcp.WithDescription("Create or add a new messenger provider token in InSales. Pass required fields using API names (language-invariant ids/handles). Returns the created object including its id."), mcp.WithToolTitle("Создать токен провайдера мессенджера"), mcp.WithString("messenger_type", mcp.Required(), mcp.Description("Тип мессенджера: whats_app, telegram, viber, sms, email, vk, instagram")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), createProviderTokenHandler(svc))

	// === ACCOUNT ===
	s.AddTool(mcp.NewTool("update_account", mcp.WithDescription("Update an existing account settings by id. Only send fields that should change. Returns the updated object as JSON."), mcp.WithToolTitle("Обновить настройки аккаунта"), mcp.WithString("client_cookies_whitelist", mcp.Description("Белый список cookies")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), updateAccountHandler(svc))

	// === COUNTS ===
	s.AddTool(mcp.NewTool("get_orders_count", mcp.WithDescription("List order from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Количество заказов"), mcp.WithString("fulfillment_status", mcp.Description("Статус выполнения")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getOrdersCountHandler(svc))
	s.AddTool(mcp.NewTool("get_products_count", mcp.WithDescription("List product from the InSales shop. Returns JSON items (envelope with total/page/per_page when paginated). Use page/per_page and filters to narrow results. IDs and handles are language-invariant."), mcp.WithToolTitle("Количество товаров"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), getProductsCountHandler(svc))

	// === PIPELINE TOOLS ===
	s.AddTool(mcp.NewTool("search_product_details",
		mcp.WithDescription("Search product and return matching data in one call. Pipeline tool: combines lookup steps for common workflows. Returns JSON items relevant to the query."), mcp.WithToolTitle("Поиск товара с полной информацией: кастомные поля, метаполя, изображения. Возвращает полную карточку товара одним вызовом."),
		mcp.WithString("query", mcp.Description("Поисковый запрос (название товара)")),
		mcp.WithString("product_id", mcp.Description("ID товара (альтернатива query)")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), searchProductDetailsHandler(svc))

	s.AddTool(mcp.NewTool("order_full_details",
		mcp.WithDescription("Manage order via InSales API. Полная информация о заказе: позиции, адрес доставки, клиент, кастомные поля. Всё одним вызовом."), mcp.WithToolTitle("Полная информация о заказе: позиции, адрес доставки, клиент, кастомные поля. Всё одним вызовом."),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("ID заказа")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), orderFullDetailsHandler(svc))

	s.AddTool(mcp.NewTool("client_full_profile",
		mcp.WithDescription("Manage client via InSales API. Полный профиль клиента: данные, последние заказы, бонусы, метаполя. Всё одним вызовом."), mcp.WithToolTitle("Полный профиль клиента: данные, последние заказы, бонусы, метаполя. Всё одним вызовом."),
		mcp.WithString("client_id", mcp.Required(), mcp.Description("ID клиента")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), clientFullProfileHandler(svc))

	s.AddTool(mcp.NewTool("product_catalog_search",
		mcp.WithDescription("Manage product via InSales API. Поиск по каталогу товаров с фильтрами, пагинацией и подсчётом."), mcp.WithToolTitle("Поиск по каталогу товаров с фильтрами, пагинацией и подсчётом."),
		mcp.WithString("query", mcp.Description("Поисковый запрос")),
		mcp.WithString("category_id", mcp.Description("Фильтр по категории")),
		mcp.WithString("collection_id", mcp.Description("Фильтр по коллекции")),
		mcp.WithNumber("page", mcp.DefaultNumber(1)),
		mcp.WithNumber("per_page", mcp.DefaultNumber(25)),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), productCatalogSearchHandler(svc))

	s.AddTool(mcp.NewTool("collection_products",
		mcp.WithDescription("Manage collection via InSales API. Товары коллекции с полной информацией: данные коллекции + все её товары."), mcp.WithToolTitle("Товары коллекции с полной информацией: данные коллекции + все её товары."),
		mcp.WithString("collection_id", mcp.Required(), mcp.Description("ID коллекции")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	), collectionProductsHandler(svc))

	// === DATA EXPORT (CSV) ===
	s.AddTool(mcp.NewTool("export_products",
		mcp.WithDescription("Export product as CSV text. Returns a CSV string with stable column headers. Use for downloads or bulk processing."), mcp.WithToolTitle("Экспорт товаров в CSV. Возвращает CSV-файл с колонками: id, title, price, sku, vendor, quantity, created_at"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), exportProductsHandler(svc))
	s.AddTool(mcp.NewTool("export_orders",
		mcp.WithDescription("Export order as CSV text. Returns a CSV string with stable column headers. Use for downloads or bulk processing."), mcp.WithToolTitle("Экспорт заказов в CSV. Возвращает CSV-файл с колонками: id, number, status, fulfillment_status, total_price, client_name, created_at"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), exportOrdersHandler(svc))
	s.AddTool(mcp.NewTool("export_clients",
		mcp.WithDescription("Export client as CSV text. Returns a CSV string with stable column headers. Use for downloads or bulk processing."), mcp.WithToolTitle("Экспорт клиентов в CSV. Возвращает CSV-файл с колонками: id, name, surname, email, phone, created_at"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	), exportClientsHandler(svc))
}

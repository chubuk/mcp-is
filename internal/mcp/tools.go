package mcp

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

func GetAllTools() []Tool {
	return []Tool{
		// === ТОВАРЫ ===
		{
			Name:        "list_products",
			Description: "Получить список товаров магазина с пагинацией и фильтрацией",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page":         map[string]interface{}{"type": "integer", "description": "Номер страницы", "default": 1},
					"per_page":     map[string]interface{}{"type": "integer", "description": "Товаров на странице (макс 250)", "default": 25},
					"category_id":  map[string]interface{}{"type": "integer", "description": "Фильтр по категории"},
					"collection_id": map[string]interface{}{"type": "integer", "description": "Фильтр по коллекции"},
					"q":            map[string]interface{}{"type": "string", "description": "Поиск по названию/артикулу"},
				},
			},
		},
		{
			Name:        "get_product",
			Description: "Получить детальную информацию о товаре по ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id": map[string]interface{}{"type": "string", "description": "ID товара"},
				},
				"required": []string{"product_id"},
			},
		},
		{
			Name:        "create_product",
			Description: "Создать новый товар в магазине",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":       map[string]interface{}{"type": "string", "description": "Название товара"},
					"description": map[string]interface{}{"type": "string", "description": "Описание товара"},
					"price":       map[string]interface{}{"type": "string", "description": "Цена (например 1000.00)"},
					"vendor":      map[string]interface{}{"type": "string", "description": "Производитель"},
					"sku":         map[string]interface{}{"type": "string", "description": "Артикул"},
					"quantity":    map[string]interface{}{"type": "integer", "description": "Количество на складе"},
				},
				"required": []string{"title", "price"},
			},
		},
		{
			Name:        "update_product",
			Description: "Обновить данные товара",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id":  map[string]interface{}{"type": "string", "description": "ID товара"},
					"title":       map[string]interface{}{"type": "string"},
					"description": map[string]interface{}{"type": "string"},
					"price":       map[string]interface{}{"type": "string"},
					"quantity":    map[string]interface{}{"type": "integer"},
				},
				"required": []string{"product_id"},
			},
		},
		{
			Name:        "delete_product",
			Description: "Удалить товар из магазина",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id": map[string]interface{}{"type": "string", "description": "ID товара"},
				},
				"required": []string{"product_id"},
			},
		},
		// === ЗАКАЗЫ ===
		{
			Name:        "list_orders",
			Description: "Получить список заказов с фильтрацией по статусу, дате, клиенту",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page":                map[string]interface{}{"type": "integer", "default": 1},
					"per_page":            map[string]interface{}{"type": "integer", "default": 25},
					"status":              map[string]interface{}{"type": "string", "enum": []string{"open", "closed", "canceled"}},
					"fulfillment_status":  map[string]interface{}{"type": "string", "enum": []string{"pending", "shipped", "delivered", "canceled"}},
					"paid":                map[string]interface{}{"type": "boolean"},
					"client_id":           map[string]interface{}{"type": "integer"},
				},
			},
		},
		{
			Name:        "get_order",
			Description: "Получить детали заказа по ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"order_id": map[string]interface{}{"type": "string", "description": "ID заказа"},
				},
				"required": []string{"order_id"},
			},
		},
		{
			Name:        "update_order",
			Description: "Обновить статус или данные заказа",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"order_id":            map[string]interface{}{"type": "string", "description": "ID заказа"},
					"status":              map[string]interface{}{"type": "string", "enum": []string{"open", "closed", "canceled"}},
					"fulfillment_status":  map[string]interface{}{"type": "string", "enum": []string{"pending", "shipped", "delivered", "canceled"}},
					"comment":             map[string]interface{}{"type": "string"},
				},
				"required": []string{"order_id"},
			},
		},
		{
			Name:        "create_order",
			Description: "Создать новый заказ вручную",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"client_id": map[string]interface{}{"type": "integer", "description": "ID клиента"},
					"comment":   map[string]interface{}{"type": "string"},
				},
				"required": []string{"client_id"},
			},
		},
		// === КЛИЕНТЫ ===
		{
			Name:        "list_clients",
			Description: "Получить список клиентов с поиском",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page":     map[string]interface{}{"type": "integer", "default": 1},
					"per_page": map[string]interface{}{"type": "integer", "default": 25},
					"q":        map[string]interface{}{"type": "string", "description": "Поиск по имени/email/телефону"},
				},
			},
		},
		{
			Name:        "get_client",
			Description: "Получить информацию о клиенте по ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"client_id": map[string]interface{}{"type": "string", "description": "ID клиента"},
				},
				"required": []string{"client_id"},
			},
		},
		{
			Name:        "create_client",
			Description: "Создать нового клиента",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"email":      map[string]interface{}{"type": "string"},
					"first_name": map[string]interface{}{"type": "string"},
					"last_name":  map[string]interface{}{"type": "string"},
					"phone":      map[string]interface{}{"type": "string"},
				},
				"required": []string{"email"},
			},
		},
		{
			Name:        "update_client",
			Description: "Обновить данные клиента",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"client_id":  map[string]interface{}{"type": "string", "description": "ID клиента"},
					"email":      map[string]interface{}{"type": "string"},
					"first_name": map[string]interface{}{"type": "string"},
					"last_name":  map[string]interface{}{"type": "string"},
					"phone":      map[string]interface{}{"type": "string"},
				},
				"required": []string{"client_id"},
			},
		},
		{
			Name:        "delete_client",
			Description: "Удалить клиента",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"client_id": map[string]interface{}{"type": "string", "description": "ID клиента"},
				},
				"required": []string{"client_id"},
			},
		},
		// === КАТЕГОРИИ ===
		{
			Name:        "list_categories",
			Description: "Получить список категорий товаров",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page":      map[string]interface{}{"type": "integer", "default": 1},
					"per_page":  map[string]interface{}{"type": "integer", "default": 25},
					"parent_id": map[string]interface{}{"type": "integer", "description": "Фильтр по родительской категории"},
				},
			},
		},
		{
			Name:        "get_category",
			Description: "Получить информацию о категории по ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"category_id": map[string]interface{}{"type": "string", "description": "ID категории"},
				},
				"required": []string{"category_id"},
			},
		},
		{
			Name:        "create_category",
			Description: "Создать новую категорию",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":     map[string]interface{}{"type": "string", "description": "Название категории"},
					"parent_id": map[string]interface{}{"type": "integer", "description": "ID родительской категории"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "update_category",
			Description: "Обновить категорию",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"category_id": map[string]interface{}{"type": "string", "description": "ID категории"},
					"title":       map[string]interface{}{"type": "string"},
				},
				"required": []string{"category_id"},
			},
		},
		{
			Name:        "delete_category",
			Description: "Удалить категорию",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"category_id": map[string]interface{}{"type": "string", "description": "ID категории"},
				},
				"required": []string{"category_id"},
			},
		},
		// === КОЛЛЕКЦИИ ===
		{
			Name:        "list_collections",
			Description: "Получить список коллекций (умных подборок)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page":     map[string]interface{}{"type": "integer", "default": 1},
					"per_page": map[string]interface{}{"type": "integer", "default": 25},
				},
			},
		},
		{
			Name:        "get_collection",
			Description: "Получить информацию о коллекции по ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"collection_id": map[string]interface{}{"type": "string", "description": "ID коллекции"},
				},
				"required": []string{"collection_id"},
			},
		},
		{
			Name:        "create_collection",
			Description: "Создать новую коллекцию",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Название коллекции"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "update_collection",
			Description: "Обновить коллекцию",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"collection_id": map[string]interface{}{"type": "string", "description": "ID коллекции"},
					"title":         map[string]interface{}{"type": "string"},
				},
				"required": []string{"collection_id"},
			},
		},
		{
			Name:        "delete_collection",
			Description: "Удалить коллекцию",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"collection_id": map[string]interface{}{"type": "string", "description": "ID коллекции"},
				},
				"required": []string{"collection_id"},
			},
		},
		// === ИЗОБРАЖЕНИЯ ===
		{
			Name:        "list_images",
			Description: "Получить список изображений товара",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id": map[string]interface{}{"type": "string", "description": "ID товара"},
				},
				"required": []string{"product_id"},
			},
		},
		{
			Name:        "add_image",
			Description: "Добавить изображение к товару",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id": map[string]interface{}{"type": "string", "description": "ID товара"},
					"url":        map[string]interface{}{"type": "string", "description": "URL изображения"},
					"filename":   map[string]interface{}{"type": "string", "description": "Имя файла"},
					"position":   map[string]interface{}{"type": "integer", "description": "Позиция"},
				},
				"required": []string{"product_id", "url"},
			},
		},
		{
			Name:        "delete_image",
			Description: "Удалить изображение товара",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"product_id": map[string]interface{}{"type": "string", "description": "ID товара"},
					"image_id":   map[string]interface{}{"type": "string", "description": "ID изображения"},
				},
				"required": []string{"product_id", "image_id"},
			},
		},
		// === МЕТАПОЛЯ ===
		{
			Name:        "list_meta_fields",
			Description: "Получить список метаполей",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"resource_type": map[string]interface{}{"type": "string", "enum": []string{"Product", "Variant", "Order", "Client", "Category", "Page"}},
					"resource_id":   map[string]interface{}{"type": "integer"},
				},
			},
		},
		{
			Name:        "create_meta_field",
			Description: "Создать метаполе для ресурса",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"resource_type": map[string]interface{}{"type": "string", "enum": []string{"Product", "Variant", "Order", "Client", "Category", "Page"}},
					"resource_id":   map[string]interface{}{"type": "integer"},
					"key":           map[string]interface{}{"type": "string"},
					"value":         map[string]interface{}{"type": "string"},
				},
				"required": []string{"resource_type", "resource_id", "key", "value"},
			},
		},
		{
			Name:        "update_meta_field",
			Description: "Обновить метаполе",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"meta_field_id": map[string]interface{}{"type": "string", "description": "ID метаполя"},
					"value":         map[string]interface{}{"type": "string"},
				},
				"required": []string{"meta_field_id"},
			},
		},
		{
			Name:        "delete_meta_field",
			Description: "Удалить метаполе",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"meta_field_id": map[string]interface{}{"type": "string", "description": "ID метаполя"},
				},
				"required": []string{"meta_field_id"},
			},
		},
		// === ВЕБХУКИ ===
		{
			Name:        "list_webhooks",
			Description: "Получить список подписанных вебхуков",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "create_webhook",
			Description: "Подписаться на событие InSales (order_create, product_update и т.д.)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"event":  map[string]interface{}{"type": "string", "enum": []string{"order_create", "order_update", "order_delete", "product_create", "product_update", "product_delete", "client_create", "client_update", "client_delete"}},
					"url":    map[string]interface{}{"type": "string", "description": "URL для получения событий"},
					"format": map[string]interface{}{"type": "string", "enum": []string{"json", "xml"}, "default": "json"},
				},
				"required": []string{"event", "url"},
			},
		},
		{
			Name:        "update_webhook",
			Description: "Обновить URL вебхука",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"webhook_id": map[string]interface{}{"type": "string", "description": "ID вебхука"},
					"url":        map[string]interface{}{"type": "string"},
				},
				"required": []string{"webhook_id"},
			},
		},
		{
			Name:        "delete_webhook",
			Description: "Отписаться от вебхука",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"webhook_id": map[string]interface{}{"type": "string", "description": "ID вебхука"},
				},
				"required": []string{"webhook_id"},
			},
		},
		// === СИСТЕМНЫЕ ===
		{
			Name:        "get_shop_info",
			Description: "Получить информацию о магазине (название, домен, валюта)",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "list_currencies",
			Description: "Получить список валют магазина",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "list_delivery_variants",
			Description: "Получить доступные способы доставки",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "list_payment_gateways",
			Description: "Получить доступные способы оплаты",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

package mcp

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func GetAllTools() []Tool {
	return []Tool{
		// === ТОВАРЫ ===
		{
			Name:        "list_products",
			Description: "Получить список товаров магазина с пагинацией и фильтрацией",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page":         map[string]any{"type": "integer", "description": "Номер страницы", "default": 1},
					"per_page":     map[string]any{"type": "integer", "description": "Товаров на странице (макс 250)", "default": 25},
					"category_id":  map[string]any{"type": "integer", "description": "Фильтр по категории"},
					"collection_id": map[string]any{"type": "integer", "description": "Фильтр по коллекции"},
					"q":            map[string]any{"type": "string", "description": "Поиск по названию/артикулу"},
				},
			},
		},
		{
			Name:        "get_product",
			Description: "Получить детальную информацию о товаре по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
				},
				"required": []string{"product_id"},
			},
		},
		{
			Name:        "create_product",
			Description: "Создать новый товар в магазине",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":       map[string]any{"type": "string", "description": "Название товара"},
					"description": map[string]any{"type": "string", "description": "Описание товара"},
					"price":       map[string]any{"type": "string", "description": "Цена (например 1000.00)"},
					"vendor":      map[string]any{"type": "string", "description": "Производитель"},
					"sku":         map[string]any{"type": "string", "description": "Артикул"},
					"quantity":    map[string]any{"type": "integer", "description": "Количество на складе"},
				},
				"required": []string{"title", "price"},
			},
		},
		{
			Name:        "update_product",
			Description: "Обновить данные товара",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id":  map[string]any{"type": "string", "description": "ID товара"},
					"title":       map[string]any{"type": "string"},
					"description": map[string]any{"type": "string"},
					"price":       map[string]any{"type": "string"},
					"quantity":    map[string]any{"type": "integer"},
				},
				"required": []string{"product_id"},
			},
		},
		{
			Name:        "delete_product",
			Description: "Удалить товар из магазина",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
				},
				"required": []string{"product_id"},
			},
		},
		// === ЗАКАЗЫ ===
		{
			Name:        "list_orders",
			Description: "Получить список заказов с фильтрацией по статусу, дате, клиенту",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page":                map[string]any{"type": "integer", "default": 1},
					"per_page":            map[string]any{"type": "integer", "default": 25},
					"status":              map[string]any{"type": "string", "enum": []string{"open", "closed", "canceled"}},
					"fulfillment_status":  map[string]any{"type": "string", "enum": []string{"pending", "shipped", "delivered", "canceled"}},
					"paid":                map[string]any{"type": "boolean"},
					"client_id":           map[string]any{"type": "integer"},
				},
			},
		},
		{
			Name:        "get_order",
			Description: "Получить детали заказа по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"order_id": map[string]any{"type": "string", "description": "ID заказа"},
				},
				"required": []string{"order_id"},
			},
		},
		{
			Name:        "update_order",
			Description: "Обновить статус или данные заказа",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"order_id":            map[string]any{"type": "string", "description": "ID заказа"},
					"status":              map[string]any{"type": "string", "enum": []string{"open", "closed", "canceled"}},
					"fulfillment_status":  map[string]any{"type": "string", "enum": []string{"pending", "shipped", "delivered", "canceled"}},
					"comment":             map[string]any{"type": "string"},
				},
				"required": []string{"order_id"},
			},
		},
		{
			Name:        "create_order",
			Description: "Создать новый заказ вручную",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"client_id": map[string]any{"type": "integer", "description": "ID клиента"},
					"comment":   map[string]any{"type": "string"},
				},
				"required": []string{"client_id"},
			},
		},
		// === КЛИЕНТЫ ===
		{
			Name:        "list_clients",
			Description: "Получить список клиентов с поиском",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page":     map[string]any{"type": "integer", "default": 1},
					"per_page": map[string]any{"type": "integer", "default": 25},
					"q":        map[string]any{"type": "string", "description": "Поиск по имени/email/телефону"},
				},
			},
		},
		{
			Name:        "get_client",
			Description: "Получить информацию о клиенте по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"client_id": map[string]any{"type": "string", "description": "ID клиента"},
				},
				"required": []string{"client_id"},
			},
		},
		{
			Name:        "create_client",
			Description: "Создать нового клиента",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"email":      map[string]any{"type": "string"},
					"first_name": map[string]any{"type": "string"},
					"last_name":  map[string]any{"type": "string"},
					"phone":      map[string]any{"type": "string"},
				},
				"required": []string{"email"},
			},
		},
		{
			Name:        "update_client",
			Description: "Обновить данные клиента",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"client_id":  map[string]any{"type": "string", "description": "ID клиента"},
					"email":      map[string]any{"type": "string"},
					"first_name": map[string]any{"type": "string"},
					"last_name":  map[string]any{"type": "string"},
					"phone":      map[string]any{"type": "string"},
				},
				"required": []string{"client_id"},
			},
		},
		{
			Name:        "delete_client",
			Description: "Удалить клиента",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"client_id": map[string]any{"type": "string", "description": "ID клиента"},
				},
				"required": []string{"client_id"},
			},
		},
		// === КАТЕГОРИИ ===
		{
			Name:        "list_categories",
			Description: "Получить список категорий товаров",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page":      map[string]any{"type": "integer", "default": 1},
					"per_page":  map[string]any{"type": "integer", "default": 25},
					"parent_id": map[string]any{"type": "integer", "description": "Фильтр по родительской категории"},
				},
			},
		},
		{
			Name:        "get_category",
			Description: "Получить информацию о категории по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"category_id": map[string]any{"type": "string", "description": "ID категории"},
				},
				"required": []string{"category_id"},
			},
		},
		{
			Name:        "create_category",
			Description: "Создать новую категорию",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":     map[string]any{"type": "string", "description": "Название категории"},
					"parent_id": map[string]any{"type": "integer", "description": "ID родительской категории"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "update_category",
			Description: "Обновить категорию",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"category_id": map[string]any{"type": "string", "description": "ID категории"},
					"title":       map[string]any{"type": "string"},
				},
				"required": []string{"category_id"},
			},
		},
		{
			Name:        "delete_category",
			Description: "Удалить категорию",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"category_id": map[string]any{"type": "string", "description": "ID категории"},
				},
				"required": []string{"category_id"},
			},
		},
		// === КОЛЛЕКЦИИ ===
		{
			Name:        "list_collections",
			Description: "Получить список коллекций (умных подборок)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page":     map[string]any{"type": "integer", "default": 1},
					"per_page": map[string]any{"type": "integer", "default": 25},
				},
			},
		},
		{
			Name:        "get_collection",
			Description: "Получить информацию о коллекции по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"collection_id": map[string]any{"type": "string", "description": "ID коллекции"},
				},
				"required": []string{"collection_id"},
			},
		},
		{
			Name:        "create_collection",
			Description: "Создать новую коллекцию",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{"type": "string", "description": "Название коллекции"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "update_collection",
			Description: "Обновить коллекцию",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"collection_id": map[string]any{"type": "string", "description": "ID коллекции"},
					"title":         map[string]any{"type": "string"},
				},
				"required": []string{"collection_id"},
			},
		},
		{
			Name:        "delete_collection",
			Description: "Удалить коллекцию",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"collection_id": map[string]any{"type": "string", "description": "ID коллекции"},
				},
				"required": []string{"collection_id"},
			},
		},
		// === ИЗОБРАЖЕНИЯ ===
		{
			Name:        "list_images",
			Description: "Получить список изображений товара",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
				},
				"required": []string{"product_id"},
			},
		},
		{
			Name:        "add_image",
			Description: "Добавить изображение к товару",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
					"url":        map[string]any{"type": "string", "description": "URL изображения"},
					"filename":   map[string]any{"type": "string", "description": "Имя файла"},
					"position":   map[string]any{"type": "integer", "description": "Позиция"},
				},
				"required": []string{"product_id", "url"},
			},
		},
		{
			Name:        "delete_image",
			Description: "Удалить изображение товара",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
					"image_id":   map[string]any{"type": "string", "description": "ID изображения"},
				},
				"required": []string{"product_id", "image_id"},
			},
		},
		// === МЕТАПОЛЯ ===
		{
			Name:        "list_meta_fields",
			Description: "Получить список метаполей",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"resource_type": map[string]any{"type": "string", "enum": []string{"Product", "Variant", "Order", "Client", "Category", "Page"}},
					"resource_id":   map[string]any{"type": "integer"},
				},
			},
		},
		{
			Name:        "create_meta_field",
			Description: "Создать метаполе для ресурса",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"resource_type": map[string]any{"type": "string", "enum": []string{"Product", "Variant", "Order", "Client", "Category", "Page"}},
					"resource_id":   map[string]any{"type": "integer"},
					"key":           map[string]any{"type": "string"},
					"value":         map[string]any{"type": "string"},
				},
				"required": []string{"resource_type", "resource_id", "key", "value"},
			},
		},
		{
			Name:        "update_meta_field",
			Description: "Обновить метаполе",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"meta_field_id": map[string]any{"type": "string", "description": "ID метаполя"},
					"value":         map[string]any{"type": "string"},
				},
				"required": []string{"meta_field_id"},
			},
		},
		{
			Name:        "delete_meta_field",
			Description: "Удалить метаполе",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"meta_field_id": map[string]any{"type": "string", "description": "ID метаполя"},
				},
				"required": []string{"meta_field_id"},
			},
		},
		// === ВЕБХУКИ ===
		{
			Name:        "list_webhooks",
			Description: "Получить список подписанных вебхуков",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "create_webhook",
			Description: "Подписаться на событие InSales (order_create, product_update и т.д.)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"event":  map[string]any{"type": "string", "enum": []string{"order_create", "order_update", "order_delete", "product_create", "product_update", "product_delete", "client_create", "client_update", "client_delete"}},
					"url":    map[string]any{"type": "string", "description": "URL для получения событий"},
					"format": map[string]any{"type": "string", "enum": []string{"json", "xml"}, "default": "json"},
				},
				"required": []string{"event", "url"},
			},
		},
		{
			Name:        "update_webhook",
			Description: "Обновить URL вебхука",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"webhook_id": map[string]any{"type": "string", "description": "ID вебхука"},
					"url":        map[string]any{"type": "string"},
				},
				"required": []string{"webhook_id"},
			},
		},
		{
			Name:        "delete_webhook",
			Description: "Отписаться от вебхука",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"webhook_id": map[string]any{"type": "string", "description": "ID вебхука"},
				},
				"required": []string{"webhook_id"},
			},
		},
		// === СИСТЕМНЫЕ ===
		{
			Name:        "get_shop_info",
			Description: "Получить информацию о магазине (название, домен, валюта)",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "list_currencies",
			Description: "Получить список валют магазина",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "list_delivery_variants",
			Description: "Получить доступные способы доставки",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "list_payment_gateways",
			Description: "Получить доступные способы оплаты",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		// === СТРАНИЦЫ ===
		{
			Name:        "list_pages",
			Description: "Получить список страниц магазина (статьи, публикации, информационные страницы)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page":     map[string]any{"type": "integer", "description": "Номер страницы", "default": 1},
					"per_page": map[string]any{"type": "integer", "description": "Страниц на экран (макс 250)", "default": 25},
				},
			},
		},
		{
			Name:        "get_page",
			Description: "Получить страницу по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page_id": map[string]any{"type": "string", "description": "ID страницы"},
				},
				"required": []string{"page_id"},
			},
		},
		{
			Name:        "create_page",
			Description: "Создать новую страницу в магазине",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":   map[string]any{"type": "string", "description": "Название страницы"},
					"content": map[string]any{"type": "string", "description": "HTML-содержимое страницы"},
					"handle":  map[string]any{"type": "string", "description": "URL-идентификатор (например 'about')"},
				},
				"required": []string{"title", "content"},
			},
		},
		{
			Name:        "update_page",
			Description: "Обновить данные страницы",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page_id": map[string]any{"type": "string", "description": "ID страницы"},
					"title":   map[string]any{"type": "string"},
					"content": map[string]any{"type": "string"},
					"handle":  map[string]any{"type": "string"},
				},
				"required": []string{"page_id"},
			},
		},
		{
			Name:        "delete_page",
			Description: "Удалить страницу",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page_id": map[string]any{"type": "string", "description": "ID страницы"},
				},
				"required": []string{"page_id"},
			},
		},
		// === АКСЕССУАРЫ ===
		{
			Name:        "list_accessories",
			Description: "Получить список аксессуаров (характеристик товаров)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page":     map[string]any{"type": "integer", "description": "Номер страницы", "default": 1},
					"per_page": map[string]any{"type": "integer", "description": "Записей на странице (макс 250)", "default": 25},
				},
			},
		},
		{
			Name:        "get_accessory",
			Description: "Получить аксессуар по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"accessory_id": map[string]any{"type": "string", "description": "ID аксессуара"},
				},
				"required": []string{"accessory_id"},
			},
		},
		{
			Name:        "create_accessory",
			Description: "Создать новый аксессуар (характеристику)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{"type": "string", "description": "Название аксессуара (например 'Цвет', 'Размер')"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "update_accessory",
			Description: "Обновить аксессуар",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"accessory_id": map[string]any{"type": "string", "description": "ID аксессуара"},
					"title":        map[string]any{"type": "string"},
					"position":     map[string]any{"type": "integer"},
				},
				"required": []string{"accessory_id"},
			},
		},
		{
			Name:        "delete_accessory",
			Description: "Удалить аксессуар",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"accessory_id": map[string]any{"type": "string", "description": "ID аксессуара"},
				},
				"required": []string{"accessory_id"},
			},
		},
		// === ВАРИАНТЫ ТОВАРОВ ===
		{
			Name:        "list_variants",
			Description: "Получить список вариантов товара (размеры, цвета и т.д.)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
				},
				"required": []string{"product_id"},
			},
		},
		{
			Name:        "get_variant",
			Description: "Получить информацию о варианте товара по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
					"variant_id": map[string]any{"type": "string", "description": "ID варианта"},
				},
				"required": []string{"product_id", "variant_id"},
			},
		},
		{
			Name:        "update_variant",
			Description: "Обновить данные варианта товара (цену, остаток, артикул)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
					"variant_id": map[string]any{"type": "string", "description": "ID варианта"},
					"price":      map[string]any{"type": "string", "description": "Цена"},
					"sku":        map[string]any{"type": "string", "description": "Артикул"},
					"quantity":   map[string]any{"type": "integer", "description": "Количество на складе"},
					"weight":     map[string]any{"type": "string", "description": "Вес"},
				},
				"required": []string{"product_id", "variant_id"},
			},
		},
		{
			Name:        "delete_variant",
			Description: "Удалить вариант товара",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"product_id": map[string]any{"type": "string", "description": "ID товара"},
					"variant_id": map[string]any{"type": "string", "description": "ID варианта"},
				},
				"required": []string{"product_id", "variant_id"},
			},
		},
		// === ЗНАЧЕНИЯ СВОЙСТВ (PROPERTY VALUES) ===
		{
			Name:        "list_property_values",
			Description: "Получить список значений аксессуара (например 'Красный', 'Синий' для аксессуара 'Цвет')",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"accessory_id": map[string]any{"type": "string", "description": "ID аксессуара"},
				},
				"required": []string{"accessory_id"},
			},
		},
		{
			Name:        "create_property_value",
			Description: "Создать новое значение для аксессуара",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"accessory_id": map[string]any{"type": "string", "description": "ID аксессуара"},
					"title":        map[string]any{"type": "string", "description": "Название значения (например 'Красный')"},
				},
				"required": []string{"accessory_id", "title"},
			},
		},
		{
			Name:        "update_property_value",
			Description: "Обновить значение аксессуара",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"property_value_id": map[string]any{"type": "string", "description": "ID значения"},
					"title":             map[string]any{"type": "string"},
				},
				"required": []string{"property_value_id"},
			},
		},
		{
			Name:        "delete_property_value",
			Description: "Удалить значение аксессуара",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"property_value_id": map[string]any{"type": "string", "description": "ID значения"},
				},
				"required": []string{"property_value_id"},
			},
		},
		// === МЕНЮ ===
		{
			Name:        "list_menus",
			Description: "Получить список навигационных меню магазина",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "get_menu",
			Description: "Получить меню по ID со списком пунктов",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"menu_id": map[string]any{"type": "string", "description": "ID меню"},
				},
				"required": []string{"menu_id"},
			},
		},
		{
			Name:        "create_menu",
			Description: "Создать новое навигационное меню",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":  map[string]any{"type": "string", "description": "Название меню"},
					"handle": map[string]any{"type": "string", "description": "URL-идентификатор"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "update_menu",
			Description: "Обновить данные меню",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"menu_id": map[string]any{"type": "string", "description": "ID меню"},
					"title":   map[string]any{"type": "string"},
					"handle":  map[string]any{"type": "string"},
				},
				"required": []string{"menu_id"},
			},
		},
		{
			Name:        "delete_menu",
			Description: "Удалить меню",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"menu_id": map[string]any{"type": "string", "description": "ID меню"},
				},
				"required": []string{"menu_id"},
			},
		},
		// === РЕДИРЕКТЫ ===
		{
			Name:        "list_redirects",
			Description: "Получить список URL-редиректов магазина",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "create_redirect",
			Description: "Создать URL-редирект (старый URL → новый URL)",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"url":    map[string]any{"type": "string", "description": "Старый URL"},
					"target": map[string]any{"type": "string", "description": "Новый URL (куда перенаправлять)"},
				},
				"required": []string{"url", "target"},
			},
		},
		{
			Name:        "delete_redirect",
			Description: "Удалить URL-редирект",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"redirect_id": map[string]any{"type": "string", "description": "ID редиректа"},
				},
				"required": []string{"redirect_id"},
			},
		},
		// === БЛОГИ ===
		{
			Name:        "list_blogs",
			Description: "Получить список блогов магазина",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "get_blog",
			Description: "Получить информацию о блоге по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"blog_id": map[string]any{"type": "string", "description": "ID блога"},
				},
				"required": []string{"blog_id"},
			},
		},
		{
			Name:        "create_blog",
			Description: "Создать новый блог в магазине",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":  map[string]any{"type": "string", "description": "Название блога"},
					"handle": map[string]any{"type": "string", "description": "URL-идентификатор"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "delete_blog",
			Description: "Удалить блог",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"blog_id": map[string]any{"type": "string", "description": "ID блога"},
				},
				"required": []string{"blog_id"},
			},
		},
		// === ЗАПИСИ БЛОГА ===
		{
			Name:        "list_blog_posts",
			Description: "Получить список записей в блоге",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"blog_id": map[string]any{"type": "string", "description": "ID блога"},
				},
				"required": []string{"blog_id"},
			},
		},
		{
			Name:        "get_blog_post",
			Description: "Получить запись блога по ID",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"blog_id": map[string]any{"type": "string", "description": "ID блога"},
					"post_id": map[string]any{"type": "string", "description": "ID записи"},
				},
				"required": []string{"blog_id", "post_id"},
			},
		},
		{
			Name:        "create_blog_post",
			Description: "Создать новую запись в блоге",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"blog_id":   map[string]any{"type": "string", "description": "ID блога"},
					"title":     map[string]any{"type": "string", "description": "Заголовок записи"},
					"body_html": map[string]any{"type": "string", "description": "HTML-содержимое записи"},
				},
				"required": []string{"blog_id", "title", "body_html"},
			},
		},
		{
			Name:        "update_blog_post",
			Description: "Обновить запись блога",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"blog_id":   map[string]any{"type": "string", "description": "ID блога"},
					"post_id":   map[string]any{"type": "string", "description": "ID записи"},
					"title":     map[string]any{"type": "string"},
					"body_html": map[string]any{"type": "string"},
				},
				"required": []string{"blog_id", "post_id"},
			},
		},
		{
			Name:        "delete_blog_post",
			Description: "Удалить запись блога",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"blog_id": map[string]any{"type": "string", "description": "ID блога"},
					"post_id": map[string]any{"type": "string", "description": "ID записи"},
				},
				"required": []string{"blog_id", "post_id"},
			},
		},
	}
}

func GetToolDescription(name string) string {
	for _, t := range GetAllTools() {
		if t.Name == name {
			return t.Description
		}
	}
	return ""
}

func GetAllToolDescriptions() map[string]string {
	result := make(map[string]string)
	for _, t := range GetAllTools() {
		result[t.Name] = t.Description
	}
	return result
}

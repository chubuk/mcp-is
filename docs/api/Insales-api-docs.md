# InSales API documentation

Документация InSales API с примерами запросов и ответов.

| | |
|--|--|
| Ресурсов | 57 |
| Операций | 327 |
| Авторизация | HTTP Basic (примеры из песочницы InSales — не боевые ключи) |
| Краткий справочник | [insales-api-reference.md](insales-api-reference.md) |
| OpenAPI MCP IS | [openapi-mcp-is.yaml](openapi-mcp-is.yaml) |

## Содержание

- [Accessories](#accessories) — 5 операций
- [AccessoryValues](#accessoryvalues) — 8 операций
- [Account](#account) — 2 операции
- [ApplicationAction](#applicationaction) — 5 операций
- [ApplicationCharge](#applicationcharge) — 4 операции
- [ApplicationWidget](#applicationwidget) — 5 операций
- [Article](#article) — 5 операций
- [Blog](#blog) — 5 операций
- [BonusSystemTransaction](#bonussystemtransaction) — 4 операции
- [Category](#category) — 5 операций
- [Characteristic](#characteristic) — 5 операций
- [Client](#client) — 6 операций
- [ClientGroup](#clientgroup) — 5 операций
- [Collect](#collect) — 7 операций
- [Collection](#collection) — 5 операций
- [CollectionField](#collectionfield) — 8 операций
- [CollectionFieldValue](#collectionfieldvalue) — 5 операций
- [CollectionFilter](#collectionfilter) — 5 операций
- [CustomStatus](#customstatus) — 5 операций
- [DeliveryVariant](#deliveryvariant) — 9 операций
- [DiscountCode](#discountcode) — 5 операций
- [DiscountCollectionsGroup](#discountcollectionsgroup) — 3 операции
- [Domain](#domain) — 5 операций
- [ExternalDiscount](#externaldiscount) — 5 операций
- [Field](#field) — 7 операций
- [File](#file) — 5 операций
- [Image](#image) — 6 операций
- [JsTag](#jstag) — 6 операций
- [Marketplace](#marketplace) — 9 операций
- [Menu](#menu) — 5 операций
- [MenuItem](#menuitem) — 8 операций
- [OAuthApp](#oauthapp) — 5 операций
- [OptionName](#optionname) — 5 операций
- [OptionValue](#optionvalue) — 8 операций
- [Order](#order) — 16 операций
- [Outlet](#outlet) — 5 операций
- [Page](#page) — 5 операций
- [PaymentGateway](#paymentgateway) — 6 операций
- [PriceKind](#pricekind) — 5 операций
- [Product](#product) — 8 операций
- [ProductField](#productfield) — 8 операций
- [ProductFieldValue](#productfieldvalue) — 9 операций
- [Property](#property) — 5 операций
- [ProviderToken](#providertoken) — 1 операция
- [RecurringApplicationCharge](#recurringapplicationcharge) — 5 операций
- [Redirect](#redirect) — 5 операций
- [Review](#review) — 5 операций
- [SalesChannel](#saleschannel) — 4 операции
- [Similar](#similar) — 7 операций
- [Supplementary](#supplementary) — 3 операции
- [Unit](#unit) — 1 операция
- [UserTask](#usertask) — 5 операций
- [Variant](#variant) — 7 операций
- [VariantField](#variantfield) — 8 операций
- [VariantFieldValue](#variantfieldvalue) — 9 операций
- [Warehouse](#warehouse) — 5 операций
- [Webhook](#webhook) — 5 операций

---

## Accessories

| Операция | Метод и путь |
|----------|--------------|
| [Create accessory](#create-accessory) | `POST /admin/accessories.json` |
| [Destroy accessory](#destroy-accessory) | `DELETE /admin/accessories/:accessory_id.json` |
| [Get accessories](#get-accessories) | `GET /admin/accessories.json` |
| [Get accessory](#get-accessory) | `GET /admin/accessories/:accessory_id.json` |
| [Update accessory](#update-accessory) | `PUT /admin/accessories/:accessory_id.json` |

### Create accessory

`POST /admin/accessories.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `accessory[name]` | да | Наименование |
| `accessory[max_count]` | да | Максимальное количество опций в заказе |
| `accessory[min_count]` | да | Минимальное количество опций в заказе |
| `accessory[permalink]` | да | Идентификатор |

**Запрос**

```http
POST /admin/accessories.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "accessory": {
    "name": "Опция 1",
    "max_count": 5,
    "min_count": 5,
    "permalink": "opcia-1"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 40,
  "account_id": 1,
  "name_translations": [
    "Опция 1"
  ],
  "min_count": 5,
  "max_count": 5,
  "created_at": "2026-07-25T00:59:28.890+03:00",
  "updated_at": "2026-07-25T00:59:28.890+03:00",
  "back_office_name": "Опция 1",
  "permalink": "opcia-1",
  "name": "Опция 1"
}
```

### Destroy accessory

`DELETE /admin/accessories/:accessory_id.json`

**Запрос**

```http
DELETE /admin/accessories/68.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 68,
  "account_id": 1,
  "name_translations": [
    "Тяжёлый Бумажный Свитер"
  ],
  "min_count": 2,
  "max_count": 7,
  "created_at": "2026-07-25T00:59:29.540+03:00",
  "updated_at": "2026-07-25T00:59:29.540+03:00",
  "back_office_name": "Тяжёлый Хлопковый Автомобиль",
  "permalink": "quibusdam_expedita",
  "name": "Тяжёлый Бумажный Свитер"
}
```

### Get accessories

`GET /admin/accessories.json`

**Запрос**

```http
GET /admin/accessories.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 7,
    "account_id": 1,
    "name_translations": [
      "Опция 0"
    ],
    "min_count": 2,
    "max_count": 7,
    "created_at": "2026-07-25T00:59:27.990+03:00",
    "updated_at": "2026-07-25T00:59:27.990+03:00",
    "back_office_name": "Потрясающий Медный Сумка",
    "permalink": "nam_et",
    "name": "Опция 0"
  },
  {
    "id": 8,
    "account_id": 1,
    "name_translations": [
      "Опция 1"
    ],
    "min_count": 5,
    "max_count": 9,
    "created_at": "2026-07-25T00:59:28.003+03:00",
    "updated_at": "2026-07-25T00:59:28.003+03:00",
    "back_office_name": "Прочный Кожаный Сумка",
    "permalink": "sunt_dolor",
    "name": "Опция 1"
  },
  {
    "id": 9,
    "account_id": 1,
    "name_translations": [
      "Опция 2"
    ],
    "min_count": 3,
    "max_count": 6,
    "created_at": "2026-07-25T00:59:28.019+03:00",
    "updated_at": "2026-07-25T00:59:28.019+03:00",
    "back_office_name": "Маленький Медный Сумка",
    "permalink": "et_a",
    "name": "Опция 2"
  },
  {
    "id": 10,
    "account_id": 1,
    "name_translations": [
      "Опция 0"
    ],
    "min_count": 3,
    "max_count": 5,
    "created_at": "2026-07-25T00:59:28.032+03:00",
    "updated_at": "2026-07-25T00:59:28.032+03:00",
    "back_office_name": "Практичный Бетонный Ремень",
    "permalink": "quo_dicta",
    "name": "Опция 0"
  },
  {
    "id": 11,
    "account_id": 1,
    "name_translations": [
      "Опция 1"
    ],
    "min_count": 4,
    "max_count": 4,
    "created_at": "2026-07-25T00:59:28.043+03:00",
    "updated_at": "2026-07-25T00:59:28.043+03:00",
    "back_office_name": "Огромный Шерстяной Ремень",
    "permalink": "minima_veniam",
    "name": "Опция 1"
  },
  {
    "id": 12,
    "account_id": 1,
    "name_translations": [
      "Опция 2"
    ],
    "min_count": 1,
    "max_count": 5,
    "created_at": "2026-07-25T00:59:28.054+03:00",
    "updated_at": "2026-07-25T00:59:28.054+03:00",
    "back_office_name": "Практичный Мраморный Часы",
    "permalink": "unde_tempora",
    "name": "Опция 2"
  }
]
```

### Get accessory

`GET /admin/accessories/:accessory_id.json`

**Запрос**

```http
GET /admin/accessories/26.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 26,
  "account_id": 1,
  "name_translations": [
    "Синергетический Мраморный Свитер"
  ],
  "min_count": 2,
  "max_count": 10,
  "created_at": "2026-07-25T00:59:28.448+03:00",
  "updated_at": "2026-07-25T00:59:28.448+03:00",
  "back_office_name": "Грубый Гранитный Бумажник",
  "permalink": "deleniti_corrupti",
  "name": "Синергетический Мраморный Свитер"
}
```

### Update accessory

`PUT /admin/accessories/:accessory_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `accessory[name]` |  | Наименование |
| `accessory[max_count]` |  | Максимальное количество опций в заказе |
| `accessory[min_count]` |  | Минимальное количество опций в заказе |
| `accessory[permalink]` |  | Идентификатор |

**Запрос**

```http
PUT /admin/accessories/54.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "accessory": {
    "name": "Опция 1",
    "max_count": 7,
    "min_count": 4,
    "permalink": "opcia-1"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "account_id": 1,
  "name": "Опция 1",
  "min_count": 4,
  "max_count": 7,
  "permalink": "opcia-1",
  "id": 54,
  "back_office_name": "Удовлетворительный Гранитный Берет",
  "name_translations": [
    "Опция 1"
  ],
  "created_at": "2026-07-25T00:59:29.183+03:00",
  "updated_at": "2026-07-25T00:59:29.255+03:00"
}
```

---

## AccessoryValues

| Операция | Метод и путь |
|----------|--------------|
| [Create accessory value](#create-accessory-value) | `POST /admin/accessories/:accessory_id/values.json` |
| [Create product accessory link](#create-product-accessory-link) | `POST /admin/products/:product_id/product_accessory_links.json` |
| [Destroy accessory value](#destroy-accessory-value) | `DELETE /admin/accessories/:accessory_id/values/:accessory_value_id.json` |
| [Destroy product accessory link](#destroy-product-accessory-link) | `DELETE /admin/products/:product_id/product_accessory_links/:product_accessory_link_id.json` |
| [Get accessory value](#get-accessory-value) | `GET /admin/accessories/:accessory_id/values/:accessory_value_id.json` |
| [Get accessory values](#get-accessory-values) | `GET /admin/accessories/:accessory_id/values.json` |
| [Get product accessory links](#get-product-accessory-links) | `GET /admin/products/:product_id/product_accessory_links.json` |
| [Update accessory value](#update-accessory-value) | `PUT /admin/accessories/:accessory_id/values/:accessory_value_id.json` |

### Create accessory value

`POST /admin/accessories/:accessory_id/values.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `accessory_value[name]` | да | Название |
| `accessory_value[price]` | да | Цена |
| `accessory_value[accounting_variant_id]` |  | Идентификатор списываемого товара |

**Запрос**

```http
POST /admin/accessories/74/values.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "accessory_value": {
    "name": "Значение опции",
    "price": 100.0,
    "accounting_variant_id": null
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 44,
  "account_id": 1,
  "accessory_id": 74,
  "name_translations": [
    "Значение опции"
  ],
  "price": "100.0",
  "created_at": "2026-07-25T00:59:30.623+03:00",
  "updated_at": "2026-07-25T00:59:30.623+03:00",
  "accounting_variant_id": null,
  "name": "Значение опции"
}
```

### Create product accessory link

`POST /admin/products/:product_id/product_accessory_links.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_accessory_link[accessory_id]` | да | Идентификатор опции |

**Запрос**

```http
POST /admin/products/293977313/product_accessory_links.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_accessory_link": {
    "accessory_id": 82
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "account_id": 1,
  "accessory_id": 82,
  "created_at": "2026-07-25T00:59:49.515+03:00",
  "updated_at": "2026-07-25T00:59:49.515+03:00",
  "product_id": 293977313,
  "accessory": {
    "id": 82,
    "account_id": 1,
    "name_translations": [
      "Опция 1"
    ],
    "min_count": 1,
    "max_count": 8,
    "created_at": "2026-07-25T00:59:49.384+03:00",
    "updated_at": "2026-07-25T00:59:49.384+03:00",
    "back_office_name": "Синергетический Бумажный Свитер",
    "permalink": "opcia-1",
    "name": "Опция 1"
  },
  "product": {
    "category_id": 1,
    "created_at": "2026-07-25T00:59:49.000+03:00",
    "updated_at": "2026-07-25T00:59:49.000+03:00",
    "is_hidden": false,
    "account_id": 1,
    "title_translations": [
      "product 4"
    ],
    "available": null,
    "archived": false,
    "key": null,
    "price_min_cached": null,
    "price_max_cached": null,
    "price_min_available_cached": null,
    "price_max_available_cached": null,
    "canonical_url_collection_id": null,
    "collections_id": [],
    "characteristics_id": [],
    "custom_layout": null,
    "custom_template": null,
    "old_price_min_cached": null,
    "old_price_max_cached": null,
    "old_price_min_available_cached": null,
    "old_price_max_available_cached": null,
    "permalink": "product-4",
    "unit": "pce",
    "deleted": false,
    "cost_price_min_cached": null,
    "cost_price_max_cached": null,
    "stock_currency_id": 1,
    "sort_weight": null,
    "bundle": null,
    "ignore_discounts": null,
    "vat": 10,
    "dimensions": null,
    "fiscal_product_type": 1,
    "popularity": null,
    "available_at_warehouses": [],
    "auto_assign_bundle_price": null,
    "demo": null,
    "sales_channels_id": null,
    "demo_thematic": null,
    "reviews_average_rating_cached": null,
    "reviews_count_cached": null,
    "merge_product_id": null,
    "merge_group_product_id": null,
    "merge_group_hidden": null,
    "no_delivery": false,
    "id": 293977313,
    "main_image_id": null,
    "title": "product 4"
  }
}
```

### Destroy accessory value

`DELETE /admin/accessories/:accessory_id/values/:accessory_value_id.json`

**Запрос**

```http
DELETE /admin/accessories/78/values/66.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 66,
  "account_id": 1,
  "accessory_id": 78,
  "name_translations": [
    "Значение опции"
  ],
  "price": "9.56",
  "created_at": "2026-07-25T00:59:31.243+03:00",
  "updated_at": "2026-07-25T00:59:31.243+03:00",
  "accounting_variant_id": null,
  "name": "Значение опции"
}
```

### Destroy product accessory link

`DELETE /admin/products/:product_id/product_accessory_links/:product_accessory_link_id.json`

**Запрос**

```http
DELETE /admin/products/293977329/product_accessory_links/4.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 4,
  "account_id": 1,
  "accessory_id": 84,
  "created_at": "2026-07-25T00:59:49.815+03:00",
  "updated_at": "2026-07-25T00:59:49.815+03:00",
  "product_id": 293977329,
  "accessory": {
    "id": 84,
    "account_id": 1,
    "name_translations": [
      "Опция 1"
    ],
    "min_count": 7,
    "max_count": 9,
    "created_at": "2026-07-25T00:59:49.761+03:00",
    "updated_at": "2026-07-25T00:59:49.761+03:00",
    "back_office_name": "Практичный Гранитный Тарелка",
    "permalink": "opcia-1",
    "name": "Опция 1"
  },
  "product": {
    "category_id": 1,
    "created_at": "2026-07-25T00:59:49.000+03:00",
    "updated_at": "2026-07-25T00:59:49.000+03:00",
    "is_hidden": false,
    "account_id": 1,
    "title_translations": [
      "product 6"
    ],
    "available": null,
    "archived": false,
    "key": null,
    "price_min_cached": null,
    "price_max_cached": null,
    "price_min_available_cached": null,
    "price_max_available_cached": null,
    "canonical_url_collection_id": null,
    "collections_id": [],
    "characteristics_id": [],
    "custom_layout": null,
    "custom_template": null,
    "old_price_min_cached": null,
    "old_price_max_cached": null,
    "old_price_min_available_cached": null,
    "old_price_max_available_cached": null,
    "permalink": "product-6",
    "unit": "pce",
    "deleted": false,
    "cost_price_min_cached": null,
    "cost_price_max_cached": null,
    "stock_currency_id": 1,
    "sort_weight": null,
    "bundle": null,
    "ignore_discounts": null,
    "vat": 10,
    "dimensions": null,
    "fiscal_product_type": 1,
    "popularity": null,
    "available_at_warehouses": [],
    "auto_assign_bundle_price": null,
    "demo": null,
    "sales_channels_id": null,
    "demo_thematic": null,
    "reviews_average_rating_cached": null,
    "reviews_count_cached": null,
    "merge_product_id": null,
    "merge_group_product_id": null,
    "merge_group_hidden": null,
    "no_delivery": false,
    "id": 293977329,
    "main_image_id": null,
    "title": "product 6"
  }
}
```

### Get accessory value

`GET /admin/accessories/:accessory_id/values/:accessory_value_id.json`

**Запрос**

```http
GET /admin/accessories/72/values/22.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 22,
  "account_id": 1,
  "accessory_id": 72,
  "name_translations": [
    "Значение опции"
  ],
  "price": "6.09",
  "created_at": "2026-07-25T00:59:30.193+03:00",
  "updated_at": "2026-07-25T00:59:30.193+03:00",
  "accounting_variant_id": null,
  "name": "Значение опции"
}
```

### Get accessory values

`GET /admin/accessories/:accessory_id/values.json`

**Запрос**

```http
GET /admin/accessories/70/values.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 8,
    "account_id": 1,
    "accessory_id": 70,
    "name_translations": [
      "Значение опции"
    ],
    "price": "6.17",
    "created_at": "2026-07-25T00:59:29.887+03:00",
    "updated_at": "2026-07-25T00:59:29.887+03:00",
    "accounting_variant_id": null,
    "name": "Значение опции"
  },
  {
    "id": 9,
    "account_id": 1,
    "accessory_id": 70,
    "name_translations": [
      "Значение опции 0"
    ],
    "price": "6.98",
    "created_at": "2026-07-25T00:59:29.893+03:00",
    "updated_at": "2026-07-25T00:59:29.893+03:00",
    "accounting_variant_id": null,
    "name": "Значение опции 0"
  },
  {
    "id": 10,
    "account_id": 1,
    "accessory_id": 70,
    "name_translations": [
      "Значение опции 1"
    ],
    "price": "5.83",
    "created_at": "2026-07-25T00:59:29.901+03:00",
    "updated_at": "2026-07-25T00:59:29.901+03:00",
    "accounting_variant_id": null,
    "name": "Значение опции 1"
  },
  {
    "id": 11,
    "account_id": 1,
    "accessory_id": 70,
    "name_translations": [
      "Значение опции 2"
    ],
    "price": "1.68",
    "created_at": "2026-07-25T00:59:29.907+03:00",
    "updated_at": "2026-07-25T00:59:29.907+03:00",
    "accounting_variant_id": null,
    "name": "Значение опции 2"
  },
  {
    "id": 12,
    "account_id": 1,
    "accessory_id": 70,
    "name_translations": [
      "Значение опции 0"
    ],
    "price": "8.4",
    "created_at": "2026-07-25T00:59:29.914+03:00",
    "updated_at": "2026-07-25T00:59:29.914+03:00",
    "accounting_variant_id": null,
    "name": "Значение опции 0"
  },
  {
    "id": 13,
    "account_id": 1,
    "accessory_id": 70,
    "name_translations": [
      "Значение опции 1"
    ],
    "price": "6.86",
    "created_at": "2026-07-25T00:59:29.921+03:00",
    "updated_at": "2026-07-25T00:59:29.921+03:00",
    "accounting_variant_id": null,
    "name": "Значение опции 1"
  },
  {
    "id": 14,
    "account_id": 1,
    "accessory_id": 70,
    "name_translations": [
      "Значение опции 2"
    ],
    "price": "5.39",
    "created_at": "2026-07-25T00:59:29.927+03:00",
    "updated_at": "2026-07-25T00:59:29.927+03:00",
    "accounting_variant_id": null,
    "name": "Значение опции 2"
  }
]
```

### Get product accessory links

`GET /admin/products/:product_id/product_accessory_links.json`

**Запрос**

```http
GET /admin/products/293977297/product_accessory_links.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[]
```

### Update accessory value

`PUT /admin/accessories/:accessory_id/values/:accessory_value_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `accessory_value[name]` |  | Название |
| `accessory_value[price]` |  | Цена |
| `accessory_value[accounting_variant_id]` |  | Идентификатор списываемого товара |

**Запрос**

```http
PUT /admin/accessories/76/values/52.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "accessory_value": {
    "name": "Обновлённое значение опции",
    "price": 200.0,
    "accounting_variant_id": null
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "accessory_id": 76,
  "name": "Обновлённое значение опции",
  "price": "200.0",
  "accounting_variant_id": null,
  "name_translations": [
    "Обновлённое значение опции"
  ],
  "id": 52,
  "account_id": 1,
  "created_at": "2026-07-25T00:59:30.888+03:00",
  "updated_at": "2026-07-25T00:59:31.012+03:00"
}
```

---

## Account

| Операция | Метод и путь |
|----------|--------------|
| [Get account](#get-account) | `GET admin/account.json` |
| [Update account](#update-account) | `PUT /admin/account.json` |

### Get account

`GET admin/account.json`

**Запрос**

```http
GET admin/account.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "blocked": false,
  "city": null,
  "contact_phone": "+79993332211",
  "country": null,
  "sms_notification_phone": null,
  "enable_cart_discounts": false,
  "enable_client_discounts": false,
  "enable_group_discounts": false,
  "enable_order_discounts": false,
  "subdomain": "quentin",
  "notification_email": null,
  "organization": "InSales, ltd",
  "hide_items_out_of_stock": false,
  "state": null,
  "inviter_id": null,
  "paid_till": "2009-03-02",
  "created_at": "2008-12-31T03:00:00.000+03:00",
  "registration_type_id": 0,
  "default_vat": 10,
  "email": "somone@anyone.com",
  "icq": null,
  "phone": "+7 812 309 51 21",
  "title": "My account",
  "minimum_items_price": null,
  "stock_currency_exchange_rate": "30.0",
  "client_cookies_whitelist": null,
  "main_host": "quentin.insales.io:3002",
  "main_host_protocol": "http",
  "next_order_number": 1006,
  "time_zone": null,
  "owner": {
    "id": 1,
    "name": "quentin",
    "email": "quentin@example.com",
    "created_at": "2026-07-20T00:58:27.000+03:00"
  }
}
```

### Update account

`PUT /admin/account.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `account[client_cookies_whitelist]` |  | cookies whitelist |

**Запрос**

```http
PUT /admin/account.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "account": {
    "client_cookies_whitelist": "my_cookie"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

---

## ApplicationAction

| Операция | Метод и путь |
|----------|--------------|
| [Create application_action](#create-application-action) | `POST /admin/application_actions.json` |
| [Destroy application_action](#destroy-application-action) | `DELETE /admin/application_actions/:application_action_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/application_actions.json` |
| [Get application_action](#get-application-action) | `GET /admin/application_actions/:application_action_id.json` |
| [Update application_action](#update-application-action) | `PUT /admin/application_actions/:application_action_id.json` |

### Create application_action

`POST /admin/application_actions.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_action[title]` | да | title to show in backoffice |
| `application_action[handle]` | да | unique string identifier |
| `application_action[entity]` | да | string, defines entity that action belongs to: order, product, client, discount_code, file, page, article, collection_filter |
| `application_action[url]` | да | URL to which POST request with serialized entities will be made |

**Запрос**

```http
POST /admin/application_actions.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "application_action": {
    "title": "backoffice title",
    "handle": "unique_handle",
    "entity": "product",
    "url": "http://www.yandex.ru"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 17,
  "title": "backoffice title",
  "url": "http://www.yandex.ru",
  "handle": "unique_handle",
  "entity": "product"
}
```

### Destroy application_action

`DELETE /admin/application_actions/:application_action_id.json`

**Запрос**

```http
DELETE /admin/application_actions/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/application_actions.json`

Content-Type: application/json Get application actions

**Запрос**

```http
GET /admin/application_actions.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "title": "test",
    "url": "http://www.google.ru",
    "handle": "test",
    "entity": "order"
  }
]
```

### Get application_action

`GET /admin/application_actions/:application_action_id.json`

**Запрос**

```http
GET /admin/application_actions/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "title": "test",
  "url": "http://www.google.ru",
  "handle": "test",
  "entity": "order"
}
```

### Update application_action

`PUT /admin/application_actions/:application_action_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_action[title]` |  | title to show in backoffice |
| `application_action[handle]` |  | unique string identifier |
| `application_action[entity]` |  | string, defines entity that action belongs to: order, product, client, discount_code, file, page, article, collection_filter |
| `application_action[url]` |  | URL to which POST request with serialized entities will be made |

**Запрос**

```http
PUT /admin/application_actions/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "application_action": {
    "title": "new  title",
    "handle": "new_unique_handle",
    "entity": "client",
    "url": "http://www.yandex.ru"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "title": "new  title",
  "url": "http://www.yandex.ru",
  "handle": "new_unique_handle",
  "entity": "client"
}
```

---

## ApplicationCharge

| Операция | Метод и путь |
|----------|--------------|
| [Create application_charge](#create-application-charge) | `POST /admin/application_charges.json` |
| [Decline application_charge](#decline-application-charge) | `POST /admin/application_charges/:application_charge_id/decline.json` |
| [Get application_charge](#get-application-charge) | `GET /admin/application_charges/:application_charge_id.json` |
| [Get application_charges](#get-application-charges) | `GET /admin/application_charges.json` |

### Create application_charge

`POST /admin/application_charges.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_charge[name]` | да | purpose of payment |
| `application_charge[price]` | да | bill amount |
| `application_charge[return_url]` |  | url for notification about payment |
| `application_charge[test]` |  | flag for debugging purpose; if true, bill can be confirmed without payment; false by default |

**Запрос**

```http
POST /admin/application_charges.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "application_charge": {
    "name": "Payment for premium",
    "price": 300,
    "return_url": "http://app.my/check_payment_url/%id_in_my_application%",
    "test": true
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 8,
  "name": "Payment for premium",
  "price": "300.0",
  "return_url": "http://app.my/check_payment_url/%id_in_my_application%",
  "test": true,
  "created_at": "2026-07-25T00:58:38.593+03:00",
  "updated_at": "2026-07-25T00:58:38.593+03:00",
  "confirmation_url": "https://quentin.insales.io:3002/admin/application_charges/8",
  "status": "pending"
}
```

### Decline application_charge

`POST /admin/application_charges/:application_charge_id/decline.json`

Only bill with status ‘pending’ can be declined.

**Запрос**

```http
POST /admin/application_charges/10/decline.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 10,
  "name": "Pay",
  "price": "180.0",
  "return_url": "http://applicaiton.my/check_payment_url/%id_in_my_application%",
  "test": null,
  "created_at": "2026-07-25T00:58:38.786+03:00",
  "updated_at": "2026-07-25T00:58:38.786+03:00",
  "confirmation_url": "https://quentin.insales.io:3002/admin2/invoices/8",
  "status": "declined"
}
```

### Get application_charge

`GET /admin/application_charges/:application_charge_id.json`

**Запрос**

```http
GET /admin/application_charges/4.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 4,
  "name": "Pay",
  "price": "180.0",
  "return_url": "http://applicaiton.my/check_payment_url/%id_in_my_application%",
  "test": null,
  "created_at": "2026-07-25T00:58:38.146+03:00",
  "updated_at": "2026-07-25T00:58:38.146+03:00",
  "confirmation_url": "https://quentin.insales.io:3002/admin2/invoices/4",
  "status": "pending"
}
```

### Get application_charges

`GET /admin/application_charges.json`

**Запрос**

```http
GET /admin/application_charges.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 2,
    "name": "Pay",
    "price": "180.0",
    "return_url": "http://applicaiton.my/check_payment_url/%id_in_my_application%",
    "test": null,
    "created_at": "2026-07-25T00:58:37.873+03:00",
    "updated_at": "2026-07-25T00:58:37.873+03:00",
    "confirmation_url": "https://quentin.insales.io:3002/admin2/invoices/2",
    "status": "pending"
  }
]
```

---

## ApplicationWidget

| Операция | Метод и путь |
|----------|--------------|
| [Create application_widget](#create-application-widget) | `POST /admin/application_widgets.json` |
| [Destroy application_widget](#destroy-application-widget) | `DELETE /admin/application_widgets/:application_widget_id.json` |
| [Get application_widget](#get-application-widget) | `GET /admin/application_widgets/:application_widget_id.json` |
| [Get application_widgets](#get-application-widgets) | `GET /admin/application_widgets.json` |
| [Update application_widget](#update-application-widget) | `PUT /admin/application_widgets/:application_widget_id.json` |

### Create application_widget

`POST /admin/application_widgets.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_widget[code]` | да | html or javascript code of widget |
| `application_widget[height]` | да | height of iframe block |
| `application_widget[page_type]` |  | product or order |

**Запрос**

```http
POST /admin/application_widgets.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "application_widget": {
    "code": "some html or javascript code",
    "height": 70,
    "page_type": "order"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "code": "some html or javascript code",
  "created_at": "2026-07-25T00:58:39.705+03:00",
  "id": 57,
  "height": 70,
  "page_type": "order"
}
```

### Destroy application_widget

`DELETE /admin/application_widgets/:application_widget_id.json`

**Запрос**

```http
DELETE /admin/application_widgets/89.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get application_widget

`GET /admin/application_widgets/:application_widget_id.json`

**Запрос**

```http
GET /admin/application_widgets/25.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "code": "html/js code",
  "created_at": "2026-07-25T00:58:39.265+03:00",
  "id": 25,
  "height": 50,
  "page_type": "product"
}
```

### Get application_widgets

`GET /admin/application_widgets.json`

**Запрос**

```http
GET /admin/application_widgets.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "code": "html/js code",
    "created_at": "2026-07-25T00:58:39.070+03:00",
    "id": 9,
    "height": 50,
    "page_type": "product"
  }
]
```

### Update application_widget

`PUT /admin/application_widgets/:application_widget_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_widget[code]` |  | html or javascript code of widget |
| `application_widget[height]` |  | height of iframe block |
| `application_widget[page_type]` |  | product or order |

**Запрос**

```http
PUT /admin/application_widgets/73.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "application_widget": {
    "code": "code changed",
    "height": 80,
    "page_type": "order"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "code": "code changed",
  "created_at": "2026-07-25T00:58:39.906+03:00",
  "id": 73,
  "height": 80,
  "page_type": "order"
}
```

---

## Article

| Операция | Метод и путь |
|----------|--------------|
| [Create article](#create-article) | `POST /admin/blogs/:blog_id/articles.json` |
| [Delete article](#delete-article) | `DELETE /admin/blogs/:blog_id/articles/:article_id.json` |
| [Get article](#get-article) | `GET /admin/blogs/:blog_id/articles/:article_id.json` |
| [Get articles list](#get-articles-list) | `GET /admin/blogs/:blog_id/articles.json` |
| [Update article](#update-article) | `PUT /admin/blogs/:blog_id/articles/:article_id.json` |

### Create article

`POST /admin/blogs/:blog_id/articles.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `article[title]` | да | article title |
| `article[content]` | да | article content |
| `article[notice]` |  | notice for article |
| `article[published_at]` | да | publication date, the article is invisible to users while publication date > current time |
| `article[author]` | да | article author name |
| `article[permalink]` |  | permalink on blog/* |
| `article[html_title]` |  | html title tag content |
| `article[meta_keywords]` |  | html meta words on page |
| `article[meta_description]` |  | html meta description |
| `article[image_attributes]` |  | image for article |
| `article[related_products_attributes]` |  | goods to attach or detach. An array of hashes with the keys: product_id, _destroy. |
| `article[all_tags]` |  | array of all tags for article, max 10 tags |
| `article[pinned]` |  | set to true to pin article: it will be always on the top of the blog; only one article can be pinned, so when you pin article, previously pinned article will be unpinned |

**Запрос**

```http
POST /admin/blogs/1/articles.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "article": {
    "title": "new article",
    "content": "<p>content of my article<p>",
    "notice": "some notice",
    "published_at": "2014-09-09",
    "author": "My Name is panda",
    "permalink": "first_article",
    "html_title": "HTML TAG TITLE",
    "meta_keywords": "new cool products",
    "meta_description": "It's a article about our new products",
    "image_attributes": {
      "src": "https://abc.myinsales.ru/images/no_image_compact.jpg"
    },
    "related_products_attributes": [
      {
        "product_id": 1
      }
    ],
    "all_tags": [
      "test3",
      "test4"
    ],
    "pinned": true
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20,
  "created_at": "2026-07-25T01:00:22.000+03:00",
  "updated_at": "2026-07-25T01:00:22.000+03:00",
  "published_at": "2014-09-09T00:00:00.000+04:00",
  "pinned": true,
  "custom_template": null,
  "title": "new article",
  "author": "My Name is panda",
  "permalink": "first_article",
  "html_title": "HTML TAG TITLE",
  "meta_keywords": "new cool products",
  "meta_description": "It's a article about our new products",
  "image": {
    "id": 17,
    "created_at": "2026-07-25T01:00:22.000+03:00",
    "original_url": "/images/loading.gif"
  },
  "tags": [
    {
      "title": "test3"
    },
    {
      "title": "test4"
    }
  ],
  "related_products": [
    {
      "category_id": 1,
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "updated_at": "2015-07-27T17:05:35.000+03:00",
      "is_hidden": false,
      "account_id": 1,
      "title_translations": [
        "Cool t-shirt"
      ],
      "available": false,
      "archived": false,
      "key": null,
      "price_min_cached": "10.0",
      "price_max_cached": "10.0",
      "price_min_available_cached": null,
      "price_max_available_cached": null,
      "canonical_url_collection_id": 1,
      "collections_id": [
        1,
        3,
        8
      ],
      "characteristics_id": [
        7,
        9,
        12,
        13
      ],
      "custom_layout": null,
      "custom_template": null,
      "old_price_min_cached": null,
      "old_price_max_cached": null,
      "old_price_min_available_cached": null,
      "old_price_max_available_cached": null,
      "permalink": "футболка",
      "unit": "pce",
      "deleted": false,
      "cost_price_min_cached": null,
      "cost_price_max_cached": null,
      "stock_currency_id": 1,
      "sort_weight": null,
      "bundle": null,
      "ignore_discounts": null,
      "vat": -1,
      "dimensions": "10x20x30",
      "fiscal_product_type": 1,
      "popularity": null,
      "available_at_warehouses": [],
      "auto_assign_bundle_price": null,
      "demo": null,
      "sales_channels_id": null,
      "demo_thematic": null,
      "reviews_average_rating_cached": null,
      "reviews_count_cached": null,
      "merge_product_id": null,
      "merge_group_product_id": null,
      "merge_group_hidden": null,
      "no_delivery": false,
      "id": 1,
      "main_image_id": null,
      "title": "Cool t-shirt"
    }
  ],
  "notice": "some notice",
  "content": "<p>content of my article<p>"
}
```

### Delete article

`DELETE /admin/blogs/:blog_id/articles/:article_id.json`

**Запрос**

```http
DELETE /admin/blogs/1/articles/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
Content-Type: application/x-www-form-urlencoded
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get article

`GET /admin/blogs/:blog_id/articles/:article_id.json`

**Запрос**

```http
GET /admin/blogs/1/articles/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2026-07-25T01:00:21.000+03:00",
  "published_at": "2017-03-01T07:53:15.000+03:00",
  "pinned": false,
  "custom_template": null,
  "title": "Второая Статья",
  "author": "Ivan",
  "permalink": "statya-2",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "image": {
    "id": 1,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "original_url": "http://minio.insales.io:9000/insales-test-sber/images/articles/1/1/1/article_image_1.jpg"
  },
  "tags": [
    {
      "title": "test"
    },
    {
      "title": "test2"
    }
  ],
  "related_products": [
    {
      "category_id": 1,
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "updated_at": "2015-07-27T17:05:35.000+03:00",
      "is_hidden": false,
      "account_id": 1,
      "title_translations": [
        "Cool t-shirt"
      ],
      "available": false,
      "archived": false,
      "key": null,
      "price_min_cached": "10.0",
      "price_max_cached": "10.0",
      "price_min_available_cached": null,
      "price_max_available_cached": null,
      "canonical_url_collection_id": 1,
      "collections_id": [
        1,
        3,
        8
      ],
      "characteristics_id": [
        7,
        9,
        12,
        13
      ],
      "custom_layout": null,
      "custom_template": null,
      "old_price_min_cached": null,
      "old_price_max_cached": null,
      "old_price_min_available_cached": null,
      "old_price_max_available_cached": null,
      "permalink": "футболка",
      "unit": "pce",
      "deleted": false,
      "cost_price_min_cached": null,
      "cost_price_max_cached": null,
      "stock_currency_id": 1,
      "sort_weight": null,
      "bundle": null,
      "ignore_discounts": null,
      "vat": -1,
      "dimensions": "10x20x30",
      "fiscal_product_type": 1,
      "popularity": null,
      "available_at_warehouses": [],
      "auto_assign_bundle_price": null,
      "demo": null,
      "sales_channels_id": null,
      "demo_thematic": null,
      "reviews_average_rating_cached": null,
      "reviews_count_cached": null,
      "merge_product_id": null,
      "merge_group_product_id": null,
      "merge_group_hidden": null,
      "no_delivery": false,
      "id": 1,
      "main_image_id": null,
      "title": "Cool t-shirt"
    }
  ],
  "notice": "Второй анонс",
  "content": "Привет"
}
```

### Get articles list

`GET /admin/blogs/:blog_id/articles.json`

**Запрос**

```http
GET /admin/blogs/1/articles.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 2,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T01:00:21.000+03:00",
    "published_at": "2017-03-01T07:53:15.000+03:00",
    "pinned": false,
    "custom_template": null,
    "title": "Второая Статья",
    "author": "Ivan",
    "permalink": "statya-2",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "image": {
      "id": 1,
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "original_url": "http://minio.insales.io:9000/insales-test-sber/images/articles/1/1/1/article_image_1.jpg"
    },
    "tags": [
      {
        "title": "test"
      },
      {
        "title": "test2"
      }
    ],
    "related_products": [
      {
        "category_id": 1,
        "created_at": "2026-07-25T00:58:09.000+03:00",
        "updated_at": "2015-07-27T17:05:35.000+03:00",
        "is_hidden": false,
        "account_id": 1,
        "title_translations": [
          "Cool t-shirt"
        ],
        "available": false,
        "archived": false,
        "key": null,
        "price_min_cached": "10.0",
        "price_max_cached": "10.0",
        "price_min_available_cached": null,
        "price_max_available_cached": null,
        "canonical_url_collection_id": 1,
        "collections_id": [
          1,
          3,
          8
        ],
        "characteristics_id": [
          7,
          9,
          12,
          13
        ],
        "custom_layout": null,
        "custom_template": null,
        "old_price_min_cached": null,
        "old_price_max_cached": null,
        "old_price_min_available_cached": null,
        "old_price_max_available_cached": null,
        "permalink": "футболка",
        "unit": "pce",
        "deleted": false,
        "cost_price_min_cached": null,
        "cost_price_max_cached": null,
        "stock_currency_id": 1,
        "sort_weight": null,
        "bundle": null,
        "ignore_discounts": null,
        "vat": -1,
        "dimensions": "10x20x30",
        "fiscal_product_type": 1,
        "popularity": null,
        "available_at_warehouses": [],
        "auto_assign_bundle_price": null,
        "demo": null,
        "sales_channels_id": null,
        "demo_thematic": null,
        "reviews_average_rating_cached": null,
        "reviews_count_cached": null,
        "merge_product_id": null,
        "merge_group_product_id": null,
        "merge_group_hidden": null,
        "no_delivery": false,
        "id": 1,
        "main_image_id": null,
        "title": "Cool t-shirt"
      }
    ],
    "notice": "Второй анонс",
    "content": "Привет"
  },
  {
    "id": 1,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "published_at": "2009-03-18T07:53:15.000+03:00",
    "pinned": false,
    "custom_template": null,
    "title": "Статья",
    "author": "Ivan",
    "permalink": "statya",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "tags": [],
    "related_products": [],
    "notice": "Первый анонс",
    "content": "Первая статья"
  }
]
```

### Update article

`PUT /admin/blogs/:blog_id/articles/:article_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `article[title]` |  | article title |
| `article[content]` |  | article content |
| `article[notice]` |  | notice for article |
| `article[author]` |  | article author name |
| `article[permalink]` |  | permalink on blog/* |
| `article[html_title]` |  | html title tag content |
| `article[meta_keywords]` |  | html meta words on page |
| `article[meta_description]` |  | html meta description |
| `article[published_at]` |  | publication date, the article is invisible to users while publication date > current time |
| `article[image_attributes]` |  | image for article |
| `article[all_tags]` |  | array of tags for article, max 10 tags: to remove tags from list - pass the array without this tag, to add tag - pass a full list with a new tag added |
| `article[pinned]` |  | set to true to pin article: it will be always on the top of the blog; only one article can be pinned, so when you pin article, previously pinned article will be unpinned |

**Запрос**

```http
PUT /admin/blogs/1/articles/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "article": {
    "title": "new title",
    "image_attributes": {
      "src": "https://abc.myinsales.ru/images/no_image_compact.jpg"
    },
    "all_tags": [
      "test",
      "test3"
    ],
    "pinned": true
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2026-07-25T01:00:22.000+03:00",
  "published_at": "2017-03-01T07:53:15.000+03:00",
  "pinned": true,
  "custom_template": null,
  "title": "new title",
  "author": "Ivan",
  "permalink": "statya-2",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "image": {
    "id": 33,
    "created_at": "2026-07-25T01:00:22.000+03:00",
    "original_url": "/images/loading.gif"
  },
  "tags": [
    {
      "title": "test"
    },
    {
      "title": "test3"
    }
  ],
  "related_products": [
    {
      "category_id": 1,
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "updated_at": "2015-07-27T17:05:35.000+03:00",
      "is_hidden": false,
      "account_id": 1,
      "title_translations": [
        "Cool t-shirt"
      ],
      "available": false,
      "archived": false,
      "key": null,
      "price_min_cached": "10.0",
      "price_max_cached": "10.0",
      "price_min_available_cached": null,
      "price_max_available_cached": null,
      "canonical_url_collection_id": 1,
      "collections_id": [
        1,
        3,
        8
      ],
      "characteristics_id": [
        7,
        9,
        12,
        13
      ],
      "custom_layout": null,
      "custom_template": null,
      "old_price_min_cached": null,
      "old_price_max_cached": null,
      "old_price_min_available_cached": null,
      "old_price_max_available_cached": null,
      "permalink": "футболка",
      "unit": "pce",
      "deleted": false,
      "cost_price_min_cached": null,
      "cost_price_max_cached": null,
      "stock_currency_id": 1,
      "sort_weight": null,
      "bundle": null,
      "ignore_discounts": null,
      "vat": -1,
      "dimensions": "10x20x30",
      "fiscal_product_type": 1,
      "popularity": null,
      "available_at_warehouses": [],
      "auto_assign_bundle_price": null,
      "demo": null,
      "sales_channels_id": null,
      "demo_thematic": null,
      "reviews_average_rating_cached": null,
      "reviews_count_cached": null,
      "merge_product_id": null,
      "merge_group_product_id": null,
      "merge_group_hidden": null,
      "no_delivery": false,
      "id": 1,
      "main_image_id": null,
      "title": "Cool t-shirt"
    }
  ],
  "notice": "Второй анонс",
  "content": "Привет"
}
```

---

## Blog

| Операция | Метод и путь |
|----------|--------------|
| [Create blog](#create-blog) | `POST /admin/blogs.json` |
| [Destroy blog](#destroy-blog) | `DELETE /admin/blogs/:blog_id.json` |
| [Get blog](#get-blog) | `GET /admin/blogs/:blog_id.json` |
| [Get blogs](#get-blogs) | `GET /admin/blogs.json` |
| [Update blog](#update-blog) | `PUT /admin/blogs/:blog_id.json` |

### Create blog

`POST /admin/blogs.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `blog[title]` | да | title |
| `blog[html_title]` |  | html title tag content |
| `blog[meta_keywords]` |  | html meta words on page |
| `blog[meta_description]` |  | html meta description |
| `blog[notify_about_comment]` |  | true by default |
| `blog[captcha_enabled]` |  | false by default |
| `blog[disallow_links_in_comments]` |  | false by default |
| `blog[commentable]` |  | default: 'no', available: no, yes, moderable |

**Запрос**

```http
POST /admin/blogs.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "blog": {
    "title": "SuperBlog",
    "html_title": "HTML TAG TITLE",
    "meta_keywords": "new cool products",
    "meta_description": "It's a blog about our new products"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 19,
  "commentable": "no",
  "notify_about_comment": true,
  "captcha_enabled": false,
  "disallow_links_in_comments": false,
  "custom_template": null,
  "handle": "superblog",
  "title": "SuperBlog",
  "html_title": "HTML TAG TITLE",
  "meta_keywords": "new cool products",
  "meta_description": "It's a blog about our new products",
  "tags": []
}
```

### Destroy blog

`DELETE /admin/blogs/:blog_id.json`

**Запрос**

```http
DELETE /admin/blogs/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get blog

`GET /admin/blogs/:blog_id.json`

**Запрос**

```http
GET /admin/blogs/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "commentable": "no",
  "notify_about_comment": true,
  "captcha_enabled": false,
  "disallow_links_in_comments": false,
  "custom_template": null,
  "handle": null,
  "title": "Блог",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "tags": [
    {
      "title": "test"
    },
    {
      "title": "test2"
    }
  ]
}
```

### Get blogs

`GET /admin/blogs.json`

**Запрос**

```http
GET /admin/blogs.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "commentable": "no",
    "notify_about_comment": true,
    "captcha_enabled": false,
    "disallow_links_in_comments": false,
    "custom_template": null,
    "handle": null,
    "title": "Блог",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "tags": [
      {
        "title": "test"
      },
      {
        "title": "test2"
      }
    ]
  },
  {
    "id": 2,
    "commentable": "no",
    "notify_about_comment": true,
    "captcha_enabled": false,
    "disallow_links_in_comments": false,
    "custom_template": null,
    "handle": "second_blog",
    "title": "Second blog",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "tags": []
  }
]
```

### Update blog

`PUT /admin/blogs/:blog_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `blog[title]` |  | title |
| `blog[html_title]` |  | html title tag content |
| `blog[meta_keywords]` |  | html meta words on page |
| `blog[meta_description]` |  | html meta description |
| `blog[notify_about_comment]` |  | true by default |
| `blog[captcha_enabled]` |  | false by default |
| `blog[disallow_links_in_comments]` |  | false by default |
| `blog[commentable]` |  | default: 'no', available: no, yes, moderable |

**Запрос**

```http
PUT /admin/blogs/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "blog": {
    "title": "WowBlog",
    "notify_about_comment": false,
    "captcha_enabled": true,
    "disallow_links_in_comments": true,
    "commentable": "yes"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "commentable": "yes",
  "notify_about_comment": false,
  "captcha_enabled": true,
  "disallow_links_in_comments": true,
  "custom_template": null,
  "handle": "wowblog",
  "title": "WowBlog",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "tags": [
    {
      "title": "test"
    },
    {
      "title": "test2"
    }
  ]
}
```

---

## BonusSystemTransaction

| Операция | Метод и путь |
|----------|--------------|
| [Create transaction for client](#create-transaction-for-client) | `POST /admin/clients/:client_id/bonus_system_transactions.json` |
| [Get bonus_system_transactions](#get-bonus-system-transactions) | `GET /admin/bonus_system_transactions.json` |
| [Get transaction for client](#get-transaction-for-client) | `GET /admin/bonus_system_transactions/:bst_id.json` |
| [Get transactions for client](#get-transactions-for-client) | `GET /admin/clients/:client_id/bonus_system_transactions.json` |

### Create transaction for client

`POST /admin/clients/:client_id/bonus_system_transactions.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `bonus_system_transaction[bonus_points]` | да | bonus points; pass negative quantity if you need to revert their addition or just remove them in some cases |
| `bonus_system_transaction[description]` |  | description |

**Запрос**

```http
POST /admin/clients/2/bonus_system_transactions.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "bonus_system_transaction": {
    "bonus_points": 100,
    "description": "compliment"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 145,
  "bonus_points": 100,
  "client_id": 2,
  "bonus_system_id": 1,
  "initiator_type": "ApplicationPassword",
  "initiator_id": 1,
  "description": "compliment",
  "status": "charged",
  "created_at": "2026-07-25T00:59:01.104+03:00",
  "updated_at": "2026-07-25T00:59:01.110+03:00"
}
```

### Get bonus_system_transactions

`GET /admin/bonus_system_transactions.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` |  | Set datetime to get only data updated after it. |
| `from_id` |  | Set id to get only data starting from it. |
| `per_page` |  | Quantity of objects per page. |

**Запрос**

```http
GET /admin/bonus_system_transactions.json?updated_since=2009-01-04+14%3A44%3A43+%2B0300
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
updated_since=2009-01-04 14:44:43 +0300
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "bonus_points": 100,
    "client_id": 7,
    "bonus_system_id": 1,
    "initiator_type": "User",
    "initiator_id": 1,
    "description": "Бонус за регистрацию",
    "status": "charged",
    "created_at": "2011-03-18T10:53:14.000+03:00",
    "updated_at": "2011-03-18T10:53:15.000+03:00"
  },
  {
    "id": 25,
    "bonus_points": 100,
    "client_id": 2,
    "bonus_system_id": 1,
    "initiator_type": "Order",
    "initiator_id": 1,
    "description": "test",
    "status": "charged",
    "created_at": "2026-07-25T00:59:00.217+03:00",
    "updated_at": "2026-07-25T00:59:00.217+03:00"
  },
  {
    "id": 33,
    "bonus_points": 100,
    "client_id": 2,
    "bonus_system_id": 1,
    "initiator_type": "User",
    "initiator_id": 1,
    "description": "test",
    "status": "charged",
    "created_at": "2026-07-25T00:59:00.227+03:00",
    "updated_at": "2026-07-25T00:59:00.227+03:00"
  }
]
```

### Get transaction for client

`GET /admin/bonus_system_transactions/:bst_id.json`

**Запрос**

```http
GET /admin/bonus_system_transactions/89.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 89,
  "bonus_points": 100,
  "client_id": 2,
  "bonus_system_id": 1,
  "initiator_type": "Order",
  "initiator_id": 1,
  "description": "test",
  "status": "charged",
  "created_at": "2026-07-25T00:59:00.725+03:00",
  "updated_at": "2026-07-25T00:59:00.725+03:00"
}
```

### Get transactions for client

`GET /admin/clients/:client_id/bonus_system_transactions.json`

**Запрос**

```http
GET /admin/clients/2/bonus_system_transactions.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 57,
    "bonus_points": 100,
    "client_id": 2,
    "bonus_system_id": 1,
    "initiator_type": "Order",
    "initiator_id": 1,
    "description": "test",
    "status": "charged",
    "created_at": "2026-07-25T00:59:00.458+03:00",
    "updated_at": "2026-07-25T00:59:00.458+03:00"
  },
  {
    "id": 65,
    "bonus_points": 100,
    "client_id": 2,
    "bonus_system_id": 1,
    "initiator_type": "User",
    "initiator_id": 1,
    "description": "test",
    "status": "charged",
    "created_at": "2026-07-25T00:59:00.469+03:00",
    "updated_at": "2026-07-25T00:59:00.469+03:00"
  }
]
```

---

## Category

| Операция | Метод и путь |
|----------|--------------|
| [Create category](#create-category) | `POST /admin/categories.json` |
| [Destroy category](#destroy-category) | `DELETE /admin/categories/:category_id.json` |
| [Get categories](#get-categories) | `GET /admin/categories.json` |
| [Get category](#get-category) | `GET /admin/categories/:category_id.json` |
| [Update category](#update-category) | `PUT /admin/categories/:category_id.json` |

### Create category

`POST /admin/categories.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `category[title]` | да | title |
| `category[parent_id]` |  | parent category id |
| `category[position]` |  | position in categories list |

**Запрос**

```http
POST /admin/categories.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "category": {
    "title": "Reebok"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20497914,
  "parent_id": 1,
  "title": "Reebok",
  "position": 1,
  "created_at": "2026-07-25T00:59:35.000+03:00",
  "updated_at": "2026-07-25T00:59:35.000+03:00"
}
```

### Destroy category

`DELETE /admin/categories/:category_id.json`

**Запрос**

```http
DELETE /admin/categories/6.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get categories

`GET /admin/categories.json`

**Запрос**

```http
GET /admin/categories.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 6,
    "parent_id": 1,
    "title": "T-shirts",
    "position": 0,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00"
  },
  {
    "id": 2,
    "parent_id": 1,
    "title": "Trousers",
    "position": 1,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00"
  },
  {
    "id": 3,
    "parent_id": 1,
    "title": "Shirts",
    "position": 2,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00"
  },
  {
    "id": 4,
    "parent_id": 1,
    "title": "Socks",
    "position": 3,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00"
  },
  {
    "id": 5,
    "parent_id": 1,
    "title": "Hats",
    "position": 4,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00"
  },
  {
    "id": 8,
    "parent_id": 5,
    "title": "Top Hat",
    "position": 1,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00"
  },
  {
    "id": 1,
    "parent_id": null,
    "title": "Склад",
    "position": 0,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00"
  }
]
```

### Get category

`GET /admin/categories/:category_id.json`

**Запрос**

```http
GET /admin/categories/6.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 6,
  "parent_id": 1,
  "title": "T-shirts",
  "position": 0,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2026-07-25T00:58:09.000+03:00"
}
```

### Update category

`PUT /admin/categories/:category_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `category[title]` |  | title |

**Запрос**

```http
PUT /admin/categories/6.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "category": {
    "title": "Adidas"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 6,
  "parent_id": 1,
  "title": "Adidas",
  "position": 0,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2026-07-25T00:59:35.000+03:00"
}
```

---

## Characteristic

| Операция | Метод и путь |
|----------|--------------|
| [Create characteristic](#create-characteristic) | `POST /admin/properties/:property_id/characteristics.json` |
| [Destroy characteristic](#destroy-characteristic) | `DELETE /admin/properties/:property_id/characteristics/:characteristic_id.json` |
| [Get characteristic](#get-characteristic) | `GET /admin/properties/:property_id/characteristics/:characteristic_id.json` |
| [Get characteristics](#get-characteristics) | `GET /admin/properties/:property_id/characteristics.json` |
| [Update characteristic](#update-characteristic) | `PUT /admin/properties/:property_id/characteristics/:characteristic_id.json` |

### Create characteristic

`POST /admin/properties/:property_id/characteristics.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `characteristic[title]` | да | title |
| `characteristic[permalink]` |  | permalink |
| `characteristic[position]` |  | position in characteristics list |

**Запрос**

```http
POST /admin/properties/3/characteristics.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "characteristic": {
    "title": "Reebok"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 898582855,
  "property_id": 3,
  "position": 3,
  "title": "Reebok",
  "permalink": "reebok"
}
```

### Destroy characteristic

`DELETE /admin/properties/:property_id/characteristics/:characteristic_id.json`

**Запрос**

```http
DELETE /admin/properties/3/characteristics/7.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get characteristic

`GET /admin/properties/:property_id/characteristics/:characteristic_id.json`

**Запрос**

```http
GET /admin/properties/3/characteristics/7.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 7,
  "property_id": 3,
  "position": 1,
  "title": "Nike",
  "permalink": "nike"
}
```

### Get characteristics

`GET /admin/properties/:property_id/characteristics.json`

**Запрос**

```http
GET /admin/properties/3/characteristics.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 7,
    "property_id": 3,
    "position": 1,
    "title": "Nike",
    "permalink": "nike"
  },
  {
    "id": 8,
    "property_id": 3,
    "position": 2,
    "title": "Puma",
    "permalink": "puma"
  }
]
```

### Update characteristic

`PUT /admin/properties/:property_id/characteristics/:characteristic_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `characteristic[title]` |  | title |

**Запрос**

```http
PUT /admin/properties/3/characteristics/7.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "characteristic": {
    "title": "Adidas"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 7,
  "property_id": 3,
  "position": 1,
  "title": "Adidas",
  "permalink": "nike"
}
```

---

## Client

| Операция | Метод и путь |
|----------|--------------|
| [Create individual client](#create-individual-client) | `POST /admin/clients.json` |
| [Create juridical client](#create-juridical-client) | `POST /admin/clients.json` |
| [Destroy client](#destroy-client) | `DELETE /admin/clients/:client_id.json` |
| [Get client](#get-client) | `GET /admin/clients/:client_id.json` |
| [Get clients](#get-clients) | `GET /admin/clients.json` |
| [Update client](#update-client) | `PUT /admin/clients/:client_id.json` |

### Create individual client

`POST /admin/clients.json`

To assign values to fields, you should create them in advance. For more information, please go to section about Field. Individual client is a natural person, as opposed to legal entity or a sole trader.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `client[name]` | да | name |
| `client[surname]` | да | surname |
| `client[middlename]` | да | middlename |
| `client[registered]` |  | if true, email and password are required |
| `client[email]` |  | email |
| `client[password]` |  | password at least 6 characters |
| `client[phone]` |  | phone number |
| `client[subscribe]` |  | set true to subscribe client for emails |
| `client[client_group_id]` |  | client group id |
| `client[type]` |  | Client::Individual, by default |
| `client[ip_addr]` |  | IP address |
| `client[correspondent_account]` |  | Correspondent account |
| `client[settlement_account]` |  | Settlement account |
| `client[fields_values_attributes]` |  | fields values array |
| `client[fields_values_attributes][value]` |  | value of field |
| `client[fields_values_attributes][field_id]` |  | field id |
| `client[skip_send_mail]` |  | set it true to skip sending email confirmation |

**Запрос**

```http
POST /admin/clients.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "client": {
    "name": "Ivan",
    "surname": "Ivanov",
    "middlename": "Ivanovich",
    "registered": true,
    "email": "test@example.com",
    "password": "123456",
    "phone": "+71111111111",
    "type": "Client::Individual",
    "fields_values_attributes": [
      {
        "field_id": 256,
        "value": "text"
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 33,
  "email": "test@example.com",
  "name": "Ivan",
  "phone": "+71111111111",
  "created_at": "2026-07-25T00:59:03.000+03:00",
  "updated_at": "2026-07-25T00:59:03.000+03:00",
  "comment": null,
  "registered": true,
  "subscribe": true,
  "client_group_id": 1,
  "bonus_points": 0,
  "type": "Client::Individual",
  "correspondent_account": null,
  "settlement_account": null,
  "consent_to_personal_data": true,
  "o_auth_provider": null,
  "messenger_subscription": null,
  "external_id": null,
  "external_service": null,
  "contact_name": "Ivan",
  "surname": "Ivanov",
  "middlename": "Ivanovich",
  "progressive_discount": "25.0",
  "group_discount": "25.0",
  "ip_addr": "",
  "fields_values": [
    {
      "id": 67,
      "value": "text",
      "created_at": "2026-07-25T00:59:03.360+03:00",
      "updated_at": "2026-07-25T00:59:03.360+03:00",
      "field_id": 256,
      "type": "Текст",
      "name": "special",
      "handle": null
    }
  ],
  "default_address": null,
  "client_tags": [],
  "full_name": "Ivanov Ivan Ivanovich"
}
```

### Create juridical client

`POST /admin/clients.json`

To assign values to fields, you should create them in advance. For more information, please go to section about Field. Juridical client is a sole trader (individual entrepreneur) or legal entity.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `client[name]` | да | name of company or a sole trader |
| `client[registered]` |  | if true, email and password are required |
| `client[email]` |  | email |
| `client[password]` |  | password at least 6 characters |
| `client[phone]` |  | phone number |
| `client[subscribe]` |  | set true to subscribe client for emails |
| `client[client_group_id]` |  | client group id |
| `client[type]` | да | Client::Juridical |
| `client[ip_addr]` |  | IP address |
| `client[juridical_address]` |  | address of company or a sole trader |
| `client[inn]` | да | INN |
| `client[kpp]` |  | KPP |
| `client[ogrn]` |  | OGRN |
| `client[okpo]` |  | OKPO |
| `client[bik]` |  | BIK |
| `client[bank_name]` |  | bank name |
| `client[correspondent_account]` |  | Correspondent account |
| `client[settlement_account]` |  | Settlement account |
| `client[fields_values_attributes]` |  | fields values array |
| `client[fields_values_attributes][field_id]` |  | field id |
| `client[fields_values_attributes][value]` |  | value of field |
| `client[skip_send_mail]` |  | set it true to skip sending email confirmation |

**Запрос**

```http
POST /admin/clients.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "client": {
    "name": "OOO 'Company'",
    "type": "Client::Juridical",
    "juridical_address": "123456, Moscow, Novaya st., 1, 222",
    "inn": 1234567890,
    "kpp": 123456789,
    "bik": 123456789,
    "fields_values_attributes": [
      {
        "field_id": 288,
        "value": "text_value"
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 49,
  "email": null,
  "name": "OOO 'Company'",
  "phone": null,
  "created_at": "2026-07-25T00:59:03.000+03:00",
  "updated_at": "2026-07-25T00:59:03.000+03:00",
  "comment": null,
  "registered": false,
  "subscribe": true,
  "client_group_id": null,
  "bonus_points": 0,
  "type": "Client::Juridical",
  "correspondent_account": null,
  "settlement_account": null,
  "consent_to_personal_data": true,
  "o_auth_provider": null,
  "messenger_subscription": null,
  "external_id": null,
  "external_service": null,
  "contact_name": null,
  "juridical_address": "123456, Moscow, Novaya st., 1, 222",
  "inn": "1234567890",
  "kpp": "123456789",
  "ogrn": null,
  "okpo": null,
  "bik": "123456789",
  "bank_name": null,
  "progressive_discount": null,
  "group_discount": null,
  "ip_addr": "",
  "fields_values": [
    {
      "id": 99,
      "value": "text_value",
      "created_at": "2026-07-25T00:59:03.774+03:00",
      "updated_at": "2026-07-25T00:59:03.774+03:00",
      "field_id": 288,
      "type": "Текст",
      "name": "field_juridical",
      "handle": null
    }
  ],
  "default_address": null,
  "client_tags": [],
  "full_name": "OOO 'Company'"
}
```

### Destroy client

`DELETE /admin/clients/:client_id.json`

**Запрос**

```http
DELETE /admin/clients/5.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
"ok"
```

### Get client

`GET /admin/clients/:client_id.json`

**Запрос**

```http
GET /admin/clients/5.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 5,
  "email": "john@smith.com",
  "name": "John Smith",
  "phone": "+7(926)111-11-11",
  "created_at": "2009-04-30T04:37:32.000+04:00",
  "updated_at": "2026-07-25T00:59:02.000+03:00",
  "comment": null,
  "registered": true,
  "subscribe": true,
  "client_group_id": null,
  "bonus_points": 1000,
  "type": "Client::Individual",
  "correspondent_account": null,
  "settlement_account": null,
  "consent_to_personal_data": null,
  "o_auth_provider": null,
  "messenger_subscription": null,
  "external_id": null,
  "external_service": null,
  "contact_name": "John Smith",
  "surname": null,
  "middlename": null,
  "progressive_discount": null,
  "group_discount": null,
  "ip_addr": "",
  "fields_values": [
    {
      "id": 35,
      "value": "old value",
      "created_at": "2026-07-25T00:59:02.847+03:00",
      "updated_at": "2026-07-25T00:59:02.847+03:00",
      "field_id": 240,
      "type": "Текст",
      "name": "special",
      "handle": null
    }
  ],
  "default_address": null,
  "client_tags": [],
  "full_name": "John Smith"
}
```

### Get clients

`GET /admin/clients.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` |  | Set datetime to get only data updated after it. |
| `from_id` |  | Set id to get only data starting from it. |
| `per_page` |  | Quantity of objects per page. |
| `q` |  | Search queue. Matches name, surname, middlename, phone, or email. |
| `phone` |  | Exact match search by phone number. |
| `email` |  | Exact match search by email address. |
| `registered` |  | Filter by registration status. Accepted values: 'yes' or 'no'. |
| `subscribe` |  | Filter by subscription permission. Accepted values: 'yes' or 'no'. |
| `rfm_group` |  | RFM group value or array. Allowed values: 'small_bills', 'not_frequent', 'not_recent', 'best', 'worst'. |

**Запрос**

```http
GET /admin/clients.json?updated_since=2009-01-04+14%3A44%3A43+%2B0300
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
updated_since=2009-01-04 14:44:43 +0300
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2009-03-31T12:20:34.000+04:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [],
    "default_address": null,
    "client_tags": [],
    "full_name": "Mortimer Hant"
  },
  {
    "id": 3,
    "email": "some-one2@yandex.ru",
    "name": "Just Man",
    "phone": "+7(926)111-22-33",
    "created_at": "2009-04-30T04:37:32.000+04:00",
    "updated_at": "2009-04-30T12:20:34.000+04:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": 1,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Just Man",
    "surname": null,
    "middlename": null,
    "progressive_discount": "25.0",
    "group_discount": "25.0",
    "ip_addr": "",
    "fields_values": [],
    "default_address": null,
    "client_tags": [],
    "full_name": "Just Man"
  },
  {
    "id": 5,
    "email": "john@smith.com",
    "name": "John Smith",
    "phone": "+7(926)111-11-11",
    "created_at": "2009-04-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:02.000+03:00",
    "comment": null,
    "registered": true,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 1000,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "John Smith",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 19,
        "value": "old value",
        "created_at": "2026-07-25T00:59:02.544+03:00",
        "updated_at": "2026-07-25T00:59:02.544+03:00",
        "field_id": 224,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "default_address": null,
    "client_tags": [],
    "full_name": "John Smith"
  }
]
```

### Update client

`PUT /admin/clients/:client_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `client[name]` |  | name |
| `client[phone]` |  | phone number |
| `client[subscribe]` |  | set true to subscribe client for emails |
| `client[client_group_id]` |  | client group id |
| `client[ip_addr]` |  | IP address |
| `client[juridical_address]` |  | address of company or a sole trader |
| `client[inn]` |  | INN |
| `client[kpp]` |  | KPP |
| `client[ogrn]` |  | OGRN |
| `client[okpo]` |  | OKPO |
| `client[bik]` |  | BIK |
| `client[bank_name]` |  | bank name |
| `client[correspondent_account]` |  | Correspondent account |
| `client[settlement_account]` |  | Settlement account |
| `client[fields_values_attributes]` |  | fields values array |

**Запрос**

```http
PUT /admin/clients/5.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "client": {
    "phone": "+7(911)111-11-11",
    "subscribe": false,
    "client_group_id": 2,
    "fields_values_attributes": [
      {
        "field_id": 304,
        "value": "new value"
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 5,
  "email": "john@smith.com",
  "name": "John Smith",
  "phone": "+79111111111",
  "created_at": "2009-04-30T04:37:32.000+04:00",
  "updated_at": "2026-07-25T00:59:04.000+03:00",
  "comment": null,
  "registered": true,
  "subscribe": false,
  "client_group_id": 2,
  "bonus_points": 1000,
  "type": "Client::Individual",
  "correspondent_account": null,
  "settlement_account": null,
  "consent_to_personal_data": null,
  "o_auth_provider": null,
  "messenger_subscription": null,
  "external_id": null,
  "external_service": null,
  "contact_name": "John Smith",
  "surname": null,
  "middlename": null,
  "progressive_discount": null,
  "group_discount": null,
  "ip_addr": "",
  "fields_values": [
    {
      "id": 115,
      "value": "new value",
      "created_at": "2026-07-25T00:59:04.074+03:00",
      "updated_at": "2026-07-25T00:59:04.228+03:00",
      "field_id": 304,
      "type": "Текст",
      "name": "special",
      "handle": null
    }
  ],
  "default_address": null,
  "client_tags": [],
  "full_name": "John Smith"
}
```

---

## ClientGroup

| Операция | Метод и путь |
|----------|--------------|
| [Create client group](#create-client-group) | `POST /admin/client_groups.json` |
| [Destroy client group](#destroy-client-group) | `DELETE /admin/client_groups/:client_group_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/client_groups/:client_group_id.json` |
| [Get client groups](#get-client-groups) | `GET /admin/client_groups.json` |
| [Update client group](#update-client-group) | `PUT /admin/client_groups/:client_group_id.json` |

### Create client group

`POST /admin/client_groups.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `client_group[title]` |  | client group title |
| `client_group[discount]` |  | discount |
| `client_group[discount_description]` |  | discount description |
| `client_group[is_default]` |  | set as default group |

**Запрос**

```http
POST /admin/client_groups.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "client_group": {
    "title": "test group",
    "discount": 0.2,
    "discount_description": "test group discount",
    "is_default": false
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 33,
  "title": "test group",
  "discount": "0.2",
  "discount_description": "test group discount",
  "is_default": false,
  "b2b": false,
  "created_at": "2026-07-25T00:59:01.765+03:00",
  "updated_at": "2026-07-25T00:59:01.769+03:00"
}
```

### Destroy client group

`DELETE /admin/client_groups/:client_group_id.json`

**Запрос**

```http
DELETE /admin/client_groups/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/client_groups/:client_group_id.json`

Content-Type: application/json Get client group

**Запрос**

```http
GET /admin/client_groups/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "title": "vip",
  "discount": "25.0",
  "discount_description": null,
  "is_default": true,
  "b2b": false,
  "created_at": "2026-04-16T00:58:08.000+03:00",
  "updated_at": "2026-07-25T00:58:09.121+03:00"
}
```

### Get client groups

`GET /admin/client_groups.json`

**Запрос**

```http
GET /admin/client_groups.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "title": "vip",
    "discount": "25.0",
    "discount_description": null,
    "is_default": true,
    "b2b": false,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.121+03:00"
  }
]
```

### Update client group

`PUT /admin/client_groups/:client_group_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `client_group[title]` |  | client group title |
| `client_group[discount]` |  | discount |
| `client_group[discount_description]` |  | discount description |
| `client_group[is_default]` |  | set as default group |

**Запрос**

```http
PUT /admin/client_groups/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "client_group": {
    "title": "test group",
    "discount": 0.2,
    "discount_description": "test group discount",
    "is_default": false
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

---

## Collect

| Операция | Метод и путь |
|----------|--------------|
| [Add product to collection](#add-product-to-collection) | `POST /admin/collects.json` |
| [Get collections by product](#get-collections-by-product) | `GET /admin/collects.json?product_id=:product_id` |
| [Get collects](#get-collects) | `GET /admin/collects.json` |
| [Get products by collection](#get-products-by-collection) | `GET /admin/collects.json?collection_id=:collection_id` |
| [Move product to another collection](#move-product-to-another-collection) | `PUT /admin/collects/:collect_id.json` |
| [Remove product from collection](#remove-product-from-collection) | `DELETE /admin/collects/:collect_id.json` |
| [Update position](#update-position) | `PUT /admin/collects/:collect_id.json` |

### Add product to collection

`POST /admin/collects.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collect[collection_id]` | да | collection id |
| `collect[product_id]` | да | product id |

**Запрос**

```http
POST /admin/collects.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collect": {
    "collection_id": 2,
    "product_id": 1
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 9854,
  "collection_id": 2,
  "product_id": 1,
  "position": 1,
  "created_at": "2026-07-25T00:59:40.000+03:00"
}
```

### Get collections by product

`GET /admin/collects.json?product_id=:product_id`

**Запрос**

```http
GET /admin/collects.json?product_id=1
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
product_id=1
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "collection_id": 1,
    "product_id": 1,
    "position": 1,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 3,
    "collection_id": 3,
    "product_id": 1,
    "position": 3,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 4,
    "collection_id": 8,
    "product_id": 1,
    "position": 1,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  }
]
```

### Get collects

`GET /admin/collects.json`

**Запрос**

```http
GET /admin/collects.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "collection_id": 1,
    "product_id": 1,
    "position": 1,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 2,
    "collection_id": 1,
    "product_id": 2,
    "position": 2,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 3,
    "collection_id": 3,
    "product_id": 1,
    "position": 3,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 4,
    "collection_id": 8,
    "product_id": 1,
    "position": 1,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 5,
    "collection_id": 8,
    "product_id": 2,
    "position": 2,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 6,
    "collection_id": 8,
    "product_id": 3,
    "position": 3,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  }
]
```

### Get products by collection

`GET /admin/collects.json?collection_id=:collection_id`

**Запрос**

```http
GET /admin/collects.json?collection_id=1
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
collection_id=1
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "collection_id": 1,
    "product_id": 1,
    "position": 1,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  },
  {
    "id": 2,
    "collection_id": 1,
    "product_id": 2,
    "position": 2,
    "created_at": "2009-03-18T07:53:15.000+03:00"
  }
]
```

### Move product to another collection

`PUT /admin/collects/:collect_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collect[collection_id]` |  | collection id |

**Запрос**

```http
PUT /admin/collects/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collect": {
    "collection_id": 3
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "collection_id": 3,
  "product_id": 2,
  "position": 2,
  "created_at": "2009-03-18T07:53:15.000+03:00"
}
```

### Remove product from collection

`DELETE /admin/collects/:collect_id.json`

**Запрос**

```http
DELETE /admin/collects/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Update position

`PUT /admin/collects/:collect_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collect[position]` |  | position |

**Запрос**

```http
PUT /admin/collects/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collect": {
    "position": 3
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "collection_id": 1,
  "product_id": 2,
  "position": 3,
  "created_at": "2009-03-18T07:53:15.000+03:00"
}
```

---

## Collection

| Операция | Метод и путь |
|----------|--------------|
| [Create collection](#create-collection) | `POST /admin/collections.json` |
| [Destroy collection](#destroy-collection) | `DELETE /admin/collections/:collection_id.json` |
| [Get collection](#get-collection) | `GET /admin/collections/:collection_id.json` |
| [Get collections](#get-collections) | `GET /admin/collections.json` |
| [Update collection](#update-collection) | `PUT /admin/collections/:collection_id.json` |

### Create collection

`POST /admin/collections.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection[title]` | да | title |
| `collection[parent_id]` | да | parent collection id |
| `collection[position]` |  | position in collections list |
| `collection[sort_type]` |  | Alphabetically: A - Z => 1, Alphabetically: Z - A => 2, By date: New - Old => 3, By date: Old - New => 4, By price: Expensive - Cheap => 5, By price: Cheap - Expensive => 6, Manual => 7, By discount: Large - Little => 8, By discount: Little - Large => 9, By sort rank: Low - High => 10, By sort rank: High - Low => 11, Sort by ascending popularity => 12, Sort by descending popularity => 13 |
| `collection[description]` |  | description |
| `collection[html_title]` |  | html title |
| `collection[seo_description]` |  | seo description |
| `collection[meta_description]` |  | meta description |
| `collection[meta_keywords]` |  | meta keywords |
| `collection[is_hidden]` |  | set true to hide collection |
| `collection[permalink]` |  | generated |
| `collection[recursive]` |  | recursive |
| `collection[image_attributes]` |  | image for collection |
| `collection[image_attributes][src]` |  | image url |
| `collection[field_values_attributes]` |  | collection field values attributes |

**Запрос**

```http
POST /admin/collections.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collection": {
    "title": "T-shirts",
    "parent_id": 1,
    "position": 1,
    "sort_type": 1,
    "description": "text",
    "image_attributes": {
      "src": "https://abc.myinsales.ru/images/no_image_compact.jpg"
    },
    "field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20924767,
  "parent_id": 1,
  "is_hidden": false,
  "is_smart": null,
  "position": 1,
  "created_at": "2026-07-25T00:59:34.000+03:00",
  "updated_at": "2026-07-25T00:59:34.000+03:00",
  "recursive": null,
  "sort_type": 1,
  "custom_template": null,
  "title": "T-shirts",
  "html_title": null,
  "meta_description": null,
  "meta_keywords": null,
  "permalink": "t-shirts",
  "url": "/collection/t-shirts",
  "image": {
    "id": 9,
    "created_at": "2026-07-25T00:59:34.000+03:00",
    "original_url": "/images/loading.gif"
  },
  "field_values": [
    {
      "id": 18,
      "collection_field_id": 1,
      "value": "12345"
    }
  ],
  "description": "text",
  "seo_description": null
}
```

### Destroy collection

`DELETE /admin/collections/:collection_id.json`

**Запрос**

```http
DELETE /admin/collections/6.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get collection

`GET /admin/collections/:collection_id.json`

**Запрос**

```http
GET /admin/collections/6.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 6,
  "parent_id": 1,
  "is_hidden": false,
  "is_smart": null,
  "position": 0,
  "created_at": "2026-04-16T00:58:08.000+03:00",
  "updated_at": "2026-07-25T00:58:09.000+03:00",
  "recursive": null,
  "sort_type": 7,
  "custom_template": null,
  "title": "all",
  "html_title": null,
  "meta_description": null,
  "meta_keywords": null,
  "permalink": "all",
  "url": "/collection/all",
  "field_values": [],
  "description": null,
  "seo_description": null
}
```

### Get collections

`GET /admin/collections.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` |  | Set datetime to get only data updated after it. |
| `from_id` |  | Set id to get only data starting from it. |
| `per_page` |  | Quantity of objects per page. |

**Запрос**

```http
GET /admin/collections.json?updated_since=2016-07-04+14%3A44%3A43+%2B0300
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
updated_since=2016-07-04 14:44:43 +0300
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "parent_id": null,
    "is_hidden": false,
    "is_smart": null,
    "position": -1,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "recursive": null,
    "sort_type": 7,
    "custom_template": null,
    "title": "T-shirts",
    "html_title": null,
    "meta_description": null,
    "meta_keywords": null,
    "permalink": "футболка",
    "url": "/collection/футболка",
    "field_values": [
      {
        "id": 1,
        "collection_field_id": 1,
        "value": "some collection field value"
      },
      {
        "id": 2,
        "collection_field_id": 2,
        "value": "some collection field value"
      }
    ],
    "description": null,
    "seo_description": null
  },
  {
    "id": 2,
    "parent_id": 1,
    "is_hidden": false,
    "is_smart": null,
    "position": 1,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "recursive": null,
    "sort_type": 7,
    "custom_template": null,
    "title": "Shirts",
    "html_title": null,
    "meta_description": null,
    "meta_keywords": null,
    "permalink": "shirts",
    "url": "/collection/shirts",
    "field_values": [],
    "description": null,
    "seo_description": null
  },
  {
    "id": 3,
    "parent_id": 1,
    "is_hidden": false,
    "is_smart": null,
    "position": 2,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "recursive": null,
    "sort_type": 7,
    "custom_template": null,
    "title": "Trousers",
    "html_title": null,
    "meta_description": null,
    "meta_keywords": null,
    "permalink": "trousers",
    "url": "/collection/trousers",
    "field_values": [],
    "description": null,
    "seo_description": null
  },
  {
    "id": 4,
    "parent_id": 1,
    "is_hidden": false,
    "is_smart": null,
    "position": 3,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "recursive": null,
    "sort_type": 7,
    "custom_template": null,
    "title": "Stylish",
    "html_title": null,
    "meta_description": null,
    "meta_keywords": null,
    "permalink": "stylish",
    "url": "/collection/stylish",
    "field_values": [],
    "description": null,
    "seo_description": null
  },
  {
    "id": 5,
    "parent_id": 1,
    "is_hidden": false,
    "is_smart": null,
    "position": 4,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "recursive": null,
    "sort_type": 7,
    "custom_template": null,
    "title": "Sport",
    "html_title": null,
    "meta_description": null,
    "meta_keywords": null,
    "permalink": "sport",
    "url": "/collection/sport",
    "field_values": [],
    "description": null,
    "seo_description": null
  },
  {
    "id": 6,
    "parent_id": 1,
    "is_hidden": false,
    "is_smart": null,
    "position": 0,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "recursive": null,
    "sort_type": 7,
    "custom_template": null,
    "title": "all",
    "html_title": null,
    "meta_description": null,
    "meta_keywords": null,
    "permalink": "all",
    "url": "/collection/all",
    "field_values": [],
    "description": null,
    "seo_description": null
  },
  {
    "id": 8,
    "parent_id": 1,
    "is_hidden": false,
    "is_smart": null,
    "position": 4,
    "created_at": "2026-04-16T00:58:08.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "recursive": null,
    "sort_type": 7,
    "custom_template": null,
    "title": "sort test",
    "html_title": null,
    "meta_description": null,
    "meta_keywords": null,
    "permalink": "sort test",
    "url": "/collection/sort test",
    "field_values": [],
    "description": null,
    "seo_description": null
  }
]
```

### Update collection

`PUT /admin/collections/:collection_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection[title]` |  | title |
| `collection[position]` |  | position in collections list |
| `collection[sort_type]` |  | Alphabetically: A - Z => 1, Alphabetically: Z - A => 2, By date: New - Old => 3, By date: Old - New => 4, By price: Expensive - Cheap => 5, By price: Cheap - Expensive => 6, Manual => 7, By discount: Large - Little => 8, By discount: Little - Large => 9, By sort rank: Low - High => 10, By sort rank: High - Low => 11, Sort by ascending popularity => 12, Sort by descending popularity => 13 |
| `collection[description]` |  | description |
| `collection[html_title]` |  | html title |
| `collection[seo_description]` |  | seo description |
| `collection[meta_description]` |  | meta description |
| `collection[meta_keywords]` |  | meta keywords |
| `collection[is_hidden]` |  | set true to hide collection |
| `collection[permalink]` |  | generated |
| `collection[recursive]` |  | recursive |
| `collection[field_values_attributes]` |  | collection field values attributes |

**Запрос**

```http
PUT /admin/collections/6.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collection": {
    "title": "Books",
    "position": 2,
    "sort_type": 11,
    "description": "Description",
    "is_hidden": true,
    "field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 6,
  "parent_id": 1,
  "is_hidden": true,
  "is_smart": null,
  "position": 2,
  "created_at": "2026-04-16T00:58:08.000+03:00",
  "updated_at": "2026-07-25T00:59:34.000+03:00",
  "recursive": null,
  "sort_type": 11,
  "custom_template": null,
  "title": "Books",
  "html_title": null,
  "meta_description": null,
  "meta_keywords": null,
  "permalink": "all",
  "url": "/collection/all",
  "field_values": [
    {
      "id": 34,
      "collection_field_id": 1,
      "value": "12345"
    }
  ],
  "description": "Description",
  "seo_description": null
}
```

---

## CollectionField

| Операция | Метод и путь |
|----------|--------------|
| [Create collection_field](#create-collection-field) | `POST /admin/collection_fields.json` |
| [Destroy collection_field by handle](#destroy-collection-field-by-handle) | `DELETE /admin/collection_fields/:collection_field_handle.json` |
| [Destroy collection_field](#destroy-collection-field) | `DELETE /admin/collection_fields/:collection_field_id.json` |
| [Get collection_field by handle](#get-collection-field-by-handle) | `GET /admin/collection_fields/:collection_field_handle.json` |
| [Get collection_field](#get-collection-field) | `GET /admin/collection_fields/:collection_field_id.json` |
| [Get collection_fields](#get-collection-fields) | `GET /admin/collection_fields.json` |
| [Update collection_field by handle](#update-collection-field-by-handle) | `PUT /admin/collection_fields/:collection_field_handle.json` |
| [Update collection_field](#update-collection-field) | `PUT /admin/collection_fields/:collection_field_id.json` |

### Create collection_field

`POST /admin/collection_fields.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_field[title]` | да | title |
| `collection_field[handle]` | да | handle used in liquid |
| `collection_field[type]` | да | CollectionField::TextArea - textarea, CollectionField::TextField - textfield |
| `collection_field[position]` |  | position in collection_fields list |
| `collection_field[is_hidden]` |  | boolean |

**Запрос**

```http
POST /admin/collection_fields.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collection_field": {
    "title": "Size",
    "handle": "size",
    "type": "CollectionField::TextArea"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 19,
  "handle": "size",
  "type": "CollectionField::TextArea",
  "application_id": 1,
  "position": 3,
  "is_hidden": null,
  "title": "Size",
  "application_title": "Цифровые товары"
}
```

### Destroy collection_field by handle

`DELETE /admin/collection_fields/:collection_field_handle.json`

**Запрос**

```http
DELETE /admin/collection_fields/handle2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Destroy collection_field

`DELETE /admin/collection_fields/:collection_field_id.json`

**Запрос**

```http
DELETE /admin/collection_fields/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get collection_field by handle

`GET /admin/collection_fields/:collection_field_handle.json`

**Запрос**

```http
GET /admin/collection_fields/handle2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "CollectionField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "title": "title2",
  "application_title": "Цифровые товары"
}
```

### Get collection_field

`GET /admin/collection_fields/:collection_field_id.json`

**Запрос**

```http
GET /admin/collection_fields/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "CollectionField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "title": "title2",
  "application_title": "Цифровые товары"
}
```

### Get collection_fields

`GET /admin/collection_fields.json`

**Запрос**

```http
GET /admin/collection_fields.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "handle": "handle",
    "type": "CollectionField::TextField",
    "application_id": null,
    "position": 1,
    "is_hidden": null,
    "title": "title",
    "application_title": null
  },
  {
    "id": 2,
    "handle": "handle2",
    "type": "CollectionField::TextArea",
    "application_id": 1,
    "position": 2,
    "is_hidden": null,
    "title": "title2",
    "application_title": "Цифровые товары"
  }
]
```

### Update collection_field by handle

`PUT /admin/collection_fields/:collection_field_handle.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_field[title]` |  | title |
| `collection_field[handle]` |  | handle used in liquid  new_handle |
| `collection_field[position]` |  | position in collection_fields list |
| `collection_field[is_hidden]` |  | boolean |

**Запрос**

```http
PUT /admin/collection_fields/handle2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collection_field": {
    "title": "Volume",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Update collection_field

`PUT /admin/collection_fields/:collection_field_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_field[title]` |  | title |
| `collection_field[handle]` |  | handle used in liquid |
| `collection_field[position]` |  | position in collection_fields list |
| `collection_field[is_hidden]` |  | boolean |

**Запрос**

```http
PUT /admin/collection_fields/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collection_field": {
    "title": "Volume",
    "handle": "volume",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

---

## CollectionFieldValue

| Операция | Метод и путь |
|----------|--------------|
| [Get collection field value by handle](#get-collection-field-value-by-handle) | `GET /admin/collections/:collection_id/collection_field_values/:field_handle.json` |
| [Get collection field value](#get-collection-field-value) | `GET /admin/collections/:collection_id/collection_field_values/:field_id.json` |
| [Get collection fields values](#get-collection-fields-values) | `GET /admin/collections/:collection_id/collection_field_values.json` |
| [Update collection field value by handle](#update-collection-field-value-by-handle) | `PUT /admin/collections/:collection_id/collection_field_values/:field_handle.json` |
| [Update collection field value](#update-collection-field-value) | `PUT /admin/collections/:collection_id/collection_field_values/:field_id.json` |

### Get collection field value by handle

`GET /admin/collections/:collection_id/collection_field_values/:field_handle.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_id` |  | collection id |
| `field_handle` |  | field handle |

**Запрос**

```http
GET /admin/collections/1/collection_field_values/handle2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "collection_field_id": 2,
  "value": "some collection field value"
}
```

### Get collection field value

`GET /admin/collections/:collection_id/collection_field_values/:field_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_id` |  | collection id |
| `field_id` |  | field_id |

**Запрос**

```http
GET /admin/collections/1/collection_field_values/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "collection_field_id": 2,
  "value": "some collection field value"
}
```

### Get collection fields values

`GET /admin/collections/:collection_id/collection_field_values.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_id` |  | collection id |

**Запрос**

```http
GET /admin/collections/1/collection_field_values.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "collection_field_id": 1,
    "value": "some collection field value"
  },
  {
    "id": 2,
    "collection_field_id": 2,
    "value": "some collection field value"
  }
]
```

### Update collection field value by handle

`PUT /admin/collections/:collection_id/collection_field_values/:field_handle.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_id` |  | collection id |
| `field_handle` |  | field handle |
| `collection_field_value[value]` |  | field value |

**Запрос**

```http
PUT /admin/collections/1/collection_field_values/handle2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collection_field_value": {
    "value": "new test value of field"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Update collection field value

`PUT /admin/collections/:collection_id/collection_field_values/:field_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_id` |  | collection id |
| `field_id` |  | field_id |
| `collection_field_value[value]` |  | field value |

**Запрос**

```http
PUT /admin/collections/1/collection_field_values/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "collection_field_value": {
    "value": "new test value of field"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

---

## CollectionFilter

| Операция | Метод и путь |
|----------|--------------|
| [Create collection filter](#create-collection-filter) | `POST /admin/collection_filters.json` |
| [Delete collection filter](#delete-collection-filter) | `DELETE /admin/collection_filters/:filter_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/collection_filters/:filter_id.json` |
| [Get collection filters](#get-collection-filters) | `GET /admin/collection_filters.json?collection_id=:collection_id` |
| [Update collection filter](#update-collection-filter) | `PUT /admin/collection_filters/:filter_id.json` |

### Create collection filter

`POST /admin/collection_filters.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `filter[title]` | да | filter title |
| `filter[collection_id]` | да | collection id |
| `filter[description]` |  | description |
| `filter[permalink]` |  | filter permalink |
| `filter[html_title]` |  | tag title |
| `filter[meta_keywords]` |  | tag keyword |
| `filter[meta_description]` |  | tag description |
| `filter[properties_attributes]` |  | properties |
| `filter[options_attributes]` |  | options |
| `filter[field_values_attributes]` |  | collection filter field values attributes |

**Запрос**

```http
POST /admin/collection_filters.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "filter": {
    "title": "new filter",
    "collection_id": 1,
    "properties_attributes": [
      {
        "title": "first property",
        "value": "first prop value"
      },
      {
        "title": "second property",
        "value": "second property value"
      }
    ],
    "options_attributes": [
      {
        "title": "first option",
        "value": "first option value"
      },
      {
        "title": "second option",
        "value": "second option value"
      }
    ],
    "field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 18,
  "collection_id": 1,
  "description": null,
  "seo_description": null,
  "title": "new filter",
  "permalink": "new-filter",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "characteristis": [
    {
      "id": 898582799,
      "property_id": 35187297,
      "position": 1,
      "title": "first prop value",
      "permalink": "first-prop-value"
    },
    {
      "id": 898582807,
      "property_id": 35187298,
      "position": 1,
      "title": "second property value",
      "permalink": "second-property-value"
    }
  ],
  "option_values": [
    {
      "id": 21063157,
      "option_name_id": 2517843,
      "position": 1,
      "title": "first option value",
      "image_url": null,
      "permalink": "first-option-value"
    },
    {
      "id": 21063158,
      "option_name_id": 2517844,
      "position": 1,
      "title": "second option value",
      "image_url": null,
      "permalink": "second-option-value"
    }
  ]
}
```

### Delete collection filter

`DELETE /admin/collection_filters/:filter_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `filter[id]` | да | filter id |

**Запрос**

```http
DELETE /admin/collection_filters/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
Content-Type: application/x-www-form-urlencoded
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/collection_filters/:filter_id.json`

Content-Type: application/json Get collection filter

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `filter_id` |  | filter id |

**Запрос**

```http
GET /admin/collection_filters/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "collection_id": 1,
  "description": null,
  "seo_description": null,
  "title": "seo filter",
  "permalink": "foobar",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "characteristis": [
    {
      "id": 1,
      "property_id": 1,
      "position": 1,
      "title": "Japan",
      "permalink": "japan"
    },
    {
      "id": 2,
      "property_id": 2,
      "position": 1,
      "title": "Germany",
      "permalink": "germany"
    }
  ],
  "option_values": []
}
```

### Get collection filters

`GET /admin/collection_filters.json?collection_id=:collection_id`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_id` |  | collection id |

**Запрос**

```http
GET /admin/collection_filters.json?collection_id=1
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
collection_id=1
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "collection_id": 1,
    "description": null,
    "seo_description": null,
    "title": "seo filter",
    "permalink": "foobar",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "characteristis": [
      {
        "id": 1,
        "property_id": 1,
        "position": 1,
        "title": "Japan",
        "permalink": "japan"
      },
      {
        "id": 2,
        "property_id": 2,
        "position": 1,
        "title": "Germany",
        "permalink": "germany"
      }
    ],
    "option_values": []
  },
  {
    "id": 2,
    "collection_id": 1,
    "description": null,
    "seo_description": null,
    "title": "seo filter 2",
    "permalink": "barfoo",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "characteristis": [],
    "option_values": [
      {
        "id": 1,
        "option_name_id": 1,
        "position": 1,
        "title": "Red",
        "image_url": null,
        "permalink": null
      },
      {
        "id": 2,
        "option_name_id": 1,
        "position": 2,
        "title": "Зеленый",
        "image_url": null,
        "permalink": null
      }
    ]
  }
]
```

### Update collection filter

`PUT /admin/collection_filters/:filter_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `filter[id]` | да | filter id |
| `filter[title]` |  | filter title |
| `filter[collection_id]` |  | collection id |
| `filter[description]` |  | description |
| `filter[permalink]` |  | filter permalink |
| `filter[html_title]` |  | tag title |
| `filter[meta_keywords]` |  | tag keyword |
| `filter[meta_description]` |  | tag description |
| `filter[properties_attributes]` |  | properties |
| `filter[options_attributes]` |  | options |
| `filter[field_values_attributes]` |  | collection filter field values attributes |

**Запрос**

```http
PUT /admin/collection_filters/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "filter": {
    "title": "new-seo-filter",
    "permalink": "new-seo-filter",
    "properties_attributes": [
      {
        "title": "first property",
        "value": "first prop value"
      },
      {
        "title": "second property",
        "value": "second property value"
      }
    ],
    "options_attributes": [
      {
        "title": "first option",
        "value": "first option value"
      },
      {
        "title": "second option",
        "value": "second option value"
      }
    ],
    "field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "collection_id": 1,
  "description": null,
  "seo_description": null,
  "title": "new-seo-filter",
  "permalink": "new-seo-filter",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "characteristis": [
    {
      "id": 898582831,
      "property_id": 35187301,
      "position": 1,
      "title": "first prop value",
      "permalink": "first-prop-value"
    },
    {
      "id": 898582839,
      "property_id": 35187302,
      "position": 1,
      "title": "second property value",
      "permalink": "second-property-value"
    }
  ],
  "option_values": [
    {
      "id": 21063161,
      "option_name_id": 2517847,
      "position": 1,
      "title": "first option value",
      "image_url": null,
      "permalink": "first-option-value"
    },
    {
      "id": 21063162,
      "option_name_id": 2517848,
      "position": 1,
      "title": "second option value",
      "image_url": null,
      "permalink": "second-option-value"
    }
  ]
}
```

---

## CustomStatus

| Операция | Метод и путь |
|----------|--------------|
| [Create custom status](#create-custom-status) | `POST /admin/custom_statuses.json` |
| [Destroy custom status](#destroy-custom-status) | `DELETE /admin/custom_statuses/:custom_status_permalink.json` |
| [Get custom status](#get-custom-status) | `GET /admin/custom_statuses/:custom_status_permalink.json` |
| [Get custom statuses](#get-custom-statuses) | `GET /admin/custom_statuses.json` |
| [Update custom status](#update-custom-status) | `PUT /admin/custom_statuses/:custom_status_permalink.json` |

### Create custom status

`POST /admin/custom_statuses.json`

For more information about custom statuses More info

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `custom_status[system_status]` | да | system status |
| `custom_status[title]` | да | description |

**Запрос**

```http
POST /admin/custom_statuses.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "custom_status": {
    "system_status": "new",
    "title": "New (ordered by phone)"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 41,
  "system_status": "new",
  "permalink": "new-ordered-by-phone",
  "position": 2,
  "is_default": false,
  "created_at": "2026-07-25T00:59:05.391+03:00",
  "updated_at": "2026-07-25T00:59:05.391+03:00",
  "title": "New (ordered by phone)",
  "color": "#DC5F59"
}
```

### Destroy custom status

`DELETE /admin/custom_statuses/:custom_status_permalink.json`

For more information about custom statuses More info

**Запрос**

```http
DELETE /admin/custom_statuses/novyy.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get custom status

`GET /admin/custom_statuses/:custom_status_permalink.json`

For more information about custom statuses More info

**Запрос**

```http
GET /admin/custom_statuses/novyy.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 9,
  "system_status": "new",
  "permalink": "novyy",
  "position": 1,
  "is_default": true,
  "created_at": "2026-07-25T00:59:04.983+03:00",
  "updated_at": "2026-07-25T00:59:04.983+03:00",
  "title": "Новый",
  "color": "#DC5F59"
}
```

### Get custom statuses

`GET /admin/custom_statuses.json`

For more information about custom statuses More info

**Запрос**

```http
GET /admin/custom_statuses.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": null,
    "system_status": "new",
    "permalink": "new",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "Новый",
    "color": "#DC5F59"
  },
  {
    "id": null,
    "system_status": "accepted",
    "permalink": "accepted",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "В обработке",
    "color": "#85C5E0"
  },
  {
    "id": null,
    "system_status": "approved",
    "permalink": "approved",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "Согласован",
    "color": "#86AA60"
  },
  {
    "id": null,
    "system_status": "dispatched",
    "permalink": "dispatched",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "Отгружен",
    "color": "#3A9FC9"
  },
  {
    "id": null,
    "system_status": "delivered",
    "permalink": "delivered",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "Доставлен",
    "color": "#249641"
  },
  {
    "id": null,
    "system_status": "returning",
    "permalink": "returning",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "В процессе возврата",
    "color": "#686A69"
  },
  {
    "id": null,
    "system_status": "declined",
    "permalink": "declined",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "Отменен",
    "color": "#929594"
  },
  {
    "id": null,
    "system_status": "returned",
    "permalink": "returned",
    "position": null,
    "is_default": true,
    "created_at": null,
    "updated_at": null,
    "title": "Возврат",
    "color": "#A69196"
  }
]
```

### Update custom status

`PUT /admin/custom_statuses/:custom_status_permalink.json`

For more information about custom statuses More info

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `custom_status[title]` |  | discount |

**Запрос**

```http
PUT /admin/custom_statuses/novyy.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "custom_status": {
    "title": "New (ordered by website)"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

---

## DeliveryVariant

| Операция | Метод и путь |
|----------|--------------|
| [Create delivery_variant External](#create-delivery-variant-external) | `POST /admin/delivery_variants.json` |
| [Create delivery_variant Fixed](#create-delivery-variant-fixed) | `POST /admin/delivery_variants.json` |
| [Create delivery_variant LocationDepend](#create-delivery-variant-locationdepend) | `POST /admin/delivery_variants.json` |
| [Create delivery_variant PickUp](#create-delivery-variant-pickup) | `POST /admin/delivery_variants.json` |
| [Create delivery_variant PriceDepend](#create-delivery-variant-pricedepend) | `POST /admin/delivery_variants.json` |
| [Destroy delivery_variant](#destroy-delivery-variant) | `DELETE /admin/delivery_variants/:delivery_variant_id.json` |
| [Get delivery_variant](#get-delivery-variant) | `GET /admin/delivery_variants/:delivery_variant_id.json` |
| [Get delivery_variants](#get-delivery-variants) | `GET /admin/delivery_variants.json` |
| [Update delivery variant](#update-delivery-variant) | `PUT /admin/delivery_variants/:delivery_variant_id.json` |

### Create delivery_variant External

`POST /admin/delivery_variants.json`

DeliveryVariant::External is an external delivery variant.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `delivery_variant[title]` | да | title |
| `delivery_variant[type]` | да | DeliveryVariant::External, DeliveryVariant::Fixed, DeliveryVariant::LocationDepend, DeliveryVariant::MapZoneDelivery, DeliveryVariant::NoDelivery, DeliveryVariant::None, DeliveryVariant::OfficialEmsShipping, DeliveryVariant::OfficialFirstClass, DeliveryVariant::OfficialRussianpost, DeliveryVariant::PickUp, DeliveryVariant::PriceDepend |
| `delivery_variant[description]` |  | description |
| `delivery_variant[position]` |  | position |
| `delivery_variant[add_payment_gateways]` |  | set true to add all payment variants |
| `delivery_variant[delivery_locations_attributes]` |  | array of delivery locations |
| `delivery_variant[delivery_locations_attributes][country]` |  | Optional. Default: RU. It should be enabled in checkout options to be binded to shipping method. |
| `delivery_variant[delivery_locations_attributes][region]` |  | Required. Field Region should be enabled in checkout options to be binded to shipping method. Examples: "Ивановская", "Омская" |
| `delivery_variant[delivery_locations_attributes][area]` |  | Optional. Ignored if city present. Field Region should be filled. Examples: "Знаменский" |
| `delivery_variant[delivery_locations_attributes][city]` |  | Optional. Field Region should be filled. Examples: "Омск" |
| `delivery_variant[inverted]` |  | false by default, if true this delivery variant will be available for all regions except passed in parameters 'country', 'region' and 'city' |
| `delivery_variant[charge_up_to]` |  | free on orders from given sum |
| `delivery_variant[javascript]` |  | if needed, paste javascript to be displayed on page of delivery variant selection |
| `delivery_variant[url]` |  | url for sending request on delivery price estimation |
| `delivery_variant[api_version]` |  | 'v1' or 'v2' (recomended) |
| `delivery_variant[customer_pickup]` |  | false by default, set it true if this delivery variant is customer pickup |

**Запрос**

```http
POST /admin/delivery_variants.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "delivery_variant": {
    "title": "express delivery",
    "type": "DeliveryVariant::External",
    "description": "text",
    "position": 1,
    "add_payment_gateways": true,
    "delivery_locations_attributes": [
      {
        "region": "A",
        "city": "B"
      }
    ],
    "inverted": false,
    "charge_up_to": 1000,
    "javascript": "<script type='text/javascript' src='http://mydelivery.ru/delivery.js'></script>\n        <script type='text/javascript'>\n          $(function(){\n          alert('Choose a point of customer pickup');\n          });\n        </script>",
    "url": "http://google.ru"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 49,
  "position": 1,
  "updated_at": "2026-07-25T00:58:49.000+03:00",
  "created_at": "2026-07-25T00:58:49.000+03:00",
  "charge_up_to": "1000.0",
  "inverted": false,
  "show_zip_code": false,
  "title": "express delivery",
  "description": "text",
  "type": "DeliveryVariant::External",
  "min_order_sum": null,
  "show_calendar_in_checkout": null,
  "delivery_date_required": null,
  "forbid_weekends": null,
  "utc_from_address": null,
  "estimated_delivery_period": null,
  "order_time_is_later": null,
  "forbidden_days": null,
  "show_time_intervals_in_checkout": null,
  "time_intervals_required": null,
  "url": "http://google.ru",
  "javascript": "<script type='text/javascript' src='http://mydelivery.ru/delivery.js'></script>\n        <script type='text/javascript'>\n          $(function(){\n          alert('Choose a point of customer pickup');\n          });\n        </script>",
  "customer_pickup": false,
  "payment_delivery_variants": [
    {
      "id": 12,
      "payment_gateway_id": 1,
      "created_at": "2026-07-25T00:58:49.000+03:00"
    },
    {
      "id": 13,
      "payment_gateway_id": 3,
      "created_at": "2026-07-25T00:58:49.000+03:00"
    },
    {
      "id": 14,
      "payment_gateway_id": 4,
      "created_at": "2026-07-25T00:58:49.000+03:00"
    }
  ],
  "delivery_locations": [
    {
      "id": 25,
      "city": "B",
      "region": "A",
      "country": "RU"
    }
  ],
  "delivery_time_intervals": []
}
```

### Create delivery_variant Fixed

`POST /admin/delivery_variants.json`

DeliveryVariant::Fixed is a delivery variant with a fixed delivery price.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `delivery_variant[title]` | да | title |
| `delivery_variant[type]` | да | DeliveryVariant::External, DeliveryVariant::Fixed, DeliveryVariant::LocationDepend, DeliveryVariant::MapZoneDelivery, DeliveryVariant::NoDelivery, DeliveryVariant::None, DeliveryVariant::OfficialEmsShipping, DeliveryVariant::OfficialFirstClass, DeliveryVariant::OfficialRussianpost, DeliveryVariant::PickUp, DeliveryVariant::PriceDepend |
| `delivery_variant[description]` |  | description |
| `delivery_variant[position]` |  | position |
| `delivery_variant[add_payment_gateways]` |  | set true to add all payment variants |
| `delivery_variant[delivery_locations_attributes]` |  | array of delivery locations |
| `delivery_variant[delivery_locations_attributes][country]` |  | if shop has many countries, fill this field;it should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][region]` |  | field Region should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][city]` |  | field Region should be filled |
| `delivery_variant[inverted]` |  | false by default, if true this delivery variant will be available for all regions except passed in parameters 'country', 'region' and 'city' |
| `delivery_variant[charge_up_to]` |  | free on orders from given sum |
| `delivery_variant[price]` | да | delivery price |
| `delivery_variant[javascript]` |  | if needed, paste javascript to be displayed on page of delivery variant selection |
| `delivery_variant[show_calendar_in_checkout]` |  | show delivery date selection in checkout |
| `delivery_variant[delivery_date_]` | да | delivery date selection is required in chekout |
| `delivery_variant[forbid_weekends]` |  | weekend days are forbidden for delivery |
| `delivery_variant[estimated_delivery_period]` |  | estimated delivery period |
| `delivery_variant[order_time_is_later]` |  | increase estimated delivery time for one day if order time is later |
| `delivery_variant[forbidden_days]` |  | forbidden days for delivery (23.02,08.03,..) |
| `delivery_variant[show_time_intervals_in_checkout]` |  | show delivery time intervals in checkout |
| `delivery_variant[time_intervals_]` | да | delivery time interval selection is required in checkout |
| `delivery_variant[time_intervals]` |  | time intervals (10-12, 12-16, ...) |

**Запрос**

```http
POST /admin/delivery_variants.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "delivery_variant": {
    "title": "express delivery",
    "type": "DeliveryVariant::Fixed",
    "description": "text",
    "position": 1,
    "add_payment_gateways": true,
    "delivery_locations_attributes": [
      {
        "region": "A",
        "city": "B"
      }
    ],
    "inverted": false,
    "charge_up_to": 1000,
    "price": 300,
    "show_calendar_in_checkout": true,
    "delivery_date_required": true,
    "forbid_weekends": true,
    "estimated_delivery_period": 2,
    "order_time_is_later": "18:00",
    "forbidden_days": "23.02,08.03,01.05",
    "show_time_intervals_in_checkout": true,
    "time_intervals_required": true,
    "time_intervals": "10-12,12-16,16-20"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 33,
  "position": 1,
  "updated_at": "2026-07-25T00:58:49.000+03:00",
  "created_at": "2026-07-25T00:58:49.000+03:00",
  "charge_up_to": "1000.0",
  "inverted": false,
  "show_zip_code": false,
  "title": "express delivery",
  "description": "text",
  "type": "DeliveryVariant::Fixed",
  "min_order_sum": null,
  "show_calendar_in_checkout": true,
  "delivery_date_required": true,
  "forbid_weekends": true,
  "utc_from_address": null,
  "estimated_delivery_period": 2,
  "order_time_is_later": "18:00",
  "forbidden_days": "23.02,08.03,01.05",
  "show_time_intervals_in_checkout": true,
  "time_intervals_required": true,
  "price": "300.0",
  "javascript": null,
  "payment_delivery_variants": [
    {
      "id": 6,
      "payment_gateway_id": 1,
      "created_at": "2026-07-25T00:58:49.000+03:00"
    },
    {
      "id": 7,
      "payment_gateway_id": 3,
      "created_at": "2026-07-25T00:58:49.000+03:00"
    },
    {
      "id": 8,
      "payment_gateway_id": 4,
      "created_at": "2026-07-25T00:58:49.000+03:00"
    }
  ],
  "delivery_locations": [
    {
      "id": 9,
      "city": "B",
      "region": "A",
      "country": "RU"
    }
  ],
  "delivery_time_intervals": [
    {
      "time_from": "10:00",
      "time_to": "12:00",
      "monday": true,
      "tuesday": true,
      "wednesday": true,
      "thursday": true,
      "friday": true,
      "saturday": true,
      "sunday": true
    },
    {
      "time_from": "12:00",
      "time_to": "16:00",
      "monday": true,
      "tuesday": true,
      "wednesday": true,
      "thursday": true,
      "friday": true,
      "saturday": true,
      "sunday": true
    },
    {
      "time_from": "16:00",
      "time_to": "20:00",
      "monday": true,
      "tuesday": true,
      "wednesday": true,
      "thursday": true,
      "friday": true,
      "saturday": true,
      "sunday": true
    }
  ]
}
```

### Create delivery_variant LocationDepend

`POST /admin/delivery_variants.json`

DeliveryVariant::LocationDepend is a delivery variant where delivery price depends on location (through delivery zones) and weight of delivered object.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `delivery_variant[title]` | да | title |
| `delivery_variant[type]` | да | DeliveryVariant::External, DeliveryVariant::Fixed, DeliveryVariant::LocationDepend, DeliveryVariant::MapZoneDelivery, DeliveryVariant::NoDelivery, DeliveryVariant::None, DeliveryVariant::OfficialEmsShipping, DeliveryVariant::OfficialFirstClass, DeliveryVariant::OfficialRussianpost, DeliveryVariant::PickUp, DeliveryVariant::PriceDepend |
| `delivery_variant[description]` |  | description |
| `delivery_variant[position]` |  | position |
| `delivery_variant[add_payment_gateways]` |  | set true to add all payment variants |
| `delivery_variant[delivery_locations_attributes]` |  | array of delivery locations |
| `delivery_variant[delivery_locations_attributes][country]` |  | if shop has many countries, fill this field;it should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][region]` |  | field Region should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][city]` |  | field Region should be filled |
| `delivery_variant[inverted]` |  | false by default, if true this delivery variant will be available for all regions except passed in parameters 'country', 'region' and 'city' |
| `delivery_variant[charge_up_to]` |  | free on orders from given sum |
| `delivery_variant[delivery_zones_attributes]` | да | array of delivery zones |
| `delivery_variant[delivery_zones_attributes][title]` |  | title of delivery zone |
| `delivery_variant[delivery_zones_attributes][tariffs_attributes]` | да | array of tariffs attributes for delivery zone |
| `delivery_variant[delivery_zones_attributes][tariffs_attributes][price]` | да | delivery tariff price |
| `delivery_variant[delivery_zones_attributes][tariffs_attributes][step]` | да | increment which increases delivery price; measured in kg |
| `delivery_variant[delivery_zones_attributes][tariffs_attributes][step_price]` | да | price of increment |
| `delivery_variant[delivery_zones_attributes][tariffs_attributes][max_weight]` | да | max weight for tariff |
| `delivery_variant[delivery_zones_attributes][locations_attributes]` |  | array of locations in delivery zone |
| `delivery_variant[delivery_zones_attributes][locations_attributes][country]` | да | field country should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_zones_attributes][locations_attributes][region]` |  | field Region should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][locations_attributes][city]` |  | field Region should be filled |

**Запрос**

```http
POST /admin/delivery_variants.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "delivery_variant": {
    "title": "express delivery",
    "type": "DeliveryVariant::LocationDepend",
    "description": "text",
    "position": 1,
    "add_payment_gateways": true,
    "delivery_locations_attributes": [
      {
        "region": "A",
        "city": "B"
      }
    ],
    "inverted": true,
    "charge_up_to": 1000,
    "delivery_zones_attributes": [
      {
        "title": "test",
        "tariffs_attributes": [
          {
            "max_weight": 30,
            "price": 100,
            "step": 2,
            "step_price": 50
          }
        ],
        "locations_attributes": [
          {
            "country": "A",
            "region": "C",
            "city": "D"
          }
        ]
      }
    ]
  },
  "delivery_variant[delivery_zones_attributes]": {
    "title": "express delivery"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 65,
  "position": 1,
  "updated_at": "2026-07-25T00:58:50.000+03:00",
  "created_at": "2026-07-25T00:58:50.000+03:00",
  "charge_up_to": "1000.0",
  "inverted": true,
  "show_zip_code": false,
  "title": "express delivery",
  "description": "text",
  "type": "DeliveryVariant::LocationDepend",
  "min_order_sum": null,
  "show_calendar_in_checkout": null,
  "delivery_date_required": null,
  "forbid_weekends": null,
  "utc_from_address": null,
  "estimated_delivery_period": null,
  "order_time_is_later": null,
  "forbidden_days": null,
  "show_time_intervals_in_checkout": null,
  "time_intervals_required": null,
  "payment_delivery_variants": [
    {
      "id": 18,
      "payment_gateway_id": 1,
      "created_at": "2026-07-25T00:58:50.000+03:00"
    },
    {
      "id": 19,
      "payment_gateway_id": 3,
      "created_at": "2026-07-25T00:58:50.000+03:00"
    },
    {
      "id": 20,
      "payment_gateway_id": 4,
      "created_at": "2026-07-25T00:58:50.000+03:00"
    }
  ],
  "delivery_locations": [
    {
      "id": 49,
      "city": "B",
      "region": "A",
      "country": "RU"
    }
  ],
  "delivery_time_intervals": [],
  "delivery_zones": [
    {
      "id": 9,
      "title": "test",
      "tariffs": [
        {
          "id": 9,
          "max_weight": "30.0",
          "price": "100.0",
          "step": "2.0",
          "step_price": "50.0"
        }
      ],
      "locations": [
        {
          "id": 57,
          "country": "A",
          "region": "C",
          "city": "D"
        }
      ]
    }
  ]
}
```

### Create delivery_variant PickUp

`POST /admin/delivery_variants.json`

DeliveryVariant::PickUp is a delivery variant with a fixed delivery price.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `delivery_variant[title]` | да | title |
| `delivery_variant[type]` | да | DeliveryVariant::External, DeliveryVariant::Fixed, DeliveryVariant::LocationDepend, DeliveryVariant::MapZoneDelivery, DeliveryVariant::NoDelivery, DeliveryVariant::None, DeliveryVariant::OfficialEmsShipping, DeliveryVariant::OfficialFirstClass, DeliveryVariant::OfficialRussianpost, DeliveryVariant::PickUp, DeliveryVariant::PriceDepend |
| `delivery_variant[description]` |  | description |
| `delivery_variant[position]` |  | position |
| `delivery_variant[add_payment_gateways]` |  | set true to add all payment variants |
| `delivery_variant[delivery_locations_attributes]` |  | array of delivery locations |
| `delivery_variant[delivery_locations_attributes][country]` |  | if shop has many countries, fill this field;it should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][region]` |  | field Region should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][city]` |  | field Region should be filled |
| `delivery_variant[inverted]` |  | false by default, if true this delivery variant will be available for all regions except passed in parameters 'country', 'region' and 'city' |
| `delivery_variant[charge_up_to]` |  | free on orders from given sum |
| `delivery_variant[javascript]` |  | if needed, paste javascript to be displayed on page of delivery variant selection |
| `delivery_variant[show_calendar_in_checkout]` |  | show delivery date selection in checkout |
| `delivery_variant[delivery_date_]` | да | delivery date selection is required in chekout |
| `delivery_variant[forbid_weekends]` |  | weekend days are forbidden for delivery |
| `delivery_variant[estimated_delivery_period]` |  | estimated delivery period |
| `delivery_variant[order_time_is_later]` |  | increase estimated delivery time for one day if order time is later |
| `delivery_variant[forbidden_days]` |  | forbidden days for delivery (23.02,08.03,..) |
| `delivery_variant[show_time_intervals_in_checkout]` |  | show delivery time intervals in checkout |
| `delivery_variant[time_intervals_]` | да | delivery time interval selection is required in checkout |
| `delivery_variant[time_intervals]` |  | time intervals (10-12, 12-16, ...) |
| `delivery_variant[pick_up_sources_attributes]` |  | array of pick up sources |
| `delivery_variant[pick_up_sources_attributes][title]` | да | title |
| `delivery_variant[pick_up_sources_attributes][http_method]` | да | http method (GET, POST) used for request |
| `delivery_variant[pick_up_sources_attributes][url]` | да | url for receiving list of delivery points |
| `delivery_variant[pick_up_sources_attributes][point_info_url]` | да | url for receiving delivery info to the selected point |

**Запрос**

```http
POST /admin/delivery_variants.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "delivery_variant": {
    "title": "express delivery",
    "type": "DeliveryVariant::PickUp",
    "description": "text",
    "position": 1,
    "add_payment_gateways": true,
    "delivery_locations_attributes": [
      {
        "region": "A",
        "city": "B"
      }
    ],
    "inverted": false,
    "charge_up_to": 1000,
    "show_calendar_in_checkout": true,
    "delivery_date_required": true,
    "forbid_weekends": true,
    "estimated_delivery_period": 2,
    "order_time_is_later": "18:00",
    "forbidden_days": "23.02,08.03,01.05",
    "show_time_intervals_in_checkout": true,
    "time_intervals_required": true,
    "time_intervals": "10-12,12-16,16-20",
    "pick_up_sources_attributes": [
      {
        "title": "My company",
        "pick_up_source_http_method": "POST",
        "url": "https://example.com/get_points",
        "point_info_url": "https://example.com/get_point_info"
      }
    ]
  },
  "delivery_variant[pick_up_sources_attributes]": {
    "title": "express delivery",
    "http_method": "post"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 97,
  "position": 1,
  "updated_at": "2026-07-25T00:58:51.000+03:00",
  "created_at": "2026-07-25T00:58:51.000+03:00",
  "charge_up_to": "1000.0",
  "inverted": false,
  "show_zip_code": false,
  "title": "express delivery",
  "description": "text",
  "type": "DeliveryVariant::PickUp",
  "min_order_sum": null,
  "show_calendar_in_checkout": true,
  "delivery_date_required": true,
  "forbid_weekends": true,
  "utc_from_address": null,
  "estimated_delivery_period": 2,
  "order_time_is_later": "18:00",
  "forbidden_days": "23.02,08.03,01.05",
  "show_time_intervals_in_checkout": true,
  "time_intervals_required": true,
  "payment_delivery_variants": [
    {
      "id": 32,
      "payment_gateway_id": 1,
      "created_at": "2026-07-25T00:58:51.000+03:00"
    },
    {
      "id": 33,
      "payment_gateway_id": 3,
      "created_at": "2026-07-25T00:58:51.000+03:00"
    },
    {
      "id": 34,
      "payment_gateway_id": 4,
      "created_at": "2026-07-25T00:58:51.000+03:00"
    },
    {
      "id": 35,
      "payment_gateway_id": 20,
      "created_at": "2026-07-25T00:58:51.000+03:00"
    }
  ],
  "delivery_locations": [
    {
      "id": 89,
      "city": "B",
      "region": "A",
      "country": "RU"
    }
  ],
  "delivery_time_intervals": [
    {
      "time_from": "10:00",
      "time_to": "12:00",
      "monday": true,
      "tuesday": true,
      "wednesday": true,
      "thursday": true,
      "friday": true,
      "saturday": true,
      "sunday": true
    },
    {
      "time_from": "12:00",
      "time_to": "16:00",
      "monday": true,
      "tuesday": true,
      "wednesday": true,
      "thursday": true,
      "friday": true,
      "saturday": true,
      "sunday": true
    },
    {
      "time_from": "16:00",
      "time_to": "20:00",
      "monday": true,
      "tuesday": true,
      "wednesday": true,
      "thursday": true,
      "friday": true,
      "saturday": true,
      "sunday": true
    }
  ],
  "pick_up_sources": [
    {
      "id": 2,
      "url": "https://example.com/get_points",
      "point_info_url": "https://example.com/get_point_info",
      "method": "POST",
      "created_at": "2026-07-25T00:58:51.417+03:00"
    }
  ],
  "pick_up_source_delivery_variants": [
    {
      "id": 2,
      "pick_up_source_id": 2,
      "created_at": "2026-07-25T00:58:51.420+03:00"
    }
  ]
}
```

### Create delivery_variant PriceDepend

`POST /admin/delivery_variants.json`

DeliveryVariant::PriceDepend is a delivery variant with rules which define delivery price depending on order sum and weight.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `delivery_variant[title]` | да | title |
| `delivery_variant[type]` | да | DeliveryVariant::External, DeliveryVariant::Fixed, DeliveryVariant::LocationDepend, DeliveryVariant::MapZoneDelivery, DeliveryVariant::NoDelivery, DeliveryVariant::None, DeliveryVariant::OfficialEmsShipping, DeliveryVariant::OfficialFirstClass, DeliveryVariant::OfficialRussianpost, DeliveryVariant::PickUp, DeliveryVariant::PriceDepend |
| `delivery_variant[description]` |  | description |
| `delivery_variant[position]` |  | position |
| `delivery_variant[add_payment_gateways]` |  | set true to add all payment variants |
| `delivery_variant[delivery_locations_attributes]` |  | array of delivery locations |
| `delivery_variant[delivery_locations_attributes][country]` |  | if shop has many countries, fill this field;it should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][region]` |  | field Region should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][city]` |  | field Region should be filled |
| `delivery_variant[inverted]` |  | false by default, if true this delivery variant will be available for all regions except passed in parameters 'country', 'region' and 'city' |
| `delivery_variant[charge_up_to]` |  | free on orders from given sum |
| `delivery_variant[min_price]` |  | minimal delivery price |
| `delivery_variant[rules_attributes]` |  | array of rules define delivery price depending on order sum and weight |
| `delivery_variant[rules_attributes][price]` |  | delivery price by a rule |
| `delivery_variant[rules_attributes][order_price]` |  | order sum by a rule |
| `delivery_variant[rules_attributes][order_weight]` |  | order weight by a rule |

**Запрос**

```http
POST /admin/delivery_variants.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "delivery_variant": {
    "title": "express delivery",
    "type": "DeliveryVariant::PriceDepend",
    "description": "text",
    "position": 1,
    "add_payment_gateways": true,
    "delivery_locations_attributes": [
      {
        "region": "A",
        "city": "B"
      }
    ],
    "inverted": false,
    "charge_up_to": 1000,
    "min_price": 300,
    "rules_attributes": [
      {
        "price": 300,
        "order_price": 10,
        "order_weight": 1
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 81,
  "position": 1,
  "updated_at": "2026-07-25T00:58:50.000+03:00",
  "created_at": "2026-07-25T00:58:50.000+03:00",
  "charge_up_to": "1000.0",
  "inverted": false,
  "show_zip_code": false,
  "title": "express delivery",
  "description": "text",
  "type": "DeliveryVariant::PriceDepend",
  "min_order_sum": null,
  "show_calendar_in_checkout": null,
  "delivery_date_required": null,
  "forbid_weekends": null,
  "utc_from_address": null,
  "estimated_delivery_period": null,
  "order_time_is_later": null,
  "forbidden_days": null,
  "show_time_intervals_in_checkout": null,
  "time_intervals_required": null,
  "min_price": "300.0",
  "payment_delivery_variants": [
    {
      "id": 24,
      "payment_gateway_id": 1,
      "created_at": "2026-07-25T00:58:50.000+03:00"
    },
    {
      "id": 25,
      "payment_gateway_id": 3,
      "created_at": "2026-07-25T00:58:50.000+03:00"
    },
    {
      "id": 26,
      "payment_gateway_id": 4,
      "created_at": "2026-07-25T00:58:50.000+03:00"
    }
  ],
  "delivery_locations": [
    {
      "id": 73,
      "city": "B",
      "region": "A",
      "country": "RU"
    }
  ],
  "delivery_time_intervals": [],
  "rules": [
    {
      "id": 24,
      "order_price": "10.0",
      "price": "300.0",
      "order_weight": "1.0"
    }
  ]
}
```

### Destroy delivery_variant

`DELETE /admin/delivery_variants/:delivery_variant_id.json`

**Запрос**

```http
DELETE /admin/delivery_variants/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get delivery_variant

`GET /admin/delivery_variants/:delivery_variant_id.json`

**Запрос**

```http
GET /admin/delivery_variants/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "position": 2,
  "updated_at": "2009-03-31T12:20:34.000+04:00",
  "created_at": "2009-03-30T04:37:32.000+04:00",
  "charge_up_to": null,
  "inverted": false,
  "show_zip_code": false,
  "title": "курьер",
  "description": "",
  "type": "DeliveryVariant::PriceDepend",
  "min_order_sum": null,
  "show_calendar_in_checkout": null,
  "delivery_date_required": null,
  "forbid_weekends": null,
  "utc_from_address": null,
  "estimated_delivery_period": null,
  "order_time_is_later": null,
  "forbidden_days": null,
  "show_time_intervals_in_checkout": null,
  "time_intervals_required": null,
  "min_price": 200.0,
  "payment_delivery_variants": [
    {
      "id": 2,
      "payment_gateway_id": 1,
      "created_at": "2026-07-25T00:58:09.000+03:00"
    }
  ],
  "delivery_locations": [],
  "delivery_time_intervals": [],
  "rules": [
    {
      "id": 1,
      "order_price": 2000.0,
      "price": 600.0,
      "order_weight": 10
    },
    {
      "id": 2,
      "order_price": 2000.0,
      "price": 700.0,
      "order_weight": 20
    },
    {
      "id": 3,
      "order_price": 2000.0,
      "price": 800.0,
      "order_weight": null
    },
    {
      "id": 4,
      "order_price": 10000.0,
      "price": 300.0,
      "order_weight": 10
    },
    {
      "id": 5,
      "order_price": 10000.0,
      "price": 400.0,
      "order_weight": 20
    },
    {
      "id": 6,
      "order_price": 10000.0,
      "price": 500.0,
      "order_weight": null
    },
    {
      "id": 7,
      "order_price": null,
      "price": 0.0,
      "order_weight": 10
    },
    {
      "id": 8,
      "order_price": null,
      "price": 100.0,
      "order_weight": 20
    }
  ]
}
```

### Get delivery_variants

`GET /admin/delivery_variants.json`

**Запрос**

```http
GET /admin/delivery_variants.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "position": 1,
    "updated_at": "2009-03-31T12:20:34.000+04:00",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "charge_up_to": null,
    "inverted": false,
    "show_zip_code": false,
    "title": "самовывоз",
    "description": "самовывоз из офиса",
    "type": "DeliveryVariant::Fixed",
    "min_order_sum": null,
    "show_calendar_in_checkout": null,
    "delivery_date_required": null,
    "forbid_weekends": null,
    "utc_from_address": null,
    "estimated_delivery_period": null,
    "order_time_is_later": null,
    "forbidden_days": null,
    "show_time_intervals_in_checkout": null,
    "time_intervals_required": null,
    "price": 100.0,
    "javascript": null,
    "payment_delivery_variants": [
      {
        "id": 1,
        "payment_gateway_id": 1,
        "created_at": "2026-07-25T00:58:09.000+03:00"
      }
    ],
    "delivery_locations": [],
    "delivery_time_intervals": []
  },
  {
    "id": 2,
    "position": 2,
    "updated_at": "2009-03-31T12:20:34.000+04:00",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "charge_up_to": null,
    "inverted": false,
    "show_zip_code": false,
    "title": "курьер",
    "description": "",
    "type": "DeliveryVariant::PriceDepend",
    "min_order_sum": null,
    "show_calendar_in_checkout": null,
    "delivery_date_required": null,
    "forbid_weekends": null,
    "utc_from_address": null,
    "estimated_delivery_period": null,
    "order_time_is_later": null,
    "forbidden_days": null,
    "show_time_intervals_in_checkout": null,
    "time_intervals_required": null,
    "min_price": 200.0,
    "payment_delivery_variants": [
      {
        "id": 2,
        "payment_gateway_id": 1,
        "created_at": "2026-07-25T00:58:09.000+03:00"
      }
    ],
    "delivery_locations": [],
    "delivery_time_intervals": [],
    "rules": [
      {
        "id": 1,
        "order_price": 2000.0,
        "price": 600.0,
        "order_weight": 10
      },
      {
        "id": 2,
        "order_price": 2000.0,
        "price": 700.0,
        "order_weight": 20
      },
      {
        "id": 3,
        "order_price": 2000.0,
        "price": 800.0,
        "order_weight": null
      },
      {
        "id": 4,
        "order_price": 10000.0,
        "price": 300.0,
        "order_weight": 10
      },
      {
        "id": 5,
        "order_price": 10000.0,
        "price": 400.0,
        "order_weight": 20
      },
      {
        "id": 6,
        "order_price": 10000.0,
        "price": 500.0,
        "order_weight": null
      },
      {
        "id": 7,
        "order_price": null,
        "price": 0.0,
        "order_weight": 10
      },
      {
        "id": 8,
        "order_price": null,
        "price": 100.0,
        "order_weight": 20
      }
    ]
  },
  {
    "id": 3,
    "position": 3,
    "updated_at": "2009-03-31T12:20:34.000+04:00",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "charge_up_to": null,
    "inverted": false,
    "show_zip_code": false,
    "title": "Внешний",
    "description": "Внешний способ",
    "type": "DeliveryVariant::External",
    "min_order_sum": null,
    "show_calendar_in_checkout": null,
    "delivery_date_required": null,
    "forbid_weekends": null,
    "utc_from_address": null,
    "estimated_delivery_period": null,
    "order_time_is_later": null,
    "forbidden_days": null,
    "show_time_intervals_in_checkout": null,
    "time_intervals_required": null,
    "url": "url",
    "javascript": null,
    "customer_pickup": false,
    "payment_delivery_variants": [],
    "delivery_locations": [],
    "delivery_time_intervals": []
  },
  {
    "id": 4,
    "position": 4,
    "updated_at": "2009-03-31T12:20:34.000+04:00",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "charge_up_to": null,
    "inverted": false,
    "show_zip_code": false,
    "title": "Внешний",
    "description": "Внешний способ",
    "type": "DeliveryVariant::MapZoneDelivery",
    "min_order_sum": null,
    "show_calendar_in_checkout": null,
    "delivery_date_required": null,
    "forbid_weekends": null,
    "utc_from_address": null,
    "estimated_delivery_period": null,
    "order_time_is_later": null,
    "forbidden_days": null,
    "show_time_intervals_in_checkout": null,
    "time_intervals_required": null,
    "payment_delivery_variants": [],
    "delivery_locations": [],
    "delivery_time_intervals": []
  }
]
```

### Update delivery variant

`PUT /admin/delivery_variants/:delivery_variant_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `delivery_variant[title]` |  | title |
| `delivery_variant[description]` |  | description |
| `delivery_variant[position]` |  | position |
| `delivery_variant[add_payment_gateways]` |  | set true to add all payment variants |
| `delivery_variant[charge_up_to]` |  | free on orders from given sum |
| `delivery_variant[payment_delivery_variants_attributes]` |  | array of available payment gateways for this delivery variant |
| `delivery_variant[pick_up_sources_attributes]` |  | array of pick up sources |
| `delivery_variant[payment_delivery_variants_attributes][_destroy]` |  | set it 1 to delete this payment variant |
| `delivery_variant[payment_delivery_variants_attributes][id]` |  | only for removal: id of delivery variant and payment gateway binding you want to delete |
| `delivery_variant[payment_delivery_variants_attributes][payment_gateway_id]` |  | for creation: id of payment gateway you want to bind with this delivery variant |

**Запрос**

```http
PUT /admin/delivery_variants/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "delivery_variant": {
    "title": "new title",
    "description": "new description",
    "position": 1,
    "charge_up_to": 2500,
    "payment_delivery_variants_attributes": [
      {
        "_destroy": 1,
        "id": 2
      },
      {
        "payment_gateway_id": 3
      }
    ],
    "pick_up_sources_attributes": [
      {
        "title": "My company",
        "pick_up_source_http_method": "POST",
        "url": "https://example.com/get_points",
        "point_info_url": "https://example.com/get_point_info"
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "position": 1,
  "updated_at": "2026-07-25T00:58:51.000+03:00",
  "created_at": "2009-03-30T04:37:32.000+04:00",
  "charge_up_to": "2500.0",
  "inverted": false,
  "show_zip_code": false,
  "title": "new title",
  "description": "new description",
  "type": "DeliveryVariant::PriceDepend",
  "min_order_sum": null,
  "show_calendar_in_checkout": null,
  "delivery_date_required": null,
  "forbid_weekends": null,
  "utc_from_address": null,
  "estimated_delivery_period": null,
  "order_time_is_later": null,
  "forbidden_days": null,
  "show_time_intervals_in_checkout": null,
  "time_intervals_required": null,
  "min_price": 200.0,
  "payment_delivery_variants": [
    {
      "id": 42,
      "payment_gateway_id": 3,
      "created_at": "2026-07-25T00:58:51.000+03:00"
    }
  ],
  "delivery_locations": [],
  "delivery_time_intervals": [],
  "rules": [
    {
      "id": 1,
      "order_price": 2000.0,
      "price": 600.0,
      "order_weight": 10
    },
    {
      "id": 2,
      "order_price": 2000.0,
      "price": 700.0,
      "order_weight": 20
    },
    {
      "id": 3,
      "order_price": 2000.0,
      "price": 800.0,
      "order_weight": null
    },
    {
      "id": 4,
      "order_price": 10000.0,
      "price": 300.0,
      "order_weight": 10
    },
    {
      "id": 5,
      "order_price": 10000.0,
      "price": 400.0,
      "order_weight": 20
    },
    {
      "id": 6,
      "order_price": 10000.0,
      "price": 500.0,
      "order_weight": null
    },
    {
      "id": 7,
      "order_price": null,
      "price": 0.0,
      "order_weight": 10
    },
    {
      "id": 8,
      "order_price": null,
      "price": 100.0,
      "order_weight": 20
    }
  ]
}
```

---

## DiscountCode

| Операция | Метод и путь |
|----------|--------------|
| [Create discount code](#create-discount-code) | `POST /admin/discount_codes.json` |
| [Destroy discount code](#destroy-discount-code) | `DELETE /admin/discount_codes/:discount_code_id.json` |
| [Get discount code](#get-discount-code) | `GET /admin/discount_codes/:discount_code_id.json` |
| [Get discount codes](#get-discount-codes) | `GET /admin/discount_codes.json` |
| [Update discount code](#update-discount-code) | `PUT /admin/discount_codes/:discount_code_id.json` |

### Create discount code

`POST /admin/discount_codes.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `discount_code[code]` | да | code |
| `discount_code[description]` |  | description |
| `discount_code[disabled]` |  | disabled |
| `discount_code[act_once]` |  | act once |
| `discount_code[act_once_for_client]` |  | act once for client |
| `discount_code[expired_at]` |  | expired at |
| `discount_code[type_id]` | да | type_id (1 - percent, 2 - money) |
| `discount_code[discount]` | да | discount |
| `discount_code[min_price]` |  | min price |

**Запрос**

```http
POST /admin/discount_codes.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "discount_code": {
    "code": "BIRTHDAY",
    "description": "Birthday discount",
    "act_once_for_client": true,
    "expired_at": "2026-08-01",
    "type_id": 1,
    "discount": 10,
    "min_price": 500
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 188,
  "code": "birthday",
  "type_id": 1,
  "discount": "10.0",
  "description": "Birthday discount",
  "disabled": true,
  "act_once": true,
  "act_once_for_client": true,
  "min_price": "500.0",
  "created_at": "2026-07-25T00:59:06.614+03:00",
  "updated_at": "2026-07-25T00:59:06.618+03:00",
  "expired_at": "2026-08-01",
  "skip_discounted_product": null,
  "skip_bundles": null,
  "worked": false,
  "discount_collections": []
}
```

### Destroy discount code

`DELETE /admin/discount_codes/:discount_code_id.json`

**Запрос**

```http
DELETE /admin/discount_codes/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get discount code

`GET /admin/discount_codes/:discount_code_id.json`

**Запрос**

```http
GET /admin/discount_codes/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "code": "dc1",
  "type_id": 1,
  "discount": "30.0",
  "description": "процентная скидка",
  "disabled": false,
  "act_once": true,
  "act_once_for_client": null,
  "min_price": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2009-03-18T10:53:15.000+03:00",
  "expired_at": "2020-12-12",
  "skip_discounted_product": null,
  "skip_bundles": null,
  "worked": false,
  "discount_collections": []
}
```

### Get discount codes

`GET /admin/discount_codes.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` |  | Set datetime to get only data updated after it. |
| `from_id` |  | Set id to get only data starting from it. |
| `q` |  | Strict search by code or regexp for description. |
| `per_page` |  | Quantity of objects per page. |
| `page` |  | Page number. |

**Запрос**

```http
GET /admin/discount_codes.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 2,
    "code": "dc2",
    "type_id": 2,
    "discount": "100.0",
    "description": "фиксированная скидка",
    "disabled": false,
    "act_once": true,
    "act_once_for_client": null,
    "min_price": null,
    "created_at": "2009-03-18T10:53:14.000+03:00",
    "updated_at": "2009-03-18T10:53:15.000+03:00",
    "expired_at": "2020-12-12",
    "skip_discounted_product": null,
    "skip_bundles": null,
    "worked": false,
    "discount_collections": []
  },
  {
    "id": 1,
    "code": "dc1",
    "type_id": 1,
    "discount": "30.0",
    "description": "процентная скидка",
    "disabled": false,
    "act_once": true,
    "act_once_for_client": null,
    "min_price": null,
    "created_at": "2009-03-18T10:53:14.000+03:00",
    "updated_at": "2009-03-18T10:53:15.000+03:00",
    "expired_at": "2020-12-12",
    "skip_discounted_product": null,
    "skip_bundles": null,
    "worked": false,
    "discount_collections": []
  }
]
```

### Update discount code

`PUT /admin/discount_codes/:discount_code_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `discount_code[discount]` |  | discount |

**Запрос**

```http
PUT /admin/discount_codes/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "discount_code": {
    "discount": 15
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

---

## DiscountCollectionsGroup

| Операция | Метод и путь |
|----------|--------------|
| [Get discount codes for discount collections group](#get-discount-codes-for-discount-collections-group) | `GET /admin/discount_collections_groups/:discount_collections_group_id/discount_codes.json` |
| [Get discount collections group](#get-discount-collections-group) | `GET /admin/discount_collections_groups/:discount_collections_group_id.json` |
| [Get discount collections groups](#get-discount-collections-groups) | `GET /admin/discount_collections_groups.json` |

### Get discount codes for discount collections group

`GET /admin/discount_collections_groups/:discount_collections_group_id/discount_codes.json`

**Запрос**

```http
GET /admin/discount_collections_groups/220/discount_codes.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 204,
    "code": "coded2430f",
    "type_id": 2,
    "discount": "10.0",
    "description": null,
    "disabled": false,
    "act_once": true,
    "act_once_for_client": null,
    "min_price": null,
    "created_at": "2026-07-25T00:59:07.836+03:00",
    "updated_at": "2026-07-25T00:59:07.842+03:00",
    "expired_at": null,
    "skip_discounted_product": null,
    "skip_bundles": null,
    "worked": true,
    "discount_collections": []
  }
]
```

### Get discount collections group

`GET /admin/discount_collections_groups/:discount_collections_group_id.json`

**Запрос**

```http
GET /admin/discount_collections_groups/204.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 204,
  "title": "Group d74c3aa2",
  "account_id": 1,
  "created_at": "2026-07-25T00:59:07.000+03:00",
  "updated_at": "2026-07-25T00:59:07.000+03:00",
  "collection_ids": [
    1
  ],
  "collections": [
    {
      "id": 1,
      "title": "T-shirts",
      "parent_id": null,
      "permalink": "футболка"
    }
  ]
}
```

### Get discount collections groups

`GET /admin/discount_collections_groups.json`

**Запрос**

```http
GET /admin/discount_collections_groups.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 172,
    "title": "Discount collection group 1",
    "account_id": 1,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2026-07-25T00:58:09.000+03:00",
    "collection_ids": [
      171,
      173,
      175
    ],
    "collections": [
      {
        "id": 171,
        "title": "Копилки",
        "parent_id": 170,
        "permalink": "order_test_collection_1"
      },
      {
        "id": 173,
        "title": "Товары для комплекта со скидкой",
        "parent_id": 170,
        "permalink": "order_test_collection_3"
      },
      {
        "id": 175,
        "title": "Комплекты со скидкой",
        "parent_id": 170,
        "permalink": "order_test_collection_5"
      }
    ]
  },
  {
    "id": 188,
    "title": "Group a4588d65",
    "account_id": 1,
    "created_at": "2026-07-25T00:59:07.000+03:00",
    "updated_at": "2026-07-25T00:59:07.000+03:00",
    "collection_ids": [
      1
    ],
    "collections": [
      {
        "id": 1,
        "title": "T-shirts",
        "parent_id": null,
        "permalink": "футболка"
      }
    ]
  }
]
```

---

## Domain

| Операция | Метод и путь |
|----------|--------------|
| [Create domain](#create-domain) | `POST /admin/domains.json` |
| [Destroy domain](#destroy-domain) | `DELETE /admin/domains/:domain_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/domains/:domain_id.json` |
| [Get domains](#get-domains) | `GET /admin/domains.json` |
| [Update domain](#update-domain) | `PUT /admin/domains/:domain_id.json` |

### Create domain

`POST /admin/domains.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `domain[domain]` | да | domain |
| `domain[main]` |  | boolean, could be false only if there are more than one domain |

**Запрос**

```http
POST /admin/domains.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "domain": {
    "domain": "another-domain.ru",
    "main": true
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 18,
  "domain": "another-domain.ru",
  "created_at": "2026-07-25T00:58:53.000+03:00",
  "updated_at": "2026-07-25T00:58:53.000+03:00",
  "main": true
}
```

### Destroy domain

`DELETE /admin/domains/:domain_id.json`

**Запрос**

```http
DELETE /admin/domains/22.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/domains/:domain_id.json`

Content-Type: application/json Get domain

**Запрос**

```http
GET /admin/domains/14.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 14,
  "domain": "some-domain.ru",
  "created_at": "2026-07-25T00:58:52.000+03:00",
  "updated_at": "2026-07-25T00:58:52.000+03:00",
  "main": true
}
```

### Get domains

`GET /admin/domains.json`

**Запрос**

```http
GET /admin/domains.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 12,
    "domain": "some-domain.ru",
    "created_at": "2026-07-25T00:58:52.000+03:00",
    "updated_at": "2026-07-25T00:58:52.000+03:00",
    "main": true
  }
]
```

### Update domain

`PUT /admin/domains/:domain_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `domain[domain]` |  | domain |
| `domain[main]` |  | boolean, could be false only if there are more than one domain |

**Запрос**

```http
PUT /admin/domains/20.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "domain": {
    "domain": "other-domain.ru"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20,
  "domain": "other-domain.ru",
  "created_at": "2026-07-25T00:58:53.000+03:00",
  "updated_at": "2026-07-25T00:58:53.000+03:00",
  "main": true
}
```

---

## ExternalDiscount

| Операция | Метод и путь |
|----------|--------------|
| [Create external discount](#create-external-discount) | `POST /admin/external_discounts.json` |
| [Destroy external discount](#destroy-external-discount) | `DELETE /admin/external_discounts/:external_discount_id.json` |
| [Get external discount](#get-external-discount) | `GET /admin/external_discounts/:external_discount_id.json` |
| [Get external discounts](#get-external-discounts) | `GET /admin/external_discounts.json` |
| [Update external discount](#update-external-discount) | `PUT /admin/external_discounts/:external_discount_id.json` |

### Create external discount

`POST /admin/external_discounts.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `external_discount[description]` |  | description |
| `external_discount[url]` |  | external service url |
| `external_discount[order_lines_depend]` |  | discount calculation depends on order lines |
| `external_discount[active]` |  | active |
| `external_discount[observed_fields]` |  | discount calculation depends on this order fields (array of field ids) |

**Запрос**

```http
POST /admin/external_discounts.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "external_discount": {
    "description": "Birthday discount",
    "url": "https://myservice.com/birthday_external_discount",
    "order_lines_depend": false,
    "active": true,
    "observed_fields": [
      1,
      2
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 57,
  "url": "https://myservice.com/birthday_external_discount",
  "order_lines_depend": false,
  "active": true,
  "created_at": "2026-07-25T00:59:08.000+03:00",
  "updated_at": "2026-07-25T00:59:08.000+03:00",
  "description": "Birthday discount",
  "observed_fields": [
    1,
    2
  ]
}
```

### Destroy external discount

`DELETE /admin/external_discounts/:external_discount_id.json`

**Запрос**

```http
DELETE /admin/external_discounts/89.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get external discount

`GET /admin/external_discounts/:external_discount_id.json`

**Запрос**

```http
GET /admin/external_discounts/25.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 25,
  "url": "https://myservice.com/external_discount",
  "order_lines_depend": null,
  "active": true,
  "created_at": "2026-07-25T00:59:08.000+03:00",
  "updated_at": "2026-07-25T00:59:08.000+03:00",
  "description": "My External Discount",
  "observed_fields": []
}
```

### Get external discounts

`GET /admin/external_discounts.json`

**Запрос**

```http
GET /admin/external_discounts.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 9,
    "url": "https://myservice.com/external_discount",
    "order_lines_depend": null,
    "active": true,
    "created_at": "2026-07-25T00:59:08.000+03:00",
    "updated_at": "2026-07-25T00:59:08.000+03:00",
    "description": "My External Discount",
    "observed_fields": []
  }
]
```

### Update external discount

`PUT /admin/external_discounts/:external_discount_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `external_discount[description]` |  | description |
| `external_discount[url]` |  | external service url |
| `external_discount[order_lines_depend]` |  | discount calculation depends on order lines |
| `external_discount[active]` |  | active |

**Запрос**

```http
PUT /admin/external_discounts/73.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "external_discount": {
    "url": "https://myservice.com/birthday_external_discount_v2"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 73,
  "url": "https://myservice.com/birthday_external_discount_v2",
  "order_lines_depend": null,
  "active": true,
  "created_at": "2026-07-25T00:59:08.000+03:00",
  "updated_at": "2026-07-25T00:59:08.000+03:00",
  "description": "My External Discount",
  "observed_fields": []
}
```

---

## Field

| Операция | Метод и путь |
|----------|--------------|
| [Create field](#create-field) | `POST /admin/fields.json` |
| [Create select field](#create-select-field) | `POST /admin/fields.json` |
| [Destroy field](#destroy-field) | `DELETE /admin/fields/:field_id.json` |
| [Get field](#get-field) | `GET /admin/fields/:field_id.json` |
| [Get fields](#get-fields) | `GET /admin/fields.json` |
| [Update field](#update-field) | `PUT /admin/fields/:non_select_field_id.json` |
| [Update select field](#update-select-field) | `PUT /admin/fields/:field_id.json` |

### Create field

`POST /admin/fields.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `field[type]` | да | Field::TextField, Field::TextArea, Field::Checkbox, Field::FileField (only for Order), Field::Delivery |
| `field[office_title]` | да | title for backoffice |
| `field[destiny]` | да | integer, defines subject that field belongs to: 1 - ShippingAddress, 2 - Client::Individual, 5 - Client::Juridical, 3 - Order, 4 - OrderLine, 6 - Client |
| `field[handle]` |  | You can create or update value by handle |
| `field[title]` |  | title |
| `field[active]` |  | true by default, defines ability of filling field by a customer |
| `field[for_buyer]` |  | true by defaut, defines visibility for a customer |
| `field[obligatory]` |  | false by default |
| `field[show_in_result]` |  | true by default, field will be shown in result |
| `field[show_in_checkout]` |  | true by default, field will be shown in checkout |
| `field[position]` |  | position in fields list |
| `field[is_indexed]` |  | boolean, true if searchable |
| `field[hide_in_backoffice]` |  | boolean, true to hide in backoffice |

**Запрос**

```http
POST /admin/fields.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "field": {
    "type": "Field::TextField",
    "office_title": "title_for_text_field",
    "destiny": 1,
    "handle": "my_awesome_field_1",
    "title": "title_for_text_field",
    "active": false,
    "for_buyer": false,
    "show_in_result": false,
    "show_in_checkout": false,
    "is_indexed": true,
    "hide_in_backoffice": true
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 192,
  "destiny": 1,
  "position": 22,
  "office_title": "title_for_text_field",
  "for_buyer": false,
  "obligatory": false,
  "active": false,
  "system_name": null,
  "created_at": "2026-07-25T00:58:54.621+03:00",
  "updated_at": "2026-07-25T00:58:54.621+03:00",
  "show_in_result": false,
  "show_in_checkout": false,
  "is_indexed": true,
  "hide_in_backoffice": true,
  "handle": "my_awesome_field_1",
  "placeholder": null,
  "title": "title_for_text_field",
  "example": null,
  "type": "Field::TextField"
}
```

### Create select field

`POST /admin/fields.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `field[type]` | да | Field::Select |
| `field[office_title]` | да | title for backoffice |
| `field[destiny]` | да | integer, defines subject that field belongs to: 1 - ShippingAddress, 2 - Client::Individual, 5 - Client::Juridical, 3 - Order, 4 - OrderLine, 6 - Client |
| `field[handle]` |  | You can create or update value by handle |
| `field[title]` |  | title |
| `field[active]` |  | true by default, defines ability of filling field by a customer |
| `field[for_buyer]` |  | true by defaut, defines visibility for a customer |
| `field[obligatory]` |  | false by default |
| `field[show_in_result]` |  | true by default, field will be shown in result |
| `field[show_in_checkout]` |  | true by default, field will be shown in checkout |
| `field[position]` |  | position in fields list |
| `field[is_indexed]` |  | boolean, true if searchable |
| `field[hide_in_backoffice]` |  | boolean, true to hide in backoffice |
| `field[field_options_attributes]` |  | array of field options objects |

**Запрос**

```http
POST /admin/fields.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "field": {
    "type": "Field::Select",
    "office_title": "title_for_select_field",
    "destiny": 1,
    "handle": "my_select_field_1",
    "title": "title_for_select_field",
    "active": false,
    "for_buyer": false,
    "show_in_result": false,
    "show_in_checkout": false,
    "is_indexed": true,
    "hide_in_backoffice": true,
    "field_options_attributes": [
      {
        "value": 1
      },
      {
        "value": 2
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 208,
  "destiny": 1,
  "position": 22,
  "office_title": "title_for_select_field",
  "for_buyer": false,
  "obligatory": false,
  "active": false,
  "system_name": null,
  "created_at": "2026-07-25T00:58:55.784+03:00",
  "updated_at": "2026-07-25T00:58:55.784+03:00",
  "show_in_result": false,
  "show_in_checkout": false,
  "is_indexed": true,
  "hide_in_backoffice": true,
  "handle": "my_select_field_1",
  "placeholder": null,
  "title": "title_for_select_field",
  "example": null,
  "type": "Field::Select",
  "field_options": [
    {
      "id": 42,
      "position": 1,
      "value": "1"
    },
    {
      "id": 50,
      "position": 2,
      "value": "2"
    }
  ]
}
```

### Destroy field

`DELETE /admin/fields/:field_id.json`

Only these types can be destroyed: Field::TextField, Field::TextArea, Field::Select,Field::Checkbox, Field::Delivery, Field::FileField (only for Order)

**Запрос**

```http
DELETE /admin/fields/23.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get field

`GET /admin/fields/:field_id.json`

**Запрос**

```http
GET /admin/fields/23.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 23,
  "destiny": 2,
  "position": 6,
  "office_title": "Селект",
  "for_buyer": true,
  "obligatory": false,
  "active": false,
  "system_name": null,
  "created_at": "2010-10-12T13:31:09.692+04:00",
  "updated_at": "2010-10-12T13:31:09.692+04:00",
  "show_in_result": true,
  "show_in_checkout": true,
  "is_indexed": null,
  "hide_in_backoffice": null,
  "handle": null,
  "placeholder": null,
  "title": "Селект",
  "example": null,
  "type": "Field::Select",
  "field_options": [
    {
      "id": 1,
      "position": 1,
      "value": "Передан"
    },
    {
      "id": 2,
      "position": 2,
      "value": "Не передан"
    }
  ]
}
```

### Get fields

`GET /admin/fields.json`

**Запрос**

```http
GET /admin/fields.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "destiny": 1,
    "position": 0,
    "office_title": "Страна",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "country",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::Country"
  },
  {
    "id": 11,
    "destiny": 2,
    "position": 1,
    "office_title": "Имя",
    "for_buyer": true,
    "obligatory": true,
    "active": true,
    "system_name": "name",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::ObligatoryTextField"
  },
  {
    "id": 2,
    "destiny": 1,
    "position": 1,
    "office_title": "Регион",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "state",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 3,
    "destiny": 1,
    "position": 2,
    "office_title": "Город",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "city",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 12,
    "destiny": 2,
    "position": 2,
    "office_title": "Телефон",
    "for_buyer": true,
    "obligatory": true,
    "active": true,
    "system_name": "phone",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::Phone"
  },
  {
    "id": 4,
    "destiny": 1,
    "position": 3,
    "office_title": "Адрес",
    "for_buyer": true,
    "obligatory": false,
    "active": true,
    "system_name": "address",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextArea"
  },
  {
    "id": 13,
    "destiny": 2,
    "position": 3,
    "office_title": "E-mail",
    "for_buyer": true,
    "obligatory": true,
    "active": true,
    "system_name": "email",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 14,
    "destiny": 2,
    "position": 4,
    "office_title": "Подтверждение e-mail",
    "for_buyer": true,
    "obligatory": true,
    "active": false,
    "system_name": "email_confirmation",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 5,
    "destiny": 1,
    "position": 4,
    "office_title": "Почтовый индекс",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "zip",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 21,
    "destiny": 2,
    "position": 4,
    "office_title": "Фамилия",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "surname",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 22,
    "destiny": 2,
    "position": 5,
    "office_title": "Отчество",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "middlename",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 23,
    "destiny": 2,
    "position": 6,
    "office_title": "Селект",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": null,
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": "Селект",
    "example": null,
    "type": "Field::Select",
    "field_options": [
      {
        "id": 1,
        "position": 1,
        "value": "Передан"
      },
      {
        "id": 2,
        "position": 2,
        "value": "Не передан"
      }
    ]
  },
  {
    "id": 172,
    "destiny": 2,
    "position": 18,
    "office_title": "Согласие на обработку персональных данных",
    "for_buyer": true,
    "obligatory": true,
    "active": false,
    "system_name": "consent_to_personal_data",
    "created_at": "2026-07-25T00:58:09.211+03:00",
    "updated_at": "2026-07-25T00:58:09.211+03:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemCheckbox"
  },
  {
    "id": 106,
    "destiny": 1,
    "position": 19,
    "office_title": "Улица",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "street",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 107,
    "destiny": 1,
    "position": 20,
    "office_title": "Дом",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "house",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 108,
    "destiny": 1,
    "position": 21,
    "office_title": "Квартира",
    "for_buyer": true,
    "obligatory": false,
    "active": false,
    "system_name": "flat",
    "created_at": "2010-10-12T13:31:09.692+04:00",
    "updated_at": "2010-10-12T13:31:09.692+04:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": null,
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::SystemTextField"
  },
  {
    "id": 174,
    "destiny": 3,
    "position": 22,
    "office_title": "АКС айди отгрузки",
    "for_buyer": true,
    "obligatory": true,
    "active": false,
    "system_name": null,
    "created_at": "2026-07-25T00:58:09.211+03:00",
    "updated_at": "2026-07-25T00:58:09.211+03:00",
    "show_in_result": true,
    "show_in_checkout": true,
    "is_indexed": null,
    "hide_in_backoffice": null,
    "handle": "aks_shipment_id",
    "placeholder": null,
    "title": null,
    "example": null,
    "type": "Field::TextField"
  }
]
```

### Update field

`PUT /admin/fields/:non_select_field_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `field[office_title]` |  | title for backoffice |
| `field[title]` |  | title |
| `field[active]` |  | true by default, defines ability of filling field by a customer |
| `field[for_buyer]` |  | true by defaut, defines visibility for a customer |
| `field[obligatory]` |  | false by default |
| `field[show_in_result]` |  | true by default, field will be shown in result |
| `field[show_in_checkout]` |  | true by default, field will be shown in checkout |
| `field[position]` |  | position in fields list |
| `field[handle]` |  | You can create or update value by handle |
| `field[is_indexed]` |  | boolean, true if searchable |
| `field[hide_in_backoffice]` |  | boolean, true to hide in backoffice |

**Запрос**

```http
PUT /admin/fields/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "field": {
    "office_title": "state_code",
    "title": "State code",
    "active": true,
    "for_buyer": true,
    "obligatory": true,
    "show_in_result": true,
    "show_in_checkout": true,
    "position": 1,
    "handle": "my_awesome_field_2",
    "is_indexed": true,
    "hide_in_backoffice": true
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "destiny": 1,
  "position": 1,
  "office_title": "state_code",
  "for_buyer": true,
  "obligatory": true,
  "active": true,
  "system_name": "state",
  "created_at": "2010-10-12T13:31:09.692+04:00",
  "updated_at": "2026-07-25T00:58:54.949+03:00",
  "show_in_result": true,
  "show_in_checkout": true,
  "is_indexed": true,
  "hide_in_backoffice": true,
  "handle": "my_awesome_field_2",
  "placeholder": null,
  "title": "State code",
  "example": null,
  "type": "Field::SystemTextField"
}
```

### Update select field

`PUT /admin/fields/:field_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `field[office_title]` |  | title for backoffice |
| `field[title]` |  | title |
| `field[active]` |  | true by default, defines ability of filling field by a customer |
| `field[for_buyer]` |  | true by defaut, defines visibility for a customer |
| `field[obligatory]` |  | false by default |
| `field[show_in_result]` |  | true by default, field will be shown in result |
| `field[show_in_checkout]` |  | true by default, field will be shown in checkout |
| `field[position]` |  | position in fields list |
| `field[handle]` |  | You can create or update value by handle |
| `field[is_indexed]` |  | boolean, true if searchable |
| `field[hide_in_backoffice]` |  | boolean, true to hide in backoffice |
| `field[field_options_attributes]` |  | array of field options objects |

**Запрос**

```http
PUT /admin/fields/23.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "field": {
    "office_title": "state_code",
    "title": "State code",
    "active": true,
    "for_buyer": true,
    "obligatory": true,
    "show_in_result": true,
    "show_in_checkout": true,
    "position": 2,
    "handle": "my_awesome_field_10",
    "is_indexed": true,
    "hide_in_backoffice": true,
    "field_options_attributes": [
      {
        "value": 10
      },
      {
        "id": 1,
        "value": 20
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 23,
  "destiny": 2,
  "position": 2,
  "office_title": "state_code",
  "for_buyer": true,
  "obligatory": true,
  "active": true,
  "system_name": null,
  "created_at": "2010-10-12T13:31:09.692+04:00",
  "updated_at": "2026-07-25T00:58:55.512+03:00",
  "show_in_result": true,
  "show_in_checkout": true,
  "is_indexed": true,
  "hide_in_backoffice": true,
  "handle": "my_awesome_field_10",
  "placeholder": null,
  "title": "State code",
  "example": null,
  "type": "Field::Select",
  "field_options": [
    {
      "id": 1,
      "position": 1,
      "value": "20"
    },
    {
      "id": 2,
      "position": 2,
      "value": "Не передан"
    },
    {
      "id": 18,
      "position": 3,
      "value": "10"
    }
  ]
}
```

---

## File

| Операция | Метод и путь |
|----------|--------------|
| [Create file from attachment](#create-file-from-attachment) | `POST /admin/files.json` |
| [Create file from url](#create-file-from-url) | `POST /admin/files.json` |
| [Destroy file](#destroy-file) | `DELETE /admin/files/:account_file_id.json` |
| [Get file](#get-file) | `GET /admin/files/:account_file_id.json` |
| [Get files](#get-files) | `GET /admin/files.json` |

### Create file from attachment

`POST /admin/files.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `file[attachment]` | да | file |
| `file[filename]` | да | file name |

**Запрос**

```http
POST /admin/files.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "file": {
    "attachment": "/9j/4AAQSkZJRgABAQEASABIAAD//gA+Q1JFQVRPUjogZ2QtanBlZyB2MS4w\nICh1c2luZyBJSkcgSlBFRyB2NjIpLCBkZWZhdWx0IHF1YWxpdHkK/9sAQwAG\nBAUGBQQGBgUGBwcGCAoQCgoJCQoUDg8MEBcUGBgXFBYWGh0lHxobIxwWFiAs\nICMmJykqKRkfLTAtKDAlKCko/9sAQwEHBwcKCAoTCgoTKBoWGigoKCgoKCgo\nKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgo/8AA\nEQgAWQBkAwEiAAIRAQMRAf/EAB0AAAEEAwEBAAAAAAAAAAAAAAgABQYHAQME\nAgn/xABEEAABAwMBBAYHBgIHCQAAAAABAgMEAAURBgcSITEIE0FRYYEUIjJx\nkbHSIzd1lKGzFRg1QkNSYqLwJSYzY3KCksLT/8QAGQEAAgMBAAAAAAAAAAAA\nAAAAAAEDBAUC/8QAIxEAAgICAgIBBQAAAAAAAAAAAAECEQMEEjEhQVEFEyIz\ngf/aAAwDAQACEQMRAD8AKmlVT7Q9oF4sGqV223NxFthpC0hxoqUSQSeO8O6o\ncvbbd23ktO/w5LiuASWFDP8AmrlzivDY+LqwiaVDwNtl5KlpDUTeQMqHobvA\nfGtLm3G9IRvlFvCOPFUZwf8AtS+7D5EEbSoYv5g7qpwIaNqUonAHo7nP/wAq\n4pHSQujDq2nFWhLiCUqSqO5kEdntV1YBVUqEtfSXvGfVVZvy7n1VpV0mL7n1\nTZfyzn1UwC6pUIKukzqDs/gv5Zz661npM6j7DZfyrn10AGFSoOVdJnU39X+C\n/lXPrr2OkjqlxIDD1jW6R/wxDdz8SrFABh0qZ9IXB67aTs1xlhAkS4bL7m4M\nJ3lIBOB2DJpUAUNt3usew67NzmRm5LLbbSeqcGQSUKA4VV8rag2JoQxY7e2w\nocFdWARU86UPUpvxVKSpbG/F6xCTglPrZwe/FVDb1Q4Vyc9HjIQyyvAku4cc\nHHhugjdB8s+NUtnHjk7mrJY5JKPFPwXBZro07aWp98YMVCsKaQ4d31f76u4d\n3fXq/SrdfYhTYoMVx5KSClZwV45lIJ9Y+GKi1usV01pdWFMOmW2jGX3DlpA7\nCByJ9/Luq2bZpyBpiE21bCHZ6R68pwZDiu0b3+h4VhuEMfn38fBS2c1A26g0\n9qELUpNtdbaPFCkxktocHhkAg+BqIXiJJknfeZWiaBhSVJwXMdo7zRoL1XLa\nixw8A2p7rG+rdxlLiBk8+Y4jj41HZ1tsGso5D8CEuafW3mMNqV78cD7+Yra1\nc7nSaKq2+PaAxP61gc6tDa7oFenrupyGHHITw6xpxQ9bHahf+NJ4E9owarAg\npUQoYIrQaa7L+PIskeSMGsVk1ikdirts/wDSDfuPyNcVd1m/pFr3K+RoA+ke\nzj7vtM/hsb9tNKls4+77TP4ZG/bTSoAHXpTuFGonhkAbkfn7lVS1tdRMkIRI\nO6CrOeQzjl8R+tXJ0qsK1M8kkZ6uPgHtOFcKHp9L7LXpLCypkEJWlQzuE8s+\nBqOcFIGX7p2/Lstl/h0Wa20pWSpDDZPDwXgDPPlXbaZl1E1wpnOSG1J+xSXd\n5RPd61UDA1M5AblJVGYeL7WG1nILTg4BYAOMgZ+NXHscntahtsx6Y71a22/R\n9xtPsrUCAodw5ms2enxtmfsY3FOTHDVe0G4yFqhz7UhEtlCmGkqGVNk8CEnx\n8fGvOkhchFEiZPhxIyXApS1kNgAcgnHEnPj8a49QliU6lltaRLDhQCsZ3TnB\nxnxzjspnf07cusT10z0pCVb286d3cAPYOVa+th4pUis+LjXRONfzEaiaSxAl\nFXo32eFew6o4OFA9vcrvobr7C6mW6hbZQQo44csH/XCr2hT4iI4bQ2SFulMh\nwD1lgJGefLPLyrfd9A2fVbyZlnddW8GwpUb2G5I5ApdPAODGCOGR7q1J6jyQ\nuL8ljTuH4+gZyhQVu4416cZcax1ja0Z5bySM0Tlk0JFtqt17diKTzjxkBCz/\nANTisk+VP0pVsYgLju2mJMiKG643JCnQoeJJ/UcakX0ebjfLyauPG8nQINd1\nm4XFr3K+Rp013ZEWPUsyPGSRCU4Vx+OfszxAz4A48qa7QP8AaDXuPyNZE4OE\nnGXaOJRcW0z6R7OPu+0z+GRv200qxs3+73TP4ZG/bTSrkQMnS6WU6kfx2Nxj\n+iqoazXRCHermgLaWnq3N44C0nsPce0HvAq9el2f955A/wCVG+SqHAAUqsTV\nkkn6fUw6gAuuQnfWZkIRkY7j2Z7+NPemrmzoqUtReffYloCHWwNzkchQwSeH\nf41Cos6TFQUR5TzSFc0oWQD5V0tyWpI3JjqgvscPH40JemRyg5KpdFkS7/Gv\nOun7rBJXbmGUuvKVkJBCcqwO/IPGsW3XLc5b5fe6tTn9mTugDuHuFQiTcYkG\nxuW+3uF56Qft3sYSEZB3R3kkDJ7hTAOfCp1kcOiJa8ZLz/C/bJMt7cZW4hS2\nXVgLQlXnnOedTS33q0Wq2b8PcMxSjiO2AXFZ7cJz29tDZbwZMF9pXEhPWDvO\nOf6Zp30/f3LddIRjoyktqQnJxneBGM+PCtjX3opRjJUmSY8Siy29Rank3+0u\nzmS4xIiK3XUJO6efAntyKY7BPdvS3WH3kNvupIbcUopSF9mTngCe3srk01MQ\n9MTujMKYkpIxzJ4EK/xA8KYWUu225vRST1gWUj3d/wAKt5MzXGS66NrBDhT9\nGvarGdbbjCc04zNb3UuNuDCknBHH4cxwPMVBbUMTmz4H5GiCv7LGudGyLbIb\n37tAiqkQXwPX+zGS0T2pKQQO44qgrYAJbZ/vA4+FYW9+2yvvY3DJb9n0X2b/\nAHeaY/DI37aaVY2bfd5pj8Mj/tppVTKZXO1nZReNZarcuMN62iKpltvckqVv\nbyc8cBJHbUI/l2vh5uWH4r+iiiKgOZFa1Ptp9pYFAAvno6Xw/wBrYB5r+isK\n6Od8PJ7T/wDn+iidVPjJ5upFa1XWEkEqkIAHeaABhV0b7+cn0jT481/RVQX6\n3x9Napm2y4W9iWuBIUy8G/VQ4U88HAIHlmjan7SdKQypK7s24ocCGULc/UDF\nBptHuEe5671BOiq3o8ia442VDBKSeHA8qQmME+5W95SPQrOmGBnO4srKvjyx\nVmaJ2IXvV+l4F9t0i0MxZYUptt8rC07qinjhBHNJ7aqXeG9x3fKi02HbR9LW\nnZpY7Xc7mmLLYQ4FpW0vdBLqyPWAxyIpgiCJ6OerEexcrInjng66OPf7FeXO\njlq1a983KyFZ5qLzuT57lE9bL/arpFTIt1wjSWFcAtpYUM93vrsEtg8nUHzp\n2xgrN9HbWbSipq72hCsYymQ8DjyRWgdGvVSTlM6whQ5EOOf/ADosw82eS0/G\nvQWk/wBYUgfkadI2x6zaUs1slKbW/DhtR3FNklJUlABIzxxwpU77w7xSoA5J\nDTqgd2mSdbZzgPV/OpNSoArefp29O56v5moxetCajuEKRG31oQ8goKkLIIBG\nOBq76VKh2CC5sX2kwlbsGe3KYHshxeFAf9wPzrhlbDNYvhx56A6uUs7yt15r\nBPwozKVFCAgOwjXa1bqbKRx9pUlvj5U7RdiG0bqgyyiNDZHAbzyVEfBNGTSp\nhQO2z7ZNqjTLMkv3Bch+SpKnDvHdGBwwPPnVhwdN3xvHWOk+dWNSpUOyLw7T\nPbx1i8+dO7ER5A9ZQpxpUxGkNKA5ilW6lQB//9k=\n",
    "filename": "some file"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 36,
  "created_at": "2026-07-25T01:00:26.288+03:00",
  "absolute_url": "http://minio.insales.io:9000/insales-test-sber/files/1/36/36/original/some_file.?1784930426"
}
```

### Create file from url

`POST /admin/files.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `file[src]` | да | file url |

**Запрос**

```http
POST /admin/files.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "file": {
    "src": "https://abc.myinsales.ru/images/no_image_compact.jpg"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20,
  "created_at": "2026-07-25T01:00:25.986+03:00",
  "absolute_url": "http://minio.insales.io:9000/insales-test-sber/files/1/20/20/original/no_image_compact.jpg?1784930425"
}
```

### Destroy file

`DELETE /admin/files/:account_file_id.json`

**Запрос**

```http
DELETE /admin/files/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get file

`GET /admin/files/:account_file_id.json`

**Запрос**

```http
GET /admin/files/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "created_at": "2009-03-18T07:53:15.000+03:00",
  "absolute_url": "http://minio.insales.io:9000/insales-test-sber/files/1/original/thumb_prayerbeads.jpg"
}
```

### Get files

`GET /admin/files.json`

**Запрос**

```http
GET /admin/files.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "created_at": "2009-03-18T07:53:15.000+03:00",
    "absolute_url": "http://minio.insales.io:9000/insales-test-sber/files/1/original/thumb_prayerbeads.jpg"
  }
]
```

---

## Image

| Операция | Метод и путь |
|----------|--------------|
| [Create image from attachment](#create-image-from-attachment) | `POST /admin/products/:product_id/images.json` |
| [Create image from src](#create-image-from-src) | `POST /admin/products/:product_id/images.json` |
| [Destroy image](#destroy-image) | `DELETE /admin/products/:product_id/images/:image_id.json` |
| [Get image](#get-image) | `GET /admin/products/:product_id/images/:image_id.json` |
| [Get images](#get-images) | `GET /admin/products/:product_id/images.json` |
| [Update image](#update-image) | `PUT /admin/products/:product_id/images/:image_id.json` |

### Create image from attachment

`POST /admin/products/:product_id/images.json`

After uploading your picture you may find ‘loading.gif’ instead of it. That’s because picture uploading and creating of thumbnails are asynchronous. When it’s uploaded and thumbnails’re created its attribute updated_at will be overwritten with new timestamp.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `image[attachment]` | да | image file |
| `image[filename]` | да | file name |
| `image[title]` |  | title |
| `image[position]` |  | position |

**Запрос**

```http
POST /admin/products/1/images.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "image": {
    "attachment": "/9j/4AAQSkZJRgABAQEASABIAAD//gA+Q1JFQVRPUjogZ2QtanBlZyB2MS4w\nICh1c2luZyBJSkcgSlBFRyB2NjIpLCBkZWZhdWx0IHF1YWxpdHkK/9sAQwAG\nBAUGBQQGBgUGBwcGCAoQCgoJCQoUDg8MEBcUGBgXFBYWGh0lHxobIxwWFiAs\nICMmJykqKRkfLTAtKDAlKCko/9sAQwEHBwcKCAoTCgoTKBoWGigoKCgoKCgo\nKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgo/8AA\nEQgAWQBkAwEiAAIRAQMRAf/EAB0AAAEEAwEBAAAAAAAAAAAAAAgABQYHAQME\nAgn/xABEEAABAwMBBAYHBgIHCQAAAAABAgMEAAURBgcSITEIE0FRYYEUIjJx\nkbHSIzd1lKGzFRg1QkNSYqLwJSYzY3KCksLT/8QAGQEAAgMBAAAAAAAAAAAA\nAAAAAAEDBAUC/8QAIxEAAgICAgIBBQAAAAAAAAAAAAECEQMEEjEhQVEFEyIz\ngf/aAAwDAQACEQMRAD8AKmlVT7Q9oF4sGqV223NxFthpC0hxoqUSQSeO8O6o\ncvbbd23ktO/w5LiuASWFDP8AmrlzivDY+LqwiaVDwNtl5KlpDUTeQMqHobvA\nfGtLm3G9IRvlFvCOPFUZwf8AtS+7D5EEbSoYv5g7qpwIaNqUonAHo7nP/wAq\n4pHSQujDq2nFWhLiCUqSqO5kEdntV1YBVUqEtfSXvGfVVZvy7n1VpV0mL7n1\nTZfyzn1UwC6pUIKukzqDs/gv5Zz661npM6j7DZfyrn10AGFSoOVdJnU39X+C\n/lXPrr2OkjqlxIDD1jW6R/wxDdz8SrFABh0qZ9IXB67aTs1xlhAkS4bL7m4M\nJ3lIBOB2DJpUAUNt3usew67NzmRm5LLbbSeqcGQSUKA4VV8rag2JoQxY7e2w\nocFdWARU86UPUpvxVKSpbG/F6xCTglPrZwe/FVDb1Q4Vyc9HjIQyyvAku4cc\nHHhugjdB8s+NUtnHjk7mrJY5JKPFPwXBZro07aWp98YMVCsKaQ4d31f76u4d\n3fXq/SrdfYhTYoMVx5KSClZwV45lIJ9Y+GKi1usV01pdWFMOmW2jGX3DlpA7\nCByJ9/Luq2bZpyBpiE21bCHZ6R68pwZDiu0b3+h4VhuEMfn38fBS2c1A26g0\n9qELUpNtdbaPFCkxktocHhkAg+BqIXiJJknfeZWiaBhSVJwXMdo7zRoL1XLa\nixw8A2p7rG+rdxlLiBk8+Y4jj41HZ1tsGso5D8CEuafW3mMNqV78cD7+Yra1\nc7nSaKq2+PaAxP61gc6tDa7oFenrupyGHHITw6xpxQ9bHahf+NJ4E9owarAg\npUQoYIrQaa7L+PIskeSMGsVk1ikdirts/wDSDfuPyNcVd1m/pFr3K+RoA+ke\nzj7vtM/hsb9tNKls4+77TP4ZG/bTSoAHXpTuFGonhkAbkfn7lVS1tdRMkIRI\nO6CrOeQzjl8R+tXJ0qsK1M8kkZ6uPgHtOFcKHp9L7LXpLCypkEJWlQzuE8s+\nBqOcFIGX7p2/Lstl/h0Wa20pWSpDDZPDwXgDPPlXbaZl1E1wpnOSG1J+xSXd\n5RPd61UDA1M5AblJVGYeL7WG1nILTg4BYAOMgZ+NXHscntahtsx6Y71a22/R\n9xtPsrUCAodw5ms2enxtmfsY3FOTHDVe0G4yFqhz7UhEtlCmGkqGVNk8CEnx\n8fGvOkhchFEiZPhxIyXApS1kNgAcgnHEnPj8a49QliU6lltaRLDhQCsZ3TnB\nxnxzjspnf07cusT10z0pCVb286d3cAPYOVa+th4pUis+LjXRONfzEaiaSxAl\nFXo32eFew6o4OFA9vcrvobr7C6mW6hbZQQo44csH/XCr2hT4iI4bQ2SFulMh\nwD1lgJGefLPLyrfd9A2fVbyZlnddW8GwpUb2G5I5ApdPAODGCOGR7q1J6jyQ\nuL8ljTuH4+gZyhQVu4416cZcax1ja0Z5bySM0Tlk0JFtqt17diKTzjxkBCz/\nANTisk+VP0pVsYgLju2mJMiKG643JCnQoeJJ/UcakX0ebjfLyauPG8nQINd1\nm4XFr3K+Rp013ZEWPUsyPGSRCU4Vx+OfszxAz4A48qa7QP8AaDXuPyNZE4OE\nnGXaOJRcW0z6R7OPu+0z+GRv200qxs3+73TP4ZG/bTSrkQMnS6WU6kfx2Nxj\n+iqoazXRCHermgLaWnq3N44C0nsPce0HvAq9el2f955A/wCVG+SqHAAUqsTV\nkkn6fUw6gAuuQnfWZkIRkY7j2Z7+NPemrmzoqUtReffYloCHWwNzkchQwSeH\nf41Cos6TFQUR5TzSFc0oWQD5V0tyWpI3JjqgvscPH40JemRyg5KpdFkS7/Gv\nOun7rBJXbmGUuvKVkJBCcqwO/IPGsW3XLc5b5fe6tTn9mTugDuHuFQiTcYkG\nxuW+3uF56Qft3sYSEZB3R3kkDJ7hTAOfCp1kcOiJa8ZLz/C/bJMt7cZW4hS2\nXVgLQlXnnOedTS33q0Wq2b8PcMxSjiO2AXFZ7cJz29tDZbwZMF9pXEhPWDvO\nOf6Zp30/f3LddIRjoyktqQnJxneBGM+PCtjX3opRjJUmSY8Siy29Rank3+0u\nzmS4xIiK3XUJO6efAntyKY7BPdvS3WH3kNvupIbcUopSF9mTngCe3srk01MQ\n9MTujMKYkpIxzJ4EK/xA8KYWUu225vRST1gWUj3d/wAKt5MzXGS66NrBDhT9\nGvarGdbbjCc04zNb3UuNuDCknBHH4cxwPMVBbUMTmz4H5GiCv7LGudGyLbIb\n37tAiqkQXwPX+zGS0T2pKQQO44qgrYAJbZ/vA4+FYW9+2yvvY3DJb9n0X2b/\nAHeaY/DI37aaVY2bfd5pj8Mj/tppVTKZXO1nZReNZarcuMN62iKpltvckqVv\nbyc8cBJHbUI/l2vh5uWH4r+iiiKgOZFa1Ptp9pYFAAvno6Xw/wBrYB5r+isK\n6Od8PJ7T/wDn+iidVPjJ5upFa1XWEkEqkIAHeaABhV0b7+cn0jT481/RVQX6\n3x9Napm2y4W9iWuBIUy8G/VQ4U88HAIHlmjan7SdKQypK7s24ocCGULc/UDF\nBptHuEe5671BOiq3o8ia442VDBKSeHA8qQmME+5W95SPQrOmGBnO4srKvjyx\nVmaJ2IXvV+l4F9t0i0MxZYUptt8rC07qinjhBHNJ7aqXeG9x3fKi02HbR9LW\nnZpY7Xc7mmLLYQ4FpW0vdBLqyPWAxyIpgiCJ6OerEexcrInjng66OPf7FeXO\njlq1a983KyFZ5qLzuT57lE9bL/arpFTIt1wjSWFcAtpYUM93vrsEtg8nUHzp\n2xgrN9HbWbSipq72hCsYymQ8DjyRWgdGvVSTlM6whQ5EOOf/ADosw82eS0/G\nvQWk/wBYUgfkadI2x6zaUs1slKbW/DhtR3FNklJUlABIzxxwpU77w7xSoA5J\nDTqgd2mSdbZzgPV/OpNSoArefp29O56v5moxetCajuEKRG31oQ8goKkLIIBG\nOBq76VKh2CC5sX2kwlbsGe3KYHshxeFAf9wPzrhlbDNYvhx56A6uUs7yt15r\nBPwozKVFCAgOwjXa1bqbKRx9pUlvj5U7RdiG0bqgyyiNDZHAbzyVEfBNGTSp\nhQO2z7ZNqjTLMkv3Bch+SpKnDvHdGBwwPPnVhwdN3xvHWOk+dWNSpUOyLw7T\nPbx1i8+dO7ER5A9ZQpxpUxGkNKA5ilW6lQB//9k=\n",
    "filename": "image.jpg"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 51,
  "product_id": 1,
  "external_id": null,
  "position": 1,
  "created_at": "2026-07-25T00:59:42.000+03:00",
  "image_processing": true,
  "title": null,
  "url": "/images/loading.gif",
  "original_url": "/images/loading.gif",
  "medium_url": "/images/loading.gif",
  "small_url": "/images/loading.gif",
  "thumb_url": "/images/loading.gif",
  "compact_url": "/images/loading.gif",
  "large_url": "/images/loading.gif",
  "filename": "image.jpg",
  "has_bitmap_image": false,
  "image_content_type": "image/jpeg"
}
```

### Create image from src

`POST /admin/products/:product_id/images.json`

After uploading your picture you may find ‘loading.gif’ instead of it. That’s because picture uploading and creating of thumbnails are asynchronous. When it’s uploaded and thumbnails’re created its attribute updated_at will be overwritten with a new timestamp.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `image[src]` | да | image url |
| `image[filename]` |  | file name |
| `image[title]` |  | title |
| `image[position]` |  | position |
| `image[external_id]` |  | ID in external system. It helps to identify image after recreation and so on |

**Запрос**

```http
POST /admin/products/1/images.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "image": {
    "src": "https://abc.myinsales.ru/images/no_image_compact.jpg"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 35,
  "product_id": 1,
  "external_id": null,
  "position": 1,
  "created_at": "2026-07-25T00:59:42.000+03:00",
  "image_processing": true,
  "title": null,
  "url": "/images/loading.gif",
  "original_url": "/images/loading.gif",
  "medium_url": "/images/loading.gif",
  "small_url": "/images/loading.gif",
  "thumb_url": "/images/loading.gif",
  "compact_url": "/images/loading.gif",
  "large_url": "/images/loading.gif",
  "filename": "no_image_compact.jpg",
  "has_bitmap_image": false,
  "image_content_type": "image/jpeg"
}
```

### Destroy image

`DELETE /admin/products/:product_id/images/:image_id.json`

**Запрос**

```http
DELETE /admin/products/1/images/83.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get image

`GET /admin/products/:product_id/images/:image_id.json`

**Запрос**

```http
GET /admin/products/1/images/19.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 19,
  "product_id": 1,
  "external_id": null,
  "position": 1,
  "created_at": "2026-07-25T00:59:41.000+03:00",
  "image_processing": true,
  "title": null,
  "url": "/images/loading.gif",
  "original_url": "/images/loading.gif",
  "medium_url": "/images/loading.gif",
  "small_url": "/images/loading.gif",
  "thumb_url": "/images/loading.gif",
  "compact_url": "/images/loading.gif",
  "large_url": "/images/loading.gif",
  "filename": "image.jpg",
  "has_bitmap_image": false,
  "image_content_type": "image/jpeg"
}
```

### Get images

`GET /admin/products/:product_id/images.json`

**Запрос**

```http
GET /admin/products/1/images.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[]
```

### Update image

`PUT /admin/products/:product_id/images/:image_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `image[position]` |  | position |
| `image[title]` |  | title |

**Запрос**

```http
PUT /admin/products/1/images/67.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "image": {
    "position": 5,
    "title": "myimage"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 67,
  "product_id": 1,
  "external_id": null,
  "position": 5,
  "created_at": "2026-07-25T00:59:42.000+03:00",
  "image_processing": true,
  "title": "myimage",
  "url": "/images/loading.gif",
  "original_url": "/images/loading.gif",
  "medium_url": "/images/loading.gif",
  "small_url": "/images/loading.gif",
  "thumb_url": "/images/loading.gif",
  "compact_url": "/images/loading.gif",
  "large_url": "/images/loading.gif",
  "filename": "image.jpg",
  "has_bitmap_image": false,
  "image_content_type": "image/jpeg"
}
```

---

## JsTag

| Операция | Метод и путь |
|----------|--------------|
| [Create js_tag from code](#create-js-tag-from-code) | `POST /admin/js_tags.json` |
| [Create js_tag from url](#create-js-tag-from-url) | `POST /admin/js_tags.json` |
| [Destroy js_tag](#destroy-js-tag) | `DELETE /admin/js_tags/:js_tag_id.json` |
| [Get js_tag](#get-js-tag) | `GET /admin/js_tags/:js_tag_id.json` |
| [Get js_tags](#get-js-tags) | `GET /admin/js_tags.json` |
| [Update js_tag](#update-js-tag) | `PUT /admin/js_tags/:js_tag_id.json` |

### Create js_tag from code

`POST /admin/js_tags.json`

Please be careful when plugging in jquery: many default themes use old version of it. Including of a new version of JQuery may result in malfunction of your store.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `js_tag[type]` | да | JsTag::TextTag |
| `js_tag[content]` | да | javascript code without tag SCRIPT |
| `js_tag[name]` |  | name for js tag |

**Запрос**

```http
POST /admin/js_tags.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "js_tag": {
    "type": "JsTag::TextTag",
    "content": "alert('test');",
    "name": "alert test"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "content": "alert('test');",
  "id": 57,
  "name": "alert test",
  "created_at": "2026-07-25T01:00:27.000+03:00",
  "updated_at": "2026-07-25T01:00:27.000+03:00",
  "type": "JsTag::TextTag"
}
```

### Create js_tag from url

`POST /admin/js_tags.json`

Please be careful while including JQuery: many default themes use an old version of it. Including of a new version of JQuery may result in malfunction of your store.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `js_tag[type]` | да | JsTag::FileTag |
| `js_tag[content]` | да | url of javascript file |
| `js_tag[name]` |  | name for js tag |

**Запрос**

```http
POST /admin/js_tags.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "js_tag": {
    "type": "JsTag::FileTag",
    "content": "http://code.jquery.com/jquery-1.9.1.min.js",
    "name": "jquery-1.9.1.min"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "content": "http://code.jquery.com/jquery-1.9.1.min.js",
  "id": 89,
  "name": "jquery-1.9.1.min",
  "created_at": "2026-07-25T01:00:27.000+03:00",
  "updated_at": "2026-07-25T01:00:27.000+03:00",
  "type": "JsTag::FileTag"
}
```

### Destroy js_tag

`DELETE /admin/js_tags/:js_tag_id.json`

**Запрос**

```http
DELETE /admin/js_tags/121.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get js_tag

`GET /admin/js_tags/:js_tag_id.json`

**Запрос**

```http
GET /admin/js_tags/25.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "content": "http://code.jquery.com/jquery-1.9.1.min.js",
  "id": 25,
  "name": "jquery1.9.1",
  "created_at": "2026-07-25T01:00:26.000+03:00",
  "updated_at": "2026-07-25T01:00:26.000+03:00",
  "type": "JsTag::FileTag"
}
```

### Get js_tags

`GET /admin/js_tags.json`

**Запрос**

```http
GET /admin/js_tags.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "content": "http://code.jquery.com/jquery-1.9.1.min.js",
    "id": 9,
    "name": "jquery1.9.1",
    "created_at": "2026-07-25T01:00:26.000+03:00",
    "updated_at": "2026-07-25T01:00:26.000+03:00",
    "type": "JsTag::FileTag"
  }
]
```

### Update js_tag

`PUT /admin/js_tags/:js_tag_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `js_tag[content]` |  | url of javascript file or javascript code without tag SCRIPT |
| `js_tag[name]` |  | name for js tag |

**Запрос**

```http
PUT /admin/js_tags/105.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "js_tag": {
    "content": "http://code.jquery.com/jquery-2.2.4.min.js",
    "name": "jquery-2.2.4"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "content": "http://code.jquery.com/jquery-2.2.4.min.js",
  "id": 105,
  "name": "jquery-2.2.4",
  "created_at": "2026-07-25T01:00:27.000+03:00",
  "updated_at": "2026-07-25T01:00:27.000+03:00",
  "type": "JsTag::FileTag"
}
```

---

## Marketplace

| Операция | Метод и путь |
|----------|--------------|
| [Create Google Merchant marketplace](#create-google-merchant-marketplace) | `POST /admin/marketplaces.json` |
| [Create Liquid for Yandex marketplace](#create-liquid-for-yandex-marketplace) | `POST /admin/marketplaces.json` |
| [Create Model for Yandex marketplace](#create-model-for-yandex-marketplace) | `POST /admin/marketplaces.json` |
| [Destroy marketplace](#destroy-marketplace) | `DELETE /admin/marketplaces/:ym_marketplace_id.json` |
| [Get marketplace](#get-marketplace) | `GET /admin/marketplaces/:ym_marketplace_id.json` |
| [Get marketplaces](#get-marketplaces) | `GET /admin/marketplaces.json` |
| [Update Google Merchant marketplace](#update-google-merchant-marketplace) | `PUT /admin/marketplaces/:gm_marketplace_id.json` |
| [Update Liquid for Yandex marketplace](#update-liquid-for-yandex-marketplace) | `PUT /admin/marketplaces/:yl_marketplace_id.json` |
| [Update Model for Yandex marketplace](#update-model-for-yandex-marketplace) | `PUT /admin/marketplaces/:ym_marketplace_id.json` |

### Create Google Merchant marketplace

`POST /admin/marketplaces.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `marketplace[name]` | да | Name |
| `marketplace[type]` | да | Marketplace::GoogleMerchant |
| `marketplace[description_type]` |  | 1 fo short description, 2 for full |
| `marketplace[local_delivery_cost]` |  | Local delivery cost |
| `marketplace[adult]` |  | 1 if products are adult-only, else 0 |
| `marketplace[utm_label]` |  | UTM label |

**Запрос**

```http
POST /admin/marketplaces.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "marketplace": {
    "name": "GM",
    "type": "Marketplace::GoogleMerchant",
    "description_type": 1,
    "adult": 0
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 73,
  "name": "GM",
  "type": "Marketplace::GoogleMerchant",
  "created_at": "2026-07-25T01:00:19.413+03:00",
  "updated_at": "2026-07-25T01:00:19.413+03:00",
  "url": "https://quentin.insales.io:3002/marketplace/73.xml",
  "collection_ids": [],
  "description_type": 1,
  "local_delivery_cost": null,
  "utm_label": null,
  "adult": false,
  "delivery_new_style": false,
  "use_stock_currency": false,
  "use_variants": false
}
```

### Create Liquid for Yandex marketplace

`POST /admin/marketplaces.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `marketplace[name]` | да | Name |
| `marketplace[type]` | да | Marketplace type, e.g. Marketplace::LiquidYandexMarket |
| `marketplace[shop_name]` |  | Shop name |
| `marketplace[use_variants]` |  | Create document for all variants, not aggregated product |
| `marketplace[shop_company]` |  | Company name |
| `marketplace[page_encoding]` | да | Encoding, windows-1251 or utf-8 |
| `marketplace[offer_liquid]` | да | Liquid template |
| `marketplace[collection_ids]` |  | Collections to load |

**Запрос**

```http
POST /admin/marketplaces.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "marketplace": {
    "name": "LYM",
    "type": "Marketplace::LiquidYandexMarket",
    "shop_name": "eBay",
    "use_variants": 1,
    "shop_company": "eBay",
    "page_encoding": "utf-8",
    "offer_liquid": "\n<offer type=\"vendor.model\" available=\"{% if product.available %}true{% else %}false{% endif %}\" id=\"{{product.id}}\" >\n  <url>{{account.url}}/product_by_id/{{product.id}}</url>\n  <price>{{product.price}}</price>\n  <currencyId>RUR</currencyId>\n  <categoryId>{{collection.id}}</categoryId>\n  {%for image in product.images limit: 10 %}\n  <picture>{{image.url}}</picture>\n  {% endfor %}\n  <store>false</store>\n  <pickup>true</pickup>\n  <delivery>true</delivery>\n  <vendor>{{product.properties.vendor}}</vendor>\n  <model>{{product.title}}</model>\n  <description>{{product.description}}</description>\n</offer>\n        ",
    "collection_ids": [
      6,
      1,
      2,
      3,
      4,
      5,
      8
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 121,
  "name": "LYM",
  "type": "Marketplace::LiquidYandexMarket",
  "created_at": "2026-07-25T01:00:20.422+03:00",
  "updated_at": "2026-07-25T01:00:20.422+03:00",
  "url": "https://quentin.insales.io:3002/marketplace/121.xml",
  "collection_ids": [],
  "local_delivery_cost": null,
  "vendor_id": null,
  "offer_liquid": "\n<offer type=\"vendor.model\" available=\"{% if product.available %}true{% else %}false{% endif %}\" id=\"{{product.id}}\" >\n  <url>{{account.url}}/product_by_id/{{product.id}}</url>\n  <price>{{product.price}}</price>\n  <currencyId>RUR</currencyId>\n  <categoryId>{{collection.id}}</categoryId>\n  {%for image in product.images limit: 10 %}\n  <picture>{{image.url}}</picture>\n  {% endfor %}\n  <store>false</store>\n  <pickup>true</pickup>\n  <delivery>true</delivery>\n  <vendor>{{product.properties.vendor}}</vendor>\n  <model>{{product.title}}</model>\n  <description>{{product.description}}</description>\n</offer>\n        ",
  "page_encoding": "utf-8",
  "shop_company": "eBay",
  "shop_name": "eBay",
  "delivery_new_style": false,
  "use_stock_currency": false,
  "use_variants": true
}
```

### Create Model for Yandex marketplace

`POST /admin/marketplaces.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `marketplace[name]` | да | Name |
| `marketplace[type]` | да | Marketplace type, e.g. Marketplace::ModelYandexMarket |
| `marketplace[shop_name]` |  | Shop name |
| `marketplace[shop_company]` |  | Company name |
| `marketplace[description_type]` |  | 1 fo short description, 2 for full |
| `marketplace[vendor_id]` | да | id of corresponding property |
| `marketplace[use_sku_as_vendor_code]` |  | Use sku as vendor code |
| `marketplace[pickup]` |  | Pickup enable |
| `marketplace[store]` |  | Offline store exists |
| `marketplace[delivery_style]` |  | Delivery pricing type, 0 for old, 1 for new |
| `marketplace[local_delivery_cost]` |  | Local delivery cost |
| `marketplace[adult]` |  | 1 if products are adult-only, else 0 |
| `marketplace[utm_label]` |  | UTM label |
| `marketplace[page_encoding]` | да | Encoding, windows-1251 or utf-8 |
| `marketplace[image_style]` | да | Image style: thumb, compact, medium, large or original |
| `marketplace[model_type]` | да | Model naming type: name or sku |
| `marketplace[collection_ids]` |  | Collections to load |

**Запрос**

```http
POST /admin/marketplaces.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "marketplace": {
    "name": "YM",
    "type": "Marketplace::ModelYandexMarket",
    "shop_name": "eBay",
    "shop_company": "eBay",
    "description_type": 1,
    "vendor_id": 6,
    "adult": 0,
    "page_encoding": "windows-1251",
    "image_style": "oroginal",
    "collection_ids": [
      6,
      1,
      2,
      3,
      4,
      5,
      8
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 57,
  "name": "YM",
  "type": "Marketplace::ModelYandexMarket",
  "created_at": "2026-07-25T01:00:19.142+03:00",
  "updated_at": "2026-07-25T01:00:19.142+03:00",
  "url": "https://quentin.insales.io:3002/marketplace/57.xml",
  "collection_ids": [],
  "delivery_price_kind": null,
  "description_type": 1,
  "local_delivery_cost": null,
  "sales_notes": null,
  "utm_label": null,
  "vendor_id": 6,
  "page_encoding": "windows-1251",
  "shop_company": "eBay",
  "shop_name": "eBay",
  "adult": false,
  "delivery": true,
  "delivery_new_style": false,
  "pickup": false,
  "store": false,
  "use_sku_as_vendor_code": false,
  "use_stock_currency": false,
  "use_variants": false
}
```

### Destroy marketplace

`DELETE /admin/marketplaces/:ym_marketplace_id.json`

**Запрос**

```http
DELETE /admin/marketplaces/153.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get marketplace

`GET /admin/marketplaces/:ym_marketplace_id.json`

**Запрос**

```http
GET /admin/marketplaces/41.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 41,
  "name": "Yandex Model",
  "type": "Marketplace::ModelYandexMarket",
  "created_at": "2026-07-25T01:00:18.713+03:00",
  "updated_at": "2026-07-25T01:00:18.713+03:00",
  "url": "https://quentin.insales.io:3002/marketplace/41.xml",
  "collection_ids": [],
  "delivery_price_kind": null,
  "description_type": 1,
  "local_delivery_cost": null,
  "sales_notes": null,
  "utm_label": null,
  "vendor_id": 6,
  "page_encoding": "windows-1251",
  "shop_company": "eBay",
  "shop_name": "eBay",
  "adult": false,
  "delivery": false,
  "delivery_new_style": false,
  "pickup": false,
  "store": false,
  "use_sku_as_vendor_code": false,
  "use_stock_currency": false,
  "use_variants": false
}
```

### Get marketplaces

`GET /admin/marketplaces.json`

**Запрос**

```http
GET /admin/marketplaces.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 17,
    "name": "Yandex Model",
    "type": "Marketplace::LiquidYandexMarket",
    "created_at": "2026-07-25T01:00:18.400+03:00",
    "updated_at": "2026-07-25T01:00:18.400+03:00",
    "url": "https://quentin.insales.io:3002/marketplace/17.xml",
    "collection_ids": [],
    "local_delivery_cost": null,
    "vendor_id": null,
    "offer_liquid": "<offer>...</offer>",
    "page_encoding": "utf-8",
    "shop_company": "eBay",
    "shop_name": "eBay",
    "delivery_new_style": false,
    "use_stock_currency": false,
    "use_variants": false
  },
  {
    "id": 25,
    "name": "Google Model",
    "type": "Marketplace::GoogleMerchant",
    "created_at": "2026-07-25T01:00:18.405+03:00",
    "updated_at": "2026-07-25T01:00:18.405+03:00",
    "url": "https://quentin.insales.io:3002/marketplace/25.xml",
    "collection_ids": [],
    "description_type": null,
    "local_delivery_cost": null,
    "utm_label": null,
    "adult": false,
    "delivery_new_style": false,
    "use_stock_currency": false,
    "use_variants": false
  }
]
```

### Update Google Merchant marketplace

`PUT /admin/marketplaces/:gm_marketplace_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `marketplace[name]` |  | Name |
| `marketplace[description_type]` |  | 1 fo short description, 2 for full |
| `marketplace[local_delivery_cost]` |  | Local delivery cost |
| `marketplace[adult]` |  | 1 if products are adult-only, else 0 |
| `marketplace[utm_label]` |  | UTM label |

**Запрос**

```http
PUT /admin/marketplaces/105.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "marketplace": {
    "name": "New GM",
    "description_type": 2,
    "adult": 0
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 105,
  "name": "New GM",
  "type": "Marketplace::GoogleMerchant",
  "created_at": "2026-07-25T01:00:20.059+03:00",
  "updated_at": "2026-07-25T01:00:20.142+03:00",
  "url": "https://quentin.insales.io:3002/marketplace/105.xml",
  "collection_ids": [],
  "description_type": 2,
  "local_delivery_cost": null,
  "utm_label": null,
  "adult": false,
  "delivery_new_style": false,
  "use_stock_currency": false,
  "use_variants": false
}
```

### Update Liquid for Yandex marketplace

`PUT /admin/marketplaces/:yl_marketplace_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `marketplace[name]` |  | Name |
| `marketplace[shop_name]` |  | Shop name |
| `marketplace[use_variants]` |  | Create document for all variants, not aggregated product |
| `marketplace[shop_company]` |  | Company name |
| `marketplace[page_encoding]` |  | Encoding, windows-1251 or utf-8 |
| `marketplace[offer_liquid]` |  | Liquid template |
| `marketplace[collection_ids]` |  | Collections to load |

**Запрос**

```http
PUT /admin/marketplaces/137.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "marketplace": {
    "name": "New LYM"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 137,
  "name": "New LYM",
  "type": "Marketplace::LiquidYandexMarket",
  "created_at": "2026-07-25T01:00:20.640+03:00",
  "updated_at": "2026-07-25T01:00:20.725+03:00",
  "url": "https://quentin.insales.io:3002/marketplace/137.xml",
  "collection_ids": [],
  "local_delivery_cost": null,
  "vendor_id": null,
  "offer_liquid": "<offer>...</offer>",
  "page_encoding": "utf-8",
  "shop_company": "eBay",
  "shop_name": "eBay",
  "delivery_new_style": false,
  "use_stock_currency": false,
  "use_variants": false
}
```

### Update Model for Yandex marketplace

`PUT /admin/marketplaces/:ym_marketplace_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `marketplace[name]` |  | Name |
| `marketplace[shop_name]` |  | Shop name |
| `marketplace[shop_company]` |  | Company name |
| `marketplace[description_type]` |  | 1 fo short description, 2 for full |
| `marketplace[vendor_id]` |  | id of corresponding property |
| `marketplace[use_sku_as_vendor_code]` |  | Use sku as vendor code |
| `marketplace[pickup]` |  | Pickup enable |
| `marketplace[store]` |  | Offline store exists |
| `marketplace[delivery_style]` |  | Delivery pricing type, 0 for old, 1 for new |
| `marketplace[local_delivery_cost]` |  | Local delivery cost |
| `marketplace[adult]` |  | 1 if products are adult-only, else 0 |
| `marketplace[utm_label]` |  | UTM label |
| `marketplace[page_encoding]` |  | Encoding, windows-1251 or utf-8 |
| `marketplace[image_style]` |  | Image style: thumb, compact, medium, large or original |
| `marketplace[model_type]` |  | Model naming type: name or sku |
| `marketplace[collection_ids]` |  | Collections to load |

**Запрос**

```http
PUT /admin/marketplaces/89.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "marketplace": {
    "name": "New YM",
    "vendor_id": 4,
    "collection_ids": [
      6,
      1,
      2,
      3,
      4,
      5,
      8
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 89,
  "name": "New YM",
  "type": "Marketplace::ModelYandexMarket",
  "created_at": "2026-07-25T01:00:19.721+03:00",
  "updated_at": "2026-07-25T01:00:19.813+03:00",
  "url": "https://quentin.insales.io:3002/marketplace/89.xml",
  "collection_ids": [
    6,
    1,
    2,
    3,
    4,
    5,
    8
  ],
  "delivery_price_kind": null,
  "description_type": 1,
  "local_delivery_cost": null,
  "sales_notes": null,
  "utm_label": null,
  "vendor_id": 4,
  "page_encoding": "windows-1251",
  "shop_company": "eBay",
  "shop_name": "eBay",
  "adult": false,
  "delivery": false,
  "delivery_new_style": false,
  "pickup": false,
  "store": false,
  "use_sku_as_vendor_code": false,
  "use_stock_currency": false,
  "use_variants": false
}
```

---

## Menu

| Операция | Метод и путь |
|----------|--------------|
| [Create menu](#create-menu) | `POST /admin/menus.json` |
| [Delete menu](#delete-menu) | `DELETE /admin/menus/:menu_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/menus/:menu_id.json` |
| [Get menus list](#get-menus-list) | `GET /admin/menus.json` |
| [Update menu](#update-menu) | `PUT /admin/menus/:menu_id.json` |

### Create menu

`POST /admin/menus.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu[name]` | да | menu text |
| `menu[handle]` | да | handle used in liquid |

**Запрос**

```http
POST /admin/menus.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "menu": {
    "name": "top menu",
    "handle": "top-menu"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 18,
  "name": "top menu",
  "handle": "top-menu",
  "used_in_theme": false
}
```

### Delete menu

`DELETE /admin/menus/:menu_id.json`

**Запрос**

```http
DELETE /admin/menus/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
Content-Type: application/x-www-form-urlencoded
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/menus/:menu_id.json`

Content-Type: application/json Get menu

**Запрос**

```http
GET /admin/menus/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "name": "Main menu",
  "handle": "main",
  "used_in_theme": false
}
```

### Get menus list

`GET /admin/menus.json`

**Запрос**

```http
GET /admin/menus.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "name": "Main menu",
    "handle": "main",
    "used_in_theme": false
  },
  {
    "id": 2,
    "name": "Secondary menu",
    "handle": "second",
    "used_in_theme": false
  }
]
```

### Update menu

`PUT /admin/menus/:menu_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu[name]` |  | menu text |
| `menu[handle]` |  | handle used in liquid |

**Запрос**

```http
PUT /admin/menus/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "menu": {
    "name": "Bottom menu"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "name": "Bottom menu",
  "handle": "main",
  "used_in_theme": false
}
```

---

## MenuItem

| Операция | Метод и путь |
|----------|--------------|
| [Add blog to menu](#add-blog-to-menu) | `POST /admin/menus/:menu_id/items.json` |
| [Add collection to menu](#add-collection-to-menu) | `POST /admin/menus/:menu_id/items.json` |
| [Add external link to menu](#add-external-link-to-menu) | `POST /admin/menus/:menu_id/items.json` |
| [Add page to menu](#add-page-to-menu) | `POST /admin/menus/:menu_id/items.json` |
| [Delete menu_item](#delete-menu-item) | `DELETE /admin/menus/:menu_id/items/:item_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/menus/:menu_id/items/:item_id.json` |
| [Get menu items list](#get-menu-items-list) | `GET /admin/menus/:menu_id/items.json` |
| [Update menu item](#update-menu-item) | `PUT /admin/menus/:menu_id/items/:item_id.json` |

### Add blog to menu

`POST /admin/menus/:menu_id/items.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu_item[title]` | да | menu item text |
| `menu_item[position]` |  | menu item position |
| `menu_item[related_object_id]` | да | referenc to blog: "Blog_#{blog_id}" |

**Запрос**

```http
POST /admin/menus/1/items.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "menu_item": {
    "title": "new item",
    "related_object_id": "Blog_1"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 60,
  "position": 5,
  "title": "new item",
  "url": "/blog",
  "related_object_id": "Blog_1"
}
```

### Add collection to menu

`POST /admin/menus/:menu_id/items.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu_item[title]` | да | menu item text |
| `menu_item[position]` |  | menu item position |
| `menu_item[related_object_id]` | да | referenc to collection: "Collection_#{collection_id}" |

**Запрос**

```http
POST /admin/menus/1/items.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "menu_item": {
    "title": "new item",
    "related_object_id": "Collection_1"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 52,
  "position": 5,
  "title": "new item",
  "url": "/collection/футболка",
  "related_object_id": "Collection_1"
}
```

### Add external link to menu

`POST /admin/menus/:menu_id/items.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu_item[title]` | да | menu item text |
| `menu_item[position]` |  | menu item position |
| `menu_item[url]` | да | menu item url |

**Запрос**

```http
POST /admin/menus/1/items.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "menu_item": {
    "title": "new item",
    "url": "http://my-another-shop.com/"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 68,
  "position": 5,
  "title": "new item",
  "url": "http://my-another-shop.com/",
  "related_object_id": "Page_0"
}
```

### Add page to menu

`POST /admin/menus/:menu_id/items.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu_item[title]` | да | menu item text |
| `menu_item[position]` |  | menu item position |
| `menu_item[related_object_id]` | да | referenc to page: "Page_#{page_id}" |

**Запрос**

```http
POST /admin/menus/1/items.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "menu_item": {
    "title": "new item",
    "related_object_id": "Page_1"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 44,
  "position": 5,
  "title": "new item",
  "url": "/",
  "related_object_id": "Page_1"
}
```

### Delete menu_item

`DELETE /admin/menus/:menu_id/items/:item_id.json`

**Запрос**

```http
DELETE /admin/menus/1/items/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
Content-Type: application/x-www-form-urlencoded
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/menus/:menu_id/items/:item_id.json`

Content-Type: application/json Get menu item

**Запрос**

```http
GET /admin/menus/1/items/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "position": 1,
  "title": "About",
  "url": "/about",
  "related_object_id": "Page_0"
}
```

### Get menu items list

`GET /admin/menus/:menu_id/items.json`

**Запрос**

```http
GET /admin/menus/1/items.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "position": 1,
    "title": "About",
    "url": "/about",
    "related_object_id": "Page_0"
  },
  {
    "id": 2,
    "position": 2,
    "title": "Доставка",
    "url": "/delivery",
    "related_object_id": "Page_0"
  },
  {
    "id": 3,
    "position": 3,
    "title": "Страница",
    "url": "/",
    "related_object_id": "Page_1"
  },
  {
    "id": 4,
    "position": 4,
    "title": "Категория",
    "url": "/collection/футболка",
    "related_object_id": "Collection_1"
  }
]
```

### Update menu item

`PUT /admin/menus/:menu_id/items/:item_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu_item[title]` |  | menu item text |
| `menu_item[url]` |  | menu item url |
| `menu_item[related_object_id]` |  | referenc to collection, page or blog |
| `menu_item[position]` |  | menu item position |

**Запрос**

```http
PUT /admin/menus/1/items/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "menu_item": {
    "title": "last menu item",
    "position": 1000
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "position": 1000,
  "title": "last menu item",
  "url": "/about",
  "related_object_id": "Page_0"
}
```

---

## OAuthApp

| Операция | Метод и путь |
|----------|--------------|
| [Create o_auth_app](#create-o-auth-app) | `POST /admin/o_auth_apps.json` |
| [Destroy o_auth_app](#destroy-o-auth-app) | `DELETE /admin/o_auth_apps/:o_auth_app_id.json` |
| [Get o_auth_app](#get-o-auth-app) | `GET /admin/o_auth_apps/:o_auth_app_id.json` |
| [Get o_auth_apps](#get-o-auth-apps) | `GET /admin/o_auth_apps.json` |
| [Update o_auth_app](#update-o-auth-app) | `PUT /admin/o_auth_apps/:o_auth_app_id.json` |

### Create o_auth_app

`POST /admin/o_auth_apps.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `o_auth_app[provider_handle]` | да | provider name |
| `o_auth_app[client_id]` | да | client id |
| `o_auth_app[client_secret]` | да | client_secret |
| `o_auth_app[issuer]` | да | issuer |

**Запрос**

```http
POST /admin/o_auth_apps.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "o_auth_app": {
    "provider_handle": "provider name",
    "client_id": "10332",
    "client_secret": "NdsaPd",
    "issuer": "/.well-known/openid-configuration"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "client_id": "10332",
  "client_secret": "NdsaPd",
  "provider": "provider name",
  "issuer": "/.well-known/openid-configuration",
  "id": 8,
  "enabled": true
}
```

### Destroy o_auth_app

`DELETE /admin/o_auth_apps/:o_auth_app_id.json`

**Запрос**

```http
DELETE /admin/o_auth_apps/12.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get o_auth_app

`GET /admin/o_auth_apps/:o_auth_app_id.json`

**Запрос**

```http
GET /admin/o_auth_apps/4.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "client_id": "953",
  "client_secret": "NdsaPd",
  "provider": "new_provider",
  "issuer": "/.well-known/openid-configuration",
  "id": 4,
  "enabled": true
}
```

### Get o_auth_apps

`GET /admin/o_auth_apps.json`

**Запрос**

```http
GET /admin/o_auth_apps.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "client_id": "953",
    "client_secret": "NdsaPd",
    "provider": "new_provider",
    "issuer": "/.well-known/openid-configuration",
    "id": 2,
    "enabled": true
  }
]
```

### Update o_auth_app

`PUT /admin/o_auth_apps/:o_auth_app_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `o_auth_app[provider_handle]` |  | provider name |
| `o_auth_app[client_id]` |  | client id |
| `o_auth_app[client_secret]` |  | client secret |
| `o_auth_app[issuer]` |  | issuer |

**Запрос**

```http
PUT /admin/o_auth_apps/10.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "o_auth_app": {
    "provider_handle": "new name provider"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "client_id": "953",
  "client_secret": "NdsaPd",
  "provider": "new name provider",
  "issuer": "/.well-known/openid-configuration",
  "id": 10,
  "enabled": true
}
```

---

## OptionName

| Операция | Метод и путь |
|----------|--------------|
| [Create option_name](#create-option-name) | `POST /admin/option_names.json` |
| [Destroy option_name](#destroy-option-name) | `DELETE /admin/option_names/:option_name_id.json` |
| [Get option_name](#get-option-name) | `GET /admin/option_names/:option_name_id.json` |
| [Get option_names](#get-option-names) | `GET /admin/option_names.json` |
| [Update option_name](#update-option-name) | `PUT /admin/option_names/:option_name_id.json` |

### Create option_name

`POST /admin/option_names.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `option_name[title]` | да | title |
| `option_name[position]` |  | position in option names list |
| `option_name[navigational]` |  | true by default; show in products filter on category page |

**Запрос**

```http
POST /admin/option_names.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "option_name": {
    "title": "Colour",
    "navigational": false
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2517850,
  "position": 3,
  "navigational": false,
  "permalink": "colour",
  "title": "Colour",
  "api_permalink": "o-colour"
}
```

### Destroy option_name

`DELETE /admin/option_names/:option_name_id.json`

**Запрос**

```http
DELETE /admin/option_names/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get option_name

`GET /admin/option_names/:option_name_id.json`

**Запрос**

```http
GET /admin/option_names/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "position": 1,
  "navigational": true,
  "permalink": null,
  "title": "Вариант",
  "api_permalink": null
}
```

### Get option_names

`GET /admin/option_names.json`

**Запрос**

```http
GET /admin/option_names.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "position": 1,
    "navigational": true,
    "permalink": null,
    "title": "Вариант",
    "api_permalink": null
  },
  {
    "id": 3,
    "position": 2,
    "navigational": true,
    "permalink": null,
    "title": "Цвет",
    "api_permalink": null
  }
]
```

### Update option_name

`PUT /admin/option_names/:option_name_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `option_name[title]` |  | title |
| `option_name[position]` |  | position |
| `option_name[navigational]` |  | true by default; show in products filter on category page |

**Запрос**

```http
PUT /admin/option_names/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "option_name": {
    "title": "Type",
    "position": 2,
    "navigational": false
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "position": 2,
  "navigational": false,
  "permalink": "type",
  "title": "Type",
  "api_permalink": "o-type"
}
```

---

## OptionValue

| Операция | Метод и путь |
|----------|--------------|
| [Create option_value](#create-option-value) | `POST /admin/option_names/:option_name_id/option_values.json` |
| [Destroy option_value image](#destroy-option-value-image) | `DELETE /admin/option_names/:option_name_id/option_values/:option_value_id/destroy_image.json` |
| [Destroy option_value](#destroy-option-value) | `DELETE /admin/option_names/:option_name_id/option_values/:option_value_id.json` |
| [Get destroy option_value status](#get-destroy-option-value-status) | `GET /admin/option_names/:option_name_id/option_values/:option_value_id/destroy_status.json` |
| [Get option_value](#get-option-value) | `GET /admin/option_names/:option_name_id/option_values/:option_value_id.json` |
| [Get option_values for all options](#get-option-values-for-all-options) | `GET /admin/option_values.json` |
| [Get option_values for specific option](#get-option-values-for-specific-option) | `GET /admin/option_names/:option_name_id/option_values.json` |
| [Update option_value](#update-option-value) | `PUT /admin/option_names/:option_name_id/option_values/:option_value_id.json` |

### Create option_value

`POST /admin/option_names/:option_name_id/option_values.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `option_value[title]` | да | title |
| `option_value[position]` |  | position in option_values list |
| `option_value[image][filename]` |  | image filename, required if attachment present |
| `option_value[image][attachment]` |  | base64 encoded image body |
| `option_value[image][src]` |  | image url |

**Запрос**

```http
POST /admin/option_names/1/option_values.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "option_value": {
    "title": "Blue",
    "image": {
      "filename": "image.jpg",
      "attachment": "/9j/4AAQSkZJRgABAQEASABIAAD//gA+Q1JFQVRPUjogZ2QtanBlZyB2MS4w\nICh1c2luZyBJSkcgSlBFRyB2NjIpLCBkZWZhdWx0IHF1YWxpdHkK/9sAQwAG\nBAUGBQQGBgUGBwcGCAoQCgoJCQoUDg8MEBcUGBgXFBYWGh0lHxobIxwWFiAs\nICMmJykqKRkfLTAtKDAlKCko/9sAQwEHBwcKCAoTCgoTKBoWGigoKCgoKCgo\nKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgo/8AA\nEQgAWQBkAwEiAAIRAQMRAf/EAB0AAAEEAwEBAAAAAAAAAAAAAAgABQYHAQME\nAgn/xABEEAABAwMBBAYHBgIHCQAAAAABAgMEAAURBgcSITEIE0FRYYEUIjJx\nkbHSIzd1lKGzFRg1QkNSYqLwJSYzY3KCksLT/8QAGQEAAgMBAAAAAAAAAAAA\nAAAAAAEDBAUC/8QAIxEAAgICAgIBBQAAAAAAAAAAAAECEQMEEjEhQVEFEyIz\ngf/aAAwDAQACEQMRAD8AKmlVT7Q9oF4sGqV223NxFthpC0hxoqUSQSeO8O6o\ncvbbd23ktO/w5LiuASWFDP8AmrlzivDY+LqwiaVDwNtl5KlpDUTeQMqHobvA\nfGtLm3G9IRvlFvCOPFUZwf8AtS+7D5EEbSoYv5g7qpwIaNqUonAHo7nP/wAq\n4pHSQujDq2nFWhLiCUqSqO5kEdntV1YBVUqEtfSXvGfVVZvy7n1VpV0mL7n1\nTZfyzn1UwC6pUIKukzqDs/gv5Zz661npM6j7DZfyrn10AGFSoOVdJnU39X+C\n/lXPrr2OkjqlxIDD1jW6R/wxDdz8SrFABh0qZ9IXB67aTs1xlhAkS4bL7m4M\nJ3lIBOB2DJpUAUNt3usew67NzmRm5LLbbSeqcGQSUKA4VV8rag2JoQxY7e2w\nocFdWARU86UPUpvxVKSpbG/F6xCTglPrZwe/FVDb1Q4Vyc9HjIQyyvAku4cc\nHHhugjdB8s+NUtnHjk7mrJY5JKPFPwXBZro07aWp98YMVCsKaQ4d31f76u4d\n3fXq/SrdfYhTYoMVx5KSClZwV45lIJ9Y+GKi1usV01pdWFMOmW2jGX3DlpA7\nCByJ9/Luq2bZpyBpiE21bCHZ6R68pwZDiu0b3+h4VhuEMfn38fBS2c1A26g0\n9qELUpNtdbaPFCkxktocHhkAg+BqIXiJJknfeZWiaBhSVJwXMdo7zRoL1XLa\nixw8A2p7rG+rdxlLiBk8+Y4jj41HZ1tsGso5D8CEuafW3mMNqV78cD7+Yra1\nc7nSaKq2+PaAxP61gc6tDa7oFenrupyGHHITw6xpxQ9bHahf+NJ4E9owarAg\npUQoYIrQaa7L+PIskeSMGsVk1ikdirts/wDSDfuPyNcVd1m/pFr3K+RoA+ke\nzj7vtM/hsb9tNKls4+77TP4ZG/bTSoAHXpTuFGonhkAbkfn7lVS1tdRMkIRI\nO6CrOeQzjl8R+tXJ0qsK1M8kkZ6uPgHtOFcKHp9L7LXpLCypkEJWlQzuE8s+\nBqOcFIGX7p2/Lstl/h0Wa20pWSpDDZPDwXgDPPlXbaZl1E1wpnOSG1J+xSXd\n5RPd61UDA1M5AblJVGYeL7WG1nILTg4BYAOMgZ+NXHscntahtsx6Y71a22/R\n9xtPsrUCAodw5ms2enxtmfsY3FOTHDVe0G4yFqhz7UhEtlCmGkqGVNk8CEnx\n8fGvOkhchFEiZPhxIyXApS1kNgAcgnHEnPj8a49QliU6lltaRLDhQCsZ3TnB\nxnxzjspnf07cusT10z0pCVb286d3cAPYOVa+th4pUis+LjXRONfzEaiaSxAl\nFXo32eFew6o4OFA9vcrvobr7C6mW6hbZQQo44csH/XCr2hT4iI4bQ2SFulMh\nwD1lgJGefLPLyrfd9A2fVbyZlnddW8GwpUb2G5I5ApdPAODGCOGR7q1J6jyQ\nuL8ljTuH4+gZyhQVu4416cZcax1ja0Z5bySM0Tlk0JFtqt17diKTzjxkBCz/\nANTisk+VP0pVsYgLju2mJMiKG643JCnQoeJJ/UcakX0ebjfLyauPG8nQINd1\nm4XFr3K+Rp013ZEWPUsyPGSRCU4Vx+OfszxAz4A48qa7QP8AaDXuPyNZE4OE\nnGXaOJRcW0z6R7OPu+0z+GRv200qxs3+73TP4ZG/bTSrkQMnS6WU6kfx2Nxj\n+iqoazXRCHermgLaWnq3N44C0nsPce0HvAq9el2f955A/wCVG+SqHAAUqsTV\nkkn6fUw6gAuuQnfWZkIRkY7j2Z7+NPemrmzoqUtReffYloCHWwNzkchQwSeH\nf41Cos6TFQUR5TzSFc0oWQD5V0tyWpI3JjqgvscPH40JemRyg5KpdFkS7/Gv\nOun7rBJXbmGUuvKVkJBCcqwO/IPGsW3XLc5b5fe6tTn9mTugDuHuFQiTcYkG\nxuW+3uF56Qft3sYSEZB3R3kkDJ7hTAOfCp1kcOiJa8ZLz/C/bJMt7cZW4hS2\nXVgLQlXnnOedTS33q0Wq2b8PcMxSjiO2AXFZ7cJz29tDZbwZMF9pXEhPWDvO\nOf6Zp30/f3LddIRjoyktqQnJxneBGM+PCtjX3opRjJUmSY8Siy29Rank3+0u\nzmS4xIiK3XUJO6efAntyKY7BPdvS3WH3kNvupIbcUopSF9mTngCe3srk01MQ\n9MTujMKYkpIxzJ4EK/xA8KYWUu225vRST1gWUj3d/wAKt5MzXGS66NrBDhT9\nGvarGdbbjCc04zNb3UuNuDCknBHH4cxwPMVBbUMTmz4H5GiCv7LGudGyLbIb\n37tAiqkQXwPX+zGS0T2pKQQO44qgrYAJbZ/vA4+FYW9+2yvvY3DJb9n0X2b/\nAHeaY/DI37aaVY2bfd5pj8Mj/tppVTKZXO1nZReNZarcuMN62iKpltvckqVv\nbyc8cBJHbUI/l2vh5uWH4r+iiiKgOZFa1Ptp9pYFAAvno6Xw/wBrYB5r+isK\n6Od8PJ7T/wDn+iidVPjJ5upFa1XWEkEqkIAHeaABhV0b7+cn0jT481/RVQX6\n3x9Napm2y4W9iWuBIUy8G/VQ4U88HAIHlmjan7SdKQypK7s24ocCGULc/UDF\nBptHuEe5671BOiq3o8ia442VDBKSeHA8qQmME+5W95SPQrOmGBnO4srKvjyx\nVmaJ2IXvV+l4F9t0i0MxZYUptt8rC07qinjhBHNJ7aqXeG9x3fKi02HbR9LW\nnZpY7Xc7mmLLYQ4FpW0vdBLqyPWAxyIpgiCJ6OerEexcrInjng66OPf7FeXO\njlq1a983KyFZ5qLzuT57lE9bL/arpFTIt1wjSWFcAtpYUM93vrsEtg8nUHzp\n2xgrN9HbWbSipq72hCsYymQ8DjyRWgdGvVSTlM6whQ5EOOf/ADosw82eS0/G\nvQWk/wBYUgfkadI2x6zaUs1slKbW/DhtR3FNklJUlABIzxxwpU77w7xSoA5J\nDTqgd2mSdbZzgPV/OpNSoArefp29O56v5moxetCajuEKRG31oQ8goKkLIIBG\nOBq76VKh2CC5sX2kwlbsGe3KYHshxeFAf9wPzrhlbDNYvhx56A6uUs7yt15r\nBPwozKVFCAgOwjXa1bqbKRx9pUlvj5U7RdiG0bqgyyiNDZHAbzyVEfBNGTSp\nhQO2z7ZNqjTLMkv3Bch+SpKnDvHdGBwwPPnVhwdN3xvHWOk+dWNSpUOyLw7T\nPbx1i8+dO7ER5A9ZQpxpUxGkNKA5ilW6lQB//9k=\n"
    }
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 21063164,
  "option_name_id": 1,
  "position": 21,
  "title": "Blue",
  "image_url": "http://minio.insales.io:9000/insales-test-sber/images/option_values/1/1532/21063164/thumb/image.jpg?1784930385",
  "permalink": "blue"
}
```

### Destroy option_value image

`DELETE /admin/option_names/:option_name_id/option_values/:option_value_id/destroy_image.json`

**Запрос**

```http
DELETE /admin/option_names/1/option_values/20/destroy_image.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20,
  "option_name_id": 1,
  "position": 20,
  "title": "Синий-10",
  "image_url": null,
  "permalink": "siniy-10"
}
```

### Destroy option_value

`DELETE /admin/option_names/:option_name_id/option_values/:option_value_id.json`

**Запрос**

```http
DELETE /admin/option_names/1/option_values/20.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20,
  "jid": 46
}
```

### Get destroy option_value status

`GET /admin/option_names/:option_name_id/option_values/:option_value_id/destroy_status.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `jid` | да | job id |

**Запрос**

```http
GET /admin/option_names/1/option_values/20/destroy_status.json?jid=101
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
jid=101
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "error",
  "error": "Not deleted"
}
```

### Get option_value

`GET /admin/option_names/:option_name_id/option_values/:option_value_id.json`

**Запрос**

```http
GET /admin/option_names/1/option_values/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "option_name_id": 1,
  "position": 2,
  "title": "Зеленый",
  "image_url": null,
  "permalink": null
}
```

### Get option_values for all options

`GET /admin/option_values.json`

**Запрос**

```http
GET /admin/option_values.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "option_name_id": 1,
    "position": 1,
    "title": "Red",
    "image_url": null,
    "permalink": null
  },
  {
    "id": 2,
    "option_name_id": 1,
    "position": 2,
    "title": "Зеленый",
    "image_url": null,
    "permalink": null
  },
  {
    "id": 3,
    "option_name_id": 1,
    "position": 3,
    "title": "Базовая",
    "image_url": null,
    "permalink": null
  },
  {
    "id": 4,
    "option_name_id": 3,
    "position": 1,
    "title": "Синий",
    "image_url": null,
    "permalink": null
  },
  {
    "id": 20,
    "option_name_id": 1,
    "position": 20,
    "title": "Синий-10",
    "image_url": null,
    "permalink": null
  }
]
```

### Get option_values for specific option

`GET /admin/option_names/:option_name_id/option_values.json`

**Запрос**

```http
GET /admin/option_names/1/option_values.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "option_name_id": 1,
    "position": 1,
    "title": "Red",
    "image_url": null,
    "permalink": null
  },
  {
    "id": 2,
    "option_name_id": 1,
    "position": 2,
    "title": "Зеленый",
    "image_url": null,
    "permalink": null
  },
  {
    "id": 3,
    "option_name_id": 1,
    "position": 3,
    "title": "Базовая",
    "image_url": null,
    "permalink": null
  },
  {
    "id": 20,
    "option_name_id": 1,
    "position": 20,
    "title": "Синий-10",
    "image_url": null,
    "permalink": null
  }
]
```

### Update option_value

`PUT /admin/option_names/:option_name_id/option_values/:option_value_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `option_value[position]` |  | position |
| `option_value[image][filename]` |  | image filename, required if attachment present |
| `option_value[image][attachment]` |  | base64 encoded image body |
| `option_value[image][src]` |  | image url |

**Запрос**

```http
PUT /admin/option_names/1/option_values/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "option_value": {
    "position": 3,
    "image": {
      "src": "https://abc.myinsales.ru/images/no_image_compact.jpg"
    }
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "option_name_id": 1,
  "position": 3,
  "title": "Зеленый",
  "image_url": "http://minio.insales.io:9000/insales-test-sber/images/option_values/1/2/2/thumb/no_image_compact.jpg?1784930386",
  "permalink": "zelenyy"
}
```

---

## Order

| Операция | Метод и путь |
|----------|--------------|
| [Create discount](#create-discount) | `PUT /admin/orders/:order_id.json` |
| [Create order line by product_id](#create-order-line-by-product-id) | `PUT /admin/orders/:order_id.json` |
| [Create order line by variant_id](#create-order-line-by-variant-id) | `PUT /admin/orders/:order_id.json` |
| [Create order with bundle](#create-order-with-bundle) | `POST /admin/orders.json` |
| [Create order with product](#create-order-with-product) | `POST /admin/orders.json` |
| [Get order](#get-order) | `GET /admin/orders/:order_id.json` |
| [Get orders count](#get-orders-count) | `GET /admin/orders/count.json` |
| [Get orders](#get-orders) | `GET /admin/orders.json` |
| [Remove discount](#remove-discount) | `PUT /admin/orders/:order_id.json` |
| [Remove order line](#remove-order-line) | `PUT /admin/orders/:order_id.json` |
| [Update custom status](#update-custom-status) | `PUT /admin/orders/:order_id.json` |
| [Update discount](#update-discount) | `PUT /admin/orders/:order_id.json` |
| [Update discount](#update-discount) | `PUT /admin/orders/:order_id.json` |
| [Update order](#update-order) | `PUT /admin/orders/:order_id.json` |
| [Update order line](#update-order-line) | `PUT /admin/orders/:order_id.json` |
| [Update shipping address](#update-shipping-address) | `PUT /admin/orders/:order_id.json` |

### Create discount

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[discounts_attributes]` | да | discounts attributes |
| `order[discounts_attributes][][discount]` | да | discount amount |
| `order[discounts_attributes][][type_id]` |  | 1 - percent (default), 2 - money |
| `order[discounts_attributes][][description]` |  | description text |
| `order[discounts_attributes][][discount_order_lines_ids]` |  | ids of affected order lines (all will be affected if not specified) |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "discounts_attributes": [
      {
        "discount": 100,
        "type_id": 2,
        "description": "-100 (absolute)"
      },
      {
        "discount": 10,
        "description": "10% only for first order line",
        "discount_order_lines_ids": [
          25
        ]
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 779,
      "value": "old value",
      "created_at": "2026-07-25T00:59:23.947+03:00",
      "updated_at": "2026-07-25T00:59:23.947+03:00",
      "field_id": 968,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 25,
      "order_id": 1,
      "sale_price": 100.0,
      "full_sale_price": 55.72,
      "total_price": 200.0,
      "full_total_price": 111.43,
      "discounts_amount": 88.57,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": "10x20x30",
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt (Red)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:24.000+03:00",
      "created_at": "2026-07-25T00:59:24.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 26,
      "order_id": 1,
      "sale_price": 150.0,
      "full_sale_price": 98.57,
      "total_price": 150.0,
      "full_total_price": 98.57,
      "discounts_amount": 51.43,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt (Зеленый / Синий)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:24.000+03:00",
      "created_at": "2026-07-25T00:59:24.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": {
    "id": 99,
    "description": "-20 (absolute)",
    "warning": null,
    "type_id": 2,
    "amount": 20.0,
    "full_amount": 20.0,
    "percent": null,
    "discount": "20.0",
    "reference_id": null,
    "reference_type": null,
    "discount_products_ids": [],
    "discount_order_lines_ids": [],
    "discount_code_id": null,
    "created_at": "2026-07-25T00:59:24.000+03:00",
    "updated_at": "2026-07-25T00:59:24.000+03:00"
  },
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 771,
        "value": "old value",
        "created_at": "2026-07-25T00:59:23.893+03:00",
        "updated_at": "2026-07-25T00:59:23.893+03:00",
        "field_id": 960,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 960,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:23.882+03:00",
        "updated_at": "2026-07-25T00:59:23.882+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:23.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 763,
        "value": "old value",
        "created_at": "2026-07-25T00:59:23.868+03:00",
        "updated_at": "2026-07-25T00:59:23.868+03:00",
        "field_id": 952,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [
    {
      "id": 99,
      "description": "-20 (absolute)",
      "warning": null,
      "type_id": 2,
      "amount": 20.0,
      "full_amount": 20.0,
      "percent": null,
      "discount": "20.0",
      "reference_id": null,
      "reference_type": null,
      "discount_products_ids": [],
      "discount_order_lines_ids": [],
      "discount_code_id": null,
      "created_at": "2026-07-25T00:59:24.000+03:00",
      "updated_at": "2026-07-25T00:59:24.000+03:00"
    },
    {
      "id": 107,
      "description": "-100 (absolute)",
      "warning": null,
      "type_id": 2,
      "amount": 100.0,
      "full_amount": 100.0,
      "percent": null,
      "discount": "100.0",
      "reference_id": null,
      "reference_type": null,
      "discount_products_ids": [],
      "discount_order_lines_ids": [],
      "discount_code_id": null,
      "created_at": "2026-07-25T00:59:24.000+03:00",
      "updated_at": "2026-07-25T00:59:24.000+03:00"
    },
    {
      "id": 115,
      "description": "10% only for first order line",
      "warning": null,
      "type_id": 1,
      "amount": 20.0,
      "full_amount": 20.0,
      "percent": "10.0",
      "discount": "10.0",
      "reference_id": null,
      "reference_type": null,
      "discount_products_ids": [
        1
      ],
      "discount_order_lines_ids": [
        25
      ],
      "discount_code_id": null,
      "created_at": "2026-07-25T00:59:24.000+03:00",
      "updated_at": "2026-07-25T00:59:24.000+03:00"
    }
  ],
  "total_price": 310.0,
  "paid_amount": 0.0,
  "items_price": 210.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:24.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "30.0",
  "warehouse_id": null
}
```

### Create order line by product_id

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[order_lines_attributes]` | да | order lines attributes |
| `order[order_lines_attributes][][quantity]` | да | quantity |
| `order[order_lines_attributes][][vat]` |  | integer, VAT for order line is equal to VAT of its product, by default, products get their value from account by default,if default in account set as 'No VAT', you need to use 'force' to change VAT of order line through API; possible values depend on country, set to -1 if this product is not a subject to VAT |
| `order[order_lines_attributes][][product_id]` | да | product id |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "order_lines_attributes": [
      {
        "product_id": 1,
        "quantity": 1,
        "vat": 22
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 707,
      "value": "old value",
      "created_at": "2026-07-25T00:59:22.519+03:00",
      "updated_at": "2026-07-25T00:59:22.519+03:00",
      "field_id": 896,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 24,
      "order_id": 1,
      "sale_price": 300.0,
      "full_sale_price": 300.0,
      "total_price": 300.0,
      "full_total_price": 300.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": "10x20x30",
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt (Red)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:22.000+03:00",
      "created_at": "2026-07-25T00:59:22.000+03:00",
      "bundle_id": null,
      "vat": 22,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 699,
        "value": "old value",
        "created_at": "2026-07-25T00:59:22.485+03:00",
        "updated_at": "2026-07-25T00:59:22.485+03:00",
        "field_id": 888,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 888,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:22.474+03:00",
        "updated_at": "2026-07-25T00:59:22.474+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:22.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 691,
        "value": "old value",
        "created_at": "2026-07-25T00:59:22.460+03:00",
        "updated_at": "2026-07-25T00:59:22.460+03:00",
        "field_id": 880,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 400.0,
  "paid_amount": 0.0,
  "items_price": 300.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:22.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "210.0",
  "warehouse_id": null
}
```

### Create order line by variant_id

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[order_lines_attributes]` | да | order lines attributes |
| `order[order_lines_attributes][][quantity]` | да | quantity |
| `order[order_lines_attributes][][vat]` |  | integer, VAT for order line is equal to VAT of its product, by default, products get their value from account by default,if default in account set as 'No VAT', you need to use 'force' to change VAT of order line through API; possible values depend on country, set to -1 if this product is not a subject to VAT |
| `order[order_lines_attributes][][variant_id]` | да | variant id |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "order_lines_attributes": [
      {
        "variant_id": 1,
        "quantity": 1,
        "vat": -1
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 683,
      "value": "old value",
      "created_at": "2026-07-25T00:59:21.981+03:00",
      "updated_at": "2026-07-25T00:59:21.981+03:00",
      "field_id": 872,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 23,
      "order_id": 1,
      "sale_price": 300.0,
      "full_sale_price": 300.0,
      "total_price": 300.0,
      "full_total_price": 300.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": "10x20x30",
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt (Red)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:22.000+03:00",
      "created_at": "2026-07-25T00:59:22.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 675,
        "value": "old value",
        "created_at": "2026-07-25T00:59:21.948+03:00",
        "updated_at": "2026-07-25T00:59:21.948+03:00",
        "field_id": 864,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 864,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:21.937+03:00",
        "updated_at": "2026-07-25T00:59:21.937+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:21.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 667,
        "value": "old value",
        "created_at": "2026-07-25T00:59:21.921+03:00",
        "updated_at": "2026-07-25T00:59:21.921+03:00",
        "field_id": 856,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 400.0,
  "paid_amount": 0.0,
  "items_price": 300.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:22.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "210.0",
  "warehouse_id": null
}
```

### Create order with bundle

`POST /admin/orders.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[order_lines_attributes]` | да | order lines attributes |
| `order[order_lines_attributes][][quantity]` | да | quantity |
| `order[order_lines_attributes][][variant_id]` | да | variant id |
| `order[order_lines_attributes][][vat]` |  | do not set vat for bundle, because after order creation bundle will break into bundle components, which have their own VAT rates |
| `order[client]` | да | client info |
| `order[client][name]` |  | client name |
| `order[client][email]` |  | client email |
| `order[client][phone]` | да | client phone |
| `order[client][consent_to_personal_data]` | да | consent to personal data |
| `order[shipping_address_attributes]` | да | array of shipping address attributes |
| `order[shipping_address_attributes][address]` |  | shipping address |
| `order[shipping_address_attributes][full_locality_name]` | да | full locality name |
| `order[delivery_variant_id]` | да | delivery variant id |
| `order[payment_gateway_id]` | да | payment gateway id |
| `order[coupon]` |  | discount code |

**Запрос**

```http
POST /admin/orders.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "order_lines_attributes": [
      {
        "variant_id": 5,
        "quantity": 2
      }
    ],
    "client": {
      "name": "Vasya",
      "email": "vasya@example.com",
      "phone": "79111112233",
      "consent_to_personal_data": true
    },
    "shipping_address_attributes": {
      "address": "Moscow, Krasnaya Presna 24",
      "full_locality_name": "Moscow"
    },
    "delivery_variant_id": 1,
    "payment_gateway_id": 1,
    "coupon": "20% discount"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [],
  "order_lines": [
    {
      "id": 13,
      "order_id": 12,
      "sale_price": 375.0,
      "full_sale_price": 375.0,
      "total_price": 750.0,
      "full_total_price": 750.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt (Зеленый / Синий)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:14.000+03:00",
      "created_at": "2026-07-25T00:59:14.000+03:00",
      "bundle_id": 5,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 14,
      "order_id": 12,
      "sale_price": 375.0,
      "full_sale_price": 375.0,
      "total_price": 750.0,
      "full_total_price": 750.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": null,
      "variant_id": 3,
      "product_id": 3,
      "sku": "425",
      "barcode": null,
      "title": "pants (Базовая)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:14.000+03:00",
      "created_at": "2026-07-25T00:59:14.000+03:00",
      "bundle_id": 5,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": {
    "id": 35,
    "description": "Скидка по купону 20% discount",
    "warning": null,
    "type_id": 1,
    "amount": 300.0,
    "full_amount": 300.0,
    "percent": "20.0",
    "discount": "20.0",
    "reference_id": 236,
    "reference_type": "DiscountCode",
    "discount_products_ids": [],
    "discount_order_lines_ids": [],
    "discount_code_id": 236,
    "created_at": "2026-07-25T00:59:14.000+03:00",
    "updated_at": "2026-07-25T00:59:14.000+03:00"
  },
  "shipping_address": {
    "id": 21,
    "fields_values": [],
    "name": "Vasya",
    "surname": null,
    "middlename": null,
    "phone": "79111112233",
    "formatted_phone": "7(911)111-22-33",
    "full_name": "Vasya",
    "full_locality_name": "",
    "full_delivery_address": "Moscow, Krasnaya Presna 24",
    "address_for_gis": "Moscow, Krasnaya Presna 24",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 552,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:13.615+03:00",
        "updated_at": "2026-07-25T00:59:13.615+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": false
      }
    ],
    "address": "Moscow, Krasnaya Presna 24",
    "country": null,
    "state": null,
    "city": null,
    "zip": null,
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": null,
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": null,
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": null,
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "Moscow, Krasnaya Presna 24",
      "street": "Moscow",
      "street_type": null,
      "house": "Krasnaya Presna 24",
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 81,
    "email": "vasya@example.com",
    "name": "Vasya",
    "phone": "79111112233",
    "created_at": "2026-07-25T00:59:14.000+03:00",
    "updated_at": "2026-07-25T00:59:14.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Vasya",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": []
  },
  "discounts": [
    {
      "id": 35,
      "description": "Скидка по купону 20% discount",
      "warning": null,
      "type_id": 1,
      "amount": 300.0,
      "full_amount": 300.0,
      "percent": "20.0",
      "discount": "20.0",
      "reference_id": 236,
      "reference_type": "DiscountCode",
      "discount_products_ids": [],
      "discount_order_lines_ids": [],
      "discount_code_id": 236,
      "created_at": "2026-07-25T00:59:14.000+03:00",
      "updated_at": "2026-07-25T00:59:14.000+03:00"
    }
  ],
  "total_price": 1300.0,
  "paid_amount": 0.0,
  "items_price": 1200.0,
  "id": 12,
  "key": "69d7bc5241e8f1169009be860ef90cfc",
  "number": 1006,
  "comment": null,
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз (самовывоз из офиса)",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "Наличными курьеру",
  "first_referer": null,
  "first_current_location": null,
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Прямой трафик",
  "referer": "app:digital-goods",
  "current_location": null,
  "query": null,
  "source_domain": null,
  "source": "Приложение",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2026-07-25T00:59:14.000+03:00",
  "updated_at": "2026-07-25T00:59:14.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": null,
  "currency_code": "RUR",
  "cookies": {},
  "account_id": 1,
  "manager_comment": null,
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": 1,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": 100.0,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": [],
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "1200.0",
  "warehouse_id": 1
}
```

### Create order with product

`POST /admin/orders.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[order_lines_attributes]` | да | order lines attributes |
| `order[order_lines_attributes][][quantity]` | да | quantity |
| `order[order_lines_attributes][][variant_id]` | да | variant id |
| `order[order_lines_attributes][][vat]` |  | integer, VAT for order line is equal to VAT of its product, by default, products get their value from account by default,if default in account set as 'No VAT', you need to use 'force' to change VAT of order line through API; possible values depend on country, set to -1 if this product is not a subject to VAT |
| `order[client]` | да | client info |
| `order[client][name]` |  | client name |
| `order[client][email]` |  | client email |
| `order[client][phone]` | да | client phone |
| `order[shipping_address_attributes]` | да | array of shipping address attributes |
| `order[shipping_address_attributes][address]` |  | shipping address |
| `order[delivery_variant_id]` | да | delivery variant id |
| `order[payment_gateway_id]` | да | payment gateway id |
| `order[coupon]` |  | discount code |

**Запрос**

```http
POST /admin/orders.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "order_lines_attributes": [
      {
        "variant_id": 3,
        "quantity": 2,
        "vat": 10
      }
    ],
    "client": {
      "name": "Vasya",
      "email": "vasya@example.com",
      "phone": "79111112233"
    },
    "shipping_address_attributes": {
      "address": "Moscow, Krasnaya Presna 24"
    },
    "delivery_variant_id": 1,
    "payment_gateway_id": 1,
    "coupon": "20% discount"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [],
  "order_lines": [
    {
      "id": 10,
      "order_id": 10,
      "sale_price": 600.0,
      "full_sale_price": 600.0,
      "total_price": 1200.0,
      "full_total_price": 1200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": null,
      "variant_id": 3,
      "product_id": 3,
      "sku": "425",
      "barcode": null,
      "title": "pants (Базовая)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:12.000+03:00",
      "created_at": "2026-07-25T00:59:12.000+03:00",
      "bundle_id": null,
      "vat": 10,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": {
    "id": 19,
    "description": "Скидка по купону 20% discount",
    "warning": null,
    "type_id": 1,
    "amount": 240.0,
    "full_amount": 240.0,
    "percent": "20.0",
    "discount": "20.0",
    "reference_id": 220,
    "reference_type": "DiscountCode",
    "discount_products_ids": [],
    "discount_order_lines_ids": [],
    "discount_code_id": 220,
    "created_at": "2026-07-25T00:59:12.000+03:00",
    "updated_at": "2026-07-25T00:59:12.000+03:00"
  },
  "shipping_address": {
    "id": 19,
    "fields_values": [],
    "name": "Vasya",
    "surname": null,
    "middlename": null,
    "phone": "79111112233",
    "formatted_phone": "7(911)111-22-33",
    "full_name": "Vasya",
    "full_locality_name": "",
    "full_delivery_address": "Moscow, Krasnaya Presna 24",
    "address_for_gis": "Moscow, Krasnaya Presna 24",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 504,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:12.296+03:00",
        "updated_at": "2026-07-25T00:59:12.296+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": false
      }
    ],
    "address": "Moscow, Krasnaya Presna 24",
    "country": null,
    "state": null,
    "city": null,
    "zip": null,
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": null,
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": null,
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": null,
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "Moscow, Krasnaya Presna 24",
      "street": "Moscow",
      "street_type": null,
      "house": "Krasnaya Presna 24",
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 65,
    "email": "vasya@example.com",
    "name": "Vasya",
    "phone": "79111112233",
    "created_at": "2026-07-25T00:59:12.000+03:00",
    "updated_at": "2026-07-25T00:59:12.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Vasya",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": []
  },
  "discounts": [
    {
      "id": 19,
      "description": "Скидка по купону 20% discount",
      "warning": null,
      "type_id": 1,
      "amount": 240.0,
      "full_amount": 240.0,
      "percent": "20.0",
      "discount": "20.0",
      "reference_id": 220,
      "reference_type": "DiscountCode",
      "discount_products_ids": [],
      "discount_order_lines_ids": [],
      "discount_code_id": 220,
      "created_at": "2026-07-25T00:59:12.000+03:00",
      "updated_at": "2026-07-25T00:59:12.000+03:00"
    }
  ],
  "total_price": 1060.0,
  "paid_amount": 0.0,
  "items_price": 960.0,
  "id": 10,
  "key": "9aabef7bbab5d0b349a1434a5cb97d78",
  "number": 1006,
  "comment": null,
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз (самовывоз из офиса)",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "Наличными курьеру",
  "first_referer": null,
  "first_current_location": null,
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Прямой трафик",
  "referer": "app:digital-goods",
  "current_location": null,
  "query": null,
  "source_domain": null,
  "source": "Приложение",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2026-07-25T00:59:12.000+03:00",
  "updated_at": "2026-07-25T00:59:12.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": null,
  "currency_code": "RUR",
  "cookies": {},
  "account_id": 1,
  "manager_comment": null,
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": 1,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": 100.0,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": [],
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "960.0",
  "warehouse_id": 1
}
```

### Get order

`GET /admin/orders/:order_id.json`

**Запрос**

```http
GET /admin/orders/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 179,
      "value": "old value",
      "created_at": "2026-07-25T00:59:09.637+03:00",
      "updated_at": "2026-07-25T00:59:09.637+03:00",
      "field_id": 368,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 1,
      "order_id": 1,
      "sale_price": 300.0,
      "full_sale_price": 300.0,
      "total_price": 300.0,
      "full_total_price": 300.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": null,
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": [
        "code1",
        "code2"
      ],
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 2,
      "order_id": 1,
      "sale_price": 600.0,
      "full_sale_price": 600.0,
      "total_price": 1200.0,
      "full_total_price": 1200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": 20,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 171,
        "value": "old value",
        "created_at": "2026-07-25T00:59:09.604+03:00",
        "updated_at": "2026-07-25T00:59:09.604+03:00",
        "field_id": 360,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 360,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:09.593+03:00",
        "updated_at": "2026-07-25T00:59:09.593+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:09.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 163,
        "value": "old value",
        "created_at": "2026-07-25T00:59:09.577+03:00",
        "updated_at": "2026-07-25T00:59:09.577+03:00",
        "field_id": 352,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 1600.0,
  "paid_amount": 0.0,
  "items_price": 1500.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:09.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "1500.0",
  "warehouse_id": null
}
```

### Get orders count

`GET /admin/orders/count.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `fulfillment_status` |  | array of statuses |
| `delivery_variant` |  | array of delivery variant ids |
| `payment_gateway_id` |  | array of payment gateway ids |

**Запрос**

```http
GET /admin/orders/count.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "count": 6
}
```

### Get orders

`GET /admin/orders.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `fulfillment_status` |  | array of statuses |
| `delivery_variant` |  | array of delivery variant ids |
| `payment_gateway_id` |  | array of payment gateway ids |
| `updated_since` |  | Set datetime to get only data updated after it. |
| `from_id` |  | Set id to get only data starting from it. |
| `per_page` |  | Quantity of objects per page. |

**Запрос**

```http
GET /admin/orders.json?fulfillment_status[]=new&fulfillment_status[]=delivered&delivery_variant[]=1&delivery_variant[]=2&payment_gateway_id[]=1&payment_gateway_id[]=3
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
fulfillment_status=["new", "delivered"]
delivery_variant=["1", "2"]
payment_gateway_id=["1", "3"]
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "fields_values": [],
    "order_lines": [],
    "order_changes": [],
    "discount": null,
    "shipping_address": {
      "id": 9,
      "fields_values": [],
      "name": "Name Name",
      "surname": null,
      "middlename": null,
      "phone": "+7(495)888-99-06",
      "formatted_phone": "+7(495)888-99-06",
      "full_name": "Name Name",
      "full_locality_name": null,
      "full_delivery_address": "567887, state, Manchester, address",
      "address_for_gis": "567887, state, Manchester, address",
      "location_valid": true,
      "recipient_fields": [
        {
          "id": 11,
          "destiny": 2,
          "position": 1,
          "office_title": "Имя",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "name",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::ObligatoryTextField"
        },
        {
          "id": 12,
          "destiny": 2,
          "position": 2,
          "office_title": "Телефон",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "phone",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Phone"
        }
      ],
      "backoffice_fields": [
        {
          "id": 4,
          "destiny": 1,
          "position": 3,
          "office_title": "Адрес",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": "address",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextArea"
        },
        {
          "id": 456,
          "destiny": 1,
          "position": 22,
          "office_title": "special",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": null,
          "created_at": "2026-07-25T00:59:11.082+03:00",
          "updated_at": "2026-07-25T00:59:11.082+03:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": "special",
          "example": null,
          "type": "Field::TextField"
        },
        {
          "id": 1,
          "destiny": 1,
          "position": 0,
          "office_title": "Страна",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "country",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Country"
        },
        {
          "id": 2,
          "destiny": 1,
          "position": 1,
          "office_title": "Регион",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "state",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 3,
          "destiny": 1,
          "position": 2,
          "office_title": "Город",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "city",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 5,
          "destiny": 1,
          "position": 4,
          "office_title": "Почтовый индекс",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "zip",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        }
      ],
      "no_delivery": false,
      "kladr_autodetected_address": "",
      "country_options": [
        {
          "code": "RU",
          "title": "Россия",
          "selected": true
        }
      ],
      "address": "address",
      "country": "RU",
      "state": "state",
      "city": "Manchester",
      "zip": "567887",
      "street": null,
      "house": null,
      "flat": null,
      "entrance": null,
      "doorphone": null,
      "floor": null,
      "comment": null,
      "kladr_json": null,
      "location": {
        "kladr_code": null,
        "zip": "567887",
        "kladr_zip": null,
        "region_zip": null,
        "country": "RU",
        "state": "state",
        "state_type": null,
        "area": null,
        "area_type": null,
        "city": "Manchester",
        "city_type": null,
        "settlement": null,
        "settlement_type": null,
        "address": "address",
        "street": "address",
        "street_type": null,
        "house": null,
        "flat": null,
        "is_kladr": false,
        "latitude": null,
        "longitude": null,
        "autodetected": null
      }
    },
    "client": {
      "id": 3,
      "email": "some-one2@yandex.ru",
      "name": "Just Man",
      "phone": "+7(926)111-22-33",
      "created_at": "2009-04-30T04:37:32.000+04:00",
      "updated_at": "2009-04-30T12:20:34.000+04:00",
      "comment": null,
      "registered": false,
      "subscribe": true,
      "client_group_id": 1,
      "bonus_points": 0,
      "type": "Client::Individual",
      "correspondent_account": null,
      "settlement_account": null,
      "consent_to_personal_data": null,
      "o_auth_provider": null,
      "messenger_subscription": null,
      "external_id": null,
      "external_service": null,
      "contact_name": "Just Man",
      "surname": null,
      "middlename": null,
      "progressive_discount": "25.0",
      "group_discount": "25.0",
      "ip_addr": "",
      "fields_values": []
    },
    "discounts": [],
    "total_price": 0.0,
    "paid_amount": 0.0,
    "items_price": 0.0,
    "id": 8,
    "key": null,
    "number": 1008,
    "comment": "comment",
    "archived": false,
    "delivery_title": "самовывоз",
    "delivery_description": "самовывоз",
    "delivery_price": 0.0,
    "full_delivery_price": 0.0,
    "payment_description": null,
    "payment_title": "наличными курьеру",
    "first_referer": null,
    "first_current_location": null,
    "first_query": null,
    "first_source_domain": null,
    "first_source": "",
    "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
    "current_location": null,
    "query": null,
    "source_domain": null,
    "source": "",
    "fulfillment_status": "new",
    "custom_status": {
      "permalink": "new",
      "title": "Новый"
    },
    "delivered_at": null,
    "accepted_at": null,
    "created_at": "2011-03-18T10:53:14.000+03:00",
    "updated_at": "2011-03-18T10:53:15.000+03:00",
    "financial_status": "paid",
    "delivery_date": null,
    "delivery_from_hour": null,
    "delivery_from_minutes": null,
    "delivery_to_hour": null,
    "delivery_to_minutes": null,
    "delivery_time": "",
    "paid_at": null,
    "delivery_variant_id": 1,
    "payment_gateway_id": 3,
    "margin": "0.0",
    "margin_amount": "0.0",
    "client_transaction_id": null,
    "currency_code": "RUR",
    "cookies": null,
    "account_id": 1,
    "manager_comment": null,
    "locale": "ru",
    "delivery_info": {
      "delivery_variant_id": null,
      "tariff_id": null,
      "title": null,
      "description": null,
      "price": null,
      "shipping_company": null,
      "shipping_company_handle": null,
      "delivery_interval": {
        "min_days": null,
        "max_days": null,
        "description": ""
      },
      "errors": [],
      "warnings": [],
      "outlet": {
        "id": null,
        "external_id": null,
        "latitude": null,
        "longitude": null,
        "type": null,
        "title": null,
        "description": null,
        "address": null,
        "payment_method": null,
        "source_id": null
      },
      "not_available": null
    },
    "responsible_user_id": null,
    "total_profit": "0.0",
    "warehouse_id": null
  },
  {
    "fields_values": [],
    "order_lines": [],
    "order_changes": [],
    "discount": null,
    "shipping_address": {
      "id": 2,
      "fields_values": [],
      "name": "JD",
      "surname": null,
      "middlename": null,
      "phone": "+7(495)111-11-11",
      "formatted_phone": "+7(495)111-11-11",
      "full_name": "JD",
      "full_locality_name": null,
      "full_delivery_address": "123456, Moscow, address",
      "address_for_gis": "123456, Moscow, address",
      "location_valid": true,
      "recipient_fields": [
        {
          "id": 11,
          "destiny": 2,
          "position": 1,
          "office_title": "Имя",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "name",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::ObligatoryTextField"
        },
        {
          "id": 12,
          "destiny": 2,
          "position": 2,
          "office_title": "Телефон",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "phone",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Phone"
        }
      ],
      "backoffice_fields": [
        {
          "id": 4,
          "destiny": 1,
          "position": 3,
          "office_title": "Адрес",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": "address",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextArea"
        },
        {
          "id": 456,
          "destiny": 1,
          "position": 22,
          "office_title": "special",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": null,
          "created_at": "2026-07-25T00:59:11.082+03:00",
          "updated_at": "2026-07-25T00:59:11.082+03:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": "special",
          "example": null,
          "type": "Field::TextField"
        },
        {
          "id": 1,
          "destiny": 1,
          "position": 0,
          "office_title": "Страна",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "country",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Country"
        },
        {
          "id": 2,
          "destiny": 1,
          "position": 1,
          "office_title": "Регион",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "state",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 3,
          "destiny": 1,
          "position": 2,
          "office_title": "Город",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "city",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 5,
          "destiny": 1,
          "position": 4,
          "office_title": "Почтовый индекс",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "zip",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        }
      ],
      "no_delivery": false,
      "kladr_autodetected_address": "",
      "country_options": [
        {
          "code": "RU",
          "title": "Россия",
          "selected": true
        }
      ],
      "address": "address",
      "country": "RU",
      "state": "Moscow",
      "city": "Moscow",
      "zip": "123456",
      "street": null,
      "house": null,
      "flat": null,
      "entrance": null,
      "doorphone": null,
      "floor": null,
      "comment": null,
      "kladr_json": null,
      "location": {
        "kladr_code": null,
        "zip": "123456",
        "kladr_zip": null,
        "region_zip": null,
        "country": "RU",
        "state": "Moscow",
        "state_type": null,
        "area": null,
        "area_type": null,
        "city": "Moscow",
        "city_type": null,
        "settlement": null,
        "settlement_type": null,
        "address": "address",
        "street": "address",
        "street_type": null,
        "house": null,
        "flat": null,
        "is_kladr": false,
        "latitude": null,
        "longitude": null,
        "autodetected": null
      }
    },
    "client": {
      "id": 7,
      "email": "john2@doe.com",
      "name": "John Doe",
      "phone": "+7(926)111-11-11",
      "created_at": "2009-03-30T04:37:32.000+04:00",
      "updated_at": "2009-03-31T12:20:34.000+04:00",
      "comment": null,
      "registered": false,
      "subscribe": true,
      "client_group_id": null,
      "bonus_points": 100,
      "type": "Client::Individual",
      "correspondent_account": null,
      "settlement_account": null,
      "consent_to_personal_data": null,
      "o_auth_provider": null,
      "messenger_subscription": null,
      "external_id": null,
      "external_service": null,
      "contact_name": "John Doe",
      "surname": null,
      "middlename": null,
      "progressive_discount": null,
      "group_discount": null,
      "ip_addr": "",
      "fields_values": []
    },
    "discounts": [],
    "total_price": 100.0,
    "paid_amount": 0.0,
    "items_price": 0.0,
    "id": 6,
    "key": "merged_client_order",
    "number": 10011,
    "comment": "comment",
    "archived": false,
    "delivery_title": "самовывоз",
    "delivery_description": "самовывоз",
    "delivery_price": 100.0,
    "full_delivery_price": 100.0,
    "payment_description": null,
    "payment_title": "наличными курьеру",
    "first_referer": "http://google.com?q=my_query",
    "first_current_location": "/my_query",
    "first_query": null,
    "first_source_domain": null,
    "first_source": "Google",
    "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
    "current_location": "/",
    "query": null,
    "source_domain": null,
    "source": "Прямой трафик",
    "fulfillment_status": "new",
    "custom_status": {
      "permalink": "new",
      "title": "Новый"
    },
    "delivered_at": null,
    "accepted_at": null,
    "created_at": "2009-03-18T10:53:14.000+03:00",
    "updated_at": "2009-03-18T10:53:15.000+03:00",
    "financial_status": "pending",
    "delivery_date": null,
    "delivery_from_hour": null,
    "delivery_from_minutes": null,
    "delivery_to_hour": null,
    "delivery_to_minutes": null,
    "delivery_time": "",
    "paid_at": null,
    "delivery_variant_id": 1,
    "payment_gateway_id": 1,
    "margin": "0.0",
    "margin_amount": "0.0",
    "client_transaction_id": null,
    "currency_code": "RUR",
    "cookies": null,
    "account_id": 1,
    "manager_comment": "manager comment1",
    "locale": "ru",
    "delivery_info": {
      "delivery_variant_id": null,
      "tariff_id": null,
      "title": null,
      "description": null,
      "price": null,
      "shipping_company": null,
      "shipping_company_handle": null,
      "delivery_interval": {
        "min_days": null,
        "max_days": null,
        "description": ""
      },
      "errors": [],
      "warnings": [],
      "outlet": {
        "id": null,
        "external_id": null,
        "latitude": null,
        "longitude": null,
        "type": null,
        "title": null,
        "description": null,
        "address": null,
        "payment_method": null,
        "source_id": null
      },
      "not_available": null
    },
    "responsible_user_id": null,
    "total_profit": "0.0",
    "warehouse_id": null
  },
  {
    "fields_values": [
      {
        "id": 1,
        "value": "1234",
        "created_at": "2026-07-25T00:58:09.208+03:00",
        "updated_at": "2026-07-25T00:58:09.208+03:00",
        "field_id": 174,
        "type": "Текст",
        "name": "АКС айди отгрузки",
        "handle": "aks_shipment_id"
      }
    ],
    "order_lines": [
      {
        "id": 7,
        "order_id": 4,
        "sale_price": 800.0,
        "full_sale_price": 643.2,
        "total_price": 800.0,
        "full_total_price": 643.2,
        "discounts_amount": 264.0,
        "quantity": 1,
        "reserved_quantity": null,
        "weight": null,
        "dimensions": null,
        "variant_id": 1,
        "product_id": null,
        "sku": "425",
        "barcode": null,
        "title": "Cool t-shirt",
        "unit": "pce",
        "comment": null,
        "updated_at": "2011-03-18T10:53:15.000+03:00",
        "created_at": "2011-03-18T10:53:14.000+03:00",
        "bundle_id": null,
        "vat": -1,
        "fiscal_product_type": 1,
        "requires_marking": null,
        "marking_codes": null,
        "accessory_lines": [],
        "external_variant_id": null
      },
      {
        "id": 8,
        "order_id": 4,
        "sale_price": 100.0,
        "full_sale_price": 80.4,
        "total_price": 200.0,
        "full_total_price": 160.8,
        "discounts_amount": 66.0,
        "quantity": 2,
        "reserved_quantity": null,
        "weight": "0.3",
        "dimensions": null,
        "variant_id": 2,
        "product_id": null,
        "sku": "426",
        "barcode": null,
        "title": "shirt",
        "unit": "pce",
        "comment": null,
        "updated_at": "2009-03-18T10:53:15.000+03:00",
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "bundle_id": null,
        "vat": -1,
        "fiscal_product_type": 1,
        "requires_marking": null,
        "marking_codes": null,
        "accessory_lines": [],
        "external_variant_id": null
      }
    ],
    "order_changes": [],
    "discount": {
      "id": 2,
      "description": "рублевая скидка",
      "warning": null,
      "type_id": 2,
      "amount": 330.0,
      "full_amount": 396.0,
      "percent": null,
      "discount": "330.0",
      "reference_id": null,
      "reference_type": null,
      "discount_products_ids": [],
      "discount_order_lines_ids": [],
      "discount_code_id": null,
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "updated_at": "2009-03-18T10:53:15.000+03:00"
    },
    "shipping_address": null,
    "client": {
      "id": 3,
      "email": "some-one2@yandex.ru",
      "name": "Just Man",
      "phone": "+7(926)111-22-33",
      "created_at": "2009-04-30T04:37:32.000+04:00",
      "updated_at": "2009-04-30T12:20:34.000+04:00",
      "comment": null,
      "registered": false,
      "subscribe": true,
      "client_group_id": 1,
      "bonus_points": 0,
      "type": "Client::Individual",
      "correspondent_account": null,
      "settlement_account": null,
      "consent_to_personal_data": null,
      "o_auth_provider": null,
      "messenger_subscription": null,
      "external_id": null,
      "external_service": null,
      "contact_name": "Just Man",
      "surname": null,
      "middlename": null,
      "progressive_discount": "25.0",
      "group_discount": "25.0",
      "ip_addr": "",
      "fields_values": []
    },
    "discounts": [
      {
        "id": 2,
        "description": "рублевая скидка",
        "warning": null,
        "type_id": 2,
        "amount": 330.0,
        "full_amount": 396.0,
        "percent": null,
        "discount": "330.0",
        "reference_id": null,
        "reference_type": null,
        "discount_products_ids": [],
        "discount_order_lines_ids": [],
        "discount_code_id": null,
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "updated_at": "2009-03-18T10:53:15.000+03:00"
      }
    ],
    "total_price": 804.0,
    "paid_amount": 0.0,
    "items_price": 804.0,
    "id": 4,
    "key": null,
    "number": 1005,
    "comment": "comment",
    "archived": false,
    "delivery_title": "самовывоз",
    "delivery_description": "самовывоз",
    "delivery_price": 0.0,
    "full_delivery_price": 0.0,
    "payment_description": null,
    "payment_title": "наличными курьеру",
    "first_referer": null,
    "first_current_location": null,
    "first_query": null,
    "first_source_domain": null,
    "first_source": "",
    "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
    "current_location": null,
    "query": null,
    "source_domain": null,
    "source": "",
    "fulfillment_status": "new",
    "custom_status": {
      "permalink": "new",
      "title": "Новый"
    },
    "delivered_at": null,
    "accepted_at": null,
    "created_at": "2011-03-18T10:53:14.000+03:00",
    "updated_at": "2011-03-18T10:53:15.000+03:00",
    "financial_status": "paid",
    "delivery_date": null,
    "delivery_from_hour": null,
    "delivery_from_minutes": null,
    "delivery_to_hour": null,
    "delivery_to_minutes": null,
    "delivery_time": "",
    "paid_at": null,
    "delivery_variant_id": 1,
    "payment_gateway_id": 3,
    "margin": "20.0",
    "margin_amount": "134.0",
    "client_transaction_id": null,
    "currency_code": "RUR",
    "cookies": null,
    "account_id": 1,
    "manager_comment": null,
    "locale": "ru",
    "delivery_info": {
      "delivery_variant_id": null,
      "tariff_id": null,
      "title": null,
      "description": null,
      "price": null,
      "shipping_company": null,
      "shipping_company_handle": null,
      "delivery_interval": {
        "min_days": null,
        "max_days": null,
        "description": ""
      },
      "errors": [],
      "warnings": [],
      "outlet": {
        "id": null,
        "external_id": null,
        "latitude": null,
        "longitude": null,
        "type": null,
        "title": null,
        "description": null,
        "address": null,
        "payment_method": null,
        "source_id": null
      },
      "not_available": null
    },
    "responsible_user_id": null,
    "total_profit": "804.0",
    "warehouse_id": null
  },
  {
    "fields_values": [
      {
        "id": 2,
        "value": "",
        "created_at": "2026-07-25T00:58:09.208+03:00",
        "updated_at": "2026-07-25T00:58:09.208+03:00",
        "field_id": 174,
        "type": "Текст",
        "name": "АКС айди отгрузки",
        "handle": "aks_shipment_id"
      }
    ],
    "order_lines": [
      {
        "id": 5,
        "order_id": 3,
        "sale_price": 800.0,
        "full_sale_price": 672.0,
        "total_price": 800.0,
        "full_total_price": 672.0,
        "discounts_amount": 240.0,
        "quantity": 1,
        "reserved_quantity": null,
        "weight": null,
        "dimensions": null,
        "variant_id": 1,
        "product_id": null,
        "sku": "425",
        "barcode": null,
        "title": "Cool t-shirt",
        "unit": "pce",
        "comment": null,
        "updated_at": "2011-03-18T10:53:15.000+03:00",
        "created_at": "2011-03-18T10:53:14.000+03:00",
        "bundle_id": null,
        "vat": -1,
        "fiscal_product_type": 1,
        "requires_marking": null,
        "marking_codes": null,
        "accessory_lines": [],
        "external_variant_id": null
      },
      {
        "id": 6,
        "order_id": 3,
        "sale_price": 100.0,
        "full_sale_price": 84.0,
        "total_price": 200.0,
        "full_total_price": 168.0,
        "discounts_amount": 60.0,
        "quantity": 2,
        "reserved_quantity": null,
        "weight": "0.3",
        "dimensions": null,
        "variant_id": 2,
        "product_id": null,
        "sku": "426",
        "barcode": null,
        "title": "shirt",
        "unit": "pce",
        "comment": null,
        "updated_at": "2009-03-18T10:53:15.000+03:00",
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "bundle_id": null,
        "vat": 0,
        "fiscal_product_type": 1,
        "requires_marking": null,
        "marking_codes": null,
        "accessory_lines": [],
        "external_variant_id": null
      }
    ],
    "order_changes": [],
    "discount": {
      "id": 1,
      "description": "процентная скидка",
      "warning": null,
      "type_id": 1,
      "amount": 300.0,
      "full_amount": 360.0,
      "percent": "30.0",
      "discount": "30.0",
      "reference_id": null,
      "reference_type": null,
      "discount_products_ids": [],
      "discount_order_lines_ids": [],
      "discount_code_id": null,
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "updated_at": "2009-03-18T10:53:15.000+03:00"
    },
    "shipping_address": {
      "id": 11,
      "fields_values": [],
      "name": "Name Name",
      "surname": null,
      "middlename": null,
      "phone": "+7(495)888-99-08",
      "formatted_phone": "+7(495)888-99-08",
      "full_name": "Name Name",
      "full_locality_name": null,
      "full_delivery_address": "567889, state, Manchester, address",
      "address_for_gis": "567889, state, Manchester, address",
      "location_valid": true,
      "recipient_fields": [
        {
          "id": 11,
          "destiny": 2,
          "position": 1,
          "office_title": "Имя",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "name",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::ObligatoryTextField"
        },
        {
          "id": 12,
          "destiny": 2,
          "position": 2,
          "office_title": "Телефон",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "phone",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Phone"
        }
      ],
      "backoffice_fields": [
        {
          "id": 4,
          "destiny": 1,
          "position": 3,
          "office_title": "Адрес",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": "address",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextArea"
        },
        {
          "id": 456,
          "destiny": 1,
          "position": 22,
          "office_title": "special",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": null,
          "created_at": "2026-07-25T00:59:11.082+03:00",
          "updated_at": "2026-07-25T00:59:11.082+03:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": "special",
          "example": null,
          "type": "Field::TextField"
        },
        {
          "id": 1,
          "destiny": 1,
          "position": 0,
          "office_title": "Страна",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "country",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Country"
        },
        {
          "id": 2,
          "destiny": 1,
          "position": 1,
          "office_title": "Регион",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "state",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 3,
          "destiny": 1,
          "position": 2,
          "office_title": "Город",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "city",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 5,
          "destiny": 1,
          "position": 4,
          "office_title": "Почтовый индекс",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "zip",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        }
      ],
      "no_delivery": false,
      "kladr_autodetected_address": "",
      "country_options": [
        {
          "code": "RU",
          "title": "Россия",
          "selected": true
        }
      ],
      "address": "address",
      "country": "RU",
      "state": "state",
      "city": "Manchester",
      "zip": "567889",
      "street": null,
      "house": null,
      "flat": null,
      "entrance": null,
      "doorphone": null,
      "floor": null,
      "comment": null,
      "kladr_json": null,
      "location": {
        "kladr_code": null,
        "zip": "567889",
        "kladr_zip": null,
        "region_zip": null,
        "country": "RU",
        "state": "state",
        "state_type": null,
        "area": null,
        "area_type": null,
        "city": "Manchester",
        "city_type": null,
        "settlement": null,
        "settlement_type": null,
        "address": "address",
        "street": "address",
        "street_type": null,
        "house": null,
        "flat": null,
        "is_kladr": false,
        "latitude": null,
        "longitude": null,
        "autodetected": null
      }
    },
    "client": {
      "id": 3,
      "email": "some-one2@yandex.ru",
      "name": "Just Man",
      "phone": "+7(926)111-22-33",
      "created_at": "2009-04-30T04:37:32.000+04:00",
      "updated_at": "2009-04-30T12:20:34.000+04:00",
      "comment": null,
      "registered": false,
      "subscribe": true,
      "client_group_id": 1,
      "bonus_points": 0,
      "type": "Client::Individual",
      "correspondent_account": null,
      "settlement_account": null,
      "consent_to_personal_data": null,
      "o_auth_provider": null,
      "messenger_subscription": null,
      "external_id": null,
      "external_service": null,
      "contact_name": "Just Man",
      "surname": null,
      "middlename": null,
      "progressive_discount": "25.0",
      "group_discount": "25.0",
      "ip_addr": "",
      "fields_values": []
    },
    "discounts": [
      {
        "id": 1,
        "description": "процентная скидка",
        "warning": null,
        "type_id": 1,
        "amount": 300.0,
        "full_amount": 360.0,
        "percent": "30.0",
        "discount": "30.0",
        "reference_id": null,
        "reference_type": null,
        "discount_products_ids": [],
        "discount_order_lines_ids": [],
        "discount_code_id": null,
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "updated_at": "2009-03-18T10:53:15.000+03:00"
      }
    ],
    "total_price": 840.0,
    "paid_amount": 0.0,
    "items_price": 840.0,
    "id": 3,
    "key": null,
    "number": 1004,
    "comment": "comment",
    "archived": false,
    "delivery_title": "самовывоз",
    "delivery_description": "самовывоз",
    "delivery_price": 0.0,
    "full_delivery_price": 0.0,
    "payment_description": null,
    "payment_title": "наличными курьеру",
    "first_referer": null,
    "first_current_location": null,
    "first_query": null,
    "first_source_domain": null,
    "first_source": "",
    "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
    "current_location": null,
    "query": null,
    "source_domain": null,
    "source": "",
    "fulfillment_status": "new",
    "custom_status": {
      "permalink": "new",
      "title": "Новый"
    },
    "delivered_at": null,
    "accepted_at": null,
    "created_at": "2011-03-18T10:53:14.000+03:00",
    "updated_at": "2011-03-18T10:53:15.000+03:00",
    "financial_status": "paid",
    "delivery_date": null,
    "delivery_from_hour": null,
    "delivery_from_minutes": null,
    "delivery_to_hour": null,
    "delivery_to_minutes": null,
    "delivery_time": "",
    "paid_at": null,
    "delivery_variant_id": 1,
    "payment_gateway_id": 3,
    "margin": "20.0",
    "margin_amount": "140.0",
    "client_transaction_id": null,
    "currency_code": "RUR",
    "cookies": null,
    "account_id": 1,
    "manager_comment": null,
    "locale": "ru",
    "delivery_info": {
      "delivery_variant_id": null,
      "tariff_id": null,
      "title": null,
      "description": null,
      "price": null,
      "shipping_company": null,
      "shipping_company_handle": null,
      "delivery_interval": {
        "min_days": null,
        "max_days": null,
        "description": ""
      },
      "errors": [],
      "warnings": [],
      "outlet": {
        "id": null,
        "external_id": null,
        "latitude": null,
        "longitude": null,
        "type": null,
        "title": null,
        "description": null,
        "address": null,
        "payment_method": null,
        "source_id": null
      },
      "not_available": null
    },
    "responsible_user_id": null,
    "total_profit": "840.0",
    "warehouse_id": null
  },
  {
    "fields_values": [
      {
        "id": 275,
        "value": "old value",
        "created_at": "2026-07-25T00:59:11.124+03:00",
        "updated_at": "2026-07-25T00:59:11.124+03:00",
        "field_id": 464,
        "type": "Текст",
        "name": "special",
        "handle": "my_awesome_field"
      }
    ],
    "order_lines": [
      {
        "id": 1,
        "order_id": 1,
        "sale_price": 300.0,
        "full_sale_price": 300.0,
        "total_price": 300.0,
        "full_total_price": 300.0,
        "discounts_amount": 0.0,
        "quantity": 1,
        "reserved_quantity": null,
        "weight": null,
        "dimensions": null,
        "variant_id": 1,
        "product_id": 1,
        "sku": "423",
        "barcode": null,
        "title": "Cool t-shirt",
        "unit": "pce",
        "comment": null,
        "updated_at": "2009-03-18T10:53:15.000+03:00",
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "bundle_id": null,
        "vat": -1,
        "fiscal_product_type": 1,
        "requires_marking": null,
        "marking_codes": [
          "code1",
          "code2"
        ],
        "accessory_lines": [],
        "external_variant_id": null
      },
      {
        "id": 2,
        "order_id": 1,
        "sale_price": 600.0,
        "full_sale_price": 600.0,
        "total_price": 1200.0,
        "full_total_price": 1200.0,
        "discounts_amount": 0.0,
        "quantity": 2,
        "reserved_quantity": null,
        "weight": "0.3",
        "dimensions": null,
        "variant_id": 2,
        "product_id": 2,
        "sku": "424",
        "barcode": null,
        "title": "shirt",
        "unit": "pce",
        "comment": null,
        "updated_at": "2009-03-18T10:53:15.000+03:00",
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "bundle_id": null,
        "vat": 20,
        "fiscal_product_type": 1,
        "requires_marking": null,
        "marking_codes": null,
        "accessory_lines": [],
        "external_variant_id": null
      }
    ],
    "order_changes": [],
    "discount": null,
    "shipping_address": {
      "id": 1,
      "fields_values": [
        {
          "id": 267,
          "value": "old value",
          "created_at": "2026-07-25T00:59:11.092+03:00",
          "updated_at": "2026-07-25T00:59:11.092+03:00",
          "field_id": 456,
          "type": "Текст",
          "name": "special",
          "handle": null
        }
      ],
      "name": "Tom Dellay",
      "surname": null,
      "middlename": null,
      "phone": "+7(495)212-85-06",
      "formatted_phone": "+7(495)212-85-06",
      "full_name": "Tom Dellay",
      "full_locality_name": null,
      "full_delivery_address": "123456, state, London, address, old value",
      "address_for_gis": "123456, state, London, address",
      "location_valid": true,
      "recipient_fields": [
        {
          "id": 11,
          "destiny": 2,
          "position": 1,
          "office_title": "Имя",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "name",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::ObligatoryTextField"
        },
        {
          "id": 12,
          "destiny": 2,
          "position": 2,
          "office_title": "Телефон",
          "for_buyer": true,
          "obligatory": true,
          "active": true,
          "system_name": "phone",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Phone"
        }
      ],
      "backoffice_fields": [
        {
          "id": 4,
          "destiny": 1,
          "position": 3,
          "office_title": "Адрес",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": "address",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextArea"
        },
        {
          "id": 456,
          "destiny": 1,
          "position": 22,
          "office_title": "special",
          "for_buyer": true,
          "obligatory": false,
          "active": true,
          "system_name": null,
          "created_at": "2026-07-25T00:59:11.082+03:00",
          "updated_at": "2026-07-25T00:59:11.082+03:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": "special",
          "example": null,
          "type": "Field::TextField"
        },
        {
          "id": 1,
          "destiny": 1,
          "position": 0,
          "office_title": "Страна",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "country",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::Country"
        },
        {
          "id": 2,
          "destiny": 1,
          "position": 1,
          "office_title": "Регион",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "state",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 3,
          "destiny": 1,
          "position": 2,
          "office_title": "Город",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "city",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        },
        {
          "id": 5,
          "destiny": 1,
          "position": 4,
          "office_title": "Почтовый индекс",
          "for_buyer": true,
          "obligatory": false,
          "active": false,
          "system_name": "zip",
          "created_at": "2010-10-12T13:31:09.692+04:00",
          "updated_at": "2010-10-12T13:31:09.692+04:00",
          "show_in_result": true,
          "show_in_checkout": true,
          "is_indexed": null,
          "hide_in_backoffice": null,
          "handle": null,
          "placeholder": null,
          "title": null,
          "example": null,
          "type": "Field::SystemTextField"
        }
      ],
      "no_delivery": false,
      "kladr_autodetected_address": "",
      "country_options": [
        {
          "code": "RU",
          "title": "Россия",
          "selected": true
        }
      ],
      "address": "address",
      "country": "RU",
      "state": "state",
      "city": "London",
      "zip": "123456",
      "street": null,
      "house": null,
      "flat": null,
      "entrance": null,
      "doorphone": null,
      "floor": null,
      "comment": null,
      "kladr_json": null,
      "location": {
        "kladr_code": null,
        "zip": "123456",
        "kladr_zip": null,
        "region_zip": null,
        "country": "RU",
        "state": "state",
        "state_type": null,
        "area": null,
        "area_type": null,
        "city": "London",
        "city_type": null,
        "settlement": null,
        "settlement_type": null,
        "address": "address",
        "street": "address",
        "street_type": null,
        "house": null,
        "flat": null,
        "is_kladr": false,
        "latitude": null,
        "longitude": null,
        "autodetected": null
      }
    },
    "client": {
      "id": 2,
      "email": "some-one@yandex.ru",
      "name": "Mortimer Hant",
      "phone": "+79162128506",
      "created_at": "2009-03-30T04:37:32.000+04:00",
      "updated_at": "2026-07-25T00:59:11.000+03:00",
      "comment": null,
      "registered": false,
      "subscribe": true,
      "client_group_id": null,
      "bonus_points": 0,
      "type": "Client::Individual",
      "correspondent_account": null,
      "settlement_account": null,
      "consent_to_personal_data": null,
      "o_auth_provider": null,
      "messenger_subscription": null,
      "external_id": null,
      "external_service": null,
      "contact_name": "Mortimer Hant",
      "surname": null,
      "middlename": null,
      "progressive_discount": null,
      "group_discount": null,
      "ip_addr": "",
      "fields_values": [
        {
          "id": 259,
          "value": "old value",
          "created_at": "2026-07-25T00:59:11.070+03:00",
          "updated_at": "2026-07-25T00:59:11.070+03:00",
          "field_id": 448,
          "type": "Текст",
          "name": "special",
          "handle": null
        }
      ]
    },
    "discounts": [],
    "total_price": 1600.0,
    "paid_amount": 0.0,
    "items_price": 1500.0,
    "id": 1,
    "key": "teest",
    "number": 1001,
    "comment": "comment",
    "archived": false,
    "delivery_title": "самовывоз",
    "delivery_description": "самовывоз",
    "delivery_price": 100.0,
    "full_delivery_price": 100.0,
    "payment_description": null,
    "payment_title": "наличными курьеру",
    "first_referer": "http://google.com?q=my_query",
    "first_current_location": "/my_query",
    "first_query": null,
    "first_source_domain": null,
    "first_source": "Google",
    "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
    "current_location": "/",
    "query": null,
    "source_domain": null,
    "source": "Прямой трафик",
    "fulfillment_status": "new",
    "custom_status": {
      "permalink": "new",
      "title": "Новый"
    },
    "delivered_at": null,
    "accepted_at": null,
    "created_at": "2009-03-18T10:53:14.000+03:00",
    "updated_at": "2026-07-25T00:59:11.000+03:00",
    "financial_status": "pending",
    "delivery_date": null,
    "delivery_from_hour": null,
    "delivery_from_minutes": null,
    "delivery_to_hour": null,
    "delivery_to_minutes": null,
    "delivery_time": "",
    "paid_at": null,
    "delivery_variant_id": 1,
    "payment_gateway_id": 1,
    "margin": "0.0",
    "margin_amount": "0.0",
    "client_transaction_id": 1,
    "currency_code": "RUR",
    "cookies": null,
    "account_id": 1,
    "manager_comment": "manager comment1",
    "locale": "ru",
    "delivery_info": {
      "delivery_variant_id": null,
      "tariff_id": null,
      "title": null,
      "description": null,
      "price": null,
      "shipping_company": null,
      "shipping_company_handle": null,
      "delivery_interval": {
        "min_days": null,
        "max_days": null,
        "description": ""
      },
      "errors": [],
      "warnings": [],
      "outlet": {
        "id": null,
        "external_id": null,
        "latitude": null,
        "longitude": null,
        "type": null,
        "title": null,
        "description": null,
        "address": null,
        "payment_method": null,
        "source_id": null
      },
      "not_available": null
    },
    "responsible_user_id": null,
    "total_profit": "1500.0",
    "warehouse_id": null
  }
]
```

### Remove discount

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[discounts_attributes]` | да | discounts attributes |
| `order[discounts_attributes][][id]` | да | discount id |
| `order[discounts_attributes][][_destroy]` | да | destroy marker (1 or true) |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "discounts_attributes": [
      {
        "id": 123,
        "_destroy": true
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 803,
      "value": "old value",
      "created_at": "2026-07-25T00:59:24.688+03:00",
      "updated_at": "2026-07-25T00:59:24.688+03:00",
      "field_id": 992,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 27,
      "order_id": 1,
      "sale_price": 100.0,
      "full_sale_price": 100.0,
      "total_price": 200.0,
      "full_total_price": 200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": "10x20x30",
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt (Red)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:25.000+03:00",
      "created_at": "2026-07-25T00:59:24.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 28,
      "order_id": 1,
      "sale_price": 150.0,
      "full_sale_price": 150.0,
      "total_price": 150.0,
      "full_total_price": 150.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt (Зеленый / Синий)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:25.000+03:00",
      "created_at": "2026-07-25T00:59:24.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 795,
        "value": "old value",
        "created_at": "2026-07-25T00:59:24.648+03:00",
        "updated_at": "2026-07-25T00:59:24.648+03:00",
        "field_id": 984,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 984,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:24.636+03:00",
        "updated_at": "2026-07-25T00:59:24.636+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:24.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 787,
        "value": "old value",
        "created_at": "2026-07-25T00:59:24.621+03:00",
        "updated_at": "2026-07-25T00:59:24.621+03:00",
        "field_id": 976,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 450.0,
  "paid_amount": 0.0,
  "items_price": 350.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:25.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "170.0",
  "warehouse_id": null
}
```

### Remove order line

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[order_lines_attributes]` | да | order lines attributes |
| `order[order_lines_attributes][][id]` | да | order line id |
| `order[order_lines_attributes][][_destroy]` | да | destroy marker (1 or true) |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "order_lines_attributes": [
      {
        "id": 1,
        "_destroy": true
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 731,
      "value": "old value",
      "created_at": "2026-07-25T00:59:23.052+03:00",
      "updated_at": "2026-07-25T00:59:23.052+03:00",
      "field_id": 920,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 2,
      "order_id": 1,
      "sale_price": 600.0,
      "full_sale_price": 600.0,
      "total_price": 1200.0,
      "full_total_price": 1200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:23.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": 20,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 723,
        "value": "old value",
        "created_at": "2026-07-25T00:59:23.020+03:00",
        "updated_at": "2026-07-25T00:59:23.020+03:00",
        "field_id": 912,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 912,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:23.010+03:00",
        "updated_at": "2026-07-25T00:59:23.010+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:22.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 715,
        "value": "old value",
        "created_at": "2026-07-25T00:59:22.993+03:00",
        "updated_at": "2026-07-25T00:59:22.993+03:00",
        "field_id": 904,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 1300.0,
  "paid_amount": 0.0,
  "items_price": 1200.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:23.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "1200.0",
  "warehouse_id": null
}
```

### Update custom status

`PUT /admin/orders/:order_id.json`

This query updates custom status. Despite the fact that “order[fulfillment_status]” parameter is permitted, using both of the parameters is strongly discouraged. Please, keep in mind that if both parameters are transfered they MUST be compatible. Please, notice that you CAN’T update custom status using order[custom_status][permalink] parameter.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[custom_status_permalink]` | да | custom status |
| `order[fulfillment_status]` |  | fulfillment status |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "custom_status_permalink": "v-obrabotke"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 659,
      "value": "old value",
      "created_at": "2026-07-25T00:59:21.459+03:00",
      "updated_at": "2026-07-25T00:59:21.459+03:00",
      "field_id": 848,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 1,
      "order_id": 1,
      "sale_price": 300.0,
      "full_sale_price": 300.0,
      "total_price": 300.0,
      "full_total_price": 300.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": null,
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": [
        "code1",
        "code2"
      ],
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 2,
      "order_id": 1,
      "sale_price": 600.0,
      "full_sale_price": 600.0,
      "total_price": 1200.0,
      "full_total_price": 1200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": 20,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 651,
        "value": "old value",
        "created_at": "2026-07-25T00:59:21.421+03:00",
        "updated_at": "2026-07-25T00:59:21.421+03:00",
        "field_id": 840,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 840,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:21.409+03:00",
        "updated_at": "2026-07-25T00:59:21.409+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:21.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 643,
        "value": "old value",
        "created_at": "2026-07-25T00:59:21.394+03:00",
        "updated_at": "2026-07-25T00:59:21.394+03:00",
        "field_id": 832,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 1600.0,
  "paid_amount": 0.0,
  "items_price": 1500.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "accepted",
  "custom_status": {
    "permalink": "accepted",
    "title": "В обработке"
  },
  "delivered_at": null,
  "accepted_at": "2026-07-25T00:59:21.000+03:00",
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:21.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "1500.0",
  "warehouse_id": null
}
```

### Update discount

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[discounts_attributes]` | да | discounts attributes |
| `order[discounts_attributes][][id]` | да | discount id |
| `order[discounts_attributes][][discount]` |  | discount amount |
| `order[discounts_attributes][][type_id]` |  | 1 - percent (default), 2 - money |
| `order[discounts_attributes][][description]` |  | description text |
| `order[discounts_attributes][][discount_order_lines_ids]` |  | ids of affected order lines (all will be affected if not specified) |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "discounts_attributes": {
      "1": {
        "id": 139,
        "discount_order_lines_ids": [
          31
        ]
      }
    }
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 851,
      "value": "old value",
      "created_at": "2026-07-25T00:59:26.250+03:00",
      "updated_at": "2026-07-25T00:59:26.250+03:00",
      "field_id": 1040,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 31,
      "order_id": 1,
      "sale_price": 100.0,
      "full_sale_price": 90.0,
      "total_price": 200.0,
      "full_total_price": 180.0,
      "discounts_amount": 20.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": "10x20x30",
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt (Red)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:26.000+03:00",
      "created_at": "2026-07-25T00:59:26.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 32,
      "order_id": 1,
      "sale_price": 150.0,
      "full_sale_price": 150.0,
      "total_price": 150.0,
      "full_total_price": 150.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt (Зеленый / Синий)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:26.000+03:00",
      "created_at": "2026-07-25T00:59:26.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": {
    "id": 139,
    "description": "-20 (absolute)",
    "warning": null,
    "type_id": 2,
    "amount": 20.0,
    "full_amount": 20.0,
    "percent": null,
    "discount": "20.0",
    "reference_id": null,
    "reference_type": null,
    "discount_products_ids": [
      1
    ],
    "discount_order_lines_ids": [
      31
    ],
    "discount_code_id": null,
    "created_at": "2026-07-25T00:59:26.000+03:00",
    "updated_at": "2026-07-25T00:59:26.000+03:00"
  },
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 843,
        "value": "old value",
        "created_at": "2026-07-25T00:59:26.218+03:00",
        "updated_at": "2026-07-25T00:59:26.218+03:00",
        "field_id": 1032,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 1032,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:26.208+03:00",
        "updated_at": "2026-07-25T00:59:26.208+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:26.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 835,
        "value": "old value",
        "created_at": "2026-07-25T00:59:26.190+03:00",
        "updated_at": "2026-07-25T00:59:26.190+03:00",
        "field_id": 1024,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [
    {
      "id": 139,
      "description": "-20 (absolute)",
      "warning": null,
      "type_id": 2,
      "amount": 20.0,
      "full_amount": 20.0,
      "percent": null,
      "discount": "20.0",
      "reference_id": null,
      "reference_type": null,
      "discount_products_ids": [
        1
      ],
      "discount_order_lines_ids": [
        31
      ],
      "discount_code_id": null,
      "created_at": "2026-07-25T00:59:26.000+03:00",
      "updated_at": "2026-07-25T00:59:26.000+03:00"
    }
  ],
  "total_price": 430.0,
  "paid_amount": 0.0,
  "items_price": 330.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:26.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "150.0",
  "warehouse_id": null
}
```

### Update discount

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[discounts_attributes]` | да | discounts attributes |
| `order[discounts_attributes][][id]` | да | discount id |
| `order[discounts_attributes][][discount]` |  | discount amount |
| `order[discounts_attributes][][type_id]` |  | 1 - percent (default), 2 - money |
| `order[discounts_attributes][][description]` |  | description text |
| `order[discounts_attributes][][discount_order_lines_ids]` |  | ids of affected order lines (all will be affected if not specified) |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "discounts_attributes": {
      "1": {
        "id": 139,
        "discount_order_lines_ids": [
          31
        ]
      }
    }
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 851,
      "value": "old value",
      "created_at": "2026-07-25T00:59:26.250+03:00",
      "updated_at": "2026-07-25T00:59:26.250+03:00",
      "field_id": 1040,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 31,
      "order_id": 1,
      "sale_price": 100.0,
      "full_sale_price": 90.0,
      "total_price": 200.0,
      "full_total_price": 180.0,
      "discounts_amount": 20.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": "10x20x30",
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt (Red)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:26.000+03:00",
      "created_at": "2026-07-25T00:59:26.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 32,
      "order_id": 1,
      "sale_price": 150.0,
      "full_sale_price": 150.0,
      "total_price": 150.0,
      "full_total_price": 150.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt (Зеленый / Синий)",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:26.000+03:00",
      "created_at": "2026-07-25T00:59:26.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": {
    "id": 139,
    "description": "-20 (absolute)",
    "warning": null,
    "type_id": 2,
    "amount": 20.0,
    "full_amount": 20.0,
    "percent": null,
    "discount": "20.0",
    "reference_id": null,
    "reference_type": null,
    "discount_products_ids": [
      1
    ],
    "discount_order_lines_ids": [
      31
    ],
    "discount_code_id": null,
    "created_at": "2026-07-25T00:59:26.000+03:00",
    "updated_at": "2026-07-25T00:59:26.000+03:00"
  },
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 843,
        "value": "old value",
        "created_at": "2026-07-25T00:59:26.218+03:00",
        "updated_at": "2026-07-25T00:59:26.218+03:00",
        "field_id": 1032,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 1032,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:26.208+03:00",
        "updated_at": "2026-07-25T00:59:26.208+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:26.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 835,
        "value": "old value",
        "created_at": "2026-07-25T00:59:26.190+03:00",
        "updated_at": "2026-07-25T00:59:26.190+03:00",
        "field_id": 1024,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [
    {
      "id": 139,
      "description": "-20 (absolute)",
      "warning": null,
      "type_id": 2,
      "amount": 20.0,
      "full_amount": 20.0,
      "percent": null,
      "discount": "20.0",
      "reference_id": null,
      "reference_type": null,
      "discount_products_ids": [
        1
      ],
      "discount_order_lines_ids": [
        31
      ],
      "discount_code_id": null,
      "created_at": "2026-07-25T00:59:26.000+03:00",
      "updated_at": "2026-07-25T00:59:26.000+03:00"
    }
  ],
  "total_price": 430.0,
  "paid_amount": 0.0,
  "items_price": 330.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:26.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "150.0",
  "warehouse_id": null
}
```

### Update order

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[client_attributes]` |  | client info |
| `order[client_attributes][name]` |  | client name |
| `order[client_attributes][email]` |  | client email |
| `order[client_attributes][phone]` |  | client phone |
| `order[delivery_variant_id]` |  | delivery variant id |
| `order[payment_gateway_id]` |  | payment gateway id |
| `order[coupon]` |  | discount code |
| `order[comment]` |  | client comment to order |
| `order[manager_comment]` |  | manager comment to order |
| `order[financial_status]` |  | pending or paid |
| `order[fields_values_attributes]` |  | fields values array |
| `order[shipping_address_attributes]` |  | array of shipping address attributes |
| `order[discounts_attributes]` |  | array of discounts attributes |
| `order[warehouse_id]` |  | id of a warehouse to reserve products on |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "client_attributes": {
      "id": 2,
      "name": "Ivan",
      "email": "ivan@example.com",
      "phone": "79111113344",
      "fields_values_attributes": [
        {
          "field_id": 808,
          "value": "new!"
        }
      ]
    },
    "delivery_variant_id": 4,
    "payment_gateway_id": 4,
    "coupon": "20% discount",
    "comment": "Please call from 9 to 18",
    "manager_comment": "New client",
    "financial_status": "paid",
    "fields_values_attributes": [
      {
        "handle": "my_awesome_field",
        "value": "new value"
      }
    ],
    "shipping_address_attributes": {
      "fields_values_attributes": [
        {
          "field_id": 816,
          "value": "new value!"
        }
      ]
    },
    "warehouse_id": 50
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 635,
      "value": "new value",
      "created_at": "2026-07-25T00:59:20.807+03:00",
      "updated_at": "2026-07-25T00:59:21.147+03:00",
      "field_id": 824,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 1,
      "order_id": 1,
      "sale_price": 300.0,
      "full_sale_price": 240.0,
      "total_price": 300.0,
      "full_total_price": 240.0,
      "discounts_amount": 60.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": null,
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": [
        "code1",
        "code2"
      ],
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 2,
      "order_id": 1,
      "sale_price": 600.0,
      "full_sale_price": 480.0,
      "total_price": 1200.0,
      "full_total_price": 960.0,
      "discounts_amount": 240.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": 20,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": {
    "id": 91,
    "description": "Скидка по купону 20% discount",
    "warning": null,
    "type_id": 1,
    "amount": 300.0,
    "full_amount": 300.0,
    "percent": "20.0",
    "discount": "20.0",
    "reference_id": 252,
    "reference_type": "DiscountCode",
    "discount_products_ids": [],
    "discount_order_lines_ids": [],
    "discount_code_id": 252,
    "created_at": "2026-07-25T00:59:21.000+03:00",
    "updated_at": "2026-07-25T00:59:21.000+03:00"
  },
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 627,
        "value": "new value!",
        "created_at": "2026-07-25T00:59:20.776+03:00",
        "updated_at": "2026-07-25T00:59:21.158+03:00",
        "field_id": 816,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+74952128506",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": "",
    "full_delivery_address": "123456, state, London, address, new value!",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 816,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:20.765+03:00",
        "updated_at": "2026-07-25T00:59:20.765+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "ivan@example.com",
    "name": "Ivan",
    "phone": "79111113344",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:21.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Ivan",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 619,
        "value": "new!",
        "created_at": "2026-07-25T00:59:20.752+03:00",
        "updated_at": "2026-07-25T00:59:21.075+03:00",
        "field_id": 808,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [
    {
      "id": 91,
      "description": "Скидка по купону 20% discount",
      "warning": null,
      "type_id": 1,
      "amount": 300.0,
      "full_amount": 300.0,
      "percent": "20.0",
      "discount": "20.0",
      "reference_id": 252,
      "reference_type": "DiscountCode",
      "discount_products_ids": [],
      "discount_order_lines_ids": [],
      "discount_code_id": 252,
      "created_at": "2026-07-25T00:59:21.000+03:00",
      "updated_at": "2026-07-25T00:59:21.000+03:00"
    }
  ],
  "total_price": 1200.0,
  "paid_amount": 1500.0,
  "items_price": 1200.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "Please call from 9 to 18",
  "archived": false,
  "delivery_title": "Внешний",
  "delivery_description": "Внешний (Внешний способ)",
  "delivery_price": 0.0,
  "full_delivery_price": 0.0,
  "payment_description": null,
  "payment_title": "Выставление счета для юр. лиц 2.",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:21.000+03:00",
  "financial_status": "paid",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": "2026-07-25T00:59:20.000+03:00",
  "delivery_variant_id": 4,
  "payment_gateway_id": 4,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "New client",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": 4,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": 0.0,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": [],
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "1200.0",
  "warehouse_id": 50
}
```

### Update order line

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[order_lines_attributes]` | да | order lines attributes |
| `order[order_lines_attributes][][id]` | да | order line id |
| `order[order_lines_attributes][][sale_price]` |  | order line sale price |
| `order[order_lines_attributes][][quantity]` |  | order line quantity |
| `order[order_lines_attributes][][comment]` |  | order line comment |
| `order[order_lines_attributes][][vat]` |  | integer, VAT for order line is equal to VAT of its product, by default, products get their value from account by default,if default in account set as 'No VAT', you need to use 'force' to change VAT of order line through API; possible values depend on country, set to -1 if this product is not a subject to VAT |
| `order[order_lines_attributes][][marking_codes]` |  | array of marking codes  array |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "order_lines_attributes": [
      {
        "id": 1,
        "sale_price": 100,
        "quantity": 2,
        "comment": "New comment",
        "vat": 0,
        "marking_codes": [
          "43c4ae4d6283b2d857ec",
          "14e10d3597d9f53392eb"
        ]
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 755,
      "value": "old value",
      "created_at": "2026-07-25T00:59:23.504+03:00",
      "updated_at": "2026-07-25T00:59:23.504+03:00",
      "field_id": 944,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 1,
      "order_id": 1,
      "sale_price": 100.0,
      "full_sale_price": 100.0,
      "total_price": 200.0,
      "full_total_price": 200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": null,
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt",
      "unit": "pce",
      "comment": "New comment",
      "updated_at": "2026-07-25T00:59:23.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": 0,
      "fiscal_product_type": 1,
      "requires_marking": true,
      "marking_codes": [
        "43c4ae4d6283b2d857ec",
        "14e10d3597d9f53392eb"
      ],
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 2,
      "order_id": 1,
      "sale_price": 600.0,
      "full_sale_price": 600.0,
      "total_price": 1200.0,
      "full_total_price": 1200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2026-07-25T00:59:23.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": 20,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 747,
        "value": "old value",
        "created_at": "2026-07-25T00:59:23.469+03:00",
        "updated_at": "2026-07-25T00:59:23.469+03:00",
        "field_id": 936,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+7(495)212-85-06",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": null,
    "full_delivery_address": "123456, state, London, address, old value",
    "address_for_gis": "123456, state, London, address",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 936,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:23.457+03:00",
        "updated_at": "2026-07-25T00:59:23.457+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "address",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "address",
      "street": "address",
      "street_type": null,
      "house": null,
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:23.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 739,
        "value": "old value",
        "created_at": "2026-07-25T00:59:23.443+03:00",
        "updated_at": "2026-07-25T00:59:23.443+03:00",
        "field_id": 928,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 1500.0,
  "paid_amount": 0.0,
  "items_price": 1400.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:23.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "1400.0",
  "warehouse_id": null
}
```

### Update shipping address

`PUT /admin/orders/:order_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[shipping_address_attributes]` | да | array of shipping address attributes |
| `order[shipping_address_attributes][address]` | да | shipping address |

**Запрос**

```http
PUT /admin/orders/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "order": {
    "shipping_address_attributes": {
      "address": "Moscow, Znamenka 12"
    }
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "fields_values": [
    {
      "id": 875,
      "value": "old value",
      "created_at": "2026-07-25T00:59:27.001+03:00",
      "updated_at": "2026-07-25T00:59:27.001+03:00",
      "field_id": 1064,
      "type": "Текст",
      "name": "special",
      "handle": "my_awesome_field"
    }
  ],
  "order_lines": [
    {
      "id": 1,
      "order_id": 1,
      "sale_price": 300.0,
      "full_sale_price": 300.0,
      "total_price": 300.0,
      "full_total_price": 300.0,
      "discounts_amount": 0.0,
      "quantity": 1,
      "reserved_quantity": null,
      "weight": null,
      "dimensions": null,
      "variant_id": 1,
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "title": "Cool t-shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": -1,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": [
        "code1",
        "code2"
      ],
      "accessory_lines": [],
      "external_variant_id": null
    },
    {
      "id": 2,
      "order_id": 1,
      "sale_price": 600.0,
      "full_sale_price": 600.0,
      "total_price": 1200.0,
      "full_total_price": 1200.0,
      "discounts_amount": 0.0,
      "quantity": 2,
      "reserved_quantity": null,
      "weight": "0.3",
      "dimensions": null,
      "variant_id": 2,
      "product_id": 2,
      "sku": "424",
      "barcode": null,
      "title": "shirt",
      "unit": "pce",
      "comment": null,
      "updated_at": "2009-03-18T10:53:15.000+03:00",
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "bundle_id": null,
      "vat": 20,
      "fiscal_product_type": 1,
      "requires_marking": null,
      "marking_codes": null,
      "accessory_lines": [],
      "external_variant_id": null
    }
  ],
  "order_changes": [],
  "discount": null,
  "shipping_address": {
    "id": 1,
    "fields_values": [
      {
        "id": 867,
        "value": "old value",
        "created_at": "2026-07-25T00:59:26.967+03:00",
        "updated_at": "2026-07-25T00:59:26.967+03:00",
        "field_id": 1056,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ],
    "name": "Tom Dellay",
    "surname": null,
    "middlename": null,
    "phone": "+74952128506",
    "formatted_phone": "+7(495)212-85-06",
    "full_name": "Tom Dellay",
    "full_locality_name": "",
    "full_delivery_address": "123456, state, London, Moscow, Znamenka 12, old value",
    "address_for_gis": "123456, state, London, Moscow, Znamenka 12",
    "location_valid": true,
    "recipient_fields": [
      {
        "id": 11,
        "destiny": 2,
        "position": 1,
        "office_title": "Имя",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "name",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::ObligatoryTextField"
      },
      {
        "id": 12,
        "destiny": 2,
        "position": 2,
        "office_title": "Телефон",
        "for_buyer": true,
        "obligatory": true,
        "active": true,
        "system_name": "phone",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Phone"
      }
    ],
    "backoffice_fields": [
      {
        "id": 4,
        "destiny": 1,
        "position": 3,
        "office_title": "Адрес",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": "address",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextArea"
      },
      {
        "id": 1056,
        "destiny": 1,
        "position": 22,
        "office_title": "special",
        "for_buyer": true,
        "obligatory": false,
        "active": true,
        "system_name": null,
        "created_at": "2026-07-25T00:59:26.955+03:00",
        "updated_at": "2026-07-25T00:59:26.955+03:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": "special",
        "example": null,
        "type": "Field::TextField"
      },
      {
        "id": 1,
        "destiny": 1,
        "position": 0,
        "office_title": "Страна",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "country",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::Country"
      },
      {
        "id": 2,
        "destiny": 1,
        "position": 1,
        "office_title": "Регион",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "state",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 3,
        "destiny": 1,
        "position": 2,
        "office_title": "Город",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "city",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      },
      {
        "id": 5,
        "destiny": 1,
        "position": 4,
        "office_title": "Почтовый индекс",
        "for_buyer": true,
        "obligatory": false,
        "active": false,
        "system_name": "zip",
        "created_at": "2010-10-12T13:31:09.692+04:00",
        "updated_at": "2010-10-12T13:31:09.692+04:00",
        "show_in_result": true,
        "show_in_checkout": true,
        "is_indexed": null,
        "hide_in_backoffice": null,
        "handle": null,
        "placeholder": null,
        "title": null,
        "example": null,
        "type": "Field::SystemTextField"
      }
    ],
    "no_delivery": false,
    "kladr_autodetected_address": "",
    "country_options": [
      {
        "code": "RU",
        "title": "Россия",
        "selected": true
      }
    ],
    "address": "Moscow, Znamenka 12",
    "country": "RU",
    "state": "state",
    "city": "London",
    "zip": "123456",
    "street": null,
    "house": null,
    "flat": null,
    "entrance": null,
    "doorphone": null,
    "floor": null,
    "comment": null,
    "kladr_json": null,
    "location": {
      "kladr_code": null,
      "zip": "123456",
      "kladr_zip": null,
      "region_zip": null,
      "country": "RU",
      "state": "state",
      "state_type": null,
      "area": null,
      "area_type": null,
      "city": "London",
      "city_type": null,
      "settlement": null,
      "settlement_type": null,
      "address": "Moscow, Znamenka 12",
      "street": "Moscow",
      "street_type": null,
      "house": "Znamenka 12",
      "flat": null,
      "is_kladr": false,
      "latitude": null,
      "longitude": null,
      "autodetected": null
    }
  },
  "client": {
    "id": 2,
    "email": "some-one@yandex.ru",
    "name": "Mortimer Hant",
    "phone": "+79162128506",
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2026-07-25T00:59:26.000+03:00",
    "comment": null,
    "registered": false,
    "subscribe": true,
    "client_group_id": null,
    "bonus_points": 0,
    "type": "Client::Individual",
    "correspondent_account": null,
    "settlement_account": null,
    "consent_to_personal_data": null,
    "o_auth_provider": null,
    "messenger_subscription": null,
    "external_id": null,
    "external_service": null,
    "contact_name": "Mortimer Hant",
    "surname": null,
    "middlename": null,
    "progressive_discount": null,
    "group_discount": null,
    "ip_addr": "",
    "fields_values": [
      {
        "id": 859,
        "value": "old value",
        "created_at": "2026-07-25T00:59:26.940+03:00",
        "updated_at": "2026-07-25T00:59:26.940+03:00",
        "field_id": 1048,
        "type": "Текст",
        "name": "special",
        "handle": null
      }
    ]
  },
  "discounts": [],
  "total_price": 1600.0,
  "paid_amount": 0.0,
  "items_price": 1500.0,
  "id": 1,
  "key": "teest",
  "number": 1001,
  "comment": "comment",
  "archived": false,
  "delivery_title": "самовывоз",
  "delivery_description": "самовывоз",
  "delivery_price": 100.0,
  "full_delivery_price": 100.0,
  "payment_description": null,
  "payment_title": "наличными курьеру",
  "first_referer": "http://google.com?q=my_query",
  "first_current_location": "/my_query",
  "first_query": null,
  "first_source_domain": null,
  "first_source": "Google",
  "referer": "http://yandex.ru/yandsearch?text=%D1%81%D0%BF%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5+%D0%BC%D0%B5%D1%88%D0%BA%D0%B8&stpar2=%2Fh1%2Ftm7%2Fs1&stpar4=%2Fs1&stpar1=%2Fu1&stpar3=%2Fm3%2Ftc14%2Fnc1",
  "current_location": "/",
  "query": null,
  "source_domain": null,
  "source": "Прямой трафик",
  "fulfillment_status": "new",
  "custom_status": {
    "permalink": "new",
    "title": "Новый"
  },
  "delivered_at": null,
  "accepted_at": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T00:59:27.000+03:00",
  "financial_status": "pending",
  "delivery_date": null,
  "delivery_from_hour": null,
  "delivery_from_minutes": null,
  "delivery_to_hour": null,
  "delivery_to_minutes": null,
  "delivery_time": "",
  "paid_at": null,
  "delivery_variant_id": 1,
  "payment_gateway_id": 1,
  "margin": "0.0",
  "margin_amount": "0.0",
  "client_transaction_id": 1,
  "currency_code": "RUR",
  "cookies": null,
  "account_id": 1,
  "manager_comment": "manager comment1",
  "locale": "ru",
  "delivery_info": {
    "delivery_variant_id": null,
    "tariff_id": null,
    "title": null,
    "description": null,
    "price": null,
    "shipping_company": null,
    "shipping_company_handle": null,
    "delivery_interval": {
      "min_days": null,
      "max_days": null,
      "description": ""
    },
    "errors": [],
    "warnings": [],
    "outlet": {
      "id": null,
      "external_id": null,
      "latitude": null,
      "longitude": null,
      "type": null,
      "title": null,
      "description": null,
      "address": null,
      "payment_method": null,
      "source_id": null
    },
    "not_available": null
  },
  "responsible_user_id": null,
  "total_profit": "1500.0",
  "warehouse_id": null
}
```

---

## Outlet

| Операция | Метод и путь |
|----------|--------------|
| [Create outlet](#create-outlet) | `POST /admin/outlets.json` |
| [Destroy outlet](#destroy-outlet) | `DELETE /admin/outlets/:outlet_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/outlets/:outlet_id.json` |
| [Get outlets](#get-outlets) | `GET /admin/outlets.json` |
| [Update outlet](#update-outlet) | `PUT /admin/outlets/:outlet_id.json` |

### Create outlet

`POST /admin/outlets.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `outlet[title]` | да | title |
| `outlet[description]` |  | description |
| `outlet[address]` | да | address in text way |
| `outlet[price]` | да | delivery price |
| `outlet[latitude]` | да | point latitude |
| `outlet[longitude]` | да | point longitude |

**Запрос**

```http
POST /admin/outlets.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "outlet": {
    "title": "Point 2",
    "description": "Point description",
    "address": "some address",
    "price": 300,
    "latitude": 55.7558,
    "longitude": 37.6173
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 8,
  "account_id": 1,
  "delivery_variant_id": 161,
  "price": "300.0",
  "title_translations": [
    "Point 2"
  ],
  "description_translations": [
    "Point description"
  ],
  "address": "some address",
  "latitude": "55.7558",
  "longitude": "37.6173",
  "created_at": "2026-07-25T00:58:56.627+03:00",
  "updated_at": "2026-07-25T00:58:56.627+03:00",
  "title": "Point 2",
  "description": "Point description"
}
```

### Destroy outlet

`DELETE /admin/outlets/:outlet_id.json`

**Запрос**

```http
DELETE /admin/outlets/12.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/outlets/:outlet_id.json`

Content-Type: application/json Get outlet

**Запрос**

```http
GET /admin/outlets/4.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 4,
  "account_id": 1,
  "delivery_variant_id": 145,
  "price": "100.0",
  "title_translations": [
    "Самовывоз от метро"
  ],
  "description_translations": [
    "Вход через арку во двор"
  ],
  "address": "ул Хмурого д1",
  "latitude": "55.7558",
  "longitude": "37.6173",
  "created_at": "2026-07-25T00:58:56.296+03:00",
  "updated_at": "2026-07-25T00:58:56.296+03:00",
  "title": "Самовывоз от метро",
  "description": "Вход через арку во двор"
}
```

### Get outlets

`GET /admin/outlets.json`

**Запрос**

```http
GET /admin/outlets.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 2,
    "account_id": 1,
    "delivery_variant_id": 129,
    "price": "100.0",
    "title_translations": [
      "Самовывоз от метро"
    ],
    "description_translations": [
      "Вход через арку во двор"
    ],
    "address": "ул Хмурого д1",
    "latitude": "55.7558",
    "longitude": "37.6173",
    "created_at": "2026-07-25T00:58:56.038+03:00",
    "updated_at": "2026-07-25T00:58:56.038+03:00",
    "title": "Самовывоз от метро",
    "description": "Вход через арку во двор"
  }
]
```

### Update outlet

`PUT /admin/outlets/:outlet_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `outlet[title]` | да | title |
| `outlet[description]` |  | description |
| `outlet[address]` | да | address in text way |
| `outlet[price]` | да | delivery price |
| `outlet[latitude]` | да | point latitude |
| `outlet[longitude]` | да | point longitude |

**Запрос**

```http
PUT /admin/outlets/10.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "outlet": {
    "title": "New point title",
    "description": "New point description",
    "address": "some new address",
    "price": 300,
    "latitude": 55.7558,
    "longitude": 37.6173
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "delivery_variant_id": 177,
  "title": "New point title",
  "description": "New point description",
  "address": "some new address",
  "price": "300.0",
  "latitude": "55.7558",
  "longitude": "37.6173",
  "title_translations": [
    "New point title"
  ],
  "description_translations": [
    "New point description"
  ],
  "id": 10,
  "account_id": 1,
  "created_at": "2026-07-25T00:58:56.817+03:00",
  "updated_at": "2026-07-25T00:58:56.887+03:00"
}
```

---

## Page

| Операция | Метод и путь |
|----------|--------------|
| [Create page](#create-page) | `POST /admin/pages.json` |
| [Delete page](#delete-page) | `DELETE /admin/pages/:page_id.json` |
| [Get page](#get-page) | `GET /admin/pages/:page_id.json` |
| [Get pages list](#get-pages-list) | `GET /admin/pages.json` |
| [Update page](#update-page) | `PUT /admin/pages/:page_id.json` |

### Create page

`POST /admin/pages.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `page[title]` | да | page title |
| `page[content]` | да | page content |
| `page[permalink]` |  | permalink on /page/* |
| `page[html_title]` |  | html title tag content |
| `page[meta_keywords]` |  | html meta words on page |
| `page[meta_description]` |  | html meta description |
| `page[created_at]` |  | publication date |

**Запрос**

```http
POST /admin/pages.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "page": {
    "title": "new page",
    "content": "<p>content of my page<p>",
    "permalink": "first_page",
    "html_title": "HTML TAG TITLE",
    "meta_keywords": "new cool products",
    "meta_description": "It's a page about our new products",
    "created_at": "2014-09-09"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 20,
  "meta_keywords": "new cool products",
  "meta_description": "It's a page about our new products",
  "is_hidden": false,
  "kind": "only_page",
  "permalink": "first_page",
  "custom_template": null,
  "created_at": "2014-09-09T00:00:00.000+04:00",
  "updated_at": "2026-07-25T01:00:31.000+03:00",
  "is_main": false,
  "title": "new page",
  "html_title": "HTML TAG TITLE",
  "content": "<p>content of my page<p>"
}
```

### Delete page

`DELETE /admin/pages/:page_id.json`

**Запрос**

```http
DELETE /admin/pages/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
Content-Type: application/x-www-form-urlencoded
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get page

`GET /admin/pages/:page_id.json`

**Запрос**

```http
GET /admin/pages/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "meta_keywords": "meta keywords",
  "meta_description": "meta description",
  "is_hidden": false,
  "kind": "only_page",
  "permalink": "страница",
  "custom_template": null,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2009-03-18T10:55:15.000+03:00",
  "is_main": true,
  "title": "any_url",
  "html_title": "html title",
  "content": "abc"
}
```

### Get pages list

`GET /admin/pages.json`

**Запрос**

```http
GET /admin/pages.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "meta_keywords": "meta keywords",
    "meta_description": "meta description",
    "is_hidden": false,
    "kind": "only_page",
    "permalink": "страница",
    "custom_template": null,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2009-03-18T10:55:15.000+03:00",
    "is_main": true,
    "title": "any_url",
    "html_title": "html title",
    "content": "abc"
  },
  {
    "id": 3,
    "meta_keywords": "meta keywords",
    "meta_description": "meta description",
    "is_hidden": false,
    "kind": "only_page",
    "permalink": "страница2",
    "custom_template": null,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2009-03-17T10:55:15.000+03:00",
    "is_main": false,
    "title": "bany_url",
    "html_title": "html title",
    "content": "abc"
  }
]
```

### Update page

`PUT /admin/pages/:page_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `page[title]` |  | page title |
| `page[content]` |  | page content |
| `page[permalink]` |  | permalink on blog/* |
| `page[html_title]` |  | html title tag content |
| `page[meta_keywords]` |  | html meta words on page |
| `page[meta_description]` |  | html meta description |
| `page[created_at]` |  | publication date |

**Запрос**

```http
PUT /admin/pages/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "page": {
    "title": "new title"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

---

## PaymentGateway

| Операция | Метод и путь |
|----------|--------------|
| [Create Cod or Custom payment_gateway](#create-cod-or-custom-payment-gateway) | `POST /admin/payment_gateways.json` |
| [Create External payment_gateway](#create-external-payment-gateway) | `POST /admin/payment_gateways.json` |
| [Destroy payment_gateway](#destroy-payment-gateway) | `DELETE /admin/payment_gateways/:payment_gateway_id_1.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/payment_gateways/:payment_gateway_id_1.json` |
| [Get payment_gateways](#get-payment-gateways) | `GET /admin/payment_gateways.json` |
| [Update payment_gateway](#update-payment-gateway) | `PUT /admin/payment_gateways/:payment_gateway_id_1.json` |

### Create Cod or Custom payment_gateway

`POST /admin/payment_gateways.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `payment_gateway[title]` | да | title |
| `payment_gateway[margin]` |  | margin measured in percentage and applied to order price with current payment method; it should fall in the range from -100 to 1000 |
| `payment_gateway[position]` |  | position in payment gateways list |
| `payment_gateway[type]` | да | available types: PaymentGateway::Cod, PaymentGateway::Custom |
| `payment_gateway[description]` |  | description |
| `payment_gateway[payment_delivery_variants_attributes]` |  | array of available delivery variants for this payment gateway; payment gateway is unavailable if list is empty |
| `payment_gateway[payment_delivery_variants_attributes][delivery_variant_id]` |  | id of delivery variant |
| `payment_gateway[add_delivery_variants]` |  | set all delivery variants available for this payment gateway |

**Запрос**

```http
POST /admin/payment_gateways.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "payment_gateway": {
    "title": "cash on delivery",
    "margin": 10,
    "type": "PaymentGateway::Cod",
    "description": "text",
    "payment_delivery_variants_attributes": [
      {
        "delivery_variant_id": 1
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 35,
  "type": "PaymentGateway::Cod",
  "position": 4,
  "created_at": "2026-07-25T00:58:58.000+03:00",
  "updated_at": "2026-07-25T00:58:58.000+03:00",
  "margin": "10.0",
  "available_for_individual_clients": true,
  "available_for_juridical_clients": true,
  "title": "cash on delivery",
  "description": "text",
  "payment_delivery_variants": [
    {
      "created_at": "2026-07-25T00:58:58.000+03:00",
      "delivery_variant_id": 1,
      "id": 58
    }
  ]
}
```

### Create External payment_gateway

`POST /admin/payment_gateways.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `payment_gateway[title]` | да | title |
| `payment_gateway[margin]` |  | margin measured in percentage and added to order price with current payment method; it should fall in the range from -100 to 1000 |
| `payment_gateway[position]` |  | position in payment gateways list |
| `payment_gateway[type]` | да | PaymentGateway::External |
| `payment_gateway[description]` |  | description |
| `payment_gateway[payment_delivery_variants_attributes]` |  | array of available delivery variants for this payment gateway; payment gateway is unavailable if list is empty |
| `payment_gateway[payment_delivery_variants_attributes][delivery_variant_id]` |  | id of delivery variant |
| `payment_gateway[url]` | да | url of external payment gateway |
| `payment_gateway[shop_id]` | да | shop id - store identificator for external systems |
| `payment_gateway[send_order]` |  | set true to send order_lines information to the "url"; parameter name - "order_json"; do not forget add "order_json" to signature validation to the end of the string. |
| `payment_gateway[add_delivery_variants]` |  | set all delivery variants available for this payment gateway |
| `payment_gateway[redirect_to_payment]` |  | true by default; if true automatically redirect to payment after successful order creation; |
| `payment_gateway[convert_currency]` |  | iso code, if you would like to convert order amount from site currency to another currency |

**Запрос**

```http
POST /admin/payment_gateways.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "payment_gateway": {
    "title": "external payment gateway",
    "margin": 10,
    "type": "PaymentGateway::External",
    "description": "text",
    "payment_delivery_variants_attributes": [
      {
        "delivery_variant_id": 1
      }
    ],
    "url": "payment-gateway.ru/transactions/create",
    "shop_id": "123",
    "redirect_to_payment": false,
    "convert_currency": "LVL"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 37,
  "type": "PaymentGateway::External",
  "position": 4,
  "created_at": "2026-07-25T00:58:58.000+03:00",
  "updated_at": "2026-07-25T00:58:58.000+03:00",
  "margin": "10.0",
  "available_for_individual_clients": true,
  "available_for_juridical_clients": true,
  "title": "external payment gateway",
  "description": "text",
  "password": "c7b0c8fc95369ad5",
  "url": "payment-gateway.ru/transactions/create",
  "shop_id": "123",
  "send_order": null,
  "convert_currency": "LVL",
  "widget_mode": null,
  "window_mode": null,
  "widget_html_code": null,
  "show_at_product": null,
  "show_at_collection": null,
  "show_at_cart": null,
  "design_version": null,
  "text_version": null,
  "rounding_version": null,
  "show_payment_fee": null,
  "info_page_id": null,
  "success_url": "http://quentin.insales.io:3002/payments/external/37/success",
  "fail_url": "http://quentin.insales.io:3002/payments/external/37/fail",
  "server_url": "http://quentin.insales.io:3002/payments/external/37/server",
  "redirect_to_payment": false,
  "payment_delivery_variants": [
    {
      "created_at": "2026-07-25T00:58:58.000+03:00",
      "delivery_variant_id": 1,
      "id": 60
    }
  ]
}
```

### Destroy payment_gateway

`DELETE /admin/payment_gateways/:payment_gateway_id_1.json`

**Запрос**

```http
DELETE /admin/payment_gateways/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/payment_gateways/:payment_gateway_id_1.json`

Content-Type: application/json Get payment_gateway

**Запрос**

```http
GET /admin/payment_gateways/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "type": "PaymentGateway::Cod",
  "position": 1,
  "created_at": "2009-03-30T04:37:32.000+04:00",
  "updated_at": "2009-03-31T12:20:34.000+04:00",
  "margin": "0.0",
  "available_for_individual_clients": true,
  "available_for_juridical_clients": true,
  "title": "Наличными курьеру",
  "description": "",
  "payment_delivery_variants": [
    {
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "delivery_variant_id": 1,
      "id": 1
    },
    {
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "delivery_variant_id": 2,
      "id": 2
    }
  ]
}
```

### Get payment_gateways

`GET /admin/payment_gateways.json`

**Запрос**

```http
GET /admin/payment_gateways.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "type": "PaymentGateway::Cod",
    "position": 1,
    "created_at": "2009-03-30T04:37:32.000+04:00",
    "updated_at": "2009-03-31T12:20:34.000+04:00",
    "margin": "0.0",
    "available_for_individual_clients": true,
    "available_for_juridical_clients": true,
    "title": "Наличными курьеру",
    "description": "",
    "payment_delivery_variants": [
      {
        "created_at": "2026-07-25T00:58:09.000+03:00",
        "delivery_variant_id": 1,
        "id": 1
      },
      {
        "created_at": "2026-07-25T00:58:09.000+03:00",
        "delivery_variant_id": 2,
        "id": 2
      }
    ]
  },
  {
    "id": 3,
    "type": "PaymentGateway::JuridicalPersonBill",
    "position": 2,
    "created_at": "2009-03-17T10:55:14.000+03:00",
    "updated_at": "2009-03-17T10:55:15.000+03:00",
    "margin": "20.0",
    "available_for_individual_clients": true,
    "available_for_juridical_clients": true,
    "title": "Выставление счета для юр. лиц.",
    "description": "",
    "redirect_to_payment": true,
    "payment_delivery_variants": []
  },
  {
    "id": 4,
    "type": "PaymentGateway::JuridicalPersonBill",
    "position": 3,
    "created_at": "2009-03-17T10:55:14.000+03:00",
    "updated_at": "2009-03-17T10:55:15.000+03:00",
    "margin": "0.0",
    "available_for_individual_clients": true,
    "available_for_juridical_clients": true,
    "title": "Выставление счета для юр. лиц 2.",
    "description": "",
    "redirect_to_payment": true,
    "payment_delivery_variants": []
  }
]
```

### Update payment_gateway

`PUT /admin/payment_gateways/:payment_gateway_id_1.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `payment_gateway[title]` |  | title |
| `payment_gateway[margin]` |  | margin measured in percentage and added to order price with current payment method; it should fall in the range from -100 to 1000 |
| `payment_gateway[position]` |  | position in payment gateways list |
| `payment_gateway[description]` |  | description |
| `payment_gateway[payment_delivery_variants_attributes]` |  | array of available delivery variants for this payment gateway. payment gateway is unavailable if list is empty |
| `payment_gateway[payment_delivery_variants_attributes][_destroy]` |  | set it 1 to delete delivery variant for this payment gateway |
| `delivery_variant[payment_delivery_variants_attributes][id]` |  | only for removal: specify id of payment gateway and delivery variant binding you want to delete |
| `delivery_variant[payment_delivery_variants_attributes][delivery_variant_id]` |  | for creation: specify id of delivery variant you want to bind with this payment gateway |

**Запрос**

```http
PUT /admin/payment_gateways/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "payment_gateway": {
    "title": "cash",
    "margin": 10,
    "description": "description",
    "payment_delivery_variants_attributes": [
      {
        "_destroy": 1,
        "id": 2
      },
      {
        "delivery_variant_id": 2
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "type": "PaymentGateway::Cod",
  "position": 1,
  "created_at": "2009-03-30T04:37:32.000+04:00",
  "updated_at": "2026-07-25T00:58:58.000+03:00",
  "margin": "10.0",
  "available_for_individual_clients": true,
  "available_for_juridical_clients": true,
  "title": "cash",
  "description": "description",
  "payment_delivery_variants": [
    {
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "delivery_variant_id": 1,
      "id": 1
    },
    {
      "created_at": "2026-07-25T00:58:09.000+03:00",
      "delivery_variant_id": 2,
      "id": 2
    }
  ]
}
```

---

## PriceKind

| Операция | Метод и путь |
|----------|--------------|
| [Create price kind](#create-price-kind) | `POST /admin/price_kinds.json` |
| [Destroy price kind](#destroy-price-kind) | `DELETE /admin/price_kinds/:price_kind_id.json` |
| [Get price kind](#get-price-kind) | `GET /admin/price_kinds/:price_kind_id.json` |
| [Get price kinds](#get-price-kinds) | `GET /admin/price_kinds.json` |
| [Update price kind](#update-price-kind) | `PUT /admin/price_kinds/:price_kind_id.json` |

### Create price kind

`POST /admin/price_kinds.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `price_kind[title]` | да | title |
| `price_kind[show_in_lists]` |  | show price in products list in back-office |

**Запрос**

```http
POST /admin/price_kinds.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "price_kind": {
    "title": "large quantity price"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 175,
  "price_index": 2,
  "title": "large quantity price",
  "show_in_lists": true,
  "variant_price_method_name": "price3"
}
```

### Destroy price kind

`DELETE /admin/price_kinds/:price_kind_id.json`

**Запрос**

```http
DELETE /admin/price_kinds/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get price kind

`GET /admin/price_kinds/:price_kind_id.json`

**Запрос**

```http
GET /admin/price_kinds/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "price_index": 1,
  "title": "от 100",
  "show_in_lists": true,
  "variant_price_method_name": "price2"
}
```

### Get price kinds

`GET /admin/price_kinds.json`

**Запрос**

```http
GET /admin/price_kinds.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "price_index": 1,
    "title": "от 100",
    "show_in_lists": true,
    "variant_price_method_name": "price2"
  }
]
```

### Update price kind

`PUT /admin/price_kinds/:price_kind_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `price_kind[title]` | да | title |
| `price_kind[show_in_lists]` |  | show price in products list in back-office |

**Запрос**

```http
PUT /admin/price_kinds/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "price_kind": {
    "title": "price for 10+ quantity"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "price_index": 1,
  "title": "price for 10+ quantity",
  "show_in_lists": true,
  "variant_price_method_name": "price2"
}
```

---

## Product

| Операция | Метод и путь |
|----------|--------------|
| [Create bundle](#create-bundle) | `POST /admin/products.json` |
| [Create product](#create-product) | `POST /admin/products.json` |
| [Destroy product](#destroy-product) | `DELETE /admin/products/:product_id.json` |
| [Get product](#get-product) | `GET /admin/products/:product_id.json` |
| [Get products count](#get-products-count) | `GET /admin/products/count.json` |
| [Get products](#get-products) | `GET /admin/products.json` |
| [Update bundle](#update-bundle) | `PUT /admin/products/4.json` |
| [Update product](#update-product) | `PUT /admin/products/:product_id.json` |

### Create bundle

`POST /admin/products.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product[category_id]` | да | category id |
| `product[title]` | да | title |
| `product[description]` |  | description |
| `product[short_description]` |  | short description |
| `product[sort_weight]` |  | sort weight |
| `product[bundle]` | да | bundle |
| `product[variants_attributes]` | да | variants attributes |
| `product[vat]` |  | do not set VAT for bundle because order lines will be created with its components, which have their own VAT |
| `product[product_bundle_components_attributes]` | да | product bundle components attributes |

**Запрос**

```http
POST /admin/products.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product": {
    "category_id": 1,
    "title": "Van Gogh Ruled Peach Notebook",
    "bundle": true,
    "variants_attributes": [
      {
        "price": 100
      }
    ],
    "product_bundle_components_attributes": [
      {
        "variant_id": 1,
        "quantity": 2,
        "free": false
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 293977361,
  "category_id": 1,
  "created_at": "2026-07-25T00:59:56.000+03:00",
  "updated_at": "2026-07-25T00:59:56.000+03:00",
  "is_hidden": false,
  "available": null,
  "canonical_url_collection_id": null,
  "archived": false,
  "unit": "pce",
  "sort_weight": null,
  "bundle": true,
  "ignore_discounts": null,
  "vat": 10,
  "dimensions": null,
  "fiscal_product_type": 1,
  "custom_template": null,
  "title": "Van Gogh Ruled Peach Notebook",
  "short_description": null,
  "permalink": "van-gogh-ruled-peach-notebook",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "currency_code": "LVL",
  "collections_ids": [],
  "sales_channels_id": null,
  "description": null,
  "images": [],
  "video_links": [],
  "option_names": [],
  "properties": [],
  "characteristics": [],
  "product_accessory_links": [],
  "product_field_values": [],
  "product_bundle_components": [
    {
      "id": 1058167308,
      "variant_id": 1,
      "quantity": "2.0",
      "free": false,
      "discount_amount": null,
      "product_id": 1,
      "discount_type": 0
    }
  ],
  "variants": [
    {
      "id": 509944659,
      "title": "",
      "product_id": 293977361,
      "sku": null,
      "barcode": null,
      "dimensions": null,
      "available": false,
      "image_ids": [],
      "image_id": null,
      "weight": null,
      "created_at": "2026-07-25T00:59:56.000+03:00",
      "updated_at": "2026-07-25T00:59:56.000+03:00",
      "quantity": 0,
      "quantity_at_warehouse0": "0.0",
      "cost_price": null,
      "cost_price_in_site_currency": null,
      "price_in_site_currency": "3000.0",
      "base_price": "100.0",
      "old_price": null,
      "price2": null,
      "price": "100.0",
      "base_price_in_site_currency": "3000.0",
      "old_price_in_site_currency": null,
      "price2_in_site_currency": null,
      "prices": [],
      "prices_in_site_currency": [],
      "variant_field_values": [],
      "option_values": []
    }
  ]
}
```

### Create product

`POST /admin/products.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product[category_id]` | да | category id |
| `product[title]` | да | title |
| `product[description]` |  | description |
| `product[short_description]` |  | short description |
| `product[variants_attributes]` | да | variants attributes |
| `product[sort_weight]` |  | sort weight |
| `product[ignore_discounts]` |  | set to 1 to make this product ignore discounts |
| `product[vat]` |  | integer, default value is set in account settings, if default in account set as 'No VAT' you cannot change VAT of products through API; possible values depend on country, set to -1 if this product is not a subject to VAT |
| `product[product_field_values_attributes]` |  | product field values attributes |

**Запрос**

```http
POST /admin/products.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product": {
    "category_id": 1,
    "title": "Van Gogh Ruled Peach Notebook",
    "variants_attributes": [
      {
        "sku": "QP021MVEN-r",
        "quantity": 1,
        "price": 740
      }
    ],
    "ignore_discounts": 1,
    "vat": -1,
    "product_field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 293977345,
  "category_id": 1,
  "created_at": "2026-07-25T00:59:55.000+03:00",
  "updated_at": "2026-07-25T00:59:55.000+03:00",
  "is_hidden": false,
  "available": true,
  "canonical_url_collection_id": null,
  "archived": false,
  "unit": "pce",
  "sort_weight": null,
  "bundle": null,
  "ignore_discounts": true,
  "vat": -1,
  "dimensions": null,
  "fiscal_product_type": 1,
  "custom_template": null,
  "title": "Van Gogh Ruled Peach Notebook",
  "short_description": null,
  "permalink": "van-gogh-ruled-peach-notebook",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "currency_code": "LVL",
  "collections_ids": [],
  "sales_channels_id": null,
  "description": null,
  "images": [],
  "video_links": [],
  "option_names": [],
  "properties": [],
  "characteristics": [],
  "product_accessory_links": [],
  "product_field_values": [
    {
      "id": 9,
      "product_field_id": 1,
      "value": "12345"
    }
  ],
  "product_bundle_components": [],
  "variants": [
    {
      "id": 509944657,
      "title": "",
      "product_id": 293977345,
      "sku": "QP021MVEN-r",
      "barcode": null,
      "dimensions": null,
      "available": true,
      "image_ids": [],
      "image_id": null,
      "weight": null,
      "created_at": "2026-07-25T00:59:55.000+03:00",
      "updated_at": "2026-07-25T00:59:55.000+03:00",
      "quantity": 1,
      "quantity_at_warehouse0": "1.0",
      "cost_price": null,
      "cost_price_in_site_currency": null,
      "price_in_site_currency": "22200.0",
      "base_price": "740.0",
      "old_price": null,
      "price2": null,
      "price": "740.0",
      "base_price_in_site_currency": "22200.0",
      "old_price_in_site_currency": null,
      "price2_in_site_currency": null,
      "prices": [],
      "prices_in_site_currency": [],
      "variant_field_values": [],
      "option_values": []
    }
  ]
}
```

### Destroy product

`DELETE /admin/products/:product_id.json`

**Запрос**

```http
DELETE /admin/products/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get product

`GET /admin/products/:product_id.json`

**Запрос**

```http
GET /admin/products/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "category_id": 1,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2015-07-27T17:05:35.000+03:00",
  "is_hidden": false,
  "available": false,
  "canonical_url_collection_id": 1,
  "archived": false,
  "unit": "pce",
  "sort_weight": null,
  "bundle": null,
  "ignore_discounts": null,
  "vat": -1,
  "dimensions": "10x20x30",
  "fiscal_product_type": 1,
  "custom_template": null,
  "title": "Cool t-shirt",
  "short_description": "its a cool t-shirt",
  "permalink": "футболка",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "currency_code": "LVL",
  "collections_ids": [
    1,
    3,
    8
  ],
  "sales_channels_id": null,
  "description": "guys, its a really cool t-shirt",
  "images": [],
  "video_links": [],
  "option_names": [
    {
      "id": 1,
      "position": 1,
      "navigational": true,
      "permalink": null,
      "title": "Вариант",
      "api_permalink": null
    }
  ],
  "properties": [
    {
      "id": 6,
      "position": 1,
      "backoffice": false,
      "is_hidden": true,
      "is_navigational": true,
      "permalink": "hidden",
      "title": "hidden"
    },
    {
      "id": 8,
      "position": 3,
      "backoffice": false,
      "is_hidden": false,
      "is_navigational": false,
      "permalink": "not_navigational",
      "title": "not_navigational"
    },
    {
      "id": 4,
      "position": 5,
      "backoffice": false,
      "is_hidden": false,
      "is_navigational": true,
      "permalink": "style",
      "title": "Style"
    },
    {
      "id": 3,
      "position": 6,
      "backoffice": false,
      "is_hidden": false,
      "is_navigational": true,
      "permalink": "vendor",
      "title": "Vendor"
    }
  ],
  "characteristics": [
    {
      "id": 7,
      "property_id": 3,
      "position": 1,
      "title": "Nike",
      "permalink": "nike"
    },
    {
      "id": 9,
      "property_id": 4,
      "position": 1,
      "title": "Old",
      "permalink": "old"
    },
    {
      "id": 12,
      "property_id": 6,
      "position": 1,
      "title": "Hidden",
      "permalink": "Hidden"
    },
    {
      "id": 13,
      "property_id": 8,
      "position": 1,
      "title": "NotNavigational",
      "permalink": "NotNavigational"
    }
  ],
  "product_accessory_links": [],
  "product_field_values": [
    {
      "id": 1,
      "product_field_id": 1,
      "value": "some product field value"
    }
  ],
  "product_bundle_components": [],
  "variants": [
    {
      "id": 1,
      "title": "Red",
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "dimensions": "10x20x30",
      "available": false,
      "image_ids": [],
      "image_id": null,
      "weight": null,
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "updated_at": "2009-03-18T10:53:14.000+03:00",
      "quantity": 0,
      "quantity_at_warehouse0": "0.0",
      "cost_price": "3.0",
      "cost_price_in_site_currency": "90.0",
      "price_in_site_currency": "300.0",
      "base_price": "10.0",
      "old_price": null,
      "price2": 12.0,
      "price": "10.0",
      "base_price_in_site_currency": "300.0",
      "old_price_in_site_currency": null,
      "price2_in_site_currency": "360.0",
      "prices": [
        12.0
      ],
      "prices_in_site_currency": [
        "360.0"
      ],
      "variant_field_values": [
        {
          "id": 1,
          "variant_field_id": 1,
          "value": "some variant field value"
        }
      ],
      "option_values": [
        {
          "id": 1,
          "option_name_id": 1,
          "position": 1,
          "title": "Red",
          "image_url": null,
          "permalink": null
        }
      ]
    }
  ]
}
```

### Get products count

`GET /admin/products/count.json`

**Запрос**

```http
GET /admin/products/count.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "count": 5
}
```

### Get products

`GET /admin/products.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` |  | Set datetime to get only data updated after it. |
| `from_id` |  | Set id to get only data starting from it. |
| `per_page` |  | Quantity of objects per page. |
| `with_deleted` |  | Set to true to get also deleted products. |
| `deleted` |  | Set to true to get only deleted products. |
| `variant_fields` |  | List of required fields separated with `,`. Possible values are: id,title,product_id,sku,barcode,dimensions,available. |

**Запрос**

```http
GET /admin/products.json?updated_since=2009-01-04+14%3A44%3A43+%2B0300
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
updated_since=2009-01-04 14:44:43 +0300
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 5,
    "category_id": 3,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2015-07-26T23:59:59.000+03:00",
    "is_hidden": false,
    "available": true,
    "canonical_url_collection_id": null,
    "archived": false,
    "unit": "pce",
    "sort_weight": null,
    "bundle": true,
    "ignore_discounts": null,
    "vat": 22,
    "dimensions": null,
    "fiscal_product_type": 1,
    "custom_template": null,
    "title": "shirt & pants",
    "short_description": null,
    "permalink": "shirt_and_pants",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "currency_code": "LVL",
    "collections_ids": [],
    "sales_channels_id": null,
    "description": null,
    "images": [],
    "video_links": [],
    "option_names": [],
    "properties": [],
    "characteristics": [],
    "product_accessory_links": [],
    "product_field_values": [],
    "product_bundle_components": [
      {
        "id": 361433619,
        "variant_id": 2,
        "quantity": "1.0",
        "free": false,
        "discount_amount": null,
        "product_id": 2,
        "discount_type": null
      },
      {
        "id": 620730336,
        "variant_id": 3,
        "quantity": "1.0",
        "free": false,
        "discount_amount": null,
        "product_id": 3,
        "discount_type": null
      }
    ],
    "variants": [
      {
        "id": 5,
        "title": "",
        "product_id": 5,
        "sku": "427",
        "barcode": null,
        "dimensions": null,
        "available": true,
        "image_ids": [],
        "image_id": null,
        "weight": "0.3",
        "created_at": "2009-03-18T12:58:14.000+03:00",
        "updated_at": "2009-03-18T13:53:14.000+03:00",
        "quantity": null,
        "quantity_at_warehouse0": null,
        "cost_price": null,
        "cost_price_in_site_currency": null,
        "price_in_site_currency": "750.0",
        "base_price": "25.0",
        "old_price": null,
        "price2": null,
        "price": "25.0",
        "base_price_in_site_currency": "750.0",
        "old_price_in_site_currency": null,
        "price2_in_site_currency": null,
        "prices": [],
        "prices_in_site_currency": [],
        "variant_field_values": [],
        "option_values": []
      }
    ]
  },
  {
    "id": 2,
    "category_id": 1,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2015-07-27T00:00:00.000+03:00",
    "is_hidden": false,
    "available": true,
    "canonical_url_collection_id": null,
    "archived": false,
    "unit": "pce",
    "sort_weight": null,
    "bundle": null,
    "ignore_discounts": null,
    "vat": -1,
    "dimensions": null,
    "fiscal_product_type": 1,
    "custom_template": null,
    "title": "shirt",
    "short_description": null,
    "permalink": "shirt",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "currency_code": "LVL",
    "collections_ids": [
      1,
      8
    ],
    "sales_channels_id": null,
    "description": null,
    "images": [
      {
        "id": 1,
        "product_id": 2,
        "external_id": null,
        "position": 1,
        "created_at": "2026-07-25T00:58:09.000+03:00",
        "image_processing": null,
        "title": null,
        "url": "http://minio.insales.io:9000/insales-test-sber/images/products/1/1/1/thumb_image_1.jpg",
        "original_url": "http://minio.insales.io:9000/insales-test-sber/images/products/1/1/1/image_1.jpg",
        "medium_url": "http://minio.insales.io:9000/insales-test-sber/images/products/1/1/1/medium_image_1.jpg",
        "small_url": "http://minio.insales.io:9000/insales-test-sber/images/products/1/1/1/micro_image_1.jpg",
        "thumb_url": "http://minio.insales.io:9000/insales-test-sber/images/products/1/1/1/thumb_image_1.jpg",
        "compact_url": "http://minio.insales.io:9000/insales-test-sber/images/products/1/1/1/compact_image_1.jpg",
        "large_url": "http://minio.insales.io:9000/insales-test-sber/images/products/1/1/1/large_image_1.jpg",
        "filename": "image_1.jpg",
        "has_bitmap_image": false,
        "image_content_type": "image/jpeg"
      }
    ],
    "video_links": [],
    "option_names": [
      {
        "id": 1,
        "position": 1,
        "navigational": true,
        "permalink": null,
        "title": "Вариант",
        "api_permalink": null
      },
      {
        "id": 3,
        "position": 2,
        "navigational": true,
        "permalink": null,
        "title": "Цвет",
        "api_permalink": null
      }
    ],
    "properties": [
      {
        "id": 4,
        "position": 5,
        "backoffice": false,
        "is_hidden": false,
        "is_navigational": true,
        "permalink": "style",
        "title": "Style"
      }
    ],
    "characteristics": [
      {
        "id": 10,
        "property_id": 4,
        "position": 2,
        "title": "New",
        "permalink": "new"
      }
    ],
    "product_accessory_links": [],
    "product_field_values": [],
    "product_bundle_components": [],
    "variants": [
      {
        "id": 2,
        "title": "Зеленый / Синий",
        "product_id": 2,
        "sku": "424",
        "barcode": null,
        "dimensions": null,
        "available": true,
        "image_ids": [],
        "image_id": null,
        "weight": "0.3",
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "updated_at": "2009-03-18T10:53:14.000+03:00",
        "quantity": 1,
        "quantity_at_warehouse0": "1.0",
        "cost_price": null,
        "cost_price_in_site_currency": null,
        "price_in_site_currency": "600.0",
        "base_price": "20.0",
        "old_price": null,
        "price2": null,
        "price": "20.0",
        "base_price_in_site_currency": "600.0",
        "old_price_in_site_currency": null,
        "price2_in_site_currency": null,
        "prices": [],
        "prices_in_site_currency": [],
        "variant_field_values": [],
        "option_values": [
          {
            "id": 2,
            "option_name_id": 1,
            "position": 2,
            "title": "Зеленый",
            "image_url": null,
            "permalink": null
          },
          {
            "id": 4,
            "option_name_id": 3,
            "position": 1,
            "title": "Синий",
            "image_url": null,
            "permalink": null
          }
        ]
      }
    ]
  },
  {
    "id": 1,
    "category_id": 1,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2015-07-27T17:05:35.000+03:00",
    "is_hidden": false,
    "available": false,
    "canonical_url_collection_id": 1,
    "archived": false,
    "unit": "pce",
    "sort_weight": null,
    "bundle": null,
    "ignore_discounts": null,
    "vat": -1,
    "dimensions": "10x20x30",
    "fiscal_product_type": 1,
    "custom_template": null,
    "title": "Cool t-shirt",
    "short_description": "its a cool t-shirt",
    "permalink": "футболка",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "currency_code": "LVL",
    "collections_ids": [
      1,
      3,
      8
    ],
    "sales_channels_id": null,
    "description": "guys, its a really cool t-shirt",
    "images": [],
    "video_links": [],
    "option_names": [
      {
        "id": 1,
        "position": 1,
        "navigational": true,
        "permalink": null,
        "title": "Вариант",
        "api_permalink": null
      }
    ],
    "properties": [
      {
        "id": 6,
        "position": 1,
        "backoffice": false,
        "is_hidden": true,
        "is_navigational": true,
        "permalink": "hidden",
        "title": "hidden"
      },
      {
        "id": 8,
        "position": 3,
        "backoffice": false,
        "is_hidden": false,
        "is_navigational": false,
        "permalink": "not_navigational",
        "title": "not_navigational"
      },
      {
        "id": 4,
        "position": 5,
        "backoffice": false,
        "is_hidden": false,
        "is_navigational": true,
        "permalink": "style",
        "title": "Style"
      },
      {
        "id": 3,
        "position": 6,
        "backoffice": false,
        "is_hidden": false,
        "is_navigational": true,
        "permalink": "vendor",
        "title": "Vendor"
      }
    ],
    "characteristics": [
      {
        "id": 7,
        "property_id": 3,
        "position": 1,
        "title": "Nike",
        "permalink": "nike"
      },
      {
        "id": 9,
        "property_id": 4,
        "position": 1,
        "title": "Old",
        "permalink": "old"
      },
      {
        "id": 12,
        "property_id": 6,
        "position": 1,
        "title": "Hidden",
        "permalink": "Hidden"
      },
      {
        "id": 13,
        "property_id": 8,
        "position": 1,
        "title": "NotNavigational",
        "permalink": "NotNavigational"
      }
    ],
    "product_accessory_links": [],
    "product_field_values": [
      {
        "id": 1,
        "product_field_id": 1,
        "value": "some product field value"
      }
    ],
    "product_bundle_components": [],
    "variants": [
      {
        "id": 1,
        "title": "Red",
        "product_id": 1,
        "sku": "423",
        "barcode": null,
        "dimensions": "10x20x30",
        "available": false,
        "image_ids": [],
        "image_id": null,
        "weight": null,
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "updated_at": "2009-03-18T10:53:14.000+03:00",
        "quantity": 0,
        "quantity_at_warehouse0": "0.0",
        "cost_price": "3.0",
        "cost_price_in_site_currency": "90.0",
        "price_in_site_currency": "300.0",
        "base_price": "10.0",
        "old_price": null,
        "price2": 12.0,
        "price": "10.0",
        "base_price_in_site_currency": "300.0",
        "old_price_in_site_currency": null,
        "price2_in_site_currency": "360.0",
        "prices": [
          12.0
        ],
        "prices_in_site_currency": [
          "360.0"
        ],
        "variant_field_values": [
          {
            "id": 1,
            "variant_field_id": 1,
            "value": "some variant field value"
          }
        ],
        "option_values": [
          {
            "id": 1,
            "option_name_id": 1,
            "position": 1,
            "title": "Red",
            "image_url": null,
            "permalink": null
          }
        ]
      }
    ]
  },
  {
    "id": 3,
    "category_id": 6,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2015-07-28T23:59:59.000+03:00",
    "is_hidden": false,
    "available": true,
    "canonical_url_collection_id": null,
    "archived": false,
    "unit": "pce",
    "sort_weight": null,
    "bundle": null,
    "ignore_discounts": null,
    "vat": -1,
    "dimensions": null,
    "fiscal_product_type": 1,
    "custom_template": null,
    "title": "pants",
    "short_description": null,
    "permalink": "pants",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "currency_code": "LVL",
    "collections_ids": [
      8
    ],
    "sales_channels_id": null,
    "description": null,
    "images": [],
    "video_links": [],
    "option_names": [
      {
        "id": 1,
        "position": 1,
        "navigational": true,
        "permalink": null,
        "title": "Вариант",
        "api_permalink": null
      }
    ],
    "properties": [],
    "characteristics": [],
    "product_accessory_links": [],
    "product_field_values": [],
    "product_bundle_components": [],
    "variants": [
      {
        "id": 3,
        "title": "Базовая",
        "product_id": 3,
        "sku": "425",
        "barcode": null,
        "dimensions": null,
        "available": true,
        "image_ids": [],
        "image_id": null,
        "weight": null,
        "created_at": "2009-03-18T10:53:14.000+03:00",
        "updated_at": "2009-03-18T10:53:14.000+03:00",
        "quantity": null,
        "quantity_at_warehouse0": null,
        "cost_price": null,
        "cost_price_in_site_currency": null,
        "price_in_site_currency": "600.0",
        "base_price": "20.0",
        "old_price": null,
        "price2": null,
        "price": "20.0",
        "base_price_in_site_currency": "600.0",
        "old_price_in_site_currency": null,
        "price2_in_site_currency": null,
        "prices": [],
        "prices_in_site_currency": [],
        "variant_field_values": [],
        "option_values": [
          {
            "id": 3,
            "option_name_id": 1,
            "position": 3,
            "title": "Базовая",
            "image_url": null,
            "permalink": null
          }
        ]
      }
    ]
  },
  {
    "id": 4,
    "category_id": 3,
    "created_at": "2026-07-25T00:58:09.000+03:00",
    "updated_at": "2015-07-29T00:00:00.000+03:00",
    "is_hidden": false,
    "available": true,
    "canonical_url_collection_id": null,
    "archived": false,
    "unit": "pce",
    "sort_weight": null,
    "bundle": true,
    "ignore_discounts": null,
    "vat": 10,
    "dimensions": null,
    "fiscal_product_type": 1,
    "custom_template": null,
    "title": "shirt & t-short",
    "short_description": null,
    "permalink": "shirt_and_tshort",
    "html_title": null,
    "meta_keywords": null,
    "meta_description": null,
    "currency_code": "LVL",
    "collections_ids": [],
    "sales_channels_id": null,
    "description": null,
    "images": [],
    "video_links": [],
    "option_names": [],
    "properties": [],
    "characteristics": [],
    "product_accessory_links": [],
    "product_field_values": [],
    "product_bundle_components": [
      {
        "id": 989437296,
        "variant_id": 2,
        "quantity": "1.0",
        "free": true,
        "discount_amount": null,
        "product_id": 2,
        "discount_type": null
      },
      {
        "id": 579560919,
        "variant_id": 1,
        "quantity": "1.0",
        "free": false,
        "discount_amount": null,
        "product_id": 1,
        "discount_type": null
      }
    ],
    "variants": [
      {
        "id": 4,
        "title": "",
        "product_id": 4,
        "sku": "426",
        "barcode": null,
        "dimensions": null,
        "available": true,
        "image_ids": [],
        "image_id": null,
        "weight": "0.3",
        "created_at": "2009-03-18T12:53:14.000+03:00",
        "updated_at": "2009-03-18T11:53:14.000+03:00",
        "quantity": null,
        "quantity_at_warehouse0": null,
        "cost_price": null,
        "cost_price_in_site_currency": null,
        "price_in_site_currency": "600.0",
        "base_price": "20.0",
        "old_price": null,
        "price2": null,
        "price": "20.0",
        "base_price_in_site_currency": "600.0",
        "old_price_in_site_currency": null,
        "price2_in_site_currency": null,
        "prices": [],
        "prices_in_site_currency": [],
        "variant_field_values": [],
        "option_values": []
      }
    ]
  }
]
```

### Update bundle

`PUT /admin/products/4.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product[title]` |  | title |
| `product[ignore_discounts]` |  | set to 1 to make this product ignore discounts |
| `product[vat]` |  | do not set VAT for bundle because order lines will be created with its components, which have their own VAT |
| `product[product_bundle_components_attributes]` |  | component without the specified identifier (id) will be created (variant_id is required); component with id will be updated (variant_id is optional); component with id and (_destroy: true) attribute will be destroyed; |

**Запрос**

```http
PUT /admin/products/4.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product": {
    "product_bundle_components_attributes": [
      {
        "variant_id": 3,
        "quantity": 2,
        "free": false
      },
      {
        "id": 1,
        "quantity": 4
      },
      {
        "id": 2,
        "_destroy": true
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 4,
  "category_id": 3,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2026-07-25T00:59:57.000+03:00",
  "is_hidden": false,
  "available": true,
  "canonical_url_collection_id": null,
  "archived": false,
  "unit": "pce",
  "sort_weight": null,
  "bundle": true,
  "ignore_discounts": null,
  "vat": 10,
  "dimensions": null,
  "fiscal_product_type": 1,
  "custom_template": null,
  "title": "shirt & t-short",
  "short_description": null,
  "permalink": "shirt_and_tshort",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "currency_code": "LVL",
  "collections_ids": [],
  "sales_channels_id": null,
  "description": null,
  "images": [],
  "video_links": [],
  "option_names": [],
  "properties": [],
  "characteristics": [],
  "product_accessory_links": [],
  "product_field_values": [],
  "product_bundle_components": [
    {
      "id": 989437296,
      "variant_id": 2,
      "quantity": "1.0",
      "free": true,
      "discount_amount": null,
      "product_id": 2,
      "discount_type": 0
    },
    {
      "id": 579560919,
      "variant_id": 1,
      "quantity": "1.0",
      "free": false,
      "discount_amount": null,
      "product_id": 1,
      "discount_type": 0
    },
    {
      "id": 1058167310,
      "variant_id": 3,
      "quantity": "2.0",
      "free": false,
      "discount_amount": null,
      "product_id": 3,
      "discount_type": 0
    }
  ],
  "variants": [
    {
      "id": 4,
      "title": "",
      "product_id": 4,
      "sku": "426",
      "barcode": null,
      "dimensions": null,
      "available": false,
      "image_ids": [],
      "image_id": null,
      "weight": "0.3",
      "created_at": "2009-03-18T12:53:14.000+03:00",
      "updated_at": "2026-07-25T00:59:57.000+03:00",
      "quantity": 0,
      "quantity_at_warehouse0": "0.0",
      "cost_price": null,
      "cost_price_in_site_currency": null,
      "price_in_site_currency": "600.0",
      "base_price": "20.0",
      "old_price": null,
      "price2": null,
      "price": "20.0",
      "base_price_in_site_currency": "600.0",
      "old_price_in_site_currency": null,
      "price2_in_site_currency": null,
      "prices": [],
      "prices_in_site_currency": [],
      "variant_field_values": [],
      "option_values": []
    }
  ]
}
```

### Update product

`PUT /admin/products/:product_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product[title]` |  | title |
| `product[ignore_discounts]` |  | set to 1 to make this product ignore discounts |
| `product[vat]` |  | integer, default value is set in account settings, if default in account set as 'No VAT' you cannot change VAT of products through API; possible values depend on country, set to -1 if this product is not a subject to VAT |
| `product[product_field_values_attributes]` |  | product field values attributes |

**Запрос**

```http
PUT /admin/products/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product": {
    "title": "Van Gogh Notebook",
    "ignore_discounts": 1,
    "vat": 10,
    "product_field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "category_id": 1,
  "created_at": "2026-07-25T00:58:09.000+03:00",
  "updated_at": "2026-07-25T00:59:57.000+03:00",
  "is_hidden": false,
  "available": false,
  "canonical_url_collection_id": 1,
  "archived": false,
  "unit": "pce",
  "sort_weight": null,
  "bundle": null,
  "ignore_discounts": true,
  "vat": 10,
  "dimensions": "10x20x30",
  "fiscal_product_type": 1,
  "custom_template": null,
  "title": "Van Gogh Notebook",
  "short_description": "its a cool t-shirt",
  "permalink": "футболка",
  "html_title": null,
  "meta_keywords": null,
  "meta_description": null,
  "currency_code": "LVL",
  "collections_ids": [
    1,
    3,
    8
  ],
  "sales_channels_id": null,
  "description": "guys, its a really cool t-shirt",
  "images": [],
  "video_links": [],
  "option_names": [
    {
      "id": 1,
      "position": 1,
      "navigational": true,
      "permalink": null,
      "title": "Вариант",
      "api_permalink": null
    }
  ],
  "properties": [
    {
      "id": 6,
      "position": 1,
      "backoffice": false,
      "is_hidden": true,
      "is_navigational": true,
      "permalink": "hidden",
      "title": "hidden"
    },
    {
      "id": 8,
      "position": 3,
      "backoffice": false,
      "is_hidden": false,
      "is_navigational": false,
      "permalink": "not_navigational",
      "title": "not_navigational"
    },
    {
      "id": 4,
      "position": 5,
      "backoffice": false,
      "is_hidden": false,
      "is_navigational": true,
      "permalink": "style",
      "title": "Style"
    },
    {
      "id": 3,
      "position": 6,
      "backoffice": false,
      "is_hidden": false,
      "is_navigational": true,
      "permalink": "vendor",
      "title": "Vendor"
    }
  ],
  "characteristics": [
    {
      "id": 7,
      "property_id": 3,
      "position": 1,
      "title": "Nike",
      "permalink": "nike"
    },
    {
      "id": 9,
      "property_id": 4,
      "position": 1,
      "title": "Old",
      "permalink": "old"
    },
    {
      "id": 12,
      "property_id": 6,
      "position": 1,
      "title": "Hidden",
      "permalink": "Hidden"
    },
    {
      "id": 13,
      "property_id": 8,
      "position": 1,
      "title": "NotNavigational",
      "permalink": "NotNavigational"
    }
  ],
  "product_accessory_links": [],
  "product_field_values": [
    {
      "id": 1,
      "product_field_id": 1,
      "value": "12345"
    }
  ],
  "product_bundle_components": [],
  "variants": [
    {
      "id": 1,
      "title": "Red",
      "product_id": 1,
      "sku": "423",
      "barcode": null,
      "dimensions": "10x20x30",
      "available": false,
      "image_ids": [],
      "image_id": null,
      "weight": null,
      "created_at": "2009-03-18T10:53:14.000+03:00",
      "updated_at": "2009-03-18T10:53:14.000+03:00",
      "quantity": 0,
      "quantity_at_warehouse0": "0.0",
      "cost_price": "3.0",
      "cost_price_in_site_currency": "90.0",
      "price_in_site_currency": "300.0",
      "base_price": "10.0",
      "old_price": null,
      "price2": 12.0,
      "price": "10.0",
      "base_price_in_site_currency": "300.0",
      "old_price_in_site_currency": null,
      "price2_in_site_currency": "360.0",
      "prices": [
        12.0
      ],
      "prices_in_site_currency": [
        "360.0"
      ],
      "variant_field_values": [
        {
          "id": 1,
          "variant_field_id": 1,
          "value": "some variant field value"
        }
      ],
      "option_values": [
        {
          "id": 1,
          "option_name_id": 1,
          "position": 1,
          "title": "Red",
          "image_url": null,
          "permalink": null
        }
      ]
    }
  ]
}
```

---

## ProductField

| Операция | Метод и путь |
|----------|--------------|
| [Create product_field](#create-product-field) | `POST /admin/product_fields.json` |
| [Destroy product_field by handle](#destroy-product-field-by-handle) | `DELETE /admin/product_fields/:product_field_handle.json` |
| [Destroy product_field](#destroy-product-field) | `DELETE /admin/product_fields/:product_field_id.json` |
| [Get product_field by handle](#get-product-field-by-handle) | `GET /admin/product_fields/:product_field_handle.json` |
| [Get product_field](#get-product-field) | `GET /admin/product_fields/:product_field_id.json` |
| [Get product_fields](#get-product-fields) | `GET /admin/product_fields.json` |
| [Update product_field by handle](#update-product-field-by-handle) | `PUT /admin/product_fields/:product_field_handle.json` |
| [Update product_field](#update-product-field) | `PUT /admin/product_fields/:product_field_id.json` |

### Create product_field

`POST /admin/product_fields.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field[title]` | да | title |
| `product_field[handle]` | да | handle used in liquid |
| `product_field[type]` | да | ProductField::TextArea - textarea, ProductField::TextField - textfield |
| `product_field[position]` |  | position in product_fields list |
| `product_field[is_hidden]` |  | boolean |
| `product_field[is_indexed]` |  | boolean |

**Запрос**

```http
POST /admin/product_fields.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_field": {
    "title": "Size",
    "handle": "size",
    "type": "ProductField::TextArea"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 8,
  "handle": "size",
  "type": "ProductField::TextArea",
  "application_id": 1,
  "position": 3,
  "is_hidden": null,
  "is_indexed": null,
  "title": "Size",
  "application_title": "Цифровые товары"
}
```

### Destroy product_field by handle

`DELETE /admin/product_fields/:product_field_handle.json`

**Запрос**

```http
DELETE /admin/product_fields/handle2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Destroy product_field

`DELETE /admin/product_fields/:product_field_id.json`

**Запрос**

```http
DELETE /admin/product_fields/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get product_field by handle

`GET /admin/product_fields/:product_field_handle.json`

**Запрос**

```http
GET /admin/product_fields/handle2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "ProductField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "is_indexed": null,
  "title": "title2",
  "application_title": "Цифровые товары"
}
```

### Get product_field

`GET /admin/product_fields/:product_field_id.json`

**Запрос**

```http
GET /admin/product_fields/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "ProductField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "is_indexed": null,
  "title": "title2",
  "application_title": "Цифровые товары"
}
```

### Get product_fields

`GET /admin/product_fields.json`

**Запрос**

```http
GET /admin/product_fields.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "handle": "handle",
    "type": "ProductField::TextField",
    "application_id": null,
    "position": 1,
    "is_hidden": null,
    "is_indexed": null,
    "title": "title",
    "application_title": null
  },
  {
    "id": 4,
    "handle": "yml_import",
    "type": "ProductField::YmlImportField",
    "application_id": null,
    "position": 1,
    "is_hidden": null,
    "is_indexed": null,
    "title": "import field",
    "application_title": null
  },
  {
    "id": 2,
    "handle": "handle2",
    "type": "ProductField::TextArea",
    "application_id": 1,
    "position": 2,
    "is_hidden": null,
    "is_indexed": null,
    "title": "title2",
    "application_title": "Цифровые товары"
  }
]
```

### Update product_field by handle

`PUT /admin/product_fields/:product_field_handle.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field[title]` |  | title |
| `product_field[handle]` |  | handle used in liquid  new_handle |
| `product_field[position]` |  | position in product_fields list |
| `product_field[is_hidden]` |  | boolean |
| `product_field[is_indexed]` |  | boolean |

**Запрос**

```http
PUT /admin/product_fields/handle2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_field": {
    "title": "Volume",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "ProductField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "is_indexed": null,
  "title": "Volume",
  "application_title": "Цифровые товары"
}
```

### Update product_field

`PUT /admin/product_fields/:product_field_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field[title]` |  | title |
| `product_field[handle]` |  | handle used in liquid |
| `product_field[position]` |  | position in product_fields list |
| `product_field[is_hidden]` |  | boolean |
| `product_field[is_indexed]` |  | boolean |

**Запрос**

```http
PUT /admin/product_fields/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_field": {
    "title": "Volume",
    "handle": "volume",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "volume",
  "type": "ProductField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "is_indexed": null,
  "title": "Volume",
  "application_title": "Цифровые товары"
}
```

---

## ProductFieldValue

| Операция | Метод и путь |
|----------|--------------|
| [Create product_field_value by handle](#create-product-field-value-by-handle) | `POST /admin/products/:product_id/product_field_values.json` |
| [Create product_field_value](#create-product-field-value) | `POST /admin/products/:product_id/product_field_values.json` |
| [Destroy product_field_value by handle](#destroy-product-field-value-by-handle) | `DELETE /admin/products/:product_id/product_field_values/:product_field_handle.json` |
| [Destroy product_field_value](#destroy-product-field-value) | `DELETE /admin/products/:product_id/product_field_values/:product_field_value_id.json` |
| [Get product_field_value by handle](#get-product-field-value-by-handle) | `GET /admin/products/:product_id/product_field_values/:product_field_handle.json` |
| [Get product_field_value](#get-product-field-value) | `GET /admin/products/:product_id/product_field_values/:product_field_value_id.json` |
| [Get product_field_values](#get-product-field-values) | `GET /admin/products/:product_id/product_field_values.json` |
| [Update product_field_value by handle](#update-product-field-value-by-handle) | `PUT /admin/products/:product_id/product_field_values/:product_field_handle.json` |
| [Update product_field_value](#update-product-field-value) | `PUT /admin/products/:product_id/product_field_values/:product_field_value_id.json` |

### Create product_field_value by handle

`POST /admin/products/:product_id/product_field_values.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field_value[product_field_handle]` | да | product field handle |
| `product_field_value[value]` |  | product field value |

**Запрос**

```http
POST /admin/products/1/product_field_values.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_field_value": {
    "product_field_handle": "handle2",
    "value": "test value"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 7,
  "product_field_id": 2,
  "value": "test value"
}
```

### Create product_field_value

`POST /admin/products/:product_id/product_field_values.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field_value[product_field_id]` | да | product field id |
| `product_field_value[value]` |  | product field value |

**Запрос**

```http
POST /admin/products/1/product_field_values.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_field_value": {
    "product_field_id": 2,
    "value": "test value"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 5,
  "product_field_id": 2,
  "value": "test value"
}
```

### Destroy product_field_value by handle

`DELETE /admin/products/:product_id/product_field_values/:product_field_handle.json`

**Запрос**

```http
DELETE /admin/products/1/product_field_values/handle.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Destroy product_field_value

`DELETE /admin/products/:product_id/product_field_values/:product_field_value_id.json`

**Запрос**

```http
DELETE /admin/products/1/product_field_values/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get product_field_value by handle

`GET /admin/products/:product_id/product_field_values/:product_field_handle.json`

**Запрос**

```http
GET /admin/products/1/product_field_values/handle.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "product_field_id": 1,
  "value": "some product field value"
}
```

### Get product_field_value

`GET /admin/products/:product_id/product_field_values/:product_field_value_id.json`

**Запрос**

```http
GET /admin/products/1/product_field_values/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "product_field_id": 1,
  "value": "some product field value"
}
```

### Get product_field_values

`GET /admin/products/:product_id/product_field_values.json`

**Запрос**

```http
GET /admin/products/1/product_field_values.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "product_field_id": 1,
    "value": "some product field value"
  }
]
```

### Update product_field_value by handle

`PUT /admin/products/:product_id/product_field_values/:product_field_handle.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field_value[value]` |  | product field value |

**Запрос**

```http
PUT /admin/products/1/product_field_values/handle.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_field_value": {
    "value": "new value"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "product_field_id": 1,
  "value": "new value"
}
```

### Update product_field_value

`PUT /admin/products/:product_id/product_field_values/:product_field_value_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field_value[value]` |  | product field value |

**Запрос**

```http
PUT /admin/products/1/product_field_values/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "product_field_value": {
    "value": "new value"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "product_field_id": 1,
  "value": "new value"
}
```

---

## Property

| Операция | Метод и путь |
|----------|--------------|
| [Create property](#create-property) | `POST /admin/properties.json` |
| [Destroy property](#destroy-property) | `DELETE /admin/properties/:property_id.json` |
| [Get properties](#get-properties) | `GET /admin/properties.json` |
| [Get property](#get-property) | `GET /admin/properties/:property_id.json` |
| [Update property](#update-property) | `PUT /admin/properties/:property_id.json` |

### Create property

`POST /admin/properties.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `property[title]` | да | title |
| `property[permalink]` |  | permalink |
| `property[position]` |  | position in property list |
| `property[backoffice]` |  | show in backoffice |
| `property[is_hidden]` |  | hide in product page on site |
| `property[is_navigational]` |  | show in products filter on category page |

**Запрос**

```http
POST /admin/properties.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "property": {
    "title": "Country",
    "permalink": "country",
    "backoffice": true,
    "is_navigational": true
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 35187304,
  "position": 104,
  "backoffice": true,
  "is_hidden": false,
  "is_navigational": true,
  "permalink": "country",
  "title": "Country"
}
```

### Destroy property

`DELETE /admin/properties/:property_id.json`

**Запрос**

```http
DELETE /admin/properties/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get properties

`GET /admin/properties.json`

**Запрос**

```http
GET /admin/properties.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": -1,
    "position": 1,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "kategoriya-yandeks-marketa",
    "title": "Категория Яндекс Маркета"
  },
  {
    "id": 6,
    "position": 1,
    "backoffice": false,
    "is_hidden": true,
    "is_navigational": true,
    "permalink": "hidden",
    "title": "hidden"
  },
  {
    "id": -5,
    "position": 1,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "kategoriya-tovara-v-wildberries",
    "title": "Категория товара в Wildberries"
  },
  {
    "id": 7,
    "position": 2,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "visible",
    "title": "visible"
  },
  {
    "id": 9,
    "position": 3,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": false,
    "permalink": "numeric",
    "title": "numeric"
  },
  {
    "id": 8,
    "position": 3,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": false,
    "permalink": "not_navigational",
    "title": "not_navigational"
  },
  {
    "id": 5,
    "position": 4,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "параметр",
    "title": "параметр"
  },
  {
    "id": 4,
    "position": 5,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "style",
    "title": "Style"
  },
  {
    "id": 3,
    "position": 6,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "vendor",
    "title": "Vendor"
  },
  {
    "id": 2,
    "position": 7,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "importer",
    "title": "Importer"
  },
  {
    "id": 1,
    "position": 8,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "producer",
    "title": "Producer"
  },
  {
    "id": 14,
    "position": 11,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "age",
    "title": "Age"
  },
  {
    "id": 15,
    "position": 12,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "age_unit",
    "title": "Age Unit"
  },
  {
    "id": 100,
    "position": 100,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "color",
    "title": "Расцветка"
  },
  {
    "id": 101,
    "position": 101,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "shape",
    "title": "Форма"
  },
  {
    "id": 102,
    "position": 102,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "size",
    "title": "Размер"
  },
  {
    "id": 103,
    "position": 103,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "material",
    "title": "Материал"
  },
  {
    "id": -7,
    "position": 10000,
    "backoffice": false,
    "is_hidden": true,
    "is_navigational": false,
    "permalink": "new-ketagoriya-tovara-v-ozon",
    "title": "Новая категория товара в OZON"
  },
  {
    "id": -6,
    "position": 10001,
    "backoffice": false,
    "is_hidden": false,
    "is_navigational": true,
    "permalink": "kategoriya-tovara-v-vk",
    "title": "Категория товара в VK"
  }
]
```

### Get property

`GET /admin/properties/:property_id.json`

**Запрос**

```http
GET /admin/properties/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "position": 8,
  "backoffice": false,
  "is_hidden": false,
  "is_navigational": true,
  "permalink": "producer",
  "title": "Producer"
}
```

### Update property

`PUT /admin/properties/:property_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `property[title]` |  | title |

**Запрос**

```http
PUT /admin/properties/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "property": {
    "title": "Brand country"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "position": 8,
  "backoffice": false,
  "is_hidden": false,
  "is_navigational": true,
  "permalink": "producer",
  "title": "Brand country"
}
```

---

## ProviderToken

| Операция | Метод и путь |
|----------|--------------|
| [Create or refresh messengers custom provider token](#create-or-refresh-messengers-custom-provider-token) | `POST /admin/dialogs/provider_tokens.json` |

### Create or refresh messengers custom provider token

`POST /admin/dialogs/provider_tokens.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `messenger_type` | да | Messenger type  ["whats_app", "telegram", "viber", "sms", "email", "vk", "instagram"] |

**Запрос**

```http
POST /admin/dialogs/provider_tokens.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "messenger_type": "sms"
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "api_token": "SOME_API_TOKEN",
  "webhook_token": "SOME_WEBHOOK_TOKEN"
}
```

---

## RecurringApplicationCharge

| Операция | Метод и путь |
|----------|--------------|
| [Create recurring_application_charge](#create-recurring-application-charge) | `POST /admin/recurring_application_charge.json` |
| [Destroy recurring_application_charge](#destroy-recurring-application-charge) | `DELETE /admin/recurring_application_charge.json` |
| [Get recurring_application_charge](#get-recurring-application-charge) | `GET /admin/recurring_application_charge.json` |
| [Update recurring_application_charge](#update-recurring-application-charge) | `PUT /admin/recurring_application_charge.json` |
| [add free days to recurring_application_charge](#add-free-days-to-recurring-application-charge) | `PATCH /admin/recurring_application_charge/add_free_days.json` |

### Create recurring_application_charge

`POST /admin/recurring_application_charge.json`

Application can have only one recurring charge at any given time.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `recurring_application_charge[monthly]` | да | amount of monthly payment |
| `recurring_application_charge[trial_expired_at]` |  | trial period ending date, by default it uses trial period from application |

**Запрос**

```http
POST /admin/recurring_application_charge.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "recurring_application_charge": {
    "monthly": 200,
    "trial_expired_at": "2017-02-03"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "monthly": "200.0",
  "trial_expired_at": "2017-02-03",
  "created_at": "2026-07-25T00:58:43.936+03:00",
  "updated_at": "2026-07-25T00:58:43.936+03:00",
  "paid_till": "2017-02-03",
  "blocked": false
}
```

### Destroy recurring_application_charge

`DELETE /admin/recurring_application_charge.json`

**Запрос**

```http
DELETE /admin/recurring_application_charge.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get recurring_application_charge

`GET /admin/recurring_application_charge.json`

**Запрос**

```http
GET /admin/recurring_application_charge.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "monthly": "1190.0",
  "trial_expired_at": "2026-08-01",
  "created_at": "2026-07-25T00:58:43.530+03:00",
  "updated_at": "2026-07-25T00:58:43.530+03:00",
  "paid_till": "2026-08-01",
  "blocked": false
}
```

### Update recurring_application_charge

`PUT /admin/recurring_application_charge.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `recurring_application_charge[monthly]` |  | amount of monthly payment |

**Запрос**

```http
PUT /admin/recurring_application_charge.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "recurring_application_charge": {
    "monthly": 300
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "monthly": "300.0",
  "trial_expired_at": "2026-08-01",
  "created_at": "2026-07-25T00:58:44.305+03:00",
  "updated_at": "2026-07-25T00:58:44.438+03:00",
  "paid_till": "2026-08-01",
  "blocked": false
}
```

### add free days to recurring_application_charge

`PATCH /admin/recurring_application_charge/add_free_days.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `recurring_application_charge[days_count]` |  | how many free days you want to add, it is possible that paid_till attribute would change but trial_expired_at would not, it is normal |

**Запрос**

```http
PATCH /admin/recurring_application_charge/add_free_days.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "recurring_application_charge": {
    "days_count": 21
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "monthly": "400.0",
  "trial_expired_at": "2026-08-23",
  "created_at": "2026-07-25T00:58:44.785+03:00",
  "updated_at": "2026-07-25T00:58:45.013+03:00",
  "paid_till": "2026-08-23",
  "blocked": false
}
```

---

## Redirect

| Операция | Метод и путь |
|----------|--------------|
| [Create redirect](#create-redirect) | `POST /admin/redirects.json` |
| [Create redirect with removing redirect chains](#create-redirect-with-removing-redirect-chains) | `POST /admin/redirects.json` |
| [Destroy redirect](#destroy-redirect) | `DELETE /admin/redirects/:redirect_id.json` |
| [Get redirect](#get-redirect) | `GET /admin/redirects/:redirect_id.json` |
| [Get redirects](#get-redirects) | `GET /admin/redirects.json` |

### Create redirect

`POST /admin/redirects.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `redirect[old_url]` | да | old url |
| `redirect[new_url]` | да | new url |

**Запрос**

```http
POST /admin/redirects.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "redirect": {
    "old_url": "/old_url",
    "new_url": "/new_url"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 8,
  "old_url": "/old_url",
  "new_url": "/new_url",
  "created_at": "2026-07-25T01:00:32.791+03:00"
}
```

### Create redirect with removing redirect chains

`POST /admin/redirects.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `redirect[old_url]` | да | old url |
| `redirect[new_url]` | да | new url |
| `redirect[remove_redirect_chains]` |  | if set to 1 and given a redirect A=>B when creating a new redirect B=>C, it will destroy old redirect (A=>B) and create new one (A=>C) |

**Запрос**

```http
POST /admin/redirects.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "redirect": {
    "old_url": "/like_new_url",
    "new_url": "/so_new_url",
    "remove_redirect_chains": 1
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 12,
  "old_url": "/like_new_url",
  "new_url": "/so_new_url",
  "created_at": "2026-07-25T01:00:33.085+03:00"
}
```

### Destroy redirect

`DELETE /admin/redirects/:redirect_id.json`

**Запрос**

```http
DELETE /admin/redirects/14.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get redirect

`GET /admin/redirects/:redirect_id.json`

**Запрос**

```http
GET /admin/redirects/4.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 4,
  "old_url": "/so_old_url",
  "new_url": "/like_new_url",
  "created_at": "2026-07-25T01:00:32.445+03:00"
}
```

### Get redirects

`GET /admin/redirects.json`

**Запрос**

```http
GET /admin/redirects.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 2,
    "old_url": "/so_old_url",
    "new_url": "/like_new_url",
    "created_at": "2026-07-25T01:00:32.199+03:00"
  }
]
```

---

## Review

| Операция | Метод и путь |
|----------|--------------|
| [Create review](#create-review) | `POST /admin/reviews.json` |
| [Destroy review](#destroy-review) | `DELETE /admin/reviews/:review_id.json` |
| [Get review](#get-review) | `GET /admin/reviews/:review_id.json` |
| [Get reviews](#get-reviews) | `GET /admin/reviews.json` |
| [Update review](#update-review) | `PUT /admin/reviews/:review_id.json` |

### Create review

`POST /admin/reviews.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `review[product_id]` | да | product id |
| `review[author]` | да | author |
| `review[content]` |  | content |
| `review[email]` |  | email |
| `review[published]` |  | should be published? |
| `review[rating]` |  | product rating (1-10) |

**Запрос**

```http
POST /admin/reviews.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "review": {
    "product_id": 1,
    "author": "Bob",
    "content": "Nice product! I'd like to buy it!",
    "email": "bob@gmail.com",
    "published": true,
    "rating": 7
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 12,
  "product_id": 1,
  "author": "Bob",
  "content": "Nice product! I'd like to buy it!",
  "email": "bob@gmail.com",
  "rating": 7,
  "manager_name": null,
  "published": true,
  "created_at": "2026-07-25T00:59:59.000+03:00",
  "updated_at": "2026-07-25T00:59:59.000+03:00",
  "manager_reply": null,
  "replied_at": null,
  "spam": false,
  "published_at": null,
  "replied_q": false,
  "title": "Cool t-shirt",
  "url": "/product/футболка",
  "first_image": {
    "id": null,
    "product_id": null,
    "external_id": null,
    "position": null,
    "created_at": null,
    "image_processing": null,
    "title": null,
    "url": "/images/no_image_thumb.jpg",
    "original_url": "/images/no_image_original.jpg",
    "medium_url": "/images/no_image_medium.jpg",
    "small_url": "/images/no_image_micro.jpg",
    "thumb_url": "/images/no_image_thumb.jpg",
    "compact_url": "/images/no_image_compact.jpg",
    "large_url": "/images/no_image_large.jpg",
    "filename": null,
    "has_bitmap_image": false,
    "image_content_type": null
  }
}
```

### Destroy review

`DELETE /admin/reviews/:review_id.json`

**Запрос**

```http
DELETE /admin/reviews/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get review

`GET /admin/reviews/:review_id.json`

**Запрос**

```http
GET /admin/reviews/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "product_id": 1,
  "author": "Автор 1",
  "content": "Отличный продукт!",
  "email": "email1@example.com",
  "rating": 5,
  "manager_name": null,
  "published": true,
  "created_at": "2023-09-01T13:00:00.000+03:00",
  "updated_at": "2023-09-01T13:00:00.000+03:00",
  "manager_reply": "nil",
  "replied_at": null,
  "spam": false,
  "published_at": null,
  "replied_q": false,
  "title": "Cool t-shirt",
  "url": "/product/футболка",
  "first_image": {
    "id": null,
    "product_id": null,
    "external_id": null,
    "position": null,
    "created_at": null,
    "image_processing": null,
    "title": null,
    "url": "/images/no_image_thumb.jpg",
    "original_url": "/images/no_image_original.jpg",
    "medium_url": "/images/no_image_medium.jpg",
    "small_url": "/images/no_image_micro.jpg",
    "thumb_url": "/images/no_image_thumb.jpg",
    "compact_url": "/images/no_image_compact.jpg",
    "large_url": "/images/no_image_large.jpg",
    "filename": null,
    "has_bitmap_image": false,
    "image_content_type": null
  },
  "image": {
    "id": 1,
    "created_at": "2023-09-01T13:00:00.000+03:00",
    "original_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/image_1.jpg",
    "medium_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/medium_image_1.jpg",
    "small_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/micro_image_1.jpg",
    "thumb_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/thumb_image_1.jpg",
    "compact_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/compact_image_1.jpg",
    "large_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/large_image_1.jpg",
    "image_content_type": "image/jpeg"
  }
}
```

### Get reviews

`GET /admin/reviews.json`

**Запрос**

```http
GET /admin/reviews.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 10,
    "product_id": 2,
    "author": "Автор 10",
    "content": "Быстрая доставка, хороший сервис.",
    "email": "email10@example.com",
    "rating": 3,
    "manager_name": null,
    "published": false,
    "created_at": "2023-09-10T18:00:00.000+03:00",
    "updated_at": "2023-09-10T18:00:00.000+03:00",
    "manager_reply": "nil",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "shirt",
    "url": "/product/shirt",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    }
  },
  {
    "id": 9,
    "product_id": 1,
    "author": "Автор 9",
    "content": "Отличный продукт, рекомендую!",
    "email": "email9@example.com",
    "rating": 5,
    "manager_name": null,
    "published": true,
    "created_at": "2023-09-09T12:00:00.000+03:00",
    "updated_at": "2023-09-09T12:00:00.000+03:00",
    "manager_reply": "Спасибо за ваш отзыв!",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "Cool t-shirt",
    "url": "/product/футболка",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    }
  },
  {
    "id": 7,
    "product_id": 7,
    "author": "Автор 7",
    "content": "Понравился, всё работает.",
    "email": "email7@example.com",
    "rating": 5,
    "manager_name": null,
    "published": true,
    "created_at": "2023-09-07T13:30:00.000+03:00",
    "updated_at": "2023-09-07T13:30:00.000+03:00",
    "manager_reply": "Спасибо за ваш отзыв!",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "fb_product",
    "url": "/product/fb_product",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    }
  },
  {
    "id": 6,
    "product_id": 6,
    "author": "Автор 6",
    "content": "Средний продукт, есть лучше.",
    "email": "email6@example.com",
    "rating": 3,
    "manager_name": null,
    "published": true,
    "created_at": "2023-09-06T11:00:00.000+03:00",
    "updated_at": "2023-09-06T11:00:00.000+03:00",
    "manager_reply": "Спасибо за ваш отзыв.",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "pants",
    "url": "/product/pants",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    }
  },
  {
    "id": 5,
    "product_id": 5,
    "author": "Автор 5",
    "content": "Отличная цена-качество!",
    "email": "email5@example.com",
    "rating": 5,
    "manager_name": null,
    "published": true,
    "created_at": "2023-09-05T17:00:00.000+03:00",
    "updated_at": "2023-09-05T17:00:00.000+03:00",
    "manager_reply": "Рады, что вам понравилось!",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "shirt & pants",
    "url": "/product/shirt_and_pants",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    }
  },
  {
    "id": 4,
    "product_id": 4,
    "author": "Автор 4",
    "content": "Продукт нормальный, но требует доработки.",
    "email": "email4@example.com",
    "rating": 4,
    "manager_name": null,
    "published": true,
    "created_at": "2023-09-04T12:00:00.000+03:00",
    "updated_at": "2023-09-04T12:00:00.000+03:00",
    "manager_reply": "Спасибо за ваш отзыв. Мы работаем над улучшениями.",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "shirt & t-short",
    "url": "/product/shirt_and_tshort",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    }
  },
  {
    "id": 2,
    "product_id": 2,
    "author": "Автор 2",
    "content": "Хороший, но есть недостатки.",
    "email": "email2@example.com",
    "rating": 3,
    "manager_name": null,
    "published": true,
    "created_at": "2023-09-02T14:00:00.000+03:00",
    "updated_at": "2023-09-02T14:00:00.000+03:00",
    "manager_reply": "Мы учтем ваши замечания и постараемся улучшить продукт.",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "shirt",
    "url": "/product/shirt",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    },
    "image": {
      "id": 2,
      "created_at": "2023-09-02T14:00:00.000+03:00",
      "original_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/2/2/image_2.jpg",
      "medium_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/2/2/medium_image_2.jpg",
      "small_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/2/2/micro_image_2.jpg",
      "thumb_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/2/2/thumb_image_2.jpg",
      "compact_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/2/2/compact_image_2.jpg",
      "large_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/2/2/large_image_2.jpg",
      "image_content_type": "image/jpeg"
    }
  },
  {
    "id": 1,
    "product_id": 1,
    "author": "Автор 1",
    "content": "Отличный продукт!",
    "email": "email1@example.com",
    "rating": 5,
    "manager_name": null,
    "published": true,
    "created_at": "2023-09-01T13:00:00.000+03:00",
    "updated_at": "2023-09-01T13:00:00.000+03:00",
    "manager_reply": "nil",
    "replied_at": null,
    "spam": false,
    "published_at": null,
    "replied_q": false,
    "title": "Cool t-shirt",
    "url": "/product/футболка",
    "first_image": {
      "id": null,
      "product_id": null,
      "external_id": null,
      "position": null,
      "created_at": null,
      "image_processing": null,
      "title": null,
      "url": "/images/no_image_thumb.jpg",
      "original_url": "/images/no_image_original.jpg",
      "medium_url": "/images/no_image_medium.jpg",
      "small_url": "/images/no_image_micro.jpg",
      "thumb_url": "/images/no_image_thumb.jpg",
      "compact_url": "/images/no_image_compact.jpg",
      "large_url": "/images/no_image_large.jpg",
      "filename": null,
      "has_bitmap_image": false,
      "image_content_type": null
    },
    "image": {
      "id": 1,
      "created_at": "2023-09-01T13:00:00.000+03:00",
      "original_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/image_1.jpg",
      "medium_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/medium_image_1.jpg",
      "small_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/micro_image_1.jpg",
      "thumb_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/thumb_image_1.jpg",
      "compact_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/compact_image_1.jpg",
      "large_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/large_image_1.jpg",
      "image_content_type": "image/jpeg"
    }
  }
]
```

### Update review

`PUT /admin/reviews/:review_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `review[content]` |  | content |

**Запрос**

```http
PUT /admin/reviews/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "review": {
    "content": "Sweet product!"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "product_id": 1,
  "author": "Автор 1",
  "content": "Sweet product!",
  "email": "email1@example.com",
  "rating": 5,
  "manager_name": null,
  "published": true,
  "created_at": "2023-09-01T13:00:00.000+03:00",
  "updated_at": "2026-07-25T00:59:59.000+03:00",
  "manager_reply": "nil",
  "replied_at": null,
  "spam": false,
  "published_at": null,
  "replied_q": false,
  "title": "Cool t-shirt",
  "url": "/product/футболка",
  "first_image": {
    "id": null,
    "product_id": null,
    "external_id": null,
    "position": null,
    "created_at": null,
    "image_processing": null,
    "title": null,
    "url": "/images/no_image_thumb.jpg",
    "original_url": "/images/no_image_original.jpg",
    "medium_url": "/images/no_image_medium.jpg",
    "small_url": "/images/no_image_micro.jpg",
    "thumb_url": "/images/no_image_thumb.jpg",
    "compact_url": "/images/no_image_compact.jpg",
    "large_url": "/images/no_image_large.jpg",
    "filename": null,
    "has_bitmap_image": false,
    "image_content_type": null
  },
  "image": {
    "id": 1,
    "created_at": "2023-09-01T13:00:00.000+03:00",
    "original_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/image_1.jpg",
    "medium_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/medium_image_1.jpg",
    "small_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/micro_image_1.jpg",
    "thumb_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/thumb_image_1.jpg",
    "compact_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/compact_image_1.jpg",
    "large_url": "http://minio.insales.io:9000/insales-test-sber/images/reviews/1/1/1/large_image_1.jpg",
    "image_content_type": "image/jpeg"
  }
}
```

---

## SalesChannel

| Операция | Метод и путь |
|----------|--------------|
| [Create sales channel](#create-sales-channel) | `POST /admin/sales_channels.json` |
| [Destroy sales channel](#destroy-sales-channel) | `DELETE /admin/sales_channels/:sales_channel_id.json` |
| [Get sales channel](#get-sales-channel) | `GET /admin/sales_channels/:sales_channel_id.json` |
| [Get sales channels](#get-sales-channels) | `GET /admin/sales_channels.json` |

### Create sales channel

`POST /admin/sales_channels.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `sales_channel[title]` | да | sales channel title |

**Запрос**

```http
POST /admin/sales_channels.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "sales_channel": {
    "title": "sales channel title"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 10,
  "title": "sales channel title",
  "state": null,
  "messages": [],
  "account_id": 1,
  "goods_setting_type": null,
  "goods_setting_id": null,
  "created_at": "2026-07-25T01:00:00.833+03:00",
  "updated_at": "2026-07-25T01:00:00.833+03:00",
  "application_id": 1,
  "last_rating_import_started_at": null,
  "last_feedbacks_import_started_at": null,
  "additional_fields": {},
  "auto_export_products": "disabled"
}
```

### Destroy sales channel

`DELETE /admin/sales_channels/:sales_channel_id.json`

**Запрос**

```http
DELETE /admin/sales_channels/3.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get sales channel

`GET /admin/sales_channels/:sales_channel_id.json`

**Запрос**

```http
GET /admin/sales_channels/3.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 3,
  "title": "test title",
  "state": null,
  "messages": [],
  "account_id": 1,
  "goods_setting_type": null,
  "goods_setting_id": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2009-03-18T10:53:14.000+03:00",
  "application_id": 1,
  "last_rating_import_started_at": null,
  "last_feedbacks_import_started_at": null,
  "additional_fields": {},
  "auto_export_products": "disabled"
}
```

### Get sales channels

`GET /admin/sales_channels.json`

**Запрос**

```http
GET /admin/sales_channels.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 3,
    "title": "test title",
    "state": null,
    "messages": [],
    "account_id": 1,
    "goods_setting_type": null,
    "goods_setting_id": null,
    "created_at": "2009-03-18T10:53:14.000+03:00",
    "updated_at": "2009-03-18T10:53:14.000+03:00",
    "application_id": 1,
    "last_rating_import_started_at": null,
    "last_feedbacks_import_started_at": null,
    "additional_fields": {},
    "auto_export_products": "disabled"
  }
]
```

---

## Similar

| Операция | Метод и путь |
|----------|--------------|
| [Create similar](#create-similar) | `POST /admin/products/:product_id/similars.json` |
| [Destroy similar](#destroy-similar) | `DELETE /admin/products/:product_id/similars/:similar_id.json` |
| [Get similars](#get-similars) | `GET /admin/products/:product_id/similars.json` |
| [Stock Currency](#stock-currency) | `POST /admin/stock_currencies.json` |
| [Get stock currencies](#get-stock-currencies) | `GET /admin/stock_currencies.json` |
| [Get stock currency](#get-stock-currency) | `GET /admin/stock_currencies/:stock_currency_id.json` |
| [Update stock_currency](#update-stock-currency) | `PUT /admin/stock_currencies/:stock_currency_id.json` |

### Create similar

`POST /admin/products/:product_id/similars.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `similar_ids` | да | similar ids array |

**Запрос**

```http
POST /admin/products/1/similars.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "similar_ids": [
    2
  ]
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Destroy similar

`DELETE /admin/products/:product_id/similars/:similar_id.json`

**Запрос**

```http
DELETE /admin/products/1/similars/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get similars

`GET /admin/products/:product_id/similars.json`

**Запрос**

```http
GET /admin/products/1/similars.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 3
  }
]
```

### Stock Currency

`POST /admin/stock_currencies.json`

Create stock_currency Maximal quantity of currencies per account is 5.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `stock_currency[code]` | да | Available codes:LVL, BYN, USD, EUR, EEK, LTL, AMD, KGS, MDL, EGP, CAD, JPY, TNB, TRY, AZN, PLN, ILS, CNY, BGN, GBP, RUR, UAH, AUD, GEL, AED, CHF, KZT, UZS, TJS, VND, TMT, MXN |
| `stock_currency[is_default]` |  | false by default, set it true to make this currency default for account |
| `stock_currency[exchange_rate_use_cb]` |  | true by default, set it false to use manual exchange rate |
| `stock_currency[exchange_rate_manual]` |  | manual exchange rate |
| `stock_currency[exchange_rate_percent]` |  | 0 by default, deviation as a percentage from defined exchange rate |

**Запрос**

```http
POST /admin/stock_currencies.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "stock_currency": {
    "code": "KZT",
    "exchange_rate_percent": 10
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 93,
  "is_default": false,
  "exchange_rate_use_cb": true,
  "exchange_rate_manual": null,
  "exchange_rate_percent": "10.0",
  "exchange_rate_round": true,
  "code": "KZT",
  "title": "Казахский тенге",
  "cb_rate": "0.21",
  "unit": "тг.",
  "exchange_rate": "0.231"
}
```

### Get stock currencies

`GET /admin/stock_currencies.json`

**Запрос**

```http
GET /admin/stock_currencies.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "is_default": true,
    "exchange_rate_use_cb": false,
    "exchange_rate_manual": "30.0",
    "exchange_rate_percent": "0.0",
    "exchange_rate_round": true,
    "code": "LVL",
    "title": "Латвийский лат",
    "cb_rate": "30.0",
    "unit": "Ls",
    "exchange_rate": "30.0"
  }
]
```

### Get stock currency

`GET /admin/stock_currencies/:stock_currency_id.json`

**Запрос**

```http
GET /admin/stock_currencies/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "is_default": true,
  "exchange_rate_use_cb": false,
  "exchange_rate_manual": "30.0",
  "exchange_rate_percent": "0.0",
  "exchange_rate_round": true,
  "code": "LVL",
  "title": "Латвийский лат",
  "cb_rate": "30.0",
  "unit": "Ls",
  "exchange_rate": "30.0"
}
```

### Update stock_currency

`PUT /admin/stock_currencies/:stock_currency_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `stock_currency[code]` |  | Available codes:LVL, BYN, USD, EUR, EEK, LTL, AMD, KGS, MDL, EGP, CAD, JPY, TNB, TRY, AZN, PLN, ILS, CNY, BGN, GBP, RUR, UAH, AUD, GEL, AED, CHF, KZT, UZS, TJS, VND, TMT, MXN |
| `stock_currency[is_default]` |  | false by default, set it true to make this currency default for account |
| `stock_currency[exchange_rate_use_cb]` |  | true by default, set it false to use manual exchange rate |
| `stock_currency[exchange_rate_manual]` |  | manual exchange rate |
| `stock_currency[exchange_rate_percent]` |  | 0 by default, deviation as a percentage from defined exchange rate |

**Запрос**

```http
PUT /admin/stock_currencies/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "stock_currency": {
    "code": "UAH",
    "is_default": true,
    "exchange_rate_use_cb": false,
    "exchange_rate_manual": 40,
    "exchange_rate_percent": 0
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "is_default": true,
  "exchange_rate_use_cb": false,
  "exchange_rate_manual": "40.0",
  "exchange_rate_percent": "0.0",
  "exchange_rate_round": true,
  "code": "UAH",
  "title": "Украинская гривна",
  "cb_rate": "5.0",
  "unit": "грн.",
  "exchange_rate": "40.0"
}
```

---

## Supplementary

| Операция | Метод и путь |
|----------|--------------|
| [Create supplementary](#create-supplementary) | `POST /admin/products/:product_id/supplementaries.json` |
| [Destroy supplementary](#destroy-supplementary) | `DELETE /admin/products/:product_id/supplementaries/:supplementary_id.json` |
| [Get supplementaries](#get-supplementaries) | `GET /admin/products/:product_id/supplementaries.json` |

### Create supplementary

`POST /admin/products/:product_id/supplementaries.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `supplementary_ids` | да | supplementary ids array |

**Запрос**

```http
POST /admin/products/1/supplementaries.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "supplementary_ids": [
    2,
    3
  ]
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Destroy supplementary

`DELETE /admin/products/:product_id/supplementaries/:supplementary_id.json`

**Запрос**

```http
DELETE /admin/products/1/supplementaries/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get supplementaries

`GET /admin/products/:product_id/supplementaries.json`

**Запрос**

```http
GET /admin/products/1/supplementaries.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 3
  }
]
```

---

## Unit

| Операция | Метод и путь |
|----------|--------------|
| [Get units](#get-units) | `GET admin/units.json?lang=ru` |

### Get units

`GET admin/units.json?lang=ru`

**Запрос**

```http
GET admin/units.json?lang=ru
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

Query:

```text
lang=ru
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "mmt": "мм",
  "cmt": "см",
  "dmt": "дм",
  "mtr": "м",
  "kmt": "км",
  "inh": "дюйм",
  "fot": "фут",
  "yrd": "ярд",
  "mtk": "м2",
  "dmk": "дм2",
  "mlt": "мл",
  "cmq": "см3",
  "ltr": "л",
  "dmq": "дм3",
  "mtq": "м3",
  "mgm": "мг",
  "grm": "г",
  "kgm": "кг",
  "tne": "т",
  "hur": "ч",
  "nbb": "боб",
  "lef": "л.",
  "set": "набор",
  "npr": "пар",
  "npl": "рул",
  "nmp": "упак",
  "dzp": "дюжина упак",
  "cnp": "100 упак",
  "pce": "шт",
  "cen": "100 шт",
  "mil": "тыс. шт",
  "box": "ящ",
  "sct": "секц",
  "btl": "флак",
  "lmt": "пог. м",
  "com": "компл",
  "pct": "кор",
  "day": "сут",
  "mth": "мес"
}
```

---

## UserTask

| Операция | Метод и путь |
|----------|--------------|
| [Create user task](#create-user-task) | `POST /admin/user_tasks.json` |
| [Delete user task](#delete-user-task) | `DELETE /admin/user_tasks/:user_task_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/user_tasks/:user_task_id.json` |
| [Get user tasks list](#get-user-tasks-list) | `GET /admin/user_tasks.json` |
| [Update user task](#update-user-task) | `PUT /admin/user_tasks/:user_task_id.json` |

### Create user task

`POST /admin/user_tasks.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `user_task[title]` | да | user task title |
| `user_task[description]` | да | user task description |
| `user_task[assignee_id]` | да | user id |
| `user_task[deadline_at]` |  | deadline |
| `user_task[complete]` |  | is task completed? |
| `user_task[order_id]` |  | order id |

**Запрос**

```http
POST /admin/user_tasks.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "user_task": {
    "title": "dispatch order N2",
    "description": "You should dispatch the order N2",
    "assignee_id": 1,
    "deadline_at": "2021-01-01",
    "order_id": 1
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 8,
  "complete": null,
  "order_id": 1,
  "assignee_id": 1,
  "title": "dispatch order N2",
  "description": "You should dispatch the order N2",
  "created_at": "2026-07-25T01:00:34.147+03:00",
  "updated_at": "2026-07-25T01:00:34.147+03:00",
  "deadline_at": "2021-01-01T00:00:00.000+03:00",
  "order": {
    "id": 1,
    "number": 1001
  }
}
```

### Delete user task

`DELETE /admin/user_tasks/:user_task_id.json`

**Запрос**

```http
DELETE /admin/user_tasks/12.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
Content-Type: application/x-www-form-urlencoded
```

**Ответ**

Status: `200`

### API-Usage-Limit: 1/500

`GET /admin/user_tasks/:user_task_id.json`

Content-Type: application/json Get user task

**Запрос**

```http
GET /admin/user_tasks/4.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 4,
  "complete": null,
  "order_id": null,
  "assignee_id": 1,
  "title": "dispatch order N1",
  "description": "You should dispatch the order N1 till noon",
  "created_at": "2026-07-25T01:00:33.780+03:00",
  "updated_at": "2026-07-25T01:00:33.780+03:00",
  "deadline_at": "2021-01-01T00:00:00.000+03:00"
}
```

### Get user tasks list

`GET /admin/user_tasks.json`

**Запрос**

```http
GET /admin/user_tasks.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 2,
    "complete": null,
    "order_id": null,
    "assignee_id": 1,
    "title": "dispatch order N1",
    "description": "You should dispatch the order N1 till noon",
    "created_at": "2026-07-25T01:00:33.558+03:00",
    "updated_at": "2026-07-25T01:00:33.558+03:00",
    "deadline_at": "2021-01-01T00:00:00.000+03:00"
  }
]
```

### Update user task

`PUT /admin/user_tasks/:user_task_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `user_task[title]` |  | user task title |
| `user_task[description]` |  | user task position |
| `user_task[complete]` |  | is user task completed? |

**Запрос**

```http
PUT /admin/user_tasks/10.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "user_task": {
    "title": "new title",
    "description": "new description",
    "complete": true
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 10,
  "complete": true,
  "order_id": null,
  "assignee_id": 1,
  "title": "new title",
  "description": "new description",
  "created_at": "2026-07-25T01:00:34.319+03:00",
  "updated_at": "2026-07-25T01:00:34.401+03:00",
  "deadline_at": "2021-01-01T00:00:00.000+03:00"
}
```

---

## Variant

| Операция | Метод и путь |
|----------|--------------|
| [Create variant](#create-variant) | `POST /admin/products/:product_id/variants.json` |
| [Destroy variant](#destroy-variant) | `DELETE /admin/products/:product_id/variants/:variant_success_id.json` |
| [Get variant](#get-variant) | `GET /admin/products/:product_id/variants/:variant_id.json` |
| [Get variants](#get-variants) | `GET /admin/products/:product_id/variants.json` |
| [Update group of variants (multi warehouses case)](#update-group-of-variants-multi-warehouses-case) | `PUT /admin/products/variants_group_update.json` |
| [Update group of variants](#update-group-of-variants) | `PUT /admin/products/variants_group_update.json` |
| [Update variant](#update-variant) | `PUT /admin/products/:product_id/variants/:variant_id.json` |

### Create variant

`POST /admin/products/:product_id/variants.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant[price]` | да | price |
| `variant[quantity]` | да | quantity |
| `variant[options]` |  | options |
| `variant[image_ids]` |  | image_ids |
| `variant[variant_field_values_attributes]` |  | variant field values attributes |

**Запрос**

```http
POST /admin/products/1/variants.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant": {
    "price": 100,
    "quantity": 1,
    "options": [
      {
        "option_name_id": 1,
        "value": "Black"
      }
    ],
    "image_ids": [
      115,
      123,
      131
    ],
    "variant_field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 509944661,
  "title": "Black",
  "product_id": 1,
  "sku": null,
  "barcode": null,
  "dimensions": "10x20x30",
  "available": true,
  "image_ids": [
    115,
    123,
    131
  ],
  "image_id": 115,
  "weight": null,
  "created_at": "2026-07-25T01:00:14.000+03:00",
  "updated_at": "2026-07-25T01:00:14.000+03:00",
  "quantity": 1,
  "quantity_at_warehouse0": "1.0",
  "quantity_at_warehouse1": "0.0",
  "cost_price": null,
  "cost_price_in_site_currency": null,
  "price_in_site_currency": "3000.0",
  "base_price": "100.0",
  "old_price": null,
  "price2": null,
  "price": "100.0",
  "base_price_in_site_currency": "3000.0",
  "old_price_in_site_currency": null,
  "price2_in_site_currency": null,
  "prices": [],
  "prices_in_site_currency": [],
  "variant_field_values": [
    {
      "id": 7,
      "variant_field_id": 1,
      "value": "12345"
    }
  ],
  "option_values": [
    {
      "id": 21063166,
      "option_name_id": 1,
      "position": 21,
      "title": "Black",
      "image_url": null,
      "permalink": "black"
    }
  ]
}
```

### Destroy variant

`DELETE /admin/products/:product_id/variants/:variant_success_id.json`

**Запрос**

```http
DELETE /admin/products/1/variants/509944663.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get variant

`GET /admin/products/:product_id/variants/:variant_id.json`

**Запрос**

```http
GET /admin/products/1/variants/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "title": "Red",
  "product_id": 1,
  "sku": "423",
  "barcode": null,
  "dimensions": "10x20x30",
  "available": false,
  "image_ids": [],
  "image_id": null,
  "weight": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2009-03-18T10:53:14.000+03:00",
  "quantity": 0,
  "quantity_at_warehouse0": "0.0",
  "quantity_at_warehouse1": "0.0",
  "cost_price": "3.0",
  "cost_price_in_site_currency": "90.0",
  "price_in_site_currency": "300.0",
  "base_price": "10.0",
  "old_price": null,
  "price2": 12.0,
  "price": "10.0",
  "base_price_in_site_currency": "300.0",
  "old_price_in_site_currency": null,
  "price2_in_site_currency": "360.0",
  "prices": [
    12.0
  ],
  "prices_in_site_currency": [
    "360.0"
  ],
  "variant_field_values": [
    {
      "id": 1,
      "variant_field_id": 1,
      "value": "some variant field value"
    }
  ],
  "option_values": [
    {
      "id": 1,
      "option_name_id": 1,
      "position": 1,
      "title": "Red",
      "image_url": null,
      "permalink": null
    }
  ]
}
```

### Get variants

`GET /admin/products/:product_id/variants.json`

**Запрос**

```http
GET /admin/products/1/variants.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "title": "Red",
    "product_id": 1,
    "sku": "423",
    "barcode": null,
    "dimensions": "10x20x30",
    "available": false,
    "image_ids": [],
    "image_id": null,
    "weight": null,
    "created_at": "2009-03-18T10:53:14.000+03:00",
    "updated_at": "2009-03-18T10:53:14.000+03:00",
    "quantity": 0,
    "quantity_at_warehouse0": "0.0",
    "quantity_at_warehouse1": "0.0",
    "cost_price": "3.0",
    "cost_price_in_site_currency": "90.0",
    "price_in_site_currency": "300.0",
    "base_price": "10.0",
    "old_price": null,
    "price2": 12.0,
    "price": "10.0",
    "base_price_in_site_currency": "300.0",
    "old_price_in_site_currency": null,
    "price2_in_site_currency": "360.0",
    "prices": [
      12.0
    ],
    "prices_in_site_currency": [
      "360.0"
    ],
    "variant_field_values": [
      {
        "id": 1,
        "variant_field_id": 1,
        "value": "some variant field value"
      }
    ],
    "option_values": [
      {
        "id": 1,
        "option_name_id": 1,
        "position": 1,
        "title": "Red",
        "image_url": null,
        "permalink": null
      }
    ]
  }
]
```

### Update group of variants (multi warehouses case)

`PUT /admin/products/variants_group_update.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variants` |  | Array with parameters for each variant  array |
| `variants[][id]` | да | Id of a variant  integer |
| `variants[][quantity]` |  | Variant quantity  float |
| `variants[][cost_price]` |  | Cost price of the variant  float |
| `variants[][old_price]` |  | Old price of the variant  float |
| `variants[][cost_price]` |  | Cost price of the variant  float |
| `variants[][price]` |  | First price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price2]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price3]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price4]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price5]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price6]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price7]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price8]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price9]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price10]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price11]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price12]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price13]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price14]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price15]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price16]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price17]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price18]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price19]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price20]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][quantity_at_warehouses]` |  | Array of quantity at each of the warehouses. Could be replaced with quantity_at_warehousesN fields, which have higher priority  array  {"type" => "float"} |
| `variants[][quantity_at_warehouse0]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse1]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse2]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse3]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse4]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse5]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse6]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse7]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse8]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse9]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse10]` |  | Quantity at n-th warehouse  float |
| `variants[][quantity_at_warehouse11]` |  | Quantity at n-th warehouse  float |

**Запрос**

```http
PUT /admin/products/variants_group_update.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variants": [
    {
      "id": 1,
      "price": 100,
      "quantity_at_warehouse0": 3,
      "quantity_at_warehouses": [
        5,
        2
      ]
    }
  ]
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[]
```

### Update group of variants

`PUT /admin/products/variants_group_update.json`

It’s possible to update at most 100 variants per request.

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variants` |  | Array with parameters for each variant  array |
| `variants[][id]` | да | Id of a variant  integer |
| `variants[][quantity]` |  | Variant quantity  float |
| `variants[][cost_price]` |  | Cost price of the variant  float |
| `variants[][old_price]` |  | Old price of the variant  float |
| `variants[][cost_price]` |  | Cost price of the variant  float |
| `variants[][price]` |  | First price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price2]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price3]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price4]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price5]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price6]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price7]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price8]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price9]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price10]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price11]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price12]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price13]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price14]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price15]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price16]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price17]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price18]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price19]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |
| `variants[][price20]` |  | n-th price of the price kinds of the variant (if price kinds are enabled)  float |

**Запрос**

```http
PUT /admin/products/variants_group_update.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variants": [
    {
      "id": 1,
      "price": 100,
      "quantity_at_warehouse0": 3
    }
  ]
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[]
```

### Update variant

`PUT /admin/products/:product_id/variants/:variant_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant[dimensions]` |  | Length x Width x Height (Ш x Г x В) |
| `variant[price]` |  | price |
| `variant[quantity_at_warehouse0]` |  | quantity at warehouse with array index 0 (require multiwarehouses enabled) |
| `variant[variant_field_values_attributes]` |  | variant field values attributes |

**Запрос**

```http
PUT /admin/products/1/variants/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant": {
    "price": 150,
    "quantity_at_warehouse0": 50,
    "variant_field_values_attributes": [
      {
        "handle": "handle",
        "value": 12345
      }
    ]
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "title": "Red",
  "product_id": 1,
  "sku": "423",
  "barcode": null,
  "dimensions": "10x20x30",
  "available": true,
  "image_ids": [],
  "image_id": null,
  "weight": null,
  "created_at": "2009-03-18T10:53:14.000+03:00",
  "updated_at": "2026-07-25T01:00:15.000+03:00",
  "quantity": 50,
  "quantity_at_warehouse0": "50.0",
  "quantity_at_warehouse1": "0.0",
  "cost_price": "3.0",
  "cost_price_in_site_currency": "90.0",
  "price_in_site_currency": "4500.0",
  "base_price": "150.0",
  "old_price": null,
  "price2": 12.0,
  "price": "150.0",
  "base_price_in_site_currency": "4500.0",
  "old_price_in_site_currency": null,
  "price2_in_site_currency": "360.0",
  "prices": [
    12.0
  ],
  "prices_in_site_currency": [
    "360.0"
  ],
  "variant_field_values": [
    {
      "id": 1,
      "variant_field_id": 1,
      "value": "12345"
    }
  ],
  "option_values": [
    {
      "id": 1,
      "option_name_id": 1,
      "position": 1,
      "title": "Red",
      "image_url": null,
      "permalink": null
    }
  ]
}
```

---

## VariantField

| Операция | Метод и путь |
|----------|--------------|
| [Create variant_field](#create-variant-field) | `POST /admin/variant_fields.json` |
| [Destroy variant_field by handle](#destroy-variant-field-by-handle) | `DELETE /admin/variant_fields/:variant_field_handle.json` |
| [Destroy variant_field](#destroy-variant-field) | `DELETE /admin/variant_fields/:variant_field_id.json` |
| [Get variant_field by handle](#get-variant-field-by-handle) | `GET /admin/variant_fields/:variant_field_handle.json` |
| [Get variant_field](#get-variant-field) | `GET /admin/variant_fields/:variant_field_id.json` |
| [Get variant_fields](#get-variant-fields) | `GET /admin/variant_fields.json` |
| [Update variant_field by handle](#update-variant-field-by-handle) | `PUT /admin/variant_fields/:variant_field_handle.json` |
| [Update variant_field](#update-variant-field) | `PUT /admin/variant_fields/:variant_field_id.json` |

### Create variant_field

`POST /admin/variant_fields.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field[title]` | да | title |
| `variant_field[handle]` | да | handle used in liquid |
| `variant_field[type]` | да | VariantField::TextArea - textarea, VariantField::TextField - textfield |
| `variant_field[position]` |  | position in variant_fields list |
| `variant_field[is_hidden]` |  | boolean |

**Запрос**

```http
POST /admin/variant_fields.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant_field": {
    "title": "Size",
    "handle": "size",
    "type": "VariantField::TextArea"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 5,
  "handle": "size",
  "type": "VariantField::TextArea",
  "application_id": 1,
  "position": 3,
  "is_hidden": null,
  "title": "Size",
  "application_title": "Цифровые товары"
}
```

### Destroy variant_field by handle

`DELETE /admin/variant_fields/:variant_field_handle.json`

**Запрос**

```http
DELETE /admin/variant_fields/handle2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Destroy variant_field

`DELETE /admin/variant_fields/:variant_field_id.json`

**Запрос**

```http
DELETE /admin/variant_fields/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get variant_field by handle

`GET /admin/variant_fields/:variant_field_handle.json`

**Запрос**

```http
GET /admin/variant_fields/handle2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "VariantField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "title": "title2",
  "application_title": "Цифровые товары"
}
```

### Get variant_field

`GET /admin/variant_fields/:variant_field_id.json`

**Запрос**

```http
GET /admin/variant_fields/2.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "VariantField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "title": "title2",
  "application_title": "Цифровые товары"
}
```

### Get variant_fields

`GET /admin/variant_fields.json`

**Запрос**

```http
GET /admin/variant_fields.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "handle": "handle",
    "type": "VariantField::TextField",
    "application_id": null,
    "position": 1,
    "is_hidden": null,
    "title": "title",
    "application_title": null
  },
  {
    "id": 2,
    "handle": "handle2",
    "type": "VariantField::TextArea",
    "application_id": 1,
    "position": 2,
    "is_hidden": null,
    "title": "title2",
    "application_title": "Цифровые товары"
  }
]
```

### Update variant_field by handle

`PUT /admin/variant_fields/:variant_field_handle.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field[title]` |  | title |
| `variant_field[handle]` |  | handle used in liquid  new_handle |
| `variant_field[position]` |  | position in variant_fields list |
| `variant_field[is_hidden]` |  | boolean |

**Запрос**

```http
PUT /admin/variant_fields/handle2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant_field": {
    "title": "Volume",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "handle2",
  "type": "VariantField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "title": "Volume",
  "application_title": "Цифровые товары"
}
```

### Update variant_field

`PUT /admin/variant_fields/:variant_field_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field[title]` |  | title |
| `variant_field[handle]` |  | handle used in liquid |
| `variant_field[position]` |  | position in variant_fields list |
| `variant_field[is_hidden]` |  | boolean |

**Запрос**

```http
PUT /admin/variant_fields/2.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant_field": {
    "title": "Volume",
    "handle": "volume",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 2,
  "handle": "volume",
  "type": "VariantField::TextArea",
  "application_id": 1,
  "position": 2,
  "is_hidden": null,
  "title": "Volume",
  "application_title": "Цифровые товары"
}
```

---

## VariantFieldValue

| Операция | Метод и путь |
|----------|--------------|
| [Create variant_field_value by handle](#create-variant-field-value-by-handle) | `POST /admin/variants/:variant_id/variant_field_values.json` |
| [Create variant_field_value](#create-variant-field-value) | `POST /admin/variants/:variant_id/variant_field_values.json` |
| [Destroy variant_field_value by handle](#destroy-variant-field-value-by-handle) | `DELETE /admin/variants/:variant_id/variant_field_values/:variant_field_handle.json` |
| [Destroy variant_field_value](#destroy-variant-field-value) | `DELETE /admin/variants/:variant_id/variant_field_values/:variant_field_value_id.json` |
| [Get variant_field_value by handle](#get-variant-field-value-by-handle) | `GET /admin/variants/:variant_id/variant_field_values/:variant_field_handle.json` |
| [Get variant_field_value](#get-variant-field-value) | `GET /admin/variants/:variant_id/variant_field_values/:variant_field_value_id.json` |
| [Get variant_field_values](#get-variant-field-values) | `GET /admin/variants/:variant_id/variant_field_values.json` |
| [Update variant_field_value by handle](#update-variant-field-value-by-handle) | `PUT /admin/variants/:variant_id/variant_field_values/:variant_field_handle.json` |
| [Update variant_field_value](#update-variant-field-value) | `PUT /admin/variants/:variant_id/variant_field_values/:variant_field_value_id.json` |

### Create variant_field_value by handle

`POST /admin/variants/:variant_id/variant_field_values.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field_value[variant_field_handle]` | да | variant field handle |
| `variant_field_value[value]` |  | variant field value |

**Запрос**

```http
POST /admin/variants/1/variant_field_values.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant_field_value": {
    "variant_field_handle": "handle2",
    "value": "test value"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 5,
  "variant_field_id": 2,
  "value": "test value"
}
```

### Create variant_field_value

`POST /admin/variants/:variant_id/variant_field_values.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field_value[variant_field_id]` | да | variant field id |
| `variant_field_value[value]` |  | variant field value |

**Запрос**

```http
POST /admin/variants/1/variant_field_values.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant_field_value": {
    "variant_field_id": 2,
    "value": "test value"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 3,
  "variant_field_id": 2,
  "value": "test value"
}
```

### Destroy variant_field_value by handle

`DELETE /admin/variants/:variant_id/variant_field_values/:variant_field_handle.json`

**Запрос**

```http
DELETE /admin/variants/1/variant_field_values/handle.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Destroy variant_field_value

`DELETE /admin/variants/:variant_id/variant_field_values/:variant_field_value_id.json`

**Запрос**

```http
DELETE /admin/variants/1/variant_field_values/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get variant_field_value by handle

`GET /admin/variants/:variant_id/variant_field_values/:variant_field_handle.json`

**Запрос**

```http
GET /admin/variants/1/variant_field_values/handle.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "variant_field_id": 1,
  "value": "some variant field value"
}
```

### Get variant_field_value

`GET /admin/variants/:variant_id/variant_field_values/:variant_field_value_id.json`

**Запрос**

```http
GET /admin/variants/1/variant_field_values/1.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "variant_field_id": 1,
  "value": "some variant field value"
}
```

### Get variant_field_values

`GET /admin/variants/:variant_id/variant_field_values.json`

**Запрос**

```http
GET /admin/variants/1/variant_field_values.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "variant_field_id": 1,
    "value": "some variant field value"
  }
]
```

### Update variant_field_value by handle

`PUT /admin/variants/:variant_id/variant_field_values/:variant_field_handle.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field_value[value]` |  | variant field value |

**Запрос**

```http
PUT /admin/variants/1/variant_field_values/handle.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant_field_value": {
    "value": "new value"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "variant_field_id": 1,
  "value": "new value"
}
```

### Update variant_field_value

`PUT /admin/variants/:variant_id/variant_field_values/:variant_field_value_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field_value[value]` |  | variant field value |

**Запрос**

```http
PUT /admin/variants/1/variant_field_values/1.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "variant_field_value": {
    "value": "new value"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 1,
  "variant_field_id": 1,
  "value": "new value"
}
```

---

## Warehouse

| Операция | Метод и путь |
|----------|--------------|
| [Create warehouse](#create-warehouse) | `POST /admin/warehouses.json` |
| [Delete warehouse](#delete-warehouse) | `DELETE /admin/warehouses/:warehouse_id.json` |
| [API-Usage-Limit: 1/500](#api-usage-limit-1-500) | `GET /admin/warehouses/:warehouse_id.json` |
| [Get warehouses list](#get-warehouses-list) | `GET /admin/warehouses.json` |
| [Update warehouse](#update-warehouse) | `PUT /admin/warehouses/:warehouse_id.json` |

### Create warehouse

`POST /admin/warehouses.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `warehouse[title]` | да | warehouse title |
| `warehouse[position]` |  | warehouse position (1 - main priority) |

**Запрос**

```http
POST /admin/warehouses.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "warehouse": {
    "title": "secondary",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 96,
  "title": "secondary",
  "position": 2,
  "array_index": 2,
  "created_at": "2026-07-25T01:00:35.320+03:00",
  "updated_at": "2026-07-25T01:00:35.320+03:00",
  "status": "creating"
}
```

### Delete warehouse

`DELETE /admin/warehouses/:warehouse_id.json`

**Запрос**

```http
DELETE /admin/warehouses/100.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
Content-Type: application/x-www-form-urlencoded
```

**Ответ**

Status: `202`

### API-Usage-Limit: 1/500

`GET /admin/warehouses/:warehouse_id.json`

Content-Type: application/json Get warehouse

**Запрос**

```http
GET /admin/warehouses/92.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 92,
  "title": "main",
  "position": 2,
  "array_index": 1,
  "created_at": "2026-07-25T01:00:35.022+03:00",
  "updated_at": "2026-07-25T01:00:35.022+03:00",
  "status": "creating"
}
```

### Get warehouses list

`GET /admin/warehouses.json`

**Запрос**

```http
GET /admin/warehouses.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "id": 1,
    "title": "Главный склад",
    "position": 1,
    "array_index": 0,
    "created_at": "2008-12-31T03:00:00.000+03:00",
    "updated_at": "2008-12-31T03:00:00.000+03:00",
    "status": "ready"
  },
  {
    "id": 90,
    "title": "main",
    "position": 2,
    "array_index": 1,
    "created_at": "2026-07-25T01:00:34.812+03:00",
    "updated_at": "2026-07-25T01:00:34.812+03:00",
    "status": "creating"
  }
]
```

### Update warehouse

`PUT /admin/warehouses/:warehouse_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `warehouse[title]` |  | warehouse title |
| `warehouse[position]` |  | warehouse position |

**Запрос**

```http
PUT /admin/warehouses/98.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "warehouse": {
    "title": "new title",
    "position": 2
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "id": 98,
  "title": "new title",
  "position": 2,
  "array_index": 1,
  "created_at": "2026-07-25T01:00:35.470+03:00",
  "updated_at": "2026-07-25T01:00:35.541+03:00",
  "status": "creating"
}
```

---

## Webhook

| Операция | Метод и путь |
|----------|--------------|
| [Create webhook](#create-webhook) | `POST /admin/webhooks.json` |
| [Destroy webhook](#destroy-webhook) | `DELETE /admin/webhooks/:webhook_id.json` |
| [Get webhook](#get-webhook) | `GET /admin/webhooks/:webhook_id.json` |
| [Get webhooks](#get-webhooks) | `GET /admin/webhooks.json` |
| [Update webhook](#update-webhook) | `PUT /admin/webhooks/:webhook_id.json` |

### Create webhook

`POST /admin/webhooks.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `webhook[address]` | да | address |
| `webhook[topic]` | да | Topics: 'orders/create', 'orders/update', 'orders/destroy', 'products/create', 'products/update', 'client/create', 'client/update' |
| `webhook[format_type]` |  | format of data sent by webhook: 'json' or 'xml' |
| `webhook[warehouse_id]` |  | For 'orders/create', 'orders/update', 'orders/destroy' topics only, filters sent orders by warehouse |
| `webhook[sales_channel_id]` |  | For 'products/update' topic only, filters sent products by sales channel |
| `webhook[batch_size]` |  | For 'products/create', 'products/update' topics only, products are sent in batches, default value: 10 |

**Запрос**

```http
POST /admin/webhooks.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "webhook": {
    "address": "http://application.com/orders/create",
    "topic": "orders/create",
    "format_type": "json"
  }
}
```

**Ответ**

Status: `201`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "address": "http://application.com/orders/create",
  "created_at": "2026-07-25T00:58:46.166+03:00",
  "id": 13,
  "topic": "orders/create",
  "format_type": "json",
  "warehouse_id": null,
  "sales_channel_id": null,
  "batch_size": null
}
```

### Destroy webhook

`DELETE /admin/webhooks/:webhook_id.json`

**Запрос**

```http
DELETE /admin/webhooks/17.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "status": "ok"
}
```

### Get webhook

`GET /admin/webhooks/:webhook_id.json`

**Запрос**

```http
GET /admin/webhooks/9.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "address": "http://app.ru/orders/create",
  "created_at": "2026-07-25T00:58:45.884+03:00",
  "id": 9,
  "topic": "orders/create",
  "format_type": "xml",
  "warehouse_id": null,
  "sales_channel_id": null,
  "batch_size": null
}
```

### Get webhooks

`GET /admin/webhooks.json`

**Запрос**

```http
GET /admin/webhooks.json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
[
  {
    "address": "https://webhook.site/4f4dd3e6-46d8-403f-8a3e-4fb0b1f108a0",
    "created_at": "2022-05-05T20:37:45.000+03:00",
    "id": 1,
    "topic": "products/update",
    "format_type": "json",
    "warehouse_id": null,
    "sales_channel_id": null,
    "batch_size": 5
  },
  {
    "address": "https://webhook.site/4f4dd3e6-46d8-403f-8a3e-4fb0b1f108a0",
    "created_at": "2022-05-05T20:37:45.000+03:00",
    "id": 3,
    "topic": "products/create",
    "format_type": "json",
    "warehouse_id": null,
    "sales_channel_id": null,
    "batch_size": 5
  },
  {
    "address": "https://webhook.site/4f4dd3e6-46d8-403f-8a3e-4fb0b1f108a0",
    "created_at": "2022-05-05T20:37:45.000+03:00",
    "id": 4,
    "topic": "client/update",
    "format_type": "json",
    "warehouse_id": null,
    "sales_channel_id": null,
    "batch_size": null
  },
  {
    "address": "https://webhook.site/4f4dd3e6-46d8-403f-8a3e-4fb0b1f108a0",
    "created_at": "2022-05-05T20:37:45.000+03:00",
    "id": 5,
    "topic": "client/create",
    "format_type": "json",
    "warehouse_id": null,
    "sales_channel_id": null,
    "batch_size": null
  },
  {
    "address": "http://app.ru/orders/create",
    "created_at": "2026-07-25T00:58:45.682+03:00",
    "id": 7,
    "topic": "orders/create",
    "format_type": "xml",
    "warehouse_id": null,
    "sales_channel_id": null,
    "batch_size": null
  }
]
```

### Update webhook

`PUT /admin/webhooks/:webhook_id.json`

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `webhook[address]` |  | address |
| `webhook[topic]` |  | Topics: 'orders/create', 'orders/update', 'orders/destroy', 'products/create', 'products/update', 'client/create', 'client/update' |
| `webhook[format_type]` |  | format of data sent by webhook: json or xml |
| `webhook[warehouse_id]` |  | For 'orders/create', 'orders/update', 'orders/destroy' topics only, filters sent orders by warehouse |
| `webhook[sales_channel_id]` |  | For 'products/update' topic only, filters sent products by sales channel |
| `webhook[batch_size]` |  | For 'products/create', 'products/update' topics only, products are sent in batches, default value: 10 |

**Запрос**

```http
PUT /admin/webhooks/15.json
Content-Type: application/json
Authorization: Basic ZGlnaXRhbC1nb29kczpjNTFjOTA3MDdhMTNjZTNmZmYyMTNhZmJiNWNkMTI3MA==
```

```json
{
  "webhook": {
    "format_type": "xml"
  }
}
```

**Ответ**

Status: `200`

```http
API-Usage-Limit: 1/500
Content-Type: application/json; charset=utf-8
```

```json
{
  "address": "http://app.ru/orders/create",
  "created_at": "2026-07-25T00:58:46.308+03:00",
  "id": 15,
  "topic": "orders/create",
  "format_type": "xml",
  "warehouse_id": null,
  "sales_channel_id": null,
  "batch_size": null
}
```

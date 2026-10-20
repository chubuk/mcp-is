# MCP IS — Архитектура

## Обзор системы

```
┌─────────────────────────────────────────────────────────────┐
│                     AI-клиенты                              │
│  Claude Desktop │ Cursor │ VS Code │ LibreChat │ Open WebUI │
└────────────────────────┬────────────────────────────────────┘
                         │ MCP (Streamable HTTP)
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    nginx (порт 443)                         │
│  • SSL/TLS (Let's Encrypt или self-signed)                  │
│  • Reverse proxy: /api/* → Go, /mcp → Go, /* → static      │
│  • Rate limiting, security headers                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│               Go API Server (порт 8080)                     │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Auth Layer  │  │ MCP Server  │  │ Web API (REST)      │ │
│  │ JWT/API Key │  │ mcp-go SDK  │  │ handlers/shops/...  │ │
│  │ Rate Limit  │  │ 282 tools   │  │ admin/auth/server   │ │
│  └─────────────┘  └──────┬──────┘  └──────────┬──────────┘ │
│                          │                     │            │
│  ┌───────────────────────┴─────────────────────┴──────────┐ │
│  │              Service Layer (service.go)                 │ │
│  └───────────────────────┬────────────────────────────────┘ │
│                          │                                  │
│  ┌───────────────────────┴────────────────────────────────┐ │
│  │         AccountManager (core/manager.go)                │ │
│  │  • Shop selection (default or by ID)                    │ │
│  │  • Read-only enforcement                                │ │
│  │  • Context propagation (shop_name, key_name)            │ │
│  └───────────────────────┬────────────────────────────────┘ │
└──────────────────────────┼──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              InSales API Client                             │
│  • HTTP Basic Auth                                          │
│  • AES-GCM encrypted passwords                              │
│  • InSales REST (products, orders, fields, …)               │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL 15                                  │
│  users • shops • ai_keys • server_settings • password_resets│
└─────────────────────────────────────────────────────────────┘
```

## Слои системы

### 1. Transport Layer
- **Streamable HTTP** (mcp-go SDK) — основной транспорт для MCP
- **SSE** — для legacy клиентов
- **CORS** — разрешённые кросс-доменные запросы

### 2. Auth Layer
- **JWT** — для веб-интерфейса (срок 24 часа)
- **API Keys** — для MCP-клиентов (SHA-256 хэши)
- **Rate Limiting** — 10 запросов/мин на login/register/forgot-password
- **Role-based access** — super_admin, admin, user

### 3. MCP Server (mcp-go SDK)
- **282 инструмента** для работы с InSales (120 read-only / 162 write)
- **Progressive discovery** — default `tools/list` ≈ core surface (~30) + `search_insales_tools` / `call_insales_tool`; полный каталог через catalog или `MCP_TOOLS_MODE=all`
- **Tool annotations** — readOnlyHint, destructiveHint, idempotentHint, openWorldHint
- **Tool filtering** — read-only режим скрывает write-команды; nil ReadOnlyHint трактуется как write
- **Server instructions** — discovery, envelope, language-invariant args, pipelines
- **Результаты** — JSON text + structuredContent (двойная форма)

### 4. Service Layer
- **Service** — бизнес-логика, вызов InSales API
- **AccountManager** — выбор магазина, read-only enforcement
- **Context propagation** — shop_name, key_name, shop_url в контексте

### 5. Data Layer
- **PostgreSQL 15** — пользователи, магазины, ключи, настройки
- **AES-GCM** — шифрование паролей InSales
- **CASCADE deletes** — каскадное удаление данных

## Потоки данных

### MCP-запрос от AI

```
AI → HTTP POST /mcp → API Key Auth → MCP Server → Service → AccountManager → InSales API
                                                        ↓
                                                   PostgreSQL (user, shop, key)
```

### Регистрация пользователя

```
Web UI → HTTP POST /api/register → Auth Handler → PostgreSQL
```

Публичная регистрация закрыта, если `REGISTRATION_ENABLED=false` (в `.env`) и в базе уже есть пользователи. Создание: `POST /api/admin/users` (admin/super_admin) — временный пароль, `require_password_change=true`.

### Добавление магазина

```
Web UI → HTTP POST /api/shops → Shop Handler → AES encrypt → PostgreSQL
```

## Безопасность

| Механизм | Реализация |
|----------|-----------|
| Пароли пользователей | bcrypt |
| Пароли InSales | AES-GCM (ENCRYPTION_KEY) |
| AI-ключи | SHA-256 хэши (сами ключи не хранятся) |
| JWT | HMAC-SHA256 (срок 24 часа) |
| Rate limiting | In-memory sliding window, 10 req/min |
| CORS | Разрешённые origins |
| DNS rebinding | Автоматическая защита (mcp-go) |
| HTTPS | Let's Encrypt или self-signed |

## Read-only режим

```
Глобальный (MCP_READ_ONLY=true в .env)
    ↓ наивысший приоритет
Per-shop (is_read_only в shops)
    ↓ если глобальный НЕ установлен
Tool filtering (tools/list возвращает только read-команды)
```

## Деплой

### Docker Compose

| Сервис | Описание | Порты |
|--------|----------|-------|
| app | Go API + MCP Server | 8080 (internal) |
| db | PostgreSQL 15 | 5432 |
| nginx | Reverse proxy + SSL | 80/443 |

### На VPS

```bash
# Установка
sudo dpkg -i mcp-is_0.1.0_amd64.deb
sudo nano /etc/mcp-is/.env
sudo systemctl start mcp-is
```

### 6. Покрытие InSales API → MCP tools (op→tool)

Сгенерировано из `docs/api/Insales-api-docs.md` vs `internal/mcp/tools.go`. Ресурсов: 57, операций API: 327, MCP tools: 282.

Статусы: **OK** — есть list/get/create/update/delete (или API не даёт большего); **частично** — нет отдельных вербов (часто в API их и нет); **out-of-scope** — ресурс не нужен / сломан на стороне InSales.

| Ресурс InSales | Ops | Methods | MCP tools | Статус |
|----------------|----:|---------|-----------|--------|
| Accessories | 5 | DELETE, GET, POST, PUT | `create_accessory`, `delete_accessory`, `get_accessory`, `list_accessories`, `list_product_accessory_links`, `update_accessory` | OK |
| AccessoryValues | 8 | DELETE, GET, POST, PUT | `create_property_value`, `delete_property_value`, `list_property_values`, `update_property_value` | OK |
| Account | 2 | GET, PUT | `get_shop_info`, `update_account` | OK (get_shop_info + update_account) |
| ApplicationAction | 5 | DELETE, GET, POST, PUT | `create_application_action`, `delete_application_action`, `get_application_action`, `update_application_action` | OK |
| ApplicationCharge | 4 | GET, POST | `create_application_charge`, `create_recurring_application_charge`, `decline_application_charge`, `delete_recurring_application_charge`, `get_application_charge`, `get_recurring_application_charge` | частично: нет update |
| ApplicationWidget | 5 | DELETE, GET, POST, PUT | `create_application_widget`, `delete_application_widget`, `get_application_widget`, `update_application_widget` | OK |
| Article | 5 | DELETE, GET, POST, PUT | `create_blog_post`, `delete_blog_post`, `get_blog_post`, `list_blog_posts`, `update_blog_post` | OK (API InSales местами 404/500) |
| Blog | 5 | DELETE, GET, POST, PUT | `create_blog`, `create_blog_post`, `delete_blog`, `delete_blog_post`, `get_blog`, `get_blog_post`, `list_blog_posts`, `update_blog`… | OK |
| BonusSystemTransaction | 4 | GET, POST | `create_bonus_system_transaction`, `get_bonus_system_transaction` | частично: нет update,delete |
| Category | 5 | DELETE, GET, POST, PUT | `create_category`, `delete_category`, `get_category`, `update_category` | OK |
| Characteristic | 5 | DELETE, GET, POST, PUT | `create_characteristic`, `delete_characteristic`, `get_characteristic`, `update_characteristic` | OK |
| Client | 6 | DELETE, GET, POST, PUT | `client_full_profile`, `create_client`, `create_client_group`, `delete_client`, `delete_client_group`, `get_client`, `get_client_group`, `list_client_bonus_transactions`… | OK |
| ClientGroup | 5 | DELETE, GET, POST, PUT | `create_client_group`, `delete_client_group`, `get_client_group`, `update_client_group` | OK |
| Collect | 7 | DELETE, GET, POST, PUT | `add_product_to_collection`, `get_collects_by_product`, `list_collects` | частично: нет update,delete |
| Collection | 5 | DELETE, GET, POST, PUT | `add_product_to_collection`, `collection_products`, `create_collection`, `create_collection_field`, `create_collection_filter`, `delete_collection`, `delete_collection_field`, `delete_collection_filter`… | OK |
| CollectionField | 8 | DELETE, GET, POST, PUT | `create_collection_field`, `delete_collection_field`, `get_collection_field`, `get_collection_field_value`, `list_collection_field_values`, `update_collection_field`, `update_collection_field_value` | OK |
| CollectionFieldValue | 5 | GET, PUT | `get_collection_field_value`, `update_collection_field_value` | OK (API: только list/get/update) |
| CollectionFilter | 5 | DELETE, GET, POST, PUT | `create_collection_filter`, `delete_collection_filter`, `get_collection_filter`, `update_collection_filter` | OK |
| CustomStatus | 5 | DELETE, GET, POST, PUT | `create_custom_status`, `delete_custom_status`, `get_custom_status`, `update_custom_status`, `update_order_custom_status` | OK |
| DeliveryVariant | 9 | DELETE, GET, POST, PUT | `create_delivery_variant`, `delete_delivery_variant`, `get_delivery_variant`, `update_delivery_variant` | OK |
| DiscountCode | 5 | DELETE, GET, POST, PUT | `create_discount_code`, `delete_discount_code`, `get_discount_code`, `update_discount_code` | OK |
| DiscountCollectionsGroup | 3 | GET | `create_discount_collections_group`, `delete_discount_collections_group`, `get_discount_collections_group`, `update_discount_collections_group` | OK |
| Domain | 5 | DELETE, GET, POST, PUT | `create_domain`, `delete_domain`, `get_domain`, `update_domain` | OK |
| ExternalDiscount | 5 | DELETE, GET, POST, PUT | `create_external_discount`, `delete_external_discount`, `get_external_discount`, `update_external_discount` | OK |
| Field | 7 | DELETE, GET, POST, PUT | — | out-of-scope: глобальные Field ≠ ProductField |
| File | 5 | DELETE, GET, POST | `create_file`, `delete_file`, `get_file` | частично: нет update |
| Image | 6 | DELETE, GET, POST, PUT | `add_image`, `delete_image`, `list_images` | OK (list/add/delete — без update в API-практике) |
| JsTag | 6 | DELETE, GET, POST, PUT | `create_js_tag`, `delete_js_tag`, `get_js_tag`, `update_js_tag` | OK |
| Marketplace | 9 | DELETE, GET, POST, PUT | `create_marketplace`, `delete_marketplace`, `get_marketplace`, `update_marketplace` | OK |
| Menu | 5 | DELETE, GET, POST, PUT | `create_menu`, `create_menu_item`, `delete_menu`, `delete_menu_item`, `get_menu`, `get_menu_item`, `list_menu_items`, `update_menu`… | OK |
| MenuItem | 8 | DELETE, GET, POST, PUT | `create_menu_item`, `delete_menu_item`, `get_menu_item`, `list_menu_items`, `update_menu_item` | OK |
| OAuthApp | 5 | DELETE, GET, POST, PUT | `create_oauth_app`, `delete_oauth_app`, `get_oauth_app`, `list_oauth_apps`, `update_oauth_app` | OK |
| OptionName | 5 | DELETE, GET, POST, PUT | `create_option_name`, `delete_option_name`, `get_option_name`, `update_option_name` | OK |
| OptionValue | 8 | DELETE, GET, POST, PUT | `create_option_value`, `delete_option_value`, `get_option_value`, `update_option_value` | OK |
| Order | 16 | GET, POST, PUT | `add_order_discount`, `add_order_line`, `create_order`, `get_order`, `order_full_details`, `remove_order_discount`, `remove_order_line`, `update_order`… | OK |
| Outlet | 5 | DELETE, GET, POST, PUT | `create_outlet`, `delete_outlet`, `get_outlet`, `update_outlet` | OK |
| Page | 5 | DELETE, GET, POST, PUT | `create_page`, `delete_page`, `get_page`, `update_page` | OK |
| PaymentGateway | 6 | DELETE, GET, POST, PUT | `create_payment_gateway`, `delete_payment_gateway`, `get_payment_gateway`, `update_payment_gateway` | OK |
| PriceKind | 5 | DELETE, GET, POST, PUT | `create_price_kind`, `delete_price_kind`, `get_price_kind`, `update_price_kind` | OK |
| Product | 8 | DELETE, GET, POST, PUT | `add_product_to_collection`, `create_product`, `create_product_field`, `create_product_field_value`, `delete_product`, `delete_product_field`, `delete_product_field_value`, `get_collects_by_product`… | OK |
| ProductField | 8 | DELETE, GET, POST, PUT | `create_product_field`, `create_product_field_value`, `delete_product_field`, `delete_product_field_value`, `get_product_field`, `get_product_field_value`, `list_product_field_values`, `update_product_field`… | OK |
| ProductFieldValue | 9 | DELETE, GET, POST, PUT | `create_product_field_value`, `delete_product_field_value`, `get_product_field_value`, `update_product_field_value` | OK |
| Property | 5 | DELETE, GET, POST, PUT | `create_property`, `create_property_value`, `delete_property`, `delete_property_value`, `get_property`, `list_property_values`, `update_property`, `update_property_value` | OK |
| ProviderToken | 1 | POST | `create_provider_token` | OK,update,delete |
| RecurringApplicationCharge | 5 | DELETE, GET, PATCH, POST, PUT | `create_recurring_application_charge`, `delete_recurring_application_charge`, `get_recurring_application_charge` | частично: нет update |
| Redirect | 5 | DELETE, GET, POST | `create_redirect`, `delete_redirect` | OK (create/delete/list — без update в API) |
| Review | 5 | DELETE, GET, POST, PUT | `create_review`, `delete_review`, `get_review`, `update_review` | OK |
| SalesChannel | 4 | DELETE, GET, POST | `create_sales_channel`, `delete_sales_channel`, `get_sales_channel`, `update_sales_channel` | OK |
| Similar | 7 | DELETE, GET, POST, PUT | `create_similar`, `delete_similar` | OK,update |
| Supplementary | 3 | DELETE, GET, POST | `create_supplementary`, `delete_supplementary`, `get_supplementary` | частично: нет update |
| Unit | 1 | GET | — | OK (в API только list) |
| UserTask | 5 | DELETE, GET, POST, PUT | `create_user_task`, `delete_user_task`, `get_user_task`, `update_user_task` | OK |
| Variant | 7 | DELETE, GET, POST, PUT | `create_delivery_variant`, `create_variant_field`, `create_variant_field_value`, `delete_delivery_variant`, `delete_variant`, `delete_variant_field`, `delete_variant_field_value`, `get_delivery_variant`… | OK |
| VariantField | 8 | DELETE, GET, POST, PUT | `create_variant_field`, `create_variant_field_value`, `delete_variant_field`, `delete_variant_field_value`, `get_variant_field`, `get_variant_field_value`, `list_variant_field_values`, `update_variant_field`… | OK |
| VariantFieldValue | 9 | DELETE, GET, POST, PUT | `create_variant_field_value`, `delete_variant_field_value`, `get_variant_field_value`, `update_variant_field_value` | OK |
| Warehouse | 5 | DELETE, GET, POST, PUT | `create_warehouse`, `delete_warehouse`, `get_warehouse`, `update_warehouse` | OK |
| Webhook | 5 | DELETE, GET, POST, PUT | `create_webhook`, `delete_webhook`, `update_webhook` | OK |

Примечания:
- `Article`/`Page` — ограничения API InSales (404/500), не баги mcp-is.
- `CollectionFieldValue` — в API только GET/PUT.
- `Field` (destiny 1–6) — отдельный ресурс от ProductField/VariantField/CollectionField.
- Нестыковки имён: `AccessoryValues` = `*_property_values`, `Article` = `blog_posts`.

# MCP IS

MCP-сервер (Model Context Protocol) для интеграции AI-агентов с интернет-магазином на платформе InSales.

## Описание

Сервер предоставляет AI-агентам (Claude Desktop, Cursor и другим MCP-клиентам) набор инструментов для управления магазином InSales через REST API. Поддерживает чтение и запись данных с настраиваемым режимом доступа.

## Архитектура

```
cmd/server/main.go          — точка входа
internal/
  config/                   — загрузка конфигурации из .env
  core/                     — бизнес-логика, менеджер аккаунтов
  handlers/                 — HTTP-обработчики (auth, shops, aikeys)
  integration/insales/      — HTTP-клиент для InSales API
  mcp/                      — обработка MCP-протокола (JSON-RPC 2.0)
  service/                  — слой сервисов
  storage/                  — работа с PostgreSQL (репозитории)
  auth/                     — JWT, AES-GCM, API-ключи, middleware
  logger/                   — структурное логирование (slog)
  metrics/                  — Prometheus метрики
  resilience/               — circuit breaker, retry с backoff
  webhooks/                 — обработка вебхуков InSales
web/
  nginx/                    — Dockerfile, nginx.conf, entrypoint.sh
  static/                   — HTML, CSS, JS (SPA фронтенд)
scripts/
  backup.sh                 — автоматический бэкап PostgreSQL
```

## Технологии

- **Язык**: Go 1.25
- **БД**: PostgreSQL 15
- **Контейнеризация**: Docker (multi-stage сборка, статическая компиляция)
- **Протокол**: MCP 2025-11-25 (через [mcp-go SDK](https://github.com/mark3labs/mcp-go))
- **Транспорт**: Streamable HTTP
- **Логирование**: slog (структурное JSON)

## Требования

- Docker и Docker Compose v2
- Ключ доступа InSales (Расширения → Разработчикам → Ключи доступа)
- AI-клиент с поддержкой MCP: Claude Desktop, Cursor, VS Code/Copilot, MiMoCode, OpenCode, Continue, LibreChat, Open WebUI, Cline, Cherry Studio, ChatGPT Desktop, Windsurf, Zed, Aider и другие совместимые MCP-клиенты

## Установка

### 1. Клонирование

```bash
git clone <url_repos>
cd mcp-is
```

### 2. Настройка

Скопируйте пример конфигурации и заполните данными вашего магазина:

```bash
cp .env.example .env
```

Содержимое `.env`:

```env
# База данных
DB_PASSWORD=password

# Сервер
SERVER_PORT=8080
LOG_LEVEL=info

# InSales — стороннее приложение (install/uninstall/autologin). Подробнее: docs/setup/INSALES-APP.md
INSALES_APP_KEY=
INSALES_SECRET_KEY=

# Режим доступа (true = только чтение)
MCP_READ_ONLY=true

# Публичная регистрация (false = пользователей создаёт только админ)
REGISTRATION_ENABLED=false

# Auth: JWT secret (обязательно сменить в production!)
JWT_SECRET=change-me-to-random-string

# Auth: AES-GCM key для шифрования паролей InSales (32 байта!)
ENCRYPTION_KEY=12345678901234567890123456789012
```

### 3. Запуск

```bash
docker compose up -d
```

### 4. Проверка

```bash
curl http://localhost:8080/health
# Ответ: ok
```

## MCP Инструменты

Реализовано через [mcp-go SDK](https://github.com/mark3labs/mcp-go) (MCP 2025-11-25). Каждый инструмент имеет аннотации `readOnlyHint`/`destructiveHint` для AI.

### Товары

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_products` | Список товаров | `page`, `per_page`, `category_id`, `collection_id`, `q` |
| `get_product` | Детали товара | `product_id` |
| `create_product` | Создать товар | `title`, `price`, `description`, `vendor`, `sku`, `quantity` |
| `update_product` | Обновить товар | `product_id`, `title`, `price`, `quantity` |
| `delete_product` | Удалить товар | `product_id` |

### Заказы

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_orders` | Список заказов | `page`, `per_page`, `status`, `fulfillment_status`, `paid` |
| `get_order` | Детали заказа | `order_id` |
| `create_order` | Создать заказ | `client_id`, `comment` |
| `update_order` | Обновить заказ | `order_id`, `status`, `fulfillment_status`, `comment` |

### Клиенты

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_clients` | Список клиентов | `page`, `per_page`, `q` |
| `get_client` | Данные клиента | `client_id` |
| `create_client` | Создать клиента | `email`, `first_name`, `last_name`, `phone` |
| `update_client` | Обновить клиента | `client_id`, `email`, `phone` |
| `delete_client` | Удалить клиента | `client_id` |

### Категории

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_categories` | Список категорий | `page`, `per_page`, `parent_id` |
| `get_category` | Данные категории | `category_id` |
| `create_category` | Создать категорию | `title`, `parent_id` |
| `update_category` | Обновить категорию | `category_id`, `title` |
| `delete_category` | Удалить категорию | `category_id` |

### Коллекции

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_collections` | Список коллекций | `page`, `per_page` |
| `get_collection` | Данные коллекции | `collection_id` |
| `create_collection` | Создать коллекцию | `title` |
| `update_collection` | Обновить коллекцию | `collection_id`, `title` |
| `delete_collection` | Удалить коллекцию | `collection_id` |

### Изображения

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_images` | Изображения товара | `product_id` |
| `add_image` | Добавить изображение | `product_id`, `url`, `filename`, `position` |
| `delete_image` | Удалить изображение | `product_id`, `image_id` |

### Метаполя

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_meta_fields` | Список метаполей | `resource_type`, `resource_id` |
| `create_meta_field` | Создать метаполе | `resource_type`, `resource_id`, `key`, `value` |
| `update_meta_field` | Обновить метаполе | `meta_field_id`, `value` |
| `delete_meta_field` | Удалить метаполе | `meta_field_id` |

### Вебхуки

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_webhooks` | Список вебхуков | — |
| `create_webhook` | Подписаться на событие | `event`, `url`, `format` |
| `update_webhook` | Обновить вебхук | `webhook_id`, `url` |
| `delete_webhook` | Удалить вебхук | `webhook_id` |

### Системные

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_currencies` | Валюты магазина | — |
| `list_delivery_variants` | Способы доставки | — |
| `list_payment_gateways` | Способы оплаты | — |

### Страницы

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_pages` | Список страниц | `page`, `per_page` |
| `get_page` | Страница по ID | `page_id` |
| `create_page` | Создать страницу | `title`, `content`, `handle` |
| `update_page` | Обновить страницу | `page_id`, `title`, `content`, `handle` |
| `delete_page` | Удалить страницу | `page_id` |

### Аксессуары

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_accessories` | Список аксессуаров | `page`, `per_page` |
| `get_accessory` | Аксессуар по ID | `accessory_id` |
| `create_accessory` | Создать аксессуар | `title` |
| `update_accessory` | Обновить аксессуар | `accessory_id`, `title`, `position` |
| `delete_accessory` | Удалить аксессуар | `accessory_id` |

### Варианты товаров

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_variants` | Список вариантов товара | `product_id` |
| `get_variant` | Вариант по ID | `product_id`, `variant_id` |
| `update_variant` | Обновить вариант | `product_id`, `variant_id`, `price`, `sku`, `quantity`, `weight` |
| `delete_variant` | Удалить вариант | `product_id`, `variant_id` |

### Значения свойств

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_property_values` | Значения аксессуара | `accessory_id` |
| `create_property_value` | Создать значение | `accessory_id`, `title` |
| `update_property_value` | Обновить значение | `property_value_id`, `title` |
| `delete_property_value` | Удалить значение | `property_value_id` |

### Меню

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_menus` | Список меню | — |
| `get_menu` | Меню по ID | `menu_id` |
| `create_menu` | Создать меню | `title`, `handle` |
| `update_menu` | Обновить меню | `menu_id`, `title`, `handle` |
| `delete_menu` | Удалить меню | `menu_id` |

### Редиректы

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_redirects` | Список редиректов | — |
| `create_redirect` | Создать редирект | `url`, `target` |
| `delete_redirect` | Удалить редирект | `redirect_id` |

### Блоги

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_blogs` | Список блогов | — |
| `get_blog` | Блог по ID | `blog_id` |
| `create_blog` | Создать блог | `title`, `handle` |
| `delete_blog` | Удалить блог | `blog_id` |

### Записи блога

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_blog_posts` | Записи блога | `blog_id` |
| `get_blog_post` | Запись по ID | `blog_id`, `post_id` |
| `create_blog_post` | Создать запись | `blog_id`, `title`, `body_html` |
| `update_blog_post` | Обновить запись | `blog_id`, `post_id`, `title`, `body_html` |
| `delete_blog_post` | Удалить запись | `blog_id`, `post_id` |

### Точки выдачи (Outlets)

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_outlets` | Список точек выдачи | — |
| `get_outlet` | Точка выдачи по ID | `outlet_id` |
| `create_outlet` | Создать точку выдачи | `title`, `address`, `price`, `latitude`, `longitude`, `description` |
| `update_outlet` | Обновить точку выдачи | `outlet_id`, `title`, `address`, `price`, `latitude`, `longitude`, `description` |
| `delete_outlet` | Удалить точку выдачи | `outlet_id` |

### Маркетплейсы

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `list_marketplaces` | Список маркетплейсов | — |
| `get_marketplace` | Маркетплейс по ID | `marketplace_id` |
| `create_marketplace` | Создать маркетплейс | `name`, `type` (обязательные), + параметры типа |
| `update_marketplace` | Обновить маркетплейс | `marketplace_id` + параметры |
| `delete_marketplace` | Удалить маркетплейс | `marketplace_id` |

**Типы маркетплейсов:**
- `Marketplace::GoogleMerchant` — Google Merchant Center
- `Marketplace::LiquidYandexMarket` — Яндекс Маркет (Liquid шаблон)
- `Marketplace::ModelYandexMarket` — Яндекс Маркет (модельный)

### Токены мессенджеров

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `create_provider_token` | Создать токен провайдера | `messenger_type` (whats_app, telegram, viber, sms, email, vk, instagram) |

### Настройки аккаунта

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `update_account` | Обновить настройки аккаунта | `client_cookies_whitelist` |

### Счётчики

| Инструмент | Описание | Параметры |
|------------|----------|-----------|
| `get_orders_count` | Количество заказов | `fulfillment_status` |
| `get_products_count` | Количество товаров | — |

### Экспорт данных (CSV)

| Инструмент | Описание | Формат |
|------------|----------|--------|
| `export_products` | Экспорт товаров | CSV: id, title, price, sku, vendor, quantity, created_at |
| `export_orders` | Экспорт заказов | CSV: id, number, fulfillment_status, total_price, client_name, created_at |
| `export_clients` | Экспорт клиентов | CSV: id, name, surname, email, phone, created_at |

## Аутентификация

### Регистрация

Публичная регистрация управляется переменной `REGISTRATION_ENABLED` в `.env`:
- `false` (рекомендуется) — `/api/register` возвращает 403, пользователей создаёт только админ из админки
- `true` — регистрация открыта

Исключение: если в базе нет ни одного пользователя, регистрация доступна даже при `false` — первый пользователь становится `super_admin` (bootstrap).

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "myuser",
    "email": "user@example.com",
    "phone": "+79991234567",
    "first_name": "Иван",
    "last_name": "Петров",
    "password": "securepassword"
  }'
```

### Вход (по nickname, email или phone)

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "myuser", "password": "securepassword"}'
```

Ответ: `{ "token": "eyJ...", "user_id": "...", "nickname": "myuser" }`

### Управление магазинами

```bash
# Добавить магазин (требуется JWT, name обязателен)
curl -X POST http://localhost:8080/api/shops \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "name": "Моя Аптека",
    "shop_url": "myshop.myinsales.ru",
    "api_key_id": "ваш_api_key_id",
    "api_password": "ваш_пароль",
    "is_default": true
  }'

# Список магазинов
curl http://localhost:8080/api/shops \
  -H "Authorization: Bearer <JWT_TOKEN>"

# Проверить доступные команды магазина
curl -X POST http://localhost:8080/api/shops/{id}/check-commands \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### Управление AI-ключами

```bash
# Сгенерировать ключ (показывается ОДИН раз!)
curl -X POST http://localhost:8080/api/ai-keys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{"name": "Cursor MCP"}'

# Список ключей (без самих ключей)
curl http://localhost:8080/api/ai-keys \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### Имя ключа в конфиге

Имя AI-ключа автоматически используется как JSON-ключ в MCP-конфиге. Slugify: кириллица → транслитерация, пробелы → дефисы, только латиница и цифры.

| Имя ключа | JSON-ключ в конфиге |
|-----------|---------------------|
| `Cursor MCP` | `"cursor-mcp"` |
| `Мой Ключ` | `"moy-klyuch"` |
| `test` | `"test"` |
| *(пустое)* | `"mcp-is"` |

Пример конфига для Claude Desktop с ключом "Cursor MCP":
```json
{
  "mcpServers": {
    "cursor-mcp": {
      "url": "https://mcp-is.ru:443/mcp",
      "headers": {
        "Authorization": "Bearer 59c213ef..."
      }
    }
  }
}
```

### MCP-запросы (требуется API-ключ)

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <AI_KEY>" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 1,
    "params": {
      "name": "list_products",
      "arguments": {"page": 1, "per_page": 5}
    }
  }'
```

## Админ-панель

### Роли
| Роль | Права |
|------|-------|
| **super_admin** | Первый пользователь. Управляет всеми. Нельзя удалить/заблокировать/сменить роль кроме него самого |
| **admin** | Управляет user-ами. Не может удалить super_admin |
| **user** | Только свои магазины и ключи |

### Страницы
| Страница | URL | Доступ |
|----------|-----|--------|
| Пользователи | `/admin/users` | admin, super_admin |
| Просмотр | `/admin/users/:id` | admin, super_admin |
| Статистика | `/admin/stats` | admin, super_admin |
| Настройки сервера | `/admin/server` | admin, super_admin |

### Настройки сервера

Суперадмин управляет доменом и HTTPS из web-интерфейса:

| Функция | Доступ | Описание |
|---------|--------|----------|
| Доменное имя | super_admin | Настраивает домен для nginx |
| HTTPS (Let's Encrypt) | super_admin | Включает/выключает SSL |
| Email для LE | super_admin | Email для уведомлений Let's Encrypt |
| Автообновление | super_admin | Галочка + интервал (14–60 дней) |
| Выпуск сертификата | admin, super_admin | Запрос на выпуск LE сертификата |
| Продление сертификата | admin, super_admin | Ручное продление |
| Статус сертификата | admin, super_admin | Тип, дата истечения, ошибки |

**Автообновление:** Go scheduler проверяет каждый час. При ошибке — повтор через 7 часов, максимум 10 попыток (~3 дня). Админы видят ошибки и могут нажать "Попробовать снова".

### Защита super_admin
- Нельзя удалить
- Нельзя заблокировать
- Нельзя изменить роль
- Только сам super_admin может менять свой пароль

## Смена и восстановление пароля

### Смена пароля
На странице `/settings` → блок "Смена пароля". Требуется текущий пароль.

### Восстановление пароля
1. На странице `/login` нажмите "Забыли пароль?"
2. Введите email
3. Получите ссылку для сброса (в логах сервера или по SMTP)
4. Задайте новый пароль

## SMTP (отправка писем)

Для отправки писем настройте SMTP в `.env`:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
SMTP_FROM=MCP IS <noreply@your-domain.com>
```

Если SMTP не настроен — токены сброса логируются в консоль.

## REST API Эндпоинты

| Метод | Путь | Описание | Авторизация |
|-------|------|----------|-------------|
| POST | `/api/register` | Регистрация (если `REGISTRATION_ENABLED=true` или база пуста) | Нет |
| POST | `/api/login` | Вход | Нет |
| POST | `/api/forgot-password` | Запрос сброса пароля | Нет |
| POST | `/api/reset-password` | Сброс по токену | Нет |
| GET | `/install` | Установка приложения InSales | Нет |
| GET | `/uninstall` | Удаление приложения InSales | Нет |
| GET | `/app` | Страница приложения InSales | Нет |
| GET | `/autologin` | Автологин из InSales | Нет |
| POST | `/api/exchange-code` | Обмен one-time code на JWT | Нет |
| POST | `/api/change-password` | Смена пароля | JWT |
| POST | `/api/set-email` | Установка email | JWT |
| POST | `/api/shops` | Добавить магазин | JWT |
| POST | `/api/shops/link-pending` | Привязать магазин из InSales | JWT |
| POST | `/api/shops/test-connection` | Проверить ключ InSales | JWT |
| POST | `/api/shops/:id/check-commands` | Проверить доступные команды | JWT |
| GET | `/api/shops` | Список магазинов | JWT |
| PUT | `/api/shops/:id` | Обновить магазин | JWT |
| DELETE | `/api/shops/:id` | Удалить магазин | JWT |
| POST | `/api/ai-keys` | Сгенерировать ключ | JWT |
| GET | `/api/ai-keys` | Список ключей | JWT |
| DELETE | `/api/ai-keys/:id` | Отозвать ключ | JWT |
| POST | `/api/admin/users` | Создать пользователя | JWT + admin |
| GET | `/api/admin/users` | Список пользователей | JWT + admin |
| PUT | `/api/admin/users/:id/role` | Изменить роль | JWT + super_admin |
| PUT | `/api/admin/users/:id/status` | Заблокировать | JWT + admin |
| DELETE | `/api/admin/users/:id` | Удалить пользователя | JWT + admin |
| GET | `/api/admin/stats` | Статистика | JWT + admin |
| GET | `/api/admin/pdn` | Настройки ПДн | JWT + admin |
| PUT | `/api/admin/pdn` | Обновить настройки ПДн | JWT + super_admin |
| GET | `/api/admin/pdn/log` | Журнал доступа к ПДн | JWT + admin |
| GET | `/api/admin/pdn/stats` | Статистика ПДн за сегодня | JWT + admin |
| POST | `/mcp` | MCP-запросы | API Key |
| GET | `/mcp` | SSE-подключение | API Key |
| GET | `/health` | Проверка здоровья | Нет |

## Безопасность

- **Пароли пользователей** — хранятся в bcrypt
- **Пароли InSales** — зашифрованы AES-GCM (ключ из `ENCRYPTION_KEY`)
- **AI-ключи** — хранятся как SHA-256 хэши (сами ключи не хранятся)
- **JWT-токены** —用于 управления магазинами (срок 24 часа)
- **MCP-эндпоинты** — защищены API-ключами

## MCP контекст

При подключении AI-агента MCP сервер возвращает контекст магазина:

```json
{
  "serverInfo": {
    "name": "mcp-is",
    "shop": "Моя Аптека",
    "key": "Cursor MCP",
    "read_only": true
  }
}
```

При `tools/list` с read-only возвращаются только 30 read-команд (без create/update/delete).

## Доступные команды

Нажмите "Команды" в списке магазинов чтобы увидеть доступные инструменты с описаниями:

| Режим | Команд | Описание |
|-------|--------|----------|
| read-only | ~42 | list_*, get_*, system, pipeline, export команды |
| read/write | 96 | Все команды включая create/update/delete + пайплайны |

## Пайплайн-инструменты

Композитные инструменты для частых сценариев — получите полную информацию за 1 вызов вместо 4-5:

| Инструмент | Что делает |
|------------|-----------|
| `search_product_details` | Поиск товара + кастомные поля + метаполя + изображения |
| `order_full_details` | Заказ + клиент + адрес + позиции + метаполя |
| `client_full_profile` | Клиент + последние заказы + бонусы + метаполя |
| `product_catalog_search` | Поиск по каталогу с фильтрами |
| `collection_products` | Коллекция + все её товары |

## Примеры запросов

### Получить список товаров

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 1,
    "params": {
      "name": "list_products",
      "arguments": {"page": 1, "per_page": 5}
    }
  }'
```

### Получить заказ по ID

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 2,
    "params": {
      "name": "get_order",
      "arguments": {"order_id": "1523259105"}
    }
  }'
```

### Список доступных инструментов

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "tools/list", "id": 3}'
```

## Режимы доступа

### Только чтение

**Глобальный (наивысший приоритет):**
```env
MCP_READ_ONLY=true
```
Блокирует ВСЕ write-операции для ВСЕХ пользователей и магазинов.

**Per-shop (если глобальный НЕ установлен):**
При добавлении магазина через Web UI или API можно установить флаг "Только чтение". Этот флаг действует ТОЛЬКО для конкретного магазина — другие магазины и пользователи продолжают работать в обычном режиме.

### Полный доступ

```env
MCP_READ_ONLY=false
```

После изменения перезапустите контейнер:

```bash
docker compose up -d --force-recreate app
```

## Развёртывание на VDS

### Требования

- Ubuntu 24.04 LTS (или Debian/Fedora)
- Docker и Docker Compose v2
- 1 CPU / 256MB RAM (минимум)

### Быстрый старт

```bash
# Клонировать проект
git clone https://github.com/chubuk/mcp-is.git
cd mcp-is

# Настроить .env
cp .env.example .env
nano .env

# Запустить
docker compose up -d

# Проверить
curl -k https://YOUR_DOMAIN/health
```

### Автоматический деплой (deploy.sh)

Скрипт `scripts/deploy.sh` автоматизирует все операции:

```bash
# Проверка системы (pre-flight)
./scripts/deploy.sh check

# Полная установка на чистый сервер
./scripts/deploy.sh install

# Обновление существующей установки
./scripts/deploy.sh update

# Ручной бэкап
./scripts/deploy.sh backup

# Проверка статуса
./scripts/deploy.sh status
```

### Автоматические бэкапы

Бэкапы настраиваются через админку (`/admin/backup`) или `.env`:

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `BACKUP_ENABLED` | Включить бэкапы | true |
| `BACKUP_METHOD` | pg_basebackup / pg_dump | pg_basebackup |
| `BACKUP_INTERVAL_HOURS` | Интервал (часов) | 24 |
| `BACKUP_RETENTION_DAYS` | Хранение файлов, журнала и инцидентов (дней) | 14 |

### CI/CD (GitHub Actions)

Автоматический деплой при пуше в main:

1. Бэкап перед деплоем
2. Сборка Docker образов
3. Деплой на VDS
4. Health check
5. Rollback при ошибке

Требуются GitHub Secrets: `SERVER_HOST`, `SERVER_USER`, `SSH_KEY`.

### Безопасность

| Мера | Описание |
|------|----------|
| **Пароли** | Генерируются автоматически deploy.sh, сохраняются в `/root/.mcp-is-credentials` |
| **Docker-сеть** | PostgreSQL доступен только из контейнера app (internal сеть) |
| **Порт 5432** | НЕ экспортирован на хост |
| **Rate limiting** | login: 5/min, API: 10/s, reset-password: 3/min, MCP: 30/min |
| **Блокировка файлов** | .env, .git, .sql, .log — 403 Forbidden |
| **Security headers** | HSTS, CSP, X-Frame-Options, Permissions-Policy |
| **Аутентификация** | JWT для API, API-ключи для MCP |
| **RBAC** | user, admin, super_admin роли |

### Systemd (опционально)

```ini
[Unit]
Description=MCP IS Server
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
WorkingDirectory=/opt/mcp-is
ExecStart=/usr/bin/docker compose up
ExecStop=/usr/bin/docker compose down
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## Команды разработки

```bash
make build          # Собрать бинарник
make run            # Запустить локально
make test           # Запустить тесты
make docker-up      # Запустить в Docker
make docker-down    # Остановить Docker
make logs           # Смотреть логи
make clean          # Очистить собранные файлы
make deb            # Собрать .deb пакет
```

## Веб-интерфейс

В проект встроен SPA-интерфейс для управления магазинами и ключами.

### Архитектура

```
Браузер → nginx (HTTPS) → статика (HTML/CSS/JS)
                        → /api/* → Go API (REST)
                        → /mcp   → Go API (MCP)
```

### Страницы

| Страница | URL | Описание | Доступ |
|----------|-----|----------|--------|
| Вход | `/login` | Вход по email/phone/nickname | Все |
| Регистрация | `/register` | Создание аккаунта (только bootstrap или `REGISTRATION_ENABLED=true`) | Все |
| Панель | `/dashboard` | Обзор магазинов и ключей | Авторизованные |
| Магазины | `/shops` | Управление магазинами InSales | Авторизованные |
| AI-ключи | `/ai-keys` | Генерация и отзыв ключей | Авторизованные |
| Настройки | `/settings` | Профиль пользователя | Авторизованные |
| Пользователи | `/admin/users` | Управление пользователями (создание, удаление) | admin, super_admin |
| Просмотр | `/admin/users/:id` | Детали пользователя | admin, super_admin |
| Статистика | `/admin/stats` | Статистика системы | admin, super_admin |

### Технологии

- HTML + CSS + vanilla JS (без фреймворков)
- SPA-роутинг через History API
- JWT токены в localStorage
- HTTPS через Let's Encrypt (или самоподписанный для разработки)
- Адаптивный дизайн (мобильный/планшет/десктоп)

### Адаптивный дизайн

Интерфейс работает на всех устройствах:

| Устройство | Ширина | Особенности |
|------------|--------|-------------|
| Десктоп | > 1024px | Полный хедер с навигацией, 2-3 колонки в сетках |
| Планшет | 640–1024px | Компактный хедер, уменьшенные отступы |
| Мобильный | < 640px | Бургер-меню, 1 колонка, скроллируемые таблицы |

**Реализация:**
- CSS media queries для 3 breakpoints
- Бургер-меню для мобильной навигации
- Таблицы с горизонтальным скроллом (`.table-wrap`)
- Адаптивные сетки статистики (`.stats-grid`)
- Увеличенные touch-target для кнопок на мобильных

### Тесты JS

Откройте `web/static/test.html` в браузере для проверки JS-модулей.

## Тестирование

```bash
# Все тесты (включая integration с PostgreSQL)
go test ./... -tags=integration -v

# Только unit-тесты (без БД)
go test ./internal/auth/... ./internal/handlers/... -v

# Только integration-тесты
go test ./internal/storage/... -tags=integration -v
```

## Логирование

Сервер использует структурное JSON-логирование через `slog`. Пример вывода:

```json
{"time":"2026-06-29T11:44:24Z","level":"INFO","msg":"starting mcp-is server","port":"8080"}
{"time":"2026-06-29T11:44:24Z","level":"INFO","msg":"insales account configured","shop":"myshop.myinsales.ru"}
{"time":"2026-06-29T11:44:24Z","level":"INFO","msg":"mcp server mode","read_only":true}
```

Уровень логирования настраивается через `LOG_LEVEL` (debug, info, warn, error).

## Лицензия

[Business Source License 1.1 (BSL 1.1)](LICENSE) — коммерческое использование запрещено без соглашения с правообладателем. Через 4 года (2030) автоматически переходит в [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

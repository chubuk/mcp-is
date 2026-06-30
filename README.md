# MCP InSales

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
  webhooks/                 — обработка вебхуков InSales
web/
  nginx/                    — Dockerfile, nginx.conf, entrypoint.sh
  static/                   — HTML, CSS, JS (SPA фронтенд)
```

## Технологии

- **Язык**: Go 1.24
- **БД**: PostgreSQL 15
- **Контейнеризация**: Docker (multi-stage сборка, статическая компиляция)
- **Протокол**: MCP (JSON-RPC 2.0 + SSE)
- **Логирование**: slog (структурное JSON)

## Требования

- Docker и Docker Compose v2
- Ключ доступа InSales (Расширения → Разработчикам → Ключи доступа)

## Установка

### 1. Клонирование

```bash
git clone <url_repos>
cd mcp-insales
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

# InSales (legacy, можно не заполнять если используете API через UI)
INSALES_SHOP=your-shop.myinsales.ru
INSALES_API_KEY=your_api_key_here
INSALES_PASSWORD=your_password_here

# Режим доступа (true = только чтение)
MCP_READ_ONLY=true

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

## Аутентификация

### Регистрация

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
# Добавить магазин (требуется JWT)
curl -X POST http://localhost:8080/api/shops \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "shop_url": "myshop.myinsales.ru",
    "api_key_id": "ваш_api_key_id",
    "api_password": "ваш_пароль",
    "is_default": true
  }'

# Список магазинов
curl http://localhost:8080/api/shops \
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
SMTP_FROM=MCP InSales <noreply@your-domain.com>
```

Если SMTP не настроен — токены сброса логируются в консоль.

## REST API Эндпоинты

| Метод | Путь | Описание | Авторизация |
|-------|------|----------|-------------|
| POST | `/api/register` | Регистрация | Нет |
| POST | `/api/login` | Вход | Нет |
| POST | `/api/forgot-password` | Запрос сброса пароля | Нет |
| POST | `/api/reset-password` | Сброс по токену | Нет |
| POST | `/api/change-password` | Смена пароля | JWT |
| POST | `/api/set-email` | Установка email | JWT |
| POST | `/api/shops` | Добавить магазин | JWT |
| GET | `/api/shops` | Список магазинов | JWT |
| PUT | `/api/shops/:id` | Обновить магазин | JWT |
| DELETE | `/api/shops/:id` | Удалить магазин | JWT |
| POST | `/api/ai-keys` | Сгенерировать ключ | JWT |
| GET | `/api/ai-keys` | Список ключей | JWT |
| DELETE | `/api/ai-keys/:id` | Отозвать ключ | JWT |
| GET | `/api/admin/users` | Список пользователей | JWT + admin |
| PUT | `/api/admin/users/:id/role` | Изменить роль | JWT + super_admin |
| PUT | `/api/admin/users/:id/status` | Заблокировать | JWT + admin |
| DELETE | `/api/admin/users/:id` | Удалить пользователя | JWT + admin |
| GET | `/api/admin/stats` | Статистика | JWT + admin |
| POST | `/mcp` | MCP-запросы | API Key |
| GET | `/mcp` | SSE-подключение | API Key |
| GET | `/health` | Проверка здоровья | Нет |

## Безопасность

- **Пароли пользователей** — хранятся в bcrypt
- **Пароли InSales** — зашифрованы AES-GCM (ключ из `ENCRYPTION_KEY`)
- **AI-ключи** — хранятся как SHA-256 хэши (сами ключи не хранятся)
- **JWT-токены** —用于 управления магазинами (срок 24 часа)
- **MCP-эндпоинты** — защищены API-ключами

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

### Только чтение (по умолчанию)

```env
MCP_READ_ONLY=true
```

Запрещены все операции записи: `create_*`, `update_*`, `delete_*`.

### Полный доступ

```env
MCP_READ_ONLY=false
```

После изменения перезапустите контейнер:

```bash
docker compose up -d --force-recreate app
```

## Развёртывание на VPS

### Требования

- Ubuntu 24.04 LTS
- Docker и Docker Compose v2
- 1 CPU / 512MB RAM (минимум)

### Установка на VPS

```bash
# Клонировать проект
git clone <url_repos>
cd mcp-insales

# Настроить .env
cp .env.example .env
nano .env

# Запустить
docker compose up -d
```

### Systemd (опционально)

```ini
[Unit]
Description=MCP InSales Server
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
WorkingDirectory=/opt/mcp-insales
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
| Регистрация | `/register` | Создание аккаунта | Все |
| Панель | `/dashboard` | Обзор магазинов и ключей | Авторизованные |
| Магазины | `/shops` | Управление магазинами InSales | Авторизованные |
| AI-ключи | `/ai-keys` | Генерация и отзыв ключей | Авторизованные |
| Настройки | `/settings` | Профиль пользователя | Авторизованные |
| Пользователи | `/admin/users` | Управление пользователями | admin, super_admin |
| Просмотр | `/admin/users/:id` | Детали пользователя | admin, super_admin |
| Статистика | `/admin/stats` | Статистика системы | admin, super_admin |

### Технологии

- HTML + CSS + vanilla JS (без фреймворков)
- SPA-роутинг через History API
- JWT токены в localStorage
- HTTPS через Let's Encrypt (или самоподписанный для разработки)

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
{"time":"2026-06-29T11:44:24Z","level":"INFO","msg":"starting mcp-insales server","port":"8080"}
{"time":"2026-06-29T11:44:24Z","level":"INFO","msg":"insales account configured","shop":"myshop.myinsales.ru"}
{"time":"2026-06-29T11:44:24Z","level":"INFO","msg":"mcp server mode","read_only":true}
```

Уровень логирования настраивается через `LOG_LEVEL` (debug, info, warn, error).

## Лицензия

[Business Source License 1.1 (BSL 1.1)](LICENSE) — коммерческое использование запрещено без соглашения с правообладателем. Через 4 года (2030) автоматически переходит в [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

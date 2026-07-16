# MCP InSales — Архитектура

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
│  │ Rate Limit  │  │ 74 tools    │  │ admin/auth/server   │ │
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
│  • 74+ REST endpoints                                       │
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
- **74 инструмента** для работы с InSales
- **Tool annotations** — readOnlyHint, destructiveHint
- **Tool filtering** — read-only режим скрывает write-команды
- **Server instructions** — контекст для AI

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
    ↓ приоритетнаибольший
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
sudo dpkg -i mcp-insales_0.1.0_amd64.deb
sudo nano /etc/mcp-insales/.env
sudo systemctl start mcp-insales
```

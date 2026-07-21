# Changelog

## [0.2.2] - 2026-07-21

## [Unreleased]

### Добавлено
- **MCP URL** — при создании AI-ключа показывается MCP URL с портом
- **Copy buttons** — копирование ключа, URL и конфигурации для Claude/VS Code
- **AI key PDN toggle** — клик по бейджу ПДн включает/выключает защиту
- **Shop edit** — редактирование магазинов (название, домен, ключ, is_read_only, pdn_enabled)
- **UI/UX скил** — универсальные паттерны (copy, toast, loading, empty states)
- **Responsive Web Design скил** — breakpoints, mobile-first, Flexbox, Grid
- **DB schema check скил** — проверка схемы БД перед коммитом
- **Shop paste feature** — быстрое заполнение ключа из InSales
- **Retry with backoff** — автоматические повторы с exponential backoff + jitter
- **Prometheus метрики** — http_requests_total, http_request_duration_seconds, mcp_tool_calls_total, active_connections
- **pprof** — профилирование CPU/memory/goroutine на порту 6060
- **OpenAPI spec** — полная спецификация REST API (28 эндпоинтов)
- **GitHub Actions CI/CD** — lint → test → build → docker → deploy
- **Backup script** — автоматический бэкап PostgreSQL с retention 7 дней

### Улучшено
- **Docker** — non-root user (appuser), resource limits (2 CPU/512MB), cap_drop: ALL, logging rotation
- **Nginx** — server_tokens off, request size limits, JSON logging, request ID, stricter login rate limit
- **PostgreSQL** — connection pool tuning (MaxConns=20, MinConns=2), pdn_access_log indexes
- **Auth** — transactional CreateUser+UpdateRole, email masking in logs
- **Config** — JWT expiry configurable (JWT_EXPIRY_HOURS)

## [0.2.0] — 2026-07-16

### Добавлено
- **Миграция на mcp-go SDK** — полная замена ручной реализации MCP
- **MCP 2025-11-25** — обновление спецификации (было 2024-11-05)
- **Streamable HTTP** — новый транспорт вместо ручного SSE
- **Tool annotations** — `readOnlyHint`, `destructiveHint` для безопасности AI
- **Server instructions** — контекст для AI о возможностях сервера
- **Recovery** — автоматическое восстановление от паник
- **24 новых MCP-инструмента**:
  - Variants: list, get, update, delete
  - PropertyValues: list, create, update, delete
  - Menus: list, get, create, update, delete
  - Redirects: list, create, delete
  - Blogs: list, get, create, delete
  - BlogPosts: list, get, create, update, delete
- **Имена магазинов/ключей** — обязательные поля для AI контекста
- **Check-commands** — проверка доступных команд для магазина
- **Server Settings UI** — домен, HTTPS, Let's Encrypt из админки
- **Rate limiting** — 10 req/min на login/register/forgot-password
- **Адаптивный дизайн** — mobile (<640px), tablet (640-1024px), desktop
- **CORS middleware** — кросс-доменные запросы
- **Документация**: Tutorial, Architecture, FAQ, Troubleshooting, Changelog

### Изменено
- Миграция с кастомной реализации MCP на mcp-go SDK
- Обновление MCP-протокола с 2024-11-05 на 2025-11-25
- Замена ручного SSE на Streamable HTTP
- Автор лицензии: Vladislav Chubuk <chubuk@mail.ru>

### Исправлено
- Test connection endpoint (shops.json → orders.json)
- Healthcheck для app контейнера

## [0.1.0] — 2026-06-29

### Добавлено
- Начальный релиз
- 49 MCP-инструментов для InSales API
- Web UI с админ-панелью (vanilla JS SPA)
- JWT аутентификация
- API-ключи для MCP-клиентов
- Docker деплой (app, db, nginx)
- .deb пакет для установки на VPS
- Read-only режим (глобальный + per-shop)
- PostgreSQL 15 для хранения данных
- SMTP для восстановления пароля
- Pагинация и поиск

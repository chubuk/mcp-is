# Changelog

## [0.2.4] - 2026-07-23

### Добавлено
- **Impersonation** — admin/super_admin может войти в аккаунт пользователя и вернуться обратно
  - super_admin → любой user или admin
  - admin → только user
  - Запрет на impersonate super_admin, себя, другого admin
  - Жёлтый баннер «Вы зашли как [имя] | Выйти»
  - Журнал действий через slog
- **MCP client config generator** — на странице AI-ключа выпадающий список из 10 MCP-систем (Claude Desktop, Cursor, VS Code, MiMoCode, OpenCode, LibreChat, Open WebUI, Cline, Cherry Studio, Continue)
  - Генерация готового конфига с подстановкой URL и ключа
  - Кнопки «Копировать» и «Скачать файл» (с правильным именем/расширением)
  - Подсказка «куда класть файл»
  - Запоминание выбранной системы в localStorage
  - Список систем хранится в БД — расширяемый без обновления кода

### Исправлено
- `.gitignore` — паттерн `server` заменён на `/server` (binaries only)
- `stop-impersonate` — убрана лишняя проверка adminMW (impersonated user имеет роль `user`)
- `server_settings` — `InitDefaults()` перезаписывала domain/https/cert_issuer при каждом старте. Заменена на `SetDefaultsIfMissing()` (вставка только отсутствующих ключей)

### Тесты
- 19 новых unit-тестов (3 JWT + 6 middleware + 10 handler)

## [0.2.3] - 2026-07-23

### Безопасность
- Убрана экспозиция PostgreSQL порта 5432 на хост
- Добавлена Docker-сеть с изоляцией БД (internal/public)
- Заблокирован /metrics эндпоинт (Prometheus)
- Rate limiting: /api/reset-password (3 req/min), /mcp (30 req/min)
- Добавлен Permissions-Policy header
- Автоматическая генерация паролей в deploy.sh
- Пароли сохраняются в /root/.mcp-is-credentials (chmod 600)
- CORS ограниченALLOWED_ORIGINS

### Деплой
- deploy.sh: install, update, check, rollback, status, backup
- backup.sh: pg_basebackup, мониторинг нагрузки, retention, удалённые хранилища
- CI/CD: бэкап перед деплоем, rollback при ошибке, проверка версии

### API
- /api/admin/backup: настройки, журнал, инциденты, статус, ручной бэкап

### Frontend
- Страница /admin/backup: статистика, настройки, журнал, инциденты

### Документация
- INSTALL.md: раздел автодеплоя и безопасности
- README.md: секция автодеплоя и безопасности
- TROUBLESHOOTING.md: раздел доступа к БД

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
- **Certbot unit tests** — 13 тестов для `internal/certbot` (GetCertInfo, CheckRenewal, ShouldRetry, incrementRetry)

### Исправлено
- **Certbot в nginx** — добавлен `certbot` в nginx Dockerfile, выпуск/продление сертификатов Let's Encrypt теперь работает
- **SettingsStore interface** — `internal/certbot` использует интерфейс вместо конкретного типа, улучшает тестируемость

### Улучшено
- **Docker** — non-root user (appuser), resource limits (2 CPU/512MB), cap_drop: ALL, logging rotation
- **Nginx** — server_tokens off, request size limits, JSON logging, request ID, stricter login rate limit, certbot для Let's Encrypt
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

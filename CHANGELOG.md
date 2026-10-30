# Changelog

## [Unreleased] — настраиваемая стоимость подписки

### Добавлено
- **«Стоимость подписки» в настройках сервера** (Админ → Сервер, super_admin): `billing_monthly` в `server_settings`, hot-reload без рестарта. `0` = бесплатно (гейт прозрачен), `>0` = сумма RecurringApplicationCharge. Env `BILLING_MONTHLY` = значение по умолчанию/фолбэк. Настройка не пересоздаёт уже созданный charge.

## [Unreleased] — favicon

### Добавлено
- `favicon.ico` (16/32/48) + `favicon.svg` (метка «M» в цветах интерфейса), `<link rel="icon">` в `index.html`. nginx/deb не менялись — `.ico` уже отдаётся статикой.

## [Unreleased] — безопасность autologin: one-time код без URL, хвосты OTP

### Безопасность
- **Утечка MCP-ключа в URL устранена**: one-time код autologin больше не содержит сырой ключ (был внутри HMAC payload) и не передаётся через URL (`/onboarding?code=…` утекал в историю/логи/Referer). Теперь серверный одноразовый стор (TTL 2 мин, атомарное сжигание — replay невозможен), передача через HttpOnly-cookie `mcp_onboarding`.
- **Порядок сжигания OTP** (`link-code/verify`): код сгорает только после проверки установки/аккаунта — сбой lookup больше не уничтожает валидный код пользователя.
- **TTL `pending_link` (48 ч)**: протухшая заявка на привязку сбрасывается (включая `user_id`) — anti-takeover: autologin не выдаёт сессию без логина по старой привязке.
- **409 «account conflict»** — редирект на SPA-экран с пояснением вместо голого текста.
- **Уведомления NIST**: письма «Магазин привязан» / «Магазин отвязан» при привязке и удалении приложения (best-effort).

### Изменено
- `POST /api/exchange-code`: код из cookie или body (fallback), одноразовый.
- `NewAppHandler`: убран параметр `hmacSecret` (HMAC-коды заменены серверным стором).

### Добавлено
- Шаблоны писем `shop_linked` / `shop_unlinked`.
- Тесты: race exchange (один победитель), CSRF-негатив, e2e cookie-handoff, anti-takeover TTL, порядок burn OTP, сид шаблонов.

## [Unreleased] — nginx config cleanup

### Исправлено
- **nginx warn'ы**: `listen … http2` (deprecated) → `http2 on;`; убран `ssl_stapling`/`ssl_stapling_verify` (Let's Encrypt серт без OCSP responder URL — stapling игнорировался). `nginx -t` без предупреждений.

## [0.7.5] — 2026-10-25 (expired-login cache leak)

### Исправлено
- **Кеш при протухшем JWT**: после истечения токена браузер мог показывать прежнюю страницу (DOM/JSON из кеша или bfcache) вместо формы логина. Теперь:
  - nginx: HTML/SPA/API/MCP/health — `Cache-Control: no-store`; JS/CSS — `private, max-age=3600`; security-заголовки (HSTS, CSP, …) больше не теряются на статике (общий `snippets/security-headers.conf`)
  - Go: `NoStoreMiddleware` на все ответы API
  - SPA: `API.isLoggedIn()` проверяет `exp` JWT (base64url); `pageshow` + bfcache restore → повторная проверка / login

### Добавлено
- `TestNoStoreMiddleware`
- Сниппет `web/nginx/snippets/security-headers.conf`

## [0.7.4] — 2026-10-23 (clear secrets on empty field + version alignment)

### Изменено
- **Стиране пароля при пустом поле**: SMTP (`pass: ""`) и API-пароль магазина (`api_password: ""`) теперь очищаются; поле не отправлено (omit) = оставить как есть. UI спрашивает подтверждение перед стиранием. Раньше пустое поле только «не меняло» пароль.
- **PostgreSQL 16**: compose / deploy.sh / документация; prod-миграция 15→16 через `pg_dump` (см. INSTALL).
- **Версии продуктов**: Go ≥1.25, golangci-lint v2.15.0, примеры/UI 0.7.4, MCP protocolVersion 2025-11-25.

### Добавлено
- `TestUpdateSMTPSettings_EmptyPassClears`, `TestShopUpdate_EmptyPassClears_OmittedKeeps`
- `has_password` в ответе `/api/shops` (для подсказки/подтверждения в UI)

## [0.7.3] — 2026-10-22 (form helper + docs)

### Изменено
- **`App.readFields` / `App.check`** — единый разбор полей формы (`bool|text|trim|int`); `saveServer` / `savePDN` / `saveSMTP` / create-edit магазинов и AI-ключей переведены с разрозненных `getElementById().checked/.value`
- **Документация**: `ui-smoke.sh` в CONTRIBUTING и README (UI-регресс рядом с `make lint` / `go test`)

## [0.7.2] — 2026-10-22 (golangci-lint gate + auth error handling)

### Исправлено
- **Смена пароля / сброс** — ошибки `UpdateRequirePasswordChange`, `ClearResetTokens` игнорировались: пользователь мог застрять с `require_password_change=true`, reset-токены не инвалидировались после использования; теперь возвращается 500
- **`CleanExpiredResetTokens` / `SendResetPassword`** — ошибки логируются (раньше молча терялись)
- **AI-ключи** — `UpdatePDNEnabled` при создании ключа: ошибка → 500 (раньше ключ создавался без PDN-флага)
- **Audit log** — `Insert` на nested write-tools теперь логирует ошибку
- **`tx.Rollback`** в `CreateSuperAdmin` — явная обработка

### Инструменты
- **golangci-lint v2.15.0 gate**: `make lint` / `make lint-fix`, `.golangci.yml` (govet, staticcheck, errcheck, ineffassign, gofmt), CI job `lint`
- Исправлены errcheck/staticcheck находки по auth/storage/backup/pdn/certbot
- `.project_rules.md`: критерий готовности — `make lint` / `go vet` / `go test` зелёные

## [0.7.1] — 2026-10-22 (UI audit fixes)

### Исправлено
- **ПДн-настройки не сохранялись** — `PUT /api/admin/pdn` только сбрасывал кэш, чекбоксы «отскакивали» после «Сохранить»; теперь пишутся в `server_settings`
- **`pdn_enabled` терялся при создании магазина** — INSERT не включал колонку; список магазинов тоже не отдавал флаг
- **Quill «Визуально» пустой** после SPA-перехода (уход/возврат на `/admin/email`) — редактор оставался привязан к удалённому DOM; HTML/Plain text при этом работали. Регрессия `b086f98`
- **Чекбоксы «выше/ниже» подписи** — `input { width:100% }` растягивал checkbox; единый hint `.form-check-hint`
- **Ошибки записи настроек** (`setSetting` / `settingsRepo.Set`) больше не игнорируются — UI получает 500 вместо ложного «сохранено»

### Добавлено
- `showPage()` в `pagination.js` — единая точка перехода по страницам (без своих `goPage`)
- `scripts/ui-smoke.sh` — регресс-guard: login, sticky чекбоксы, Quill после SPA re-mount
- Регресс-тесты: PDN persist, shop create `pdn_enabled`
- CI job `integration` (Postgres) в `.github/workflows/ci.yml`

## [0.7.0] — 2026-10-20 (backlog A/C)

### Добавлено
- **outputSchema** на 11 core tools (`WithOutputSchema`); `list_products`/`list_orders`/`list_clients` → envelope `{items,total,page,per_page}`
- **Audit logging** write-tools: таблица `audit_log` (args только hash — PDN), middleware + `call_insales_tool` nested path
- **Интеграционные тесты PostgreSQL** (`-tags integration`): app_installs, email_templates, ai_keys.shop_id ON DELETE SET NULL, audit_log; CI job `integration` с postgres:16
- **server.json** (MCP Registry remote-only card) + тест синхронизации версии
- **BILLING_MONTHLY** — опциональный RecurringApplicationCharge при установке InSales-приложения (0 = бесплатно)
- RecurringApplicationCharge: пути приведены к singleton API InSales (`/recurring_application_charge.json`)

### Документация
- ARCHITECTURE: «частично» = лимиты InSales API; roadmap MCP 2026-07-28 / OAuth 2.1 / rename
- INSALES-APP: раздел «Биллинг»

## [0.6.0] — 2026-10-20 (MCP coverage + discovery)

### Добавлено
- **284 MCP tools** (120 read-only): поля товаров/вариантов/коллекций (CRUD значений id **и** handle), order lines/discounts/shipping/custom status, menu items, delivery/payment CRUD, accessory links, collect move/position/remove, update recurring application charge
- **Progressive discovery**: default `tools/list` ≈ core surface (~30) + `search_insales_tools` / `call_insales_tool`; `MCP_TOOLS_MODE=all` — полный каталог
- **LLM-формат**: результат JSON text + `structuredContent`; envelope `{items,total,page,per_page}`; id/`product_field_id` как string
- **Описания tools**: EN description + RU `title` (UI); language-invariant args
- **Аннотации** у всех tools: `readOnlyHint` / `destructiveHint` / `idempotentHint` / `openWorldHint`
- **Матрица покрытия** InSales API → tools в `ARCHITECTURE.md` (57 ресурсов / 327 ops)

### Исправлено
- **`create_order` терял payload** — строки/адрес/скидки не доходили до InSales
- **`update_order`** — `order_id` утекал в body; `status` теперь alias `financial_status`
- **`toolFilter`** — tool без `ReadOnlyHint` больше не считается read-only
- **`field_values` в `get_product`** — join `handle`/`title` из справочника полей
- **docs** — counts 74/96/225 → 284; пример `initialize` без `2024-11-05`/`0.1.0`; матрица Collect/RecurringApplicationCharge → OK

## [0.5.3] - 2026-10-19

### Исправлено
- **Галочка «Включить отправку почты»** — больше не слетает после сохранения; enable без хоста → явная ошибка
- **Тест-письмо не врёт** — при выключенном SMTP ответ «SMTP выключен», а не ложное «отправлено»
- **SSL 465** — используется `WithSSL` (implicit TLS); AUTH auto-discover вместо только PLAIN
- **UI**: нет 👁 у SMTP-пароля (сохранённый пароль не отдаётся в браузер); статус «Отправка: включена/выключена»

### Документация
- **mail.ru / biz.mail.ru**: пароль приложений (не основной), ошибка 535, закрытые SMTP-порты на VPS (`EMAIL-DNS.md`, `INSTALL.md`, `README.md`, подсказка в админке)

## [0.5.2] - 2026-10-17

### Исправлено
- **Редактор шаблонов: пустое «Визуально»** — Quill не учитывал `innerHTML` (Delta оставался пустым); загрузка через `clipboard.dangerouslyPasteHTML`
- **Сохранение не стирает тело шаблона** — защита: пустой редактор больше не затирает непустой HTML
- **UI SMTP**: чекбокс «Включить отправку почты» в отдельном блоке; у переменных `{{Name}}` / `{{ResetURL}}` / `{{VerifyURL}}` / `{{MCPURL}}` всплывающие подсказки с примерами

## [0.5.1] - 2026-10-17

### Добавлено
- **Почта: шаблоны писем в админке** (`/admin/email`) — редактирование темы и тела, Quill-редактор, HTML/plain-режимы, переменные, preview (sandboxed iframe), тест-письмо, сброс к умолчанию. Роли: admin и super_admin
- **SMTP-настройки без рестарта** — `GET/PUT /api/admin/smtp/settings` (host/port/user/pass/from/tls_mode), hot-reload, пароль AES-GCM, env как фоллбэк
- **go-mail** вместо замороженного `net/smtp`: RFC 2047 (кириллица в теме), multipart plain+html, TLS-режимы (ssl/starttls/none)
- **Ретраи отправки** — 3 попытки (1s/5s/30s) на транзиентные ошибки SMTP
- Таблица `email_templates` (сид 4 шаблона: reset, verify, welcome, test)

### Изменено
- `SendVerifyEmail` удалён из API сервиса (шаблон `verify_email` остался в реестре)
- Валидация шаблонов на PUT: `template.Parse` + sample-exec → 400 без panic

## [0.5.0] - 2026-10-17

### Добавлено
- **Стороннее приложение InSales** — установка из бэк-офиса магазина: `GET /install` (пароль = `MD5(token + secret_key)`), `GET /uninstall`, `GET /app` → автологин `GET /autologin` (проверка `token3`). Переменные `INSALES_APP_KEY` / `INSALES_SECRET_KEY`
- **Гибридная выдача доступа** — новый email: автосоздание аккаунта + MCP-ключ на экране `/onboarding`; существующий email: привязка магазина только после входа паролем (защита от захвата аккаунта)
- **Ключи, привязанные к магазину** — `ai_keys.shop_id`; запросы MCP-ключа идут в свой магазин (а не только в default)
- Таблица `app_installs` (журнал установок), one-time code `POST /api/exchange-code`, `POST /api/shops/link-pending`
- Документация: `docs/setup/INSALES-APP.md`, `docs/insales-app-form.md`, `docs/api/api_and_apps.md`

### Изменено
- Удалены неиспользуемые `INSALES_SHOP` / `INSALES_API_KEY` / `INSALES_PASSWORD` (ключ магазина добавляется в веб-интерфейсе)
- `REGISTRATION_ENABLED` без значения в `.env` теперь **закрывает** регистрацию (раньше код считал её открытой — расходился с документацией и `.env.example`)

## [0.4.7] - 2026-10-17

### Добавлено
- **Блокировка публичной регистрации** — переменная `REGISTRATION_ENABLED` в `.env`. При `false` `POST /api/register` возвращает 403. Исключение: если пользователей нет, первый регистрацией становится `super_admin` (bootstrap)
- **Создание пользователя из админки** — `POST /api/admin/users` (admin/super_admin) и кнопка «Добавить пользователя» на `/admin/users`. Временный пароль, при первом входе принудительная смена

## [0.4.6] - 2026-10-16

### Добавлено
- **Пагинация журнала бэкапов и инцидентов** — `/api/admin/backup/log` и `/incidents` принимают `page`/`per_page` и возвращают `total`; в админке кнопки «Назад/Далее»
- **Ротация журнала и инцидентов** — `cleanup_old_backups` удаляет из `backup_log`/`backup_incidents` записи старше `retention_days` (как и файлы бэкапов)
- **Ротация системных логов** — logrotate для `/var/log/mcp-is-backup.log`, лимиты docker-логов (`10m×5`) для `db` и `nginx`

### Исправлено
- **Журнал бэкапов был пуст при `total > 0`** — pgx не сканирует `timestamptz` в `string`; добавлены `ts::text` и `COALESCE` (как в user_repo/pdn)
- **Чистка не выполнялась при пропуске бэкапа** — `cleanup_old_backups` вызывается до проверок окна/нагрузки в `--scheduled`

## [0.4.5] - 2026-10-16

### Исправлено
- **Бэкапы не работали (0 бэкапов, «backup failed»)** — три причины:
  1. кнопка «Бэкап сейчас» падала в app-контейнере: нет `bash` (добавлен в образ)
  2. ночные бэкапы (хост cron) всегда откладывались: целочисленная арифметика давала порог 0 на 1-ядерном VDS (`cores*pct/100 = 0`), условие «нагрузка >= 0» всегда истинно — проверка теперь в % от ядер
  3. скрипт использовал `docker compose exec`, которого нет в app-контейнере — заменено на `docker exec` (работает и на хосте, и в контейнере)

## [0.4.4] - 2026-10-16

### Исправлено
- **Auto-renew Let's Encrypt** — `certbot.go` вызывал `docker exec` из app-контейнера, где нет docker-клиента: все 10 попыток обновления падали. В образ добавлен `docker-cli`, в compose — `docker.sock` и монтирование `/etc/letsencrypt` (чтение статуса сертификата). Отключён random-delay certbot (`--no-random-sleep-on-renew`)

## [0.4.3] - 2026-10-16

### Добавлено
- **Показ/скрытие пароля** — кнопка «глаз» на всех полях ввода пароля (вход, регистрация, сброс, смена пароля, API-пароль магазина)

### Исправлено
- **Healthcheck nginx** — `localhost` → `127.0.0.1` (wget уходил в IPv6 ::1) и HTTPS (порт 80 отдаёт 301) — контейнер 2 месяца был unhealthy

## [0.4.2] - 2026-08-02

### Добавлено
- **SMTP тест для админа** — `POST /api/admin/smtp/test` отправляет тестовое письмо (super_admin). `GET /api/admin/smtp/status` показывает статус SMTP-конфигурации
- **Панель SMTP в админке** — на странице «Настройки сервера» отображается статус SMTP (включён/выключен, сервер, отправитель) и кнопка отправки тестового письма
- **Логирование URL восстановления** — когда SMTP выключен, URL сброса пароля выводится в логи сервера для ручной передачи пользователю
- **Очистка протухших токенов** — `CleanExpiredResetTokens()` автоматически удаляет expired токены из `password_resets` при каждом запросе сброса пароля

### Исправлено
- **INSTALL.md** — исправлены имена SMTP-переменных (`SMTP_USERNAME` → `SMTP_USER`, `SMTP_PASSWORD` → `SMTP_PASS`)

## [0.4.1] - 2026-07-29

### Добавлено
- **AI Key Name в конфигах** — имя AI-ключа теперь используется как JSON-ключ в MCP-конфиге. Ключ "Cursor MCP" → конфиг с `"cursor-mcp"`. Slugify: кириллица → транслитерация, пробелы → дефисы, только латиница и цифры. Fallback на "mcp-is" если имя пустое
- **Outlet (точки выдачи)** — CRUD: `list_outlets`, `get_outlet`, `create_outlet`, `update_outlet`, `delete_outlet`
- **Marketplace (маркетплейсы)** — CRUD: `list_marketplaces`, `get_marketplace`, `create_marketplace`, `update_marketplace`, `delete_marketplace`. Поддержка 3 типов: Google Merchant, Liquid Yandex Market, Model Yandex Market
- **ProviderToken** — `create_provider_token` для генерации токенов мессенджеров (WhatsApp, Telegram, Viber, SMS, Email, VK, Instagram)
- **Account** — `update_account` для обновления настроек аккаунта
- **Count endpoints** — `get_orders_count` и `get_products_count` для быстрого получения количества
- **Пайплайн `collection_products`** — товары коллекции с полной информацией одним вызовом
- **Data Export (CSV)** — экспорт данных в CSV-формате:
  - `export_products` — товары (id, title, price, sku, vendor, quantity, created_at)
  - `export_orders` — заказы (id, number, fulfillment_status, total_price, client_name, created_at)
  - `export_clients` — клиенты (id, name, surname, email, phone, created_at)
- **Integration tests** — тесты для InSales client (Outlet, Marketplace, ProviderToken, Account, Counts), MCP handlers (CSV, integration), storage (Settings, MCPClientTypes)

### Статистика
- 18 новых MCP-инструментов (74 → 92)
- 25+ новых unit/integration тестов

## [0.3.2] - 2026-07-25

### Добавлено
- **Пайплайн-инструменты** — композитные инструменты для частых сценариев:
  - `search_product_details` — поиск товара с кастомными полями, метаполями и изображениями за 1 вызов
  - `order_full_details` — полная информация о заказе (клиент, адрес, позиции, метаполя) за 1 вызов
  - `client_full_profile` — полный профиль клиента (заказы, бонусы, метаполя) за 1 вызов
  - `product_catalog_search` — поиск по каталогу с фильтрами и пагинацией
- **product_field_values** — `get_product` теперь возвращает кастомные поля товара (Противопоказания, Состав и т.д.)
- **Расширены типы Product** — добавлены description, short_description, permalink, is_hidden, available, properties, characteristics, option_names, images, field_values
- **Расширены типы Order** — добавлены shipping_address, fields_values
- **Расширены типы Client** — добавлены bonus_points, progressive_discount, group_discount, client_tags, full_name, fields_values

### Исправлено
- **URL-кодирование** — ListProducts и ListOrders теперь кодируют кириллицу и спецсимволы через `url.QueryEscape()`
- **Meta fields фильтрация** — `list_meta_fields` теперь поддерживает фильтры `resource_type` и `resource_id`

### Тесты
- 5 тестов URL-кодирования (кириллица, пробелы, спецсимволы)
- 4 теста мапперов (Product field values, Order shipping address, Client fields)
- 2 теста MetaFields фильтрации

## [0.3.0] - 2026-07-23

### Добавлено
- **~40 новых MCP-инструментов** по полному API InSales (225 инструментов итого)
- **DiscountCodes** (промокоды) — list, get, create, update, delete
- **Collects** (товар-коллекция) — list, getByProduct, add, remove
- **Properties** (свойства товаров) — list, get, create, update, delete
- **Characteristics** (характеристики свойств) — list, get, create, update, delete
- **Reviews** (отзывы) — list, get, create, update, delete
- **Files** (файлы) — list, get, create (URL), delete
- **Domains** — CRUD
- **ClientGroups** — CRUD
- **CustomStatuses** — CRUD
- **ProductFields, VariantFields, CollectionFields** — CRUD
- **CollectionFilters** — CRUD
- **JsTags** — CRUD
- **ExternalDiscounts** — CRUD
- **DiscountCollectionsGroups** — CRUD
- **Warehouses** — CRUD
- **SalesChannels** — CRUD
- **StockCurrencies** — CRUD
- **PriceKinds** — CRUD
- **OptionNames** — CRUD
- **OptionValues** — CRUD (nested under option_names)
- **ApplicationActions** — CRUD
- **ApplicationCharges** — list, get, create, decline
- **ApplicationWidgets** — CRUD
- **RecurringApplicationCharges** — CRUD
- **BonusSystemTransactions** — list, get, create, listByClient
- **OAuthApps** — CRUD
- **UserTasks** — CRUD
- **Supplementaries** — list, get, create, delete
- **Similars** — listByProduct, create, delete
- **Units** — list (read-only)
- **update_blog** — инструмент для обновления блога (client метод уже был)

### Исправлено
- **UpdatePropertyValue/DeletePropertyValue** — URL изменён с legacy `/accessory_values/:id.json` на официальный `/accessories/:accessory_id/values/:value_id.json` (требуется accessory_id)
- **CreateAccessory** — добавлены обязательные поля `max_count`, `min_count`, `permalink`
- **CreateClient** — обязательные поля заменены с `first_name/last_name` на `name/surname/middlename` (по official API)
- **list_orders** — добавлены фильтры `fulfillment_status`, `updated_since`, `from_id`
- **list_products** — добавлены фильтры `updated_since`, `from_id`, `with_deleted`, `deleted`
- Generic CRUD helpers для клиента, сервиса и хендлеров — сокращение дублирования кода

### Файлы
- Добавлен `docs/api/Insales-api-docs.txt` — полная документация API InSales
- `docs/api/insales-api-reference.md` — структурированный MD-референс (в процессе)

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
- **Bot access control** — переменные `BLOCK_AI_BOTS` и `BLOCK_SEARCH_BOTS` в .env
  - Блокировка AI-скраперов (GPTBot, CCBot, Anthropic-AI, Bytespider и др.) через nginx User-Agent
  - Генерация `robots.txt` из переменных окружения
  - 21 unit-тест для entrypoint.sh
- **Расширен список совместимых MCP-клиентов** — 14 систем в документации (Claude Desktop, Cursor, VS Code, MiMoCode, OpenCode, Continue, LibreChat, Open WebUI, Cline, Cherry Studio, ChatGPT Desktop, Windsurf, Zed, Aider)

### Исправлено
- `.gitignore` — паттерн `server` заменён на `/server` (binaries only)
- `stop-impersonate` — убрана лишняя проверка adminMW (impersonated user имеет роль `user`)
- `server_settings` — `InitDefaults()` перезаписывала domain/https/cert_issuer при каждом старте. Заменена на `SetDefaultsIfMissing()` (вставка только отсутствующих ключей)
- **MCP read-only фильтр** — mcp-go по умолчанию считает `ReadOnlyHint=false`, инструменты без аннотации фильтровались в read-only режиме. Добавлены `WithReadOnlyHintAnnotation(true)` к 31 read-инструменту
- **HTML-ошибки InSales API** — MCP-сервер возвращал сырой HTML при 404/500 ошибках API. Теперь возвращает чистые JSON-ошибки
- **get_shop_info** — парсил storefront HTML. Теперь использует `GET /admin/account.json` (официальный API)
- **PropertyValue URL** — использовал legacy `/property_values.json`. Исправлен на `/accessories/:id/values.json` (официальный API)
- **Accessory** — поле `title` заменено на `name` (соответствует официальному API)
- **BlogPost URL** — `/blog_posts.json` заменён на `/articles.json` (официальный API)
- **Blog articles** — URL исправлен на официальный `/blogs/:id/articles.json` (был `/articles.json`)
- **Create Article** — добавлены обязательные поля `content`, `published_at`, `author`
- **Update Blog** — добавлен новый инструмент (PUT /blogs/:id.json)
- blog_posts/articles — InSales API возвращает 404 для read-only ключа (ограничение API)

### Тесты
- 19 новых unit-тестов (3 JWT + 6 middleware + 10 handler)
- 21 bash-тест для entrypoint.sh (robots.txt, nginx config)

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

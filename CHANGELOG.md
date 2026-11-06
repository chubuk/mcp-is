# Changelog

## [0.10.0] — 2026-11-06 (capabilities probe, политика ИИ-ключей, диагностика для LLM)

### Добавлено
- **Живая проверка прав ключа магазина (capabilities probe)**: 25 GET-доменов InSales (`?per_page=1`), пул 4 воркера, таймаут 30s. Снимок в `shops.capabilities` (JSON с `since`/`source`), TTL 24h, cooldown 1h. Probe только в 4 точках: создание/смена кредов, кнопка «Проверить доступы», TTL-лениво, первый 403. `tools/list` — ноль запросов к InSales.
- **Фильтрация инструментов по правам**: `toolFilter` скрывает инструменты denied-доменов; `call_insales_tool` — remediation-сообщение при denied; mcp-go `ToolFilterFunc` принудительно применяется и к `tools/list`, и к `tools/call`.
- **Политика ИИ-ключа (`ai_keys.command_policy`)**: JSON `{default, groups:{read/write/delete}, commands:{override}}`; приоритет override команды > класс в группе > default; deny побеждает. UI — редактор политики в создании/редактировании ИИ-ключа: пресеты, поиск, аккордеон групп, классы «Чтение/Изменение/Удаление», per-command override, summary-чипсы, diff перед сохранением, серые позиции «нет прав магазина с <дата>» по capabilities-снимку. `GET /api/tool-meta` — метаданные инструментов (name/group/class/domain) для редактора.
- **outputSchema для мутационных инструментов**: `addTool` автоматически ставит схему `MutationResult` (`update_*`/`delete_*`/`add_*`/`decline_*`) — типизированные MCP-клиенты получают контракт результата.
- **Модель состояний (политика × капабилити)**: `enabled` / `off_by_policy` (постоянно) / `unavailable_no_scope` (временно, авто-восстановление перепроверкой) / `unavailable_unknown` (fail-open); denied несёт `since`/`source` (probe|runtime_403).
- **`POST /api/shops/{id}/check-access`** (вместо `check-commands`): probe + сохранение снимка + ответ с доменами/инструментами.
- **`store_info`**: блок `capabilities` (checked_at, stale, credentials_ok, allowed/denied/unknown) + `direct_rest` (прямой REST недоступен).
- **Инструкции соединения**: блок о правах ключа, предупреждение про затирание пустыми значениями, `clear_fields`.
- **Маппинг инструментов (`tool_domains.go`)**: 287 инструментов (123 read / 110 write / 54 delete) → {probe-домен, функциональная группа, класс операции}; drift-тест `TestToolDomainCoverage`.
- **Mutation result (3c)**: delete-обработчики возвращают `{ok, action, resource, id, deleted:<снимок>}` с pre-fetch; update-обработчики — `{ok, action, resource, id, changed_fields:{old,new}, ignored_keys, data}` (pre-fetch +1 GET на мутацию, осознанно); `toolError` — remediation для 403/401.
- **handle/permalink унификация**: `handle` — единый публичный ключ, `permalink` — алиас (дублируется в выдаче).
- **ignored_keys на все мутации**: whitelist аргументов извлекается автоматически из `inputSchema` при регистрации (`addTool`); незнакомые ключи не уходят в InSales и возвращаются списком `ignored_keys` для всех create/update/delete.
- **inputSchema**: форматы полей (HTML/текст/slug/число) в description create/update-инструментов.
- **Лимит ИИ-ключей**: `MAX_AI_KEYS_PER_USER` 20 → 50.

### Изменено
- **`ai_keys.available_commands` — legacy**: производное чтение (политика × капабилити), зеркало не пишется; `UpdateCommands`/`UpdateCommandsByUserID` удалены.
- **Версия**: 0.9.7 → 0.10.0 (скрытие инструментов + переименование роута = breaking-семантика).

## [0.9.7] — 2026-11-04 (security-гейт gosec/govulncheck, закрытые уязвимости, WCAG-контраст)

### Добавлено
- **Go security-гейт** — `govulncheck ./...` в `make lint` и CI (анализ достижимости уязвимостей), `gosec` в golangci-lint (G104 исключён как дубль errcheck, тесты исключены)

### Безопасность
- **Достижимые уязвимости закрыты** (найдены govulncheck): `golang.org/x/text` v0.40.0→v0.41.0 (panic на crafted input в precis), `jackc/pgx/v5` v5.7.4→v5.9.2 (SQL injection через dollar-quoted placeholders)
- **`normalizeShopHost` ужесточён** — строгий шаблон host[:port] (закрывает smuggle `/@\\` в redirect-URL install-потока)
- **WCAG AA контраст в UI** — акценты затемнены (`#e67e22`→`#b85c15`, `#4a6cf7`→`#3b5bdb`), текст футера осветлён (`#666`→`#8a8a9c`), приглушённый токен `#888`→`#707070` (~28 вхождений); все пары ≥4.5:1 (impeccable, живой рендер — 0 находок)

### Исправлено
- **Мёртвый импорт `net/http/pprof`** — pprof-роуты регистрировались на DefaultServeMux, который не обслуживался; удалён
- **`readFile` в backup** — чтение через подпроцесс `cat` заменено на `os.ReadFile`
- **test.html: скрипт-ошибка при загрузке** — `app.js` route() ожидает `#header`/`#app`/`#impersonate-*`; добавлены скрытые DOM-стабы харнеса (61/61 assert зелёные)

## [0.9.6] — 2026-11-04 (FE-tooling, DOM-XSS-фиксы, dockerproxy, restore-фиксы)

### Добавлено
- **Lint-гейты инфраструктуры**: `make lint-nginx` (реальные тесты entrypoint + `nginx -t` + gixy), `make lint-docker` (hadolint + trivy config/secret + `docker compose config`), `make lint-frontend` (ESLint, safety-правила), `make lint-all`; в CI — job `lint-infra` (gixy HIGH — падение, hadolint/eslint — warn).
- **FE-тулчейн (JS/CSS/HTML)**: `eslint-plugin-no-unsanitized` (DOM-XSS), stylelint 18 + `stylelint-config-recommended`, html-validate, Lighthouse CI (бюджеты perf/BP/a11y/SEO, job `lhci` в CI), retire.js (CVE vendored-библиотек). `make lint-frontend` гоняет eslint+stylelint+html-validate+retire; lhci — `cd web && npx lhci autorun` (бюджеты в `web/lighthouserc.json`).

### Безопасность
- **DOM-XSS во фронтенде (41 сайт `innerHTML`/`insertAdjacentHTML`) закрыт**: добавлен `esc()` (HTML-экранирование, `api.js`) и экранированы данные API/ошибок/имена файлов перед вставкой в HTML; в `confirmRestore` исправлена JS-инъекция через имя файла в inline `onclick` (экранирование в JS-контексте через `JSON.stringify`); `billing_pay_url` в баннере биллинга проверяется на http(s)-схему. Осознанные присваивания HTML собственных render-функций помечены `eslint-disable-next-line` с причиной; правило `no-unsanitized` включено как `error`.
- **Quill 2.0.3 — известная уязвимость CVE-2025-15056 (low, XSS в HTML-export)**: фикса в upstream нет (2.0.3 — последняя версия); риск принят осознанно (редактор email-шаблонов — только для админов), retire печатает её informational (`--severity medium`).
- **docker.sock убран из app-контейнера**: Docker доступен только через allowlist-прокси `dockerproxy` (tecnativa/docker-socket-proxy v0.5.0, internal-сеть): разрешены только containers/exec/allow_start/allow_stop/images + POST — операции certbot и backup/restore; restore-исполнитель тоже без сокета (через `DOCKER_HOST`).
- **`.bak-*` больше не отдаются по HTTP**: deny-regex теперь покрывает `bak-<timestamp>` (раньше `*.bak-20261025-225344` были доступны); устаревшие `.bak-*` файлы удалены.

### Исправлено
- **`entrypoint_test.sh` реально тестирует `entrypoint.sh`** (раньше проверял собственные heredoc-заготовки вместо результата): 4 сценария / 30 проверок внутри nginx-образа.
- **Restore через API не запускался вообще**: `docker run image backup.sh --restore ...` упирался в `ENTRYPOINT ["./mcp-is"]` — сервер получал скрипт аргументами; добавлен `--entrypoint` (проверено живым restore через `/api/admin/backup/restore`).
- **dockerproxy allowlist: `ALLOW_STOP`** — restore-процедуре нужен `docker stop` контейнеров при свапе PGDATA (без флага свап шёл бы при живой БД).
- **CI `lint-infra`**: `nginx -t` получил `--add-host app:127.0.0.1` (на раннере с `--no-deps` upstream не резолвится); gixy — Python 3.12 (в 3.13 удалён `sre_constants`).

## [0.9.5] — 2026-11-02 (безопасная запись payload, честная пагинация, SEO-тулы, фиксы create/add)

### Добавлено
- **SEO-инструменты `get_seo` / `update_seo`** (core surface): просмотр и правка `html_title`/`meta_keywords`/`meta_description` у товаров, страниц, записей блога, блогов и коллекций (у коллекций также `seo_description`), очистка через `clear_fields`. SEO главной страницы через API недоступен (задокументировано в README).
- **`ignored_keys` в ответах `update_product`/`update_seo`**: лишние ключи аргументов не отправляются в InSales и возвращаются списком — агент видит, что его параметр проигнорирован.

### Исправлено
- **Пустые значения больше не затирают данные**: исходящие payload'ы чистятся от `nil`/пустых строк (числа `0`/`false` сохраняются — иначе не записать нулевой остаток); эхо-объект `get_product` с пустыми полями больше не перезаписывает товар. Явная очистка поля — только `clear_fields` (с валидацией имён).
- **`update_product` по whitelist**: схема расширена (`short_description`, `vendor`, `sku`, SEO-поля, `category_id`, `clear_fields`), `quantity` и `price` убраны — цена/остатки живут в вариантах и меняются через `update_variant` (попытка записи через `update_product` была тихим no-op); алиас `keywords`→`meta_keywords`.
- **Честная пагинация**: значения `per_page` 1–9 (InSales клампит снизу до 10) эмулируются окном из 250-страниц; `total` в конверте = число элементов в ответе (задокументировано в описаниях и инструкциях); `per_page` принимает и строку, и число; `per_page` < 1 трактуется как не заданный (дефолт 25, раньше молча получали 1).
- **Жёсткое отсутствие кеширования**: исходящие запросы к InSales — `Cache-Control: no-cache, no-store`; ответы `/mcp` — `no-store`.
- `create_collection` принимает `parent_id` (без него InSales отвечал 422 «parent_id: Не заполнены обязательные поля»).
- **`create_product` теперь реально создаёт товары**: payload собирался без обязательного `variants_attributes` (InSales отвечал 422). Цена/артикул/остаток уходят в вариант, как требует API.
- **`add_image` прикрепляет изображения**: payload слал `image[url]`, а InSales принимает только `image[src]` (422 «Не заполнены обязательные поля»).

## [0.9.4] — 2026-11-02 (429 = JSON-RPC error + Retry-After, лимиты в initialize.instructions)

### Исправлено
- **HTTP 429 при переполнении очереди лимита теперь понятен агентам**: вместо nginx-HTML возвращается JSON-RPC error (`code:-32000`, `message` с инструкцией подождать и повторить, `data.retry_after`) + заголовок `Retry-After: 30` (спека MCP разрешает JSON-RPC error body в HTTP-ошибках; текст actionable для LLM self-correction).
- **Лимиты объявлены в `initialize.instructions`** (раздел Rate limits, EN): агенты знают про 300 req/min per key, замедление при приближении к лимиту и поведение при 429 — до первого вызова.

## [0.9.3] — 2026-11-02 (rate limiting per-API-key, замедление вместо 503)

### Исправлено
- **Rate limiting MCP переписан на per-API-ключ** (nginx): раньше лимит 30/мин был на IP-адрес — один агент выбивал бюджет всем за одним NAT и получал HTTP 503. Теперь 300/мин на каждый API-ключ (как лимиты InSales), burst 60 мгновенно, плавное замедление выдачи до очереди 200 (nginx `delay`), отказ — 429 вместо 503. Отдельный щит от флуда 100/с на IP.

## [0.9.2] — 2026-11-01 (фикс: boolean-схемы в outputSchema реально применяются)

### Исправлено
- **Нормализация outputSchema теперь покрывает все 285 инструментов**: предыдущий фикс (0.9.1) был в обёртке `registerTool`, через которую шёл лишь 1 инструмент из 282 — остальные регистрировались напрямую через `s.AddTool` и отдавали boolean-схемы (`items:true`, `additionalProperties:true`), из-за чего строгие типизированные MCP-клиенты (Hermes) падали на discovery. Единая точка `addTool` + тест-регрессия `TestNoBoolSchemasAfterRegistration`.

## [0.9.1] — 2026-11-01 (outputSchema-комплаенс, рантайм-валидация, ремонт маскирования ПДн)

### Исправлено
- **Маскирование ПДн (152-ФЗ) реально работает на всех путях**: раньше типизированные результаты (list_orders, get_order) и часть инструментов (create/update_order, order_sub_ops, order_full_details, client_full_profile) отдавали ИИ реальные email/телефоны/ФИО при включённом режиме. Теперь маска рекурсивна (вложенные client/shipping_address/элементы списков), покрывает name/surname/middlename/full_name/formatted_phone/client_name, не теряет не-строковые значения, применяется во всех order-клиентских путях и в CSV-экспортах (export_orders/export_clients).
- **outputSchema без boolean-схем**: свойства `true` (генерировались для полей `any`) заменены на `{}` — строгие типизированные MCP-клиенты больше не отваливаются на discovery.
- **Рантайм-валидация structuredContent против outputSchema** (opt-in mcp-go): расхождения теперь видны как isError вместо молчаливой рассинхронизации.

## [0.9.0] — 2026-11-01 (MCP-compliance: structuredContent всегда объект)

### Исправлено
- **`structuredContent` по MCP-спеке — всегда JSON-объект** (CallToolResult, 2025-06-18): 52 инструмента списка (list_categories, list_variants, list_images и др.) отдавали массив и ломали строгие MCP-клиенты. Единая обёртка `wrapEnvelope` в `toolResult`: массив → `{items, total}`, скаляры/прочее не-объектное → `{value}`. Формат текстового контента меняется тем же образом (двойная форма идентична). `call_insales_tool` чинится автоматически (passthrough).
- Формат ответов list-инструментов единообразен с домашним конвертом `{items, total}` (семантика `total` = число возвращённых элементов, как у ProductsPage).

## [0.8.0] — 2026-10-31 (обязательная привязка ИИ-ключей к магазинам, именование соединений, лимиты, перевыпуск)

### Добавлено
- **Обязательная привязка ИИ-ключа к магазину**: `POST /api/ai-keys` требует `shop_id` (магазин пользователя); без магазинов в аккаунте создание запрещено (400 «Сначала добавьте магазин», кнопка в UI скрыта). Один ключ — один магазин.
- **Редактор легаси-ключей**: `PUT /api/ai-keys/{id}` (merge-patch) не сохранит непривязанный ключ без `shop_id`; для привязанных частичные апдейты (toggle ПДн) работают как раньше. Update/Delete/Rotate проверяют владельца ключа (раньше — нет).
- **`connection_name`** — алиас соединения ИИ-ключа: дефолт `slug(магазин)`, редактируемый, `^[a-z0-9-]{1,64}$`, уникален в аккаунте (авто-суффикс `-2` + partial unique index в БД). Подставляется в примеры конфигов MCP-клиентов вместо slug названия ключа.
- **Идентификация соединения для ИИ**: персональные `instructions` (алиас, магазин, правило «не угадывай запись») через `OnAfterInitialize`; инструмент `store_info` — ground truth (алиас/магазин/домен/привязка) без вызова InSales.
- **Непривязанные ключи — fail-closed**: data-инструменты возвращают модели actionable-ошибку («попросите владельца привязать ключ…») вместо молчаливого фолбэка на дефолт-магазин.
- **Перевыпуск ключа** `POST /api/ai-keys/{id}/rotate` + кнопка в редакторе: новый секрет, старый перестаёт работать немедленно, имя/привязка/алиас сохраняются. Штатное восстановление потерянного ключа.
- **Лимиты**: `MAX_SHOPS_PER_USER` / `MAX_AI_KEYS_PER_USER` (env, дефолт 20 каждая; превышение — 403). Автопровижининг установки InSales — вне квот.
- **Страница установки приложения**: кроме ключа и URL — примеры конфига с алиасом и рабочая кнопка «Скачать файл»; автоключ создаётся с `connection_name`, ответы `ExchangeCode`/`LinkPending` возвращают его.

## [0.7.6] — 2026-10-30 (billing gate, безопасность аккаунтов, restore из бэкапа, настраиваемая цена)

### Добавлено
- **Биллинг-гейт MCP-данных**: вход в приложение всегда бесплатный; при неоплаченном периоде вызовы MCP-инструментов магазина возвращают мягкое пояснение (ИИ сообщает о периоде вместо данных). Статус проверяется лениво (TTL `active` 24ч / `unpaid` 1ч), `billing_until` — локально на каждом вызове (отсечение в середине дня). Баннер «Оплатить/Продлить» (`billing_pay_url`).
- **«Стоимость подписки» в настройках сервера** (Админ → Сервер, super_admin): `billing_monthly` в `server_settings`, hot-reload без рестарта. `0` = бесплатно (гейт прозрачен), `>0` = сумма RecurringApplicationCharge. Env `BILLING_MONTHLY` = значение по умолчанию/фолбэк. Настройка не пересоздаёт уже созданный charge.
- **Восстановление из бэкапа в админке** («Бэкапы», super_admin): список файлов на диске (дата/тип/размер/sha256), «Проверить» (целостность без изменений в БД), «Восстановить» (модалка с предупреждением о простое), «Копировать команду» (CLI-фолбэк), поллинг статуса, блок «Как это работает».
- **Restore-скрипт переписан (steps 0–6)**: проверка целостности (sha256/gzip/tar/backup_label, обе раскладки тара) → страховочный дамп `pre_restore_*.sql.gz` (tripwire >1KB + users) → свап → верификация → **автоматический откат** при любом сбое (в т.ч. kill mid-swap). Статус-файл `restore_status.json`, `.sha256`-сайдкары на новых бэкапах.
- API: `GET /api/admin/backup/files`, `POST …/files/{name}/verify`, `POST …/restore` (super_admin, single-flight 409, preflight образа), `GET …/restore/status`. Path-escape юнит-тесты.
- `backup.sh --verify <file>` → `readable|corrupt`.
- **Привязка магазина к существующему аккаунту через OTP-код** (`/api/link-code/send|verify`): 6 цифр, TTL 10 мин, одноразовый, ≤5 попыток → инвалидация, cooldown ресенда ~90с. Верифицированный аккаунт линкуется без смены пароля (вариант B); «мёртвый» (неподтверждённая регистрация) — пароль отзывается + `require_password_change`.
- `favicon.ico` (16/32/48) + `favicon.svg` (метка «M» в цветах интерфейса), `<link rel="icon">` в `index.html`. nginx/deb не менялись — `.ico` уже отдаётся статикой.
- Шаблоны писем `login_code`, `shop_linked` / `shop_unlinked`.
- INSTALL.md: раздел «Если восстановление не удалось»; убрана ложь «дамп поверх работающей БД (без остановки)». TROUBLESHOOTING.md: «Битый бэкап / нет доступа в админку».
- Тесты: race exchange (один победитель), CSRF-негатив, e2e cookie-handoff, anti-takeover TTL, порядок burn OTP, path-escape backup API, сид шаблонов.

### Безопасность
- **Утечка MCP-ключа в URL устранена**: one-time код autologin больше не содержит сырой ключ (был внутри HMAC payload) и не передаётся через URL (`/onboarding?code=…` утекал в историю/логи/Referer). Теперь серверный одноразовый стор (TTL 2 мин, атомарное сжигание — replay невозможен), передача через HttpOnly-cookie `mcp_onboarding`.
- **Verify-before-use для публичной регистрации** (инвариант «пароль = владелец ящика»): `users.email_verified` — Register создаёт аккаунт с нерабочим паролем, пароль ставит подтвердивший (форма «код + пароль»); login-gate 403 `require_email`; lazy-чистка неподтверждённых регистраций (24ч).
- **Смена e-mail пользователем удалена** (anti pre-hijack, USENIX Classic-Federated Merge) — только вручную администратором; e-mail нормализуется lower+trim при записи/поиске.
- **Порядок сжигания OTP** (`link-code/verify`): код сгорает только после проверки установки/аккаунта — сбой lookup больше не уничтожает валидный код пользователя.
- **TTL `pending_link` (48 ч)**: протухшая заявка на привязку сбрасывается (включая `user_id`) — anti-takeover: autologin не выдаёт сессию без логина по старой привязке.
- **409 «account conflict»** — редирект на SPA-экран с пояснением вместо голого текста.
- **Уведомления NIST**: письма «Магазин привязан» / «Магазин отвязан» при привязке и удалении приложения (best-effort).
- **Флуд-защита почтовых операций**: отдельный rate limit (5/ч) на register/forgot-password/verify/link-code.

### Изменено
- `POST /api/exchange-code`: код из cookie или body (fallback), одноразовый.
- `NewAppHandler`: убран параметр `hmacSecret` (HMAC-коды заменены серверным стором).

### Исправлено
- **Restore был сломан**: несовпадение раскладки тара (`pg_basebackup/` vs корень), `rm -rf` данных на работающей БД, sql-наложение без DROP SCHEMA, WAL распаковывался в корень PGDATA вместо `pg_wal/`.
- **`email_verified` не выбирался в `GetByID`/`GetByEmail`** — zero-value false запускал условный seize верифицированных пользователей при OTP-линке.
- Общий lock backup+restore в `$BACKUP_DIR/.backup.lock` (mkdir-атомарный, stale >2ч) — взаимное исключение восстановления и бэкапа.
- `pre_restore_*` исключены из `find -size -1k`-чистки.
- **nginx warn'ы**: `listen … http2` (deprecated) → `http2 on;`; убран `ssl_stapling`/`ssl_stapling_verify` (Let's Encrypt серт без OCSP responder URL — stapling игнорировался). `nginx -t` без предупреждений.

### GATE
- 13/13 сценариев restore на throwaway-compose (happy path tar+sql, битые архивы, отказ дампа, kill mid-swap, мусорный PGDATA, lock, 409, disk-full, обрыв sql) — см. план `1793355392246-quick-cactus`.

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
- **OpenAPI spec** — полная спецификация REST API (28 эндпоинтов)
- **GitHub Actions CI/CD** — lint → test → build → docker → deploy
- **Backup script** — автоматический бэкап PostgreSQL с retention 7 дней
- **Certbot unit tests** — 13 тестов для `internal/certbot` (GetCertInfo, CheckRenewal, ShouldRetry, incrementRetry)

### Безопасность
- **Бэкапы закрыты от world-read** — новые архивы в `backups/` (включая `pre_restore_*`) сразу `600` с владельцем `appuser` (uid 1001:1001) — их читает restore-контейнер; каталог `backups/` — `700`. Stale-lock бэкапа/restore теперь снимается и root'ом, и appuser'ом. На установленных серверах права выставляются автоматически при первом прогоне `backup.sh`.
- **deploy.sh** — копия старого `.env` (`.env.backup.*`) теперь `600`, а не `644` (раньше утечка секретов в world-readable файл).

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

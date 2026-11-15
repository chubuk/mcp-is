# MCP IS — документация и возможности

MCP IS — MCP-сервер для интеграции ИИ-агентов с платформой InSales (товары, заказы, клиенты, блог, SEO) + веб-админка, ИИ-ключи, биллинг и контент-движок PageTree. Этот документ — карта документации: что где описано.

## Где какая документация

| Документ | Содержание |
|---|---|
| `README.md` | Обзор продукта, быстрый старт, лицензия (BSL 1.1) |
| `CHANGELOG.md` | История версий |
| `docs/setup/INSTALL.md` | Установка (Docker, .deb-пакет, переменные окружения) |
| `docs/setup/ARCHITECTURE.md` | Архитектура: сервисы, БД, схема контуров |
| `docs/setup/API.md` | HTTP API сервера: авторизация, ключи, админ-эндпоинты |
| `docs/setup/FAQ.md` | Частые вопросы |
| `docs/setup/INSALES-APP.md` | Установка как приложение InSales (app-mode, автологин) |
| `docs/setup/EMAIL-DNS.md` | Почта и DNS-настройки |
| `docs/api/openapi-mcp-is.yaml` | OpenAPI-спецификация HTTP API |
| `docs/api/insales-api-reference.md` | Справочник REST API InSales |
| `docs/mcp-acceptance-checklist.md` | Чеклист приёмки MCP-инструментов |
| `server.json` | Манифест MCP-сервера для реестров |
| `mcp-pagetree.md`, `api-pagetree.md` | Контент-движок PageTree: MCP и REST API (в репо `chubuk/pagetree`, см. ссылки ниже) |

## Возможности

### MCP-инструменты InSales
- **49 core-тулов** в `tools/list` по умолчанию (товары, заказы, клиенты, коллекции, блог, SEO, доп. поля); полный каталог (~287) — через `search_insales_tools` + `call_insales_tool` (progressive discovery).
- Классы операций: чтение / запись / удаление — политика ключа может запрещать классы целиком.
- Загрузка картинок байтами (`image_data`), идемпотентность create-операций, фильтрация read-only ключей.

### ИИ-ключи и безопасность
- Ключ на привязан к одному магазину (`shop_id`); политика команд (read/write/delete), ротация, лимиты (50 на пользователя).
- Rate limiting `/mcp` — per-ключ; JWT-аутентификация админки; PDN-маскировка персональных данных в выгрузках.

### Биллинг и приложение InSales
- Биллинг-гейт (RecurringApplicationCharge), app-mode установка/автологин, scopes карточки приложения.

### Контент-движок PageTree
- Страницы и новости сайта (`*.md` + sidecar `*.json`), редактор в админке («Страницы»), медиа, SEO, меню.
- Доступ агентов: REST `/api/agent/v1/` и MCP `/api/agent/mcp` (Bearer `ptk_*`), управление ключами — страница «Ключи и доступы».
- Подробно: [REST API PageTree](https://mcp-is.ru/pages/docs/api-pagetree), [MCP API PageTree](https://mcp-is.ru/pages/docs/mcp-pagetree).

### Эксплуатация
- Резервное копирование/восстановление (backup.sh, pg_basebackup), метрики Prometheus `/metrics`, health `/health`.
- Деплой: Docker-образы (app, nginx, db, socket-proxy), ghcr.io; статика веба — volume, не в образе.

## Быстрые ссылки

- Публичный сайт и админка: `/` (SPA), `/content/admin/` (редактор страниц), `/content/admin/keys.html` (ключи PageTree)
- MCP-эндпоинт InSales: `/mcp`; PageTree: `/api/agent/mcp`
- Здоровье сервиса: `GET /health`

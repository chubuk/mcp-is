# MCP IS — Пошаговое руководство

## Что такое MCP IS

MCP IS — это сервер, который подключает AI-агентов (Claude, Cursor, ChatGPT и др.) к вашему интернет-магазину на платформе InSales. Через MCP-протокол AI получает доступ к товарам, заказам, клиентам и другим данным магазина.

## Предварительные требования

- Docker и Docker Compose v2
- Ключ доступа InSales (Расширения → Разработчикам → Ключи доступа)
- AI-клиент: Claude Desktop, Cursor, VS Code/Copilot, Continue, LibreChat или Open WebUI

## Шаг 1: Установка

```bash
# Клонировать проект
git clone https://github.com/chubuk/MCP-IS.git
cd MCP-IS

# Настроить переменные окружения
cp .env.example .env
nano .env  # Заполните INSALES_SHOP, INSALES_API_KEY, INSALES_PASSWORD
```

## Шаг 2: Запуск

```bash
docker compose up -d
```

Проверьте работоспособность:
```bash
curl http://localhost:8080/health
# Ответ: ok
```

## Шаг 3: Первый вход

1. Откройте `https://localhost:8443` (самоподписанный сертификат)
2. Зарегистрируйтесь или войдите (по умолчанию: admin / password)
3. Вы попадёте на панель управления

## Шаг 4: Добавление магазина

1. Перейдите в раздел **Магазины**
2. Нажмите **+ Добавить**
3. Заполните поля:
   - **Название** (обязательно): например "Мой магазин"
   - **Домен**: `myshop.myinsales.ru`
   - **ID ключа**: из панели InSales
   - **Пароль ключа**: из панели InSales
4. Нажмите **Проверить соединение** — убедитесь что всё работает
5. Нажмите **Сохранить**

## Шаг 5: Генерация AI-ключа

1. Перейдите в раздел **AI-ключи**
2. Нажмите **+ Создать ключ**
3. Введите название (например "Cursor MCP")
4. **Скопируйте ключ** — он показывается только один раз!

## Шаг 6: Подключение AI-клиента

### Claude Desktop

Добавьте в `~/.config/claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "insales": {
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <ВАШ_КЛЮЧ>"
      }
    }
  }
}
```

### Cursor

В настройках MCP:
- URL: `http://localhost:8080/mcp`
- Заголовок: `Authorization: Bearer <ВАШ_КЛЮЧ>`

### VS Code / Copilot

Создайте `.vscode/mcp.json`:
```json
{
  "servers": {
    "mcp-is": {
      "type": "http",
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <ВАШ_КЛЮЧ>"
      }
    }
  }
}
```

### LibreChat

В `librechat.yaml`:
```yaml
mcpServers:
  mcp-is:
    type: streamable-http
    url: http://localhost:8080/mcp
    headers:
      Authorization: 'Bearer ${INSALES_MCP_API_KEY}'
```

## Шаг 7: Первые запросы

Попросите AI:
- "Покажи список товаров в моём магазине"
- "Сколько заказов за последний месяц?"
- "Какие категории товаров есть?"
- "Покажи информацию о магазине"

## Шаг 8: Управление через Web UI

| Раздел | Что можно делать |
|--------|------------------|
| Дашборд | Обзор магазинов и ключей |
| Магазины | Добавлять/удалять магазины, проверять команды |
| AI-ключи | Создавать/отзывать ключи |
| Настройки | Смена пароля, email |
| Пользователи | Управление пользователями (admin) |
| Сервер | Домен, HTTPS, сертификаты (super_admin) |

## Полезные ссылки

- [INSTALL.md](INSTALL.md) — подробная установка
- [USAGE.md](USAGE.md) — справка по использованию
- [API.md](API.md) — справка по API
- [ARCHITECTURE.md](ARCHITECTURE.md) — архитектура системы

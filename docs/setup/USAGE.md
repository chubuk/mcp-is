# Использование MCP InSales

## Быстрый старт

### 1. Регистрация

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "myuser",
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

### 2. Вход

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "myuser", "password": "securepassword"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
```

### 3. Добавление магазина

```bash
curl -X POST http://localhost:8080/api/shops \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "shop_url": "myshop.myinsales.ru",
    "api_key_id": "ваш_ключ",
    "api_password": "ваш_пароль",
    "is_default": true
  }'
```

### 4. Генерация AI-ключа

```bash
AI_KEY=$(curl -s -X POST http://localhost:8080/api/ai-keys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Cursor MCP"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['key'])")
echo "Сохраните ключ: $AI_KEY"
```

### 5. MCP-запрос

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AI_KEY" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

## Подключение MCP-клиента

### Claude Desktop

Добавьте в `~/.config/claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "insales": {
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <ВАШ_AI_КЛЮЧ>"
      }
    }
  }
}
```

### Cursor / VS Code (MCP-расширение)

В настройках MCP укажите:
- URL: `http://localhost:8080/mcp`
- Заголовок: `Authorization: Bearer <ВАШ_AI_КЛЮЧ>`

### VS Code / Copilot

Создайте файл `.vscode/mcp.json` в корне проекта:

```json
{
  "servers": {
    "mcp-insales": {
      "type": "http",
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <ВАШ_AI_КЛЮЧ>"
      }
    }
  }
}
```

Или через Command Palette: **MCP: Add Server** → выбрать **Workspace** → ввести URL и заголовок.

### Continue

Создайте файл `.continue/mcpServers/mcp-insales.yaml`:

```yaml
name: MCP InSales
version: 0.1.0
schema: v1
mcpServers:
  - name: InSales
    type: streamable-http
    url: http://localhost:8080/mcp
    env:
      MCP_API_KEY: <ВАШ_AI_КЛЮЧ>
```

Или добавьте в `config.yaml`:

```yaml
mcpServers:
  - name: InSales
    type: streamable-http
    url: http://localhost:8080/mcp
    env:
      MCP_API_KEY: <ВАШ_AI_КЛЮЧ>
```

### LibreChat

Добавьте в `librechat.yaml`:

```yaml
mcpServers:
  mcp-insales:
    type: streamable-http
    url: http://localhost:8080/mcp
    headers:
      Authorization: 'Bearer ${INSALES_MCP_API_KEY}'
    requiresOAuth: false
    serverInstructions: true
    timeout: 30000
```

В `.env` LibreChat добавьте:
```bash
INSALES_MCP_API_KEY=<ваш_AI_ключ_из_mcp-insales>
```

### Open WebUI (v0.6.31+)

**Через Web UI:**
1. Admin Settings → External Tools → + (Add Server)
2. Type: **MCP (Streamable HTTP)**
3. URL: `http://localhost:8080/mcp`
4. Auth: **Bearer** → ввести API-ключ mcp-insales
5. Save → Restart

**Если Open WebUI в Docker:**
- URL: `http://host.docker.internal:8080/mcp`

**Через конфигурацию:**
```yaml
mcpServers:
  mcp-insales:
    type: streamable-http
    url: http://host.docker.internal:8080/mcp
    headers:
      Authorization: 'Bearer ${INSALES_MCP_API_KEY}'
```

## Примеры запросов через curl

### Список товаров

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

Ответ:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [{
      "type": "text",
      "text": "[{\"id\":\"275637484\",\"title\":\"ЛЕВЕТОН П 60...\"}]"
    }]
  }
}
```

### Поиск товара

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 2,
    "params": {
      "name": "list_products",
      "arguments": {"q": "витамин", "per_page": 3}
    }
  }'
```

### Список заказов

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 3,
    "params": {
      "name": "list_orders",
      "arguments": {"page": 1, "per_page": 3}
    }
  }'
```

### Заказы с фильтрацией

```bash
# Только оплаченные
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 4,
    "params": {
      "name": "list_orders",
      "arguments": {"status": "open", "per_page": 10}
    }
  }'
```

### Детали заказа

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 5,
    "params": {
      "name": "get_order",
      "arguments": {"order_id": "1523259105"}
    }
  }'
```

### Список клиентов

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 6,
    "params": {
      "name": "list_clients",
      "arguments": {"page": 1, "per_page": 5}
    }
  }'
```

### Поиск клиента

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 7,
    "params": {
      "name": "list_clients",
      "arguments": {"q": "ivanov"}
    }
  }'
```

### Категории

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 8,
    "params": {
      "name": "list_categories",
      "arguments": {"per_page": 10}
    }
  }'
```

### Коллекции

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 9,
    "params": {
      "name": "list_collections",
      "arguments": {"per_page": 10}
    }
  }'
```

### Валюты магазина

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 10,
    "params": {
      "name": "list_currencies",
      "arguments": {}
    }
  }'
```

### Способы доставки

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 11,
    "params": {
      "name": "list_delivery_variants",
      "arguments": {}
    }
  }'
```

### Способы оплаты

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 12,
    "params": {
      "name": "list_payment_gateways",
      "arguments": {}
    }
  }'
```

## Управление аккаунтом

### Смена пароля

```bash
curl -X POST http://localhost:8080/api/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "old_password": "старый_пароль",
    "new_password": "новый_пароль"
  }'
```

### Восстановление пароля

1. Запросите сброс:
```bash
curl -X POST http://localhost:8080/api/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

2. Используйте токен из письма (или логов, если SMTP не настроен):
```bash
curl -X POST http://localhost:8080/api/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "токен_из_письма",
    "new_password": "новый_пароль"
  }'
```

## Админ-панель

### Через веб-интерфейс

1. Войдите как super_admin на `https://localhost:8443`
2. Перейдите в раздел «Админ»
3. Доступные действия:
   - Просмотр списка пользователей
   - Просмотр деталей пользователя
   - Смена роли (только super_admin)
   - Блокировка/разблокировка
   - Удаление пользователя
   - Статистика системы

### Через API

```bash
# Список пользователей
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/admin/users

# Статистика
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/admin/stats

# Смена роли
curl -X PUT http://localhost:8080/api/admin/users/USER_ID/role \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"role": "admin"}'

# Блокировка
curl -X PUT http://localhost:8080/api/admin/users/USER_ID/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"is_active": false}'
```

### Роли

| Роль | Возможности |
|------|-------------|
| `super_admin` | Полный доступ, управление ролями, защищён от удаления/блокировки |
| `admin` | Управление пользователями (кроме super_admin) |
| `user` | Только собственные данные |

## Примеры с AI-агентом

При подключении через Claude Desktop или Cursor можно задавать вопросы естественным языком:

- "Покажи последние 5 заказов"
- "Сколько товаров в категории Электроника?"
- "Какие способы доставки доступны?"
- "Найди клиента с email ivanov@mail.ru"
- "Какова общая сумма заказа номер 69196?"
- "Покажи товары с остатком менее 10 штук"

## Формат ответа

Все ответы возвращаются в формате JSON-RPC 2.0:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "... (JSON с данными)"
      }
    ]
  }
}
```

При ошибке:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "описание ошибки"
  }
}
```

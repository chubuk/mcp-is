# MCP InSales — Решение проблем

## Docker

### Контейнер не запускается

**Симптом:** `docker compose up` выдаёт ошибку

**Решение:**
```bash
# Проверить Docker
docker --version
docker compose version

# Проверить порты
sudo lsof -i :80
sudo lsof -i :443
sudo lsof -i :5432

# Если порт занят — измените в .env:
# NGINX_HTTP_PORT=8888
# NGINX_HTTPS_PORT=8443
```

### Nginx показывает "unhealthy"

**Симптом:** `docker compose ps` показывает `unhealthy` для nginx

**Причина:** Healthcheck проверяет HTTP на порту 80, но nginx перенаправляет на 443. Это нормально для dev-окружения.

**Решение:** Игнорировать для dev. Для production настройте SSL.

## Авторизация

### Ошибка "invalid API key"

**Причина:** Ключ невалидный или отозван.

**Решение:**
1. Проверьте что ключ скопирован полностью
2. Проверьте что ключ не отозван в AI-ключи
3. Формат заголовка: `Authorization: Bearer <ключ>`

### Ошибка "missing token"

**Причина:** Не передан JWT-токен.

**Решение:**
```bash
# Получить токен
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "admin", "password": "password"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
```

## MCP

### AI не видит инструменты

**Причины:**
1. Сервер не запущен
2. Неверный URL
3. Невалидный API-ключ

**Решение:**
```bash
# Проверить сервер
curl http://localhost:8080/health

# Проверить MCP
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ключ>" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

### AI получает ошибку "shop not found"

**Причина:** Не добавлен магазин или неверный ключ.

**Решение:**
1. Добавьте магазин в Web UI
2. Проверьте что магазин привязан к пользователю
3. Убедитесь что используется правильный AI-ключ

## База данных

### Ошибка "connection refused"

**Решение:**
```bash
# Проверить контейнер БД
docker compose ps db

# Проверить логи
docker compose logs db

# Перезапустить
docker compose restart db
```

### Как сбросить пароль admin

```bash
# Сгенерировать хэш
HASH=$(go run -e - <<'EOF'
package main
import ("fmt"; "golang.org/x/crypto/bcrypt")
func main() {
    h, _ := bcrypt.GenerateFromPassword([]byte("newpassword"), bcrypt.DefaultCost)
    fmt.Print(string(h))
}
EOF
)

# Обновить в БД
docker exec mcp_insales_db psql -U admin -d mcp_insales \
  -c "UPDATE users SET password_hash = '$HASH' WHERE nickname = 'admin';"
```

## HTTPS

### Ошибка TLS "certificate has expired"

**Решение:**
1. Web UI → Сервер → Продлить сертификат
2. Или: `docker exec mcp_insales_nginx certbot renew`

### Let's Encrypt не работает на localhost

**Нормально.** Let's Encrypt требует публичный домен. На localhost используется self-signed сертификат.

## Логи

```bash
# Все логи
docker compose logs -f

# Только ошибки
docker compose logs -f --since 10m app

# Последние 100 строк
docker compose logs --tail=100 app
```

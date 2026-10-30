# MCP IS — Решение проблем

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
  -d '{"login": "admin@mcp-is.local", "password": "Admin123"}' | \
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
docker exec mcp_is_db psql -U admin -d mcp_is \
  -c "UPDATE users SET password_hash = '$HASH' WHERE nickname = 'admin';"
```

## HTTPS

### Ошибка TLS "certificate has expired"

**Решение:**
1. Web UI → Сервер → Продлить сертификат
2. Или: `docker exec mcp_is_nginx certbot renew`

### Let's Encrypt не работает на localhost

**Нормально.** Let's Encrypt требует публичный домен. На localhost используется self-signed сертификат.

### Продление Let's Encrypt «висит» несколько минут

**Нормально.** certbot в неинтерактивном режиме делает случайную паузу до ~5 минут. Для ручного/скриптового продления используйте `--no-random-sleep-on-renew`:

```bash
docker exec mcp_is_nginx certbot renew --no-random-sleep-on-renew --non-interactive
```

## Бэкапы

### Кнопка «Бэкап сейчас» → "backup failed" / exit 127

Скрипту нужен `bash` в app-контейнере (образы до v0.4.5 bash не содержали). Обновите образ до актуального.

### Бэкапы не выполняются по расписанию, «Нагрузка слишком высокая»

Порог `BACKUP_MAX_LOAD_PCT` считается в % от ядер: на 1-ядерном сервере 50% = load 0.5. Если фоновая нагрузка машины выше порога, бэкап будет откладываться каждую ночь. Поднимите порог в админке (например, 150) или посмотрите лог: `/var/log/mcp-is-backup.log`.

### Permission denied при записи в /opt/mcp-is/backups

Каталог должен принадлежать appuser (uid 1001):

```bash
chown -R 1001:1001 /opt/mcp-is/backups
```

### Восстановление из бэкапа

См. INSTALL.md → «Восстановление из бэкапа». Restore выполняется **с хоста**:
`/opt/mcp-is/scripts/backup.sh --restore <файл>` или из админки (раздел «Бэкапы», super_admin).

### Битый бэкап / «Нет доступа в админку» после восстановления

1. Проверьте целостность файла заранее: `backup.sh --verify <файл>` → `readable|corrupt`.
   Повреждённый архив (оборванная копия, битый диск) отклоняется на первом шаге — БД не трогается.
2. Если восстановление уже началось и пошло не так: `cat /opt/mcp-is/backups/restore_status.json`.
   `restored-previous` = автоматический откат выполнен, БД в исходном состоянии.
3. `FATAL` — см. INSTALL.md → «Если восстановление не удалось» (откат вручную из
   `pre_restore_*.sql.gz`, битый PGDATA отодвинут в `data.bad-<timestamp>`).

## Доступ к базе данных

### PostgreSQL не доступен снаружи (port 5432)

**Нормально.** PostgreSQL intentionally не экспортирован на хост для безопасности.

**Доступ к БД:**

```bash
# Через Docker exec
docker compose exec db psql -U admin -d mcp_is

# С паролем
PGASSWORD=$(grep DB_PASSWORD .env | cut -d= -f2) docker compose exec db psql -U admin -d mcp_is
```

### Ошибка "connection refused" при подключении к БД

**Причина:** PostgreSQL контейнер не запущен или не healthy.

**Решение:**
```bash
docker compose ps
docker compose logs db --tail=20
docker compose restart db
```

### Ошибка "password authentication failed"

**Причина:** Неверный пароль в DB_URL или .env.

**Решение:**
```bash
# Проверить пароль в .env
grep DB_PASSWORD .env

# Проверить переменную в контейнере
docker compose exec db env | grep POSTGRES
```

## Логи

```bash
# Все логи
docker compose logs -f

# Только ошибки
docker compose logs -f --since 10m app

# Последние 100 строк
docker compose logs --tail=100 app
```

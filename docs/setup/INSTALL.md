# Установка MCP IS

## Быстрый старт

### 1. Установка Docker

**Linux (Alt Linux / Ubuntu / Debian):**

```bash
# Docker
sudo apt update
sudo apt install docker.io docker-compose-v2

# Добавить пользователя в группу docker
sudo usermod -aG docker $USER

# Перелогиниться для применения
```

**Проверка:**

```bash
docker --version
docker compose version
```

### 2. Клонирование проекта

```bash
git clone <url_repos>
cd mcp-is
```

### 3. Получение ключа InSales

1. Войдите в панель администратора вашего магазина InSales
2. Перейдите: **Расширения → Разработчикам → Ключи доступа**
3. Нажмите **Создать ключ доступа**
4. Скопируйте:
   - **Идентификатор** (ID) — это будет `INSALES_API_KEY`
   - **Пароль** — это будет `INSALES_PASSWORD`
5. Технический домен магазина (виден в адресной строке) — это `INSALES_SHOP`

### 4. Настройка

```bash
cp .env.example .env
```

Отредактируйте `.env`:

```env
# База данных (пароль можно оставить по умолчанию)
DB_PASSWORD=password

# Сервер
SERVER_PORT=8080
LOG_LEVEL=info

# InSales (обязательно заполнить!)
INSALES_SHOP=your-shop.myinsales.ru
INSALES_API_KEY=ваш_идентификатор
INSALES_PASSWORD=ваш_пароль

# Режим доступа
MCP_READ_ONLY=true

# Auth: JWT secret (обязательно сменить в production!)
JWT_SECRET=change-me-to-random-string

# Auth: AES-GCM key для шифрования паролей InSales (32 байта!)
ENCRYPTION_KEY=12345678901234567890123456789012
```

### 5. Запуск

```bash
docker compose up -d
```

Первый запуск занимает 1-2 минуты (скачивание образов).

### 6. Проверка

```bash
# Проверка здоровья сервера
curl http://localhost:8080/health

# Просмотр логов
docker compose logs -f app
```

## Web-интерфейс

После запуска доступен веб-интерфейс:

- **HTTP:** `http://localhost:8888` (перенаправляет на HTTPS)
- **HTTPS:** `https://localhost:8443` (самоподписанный сертификат)

> При первом входе браузер покажет предупреждение о сертификате — это нормально для самоподписанного сертификата.

### Порты

По умолчанию nginx слушает порты 80/443. Если они заняты, измените в `.env`:

```env
NGINX_HTTP_PORT=8888
NGINX_HTTPS_PORT=8443
```

### Навигация

- **Панель** — общая статистика
- **Магазины** — управление подключёнными магазинами InSales
- **AI-ключи** — генерация и управление ключами для AI-агентов
- **Настройки** — профиль, смена пароля
- **Админ** — (только для admin/super_admin) управление пользователями

## SMTP (опционально)

Для отправки писем восстановления пароля настройте SMTP в `.env`:

```env
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=your-password
SMTP_FROM=noreply@example.com
```

Если SMTP не настроен, письма будут только логироваться в консоль сервера.

## Обновление

### Через deb-пакет

```bash
# Скачайте новый deb-пакет
sudo dpkg -i mcp-is_X.Y.Z_amd64.deb
```

Скрипт автоматически:
- Создаст бэкап БД и файлов (с вашим подтверждением)
- Остановит сервис
- Обновит файлы
- Выполнит миграции БД
- Проверит работоспособность
- Запустит сервис

### Через git (для разработчиков)

```bash
cd mcp-is
git pull
docker compose up -d --build
```

## Восстановление из бэкапа

Бэкапы хранятся в `/opt/mcp-is/backups/`:
- `db_YYYYMMDD_HHMMSS.sql` — дамп базы данных
- `files_YYYYMMDD_HHMMSS.tar.gz` — файлы проекта

### Восстановление БД

```bash
# 1. Остановите сервис
sudo systemctl stop mcp-is

# 2. Восстановите БД
sudo docker compose -f /opt/mcp-is/docker-compose.yaml \
  exec -T db psql -U admin -d mcp_is < /opt/mcp-is/backups/db_YYYYMMDD_HHMMSS.sql

# 3. Запустите сервис
sudo systemctl start mcp-is
```

### Восстановление файлов

```bash
# 1. Остановите сервис
sudo systemctl stop mcp-is

# 2. Распакуйте бэкап
sudo tar xzf /opt/mcp-is/backups/files_YYYYMMDD_HHMMSS.tar.gz -C /opt/mcp-is/

# 3. Пересоберите контейнеры
cd /opt/mcp-is && sudo docker compose up -d --build

# 4. Запустите сервис
sudo systemctl start mcp-is
```

### Полное восстановление

```bash
# 1. Остановите сервис
sudo systemctl stop mcp-is

# 2. Восстановите файлы
sudo tar xzf /opt/mcp-is/backups/files_YYYYMMDD_HHMMSS.tar.gz -C /opt/mcp-is/

# 3. Пересоберите и запустите
cd /opt/mcp-is && sudo docker compose up -d --build

# 4. Восстановите БД (если нужно)
sudo docker compose -f /opt/mcp-is/docker-compose.yaml \
  exec -T db psql -U admin -d mcp_is < /opt/mcp-is/backups/db_YYYYMMDD_HHMMSS.sql

# 5. Запустите сервис
sudo systemctl start mcp-is
```

## Остановка

```bash
docker compose down
```

## Удаление данных БД

```bash
docker compose down -v
```

## Частые проблемы

### Ошибка "connection refused"

Контейнер ещё не запустился. Подождите 5-10 секунд и повторите запрос.

### Ошибка "insales credentials not configured"

Не заполнены `INSALES_SHOP` и `INSALES_API_KEY` в `.env`. Проверьте файл:

```bash
cat .env | grep INSALES
```

### Ошибка "API error 401"

Неверный ключ доступа. Проверьте `INSALES_API_KEY` и `INSALES_PASSWORD` в `.env`.

### Предупреждение "INSALES_TOKEN variable is not set"

Можно игнорировать. Эта переменная больше не используется.

---

## Автоматический деплой на VDS

### Быстрый старт (чистый сервер)

```bash
# Подключиться к VDS
ssh root@YOUR_SERVER_IP

# Клонировать и установить
cd /opt
git clone https://github.com/chubuk/mcp-is.git
cd mcp-is
cp .env.example .env
nano .env  # Заполнить настройки

# Запустить
docker compose up -d

# Проверить
curl -k https://YOUR_DOMAIN/health
```

### Использование deploy.sh

Скрипт `scripts/deploy.sh` автоматизирует все операции:

```bash
# Проверка системы (pre-flight)
./scripts/deploy.sh check

# Полная установка на чистый сервер
./scripts/deploy.sh install

# Обновление существующей установки
./scripts/deploy.sh update

# Ручной бэкап
./scripts/deploy.sh backup

# Проверка статуса
./scripts/deploy.sh status
```

### Автоматические бэкапы

Бэкапы настраиваются через админку (`/admin/backup`) или `.env`:

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `BACKUP_ENABLED` | Включить бэкапы | true |
| `BACKUP_METHOD` | pg_basebackup / pg_dump | pg_basebackup |
| `BACKUP_INTERVAL_HOURS` | Интервал (часов) | 24 |
| `BACKUP_RETENTION_DAYS` | Хранение (дней) | 14 |

Cron настроен автоматически при установке через deploy.sh.

### Безопасность

- `.env` файл заблокирован в nginx (403 Forbidden)
- Скрытые файлы (`.git`, `.env`) недоступны через веб
- Доступ к админке только для admin/super_admin ролей

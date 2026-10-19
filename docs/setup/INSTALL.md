# Установка MCP IS

> **Лицензия (BSL 1.1).** Коммерческое использование проекта (в том числе отдельный платный сервис или платное приложение InSales) **не разрешено без согласия разработчика** — Vladislav Chubuk &lt;chubuk@mail.ru&gt;. Личное, учебное и тестовое использование — свободно. Подробнее: `README.md`, `LICENSE`.

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

### 3. Ключи InSales

Ключи доступа к магазину добавляются **в веб-интерфейсе** (раздел «Магазины») после первого входа: идентификатор + пароль из **Расширения → Разработчикам → Ключи доступа**.

Если планируете режим **стороннего приложения InSales** (установка из бэк-офиса магазина), дополнительно понадобятся идентификатор и секрет самого приложения — см. `docs/setup/INSALES-APP.md` и `docs/insales-app-form.md`.

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

# InSales — стороннее приложение (необязательно, для установки из ЛК InSales)
INSALES_APP_KEY=
INSALES_SECRET_KEY=

# Регистрация: false = только админ создаёт пользователей.
# Для отдельного сервиса, где пользователи сами регистрируются на сайте, поставьте true.
REGISTRATION_ENABLED=false

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
- **Сервер** — (admin/super_admin) домен, HTTPS, сертификаты
- **ПДн** — (admin/super_admin) защита персональных данных
- **Бэкапы** — (admin/super_admin) резервные копии
- **Почта** — (admin/super_admin) SMTP, шаблоны писем, тест-письмо

## SMTP (опционально)

Для отправки писем (сброс пароля, уведомления) настройте SMTP в `.env` или в админке **Почта** (`/admin/email`) — изменения вступают в силу без рестарта:

```env
SMTP_HOST=smtp.mail.ru
SMTP_PORT=465          # 465 = SSL (mail.ru), 587 = STARTTLS
SMTP_USER=user@mail.ru
SMTP_PASS=пароль-приложений   # НЕ основной пароль ящика
SMTP_FROM=MCP IS <user@mail.ru>
```

**Mail.ru:** обычный пароль не работает — создайте [пароль приложений](https://account.mail.ru/user/2-step-auth/passwords/) (Настройки → Безопасность → Пароли для внешних приложений).  
Если SMTP timeout — на VPS могут быть закрыты порты 25/465/587 (поддержка хостера).

Если SMTP не настроен, письма будут только логироваться в консоль сервера. При запросе сброса пароля URL восстановления будет выведен в логи — администратор может передать его пользователю вручную.

### DNS-записи для доставки (SPF / DKIM / DMARC)

Чтобы письма не попадали в спам, в DNS домена должны быть:

| Тип | Имя | Значение |
|-----|-----|----------|
| TXT | `@` | `v=spf1 include:_spf.mail.ru ~all` |
| TXT | `_dmarc` | `v=DMARC1; p=none; rua=mailto:dmarc@ваш-домен` |
| TXT | `<selector>._domainkey` | `v=DKIM1; k=rsa; p=…` (из панели почты) |

**Вариант без DNS:** `SMTP_FROM` с адресом провайдера (`…@mail.ru`) — письма дойдут, но «от кого» не ваш домен.

Подробнее: `docs/setup/EMAIL-DNS.md`. Проверка: `./scripts/check-email-dns.sh mcp-is.ru`.  
Настройки SMTP и шаблоны писем: **Почта** (`/admin/email`, admin/super_admin).

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
- `db_YYYYMMDD_HHMMSS.tar.gz` — файловый снэпшот pg_basebackup (метод по умолчанию)
- `db_YYYYMMDD_HHMMSS.sql.gz` — SQL-дамп pg_dump (если в настройках выбран pg_dump)

Выполняйте команды **на хосте от root** (не из app-контейнера — restore останавливает контейнеры).

### Восстановление из снэпшота (tar.gz)

```bash
# Скрипт остановит app/nginx, распакует снэпшот в БД, запустит всё обратно
sudo /opt/mcp-is/scripts/backup.sh --restore /opt/mcp-is/backups/db_YYYYMMDD_HHMMSS.tar.gz
```

### Восстановление из SQL-дампа (sql.gz)

```bash
# Дамп применяется поверх работающей БД (без остановки)
gunzip -c /opt/mcp-is/backups/db_YYYYMMDD_HHMMSS.sql.gz | \
  docker exec -i mcp_is_db psql -U admin -d mcp_is
```

### Файлы проекта

Отдельного «файлового» бэкапа нет: код восстанавливается через `git pull` (или повторную установку), затем `docker compose up -d`. Не забудьте про `.env` (пароли БД и ключи) — храните его копию отдельно.

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

У магазина не заполнены ключи доступа. Добавьте магазин в веб-интерфейсе (раздел «Магазины»): идентификатор и пароль из **Расширения → Разработчикам → Ключи доступа**.

### Ошибка "API error 401"

Неверный ключ доступа магазина. Проверьте идентификатор и пароль в карточке магазина (веб-интерфейс).

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

Настройки хранятся в БД и редактируются в админке (Настройки сервера → Бэкапы); значения из `.env` используются как fallback:

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `BACKUP_ENABLED` | Включить бэкапы | true |
| `BACKUP_METHOD` | pg_basebackup / pg_dump | pg_basebackup |
| `BACKUP_INTERVAL_HOURS` | Интервал (часов) | 24 |
| `BACKUP_HOUR_START` / `BACKUP_HOUR_END` | Окно запуска (часы) | 1–6 |
| `BACKUP_MAX_LOAD_PCT` | Макс. нагрузка, % от ядер | 50 |
| `BACKUP_RETENTION_DAYS` | Хранение (дней) | 14 |

`BACKUP_RETENTION_DAYS` действует и на файлы в `/opt/mcp-is/backups/`, и на записи журнала (`backup_log`), и на инциденты (`backup_incidents`) — всё старше срока удаляется при каждом запуске `backup.sh`, даже если ночной бэкап пропущен из-за нагрузки.

Журнал и инциденты в админке листаются постранично (кнопки «Назад/Далее»).

Расписание: cron на хосте (`/etc/cron.d/mcp-is-backup`) запускает `backup.sh --scheduled` каждый час; бэкап выполняется в пределах окна, если прошёл интервал и нагрузка ниже порога. Порог считается в % от ядра (на 1-ядерном VDS 50% = load 0.5; если машина почти не грузится в среднем выше — поднимите порог, иначе бэкапы будут откладываться).

Лог `/var/log/mcp-is-backup.log` крутится через logrotate (`/etc/logrotate.d/mcp-is-backup`, weekly ×4). Docker-логи контейнеров ограничены `10m × 5` в `docker-compose.yaml`.

Если при запуске из админки («Бэкап сейчас») ошибка прав на запись в `/opt/mcp-is/backups`:

```bash
chown -R 1001:1001 /opt/mcp-is/backups   # appuser контейнера app
```

### Безопасность

#### Пароли

При установке через `deploy.sh install` пароли генерируются автоматически и сохраняются в:

```bash
/root/.mcp-is-credentials  # chmod 600
```

**НИКОГДА не используйте дефолтные пароли в production!**

#### Docker-сеть

- PostgreSQL доступен только из контейнера `app` (internal сеть)
- Порт 5432 НЕ экспортирован на хост
- Доступ к БД: `docker compose exec db psql -U admin -d mcp_is`

#### Nginx

- `.env`, `.git`, `.sql`, `.log` — заблокированы (403)
- `/metrics` — заблокирован (Prometheus не доступен снаружи)
- Rate limiting: login (5 req/min), API (10 req/s), reset-password (3 req/min), MCP (30 req/min)
- Security headers: HSTS, CSP, X-Frame-Options, Permissions-Policy

#### API

- JWT аутентификация для всех `/api/` эндпоинтов
- API-ключи для `/mcp` эндпоинта
- Role-based access control (user, admin, super_admin)
- Rate limiting на auth-эндпоинтах

#!/usr/bin/env bash
# backup.sh — Бэкап PostgreSQL для MCP-IS
#
# Поддерживает pg_basebackup (файловый снэпшот) и pg_dump (SQL-дамп).
# Мониторинг нагрузки, удалённые хранилища, инциденты.
#
# Usage:
#   ./scripts/backup.sh                   # ручной бэкап
#   ./scripts/backup.sh --scheduled       # запуск по расписанию (cron)
#   ./scripts/backup.sh --restore <file>  # восстановление из бэкапа (см. план GATE!)
#   ./scripts/backup.sh --verify <file>   # проверка целостности (readable|corrupt)
#   ./scripts/backup.sh --status          # проверка статуса
#
# Environment:
#   BACKUP_DIR               Директория бэкапов (default: /opt/mcp-is/backups)
#   BACKUP_METHOD            pg_basebackup | pg_dump (default: pg_basebackup)
#   BACKUP_INTERVAL_HOURS    Интервал часов (default: 12)
#   BACKUP_RETENTION_DAYS    Хранение дней (default: 14)
#   BACKUP_HOUR_START        Начало окна (default: 1)
#   BACKUP_HOUR_END          Конец окна (default: 6)
#   BACKUP_MAX_LOAD_PCT      Макс. CPU load % (default: 50)
#   BACKUP_MAX_IOWAIT_PCT    Макс. I/O wait % (default: 15)
#   BACKUP_REMOTE_ENABLED    Удалённый бэкап (default: false)
#   BACKUP_REMOTE_TYPE       rsync | s3 | rclone (default: rsync)
#   BACKUP_REMOTE_HOST       SSH host
#   BACKUP_REMOTE_PORT       SSH port (default: 22)
#   BACKUP_REMOTE_USER       SSH user
#   BACKUP_REMOTE_PATH       Путь на удалённом сервере
#   BACKUP_REMOTE_KEY        SSH ключ (default: ~/.ssh/id_ed25519)

set -euo pipefail

# ============================================================================
# Цвета и форматирование
# ============================================================================

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    CYAN='\033[0;36m'
    BOLD='\033[1m'
    NC='\033[0m'
else
    RED='' GREEN='' YELLOW='' BLUE='' CYAN='' BOLD='' NC=''
fi

ok()    { echo -e "${GREEN}[OK]${NC}   $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
fail()  { echo -e "${RED}[FAIL]${NC} $1"; }
info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
step()  { echo -e "\n${BOLD}==>${NC} ${CYAN}$1${NC}"; }

# ============================================================================
# Конфигурация
# ============================================================================

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
readonly BACKUP_DIR="${BACKUP_DIR:-/opt/mcp-is/backups}"
readonly BACKUP_PID_FILE="/tmp/mcp-is-backup.pid"
# Общий lock backup+restore живёт в BACKUP_DIR (общий том), НЕ в /tmp:
# PID-неймспейсы контейнеров ломают kill -0, а mkdir атомарен через хост-ФС.
readonly BACKUP_LOCK_FILE="${BACKUP_DIR}/.backup.lock"
readonly DB_CONTAINER="${DB_CONTAINER:-mcp_is_db}"
readonly APP_CONTAINER="${APP_CONTAINER:-mcp_is_app}"
readonly NGINX_CONTAINER="${NGINX_CONTAINER:-mcp_is_nginx}"
readonly DB_NAME="${DB_NAME:-mcp_is}"
readonly DB_USER="${DB_USER:-admin}"
readonly RESTORE_TMP="${BACKUP_DIR}/.restore-tmp"
readonly RESTORE_STATUS_FILE="${BACKUP_DIR}/restore_status.json"
readonly ONESHOT_PG_IMAGE="postgres:16-alpine"

# Значения по умолчанию
BACKUP_METHOD="${BACKUP_METHOD:-pg_basebackup}"
BACKUP_INTERVAL_HOURS="${BACKUP_INTERVAL_HOURS:-24}"
BACKUP_RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-14}"
BACKUP_HOUR_START="${BACKUP_HOUR_START:-1}"
BACKUP_HOUR_END="${BACKUP_HOUR_END:-6}"
BACKUP_MAX_LOAD_PCT="${BACKUP_MAX_LOAD_PCT:-50}"
BACKUP_MAX_IOWAIT_PCT="${BACKUP_MAX_IOWAIT_PCT:-15}"
BACKUP_REMOTE_ENABLED="${BACKUP_REMOTE_ENABLED:-false}"
BACKUP_REMOTE_TYPE="${BACKUP_REMOTE_TYPE:-rsync}"
BACKUP_REMOTE_HOST="${BACKUP_REMOTE_HOST:-}"
BACKUP_REMOTE_PORT="${BACKUP_REMOTE_PORT:-22}"
BACKUP_REMOTE_USER="${BACKUP_REMOTE_USER:-}"
BACKUP_REMOTE_PATH="${BACKUP_REMOTE_PATH:-/backups/mcp-is}"
BACKUP_REMOTE_KEY="${BACKUP_REMOTE_KEY:-~/.ssh/id_ed25519}"

# ============================================================================
# Загрузка конфигурации (БД → .env → defaults)
# ============================================================================

load_backup_config() {
    local key value
    
    # Все параметры для загрузки из БД
    local -A db_keys=(
        ["backup_method"]="$BACKUP_METHOD"
        ["backup_interval_hours"]="$BACKUP_INTERVAL_HOURS"
        ["backup_retention_days"]="$BACKUP_RETENTION_DAYS"
        ["backup_hour_start"]="$BACKUP_HOUR_START"
        ["backup_hour_end"]="$BACKUP_HOUR_END"
        ["backup_max_load_pct"]="$BACKUP_MAX_LOAD_PCT"
        ["backup_max_iowait_pct"]="$BACKUP_MAX_IOWAIT_PCT"
        ["backup_dir"]="$BACKUP_DIR"
        ["backup_remote_enabled"]="$BACKUP_REMOTE_ENABLED"
        ["backup_remote_type"]="$BACKUP_REMOTE_TYPE"
        ["backup_remote_host"]="$BACKUP_REMOTE_HOST"
        ["backup_remote_port"]="$BACKUP_REMOTE_PORT"
        ["backup_remote_user"]="$BACKUP_REMOTE_USER"
        ["backup_remote_path"]="$BACKUP_REMOTE_PATH"
        ["backup_remote_key"]="$BACKUP_REMOTE_KEY"
    )
    
    local config_source="defaults"
    
    for key in "${!db_keys[@]}"; do
        # 1. Пробуем из БД
        value=$(docker exec -i "$DB_CONTAINER" \
            psql -U "$DB_USER" -d "$DB_NAME" -t -A \
            -c "SELECT value FROM server_settings WHERE key = '$key'" 2>/dev/null || true)
        
        if [ -n "$value" ] && [ "$value" != "" ]; then
            export "$(echo "${key^^}" | tr '-' '_' | sed 's/^BACKUP_/BACKUP_/')=$value"
            config_source="db"
            continue
        fi
        
        # 2. Проверяем .env
        local env_value=$(grep -E "^${key^^}=" "${PROJECT_DIR}/.env" 2>/dev/null | cut -d'=' -f2 || true)
        if [ -n "$env_value" ]; then
            export "$(echo "${key^^}" | tr '-' '_' | sed 's/^BACKUP_/BACKUP_/')=$env_value"
            config_source="env"
        fi
        # 3. Используем default (уже установлен)
    done
    
    info "Конфигурация загружена (источник: $config_source)"
}

# ============================================================================
# Мониторинг нагрузки
# ============================================================================

check_load() {
    local cores
    cores=$(nproc)
    
    # Чтение load average
    local one five fifteen
    read -r one five fifteen _ _ < /proc/loadavg
    # % загрузки относительно числа ядер: на 1 ядре 50% = load 0.5
    # (старая формула cores*pct/100 давала 0 на 1 ядре → условие всегда true)
    local load_pct
    load_pct=$(awk -v l="$five" -v c="$cores" 'BEGIN{printf "%.0f", l*100/c}')

    info "Текущая нагрузка: 1min=$one 5min=$five (~${load_pct}% of ${cores} cores, порог=${BACKUP_MAX_LOAD_PCT}%)"

    if [ "$load_pct" -ge "$BACKUP_MAX_LOAD_PCT" ]; then
        warn "Нагрузка слишком высокая: 5min=${load_pct}% >= ${BACKUP_MAX_LOAD_PCT}%"
        return 1
    fi
    
    # Проверка I/O wait
    local iowait
    iowait=$(awk '/^cpu /{printf "%.0f", $6/($2+$3+$4+$5+$6+$7+$8)*100}' /proc/stat)
    info "I/O wait: ${iowait}% (порог=${BACKUP_MAX_IOWAIT_PCT}%)"
    
    if [ "$iowait" -ge "$BACKUP_MAX_IOWAIT_PCT" ]; then
        warn "I/O wait слишком высокий: ${iowait}% >= ${BACKUP_MAX_IOWAIT_PCT}%"
        return 1
    fi
    
    return 0
}

# ============================================================================
# Проверка временного окна
# ============================================================================

check_backup_window() {
    local hour
    hour=$(date +%H)
    
    if [ "$hour" -ge "$BACKUP_HOUR_START" ] && [ "$hour" -lt "$BACKUP_HOUR_END" ]; then
        return 0  # Внутри окна
    fi
    
    return 1  # Вне окна
}

# ============================================================================
# Проверка интервала
# ============================================================================

check_backup_interval() {
    local last_backup
    last_backup=$(ls -t "$BACKUP_DIR"/db_*.tar.gz "$BACKUP_DIR"/db_*.sql.gz 2>/dev/null | head -1 || true)
    
    if [ -z "$last_backup" ]; then
        info "Бэкапов не найдено — можно делать"
        return 0
    fi
    
    local last_time
    last_time=$(stat -c %Y "$last_backup" 2>/dev/null || echo 0)
    local now
    now=$(date +%s)
    local elapsed=$((now - last_time))
    local interval_sec=$((BACKUP_INTERVAL_HOURS * 3600))
    
    if [ "$elapsed" -ge "$interval_sec" ]; then
        local hours_ago=$((elapsed / 3600))
        info "Последний бэкап: ${hours_ago}ч назад (интервал: ${BACKUP_INTERVAL_HOURS}ч)"
        return 0
    fi
    
    local remaining=$(( (interval_sec - elapsed) / 3600 ))
    info "Следующий бэкап через: ${remaining}ч"
    return 1
}

# ============================================================================
# Обработка "застрявших" бэкапов
# ============================================================================

check_stuck_backup() {
    if [ ! -f "$BACKUP_PID_FILE" ]; then
        return 0
    fi
    
    local old_pid
    old_pid=$(cat "$BACKUP_PID_FILE")
    local old_time
    old_time=$(stat -c %Y "$BACKUP_PID_FILE" 2>/dev/null || echo 0)
    local now
    now=$(date +%s)
    local age=$((now - old_time))
    local max_age=$((BACKUP_INTERVAL_HOURS * 3600))
    
    # Если PID жив И бэкап старше интервала — это застрявший бэкап
    if kill -0 "$old_pid" 2>/dev/null && [ "$age" -gt "$max_age" ]; then
        warn "Обнаружен застрявший бэкап (PID=$old_pid, возраст=${age}с)"
        
        # Остановить застрявший процесс
        kill -TERM "$old_pid" 2>/dev/null || true
        sleep 5
        kill -9 "$old_pid" 2>/dev/null || true
        
        # Удалить частичные данные
        cleanup_partial_backup
        
        # Записать инцидент в БД
        log_incident "stuck_backup" "Застрявший бэкап PID=$old_pid остановлен. Частичные данные удалены." "warning"
        
        ok "Застрявший бэкап остановлен, частичные данные очищены"
        rm -f "$BACKUP_PID_FILE"
        rm -rf "$BACKUP_LOCK_FILE"
    fi
    
    return 0
}

cleanup_partial_backup() {
    # Удалить неполные .tar.gz файлы (меньше 1KB и старше 5 минут)
    # pre_restore_* исключены: это страховочные дампы, они могут быть мелкими
    find "$BACKUP_DIR" -name "*.tar.gz" -size -1k -mmin +5 ! -name "pre_restore_*" -delete 2>/dev/null || true
    find "$BACKUP_DIR" -name "*.sql.gz" -size -1k -mmin +5 ! -name "pre_restore_*" -delete 2>/dev/null || true
    rm -rf "$BACKUP_LOCK_FILE"
}

# ============================================================================
# Блокировка запуска нового бэкапа / восстановления (общая)
# ============================================================================

# acquire_shared_lock: mkdir-атомарный lock в общем томе BACKUP_DIR.
# Используется и бэкапом, и restore — взаимное исключение обязательно.
acquire_shared_lock() {
    if ! mkdir "$BACKUP_LOCK_FILE" 2>/dev/null; then
        # Lock существует: stale старше 2 часов — снимаем, иначе отказ
        local lock_age
        lock_age=$(( $(date +%s) - $(stat -c %Y "$BACKUP_LOCK_FILE" 2>/dev/null || echo 0) ))
        if [ "$lock_age" -gt 7200 ]; then
            warn "Снимаю stale-lock (возраст ${lock_age}с)"
            rm -rf "$BACKUP_LOCK_FILE"
            mkdir "$BACKUP_LOCK_FILE" 2>/dev/null || { warn "Бэкап/восстановление уже выполняется"; return 1; }
        else
            warn "Бэкап или восстановление уже выполняется (lock возраст ${lock_age}с)"
            return 1
        fi
    fi
    echo "$(hostname) $$ $(date +%s)" > "$BACKUP_LOCK_FILE/owner"
    trap 'rm -rf "$BACKUP_LOCK_FILE" "$BACKUP_PID_FILE"' EXIT
}

acquire_backup_lock() {
    acquire_shared_lock || return 1
    echo $$ > "$BACKUP_PID_FILE"
}

# ============================================================================
# Бэкап pg_basebackup
# ============================================================================

write_sha256() {
    local file="$1"
    if command -v sha256sum >/dev/null 2>&1; then
        (cd "$(dirname "$file")" && sha256sum "$(basename "$file")" > "$(basename "$file").sha256") 2>/dev/null || true
    fi
}

do_backup_pg_basebackup() {
    local timestamp
    timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="${BACKUP_DIR}/db_${timestamp}.tar.gz"
    
    info "Создание файлового снэпшота через pg_basebackup..."
    local start_time
    start_time=$(date +%s)
    
    # pg_basebackup: tar format + gzip, без остановки сервера
    # -Ft = tar format, -z = gzip, -Xs = stream WAL, -D = output dir
    # очистка на случай, если прошлая попытка оборвалась после создания /tmp/pgbackup
    docker exec "$DB_CONTAINER" rm -rf /tmp/pgbackup
    if docker exec -i "$DB_CONTAINER" \
        pg_basebackup -U "$DB_USER" -Ft -z -Xs -D /tmp/pgbackup 2>&1; then
        
        # Скопировать из контейнера
        docker cp "$DB_CONTAINER:/tmp/pgbackup/" "/tmp/pgbackup/"
        tar -czf "$backup_file" -C /tmp/pgbackup .
        rm -rf /tmp/pgbackup/
        docker exec -i "$DB_CONTAINER" rm -rf /tmp/pgbackup
        write_sha256 "$backup_file"
        
        local end_time
        end_time=$(date +%s)
        local duration_ms=$(( (end_time - start_time) * 1000 ))
        local size
        size=$(du -h "$backup_file" | cut -f1)
        
        ok "Снэпшот создан: $backup_file ($size, ${duration_ms}мс)"
        log_backup_entry "ok" "$size" "$backup_file" "$duration_ms"
        
        # Удалённый бэкап
        if [ "$BACKUP_REMOTE_ENABLED" = "true" ]; then
            sync_to_remote "$backup_file"
        fi
        
        return 0
    else
        local end_time
        end_time=$(date +%s)
        local duration_ms=$(( (end_time - start_time) * 1000 ))
        
        fail "Ошибка pg_basebackup"
        log_backup_entry "error" "0" "pg_basebackup failed" "$duration_ms"
        log_incident "backup_error" "pg_basebackup завершился с ошибкой" "critical"
        return 1
    fi
}

# ============================================================================
# Бэкап pg_dump
# ============================================================================

do_backup_pg_dump() {
    local timestamp
    timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="${BACKUP_DIR}/db_${timestamp}.sql.gz"
    
    info "Создание SQL-дампа через pg_dump..."
    local start_time
    start_time=$(date +%s)
    
    if docker exec -i "$DB_CONTAINER" \
        pg_dump -U "$DB_USER" "$DB_NAME" 2>/dev/null | gzip > "$backup_file"; then
        write_sha256 "$backup_file"
        
        local end_time
        end_time=$(date +%s)
        local duration_ms=$(( (end_time - start_time) * 1000 ))
        local size
        size=$(du -h "$backup_file" | cut -f1)
        
        ok "Дамп создан: $backup_file ($size, ${duration_ms}мс)"
        log_backup_entry "ok" "$size" "$backup_file" "$duration_ms"
        
        # Удалённый бэкап
        if [ "$BACKUP_REMOTE_ENABLED" = "true" ]; then
            sync_to_remote "$backup_file"
        fi
        
        return 0
    else
        local end_time
        end_time=$(date +%s)
        local duration_ms=$(( (end_time - start_time) * 1000 ))
        
        fail "Ошибка pg_dump"
        log_backup_entry "error" "0" "pg_dump failed" "$duration_ms"
        log_incident "backup_error" "pg_dump завершился с ошибкой" "critical"
        return 1
    fi
}

# ============================================================================
# Удалённый бэкап
# ============================================================================

sync_to_remote() {
    local file="$1"
    local filename
    filename=$(basename "$file")
    
    info "Синхронизация с удалённым хранилищем ($BACKUP_REMOTE_TYPE)..."
    
    case "$BACKUP_REMOTE_TYPE" in
        rsync)
            if [ -z "$BACKUP_REMOTE_HOST" ] || [ -z "$BACKUP_REMOTE_USER" ]; then
                warn "Удалённое хранилище не настроено (HOST/USER)"
                return 1
            fi
            
            rsync -avz -e "ssh -i $BACKUP_REMOTE_KEY -p $BACKUP_REMOTE_PORT" \
                "$file" "${BACKUP_REMOTE_USER}@${BACKUP_REMOTE_HOST}:${BACKUP_REMOTE_PATH}/" 2>&1 && \
                ok "Синхронизировано: $filename → $BACKUP_REMOTE_HOST" || \
                warn "Ошибка синхронизации с $BACKUP_REMOTE_HOST"
            ;;
            
        s3)
            if command -v aws &>/dev/null; then
                aws s3 cp "$file" "s3://${BACKUP_REMOTE_PATH}/$filename" 2>&1 && \
                    ok "Загружено в S3: $filename" || \
                    warn "Ошибка загрузки в S3"
            elif command -v mc &>/dev/null; then
                mc cp "$file" "${BACKUP_REMOTE_PATH}/$filename" 2>&1 && \
                    ok "Загружено в MinIO: $filename" || \
                    warn "Ошибка загрузки в MinIO"
            else
                warn "aws или mc не установлены"
                return 1
            fi
            ;;
            
        rclone)
            if command -v rclone &>/dev/null; then
                rclone copy "$file" "${BACKUP_REMOTE_PATH}/" 2>&1 && \
                    ok "Загружено через rclone: $filename" || \
                    warn "Ошибка загрузки через rclone"
            else
                warn "rclone не установлен"
                return 1
            fi
            ;;
            
        *)
            warn "Неизвестный тип удалённого хранилища: $BACKUP_REMOTE_TYPE"
            return 1
            ;;
    esac
}

# ============================================================================
# Логирование в БД
# ============================================================================

log_backup_entry() {
    local status="$1"
    local size="$2"
    local file="$3"
    local duration_ms="${4:-0}"
    
    local one five fifteen iowait
    read -r one five fifteen _ _ < /proc/loadavg
    iowait=$(awk '/^cpu /{printf "%.2f", $6/($2+$3+$4+$5+$6+$7+$8)*100}' /proc/stat)
    
    local size_bytes=0
    if [ -f "$file" ]; then
        size_bytes=$(stat -c %s "$file" 2>/dev/null || echo 0)
    fi
    
    docker exec -i "$DB_CONTAINER" \
        psql -U "$DB_USER" -d "$DB_NAME" -c \
        "INSERT INTO backup_log (ts, status, size_bytes, error_msg, load_1m, load_5m, load_15m, iowait_pct, duration_ms)
         VALUES (NOW(), '$status', $size_bytes, '', $one, $five, $fifteen, $iowait, $duration_ms)" 2>/dev/null || true
}

log_incident() {
    local type="$1"
    local message="$2"
    local severity="$3"
    
    docker exec -i "$DB_CONTAINER" \
        psql -U "$DB_USER" -d "$DB_NAME" -c \
        "INSERT INTO backup_incidents (ts, type, message, severity)
         VALUES (NOW(), '$type', '$message', '$severity')" 2>/dev/null || true
}

# ============================================================================
# Очистка старых бэкапов
# ============================================================================

cleanup_old_backups() {
    local count_before
    count_before=$(find "$BACKUP_DIR" -name "db_*" -type f | wc -l)
    
    find "$BACKUP_DIR" -name "db_*.tar.gz" -mtime +$BACKUP_RETENTION_DAYS -delete 2>/dev/null || true
    find "$BACKUP_DIR" -name "db_*.sql.gz" -mtime +$BACKUP_RETENTION_DAYS -delete 2>/dev/null || true
    # sha256-сайдкары сиротеют без своих архивов — чистим по той же ретенции
    find "$BACKUP_DIR" -name "db_*.sha256" -mtime +$BACKUP_RETENTION_DAYS -delete 2>/dev/null || true
    
    local count_after
    count_after=$(find "$BACKUP_DIR" -name "db_*" -type f | wc -l)
    local deleted=$((count_before - count_after))
    
    if [ "$deleted" -gt 0 ]; then
        info "Удалено старых бэкапов: $deleted (хранение: $BACKUP_RETENTION_DAYS дней)"
    fi
    
    ok "Бэкапов хранится: $count_after"

    # Чистка журнала и инцидентов за пределами retention (интерполяция в SQL — только число)
    local days="$BACKUP_RETENTION_DAYS"
    [[ "$days" =~ ^[0-9]+$ ]] || days=14
    docker exec -i "$DB_CONTAINER" \
        psql -U "$DB_USER" -d "$DB_NAME" -c \
        "DELETE FROM backup_log WHERE ts < NOW() - INTERVAL '$days days';
         DELETE FROM backup_incidents WHERE ts < NOW() - INTERVAL '$days days';" 2>/dev/null || true
}

# ============================================================================
# Уведомления о пропусках
# ============================================================================

check_missed_backups() {
    local last_backup
    last_backup=$(ls -t "$BACKUP_DIR"/db_*.tar.gz "$BACKUP_DIR"/db_*.sql.gz 2>/dev/null | head -1 || true)
    
    if [ -z "$last_backup" ]; then
        log_incident "missed_backup" "Бэкапов нет вообще!" "critical"
        warn "КРИТИЧНО: Бэкапов нет вообще!"
        return
    fi
    
    local last_time
    last_time=$(stat -c %Y "$last_backup" 2>/dev/null || echo 0)
    local now
    now=$(date +%s)
    local elapsed=$((now - last_time))
    local max_elapsed=$((BACKUP_INTERVAL_HOURS * 3600 * 2))
    
    if [ "$elapsed" -gt "$max_elapsed" ]; then
        local hours_ago=$((elapsed / 3600))
        log_incident "missed_backup" "Последний бэкап старше ${hours_ago}ч (норма: ${BACKUP_INTERVAL_HOURS}ч)" "critical"
        warn "ПРЕДУПРЕЖДЕНИЕ: Последний бэкап старше ${hours_ago}ч"
    fi
}

# ============================================================================
# Восстановление из бэкапа (steps 0-6, план eager-canyon rev.4)
# ============================================================================

# restore_status пишет машиночитаемый статус для UI-поллинга.
restore_status() {
    local state="$1" phase="$2" error="${3:-}"
    local tmp="${RESTORE_STATUS_FILE}.tmp"
    cat > "$tmp" <<EOF
{"state":"$state","file":"$(basename "${RESTORE_FILE:-}")","phase":"$phase","started_at":"${RESTORE_STARTED_AT:-}","finished_at":"$(date -u +%Y-%m-%dT%H:%M:%SZ)","error":"$error"}
EOF
    mv "$tmp" "$RESTORE_STATUS_FILE"
}

# step0_verify: проверка целостности БЕЗ побочек в БД.
# 0 = ok; 1 = повреждён/нечитаем. Поддерживает обе раскладки тара
# (корневая `tar -C dir .` и вложенная `pg_basebackup/`) — формат писался
# разными версиями скрипта, старые бэкапы уже лежат на диске.
step0_verify() {
    local backup_file="$1"

    [ -f "$backup_file" ] || { fail "Файл не найден: $backup_file"; return 1; }

    # sha256-сайдкар (если есть) — быстрый признак порчи
    if [ -f "$backup_file.sha256" ]; then
        if ! (cd "$(dirname "$backup_file")" && sha256sum -c "$(basename "$backup_file").sha256" >/dev/null 2>&1); then
            fail "sha256 не совпал — файл повреждён"
            return 1
        fi
        ok "sha256 совпал"
    fi

    # Свободное место ≥ 2× файла (распаковка + страховочный дамп)
    local need_kb avail_kb
    need_kb=$(( $(stat -c %s "$backup_file") / 1024 * 2 ))
    avail_kb=$(df -P "$BACKUP_DIR" | awk 'NR==2{print $4}')
    if [ "$avail_kb" -lt "$need_kb" ]; then
        fail "Недостаточно места: нужно ~${need_kb}KB, доступно ${avail_kb}KB"
        return 1
    fi

    gzip -t "$backup_file" 2>/dev/null || { fail "gzip -t: архив повреждён"; return 1; }

    case "$backup_file" in
        *.tar.gz)
            tar -tzf "$backup_file" >/dev/null 2>&1 || { fail "tar -tzf: архив повреждён"; return 1; }
            rm -rf "$RESTORE_TMP"; mkdir -p "$RESTORE_TMP"
            tar -xzf "$backup_file" -C "$RESTORE_TMP" 2>/dev/null || { fail "распаковка внешнего архива"; rm -rf "$RESTORE_TMP"; return 1; }

            # Обе раскладки: base.tar.gz в корне или в pg_basebackup/
            local base_dir=""
            if [ -f "$RESTORE_TMP/base.tar.gz" ]; then
                base_dir="$RESTORE_TMP"
            elif [ -f "$RESTORE_TMP/pg_basebackup/base.tar.gz" ]; then
                base_dir="$RESTORE_TMP/pg_basebackup"
            else
                fail "base.tar.gz не найден в архиве (ни в корне, ни в pg_basebackup/)"
                rm -rf "$RESTORE_TMP"; return 1
            fi
            gzip -t "$base_dir/base.tar.gz" 2>/dev/null || { fail "inner base.tar.gz повреждён"; rm -rf "$RESTORE_TMP"; return 1; }
            # grep -c (не -q): при pipefail grep -q даёт SIGPIPE на tar → ложный fail
            [ "$(tar -tzf "$base_dir/base.tar.gz" 2>/dev/null | grep -c "backup_label")" -gt 0 ] || {
                fail "в base.tar.gz нет backup_label — это не снэпшот PostgreSQL"
                rm -rf "$RESTORE_TMP"; return 1
            }
            echo "$base_dir" > "$RESTORE_TMP/.base_dir"
            ok "Целостность подтверждена (раскладка: $(basename "$base_dir"))"
            ;;
        *.sql.gz)
            gunzip -t "$backup_file" 2>/dev/null || { fail "gunzip -t: дамп повреждён"; return 1; }
            ok "Целостность подтверждена (sql.gz)"
            ;;
        *)
            fail "Неизвестный формат бэкапа"
            return 1
            ;;
    esac
    return 0
}

# verify_only — режим --verify: только шаг 0, ответ для API (readable|corrupt).
verify_only() {
    if step0_verify "$1"; then
        echo "readable"
        rm -rf "$RESTORE_TMP"
        return 0
    fi
    echo "corrupt"
    rm -rf "$RESTORE_TMP"
    return 1
}

# pre_restore_dump: страховочный дамп текущей БД перед свапом.
# Tripwire: >1KB и в нём есть строки users — иначе дамп бессмысленен.
# Путь дампа → глобальная PRE_RESTORE_DUMP (не stdout: пайпы ломают exit code).
PRE_RESTORE_DUMP=""
pre_restore_dump() {
    local ts dump
    ts=$(date +%Y%m%d_%H%M%S)
    dump="${BACKUP_DIR}/pre_restore_${ts}.sql.gz"
    info "Страховочный дамп текущей БД → $(basename "$dump")"
    if ! docker exec -i "$DB_CONTAINER" pg_dump -U "$DB_USER" "$DB_NAME" 2>/dev/null | gzip > "$dump"; then
        rm -f "$dump"
        fail "pg_dump не прошёл — восстановление отменено"
        return 1
    fi
    if ! gunzip -t "$dump" 2>/dev/null || [ "$(stat -c %s "$dump")" -lt 1024 ]; then
        rm -f "$dump"
        fail "страховочный дамп пустой/битый — восстановление отменено"
        return 1
    fi
    # grep -c (не -q): при pipefail grep -q даёт SIGPIPE на gunzip → ложный fail
    if [ "$(gunzip -c "$dump" | grep -c "COPY public.users\|INSERT INTO public.users\|COPY users\|INSERT INTO users")" -eq 0 ]; then
        fail "в страховочном дампе нет таблицы users — восстановление отменено"
        return 1
    fi
    # Храним только 3 последних
    ls -1t "${BACKUP_DIR}"/pre_restore_*.sql.gz 2>/dev/null | tail -n +4 | xargs -r rm -f
    PRE_RESTORE_DUMP="$dump"
    return 0
}

# rollback_previous: откат на pre_restore дамп (шаг 6).
# mode=tar  — PGDATA заменён снэпшотом и не стартует: отодвигаем данные,
#             initdb через mount schema.sql, применяем дамп.
# mode=sql  — БД жива, схема public сломана применением дампа: только
#             пересоздание схемы + дамп, PGDATA не трогаем!
ROLLBACK_DONE=0
rollback_previous() {
    local dump="$1" mode="${2:-tar}"
    ROLLBACK_DONE=1
    warn "ОТКАТ ($mode): восстанавливаю предыдущее состояние из $(basename "$dump")"

    if [ "$mode" = "tar" ]; then
        # PGDATA — точка монтирования volume: её НЕЛЬЗЯ переименовать (EBUSY).
        # Выносим СОДЕРЖИМОЕ в evidence-каталог на постоянном томе BACKUP_DIR.
        local bad="${BACKUP_DIR}/pgdata.bad-$(date +%Y%m%d_%H%M%S)"
        docker run --rm --volumes-from "$DB_CONTAINER" -v "$BACKUP_DIR":"$BACKUP_DIR" "$ONESHOT_PG_IMAGE" \
            sh -c "mkdir -p '$bad' && find /var/lib/postgresql/data -mindepth 1 -maxdepth 1 -exec mv {} '$bad/' \; && chown -R postgres:postgres /var/lib/postgresql/data" \
            || { restore_status "FATAL" "rollback" "не удалось отодвинуть PGDATA"; return 1; }

        # initdb применит schema.sql из mount /docker-entrypoint-initdb.d
        docker start "$DB_CONTAINER" >/dev/null 2>&1 || true
        local i
        for i in $(seq 1 30); do
            docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1 && break
            sleep 2
        done
    fi

    # Общая часть: чистая схема + дамп (в дампе свои CREATE — схема не должна
    # мешать; pgcrypto пересоздаём явно — DROP SCHEMA ломает его функции).
    docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -q <<'SQL' || true
DROP EXTENSION IF EXISTS pgcrypto CASCADE;
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
GRANT ALL ON SCHEMA public TO public;
SQL

    if gunzip -c "$dump" | docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -q; then
        restore_status "restored-previous" "rollback"
        ok "Откат завершён: БД = состояние до восстановления"
        docker start "$APP_CONTAINER" "$NGINX_CONTAINER" >/dev/null 2>&1 || true
        return 0
    fi
    restore_status "FATAL" "rollback" "psql pre_restore завершился с ошибкой — см. INSTALL «Если восстановление не удалось»"
    fail "ОТКАТ НЕ УДАЛСЯ. Инструкция: docs/setup/INSTALL.md, раздел «Если восстановление не удалось»"
    return 1
}

do_restore() {
    local backup_file="$1"
    
    if [ -z "$backup_file" ] || [ ! -f "$backup_file" ]; then
        fail "Файл бэкапа не найден: $backup_file"
        return 1
    fi
    # Файл обязан лежать внутри BACKUP_DIR (защита от path traversal и сюрпризов)
    case "$(realpath "$backup_file" 2>/dev/null || echo "$backup_file")" in
        "$(realpath "$BACKUP_DIR" 2>/dev/null || echo "$BACKUP_DIR")"/*) ;;
        *) fail "Файл вне $BACKUP_DIR — отказ"; return 1 ;;
    esac

    RESTORE_FILE="$backup_file"
    RESTORE_STARTED_AT=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    step "Восстановление из: $(basename "$backup_file")"

    # Lock: restore и backup взаимно исключаемы (shared lock в общем томе)
    if ! acquire_shared_lock; then
        restore_status "failed" "lock" "бэкап или восстановление уже выполняется"
        return 1
    fi

    # ── Шаг 0: целостность (без побочек в БД) ──
    restore_status "running" "verify"
    if ! step0_verify "$backup_file"; then
        restore_status "corrupt" "verify" "архив не прошёл проверку целостности"
        rm -rf "$RESTORE_TMP"
        return 1
    fi

    # ── Шаг 2: страховочный дамп (до любых разрушений) ──
    restore_status "running" "pre_restore_dump"
    if ! pre_restore_dump; then
        restore_status "failed" "pre_restore_dump" "страховочный дамп не создан"
        rm -rf "$RESTORE_TMP"
        return 1
    fi
    local safety_dump="$PRE_RESTORE_DUMP"

    # trap: при обрыве ПОСЛЕ начала свапа (SIGTERM/EXIT) — автоматический откат
    # к страховочному дампу; иначе просто вернуть контейнеры и снять lock.
    # Используем PRE_RESTORE_DUMP (глобальная): local-scope do_restore к моменту
    # срабатывания EXIT-трапа уже распакован (set -e внутри функции → shell exit).
    SWAP_STARTED=0
    RESTORE_DONE=0
    RESTORE_MODE="tar"
    on_restore_exit() {
        if [ "$SWAP_STARTED" = "1" ] && [ "$RESTORE_DONE" != "1" ] && [ "$ROLLBACK_DONE" != "1" ]; then
            warn "Восстановление прервано во время свапа — автоматический откат"
            rollback_previous "$PRE_RESTORE_DUMP" "$RESTORE_MODE" || true
        else
            docker start "$DB_CONTAINER" "$APP_CONTAINER" "$NGINX_CONTAINER" >/dev/null 2>&1 || true
        fi
        rm -rf "$BACKUP_LOCK_FILE" "$RESTORE_TMP"
    }
    trap on_restore_exit EXIT
    # bash не выполняет EXIT-трап при «голом» SIGTERM/SIGINT — переводим их в exit
    trap 'exit 143' TERM INT

    case "$backup_file" in
        *.tar.gz)
            local base_dir
            base_dir=$(cat "$RESTORE_TMP/.base_dir")

            # ── Шаг 3+4: подготовка и переключение ──
            restore_status "running" "swap"
            SWAP_STARTED=1
            info "Остановка контейнеров (app → nginx → db)..."
            docker stop "$APP_CONTAINER" "$NGINX_CONTAINER" "$DB_CONTAINER" >/dev/null 2>&1 || true

            info "Установка снэпшота в volume db_data (one-shot контейнер)..."
            # Раскладка тара зависит от версии бэкапа — приводим к единой в RESTORE_TMP
            # pg_wal.tar.gz содержит содержимое pg_wal/ БЕЗ префикса — распаковывать
            # внутрь $PGDATA/pg_wal/ (не в корень: иначе WAL не находит postgres).
            local wal_arg="" wal_src=""
            if [ -f "$base_dir/pg_wal.tar.gz" ]; then
                wal_src="$base_dir/pg_wal.tar.gz"
            elif [ -f "$RESTORE_TMP/pg_wal.tar.gz" ]; then
                wal_src="$RESTORE_TMP/pg_wal.tar.gz"
            fi
            if [ -n "$wal_src" ]; then
                [ "$wal_src" = "$RESTORE_TMP/pg_wal.tar.gz" ] || cp "$wal_src" "$RESTORE_TMP/pg_wal.tar.gz"
                wal_arg="mkdir -p /var/lib/postgresql/data/pg_wal && tar -xzf $RESTORE_TMP/pg_wal.tar.gz -C /var/lib/postgresql/data/pg_wal &&"
            fi
            if ! docker run --rm --volumes-from "$DB_CONTAINER" -v "$RESTORE_TMP":"$RESTORE_TMP" "$ONESHOT_PG_IMAGE" \
                sh -c "rm -rf /var/lib/postgresql/data/* &&
                       tar -xzf $base_dir/base.tar.gz -C /var/lib/postgresql/data &&
                       $wal_arg
                       chown -R postgres:postgres /var/lib/postgresql/data"; then
                warn "Установка снэпшота не удалась — откат"
                rollback_previous "$safety_dump"
                return 1
            fi

            # ── Шаг 5: старт + проверка ──
            restore_status "running" "start"
            info "Запуск базы..."
            docker start "$DB_CONTAINER" >/dev/null 2>&1
            local up=0
            for i in $(seq 1 30); do
                if docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1; then up=1; break; fi
                sleep 2
            done
            if [ "$up" != "1" ] || \
               ! docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -t -A -c "SELECT count(*) FROM users" >/dev/null 2>&1; then
                warn "БД не стартовала или users недоступна — откат"
                rollback_previous "$safety_dump"
                return 1
            fi
            info "Запуск app и nginx..."
            docker start "$APP_CONTAINER" "$NGINX_CONTAINER" >/dev/null 2>&1 || true
            sleep 5
            if docker exec "$APP_CONTAINER" wget -qO- http://localhost:8080/health >/dev/null 2>&1; then
                ok "Healthcheck прошёл"
            else
                warn "Healthcheck не прошёл (app поднимается дольше) — проверьте вручную"
            fi
            ;;

        *.sql.gz)
            # sql.gz-путь: останавливаем только app, БД живая.
            # DROP SCHEMA обязателен: иначе наложение дампа на живую схему коллизит.
            restore_status "running" "sql_restore"
            RESTORE_MODE="sql"
            SWAP_STARTED=1
            info "Восстановление из SQL-дампа (остановка app)..."
            docker stop "$APP_CONTAINER" >/dev/null 2>&1 || true

            if ! docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -q <<'SQL'
DROP EXTENSION IF EXISTS pgcrypto CASCADE;
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
GRANT ALL ON SCHEMA public TO public;
SQL
            then
                warn "DROP SCHEMA не прошёл — откат"
                rollback_previous "$safety_dump" "sql"
                return 1
            fi

            if gunzip -c "$backup_file" | docker exec -i "$DB_CONTAINER" \
                psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -q; then
                ok "SQL-дамп применён"
            else
                warn "psql завершился с ошибкой (ON_ERROR_STOP) — откат"
                rollback_previous "$safety_dump" "sql"
                return 1
            fi

            restore_status "running" "start"
            docker start "$APP_CONTAINER" >/dev/null 2>&1 || true
            sleep 5
            if docker exec "$APP_CONTAINER" wget -qO- http://localhost:8080/health >/dev/null 2>&1; then
                ok "Healthcheck прошёл"
            else
                warn "Healthcheck не прошёл — проверьте вручную"
            fi
            ;;
    esac

    rm -rf "$RESTORE_TMP"
    RESTORE_DONE=1
    restore_status "restored" "done"
    ok "Восстановление завершено успешно"
}

# ============================================================================
# Статус
# ============================================================================

show_status() {
    step "Статус системы"
    
    # Последний бэкап
    local last_backup
    last_backup=$(ls -t "$BACKUP_DIR"/db_*.tar.gz "$BACKUP_DIR"/db_*.sql.gz 2>/dev/null | head -1 || true)
    
    if [ -n "$last_backup" ]; then
        local last_time
        last_time=$(stat -c %Y "$last_backup")
        local now
        now=$(date +%s)
        local age_hours=$(( (now - last_time) / 3600 ))
        local size
        size=$(du -h "$last_backup" | cut -f1)
        ok "Последний бэкап: $(basename "$last_backup") ($size, ${age_hours}ч назад)"
    else
        warn "Бэкапов не найдено"
    fi
    
    # Количество бэкапов
    local count
    count=$(find "$BACKUP_DIR" -name "db_*" -type f | wc -l)
    info "Всего бэкапов: $count"
    
    # Нагрузка
    local one five fifteen
    read -r one five fifteen _ _ < /proc/loadavg
    info "Нагрузка: 1min=$one 5min=$five 15min=$fifteen"
    
    # Диск
    local disk
    disk=$(df -h "$BACKUP_DIR" | awk 'NR==2{print $4}')
    info "Свободное место: $disk"
}

# ============================================================================
# Основная функция
# ============================================================================

main() {
    local command="${1:-}"
    
    case "$command" in
        --scheduled)
            step "Запуск по расписанию"
            mkdir -p "$BACKUP_DIR"
            load_backup_config
            check_stuck_backup
            
            # Ротация всегда — иначе при пропуске бэкапа (load_high) журнал и файлы копятся
            cleanup_old_backups
            
            if ! check_backup_window; then
                info "Вне временного окна ($BACKUP_HOUR_START:00 - $BACKUP_HOUR_END:00)"
                return 0
            fi
            
            if ! check_backup_interval; then
                return 0
            fi
            
            if ! check_load; then
                log_incident "load_high" "Нагрузка слишком высока, бэкап отложен" "warning"
                return 0
            fi
            
            acquire_backup_lock || return 0
            
            if [ "$BACKUP_METHOD" = "pg_basebackup" ]; then
                do_backup_pg_basebackup
            else
                do_backup_pg_dump
            fi
            
            check_missed_backups
            ;;
            
        --restore)
            do_restore "${2:-}"
            ;;
            
        --verify)
            verify_only "${2:-}"
            ;;
            
        --status)
            show_status
            ;;
            
        "")
            step "Ручной бэкап"
            mkdir -p "$BACKUP_DIR"
            load_backup_config
            
            if ! acquire_backup_lock; then
                return 1
            fi
            
            if [ "$BACKUP_METHOD" = "pg_basebackup" ]; then
                do_backup_pg_basebackup
            else
                do_backup_pg_dump
            fi
            
            cleanup_old_backups
            ;;
            
        *)
            echo "Usage: $0 [--scheduled|--restore <file>|--verify <file>|--status]"
            return 1
            ;;
    esac
}

main "$@"

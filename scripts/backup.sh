#!/usr/bin/env bash
# backup.sh — Бэкап PostgreSQL для MCP-IS
#
# Поддерживает pg_basebackup (файловый снэпшот) и pg_dump (SQL-дамп).
# Мониторинг нагрузки, удалённые хранилища, инциденты.
#
# Usage:
#   ./scripts/backup.sh                   # ручной бэкап
#   ./scripts/backup.sh --scheduled       # запуск по расписанию (cron)
#   ./scripts/backup.sh --restore <file>  # восстановление из бэкапа
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
readonly BACKUP_LOCK_FILE="/tmp/mcp-is-backup.lock"
readonly DB_CONTAINER="mcp_is_db"
readonly DB_NAME="mcp_is"
readonly DB_USER="admin"

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
        rm -f "$BACKUP_PID_FILE" "$BACKUP_LOCK_FILE"
    fi
    
    return 0
}

cleanup_partial_backup() {
    # Удалить неполные .tar.gz файлы (меньше 1KB и старше 5 минут)
    find "$BACKUP_DIR" -name "*.tar.gz" -size -1k -mmin +5 -delete 2>/dev/null || true
    find "$BACKUP_DIR" -name "*.sql.gz" -size -1k -mmin +5 -delete 2>/dev/null || true
    rm -f "$BACKUP_LOCK_FILE"
}

# ============================================================================
# Блокировка запуска нового бэкапа
# ============================================================================

acquire_backup_lock() {
    if [ -f "$BACKUP_LOCK_FILE" ]; then
        local lock_pid
        lock_pid=$(cat "$BACKUP_LOCK_FILE")
        if kill -0 "$lock_pid" 2>/dev/null; then
            warn "Бэкап уже выполняется (PID=$lock_pid)"
            return 1
        fi
        rm -f "$BACKUP_LOCK_FILE"
    fi
    
    echo $$ > "$BACKUP_LOCK_FILE"
    echo $$ > "$BACKUP_PID_FILE"
    trap 'rm -f "$BACKUP_LOCK_FILE" "$BACKUP_PID_FILE"' EXIT
}

# ============================================================================
# Бэкап pg_basebackup
# ============================================================================

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
        docker cp "mcp_is_db:/tmp/pgbackup/" "/tmp/pgbackup/"
        tar -czf "$backup_file" -C /tmp/pgbackup .
        rm -rf /tmp/pgbackup/
        docker exec -i "$DB_CONTAINER" rm -rf /tmp/pgbackup
        
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
# Восстановление из бэкапа
# ============================================================================

do_restore() {
    local backup_file="$1"
    
    if [ -z "$backup_file" ] || [ ! -f "$backup_file" ]; then
        fail "Файл бэкапа не найден: $backup_file"
        return 1
    fi
    
    step "Восстановление из: $(basename "$backup_file")"
    
    # shortcut: restore останавливает собственный контейнер app, если скрипт
    # запущен внутри него — восстанавливать с хоста: /opt/mcp-is/scripts/backup.sh --restore <file>
    info "Остановка контейнеров..."
    docker stop mcp_is_app mcp_is_nginx
    
    # Распаковать бэкап
    if [[ "$backup_file" == *.tar.gz ]]; then
        info "Распаковка tar.gz бэкапа..."
        docker exec -i "$DB_CONTAINER" sh -c "rm -rf /var/lib/postgresql/data/*"
        tar -xzf "$backup_file" -C /tmp/
        docker cp /tmp/pg_basebackup/. "${DB_CONTAINER}:/var/lib/postgresql/data/"
        rm -rf /tmp/pg_basebackup/
        
        # Настроить восстановление
        docker exec -i "$DB_CONTAINER" \
            sh -c "touch /var/lib/postgresql/data/recovery.signal"
    elif [[ "$backup_file" == *.sql.gz ]]; then
        info "Восстановление из SQL-дампа..."
        gunzip -c "$backup_file" | docker exec -i "$DB_CONTAINER" \
            psql -U "$DB_USER" -d "$DB_NAME" -q
    else
        fail "Неизвестный формат бэкапа"
        return 1
    fi
    
    # Запустить контейнеры
    info "Запуск контейнеров..."
    docker start "$DB_CONTAINER"
    sleep 10
    docker start mcp_is_app mcp_is_nginx
    
    # Проверка здоровья
    sleep 10
    if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        ok "Восстановление завершено успешно"
    else
        warn "Восстановление завершено, но healthcheck не проходит"
    fi
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
            echo "Usage: $0 [--scheduled|--restore <file>|--status]"
            return 1
            ;;
    esac
}

main "$@"

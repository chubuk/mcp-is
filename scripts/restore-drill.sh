#!/usr/bin/env bash
# restore-drill.sh — локальное учение backup → verify → restore.
#
# Полностью изолировано: throwaway-контейнер mcpisgate_db и каталог
# /tmp/mcpis-drill-backups. Штатные mcp_is_* контейнеры и прод НЕ трогаются.
#
# Usage:
#   ./scripts/restore-drill.sh            # учение (default: BACKUP_METHOD=pg_dump)
#   BACKUP_METHOD=pg_basebackup ./scripts/restore-drill.sh
#   ./scripts/restore-drill.sh --keep     # не удалять контейнер/артефакты в конце
#
# Проверяет: бэкап создаётся, --verify даёт readable, битый файл даёт corrupt,
# --restore возвращает данные (маркерная таблица), restore_status.json пишется.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# ── Изоляция: только наши имена, никаких mcp_is_* ──
export DB_CONTAINER="mcpisgate_db"
export APP_CONTAINER="mcpisgate_app"   # не существует — backup.sh переживает это с || true
export NGINX_CONTAINER="mcpisgate_nginx"
export DB_NAME="mcp_is"
export DB_USER="admin"
export BACKUP_DIR="/tmp/mcpis-drill-backups"
export BACKUP_METHOD="${BACKUP_METHOD:-pg_dump}"
DB_IMAGE="postgres:16-alpine"
DB_PASSWORD="drill-password-not-a-secret"

KEEP=0
[ "${1:-}" = "--keep" ] && KEEP=1

info() { echo "[INFO] $1"; }
ok()   { echo "[OK]   $1"; }
fail() { echo "[FAIL] $1"; exit 1; }

psql_drill() {
    docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -q -t -A "$@"
}

cleanup() {
    if [ "$KEEP" = "1" ]; then
        info "--keep: контейнер $DB_CONTAINER и $BACKUP_DIR оставлены"
        return
    fi
    docker rm -f "$DB_CONTAINER" >/dev/null 2>&1 || true
    rm -rf "$BACKUP_DIR"
    info "контейнер и артефакты удалены"
}
trap cleanup EXIT

step() { echo; echo "==> $1"; }

# ── Шаг 0: throwaway-БД со схемой ──
step "Throwaway-БД ($DB_CONTAINER, метод: $BACKUP_METHOD)"
docker rm -f "$DB_CONTAINER" >/dev/null 2>&1 || true
rm -rf "$BACKUP_DIR"
mkdir -p "$BACKUP_DIR"
docker run -d --name "$DB_CONTAINER" \
    -e POSTGRES_USER="$DB_USER" \
    -e POSTGRES_PASSWORD="$DB_PASSWORD" \
    -e POSTGRES_DB="$DB_NAME" \
    "$DB_IMAGE" >/dev/null
for i in $(seq 1 30); do
    docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1 && break
    sleep 1
done
docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1 \
    || fail "throwaway-БД не стартовала"
docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -q < "$PROJECT_DIR/internal/storage/schema.sql" \
    || fail "schema.sql не применился"
psql_drill -c "CREATE TABLE IF NOT EXISTS drill_marker(v text); INSERT INTO drill_marker VALUES ('before-backup');"
ok "схема применена, маркер записан"

# ── Шаг 1: бэкап ──
step "Бэкап"
"$PROJECT_DIR/scripts/backup.sh" || fail "backup.sh завершился с ошибкой"
backup_file=$(ls -t "$BACKUP_DIR" | grep -E '^db_.*\.(tar|sql)\.gz$' | head -1 || true)
[ -n "$backup_file" ] && [ -s "$BACKUP_DIR/$backup_file" ] || fail "файл бэкапа не создан"
backup_file="$BACKUP_DIR/$backup_file"
ok "бэкап: $backup_file"

# ── Шаг 2: verify ──
step "Verify (ожидается readable)"
verify_out=$("$PROJECT_DIR/scripts/backup.sh" --verify "$backup_file" | tail -1 || true)
echo "$verify_out" | grep -q "readable" || fail "verify: $verify_out"
ok "readable"

# ── Шаг 3: битый файл → corrupt ──
step "Verify битого файла (ожидается corrupt)"
echo "not-a-backup" > "$BACKUP_DIR/db_corrupt.sql.gz"
corrupt_out=$("$PROJECT_DIR/scripts/backup.sh" --verify "$BACKUP_DIR/db_corrupt.sql.gz" | tail -1 || true)
echo "$corrupt_out" | grep -q "corrupt" || fail "verify битого: $corrupt_out"
ok "corrupt"

# ── Шаг 4: портим данные ──
step "Имитация потери данных"
psql_drill -c "DROP TABLE drill_marker;" || fail "не удалось удалить маркер"
ok "маркер удалён"

# ── Шаг 5: restore ──
step "Restore"
"$PROJECT_DIR/scripts/backup.sh" --restore "$backup_file" || fail "restore завершился с ошибкой"
[ -f "$BACKUP_DIR/restore_status.json" ] || fail "restore_status.json не создан"
cat "$BACKUP_DIR/restore_status.json"
ok "restore отработал"

# ── Шаг 6: данные вернулись ──
step "Проверка данных после restore"
marker=$(psql_drill -c "SELECT v FROM drill_marker LIMIT 1")
[ "$marker" = "before-backup" ] || fail "маркер не вернулся: '$marker'"
ok "маркер восстановлен: $marker"

echo
ok "УЧЕНИЕ ПРОЙДЕНО: backup → verify → corrupt-detect → restore → данные на месте"

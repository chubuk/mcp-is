#!/bin/bash
# Schema check script — compares Go structs with DB columns
# Usage: ./scripts/check-schema.sh

set -e

DB_CONTAINER="mcp_is_db"
DB_NAME="mcp_is"
DB_USER="admin"

echo "=== Database Schema Check ==="
echo ""

# Check if DB is running
if ! docker ps --format '{{.Names}}' | grep -q "$DB_CONTAINER"; then
    echo "ERROR: Container $DB_CONTAINER is not running"
    exit 1
fi

# Get DB columns
echo "DB columns (users, shops, ai_keys):"
docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -t -A -c "
    SELECT table_name || '.' || column_name
    FROM information_schema.columns
    WHERE table_name IN ('users','shops','ai_keys')
    ORDER BY table_name, ordinal_position;
" 2>/dev/null

echo ""
echo "=== Checking for missing migrations ==="

# Check Go struct fields
for table in users shops ai_keys; do
    echo ""
    echo "--- $table ---"
    echo "DB columns:"
    docker exec "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -t -A -c "
        SELECT column_name FROM information_schema.columns
        WHERE table_name='$table' ORDER BY ordinal_position;
    " 2>/dev/null | tr '\n' ' '
    echo ""
    echo "schema.sql:"
    grep -A 20 "CREATE TABLE.*$table" internal/storage/schema.sql 2>/dev/null | grep -oP '\w+\s+(TEXT|BOOLEAN|UUID|INTEGER|TIMESTAMP)' | awk '{print $1}' | tr '\n' ' '
    echo ""
done

echo ""
echo "=== Done ==="

#!/bin/bash
# PostgreSQL backup script for MCP IS
# Usage: ./scripts/backup.sh [backup_dir]

set -e

BACKUP_DIR="${1:-/opt/mcp-is/backups}"
DATE=$(date +%Y%m%d_%H%M%S)
CONTAINER="mcp_is_db"
DB_NAME="mcp_is"
DB_USER="admin"
RETENTION_DAYS=7

# Create backup directory
mkdir -p "$BACKUP_DIR"

echo "[$(date)] Starting backup..."

# Dump database
docker exec "$CONTAINER" pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "$BACKUP_DIR/db_$DATE.sql.gz"

# Get backup size
SIZE=$(du -h "$BACKUP_DIR/db_$DATE.sql.gz" | cut -f1)
echo "[$(date)] Backup completed: db_$DATE.sql.gz ($SIZE)"

# Cleanup old backups
find "$BACKUP_DIR" -name "db_*.sql.gz" -mtime +$RETENTION_DAYS -delete 2>/dev/null
REMAINING=$(find "$BACKUP_DIR" -name "db_*.sql.gz" | wc -l)
echo "[$(date)] Retention: $REMAINING backups kept (retention: $RETENTION_DAYS days)"

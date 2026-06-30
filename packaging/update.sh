#!/bin/bash
set -e

INSTALL_DIR="/opt/mcp-insales"
CONF_DIR="/etc/mcp-insales"
ENV_FILE="${CONF_DIR}/.env"
VERSION_FILE="${INSTALL_DIR}/VERSION"
BACKUP_DIR="${INSTALL_DIR}/backups"

echo ""
echo "====================================="
echo "  MCP InSales Manual Update"
echo "====================================="
echo ""

# Check current version
if [ -f "${VERSION_FILE}" ]; then
    CURRENT_VERSION=$(cat "${VERSION_FILE}")
    echo "Current version: ${CURRENT_VERSION}"
else
    echo "Current version: unknown (first install?)"
fi

# Check if new package is available
NEW_VERSION=$(dpkg -s mcp-insales 2>/dev/null | grep "^Version:" | awk '{print $2}' || echo "unknown")
echo "Package version: ${NEW_VERSION}"
echo ""

if [ "${CURRENT_VERSION}" = "${NEW_VERSION}" ]; then
    echo "Already up to date."
    echo ""
    read -p "Force reinstall? [y/N]: " FORCE
    if [[ ! "${FORCE}" =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

# Create backup
echo ""
echo "Creating backup..."
mkdir -p "${BACKUP_DIR}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Backup database
if docker ps --format '{{.Names}}' | grep -q "mcp_insales_db"; then
    echo "  Backing up database..."
    cd "${INSTALL_DIR}"
    docker compose exec -T db pg_dump -U admin mcp_insales > "${BACKUP_DIR}/db_${TIMESTAMP}.sql" 2>/dev/null || true
    if [ -s "${BACKUP_DIR}/db_${TIMESTAMP}.sql" ]; then
        echo "  Database backup: ${BACKUP_DIR}/db_${TIMESTAMP}.sql"
    else
        echo "  WARNING: Database backup failed"
    fi
else
    echo "  Database container not running, skipping DB backup"
fi

# Backup files
echo "  Backing up project files..."
cd "${INSTALL_DIR}"
tar czf "${BACKUP_DIR}/files_${TIMESTAMP}.tar.gz" \
    --exclude='backups' \
    --exclude='.env' \
    internal/ web/ docker-compose.yaml Dockerfile .env.example 2>/dev/null || true

if [ -s "${BACKUP_DIR}/files_${TIMESTAMP}.tar.gz" ]; then
    echo "  Files backup: ${BACKUP_DIR}/files_${TIMESTAMP}.tar.gz"
fi

echo ""
echo "Backup completed."
echo ""

# Stop service
echo "Stopping service..."
systemctl stop mcp-insales 2>/dev/null || true

# Rebuild containers
echo "Rebuilding containers..."
cd "${INSTALL_DIR}"
docker compose up -d --build --remove-orphans

# Wait for database
echo "Waiting for database..."
WAIT_COUNT=0
MAX_WAIT=30
while [ ${WAIT_COUNT} -lt ${MAX_WAIT} ]; do
    if docker compose exec -T db pg_isready -U admin -d mcp_insales &>/dev/null; then
        echo "Database is ready."
        break
    fi
    sleep 2
    WAIT_COUNT=$((WAIT_COUNT + 1))
done

# Run migrations
echo ""
echo "Running migrations..."
MIGRATIONS_FILE="${INSTALL_DIR}/internal/storage/migrations.sql"
if [ -f "${MIGRATIONS_FILE}" ]; then
    docker compose exec -T db psql -U admin -d mcp_insales -f /dev/stdin < "${MIGRATIONS_FILE}" 2>/dev/null && \
        echo "Migrations applied." || \
        echo "WARNING: Some migrations may have failed"
fi

# Health check
echo ""
echo "Running health check..."
sleep 3
if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
    echo "Health check passed!"
else
    echo "WARNING: Health check failed"
fi

# Save version
echo "${NEW_VERSION}" > "${VERSION_FILE}"

# Start service
echo ""
echo "Starting service..."
systemctl start mcp-insales

echo ""
echo "====================================="
echo "  Update complete!"
echo "====================================="
echo ""
echo "  Version: ${NEW_VERSION}"
echo "  Backup:  ${BACKUP_DIR}/"
echo ""

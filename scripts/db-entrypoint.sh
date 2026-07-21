#!/bin/bash
set -e

# Run standard entrypoint first
docker-entrypoint.sh postgres &

# Wait for PostgreSQL to be ready
until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" > /dev/null 2>&1; do
    sleep 1
done

echo "PostgreSQL is ready. Running migrations..."

# Run migrations (safe to run multiple times due to IF NOT EXISTS)
if [ -f /docker-entrypoint-initdb.d/02-migrations.sql ]; then
    psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/02-migrations.sql
    echo "Migrations completed."
fi

# Keep container running
wait

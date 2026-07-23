-- MCP IS - Database Migrations
-- Этот файл выполняется при каждом обновлении deb-пакета.
-- Используйте ADD COLUMN IF NOT EXISTS для безопасного повторного выполнения.

-- Миграция v0.1.0 (начальная версия)
-- (schema.sql уже содержит базовую схему)

-- Миграция v0.2.0 - Per-shop read-only
ALTER TABLE shops ADD COLUMN IF NOT EXISTS is_read_only BOOLEAN DEFAULT FALSE;

-- Миграция v0.2.1 - Jurisdiction for shops
ALTER TABLE shops ADD COLUMN IF NOT EXISTS jurisdiction TEXT DEFAULT 'AUTO';

-- Миграция v0.2.2 - PDN enabled for shops
ALTER TABLE shops ADD COLUMN IF NOT EXISTS pdn_enabled BOOLEAN DEFAULT FALSE;

-- Миграция v0.3.0 (будущее)
-- ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;

-- Миграция v0.4.0 - Require password change flag
ALTER TABLE users ADD COLUMN IF NOT EXISTS require_password_change BOOLEAN DEFAULT FALSE;

-- Миграция v0.5.0 - PDN enabled flag for AI keys
ALTER TABLE ai_keys ADD COLUMN IF NOT EXISTS pdn_enabled BOOLEAN DEFAULT FALSE;

-- Миграция v0.6.0 - Backup tables
CREATE TABLE IF NOT EXISTS backup_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ts TIMESTAMPTZ DEFAULT NOW(),
    status TEXT NOT NULL,
    size_bytes BIGINT DEFAULT 0,
    error_msg TEXT DEFAULT '',
    load_1m NUMERIC(4,2),
    load_5m NUMERIC(4,2),
    load_15m NUMERIC(4,2),
    iowait_pct NUMERIC(5,2),
    duration_ms INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_backup_log_ts ON backup_log(ts DESC);

CREATE TABLE IF NOT EXISTS backup_incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ts TIMESTAMPTZ DEFAULT NOW(),
    type TEXT NOT NULL,
    message TEXT NOT NULL,
    severity TEXT NOT NULL,
    resolved BOOLEAN DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_backup_incidents_ts ON backup_incidents(ts DESC);
CREATE INDEX IF NOT EXISTS idx_backup_incidents_severity ON backup_incidents(severity);

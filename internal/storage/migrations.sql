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

-- Миграция v0.7.0 - MCP client types
CREATE TABLE IF NOT EXISTS mcp_client_types (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name TEXT NOT NULL UNIQUE,
    config_key TEXT NOT NULL,
    config_format TEXT NOT NULL DEFAULT 'json',
    config_template TEXT NOT NULL,
    file_name TEXT NOT NULL DEFAULT '',
    file_description TEXT NOT NULL DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Миграция v0.8.0 - InSales app installs (external application mode)
ALTER TABLE ai_keys ADD COLUMN IF NOT EXISTS shop_id UUID REFERENCES shops(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS app_installs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    insales_id TEXT UNIQUE NOT NULL,
    shop_host TEXT NOT NULL,
    api_password_encrypted TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'installed',
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    shop_id UUID REFERENCES shops(id) ON DELETE SET NULL,
    installed_at TIMESTAMPTZ DEFAULT NOW(),
    uninstalled_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_app_installs_user ON app_installs(user_id);
CREATE INDEX IF NOT EXISTS idx_app_installs_status ON app_installs(status);

-- Миграция v0.9.0 - Email templates (runtime edit)
CREATE TABLE IF NOT EXISTS email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    subject TEXT NOT NULL,
    html_body TEXT NOT NULL,
    text_body TEXT DEFAULT '',
    description TEXT DEFAULT '',
    required_vars TEXT[] DEFAULT '{}',
    is_modified BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

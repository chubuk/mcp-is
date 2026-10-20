CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nickname TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    phone TEXT UNIQUE,
    first_name TEXT,
    last_name TEXT,
    password_hash TEXT NOT NULL,
    role TEXT DEFAULT 'user' CHECK (role IN ('user', 'admin', 'super_admin')),
    is_active BOOLEAN DEFAULT TRUE,
    require_password_change BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone) WHERE phone IS NOT NULL;

-- Shops (InSales credentials)
CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT DEFAULT '',
    shop_url TEXT NOT NULL,
    api_key_id TEXT NOT NULL,
    api_password_encrypted TEXT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    is_read_only BOOLEAN DEFAULT FALSE,
    jurisdiction TEXT DEFAULT 'AUTO',
    pdn_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_shops_user_id ON shops(user_id);

-- AI Keys (for MCP clients)
CREATE TABLE IF NOT EXISTS ai_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    key_hash TEXT UNIQUE NOT NULL,
    name TEXT,
    available_commands TEXT DEFAULT '',
    commands_updated_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_keys_hash ON ai_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_ai_keys_user_id ON ai_keys(user_id);

-- Password reset tokens
CREATE TABLE IF NOT EXISTS password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets(token);
CREATE INDEX IF NOT EXISTS idx_password_resets_user_id ON password_resets(user_id);

-- Server settings (domain, HTTPS, cert info)
CREATE TABLE IF NOT EXISTS server_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- PDN access log (152-ФЗ compliance)
CREATE TABLE IF NOT EXISTS pdn_access_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ts TIMESTAMPTZ DEFAULT NOW(),
    user_id UUID,
    shop_id UUID,
    tool TEXT,
    data_types TEXT[],
    jurisdiction TEXT,
    llm_provider TEXT,
    masked BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_pdn_user_id ON pdn_access_log(user_id);
CREATE INDEX IF NOT EXISTS idx_pdn_shop_id ON pdn_access_log(shop_id);
CREATE INDEX IF NOT EXISTS idx_pdn_ts ON pdn_access_log(ts);

-- AI keys bound to a specific shop (app installs). NULL = user's default shop (personal mode).
ALTER TABLE ai_keys ADD COLUMN IF NOT EXISTS shop_id UUID REFERENCES shops(id) ON DELETE SET NULL;

-- InSales external-app installs (protocol: password = MD5(token + secret_key))
CREATE TABLE IF NOT EXISTS app_installs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    insales_id TEXT UNIQUE NOT NULL,
    shop_host TEXT NOT NULL,
    api_password_encrypted TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'installed', -- installed | uninstalled | pending_link
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    shop_id UUID REFERENCES shops(id) ON DELETE SET NULL,
    installed_at TIMESTAMPTZ DEFAULT NOW(),
    uninstalled_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_app_installs_user ON app_installs(user_id);
CREATE INDEX IF NOT EXISTS idx_app_installs_status ON app_installs(status);

-- Email templates (runtime-editable, seeded from embed defaults)
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

-- Backup log
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

-- Backup incidents
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

-- Audit log for MCP write-tool calls (args stored as hash only — PDN)
CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    request_id TEXT,
    user_id TEXT,
    api_key_id TEXT,
    shop_id TEXT,
    tool_name TEXT NOT NULL,
    action TEXT NOT NULL,
    args_hash TEXT,
    result_ok BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_audit_log_created ON audit_log(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_tool ON audit_log(tool_name);

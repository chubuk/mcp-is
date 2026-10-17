package storage

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) {
	migrations := []string{
		`ALTER TABLE shops ADD COLUMN IF NOT EXISTS is_read_only BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE shops ADD COLUMN IF NOT EXISTS jurisdiction TEXT DEFAULT 'AUTO'`,
		`ALTER TABLE shops ADD COLUMN IF NOT EXISTS pdn_enabled BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS require_password_change BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE ai_keys ADD COLUMN IF NOT EXISTS pdn_enabled BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE ai_keys ADD COLUMN IF NOT EXISTS shop_id UUID REFERENCES shops(id) ON DELETE SET NULL`,
		`CREATE TABLE IF NOT EXISTS app_installs (
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
		)`,
		`CREATE INDEX IF NOT EXISTS idx_app_installs_user ON app_installs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_app_installs_status ON app_installs(status)`,
		`CREATE TABLE IF NOT EXISTS email_templates (
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
		)`,
		`CREATE TABLE IF NOT EXISTS mcp_client_types (
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
		)`,
	}

	for i, migration := range migrations {
		_, err := pool.Exec(context.Background(), migration)
		if err != nil {
			slog.Error("migration failed", "index", i+1, "error", err)
		} else {
			slog.Debug("migration applied", "index", i+1)
		}
	}

	// Also run from migrations.sql file if exists
	if data, err := os.ReadFile("internal/storage/migrations.sql"); err == nil {
		_, err := pool.Exec(context.Background(), string(data))
		if err != nil {
			slog.Debug("migrations.sql already applied or error", "error", err)
		}
	}

	slog.Info("database migrations completed")
}

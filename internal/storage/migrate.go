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
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS require_password_change BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE ai_keys ADD COLUMN IF NOT EXISTS pdn_enabled BOOLEAN DEFAULT FALSE`,
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

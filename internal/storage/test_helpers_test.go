//go:build integration

package storage

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:password@localhost:5432/mcp_is?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Fatalf("ping test db: %v", err)
	}

	// Ensure base schema exists (schema.sql is idempotent enough for CREATE IF NOT EXISTS paths).
	var one int
	if err := pool.QueryRow(context.Background(), `SELECT 1 FROM information_schema.tables WHERE table_name='users'`).Scan(&one); err != nil {
		schema, err2 := os.ReadFile("schema.sql")
		if err2 != nil {
			// allow running from repo root
			schema, err2 = os.ReadFile("internal/storage/schema.sql")
		}
		if err2 != nil {
			t.Fatalf("users table missing and schema.sql unreadable: %v", err2)
		}
		if _, err3 := pool.Exec(context.Background(), string(schema)); err3 != nil {
			t.Fatalf("apply schema.sql: %v", err3)
		}
	}
	// idempotent column/table migrations (pdn_enabled, mcp_client_types, …)
	RunMigrations(pool)

	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM audit_log")
		pool.Exec(context.Background(), "DELETE FROM password_resets")
		pool.Exec(context.Background(), "DELETE FROM ai_keys")
		pool.Exec(context.Background(), "DELETE FROM app_installs")
		pool.Exec(context.Background(), "DELETE FROM email_templates")
		pool.Exec(context.Background(), "DELETE FROM shops")
		pool.Exec(context.Background(), "DELETE FROM users")
		pool.Close()
	})

	return pool
}

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
		dbURL = "postgres://admin:password@localhost:5432/mcp_insales?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Fatalf("ping test db: %v", err)
	}

	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM ai_keys")
		pool.Exec(context.Background(), "DELETE FROM shops")
		pool.Exec(context.Background(), "DELETE FROM users")
		pool.Close()
	})

	return pool
}

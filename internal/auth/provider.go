package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider interface {
	Validate(apiKey string) (string, error)
}

type DBProvider struct {
	db *pgxpool.Pool
}

func NewDBProvider(db *pgxpool.Pool) *DBProvider {
	return &DBProvider{db: db}
}

func (p *DBProvider) Validate(apiKey string) (string, error) {
	var userID string
	err := p.db.QueryRow(context.Background(),
		`SELECT id FROM users WHERE password_hash = $1`, apiKey,
	).Scan(&userID)

	if err != nil {
		return "", fmt.Errorf("invalid api key: %w", err)
	}
	return userID, nil
}

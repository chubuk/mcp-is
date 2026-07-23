package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingsRepository struct {
	db *pgxpool.Pool
}

func NewSettingsRepository(db *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(key string) (string, error) {
	var value string
	err := r.db.QueryRow(context.Background(),
		`SELECT value FROM server_settings WHERE key = $1`, key,
	).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("setting not found: %w", err)
	}
	return value, nil
}

func (r *SettingsRepository) GetDefault(key, fallback string) string {
	value, err := r.Get(key)
	if err != nil {
		return fallback
	}
	return value
}

func (r *SettingsRepository) Set(key, value string) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO server_settings (key, value, updated_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()`,
		key, value,
	)
	return err
}

func (r *SettingsRepository) GetAll() (map[string]string, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT key, value FROM server_settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

func (r *SettingsRepository) InitDefaults() error {
	defaults := map[string]string{
		"domain":              "localhost",
		"https_enabled":       "false",
		"auto_renew":          "true",
		"renew_days_before":   "30",
		"cert_expires_at":     "",
		"cert_issuer":         "",
		"cert_needs_attention": "false",
		"last_renew_error":    "",
		"retry_count":         "0",
		"last_retry_at":       "",
		"letsencrypt_email":   "",
	}
	for key, value := range defaults {
		_ = r.Set(key, value)
	}
	return nil
}

func (r *SettingsRepository) SetDefaultsIfMissing() error {
	defaults := map[string]string{
		"domain":              "localhost",
		"https_enabled":       "false",
		"auto_renew":          "true",
		"renew_days_before":   "30",
		"cert_expires_at":     "",
		"cert_issuer":         "",
		"cert_needs_attention": "false",
		"last_renew_error":    "",
		"retry_count":         "0",
		"last_retry_at":       "",
		"letsencrypt_email":   "",
	}
	for key, value := range defaults {
		_, err := r.db.Exec(context.Background(),
			`INSERT INTO server_settings (key, value, updated_at)
			 SELECT $1, $2, NOW()
			 WHERE NOT EXISTS (SELECT 1 FROM server_settings WHERE key = $1)`,
			key, value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

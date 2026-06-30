package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AIKey struct {
	ID         string
	UserID     string
	KeyHash    string
	Name       string
	LastUsedAt string
	CreatedAt  string
}

type AIKeyRepository struct {
	db *pgxpool.Pool
}

func NewAIKeyRepository(db *pgxpool.Pool) *AIKeyRepository {
	return &AIKeyRepository{db: db}
}

func (r *AIKeyRepository) Create(userID, keyHash, name string) (*AIKey, error) {
	var k AIKey
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO ai_keys (user_id, key_hash, name)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, key_hash, COALESCE(name,''), created_at::text`,
		userID, keyHash, nullStr(name),
	).Scan(&k.ID, &k.UserID, &k.KeyHash, &k.Name, &k.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create ai key: %w", err)
	}
	return &k, nil
}

func (r *AIKeyRepository) GetByHash(keyHash string) (*AIKey, error) {
	var k AIKey
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, key_hash, COALESCE(name,''), COALESCE(last_used_at::text,''), created_at::text
		 FROM ai_keys WHERE key_hash = $1`, keyHash,
	).Scan(&k.ID, &k.UserID, &k.KeyHash, &k.Name, &k.LastUsedAt, &k.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("ai key not found: %w", err)
	}
	return &k, nil
}

func (r *AIKeyRepository) ListByUserID(userID string) ([]AIKey, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, key_hash, COALESCE(name,''), COALESCE(last_used_at::text,''), created_at::text
		 FROM ai_keys WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list ai keys: %w", err)
	}
	defer rows.Close()

	var keys []AIKey
	for rows.Next() {
		var k AIKey
		if err := rows.Scan(&k.ID, &k.UserID, &k.KeyHash, &k.Name, &k.LastUsedAt, &k.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ai key: %w", err)
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func (r *AIKeyRepository) UpdateLastUsed(id string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE ai_keys SET last_used_at = NOW() WHERE id = $1`, id,
	)
	return err
}

func (r *AIKeyRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM ai_keys WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete ai key: %w", err)
	}
	return nil
}

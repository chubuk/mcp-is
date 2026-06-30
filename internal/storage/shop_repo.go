package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Shop struct {
	ID                   string
	UserID               string
	ShopURL              string
	APIKeyID             string
	APIPasswordEncrypted string
	IsDefault            bool
	CreatedAt            string
	UpdatedAt            string
}

type ShopRepository struct {
	db *pgxpool.Pool
}

func NewShopRepository(db *pgxpool.Pool) *ShopRepository {
	return &ShopRepository{db: db}
}

func (r *ShopRepository) Create(userID, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault bool) (*Shop, error) {
	var s Shop
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO shops (user_id, shop_url, api_key_id, api_password_encrypted, is_default)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, shop_url, api_key_id, is_default, created_at::text`,
		userID, shopURL, apiKeyID, apiPasswordEncrypted, isDefault,
	).Scan(&s.ID, &s.UserID, &s.ShopURL, &s.APIKeyID, &s.IsDefault, &s.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create shop: %w", err)
	}
	return &s, nil
}

func (r *ShopRepository) GetByID(id string) (*Shop, error) {
	var s Shop
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, shop_url, api_key_id, api_password_encrypted, is_default, created_at::text, updated_at::text
		 FROM shops WHERE id = $1`, id,
	).Scan(&s.ID, &s.UserID, &s.ShopURL, &s.APIKeyID, &s.APIPasswordEncrypted, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	return &s, nil
}

func (r *ShopRepository) ListByUserID(userID string) ([]Shop, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, shop_url, api_key_id, is_default, created_at::text, updated_at::text
		 FROM shops WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list shops: %w", err)
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var s Shop
		if err := rows.Scan(&s.ID, &s.UserID, &s.ShopURL, &s.APIKeyID, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shop: %w", err)
		}
		shops = append(shops, s)
	}
	return shops, nil
}

func (r *ShopRepository) GetDefault(userID string) (*Shop, error) {
	var s Shop
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, shop_url, api_key_id, api_password_encrypted, is_default, created_at::text, updated_at::text
		 FROM shops WHERE user_id = $1 AND is_default = TRUE`, userID,
	).Scan(&s.ID, &s.UserID, &s.ShopURL, &s.APIKeyID, &s.APIPasswordEncrypted, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("default shop not found: %w", err)
	}
	return &s, nil
}

func (r *ShopRepository) Update(id, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault bool) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE shops SET shop_url=$1, api_key_id=$2, api_password_encrypted=$3, is_default=$4, updated_at=NOW()
		 WHERE id=$5`, shopURL, apiKeyID, apiPasswordEncrypted, isDefault, id,
	)
	if err != nil {
		return fmt.Errorf("update shop: %w", err)
	}
	return nil
}

func (r *ShopRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM shops WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete shop: %w", err)
	}
	return nil
}

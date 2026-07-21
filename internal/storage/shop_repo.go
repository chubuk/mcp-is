package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Shop struct {
	ID                   string
	UserID               string
	Name                 string
	ShopURL              string
	APIKeyID             string
	APIPasswordEncrypted string
	IsDefault            bool
	IsReadOnly           bool
	Jurisdiction         string
	PDNEnabled           bool
	CreatedAt            string
	UpdatedAt            string
}

type ShopRepository struct {
	db *pgxpool.Pool
}

func NewShopRepository(db *pgxpool.Pool) *ShopRepository {
	return &ShopRepository{db: db}
}

func (r *ShopRepository) Create(userID, name, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault, isReadOnly bool) (*Shop, error) {
	var s Shop
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO shops (user_id, name, shop_url, api_key_id, api_password_encrypted, is_default, is_read_only)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, user_id, name, shop_url, api_key_id, is_default, is_read_only, COALESCE(jurisdiction,'AUTO'), pdn_enabled, created_at::text`,
		userID, name, shopURL, apiKeyID, apiPasswordEncrypted, isDefault, isReadOnly,
	).Scan(&s.ID, &s.UserID, &s.Name, &s.ShopURL, &s.APIKeyID, &s.IsDefault, &s.IsReadOnly, &s.Jurisdiction, &s.PDNEnabled, &s.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create shop: %w", err)
	}
	return &s, nil
}

func (r *ShopRepository) GetByID(id string) (*Shop, error) {
	var s Shop
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, COALESCE(name,''), shop_url, api_key_id, api_password_encrypted, is_default, is_read_only, COALESCE(jurisdiction,'AUTO'), pdn_enabled, created_at::text, updated_at::text
		 FROM shops WHERE id = $1`, id,
	).Scan(&s.ID, &s.UserID, &s.Name, &s.ShopURL, &s.APIKeyID, &s.APIPasswordEncrypted, &s.IsDefault, &s.IsReadOnly, &s.Jurisdiction, &s.PDNEnabled, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	return &s, nil
}

func (r *ShopRepository) ListByUserID(userID string) ([]Shop, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, COALESCE(name,''), shop_url, api_key_id, is_default, is_read_only, created_at::text, updated_at::text
		 FROM shops WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list shops: %w", err)
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var s Shop
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.ShopURL, &s.APIKeyID, &s.IsDefault, &s.IsReadOnly, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan shop: %w", err)
		}
		shops = append(shops, s)
	}
	return shops, nil
}

func (r *ShopRepository) ListByUserIDFiltered(userID string, page, perPage int) ([]Shop, int, error) {
	offset := (page - 1) * perPage
	if offset < 0 {
		offset = 0
	}

	var total int
	err := r.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM shops WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count shops: %w", err)
	}

	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, COALESCE(name,''), shop_url, api_key_id, is_default, is_read_only, created_at::text, updated_at::text
		 FROM shops WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC LIMIT $2 OFFSET $3`,
		userID, perPage, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list shops filtered: %w", err)
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var s Shop
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.ShopURL, &s.APIKeyID, &s.IsDefault, &s.IsReadOnly, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan shop: %w", err)
		}
		shops = append(shops, s)
	}
	return shops, total, nil
}

func (r *ShopRepository) GetDefault(userID string) (*Shop, error) {
	var s Shop
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, COALESCE(name,''), shop_url, api_key_id, api_password_encrypted, is_default, is_read_only, COALESCE(jurisdiction,'AUTO'), pdn_enabled, created_at::text, updated_at::text
		 FROM shops WHERE user_id = $1 AND is_default = TRUE`, userID,
	).Scan(&s.ID, &s.UserID, &s.Name, &s.ShopURL, &s.APIKeyID, &s.APIPasswordEncrypted, &s.IsDefault, &s.IsReadOnly, &s.Jurisdiction, &s.PDNEnabled, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("default shop not found: %w", err)
	}
	return &s, nil
}

func (r *ShopRepository) Update(id, name, shopURL, apiKeyID, apiPasswordEncrypted string, isDefault, isReadOnly, pdnEnabled bool) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE shops SET name=$1, shop_url=$2, api_key_id=$3, api_password_encrypted=$4, is_default=$5, is_read_only=$6, pdn_enabled=$7, updated_at=NOW()
		 WHERE id=$8`, name, shopURL, apiKeyID, apiPasswordEncrypted, isDefault, isReadOnly, pdnEnabled, id,
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

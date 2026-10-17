package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppInstall struct {
	ID                   string
	InsalesID            string
	ShopHost             string
	APIPasswordEncrypted string
	Status               string // installed | uninstalled | pending_link
	UserID               string
	ShopID               string
	InstalledAt          string
	UninstalledAt        string
	UpdatedAt            string
}

type AppInstallRepository struct {
	db *pgxpool.Pool
}

func NewAppInstallRepository(db *pgxpool.Pool) *AppInstallRepository {
	return &AppInstallRepository{db: db}
}

const appInstallCols = `id, insales_id, shop_host, api_password_encrypted, status,
	COALESCE(user_id::text,''), COALESCE(shop_id::text,''), installed_at::text,
	COALESCE(uninstalled_at::text,''), updated_at::text`

// UpsertInstall stores the per-shop app password from an install (or reinstall) callback.
func (r *AppInstallRepository) UpsertInstall(insalesID, shopHost, apiPasswordEncrypted string) (*AppInstall, error) {
	var a AppInstall
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO app_installs (insales_id, shop_host, api_password_encrypted, status)
		 VALUES ($1, $2, $3, 'installed')
		 ON CONFLICT (insales_id) DO UPDATE SET
		   shop_host = EXCLUDED.shop_host,
		   api_password_encrypted = EXCLUDED.api_password_encrypted,
		   status = 'installed',
		   uninstalled_at = NULL,
		   updated_at = NOW()
		 RETURNING `+appInstallCols,
		insalesID, shopHost, apiPasswordEncrypted,
	).Scan(&a.ID, &a.InsalesID, &a.ShopHost, &a.APIPasswordEncrypted, &a.Status,
		&a.UserID, &a.ShopID, &a.InstalledAt, &a.UninstalledAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert app install: %w", err)
	}
	return &a, nil
}

func (r *AppInstallRepository) GetByInsalesID(insalesID string) (*AppInstall, error) {
	var a AppInstall
	err := r.db.QueryRow(context.Background(),
		`SELECT `+appInstallCols+` FROM app_installs WHERE insales_id = $1`, insalesID,
	).Scan(&a.ID, &a.InsalesID, &a.ShopHost, &a.APIPasswordEncrypted, &a.Status,
		&a.UserID, &a.ShopID, &a.InstalledAt, &a.UninstalledAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("app install not found: %w", err)
	}
	return &a, nil
}

// MarkUninstalled flags the install and clears the shop link (shop row deleted by caller).
func (r *AppInstallRepository) MarkUninstalled(insalesID string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE app_installs SET status = 'uninstalled', uninstalled_at = NOW(), shop_id = NULL, updated_at = NOW()
		 WHERE insales_id = $1`, insalesID)
	if err != nil {
		return fmt.Errorf("mark uninstalled: %w", err)
	}
	return nil
}

// SetPendingLink records that autologin matched an existing user — link completes after login.
func (r *AppInstallRepository) SetPendingLink(insalesID, userID string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE app_installs SET status = 'pending_link', user_id = $2, updated_at = NOW()
		 WHERE insales_id = $1`, insalesID, userID)
	if err != nil {
		return fmt.Errorf("set pending link: %w", err)
	}
	return nil
}

// LinkShop binds the install to a provisioned user + shop.
func (r *AppInstallRepository) LinkShop(insalesID, userID, shopID string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE app_installs SET status = 'installed', user_id = $2, shop_id = $3, updated_at = NOW()
		 WHERE insales_id = $1`, insalesID, userID, shopID)
	if err != nil {
		return fmt.Errorf("link shop: %w", err)
	}
	return nil
}

// SetUser sets the owning user without touching shop binding.
func (r *AppInstallRepository) SetUser(insalesID, userID string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE app_installs SET user_id = $2, updated_at = NOW() WHERE insales_id = $1`, insalesID, userID)
	if err != nil {
		return fmt.Errorf("set user: %w", err)
	}
	return nil
}

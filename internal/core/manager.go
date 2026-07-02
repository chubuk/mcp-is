package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/your-project/mcp-insales/internal/auth"
	"github.com/your-project/mcp-insales/internal/integration/insales"
	"github.com/your-project/mcp-insales/internal/storage"
)

type contextKey string

const (
	ShopReadOnlyKey contextKey = "shop_read_only"
	ShopNameKey     contextKey = "shop_name"
	ShopURLKey      contextKey = "shop_url"
)

type AccountManager struct {
	shopRepo     *storage.ShopRepository
	aes          *auth.AESCipher
	clients      map[string]*insales.InSalesClient
	globalReadOnly bool
}

func NewAccountManager(shopRepo *storage.ShopRepository, aes *auth.AESCipher, globalReadOnly bool) *AccountManager {
	return &AccountManager{
		shopRepo:       shopRepo,
		aes:            aes,
		clients:        make(map[string]*insales.InSalesClient),
		globalReadOnly: globalReadOnly,
	}
}

func (m *AccountManager) GetClient(ctx context.Context) (*insales.InSalesClient, context.Context, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, ctx, errors.New("user_id missing in context")
	}

	shopID, _ := ctx.Value(auth.ShopIDKey).(string)

	var shop *storage.Shop
	var err error

	if shopID != "" {
		shop, err = m.shopRepo.GetByID(shopID)
	} else {
		shop, err = m.shopRepo.GetDefault(userID)
	}

	if err != nil {
		return nil, ctx, fmt.Errorf("shop not found: %w", err)
	}

	shopReadOnly := m.globalReadOnly || shop.IsReadOnly
	ctx = context.WithValue(ctx, ShopReadOnlyKey, shopReadOnly)
	ctx = context.WithValue(ctx, ShopNameKey, shop.Name)
	ctx = context.WithValue(ctx, ShopURLKey, shop.ShopURL)

	cacheKey := shop.ID
	if client, exists := m.clients[cacheKey]; exists {
		return client, ctx, nil
	}

	password, err := m.aes.Decrypt(shop.APIPasswordEncrypted)
	if err != nil {
		return nil, ctx, fmt.Errorf("decrypt password: %w", err)
	}

	client, err := insales.NewClient(shop.ShopURL, shop.APIKeyID, password)
	if err != nil {
		return nil, ctx, fmt.Errorf("create client: %w", err)
	}

	m.clients[cacheKey] = client
	return client, ctx, nil
}

package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/your-project/mcp-insales/internal/auth"
	"github.com/your-project/mcp-insales/internal/integration/insales"
	"github.com/your-project/mcp-insales/internal/storage"
)

type AccountManager struct {
	shopRepo *storage.ShopRepository
	aes      *auth.AESCipher
	clients  map[string]*insales.InSalesClient
}

func NewAccountManager(shopRepo *storage.ShopRepository, aes *auth.AESCipher) *AccountManager {
	return &AccountManager{
		shopRepo: shopRepo,
		aes:      aes,
		clients:  make(map[string]*insales.InSalesClient),
	}
}

func (m *AccountManager) GetClient(ctx context.Context) (*insales.InSalesClient, error) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("user_id missing in context")
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
		return nil, fmt.Errorf("shop not found: %w", err)
	}

	cacheKey := shop.ID
	if client, exists := m.clients[cacheKey]; exists {
		return client, nil
	}

	password, err := m.aes.Decrypt(shop.APIPasswordEncrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt password: %w", err)
	}

	client, err := insales.NewClient(shop.ShopURL, shop.APIKeyID, password)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	m.clients[cacheKey] = client
	return client, nil
}

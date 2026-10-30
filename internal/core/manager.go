package core

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/chubuk/mcp-is/internal/storage"
)

type contextKey string

const (
	ShopReadOnlyKey contextKey = "shop_read_only"
	ShopNameKey     contextKey = "shop_name"
	ShopURLKey      contextKey = "shop_url"
)

// SettingsGetter is the settings-repository slice the billing gate needs
// (defined at consumer; *storage.SettingsRepository satisfies it).
type SettingsGetter interface {
	GetDefault(key, fallback string) string
}

type AccountManager struct {
	shopRepo       *storage.ShopRepository
	aes            *auth.AESCipher
	clients        map[string]*insales.InSalesClient
	globalReadOnly bool
	billingMonthly int // env fallback; server_settings billing_monthly overrides
	settings       SettingsGetter
}

func NewAccountManager(shopRepo *storage.ShopRepository, aes *auth.AESCipher, globalReadOnly bool, billingMonthly int, settings SettingsGetter) *AccountManager {
	return &AccountManager{
		shopRepo:       shopRepo,
		aes:            aes,
		clients:        make(map[string]*insales.InSalesClient),
		globalReadOnly: globalReadOnly,
		billingMonthly: billingMonthly,
		settings:       settings,
	}
}

// billingPrice mirrors AppHandler.billingPrice: settings override, env fallback.
func (m *AccountManager) billingPrice() int {
	if m.settings == nil {
		return m.billingMonthly
	}
	v := m.settings.GetDefault("billing_monthly", strconv.Itoa(m.billingMonthly))
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return m.billingMonthly
	}
	return n
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
		if err == nil && shop.UserID != userID {
			return nil, ctx, fmt.Errorf("shop %s does not belong to user", shopID)
		}
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
	client, exists := m.clients[cacheKey]
	if !exists {
		password, derr := m.aes.Decrypt(shop.APIPasswordEncrypted)
		if derr != nil {
			return nil, ctx, fmt.Errorf("decrypt password: %w", derr)
		}
		client, err = insales.NewClient(shop.ShopURL, shop.APIKeyID, password)
		if err != nil {
			return nil, ctx, fmt.Errorf("create client: %w", err)
		}
		m.clients[cacheKey] = client
	}

	// Billing gate: login is free, MCP data access is not. Lazy TTL refresh —
	// at most one InSales call per shop per TTL window. billingMonthly=0 = free.
	if m.billingPrice() > 0 {
		if berr := CheckBilling(shop, time.Now(),
			func() (map[string]any, error) {
				raw, ferr := client.GetRecurringApplicationCharge(context.Background(), "")
				if ferr != nil {
					return nil, ferr
				}
				charge, _ := raw.(map[string]any)
				return charge, nil
			},
			m.shopRepo.UpdateBillingStatus,
		); berr != nil {
			return nil, ctx, berr
		}
	}

	return client, ctx, nil
}

// InvalidateShop drops the cached API client for a shop (install/uninstall/password rotate).
func (m *AccountManager) InvalidateShop(shopID string) {
	if m == nil {
		return
	}
	delete(m.clients, shopID)
}

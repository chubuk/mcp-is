package core

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/chubuk/mcp-is/internal/storage"
)

// Sentinel billing errors — matched via errors.Is in the MCP tool-error layer.
var (
	ErrBillingExpired = errors.New("billing period expired")
	ErrBillingUnpaid  = errors.New("billing period unpaid")
)

// Billing status values stored in shops.billing_status.
const (
	BillingActive  = "active"
	BillingUnpaid  = "unpaid"
	BillingUnknown = "unknown" // app shop whose charge check failed at login — retry fast
)

// Billing TTLs for lazy refresh. Unpaid, unknown and expired shops recover fast
// (1h); active shops are re-checked once a day per InSales guidance.
const (
	billingTTLActive = 24 * time.Hour
	billingTTLFast   = time.Hour
)

// msck is the InSales timezone — API timestamps are +03:00.
var msck = time.FixedZone("MSK", 3*3600)

// ParseInSalesDeadline parses a charge date field. A full RFC3339 timestamp is
// used as-is; a bare date "2006-01-02" means the shop is accessible the whole
// day, so the deadline is the start of the next day in MSK (exclusive bound —
// no 1-second gap at 23:59:59).
func ParseInSalesDeadline(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	if d, err := time.ParseInLocation("2006-01-02", s, msck); err == nil {
		return d.AddDate(0, 0, 1), true
	}
	return time.Time{}, false
}

// MapCharge maps a recurring_application_charge payload to (status, until, payURL).
// blocked=true wins even when paid_till is in the future: a chargeback or failed
// recurring payment must not leave the period open.
func MapCharge(m map[string]any, now time.Time) (string, *time.Time, string) {
	var until *time.Time
	for _, k := range []string{"paid_till", "trial_expired_at"} {
		s, _ := m[k].(string)
		t, ok := ParseInSalesDeadline(s)
		if !ok {
			continue
		}
		if until == nil || t.After(*until) {
			until = &t
		}
	}

	payURL := ""
	for _, k := range []string{"confirmation_url", "url"} {
		if s, _ := m[k].(string); len(s) > 4 && s[:4] == "http" {
			payURL = s
			break
		}
	}

	if blocked, _ := m["blocked"].(bool); blocked {
		return BillingUnpaid, until, payURL
	}
	if until != nil && !now.Before(*until) {
		return BillingUnpaid, until, payURL
	}
	return BillingActive, until, payURL
}

// BillingTTL returns the lazy-refresh interval. Unpaid/unknown and shops whose
// deadline has passed recover fast — the user may have just paid.
func BillingTTL(status string, until *time.Time, now time.Time) time.Duration {
	if status == BillingUnpaid || status == BillingUnknown {
		return billingTTLFast
	}
	if until != nil && !now.Before(*until) {
		return billingTTLFast
	}
	return billingTTLActive
}

// NeedsBillingRefresh reports whether the cached billing state is stale enough
// to re-check with InSales. Personal shops (no app install) never refresh —
// they have no RecurringApplicationCharge. App installs refresh immediately on
// the first visit, then per TTL.
func NeedsBillingRefresh(shop *storage.Shop, now time.Time) bool {
	if !shop.IsAppInstall {
		return false
	}
	if shop.BillingCheckedAt == nil {
		return true
	}
	return now.Sub(*shop.BillingCheckedAt) > BillingTTL(shop.BillingStatus, shop.BillingUntil, now)
}

// BillingGateError returns the sentinel for a shop's cached billing state, or
// nil when data access is allowed. Runs after a possible refresh.
func BillingGateError(shop *storage.Shop, now time.Time) error {
	if shop.BillingUntil != nil && !now.Before(*shop.BillingUntil) {
		return fmt.Errorf("shop %s: %w", shop.ID, ErrBillingExpired)
	}
	if shop.BillingStatus == BillingUnpaid {
		return fmt.Errorf("shop %s: %w", shop.ID, ErrBillingUnpaid)
	}
	return nil
}

// RefreshBilling re-checks the charge via fetch and persists the result.
// On fetch/write failure the previous cache is kept (fail-open) and only
// billing_checked_at is touched so the retry lands one TTL later, not on
// every call.
func RefreshBilling(shop *storage.Shop, now time.Time,
	fetch func() (map[string]any, error),
	write func(id, status string, until *time.Time, payURL string) error,
) error {
	charge, err := fetch()
	if err != nil {
		// touch checked_at only — keeps old status, spaces retries by TTL
		return write(shop.ID, shop.BillingStatus, shop.BillingUntil, shop.BillingPayURL)
	}
	status, until, payURL := MapCharge(charge, now)
	if err := write(shop.ID, status, until, payURL); err != nil {
		return err
	}
	shop.BillingStatus, shop.BillingUntil, shop.BillingPayURL = status, until, payURL
	checked := now
	shop.BillingCheckedAt = &checked
	return nil
}

// CheckBilling refreshes stale billing state (at most one InSales call per TTL)
// and returns the gate error, or nil when data access is allowed. Fetch/write
// failures are logged and ignored — an InSales outage must not block the shop.
func CheckBilling(shop *storage.Shop, now time.Time,
	fetch func() (map[string]any, error),
	write func(id, status string, until *time.Time, payURL string) error,
) error {
	if NeedsBillingRefresh(shop, now) {
		if err := RefreshBilling(shop, now, fetch, write); err != nil {
			slog.Warn("billing refresh failed, using cached status", "shop_id", shop.ID, "error", err)
		}
	}
	return BillingGateError(shop, now)
}

package core

import (
	"errors"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/storage"
)

func TestParseInSalesDeadline(t *testing.T) {
	// bare date: accessible the whole day → boundary is next-day midnight MSK
	until, ok := ParseInSalesDeadline("2026-10-27")
	if !ok {
		t.Fatal("date not parsed")
	}
	want := time.Date(2026, 10, 28, 0, 0, 0, 0, msck)
	if !until.Equal(want) {
		t.Errorf("date-only boundary = %v, want %v", until, want)
	}

	// full timestamp: used as-is
	until, ok = ParseInSalesDeadline("2026-10-27T15:00:00+03:00")
	if !ok {
		t.Fatal("timestamp not parsed")
	}
	if until.Hour() != 15 {
		t.Errorf("timestamp = %v, want 15:00 as-is", until)
	}

	if _, ok := ParseInSalesDeadline(""); ok {
		t.Error("empty string must not parse")
	}
	if _, ok := ParseInSalesDeadline("garbage"); ok {
		t.Error("garbage must not parse")
	}
}

func TestMapCharge_BlockedWins(t *testing.T) {
	// blocked=true + paid_till in the future → unpaid (critic #2, test 15)
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	status, _, payURL := MapCharge(map[string]any{
		"blocked":          true,
		"paid_till":        "2026-12-01",
		"confirmation_url": "https://x.myinsales.ru/confirm/1",
	}, now)
	if status != BillingUnpaid {
		t.Errorf("blocked+future paid_till = %s, want unpaid", status)
	}
	if payURL != "https://x.myinsales.ru/confirm/1" {
		t.Errorf("payURL = %q", payURL)
	}
}

func TestMapCharge_Active(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	status, until, _ := MapCharge(map[string]any{"paid_till": "2026-12-01"}, now)
	if status != BillingActive {
		t.Errorf("future paid_till = %s, want active", status)
	}
	if until == nil || until.Month() != 12 {
		t.Errorf("until = %v, want Dec boundary", until)
	}

	// trial counts too — max(paid_till, trial_expired_at)
	status, until, _ = MapCharge(map[string]any{
		"paid_till":        "2026-11-01",
		"trial_expired_at": "2026-11-15",
	}, now)
	if status != BillingActive || until == nil || until.Day() != 16 {
		t.Errorf("max() = (%s,%v), want active until Nov 16 boundary", status, until)
	}
}

func TestMapCharge_Expired(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	status, _, _ := MapCharge(map[string]any{"paid_till": "2026-10-25"}, now)
	if status != BillingUnpaid {
		t.Errorf("past paid_till = %s, want unpaid", status)
	}
	// no dates at all → no deadline, treated as active (fail-open)
	status, until, _ := MapCharge(map[string]any{}, now)
	if status != BillingActive || until != nil {
		t.Errorf("empty charge = (%s,%v), want active/nil", status, until)
	}
}

func TestBillingGateError(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	// test 8: now >= billing_until → immediate refusal even with fresh TTL
	err := BillingGateError(&storage.Shop{ID: "s1", BillingStatus: BillingActive, BillingUntil: &past}, now)
	if !errors.Is(err, ErrBillingExpired) {
		t.Errorf("expired until: %v, want ErrBillingExpired", err)
	}

	err = BillingGateError(&storage.Shop{ID: "s1", BillingStatus: BillingUnpaid, BillingUntil: &future}, now)
	if !errors.Is(err, ErrBillingUnpaid) {
		t.Errorf("unpaid: %v, want ErrBillingUnpaid", err)
	}

	// active + future until → allowed
	if err := BillingGateError(&storage.Shop{ID: "s1", BillingStatus: BillingActive, BillingUntil: &future}, now); err != nil {
		t.Errorf("active shop gated: %v", err)
	}
	// unknown → fail-open (charge check failed, InSales may be down)
	if err := BillingGateError(&storage.Shop{ID: "s1", BillingStatus: BillingUnknown}, now); err != nil {
		t.Errorf("unknown shop gated: %v", err)
	}
}

func TestNeedsBillingRefresh(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	recent := now.Add(-time.Minute)
	old25h := now.Add(-25 * time.Hour)
	old2h := now.Add(-2 * time.Hour)

	// personal shop (no app install) → no refresh, ever (critic #6)
	if NeedsBillingRefresh(&storage.Shop{BillingStatus: BillingActive}, now) {
		t.Error("personal shop must not refresh")
	}
	// app install awaiting first check → refresh at once
	if !NeedsBillingRefresh(&storage.Shop{IsAppInstall: true, BillingStatus: BillingUnknown}, now) {
		t.Error("app install must refresh immediately on first visit")
	}
	// test 10: fresh status → 0 InSales calls
	if NeedsBillingRefresh(&storage.Shop{IsAppInstall: true, BillingStatus: BillingActive, BillingCheckedAt: &recent}, now) {
		t.Error("fresh active must not refresh")
	}
	// test 11: stale active (24h TTL) → refresh
	if !NeedsBillingRefresh(&storage.Shop{IsAppInstall: true, BillingStatus: BillingActive, BillingCheckedAt: &old25h}, now) {
		t.Error("stale active must refresh")
	}
	// test 12: unpaid recovers fast (1h TTL)
	if NeedsBillingRefresh(&storage.Shop{IsAppInstall: true, BillingStatus: BillingUnpaid, BillingCheckedAt: &recent}, now) {
		t.Error("fresh unpaid must not refresh")
	}
	if !NeedsBillingRefresh(&storage.Shop{IsAppInstall: true, BillingStatus: BillingUnpaid, BillingCheckedAt: &old2h}, now) {
		t.Error("unpaid after 1h must refresh")
	}
}

func TestRefreshBilling(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	shop := &storage.Shop{ID: "s1", BillingStatus: BillingUnknown}
	shop2 := &storage.Shop{ID: "s2", BillingStatus: BillingActive}

	byID := map[string]*storage.Shop{"s1": shop, "s2": shop2}
	write := func(id, status string, until *time.Time, payURL string) error {
		s := byID[id]
		s.BillingStatus, s.BillingUntil, s.BillingPayURL = status, until, payURL
		c := now
		s.BillingCheckedAt = &c
		return nil
	}

	// test 13: mapping result persisted
	err := RefreshBilling(shop, now,
		func() (map[string]any, error) {
			return map[string]any{"blocked": true, "confirmation_url": "https://pay/1"}, nil
		}, write)
	if err != nil || shop.BillingStatus != BillingUnpaid || shop.BillingPayURL != "https://pay/1" {
		t.Errorf("blocked charge: err=%v status=%s payURL=%s", err, shop.BillingStatus, shop.BillingPayURL)
	}

	// test 14: API error → previous cache kept, checked_at touched (retry spaced)
	err = RefreshBilling(shop2, now,
		func() (map[string]any, error) { return nil, errors.New("insales down") }, write)
	if err != nil || shop2.BillingStatus != BillingActive {
		t.Errorf("fetch error: err=%v status=%s, want keep active", err, shop2.BillingStatus)
	}
	if shop2.BillingCheckedAt == nil {
		t.Error("fetch error must still touch checked_at")
	}
}

func TestBillingTTLBoundary(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Second)
	// expiry approaching / past → fast TTL so payment recovery is quick
	if d := BillingTTL(BillingActive, &past, now); d != billingTTLFast {
		t.Errorf("TTL with past until = %v, want fast", d)
	}
	future := now.Add(48 * time.Hour)
	if d := BillingTTL(BillingActive, &future, now); d != billingTTLActive {
		t.Errorf("TTL active = %v, want 24h", d)
	}
}

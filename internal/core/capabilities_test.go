package core

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestStatusFromHTTP(t *testing.T) {
	cases := []struct {
		code int
		want string
	}{
		{http.StatusOK, CapStatusAllowed},
		{http.StatusForbidden, CapStatusDenied},
		{http.StatusNotFound, CapStatusUnknown},
		{http.StatusInternalServerError, CapStatusUnknown},
		{http.StatusTooManyRequests, CapStatusUnknown},
		{http.StatusUnauthorized, CapStatusUnknown},
	}
	for _, c := range cases {
		if got := StatusFromHTTP(c.code); got != c.want {
			t.Errorf("StatusFromHTTP(%d) = %q, want %q", c.code, got, c.want)
		}
	}
}

func TestParseSnapshot_Empty(t *testing.T) {
	if s := ParseSnapshot(""); s != nil {
		t.Errorf("ParseSnapshot(\"\") = %v, want nil", s)
	}
}

func TestParseSnapshot_Valid(t *testing.T) {
	raw := `{"credentials_ok":true,"domains":{"orders":{"status":"denied","since":"2026-01-01T00:00:00Z","source":"probe"}}}`
	s := ParseSnapshot(raw)
	if s == nil {
		t.Fatal("ParseSnapshot returned nil")
	}
	if !s.CredentialsOK {
		t.Error("CredentialsOK = false, want true")
	}
	if dc := s.Domains["orders"]; dc.Status != CapStatusDenied || dc.Source != CapSourceProbe {
		t.Errorf("orders cap = %+v, want denied/probe", dc)
	}
}

func TestParseSnapshot_Corrupt(t *testing.T) {
	if s := ParseSnapshot("{broken"); s != nil {
		t.Errorf("ParseSnapshot(corrupt) = %v, want nil", s)
	}
}

func TestProbeAll_StatusMapping(t *testing.T) {
	statuses := map[string]int{
		"/account.json":    http.StatusOK,
		"/products.json":   http.StatusForbidden,
		"/orders.json":     http.StatusNotFound,
		"/clients.json":    http.StatusInternalServerError,
		"/categories.json": http.StatusTooManyRequests,
	}
	p := NewProber(func(_ context.Context, endpoint string) (int, error) {
		if code, ok := statuses[endpoint]; ok {
			return code, nil
		}
		return http.StatusOK, nil
	})

	snap := p.ProbeAll(context.Background(), nil)
	if !snap.CredentialsOK {
		t.Error("CredentialsOK = false, want true")
	}
	if got := DomainStatus(snap, "account"); got != CapStatusAllowed {
		t.Errorf("account = %q, want allowed", got)
	}
	if got := DomainStatus(snap, "products"); got != CapStatusDenied {
		t.Errorf("products = %q, want denied", got)
	}
	if dc := snap.Domains["products"]; dc.Source != CapSourceProbe {
		t.Errorf("products source = %q, want probe", dc.Source)
	}
	// 404/5xx/429 → unknown
	for _, d := range []string{"orders", "clients", "categories"} {
		if got := DomainStatus(snap, d); got != CapStatusUnknown {
			t.Errorf("%s = %q, want unknown", d, got)
		}
	}
}

func TestProbeAll_401(t *testing.T) {
	p := NewProber(func(_ context.Context, _ string) (int, error) {
		return http.StatusUnauthorized, nil
	})
	snap := p.ProbeAll(context.Background(), nil)
	if snap.CredentialsOK {
		t.Error("CredentialsOK = true, want false on 401")
	}
	if snap.Error == "" {
		t.Error("Error empty, want 401 message")
	}
}

func TestProbeAll_PreservesOldOnUnknown(t *testing.T) {
	prev := &Snapshot{
		CredentialsOK: true,
		Domains: map[string]DomainCap{
			"orders": {Status: CapStatusAllowed},
		},
	}
	// orders → 500 (unknown), старое allowed должно сохраниться
	p := NewProber(func(_ context.Context, endpoint string) (int, error) {
		if endpoint == "/orders.json" {
			return http.StatusInternalServerError, nil
		}
		return http.StatusOK, nil
	})
	snap := p.ProbeAll(context.Background(), prev)
	if got := DomainStatus(snap, "orders"); got != CapStatusAllowed {
		t.Errorf("orders = %q, want preserved allowed", got)
	}
}

func TestProbeAll_RateLimited(t *testing.T) {
	p := NewProber(func(_ context.Context, _ string) (int, error) {
		return http.StatusTooManyRequests, nil
	})
	snap := p.ProbeAll(context.Background(), nil)
	if snap.Error == "" {
		t.Error("Error empty, want rate-limit message")
	}
}

func TestShouldProbe_CooldownAndSingleflight(t *testing.T) {
	p := NewProber(func(_ context.Context, _ string) (int, error) {
		return http.StatusOK, nil
	})
	if !p.ShouldProbe("shop1") {
		t.Fatal("first ShouldProbe = false, want true")
	}
	p.MarkProbed("shop1")
	if p.ShouldProbe("shop1") {
		t.Error("in-flight ShouldProbe = true, want false")
	}
	p.MarkProbeDone("shop1")
	if p.ShouldProbe("shop1") {
		t.Error("cooldown ShouldProbe = true, want false")
	}
	// другой магазин не блокируется
	if !p.ShouldProbe("shop2") {
		t.Error("other shop ShouldProbe = false, want true")
	}
}

func TestIsStale(t *testing.T) {
	if !IsStale(nil) {
		t.Error("IsStale(nil) = false, want true")
	}
	old := time.Now().Add(-25 * time.Hour)
	if !IsStale(&old) {
		t.Error("IsStale(25h old) = false, want true")
	}
	fresh := time.Now()
	if IsStale(&fresh) {
		t.Error("IsStale(now) = true, want false")
	}
}

func TestMarkDomainDenied403(t *testing.T) {
	snap := MarkDomainDenied403(nil, "orders")
	if got := DomainStatus(snap, "orders"); got != CapStatusDenied {
		t.Errorf("orders = %q, want denied", got)
	}
	if dc := snap.Domains["orders"]; dc.Source != CapSourceRuntime403 {
		t.Errorf("source = %q, want runtime_403", dc.Source)
	}
	if dc := snap.Domains["orders"]; dc.Since == "" {
		t.Error("since empty, want timestamp")
	}
	// prev не мутируется
	if snap2 := MarkDomainDenied403(snap, "products"); DomainStatus(snap2, "products") != CapStatusDenied {
		t.Error("second mark failed")
	}
	if DomainStatus(snap, "products") != CapStatusUnknown {
		t.Error("prev mutated")
	}
}

func TestEffectiveStatus_TwoLayers(t *testing.T) {
	cases := []struct {
		policyOn bool
		cap      string
		want     string
	}{
		{true, CapStatusAllowed, EffEnabled},
		{true, CapStatusDenied, EffUnavailableNoScope},
		{true, CapStatusUnknown, EffUnavailableUnknown},
		{false, CapStatusAllowed, EffOffByPolicy},
		{false, CapStatusDenied, EffOffByPolicy},
		{false, CapStatusUnknown, EffOffByPolicy},
	}
	for _, c := range cases {
		if got := EffectiveStatus(c.policyOn, c.cap); got != c.want {
			t.Errorf("EffectiveStatus(%v, %q) = %q, want %q", c.policyOn, c.cap, got, c.want)
		}
	}
}

func TestProbeAll_ConcurrentSafety(t *testing.T) {
	p := NewProber(func(_ context.Context, _ string) (int, error) {
		time.Sleep(time.Millisecond)
		return http.StatusOK, nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = p.ProbeAll(context.Background(), nil)
		}()
	}
	wg.Wait()
}

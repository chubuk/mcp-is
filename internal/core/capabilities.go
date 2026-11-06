package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Статусы прав домена.
const (
	CapStatusAllowed = "allowed"
	CapStatusDenied  = "denied"
	CapStatusUnknown = "unknown"
)

// Источники пометки denied.
const (
	CapSourceProbe      = "probe"
	CapSourceRuntime403 = "runtime_403"
)

// capabilitiesTTL — устаревание снимка (ленивый re-probe).
const capabilitiesTTL = 24 * time.Hour

// probeCooldown — минимальный интервал автозапусков probe на магазин.
const probeCooldown = 1 * time.Hour

// DomainCap — права одного домена.
type DomainCap struct {
	Status string `json:"status"` // allowed | denied | unknown
	Since  string `json:"since,omitempty"`
	Source string `json:"source,omitempty"` // probe | runtime_403
}

// Snapshot — JSON-снимок прав магазина.
type Snapshot struct {
	CredentialsOK bool                 `json:"credentials_ok"`
	Error         string               `json:"error,omitempty"`
	Domains       map[string]DomainCap `json:"domains"`
}

// DomainProbe — описание probe-эндпоинта домена.
type DomainProbe struct {
	Domain   string
	Endpoint string
}

// ProbeDomains — 25 доменов InSales (GET ?per_page=1).
var ProbeDomains = []DomainProbe{
	{Domain: "account", Endpoint: "/account.json"},
	{Domain: "products", Endpoint: "/products.json"},
	{Domain: "orders", Endpoint: "/orders.json"},
	{Domain: "clients", Endpoint: "/clients.json"},
	{Domain: "categories", Endpoint: "/categories.json"},
	{Domain: "collections", Endpoint: "/collections.json"},
	{Domain: "webhooks", Endpoint: "/webhooks.json"},
	{Domain: "delivery", Endpoint: "/delivery_variants.json"},
	{Domain: "payments", Endpoint: "/payment_gateways.json"},
	{Domain: "pages", Endpoint: "/pages.json"},
	{Domain: "blogs", Endpoint: "/blogs.json"},
	{Domain: "menus", Endpoint: "/menus.json"},
	{Domain: "domains", Endpoint: "/domains.json"},
	{Domain: "files", Endpoint: "/files.json"},
	{Domain: "reviews", Endpoint: "/reviews.json"},
	{Domain: "discounts", Endpoint: "/discount_codes.json"},
	{Domain: "properties", Endpoint: "/properties.json"},
	{Domain: "fields", Endpoint: "/product_fields.json"},
	{Domain: "warehouse", Endpoint: "/warehouses.json"},
	{Domain: "outlets", Endpoint: "/outlets.json"},
	{Domain: "marketplaces", Endpoint: "/marketplaces.json"},
	{Domain: "js_tags", Endpoint: "/js_tags.json"},
	{Domain: "redirects", Endpoint: "/redirects.json"},
	{Domain: "apps", Endpoint: "/o_auth_apps.json"},
	{Domain: "user_tasks", Endpoint: "/user_tasks.json"},
}

// StatusFromHTTP маппит HTTP-статус probe в статус домена.
func StatusFromHTTP(code int) string {
	switch code {
	case http.StatusOK:
		return CapStatusAllowed
	case http.StatusForbidden:
		return CapStatusDenied
	default:
		return CapStatusUnknown
	}
}

// ParseSnapshot разбирает JSON снимка; пустая строка → nil.
func ParseSnapshot(raw string) *Snapshot {
	if raw == "" {
		return nil
	}
	var s Snapshot
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		slog.Warn("parse capabilities snapshot", "error", err)
		return nil
	}
	if s.Domains == nil {
		s.Domains = make(map[string]DomainCap)
	}
	return &s
}

// MarshalSnapshot сериализует снимок; nil → "".
func MarshalSnapshot(s *Snapshot) string {
	if s == nil {
		return ""
	}
	b, err := json.Marshal(s)
	if err != nil {
		slog.Error("marshal capabilities snapshot", "error", err)
		return ""
	}
	return string(b)
}

// NewUnknownSnapshot — все домены unknown (смена кредов / первый запуск).
func NewUnknownSnapshot() *Snapshot {
	return &Snapshot{
		CredentialsOK: true,
		Domains:       make(map[string]DomainCap, len(ProbeDomains)),
	}
}

// ProbeResult — результат пробы одного домена.
type ProbeResult struct {
	Domain     string
	HTTPStatus int
	Err        error
}

// Prober выполняет probe доменов InSales.
type Prober struct {
	// DoGet — функция GET-запроса (внедряется для тестов); возвращает status.
	DoGet func(ctx context.Context, endpoint string) (int, error)

	// mu защищает запуски probe per-shop (singleflight + cooldown).
	mu        sync.Mutex
	lastProbe map[string]time.Time
	inFlight  map[string]bool
}

// NewProber создаёт Prober с дефолтным DoGet (через HTTP-клиент InSales).
func NewProber(doGet func(ctx context.Context, endpoint string) (int, error)) *Prober {
	return &Prober{
		DoGet:     doGet,
		lastProbe: make(map[string]time.Time),
		inFlight:  make(map[string]bool),
	}
}

// ProbeAllProbe запускает probe всех доменов с пулом воркеров.
// При 401 — весь probe credentials_ok:false, домены не сохраняются.
// При 429 — прерывает остаток, старые значения доменов не затирает (merge).
func (p *Prober) ProbeAll(ctx context.Context, prev *Snapshot) *Snapshot {
	const workers = 4

	results := make(chan ProbeResult, len(ProbeDomains))
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	probeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for _, dp := range ProbeDomains {
		wg.Add(1)
		go func(dp DomainProbe) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			status, err := p.DoGet(probeCtx, dp.Endpoint)
			results <- ProbeResult{Domain: dp.Domain, HTTPStatus: status, Err: err}
		}(dp)
	}

	wg.Wait()
	close(results)

	snap := NewUnknownSnapshot()
	rateLimited := false
	badCreds := false

	for r := range results {
		if r.HTTPStatus == http.StatusUnauthorized {
			badCreds = true
		}
		if r.HTTPStatus == http.StatusTooManyRequests {
			rateLimited = true
		}
		if r.HTTPStatus == 0 && r.Err != nil {
			// timeout / network — domain остаётся unknown
			continue
		}

		st := StatusFromHTTP(r.HTTPStatus)
		if st == CapStatusUnknown && prev != nil {
			// не затираем старое значение при частичном провале
			if old, ok := prev.Domains[r.Domain]; ok && old.Status != CapStatusUnknown {
				snap.Domains[r.Domain] = old
				continue
			}
		}

		dc := DomainCap{Status: st}
		if st == CapStatusDenied {
			dc.Since = time.Now().UTC().Format(time.RFC3339)
			dc.Source = CapSourceProbe
		}
		snap.Domains[r.Domain] = dc
	}

	if badCreds {
		snap.CredentialsOK = false
		snap.Error = "неверные учётные данные InSales (401)"
		return snap
	}
	snap.CredentialsOK = true
	if rateLimited {
		snap.Error = "rate limit (429) — часть доменов не проверена"
	}
	return snap
}

// ShouldProbe — можно ли запускать автопrobe сейчас (cooldown + singleflight).
func (p *Prober) ShouldProbe(shopID string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.inFlight[shopID] {
		return false
	}
	if last, ok := p.lastProbe[shopID]; ok && time.Since(last) < probeCooldown {
		return false
	}
	return true
}

// MarkProbed фиксирует время запуска probe.
func (p *Prober) MarkProbed(shopID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastProbe[shopID] = time.Now()
	p.inFlight[shopID] = true
}

// MarkProbeDone снимает флаг in-flight (по завершении).
func (p *Prober) MarkProbeDone(shopID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.inFlight[shopID] = false
}

// IsStale — снимок старше TTL.
func IsStale(checkedAt *time.Time) bool {
	if checkedAt == nil {
		return true
	}
	return time.Since(*checkedAt) > capabilitiesTTL
}

// MarkDomainDenied403 — оптимистичная пометка домена denied при runtime-403.
// Возвращает обновлённый снимок (prev не мутирует).
func MarkDomainDenied403(prev *Snapshot, domain string) *Snapshot {
	snap := &Snapshot{}
	if prev != nil {
		snap.CredentialsOK = prev.CredentialsOK
		snap.Domains = make(map[string]DomainCap, len(prev.Domains))
		for k, v := range prev.Domains {
			snap.Domains[k] = v
		}
	} else {
		snap.Domains = make(map[string]DomainCap)
		snap.CredentialsOK = true
	}
	snap.Domains[domain] = DomainCap{
		Status: CapStatusDenied,
		Since:  time.Now().UTC().Format(time.RFC3339),
		Source: CapSourceRuntime403,
	}
	return snap
}

// ResetForCredsRotation — сброс снимка в unknown при смене кредов.
func ResetForCredsRotation() *Snapshot {
	return NewUnknownSnapshot()
}

// DomainStatus возвращает статус домена из снимка; нет снимка → unknown.
func DomainStatus(snap *Snapshot, domain string) string {
	if snap == nil {
		return CapStatusUnknown
	}
	dc, ok := snap.Domains[domain]
	if !ok {
		return CapStatusUnknown
	}
	return dc.Status
}

// Summary — краткая сводка для UI.
type Summary struct {
	CheckedAt      string   `json:"checked_at"`
	Stale          bool     `json:"stale"`
	CredentialsOK  bool     `json:"credentials_ok"`
	AllowedDomains []string `json:"allowed"`
	DeniedDomains  []string `json:"denied"`
	UnknownDomains []string `json:"unknown"`
	Error          string   `json:"error,omitempty"`
}

// BuildSummary строит Summary из снимка.
func BuildSummary(snap *Snapshot, checkedAt *time.Time) Summary {
	s := Summary{
		CredentialsOK: snap != nil && snap.CredentialsOK,
		Stale:         IsStale(checkedAt),
	}
	if checkedAt != nil {
		s.CheckedAt = checkedAt.UTC().Format(time.RFC3339)
	}
	if snap != nil {
		s.Error = snap.Error
		for _, dp := range ProbeDomains {
			st := DomainStatus(snap, dp.Domain)
			switch st {
			case CapStatusAllowed:
				s.AllowedDomains = append(s.AllowedDomains, dp.Domain)
			case CapStatusDenied:
				s.DeniedDomains = append(s.DeniedDomains, dp.Domain)
			default:
				s.UnknownDomains = append(s.UnknownDomains, dp.Domain)
			}
		}
	} else {
		for _, dp := range ProbeDomains {
			s.UnknownDomains = append(s.UnknownDomains, dp.Domain)
		}
	}
	return s
}

// EffectiveStatus — статус команды (политика × капабилити), см. план 3d.
const (
	EffEnabled            = "enabled"
	EffOffByPolicy        = "off_by_policy"
	EffUnavailableNoScope = "unavailable_no_scope"
	EffUnavailableUnknown = "unavailable_unknown"
)

// EffectiveStatus вычисляет статус команды из двух слоёв.
func EffectiveStatus(policyOn bool, capStatus string) string {
	if !policyOn {
		return EffOffByPolicy
	}
	switch capStatus {
	case CapStatusAllowed:
		return EffEnabled
	case CapStatusDenied:
		return EffUnavailableNoScope
	default:
		return EffUnavailableUnknown
	}
}

// EnsureSnapshot возвращает непустой снимок (fail-open: nil → все unknown).
func EnsureSnapshot(raw string) *Snapshot {
	if s := ParseSnapshot(raw); s != nil {
		return s
	}
	return NewUnknownSnapshot()
}

// String — для отладки.
func (s Snapshot) String() string {
	return fmt.Sprintf("Snapshot{creds_ok=%v, domains=%d, err=%q}", s.CredentialsOK, len(s.Domains), s.Error)
}

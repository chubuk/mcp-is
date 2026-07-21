package pdn

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PDNSettingsCache struct {
	mu      sync.RWMutex
	cache   *PDNSettings
	expires time.Time
	ttl     time.Duration
}

type PDNSettings struct {
	Enabled     bool
	Log         bool
	MaskPII     bool
	AutoProtect bool
}

func NewPDNSettingsCache() *PDNSettingsCache {
	return &PDNSettingsCache{
		ttl: 5 * time.Minute,
	}
}

func (c *PDNSettingsCache) Get(db *pgxpool.Pool) *PDNSettings {
	c.mu.RLock()
	if c.cache != nil && time.Now().Before(c.expires) {
		result := *c.cache
		c.mu.RUnlock()
		return &result
	}
	c.mu.RUnlock()

	settings := c.loadFromDB(db)

	c.mu.Lock()
	c.cache = settings
	c.expires = time.Now().Add(c.ttl)
	c.mu.Unlock()

	return settings
}

func (c *PDNSettingsCache) loadFromDB(db *pgxpool.Pool) *PDNSettings {
	s := &PDNSettings{
		Enabled:     false,
		Log:         true,
		MaskPII:     false,
		AutoProtect: true,
	}

	rows, err := db.Query(context.Background(),
		`SELECT key, value FROM server_settings WHERE key IN ('pdn_control_enabled', 'pdn_log_access', 'pdn_mask_pii', 'pdn_auto_protect')`,
	)
	if err != nil {
		return s
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}
		switch key {
		case "pdn_control_enabled":
			s.Enabled = value == "true"
		case "pdn_log_access":
			s.Log = value == "true"
		case "pdn_mask_pii":
			s.MaskPII = value == "true"
		case "pdn_auto_protect":
			s.AutoProtect = value == "true"
		}
	}

	return s
}

func (c *PDNSettingsCache) Invalidate() {
	c.mu.Lock()
	c.cache = nil
	c.mu.Unlock()
}

func ParseBool(s string) bool {
	return strings.ToLower(s) == "true"
}

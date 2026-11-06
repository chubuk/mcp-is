package config

import "testing"

// Missing/invalid quota env vars fall back to 20.
func TestConfigQuotaDefaults(t *testing.T) {
	t.Setenv("MAX_SHOPS_PER_USER", "")
	t.Setenv("MAX_AI_KEYS_PER_USER", "junk")

	cfg := Load()
	if cfg.MaxShopsPerUser != 20 {
		t.Errorf("MaxShopsPerUser = %d, want 20", cfg.MaxShopsPerUser)
	}
	if cfg.MaxAIKeysPerUser != 50 {
		t.Errorf("MaxAIKeysPerUser = %d, want 50", cfg.MaxAIKeysPerUser)
	}
}

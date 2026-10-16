//go:build integration

package storage

import (
	"context"
	"testing"
)

func TestSettingsStore_GetSet(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewSettingsRepository(pool)

	// Set a value
	err := repo.Set("test_key", "test_value")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get the value
	value, err := repo.Get("test_key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if value != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", value)
	}
}

func TestSettingsStore_GetDefault(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewSettingsRepository(pool)

	// Non-existent key should return fallback
	value := repo.GetDefault("nonexistent", "fallback")
	if value != "fallback" {
		t.Errorf("expected 'fallback', got '%s'", value)
	}

	// Set and get
	repo.Set("existing", "real")
	value = repo.GetDefault("existing", "fallback")
	if value != "real" {
		t.Errorf("expected 'real', got '%s'", value)
	}
}

func TestSettingsStore_Upsert(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewSettingsRepository(pool)

	// First insert
	err := repo.Set("upsert_key", "value1")
	if err != nil {
		t.Fatalf("first Set failed: %v", err)
	}

	// Upsert (update)
	err = repo.Set("upsert_key", "value2")
	if err != nil {
		t.Fatalf("second Set failed: %v", err)
	}

	// Verify updated value
	value, err := repo.Get("upsert_key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if value != "value2" {
		t.Errorf("expected 'value2', got '%s'", value)
	}
}

func TestSettingsStore_SetDefaultsIfMissing(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewSettingsRepository(pool)

	// Set a custom value first
	repo.Set("domain", "custom.domain.com")

	// SetDefaultsIfMissing should NOT overwrite existing
	err := repo.SetDefaultsIfMissing()
	if err != nil {
		t.Fatalf("SetDefaultsIfMissing failed: %v", err)
	}

	val, _ := repo.Get("domain")
	if val != "custom.domain.com" {
		t.Errorf("expected 'custom.domain.com', got '%s'", val)
	}

	// But missing keys should be set
	val = repo.GetDefault("https_enabled", "missing")
	if val == "missing" {
		t.Error("expected https_enabled to be set by defaults")
	}
}

func TestMCPClientTypeRepo_ListActive(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewMCPClientTypeRepository(pool)

	// Seeds are applied by Seed(), not migrations
	if err := repo.Seed(); err != nil {
		t.Fatalf("Seed failed: %v", err)
	}
	types, err := repo.ListActive()
	if err != nil {
		t.Fatalf("ListActive failed: %v", err)
	}

	// At least Claude Desktop should exist
	found := false
	for _, ct := range types {
		if ct.Name == "Claude Desktop" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'Claude Desktop' in client types")
	}
}

func TestMCPClientTypeRepo_Seed(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewMCPClientTypeRepository(pool)

	// Seed should be idempotent
	err := repo.Seed()
	if err != nil {
		t.Fatalf("Seed failed: %v", err)
	}

	// After seed, should have types
	types, err := repo.ListActive()
	if err != nil {
		t.Fatalf("ListActive after Seed failed: %v", err)
	}
	if len(types) == 0 {
		t.Error("expected client types after seed, got 0")
	}
}

func TestMCPClientTypeRepo_TemplatesContainKeyName(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewMCPClientTypeRepository(pool)

	// Seed templates
	err := repo.Seed()
	if err != nil {
		t.Fatalf("Seed failed: %v", err)
	}

	types, err := repo.ListActive()
	if err != nil {
		t.Fatalf("ListActive failed: %v", err)
	}

	for _, ct := range types {
		if ct.ConfigFormat == "json" {
			if !contains(ct.ConfigTemplate, "{{KEY_NAME}}") {
				t.Errorf("template for %q does not contain {{KEY_NAME}}: %s", ct.Name, ct.ConfigTemplate[:80])
			}
		}
		if ct.ConfigFormat == "yaml" {
			if !contains(ct.ConfigTemplate, "{{KEY_NAME}}") {
				t.Errorf("YAML template for %q does not contain {{KEY_NAME}}", ct.Name)
			}
		}
		// Should NOT contain hardcoded "mcp-is" as JSON key
		if contains(ct.ConfigTemplate, `"mcp-is"`) {
			t.Errorf("template for %q still contains hardcoded \"mcp-is\"", ct.Name)
		}
	}
}

func TestMCPClientTypeRepo_SeedUpsertsStaleTemplates(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewMCPClientTypeRepository(pool)

	// Seed once to create rows, then poison a real seed row with old "mcp-is" format
	if err := repo.Seed(); err != nil {
		t.Fatalf("Seed failed: %v", err)
	}
	types, err := repo.ListActive()
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if len(types) == 0 {
		t.Fatal("no client types after Seed")
	}
	seedName := types[0].Name

	_, err = pool.Exec(context.Background(),
		`UPDATE mcp_client_types SET config_template = '{"mcpServers":{"mcp-is":{"url":"{{URL}}"}}}' WHERE name = $1`, seedName)
	if err != nil {
		t.Fatalf("poison template: %v", err)
	}

	// Seed() must UPSERT: restore template from code even though the row exists
	if err := repo.Seed(); err != nil {
		t.Fatalf("second Seed failed: %v", err)
	}

	types, err = repo.ListActive()
	if err != nil {
		t.Fatalf("ListActive after re-seed: %v", err)
	}
	for _, ct := range types {
		if ct.Name != seedName {
			continue
		}
		if ct.ConfigFormat == "json" && !contains(ct.ConfigTemplate, "{{KEY_NAME}}") {
			t.Errorf("Seed did not upsert template: %s", ct.ConfigTemplate)
		}
		if contains(ct.ConfigTemplate, `"mcp-is"`) {
			t.Errorf("template still has hardcoded mcp-is after Seed")
		}
		return
	}
	t.Errorf("seed row %q not found after re-seed", seedName)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

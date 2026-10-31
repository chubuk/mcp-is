//go:build integration

package storage

import (
	"testing"
)

func TestAIKeyRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("keyuser", "", "", "", "", "hash")

	key, err := keyRepo.Create(user.ID, "keyhash123", "Cursor MCP")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if key.ID == "" {
		t.Error("expected non-empty ID")
	}
	if key.KeyHash != "keyhash123" {
		t.Errorf("expected key_hash 'keyhash123', got '%s'", key.KeyHash)
	}
}

func TestAIKeyRepo_GetByHash(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("hashuser", "", "", "", "", "hash")
	keyRepo.Create(user.ID, "myhash", "test")

	found, err := keyRepo.GetByHash("myhash")
	if err != nil {
		t.Fatalf("GetByHash failed: %v", err)
	}
	if found.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", found.Name)
	}
}

func TestAIKeyRepo_GetByHash_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	keyRepo := NewAIKeyRepository(pool)

	_, err := keyRepo.GetByHash("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestAIKeyRepo_ListByUserID(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("listkey", "", "", "", "", "hash")
	keyRepo.Create(user.ID, "h1", "key1")
	keyRepo.Create(user.ID, "h2", "key2")

	keys, err := keyRepo.ListByUserID(user.ID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestAIKeyRepo_Delete(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("delkey", "", "", "", "", "hash")
	key, _ := keyRepo.Create(user.ID, "delhash", "to-delete")

	err := keyRepo.Delete(key.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = keyRepo.GetByHash("delhash")
	if err == nil {
		t.Error("expected error for deleted key")
	}
}

func TestAIKeyRepo_UpdateLastUsed(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("lastused", "", "", "", "", "hash")
	key, _ := keyRepo.Create(user.ID, "luhash", "test")

	err := keyRepo.UpdateLastUsed(key.ID)
	if err != nil {
		t.Fatalf("UpdateLastUsed failed: %v", err)
	}

	found, _ := keyRepo.GetByHash("luhash")
	if found.LastUsedAt == "" {
		t.Error("expected last_used_at to be set")
	}
}

func TestAIKeyRepo_ListByUserIDFiltered(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("filtkey", "", "", "", "", "hash")
	keyRepo.Create(user.ID, "h1", "Cursor")
	keyRepo.Create(user.ID, "h2", "Windsurf")
	keyRepo.Create(user.ID, "h3", "Roo Code")

	keys, total, err := keyRepo.ListByUserIDFiltered(user.ID, 1, 2)
	if err != nil {
		t.Fatalf("ListByUserIDFiltered failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys on page 1, got %d", len(keys))
	}

	keys2, _, _ := keyRepo.ListByUserIDFiltered(user.ID, 2, 2)
	if len(keys2) != 1 {
		t.Errorf("expected 1 key on page 2, got %d", len(keys2))
	}
}

func TestAIKeyRepo_UpdateConnectionName(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("connname", "", "", "", "", "hash")
	key, _ := keyRepo.Create(user.ID, "connhash", "test")

	if err := keyRepo.UpdateConnectionName(key.ID, "my-shop"); err != nil {
		t.Fatalf("UpdateConnectionName failed: %v", err)
	}

	found, err := keyRepo.GetByHash("connhash")
	if err != nil {
		t.Fatalf("GetByHash failed: %v", err)
	}
	if got, want := found.ConnectionName, "my-shop"; got != want {
		t.Errorf("ConnectionName = %q, want %q", got, want)
	}

	// empty string clears back to NULL
	if err := keyRepo.UpdateConnectionName(key.ID, ""); err != nil {
		t.Fatalf("UpdateConnectionName clear failed: %v", err)
	}
	found, _ = keyRepo.GetByHash("connhash")
	if got, want := found.ConnectionName, ""; got != want {
		t.Errorf("ConnectionName after clear = %q, want %q", got, want)
	}
}

func TestAIKeyRepo_UpdateShopBinding(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("binduser", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "Bind Shop", "bind.myinsales.ru", "ik", "ip", true, false, false)
	key, _ := keyRepo.Create(user.ID, "bindhash", "test")

	if err := keyRepo.UpdateShopBinding(key.ID, shop.ID, "bind-alias"); err != nil {
		t.Fatalf("UpdateShopBinding failed: %v", err)
	}

	found, err := keyRepo.GetByID(key.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got, want := found.ShopID, shop.ID; got != want {
		t.Errorf("ShopID = %q, want %q", got, want)
	}
	if got, want := found.ConnectionName, "bind-alias"; got != want {
		t.Errorf("ConnectionName = %q, want %q", got, want)
	}
}

func TestAIKeyRepo_LegacyNullConnectionName(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("legacyconn", "", "", "", "", "hash")
	key, _ := keyRepo.Create(user.ID, "legacyhash", "legacy")

	// simulate pre-migration row: connection_name left NULL
	if _, err := pool.Exec(t.Context(),
		`UPDATE ai_keys SET connection_name = NULL WHERE id = $1`, key.ID); err != nil {
		t.Fatalf("force NULL connection_name: %v", err)
	}

	found, err := keyRepo.GetByHash("legacyhash")
	if err != nil {
		t.Fatalf("GetByHash on legacy row: %v", err)
	}
	if got, want := found.ConnectionName, ""; got != want {
		t.Errorf("ConnectionName = %q, want %q", got, want)
	}

	keys, err := keyRepo.ListByUserID(user.ID)
	if err != nil {
		t.Fatalf("ListByUserID on legacy row: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	if got, want := keys[0].ConnectionName, ""; got != want {
		t.Errorf("ListByUserID ConnectionName = %q, want %q", got, want)
	}
}

func TestAIKeyRepo_ShopNameJoin(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, _ := userRepo.Create("shopname", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "My Shop", "name.myinsales.ru", "ik", "ip", true, false, false)
	key, _ := keyRepo.CreateForShop(user.ID, shop.ID, "namehash", "test")

	tests := []struct {
		name string
		get  func() (string, string, error)
	}{
		{"GetByHash", func() (string, string, error) {
			k, err := keyRepo.GetByHash("namehash")
			if err != nil {
				return "", "", err
			}
			return k.ShopName, k.ShopID, nil
		}},
		{"GetByID", func() (string, string, error) {
			k, err := keyRepo.GetByID(key.ID)
			if err != nil {
				return "", "", err
			}
			return k.ShopName, k.ShopID, nil
		}},
		{"ListByUserID", func() (string, string, error) {
			ks, err := keyRepo.ListByUserID(user.ID)
			if err != nil {
				return "", "", err
			}
			return ks[0].ShopName, ks[0].ShopID, nil
		}},
		{"ListByUserIDFiltered", func() (string, string, error) {
			ks, _, err := keyRepo.ListByUserIDFiltered(user.ID, 1, 10)
			if err != nil {
				return "", "", err
			}
			return ks[0].ShopName, ks[0].ShopID, nil
		}},
	}

	for _, tt := range tests {
		shopName, shopID, err := tt.get()
		if err != nil {
			t.Fatalf("%s failed: %v", tt.name, err)
		}
		if got, want := shopName, "My Shop"; got != want {
			t.Errorf("%s ShopName = %q, want %q", tt.name, got, want)
		}
		if got, want := shopID, shop.ID; got != want {
			t.Errorf("%s ShopID = %q, want %q", tt.name, got, want)
		}
	}

	// unbound key: JOIN yields empty shop name, no crash
	if _, err := keyRepo.Create(user.ID, "unboundhash", "unbound"); err != nil {
		t.Fatalf("Create unbound: %v", err)
	}
	unbound, err := keyRepo.GetByHash("unboundhash")
	if err != nil {
		t.Fatalf("GetByHash unbound: %v", err)
	}
	if got, want := unbound.ShopName, ""; got != want {
		t.Errorf("unbound ShopName = %q, want %q", got, want)
	}
}

// connection_name is unique per user at the DB level (partial unique index).
func TestAIKeyRepo_ConnectionNameUniquePerUser(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	u1, err := userRepo.Create("uniq1", "", "", "", "", "hash")
	if err != nil {
		t.Fatal(err)
	}
	u2, err := userRepo.Create("uniq2", "", "", "", "", "hash")
	if err != nil {
		t.Fatal(err)
	}
	k1, _ := keyRepo.Create(u1.ID, "uhash1", "k1")
	k2, _ := keyRepo.Create(u1.ID, "uhash2", "k2")
	k3, _ := keyRepo.Create(u2.ID, "uhash3", "k3")

	if err := keyRepo.UpdateConnectionName(k1.ID, "alias"); err != nil {
		t.Fatalf("first store: %v", err)
	}
	if err := keyRepo.UpdateConnectionName(k2.ID, "alias"); err == nil {
		t.Error("duplicate connection_name within one user must fail at DB level")
	}
	if err := keyRepo.UpdateConnectionName(k3.ID, "alias"); err != nil {
		t.Errorf("same alias for another user must be allowed: %v", err)
	}
}

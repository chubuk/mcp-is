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

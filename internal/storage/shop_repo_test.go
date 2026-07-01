//go:build integration

package storage

import (
	"testing"
)

func TestShopRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("shopuser", "", "", "", "", "hash")

	shop, err := shopRepo.Create(user.ID, "myshop.myinsales.ru", "key123", "encrypted_pass", true)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if shop.ID == "" {
		t.Error("expected non-empty ID")
	}
	if shop.ShopURL != "myshop.myinsales.ru" {
		t.Errorf("expected shop_url 'myshop.myinsales.ru', got '%s'", shop.ShopURL)
	}
	if !shop.IsDefault {
		t.Error("expected is_default=true")
	}
}

func TestShopRepo_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("getshop", "", "", "", "", "hash")
	created, _ := shopRepo.Create(user.ID, "shop.ru", "k", "enc", false)

	found, err := shopRepo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found.ShopURL != "shop.ru" {
		t.Errorf("expected 'shop.ru', got '%s'", found.ShopURL)
	}
}

func TestShopRepo_ListByUserID(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("listshop", "", "", "", "", "hash")
	shopRepo.Create(user.ID, "s1.ru", "k1", "e1", true)
	shopRepo.Create(user.ID, "s2.ru", "k2", "e2", false)

	shops, err := shopRepo.ListByUserID(user.ID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(shops) != 2 {
		t.Errorf("expected 2 shops, got %d", len(shops))
	}
}

func TestShopRepo_GetDefault(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("defshop", "", "", "", "", "hash")
	shopRepo.Create(user.ID, "other.ru", "k", "e", false)
	shopRepo.Create(user.ID, "default.ru", "k", "e", true)

	shop, err := shopRepo.GetDefault(user.ID)
	if err != nil {
		t.Fatalf("GetDefault failed: %v", err)
	}
	if shop.ShopURL != "default.ru" {
		t.Errorf("expected 'default.ru', got '%s'", shop.ShopURL)
	}
}

func TestShopRepo_Update(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("updshop", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "old.ru", "k", "e", false)

	err := shopRepo.Update(shop.ID, "new.ru", "k2", "e2", true)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := shopRepo.GetByID(shop.ID)
	if updated.ShopURL != "new.ru" {
		t.Errorf("expected 'new.ru', got '%s'", updated.ShopURL)
	}
	if !updated.IsDefault {
		t.Error("expected is_default=true after update")
	}
}

func TestShopRepo_Delete(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("delshop", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "del.ru", "k", "e", false)

	err := shopRepo.Delete(shop.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = shopRepo.GetByID(shop.ID)
	if err == nil {
		t.Error("expected error for deleted shop")
	}
}

func TestShopRepo_ListByUserIDFiltered(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("filtshop", "", "", "", "", "hash")
	shopRepo.Create(user.ID, "s1.ru", "k1", "e1", true, false)
	shopRepo.Create(user.ID, "s2.ru", "k2", "e2", false, false)
	shopRepo.Create(user.ID, "s3.ru", "k3", "e3", false, true)

	shops, total, err := shopRepo.ListByUserIDFiltered(user.ID, 1, 2)
	if err != nil {
		t.Fatalf("ListByUserIDFiltered failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(shops) != 2 {
		t.Errorf("expected 2 shops on page 1, got %d", len(shops))
	}
	if !shops[0].IsDefault {
		t.Error("expected first shop to be default (sorted by is_default DESC)")
	}

	shops2, _, _ := shopRepo.ListByUserIDFiltered(user.ID, 2, 2)
	if len(shops2) != 1 {
		t.Errorf("expected 1 shop on page 2, got %d", len(shops2))
	}
}

func TestShopRepo_IsReadOnly(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("roshop", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "ro.ru", "k", "e", false, true)

	found, _ := shopRepo.GetByID(shop.ID)
	if !found.IsReadOnly {
		t.Error("expected is_read_only=true")
	}
}

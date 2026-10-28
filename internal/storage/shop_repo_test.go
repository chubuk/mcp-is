//go:build integration

package storage

import (
	"testing"
	"time"
)

func TestShopRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("shopuser", "", "", "", "", "hash")

	shop, err := shopRepo.Create(user.ID, "myshop.myinsales.ru", "myshop.myinsales.ru", "key123", "encrypted_pass", true, false, false)
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
	created, _ := shopRepo.Create(user.ID, "shop.ru", "shop.ru", "k", "enc", false, false, false)

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
	shopRepo.Create(user.ID, "s1.ru", "s1.ru", "k1", "e1", true, false, false)
	shopRepo.Create(user.ID, "s2.ru", "s2.ru", "k2", "e2", false, false, false)

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
	shopRepo.Create(user.ID, "other.ru", "other.ru", "k", "e", false, false, false)
	shopRepo.Create(user.ID, "default.ru", "default.ru", "k", "e", true, false, false)

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
	shop, _ := shopRepo.Create(user.ID, "old.ru", "old.ru", "k", "e", false, false, false)

	err := shopRepo.Update(shop.ID, "new.ru", "new.ru", "k2", "e2", true, false, false)
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

func TestShopRepo_UpdateAPIPassword(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("updpass", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "p.ru", "p.ru", "k", "old_enc", false, false, false)

	if err := shopRepo.UpdateAPIPassword(shop.ID, "new_enc"); err != nil {
		t.Fatalf("UpdateAPIPassword failed: %v", err)
	}
	updated, _ := shopRepo.GetByID(shop.ID)
	if updated.APIPasswordEncrypted != "new_enc" {
		t.Errorf("api_password_encrypted = %q, want 'new_enc'", updated.APIPasswordEncrypted)
	}

	// vanished shop → error, not a silent 0-row update
	if err := shopRepo.UpdateAPIPassword("no-such-id", "x"); err == nil {
		t.Error("expected error for unknown shop id")
	}
}

func TestShopRepo_Delete(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("delshop", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "del.ru", "del.ru", "k", "e", false, false, false)

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
	shopRepo.Create(user.ID, "s1.ru", "s1.ru", "k1", "e1", true, false, false)
	shopRepo.Create(user.ID, "s2.ru", "s2.ru", "k2", "e2", false, false, false)
	shopRepo.Create(user.ID, "s3.ru", "s3.ru", "k3", "e3", false, true, false)

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
	shop, _ := shopRepo.Create(user.ID, "ro.ru", "ro.ru", "k", "e", false, true, false)

	found, _ := shopRepo.GetByID(shop.ID)
	if !found.IsReadOnly {
		t.Error("expected is_read_only=true")
	}
}

func TestShopRepo_UpdateBillingStatus(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	user, _ := userRepo.Create("billuser", "", "", "", "", "hash")
	shop, _ := shopRepo.Create(user.ID, "bill.ru", "bill.ru", "k", "e", true, false, false)

	// default: active, never checked
	found, _ := shopRepo.GetByID(shop.ID)
	if found.BillingStatus != "active" || found.BillingCheckedAt != nil {
		t.Errorf("default billing = (%s, %v), want (active, nil)", found.BillingStatus, found.BillingCheckedAt)
	}

	until := time.Now().Add(24 * time.Hour)
	if err := shopRepo.UpdateBillingStatus(shop.ID, "unpaid", &until, "https://pay/1"); err != nil {
		t.Fatal(err)
	}
	found, _ = shopRepo.GetByID(shop.ID)
	if found.BillingStatus != "unpaid" || found.BillingUntil == nil || found.BillingCheckedAt == nil || found.BillingPayURL != "https://pay/1" {
		t.Errorf("after update = (%s, %v, %v, %q)", found.BillingStatus, found.BillingUntil, found.BillingCheckedAt, found.BillingPayURL)
	}

	// unknown shop id → error
	if err := shopRepo.UpdateBillingStatus("00000000-0000-0000-0000-000000000000", "active", nil, ""); err == nil {
		t.Error("expected error for missing shop")
	}

	// test 16: uninstall deletes the shop row — billing columns go with it
	if err := shopRepo.Delete(shop.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := shopRepo.GetByID(shop.ID); err == nil {
		t.Error("shop must be gone after delete")
	}
}

func TestShopRepo_IsAppInstallFlag(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)
	installRepo := NewAppInstallRepository(pool)

	user, _ := userRepo.Create("appflag", "", "", "", "", "hash")

	// personal shop: no app_installs row → IsAppInstall=false
	shop, _ := shopRepo.Create(user.ID, "p.ru", "p.ru", "k", "e", true, false, false)
	found, _ := shopRepo.GetByID(shop.ID)
	if found.IsAppInstall {
		t.Error("personal shop must have IsAppInstall=false")
	}

	// app-provisioned shop: installed app_installs row → true
	inst, err := installRepo.UpsertInstall("ins-1", "a.myinsales.ru", "enc")
	if err != nil {
		t.Fatal(err)
	}
	if err := installRepo.LinkShop(inst.InsalesID, user.ID, shop.ID); err != nil {
		t.Fatal(err)
	}
	found, _ = shopRepo.GetByID(shop.ID)
	if !found.IsAppInstall {
		t.Error("app-installed shop must have IsAppInstall=true")
	}

	// uninstalled → flag drops (refresh must stop)
	if err := installRepo.MarkUninstalled("ins-1"); err != nil {
		t.Fatal(err)
	}
	found, _ = shopRepo.GetByID(shop.ID)
	if found.IsAppInstall {
		t.Error("uninstalled app must have IsAppInstall=false")
	}
}

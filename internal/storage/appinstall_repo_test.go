//go:build integration

package storage

import (
	"testing"
)

func TestAppInstallRepo_UpsertAndReload(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewAppInstallRepository(pool)

	a, err := repo.UpsertInstall("101", "shop1.myinsales.ru", "enc-pass-1")
	if err != nil {
		t.Fatalf("UpsertInstall: %v", err)
	}
	if a.Status != "installed" || a.InsalesID != "101" {
		t.Errorf("got %+v", a)
	}

	// reinstall updates password and clears uninstalled_at
	a2, err := repo.UpsertInstall("101", "shop1.myinsales.ru", "enc-pass-2")
	if err != nil {
		t.Fatalf("UpsertInstall reload: %v", err)
	}
	if a2.APIPasswordEncrypted != "enc-pass-2" {
		t.Errorf("password not updated: %s", a2.APIPasswordEncrypted)
	}

	found, err := repo.GetByInsalesID("101")
	if err != nil {
		t.Fatalf("GetByInsalesID: %v", err)
	}
	if found.ShopHost != "shop1.myinsales.ru" {
		t.Errorf("shop_host = %s", found.ShopHost)
	}
}

func TestAppInstallRepo_UninstallAndPendingLink(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewAppInstallRepository(pool)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)

	if _, err := repo.UpsertInstall("202", "shop2.myinsales.ru", "enc"); err != nil {
		t.Fatal(err)
	}
	user, err := userRepo.Create("appuser", "", "", "", "", "hash")
	if err != nil {
		t.Fatal(err)
	}
	shop, err := shopRepo.Create(user.ID, "App Shop", "appshop.myinsales.ru", "key", "pass", true, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.LinkShop("202", user.ID, shop.ID); err != nil {
		t.Fatalf("LinkShop: %v", err)
	}
	got, _ := repo.GetByInsalesID("202")
	if got.ShopID != shop.ID || got.UserID != user.ID {
		t.Errorf("link mismatch: %+v", got)
	}

	if err := repo.MarkUninstalled("202"); err != nil {
		t.Fatalf("MarkUninstalled: %v", err)
	}
	got, _ = repo.GetByInsalesID("202")
	if got.Status != "uninstalled" || got.ShopID != "" {
		t.Errorf("after uninstall: %+v", got)
	}

	if err := repo.SetPendingLink("202", user.ID); err != nil {
		t.Fatalf("SetPendingLink: %v", err)
	}
	got, _ = repo.GetByInsalesID("202")
	if got.Status != "pending_link" {
		t.Errorf("pending_link status = %s", got.Status)
	}
}

// SetPendingLink must never rewrite ownership: COALESCE(NULLIF(user_id,”), $2).
func TestAppInstallRepo_SetPendingLink_KeepsUser(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewAppInstallRepository(pool)
	userRepo := NewUserRepository(pool)

	if _, err := repo.UpsertInstall("303", "shop3.myinsales.ru", "enc"); err != nil {
		t.Fatal(err)
	}
	userA, _ := userRepo.Create("ownA", "", "", "", "", "hash")
	userB, _ := userRepo.Create("ownB", "", "", "", "", "hash")

	if err := repo.SetUser("303", userA.ID); err != nil {
		t.Fatal(err)
	}
	// conflicting link attempt must not steal ownership
	if err := repo.SetPendingLink("303", userB.ID); err != nil {
		t.Fatal(err)
	}
	got, _ := repo.GetByInsalesID("303")
	if got.UserID != userA.ID {
		t.Errorf("user_id rewritten to %s, want %s", got.UserID, userA.ID)
	}
	if got.Status != "pending_link" {
		t.Errorf("status = %s, want pending_link", got.Status)
	}

	// empty user_id → set
	if _, err := repo.UpsertInstall("404", "shop4.myinsales.ru", "enc"); err != nil {
		t.Fatal(err)
	}
	if err := repo.SetPendingLink("404", userA.ID); err != nil {
		t.Fatal(err)
	}
	got, _ = repo.GetByInsalesID("404")
	if got.UserID != userA.ID {
		t.Errorf("empty user_id not set: %+v", got)
	}
}

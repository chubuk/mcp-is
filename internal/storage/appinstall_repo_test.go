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

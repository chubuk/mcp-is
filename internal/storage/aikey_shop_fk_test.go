//go:build integration

package storage

import "testing"

// ai_keys.shop_id must survive shop delete (ON DELETE SET NULL).
func TestAIKeyRepo_ShopDeleteSetsNull(t *testing.T) {
	pool := setupTestDB(t)
	userRepo := NewUserRepository(pool)
	shopRepo := NewShopRepository(pool)
	keyRepo := NewAIKeyRepository(pool)

	user, err := userRepo.Create("fkuser", "", "", "", "", "hash")
	if err != nil {
		t.Fatal(err)
	}
	shop, err := shopRepo.Create(user.ID, "FK Shop", "fk.myinsales.ru", "ik", "ip", true, false)
	if err != nil {
		t.Fatal(err)
	}
	key, err := keyRepo.CreateForShop(user.ID, shop.ID, "fkhash", "FK key")
	if err != nil {
		t.Fatal(err)
	}
	if key.ShopID != shop.ID {
		t.Fatalf("CreateForShop shop_id = %q, want %q", key.ShopID, shop.ID)
	}

	if err := shopRepo.Delete(shop.ID); err != nil {
		t.Fatalf("Delete shop: %v", err)
	}

	found, err := keyRepo.GetByHash("fkhash")
	if err != nil {
		t.Fatalf("key must survive shop delete: %v", err)
	}
	if found.ShopID != "" {
		t.Errorf("shop_id should be NULL after shop delete, got %q", found.ShopID)
	}
}

func TestAuditRepo_InsertAndList(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewAuditRepository(pool)

	err := repo.Insert(t.Context(), AuditLogEntry{
		UserID: "u1", ToolName: "create_product", Action: "call",
		ArgsHash: "abc", ResultOK: true,
	})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}
	list, err := repo.ListRecent(t.Context(), 10)
	if err != nil {
		t.Fatalf("ListRecent: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("empty audit list")
	}
	if list[0].ToolName != "create_product" {
		t.Errorf("tool = %s", list[0].ToolName)
	}
}

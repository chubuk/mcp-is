package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupAIKeyHandler() (*AIKeyHandler, *mockAIKeyRepo, *mockShopRepo) {
	keys := newMockAIKeyRepo()
	shops := newMockShopRepo()
	return NewAIKeyHandler(keys, shops, 20), keys, shops
}

func ptr(s string) *string { return &s }

func boolPtr(b bool) *bool { return &b }

func createAIKey(t *testing.T, h *AIKeyHandler, userID, shopID string) AIKeyResponse {
	t.Helper()
	body, _ := json.Marshal(AIKeyRequest{Name: ptr("Cursor MCP"), ShopID: ptr(shopID)})
	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), userID))
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create key: got %d: %s", w.Code, w.Body.String())
	}
	var resp AIKeyResponse
	json.NewDecoder(w.Body).Decode(&resp)
	return resp
}

func TestAIKeyCreate_Success(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "romashka.myinsales.ru", "ik", "ip", true, false, false)

	resp := createAIKey(t, h, "user-1", shop.ID)

	if resp.Key == "" {
		t.Error("expected raw key in response (shown once)")
	}
	if resp.ShopID != shop.ID {
		t.Errorf("shop_id = %q, want %q", resp.ShopID, shop.ID)
	}
	if resp.ShopName != "Ромашка" {
		t.Errorf("shop_name = %q, want Ромашка", resp.ShopName)
	}
	// default connection_name = slug(shop name)
	if resp.ConnectionName != "romashka" {
		t.Errorf("connection_name = %q, want romashka", resp.ConnectionName)
	}
}

func TestAIKeyCreate_NoShopID(t *testing.T) {
	h, _, shops := setupAIKeyHandler()

	body, _ := json.Marshal(AIKeyRequest{Name: ptr("k")})
	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want 400", w.Code)
	}

	// zero shops → "Сначала добавьте магазин"
	if !strings.Contains(w.Body.String(), "Сначала добавьте магазин") {
		t.Errorf("zero shops should hint to add a shop first: %s", w.Body.String())
	}

	// with an existing shop → generic "Укажите магазин"
	shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req2 = req2.WithContext(withUserID(req2.Context(), "user-1"))
	h.Create(w2, req2)
	if !strings.Contains(w2.Body.String(), "Укажите магазин") {
		t.Errorf("want generic shop_id hint: %s", w2.Body.String())
	}
}

func TestAIKeyCreate_ForeignShop(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-2", "Чужой", "foreign.myinsales.ru", "ik", "ip", true, false, false)

	body, _ := json.Marshal(AIKeyRequest{Name: ptr("k"), ShopID: ptr(shop.ID)})
	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want 400", w.Code)
	}
}

func TestAIKeyCreate_InvalidConnectionName(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)

	body, _ := json.Marshal(AIKeyRequest{Name: ptr("k"), ShopID: ptr(shop.ID), ConnectionName: ptr("Bad_Name!")})
	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want 400 for invalid connection_name", w.Code)
	}
}

func TestAIKeyCreate_ConnectionNameCollisionGetsSuffix(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)

	first := createAIKey(t, h, "user-1", shop.ID)
	second := createAIKey(t, h, "user-1", shop.ID)

	if first.ConnectionName != "romashka" {
		t.Errorf("first = %q, want romashka", first.ConnectionName)
	}
	if second.ConnectionName != "romashka-2" {
		t.Errorf("second = %q, want romashka-2", second.ConnectionName)
	}
}

func TestAIKeyList_IncludesShopFields(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	createAIKey(t, h, "user-1", shop.ID)

	req := httptest.NewRequest(http.MethodGet, "/api/ai-keys", nil)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d", w.Code)
	}
	var result map[string]any
	json.NewDecoder(w.Body).Decode(&result)
	items := result["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 key, got %d", len(items))
	}
	item := items[0].(map[string]any)
	if item["shop_id"] != shop.ID {
		t.Errorf("item shop_id = %v, want %s", item["shop_id"], shop.ID)
	}
	if item["connection_name"] != "romashka" {
		t.Errorf("item connection_name = %v, want romashka", item["connection_name"])
	}
}

func TestAIKeyUpdate_UnboundRequiresShopID(t *testing.T) {
	h, keys, _ := setupAIKeyHandler()
	key, _ := keys.Create("user-1", "hash", "legacy") // unbound legacy key

	// any change without shop_id → 400
	body, _ := json.Marshal(AIKeyRequest{PDNEnabled: boolPtr(true)})
	req := httptest.NewRequest(http.MethodPut, "/api/ai-keys/"+key.ID, bytes.NewReader(body))
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("unbound PUT without shop_id: got %d, want 400", w.Code)
	}
}

func TestAIKeyUpdate_BindLegacyKey(t *testing.T) {
	h, keys, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	key, _ := keys.Create("user-1", "hash", "legacy")

	body, _ := json.Marshal(AIKeyRequest{ShopID: ptr(shop.ID)})
	req := httptest.NewRequest(http.MethodPut, "/api/ai-keys/"+key.ID, bytes.NewReader(body))
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d: %s", w.Code, w.Body.String())
	}
	var resp AIKeyResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ShopID != shop.ID {
		t.Errorf("shop_id = %q, want %q", resp.ShopID, shop.ID)
	}
	if resp.ConnectionName != "romashka" {
		t.Errorf("connection_name = %q, want romashka (default on first edit)", resp.ConnectionName)
	}
}

func TestAIKeyUpdate_TogglePDNOnBoundKey(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	key := createAIKey(t, h, "user-1", shop.ID)

	body, _ := json.Marshal(AIKeyRequest{PDNEnabled: boolPtr(true)})
	req := httptest.NewRequest(http.MethodPut, "/api/ai-keys/"+key.ID, bytes.NewReader(body))
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("togglePDN on bound key: got %d: %s", w.Code, w.Body.String())
	}
	var resp AIKeyResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.PDNEnabled {
		t.Error("pdn_enabled should be true")
	}
	if resp.Name != "Cursor MCP" {
		t.Errorf("name changed by partial update: %q", resp.Name)
	}
	if resp.ShopID != shop.ID {
		t.Errorf("binding lost after partial update: %q", resp.ShopID)
	}
}

func TestAIKeyUpdate_ForeignUser(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	key := createAIKey(t, h, "user-1", shop.ID)

	body, _ := json.Marshal(AIKeyRequest{Name: ptr("hacked")})
	req := httptest.NewRequest(http.MethodPut, "/api/ai-keys/"+key.ID, bytes.NewReader(body))
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-2"))
	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("foreign update: got %d, want 404", w.Code)
	}
}

func TestAIKeyDelete_Success(t *testing.T) {
	h, keys, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	key := createAIKey(t, h, "user-1", shop.ID)

	req := httptest.NewRequest(http.MethodDelete, "/api/ai-keys/"+key.ID, nil)
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
	if _, err := keys.GetByID(key.ID); err == nil {
		t.Error("expected key to be deleted")
	}
}

func TestAIKeyDelete_ForeignUser(t *testing.T) {
	h, keys, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	key := createAIKey(t, h, "user-1", shop.ID)

	req := httptest.NewRequest(http.MethodDelete, "/api/ai-keys/"+key.ID, nil)
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-2"))
	w := httptest.NewRecorder()
	h.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("foreign delete: got %d, want 404", w.Code)
	}
	if _, err := keys.GetByID(key.ID); err != nil {
		t.Error("key must survive foreign delete attempt")
	}
}

func TestAIKeyCreate_Limit(t *testing.T) {
	keys := newMockAIKeyRepo()
	shops := newMockShopRepo()
	h := NewAIKeyHandler(keys, shops, 2)
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)

	for i := 0; i < 2; i++ {
		createAIKey(t, h, "user-1", shop.ID)
	}

	body, _ := json.Marshal(AIKeyRequest{Name: ptr("k3"), ShopID: ptr(shop.ID)})
	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("third create at limit=2: got %d, want 403", w.Code)
	}
	if !strings.Contains(w.Body.String(), "лимит") {
		t.Errorf("403 body should mention the limit: %s", w.Body.String())
	}
}

func TestAIKeyRotate_NewSecretOldDies(t *testing.T) {
	h, keys, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	key := createAIKey(t, h, "user-1", shop.ID)

	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys/"+key.ID+"/rotate", nil)
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()
	h.Rotate(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("rotate: got %d: %s", w.Code, w.Body.String())
	}
	var resp AIKeyResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Key == "" || resp.Key == key.Key {
		t.Errorf("rotate must return a fresh raw key, got %q", resp.Key)
	}
	if resp.Name != key.Name || resp.ShopID != key.ShopID || resp.ConnectionName != key.ConnectionName {
		t.Errorf("name/binding must survive rotation: %+v", resp)
	}

	// old hash must stop authenticating
	if _, err := keys.GetByHash(key.KeyHash); err == nil {
		t.Error("old key hash must not resolve after rotation")
	}
	// new hash resolves
	if _, err := keys.GetByHash(resp.KeyHash); err != nil {
		t.Errorf("new key hash must resolve: %v", err)
	}
}

func TestAIKeyRotate_ForeignUser(t *testing.T) {
	h, _, shops := setupAIKeyHandler()
	shop, _ := shops.Create("user-1", "Ромашка", "r.myinsales.ru", "ik", "ip", true, false, false)
	key := createAIKey(t, h, "user-1", shop.ID)

	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys/"+key.ID+"/rotate", nil)
	req.SetPathValue("id", key.ID)
	req = req.WithContext(withUserID(req.Context(), "user-2"))
	w := httptest.NewRecorder()
	h.Rotate(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("foreign rotate: got %d, want 404", w.Code)
	}
}

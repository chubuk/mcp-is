package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-project/mcp-insales/internal/auth"
)

func setupShopHandler() (*ShopHandler, *mockShopRepo) {
	repo := newMockShopRepo()
	aiKeyRepo := newMockAIKeyRepo()
	aes, _ := auth.NewAESCipher("12345678901234567890123456789012")
	return NewShopHandler(repo, aiKeyRepo, aes, nil), repo
}

func withUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, auth.UserIDKey, userID)
}

func TestShopCreate_Success(t *testing.T) {
	h, _ := setupShopHandler()

	body, _ := json.Marshal(ShopRequest{
		Name:        "Мой магазин",
		ShopURL:     "myshop.myinsales.ru",
		APIKeyID:    "key123",
		APIPassword: "secret",
		IsDefault:   true,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/shops", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp ShopResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID == "" {
		t.Error("expected shop ID")
	}
	if resp.ShopURL != "myshop.myinsales.ru" {
		t.Errorf("expected shop_url 'myshop.myinsales.ru', got '%s'", resp.ShopURL)
	}
}

func TestShopCreate_MissingFields(t *testing.T) {
	h, _ := setupShopHandler()

	body, _ := json.Marshal(ShopRequest{ShopURL: "url"})
	req := httptest.NewRequest(http.MethodPost, "/api/shops", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestShopCreate_EncryptsPassword(t *testing.T) {
	h, repo := setupShopHandler()

	body, _ := json.Marshal(ShopRequest{
		Name:        "Тестовый",
		ShopURL:     "test.myinsales.ru",
		APIKeyID:    "key",
		APIPassword: "my-password",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/shops", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	shops, _ := repo.ListByUserID("user-1")
	if len(shops) != 1 {
		t.Fatalf("expected 1 shop, got %d", len(shops))
	}
	if shops[0].APIPasswordEncrypted == "my-password" {
		t.Error("password should be encrypted, not plaintext")
	}
}

func TestShopList_Success(t *testing.T) {
	h, repo := setupShopHandler()

	repo.Create("user-1", "Магазин 1", "shop1.ru", "k1", "enc1", true, false)
	repo.Create("user-1", "Магазин 2", "shop2.ru", "k2", "enc2", false, false)

	req := httptest.NewRequest(http.MethodGet, "/api/shops", nil)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result map[string]any
	json.NewDecoder(w.Body).Decode(&result)
	items := result["items"].([]any)
	if len(items) != 2 {
		t.Errorf("expected 2 shops, got %d", len(items))
	}
}

func TestShopList_Empty(t *testing.T) {
	h, _ := setupShopHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/shops", nil)
	req = req.WithContext(withUserID(req.Context(), "user-nobody"))
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []ShopResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 0 {
		t.Errorf("expected 0 shops, got %d", len(resp))
	}
}

func TestShopUpdate_Success(t *testing.T) {
	h, repo := setupShopHandler()

	shop, _ := repo.Create("user-1", "Старое имя", "old.ru", "k1", "enc1", false, false)

	body, _ := json.Marshal(ShopRequest{
		Name:     "Новое имя",
		ShopURL:  "new.ru",
		APIKeyID: "k2",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/shops/"+shop.ID, bytes.NewReader(body))
	req.SetPathValue("id", shop.ID)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestShopUpdate_NotFound(t *testing.T) {
	h, _ := setupShopHandler()

	body, _ := json.Marshal(ShopRequest{Name: "Test", ShopURL: "x"})
	req := httptest.NewRequest(http.MethodPut, "/api/shops/nonexistent", bytes.NewReader(body))
	req.SetPathValue("id", "nonexistent")
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestShopDelete_Success(t *testing.T) {
	h, repo := setupShopHandler()

	shop, _ := repo.Create("user-1", "Удаляемый", "del.ru", "k", "enc", false, false)

	req := httptest.NewRequest(http.MethodDelete, "/api/shops/"+shop.ID, nil)
	req.SetPathValue("id", shop.ID)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTestConnection_MissingFields(t *testing.T) {
	h, _ := setupShopHandler()

	body, _ := json.Marshal(TestConnectionRequest{ShopURL: "test.ru"})
	req := httptest.NewRequest(http.MethodPost, "/api/shops/test-connection", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.TestConnection(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp TestConnectionResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error == "" {
		t.Error("expected error message")
	}
}

func TestTestConnection_InvalidJSON(t *testing.T) {
	h, _ := setupShopHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/shops/test-connection", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	h.TestConnection(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTestConnection_InvalidDomain(t *testing.T) {
	h, _ := setupShopHandler()

	body, _ := json.Marshal(TestConnectionRequest{
		ShopURL:     "nonexistent.invalid",
		APIKeyID:    "fakekey",
		APIPassword: "fakepass",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/shops/test-connection", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.TestConnection(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 (error returned in body), got %d", w.Code)
	}

	var resp TestConnectionResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.OK {
		t.Error("expected ok=false for invalid domain")
	}
	if resp.Error == "" {
		t.Error("expected error message for invalid domain")
	}
}

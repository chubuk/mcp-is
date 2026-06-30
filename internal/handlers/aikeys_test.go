package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupAIKeyHandler() (*AIKeyHandler, *mockAIKeyRepo) {
	repo := newMockAIKeyRepo()
	return NewAIKeyHandler(repo), repo
}

func TestAIKeyCreate_Success(t *testing.T) {
	h, _ := setupAIKeyHandler()

	body, _ := json.Marshal(AIKeyRequest{Name: "Cursor MCP"})
	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp AIKeyResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Key == "" {
		t.Error("expected raw key in response (shown once)")
	}
	if resp.ID == "" {
		t.Error("expected key ID")
	}
	if resp.Name != "Cursor MCP" {
		t.Errorf("expected name 'Cursor MCP', got '%s'", resp.Name)
	}
}

func TestAIKeyCreate_EmptyName(t *testing.T) {
	h, _ := setupAIKeyHandler()

	body, _ := json.Marshal(AIKeyRequest{})
	req := httptest.NewRequest(http.MethodPost, "/api/ai-keys", bytes.NewReader(body))
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAIKeyList_Success(t *testing.T) {
	h, repo := setupAIKeyHandler()

	repo.Create("user-1", "hash1", "key1")
	repo.Create("user-1", "hash2", "key2")

	req := httptest.NewRequest(http.MethodGet, "/api/ai-keys", nil)
	req = req.WithContext(withUserID(req.Context(), "user-1"))
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	items := result["items"].([]interface{})
	if len(items) != 2 {
		t.Errorf("expected 2 keys, got %d", len(items))
	}
}

func TestAIKeyList_Empty(t *testing.T) {
	h, _ := setupAIKeyHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/ai-keys", nil)
	req = req.WithContext(withUserID(req.Context(), "user-nobody"))
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []AIKeyResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 0 {
		t.Errorf("expected 0 keys, got %d", len(resp))
	}
}

func TestAIKeyDelete_Success(t *testing.T) {
	h, repo := setupAIKeyHandler()

	key, _ := repo.Create("user-1", "hash", "name")

	req := httptest.NewRequest(http.MethodDelete, "/api/ai-keys/"+key.ID, nil)
	req.SetPathValue("id", key.ID)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Verify deleted
	_, err := repo.GetByHash("hash")
	if err == nil {
		t.Error("expected key to be deleted")
	}
}

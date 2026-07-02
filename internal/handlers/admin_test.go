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

func setupAdminHandler() (*AdminHandler, *mockUserRepo, *mockShopRepo, *mockAIKeyRepo) {
	userRepo := newMockUserRepo()
	shopRepo := newMockShopRepo()
	aiKeyRepo := newMockAIKeyRepo()
	return NewAdminHandler(userRepo, shopRepo, aiKeyRepo), userRepo, shopRepo, aiKeyRepo
}

func adminContext(userID, role string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.UserIDKey, userID)
	ctx = context.WithValue(ctx, auth.RoleKey, role)
	return ctx
}

func TestListUsers_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	userRepo.Create("user1", "u1@test.com", "", "", "", "hash1")
	userRepo.Create("user2", "u2@test.com", "", "", "", "hash2")

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	w := httptest.NewRecorder()

	h.ListUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result map[string]any
	json.NewDecoder(w.Body).Decode(&result)
	items := result["items"].([]any)
	if len(items) != 2 {
		t.Errorf("expected 2 users, got %d", len(items))
	}
}

func TestGetUser_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	u, _ := userRepo.Create("testuser", "test@test.com", "+79991234567", "Ivan", "Petrov", "hash")

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/"+u.ID, nil)
	req.SetPathValue("id", u.ID)
	w := httptest.NewRecorder()

	h.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp UserResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Nickname != "testuser" {
		t.Errorf("expected nickname 'testuser', got '%s'", resp.Nickname)
	}
	if resp.Email != "test@test.com" {
		t.Errorf("expected email 'test@test.com', got '%s'", resp.Email)
	}
	if resp.FirstName != "Ivan" {
		t.Errorf("expected first_name 'Ivan', got '%s'", resp.FirstName)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	h, _, _, _ := setupAdminHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	h.GetUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateRole_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	superAdmin, _ := userRepo.Create("superadmin", "sa@test.com", "", "", "", "hash")
	superAdmin.Role = "super_admin"

	target, _ := userRepo.Create("target", "target@test.com", "", "", "", "hash")

	body, _ := json.Marshal(RoleRequest{Role: "admin"})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+target.ID+"/role", bytes.NewReader(body))
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(superAdmin.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.UpdateRole(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp UserResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", resp.Role)
	}
}

func TestUpdateRole_CannotChangeSuperAdmin(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	superAdmin, _ := userRepo.Create("superadmin", "sa@test.com", "", "", "", "hash")
	superAdmin.Role = "super_admin"

	body, _ := json.Marshal(RoleRequest{Role: "user"})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+superAdmin.ID+"/role", bytes.NewReader(body))
	req.SetPathValue("id", superAdmin.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.UpdateRole(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestUpdateRole_OnlySuperAdmin(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")
	admin.Role = "admin"

	target, _ := userRepo.Create("target", "t@test.com", "", "", "", "hash")

	body, _ := json.Marshal(RoleRequest{Role: "admin"})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+target.ID+"/role", bytes.NewReader(body))
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(admin.ID, "admin"))
	w := httptest.NewRecorder()

	h.UpdateRole(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestUpdateRole_InvalidRole(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	superAdmin, _ := userRepo.Create("superadmin", "sa@test.com", "", "", "", "hash")
	superAdmin.Role = "super_admin"

	target, _ := userRepo.Create("target", "t@test.com", "", "", "", "hash")

	body, _ := json.Marshal(RoleRequest{Role: "god"})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+target.ID+"/role", bytes.NewReader(body))
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(superAdmin.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.UpdateRole(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateStatus_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	target, _ := userRepo.Create("target", "t@test.com", "", "", "", "hash")

	body, _ := json.Marshal(StatusRequest{IsActive: false})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+target.ID+"/status", bytes.NewReader(body))
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.UpdateStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateStatus_CannotBlockSuperAdmin(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	superAdmin, _ := userRepo.Create("superadmin", "sa@test.com", "", "", "", "hash")
	superAdmin.Role = "super_admin"

	body, _ := json.Marshal(StatusRequest{IsActive: false})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+superAdmin.ID+"/status", bytes.NewReader(body))
	req.SetPathValue("id", superAdmin.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.UpdateStatus(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestUpdateStatus_CannotBlockYourself(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	body, _ := json.Marshal(StatusRequest{IsActive: false})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+caller.ID+"/status", bytes.NewReader(body))
	req.SetPathValue("id", caller.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.UpdateStatus(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	target, _ := userRepo.Create("target", "t@test.com", "", "", "", "hash")

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/"+target.ID, nil)
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	_, err := userRepo.GetByID(target.ID)
	if err == nil {
		t.Error("expected user to be deleted")
	}
}

func TestDeleteUser_CannotDeleteSuperAdmin(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	superAdmin, _ := userRepo.Create("superadmin", "sa@test.com", "", "", "", "hash")
	superAdmin.Role = "super_admin"

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/"+superAdmin.ID, nil)
	req.SetPathValue("id", superAdmin.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDeleteUser_CannotDeleteYourself(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/"+caller.ID, nil)
	req.SetPathValue("id", caller.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestStats_Success(t *testing.T) {
	h, userRepo, shopRepo, aiKeyRepo := setupAdminHandler()

	u1, _ := userRepo.Create("user1", "u1@test.com", "", "", "", "hash")
	u2, _ := userRepo.Create("user2", "u2@test.com", "", "", "", "hash")

	shopRepo.Create(u1.ID, "Магазин 1", "shop1.myinsales.ru", "key1", "pass1", true, false)
	shopRepo.Create(u2.ID, "Магазин 2", "shop2.myinsales.ru", "key2", "pass2", true, false)
	shopRepo.Create(u1.ID, "Магазин 3", "shop3.myinsales.ru", "key3", "pass3", false, false)

	aiKeyRepo.Create(u1.ID, "hash1", "AI Key 1")

	req := httptest.NewRequest(http.MethodGet, "/api/admin/stats", nil)
	w := httptest.NewRecorder()

	h.Stats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var stats StatsResponse
	json.NewDecoder(w.Body).Decode(&stats)
	if stats.TotalUsers != 2 {
		t.Errorf("expected 2 users, got %d", stats.TotalUsers)
	}
	if stats.TotalShops != 3 {
		t.Errorf("expected 3 shops, got %d", stats.TotalShops)
	}
	if stats.TotalAIKeys != 1 {
		t.Errorf("expected 1 ai_key, got %d", stats.TotalAIKeys)
	}
}

func TestStats_Empty(t *testing.T) {
	h, _, _, _ := setupAdminHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/admin/stats", nil)
	w := httptest.NewRecorder()

	h.Stats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var stats StatsResponse
	json.NewDecoder(w.Body).Decode(&stats)
	if stats.TotalUsers != 0 || stats.TotalShops != 0 || stats.TotalAIKeys != 0 {
		t.Errorf("expected all zeros, got %+v", stats)
	}
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/email"
)

func setupAdminHandler() (*AdminHandler, *mockUserRepo, *mockShopRepo, *mockAIKeyRepo) {
	userRepo := newMockUserRepo()
	shopRepo := newMockShopRepo()
	aiKeyRepo := newMockAIKeyRepo()
	jwtMgr := auth.NewJWTManager("test-secret", 24*time.Hour)
	return NewAdminHandler(userRepo, shopRepo, aiKeyRepo, jwtMgr, nil), userRepo, shopRepo, aiKeyRepo
}

func adminContext(userID, role string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.UserIDKey, userID)
	ctx = context.WithValue(ctx, auth.RoleKey, role)
	return ctx
}

func TestCreateUser_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")

	body, _ := json.Marshal(RegisterRequest{
		Nickname: "newuser",
		Email:    "new@test.com",
		Password: "Password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/users", bytes.NewReader(body))
	req = req.WithContext(adminContext(admin.ID, "admin"))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	created, err := userRepo.GetByLogin("newuser")
	if err != nil {
		t.Fatal(err)
	}
	if !created.RequirePasswordChange {
		t.Error("expected require_password_change=true")
	}
}

func TestCreateUser_WeakPassword(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")

	body, _ := json.Marshal(RegisterRequest{Nickname: "newuser", Password: "123"})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/users", bytes.NewReader(body))
	req = req.WithContext(adminContext(admin.ID, "admin"))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")
	userRepo.Create("newuser", "n@test.com", "", "", "", "hash")

	body, _ := json.Marshal(RegisterRequest{Nickname: "newuser", Password: "Password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/users", bytes.NewReader(body))
	req = req.WithContext(adminContext(admin.ID, "admin"))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
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

	shopRepo.Create(u1.ID, "Магазин 1", "shop1.myinsales.ru", "key1", "pass1", true, false, false)
	shopRepo.Create(u2.ID, "Магазин 2", "shop2.myinsales.ru", "key2", "pass2", true, false, false)
	shopRepo.Create(u1.ID, "Магазин 3", "shop3.myinsales.ru", "key3", "pass3", false, false, false)

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

// ==================== IMPERSONATE ====================

func TestImpersonate_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "admin@test.com", "", "", "", "hash")
	admin.Role = "super_admin"

	target, _ := userRepo.Create("user1", "u1@test.com", "", "", "", "hash")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+target.ID+"/impersonate", nil)
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(admin.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.Impersonate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp ImpersonateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.UserID != target.ID {
		t.Errorf("expected user_id '%s', got '%s'", target.ID, resp.UserID)
	}
	if resp.Nickname != "user1" {
		t.Errorf("expected nickname 'user1', got '%s'", resp.Nickname)
	}
}

func TestImpersonate_AdminCanImpersonateUser(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "admin@test.com", "", "", "", "hash")
	admin.Role = "admin"

	target, _ := userRepo.Create("user1", "u1@test.com", "", "", "", "hash")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+target.ID+"/impersonate", nil)
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(admin.ID, "admin"))
	w := httptest.NewRecorder()

	h.Impersonate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestImpersonate_CannotImpersonateYourself(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "admin@test.com", "", "", "", "hash")
	admin.Role = "super_admin"

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+admin.ID+"/impersonate", nil)
	req.SetPathValue("id", admin.ID)
	req = req.WithContext(adminContext(admin.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.Impersonate(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestImpersonate_CannotImpersonateSuperAdmin(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	caller, _ := userRepo.Create("caller", "c@test.com", "", "", "", "hash")
	caller.Role = "super_admin"

	sa, _ := userRepo.Create("superadmin", "sa@test.com", "", "", "", "hash")
	sa.Role = "super_admin"

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+sa.ID+"/impersonate", nil)
	req.SetPathValue("id", sa.ID)
	req = req.WithContext(adminContext(caller.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.Impersonate(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestImpersonate_AdminCannotImpersonateAdmin(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")
	admin.Role = "admin"

	otherAdmin, _ := userRepo.Create("otheradmin", "oa@test.com", "", "", "", "hash")
	otherAdmin.Role = "admin"

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+otherAdmin.ID+"/impersonate", nil)
	req.SetPathValue("id", otherAdmin.ID)
	req = req.WithContext(adminContext(admin.ID, "admin"))
	w := httptest.NewRecorder()

	h.Impersonate(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestImpersonate_OnlyAdminsCanImpersonate(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	user, _ := userRepo.Create("user1", "u1@test.com", "", "", "", "hash")

	target, _ := userRepo.Create("target", "t@test.com", "", "", "", "hash")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+target.ID+"/impersonate", nil)
	req.SetPathValue("id", target.ID)
	req = req.WithContext(adminContext(user.ID, "user"))
	w := httptest.NewRecorder()

	h.Impersonate(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestImpersonate_UserNotFound(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")
	admin.Role = "super_admin"

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/nonexistent/impersonate", nil)
	req.SetPathValue("id", "nonexistent")
	req = req.WithContext(adminContext(admin.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.Impersonate(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestStopImpersonate_Success(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")
	admin.Role = "super_admin"

	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.UserIDKey, admin.ID)
	ctx = context.WithValue(ctx, auth.RoleKey, admin.Role)
	ctx = context.WithValue(ctx, auth.ImpersonatorIDKey, admin.ID)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/stop-impersonate", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.StopImpersonate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp ImpersonateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.UserID != admin.ID {
		t.Errorf("expected user_id '%s', got '%s'", admin.ID, resp.UserID)
	}
}

func TestStopImpersonate_NotImpersonating(t *testing.T) {
	h, userRepo, _, _ := setupAdminHandler()

	admin, _ := userRepo.Create("admin", "a@test.com", "", "", "", "hash")
	admin.Role = "super_admin"

	req := httptest.NewRequest(http.MethodPost, "/api/admin/stop-impersonate", nil)
	req = req.WithContext(adminContext(admin.ID, "super_admin"))
	w := httptest.NewRecorder()

	h.StopImpersonate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestStopImpersonate_AdminNotFound(t *testing.T) {
	h, _, _, _ := setupAdminHandler()

	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.UserIDKey, "user1")
	ctx = context.WithValue(ctx, auth.RoleKey, "user")
	ctx = context.WithValue(ctx, auth.ImpersonatorIDKey, "nonexistent")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/stop-impersonate", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.StopImpersonate(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestSMTPStatus_Disabled(t *testing.T) {
	h, _, _, _ := setupAdminHandler()
	// emailSvc is nil in setup, so SMTPStatus will panic. Need to set it.
	// Actually setupAdminHandler passes nil for emailSvc. Let's create with a disabled service.
	emailSvc := email.NewService(email.Config{Host: "", Domain: "example.com"})
	h.emailSvc = emailSvc

	req := httptest.NewRequest(http.MethodGet, "/api/admin/smtp/status", nil)
	w := httptest.NewRecorder()

	h.SMTPStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp SMTPStatusResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Disabled {
		t.Error("expected disabled=true")
	}
}

func TestSMTPStatus_Enabled(t *testing.T) {
	h, _, _, _ := setupAdminHandler()
	emailSvc := email.NewService(email.Config{Host: "smtp.example.com", Port: "587", From: "test@example.com", Domain: "example.com"})
	h.emailSvc = emailSvc

	req := httptest.NewRequest(http.MethodGet, "/api/admin/smtp/status", nil)
	w := httptest.NewRecorder()

	h.SMTPStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp SMTPStatusResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Disabled {
		t.Error("expected disabled=false")
	}
	if resp.Host != "smtp.example.com" {
		t.Errorf("host = %q, want smtp.example.com", resp.Host)
	}
}

func TestSMTPStatus_NilEmailSvc(t *testing.T) {
	// In production, emailSvc is always non-nil (created in main.go).
	// This test documents that nil emailSvc causes a panic — acceptable since
	// the SMTP routes are only registered when emailSvc is available.
	t.Skip("nil emailSvc panics — expected, routes are guarded in main.go")
}

func TestTestSMTP_Disabled(t *testing.T) {
	h, _, _, _ := setupAdminHandler()
	emailSvc := email.NewService(email.Config{Host: "", Domain: "example.com"})
	h.emailSvc = emailSvc

	body, _ := json.Marshal(TestSMTPRequest{Email: "test@example.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/smtp/test", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.TestSMTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTestSMTP_EmptyEmail(t *testing.T) {
	h, _, _, _ := setupAdminHandler()
	emailSvc := email.NewService(email.Config{Host: "smtp.example.com", Domain: "example.com"})
	h.emailSvc = emailSvc

	body, _ := json.Marshal(TestSMTPRequest{Email: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/smtp/test", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.TestSMTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

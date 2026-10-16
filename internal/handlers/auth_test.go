package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthHandler() (*AuthHandler, *mockUserRepo) {
	repo := newMockUserRepo()
	jwtMgr := auth.NewJWTManager("test-secret-key-for-testing1234", 0)
	return NewAuthHandler(repo, jwtMgr, nil), repo
}

func TestRegister_Success(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{
		Nickname: "testuser",
		Email:    "test@example.com",
		Phone:    "+79991234567",
		Password: "Password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.UserID == "" {
		t.Error("expected user_id in response")
	}
	if resp.Nickname != "testuser" {
		t.Errorf("expected nickname 'testuser', got '%s'", resp.Nickname)
	}
}

func TestRegister_MissingNickname(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Password: "Password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "user", Password: "123"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_DuplicateNickname(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "user1", Password: "Password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("first register failed: %d", w.Code)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w2 := httptest.NewRecorder()
	h.Register(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate, got %d", w2.Code)
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	h, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLogin_Success_ByNickname(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "myuser", Password: "Secret123"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	loginBody, _ := json.Marshal(LoginRequest{Login: "myuser", Password: "Secret123"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", loginW.Code, loginW.Body.String())
	}

	var resp AuthResponse
	json.NewDecoder(loginW.Body).Decode(&resp)
	if resp.Token == "" {
		t.Error("expected token in response")
	}
}

func TestLogin_Success_ByEmail(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "emailuser", Email: "me@test.com", Password: "Pass1234"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	loginBody, _ := json.Marshal(LoginRequest{Login: "me@test.com", Password: "Pass1234"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", loginW.Code, loginW.Body.String())
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "user2", Password: "Correct1"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	loginBody, _ := json.Marshal(LoginRequest{Login: "user2", Password: "wrong"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Login(loginW, loginReq)

	if loginW.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", loginW.Code)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(LoginRequest{Login: "nobody", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(LoginRequest{Login: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestChangePassword_Success(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "chpwd", Password: "Oldpass123"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	chgBody, _ := json.Marshal(ChangePasswordRequest{OldPassword: "Oldpass123", NewPassword: "Newpass123"})
	chgReq := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewReader(chgBody))
	chgReq = chgReq.WithContext(withUserID(chgReq.Context(), "user-1"))
	chgW := httptest.NewRecorder()
	h.ChangePassword(chgW, chgReq)

	if chgW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", chgW.Code, chgW.Body.String())
	}
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "chpwd2", Password: "Correct1"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	chgBody, _ := json.Marshal(ChangePasswordRequest{OldPassword: "wrong", NewPassword: "Newpass1"})
	chgReq := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewReader(chgBody))
	chgReq = chgReq.WithContext(withUserID(chgReq.Context(), "user-1"))
	chgW := httptest.NewRecorder()
	h.ChangePassword(chgW, chgReq)

	if chgW.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", chgW.Code)
	}
}

func TestForgotPassword_Success(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "forgot", Email: "forgot@test.com", Password: "Pass1234"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	forgotBody, _ := json.Marshal(ForgotPasswordRequest{Email: "forgot@test.com"})
	forgotReq := httptest.NewRequest(http.MethodPost, "/api/forgot-password", bytes.NewReader(forgotBody))
	forgotW := httptest.NewRecorder()
	h.ForgotPassword(forgotW, forgotReq)

	if forgotW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", forgotW.Code)
	}
}

func TestSetEmail_Success(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "setemail", Password: "Pass1234"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	setBody, _ := json.Marshal(SetEmailRequest{Email: "new@test.com"})
	setReq := httptest.NewRequest(http.MethodPost, "/api/set-email", bytes.NewReader(setBody))
	setReq = setReq.WithContext(withUserID(setReq.Context(), "user-1"))
	setW := httptest.NewRecorder()
	h.SetEmail(setW, setReq)

	if setW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", setW.Code, setW.Body.String())
	}
}

func TestLogin_RequirePasswordChange(t *testing.T) {
	h, repo := setupAuthHandler()

	hash, _ := bcrypt.GenerateFromPassword([]byte("Pass1234"), bcrypt.DefaultCost)
	user := &storage.User{
		ID:                    "sa-1",
		Nickname:              "admin",
		Email:                 "admin@test.com",
		PasswordHash:          string(hash),
		Role:                  "super_admin",
		IsActive:              true,
		RequirePasswordChange: true,
	}
	repo.users[user.ID] = user

	loginBody, _ := json.Marshal(LoginRequest{Login: "admin", Password: "Pass1234"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", loginW.Code, loginW.Body.String())
	}

	var resp AuthResponse
	json.NewDecoder(loginW.Body).Decode(&resp)
	if !resp.RequirePasswordChange {
		t.Error("expected require_password_change=true")
	}
	if resp.Token == "" {
		t.Error("expected token in response")
	}
}

func TestLogin_NoForceChange_WhenFlagFalse(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "nouser", Password: "Pass1234"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	loginBody, _ := json.Marshal(LoginRequest{Login: "nouser", Password: "Pass1234"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", loginW.Code, loginW.Body.String())
	}

	var resp AuthResponse
	json.NewDecoder(loginW.Body).Decode(&resp)
	if resp.RequirePasswordChange {
		t.Error("expected require_password_change=false (not set)")
	}
}

func TestChangePassword_ClearsForceFlag(t *testing.T) {
	h, repo := setupAuthHandler()

	hash, _ := bcrypt.GenerateFromPassword([]byte("Old12345"), bcrypt.DefaultCost)
	user := &storage.User{
		ID:                    "cp-1",
		Nickname:              "chpwd",
		PasswordHash:          string(hash),
		Role:                  "super_admin",
		IsActive:              true,
		RequirePasswordChange: true,
	}
	repo.users[user.ID] = user

	chgBody, _ := json.Marshal(ChangePasswordRequest{OldPassword: "Old12345", NewPassword: "New12345"})
	chgReq := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewReader(chgBody))
	chgReq = chgReq.WithContext(withUserID(chgReq.Context(), "cp-1"))
	chgW := httptest.NewRecorder()
	h.ChangePassword(chgW, chgReq)

	if chgW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", chgW.Code, chgW.Body.String())
	}

	if user.RequirePasswordChange {
		t.Error("expected require_password_change to be cleared after password change")
	}
}

func TestResetPassword_MissingToken(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(ResetPasswordRequest{NewPassword: "NewPass123"})
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ResetPassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestResetPassword_MissingPassword(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(ResetPasswordRequest{Token: "sometoken"})
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ResetPassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestResetPassword_ShortPassword(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(ResetPasswordRequest{Token: "sometoken", NewPassword: "short"})
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ResetPassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestResetPassword_WeakPassword(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(ResetPasswordRequest{Token: "sometoken", NewPassword: "alllowercase1"})
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ResetPassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestResetPassword_InvalidToken(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(ResetPasswordRequest{Token: "invalidtoken", NewPassword: "NewPass123"})
	req := httptest.NewRequest(http.MethodPost, "/api/reset-password", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ResetPassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestForgotPassword_EmptyEmail(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(ForgotPasswordRequest{Email: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/forgot-password", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ForgotPassword(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestForgotPassword_UnknownEmail(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(ForgotPasswordRequest{Email: "unknown@test.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/forgot-password", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ForgotPassword(w, req)

	// Should return 200 even for unknown email (prevents user enumeration)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

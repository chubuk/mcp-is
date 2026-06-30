package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-project/mcp-insales/internal/auth"
)

func setupAuthHandler() (*AuthHandler, *mockUserRepo) {
	repo := newMockUserRepo()
	jwtMgr := auth.NewJWTManager("test-secret-key-for-testing1234")
	return NewAuthHandler(repo, jwtMgr, nil), repo
}

func TestRegister_Success(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{
		Nickname: "testuser",
		Email:    "test@example.com",
		Phone:    "+79991234567",
		Password: "password123",
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

	body, _ := json.Marshal(RegisterRequest{Password: "password123"})
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

	body, _ := json.Marshal(RegisterRequest{Nickname: "user1", Password: "password123"})
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

	body, _ := json.Marshal(RegisterRequest{Nickname: "myuser", Password: "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	loginBody, _ := json.Marshal(LoginRequest{Login: "myuser", Password: "secret123"})
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

	body, _ := json.Marshal(RegisterRequest{Nickname: "emailuser", Email: "me@test.com", Password: "pass1234"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	loginBody, _ := json.Marshal(LoginRequest{Login: "me@test.com", Password: "pass1234"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()
	h.Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", loginW.Code, loginW.Body.String())
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	h, _ := setupAuthHandler()

	body, _ := json.Marshal(RegisterRequest{Nickname: "user2", Password: "correct"})
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

	body, _ := json.Marshal(RegisterRequest{Nickname: "chpwd", Password: "oldpass123"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	chgBody, _ := json.Marshal(ChangePasswordRequest{OldPassword: "oldpass123", NewPassword: "newpass123"})
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

	body, _ := json.Marshal(RegisterRequest{Nickname: "chpwd2", Password: "correct"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	chgBody, _ := json.Marshal(ChangePasswordRequest{OldPassword: "wrong", NewPassword: "newpass"})
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

	body, _ := json.Marshal(RegisterRequest{Nickname: "forgot", Email: "forgot@test.com", Password: "pass1234"})
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

	body, _ := json.Marshal(RegisterRequest{Nickname: "setemail", Password: "pass1234"})
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

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/email"
	"github.com/chubuk/mcp-is/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo            storage.UserRepositoryInterface
	jwtMgr              *auth.JWTManager
	emailSvc            *email.Service
	registrationEnabled bool
}

func NewAuthHandler(userRepo storage.UserRepositoryInterface, jwtMgr *auth.JWTManager, emailSvc *email.Service, registrationEnabled bool) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtMgr: jwtMgr, emailSvc: emailSvc, registrationEnabled: registrationEnabled}
}

type RegisterRequest struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token                 string `json:"token,omitempty"`
	UserID                string `json:"user_id,omitempty"`
	Nickname              string `json:"nickname,omitempty"`
	Role                  string `json:"role,omitempty"`
	RequireEmail          bool   `json:"require_email,omitempty"`
	RequirePasswordChange bool   `json:"require_password_change,omitempty"`
	Message               string `json:"message,omitempty"`
	Error                 string `json:"error,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type SetEmailRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if msg := validateNewUser(req.Nickname, req.Password); msg != "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: msg})
		return
	}

	count, err := h.userRepo.Count()
	if err != nil {
		slog.Error("count users", "error", err)
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}

	if count > 0 && !h.registrationEnabled {
		writeJSON(w, http.StatusForbidden, AuthResponse{Error: "регистрация закрыта"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("hash password", "error", err)
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}

	var user *storage.User
	if count == 0 {
		user, err = h.userRepo.CreateSuperAdmin(req.Nickname, req.Email, req.Phone, req.FirstName, req.LastName, string(hash))
		if err != nil {
			slog.Error("create super_admin", "error", err)
			writeJSON(w, http.StatusConflict, AuthResponse{Error: "user with this nickname, email or phone already exists"})
			return
		}
	} else {
		user, err = h.userRepo.Create(req.Nickname, req.Email, req.Phone, req.FirstName, req.LastName, string(hash))
		if err != nil {
			slog.Warn("register failed", "error", err)
			writeJSON(w, http.StatusConflict, AuthResponse{Error: "user with this nickname, email or phone already exists"})
			return
		}
	}

	slog.Info("user registered", "id", user.ID, "nickname", user.Nickname, "role", user.Role)
	writeJSON(w, http.StatusCreated, AuthResponse{UserID: user.ID, Nickname: user.Nickname})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if req.Login == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "login and password are required"})
		return
	}

	user, err := h.userRepo.GetByLogin(req.Login)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, AuthResponse{Error: "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeJSON(w, http.StatusUnauthorized, AuthResponse{Error: "invalid credentials"})
		return
	}

	if user.Email == "" {
		token, err := h.jwtMgr.GenerateToken(user.ID, user.Nickname, user.Role)
		if err != nil {
			slog.Error("generate token", "error", err)
			writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
			return
		}
		writeJSON(w, http.StatusOK, AuthResponse{
			Token:        token,
			UserID:       user.ID,
			Nickname:     user.Nickname,
			Role:         user.Role,
			RequireEmail: true,
		})
		return
	}

	if user.RequirePasswordChange {
		token, err := h.jwtMgr.GenerateToken(user.ID, user.Nickname, user.Role)
		if err != nil {
			slog.Error("generate token", "error", err)
			writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
			return
		}
		writeJSON(w, http.StatusOK, AuthResponse{
			Token:                 token,
			UserID:                user.ID,
			Nickname:              user.Nickname,
			Role:                  user.Role,
			RequirePasswordChange: true,
		})
		return
	}

	token, err := h.jwtMgr.GenerateToken(user.ID, user.Nickname, user.Role)
	if err != nil {
		slog.Error("generate token", "error", err)
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}

	slog.Info("user logged in", "id", user.ID, "nickname", user.Nickname, "role", user.Role)
	writeJSON(w, http.StatusOK, AuthResponse{Token: token, UserID: user.ID, Nickname: user.Nickname, Role: user.Role})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSON(w, http.StatusUnauthorized, AuthResponse{Error: "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "old and new password are required"})
		return
	}

	if len(req.NewPassword) < 8 {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "new password must be at least 8 characters"})
		return
	}

	if !hasUpperCase(req.NewPassword) || !hasDigit(req.NewPassword) {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "new password must contain at least one uppercase letter and one digit"})
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, AuthResponse{Error: "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		writeJSON(w, http.StatusUnauthorized, AuthResponse{Error: "incorrect current password"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}

	if err := h.userRepo.UpdatePassword(userID, string(hash)); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to update password"})
		return
	}

	h.userRepo.UpdateRequirePasswordChange(userID, false)

	slog.Info("password changed", "user_id", userID)
	writeJSON(w, http.StatusOK, AuthResponse{UserID: userID})
}

func (h *AuthHandler) ForceChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSON(w, http.StatusUnauthorized, AuthResponse{Error: "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "new password is required"})
		return
	}

	if len(req.NewPassword) < 8 {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "new password must be at least 8 characters"})
		return
	}

	if !hasUpperCase(req.NewPassword) || !hasDigit(req.NewPassword) {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "new password must contain at least one uppercase letter and one digit"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}

	if err := h.userRepo.UpdatePassword(userID, string(hash)); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to update password"})
		return
	}

	h.userRepo.UpdateRequirePasswordChange(userID, false)

	slog.Info("password force changed", "user_id", userID)
	writeJSON(w, http.StatusOK, AuthResponse{UserID: userID})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if req.Email == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "email is required"})
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		writeJSON(w, http.StatusOK, AuthResponse{Message: "If the email exists, a reset link has been sent"})
		return
	}

	token, err := generateResetToken()
	if err != nil {
		slog.Error("generate reset token", "error", err)
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}
	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	h.userRepo.CleanExpiredResetTokens()

	if err := h.userRepo.CreateResetToken(user.ID, token, expiresAt); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to create reset token"})
		return
	}

	slog.Info("password reset requested", "user_id", user.ID, "email", maskEmail(user.Email))

	if h.emailSvc != nil {
		name := user.Nickname
		if user.FirstName != "" {
			name = user.FirstName
		}
		go h.emailSvc.SendResetPassword(user.Email, name, token)
	}

	writeJSON(w, http.StatusOK, AuthResponse{Message: "If the email exists, a reset link has been sent"})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "token and new password are required"})
		return
	}

	if len(req.NewPassword) < 8 {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "password must be at least 8 characters"})
		return
	}

	if !hasUpperCase(req.NewPassword) || !hasDigit(req.NewPassword) {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "password must contain at least one uppercase letter and one digit"})
		return
	}

	user, err := h.userRepo.GetByResetToken(req.Token)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid or expired token"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}

	if err := h.userRepo.UpdatePassword(user.ID, string(hash)); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to update password"})
		return
	}

	h.userRepo.UpdateRequirePasswordChange(user.ID, false)
	h.userRepo.ClearResetTokens(user.ID)
	slog.Info("password reset", "user_id", user.ID)
	writeJSON(w, http.StatusOK, AuthResponse{Message: "Password updated successfully"})
}

func (h *AuthHandler) SetEmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSON(w, http.StatusUnauthorized, AuthResponse{Error: "unauthorized"})
		return
	}

	var req SetEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if req.Email == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "email is required"})
		return
	}

	if !isValidEmail(req.Email) {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid email format"})
		return
	}

	if err := h.userRepo.UpdateEmail(userID, req.Email); err != nil {
		writeJSON(w, http.StatusConflict, AuthResponse{Error: "email already exists"})
		return
	}

	slog.Info("email set", "user_id", userID, "email", maskEmail(req.Email))
	writeJSON(w, http.StatusOK, AuthResponse{UserID: userID})
}

func generateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate reset token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func validateNewUser(nickname, password string) string {
	if nickname == "" || password == "" {
		return "укажите никнейм и пароль"
	}
	if len(password) < 8 {
		return "пароль должен быть не менее 8 символов"
	}
	if !hasUpperCase(password) || !hasDigit(password) {
		return "пароль должен содержать заглавную букву и цифру"
	}
	return ""
}

func hasUpperCase(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func hasDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && strings.Contains(email, "@")
}

func maskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return "***"
	}
	local := parts[0]
	if len(local) > 2 {
		return local[:2] + "***@" + parts[1]
	}
	return "***@" + parts[1]
}

package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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
	codes               *codeStore
}

func NewAuthHandler(userRepo storage.UserRepositoryInterface, jwtMgr *auth.JWTManager, emailSvc *email.Service, registrationEnabled bool) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtMgr: jwtMgr, emailSvc: emailSvc, registrationEnabled: registrationEnabled, codes: newCodeStore()}
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

type VerifyEmailRequest struct {
	UserID   string `json:"user_id"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "invalid request body"})
		return
	}

	if req.Nickname == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "укажите никнейм"})
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

	// Анти-разбухание базы: неподтверждённые аккаунты старше 24ч удаляются.
	if n, err := h.userRepo.PurgeUnverified(24 * 3600); err != nil {
		slog.Warn("purge unverified users", "error", err)
	} else if n > 0 {
		slog.Info("purged unverified users", "count", n)
	}

	var user *storage.User
	if count == 0 {
		// Bootstrap: первый super_admin — пароль сразу (SMTP может быть не настроен).
		if msg := validateNewUser(req.Nickname, req.Password); msg != "" {
			writeJSON(w, http.StatusBadRequest, AuthResponse{Error: msg})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("hash password", "error", err)
			writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
			return
		}
		user, err = h.userRepo.CreateSuperAdmin(req.Nickname, req.Email, req.Phone, req.FirstName, req.LastName, string(hash))
		if err != nil {
			slog.Error("create super_admin", "error", err)
			writeJSON(w, http.StatusConflict, AuthResponse{Error: "user with this nickname, email or phone already exists"})
			return
		}
		slog.Info("user registered", "id", user.ID, "nickname", user.Nickname, "role", user.Role)
		writeJSON(w, http.StatusCreated, AuthResponse{UserID: user.ID, Nickname: user.Nickname})
		return
	}

	// Verify-before-use: пароль ставит подтвердивший почту (код + пароль в одной
	// форме) — до верификации аккаунт с нерабочим паролем, вход закрыт.
	if req.Email == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "для регистрации нужен e-mail"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(randomPassword()), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("hash password", "error", err)
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}
	user, err = h.userRepo.Create(req.Nickname, req.Email, req.Phone, req.FirstName, req.LastName, string(hash))
	if err != nil {
		slog.Warn("register failed", "error", err)
		writeJSON(w, http.StatusConflict, AuthResponse{Error: "user with this nickname, email or phone already exists"})
		return
	}
	if err := h.userRepo.UpdateEmailVerified(user.ID, false); err != nil {
		slog.Error("set email unverified", "error", err)
	}
	user.EmailVerified = false
	h.sendVerifyCode(user)
	slog.Info("user registered, email pending", "id", user.ID, "nickname", user.Nickname)
	writeJSON(w, http.StatusCreated, AuthResponse{
		UserID:       user.ID,
		Nickname:     user.Nickname,
		RequireEmail: true,
		Message:      "Подтвердите e-mail — код отправлен на почту. Пароль задаётся при подтверждении.",
	})
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

	// Verify-before-use: неподтверждённый e-mail не даёт доступа (анти pre-hijack).
	if user.Email != "" && !user.EmailVerified {
		h.sendVerifyCode(user)
		writeJSON(w, http.StatusForbidden, AuthResponse{
			RequireEmail: true,
			UserID:       user.ID,
			Nickname:     user.Nickname,
			Message:      "e-mail не подтверждён. Код подтверждения отправлен на почту.",
		})
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

	if err := h.userRepo.UpdateRequirePasswordChange(userID, false); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to update password"})
		return
	}

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

	if err := h.userRepo.UpdateRequirePasswordChange(userID, false); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to update password"})
		return
	}

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

	if err := h.userRepo.CleanExpiredResetTokens(); err != nil {
		slog.Warn("clean expired reset tokens", "error", err)
	}

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
		go func() {
			if err := h.emailSvc.SendResetPassword(user.Email, name, token); err != nil {
				slog.Error("send reset password email", "error", err, "user_id", user.ID)
			}
		}()
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

	if err := h.userRepo.UpdateRequirePasswordChange(user.ID, false); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to update password"})
		return
	}
	if err := h.userRepo.ClearResetTokens(user.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to update password"})
		return
	}
	slog.Info("password reset", "user_id", user.ID)
	writeJSON(w, http.StatusOK, AuthResponse{Message: "Password updated successfully"})
}

// sendVerifyCode issues a one-time code and mails it (async, cooldown-safe).
func (h *AuthHandler) sendVerifyCode(user *storage.User) {
	if user.Email == "" {
		return
	}
	code, err := h.codes.issue("register:" + user.ID)
	if err != nil {
		if err != ErrCodeCooldown {
			slog.Error("issue verify code", "error", err, "user_id", user.ID)
		}
		return
	}
	if h.emailSvc == nil {
		return
	}
	go func() {
		if err := h.emailSvc.SendTemplate(context.Background(), "login_code", user.Email,
			map[string]any{"Name": user.Nickname, "Code": code}); err != nil {
			slog.Error("send verify code email", "error", err, "user_id", user.ID)
		}
	}()
}

// VerifyEmail handles POST /api/verify-email — подтверждение e-mail кодом и
// установка пароля подтвердившим (инвариант: пароль ставит только владелец ящика).
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Code == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "user_id и code обязательны"})
		return
	}
	// Пароль проверяем ДО сжигания кода, чтобы нестандартный пароль не отнимал попытку.
	user, err := h.userRepo.GetByID(req.UserID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, AuthResponse{Error: "пользователь не найден"})
		return
	}
	if msg := validateNewUser(user.Nickname, req.Password); msg != "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: msg})
		return
	}
	if err := h.codes.verify("register:"+req.UserID, req.Code); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "неверный или истёкший код"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}
	if err := h.userRepo.UpdatePassword(user.ID, string(hash)); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}
	if err := h.userRepo.UpdateRequirePasswordChange(user.ID, false); err != nil {
		slog.Warn("clear require_password_change", "error", err, "user_id", user.ID)
	}
	if err := h.userRepo.UpdateEmailVerified(user.ID, true); err != nil {
		slog.Error("set email verified", "error", err, "user_id", user.ID)
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}
	slog.Info("email verified", "user_id", user.ID)
	writeJSON(w, http.StatusOK, AuthResponse{UserID: user.ID, Message: "e-mail подтверждён. Теперь можно войти."})
}

// ResendVerifyEmail handles POST /api/verify-email/resend.
func (h *AuthHandler) ResendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "user_id обязателен"})
		return
	}
	user, err := h.userRepo.GetByID(req.UserID)
	if err != nil {
		// Не раскрываем существование аккаунта.
		writeJSON(w, http.StatusOK, AuthResponse{Message: "Если аккаунт существует, код отправлен."})
		return
	}
	if user.Email == "" || user.EmailVerified {
		writeJSON(w, http.StatusOK, AuthResponse{Message: "Если аккаунт существует, код отправлен."})
		return
	}
	code, err := h.codes.issue("register:" + user.ID)
	if err != nil {
		if err == ErrCodeCooldown {
			writeJSON(w, http.StatusTooManyRequests, AuthResponse{Error: "Слишком часто. Попробуйте позже."})
			return
		}
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}
	go func() {
		if err := h.emailSvc.SendTemplate(context.Background(), "login_code", user.Email,
			map[string]any{"Name": user.Nickname, "Code": code}); err != nil {
			slog.Error("send verify code email", "error", err, "user_id", user.ID)
		}
	}()
	writeJSON(w, http.StatusOK, AuthResponse{Message: "Если аккаунт существует, код отправлен."})
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

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/email"
	"github.com/chubuk/mcp-is/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type AuthHandler struct {
	userRepo storage.UserRepositoryInterface
	jwtMgr   *auth.JWTManager
	emailSvc *email.Service
}

func NewAuthHandler(userRepo storage.UserRepositoryInterface, jwtMgr *auth.JWTManager, emailSvc *email.Service) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtMgr: jwtMgr, emailSvc: emailSvc}
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
	Token                string `json:"token,omitempty"`
	UserID               string `json:"user_id,omitempty"`
	Nickname             string `json:"nickname,omitempty"`
	Role                 string `json:"role,omitempty"`
	RequireEmail         bool   `json:"require_email,omitempty"`
	RequirePasswordChange bool  `json:"require_password_change,omitempty"`
	Message              string `json:"message,omitempty"`
	Error                string `json:"error,omitempty"`
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

	if req.Nickname == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "nickname and password are required"})
		return
	}

	if len(req.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "password must be at least 8 characters"})
		return
	}

	if !hasUpperCase(req.Password) || !hasDigit(req.Password) {
		writeJSON(w, http.StatusBadRequest, AuthResponse{Error: "password must contain at least one uppercase letter and one digit"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("hash password", "error", err)
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "internal error"})
		return
	}

	count, _ := h.userRepo.Count()
	role := "user"
	if count == 0 {
		role = "super_admin"
	}

	user, err := h.userRepo.Create(req.Nickname, req.Email, req.Phone, req.FirstName, req.LastName, string(hash))
	if err != nil {
		slog.Warn("register failed", "error", err)
		writeJSON(w, http.StatusConflict, AuthResponse{Error: "user with this nickname, email or phone already exists"})
		return
	}

	if role == "super_admin" {
		h.userRepo.UpdateRole(user.ID, role)
	}

	slog.Info("user registered", "id", user.ID, "nickname", user.Nickname, "role", role)
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
		token, _ := h.jwtMgr.GenerateToken(user.ID, user.Nickname, user.Role)
		writeJSON(w, http.StatusOK, AuthResponse{
			Token:        token,
			UserID:       user.ID,
			Nickname:     user.Nickname,
			Role:         user.Role,
			RequireEmail: true,
		})
		return
	}

	if user.Role == "super_admin" && isDefaultPassword(user.PasswordHash) {
		token, _ := h.jwtMgr.GenerateToken(user.ID, user.Nickname, user.Role)
		writeJSON(w, http.StatusOK, AuthResponse{
			Token:                token,
			UserID:               user.ID,
			Nickname:             user.Nickname,
			Role:                 user.Role,
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
	userID := r.Context().Value(auth.UserIDKey).(string)

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

	slog.Info("password changed", "user_id", userID)
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

	token := generateResetToken()
	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	if err := h.userRepo.CreateResetToken(user.ID, token, expiresAt); err != nil {
		writeJSON(w, http.StatusInternalServerError, AuthResponse{Error: "failed to create reset token"})
		return
	}

	slog.Info("password reset requested", "user_id", user.ID, "email", user.Email)

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

	h.userRepo.ClearResetTokens(user.ID)
	slog.Info("password reset", "user_id", user.ID)
	writeJSON(w, http.StatusOK, AuthResponse{Message: "Password updated successfully"})
}

func (h *AuthHandler) SetEmail(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

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

	slog.Info("email set", "user_id", userID, "email", req.Email)
	writeJSON(w, http.StatusOK, AuthResponse{UserID: userID})
}

func generateResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
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
	return emailRegex.MatchString(email)
}

func isDefaultPassword(hash string) bool {
	defaultHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	return hash == string(defaultHash)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

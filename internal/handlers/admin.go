package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
)

type AdminHandler struct {
	userRepo storage.UserRepositoryInterface
	shopRepo storage.ShopRepositoryInterface
	aiKeyRepo storage.AIKeyRepositoryInterface
	jwtMgr   *auth.JWTManager
}

func NewAdminHandler(userRepo storage.UserRepositoryInterface, shopRepo storage.ShopRepositoryInterface, aiKeyRepo storage.AIKeyRepositoryInterface, jwtMgr *auth.JWTManager) *AdminHandler {
	return &AdminHandler{userRepo: userRepo, shopRepo: shopRepo, aiKeyRepo: aiKeyRepo, jwtMgr: jwtMgr}
}

type UserResponse struct {
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at,omitempty"`
	Error     string `json:"error,omitempty"`
}

type RoleRequest struct {
	Role string `json:"role"`
}

type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

type StatsResponse struct {
	TotalUsers    int `json:"total_users"`
	TotalShops    int `json:"total_shops"`
	TotalAIKeys   int `json:"total_ai_keys"`
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	query := r.URL.Query().Get("q")

	users, total, err := h.userRepo.ListFiltered(page, perPage, query)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
		return
	}

	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = UserResponse{
			ID: u.ID, Nickname: u.Nickname, Email: u.Email, Phone: u.Phone,
			FirstName: u.FirstName, LastName: u.LastName, Role: u.Role,
			IsActive: u.IsActive, CreatedAt: u.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":    result,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, UserResponse{Error: "user not found"})
		return
	}

	writeJSON(w, http.StatusOK, UserResponse{
		ID: user.ID, Nickname: user.Nickname, Email: user.Email, Phone: user.Phone,
		FirstName: user.FirstName, LastName: user.LastName, Role: user.Role,
		IsActive: user.IsActive, CreatedAt: user.CreatedAt,
	})
}

func (h *AdminHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	targetID := r.PathValue("id")
	callerID := r.Context().Value(auth.UserIDKey).(string)
	callerRole := r.Context().Value(auth.RoleKey).(string)

	target, err := h.userRepo.GetByID(targetID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, UserResponse{Error: "user not found"})
		return
	}

	if target.Role == "super_admin" && targetID != callerID {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "cannot change super_admin role"})
		return
	}

	if callerRole != "super_admin" {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "only super_admin can change roles"})
		return
	}

	var req RoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, UserResponse{Error: "invalid request body"})
		return
	}

	if req.Role != "user" && req.Role != "admin" && req.Role != "super_admin" {
		writeJSON(w, http.StatusBadRequest, UserResponse{Error: "invalid role"})
		return
	}

	if err := h.userRepo.UpdateRole(targetID, req.Role); err != nil {
		writeJSON(w, http.StatusInternalServerError, UserResponse{Error: "failed to update role"})
		return
	}

	slog.Info("role updated", "target", targetID, "new_role", req.Role, "by", callerID)
	writeJSON(w, http.StatusOK, UserResponse{ID: targetID, Role: req.Role})
}

func (h *AdminHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	targetID := r.PathValue("id")
	callerID := r.Context().Value(auth.UserIDKey).(string)

	target, err := h.userRepo.GetByID(targetID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, UserResponse{Error: "user not found"})
		return
	}

	if target.Role == "super_admin" {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "cannot block super_admin"})
		return
	}

	if targetID == callerID {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "cannot block yourself"})
		return
	}

	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, UserResponse{Error: "invalid request body"})
		return
	}

	if err := h.userRepo.UpdateStatus(targetID, req.IsActive); err != nil {
		writeJSON(w, http.StatusInternalServerError, UserResponse{Error: "failed to update status"})
		return
	}

	slog.Info("status updated", "target", targetID, "active", req.IsActive, "by", callerID)
	writeJSON(w, http.StatusOK, UserResponse{ID: targetID, IsActive: req.IsActive})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	targetID := r.PathValue("id")
	callerID := r.Context().Value(auth.UserIDKey).(string)

	target, err := h.userRepo.GetByID(targetID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, UserResponse{Error: "user not found"})
		return
	}

	if target.Role == "super_admin" {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "cannot delete super_admin"})
		return
	}

	if targetID == callerID {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "cannot delete yourself"})
		return
	}

	if err := h.userRepo.Delete(targetID); err != nil {
		writeJSON(w, http.StatusInternalServerError, UserResponse{Error: "failed to delete user"})
		return
	}

	slog.Info("user deleted", "target", targetID, "nickname", target.Nickname, "by", callerID)
	writeJSON(w, http.StatusOK, UserResponse{ID: targetID})
}

func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	users, _ := h.userRepo.Count()
	shops := 0
	keys := 0

	allUsers, err := h.userRepo.ListAll()
	if err == nil {
		for _, u := range allUsers {
			s, _ := h.shopRepo.ListByUserID(u.ID)
			shops += len(s)
			k, _ := h.aiKeyRepo.ListByUserID(u.ID)
			keys += len(k)
		}
	}

	writeJSON(w, http.StatusOK, StatsResponse{
		TotalUsers:  users,
		TotalShops:  shops,
		TotalAIKeys: keys,
	})
}

type ImpersonateResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

func (h *AdminHandler) Impersonate(w http.ResponseWriter, r *http.Request) {
	callerID := r.Context().Value(auth.UserIDKey).(string)
	callerRole := r.Context().Value(auth.RoleKey).(string)

	if callerRole != "super_admin" && callerRole != "admin" {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "only admins can impersonate"})
		return
	}

	targetID := r.PathValue("id")
	if targetID == callerID {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "cannot impersonate yourself"})
		return
	}

	target, err := h.userRepo.GetByID(targetID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, UserResponse{Error: "user not found"})
		return
	}

	if target.Role == "super_admin" {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "cannot impersonate super_admin"})
		return
	}

	if callerRole == "admin" && target.Role == "admin" {
		writeJSON(w, http.StatusForbidden, UserResponse{Error: "admin cannot impersonate another admin"})
		return
	}

	token, err := h.jwtMgr.GenerateTokenWithImpersonator(target.ID, target.Nickname, target.Role, callerID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, UserResponse{Error: "failed to generate token"})
		return
	}

	slog.Info("impersonate started", "target", targetID, "target_nickname", target.Nickname, "by", callerID, "by_role", callerRole)
	writeJSON(w, http.StatusOK, ImpersonateResponse{
		Token:    token,
		UserID:   target.ID,
		Nickname: target.Nickname,
		Role:     target.Role,
	})
}

func (h *AdminHandler) StopImpersonate(w http.ResponseWriter, r *http.Request) {
	callerID := r.Context().Value(auth.UserIDKey).(string)
	impersonatorID, ok := r.Context().Value(auth.ImpersonatorIDKey).(string)
	if !ok || impersonatorID == "" {
		writeJSON(w, http.StatusBadRequest, UserResponse{Error: "not impersonating"})
		return
	}

	admin, err := h.userRepo.GetByID(impersonatorID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, UserResponse{Error: "admin not found"})
		return
	}

	token, err := h.jwtMgr.GenerateToken(admin.ID, admin.Nickname, admin.Role)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, UserResponse{Error: "failed to generate token"})
		return
	}

	slog.Info("impersonate stopped", "impersonator", impersonatorID, "was_using", callerID)
	writeJSON(w, http.StatusOK, ImpersonateResponse{
		Token:    token,
		UserID:   admin.ID,
		Nickname: admin.Nickname,
		Role:     admin.Role,
	})
}

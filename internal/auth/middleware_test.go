package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJWTMiddleware_PassesImpersonatorID(t *testing.T) {
	mgr := NewJWTManager("test-secret", 24*time.Hour)

	token, _ := mgr.GenerateTokenWithImpersonator("user-1", "user1", "user", "admin-1")

	var gotUserID, gotImpersonatorID string
	var hasImpersonator bool

	handler := JWTMiddleware(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = r.Context().Value(UserIDKey).(string)
		gotImpersonatorID, hasImpersonator = r.Context().Value(ImpersonatorIDKey).(string)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if gotUserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", gotUserID)
	}
	if !hasImpersonator {
		t.Error("expected ImpersonatorID in context")
	}
	if gotImpersonatorID != "admin-1" {
		t.Errorf("expected ImpersonatorID 'admin-1', got '%s'", gotImpersonatorID)
	}
}

func TestJWTMiddleware_NoImpersonatorID(t *testing.T) {
	mgr := NewJWTManager("test-secret", 24*time.Hour)

	token, _ := mgr.GenerateToken("user-1", "user1", "user")

	var hasImpersonator bool

	handler := JWTMiddleware(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasImpersonator = r.Context().Value(ImpersonatorIDKey).(string)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if hasImpersonator {
		t.Error("expected no ImpersonatorID in context for regular token")
	}
}

func TestJWTMiddleware_EmptyImpersonatorID(t *testing.T) {
	mgr := NewJWTManager("test-secret", 24*time.Hour)

	token, _ := mgr.GenerateTokenWithImpersonator("user-1", "user1", "user", "")

	var hasImpersonator bool

	handler := JWTMiddleware(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasImpersonator = r.Context().Value(ImpersonatorIDKey).(string)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if hasImpersonator {
		t.Error("expected no ImpersonatorID in context for empty impersonator_id")
	}
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	mgr := NewJWTManager("test-secret", 24*time.Hour)

	handler := JWTMiddleware(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	mgr := NewJWTManager("test-secret", 24*time.Hour)

	handler := JWTMiddleware(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestRequireRole_ImpersonatedUser(t *testing.T) {
	handler := RequireRole("admin", "super_admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, RoleKey, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for user role, got %d", w.Code)
	}
}

func TestNoStoreMiddleware(t *testing.T) {
	handler := NoStoreMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/shops", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if got := w.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control: got %q, want %q", got, "no-store")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

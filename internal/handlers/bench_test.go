package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkLogin(b *testing.B) {
	h, _ := setupAuthHandler()
	body, _ := json.Marshal(LoginRequest{Login: "testuser", Password: "Password123"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h.Login(w, req)
	}
}

func BenchmarkLoginParallel(b *testing.B) {
	h, _ := setupAuthHandler()
	body, _ := json.Marshal(LoginRequest{Login: "testuser", Password: "Password123"})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
			w := httptest.NewRecorder()
			h.Login(w, req)
		}
	})
}

func BenchmarkRegister(b *testing.B) {
	h, _ := setupAuthHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body, _ := json.Marshal(RegisterRequest{
			Nickname: "benchuser",
			Password: "Password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h.Register(w, req)
	}
}

func BenchmarkChangePassword(b *testing.B) {
	h, _ := setupAuthHandler()
	body, _ := json.Marshal(ChangePasswordRequest{
		OldPassword: "Password123",
		NewPassword: "NewPass456",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewReader(body))
		req = req.WithContext(withUserID(req.Context(), "user-1"))
		w := httptest.NewRecorder()
		h.ChangePassword(w, req)
	}
}

func BenchmarkWriteJSON(b *testing.B) {
	data := map[string]string{"key": "value", "status": "ok"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		writeJSON(w, 200, data)
	}
}

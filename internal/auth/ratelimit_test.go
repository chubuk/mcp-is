package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowsWithinLimit(t *testing.T) {
	rl := NewRateLimiter(3, 1*time.Minute)

	if !rl.isAllowed("ip1") {
		t.Error("request 1 should be allowed")
	}
	if !rl.isAllowed("ip1") {
		t.Error("request 2 should be allowed")
	}
	if !rl.isAllowed("ip1") {
		t.Error("request 3 should be allowed")
	}
}

func TestRateLimiterBlocksOverLimit(t *testing.T) {
	rl := NewRateLimiter(2, 1*time.Minute)

	rl.isAllowed("ip2")
	rl.isAllowed("ip2")

	if rl.isAllowed("ip2") {
		t.Error("request 3 should be blocked")
	}
}

func TestRateLimiter_SeparateKeys(t *testing.T) {
	rl := NewRateLimiter(2, 1*time.Minute)

	rl.isAllowed("ip-a")
	rl.isAllowed("ip-a")
	rl.isAllowed("ip-b")

	if rl.isAllowed("ip-a") {
		t.Error("ip-a should be blocked after 2 requests")
	}
	if !rl.isAllowed("ip-b") {
		t.Error("ip-b should still be allowed (only 1 request)")
	}
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	rl := NewRateLimiter(1, 10*time.Millisecond)

	rl.isAllowed("ip3")
	rl.isAllowed("ip3")

	time.Sleep(15 * time.Millisecond)

	if !rl.isAllowed("ip3") {
		t.Error("request should be allowed after window expiry")
	}
}

func TestRateLimiter_MiddlewareBlocks(t *testing.T) {
	rl := NewRateLimiter(1, 1*time.Minute)

	called := false
	handler := rl.Middleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "1.2.3.4:1234"

	w1 := httptest.NewRecorder()
	handler(w1, req)
	if !called {
		t.Error("first request should pass through")
	}

	called = false
	w2 := httptest.NewRecorder()
	handler(w2, req)
	if called {
		t.Error("second request should be blocked")
	}
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w2.Code)
	}
}

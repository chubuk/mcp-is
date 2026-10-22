package auth

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Cleanup if too many keys
	if len(rl.requests) > 10000 {
		now := time.Now()
		windowStart := now.Add(-rl.window)
		for k, reqs := range rl.requests {
			if len(reqs) == 0 || reqs[len(reqs)-1].Before(windowStart) {
				delete(rl.requests, k)
			}
		}
	}

	now := time.Now()
	windowStart := now.Add(-rl.window)

	reqs := rl.requests[key]
	filtered := make([]time.Time, 0, len(reqs))
	for _, t := range reqs {
		if t.After(windowStart) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= rl.limit {
		rl.requests[key] = filtered
		return false
	}

	rl.requests[key] = append(filtered, now)
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)
		for key, reqs := range rl.requests {
			filtered := make([]time.Time, 0, len(reqs))
			for _, t := range reqs {
				if t.After(windowStart) {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = filtered
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			// Take first IP (client) from X-Forwarded-For
			if idx := strings.Index(fwd, ","); idx > 0 {
				fwd = strings.TrimSpace(fwd[:idx])
			}
			// Basic validation: must be a valid IP
			if net.ParseIP(fwd) != nil {
				ip = fwd
			}
		}

		if !rl.isAllowed(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"too many requests, try again later"}`))
			return
		}

		next(w, r)
	}
}

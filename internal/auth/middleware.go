package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const (
	UserIDKey         contextKey = "user_id"
	NicknameKey       contextKey = "nickname"
	RoleKey           contextKey = "role"
	ShopIDKey         contextKey = "shop_id"
	KeyNameKey        contextKey = "key_name"
	KeyIDKey          contextKey = "key_id"
	KeyPDNKey         contextKey = "ai_key_pdn_enabled"
	RequestIDKey      contextKey = "request_id"
	ImpersonatorIDKey contextKey = "impersonator_id"
)

var (
	lastKeyUpdate   = make(map[string]time.Time)
	lastKeyUpdateMu sync.Mutex
)

func shouldUpdateLastUsed(keyID string) bool {
	lastKeyUpdateMu.Lock()
	defer lastKeyUpdateMu.Unlock()

	// Cleanup entries older than 10 minutes periodically
	if len(lastKeyUpdate) > 1000 {
		now := time.Now()
		for k, v := range lastKeyUpdate {
			if now.Sub(v) > 10*time.Minute {
				delete(lastKeyUpdate, k)
			}
		}
	}

	last, exists := lastKeyUpdate[keyID]
	if !exists {
		lastKeyUpdate[keyID] = time.Now()
		return true
	}
	if time.Since(last) > 5*time.Minute {
		lastKeyUpdate[keyID] = time.Now()
		return true
	}
	return false
}

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 16)
			rand.Read(b)
			id = hex.EncodeToString(b)
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "error", rec, "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// NoStoreMiddleware marks responses as non-cacheable (JWT/API payloads are sensitive).
func NoStoreMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func JWTMiddleware(jwtMgr *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}

			if tokenStr == "" {
				http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
				return
			}

			claims, err := jwtMgr.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, NicknameKey, claims.Nickname)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			if claims.ImpersonatorID != "" {
				ctx = context.WithValue(ctx, ImpersonatorIDKey, claims.ImpersonatorID)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool)
	for _, r := range roles {
		allowed[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(RoleKey).(string)
			if role == "" || !allowed[role] {
				http.Error(w, `{"error":"insufficient permissions"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func APIKeyMiddleware(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := ""
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
			if apiKey == "" {
				apiKey = r.Header.Get("X-API-Key")
			}

			if apiKey == "" {
				http.Error(w, `{"error":"missing API key"}`, http.StatusUnauthorized)
				return
			}

			keyHash := HashAPIKey(apiKey)

			var userID, keyID, keyName, shopID string
			var aiKeyPDN bool
			err := db.QueryRow(r.Context(),
				`SELECT id, user_id, COALESCE(name,''), pdn_enabled, COALESCE(shop_id::text,'') FROM ai_keys WHERE key_hash = $1`, keyHash,
			).Scan(&keyID, &userID, &keyName, &aiKeyPDN, &shopID)

			if err != nil {
				slog.Warn("invalid api key", "hash", keyHash[:8]+"...")
				http.Error(w, `{"error":"invalid API key"}`, http.StatusUnauthorized)
				return
			}

			if shouldUpdateLastUsed(keyID) {
				go func() {
					if _, err := db.Exec(context.Background(), `UPDATE ai_keys SET last_used_at = NOW() WHERE id = $1`, keyID); err != nil {
						slog.Warn("update ai key last_used_at", "error", err, "key_id", keyID)
					}
				}()
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, KeyNameKey, keyName)
			ctx = context.WithValue(ctx, KeyIDKey, keyID)
			ctx = context.WithValue(ctx, KeyPDNKey, aiKeyPDN)
			// App-installed keys are bound to a shop; personal keys fall back to the default shop.
			if shopID != "" {
				ctx = context.WithValue(ctx, ShopIDKey, shopID)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

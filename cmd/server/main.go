package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-project/mcp-insales/internal/auth"
	"github.com/your-project/mcp-insales/internal/certbot"
	"github.com/your-project/mcp-insales/internal/config"
	"github.com/your-project/mcp-insales/internal/core"
	"github.com/your-project/mcp-insales/internal/email"
	"github.com/your-project/mcp-insales/internal/handlers"
	"github.com/your-project/mcp-insales/internal/logger"
	"github.com/your-project/mcp-insales/internal/mcp"
	"github.com/your-project/mcp-insales/internal/service"
	"github.com/your-project/mcp-insales/internal/storage"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.LogLevel)

	slog.Info("starting mcp-insales server", "port", cfg.ServerPort)

	db, err := storage.NewDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("database connected")

	pool := db.Pool()

	userRepo := storage.NewUserRepository(pool)
	shopRepo := storage.NewShopRepository(pool)
	aiKeyRepo := storage.NewAIKeyRepository(pool)
	settingsRepo := storage.NewSettingsRepository(pool)
	_ = settingsRepo.InitDefaults()

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret)
	aesCipher, err := auth.NewAESCipher(cfg.EncryptionKey)
	if err != nil {
		slog.Error("invalid encryption key", "error", err)
		os.Exit(1)
	}

	manager := core.NewAccountManager(shopRepo, aesCipher, cfg.ReadOnly)
	svc := service.NewService(manager)
	mcpServer := mcp.NewServer(svc, cfg.ReadOnly)

	emailSvc := email.NewService(email.Config{
		Host:   cfg.SMTPHost,
		Port:   cfg.SMTPPort,
		User:   cfg.SMTPUser,
		Pass:   cfg.SMTPPass,
		From:   cfg.SMTPFrom,
		Domain: cfg.Domain,
	})

	authHandler := handlers.NewAuthHandler(userRepo, jwtMgr, emailSvc)
	shopHandler := handlers.NewShopHandler(shopRepo, aesCipher)
	aiKeyHandler := handlers.NewAIKeyHandler(aiKeyRepo)
	adminHandler := handlers.NewAdminHandler(userRepo, shopRepo, aiKeyRepo)
	serverHandler := handlers.NewServerHandler(settingsRepo)

	certbotSvc := certbot.NewCertbotService(settingsRepo)
	scheduler := certbot.NewScheduler(certbotSvc, settingsRepo)

	authRateLimiter := auth.NewRateLimiter(10, 1*time.Minute)

	mux := http.NewServeMux()

	mux.Handle("POST /api/register", authRateLimiter.Middleware(authHandler.Register))
	mux.Handle("POST /api/login", authRateLimiter.Middleware(authHandler.Login))
	mux.Handle("POST /api/forgot-password", authRateLimiter.Middleware(authHandler.ForgotPassword))
	mux.HandleFunc("POST /api/reset-password", authHandler.ResetPassword)

	jwtMW := auth.JWTMiddleware(jwtMgr)
	mux.Handle("POST /api/change-password", jwtMW(http.HandlerFunc(authHandler.ChangePassword)))
	mux.Handle("POST /api/set-email", jwtMW(http.HandlerFunc(authHandler.SetEmail)))

	mux.Handle("POST /api/shops", jwtMW(http.HandlerFunc(shopHandler.Create)))
	mux.Handle("GET /api/shops", jwtMW(http.HandlerFunc(shopHandler.List)))
	mux.Handle("POST /api/shops/test-connection", jwtMW(http.HandlerFunc(shopHandler.TestConnection)))
	mux.Handle("PUT /api/shops/{id}", jwtMW(http.HandlerFunc(shopHandler.Update)))
	mux.Handle("DELETE /api/shops/{id}", jwtMW(http.HandlerFunc(shopHandler.Delete)))

	mux.Handle("POST /api/ai-keys", jwtMW(http.HandlerFunc(aiKeyHandler.Create)))
	mux.Handle("GET /api/ai-keys", jwtMW(http.HandlerFunc(aiKeyHandler.List)))
	mux.Handle("DELETE /api/ai-keys/{id}", jwtMW(http.HandlerFunc(aiKeyHandler.Delete)))

	adminMW := auth.RequireRole("admin", "super_admin")
	mux.Handle("GET /api/admin/users", jwtMW(adminMW(http.HandlerFunc(adminHandler.ListUsers))))
	mux.Handle("GET /api/admin/users/{id}", jwtMW(adminMW(http.HandlerFunc(adminHandler.GetUser))))
	mux.Handle("PUT /api/admin/users/{id}/role", jwtMW(auth.RequireRole("super_admin")(http.HandlerFunc(adminHandler.UpdateRole))))
	mux.Handle("PUT /api/admin/users/{id}/status", jwtMW(adminMW(http.HandlerFunc(adminHandler.UpdateStatus))))
	mux.Handle("DELETE /api/admin/users/{id}", jwtMW(adminMW(http.HandlerFunc(adminHandler.DeleteUser))))
	mux.Handle("GET /api/admin/stats", jwtMW(adminMW(http.HandlerFunc(adminHandler.Stats))))

	superAdminMW := auth.RequireRole("super_admin")
	mux.Handle("GET /api/admin/server", jwtMW(adminMW(http.HandlerFunc(serverHandler.GetSettings))))
	mux.Handle("PUT /api/admin/server", jwtMW(superAdminMW(http.HandlerFunc(serverHandler.UpdateSettings))))
	mux.Handle("POST /api/admin/server/cert-issue", jwtMW(adminMW(http.HandlerFunc(serverHandler.IssueCert))))
	mux.Handle("POST /api/admin/server/cert-renew", jwtMW(adminMW(http.HandlerFunc(serverHandler.RenewCert))))

	apiKeyMW := auth.APIKeyMiddleware(pool)
	mux.Handle("POST /mcp", apiKeyMW(http.HandlerFunc(mcpServer.HandleMCP)))
	mux.Handle("GET /mcp", apiKeyMW(http.HandlerFunc(mcpServer.HandleSSE)))

	mux.HandleFunc("POST /webhook/{event}", mcpServer.HandleWebhook)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	httpServer := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      auth.CORSMiddleware(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("http server listening", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	scheduler.Start()

	<-ctx.Done()
	slog.Info("shutting down...")

	scheduler.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}

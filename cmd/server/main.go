package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/certbot"
	"github.com/chubuk/mcp-is/internal/config"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/email"
	"github.com/chubuk/mcp-is/internal/handlers"
	"github.com/chubuk/mcp-is/internal/logger"
	"github.com/chubuk/mcp-is/internal/mcp"
	"github.com/chubuk/mcp-is/internal/metrics"
	"github.com/chubuk/mcp-is/internal/service"
	"github.com/chubuk/mcp-is/internal/storage"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.LogLevel)

	slog.Info("starting mcp-is server", "port", cfg.ServerPort)

	db, err := storage.NewDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("database connected")

	pool := db.Pool()

	storage.RunMigrations(pool)

	userRepo := storage.NewUserRepository(pool)
	shopRepo := storage.NewShopRepository(pool)
	aiKeyRepo := storage.NewAIKeyRepository(pool)
	settingsRepo := storage.NewSettingsRepository(pool)
	_ = settingsRepo.SetDefaultsIfMissing()

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)
	aesCipher, err := auth.NewAESCipher(cfg.EncryptionKey)
	if err != nil {
		slog.Error("invalid encryption key", "error", err)
		os.Exit(1)
	}

	manager := core.NewAccountManager(shopRepo, aesCipher, cfg.ReadOnly)
	svc := service.NewService(manager)
	mcpServer := mcp.NewMCPServer(svc, cfg.ReadOnly, cfg.PDNAutoProtect, cfg.ServerCountry)

	emailSvc := email.NewService(email.Config{
		Host:   cfg.SMTPHost,
		Port:   cfg.SMTPPort,
		User:   cfg.SMTPUser,
		Pass:   cfg.SMTPPass,
		From:   cfg.SMTPFrom,
		Domain: cfg.Domain,
	})

	authHandler := handlers.NewAuthHandler(userRepo, jwtMgr, emailSvc)
	shopHandler := handlers.NewShopHandler(shopRepo, aiKeyRepo, aesCipher, mcpServer)
	aiKeyHandler := handlers.NewAIKeyHandler(aiKeyRepo)
	adminHandler := handlers.NewAdminHandler(userRepo, shopRepo, aiKeyRepo, jwtMgr)
	serverHandler := handlers.NewServerHandler(settingsRepo)
	pdnHandler := handlers.NewPDNHandler(pool)

	certbotSvc := certbot.NewCertbotService(settingsRepo)
	scheduler := certbot.NewScheduler(certbotSvc, settingsRepo)

	authRateLimiter := auth.NewRateLimiter(10, 1*time.Minute)

	mux := http.NewServeMux()

	mux.Handle("POST /api/register", authRateLimiter.Middleware(authHandler.Register))
	mux.Handle("POST /api/login", authRateLimiter.Middleware(authHandler.Login))
	mux.Handle("POST /api/forgot-password", authRateLimiter.Middleware(authHandler.ForgotPassword))
	mux.Handle("POST /api/reset-password", authRateLimiter.Middleware(authHandler.ResetPassword))

	jwtMW := auth.JWTMiddleware(jwtMgr)
	mux.Handle("POST /api/change-password", jwtMW(http.HandlerFunc(authHandler.ChangePassword)))
	mux.Handle("POST /api/force-change-password", jwtMW(http.HandlerFunc(authHandler.ForceChangePassword)))
	mux.Handle("POST /api/set-email", jwtMW(http.HandlerFunc(authHandler.SetEmail)))

	mux.Handle("POST /api/shops", jwtMW(http.HandlerFunc(shopHandler.Create)))
	mux.Handle("GET /api/shops", jwtMW(http.HandlerFunc(shopHandler.List)))
	mux.Handle("POST /api/shops/test-connection", jwtMW(http.HandlerFunc(shopHandler.TestConnection)))
	mux.Handle("PUT /api/shops/{id}", jwtMW(http.HandlerFunc(shopHandler.Update)))
	mux.Handle("DELETE /api/shops/{id}", jwtMW(http.HandlerFunc(shopHandler.Delete)))
	mux.Handle("POST /api/shops/{id}/check-commands", jwtMW(http.HandlerFunc(shopHandler.CheckCommands)))

	mux.Handle("POST /api/ai-keys", jwtMW(http.HandlerFunc(aiKeyHandler.Create)))
	mux.Handle("GET /api/ai-keys", jwtMW(http.HandlerFunc(aiKeyHandler.List)))
	mux.Handle("PUT /api/ai-keys/{id}", jwtMW(http.HandlerFunc(aiKeyHandler.Update)))
	mux.Handle("DELETE /api/ai-keys/{id}", jwtMW(http.HandlerFunc(aiKeyHandler.Delete)))

	adminMW := auth.RequireRole("admin", "super_admin")
	mux.Handle("GET /api/admin/users", jwtMW(adminMW(http.HandlerFunc(adminHandler.ListUsers))))
	mux.Handle("GET /api/admin/users/{id}", jwtMW(adminMW(http.HandlerFunc(adminHandler.GetUser))))
	mux.Handle("PUT /api/admin/users/{id}/role", jwtMW(auth.RequireRole("super_admin")(http.HandlerFunc(adminHandler.UpdateRole))))
	mux.Handle("PUT /api/admin/users/{id}/status", jwtMW(adminMW(http.HandlerFunc(adminHandler.UpdateStatus))))
	mux.Handle("DELETE /api/admin/users/{id}", jwtMW(adminMW(http.HandlerFunc(adminHandler.DeleteUser))))
	mux.Handle("GET /api/admin/stats", jwtMW(adminMW(http.HandlerFunc(adminHandler.Stats))))
	mux.Handle("POST /api/admin/users/{id}/impersonate", jwtMW(adminMW(http.HandlerFunc(adminHandler.Impersonate))))
	mux.Handle("POST /api/admin/stop-impersonate", jwtMW(http.HandlerFunc(adminHandler.StopImpersonate)))

	superAdminMW := auth.RequireRole("super_admin")
	mux.Handle("GET /api/admin/server", jwtMW(adminMW(http.HandlerFunc(serverHandler.GetSettings))))
	mux.Handle("PUT /api/admin/server", jwtMW(superAdminMW(http.HandlerFunc(serverHandler.UpdateSettings))))
	mux.Handle("POST /api/admin/server/cert-issue", jwtMW(adminMW(http.HandlerFunc(serverHandler.IssueCert))))
	mux.Handle("POST /api/admin/server/cert-renew", jwtMW(adminMW(http.HandlerFunc(serverHandler.RenewCert))))

	mux.Handle("GET /api/admin/pdn", jwtMW(adminMW(http.HandlerFunc(pdnHandler.GetSettings))))
	mux.Handle("PUT /api/admin/pdn", jwtMW(superAdminMW(http.HandlerFunc(pdnHandler.UpdateSettings))))
	mux.Handle("GET /api/admin/pdn/log", jwtMW(adminMW(http.HandlerFunc(pdnHandler.GetLog))))
	mux.Handle("GET /api/admin/pdn/stats", jwtMW(adminMW(http.HandlerFunc(pdnHandler.GetStats))))
	mux.Handle("GET /api/admin/pdn/jurisdictions", jwtMW(adminMW(http.HandlerFunc(pdnHandler.GetJurisdictions))))

	backupHandler := handlers.NewBackupHandler(pool)
	mcpClientTypeRepo := storage.NewMCPClientTypeRepository(pool)
	_ = mcpClientTypeRepo.Seed()
	mcpClientTypeHandler := handlers.NewMCPClientTypeHandler(mcpClientTypeRepo)
	mux.Handle("GET /api/admin/backup", jwtMW(adminMW(http.HandlerFunc(backupHandler.GetSettings))))
	mux.Handle("PUT /api/admin/backup", jwtMW(superAdminMW(http.HandlerFunc(backupHandler.UpdateSettings))))
	mux.Handle("GET /api/admin/backup/log", jwtMW(adminMW(http.HandlerFunc(backupHandler.GetLog))))
	mux.Handle("GET /api/admin/backup/incidents", jwtMW(adminMW(http.HandlerFunc(backupHandler.GetIncidents))))
	mux.Handle("GET /api/admin/backup/status", jwtMW(adminMW(http.HandlerFunc(backupHandler.GetStatus))))
	mux.Handle("POST /api/admin/backup/run", jwtMW(adminMW(http.HandlerFunc(backupHandler.RunBackup))))

	mux.HandleFunc("GET /api/mcp-client-types", mcpClientTypeHandler.List)

	apiKeyMW := auth.APIKeyMiddleware(pool)
	mux.Handle("/mcp", apiKeyMW(http.HandlerFunc(mcpServer.ServeHTTP)))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.Handle("GET /metrics", metrics.Handler())

	startupTime := time.Now()
	metrics.StartupTime.Set(time.Since(startupTime).Seconds())

	httpServer := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      auth.Chain(mux, auth.RecoveryMiddleware, auth.RequestIDMiddleware, auth.CORSMiddleware(cfg.AllowedOrigins)),
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

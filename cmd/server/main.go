package main

import (
	"context"
	"log/slog"
	"net/http"
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
	"github.com/chubuk/pagetree"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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

	manager := core.NewAccountManager(shopRepo, aesCipher, cfg.ReadOnly, cfg.BillingMonthly, settingsRepo)
	svc := service.NewService(manager)
	auditRepo := storage.NewAuditRepository(pool)
	idemRepo := storage.NewIdempotencyRepository(pool)
	mcp.SetPDNLog(storage.NewPDNLogRepository(pool), pool)
	mcpServer := mcp.NewMCPServer(svc, cfg.ReadOnly, cfg.PDNAutoProtect, cfg.ServerCountry, auditRepo, idemRepo)

	emailCfg := handlers.LoadSMTPConfigFromSettings(settingsRepo, email.Config{
		Host:   cfg.SMTPHost,
		Port:   cfg.SMTPPort,
		User:   cfg.SMTPUser,
		Pass:   cfg.SMTPPass,
		From:   cfg.SMTPFrom,
		Domain: cfg.Domain,
	}, aesCipher)
	emailSvc := email.NewService(emailCfg)

	tplRepo := storage.NewEmailTemplateRepository(pool)
	handlers.SeedEmailTemplates(tplRepo)
	emailSvc.SetTemplateStore(email.FuncStore{
		GetFunc: func(key string) (email.Template, error) {
			t, err := tplRepo.GetByKey(key)
			if err != nil {
				return email.Default(key)
			}
			return email.Template{
				Key: t.Key, Name: t.Name, Subject: t.Subject, HTMLBody: t.HTMLBody,
				TextBody: t.TextBody, Description: t.Description, RequiredVars: t.RequiredVars,
				IsModified: t.IsModified,
			}, nil
		},
		ListFunc: func() ([]email.Template, error) {
			items, err := tplRepo.List()
			if err != nil {
				return nil, err
			}
			out := make([]email.Template, 0, len(items))
			for _, t := range items {
				out = append(out, email.Template{
					Key: t.Key, Name: t.Name, Subject: t.Subject, HTMLBody: t.HTMLBody,
					TextBody: t.TextBody, Description: t.Description, RequiredVars: t.RequiredVars,
				})
			}
			return out, nil
		},
	})
	emailAdmin := handlers.NewEmailAdminHandler(tplRepo, settingsRepo, emailSvc, aesCipher)

	authHandler := handlers.NewAuthHandler(userRepo, jwtMgr, emailSvc, cfg.RegistrationEnabled)
	shopHandler := handlers.NewShopHandler(shopRepo, aiKeyRepo, aesCipher, mcpServer, cfg.MaxShopsPerUser)
	aiKeyHandler := handlers.NewAIKeyHandler(aiKeyRepo, shopRepo, cfg.MaxAIKeysPerUser)
	adminHandler := handlers.NewAdminHandler(userRepo, shopRepo, aiKeyRepo, jwtMgr, emailSvc)
	serverHandler := handlers.NewServerHandler(settingsRepo, cfg.BillingMonthly)
	pdnHandler := handlers.NewPDNHandler(pool, settingsRepo)

	appInstallRepo := storage.NewAppInstallRepository(pool)
	appHandler := handlers.NewAppHandler(
		appInstallRepo, userRepo, shopRepo, aiKeyRepo,
		aesCipher, jwtMgr, manager,
		cfg.InSalesAppKey, cfg.InSalesSecretKey, cfg.Domain,
		cfg.BillingMonthly, settingsRepo,
	)
	appHandler.SetEmailService(emailSvc)
	if cfg.AppModeEnabled() {
		slog.Info("InSales app mode enabled", "app_key", cfg.InSalesAppKey)
	} else {
		slog.Info("InSales app mode disabled (set INSALES_APP_KEY and INSALES_SECRET_KEY)")
	}

	certbotSvc := certbot.NewCertbotService(settingsRepo)
	scheduler := certbot.NewScheduler(certbotSvc, settingsRepo)

	authRateLimiter := auth.NewRateLimiter(10, 1*time.Minute)
	// Флуд-защита почтовых операций (регистрация, коды, сброс пароля): 5/час на IP.
	mailRateLimiter := auth.NewRateLimiter(5, 1*time.Hour)

	mux := http.NewServeMux()

	mux.Handle("POST /api/register", mailRateLimiter.Middleware(authHandler.Register))
	mux.Handle("POST /api/login", authRateLimiter.Middleware(authHandler.Login))
	mux.Handle("POST /api/forgot-password", mailRateLimiter.Middleware(authHandler.ForgotPassword))
	mux.Handle("POST /api/reset-password", authRateLimiter.Middleware(authHandler.ResetPassword))
	mux.Handle("POST /api/verify-email", mailRateLimiter.Middleware(authHandler.VerifyEmail))
	mux.Handle("POST /api/verify-email/resend", mailRateLimiter.Middleware(authHandler.ResendVerifyEmail))

	// InSales external-app callbacks (public)
	mux.Handle("GET /install", authRateLimiter.Middleware(appHandler.Install))
	mux.Handle("GET /uninstall", authRateLimiter.Middleware(appHandler.Uninstall))
	mux.Handle("GET /app", authRateLimiter.Middleware(appHandler.AppLanding))
	mux.Handle("GET /autologin", authRateLimiter.Middleware(appHandler.Autologin))
	mux.Handle("POST /api/exchange-code", authRateLimiter.Middleware(appHandler.ExchangeCode))

	jwtMW := auth.JWTMiddleware(jwtMgr)
	mux.Handle("POST /api/change-password", jwtMW(http.HandlerFunc(authHandler.ChangePassword)))
	mux.Handle("POST /api/force-change-password", jwtMW(http.HandlerFunc(authHandler.ForceChangePassword)))
	mux.Handle("POST /api/shops/link-pending", jwtMW(http.HandlerFunc(appHandler.LinkPending)))
	mux.Handle("POST /api/link-code/send", mailRateLimiter.Middleware(appHandler.LinkCodeSend))
	mux.Handle("POST /api/link-code/verify", mailRateLimiter.Middleware(appHandler.LinkCodeVerify))

	mux.Handle("POST /api/shops", jwtMW(http.HandlerFunc(shopHandler.Create)))
	mux.Handle("GET /api/shops", jwtMW(http.HandlerFunc(shopHandler.List)))
	mux.Handle("POST /api/shops/test-connection", jwtMW(http.HandlerFunc(shopHandler.TestConnection)))
	mux.Handle("PUT /api/shops/{id}", jwtMW(http.HandlerFunc(shopHandler.Update)))
	mux.Handle("DELETE /api/shops/{id}", jwtMW(http.HandlerFunc(shopHandler.Delete)))
	mux.Handle("POST /api/shops/{id}/check-access", jwtMW(http.HandlerFunc(shopHandler.CheckAccess)))

	mux.Handle("GET /api/tool-meta", jwtMW(http.HandlerFunc(handlers.ToolMeta)))
	mux.Handle("POST /api/ai-keys", jwtMW(http.HandlerFunc(aiKeyHandler.Create)))
	mux.Handle("GET /api/ai-keys", jwtMW(http.HandlerFunc(aiKeyHandler.List)))
	mux.Handle("PUT /api/ai-keys/{id}", jwtMW(http.HandlerFunc(aiKeyHandler.Update)))
	mux.Handle("POST /api/ai-keys/{id}/rotate", jwtMW(http.HandlerFunc(aiKeyHandler.Rotate)))
	mux.Handle("DELETE /api/ai-keys/{id}", jwtMW(http.HandlerFunc(aiKeyHandler.Delete)))

	adminMW := auth.RequireRole("admin", "super_admin")
	mux.Handle("POST /api/admin/users", jwtMW(adminMW(http.HandlerFunc(adminHandler.CreateUser))))
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

	mux.Handle("GET /api/admin/smtp/status", jwtMW(adminMW(http.HandlerFunc(adminHandler.SMTPStatus))))
	mux.Handle("POST /api/admin/smtp/test", jwtMW(adminMW(http.HandlerFunc(adminHandler.TestSMTP))))
	mux.Handle("GET /api/admin/smtp/settings", jwtMW(adminMW(http.HandlerFunc(emailAdmin.GetSMTPSettings))))
	mux.Handle("PUT /api/admin/smtp/settings", jwtMW(adminMW(http.HandlerFunc(emailAdmin.UpdateSMTPSettings))))

	mux.Handle("GET /api/admin/email/templates", jwtMW(adminMW(http.HandlerFunc(emailAdmin.ListTemplates))))
	mux.Handle("GET /api/admin/email/templates/{key}", jwtMW(adminMW(http.HandlerFunc(emailAdmin.GetTemplate))))
	mux.Handle("PUT /api/admin/email/templates/{key}", jwtMW(adminMW(http.HandlerFunc(emailAdmin.SaveTemplate))))
	mux.Handle("POST /api/admin/email/templates/{key}/reset", jwtMW(adminMW(http.HandlerFunc(emailAdmin.ResetTemplate))))
	mux.Handle("POST /api/admin/email/templates/{key}/preview", jwtMW(adminMW(http.HandlerFunc(emailAdmin.PreviewTemplate))))
	mux.Handle("POST /api/admin/email/templates/{key}/test", jwtMW(adminMW(http.HandlerFunc(emailAdmin.TestTemplate))))

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
	mux.Handle("GET /api/admin/backup/files", jwtMW(adminMW(http.HandlerFunc(backupHandler.ListFiles))))
	mux.Handle("POST /api/admin/backup/files/{name}/verify", jwtMW(adminMW(http.HandlerFunc(backupHandler.VerifyFile))))
	mux.Handle("POST /api/admin/backup/restore", jwtMW(superAdminMW(http.HandlerFunc(backupHandler.Restore))))
	mux.Handle("GET /api/admin/backup/restore/status", jwtMW(adminMW(http.HandlerFunc(backupHandler.GetRestoreStatus))))

	// PageTree — file-based markdown pages/news engine
	ptContentDir := getEnvDefault("PAGETREE_CONTENT_DIR", "./content")
	ptDataDir := getEnvDefault("PAGETREE_DATA_DIR", "./data")
	ptSiteURL := "https://" + cfg.Domain
	var ptMasterKey, ptPrevKey []byte
	if mk := os.Getenv("PAGETREE_MASTER_KEY"); mk != "" {
		ptMasterKey = []byte(mk)
	}
	if pk := os.Getenv("PAGETREE_MASTER_KEY_PREVIOUS"); pk != "" {
		ptPrevKey = []byte(pk)
	}
	ptRenderer := &mdRenderer{}
	ptEngine, ptErr := pagetree.New(pagetree.Config{
		Root:     ptContentDir,
		DataDir:  ptDataDir,
		SiteURL:  ptSiteURL,
		Renderer: ptRenderer,
		Log:      slog.Default(),
		MasterKey: ptMasterKey,
		PrevKey:   ptPrevKey,
	})
	if ptErr != nil {
		slog.Warn("pagetree init failed (content features disabled)", "error", ptErr)
	} else {
		// Public — no auth
		mux.Handle("GET /pages/", ptEngine.PageHandler())
		mux.Handle("GET /news/", ptEngine.NewsHandler())
		mux.Handle("GET /media/", ptEngine.MediaHandler())
		mux.Handle("GET /api/content/menu", ptEngine.MenuHandler())
		// Admin — JWT + role
		mux.Handle("/api/admin/content/", jwtMW(adminMW(ptEngine.ContentAPI())))
		mux.Handle("/api/admin/content/keys/", jwtMW(adminMW(ptEngine.KeysAPI())))
		// Agent API — Bearer ptk_* (module-owned auth)
		mux.Handle("/api/agent/v1/", ptEngine.AgentAPI())
		mux.Handle("/api/agent/mcp", ptEngine.AgentMCPHandler())
		// Admin UI — static shell (data via XHR with JWT)
		mux.Handle("/content/admin/", ptEngine.AdminUIHandler())
		slog.Info("pagetree mounted", "content", ptContentDir, "site_url", ptSiteURL)
	}

	mux.HandleFunc("GET /api/mcp-client-types", mcpClientTypeHandler.List)

	apiKeyMW := auth.APIKeyMiddleware(pool)
	mux.Handle("/mcp", apiKeyMW(http.HandlerFunc(mcpServer.ServeHTTP)))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.Handle("GET /metrics", metrics.Handler())

	startupTime := time.Now()
	metrics.StartupTime.Set(time.Since(startupTime).Seconds())

	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: auth.Chain(mux, auth.RecoveryMiddleware, auth.RequestIDMiddleware, auth.NoStoreMiddleware, auth.CORSMiddleware(cfg.AllowedOrigins)),
		// 20-МБ image_data: клиент → nginx → app → InSales укладывается не всегда
		// в 30с (замер: полный update_blog_post с аплоадом ~37с) — WriteTimeout
		// 30с рвал соединение уже после успешного вызова, ответ не доходил.
		// ReadHeaderTimeout держит защиту от slowloris.
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
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

// mdRenderer adapts gomarkdown to pagetree.Renderer interface.
type mdRenderer struct{}

func (r *mdRenderer) Render(md []byte) ([]byte, error) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer), nil
}

func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

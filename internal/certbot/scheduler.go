package certbot

import (
	"log/slog"
	"time"

	"github.com/your-project/mcp-insales/internal/storage"
)

type Scheduler struct {
	certbot      *CertbotService
	settingsRepo *storage.SettingsRepository
	stopCh       chan struct{}
}

func NewScheduler(certbot *CertbotService, settingsRepo *storage.SettingsRepository) *Scheduler {
	return &Scheduler{
		certbot:      certbot,
		settingsRepo: settingsRepo,
		stopCh:       make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	slog.Info("certificate scheduler started", "interval", CertCheckInterval)

	go func() {
		ticker := time.NewTicker(CertCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.checkAndRenew()
			case <-s.stopCh:
				slog.Info("certificate scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) checkAndRenew() {
	autoRenew := s.settingsRepo.GetDefault("auto_renew", "true")
	if autoRenew != "true" {
		return
	}

	domain := s.settingsRepo.GetDefault("domain", "localhost")
	if domain == "localhost" {
		return
	}

	needsAttention := s.settingsRepo.GetDefault("cert_needs_attention", "false")
	if needsAttention == "true" {
		if !s.certbot.ShouldRetry() {
			slog.Debug("cert needs attention but retry limit reached", "domain", domain)
			return
		}
		s.attemptRenew(domain)
		return
	}

	needsRenewal, daysLeft, err := s.certbot.CheckRenewal(domain)
	if err != nil {
		slog.Warn("failed to check certificate", "domain", domain, "error", err)
		return
	}

	if !needsRenewal {
		slog.Debug("certificate is fine", "domain", domain, "days_left", daysLeft)
		return
	}

	slog.Info("certificate needs renewal", "domain", domain, "days_left", daysLeft)
	s.attemptRenew(domain)
}

func (s *Scheduler) attemptRenew(domain string) {
	slog.Info("attempting certificate renewal", "domain", domain)

	err := s.certbot.RenewCert(domain)
	if err != nil {
		retryCount := s.settingsRepo.GetDefault("retry_count", "0")
		slog.Error("certificate renewal failed",
			"domain", domain,
			"error", err,
			"retry_count", retryCount,
		)
		return
	}

	slog.Info("certificate renewed successfully", "domain", domain)
}

package certbot

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	MaxRetryCount     = 10
	RetryInterval     = 7 * time.Hour
	CertCheckInterval = 1 * time.Hour
)

type CertInfo struct {
	Domain    string
	Issuer    string
	ExpiresAt time.Time
	IsValid   bool
}

// SettingsStore is the interface certbot needs from the settings repository.
type SettingsStore interface {
	GetDefault(key, fallback string) string
	Set(key, value string) error
}

type CertbotService struct {
	settingsRepo SettingsStore
	certPath     string
	keyPath      string
}

func NewCertbotService(settingsRepo SettingsStore) *CertbotService {
	return &CertbotService{
		settingsRepo: settingsRepo,
		certPath:     "/etc/letsencrypt/live",
		keyPath:      "/etc/letsencrypt/live",
	}
}

func (s *CertbotService) IssueCert(domain, email string) error {
	slog.Info("issuing certificate via certbot", "domain", domain, "email", email)

	cmd := exec.Command("docker", "exec", "mcp_is_nginx",
		"certbot", "certonly", "--webroot",
		"--webroot-path=/var/www/certbot",
		"-d", domain,
		"--email", email,
		"--agree-tos", "--non-interactive",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("certbot issue failed", "error", err, "output", string(output))
		s.settingsRepo.Set("last_renew_error", string(output))
		return fmt.Errorf("certbot issue: %w", err)
	}

	slog.Info("certificate issued successfully", "domain", domain)
	s.settingsRepo.Set("cert_issuer", "Let's Encrypt")
	s.settingsRepo.Set("cert_needs_attention", "false")
	s.settingsRepo.Set("last_renew_error", "")
	s.settingsRepo.Set("retry_count", "0")

	return nil
}

func (s *CertbotService) RenewCert(domain string) error {
	slog.Info("renewing certificate via certbot", "domain", domain)

	cmd := exec.Command("docker", "exec", "mcp_is_nginx",
		"certbot", "renew", "--non-interactive",
		// certbot в неинтерактивном режиме спит случайное время (до ~8 мин) перед renew
		"--no-random-sleep-on-renew",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("certbot renew failed", "error", err, "output", string(output))
		s.incrementRetry(string(output))
		return fmt.Errorf("certbot renew: %w", err)
	}

	slog.Info("certificate renewed successfully", "domain", domain)
	s.settingsRepo.Set("cert_needs_attention", "false")
	s.settingsRepo.Set("last_renew_error", "")
	s.settingsRepo.Set("retry_count", "0")
	s.settingsRepo.Set("cert_expires_at", time.Now().Add(90*24*time.Hour).Format(time.RFC3339))

	s.reloadNginx()
	return nil
}

func (s *CertbotService) GetCertInfo(domain string) (*CertInfo, error) {
	certFile := filepath.Join(s.certPath, domain, "fullchain.pem")

	data, err := os.ReadFile(certFile)
	if err != nil {
		return &CertInfo{
			Domain:  domain,
			Issuer:  "Self-signed",
			IsValid: false,
		}, nil
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &CertInfo{
		Domain:    domain,
		Issuer:    cert.Issuer.CommonName,
		ExpiresAt: cert.NotAfter,
		IsValid:   time.Now().Before(cert.NotAfter),
	}, nil
}

func (s *CertbotService) CheckRenewal(domain string) (needsRenewal bool, daysLeft int, err error) {
	info, err := s.GetCertInfo(domain)
	if err != nil {
		return false, 0, err
	}

	daysLeft = int(time.Until(info.ExpiresAt).Hours() / 24)
	renewDaysStr := s.settingsRepo.GetDefault("renew_days_before", "30")
	renewDays := 30
	fmt.Sscanf(renewDaysStr, "%d", &renewDays)

	return daysLeft <= renewDays, daysLeft, nil
}

func (s *CertbotService) ShouldRetry() bool {
	retryCountStr := s.settingsRepo.GetDefault("retry_count", "0")
	retryCount := 0
	fmt.Sscanf(retryCountStr, "%d", &retryCount)

	if retryCount >= MaxRetryCount {
		return false
	}

	lastRetryStr := s.settingsRepo.GetDefault("last_retry_at", "")
	if lastRetryStr == "" {
		return true
	}

	lastRetry, err := time.Parse(time.RFC3339, lastRetryStr)
	if err != nil {
		return true
	}

	return time.Since(lastRetry) >= RetryInterval
}

func (s *CertbotService) incrementRetry(errorMsg string) {
	retryCountStr := s.settingsRepo.GetDefault("retry_count", "0")
	retryCount := 0
	fmt.Sscanf(retryCountStr, "%d", &retryCount)
	retryCount++

	s.settingsRepo.Set("retry_count", fmt.Sprintf("%d", retryCount))
	s.settingsRepo.Set("last_retry_at", time.Now().Format(time.RFC3339))
	s.settingsRepo.Set("last_renew_error", errorMsg)

	if retryCount >= MaxRetryCount {
		s.settingsRepo.Set("cert_needs_attention", "true")
		slog.Error("certificate renewal exhausted all retries", "count", retryCount)
	}
}

func (s *CertbotService) reloadNginx() {
	cmd := exec.Command("docker", "exec", "mcp_is_nginx", "nginx", "-s", "reload")
	if err := cmd.Run(); err != nil {
		slog.Error("failed to reload nginx", "error", err)
	} else {
		slog.Info("nginx reloaded successfully")
	}
}

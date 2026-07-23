package certbot

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type mockSettings struct {
	data map[string]string
}

func newMockSettings() *mockSettings {
	return &mockSettings{data: make(map[string]string)}
}

func (m *mockSettings) GetDefault(key, fallback string) string {
	if v, ok := m.data[key]; ok {
		return v
	}
	return fallback
}

func (m *mockSettings) Set(key, value string) error {
	m.data[key] = value
	return nil
}

func generateTestCert(t *testing.T, dir string, domain string, notAfter time.Time) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: domain},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     notAfter,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}

	certPath := filepath.Join(dir, domain)
	if err := os.MkdirAll(certPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	certFile := filepath.Join(certPath, "fullchain.pem")
	f, err := os.Create(certFile)
	if err != nil {
		t.Fatalf("create cert file: %v", err)
	}
	defer f.Close()
	pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
}

func TestGetCertInfo_ValidCert(t *testing.T) {
	dir := t.TempDir()
	domain := "example.com"
	expiresAt := time.Now().Add(80 * 24 * time.Hour)
	generateTestCert(t, dir, domain, expiresAt)

	svc := &CertbotService{certPath: dir}
	info, err := svc.GetCertInfo(domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Issuer != domain {
		t.Errorf("expected issuer %q, got %q", domain, info.Issuer)
	}
	if !info.IsValid {
		t.Error("expected IsValid=true")
	}
	if info.ExpiresAt.Before(time.Now()) {
		t.Error("expected expiry in the future")
	}
}

func TestGetCertInfo_ExpiredCert(t *testing.T) {
	dir := t.TempDir()
	domain := "expired.com"
	expiresAt := time.Now().Add(-24 * time.Hour)
	generateTestCert(t, dir, domain, expiresAt)

	svc := &CertbotService{certPath: dir}
	info, err := svc.GetCertInfo(domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.IsValid {
		t.Error("expected IsValid=false for expired cert")
	}
	if info.Issuer != domain {
		t.Errorf("expected issuer %q, got %q", domain, info.Issuer)
	}
}

func TestGetCertInfo_SelfSigned(t *testing.T) {
	dir := t.TempDir()

	svc := &CertbotService{certPath: dir}
	info, err := svc.GetCertInfo("missing.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Issuer != "Self-signed" {
		t.Errorf("expected issuer %q, got %q", "Self-signed", info.Issuer)
	}
	if info.IsValid {
		t.Error("expected IsValid=false for missing cert")
	}
}

func TestGetCertInfo_InvalidPEM(t *testing.T) {
	dir := t.TempDir()
	domain := "bad.com"
	certPath := filepath.Join(dir, domain)
	os.MkdirAll(certPath, 0755)
	os.WriteFile(filepath.Join(certPath, "fullchain.pem"), []byte("not a cert"), 0644)

	svc := &CertbotService{certPath: dir}
	_, err := svc.GetCertInfo(domain)
	if err == nil {
		t.Error("expected error for invalid PEM")
	}
}

func TestCheckRenewal_NeedsRenewal(t *testing.T) {
	dir := t.TempDir()
	domain := "renew.com"
	expiresAt := time.Now().Add(10 * 24 * time.Hour)
	generateTestCert(t, dir, domain, expiresAt)

	settings := newMockSettings()
	settings.Set("renew_days_before", "30")

	svc := &CertbotService{settingsRepo: settings, certPath: dir}
	needsRenewal, daysLeft, err := svc.CheckRenewal(domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !needsRenewal {
		t.Error("expected needsRenewal=true (10 days < 30 threshold)")
	}
	if daysLeft < 9 || daysLeft > 10 {
		t.Errorf("expected ~10 days left, got %d", daysLeft)
	}
}

func TestCheckRenewal_NoRenewalNeeded(t *testing.T) {
	dir := t.TempDir()
	domain := "fine.com"
	expiresAt := time.Now().Add(60 * 24 * time.Hour)
	generateTestCert(t, dir, domain, expiresAt)

	settings := newMockSettings()
	settings.Set("renew_days_before", "30")

	svc := &CertbotService{settingsRepo: settings, certPath: dir}
	needsRenewal, daysLeft, err := svc.CheckRenewal(domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if needsRenewal {
		t.Error("expected needsRenewal=false (60 days > 30 threshold)")
	}
	if daysLeft < 59 || daysLeft > 60 {
		t.Errorf("expected ~60 days left, got %d", daysLeft)
	}
}

func TestShouldRetry_NoRetries(t *testing.T) {
	settings := newMockSettings()
	svc := &CertbotService{settingsRepo: settings}

	if !svc.ShouldRetry() {
		t.Error("expected ShouldRetry=true when no retries recorded")
	}
}

func TestShouldRetry_UnderLimit(t *testing.T) {
	settings := newMockSettings()
	settings.Set("retry_count", "5")
	settings.Set("last_retry_at", time.Now().Add(-8*time.Hour).Format(time.RFC3339))

	svc := &CertbotService{settingsRepo: settings}

	if !svc.ShouldRetry() {
		t.Error("expected ShouldRetry=true (5 < 10, last retry > 7h ago)")
	}
}

func TestShouldRetry_AtLimit(t *testing.T) {
	settings := newMockSettings()
	settings.Set("retry_count", "10")

	svc := &CertbotService{settingsRepo: settings}

	if svc.ShouldRetry() {
		t.Error("expected ShouldRetry=false (10 >= 10)")
	}
}

func TestShouldRetry_RecentRetry(t *testing.T) {
	settings := newMockSettings()
	settings.Set("retry_count", "3")
	settings.Set("last_retry_at", time.Now().Format(time.RFC3339))

	svc := &CertbotService{settingsRepo: settings}

	if svc.ShouldRetry() {
		t.Error("expected ShouldRetry=false (last retry too recent)")
	}
}

func TestIncrementRetry(t *testing.T) {
	settings := newMockSettings()
	svc := &CertbotService{settingsRepo: settings}

	svc.incrementRetry("test error")

	if settings.data["retry_count"] != "1" {
		t.Errorf("expected retry_count=1, got %q", settings.data["retry_count"])
	}
	if settings.data["last_renew_error"] != "test error" {
		t.Errorf("expected error msg preserved, got %q", settings.data["last_renew_error"])
	}
	if settings.data["last_retry_at"] == "" {
		t.Error("expected last_retry_at to be set")
	}
}

func TestIncrementRetry_ExhaustsRetries(t *testing.T) {
	settings := newMockSettings()
	settings.Set("retry_count", "9")

	svc := &CertbotService{settingsRepo: settings}
	svc.incrementRetry("final error")

	if settings.data["retry_count"] != "10" {
		t.Errorf("expected retry_count=10, got %q", settings.data["retry_count"])
	}
	if settings.data["cert_needs_attention"] != "true" {
		t.Error("expected cert_needs_attention=true after max retries")
	}
}

func TestConstants(t *testing.T) {
	if MaxRetryCount != 10 {
		t.Errorf("expected MaxRetryCount=10, got %d", MaxRetryCount)
	}
	if RetryInterval != 7*time.Hour {
		t.Errorf("expected RetryInterval=7h, got %v", RetryInterval)
	}
	if CertCheckInterval != 1*time.Hour {
		t.Errorf("expected CertCheckInterval=1h, got %v", CertCheckInterval)
	}
}

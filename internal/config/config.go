package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	LogLevel    string
	// App mode (InSales external application). Both must be set to enable install/uninstall/autologin.
	InSalesAppKey       string
	InSalesSecretKey    string
	ReadOnly            bool
	JWTSecret           string
	EncryptionKey       string
	SMTPHost            string
	SMTPPort            string
	SMTPUser            string
	SMTPPass            string
	SMTPFrom            string
	Domain              string
	AllowedOrigins      []string
	MaxBodySize         int64
	JWTExpiry           time.Duration
	PDNAutoProtect      bool
	ServerCountry       string
	RegistrationEnabled bool
	BillingMonthly      int // RUB/month for InSales app; 0 = free
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:      getEnv("DB_URL", "postgres://admin:password@localhost:5432/mcp_is?sslmode=disable"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		InSalesAppKey:    getEnv("INSALES_APP_KEY", ""),
		InSalesSecretKey: getEnv("INSALES_SECRET_KEY", ""),
		ReadOnly:         getEnv("MCP_READ_ONLY", "true") == "true",
		JWTSecret:        getEnv("JWT_SECRET", "change-me-in-production"),
		EncryptionKey:    getEnv("ENCRYPTION_KEY", "change-me-32-bytes-key!!!"),
		SMTPHost:         getEnv("SMTP_HOST", ""),
		SMTPPort:         getEnv("SMTP_PORT", "587"),
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPass:         getEnv("SMTP_PASS", ""),
		SMTPFrom:         getEnv("SMTP_FROM", "MCP IS <noreply@localhost>"),
		Domain:           getEnv("DOMAIN", "localhost"),
		AllowedOrigins:   parseOrigins(getEnv("ALLOWED_ORIGINS", "*")),
		MaxBodySize:      1 << 20, // 1MB default
		ServerCountry:    getEnv("SERVER_COUNTRY", "RU"),
		PDNAutoProtect:   getEnv("PDN_AUTO_PROTECT", "true") == "true",
		JWTExpiry:        parseDuration(getEnv("JWT_EXPIRY_HOURS", "24"), time.Hour),
		// Default false: registration is admin-only unless explicitly enabled (see .env.example).
		RegistrationEnabled: getEnv("REGISTRATION_ENABLED", "false") == "true",
		BillingMonthly:      getEnvInt("BILLING_MONTHLY", 0),
	}

	// Security warnings
	if cfg.JWTSecret == "change-me-in-production" {
		slog.Warn("SECURITY: JWT_SECRET is using default value! Change it in .env for production!")
	}
	if cfg.EncryptionKey == "change-me-32-bytes-key!!!" {
		slog.Warn("SECURITY: ENCRYPTION_KEY is using default value! Change it in .env for production!")
	}
	if cfg.Domain == "localhost" {
		slog.Warn("SECURITY: DOMAIN is localhost — HTTPS will use self-signed certificate")
	}

	return cfg
}

// AppModeEnabled reports whether the InSales external-application mode is configured.
func (c *Config) AppModeEnabled() bool {
	return c.InSalesAppKey != "" && c.InSalesSecretKey != ""
}

func parseOrigins(origins string) []string {
	if origins == "*" {
		return []string{"*"}
	}
	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string, unit time.Duration) time.Duration {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 24 * time.Hour
	}
	return time.Duration(n) * unit
}

func getEnvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}

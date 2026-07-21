package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	ServerPort       string
	LogLevel         string
	InSalesShop      string
	InSalesAPIKey    string
	InSalesPassword  string
	ReadOnly         bool
	JWTSecret        string
	EncryptionKey    string
	SMTPHost         string
	SMTPPort         string
	SMTPUser         string
	SMTPPass         string
	SMTPFrom         string
	Domain           string
	AllowedOrigins   []string
	MaxBodySize      int64
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:    getEnv("DB_URL", "postgres://admin:password@localhost:5432/mcp_is?sslmode=disable"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		InSalesShop:    getEnv("INSALES_SHOP", ""),
		InSalesAPIKey:  getEnv("INSALES_API_KEY", ""),
		InSalesPassword: getEnv("INSALES_PASSWORD", ""),
		ReadOnly:       getEnv("MCP_READ_ONLY", "true") == "true",
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		EncryptionKey:  getEnv("ENCRYPTION_KEY", "change-me-32-bytes-key!!!"),
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASS", ""),
		SMTPFrom:       getEnv("SMTP_FROM", "MCP IS <noreply@localhost>"),
		Domain:         getEnv("DOMAIN", "localhost"),
		AllowedOrigins: parseOrigins(getEnv("ALLOWED_ORIGINS", "*")),
		MaxBodySize:    1 << 20, // 1MB default
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

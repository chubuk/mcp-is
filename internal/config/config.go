package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	ServerPort      string
	LogLevel        string
	InSalesShop     string
	InSalesAPIKey   string
	InSalesPassword string
	ReadOnly        bool
	JWTSecret       string
	EncryptionKey   string
	SMTPHost        string
	SMTPPort        string
	SMTPUser        string
	SMTPPass        string
	SMTPFrom        string
	Domain          string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		DatabaseURL:     getEnv("DB_URL", "postgres://admin:password@localhost:5432/mcp_insales?sslmode=disable"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		InSalesShop:     getEnv("INSALES_SHOP", ""),
		InSalesAPIKey:   getEnv("INSALES_API_KEY", ""),
		InSalesPassword: getEnv("INSALES_PASSWORD", ""),
		ReadOnly:        getEnv("MCP_READ_ONLY", "true") == "true",
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
		EncryptionKey:   getEnv("ENCRYPTION_KEY", "change-me-32-bytes-key!!!"),
		SMTPHost:        getEnv("SMTP_HOST", ""),
		SMTPPort:        getEnv("SMTP_PORT", "587"),
		SMTPUser:        getEnv("SMTP_USER", ""),
		SMTPPass:        getEnv("SMTP_PASS", ""),
		SMTPFrom:        getEnv("SMTP_FROM", "MCP InSales <noreply@localhost>"),
		Domain:          getEnv("DOMAIN", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

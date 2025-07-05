package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	JWTSecret      []byte
	DatabaseURL    string
	ValkeyURL      string
	Port           string
	AllowedOrigins []string
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPass       string
	SessionDomain  string
	SessionSecure  bool
}

var AppConfig *Config

func Init() {
	AppConfig = &Config{
		JWTSecret:      getRequiredEnvBytes("JWT_SECRET"),
		DatabaseURL:    getEnvOrDefault("COUCHDB_URL", "http://admin:password@localhost:5984/"),
		ValkeyURL:      getEnvOrDefault("VALKEY_URL", "valkey://valkeypassword@localhost:6379"),
		Port:           getEnvOrDefault("PORT", "8080"),
		AllowedOrigins: strings.Split(getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		SMTPHost:       getEnvOrDefault("SMTP_HOST", "localhost"),
		SMTPPort:       getEnvOrDefault("SMTP_PORT", "1025"),
		SMTPUser:       getEnvOrDefault("SMTP_USER", "hello@webenable.asia"),
		SMTPPass:       os.Getenv("SMTP_PASS"),
		SessionDomain:  getEnvOrDefault("SESSION_DOMAIN", ""),
		SessionSecure:  getEnvOrDefault("SESSION_SECURE", "false") == "true",
	}
}

func getRequiredEnvBytes(key string) []byte {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable is required", key)
	}
	return []byte(value)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

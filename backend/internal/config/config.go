package config

import (
	"os"
	"strings"
)

type Config struct {
	Port            string
	DBPath          string
	TurnstileSecret string
	CORSOrigins     []string
	ValidSources    []string
	ResendAPIKey    string
	ResendFrom      string
	SummaryEmail    string
	SummarySchedule string
}

func Load() *Config {
	return &Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		DBPath:          getEnvOrDefault("DB_PATH", "/data/dropmail.db"),
		TurnstileSecret: getEnvOrDefault("TURNSTILE_SECRET", "1x0000000000000000000000000000000AA"),
		CORSOrigins:     parseCSV(getEnvOrDefault("CORS_ORIGINS", "http://localhost:3000/")),
		ValidSources:    parseCSV(os.Getenv("VALID_SOURCES")),
		ResendAPIKey:    os.Getenv("RESEND_API_KEY"),
		ResendFrom:      os.Getenv("RESEND_FROM"),
		SummaryEmail:    os.Getenv("SUMMARY_EMAIL"),
		SummarySchedule: getEnvOrDefault("SUMMARY_SCHEDULE", "08:00"),
	}
}

// getEnvOrDefault returns the environment variable value or a default.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// MaskedTurnstileSecret returns a masked version of the secret for logging.
func (c *Config) MaskedTurnstileSecret() string {
	return maskSecret(c.TurnstileSecret)
}

// MaskedResendAPIKey returns a masked version of the Resend API key for logging.
func (c *Config) MaskedResendAPIKey() string {
	return maskSecret(c.ResendAPIKey)
}

func maskSecret(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***"
}

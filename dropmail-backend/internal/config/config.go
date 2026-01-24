package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port            string
	DBPath          string
	TurnstileSecret string
	CORSOrigins     []string
}

// Load loads configuration from environment variables.
// Returns an error if required environment variables are missing.
func Load() (*Config, error) {
	cfg := &Config{
		Port:   getEnvOrDefault("PORT", "8080"),
		DBPath: getEnvOrDefault("DB_PATH", "/data/dropmail.db"),
	}

	// Required: TURNSTILE_SECRET
	cfg.TurnstileSecret = os.Getenv("TURNSTILE_SECRET")
	if cfg.TurnstileSecret == "" {
		return nil, fmt.Errorf("TURNSTILE_SECRET is required\n  Get your secret key from: https://dash.cloudflare.com/turnstile\n  Example: TURNSTILE_SECRET=0x4AAAAAAxxxxxxxxxxxxxxxx")
	}

	// Required: CORS_ORIGINS
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		return nil, fmt.Errorf("CORS_ORIGINS is required\n  Provide comma-separated list of allowed origins\n  Example: CORS_ORIGINS=https://example.com,https://notion.so")
	}
	cfg.CORSOrigins = parseCORSOrigins(corsOrigins)

	return cfg, nil
}

// getEnvOrDefault returns the environment variable value or a default.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseCORSOrigins parses comma-separated CORS origins and trims whitespace.
func parseCORSOrigins(origins string) []string {
	parts := strings.Split(origins, ",")
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
	if len(c.TurnstileSecret) <= 8 {
		return "***"
	}
	return c.TurnstileSecret[:4] + "***"
}

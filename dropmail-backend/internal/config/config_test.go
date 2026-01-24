package config

import (
	"os"
	"testing"
)

func TestLoad_MissingTurnstileSecret(t *testing.T) {
	// Clear environment
	os.Unsetenv("TURNSTILE_SECRET")
	os.Unsetenv("CORS_ORIGINS")

	_, err := Load()
	if err == nil {
		t.Error("expected error for missing TURNSTILE_SECRET, got nil")
	}
}

func TestLoad_MissingCORSOrigins(t *testing.T) {
	// Set TURNSTILE_SECRET but not CORS_ORIGINS
	os.Setenv("TURNSTILE_SECRET", "test-secret")
	os.Unsetenv("CORS_ORIGINS")
	defer os.Unsetenv("TURNSTILE_SECRET")

	_, err := Load()
	if err == nil {
		t.Error("expected error for missing CORS_ORIGINS, got nil")
	}
}

func TestLoad_ValidConfiguration(t *testing.T) {
	// Set required environment variables
	os.Setenv("TURNSTILE_SECRET", "test-secret-key")
	os.Setenv("CORS_ORIGINS", "https://example.com,https://notion.so")
	defer func() {
		os.Unsetenv("TURNSTILE_SECRET")
		os.Unsetenv("CORS_ORIGINS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.TurnstileSecret != "test-secret-key" {
		t.Errorf("expected TurnstileSecret 'test-secret-key', got '%s'", cfg.TurnstileSecret)
	}

	if len(cfg.CORSOrigins) != 2 {
		t.Errorf("expected 2 CORS origins, got %d", len(cfg.CORSOrigins))
	}

	if cfg.CORSOrigins[0] != "https://example.com" {
		t.Errorf("expected first CORS origin 'https://example.com', got '%s'", cfg.CORSOrigins[0])
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	// Set required environment variables
	os.Setenv("TURNSTILE_SECRET", "test-secret")
	os.Setenv("CORS_ORIGINS", "https://example.com")
	os.Unsetenv("PORT")
	os.Unsetenv("DB_PATH")
	defer func() {
		os.Unsetenv("TURNSTILE_SECRET")
		os.Unsetenv("CORS_ORIGINS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default Port '8080', got '%s'", cfg.Port)
	}

	if cfg.DBPath != "/data/dropmail.db" {
		t.Errorf("expected default DBPath '/data/dropmail.db', got '%s'", cfg.DBPath)
	}
}

func TestParseCORSOrigins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single origin",
			input:    "https://example.com",
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple origins",
			input:    "https://example.com,https://notion.so",
			expected: []string{"https://example.com", "https://notion.so"},
		},
		{
			name:     "origins with whitespace",
			input:    "https://example.com , https://notion.so , https://test.com",
			expected: []string{"https://example.com", "https://notion.so", "https://test.com"},
		},
		{
			name:     "empty parts filtered",
			input:    "https://example.com,,https://notion.so",
			expected: []string{"https://example.com", "https://notion.so"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCORSOrigins(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d origins, got %d", len(tt.expected), len(result))
				return
			}
			for i, origin := range result {
				if origin != tt.expected[i] {
					t.Errorf("expected origin[%d] '%s', got '%s'", i, tt.expected[i], origin)
				}
			}
		})
	}
}

func TestMaskedTurnstileSecret(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		expected string
	}{
		{
			name:     "short secret",
			secret:   "abc",
			expected: "***",
		},
		{
			name:     "long secret",
			secret:   "0x4AAAAAAxxxxxxxxxxxxxxxx",
			expected: "0x4A***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{TurnstileSecret: tt.secret}
			result := cfg.MaskedTurnstileSecret()
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultTurnstileSecret(t *testing.T) {
	// Clear environment to test default
	os.Unsetenv("TURNSTILE_SECRET")
	defer os.Unsetenv("TURNSTILE_SECRET")

	cfg := Load()
	if cfg.TurnstileSecret != "1x0000000000000000000000000000000AA" {
		t.Errorf("expected default Turnstile secret, got '%s'", cfg.TurnstileSecret)
	}
}

func TestLoad_DefaultCORSOrigins(t *testing.T) {
	// Clear environment to test default
	os.Unsetenv("CORS_ORIGINS")
	defer os.Unsetenv("CORS_ORIGINS")

	cfg := Load()
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "http://localhost:3000/" {
		t.Errorf("expected default CORS origin 'http://localhost:3000/', got %v", cfg.CORSOrigins)
	}
}

func TestLoad_ValidConfiguration(t *testing.T) {
	// Set environment variables
	os.Setenv("TURNSTILE_SECRET", "test-secret-key")
	os.Setenv("CORS_ORIGINS", "https://example.com,https://notion.so")
	defer func() {
		os.Unsetenv("TURNSTILE_SECRET")
		os.Unsetenv("CORS_ORIGINS")
	}()

	cfg := Load()

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
	// Clear optional environment variables to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DB_PATH")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_PATH")
	}()

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default Port '8080', got '%s'", cfg.Port)
	}

	if cfg.DBPath != "/data/dropmail.db" {
		t.Errorf("expected default DBPath '/data/dropmail.db', got '%s'", cfg.DBPath)
	}
}

func TestParseCSV(t *testing.T) {
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
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCSV(tt.input)
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

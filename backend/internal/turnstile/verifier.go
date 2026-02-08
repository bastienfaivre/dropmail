package turnstile

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// VerifyURL is the Cloudflare Turnstile verification endpoint.
	VerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	// DefaultTimeout is the default timeout for verification requests.
	DefaultTimeout = 5 * time.Second
)

// VerifyResponse represents the response from Cloudflare's siteverify API.
type VerifyResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes,omitempty"`
	Hostname   string   `json:"hostname,omitempty"`
	Action     string   `json:"action,omitempty"`
	CData      string   `json:"cdata,omitempty"`
}

// Verifier handles Turnstile token verification.
type Verifier struct {
	secret     string
	httpClient *http.Client
	verifyURL  string // allow override for testing
}

// NewVerifier creates a new Turnstile verifier.
func NewVerifier(secret string) *Verifier {
	return &Verifier{
		secret: secret,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		verifyURL: VerifyURL,
	}
}

// VerifyToken verifies a Turnstile token with Cloudflare's siteverify API.
// Returns nil on success, or an error on failure.
func (v *Verifier) VerifyToken(ctx context.Context, token, remoteIP string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}

	data := url.Values{}
	data.Set("secret", v.secret)
	data.Set("response", token)
	if remoteIP != "" {
		data.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.verifyURL, strings.NewReader(data.Encode()))
	if err != nil {
		slog.Error("failed to create turnstile verification request", "error", err)
		return fmt.Errorf("verification service unavailable")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		slog.Error("turnstile verification request failed",
			"error", err,
			"remoteIP", remoteIP,
		)
		return fmt.Errorf("verification service unavailable")
	}
	defer resp.Body.Close()

	var verifyResp VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		slog.Error("failed to parse turnstile verification response",
			"error", err,
			"status", resp.StatusCode,
		)
		return fmt.Errorf("verification service unavailable")
	}

	if !verifyResp.Success {
		slog.Warn("turnstile verification failed",
			"error_codes", verifyResp.ErrorCodes,
			"remoteIP", remoteIP,
		)
		return fmt.Errorf("turnstile verification failed")
	}

	slog.Debug("turnstile verification successful",
		"hostname", verifyResp.Hostname,
		"remoteIP", remoteIP,
	)
	return nil
}

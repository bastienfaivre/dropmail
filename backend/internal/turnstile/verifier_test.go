package turnstile

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestVerifyToken_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		r.ParseForm()
		if r.Form.Get("secret") != "test-secret" {
			t.Errorf("expected secret 'test-secret', got '%s'", r.Form.Get("secret"))
		}
		if r.Form.Get("response") != "valid-token" {
			t.Errorf("expected response 'valid-token', got '%s'", r.Form.Get("response"))
		}
		if r.Form.Get("remoteip") != "192.168.1.1" {
			t.Errorf("expected remoteip '192.168.1.1', got '%s'", r.Form.Get("remoteip"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VerifyResponse{
			Success:  true,
			Hostname: "example.com",
		})
	}))
	defer server.Close()

	verifier := NewVerifier("test-secret")
	verifier.verifyURL = server.URL

	err := verifier.VerifyToken(context.Background(), "valid-token", "192.168.1.1")
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(VerifyResponse{
			Success:    false,
			ErrorCodes: []string{"invalid-input-response"},
		})
	}))
	defer server.Close()

	verifier := NewVerifier("test-secret")
	verifier.verifyURL = server.URL

	err := verifier.VerifyToken(context.Background(), "invalid-token", "192.168.1.1")
	if err == nil {
		t.Error("expected error for invalid token")
	}
	if !strings.Contains(err.Error(), "turnstile verification failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestVerifyToken_EmptyToken(t *testing.T) {
	verifier := NewVerifier("test-secret")

	err := verifier.VerifyToken(context.Background(), "", "192.168.1.1")
	if err == nil {
		t.Error("expected error for empty token")
	}
	if !strings.Contains(err.Error(), "token is empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestVerifyToken_ServiceUnavailable(t *testing.T) {
	// Create a server that always returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	verifier := NewVerifier("test-secret")
	verifier.verifyURL = server.URL

	err := verifier.VerifyToken(context.Background(), "valid-token", "192.168.1.1")
	if err == nil {
		t.Error("expected error for service unavailable")
	}
	if !strings.Contains(err.Error(), "verification service unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestVerifyToken_Timeout(t *testing.T) {
	// Create a server that delays longer than timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		json.NewEncoder(w).Encode(VerifyResponse{Success: true})
	}))
	defer server.Close()

	verifier := NewVerifier("test-secret")
	verifier.verifyURL = server.URL
	verifier.httpClient.Timeout = 10 * time.Millisecond // Very short timeout

	err := verifier.VerifyToken(context.Background(), "valid-token", "192.168.1.1")
	if err == nil {
		t.Error("expected error for timeout")
	}
	if !strings.Contains(err.Error(), "verification service unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestVerifyToken_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	verifier := NewVerifier("test-secret")
	verifier.verifyURL = server.URL

	err := verifier.VerifyToken(context.Background(), "valid-token", "192.168.1.1")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "verification service unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

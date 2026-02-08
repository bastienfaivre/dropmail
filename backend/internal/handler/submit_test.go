package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"backend/internal/database"
	"backend/internal/model"
)

// mockVerifier is a mock TurnstileVerifier for testing.
type mockVerifier struct {
	shouldFail bool
	failErr    error
}

func (m *mockVerifier) VerifyToken(ctx context.Context, token, remoteIP string) error {
	if m.shouldFail {
		if m.failErr != nil {
			return m.failErr
		}
		return fmt.Errorf("turnstile verification failed")
	}
	return nil
}

var testValidSources = []string{"test", "first-source", "second-source"}

func setupTestDB(t *testing.T) (*database.DB, func()) {
	t.Helper()

	// Create temp directory for test database
	tempDir, err := os.MkdirTemp("", "dropmail-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to open database: %v", err)
	}

	// Create submissions table for testing
	_, err = db.Conn().Exec(`
		CREATE TABLE IF NOT EXISTS submissions (
			email TEXT PRIMARY KEY NOT NULL,
			source TEXT NOT NULL DEFAULT 'unknown',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		db.Close()
		os.RemoveAll(tempDir)
		t.Fatalf("failed to create table: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.RemoveAll(tempDir)
	}

	return db, cleanup
}

func TestSubmitHandler_ValidSubmission(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `{"email": "test@example.com", "source": "test", "turnstile_token": "valid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp model.SubmitResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if resp.Message != "Email submitted successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}

	// Verify data was stored
	var count int
	err := db.Conn().QueryRowContext(context.Background(), "SELECT COUNT(*) FROM submissions").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 submission, got %d", count)
	}
}

func TestSubmitHandler_InvalidEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `{"email": "not-an-email", "source": "test", "turnstile_token": "valid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "Invalid email format" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}

	// Verify no data was stored
	var count int
	db.Conn().QueryRowContext(context.Background(), "SELECT COUNT(*) FROM submissions").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 submissions, got %d", count)
	}
}

func TestSubmitHandler_MissingEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `{"source": "test", "turnstile_token": "valid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "Email is required" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}
}

func TestSubmitHandler_MissingTurnstileToken(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `{"email": "test@example.com", "source": "test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "Turnstile token required" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}
}

func TestSubmitHandler_InvalidJSON(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `not valid json`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "Invalid JSON body" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}
}

func TestSubmitHandler_MethodNotAllowed(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	req := httptest.NewRequest(http.MethodGet, "/api/submit", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestSubmitHandler_MissingSource(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `{"email": "test@example.com", "turnstile_token": "valid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "Source is required" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}
}

func TestSubmitHandler_TurnstileVerificationFailed(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: true}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `{"email": "test@example.com", "source": "test", "turnstile_token": "invalid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "Turnstile verification failed" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}

	// Verify no data was stored
	var count int
	db.Conn().QueryRowContext(context.Background(), "SELECT COUNT(*) FROM submissions").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 submissions, got %d", count)
	}
}

func TestSubmitHandler_TurnstileServiceUnavailable(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: true, failErr: fmt.Errorf("verification service unavailable")}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	body := `{"email": "test@example.com", "source": "test", "turnstile_token": "valid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}

	var resp model.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "Verification service unavailable" {
		t.Errorf("unexpected error message: %s", resp.Error)
	}

	// Verify no data was stored
	var count int
	db.Conn().QueryRowContext(context.Background(), "SELECT COUNT(*) FROM submissions").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 submissions, got %d", count)
	}
}

func TestSubmitHandler_DuplicateEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(db, verifier, testValidSources)

	// First submission
	body1 := `{"email": "duplicate@example.com", "source": "first-source", "turnstile_token": "valid-token"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body1))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()

	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("first submission: expected status %d, got %d", http.StatusOK, rec1.Code)
	}

	// Get original source and created_at
	var originalSource string
	var originalCreatedAt string
	err := db.Conn().QueryRowContext(context.Background(),
		"SELECT source, created_at FROM submissions WHERE email = ?",
		"duplicate@example.com").Scan(&originalSource, &originalCreatedAt)
	if err != nil {
		t.Fatalf("failed to query original submission: %v", err)
	}

	if originalSource != "first-source" {
		t.Errorf("expected source 'first-source', got '%s'", originalSource)
	}

	// Second submission with same email but different source
	body2 := `{"email": "duplicate@example.com", "source": "second-source", "turnstile_token": "valid-token"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	// Should still return success
	if rec2.Code != http.StatusOK {
		t.Errorf("second submission: expected status %d, got %d", http.StatusOK, rec2.Code)
	}

	var resp model.SubmitResponse
	if err := json.NewDecoder(rec2.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true for duplicate submission")
	}

	// Verify only one record exists
	var count int
	err = db.Conn().QueryRowContext(context.Background(), "SELECT COUNT(*) FROM submissions WHERE email = ?", "duplicate@example.com").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 submission, got %d", count)
	}

	// Verify original data was preserved (source and created_at not overwritten)
	var currentSource string
	var currentCreatedAt string
	err = db.Conn().QueryRowContext(context.Background(),
		"SELECT source, created_at FROM submissions WHERE email = ?",
		"duplicate@example.com").Scan(&currentSource, &currentCreatedAt)
	if err != nil {
		t.Fatalf("failed to query current submission: %v", err)
	}

	if currentSource != originalSource {
		t.Errorf("source was overwritten: expected '%s', got '%s'", originalSource, currentSource)
	}

	if currentCreatedAt != originalCreatedAt {
		t.Errorf("created_at was overwritten: expected '%s', got '%s'", originalCreatedAt, currentCreatedAt)
	}
}

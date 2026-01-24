package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bastien/dropmail-backend/internal/database"
	"github.com/bastien/dropmail-backend/internal/model"
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
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

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(repo, verifier)

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

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(repo, verifier)

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

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(repo, verifier)

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

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(repo, verifier)

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

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(repo, verifier)

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

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(repo, verifier)

	req := httptest.NewRequest(http.MethodGet, "/api/submit", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestSubmitHandler_DefaultSource(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: false}
	handler := NewSubmitHandler(repo, verifier)

	// Submit without source
	body := `{"email": "test@example.com", "turnstile_token": "valid-token"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify source was set to "direct"
	var source string
	err := db.Conn().QueryRowContext(context.Background(), "SELECT source FROM submissions WHERE email = ?", "test@example.com").Scan(&source)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Fatal("submission not found in database")
		}
		t.Fatalf("failed to query database: %v", err)
	}

	if source != "direct" {
		t.Errorf("expected source 'direct', got '%s'", source)
	}
}

func TestSubmitHandler_TurnstileVerificationFailed(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: true}
	handler := NewSubmitHandler(repo, verifier)

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

	repo := database.NewSubmissionRepository(db)
	verifier := &mockVerifier{shouldFail: true, failErr: fmt.Errorf("verification service unavailable")}
	handler := NewSubmitHandler(repo, verifier)

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

package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/bastien/dropmail-backend/internal/database"
	"github.com/bastien/dropmail-backend/internal/model"
)

// TurnstileVerifier defines the interface for Turnstile token verification.
type TurnstileVerifier interface {
	VerifyToken(ctx context.Context, token, remoteIP string) error
}

// SubmitHandler handles email submission requests.
type SubmitHandler struct {
	repo     *database.SubmissionRepository
	verifier TurnstileVerifier
}

// NewSubmitHandler creates a new submit handler.
func NewSubmitHandler(repo *database.SubmissionRepository, verifier TurnstileVerifier) *SubmitHandler {
	return &SubmitHandler{repo: repo, verifier: verifier}
}

// ServeHTTP implements http.Handler.
func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req model.SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Trim whitespace from email
	req.Email = strings.TrimSpace(req.Email)

	// Validate required fields
	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if req.TurnstileToken == "" {
		writeError(w, http.StatusBadRequest, "Turnstile token required")
		return
	}

	// Validate email format
	if !model.ValidateEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Default source to "direct" if not provided
	source := req.Source
	if source == "" {
		source = "direct"
	}

	// Verify Turnstile token
	clientIP := getClientIP(r)
	if err := h.verifier.VerifyToken(r.Context(), req.TurnstileToken, clientIP); err != nil {
		errMsg := err.Error()
		// Map internal error messages to user-friendly responses
		if strings.Contains(errMsg, "verification service unavailable") {
			writeError(w, http.StatusServiceUnavailable, "Verification service unavailable")
		} else {
			writeError(w, http.StatusBadRequest, "Turnstile verification failed")
		}
		return
	}

	// Store submission
	submission, err := h.repo.Insert(r.Context(), req.Email, source)
	if err != nil {
		slog.Error("failed to store submission",
			"error", err,
			"email", req.Email,
			"source", source,
		)
		writeError(w, http.StatusInternalServerError, "Failed to store submission")
		return
	}

	// Log successful submission
	slog.Info("submission received",
		"email", req.Email,
		"source", source,
		"ip", r.RemoteAddr,
		"id", submission.ID,
		"created_at", submission.CreatedAt,
	)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SubmitResponse{
		Success: true,
		Message: "Email submitted successfully",
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{Error: message})
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxied requests)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs; the first is the client
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

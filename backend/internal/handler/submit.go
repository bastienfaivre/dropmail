package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/model"
)

// TurnstileVerifier defines the interface for Turnstile token verification.
type TurnstileVerifier interface {
	VerifyToken(ctx context.Context, token, remoteIP string) error
}

type SubmitHandler struct {
	db           *database.DB
	verifier     TurnstileVerifier
	validSources map[string]struct{}
}

func NewSubmitHandler(db *database.DB, verifier TurnstileVerifier, validSources []string) *SubmitHandler {
	sourceMap := make(map[string]struct{}, len(validSources))
	for _, s := range validSources {
		sourceMap[s] = struct{}{}
	}
	return &SubmitHandler{db: db, verifier: verifier, validSources: sourceMap}
}

const maxBodySize = 2048

func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req model.SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if req.TurnstileToken == "" {
		writeError(w, http.StatusBadRequest, "Turnstile token required")
		return
	}

	if len(req.Email) > 254 {
		writeError(w, http.StatusBadRequest, "Email too long")
		return
	}

	if !model.ValidateEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	source := strings.TrimSpace(req.Source)
	if source == "" {
		writeError(w, http.StatusBadRequest, "Source is required")
		return
	}
	if len(source) > 32 {
		writeError(w, http.StatusBadRequest, "Source too long")
		return
	}

	if _, ok := h.validSources[source]; !ok {
		writeError(w, http.StatusBadRequest, "Invalid source")
		return
	}

	clientIP := getClientIP(r)
	if err := h.verifier.VerifyToken(r.Context(), req.TurnstileToken, clientIP); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "verification service unavailable") {
			writeError(w, http.StatusServiceUnavailable, "Verification service unavailable")
		} else {
			writeError(w, http.StatusBadRequest, "Turnstile verification failed")
		}
		return
	}

	reqID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	submission, err := h.db.InsertSubmission(r.Context(), req.Email, source)
	if err != nil {
		slog.Error("failed to store submission",
			"request_id", reqID,
			"error", err,
			"email", req.Email,
			"source", source,
		)
		writeError(w, http.StatusInternalServerError, "Failed to store submission")
		return
	}

	slog.Info("submission received",
		"request_id", reqID,
		"email", req.Email,
		"source", source,
		"ip", r.RemoteAddr,
		"created_at", submission.CreatedAt,
	)

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
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

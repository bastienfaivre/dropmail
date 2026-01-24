package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the health check response body.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler handles health check requests.
type HealthHandler struct{}

// NewHealthHandler creates a new health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP implements http.Handler for health checks.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{Status: "healthy"}
	json.NewEncoder(w).Encode(response)
}

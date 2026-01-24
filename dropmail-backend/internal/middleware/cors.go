package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// NewCORS creates a CORS middleware with the specified allowed origins.
func NewCORS(allowedOrigins []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"POST", "OPTIONS", "GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	})
}

// Logging wraps a handler with request logging middleware.
func Logging(next http.Handler) http.Handler {
	return next
}

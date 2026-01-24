package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bastien/dropmail-backend/internal/config"
	"github.com/bastien/dropmail-backend/internal/database"
	"github.com/bastien/dropmail-backend/internal/handler"
	"github.com/bastien/dropmail-backend/internal/middleware"
	"github.com/bastien/dropmail-backend/internal/turnstile"
)

func main() {
	// Initialize structured JSON logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration - fail fast if invalid
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Configuration validation failed\n\n%s\n", err)
		os.Exit(1)
	}

	// Log configuration (without exposing secrets)
	slog.Info("configuration loaded",
		"port", cfg.Port,
		"db_path", cfg.DBPath,
		"cors_origins", cfg.CORSOrigins,
		"turnstile_secret", cfg.MaskedTurnstileSecret(),
	)

	// Initialize database
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	// Look for migrations in these locations (in order):
	// 1. /migrations (Docker container)
	// 2. ./migrations (local development)
	// 3. Next to executable
	var migrationsPath string
	for _, path := range []string{"/migrations", "migrations"} {
		if _, err := os.Stat(path); err == nil {
			migrationsPath = path
			break
		}
	}
	if migrationsPath == "" {
		execPath, err := os.Executable()
		if err != nil {
			slog.Error("failed to get executable path", "error", err)
			os.Exit(1)
		}
		migrationsPath = filepath.Join(filepath.Dir(execPath), "migrations")
	}

	// Convert to file:// URL for golang-migrate
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		slog.Error("failed to get absolute path for migrations", "error", err)
		os.Exit(1)
	}
	migrationsURL := "file://" + absPath

	slog.Info("running database migrations", "path", migrationsURL)
	if err := db.RunMigrations(migrationsURL); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Create submission repository
	submissionRepo := database.NewSubmissionRepository(db)

	// Create Turnstile verifier
	turnstileVerifier := turnstile.NewVerifier(cfg.TurnstileSecret)

	// Create router
	mux := http.NewServeMux()

	// Register handlers
	healthHandler := handler.NewHealthHandler()
	mux.Handle("/health", healthHandler)

	submitHandler := handler.NewSubmitHandler(submissionRepo, turnstileVerifier)
	mux.Handle("/api/submit", submitHandler)

	// Create middleware chain: CORS → Rate Limit → Handler
	corsMiddleware := middleware.NewCORS(cfg.CORSOrigins)
	rateLimiter := middleware.NewRateLimiter(10, 5) // 10 requests per minute, burst of 5

	// Apply middleware chain
	var finalHandler http.Handler = mux
	finalHandler = rateLimiter.Middleware(finalHandler)
	finalHandler = corsMiddleware.Handler(finalHandler)

	slog.Info("middleware configured",
		"cors_origins", cfg.CORSOrigins,
		"rate_limit", "10/min",
		"rate_burst", 5,
	)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors from the server
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		slog.Info("starting HTTP server", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}

	case sig := <-shutdown:
		slog.Info("shutdown signal received", "signal", sig)

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			server.Close()
		}
	}

	slog.Info("server stopped")
}

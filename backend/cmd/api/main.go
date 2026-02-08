package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/email"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/scheduler"
	"backend/internal/turnstile"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	slog.Info("configuration loaded",
		"port", cfg.Port,
		"db_path", cfg.DBPath,
		"cors_origins", cfg.CORSOrigins,
		"valid_sources", cfg.ValidSources,
		"turnstile_secret", cfg.MaskedTurnstileSecret(),
		"resend_api_key", cfg.MaskedResendAPIKey(),
		"resend_from", cfg.ResendFrom,
		"summary_email", cfg.SummaryEmail,
		"summary_schedule", cfg.SummarySchedule,
	)

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	// Look for migrations in these locations (in order):
	// 1. /migrations (production)
	// 2. ./migrations (development)
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

	turnstileVerifier := turnstile.NewVerifier(cfg.TurnstileSecret)

	mux := http.NewServeMux()

	submitHandler := handler.NewSubmitHandler(db, turnstileVerifier, cfg.ValidSources)
	mux.Handle("/api/submit", submitHandler)

	corsMiddleware := middleware.NewCORS(cfg.CORSOrigins)
	finalHandler := middleware.RequestID(middleware.SecurityHeaders(corsMiddleware.Handler(mux)))

	// Create server
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           finalHandler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start daily summary scheduler if Resend is configured
	var cancelScheduler context.CancelFunc
	if cfg.ResendAPIKey != "" && cfg.SummaryEmail != "" && cfg.ResendFrom != "" {
		mailer := email.NewClient(cfg.ResendAPIKey, cfg.ResendFrom, cfg.SummaryEmail)
		ds, err := scheduler.NewDailySummary(db, mailer, cfg.SummarySchedule)
		if err != nil {
			slog.Error("failed to create daily summary scheduler", "error", err)
			os.Exit(1)
		}
		var schedulerCtx context.Context
		schedulerCtx, cancelScheduler = context.WithCancel(context.Background())
		go ds.Start(schedulerCtx)
		slog.Info("daily summary scheduler started")
	} else {
		slog.Info("daily summary scheduler disabled (RESEND_API_KEY, RESEND_FROM, or SUMMARY_EMAIL not set)")
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

		if cancelScheduler != nil {
			cancelScheduler()
		}

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

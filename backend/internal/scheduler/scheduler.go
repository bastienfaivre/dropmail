package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"backend/internal/database"
	"backend/internal/email"
)

// DailySummary runs a daily job that queries recent submissions and emails a summary.
type DailySummary struct {
	db     *database.DB
	mailer *email.Client
	hour   int
	minute int
}

// NewDailySummary creates a new daily summary scheduler.
// schedule must be in "HH:MM" format (24h, UTC).
func NewDailySummary(db *database.DB, mailer *email.Client, schedule string) (*DailySummary, error) {
	var hour, minute int
	if _, err := fmt.Sscanf(schedule, "%d:%d", &hour, &minute); err != nil {
		return nil, fmt.Errorf("invalid schedule format %q (expected HH:MM): %w", schedule, err)
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return nil, fmt.Errorf("invalid schedule time %q: hour must be 0-23, minute must be 0-59", schedule)
	}

	return &DailySummary{
		db:     db,
		mailer: mailer,
		hour:   hour,
		minute: minute,
	}, nil
}

// Start launches the scheduler in a background goroutine.
// It blocks until ctx is cancelled.
func (s *DailySummary) Start(ctx context.Context) {
	for {
		next := s.nextRun(time.Now().UTC())
		delay := time.Until(next)

		slog.Info("daily summary scheduled", "next_run", next.Format(time.RFC3339), "delay", delay)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			slog.Info("daily summary scheduler stopped")
			return
		case <-timer.C:
			s.run(ctx)
		}
	}
}

func (s *DailySummary) nextRun(now time.Time) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, 0, 0, time.UTC)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func (s *DailySummary) run(ctx context.Context) {
	slog.Info("running daily summary job")

	since := time.Now().UTC().Add(-24 * time.Hour)
	submissions, err := s.db.GetRecentSubmissions(ctx, since)
	if err != nil {
		slog.Error("daily summary: failed to query submissions", "error", err)
		return
	}

	if len(submissions) == 0 {
		slog.Info("daily summary: no submissions to report, skipping email")
		return
	}

	if err := s.mailer.SendDailySummary(submissions); err != nil {
		slog.Error("daily summary: failed to send email", "error", err)
		return
	}

	slog.Info("daily summary sent", "count", len(submissions))
}

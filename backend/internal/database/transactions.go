package database

import (
	"context"
	"fmt"
	"time"

	"backend/internal/model"
)

// InsertSubmission creates a new submission in the database.
// Uses parameterized queries to prevent SQL injection.
// If email already exists, returns success without overwriting existing data.
func (db *DB) InsertSubmission(ctx context.Context, email, source string) (*model.Submission, error) {
	query := `
		INSERT INTO submissions (email, source) VALUES (?, ?)
		ON CONFLICT(email) DO UPDATE SET email = email
		RETURNING created_at`

	var submission model.Submission
	submission.Email = email
	submission.Source = source
	err := db.conn.QueryRowContext(ctx, query, email, source).Scan(&submission.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert submission: %w", err)
	}

	return &submission, nil
}

// GetRecentSubmissions returns all submissions created since the given time.
func (db *DB) GetRecentSubmissions(ctx context.Context, since time.Time) ([]model.Submission, error) {
	query := `SELECT email, source, created_at FROM submissions WHERE created_at >= ? ORDER BY created_at DESC`

	rows, err := db.conn.QueryContext(ctx, query, since.UTC().Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, fmt.Errorf("failed to query recent submissions: %w", err)
	}
	defer rows.Close()

	var submissions []model.Submission
	for rows.Next() {
		var s model.Submission
		if err := rows.Scan(&s.Email, &s.Source, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating submissions: %w", err)
	}

	return submissions, nil
}

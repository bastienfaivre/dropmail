package database

import (
	"context"
	"fmt"

	"github.com/bastien/dropmail-backend/internal/model"
)

// SubmissionRepository handles database operations for submissions.
type SubmissionRepository struct {
	db *DB
}

// NewSubmissionRepository creates a new submission repository.
func NewSubmissionRepository(db *DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

// Insert creates a new submission in the database.
// Uses parameterized queries to prevent SQL injection.
func (r *SubmissionRepository) Insert(ctx context.Context, email, source string) (*model.Submission, error) {
	query := `INSERT INTO submissions (email, source) VALUES (?, ?)`

	result, err := r.db.conn.ExecContext(ctx, query, email, source)
	if err != nil {
		return nil, fmt.Errorf("failed to insert submission: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Retrieve the created submission to get the timestamp
	var submission model.Submission
	selectQuery := `SELECT id, email, source, created_at FROM submissions WHERE id = ?`
	err = r.db.conn.QueryRowContext(ctx, selectQuery, id).Scan(
		&submission.ID,
		&submission.Email,
		&submission.Source,
		&submission.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created submission: %w", err)
	}

	return &submission, nil
}

package model

import (
	"net/mail"
	"time"
)

// Submission represents an email subscription record.
type Submission struct {
	Email     string    `json:"email"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// SubmitRequest represents the incoming submission request.
type SubmitRequest struct {
	Email          string `json:"email"`
	Source         string `json:"source"`
	TurnstileToken string `json:"turnstile_token"`
}

// SubmitResponse represents a successful submission response.
type SubmitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidateEmail checks if the email format is valid.
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

package digikey

import (
	"errors"
	"fmt"
)

var (
	// ErrUnauthorized indicates invalid or expired credentials.
	ErrUnauthorized = errors.New("digikey: unauthorized")

	// ErrForbidden indicates the request was forbidden.
	ErrForbidden = errors.New("digikey: forbidden")

	// ErrNotFound indicates the requested resource was not found.
	ErrNotFound = errors.New("digikey: not found")

	// ErrRateLimitExceeded indicates the rate limit has been exceeded.
	ErrRateLimitExceeded = errors.New("digikey: rate limit exceeded")

	// ErrInvalidRequest indicates an invalid request.
	ErrInvalidRequest = errors.New("digikey: invalid request")

	// ErrServerError indicates a server-side error.
	ErrServerError = errors.New("digikey: server error")
)

// APIError represents an error returned by the Digi-Key API.
type APIError struct {
	StatusCode int
	Message    string
	Details    string
	RequestID  string
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("digikey: %s (status %d): %s", e.Message, e.StatusCode, e.Details)
	}
	return fmt.Sprintf("digikey: %s (status %d)", e.Message, e.StatusCode)
}

// Unwrap returns the underlying error type for errors.Is matching.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 429:
		return ErrRateLimitExceeded
	case 400:
		return ErrInvalidRequest
	default:
		if e.StatusCode >= 500 {
			return ErrServerError
		}
		return nil
	}
}

// AuthError represents an authentication error.
type AuthError struct {
	Err         string `json:"error"`
	Description string `json:"error_description"`
}

func (e *AuthError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("digikey: auth error: %s: %s", e.Err, e.Description)
	}
	return fmt.Sprintf("digikey: auth error: %s", e.Err)
}

// Unwrap returns the underlying unauthorized error.
func (e *AuthError) Unwrap() error {
	return ErrUnauthorized
}

// RateLimitError provides details about rate limit violations.
type RateLimitError struct {
	Limit     int
	Remaining int
	ResetAt   string
	Type      string // "minute" or "day"
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("digikey: rate limit exceeded (%s limit: %d, remaining: %d, resets at: %s)",
		e.Type, e.Limit, e.Remaining, e.ResetAt)
}

// Unwrap returns the underlying rate limit error.
func (e *RateLimitError) Unwrap() error {
	return ErrRateLimitExceeded
}

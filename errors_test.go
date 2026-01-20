package digikey

import (
	"errors"
	"testing"
)

// TestAPIErrorError tests the Error method of APIError.
func TestAPIErrorError(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Message:    "bad request",
	}

	expected := "digikey: bad request (status 400)"
	if err.Error() != expected {
		t.Errorf("expected error string %q, got %q", expected, err.Error())
	}
}

// TestAPIErrorErrorWithDetails tests Error method with details.
func TestAPIErrorErrorWithDetails(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Message:    "bad request",
		Details:    "missing required field",
	}

	errStr := err.Error()
	if errStr == "" {
		t.Fatal("expected non-empty error string")
	}

	if !contains(errStr, "bad request") {
		t.Errorf("expected error to contain message")
	}
	if !contains(errStr, "400") {
		t.Errorf("expected error to contain status code")
	}
}

// TestAPIErrorUnwrap tests the Unwrap method.
func TestAPIErrorUnwrap(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   error
	}{
		{401, ErrUnauthorized},
		{403, ErrForbidden},
		{404, ErrNotFound},
		{429, ErrRateLimitExceeded},
		{400, ErrInvalidRequest},
		{500, ErrServerError},
	}

	for _, tc := range testCases {
		err := &APIError{
			StatusCode: tc.statusCode,
			Message:    "test",
		}

		unwrapped := err.Unwrap()
		if !errors.Is(err, tc.expected) {
			t.Errorf("expected %v, got %v for status %d", tc.expected, unwrapped, tc.statusCode)
		}
	}
}

// TestAPIErrorUnwrapUnknownStatus tests Unwrap with unknown status codes.
func TestAPIErrorUnwrapUnknownStatus(t *testing.T) {
	err := &APIError{
		StatusCode: 418,
		Message:    "test",
	}

	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Errorf("expected nil for unknown status code, got %v", unwrapped)
	}
}

// TestAuthError tests AuthError type.
func TestAuthError(t *testing.T) {
	err := &AuthError{
		Err:         "invalid_client",
		Description: "Client credentials not found",
	}

	errStr := err.Error()
	if !contains(errStr, "invalid_client") {
		t.Errorf("expected auth error to contain error code")
	}

	if !contains(errStr, "Client credentials") {
		t.Errorf("expected auth error to contain description")
	}
}

// TestAuthErrorWithoutDescription tests AuthError without description.
func TestAuthErrorWithoutDescription(t *testing.T) {
	err := &AuthError{
		Err: "invalid_client",
	}

	errStr := err.Error()
	if !contains(errStr, "invalid_client") {
		t.Errorf("expected auth error to contain error code")
	}
}

// TestAuthErrorUnwrap tests AuthError Unwrap.
func TestAuthErrorUnwrap(t *testing.T) {
	err := &AuthError{
		Err: "invalid_client",
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Error("expected auth error to unwrap to ErrUnauthorized")
	}
}

// TestRateLimitError tests RateLimitError type.
func TestRateLimitError(t *testing.T) {
	err := &RateLimitError{
		Limit:     1000,
		Remaining: 50,
		ResetAt:   "2024-01-20T12:30:00Z",
		Type:      "minute",
	}

	errStr := err.Error()
	if !contains(errStr, "rate limit exceeded") {
		t.Errorf("expected rate limit error message")
	}

	if !contains(errStr, "1000") {
		t.Errorf("expected rate limit error to contain limit")
	}

	if !contains(errStr, "50") {
		t.Errorf("expected rate limit error to contain remaining")
	}
}

// TestRateLimitErrorUnwrap tests RateLimitError Unwrap.
func TestRateLimitErrorUnwrap(t *testing.T) {
	err := &RateLimitError{
		Limit: 1000,
		Type:  "minute",
	}

	if !errors.Is(err, ErrRateLimitExceeded) {
		t.Error("expected rate limit error to unwrap to ErrRateLimitExceeded")
	}
}

// TestErrorVariables tests that error variables are distinct.
func TestErrorVariables(t *testing.T) {
	errors := []error{
		ErrUnauthorized,
		ErrForbidden,
		ErrNotFound,
		ErrRateLimitExceeded,
		ErrInvalidRequest,
		ErrServerError,
	}

	for i, err1 := range errors {
		for j, err2 := range errors {
			if i != j && err1 == err2 {
				t.Errorf("error %d should not equal error %d", i, j)
			}
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substring string) bool {
	for i := 0; i <= len(s)-len(substring); i++ {
		if s[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

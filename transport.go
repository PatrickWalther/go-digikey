package digikey

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// do performs an HTTP request with authentication, rate limiting, and retries.
func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	return c.doWithRetry(ctx, method, path, body, result, false)
}

// doWithRetry performs an HTTP request with retry logic.
func (c *Client) doWithRetry(ctx context.Context, method, path string, body interface{}, result interface{}, isRetryAfter401 bool) error {
	var lastErr error
	maxAttempts := c.retryConfig.MaxRetries + 1

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			backoff := c.retryConfig.calculateBackoff(attempt - 1)
			if err := sleep(ctx, backoff); err != nil {
				return err
			}
		}

		statusCode, shouldRetryRequest, err := c.doOnce(ctx, method, path, body, result)
		if err == nil {
			return nil
		}

		lastErr = err

		// Handle 401: refresh token and retry once
		if statusCode == http.StatusUnauthorized && !isRetryAfter401 {
			c.tokenManager.invalidate()
			return c.doWithRetry(ctx, method, path, body, result, true)
		}

		// Don't retry if not retryable
		if !shouldRetryRequest {
			return err
		}

		// Don't retry on last attempt
		if attempt >= maxAttempts-1 {
			return err
		}
	}

	return lastErr
}

// doOnce performs a single HTTP request attempt.
// Returns (statusCode, shouldRetry, error).
func (c *Client) doOnce(ctx context.Context, method, path string, body interface{}, result interface{}) (int, bool, error) {
	if err := c.rateLimiter.Allow(); err != nil {
		return 0, false, err
	}

	token, err := c.tokenManager.getToken(ctx)
	if err != nil {
		return 0, shouldRetry(err, 0), err
	}

	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return 0, false, fmt.Errorf("digikey: failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(bodyBytes))
	if err != nil {
		return 0, false, fmt.Errorf("digikey: failed to create request: %w", err)
	}

	locale := c.getLocale()
	c.setHeaders(req, token, locale)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, shouldRetry(err, 0), fmt.Errorf("digikey: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, false, fmt.Errorf("digikey: failed to read response: %w", err)
	}

	// Sync rate limiter from response headers on every response.
	c.rateLimiter.UpdateFromHeaders(resp.Header)

	// Handle rate limiting (429)
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		apiErr := c.handleErrorResponse(resp.StatusCode, respBody, resp.Header)
		rateErr := c.buildRateLimitError(resp.Header, retryAfter, apiErr)
		c.rateLimiter.UpdateFromRateLimitError(rateErr)
		return resp.StatusCode, shouldRetryRateLimit(rateErr), rateErr
	}

	// Handle unauthorized (401)
	if resp.StatusCode == http.StatusUnauthorized {
		return resp.StatusCode, false, &APIError{ // Don't retry here; handled in doWithRetry
			StatusCode: resp.StatusCode,
			Message:    "unauthorized",
			Details:    string(respBody),
			RequestID:  resp.Header.Get("X-Request-Id"),
		}
	}

	// Handle other errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := c.handleErrorResponse(resp.StatusCode, respBody, resp.Header)
		return resp.StatusCode, shouldRetry(nil, resp.StatusCode), apiErr
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return resp.StatusCode, false, fmt.Errorf("digikey: failed to parse response: %w", err)
		}
	}

	return resp.StatusCode, false, nil
}

// setHeaders sets the required headers for Digi-Key API requests.
func (c *Client) setHeaders(req *http.Request, token string, locale Locale) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-DIGIKEY-Client-Id", c.clientID)
	req.Header.Set("X-DIGIKEY-Locale-Site", locale.Site)
	req.Header.Set("X-DIGIKEY-Locale-Language", locale.Language)
	req.Header.Set("X-DIGIKEY-Locale-Currency", locale.Currency)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.customerID != "" {
		req.Header.Set("X-DIGIKEY-Customer-Id", c.customerID)
	}
}

// handleErrorResponse parses error responses from the API.
func (c *Client) handleErrorResponse(statusCode int, body []byte, headers http.Header) error {
	apiErr := &APIError{
		StatusCode: statusCode,
		RequestID:  headers.Get("X-Request-Id"),
	}

	var errResp struct {
		Message      string `json:"message"`
		Details      string `json:"details"`
		ErrorMessage string `json:"ErrorMessage"`
		ErrorDetails string `json:"ErrorDetails"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil {
		apiErr.Message = errResp.Message
		apiErr.Details = errResp.Details
		if apiErr.Message == "" {
			apiErr.Message = errResp.ErrorMessage
		}
		if apiErr.Details == "" {
			apiErr.Details = errResp.ErrorDetails
		}
	} else {
		apiErr.Message = http.StatusText(statusCode)
		apiErr.Details = string(body)
	}

	return apiErr
}

func shouldRetryRateLimit(rateErr *RateLimitError) bool {
	return !strings.EqualFold(strings.TrimSpace(rateErr.Type), "day")
}

func (c *Client) buildRateLimitError(headers http.Header, retryAfter int, apiErr error) *RateLimitError {
	apiTyped, _ := apiErr.(*APIError)
	message := ""
	if apiTyped != nil {
		message = strings.TrimSpace(apiTyped.Message)
	}

	rateType := detectRateLimitType(headers, message)
	limit, remaining := extractRateLimitWindow(headers, rateType, c.RateLimitStats())
	resetAt := resolveRateLimitResetAt(headers, rateType, retryAfter)

	return &RateLimitError{
		Limit:     limit,
		Remaining: remaining,
		ResetAt:   resetAt,
		Type:      rateType,
	}
}

func detectRateLimitType(headers http.Header, message string) string {
	if hasAnyHeader(headers, "X-BurstLimit-Limit", "X-BurstLimit-Remaining", "X-BurstLimit-Reset", "X-BurstLimit-ResetTime") {
		return "minute"
	}

	normalizedMessage := strings.ToLower(strings.TrimSpace(message))
	if strings.Contains(normalizedMessage, "burstlimit") {
		return "minute"
	}

	if hasAnyHeader(headers, "X-RateLimit-Reset", "X-RateLimit-ResetTime") {
		return "day"
	}
	if strings.Contains(normalizedMessage, "daily ratelimit exceeded") {
		return "day"
	}

	// Default to burst/minute behavior when provider context is incomplete.
	return "minute"
}

func extractRateLimitWindow(headers http.Header, rateType string, stats RateLimitStats) (int, int) {
	switch rateType {
	case "day":
		limit := parseHeaderInt(headers, "X-RateLimit-Limit")
		remaining := parseHeaderInt(headers, "X-RateLimit-Remaining")
		if limit <= 0 {
			limit = stats.DayLimit
		}
		if remaining < 0 {
			remaining = 0
		}
		return limit, remaining
	default:
		limit := parseHeaderInt(headers, "X-BurstLimit-Limit")
		remaining := parseHeaderInt(headers, "X-BurstLimit-Remaining")
		if limit <= 0 {
			limit = stats.MinuteLimit
		}
		if remaining < 0 {
			remaining = 0
		}
		return limit, remaining
	}
}

func resolveRateLimitResetAt(headers http.Header, rateType string, retryAfter int) string {
	now := time.Now()
	if ts := parseHeaderResetTime(headers, rateType); !ts.IsZero() {
		return ts.UTC().Format(time.RFC3339)
	}

	if seconds := parseHeaderResetSeconds(headers, rateType); seconds > 0 {
		return now.Add(time.Duration(seconds) * time.Second).UTC().Format(time.RFC3339)
	}

	if retryAfter > 0 {
		return now.Add(time.Duration(retryAfter) * time.Second).UTC().Format(time.RFC3339)
	}
	return ""
}

func parseHeaderResetTime(headers http.Header, rateType string) time.Time {
	var value string
	switch rateType {
	case "day":
		value = strings.TrimSpace(headers.Get("X-RateLimit-ResetTime"))
	default:
		value = strings.TrimSpace(headers.Get("X-BurstLimit-ResetTime"))
	}
	if value == "" {
		return time.Time{}
	}

	layouts := []string{time.RFC3339, time.RFC1123, time.RFC1123Z}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func parseHeaderResetSeconds(headers http.Header, rateType string) int {
	switch rateType {
	case "day":
		return parseHeaderInt(headers, "X-RateLimit-Reset")
	default:
		return parseHeaderInt(headers, "X-BurstLimit-Reset")
	}
}

func hasAnyHeader(headers http.Header, keys ...string) bool {
	for _, key := range keys {
		if strings.TrimSpace(headers.Get(key)) != "" {
			return true
		}
	}
	return false
}

func parseHeaderInt(headers http.Header, key string) int {
	raw := strings.TrimSpace(headers.Get(key))
	if raw == "" {
		return -1
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return -1
	}
	return value
}

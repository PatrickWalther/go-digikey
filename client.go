package digikey

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	defaultBaseURL = "https://api.digikey.com"
	defaultTimeout = 30 * time.Second
)

// Client is the Digi-Key API client.
type Client struct {
	httpClient   *http.Client
	baseURL      string
	clientID     string
	tokenManager *tokenManager
	rateLimiter  *RateLimiter
	retryConfig  RetryConfig
	cache        Cache
	cacheConfig  CacheConfig
	locale       Locale
	localeMu     sync.RWMutex
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL (useful for testing).
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithLocale sets the locale for API requests.
func WithLocale(locale Locale) ClientOption {
	return func(c *Client) {
		c.locale = locale
	}
}

// WithRateLimiter sets a custom rate limiter.
func WithRateLimiter(rateLimiter *RateLimiter) ClientOption {
	return func(c *Client) {
		c.rateLimiter = rateLimiter
	}
}

// WithTokenURL sets a custom token URL (useful for testing).
func WithTokenURL(tokenURL string) ClientOption {
	return func(c *Client) {
		if c.tokenManager != nil {
			c.tokenManager.tokenURL = tokenURL
		}
	}
}

// WithRetryConfig sets the retry configuration.
func WithRetryConfig(config RetryConfig) ClientOption {
	return func(c *Client) {
		c.retryConfig = config
	}
}

// WithCache sets a custom cache implementation.
func WithCache(cache Cache) ClientOption {
	return func(c *Client) {
		c.cache = cache
	}
}

// WithCacheConfig sets the cache configuration.
func WithCacheConfig(config CacheConfig) ClientOption {
	return func(c *Client) {
		c.cacheConfig = config
	}
}

// WithoutCache disables caching.
func WithoutCache() ClientOption {
	return func(c *Client) {
		c.cacheConfig.Enabled = false
	}
}

// WithoutRetry disables retries.
func WithoutRetry() ClientOption {
	return func(c *Client) {
		c.retryConfig = NoRetry()
	}
}

// NewClient creates a new Digi-Key API client.
func NewClient(clientID, clientSecret string, opts ...ClientOption) *Client {
	cacheConfig := DefaultCacheConfig()

	c := &Client{
		httpClient:  &http.Client{Timeout: defaultTimeout},
		baseURL:     defaultBaseURL,
		clientID:    clientID,
		locale:      DefaultLocale(),
		rateLimiter: NewRateLimiter(),
		retryConfig: DefaultRetryConfig(),
		cacheConfig: cacheConfig,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.tokenManager = newTokenManager(c.httpClient, clientID, clientSecret, "")

	// Initialize default cache if caching is enabled and no custom cache was provided
	if c.cacheConfig.Enabled && c.cache == nil {
		c.cache = NewMemoryCache(c.cacheConfig.DetailsTTL)
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetLocale updates the locale for subsequent requests.
func (c *Client) SetLocale(locale Locale) {
	c.localeMu.Lock()
	defer c.localeMu.Unlock()
	c.locale = locale
}

// getLocale returns the current locale (thread-safe).
func (c *Client) getLocale() Locale {
	c.localeMu.RLock()
	defer c.localeMu.RUnlock()
	return c.locale
}

// RateLimitStats returns current rate limit usage statistics.
func (c *Client) RateLimitStats() RateLimitStats {
	return c.rateLimiter.Stats()
}

// ClearCache clears all cached responses.
func (c *Client) ClearCache() {
	if mc, ok := c.cache.(*MemoryCache); ok {
		mc.Clear()
	}
}

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
// Returns (error, statusCode, shouldRetry).
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

	// Handle rate limiting (429)
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		c.rateLimiter.UpdateFromResponse(retryAfter)
		apiErr := c.handleErrorResponse(resp.StatusCode, respBody, resp.Header)
		return resp.StatusCode, true, apiErr
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
}

// handleErrorResponse parses error responses from the API.
func (c *Client) handleErrorResponse(statusCode int, body []byte, headers http.Header) error {
	apiErr := &APIError{
		StatusCode: statusCode,
		RequestID:  headers.Get("X-Request-Id"),
	}

	var errResp struct {
		Message string `json:"message"`
		Details string `json:"details"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil {
		apiErr.Message = errResp.Message
		apiErr.Details = errResp.Details
	} else {
		apiErr.Message = http.StatusText(statusCode)
		apiErr.Details = string(body)
	}

	return apiErr
}

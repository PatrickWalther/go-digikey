package digikey

import (
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
	customerID   string
	tokenURL     string
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
		c.tokenURL = tokenURL
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

// WithCustomerID sets the customer ID for API requests.
// When set, the X-DIGIKEY-Customer-Id header is included in all requests.
func WithCustomerID(customerID string) ClientOption {
	return func(c *Client) {
		c.customerID = customerID
	}
}

// NewClient creates a new Digi-Key API client.
func NewClient(clientID, clientSecret string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient:  &http.Client{Timeout: defaultTimeout},
		baseURL:     defaultBaseURL,
		clientID:    clientID,
		locale:      DefaultLocale(),
		rateLimiter: NewRateLimiter(),
		retryConfig: DefaultRetryConfig(),
		cacheConfig: DefaultCacheConfig(),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.tokenManager = newTokenManager(c.httpClient, clientID, clientSecret, c.tokenURL)

	if c.cacheConfig.Enabled && c.cache == nil {
		c.cache = NewMemoryCache(c.cacheConfig.DetailsTTL)
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

// Close releases resources held by the client.
// Always call Close when done with the client to prevent goroutine leaks.
func (c *Client) Close() error {
	if mc, ok := c.cache.(*MemoryCache); ok {
		mc.Close()
	}
	return nil
}

// ClearCache clears all cached responses.
func (c *Client) ClearCache() {
	if mc, ok := c.cache.(*MemoryCache); ok {
		mc.Clear()
	}
}


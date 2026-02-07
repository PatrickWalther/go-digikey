package digikey

import (
	"context"
	"net/http"
	"testing"
	"time"
)

// TestNewClient tests client creation with default options.
func TestNewClient(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()

	if client == nil {
		t.Fatal("expected non-nil client")
	}

	if client.baseURL != defaultBaseURL {
		t.Errorf("expected base URL %s, got %s", defaultBaseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Fatal("expected non-nil HTTP client")
	}

	if client.rateLimiter == nil {
		t.Fatal("expected non-nil rate limiter")
	}

	if client.clientID != "test-id" {
		t.Errorf("expected client ID test-id, got %s", client.clientID)
	}
}

// TestNewClientWithCustomHTTPClient tests client creation with custom HTTP client.
func TestNewClientWithCustomHTTPClient(t *testing.T) {
	customHTTPClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	client := NewClient("test-id", "test-secret", WithHTTPClient(customHTTPClient))
	defer client.Close()

	if client.httpClient != customHTTPClient {
		t.Error("expected same HTTP client instance")
	}
}

// TestNewClientWithRetryConfig tests client creation with custom retry config.
func TestNewClientWithRetryConfig(t *testing.T) {
	config := RetryConfig{
		MaxRetries:     5,
		InitialBackoff: 100 * time.Millisecond,
	}

	client := NewClient("test-id", "test-secret", WithRetryConfig(config))
	defer client.Close()

	if client.retryConfig.MaxRetries != 5 {
		t.Errorf("expected max retries 5, got %d", client.retryConfig.MaxRetries)
	}
}

// TestNewClientWithLocale tests client creation with custom locale.
func TestNewClientWithLocale(t *testing.T) {
	locale := Locale{
		Site:     "en-US",
		Language: "en",
		Currency: "USD",
	}

	client := NewClient("test-id", "test-secret", WithLocale(locale))
	defer client.Close()

	currentLocale := client.getLocale()
	if currentLocale.Site != locale.Site {
		t.Errorf("expected locale site %s, got %s", locale.Site, currentLocale.Site)
	}
}

// TestNewClientWithCache tests client creation with custom cache.
func TestNewClientWithCache(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	client := NewClient("test-id", "test-secret", WithCache(cache))
	defer client.Close()

	if client.cache == nil {
		t.Fatal("expected non-nil cache")
	}

	if client.cache != cache {
		t.Error("expected same cache instance")
	}
}

// TestNewClientWithCacheDisabled tests client creation with cache disabled.
func TestNewClientWithCacheDisabled(t *testing.T) {
	client := NewClient("test-id", "test-secret", WithoutCache())
	defer client.Close()

	if client.cacheConfig.Enabled {
		t.Error("expected cache to be disabled")
	}
}

// TestNewClientWithRetryDisabled tests client creation with retry disabled.
func TestNewClientWithRetryDisabled(t *testing.T) {
	client := NewClient("test-id", "test-secret", WithoutRetry())
	defer client.Close()

	if client.retryConfig.MaxRetries != 0 {
		t.Errorf("expected max retries 0, got %d", client.retryConfig.MaxRetries)
	}
}

// TestSetLocale tests locale changes.
func TestSetLocale(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()

	locale1 := client.getLocale()

	newLocale := Locale{
		Site:     "de-DE",
		Language: "de",
		Currency: "EUR",
	}
	client.SetLocale(newLocale)

	locale2 := client.getLocale()
	if locale2.Site == locale1.Site {
		t.Error("expected locale to change")
	}

	if locale2.Site != "de-DE" {
		t.Errorf("expected locale site de-DE, got %s", locale2.Site)
	}
}

// TestRateLimitStats tests rate limit stats retrieval.
func TestRateLimitStats(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()

	stats := client.RateLimitStats()
	if stats.MinuteUsed < 0 {
		t.Errorf("expected non-negative minute used, got %d", stats.MinuteUsed)
	}
}

// TestClearCache tests cache clearing.
func TestClearCache(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()

	// Should not panic
	client.ClearCache()
}

// TestClientWithBaseURL tests WithBaseURL option
func TestClientWithBaseURL(t *testing.T) {
	client := NewClient("id", "secret", WithBaseURL("https://custom.example.com"))
	defer client.Close()
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// TestClientWithRateLimiter tests WithRateLimiter option
func TestClientWithRateLimiter(t *testing.T) {
	limiter := NewRateLimiterWithLimits(100, 1000)
	client := NewClient("id", "secret", WithRateLimiter(limiter))
	defer client.Close()
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// TestClientWithTokenURL tests WithTokenURL option
func TestClientWithTokenURL(t *testing.T) {
	customURL := "https://custom.example.com/token"
	client := NewClient("id", "secret", WithTokenURL(customURL))
	defer client.Close()
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	// Verify the token URL was passed through to the token manager
	if client.tokenManager.tokenURL != customURL {
		t.Errorf("expected token URL %s, got %s", customURL, client.tokenManager.tokenURL)
	}
}

// TestNewClientOptionsAppliedOnce tests that options are not applied twice
func TestNewClientOptionsAppliedOnce(t *testing.T) {
	callCount := 0
	counterOpt := func(c *Client) {
		callCount++
	}
	client := NewClient("id", "secret", ClientOption(counterOpt))
	defer client.Close()
	if callCount != 1 {
		t.Errorf("expected option to be applied once, got %d times", callCount)
	}
}

// TestClientWithCacheConfig tests WithCacheConfig option
func TestClientWithCacheConfig(t *testing.T) {
	config := CacheConfig{
		Enabled:    true,
		SearchTTL:  10 * time.Minute,
		DetailsTTL: 5 * time.Minute,
	}
	client := NewClient("id", "secret", WithCacheConfig(config))
	defer client.Close()
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// TestNewClientWithCustomerID tests WithCustomerID option.
func TestNewClientWithCustomerID(t *testing.T) {
	client := NewClient("test-id", "test-secret", WithCustomerID("12345"))
	defer client.Close()

	if client.customerID != "12345" {
		t.Errorf("expected customer ID '12345', got '%s'", client.customerID)
	}
}

// TestNewClientWithoutCustomerID tests that customer ID defaults to empty.
func TestNewClientWithoutCustomerID(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()

	if client.customerID != "" {
		t.Errorf("expected empty customer ID, got '%s'", client.customerID)
	}
}

// TestClientClose tests that Close works without error.
func TestClientClose(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	err := client.Close()
	if err != nil {
		t.Errorf("expected no error from Close, got %v", err)
	}
}

// TestClientCloseWithoutCache tests Close on a client with cache disabled.
func TestClientCloseWithoutCache(t *testing.T) {
	client := NewClient("test-id", "test-secret", WithoutCache())
	err := client.Close()
	if err != nil {
		t.Errorf("expected no error from Close, got %v", err)
	}
}

// TestContextTimeout tests that context timeout is respected.
func TestContextTimeout(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// This will fail because we don't have real credentials,
	// but we want to verify context timeout handling
	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{
		Keywords: "test",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

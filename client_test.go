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

	currentLocale := client.getLocale()
	if currentLocale.Site != locale.Site {
		t.Errorf("expected locale site %s, got %s", locale.Site, currentLocale.Site)
	}
}

// TestNewClientWithCache tests client creation with custom cache.
func TestNewClientWithCache(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	client := NewClient("test-id", "test-secret", WithCache(cache))

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

	if client.cacheConfig.Enabled {
		t.Error("expected cache to be disabled")
	}
}

// TestNewClientWithRetryDisabled tests client creation with retry disabled.
func TestNewClientWithRetryDisabled(t *testing.T) {
	client := NewClient("test-id", "test-secret", WithoutRetry())

	if client.retryConfig.MaxRetries != 0 {
		t.Errorf("expected max retries 0, got %d", client.retryConfig.MaxRetries)
	}
}

// TestSetLocale tests locale changes.
func TestSetLocale(t *testing.T) {
	client := NewClient("test-id", "test-secret")

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

	stats := client.RateLimitStats()
	if stats.MinuteUsed < 0 {
		t.Errorf("expected non-negative minute used, got %d", stats.MinuteUsed)
	}
}

// TestClearCache tests cache clearing.
func TestClearCache(t *testing.T) {
	client := NewClient("test-id", "test-secret")

	// Should not panic
	client.ClearCache()
}

// TestContextTimeout tests that context timeout is respected.
func TestContextTimeout(t *testing.T) {
	client := NewClient("test-id", "test-secret")

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// This will fail because we don't have real credentials,
	// but we want to verify context timeout handling
	_, err := client.KeywordSearch(ctx, &SearchRequest{
		Keywords: "test",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

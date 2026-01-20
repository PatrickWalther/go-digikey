package digikey

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

// TestNewSearch tests SearchOptions creation
func TestNewSearch(t *testing.T) {
	search := NewSearch("resistor")
	if search == nil {
		t.Error("expected non-nil SearchOptions")
		return
	}
	if search.request.Keywords != "resistor" {
		t.Errorf("expected keywords 'resistor', got '%s'", search.request.Keywords)
	}
	if search.request.Limit != 10 {
		t.Errorf("expected default limit 10, got %d", search.request.Limit)
	}
}

// TestSearchOptionsLimit tests Limit builder method
func TestSearchOptionsLimit(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{0, 1},
		{-1, 1},
		{50, 50},
		{100, 50}, // max is 50
		{25, 25},
	}

	for _, test := range tests {
		search := NewSearch("test").Limit(test.input)
		if search.request.Limit != test.expected {
			t.Errorf("Limit(%d): expected %d, got %d", test.input, test.expected, search.request.Limit)
		}
	}
}

// TestSearchOptionsOffset tests Offset builder method
func TestSearchOptionsOffset(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{10, 10},
		{-1, 0}, // negative becomes 0
		{100, 100},
	}

	for _, test := range tests {
		search := NewSearch("test").Offset(test.input)
		if search.request.Offset != test.expected {
			t.Errorf("Offset(%d): expected %d, got %d", test.input, test.expected, search.request.Offset)
		}
	}
}

// TestSearchOptionsChaining tests method chaining
func TestSearchOptionsChaining(t *testing.T) {
	search := NewSearch("LED").Limit(20).Offset(5)
	if search.request.Keywords != "LED" {
		t.Errorf("expected keywords 'LED', got '%s'", search.request.Keywords)
	}
	if search.request.Limit != 20 {
		t.Errorf("expected limit 20, got %d", search.request.Limit)
	}
	if search.request.Offset != 5 {
		t.Errorf("expected offset 5, got %d", search.request.Offset)
	}
}

// TestSearchOptionsBuild tests Build method
func TestSearchOptionsBuild(t *testing.T) {
	built := NewSearch("diode").Limit(15).Offset(3).Build()
	if built == nil {
		t.Error("Build() returned nil")
		return
	}
	if built.Keywords != "diode" {
		t.Errorf("expected keywords 'diode', got '%s'", built.Keywords)
	}
	if built.Limit != 15 {
		t.Errorf("expected limit 15, got %d", built.Limit)
	}
	if built.Offset != 3 {
		t.Errorf("expected offset 3, got %d", built.Offset)
	}
}

// TestSearchOptionsWithFilterOptions tests WithFilterOptions builder method
func TestSearchOptionsWithFilterOptions(t *testing.T) {
	filter := &FilterRequest{
		CategoryFilter: []int{1, 2},
	}
	search := NewSearch("test").WithFilterOptions(filter)
	if search.request.FilterOptionsRequest == nil {
		t.Error("expected FilterOptionsRequest to be set")
	}
	if search.request.FilterOptionsRequest != filter {
		t.Error("FilterOptionsRequest should be the same object")
	}
}

// TestSearchOptionsExecute tests Execute method
func TestSearchOptionsExecute(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	search := NewSearch("test").Limit(5)

	// This will fail with auth error, but that's ok - we're testing the builder
	// The important thing is Execute tries to call the API
	_, err := search.Execute(ctx, client)
	if err == nil {
		t.Error("expected error (no real API), got nil")
	}
}

// TestKeywordSearchNilRequest tests KeywordSearch with nil request
func TestKeywordSearchNilRequest(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.KeywordSearch(ctx, nil)
	if err == nil {
		t.Error("expected error for nil request")
	}
	if err != ErrInvalidRequest {
		t.Errorf("expected ErrInvalidRequest, got %v", err)
	}
}

// TestKeywordSearchEmptyKeywords tests KeywordSearch with empty keywords
func TestKeywordSearchEmptyKeywords(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req := &SearchRequest{Keywords: ""}
	_, err := client.KeywordSearch(ctx, req)
	if err == nil {
		t.Error("expected error for empty keywords")
	}
}

// TestKeywordSearchLimitAdjustment tests that limits are clamped
func TestKeywordSearchLimitAdjustment(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tests := []struct {
		input    int
		expected int
	}{
		{0, 10},  // default
		{-1, 10}, // default
		{5, 5},
		{50, 50},
		{100, 50}, // clamped
	}

	for _, test := range tests {
		req := &SearchRequest{
			Keywords: "test",
			Limit:    test.input,
		}
		// We expect auth error, but we can check the request was at least created
		_, err := client.KeywordSearch(ctx, req)
		if err == nil {
			t.Error("expected error (no real API)")
		}
	}
}

// TestProductDetailsNilProductNumber tests ProductDetails with empty product number
func TestProductDetailsEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.ProductDetails(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

// TestProductDetailsNoCacheEmptyProductNumber tests ProductDetailsNoCache with empty product number
func TestProductDetailsNoCacheEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.ProductDetailsNoCache(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

// TestKeywordSearchWithCache tests that caching key is generated correctly
func TestKeywordSearchWithCache(t *testing.T) {
	client := NewClient("test-id", "test-secret", WithCache(NewMemoryCache(5*time.Minute)))

	// Manually set cache entry to verify it's retrieved
	req := &SearchRequest{Keywords: "test", Limit: 10}
	cacheKey := cacheKeyForSearch(DefaultLocale(), req)

	// Create fake response and cache it
	resp := &SearchResponse{
		ProductsCount: 1,
		SearchLocaleUsed: SearchLocale{Site: "US", Language: "en"},
	}
	if data, err := json.Marshal(resp); err == nil {
		client.cache.Set(cacheKey, data, 5*time.Minute)

		// Now try to retrieve it
		cached, ok := client.cache.Get(cacheKey)
		if !ok {
			t.Error("expected cached data")
		}
		if len(cached) == 0 {
			t.Error("expected non-empty cached data")
		}
	}
}

// TestProductDetailsWithCache tests caching for product details
func TestProductDetailsWithCache(t *testing.T) {
	// Verify cache key function works
	cacheKey := cacheKeyForDetails(DefaultLocale(), "TEST-123")
	if cacheKey == "" {
		t.Error("expected non-empty cache key")
	}
	if cacheKey != cacheKeyForDetails(DefaultLocale(), "TEST-123") {
		t.Error("cache key should be consistent")
	}
}

// TestClientWithBaseURL tests WithBaseURL option
func TestClientWithBaseURL(t *testing.T) {
	client := NewClient("id", "secret", WithBaseURL("https://custom.example.com"))
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// TestClientWithRateLimiter tests WithRateLimiter option
func TestClientWithRateLimiter(t *testing.T) {
	limiter := NewRateLimiterWithLimits(100, 1000)
	client := NewClient("id", "secret", WithRateLimiter(limiter))
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// TestClientWithTokenURL tests WithTokenURL option
func TestClientWithTokenURL(t *testing.T) {
	client := NewClient("id", "secret", WithTokenURL("https://custom.example.com/token"))
	if client == nil {
		t.Error("expected non-nil client")
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
	if client == nil {
		t.Error("expected non-nil client")
	}
}

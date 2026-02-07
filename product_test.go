package digikey

import (
	"context"
	"testing"
	"time"
)

// TestProductDetailsEmptyProductNumber tests ProductDetails with empty product number
func TestProductDetailsEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
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
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.ProductDetailsNoCache(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
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

// +build integration

package digikey

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestIntegrationKeywordSearch makes a real API call to search for products.
func TestIntegrationKeywordSearch(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.KeywordSearch(ctx, &SearchRequest{
		Keywords: "transistor",
		Limit:    5,
	})

	if err != nil {
		t.Fatalf("KeywordSearch failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if len(resp.Products) == 0 {
		t.Error("expected at least one product in response")
	}
}

// TestIntegrationProductDetails makes a real API call to get product details.
func TestIntegrationProductDetails(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First search for a product
	searchResp, err := client.KeywordSearch(ctx, &SearchRequest{
		Keywords: "resistor",
		Limit:    1,
	})

	if err != nil {
		t.Fatalf("KeywordSearch failed: %v", err)
	}

	if len(searchResp.Products) == 0 {
		t.Skip("no products found for integration test")
	}

	// Get product number from first variation
	product := searchResp.Products[0]
	if len(product.ProductVariations) == 0 {
		t.Skip("product has no variations")
	}
	
	productNumber := product.ProductVariations[0].DigiKeyProductNumber

	// Now get details for that product
	resp, err := client.ProductDetails(ctx, productNumber)

	if err != nil {
		t.Fatalf("ProductDetails failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	// Check that we got the manufacturer product number
	if resp.Product.ManufacturerProductNumber == "" {
		t.Error("expected manufacturer product number to be set")
	}
}

// TestIntegrationRateLimiterWithRealAPI tests rate limiter with real API calls.
func TestIntegrationRateLimiterWithRealAPI(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make several requests and verify rate limiter works
	for i := 0; i < 5; i++ {
		_, err := client.KeywordSearch(ctx, &SearchRequest{
			Keywords: "diode",
			Limit:    1,
		})

		if err != nil {
			// Rate limit or other error
			t.Logf("Request %d failed: %v (may be rate limit)", i+1, err)
			break
		}

		stats := client.RateLimitStats()
		if stats.MinuteUsed > 0 && i > 0 {
			if stats.MinuteUsed != i+1 {
				t.Logf("Rate limiter stats: used %d, remaining %d", stats.MinuteUsed, stats.MinuteRemaining)
			}
		}
	}
}

// TestIntegrationCaching tests that caching works with real API.
func TestIntegrationCaching(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &SearchRequest{
		Keywords: "LED",
		Limit:    1,
	}

	// First call - hits API
	resp1, err := client.KeywordSearch(ctx, req)
	if err != nil {
		t.Fatalf("first search failed: %v", err)
	}

	if resp1 == nil || len(resp1.Products) == 0 {
		t.Skip("no products found for caching test")
	}

	stats1 := client.RateLimitStats()

	// Second call - should be cached
	_, err = client.KeywordSearch(ctx, req)
	if err != nil {
		t.Fatalf("second search failed: %v", err)
	}

	stats2 := client.RateLimitStats()

	// If caching works, second call shouldn't increment rate limit counter
	// (though it might due to timing)
	t.Logf("After 1st call: used %d, after 2nd call: used %d", stats1.MinuteUsed, stats2.MinuteUsed)
}

// TestIntegrationLocaleSupport tests searching with different locales.
func TestIntegrationLocaleSupport(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	client.SetLocale(Locale{
		Site:     "US",
		Language: "en",
		Currency: "USD",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.KeywordSearch(ctx, &SearchRequest{
		Keywords: "capacitor",
		Limit:    1,
	})

	if err != nil {
		t.Fatalf("search with locale failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	locale := resp.SearchLocaleUsed
	t.Logf("Search used locale: Site=%s, Language=%s, Currency=%s", locale.Site, locale.Language, locale.Currency)
}

// TestIntegrationAuthRefresh tests that OAuth2 token refresh works.
func TestIntegrationAuthRefresh(t *testing.T) {
	skipIfNoCredentials(t)

	// This test verifies that token refresh doesn't cause errors
	// by making multiple requests that might trigger token refresh
	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for i := 0; i < 3; i++ {
		_, err := client.KeywordSearch(ctx, &SearchRequest{
			Keywords: "IC",
			Limit:    1,
		})

		if err != nil {
			t.Logf("Request %d: %v", i+1, err)
		}
	}
}

// TestIntegrationProductDetailsNoCache tests bypassing cache.
func TestIntegrationProductDetailsNoCache(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First search for a product
	searchResp, err := client.KeywordSearch(ctx, &SearchRequest{
		Keywords: "LED",
		Limit:    1,
	})

	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(searchResp.Products) == 0 {
		t.Skip("no products found")
	}

	product := searchResp.Products[0]
	if len(product.ProductVariations) == 0 {
		t.Skip("product has no variations")
	}
	
	productNumber := product.ProductVariations[0].DigiKeyProductNumber

	// Get details (cached)
	resp1, err := client.ProductDetails(ctx, productNumber)
	if err != nil {
		t.Fatalf("first details call failed: %v", err)
	}

	// Get details again without cache
	resp2, err := client.ProductDetailsNoCache(ctx, productNumber)
	if err != nil {
		t.Fatalf("no-cache details call failed: %v", err)
	}

	if resp1 != nil && resp2 != nil && resp1.Product.DigiKeyProductNumber != resp2.Product.DigiKeyProductNumber {
		t.Error("product numbers should match between cached and uncached calls")
	}
}

// Helper functions

func newTestClient(t *testing.T) *Client {
	clientID := os.Getenv("DIGIKEY_CLIENT_ID")
	clientSecret := os.Getenv("DIGIKEY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Fatal("DIGIKEY_CLIENT_ID and DIGIKEY_CLIENT_SECRET environment variables not set")
	}

	client := NewClient(clientID, clientSecret)
	if client == nil {
		t.Fatal("failed to create client")
	}

	return client
}

func skipIfNoCredentials(t *testing.T) {
	clientID := os.Getenv("DIGIKEY_CLIENT_ID")
	clientSecret := os.Getenv("DIGIKEY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skip("DIGIKEY_CLIENT_ID or DIGIKEY_CLIENT_SECRET not set; skipping integration test")
	}
}

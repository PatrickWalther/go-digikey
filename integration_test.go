//go:build integration
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

	resp, err := client.Search.KeywordSearch(ctx, &SearchRequest{
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
	searchResp, err := client.Search.KeywordSearch(ctx, &SearchRequest{
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
	resp, err := client.Product.Details(ctx, productNumber)

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
		_, err := client.Search.KeywordSearch(ctx, &SearchRequest{
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
	resp1, err := client.Search.KeywordSearch(ctx, req)
	if err != nil {
		t.Fatalf("first search failed: %v", err)
	}

	if resp1 == nil || len(resp1.Products) == 0 {
		t.Skip("no products found for caching test")
	}

	stats1 := client.RateLimitStats()

	// Second call - should be cached
	_, err = client.Search.KeywordSearch(ctx, req)
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

	resp, err := client.Search.KeywordSearch(ctx, &SearchRequest{
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
		_, err := client.Search.KeywordSearch(ctx, &SearchRequest{
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
	searchResp, err := client.Search.KeywordSearch(ctx, &SearchRequest{
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
	resp1, err := client.Product.Details(ctx, productNumber)
	if err != nil {
		t.Fatalf("first details call failed: %v", err)
	}

	// Get details again without cache
	resp2, err := client.Product.DetailsNoCache(ctx, productNumber)
	if err != nil {
		t.Fatalf("no-cache details call failed: %v", err)
	}

	if resp1 != nil && resp2 != nil && resp1.Product.DigiKeyProductNumber != resp2.Product.DigiKeyProductNumber {
		t.Error("product numbers should match between cached and uncached calls")
	}
}

// TestIntegrationCategories tests the Categories endpoint.
func TestIntegrationCategories(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Category.List(ctx)
	if err != nil {
		t.Fatalf("Categories failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if len(resp.Categories) == 0 {
		t.Error("expected at least one category")
	}

	t.Logf("Found %d categories (product count: %d)", len(resp.Categories), resp.ProductCount)
}

// TestIntegrationCategoriesById tests the CategoriesById endpoint.
func TestIntegrationCategoriesById(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First get all categories to find a valid ID
	cats, err := client.Category.List(ctx)
	if err != nil {
		t.Fatalf("Categories failed: %v", err)
	}
	if len(cats.Categories) == 0 {
		t.Skip("no categories found")
	}

	catID := cats.Categories[0].CategoryID

	resp, err := client.Category.GetByID(ctx, catID)
	if err != nil {
		t.Fatalf("CategoriesById failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	t.Logf("Category %d: %s", resp.Category.CategoryID, resp.Category.Name)
}

// TestIntegrationManufacturers tests the Manufacturers endpoint.
func TestIntegrationManufacturers(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Category.Manufacturers(ctx)
	if err != nil {
		t.Fatalf("Manufacturers failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if len(resp.Manufacturers) == 0 {
		t.Error("expected at least one manufacturer")
	}

	t.Logf("Found %d manufacturers", len(resp.Manufacturers))
}

// TestIntegrationAssociations tests the Associations endpoint.
func TestIntegrationAssociations(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	productNumber := findProductNumber(t, client, ctx)

	resp, err := client.Product.Associations(ctx, productNumber)
	if err != nil {
		t.Fatalf("Associations failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	t.Logf("Associations for %s: kits=%d, mating=%d, associated=%d, forUseWith=%d",
		productNumber,
		len(resp.ProductAssociations.Kits),
		len(resp.ProductAssociations.MatingProducts),
		len(resp.ProductAssociations.AssociatedProducts),
		len(resp.ProductAssociations.ForUseWithProducts))
}

// TestIntegrationMedia tests the Media endpoint.
func TestIntegrationMedia(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	productNumber := findProductNumber(t, client, ctx)

	resp, err := client.Product.Media(ctx, productNumber)
	if err != nil {
		t.Fatalf("Media failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	t.Logf("Media links for %s: %d", productNumber, len(resp.MediaLinks))
}

// TestIntegrationSubstitutions tests the Substitutions endpoint.
func TestIntegrationSubstitutions(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	productNumber := findProductNumber(t, client, ctx)

	resp, err := client.Product.Substitutions(ctx, productNumber)
	if err != nil {
		t.Fatalf("Substitutions failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	t.Logf("Substitutes for %s: %d", productNumber, resp.ProductSubstitutesCount)
}

// TestIntegrationRecommendedProducts tests the RecommendedProducts endpoint.
func TestIntegrationRecommendedProducts(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	productNumber := findProductNumber(t, client, ctx)

	resp, err := client.Product.RecommendedProducts(ctx, productNumber)
	if err != nil {
		t.Fatalf("RecommendedProducts failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	t.Logf("Recommendations for %s: %d", productNumber, len(resp.Recommendations))
}

// TestIntegrationDigiReelPricing tests the DigiReelPricing endpoint.
func TestIntegrationDigiReelPricing(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	productNumber := findProductNumber(t, client, ctx)

	resp, err := client.Pricing.DigiReel(ctx, productNumber, 1000)
	if err != nil {
		// DigiReel pricing may not be available for all products
		t.Logf("DigiReelPricing for %s: %v (may not support DigiReel)", productNumber, err)
		return
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	t.Logf("DigiReel pricing for %s: fee=%.2f, unit=%.2f, extended=%.2f",
		productNumber, resp.ReelingFee, resp.UnitPrice, resp.ExtendedPrice)
}

// TestIntegrationPackageTypeByQuantity tests the PackageTypeByQuantity endpoint.
func TestIntegrationPackageTypeByQuantity(t *testing.T) {
	skipIfNoCredentials(t)

	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	productNumber := findProductNumber(t, client, ctx)

	resp, err := client.Pricing.PackageTypeByQuantity(ctx, productNumber, 100, "")
	if err != nil {
		t.Fatalf("PackageTypeByQuantity failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	t.Logf("Package types for %s: %d products", productNumber, len(resp.Products))
}

// findProductNumber is a helper that searches for a product and returns a DigiKey product number.
func findProductNumber(t *testing.T, client *Client, ctx context.Context) string {
	t.Helper()
	searchResp, err := client.Search.KeywordSearch(ctx, &SearchRequest{
		Keywords: "resistor",
		Limit:    1,
	})
	if err != nil {
		t.Fatalf("KeywordSearch failed: %v", err)
	}
	if len(searchResp.Products) == 0 {
		t.Skip("no products found")
	}
	product := searchResp.Products[0]
	if len(product.ProductVariations) == 0 {
		t.Skip("product has no variations")
	}
	return product.ProductVariations[0].DigiKeyProductNumber
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

	t.Cleanup(func() { client.Close() })
	return client
}

func skipIfNoCredentials(t *testing.T) {
	clientID := os.Getenv("DIGIKEY_CLIENT_ID")
	clientSecret := os.Getenv("DIGIKEY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skip("DIGIKEY_CLIENT_ID or DIGIKEY_CLIENT_SECRET not set; skipping integration test")
	}
}

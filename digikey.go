// Package digikey provides a Go client for the Digi-Key API v4.
//
// The client supports OAuth 2.0 client credentials flow authentication,
// automatic token caching and refresh, rate limiting, and locale support.
//
// # Quick Start
//
// Create a client with your Digi-Key API credentials:
//
//	client := digikey.NewClient(
//	    os.Getenv("DIGIKEY_CLIENT_ID"),
//	    os.Getenv("DIGIKEY_CLIENT_SECRET"),
//	)
//
// Search for products:
//
//	results, err := client.KeywordSearch(ctx, &digikey.SearchRequest{
//	    Keywords:    "STM32F4",
//	    RecordCount: 10,
//	})
//
// Or use the fluent search builder:
//
//	results, err := digikey.NewSearch("STM32F4").
//	    Limit(10).
//	    FilterByManufacturer(1).
//	    Execute(ctx, client)
//
// Get product details:
//
//	details, err := client.ProductDetails(ctx, "497-15360-ND")
//
// # Authentication
//
// The client uses OAuth 2.0 client credentials flow (2-legged authentication).
// Tokens are automatically cached and refreshed before expiry.
//
// # Rate Limiting
//
// The client tracks rate limits (120 requests/minute, 1000 requests/day)
// and returns ErrRateLimitExceeded when limits are reached. Check current
// usage with:
//
//	stats := client.RateLimitStats()
//
// # Locale Support
//
// Set the locale for pricing and availability:
//
//	client.SetLocale(digikey.Locale{
//	    Site:     "DE",
//	    Language: "de",
//	    Currency: "EUR",
//	})
//
// # Error Handling
//
// The package provides typed errors for common error conditions:
//
//	if errors.Is(err, digikey.ErrRateLimitExceeded) {
//	    // Handle rate limit
//	}
//	if errors.Is(err, digikey.ErrUnauthorized) {
//	    // Check credentials
//	}
package digikey

// Version is the current version of the go-digikey package.
const Version = "0.1.0"

// Package digikey provides a Go client for the Digi-Key Product Search API v4.
//
// The client supports all 11 endpoints of the DigiKey Product Search API,
// OAuth 2.0 client credentials flow authentication, automatic token caching
// and refresh, rate limiting, in-memory response caching, and locale support.
//
// Endpoints are organized into service groups accessed via the client:
//
//   - client.Search  — keyword search
//   - client.Product — product details, associations, media, substitutions, recommendations
//   - client.Category — categories and manufacturers
//   - client.Pricing — DigiReel pricing and package type by quantity
//
// # Quick Start
//
// Create a client with your Digi-Key API credentials:
//
//	client := digikey.NewClient(
//	    os.Getenv("DIGIKEY_CLIENT_ID"),
//	    os.Getenv("DIGIKEY_CLIENT_SECRET"),
//	)
//	defer client.Close()
//
// Search for products:
//
//	results, err := client.Search.KeywordSearch(ctx, &digikey.SearchRequest{
//	    Keywords: "STM32F4",
//	    Limit:    10,
//	})
//
// Or use the fluent search builder:
//
//	results, err := digikey.NewSearch("STM32F4").
//	    Limit(10).
//	    Execute(ctx, client)
//
// Get product details:
//
//	details, err := client.Product.Details(ctx, "497-15360-ND")
//
// Browse categories:
//
//	categories, err := client.Category.List(ctx)
//
// Get manufacturers:
//
//	manufacturers, err := client.Category.Manufacturers(ctx)
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
const Version = "1.0.0"

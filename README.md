# go-digikey

[![Go Reference](https://pkg.go.dev/badge/github.com/PatrickWalther/go-digikey.svg)](https://pkg.go.dev/github.com/PatrickWalther/go-digikey)
[![Go Report Card](https://goreportcard.com/badge/github.com/PatrickWalther/go-digikey)](https://goreportcard.com/report/github.com/PatrickWalther/go-digikey)
[![Tests](https://github.com/PatrickWalther/go-digikey/actions/workflows/test.yml/badge.svg)](https://github.com/PatrickWalther/go-digikey/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the [Digi-Key](https://www.digikey.com) Product Search API v4. **100% API coverage** — all 11 endpoints supported.

## Requirements

- **Go 1.22+** (tested on Go 1.22 and 1.23)
- **Digi-Key API credentials** (OAuth2 client ID and secret)
- No external dependencies beyond stdlib

## Features

- **100% API coverage** — all 11 Product Search API v4 endpoints
- OAuth 2.0 client credentials flow (2-legged authentication)
- Automatic token caching and refresh with 401 auto-retry
- In-memory response caching with configurable TTL
- Automatic retries with exponential backoff for transient errors
- Rate limiting (120 requests/minute, 1000 requests/day)
- Locale support (site, language, currency)
- Customer ID header support
- No external dependencies beyond stdlib
- Thread-safe for concurrent use

## Installation

```bash
go get github.com/PatrickWalther/go-digikey
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/PatrickWalther/go-digikey"
)

func main() {
    client := digikey.NewClient(
        os.Getenv("DIGIKEY_CLIENT_ID"),
        os.Getenv("DIGIKEY_CLIENT_SECRET"),
    )
    defer client.Close()

    ctx := context.Background()

    // Search for products
    results, err := client.KeywordSearch(ctx, &digikey.SearchRequest{
        Keywords: "STM32F4",
        Limit:    10,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, product := range results.Products {
        fmt.Printf("%s - %s\n",
            product.ManufacturerProductNumber,
            product.Description.ProductDescription)
    }
}
```

## Usage

### Creating a Client

```go
// Basic client with default options
client := digikey.NewClient(clientID, clientSecret)
defer client.Close()

// Client with custom options
client := digikey.NewClient(
    clientID,
    clientSecret,
    digikey.WithLocale(digikey.Locale{
        Site:     "DE",
        Language: "de",
        Currency: "EUR",
    }),
    digikey.WithCustomerID("12345"),
    digikey.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)
```

### Keyword Search

```go
// Basic search
results, err := client.KeywordSearch(ctx, &digikey.SearchRequest{
    Keywords: "STM32F4",
    Limit:    10,
})

// Using the fluent builder
results, err := digikey.NewSearch("STM32F4").
    Limit(20).
    Offset(0).
    WithIncludes("DigiKeyProductNumber").
    Execute(ctx, client)

// Search with filters (using spec-compliant FilterId types)
results, err := client.KeywordSearch(ctx, &digikey.SearchRequest{
    Keywords: "resistor",
    Limit:    10,
    FilterOptionsRequest: &digikey.FilterRequest{
        CategoryFilter:     digikey.NewFilterIds(53),
        ManufacturerFilter: digikey.NewFilterIds(1, 2, 3),
        StatusFilter:       digikey.NewFilterIds(0),
        MinimumQuantityAvailable: 100,
        SearchOptions:       []string{"InStock"},
    },
})
```

### Product Details

```go
details, err := client.ProductDetails(ctx, "497-15360-ND")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Product: %s\n", details.Product.ManufacturerProductNumber)
fmt.Printf("Price: $%.2f\n", details.Product.UnitPrice)
fmt.Printf("Stock: %d\n", details.Product.QuantityAvailable)
```

### Categories

```go
// Get all categories
categories, err := client.Categories(ctx)
for _, cat := range categories.Categories {
    fmt.Printf("%s (%d products)\n", cat.Name, cat.ProductCount)
}

// Get a specific category by ID
category, err := client.CategoriesById(ctx, 42)
fmt.Printf("Category: %s\n", category.Category.Name)
```

### Manufacturers

```go
manufacturers, err := client.Manufacturers(ctx)
for _, mfg := range manufacturers.Manufacturers {
    fmt.Printf("%s (ID: %d)\n", mfg.Name, mfg.ID)
}
```

### Associations

```go
assoc, err := client.Associations(ctx, "497-15360-ND")
for _, mate := range assoc.ProductAssociations.MatingProducts {
    fmt.Printf("Mating: %s\n", mate.ManufacturerProductNumber)
}
```

### Substitutions

```go
subs, err := client.Substitutions(ctx, "497-15360-ND")
for _, sub := range subs.ProductSubstitutes {
    fmt.Printf("%s (%s): %s\n", sub.ManufacturerProductNumber, sub.SubstituteType, sub.UnitPrice)
}
```

### Media

```go
media, err := client.Media(ctx, "497-15360-ND")
for _, link := range media.MediaLinks {
    fmt.Printf("%s: %s\n", link.MediaType, link.URL)
}
```

### Recommended Products

```go
recs, err := client.RecommendedProducts(ctx, "497-15360-ND")
for _, rec := range recs.Recommendations {
    for _, p := range rec.RecommendedProducts {
        fmt.Printf("Recommended: %s ($%.2f)\n", p.DigiKeyProductNumber, p.UnitPrice)
    }
}
```

### DigiReel Pricing

```go
pricing, err := client.DigiReelPricing(ctx, "497-15360-ND", 1000)
fmt.Printf("Reeling fee: $%.2f, Unit: $%.4f, Extended: $%.2f\n",
    pricing.ReelingFee, pricing.UnitPrice, pricing.ExtendedPrice)
```

### Package Type by Quantity

```go
pkgType, err := client.PackageTypeByQuantity(ctx, "497-15360-ND", 100, "CT")
for _, p := range pkgType.Products {
    fmt.Printf("%s: %d available, types: %v\n",
        p.DigiKeyProductNumber, p.QuantityAvailable, p.PackageTypes)
}
```

### Locale Support

```go
// Set locale for European site
client.SetLocale(digikey.Locale{
    Site:     "DE",
    Language: "de",
    Currency: "EUR",
})
```

### Rate Limit Monitoring

```go
stats := client.RateLimitStats()
fmt.Printf("Minute: %d/%d remaining\n", stats.MinuteRemaining, stats.MinuteLimit)
fmt.Printf("Day: %d/%d remaining\n", stats.DayRemaining, stats.DayLimit)
```

### Error Handling

```go
import "errors"

results, err := client.KeywordSearch(ctx, req)
if err != nil {
    if errors.Is(err, digikey.ErrRateLimitExceeded) {
        // Wait and retry
    }
    if errors.Is(err, digikey.ErrUnauthorized) {
        // Check credentials
    }
    if errors.Is(err, digikey.ErrNotFound) {
        // Product not found
    }

    // Check for API error details
    var apiErr *digikey.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: %s (status %d)\n", apiErr.Message, apiErr.StatusCode)
    }
}
```

## API Coverage

| Endpoint | Method | Go Method |
|----------|--------|-----------|
| `/products/v4/search/keyword` | POST | `KeywordSearch()` |
| `/products/v4/search/{productNumber}/productdetails` | GET | `ProductDetails()` |
| `/products/v4/search/categories` | GET | `Categories()` |
| `/products/v4/search/categories/{categoryId}` | GET | `CategoriesById()` |
| `/products/v4/search/manufacturers` | GET | `Manufacturers()` |
| `/products/v4/search/{productNumber}/associations` | GET | `Associations()` |
| `/products/v4/search/{productNumber}/substitutions` | GET | `Substitutions()` |
| `/products/v4/search/{productNumber}/media` | GET | `Media()` |
| `/products/v4/search/{productNumber}/recommendedproducts` | GET | `RecommendedProducts()` |
| `/products/v4/search/{productNumber}/digireelpricing` | GET | `DigiReelPricing()` |
| `/products/v4/search/packagetypebyquantity/{productNumber}` | GET | `PackageTypeByQuantity()` |

**11/11 endpoints** (100% coverage)

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| `DIGIKEY_CLIENT_ID` | OAuth 2.0 client ID |
| `DIGIKEY_CLIENT_SECRET` | OAuth 2.0 client secret |

### Client Options

| Option | Description |
|--------|-------------|
| `WithHTTPClient` | Custom HTTP client |
| `WithBaseURL` | Custom base URL (for testing) |
| `WithLocale` | Set request locale |
| `WithCustomerID` | Set customer ID for `X-DIGIKEY-Customer-Id` header |
| `WithRateLimiter` | Custom rate limiter |
| `WithTokenURL` | Custom OAuth token URL |
| `WithCache` | Custom cache implementation |
| `WithCacheConfig` | Configure cache TTLs |
| `WithoutCache` | Disable caching |
| `WithRetryConfig` | Custom retry configuration |
| `WithoutRetry` | Disable retries |

### Client Methods

| Method | Description |
|--------|-------------|
| `Close()` | Release resources (always call with `defer`) |
| `KeywordSearch()` | Search products by keyword |
| `ProductDetails()` | Get product details (cached) |
| `ProductDetailsNoCache()` | Get product details (bypass cache) |
| `Categories()` | Get all product categories (cached) |
| `CategoriesById()` | Get a specific category (cached) |
| `Manufacturers()` | Get all manufacturers (cached) |
| `Associations()` | Get product associations (cached) |
| `Substitutions()` | Get product substitutes (cached) |
| `Media()` | Get product media links (cached) |
| `RecommendedProducts()` | Get recommended products (cached) |
| `DigiReelPricing()` | Get DigiReel pricing (not cached) |
| `PackageTypeByQuantity()` | Get package types by quantity (not cached) |
| `SetLocale()` | Change request locale |
| `RateLimitStats()` | Get current rate limit usage |
| `ClearCache()` | Clear all cached responses |

## Caching

The client includes in-memory caching with configurable TTL:

```go
// Default: caching enabled with 5min search, 10min details, 15min lookup TTL
client := digikey.NewClient(clientID, clientSecret)
defer client.Close()

// Custom cache configuration
client := digikey.NewClient(
    clientID,
    clientSecret,
    digikey.WithCacheConfig(digikey.CacheConfig{
        Enabled:    true,
        SearchTTL:  2 * time.Minute,
        DetailsTTL: 5 * time.Minute,
        LookupTTL:  15 * time.Minute,
    }),
)

// Disable caching
client := digikey.NewClient(clientID, clientSecret, digikey.WithoutCache())

// Force refresh (bypass cache)
details, err := client.ProductDetailsNoCache(ctx, "497-15360-ND")

// Clear all cached data
client.ClearCache()
```

**Caching behavior by endpoint:**
- **Cached (SearchTTL):** `KeywordSearch`
- **Cached (DetailsTTL):** `ProductDetails`
- **Cached (LookupTTL):** `Categories`, `CategoriesById`, `Manufacturers`, `Associations`, `Substitutions`, `Media`, `RecommendedProducts`
- **Not cached:** `DigiReelPricing`, `PackageTypeByQuantity` (pricing is time-sensitive)

## Retries

The client automatically retries failed requests with exponential backoff:

- Retries on: 429 (rate limit), 500, 502, 503, 504, network timeouts
- Does not retry: 400, 401, 403, 404
- 401 errors trigger automatic token refresh and single retry
- Default: 3 retries with 500ms initial backoff, 2x multiplier

```go
// Custom retry configuration
client := digikey.NewClient(
    clientID,
    clientSecret,
    digikey.WithRetryConfig(digikey.RetryConfig{
        MaxRetries:     5,
        InitialBackoff: time.Second,
        MaxBackoff:     time.Minute,
        Multiplier:     2.0,
        Jitter:         0.1,
    }),
)

// Disable retries
client := digikey.NewClient(clientID, clientSecret, digikey.WithoutRetry())
```

## Rate Limits

Digi-Key API enforces the following rate limits:

- **Per Minute**: 120 requests
- **Per Day**: 1000 requests

The client tracks these limits locally and returns `ErrRateLimitExceeded` before making requests that would exceed them.

## Breaking Changes in v0.2.0

- `FilterRequest.CategoryFilter`, `ManufacturerFilter`, `StatusFilter` changed from `[]int` to `[]FilterId`
- `FilterRequest.PackageTypeFilter` renamed to `PackagingFilter` and changed from `[]int` to `[]FilterId`
- `FilterRequest.ParameterFilterRequest` changed from `[]ParameterFilterRequest` to `*ParameterFilterRequest`
- `ParameterFilterRequest` restructured with `CategoryFilter *FilterId` and `ParameterFilters []ParametricFilter`
- `ParametricFilter` restructured with `FilterValues []FilterId` (was `ValueIDs []string`)
- Removed `Filters` and `SortOptions` types (unused)
- Use `NewFilterId(id)` and `NewFilterIds(ids...)` convenience constructors

## Testing

### Unit Tests (Fast, No API Calls)

Run fast unit tests that don't require API credentials:

```bash
# Run all unit tests
go test -v -short ./...

# Run only short tests (excludes integration tests)
go test -v -short -run TestMemoryCache ./...
```

### Integration Tests (Real API Calls)

Run against real Digi-Key API with actual credentials:

#### 1. Setup Credentials

Create a `.env` file (copy from `.env.example`):

```bash
cp .env.example .env
# Edit .env and add your credentials from https://developer.digikey.com/
```

#### 2. Run Integration Tests Locally

```bash
# Load credentials from .env and run integration tests
source .env  # On Windows: set -a; source .env; set +a
go test -v -tags=integration -run Integration ./...

# Or pass credentials directly
DIGIKEY_CLIENT_ID=your-id DIGIKEY_CLIENT_SECRET=your-secret \
  go test -v -tags=integration ./...
```

#### 3. GitHub Actions Integration Tests

Set up GitHub repository secrets:

1. Go to: **Settings** > **Secrets and variables** > **Actions** > **New repository secret**
2. Add two secrets:
   - `DIGIKEY_CLIENT_ID`: Your OAuth client ID
   - `DIGIKEY_CLIENT_SECRET`: Your OAuth client secret

Integration tests run automatically on push to main/develop branches.

## License

MIT License - see [LICENSE](LICENSE) for details.

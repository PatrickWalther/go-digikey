# go-digikey

A Go client library for the Digi-Key API v4.

## Features

- OAuth 2.0 client credentials flow (2-legged authentication)
- Automatic token caching and refresh with 401 auto-retry
- In-memory response caching with configurable TTL
- Automatic retries with exponential backoff for transient errors
- Rate limiting (120 requests/minute, 1000 requests/day)
- Locale support (site, language, currency)
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

    ctx := context.Background()

    // Search for products
    results, err := client.KeywordSearch(ctx, &digikey.SearchRequest{
        Keywords:    "STM32F4",
        RecordCount: 10,
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

// Client with custom options
client := digikey.NewClient(
    clientID,
    clientSecret,
    digikey.WithLocale(digikey.Locale{
        Site:     "DE",
        Language: "de",
        Currency: "EUR",
    }),
    digikey.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)
```

### Keyword Search

```go
// Basic search
results, err := client.KeywordSearch(ctx, &digikey.SearchRequest{
    Keywords:    "STM32F4",
    RecordCount: 10,
})

// Using the fluent builder
results, err := digikey.NewSearch("STM32F4").
    Limit(20).
    Offset(0).
    FilterByManufacturer(1).
    WithSearchOptions("InStock").
    Execute(ctx, client)
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

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/products/v4/search/keyword` | POST | Keyword search |
| `/products/v4/search/{productNumber}/productdetails` | GET | Product details |

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
| `WithRateLimiter` | Custom rate limiter |
| `WithTokenURL` | Custom OAuth token URL |
| `WithCache` | Custom cache implementation |
| `WithCacheConfig` | Configure cache TTLs |
| `WithoutCache` | Disable caching |
| `WithRetryConfig` | Custom retry configuration |
| `WithoutRetry` | Disable retries |

## Caching

The client includes in-memory caching with configurable TTL:

```go
// Default: caching enabled with 5min search TTL, 10min details TTL
client := digikey.NewClient(clientID, clientSecret)

// Custom cache configuration
client := digikey.NewClient(
    clientID,
    clientSecret,
    digikey.WithCacheConfig(digikey.CacheConfig{
        Enabled:    true,
        SearchTTL:  2 * time.Minute,
        DetailsTTL: 5 * time.Minute,
    }),
)

// Disable caching
client := digikey.NewClient(clientID, clientSecret, digikey.WithoutCache())

// Force refresh (bypass cache)
details, err := client.ProductDetailsNoCache(ctx, "497-15360-ND")

// Clear all cached data
client.ClearCache()
```

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

## License

MIT License - see [LICENSE](LICENSE) for details.

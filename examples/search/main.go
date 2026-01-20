// Command search demonstrates the Digi-Key API client.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/PatrickWalther/go-digikey"
)

func main() {
	var (
		keywords = flag.String("q", "", "Search keywords (required)")
		limit    = flag.Int("limit", 10, "Maximum results to return (1-50)")
		site     = flag.String("site", "US", "Locale site (US, DE, etc.)")
		lang     = flag.String("lang", "en", "Locale language (en, de, etc.)")
		currency = flag.String("currency", "USD", "Locale currency (USD, EUR, etc.)")
	)
	flag.Parse()

	if *keywords == "" {
		fmt.Fprintln(os.Stderr, "Error: search keywords required (-q)")
		flag.Usage()
		os.Exit(1)
	}

	clientID := os.Getenv("DIGIKEY_CLIENT_ID")
	clientSecret := os.Getenv("DIGIKEY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Fprintln(os.Stderr, "Error: DIGIKEY_CLIENT_ID and DIGIKEY_CLIENT_SECRET environment variables required")
		os.Exit(1)
	}

	client := digikey.NewClient(
		clientID,
		clientSecret,
		digikey.WithLocale(digikey.Locale{
			Site:     *site,
			Language: *lang,
			Currency: *currency,
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("Searching for: %s\n\n", *keywords)

	results, err := digikey.NewSearch(*keywords).
		Limit(*limit).
		Execute(ctx, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d products:\n\n", results.ProductsCount)

	for i, product := range results.Products {
		fmt.Printf("%d. %s\n", i+1, product.ManufacturerProductNumber)
		fmt.Printf("   Manufacturer: %s\n", product.Manufacturer.Name)
		fmt.Printf("   Description:  %s\n", product.Description.ProductDescription)
		fmt.Printf("   Price:        $%.4f\n", product.UnitPrice)
		fmt.Printf("   Stock:        %d\n", product.QuantityAvailable)
		if product.DatasheetURL != "" {
			fmt.Printf("   Datasheet:    %s\n", product.DatasheetURL)
		}
		fmt.Println()
	}

	stats := client.RateLimitStats()
	fmt.Printf("Rate Limit Status:\n")
	fmt.Printf("  Minute: %d/%d remaining\n", stats.MinuteRemaining, stats.MinuteLimit)
	fmt.Printf("  Day:    %d/%d remaining\n", stats.DayRemaining, stats.DayLimit)
}

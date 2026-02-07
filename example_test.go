package digikey_test

import (
	"context"
	"fmt"
	"os"

	"github.com/PatrickWalther/go-digikey"
)

func ExampleNewClient() {
	client := digikey.NewClient(
		os.Getenv("DIGIKEY_CLIENT_ID"),
		os.Getenv("DIGIKEY_CLIENT_SECRET"),
	)
	defer client.Close()

	_ = client // use client for API calls
}

func ExampleNewSearch() {
	client := digikey.NewClient(
		os.Getenv("DIGIKEY_CLIENT_ID"),
		os.Getenv("DIGIKEY_CLIENT_SECRET"),
	)
	defer client.Close()

	results, err := digikey.NewSearch("STM32F4").
		Limit(10).
		Execute(context.Background(), client)
	if err != nil {
		fmt.Println("search error:", err)
		return
	}

	for _, p := range results.Products {
		fmt.Println(p.ManufacturerProductNumber)
	}
}

func ExampleClient_SetLocale() {
	client := digikey.NewClient(
		os.Getenv("DIGIKEY_CLIENT_ID"),
		os.Getenv("DIGIKEY_CLIENT_SECRET"),
	)
	defer client.Close()

	client.SetLocale(digikey.Locale{
		Site:     "DE",
		Language: "de",
		Currency: "EUR",
	})

	_ = client // subsequent API calls will use the DE locale
}

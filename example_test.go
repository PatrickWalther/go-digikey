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

func ExampleCategoryService_List() {
	client := digikey.NewClient(
		os.Getenv("DIGIKEY_CLIENT_ID"),
		os.Getenv("DIGIKEY_CLIENT_SECRET"),
	)
	defer client.Close()

	resp, err := client.Category.List(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, cat := range resp.Categories {
		fmt.Printf("%s (%d products)\n", cat.Name, cat.ProductCount)
	}
}

func ExampleCategoryService_Manufacturers() {
	client := digikey.NewClient(
		os.Getenv("DIGIKEY_CLIENT_ID"),
		os.Getenv("DIGIKEY_CLIENT_SECRET"),
	)
	defer client.Close()

	resp, err := client.Category.Manufacturers(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, mfg := range resp.Manufacturers {
		fmt.Printf("%s (ID: %d)\n", mfg.Name, mfg.ID)
	}
}

func ExampleProductService_Associations() {
	client := digikey.NewClient(
		os.Getenv("DIGIKEY_CLIENT_ID"),
		os.Getenv("DIGIKEY_CLIENT_SECRET"),
	)
	defer client.Close()

	resp, err := client.Product.Associations(context.Background(), "497-15360-ND")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, mate := range resp.ProductAssociations.MatingProducts {
		fmt.Printf("Mating: %s\n", mate.ManufacturerProductNumber)
	}
}

func ExamplePricingService_DigiReel() {
	client := digikey.NewClient(
		os.Getenv("DIGIKEY_CLIENT_ID"),
		os.Getenv("DIGIKEY_CLIENT_SECRET"),
	)
	defer client.Close()

	resp, err := client.Pricing.DigiReel(context.Background(), "497-15360-ND", 1000)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Reeling fee: $%.2f, Unit price: $%.4f\n", resp.ReelingFee, resp.UnitPrice)
}

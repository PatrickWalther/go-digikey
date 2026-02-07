package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// --- DigiReel Pricing tests ---

func TestDigiReelPricing(t *testing.T) {
	resp := DigiReelPricingResponse{
		ReelingFee:        7.00,
		UnitPrice:         0.10,
		ExtendedPrice:     107.00,
		RequestedQuantity: 1000,
		SearchLocaleUsed:  SearchLocale{Site: "US"},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/ABC-123/digireelpricing" {
			t.Errorf("expected path /products/v4/search/ABC-123/digireelpricing, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("requestedQuantity") != "1000" {
			t.Errorf("expected requestedQuantity=1000, got %s", r.URL.Query().Get("requestedQuantity"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Pricing.DigiReel(ctx, "ABC-123", 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ReelingFee != 7.00 {
		t.Errorf("expected reeling fee 7.00, got %f", result.ReelingFee)
	}
	if result.RequestedQuantity != 1000 {
		t.Errorf("expected requested quantity 1000, got %d", result.RequestedQuantity)
	}
}

func TestDigiReelPricingEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Pricing.DigiReel(ctx, "", 100)
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestDigiReelPricingInvalidQuantity(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Pricing.DigiReel(ctx, "ABC-123", 0)
	if err == nil {
		t.Error("expected error for zero quantity")
	}

	_, err = client.Pricing.DigiReel(ctx, "ABC-123", -1)
	if err == nil {
		t.Error("expected error for negative quantity")
	}
}

func TestDigiReelPricingNotCached(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ReelingFee":7,"UnitPrice":0.1,"ExtendedPrice":107,"RequestedQuantity":1000}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Pricing.DigiReel(ctx, "TEST-1", 1000)
	_, _ = client.Pricing.DigiReel(ctx, "TEST-1", 1000)

	if callCount != 2 {
		t.Errorf("expected 2 API calls (not cached), got %d", callCount)
	}
}

// --- Package Type by Quantity tests ---

func TestPackageTypeByQuantity(t *testing.T) {
	resp := PackageTypeByQuantityResponse{
		Products: []PackageTypeByQuantityProduct{
			{
				DigiKeyProductNumber:      "123-ND",
				ManufacturerProductNumber: "ABC",
				QuantityAvailable:         500,
				PackageTypes:              []string{"Cut Tape", "Digi-Reel"},
			},
		},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/packagetypebyquantity/ABC-123" {
			t.Errorf("expected path /products/v4/search/packagetypebyquantity/ABC-123, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("requestedQuantity") != "500" {
			t.Errorf("expected requestedQuantity=500, got %s", r.URL.Query().Get("requestedQuantity"))
		}
		if r.URL.Query().Get("packagingPreference") != "CT" {
			t.Errorf("expected packagingPreference=CT, got %s", r.URL.Query().Get("packagingPreference"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Pricing.PackageTypeByQuantity(ctx, "ABC-123", 500, "CT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Products) != 1 {
		t.Errorf("expected 1 product, got %d", len(result.Products))
	}
	if len(result.Products[0].PackageTypes) != 2 {
		t.Errorf("expected 2 package types, got %d", len(result.Products[0].PackageTypes))
	}
}

func TestPackageTypeByQuantityNoPackagingPreference(t *testing.T) {
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("packagingPreference") != "" {
			t.Errorf("expected no packagingPreference param, got '%s'", r.URL.Query().Get("packagingPreference"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Products":[]}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, err := client.Pricing.PackageTypeByQuantity(ctx, "ABC-123", 100, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPackageTypeByQuantityEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Pricing.PackageTypeByQuantity(ctx, "", 100, "CT")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestPackageTypeByQuantityInvalidQuantity(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Pricing.PackageTypeByQuantity(ctx, "ABC-123", 0, "CT")
	if err == nil {
		t.Error("expected error for zero quantity")
	}

	_, err = client.Pricing.PackageTypeByQuantity(ctx, "ABC-123", -1, "")
	if err == nil {
		t.Error("expected error for negative quantity")
	}
}

func TestPackageTypeByQuantityNotCached(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Products":[]}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Pricing.PackageTypeByQuantity(ctx, "TEST-1", 100, "")
	_, _ = client.Pricing.PackageTypeByQuantity(ctx, "TEST-1", 100, "")

	if callCount != 2 {
		t.Errorf("expected 2 API calls (not cached), got %d", callCount)
	}
}

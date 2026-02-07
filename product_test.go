package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// --- Product Details tests ---

func TestProductDetailsEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.Product.Details(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestProductDetailsNoCacheEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.Product.DetailsNoCache(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestProductDetailsWithCache(t *testing.T) {
	// Verify cache key function works
	cacheKey := cacheKeyForDetails(DefaultLocale(), "TEST-123")
	if cacheKey == "" {
		t.Error("expected non-empty cache key")
	}
	if cacheKey != cacheKeyForDetails(DefaultLocale(), "TEST-123") {
		t.Error("cache key should be consistent")
	}
}

// --- Associations tests ---

func TestAssociations(t *testing.T) {
	resp := ProductAssociationsResponse{
		ProductAssociations: ProductAssociations{
			MatingProducts: []ProductSummary{
				{ManufacturerProductNumber: "MATE-001", UnitPrice: "1.50"},
			},
		},
		SearchLocaleUsed: SearchLocale{Site: "US"},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/ABC-123/associations" {
			t.Errorf("expected path /products/v4/search/ABC-123/associations, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Product.Associations(ctx, "ABC-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.ProductAssociations.MatingProducts) != 1 {
		t.Errorf("expected 1 mating product, got %d", len(result.ProductAssociations.MatingProducts))
	}
}

func TestAssociationsEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Product.Associations(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestAssociationsPathEscaping(t *testing.T) {
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Product number with special chars should be escaped in the raw path
		expectedRaw := "/products/v4/search/ABC%2F123/associations"
		if r.URL.RawPath != expectedRaw {
			t.Errorf("expected raw path %s, got %s", expectedRaw, r.URL.RawPath)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ProductAssociations":{}}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Product.Associations(ctx, "ABC/123")
}

func TestAssociationsCaching(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ProductAssociations":{}}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Product.Associations(ctx, "TEST-1")
	_, _ = client.Product.Associations(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

// --- Media tests ---

func TestMedia(t *testing.T) {
	resp := MediaResponse{
		MediaLinks: []MediaLink{
			{MediaType: "Photo", Title: "Front View", URL: "https://example.com/photo.jpg"},
			{MediaType: "Datasheet", Title: "Datasheet", URL: "https://example.com/ds.pdf"},
		},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/ABC-123/media" {
			t.Errorf("expected path /products/v4/search/ABC-123/media, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Product.Media(ctx, "ABC-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.MediaLinks) != 2 {
		t.Errorf("expected 2 media links, got %d", len(result.MediaLinks))
	}
	if result.MediaLinks[0].MediaType != "Photo" {
		t.Errorf("expected media type 'Photo', got '%s'", result.MediaLinks[0].MediaType)
	}
}

func TestMediaEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Product.Media(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestMediaCaching(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"MediaLinks":[]}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Product.Media(ctx, "TEST-1")
	_, _ = client.Product.Media(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

// --- Substitutions tests ---

func TestSubstitutions(t *testing.T) {
	resp := ProductSubstitutesResponse{
		ProductSubstitutesCount: 1,
		ProductSubstitutes: []ProductSubstitute{
			{
				SubstituteType:            "Direct",
				ManufacturerProductNumber: "ALT-001",
				UnitPrice:                 "2.50",
				QuantityAvailable:         300,
			},
		},
		SearchLocaleUsed: SearchLocale{Site: "US"},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/ABC-123/substitutions" {
			t.Errorf("expected path /products/v4/search/ABC-123/substitutions, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Product.Substitutions(ctx, "ABC-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ProductSubstitutesCount != 1 {
		t.Errorf("expected count 1, got %d", result.ProductSubstitutesCount)
	}
	if result.ProductSubstitutes[0].SubstituteType != "Direct" {
		t.Error("substitute type mismatch")
	}
}

func TestSubstitutionsEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Product.Substitutions(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestSubstitutionsCaching(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ProductSubstitutesCount":0,"ProductSubstitutes":[]}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Product.Substitutions(ctx, "TEST-1")
	_, _ = client.Product.Substitutions(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

// --- Recommendations tests ---

func TestRecommendedProducts(t *testing.T) {
	resp := RecommendedProductsResponse{
		Recommendations: []Recommendation{
			{
				ProductNumber: "ABC-123",
				RecommendedProducts: []RecommendedProduct{
					{DigiKeyProductNumber: "REC-001", UnitPrice: 1.25, QuantityAvailable: 200},
				},
				SearchLocaleUsed: SearchLocale{Site: "US"},
			},
		},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/ABC-123/recommendedproducts" {
			t.Errorf("expected path /products/v4/search/ABC-123/recommendedproducts, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Product.RecommendedProducts(ctx, "ABC-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Recommendations) != 1 {
		t.Errorf("expected 1 recommendation, got %d", len(result.Recommendations))
	}
	if result.Recommendations[0].ProductNumber != "ABC-123" {
		t.Error("product number mismatch")
	}
}

func TestRecommendedProductsEmptyProductNumber(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Product.RecommendedProducts(ctx, "")
	if err == nil {
		t.Error("expected error for empty product number")
	}
}

func TestRecommendedProductsCaching(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Recommendations":[]}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Product.RecommendedProducts(ctx, "TEST-1")
	_, _ = client.Product.RecommendedProducts(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

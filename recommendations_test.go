package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

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

	result, err := client.RecommendedProducts(ctx, "ABC-123")
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

	_, err := client.RecommendedProducts(ctx, "")
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

	_, _ = client.RecommendedProducts(ctx, "TEST-1")
	_, _ = client.RecommendedProducts(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

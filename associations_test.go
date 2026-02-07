package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

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

	result, err := client.Associations(ctx, "ABC-123")
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

	_, err := client.Associations(ctx, "")
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

	_, _ = client.Associations(ctx, "ABC/123")
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

	_, _ = client.Associations(ctx, "TEST-1")
	_, _ = client.Associations(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

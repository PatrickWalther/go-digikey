package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

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

	result, err := client.Substitutions(ctx, "ABC-123")
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

	_, err := client.Substitutions(ctx, "")
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

	_, _ = client.Substitutions(ctx, "TEST-1")
	_, _ = client.Substitutions(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

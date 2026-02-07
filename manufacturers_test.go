package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestManufacturers(t *testing.T) {
	resp := ManufacturersResponse{
		Manufacturers: []Manufacturer{
			{ID: 1, Name: "Texas Instruments"},
			{ID: 2, Name: "STMicroelectronics"},
		},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/manufacturers" {
			t.Errorf("expected path /products/v4/search/manufacturers, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Manufacturers(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Manufacturers) != 2 {
		t.Errorf("expected 2 manufacturers, got %d", len(result.Manufacturers))
	}
	if result.Manufacturers[0].Name != "Texas Instruments" {
		t.Errorf("expected 'Texas Instruments', got '%s'", result.Manufacturers[0].Name)
	}
}

func TestManufacturersCaching(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Manufacturers":[{"Id":1,"Name":"TI"}]}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, err := client.Manufacturers(ctx)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	_, err = client.Manufacturers(ctx)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

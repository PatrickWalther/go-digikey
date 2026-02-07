package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestCategories(t *testing.T) {
	resp := CategoriesResponse{
		ProductCount: 1000,
		Categories: []Category{
			{CategoryID: 1, Name: "Resistors"},
			{CategoryID: 2, Name: "Capacitors"},
		},
		SearchLocaleUsed: SearchLocale{Site: "US"},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/categories" {
			t.Errorf("expected path /products/v4/search/categories, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Category.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ProductCount != 1000 {
		t.Errorf("expected product count 1000, got %d", result.ProductCount)
	}
	if len(result.Categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(result.Categories))
	}
}

func TestCategoriesCaching(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ProductCount":1,"Categories":[]}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, err := client.Category.List(ctx)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	_, err = client.Category.List(ctx)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

func TestCategoriesById(t *testing.T) {
	resp := CategoryResponse{
		Category:         Category{CategoryID: 42, Name: "Op Amps"},
		SearchLocaleUsed: SearchLocale{Site: "US"},
	}
	respJSON, _ := json.Marshal(resp)

	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/products/v4/search/categories/42" {
			t.Errorf("expected path /products/v4/search/categories/42, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJSON)
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Category.GetByID(ctx, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Category.CategoryID != 42 {
		t.Errorf("expected category ID 42, got %d", result.Category.CategoryID)
	}
	if result.Category.Name != "Op Amps" {
		t.Errorf("expected category name 'Op Amps', got '%s'", result.Category.Name)
	}
}

func TestCategoriesByIdInvalidId(t *testing.T) {
	client := NewClient("test-id", "test-secret")
	defer client.Close()
	ctx := context.Background()

	_, err := client.Category.GetByID(ctx, 0)
	if err == nil {
		t.Error("expected error for categoryId 0")
	}

	_, err = client.Category.GetByID(ctx, -1)
	if err == nil {
		t.Error("expected error for negative categoryId")
	}
}

func TestCategoriesByIdCaching(t *testing.T) {
	callCount := 0
	server := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Category":{"CategoryId":1,"Name":"Test"}}`))
	})
	defer server.Close()

	client := newMockClient(t, server)
	ctx := context.Background()

	_, _ = client.Category.GetByID(ctx, 1)
	_, _ = client.Category.GetByID(ctx, 1)

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

// --- Manufacturers tests ---

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

	result, err := client.Category.Manufacturers(ctx)
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

	_, err := client.Category.Manufacturers(ctx)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	_, err = client.Category.Manufacturers(ctx)
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

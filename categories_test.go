package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newMockServer creates a test server that handles both OAuth token and API requests.
// The apiHandler is called for non-token API requests.
func newMockServer(t *testing.T, apiHandler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}
		apiHandler(w, r)
	}))
}

// newMockClient creates a test client configured to use the mock server.
func newMockClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	client := NewClient("test-id", "test-secret",
		WithBaseURL(server.URL),
		WithTokenURL(server.URL+"/v1/oauth2/token"),
		WithoutRetry(),
	)
	t.Cleanup(func() { client.Close() })
	return client
}

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

	result, err := client.Categories(ctx)
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

	_, err := client.Categories(ctx)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	_, err = client.Categories(ctx)
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

	result, err := client.CategoriesById(ctx, 42)
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

	_, err := client.CategoriesById(ctx, 0)
	if err == nil {
		t.Error("expected error for categoryId 0")
	}

	_, err = client.CategoriesById(ctx, -1)
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

	_, _ = client.CategoriesById(ctx, 1)
	_, _ = client.CategoriesById(ctx, 1)

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

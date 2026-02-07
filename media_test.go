package digikey

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

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

	result, err := client.Media(ctx, "ABC-123")
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

	_, err := client.Media(ctx, "")
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

	_, _ = client.Media(ctx, "TEST-1")
	_, _ = client.Media(ctx, "TEST-1")

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

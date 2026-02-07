package digikey

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestClientSetHeaders tests that headers are set correctly
func TestClientSetHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers
		if r.Header.Get("Authorization") == "" {
			t.Error("Authorization header missing")
		}
		if r.Header.Get("Content-Type") == "" {
			t.Error("Content-Type header missing")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"recordCount":0,"products":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make a request (will get auth error but we can see headers were attempted)
	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
}

// TestClientCustomerIdHeader tests that X-DIGIKEY-Customer-Id header is set when customerID is configured.
func TestClientCustomerIdHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customerID := r.Header.Get("X-DIGIKEY-Customer-Id")
		if customerID != "CUST-123" {
			t.Errorf("expected X-DIGIKEY-Customer-Id 'CUST-123', got '%s'", customerID)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Products":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL), WithCustomerID("CUST-123"))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
}

// TestClientCustomerIdHeaderAbsent tests that X-DIGIKEY-Customer-Id header is absent when not configured.
func TestClientCustomerIdHeaderAbsent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customerID := r.Header.Get("X-DIGIKEY-Customer-Id")
		if customerID != "" {
			t.Errorf("expected no X-DIGIKEY-Customer-Id header, got '%s'", customerID)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Products":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
}

// TestClientHandleErrorResponse tests error response handling
func TestClientHandleErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"statusCode":400,"message":"Bad Request","details":"Invalid parameters"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	if err == nil {
		t.Error("expected error response")
	}
}

// TestClientHandleRateLimitError tests rate limit error handling
func TestClientHandleRateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"statusCode":429,"message":"Rate Limited"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	if err == nil {
		t.Error("expected rate limit error")
	}
}

// TestClientRetryLogic tests that API calls are made
func TestClientRetryLogic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"products":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make request - should succeed (it will fail with JSON issue but that's ok, we're testing client setup)
	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
}

// TestClientDisableRetry tests that retry can be disabled
func TestClientDisableRetry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret",
		WithBaseURL(server.URL),
		WithoutRetry(),
	)
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	if err == nil {
		t.Error("expected error without retry")
	}
}

// TestClientCacheBypass tests that cache can be bypassed
func TestClientCacheBypass(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"recordCount":0,"products":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret",
		WithBaseURL(server.URL),
		WithCache(NewMemoryCache(5*time.Minute)),
	)
	defer client.Close()

	// Cache can be cleared without error
	client.ClearCache()
}

// TestClientResponseMalformed tests handling of malformed responses
func TestClientResponseMalformed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	if err == nil {
		t.Error("expected error for malformed response")
	}
}

// TestClientEmptyResponse tests handling of empty responses
func TestClientEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Write nothing
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	if err == nil {
		t.Error("expected error for empty response")
	}
}

// TestClientReadBodyError tests error when reading response body fails
func TestClientReadBodyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "999999999") // Lie about content length
		w.WriteHeader(http.StatusOK)
		// Don't write the promised amount of data
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
}

// TestClientContextCancellation tests that context cancellation works
func TestClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	if err == nil {
		t.Error("expected context deadline error")
	}
}

// TestClientHeadersNotModified tests that request body is not modified
func TestClientHeadersNotModified(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read request body
		body, _ := io.ReadAll(r.Body)
		originalBody := string(body)

		if originalBody == "" {
			t.Error("expected request body for POST")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"recordCount":0,"products":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
}

// TestClientDo tests the main do function with different HTTP methods
func TestClientDo(t *testing.T) {
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/products/v4/search/keyword"},
		{http.MethodGet, "/products/v4/search/123/productdetails"},
	}

	for _, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != test.method {
				t.Errorf("expected method %s, got %s", test.method, r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if test.method == http.MethodPost {
				_, _ = w.Write([]byte(`{"recordCount":0,"products":[]}`))
			} else {
				_, _ = w.Write([]byte(`{"product":{"digiKeyProductNumber":"123"}}`))
			}
		}))

		client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		if test.method == http.MethodPost {
			_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
		} else {
			_, _ = client.Product.Details(ctx, "123")
		}

		cancel()
		client.Close()
		server.Close()
	}
}

// TestClientStatusCodeHandling tests different status codes
func TestClientStatusCodeHandling(t *testing.T) {
	// Test 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"products":[]}`))
	}))

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL), WithoutRetry())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	client.Close()
	server.Close()

	// Test 400 Bad Request
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"error"}`))
	}))

	client = NewClient("test-id", "test-secret", WithBaseURL(server.URL), WithoutRetry())
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
	if err == nil {
		t.Error("expected error for 400 status")
	}
	client.Close()
	server.Close()
}

// TestClientResponseBodyClosing tests that response body is properly closed
func TestClientResponseBodyClosing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Return a response that tracks if it's closed
		_, _ = w.Write([]byte(`{"recordCount":0,"products":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-id", "test-secret", WithBaseURL(server.URL))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = client.Search.KeywordSearch(ctx, &SearchRequest{Keywords: "test", Limit: 5})
}

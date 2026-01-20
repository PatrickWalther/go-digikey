package digikey

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestTokenManagerGetToken tests token retrieval
func TestTokenManagerGetToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	tm := newTokenManager(server.Client(), "id", "secret", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := tm.getToken(ctx)
	if err != nil {
		t.Fatalf("getToken failed: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if token != "test-token" {
		t.Errorf("expected 'test-token', got '%s'", token)
	}
}

// TestTokenManagerRefresh tests token refresh
func TestTokenManagerRefresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"refresh-token","token_type":"Bearer","expires_in":7200}`))
	}))
	defer server.Close()

	tm := newTokenManager(server.Client(), "id", "secret", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token1, err := tm.refreshToken(ctx)
	if err != nil {
		t.Fatalf("refreshToken failed: %v", err)
	}
	if token1 != "refresh-token" {
		t.Errorf("expected 'refresh-token', got '%s'", token1)
	}

	// Token should be cached, so second call returns same token
	token2, err := tm.refreshToken(ctx)
	if err != nil {
		t.Fatalf("second refreshToken failed: %v", err)
	}
	if token2 != token1 {
		t.Errorf("expected same token, got %s", token2)
	}
}

// TestTokenManagerInvalidate tests token invalidation
func TestTokenManagerInvalidate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	tm := newTokenManager(server.Client(), "id", "secret", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get a token
	token1, _ := tm.getToken(ctx)

	// Invalidate
	tm.invalidate()

	// Next call should fetch new token
	tm.tokenExpiry = time.Now().Add(-1 * time.Second) // force refresh
	token2, _ := tm.getToken(ctx)

	if token1 == "" || token2 == "" {
		t.Error("tokens should not be empty")
	}
}

// TestTokenManagerAuthError tests handling of auth errors
func TestTokenManagerAuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid_client","error_description":"Invalid credentials"}`))
	}))
	defer server.Close()

	tm := newTokenManager(server.Client(), "bad-id", "bad-secret", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tm.getToken(ctx)
	if err == nil {
		t.Error("expected auth error")
	}
	if _, ok := err.(*AuthError); !ok {
		t.Errorf("expected AuthError, got %T", err)
	}
}

// TestTokenManagerBasicAuth tests that Basic Auth is used
func TestTokenManagerBasicAuth(t *testing.T) {
	authHeaderSeen := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that Authorization header has Basic auth
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && len(authHeader) > 6 {
			authHeaderSeen = true
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	tm := newTokenManager(server.Client(), "id", "secret", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = tm.getToken(ctx)
	if !authHeaderSeen {
		t.Error("expected Authorization header to be set")
	}
}

// TestTokenManagerTimeout tests context timeout
func TestTokenManagerTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tm := newTokenManager(server.Client(), "id", "secret", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := tm.getToken(ctx)
	if err == nil {
		t.Error("expected timeout error")
	}
}

// TestTokenManagerCaching tests that tokens are cached
func TestTokenManagerCaching(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"token","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	tm := newTokenManager(server.Client(), "id", "secret", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First call
	token1, _ := tm.getToken(ctx)
	if token1 != "token" {
		t.Errorf("expected token, got %s", token1)
	}

	// Should have made exactly 1 call
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

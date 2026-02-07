package digikey

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

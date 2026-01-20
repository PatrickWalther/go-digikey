package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultTokenURL = "https://api.digikey.com/v1/oauth2/token"
	tokenExpiryBuffer = 60 * time.Second
)

// tokenResponse represents the OAuth2 token response from Digi-Key.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// tokenManager handles OAuth2 token caching and refresh.
type tokenManager struct {
	mu           sync.RWMutex
	httpClient   *http.Client
	clientID     string
	clientSecret string
	tokenURL     string
	accessToken  string
	tokenExpiry  time.Time
}

func newTokenManager(httpClient *http.Client, clientID, clientSecret, tokenURL string) *tokenManager {
	if tokenURL == "" {
		tokenURL = defaultTokenURL
	}
	return &tokenManager{
		httpClient:   httpClient,
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     tokenURL,
	}
}

// getToken returns a valid access token, refreshing if necessary.
func (tm *tokenManager) getToken(ctx context.Context) (string, error) {
	tm.mu.RLock()
	token := tm.accessToken
	expiry := tm.tokenExpiry
	tm.mu.RUnlock()

	if token != "" && time.Now().Before(expiry.Add(-tokenExpiryBuffer)) {
		return token, nil
	}

	return tm.refreshToken(ctx)
}

// refreshToken obtains a new access token from the OAuth2 endpoint.
func (tm *tokenManager) refreshToken(ctx context.Context) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.accessToken != "" && time.Now().Before(tm.tokenExpiry.Add(-tokenExpiryBuffer)) {
		return tm.accessToken, nil
	}

	data := url.Values{
		"client_id":     {tm.clientID},
		"client_secret": {tm.clientSecret},
		"grant_type":    {"client_credentials"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tm.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("digikey: failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("digikey: token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("digikey: failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var authErr AuthError
		if err := json.Unmarshal(body, &authErr); err == nil && authErr.Err != "" {
			return "", &authErr
		}
		return "", &APIError{
			StatusCode: resp.StatusCode,
			Message:    "token request failed",
			Details:    string(body),
		}
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("digikey: failed to parse token response: %w", err)
	}

	tm.accessToken = tokenResp.AccessToken
	tm.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return tm.accessToken, nil
}

// invalidate clears the cached token.
func (tm *tokenManager) invalidate() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.accessToken = ""
	tm.tokenExpiry = time.Time{}
}

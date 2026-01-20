package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	searchBasePath = "/products/v4/search"
)

// KeywordSearch searches for products using keywords.
func (c *Client) KeywordSearch(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.Keywords == "" {
		return nil, fmt.Errorf("%w: keywords are required", ErrInvalidRequest)
	}

	// Create a copy to avoid mutating the caller's request
	searchReq := *req
	if searchReq.Limit <= 0 {
		searchReq.Limit = 10
	}
	if searchReq.Limit > 50 {
		searchReq.Limit = 50
	}

	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForSearch(c.getLocale(), &searchReq)
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp SearchResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	var resp SearchResponse
	err := c.do(ctx, http.MethodPost, searchBasePath+"/keyword", &searchReq, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForSearch(c.getLocale(), &searchReq)
			c.cache.Set(cacheKey, data, c.cacheConfig.SearchTTL)
		}
	}

	return &resp, nil
}

// ProductDetails retrieves detailed information about a specific product.
func (c *Client) ProductDetails(ctx context.Context, productNumber string) (*ProductDetailsResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForDetails(c.getLocale(), productNumber)
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp ProductDetailsResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	path := fmt.Sprintf("%s/%s/productdetails", searchBasePath, url.PathEscape(productNumber))

	var resp ProductDetailsResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForDetails(c.getLocale(), productNumber)
			c.cache.Set(cacheKey, data, c.cacheConfig.DetailsTTL)
		}
	}

	return &resp, nil
}

// ProductDetailsNoCache retrieves product details bypassing the cache.
// Use this for explicit pricing refresh operations.
func (c *Client) ProductDetailsNoCache(ctx context.Context, productNumber string) (*ProductDetailsResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	path := fmt.Sprintf("%s/%s/productdetails", searchBasePath, url.PathEscape(productNumber))

	var resp ProductDetailsResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	// Update cache with fresh data
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForDetails(c.getLocale(), productNumber)
			c.cache.Set(cacheKey, data, c.cacheConfig.DetailsTTL)
		}
	}

	return &resp, nil
}

// SearchOptions provides a builder pattern for constructing search requests.
type SearchOptions struct {
	request SearchRequest
}

// NewSearch creates a new search options builder.
func NewSearch(keywords string) *SearchOptions {
	return &SearchOptions{
		request: SearchRequest{
			Keywords: keywords,
			Limit:    10,
		},
	}
}

// Limit sets the maximum number of results to return (1-50).
func (s *SearchOptions) Limit(count int) *SearchOptions {
	if count < 1 {
		count = 1
	}
	if count > 50 {
		count = 50
	}
	s.request.Limit = count
	return s
}

// Offset sets the starting position for results.
func (s *SearchOptions) Offset(position int) *SearchOptions {
	if position < 0 {
		position = 0
	}
	s.request.Offset = position
	return s
}

// WithFilterOptions sets filter options.
func (s *SearchOptions) WithFilterOptions(filterRequest *FilterRequest) *SearchOptions {
	s.request.FilterOptionsRequest = filterRequest
	return s
}

// Build returns the constructed SearchRequest.
func (s *SearchOptions) Build() *SearchRequest {
	return &s.request
}

// Execute performs the search using the provided client.
func (s *SearchOptions) Execute(ctx context.Context, client *Client) (*SearchResponse, error) {
	return client.KeywordSearch(ctx, &s.request)
}

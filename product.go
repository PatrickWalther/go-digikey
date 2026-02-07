package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ProductService handles product-related operations.
type ProductService service

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
	err := c.do(ctx, "GET", path, nil, &resp)
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
	err := c.do(ctx, "GET", path, nil, &resp)
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

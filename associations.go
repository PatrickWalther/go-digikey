package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Associations retrieves product associations for a given product number.
func (c *Client) Associations(ctx context.Context, productNumber string) (*ProductAssociationsResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForLookup(c.getLocale(), "associations", productNumber)
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp ProductAssociationsResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	path := fmt.Sprintf("%s/%s/associations", searchBasePath, url.PathEscape(productNumber))

	var resp ProductAssociationsResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForLookup(c.getLocale(), "associations", productNumber)
			c.cache.Set(cacheKey, data, c.cacheConfig.LookupTTL)
		}
	}

	return &resp, nil
}

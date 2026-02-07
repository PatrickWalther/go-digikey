package digikey

import (
	"context"
	"encoding/json"
	"net/http"
)

// Manufacturers retrieves all manufacturers.
func (c *Client) Manufacturers(ctx context.Context) (*ManufacturersResponse, error) {
	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForLookup(c.getLocale(), "manufacturers", "all")
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp ManufacturersResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	var resp ManufacturersResponse
	err := c.do(ctx, http.MethodGet, searchBasePath+"/manufacturers", nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForLookup(c.getLocale(), "manufacturers", "all")
			c.cache.Set(cacheKey, data, c.cacheConfig.LookupTTL)
		}
	}

	return &resp, nil
}

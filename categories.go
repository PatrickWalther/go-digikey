package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CategoryService handles category and manufacturer operations.
type CategoryService service

// Categories retrieves all product categories.
func (c *Client) Categories(ctx context.Context) (*CategoriesResponse, error) {
	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForLookup(c.getLocale(), "categories", "all")
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp CategoriesResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	var resp CategoriesResponse
	err := c.do(ctx, http.MethodGet, searchBasePath+"/categories", nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForLookup(c.getLocale(), "categories", "all")
			c.cache.Set(cacheKey, data, c.cacheConfig.LookupTTL)
		}
	}

	return &resp, nil
}

// CategoriesById retrieves a specific category by its ID.
func (c *Client) CategoriesById(ctx context.Context, categoryId int) (*CategoryResponse, error) {
	if categoryId <= 0 {
		return nil, fmt.Errorf("%w: categoryId must be greater than 0", ErrInvalidRequest)
	}

	identifier := fmt.Sprintf("%d", categoryId)

	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForLookup(c.getLocale(), "category", identifier)
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp CategoryResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	path := fmt.Sprintf("%s/categories/%d", searchBasePath, categoryId)

	var resp CategoryResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForLookup(c.getLocale(), "category", identifier)
			c.cache.Set(cacheKey, data, c.cacheConfig.LookupTTL)
		}
	}

	return &resp, nil
}

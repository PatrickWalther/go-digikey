package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CategoryService handles category and manufacturer operations.
type CategoryService service

// --- Category response types ---

// CategoriesResponse represents a categories list response.
type CategoriesResponse struct {
	ProductCount     int          `json:"ProductCount"`
	Categories       []Category   `json:"Categories"`
	SearchLocaleUsed SearchLocale `json:"SearchLocaleUsed"`
}

// CategoryResponse represents a single category response.
type CategoryResponse struct {
	Category         Category     `json:"Category"`
	SearchLocaleUsed SearchLocale `json:"SearchLocaleUsed"`
}

// ManufacturersResponse represents a manufacturers list response.
type ManufacturersResponse struct {
	Manufacturers []Manufacturer `json:"Manufacturers"`
}

// --- Category methods ---

// List retrieves all product categories.
func (s *CategoryService) List(ctx context.Context) (*CategoriesResponse, error) {
	c := s.client

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

// GetByID retrieves a specific category by its ID.
func (s *CategoryService) GetByID(ctx context.Context, categoryId int) (*CategoryResponse, error) {
	if categoryId <= 0 {
		return nil, fmt.Errorf("%w: categoryId must be greater than 0", ErrInvalidRequest)
	}

	c := s.client
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

// Manufacturers retrieves all manufacturers.
func (s *CategoryService) Manufacturers(ctx context.Context) (*ManufacturersResponse, error) {
	c := s.client

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

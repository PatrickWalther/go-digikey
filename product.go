package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ProductService handles product-related operations.
type ProductService service

// --- Product response types ---

// ProductDetailsResponse represents a product details response.
type ProductDetailsResponse struct {
	Product          Product      `json:"Product"`
	SearchLocaleUsed SearchLocale `json:"SearchLocaleUsed"`
}

// ProductSummary represents a summarized product in an association.
type ProductSummary struct {
	ProductUrl                string       `json:"ProductUrl"`
	Description               string       `json:"Description"`
	Manufacturer              Manufacturer `json:"Manufacturer"`
	ManufacturerProductNumber string       `json:"ManufacturerProductNumber"`
	UnitPrice                 string       `json:"UnitPrice"`
	QuantityAvailable         int          `json:"QuantityAvailable"`
}

// ProductAssociations represents product association groups.
type ProductAssociations struct {
	Kits               []ProductSummary `json:"Kits"`
	MatingProducts     []ProductSummary `json:"MatingProducts"`
	AssociatedProducts []ProductSummary `json:"AssociatedProducts"`
	ForUseWithProducts []ProductSummary `json:"ForUseWithProducts"`
}

// ProductAssociationsResponse represents a product associations response.
type ProductAssociationsResponse struct {
	ProductAssociations ProductAssociations `json:"ProductAssociations"`
	SearchLocaleUsed    SearchLocale        `json:"SearchLocaleUsed"`
}

// MediaResponse represents a product media response.
type MediaResponse struct {
	MediaLinks []MediaLink `json:"MediaLinks"`
}

// ProductSubstitute represents a product substitute.
type ProductSubstitute struct {
	SubstituteType            string       `json:"SubstituteType"`
	ProductUrl                string       `json:"ProductUrl"`
	Description               string       `json:"Description"`
	Manufacturer              Manufacturer `json:"Manufacturer"`
	ManufacturerProductNumber string       `json:"ManufacturerProductNumber"`
	UnitPrice                 string       `json:"UnitPrice"`
	QuantityAvailable         int          `json:"QuantityAvailable"`
}

// ProductSubstitutesResponse represents a product substitutes response.
type ProductSubstitutesResponse struct {
	ProductSubstitutesCount int                 `json:"ProductSubstitutesCount"`
	ProductSubstitutes      []ProductSubstitute `json:"ProductSubstitutes"`
	SearchLocaleUsed        SearchLocale        `json:"SearchLocaleUsed"`
}

// RecommendedProduct represents a recommended product.
type RecommendedProduct struct {
	DigiKeyProductNumber      string  `json:"DigiKeyProductNumber"`
	ManufacturerProductNumber string  `json:"ManufacturerProductNumber"`
	ManufacturerName          string  `json:"ManufacturerName"`
	PrimaryPhoto              string  `json:"PrimaryPhoto"`
	ProductDescription        string  `json:"ProductDescription"`
	QuantityAvailable         int     `json:"QuantityAvailable"`
	UnitPrice                 float64 `json:"UnitPrice"`
	ProductUrl                string  `json:"ProductUrl"`
}

// Recommendation represents a recommendation group.
type Recommendation struct {
	ProductNumber       string               `json:"ProductNumber"`
	RecommendedProducts []RecommendedProduct `json:"RecommendedProducts"`
	SearchLocaleUsed    SearchLocale         `json:"SearchLocaleUsed"`
}

// RecommendedProductsResponse represents a recommended products response.
type RecommendedProductsResponse struct {
	Recommendations []Recommendation `json:"Recommendations"`
}

// --- Product methods ---

// Details retrieves detailed information about a specific product.
func (s *ProductService) Details(ctx context.Context, productNumber string) (*ProductDetailsResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	c := s.client

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

// DetailsNoCache retrieves product details bypassing the cache.
// Use this for explicit pricing refresh operations.
func (s *ProductService) DetailsNoCache(ctx context.Context, productNumber string) (*ProductDetailsResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	c := s.client
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

// Associations retrieves product associations for a given product number.
func (s *ProductService) Associations(ctx context.Context, productNumber string) (*ProductAssociationsResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	c := s.client

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

// Media retrieves media links for a given product number.
func (s *ProductService) Media(ctx context.Context, productNumber string) (*MediaResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	c := s.client

	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForLookup(c.getLocale(), "media", productNumber)
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp MediaResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	path := fmt.Sprintf("%s/%s/media", searchBasePath, url.PathEscape(productNumber))

	var resp MediaResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForLookup(c.getLocale(), "media", productNumber)
			c.cache.Set(cacheKey, data, c.cacheConfig.LookupTTL)
		}
	}

	return &resp, nil
}

// Substitutions retrieves product substitutes for a given product number.
func (s *ProductService) Substitutions(ctx context.Context, productNumber string) (*ProductSubstitutesResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	c := s.client

	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForLookup(c.getLocale(), "substitutions", productNumber)
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp ProductSubstitutesResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	path := fmt.Sprintf("%s/%s/substitutions", searchBasePath, url.PathEscape(productNumber))

	var resp ProductSubstitutesResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForLookup(c.getLocale(), "substitutions", productNumber)
			c.cache.Set(cacheKey, data, c.cacheConfig.LookupTTL)
		}
	}

	return &resp, nil
}

// RecommendedProducts retrieves recommended products for a given product number.
func (s *ProductService) RecommendedProducts(ctx context.Context, productNumber string) (*RecommendedProductsResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}

	c := s.client

	// Check cache
	if c.cacheConfig.Enabled && c.cache != nil {
		cacheKey := cacheKeyForLookup(c.getLocale(), "recommendations", productNumber)
		if cached, ok := c.cache.Get(cacheKey); ok {
			var resp RecommendedProductsResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	path := fmt.Sprintf("%s/%s/recommendedproducts", searchBasePath, url.PathEscape(productNumber))

	var resp RecommendedProductsResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if c.cacheConfig.Enabled && c.cache != nil {
		if data, err := json.Marshal(resp); err == nil {
			cacheKey := cacheKeyForLookup(c.getLocale(), "recommendations", productNumber)
			c.cache.Set(cacheKey, data, c.cacheConfig.LookupTTL)
		}
	}

	return &resp, nil
}

package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// --- Search request/response types ---

// SearchRequest represents a keyword search request.
type SearchRequest struct {
	Keywords             string         `json:"Keywords"`
	Limit                int            `json:"Limit,omitempty"`
	Offset               int            `json:"Offset,omitempty"`
	FilterOptionsRequest *FilterRequest `json:"FilterOptionsRequest,omitempty"`
	Includes             string         `json:"-"` // query param, not body
}

// FilterId represents a filter identifier used in API filter requests.
// The DigiKey API expects filter IDs as objects with a string Id field.
type FilterId struct {
	Id string `json:"Id"`
}

// NewFilterId creates a FilterId from an integer.
func NewFilterId(id int) FilterId {
	return FilterId{Id: fmt.Sprintf("%d", id)}
}

// NewFilterIds creates a slice of FilterId from integers.
func NewFilterIds(ids ...int) []FilterId {
	result := make([]FilterId, len(ids))
	for i, id := range ids {
		result[i] = NewFilterId(id)
	}
	return result
}

// FilterRequest represents a filter options request.
type FilterRequest struct {
	CategoryFilter           []FilterId              `json:"CategoryFilter,omitempty"`
	ManufacturerFilter       []FilterId              `json:"ManufacturerFilter,omitempty"`
	StatusFilter             []FilterId              `json:"StatusFilter,omitempty"`
	PackagingFilter          []FilterId              `json:"PackagingFilter,omitempty"`
	MarketPlaceFilter        string                  `json:"MarketPlaceFilter,omitempty"`
	SeriesFilter             []FilterId              `json:"SeriesFilter,omitempty"`
	MinimumQuantityAvailable int                     `json:"MinimumQuantityAvailable,omitempty"`
	SearchOptions            []string                `json:"SearchOptions,omitempty"`
	ParameterFilterRequest   *ParameterFilterRequest `json:"ParameterFilterRequest,omitempty"`
}

// ParameterFilterRequest represents a parameter filter request.
type ParameterFilterRequest struct {
	CategoryFilter   *FilterId          `json:"CategoryFilter,omitempty"`
	ParameterFilters []ParametricFilter `json:"ParameterFilters,omitempty"`
}

// ParametricFilter represents a parametric filter within a ParameterFilterRequest.
type ParametricFilter struct {
	ParameterID  int        `json:"ParameterId"`
	FilterValues []FilterId `json:"FilterValues,omitempty"`
}

// SearchResponse represents a keyword search response.
type SearchResponse struct {
	Products                 []Product       `json:"Products"`
	ProductsCount            int             `json:"ProductsCount"`
	ExactMatches             []Product       `json:"ExactMatches"`
	ExactMatchesCount        int             `json:"ExactMatchCount"`
	FilterOptions            FilterOptions   `json:"FilterOptions"`
	SearchLocaleUsed         SearchLocale    `json:"SearchLocaleUsed"`
	AppliedParametricFilters []AppliedFilter `json:"AppliedParametricFilters"`
}

// FilterOptions represents available filter options.
type FilterOptions struct {
	Categories        []CategoryFilter         `json:"Categories"`
	Manufacturers     []ManufacturerFilter     `json:"Manufacturers"`
	Status            []StatusFilter           `json:"Status"`
	PackageTypes      []PackageTypeFilter      `json:"PackageTypes"`
	ParametricFilters []ParametricFilterOption `json:"ParametricFilters"`
}

// CategoryFilter represents a category filter option.
type CategoryFilter struct {
	Category     Category `json:"Category"`
	ProductCount int      `json:"ProductCount"`
}

// ManufacturerFilter represents a manufacturer filter option.
type ManufacturerFilter struct {
	Manufacturer Manufacturer `json:"Manufacturer"`
	ProductCount int          `json:"ProductCount"`
}

// StatusFilter represents a status filter option.
type StatusFilter struct {
	StatusID     int    `json:"StatusId"`
	StatusName   string `json:"StatusName"`
	ProductCount int    `json:"ProductCount"`
}

// PackageTypeFilter represents a package type filter option.
type PackageTypeFilter struct {
	PackageType  PackageType `json:"PackageType"`
	ProductCount int         `json:"ProductCount"`
}

// ParametricFilterOption represents a parametric filter option.
type ParametricFilterOption struct {
	ParameterID   int           `json:"ParameterId"`
	ParameterName string        `json:"ParameterName"`
	Values        []FilterValue `json:"Values"`
}

// FilterValue represents a filter value option.
type FilterValue struct {
	ValueID      string `json:"ValueId"`
	ValueText    string `json:"ValueText"`
	ProductCount int    `json:"ProductCount"`
}

// AppliedFilter represents an applied parametric filter.
type AppliedFilter struct {
	ParameterID   int    `json:"ParameterId"`
	ParameterName string `json:"ParameterName"`
	ValueID       string `json:"ValueId"`
	ValueText     string `json:"ValueText"`
}

// SearchService handles keyword search operations.
type SearchService service

const (
	searchBasePath = "/products/v4/search"
)

// KeywordSearch searches for products using keywords.
func (s *SearchService) KeywordSearch(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.Keywords == "" {
		return nil, fmt.Errorf("%w: keywords are required", ErrInvalidRequest)
	}

	c := s.client

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

	path := searchBasePath + "/keyword"
	if searchReq.Includes != "" {
		path += "?includes=" + searchReq.Includes
	}

	var resp SearchResponse
	err := c.do(ctx, http.MethodPost, path, &searchReq, &resp)
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

// WithIncludes sets the includes query parameter for additional response data.
func (s *SearchOptions) WithIncludes(includes string) *SearchOptions {
	s.request.Includes = includes
	return s
}

// Build returns the constructed SearchRequest.
func (s *SearchOptions) Build() *SearchRequest {
	return &s.request
}

// Execute performs the search using the provided client.
func (s *SearchOptions) Execute(ctx context.Context, client *Client) (*SearchResponse, error) {
	return client.Search.KeywordSearch(ctx, &s.request)
}

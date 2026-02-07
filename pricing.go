package digikey

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PricingService handles pricing-related operations.
type PricingService service

// --- Pricing response types ---

// DigiReelPricingResponse represents a DigiReel pricing response.
type DigiReelPricingResponse struct {
	ReelingFee        float64      `json:"ReelingFee"`
	UnitPrice         float64      `json:"UnitPrice"`
	ExtendedPrice     float64      `json:"ExtendedPrice"`
	RequestedQuantity int          `json:"RequestedQuantity"`
	SearchLocaleUsed  SearchLocale `json:"SearchLocaleUsed"`
}

// PackageTypeByQuantityProduct represents a product in a package type by quantity response.
type PackageTypeByQuantityProduct struct {
	RecommendedQuantity           int          `json:"RecommendedQuantity"`
	DigiKeyProductNumber          string       `json:"DigiKeyProductNumber"`
	QuantityAvailable             int          `json:"QuantityAvailable"`
	ProductDescription            string       `json:"ProductDescription"`
	DetailedDescription           string       `json:"DetailedDescription"`
	ManufacturerName              string       `json:"ManufacturerName"`
	ManufacturerProductNumber     string       `json:"ManufacturerProductNumber"`
	MinimumOrderQuantity          int          `json:"MinimumOrderQuantity"`
	PrimaryDatasheetUrl           string       `json:"PrimaryDatasheetUrl"`
	PrimaryPhotoUrl               string       `json:"PrimaryPhotoUrl"`
	ProductStatus                 string       `json:"ProductStatus"`
	ManufacturerLeadWeeks         string       `json:"ManufacturerLeadWeeks"`
	ManufacturerWarehouseQuantity int          `json:"ManufacturerWarehouseQuantity"`
	RohsStatus                    string       `json:"RohsStatus"`
	RoHSCompliant                 bool         `json:"RoHSCompliant"`
	QuantityOnOrder               int          `json:"QuantityOnOrder"`
	StandardPricing               []PriceBreak `json:"StandardPricing"`
	MyPricing                     []PriceBreak `json:"MyPricing"`
	ProductUrl                    string       `json:"ProductUrl"`
	MarketPlace                   bool         `json:"MarketPlace"`
	Supplier                      string       `json:"Supplier"`
	StockNote                     string       `json:"StockNote"`
	PackageTypes                  []string     `json:"PackageTypes"`
}

// PackageTypeByQuantityResponse represents a package type by quantity response.
type PackageTypeByQuantityResponse struct {
	Products []PackageTypeByQuantityProduct `json:"Products"`
}

// --- Pricing methods ---

// DigiReel retrieves DigiReel pricing for a given product number and quantity.
// This endpoint is not cached because pricing is time-sensitive.
func (s *PricingService) DigiReel(ctx context.Context, productNumber string, requestedQuantity int) (*DigiReelPricingResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}
	if requestedQuantity <= 0 {
		return nil, fmt.Errorf("%w: requestedQuantity must be greater than 0", ErrInvalidRequest)
	}

	path := fmt.Sprintf("%s/%s/digireelpricing?requestedQuantity=%d",
		searchBasePath, url.PathEscape(productNumber), requestedQuantity)

	var resp DigiReelPricingResponse
	err := s.client.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// PackageTypeByQuantity retrieves package type options for a product at a given quantity.
// The packagingPreference parameter is optional; pass an empty string to omit it.
// This endpoint is not cached because pricing is time-sensitive.
func (s *PricingService) PackageTypeByQuantity(ctx context.Context, productNumber string, requestedQuantity int, packagingPreference string) (*PackageTypeByQuantityResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}
	if requestedQuantity <= 0 {
		return nil, fmt.Errorf("%w: requestedQuantity must be greater than 0", ErrInvalidRequest)
	}

	path := fmt.Sprintf("%s/packagetypebyquantity/%s?requestedQuantity=%d",
		searchBasePath, url.PathEscape(productNumber), requestedQuantity)

	if packagingPreference != "" {
		path += "&packagingPreference=" + url.QueryEscape(packagingPreference)
	}

	var resp PackageTypeByQuantityResponse
	err := s.client.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

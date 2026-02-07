package digikey

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PricingService handles pricing-related operations.
type PricingService service

// DigiReelPricing retrieves DigiReel pricing for a given product number and quantity.
// This endpoint is not cached because pricing is time-sensitive.
func (c *Client) DigiReelPricing(ctx context.Context, productNumber string, requestedQuantity int) (*DigiReelPricingResponse, error) {
	if productNumber == "" {
		return nil, fmt.Errorf("%w: product number is required", ErrInvalidRequest)
	}
	if requestedQuantity <= 0 {
		return nil, fmt.Errorf("%w: requestedQuantity must be greater than 0", ErrInvalidRequest)
	}

	path := fmt.Sprintf("%s/%s/digireelpricing?requestedQuantity=%d",
		searchBasePath, url.PathEscape(productNumber), requestedQuantity)

	var resp DigiReelPricingResponse
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

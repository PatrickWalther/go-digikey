package digikey

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PackageTypeByQuantity retrieves package type options for a product at a given quantity.
// The packagingPreference parameter is optional; pass an empty string to omit it.
// This endpoint is not cached because pricing is time-sensitive.
func (c *Client) PackageTypeByQuantity(ctx context.Context, productNumber string, requestedQuantity int, packagingPreference string) (*PackageTypeByQuantityResponse, error) {
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
	err := c.do(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

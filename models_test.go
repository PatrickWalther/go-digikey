package digikey

import (
	"encoding/json"
	"testing"
)

// TestProductStructure tests complete Product structure.
func TestProductStructure(t *testing.T) {
	product := &Product{
		ManufacturerProductNumber: "TL072CP",
		DigiKeyProductNumber:      "TL072CP-ND",
		UnitPrice:                 0.55,
		QuantityAvailable:         1000,
		ProductStatus: ProductStatus{
			Id:   0,
			Text: "Active",
		},
		Description: Description{
			ProductDescription: "Op Amp General Purpose",
		},
		Manufacturer: Manufacturer{
			ID:   123,
			Name: "Texas Instruments",
		},
	}

	if product.DigiKeyProductNumber != "TL072CP-ND" {
		t.Error("product number mismatch")
	}
	if product.Manufacturer.Name != "Texas Instruments" {
		t.Error("manufacturer mismatch")
	}
	if product.QuantityAvailable != 1000 {
		t.Error("quantity mismatch")
	}
}

// TestProductVariationStructure tests ProductVariation structure.
func TestProductVariationStructure(t *testing.T) {
	pv := ProductVariation{
		DigiKeyProductNumber: "TL072CP-ND",
		MinimumOrderQuantity: 1,
		QuantityAvailable:    5000,
		StandardPricing: []PriceBreak{
			{BreakQuantity: 1, UnitPrice: 0.55, TotalPrice: 0.55},
			{BreakQuantity: 10, UnitPrice: 0.50, TotalPrice: 5.00},
		},
	}

	if pv.MinimumOrderQuantity != 1 {
		t.Error("minimum order quantity mismatch")
	}
	if len(pv.StandardPricing) != 2 {
		t.Error("pricing breaks mismatch")
	}
}

// TestPriceBreak tests PriceBreak structure.
func TestPriceBreak(t *testing.T) {
	pb := PriceBreak{
		BreakQuantity: 100,
		UnitPrice:     9.99,
		TotalPrice:    999.00,
	}

	if pb.BreakQuantity != 100 {
		t.Errorf("expected break quantity 100, got %d", pb.BreakQuantity)
	}
	if pb.UnitPrice != 9.99 {
		t.Errorf("expected unit price 9.99, got %f", pb.UnitPrice)
	}
}

// TestParameter tests Parameter structure.
func TestParameter(t *testing.T) {
	param := Parameter{
		ParameterID:   1,
		ParameterText: "Operating Temperature",
		ValueID:       "100",
		ValueText:     "-40°C to +85°C",
	}

	if param.ParameterText != "Operating Temperature" {
		t.Errorf("expected parameter name Operating Temperature, got %s", param.ParameterText)
	}
	if param.ValueText != "-40°C to +85°C" {
		t.Errorf("expected parameter value -40°C to +85°C, got %s", param.ValueText)
	}
}

// TestCategory tests Category structure.
func TestCategory(t *testing.T) {
	cat := Category{
		CategoryID:   12,
		Name:         "Integrated Circuits",
		ProductCount: 15000,
	}

	if cat.CategoryID != 12 {
		t.Errorf("expected category ID 12, got %d", cat.CategoryID)
	}
	if cat.Name != "Integrated Circuits" {
		t.Errorf("expected category name Integrated Circuits, got %s", cat.Name)
	}
}

// TestManufacturer tests Manufacturer structure.
func TestManufacturer(t *testing.T) {
	mfg := Manufacturer{
		ID:   10,
		Name: "Texas Instruments",
	}

	if mfg.ID != 10 {
		t.Errorf("expected manufacturer ID 10, got %d", mfg.ID)
	}
	if mfg.Name != "Texas Instruments" {
		t.Errorf("expected manufacturer name Texas Instruments, got %s", mfg.Name)
	}
}

// TestLocale tests Locale structure.
func TestLocale(t *testing.T) {
	locale := Locale{
		Site:     "US",
		Language: "en",
		Currency: "USD",
	}

	if locale.Site != "US" {
		t.Errorf("expected site US, got %s", locale.Site)
	}
	if locale.Language != "en" {
		t.Errorf("expected language en, got %s", locale.Language)
	}
	if locale.Currency != "USD" {
		t.Errorf("expected currency USD, got %s", locale.Currency)
	}
}

// TestDefaultLocale tests DefaultLocale function.
func TestDefaultLocale(t *testing.T) {
	locale := DefaultLocale()

	if locale.Site != "US" {
		t.Errorf("expected default site US, got %s", locale.Site)
	}
	if locale.Language != "en" {
		t.Errorf("expected default language en, got %s", locale.Language)
	}
	if locale.Currency != "USD" {
		t.Errorf("expected default currency USD, got %s", locale.Currency)
	}
}

// TestSearchRequest tests SearchRequest structure.
func TestSearchRequest(t *testing.T) {
	req := SearchRequest{
		Keywords: "transistor",
		Limit:    20,
		Offset:   0,
	}

	if req.Keywords != "transistor" {
		t.Errorf("expected keywords transistor, got %s", req.Keywords)
	}
	if req.Limit != 20 {
		t.Errorf("expected limit 20, got %d", req.Limit)
	}
}

// TestSearchRequest JSON marshaling.
func TestSearchRequestJSON(t *testing.T) {
	req := SearchRequest{
		Keywords: "transistor",
		Limit:    10,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal search request: %v", err)
	}

	var decoded SearchRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal search request: %v", err)
	}

	if decoded.Keywords != "transistor" {
		t.Errorf("expected keywords transistor, got %s", decoded.Keywords)
	}
}

// TestSearchResponse tests SearchResponse structure.
func TestSearchResponse(t *testing.T) {
	resp := &SearchResponse{
		Products:      []Product{{DigiKeyProductNumber: "001"}, {DigiKeyProductNumber: "002"}},
		ProductsCount: 2,
	}

	if len(resp.Products) != 2 {
		t.Errorf("expected 2 products, got %d", len(resp.Products))
	}
	if resp.ProductsCount != 2 {
		t.Errorf("expected product count 2, got %d", resp.ProductsCount)
	}
}

// TestProductDetailsResponse tests ProductDetailsResponse structure.
func TestProductDetailsResponse(t *testing.T) {
	resp := &ProductDetailsResponse{
		Product: Product{
			DigiKeyProductNumber: "TL072CP-ND",
		},
		SearchLocaleUsed: SearchLocale{
			Site:     "US",
			Language: "en",
			Currency: "USD",
		},
	}

	if resp.Product.DigiKeyProductNumber != "TL072CP-ND" {
		t.Error("product number mismatch")
	}
	if resp.SearchLocaleUsed.Site != "US" {
		t.Error("locale mismatch")
	}
}

// TestFilters tests Filters structure.
func TestFilters(t *testing.T) {
	filters := Filters{
		CategoryIDs:     []int{1, 2, 3},
		ManufacturerIDs: []int{10, 20},
		StatusIDs:       []int{1},
		PackageTypeIDs:  []int{5, 6},
	}

	if len(filters.CategoryIDs) != 3 {
		t.Errorf("expected 3 categories, got %d", len(filters.CategoryIDs))
	}
	if len(filters.ManufacturerIDs) != 2 {
		t.Errorf("expected 2 manufacturers, got %d", len(filters.ManufacturerIDs))
	}
}

// TestSortOptions tests SortOptions structure.
func TestSortOptions(t *testing.T) {
	sort := SortOptions{
		Field:     "DateAdded",
		Direction: "Ascending",
	}

	if sort.Field != "DateAdded" {
		t.Errorf("expected sort field DateAdded, got %s", sort.Field)
	}
	if sort.Direction != "Ascending" {
		t.Errorf("expected sort direction Ascending, got %s", sort.Direction)
	}
}

// TestParametricFilter tests ParametricFilter structure.
func TestParametricFilter(t *testing.T) {
	pf := ParametricFilter{
		ParameterID: 100,
		ValueIDs:    []string{"val1", "val2"},
	}

	if pf.ParameterID != 100 {
		t.Errorf("expected parameter ID 100, got %d", pf.ParameterID)
	}
	if len(pf.ValueIDs) != 2 {
		t.Errorf("expected 2 value IDs, got %d", len(pf.ValueIDs))
	}
}

// TestSearchLocale tests SearchLocale structure.
func TestSearchLocale(t *testing.T) {
	locale := SearchLocale{
		Site:     "DE",
		Language: "de",
		Currency: "EUR",
	}

	if locale.Site != "DE" {
		t.Errorf("expected site DE, got %s", locale.Site)
	}
	if locale.Language != "de" {
		t.Errorf("expected language de, got %s", locale.Language)
	}
	if locale.Currency != "EUR" {
		t.Errorf("expected currency EUR, got %s", locale.Currency)
	}
}

// TestMediaLink tests MediaLink structure.
func TestMediaLink(t *testing.T) {
	link := MediaLink{
		MediaType: "Photo",
		Title:     "Component Photo",
		URL:       "https://example.com/photo.jpg",
	}

	if link.MediaType != "Photo" {
		t.Errorf("expected media type Photo, got %s", link.MediaType)
	}
	if link.URL != "https://example.com/photo.jpg" {
		t.Errorf("expected URL, got %s", link.URL)
	}
}

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

// TestNewFilterId tests NewFilterId constructor.
func TestNewFilterId(t *testing.T) {
	f := NewFilterId(42)
	if f.Id != "42" {
		t.Errorf("expected Id '42', got '%s'", f.Id)
	}

	f0 := NewFilterId(0)
	if f0.Id != "0" {
		t.Errorf("expected Id '0', got '%s'", f0.Id)
	}
}

// TestNewFilterIds tests NewFilterIds constructor.
func TestNewFilterIds(t *testing.T) {
	ids := NewFilterIds(1, 2, 3)
	if len(ids) != 3 {
		t.Fatalf("expected 3 filter IDs, got %d", len(ids))
	}
	if ids[0].Id != "1" || ids[1].Id != "2" || ids[2].Id != "3" {
		t.Errorf("unexpected IDs: %v", ids)
	}

	empty := NewFilterIds()
	if len(empty) != 0 {
		t.Errorf("expected 0 filter IDs, got %d", len(empty))
	}
}

// TestFilterIdJSON tests FilterId JSON round-trip.
func TestFilterIdJSON(t *testing.T) {
	f := NewFilterId(100)
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("failed to marshal FilterId: %v", err)
	}

	expected := `{"Id":"100"}`
	if string(data) != expected {
		t.Errorf("expected JSON %s, got %s", expected, string(data))
	}

	var decoded FilterId
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal FilterId: %v", err)
	}
	if decoded.Id != "100" {
		t.Errorf("expected Id '100', got '%s'", decoded.Id)
	}
}

// TestFilterRequestJSON tests FilterRequest JSON round-trip with all fields.
func TestFilterRequestJSON(t *testing.T) {
	req := FilterRequest{
		CategoryFilter:           NewFilterIds(1, 2),
		ManufacturerFilter:       NewFilterIds(10),
		StatusFilter:             NewFilterIds(0),
		PackagingFilter:          NewFilterIds(5, 6),
		MarketPlaceFilter:        "US",
		SeriesFilter:             NewFilterIds(99),
		MinimumQuantityAvailable: 100,
		SearchOptions:            []string{"InStock"},
		ParameterFilterRequest: &ParameterFilterRequest{
			CategoryFilter:   &FilterId{Id: "1"},
			ParameterFilters: []ParametricFilter{
				{ParameterID: 100, FilterValues: NewFilterIds(200, 201)},
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal FilterRequest: %v", err)
	}

	var decoded FilterRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal FilterRequest: %v", err)
	}

	if len(decoded.CategoryFilter) != 2 {
		t.Errorf("expected 2 category filters, got %d", len(decoded.CategoryFilter))
	}
	if decoded.MarketPlaceFilter != "US" {
		t.Errorf("expected MarketPlaceFilter 'US', got '%s'", decoded.MarketPlaceFilter)
	}
	if decoded.MinimumQuantityAvailable != 100 {
		t.Errorf("expected MinimumQuantityAvailable 100, got %d", decoded.MinimumQuantityAvailable)
	}
	if len(decoded.SearchOptions) != 1 || decoded.SearchOptions[0] != "InStock" {
		t.Errorf("unexpected SearchOptions: %v", decoded.SearchOptions)
	}
	if decoded.ParameterFilterRequest == nil {
		t.Fatal("expected non-nil ParameterFilterRequest")
	}
	if decoded.ParameterFilterRequest.CategoryFilter == nil || decoded.ParameterFilterRequest.CategoryFilter.Id != "1" {
		t.Error("unexpected CategoryFilter in ParameterFilterRequest")
	}
}

// TestFilterRequestOmitempty tests that empty fields are omitted in JSON.
func TestFilterRequestOmitempty(t *testing.T) {
	req := FilterRequest{}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal empty FilterRequest: %v", err)
	}
	if string(data) != "{}" {
		t.Errorf("expected empty JSON object, got %s", string(data))
	}
}

// TestParametricFilterStructure tests ParametricFilter structure.
func TestParametricFilterStructure(t *testing.T) {
	pf := ParametricFilter{
		ParameterID:  100,
		FilterValues: NewFilterIds(200, 201),
	}

	if pf.ParameterID != 100 {
		t.Errorf("expected parameter ID 100, got %d", pf.ParameterID)
	}
	if len(pf.FilterValues) != 2 {
		t.Errorf("expected 2 filter values, got %d", len(pf.FilterValues))
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

// TestProductAssociationsResponseJSON tests ProductAssociationsResponse JSON round-trip.
func TestProductAssociationsResponseJSON(t *testing.T) {
	resp := ProductAssociationsResponse{
		ProductAssociations: ProductAssociations{
			MatingProducts: []ProductSummary{
				{ManufacturerProductNumber: "ABC-123", UnitPrice: "1.50", QuantityAvailable: 100},
			},
		},
		SearchLocaleUsed: SearchLocale{Site: "US", Language: "en", Currency: "USD"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ProductAssociationsResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.ProductAssociations.MatingProducts) != 1 {
		t.Errorf("expected 1 mating product, got %d", len(decoded.ProductAssociations.MatingProducts))
	}
	if decoded.ProductAssociations.MatingProducts[0].UnitPrice != "1.50" {
		t.Error("unit price mismatch")
	}
}

// TestCategoriesResponseJSON tests CategoriesResponse JSON round-trip.
func TestCategoriesResponseJSON(t *testing.T) {
	resp := CategoriesResponse{
		ProductCount: 5000,
		Categories: []Category{
			{CategoryID: 1, Name: "Resistors", ProductCount: 1000},
		},
		SearchLocaleUsed: SearchLocale{Site: "US"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded CategoriesResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ProductCount != 5000 {
		t.Errorf("expected product count 5000, got %d", decoded.ProductCount)
	}
	if len(decoded.Categories) != 1 || decoded.Categories[0].Name != "Resistors" {
		t.Error("category mismatch")
	}
}

// TestCategoryResponseJSON tests CategoryResponse JSON round-trip.
func TestCategoryResponseJSON(t *testing.T) {
	resp := CategoryResponse{
		Category:         Category{CategoryID: 42, Name: "Capacitors"},
		SearchLocaleUsed: SearchLocale{Site: "DE"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded CategoryResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Category.CategoryID != 42 {
		t.Errorf("expected category ID 42, got %d", decoded.Category.CategoryID)
	}
}

// TestManufacturersResponseJSON tests ManufacturersResponse JSON round-trip.
func TestManufacturersResponseJSON(t *testing.T) {
	resp := ManufacturersResponse{
		Manufacturers: []Manufacturer{
			{ID: 1, Name: "Texas Instruments"},
			{ID: 2, Name: "STMicroelectronics"},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ManufacturersResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.Manufacturers) != 2 {
		t.Errorf("expected 2 manufacturers, got %d", len(decoded.Manufacturers))
	}
}

// TestMediaResponseJSON tests MediaResponse JSON round-trip.
func TestMediaResponseJSON(t *testing.T) {
	resp := MediaResponse{
		MediaLinks: []MediaLink{
			{MediaType: "Photo", Title: "Front", URL: "https://example.com/photo.jpg"},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded MediaResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.MediaLinks) != 1 || decoded.MediaLinks[0].MediaType != "Photo" {
		t.Error("media link mismatch")
	}
}

// TestDigiReelPricingResponseJSON tests DigiReelPricingResponse JSON round-trip.
func TestDigiReelPricingResponseJSON(t *testing.T) {
	resp := DigiReelPricingResponse{
		ReelingFee:        7.00,
		UnitPrice:         0.10,
		ExtendedPrice:     107.00,
		RequestedQuantity: 1000,
		SearchLocaleUsed:  SearchLocale{Site: "US", Currency: "USD"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded DigiReelPricingResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ReelingFee != 7.00 {
		t.Errorf("expected reeling fee 7.00, got %f", decoded.ReelingFee)
	}
	if decoded.RequestedQuantity != 1000 {
		t.Errorf("expected requested quantity 1000, got %d", decoded.RequestedQuantity)
	}
}

// TestPackageTypeByQuantityResponseJSON tests PackageTypeByQuantityResponse JSON round-trip.
func TestPackageTypeByQuantityResponseJSON(t *testing.T) {
	resp := PackageTypeByQuantityResponse{
		Products: []PackageTypeByQuantityProduct{
			{
				DigiKeyProductNumber:      "123-ND",
				ManufacturerProductNumber: "ABC",
				QuantityAvailable:         500,
				RoHSCompliant:             true,
				StandardPricing:           []PriceBreak{{BreakQuantity: 1, UnitPrice: 0.50}},
				PackageTypes:              []string{"Cut Tape", "Digi-Reel"},
			},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded PackageTypeByQuantityResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.Products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(decoded.Products))
	}
	p := decoded.Products[0]
	if !p.RoHSCompliant {
		t.Error("expected RoHSCompliant true")
	}
	if len(p.PackageTypes) != 2 {
		t.Errorf("expected 2 package types, got %d", len(p.PackageTypes))
	}
}

// TestRecommendedProductsResponseJSON tests RecommendedProductsResponse JSON round-trip.
func TestRecommendedProductsResponseJSON(t *testing.T) {
	resp := RecommendedProductsResponse{
		Recommendations: []Recommendation{
			{
				ProductNumber: "123-ND",
				RecommendedProducts: []RecommendedProduct{
					{DigiKeyProductNumber: "456-ND", UnitPrice: 1.25, QuantityAvailable: 200},
				},
				SearchLocaleUsed: SearchLocale{Site: "US"},
			},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded RecommendedProductsResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(decoded.Recommendations))
	}
	if decoded.Recommendations[0].ProductNumber != "123-ND" {
		t.Error("product number mismatch")
	}
}

// TestProductSubstitutesResponseJSON tests ProductSubstitutesResponse JSON round-trip.
func TestProductSubstitutesResponseJSON(t *testing.T) {
	resp := ProductSubstitutesResponse{
		ProductSubstitutesCount: 1,
		ProductSubstitutes: []ProductSubstitute{
			{
				SubstituteType:            "Direct",
				ManufacturerProductNumber: "ALT-123",
				UnitPrice:                 "2.50",
				QuantityAvailable:         300,
				Manufacturer:              Manufacturer{ID: 5, Name: "Murata"},
			},
		},
		SearchLocaleUsed: SearchLocale{Site: "US"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ProductSubstitutesResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ProductSubstitutesCount != 1 {
		t.Errorf("expected count 1, got %d", decoded.ProductSubstitutesCount)
	}
	if decoded.ProductSubstitutes[0].SubstituteType != "Direct" {
		t.Error("substitute type mismatch")
	}
}

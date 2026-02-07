package digikey

import "fmt"

// Product represents a Digi-Key product.
type Product struct {
	Description               Description        `json:"Description"`
	Manufacturer              Manufacturer       `json:"Manufacturer"`
	ManufacturerProductNumber string             `json:"ManufacturerProductNumber"`
	UnitPrice                 float64            `json:"UnitPrice"`
	ProductURL                string             `json:"ProductUrl"`
	DatasheetURL              string             `json:"DatasheetUrl"`
	PhotoURL                  string             `json:"PhotoUrl"`
	ProductVariations         []ProductVariation `json:"ProductVariations"`
	QuantityAvailable         int                `json:"QuantityAvailable"`
	Parameters                []Parameter        `json:"Parameters"`
	Category                  Category           `json:"Category"`
	DigiKeyProductNumber      string             `json:"DigiKeyProductNumber"`
	ProductStatus             ProductStatus      `json:"ProductStatus"`
	DateLastBuyChance         string             `json:"DateLastBuyChance"`
	AlternatePackaging        []AlternatePackage `json:"AlternatePackaging"`
	DetailedDescription       string             `json:"DetailedDescription"`
	TariffDescription         string             `json:"TariffDescription"`
	StandardPackage           int                `json:"StandardPackage"`
	LimitedTaxonomy           LimitedTaxonomy    `json:"LimitedTaxonomy"`
	Kits                      []Kit              `json:"Kits"`
	KitContents               []KitContent       `json:"KitContents"`
	MatingProducts            []MatingProduct    `json:"MatingProducts"`
	SearchLocaleUsed          SearchLocale       `json:"SearchLocaleUsed"`
	RohsInfo                  string             `json:"RohsInfo"`
	LeadStatus                string             `json:"LeadStatus"`
	ReachInfo                 string             `json:"ReachInfo"`
	ExportInformation         string             `json:"ExportInformation"`
	PrimaryPhoto              MediaLink          `json:"PrimaryPhoto"`
	MediaLinks                []MediaLink        `json:"MediaLinks"`
	Series                    Series             `json:"Series"`
	Classifications           Classifications    `json:"Classifications"`
}

// ProductStatus represents product status information.
type ProductStatus struct {
	Id   int    `json:"Id"`
	Text string `json:"Text"`
}

// ProductVariation represents a product packaging variation.
type ProductVariation struct {
	DigiKeyProductNumber   string       `json:"DigiKeyProductNumber"`
	PackageType            PackageType  `json:"PackageType"`
	StandardPricing        []PriceBreak `json:"StandardPricing"`
	QuantityAvailable      int          `json:"QuantityAvailableforPackageType"`
	MinimumOrderQuantity   int          `json:"MinimumOrderQuantity"`
	StandardPackage        int          `json:"StandardPackage"`
	DigiReelFee            float64      `json:"DigiReelFee"`
	MyPricing              []PriceBreak `json:"MyPricing"`
	MarketplaceRestriction bool         `json:"MarketplaceRestriction"`
}

// PriceBreak represents a quantity-based pricing tier.
type PriceBreak struct {
	BreakQuantity int     `json:"BreakQuantity"`
	UnitPrice     float64 `json:"UnitPrice"`
	TotalPrice    float64 `json:"TotalPrice"`
}

// Parameter represents a product parameter/specification.
type Parameter struct {
	ParameterID   int    `json:"ParameterId"`
	ParameterText string `json:"ParameterText"`
	ValueID       string `json:"ValueId"`
	ValueText     string `json:"ValueText"`
}

// Category represents a product category.
type Category struct {
	CategoryID      int        `json:"CategoryId"`
	ParentID        int        `json:"ParentId"`
	Name            string     `json:"Name"`
	ProductCount    int        `json:"ProductCount"`
	NewProductCount int        `json:"NewProductCount"`
	ImageURL        string     `json:"ImageUrl"`
	SeoDescription  string     `json:"SeoDescription"`
	ChildCategories []Category `json:"ChildCategories"`
}

// Manufacturer represents a manufacturer.
type Manufacturer struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// Description represents a product description.
type Description struct {
	ProductDescription  string `json:"ProductDescription"`
	DetailedDescription string `json:"DetailedDescription"`
}

// PackageType represents packaging information.
type PackageType struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// AlternatePackage represents an alternate packaging option.
type AlternatePackage struct {
	DigiKeyProductNumber string      `json:"DigiKeyProductNumber"`
	QuantityAvailable    int         `json:"QuantityAvailable"`
	UnitPrice            float64     `json:"UnitPrice"`
	PackageType          PackageType `json:"PackageType"`
}

// LimitedTaxonomy represents limited taxonomy information.
type LimitedTaxonomy struct {
	Children []LimitedTaxonomy `json:"Children"`
	Value    string            `json:"Value"`
	ID       int               `json:"Id"`
}

// Kit represents a kit product.
type Kit struct {
	DigiKeyProductNumber   string `json:"DigiKeyProductNumber"`
	ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
	QuantityInKit          int    `json:"QuantityInKit"`
}

// KitContent represents content of a kit.
type KitContent struct {
	DigiKeyProductNumber   string `json:"DigiKeyProductNumber"`
	ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
	QuantityInKit          int    `json:"QuantityInKit"`
}

// MatingProduct represents a mating/compatible product.
type MatingProduct struct {
	DigiKeyProductNumber   string `json:"DigiKeyProductNumber"`
	ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
}

// SearchLocale represents the locale used for a search.
type SearchLocale struct {
	Site     string `json:"Site"`
	Language string `json:"Language"`
	Currency string `json:"Currency"`
}

// MediaLink represents a media resource.
type MediaLink struct {
	MediaType  string `json:"MediaType"`
	Title      string `json:"Title"`
	SmallPhoto string `json:"SmallPhoto"`
	Thumbnail  string `json:"Thumbnail"`
	URL        string `json:"Url"`
}

// Series represents a product series.
type Series struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// Classifications represents product classifications.
type Classifications struct {
	ReachStatus              string `json:"ReachStatus"`
	RohsStatus               string `json:"RohsStatus"`
	MoistureSensitivityLevel string `json:"MoistureSensitivityLevel"`
	ExportControlClassNumber string `json:"ExportControlClassNumber"`
	HTSUSCode                string `json:"HtsusCode"`
}

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
	CategoryFilter         []FilterId               `json:"CategoryFilter,omitempty"`
	ManufacturerFilter     []FilterId               `json:"ManufacturerFilter,omitempty"`
	StatusFilter           []FilterId               `json:"StatusFilter,omitempty"`
	PackagingFilter        []FilterId               `json:"PackagingFilter,omitempty"`
	MarketPlaceFilter      string                   `json:"MarketPlaceFilter,omitempty"`
	SeriesFilter           []FilterId               `json:"SeriesFilter,omitempty"`
	MinimumQuantityAvailable int                    `json:"MinimumQuantityAvailable,omitempty"`
	SearchOptions          []string                 `json:"SearchOptions,omitempty"`
	ParameterFilterRequest *ParameterFilterRequest  `json:"ParameterFilterRequest,omitempty"`
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

// ProductDetailsResponse represents a product details response.
type ProductDetailsResponse struct {
	Product          Product      `json:"Product"`
	SearchLocaleUsed SearchLocale `json:"SearchLocaleUsed"`
}

// --- Associations ---

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

// --- Categories ---

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

// --- Manufacturers ---

// ManufacturersResponse represents a manufacturers list response.
type ManufacturersResponse struct {
	Manufacturers []Manufacturer `json:"Manufacturers"`
}

// --- Media ---

// MediaResponse represents a product media response.
type MediaResponse struct {
	MediaLinks []MediaLink `json:"MediaLinks"`
}

// --- DigiReel Pricing ---

// DigiReelPricingResponse represents a DigiReel pricing response.
type DigiReelPricingResponse struct {
	ReelingFee        float64      `json:"ReelingFee"`
	UnitPrice         float64      `json:"UnitPrice"`
	ExtendedPrice     float64      `json:"ExtendedPrice"`
	RequestedQuantity int          `json:"RequestedQuantity"`
	SearchLocaleUsed  SearchLocale `json:"SearchLocaleUsed"`
}

// --- Package Type by Quantity ---

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

// --- Recommended Products ---

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

// --- Substitutions ---

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

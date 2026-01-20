package digikey

// Locale specifies the locale for API requests.
type Locale struct {
	Site     string // US, DE, JP, etc.
	Language string // en, de, ja, etc.
	Currency string // USD, EUR, JPY, etc.
}

// DefaultLocale returns the default US locale.
func DefaultLocale() Locale {
	return Locale{
		Site:     "US",
		Language: "en",
		Currency: "USD",
	}
}

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
	ProductStatus             string             `json:"ProductStatus"`
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

// ProductVariation represents a product packaging variation.
type ProductVariation struct {
	DigiKeyProductNumber  string      `json:"DigiKeyProductNumber"`
	PackageType           PackageType `json:"PackageType"`
	StandardPricing       []PriceBreak `json:"StandardPricing"`
	QuantityAvailable     int         `json:"QuantityAvailableforPackageType"`
	MinimumOrderQuantity  int         `json:"MinimumOrderQuantity"`
	StandardPackage       int         `json:"StandardPackage"`
	DigiReelFee           float64     `json:"DigiReelFee"`
	MyPricing             []PriceBreak `json:"MyPricing"`
	MarketplaceRestriction bool       `json:"MarketplaceRestriction"`
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
	CategoryID   int    `json:"CategoryId"`
	ParentID     int    `json:"ParentId"`
	Name         string `json:"Name"`
	ProductCount int    `json:"ProductCount"`
	NewProductCount int `json:"NewProductCount"`
	ImageURL     string `json:"ImageUrl"`
	SeoDescription string `json:"SeoDescription"`
	ChildCategories []Category `json:"ChildCategories"`
}

// Manufacturer represents a manufacturer.
type Manufacturer struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// Description represents a product description.
type Description struct {
	ProductDescription string `json:"ProductDescription"`
	DetailedDescription string `json:"DetailedDescription"`
}

// PackageType represents packaging information.
type PackageType struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// AlternatePackage represents an alternate packaging option.
type AlternatePackage struct {
	DigiKeyProductNumber string  `json:"DigiKeyProductNumber"`
	QuantityAvailable    int     `json:"QuantityAvailable"`
	UnitPrice            float64 `json:"UnitPrice"`
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
	DigiKeyProductNumber string `json:"DigiKeyProductNumber"`
	ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
	QuantityInKit        int    `json:"QuantityInKit"`
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
	MediaType string `json:"MediaType"`
	Title     string `json:"Title"`
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
	ReachStatus            string `json:"ReachStatus"`
	RohsStatus             string `json:"RohsStatus"`
	MoistureSensitivityLevel string `json:"MoistureSensitivityLevel"`
	ExportControlClassNumber string `json:"ExportControlClassNumber"`
	HTSUSCode              string `json:"HtsusCode"`
}

// SearchRequest represents a keyword search request.
type SearchRequest struct {
	Keywords              string         `json:"Keywords"`
	RecordCount           int            `json:"RecordCount,omitempty"`
	RecordStartPosition   int            `json:"RecordStartPosition,omitempty"`
	Filters               *Filters       `json:"Filters,omitempty"`
	Sort                  *SortOptions   `json:"Sort,omitempty"`
	RequestedQuantity     int            `json:"RequestedQuantity,omitempty"`
	SearchOptions         []string       `json:"SearchOptions,omitempty"`
	FilterOptionsRequest  *FilterRequest `json:"FilterOptionsRequest,omitempty"`
}

// Filters represents search filters.
type Filters struct {
	CategoryIDs        []int             `json:"CategoryIds,omitempty"`
	ManufacturerIDs    []int             `json:"ManufacturerIds,omitempty"`
	StatusIDs          []int             `json:"StatusIds,omitempty"`
	PackageTypeIDs     []int             `json:"PackageTypeIds,omitempty"`
	ParametricFilters  []ParametricFilter `json:"ParametricFilters,omitempty"`
}

// ParametricFilter represents a parametric filter.
type ParametricFilter struct {
	ParameterID int      `json:"ParameterId"`
	ValueID     string   `json:"ValueId,omitempty"`
	ValueIDs    []string `json:"ValueIds,omitempty"`
}

// SortOptions represents sorting options.
type SortOptions struct {
	Field     string `json:"SortOption"`
	Direction string `json:"Direction"`
}

// FilterRequest represents a filter options request.
type FilterRequest struct {
	CategoryFilter      []int `json:"CategoryFilter,omitempty"`
	ManufacturerFilter  []int `json:"ManufacturerFilter,omitempty"`
	StatusFilter        []int `json:"StatusFilter,omitempty"`
	PackageTypeFilter   []int `json:"PackageTypeFilter,omitempty"`
	ParameterFilterRequest []ParameterFilterRequest `json:"ParameterFilterRequest,omitempty"`
}

// ParameterFilterRequest represents a parameter filter request.
type ParameterFilterRequest struct {
	ParameterID int      `json:"ParameterId"`
	ValueIDs    []string `json:"ValueIds,omitempty"`
}

// SearchResponse represents a keyword search response.
type SearchResponse struct {
	Products              []Product      `json:"Products"`
	ProductsCount         int            `json:"ProductsCount"`
	ExactMatches          []Product      `json:"ExactMatches"`
	ExactMatchesCount     int            `json:"ExactMatchCount"`
	FilterOptions         FilterOptions  `json:"FilterOptions"`
	SearchLocaleUsed      SearchLocale   `json:"SearchLocaleUsed"`
	AppliedParametricFilters []AppliedFilter `json:"AppliedParametricFilters"`
}

// FilterOptions represents available filter options.
type FilterOptions struct {
	Categories           []CategoryFilter    `json:"Categories"`
	Manufacturers        []ManufacturerFilter `json:"Manufacturers"`
	Status               []StatusFilter      `json:"Status"`
	PackageTypes         []PackageTypeFilter `json:"PackageTypes"`
	ParametricFilters    []ParametricFilterOption `json:"ParametricFilters"`
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

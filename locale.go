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

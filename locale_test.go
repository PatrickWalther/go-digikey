package digikey

import "testing"

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

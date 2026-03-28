package phone

import (
	"testing"

	"github.com/spf13/cobra"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestPhoneCmd_Metadata(t *testing.T) {
	if phoneCmd.Use != "phone" {
		t.Errorf("unexpected Use: %s", phoneCmd.Use)
	}
	if phoneCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestPhoneCmd_Subcommands(t *testing.T) {
	names := map[string]bool{}
	for _, sub := range phoneCmd.Commands() {
		names[sub.Name()] = true
	}
	for _, want := range []string{"format", "validate", "country"} {
		if !names[want] {
			t.Errorf("missing subcommand: %s", want)
		}
	}
}

// ── stripNonDigits ─────────────────────────────────────────────────────────

func TestStripNonDigits(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"+1 (202) 555-1234", "12025551234"},
		{"0312345678", "0312345678"},
		{"abc", ""},
		{"", ""},
		{"123", "123"},
	}
	for _, tc := range tests {
		got := stripNonDigits(tc.input)
		if got != tc.want {
			t.Errorf("stripNonDigits(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// ── lookupCountry ──────────────────────────────────────────────────────────

func TestLookupCountry_ByISO(t *testing.T) {
	info, err := lookupCountry("US")
	if err != nil {
		t.Fatal(err)
	}
	if info.DialCode != "1" {
		t.Errorf("US dial code: got %q, want 1", info.DialCode)
	}
	if info.Name != "United States" {
		t.Errorf("US name: got %q", info.Name)
	}
}

func TestLookupCountry_ByName(t *testing.T) {
	info, err := lookupCountry("Japan")
	if err != nil {
		t.Fatal(err)
	}
	if info.ISO != "JP" {
		t.Errorf("Japan ISO: got %q, want JP", info.ISO)
	}
	if info.DialCode != "81" {
		t.Errorf("Japan dial code: got %q, want 81", info.DialCode)
	}
}

func TestLookupCountry_CaseInsensitive(t *testing.T) {
	info, err := lookupCountry("united kingdom")
	if err != nil {
		t.Fatal(err)
	}
	if info.DialCode != "44" {
		t.Errorf("UK dial code: got %q, want 44", info.DialCode)
	}
}

func TestLookupCountry_NotFound(t *testing.T) {
	_, err := lookupCountry("Atlantis")
	if err == nil {
		t.Error("expected error for unknown country")
	}
}

// ── formatNumber ───────────────────────────────────────────────────────────

func TestFormatNumber_US(t *testing.T) {
	r, err := formatNumber("2025551234", "US")
	if err != nil {
		t.Fatal(err)
	}
	if r.E164 != "+12025551234" {
		t.Errorf("E.164: got %q, want +12025551234", r.E164)
	}
	if r.International != "+1 202-555-1234" {
		t.Errorf("International: got %q, want +1 202-555-1234", r.International)
	}
	if r.National != "(202) 555-1234" {
		t.Errorf("National: got %q, want (202) 555-1234", r.National)
	}
}

func TestFormatNumber_JP(t *testing.T) {
	r, err := formatNumber("03-1234-5678", "JP")
	if err != nil {
		t.Fatal(err)
	}
	// After normalization, leading 0 is stripped: 312345678 (9 digits).
	if r.E164 != "+81312345678" {
		t.Errorf("E.164: got %q", r.E164)
	}
}

func TestFormatNumber_WithCountryPrefix(t *testing.T) {
	r, err := formatNumber("+44 20 7946 0958", "GB")
	if err != nil {
		t.Fatal(err)
	}
	if r.E164 != "+442079460958" {
		t.Errorf("E.164: got %q, want +442079460958", r.E164)
	}
}

func TestFormatNumber_EmptyInput(t *testing.T) {
	_, err := formatNumber("", "US")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestFormatNumber_InvalidCountry(t *testing.T) {
	_, err := formatNumber("1234567890", "ZZ")
	if err == nil {
		t.Error("expected error for unknown country")
	}
}

// ── validateNumber ─────────────────────────────────────────────────────────

func TestValidateNumber_ValidUS(t *testing.T) {
	r := validateNumber("2025551234", "US")
	if !r.Valid {
		t.Errorf("expected valid, got invalid: %s", r.Reason)
	}
}

func TestValidateNumber_TooShort(t *testing.T) {
	r := validateNumber("12345", "US")
	if r.Valid {
		t.Error("expected invalid for too-short number")
	}
	if r.Reason == "" {
		t.Error("reason should not be empty for invalid number")
	}
}

func TestValidateNumber_TooLong(t *testing.T) {
	r := validateNumber("123456789012345", "US")
	if r.Valid {
		t.Error("expected invalid for too-long number")
	}
}

func TestValidateNumber_UnknownCountry(t *testing.T) {
	r := validateNumber("1234567890", "ZZ")
	if r.Valid {
		t.Error("expected invalid for unknown country")
	}
}

func TestValidateNumber_FR(t *testing.T) {
	r := validateNumber("612345678", "FR")
	if !r.Valid {
		t.Errorf("expected valid FR number, got: %s", r.Reason)
	}
}

// ── applyPattern ───────────────────────────────────────────────────────────

func TestApplyPattern(t *testing.T) {
	got := applyPattern("(XXX) XXX-XXXX", "2025551234")
	want := "(202) 555-1234"
	if got != want {
		t.Errorf("applyPattern: got %q, want %q", got, want)
	}
}

func TestApplyPattern_ExtraDigits(t *testing.T) {
	// If more digits than X placeholders, extras are appended.
	got := applyPattern("XX-XX", "123456")
	if got != "12-3456" {
		t.Errorf("applyPattern with extra digits: got %q, want 12-3456", got)
	}
}

// ── flags ──────────────────────────────────────────────────────────────────

func TestFormatCmd_HasCountryFlag(t *testing.T) {
	f := formatCmd.Flags().Lookup("country")
	if f == nil {
		t.Error("format should have --country flag")
	}
	if f.DefValue != "US" {
		t.Errorf("country default: got %q, want US", f.DefValue)
	}
}

func TestCountryCmd_AcceptsExactlyOneArg(t *testing.T) {
	validator := cobra.ExactArgs(1)
	if err := validator(countryCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(countryCmd, []string{"US"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFormatCmd_HasJSONFlag(t *testing.T) {
	f := formatCmd.Flags().Lookup("json")
	if f == nil {
		t.Error("format should have --json flag")
	}
	if f.Shorthand != "j" {
		t.Errorf("--json shorthand: got %q, want 'j'", f.Shorthand)
	}
}

// ── country data coverage ──────────────────────────────────────────────────

func TestCountryData_AtLeast30(t *testing.T) {
	if len(countries) < 30 {
		t.Errorf("expected at least 30 countries, got %d", len(countries))
	}
}

func TestCountryData_DialCodesNotEmpty(t *testing.T) {
	for iso, info := range countries {
		if info.DialCode == "" {
			t.Errorf("country %s has empty dial code", iso)
		}
		if info.Name == "" {
			t.Errorf("country %s has empty name", iso)
		}
	}
}

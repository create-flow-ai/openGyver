package phone

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var jsonOut bool

// ── country data ───────────────────────────────────────────────────────────

// CountryInfo holds phone metadata for a single country.
type CountryInfo struct {
	Name       string `json:"name"`
	ISO        string `json:"iso"`
	DialCode   string `json:"dialCode"`
	MinDigits  int    `json:"minDigits"`  // minimum digits in national number (excluding dial code)
	MaxDigits  int    `json:"maxDigits"`  // maximum digits in national number
	NatPattern string `json:"natPattern"` // national format pattern using X as digit placeholder
	IntPattern string `json:"intPattern"` // international format pattern
}

var countries = map[string]CountryInfo{
	"US": {Name: "United States", ISO: "US", DialCode: "1", MinDigits: 10, MaxDigits: 10, NatPattern: "(XXX) XXX-XXXX", IntPattern: "+1 XXX-XXX-XXXX"},
	"CA": {Name: "Canada", ISO: "CA", DialCode: "1", MinDigits: 10, MaxDigits: 10, NatPattern: "(XXX) XXX-XXXX", IntPattern: "+1 XXX-XXX-XXXX"},
	"GB": {Name: "United Kingdom", ISO: "GB", DialCode: "44", MinDigits: 10, MaxDigits: 10, NatPattern: "0XXXX XXXXXX", IntPattern: "+44 XXXX XXXXXX"},
	"UK": {Name: "United Kingdom", ISO: "GB", DialCode: "44", MinDigits: 10, MaxDigits: 10, NatPattern: "0XXXX XXXXXX", IntPattern: "+44 XXXX XXXXXX"},
	"DE": {Name: "Germany", ISO: "DE", DialCode: "49", MinDigits: 10, MaxDigits: 11, NatPattern: "0XXX XXXXXXX", IntPattern: "+49 XXX XXXXXXX"},
	"FR": {Name: "France", ISO: "FR", DialCode: "33", MinDigits: 9, MaxDigits: 9, NatPattern: "0X XX XX XX XX", IntPattern: "+33 X XX XX XX XX"},
	"JP": {Name: "Japan", ISO: "JP", DialCode: "81", MinDigits: 10, MaxDigits: 10, NatPattern: "0XX-XXXX-XXXX", IntPattern: "+81 XX-XXXX-XXXX"},
	"KR": {Name: "South Korea", ISO: "KR", DialCode: "82", MinDigits: 10, MaxDigits: 11, NatPattern: "0XX-XXXX-XXXX", IntPattern: "+82 XX-XXXX-XXXX"},
	"CN": {Name: "China", ISO: "CN", DialCode: "86", MinDigits: 11, MaxDigits: 11, NatPattern: "0XXX XXXX XXXX", IntPattern: "+86 XXX XXXX XXXX"},
	"IN": {Name: "India", ISO: "IN", DialCode: "91", MinDigits: 10, MaxDigits: 10, NatPattern: "XXXXX XXXXX", IntPattern: "+91 XXXXX XXXXX"},
	"AU": {Name: "Australia", ISO: "AU", DialCode: "61", MinDigits: 9, MaxDigits: 9, NatPattern: "0XXX XXX XXX", IntPattern: "+61 XXX XXX XXX"},
	"BR": {Name: "Brazil", ISO: "BR", DialCode: "55", MinDigits: 10, MaxDigits: 11, NatPattern: "(XX) XXXXX-XXXX", IntPattern: "+55 XX XXXXX-XXXX"},
	"MX": {Name: "Mexico", ISO: "MX", DialCode: "52", MinDigits: 10, MaxDigits: 10, NatPattern: "XX XXXX XXXX", IntPattern: "+52 XX XXXX XXXX"},
	"IT": {Name: "Italy", ISO: "IT", DialCode: "39", MinDigits: 9, MaxDigits: 11, NatPattern: "XXX XXX XXXX", IntPattern: "+39 XXX XXX XXXX"},
	"ES": {Name: "Spain", ISO: "ES", DialCode: "34", MinDigits: 9, MaxDigits: 9, NatPattern: "XXX XXX XXX", IntPattern: "+34 XXX XXX XXX"},
	"RU": {Name: "Russia", ISO: "RU", DialCode: "7", MinDigits: 10, MaxDigits: 10, NatPattern: "8 (XXX) XXX-XX-XX", IntPattern: "+7 XXX XXX-XX-XX"},
	"NL": {Name: "Netherlands", ISO: "NL", DialCode: "31", MinDigits: 9, MaxDigits: 9, NatPattern: "0XX XXX XXXX", IntPattern: "+31 XX XXX XXXX"},
	"SE": {Name: "Sweden", ISO: "SE", DialCode: "46", MinDigits: 9, MaxDigits: 9, NatPattern: "0XX-XXX XX XX", IntPattern: "+46 XX-XXX XX XX"},
	"NO": {Name: "Norway", ISO: "NO", DialCode: "47", MinDigits: 8, MaxDigits: 8, NatPattern: "XXX XX XXX", IntPattern: "+47 XXX XX XXX"},
	"DK": {Name: "Denmark", ISO: "DK", DialCode: "45", MinDigits: 8, MaxDigits: 8, NatPattern: "XX XX XX XX", IntPattern: "+45 XX XX XX XX"},
	"FI": {Name: "Finland", ISO: "FI", DialCode: "358", MinDigits: 9, MaxDigits: 10, NatPattern: "0XX XXX XXXX", IntPattern: "+358 XX XXX XXXX"},
	"PL": {Name: "Poland", ISO: "PL", DialCode: "48", MinDigits: 9, MaxDigits: 9, NatPattern: "XXX XXX XXX", IntPattern: "+48 XXX XXX XXX"},
	"CH": {Name: "Switzerland", ISO: "CH", DialCode: "41", MinDigits: 9, MaxDigits: 9, NatPattern: "0XX XXX XX XX", IntPattern: "+41 XX XXX XX XX"},
	"AT": {Name: "Austria", ISO: "AT", DialCode: "43", MinDigits: 10, MaxDigits: 11, NatPattern: "0XXX XXXXXXX", IntPattern: "+43 XXX XXXXXXX"},
	"NZ": {Name: "New Zealand", ISO: "NZ", DialCode: "64", MinDigits: 8, MaxDigits: 9, NatPattern: "0XX XXX XXXX", IntPattern: "+64 XX XXX XXXX"},
	"SG": {Name: "Singapore", ISO: "SG", DialCode: "65", MinDigits: 8, MaxDigits: 8, NatPattern: "XXXX XXXX", IntPattern: "+65 XXXX XXXX"},
	"HK": {Name: "Hong Kong", ISO: "HK", DialCode: "852", MinDigits: 8, MaxDigits: 8, NatPattern: "XXXX XXXX", IntPattern: "+852 XXXX XXXX"},
	"TW": {Name: "Taiwan", ISO: "TW", DialCode: "886", MinDigits: 9, MaxDigits: 9, NatPattern: "0XXX-XXX-XXX", IntPattern: "+886 XXX-XXX-XXX"},
	"TH": {Name: "Thailand", ISO: "TH", DialCode: "66", MinDigits: 9, MaxDigits: 9, NatPattern: "0XX XXX XXXX", IntPattern: "+66 XX XXX XXXX"},
	"PH": {Name: "Philippines", ISO: "PH", DialCode: "63", MinDigits: 10, MaxDigits: 10, NatPattern: "0XXX XXX XXXX", IntPattern: "+63 XXX XXX XXXX"},
	"ID": {Name: "Indonesia", ISO: "ID", DialCode: "62", MinDigits: 10, MaxDigits: 12, NatPattern: "0XXX-XXXX-XXXX", IntPattern: "+62 XXX-XXXX-XXXX"},
}

// ── helpers ────────────────────────────────────────────────────────────────

// stripNonDigits removes everything except digits from a string.
func stripNonDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// lookupCountry resolves a country name or ISO code to a CountryInfo.
func lookupCountry(input string) (*CountryInfo, error) {
	upper := strings.ToUpper(strings.TrimSpace(input))

	// Try direct ISO code match.
	if info, ok := countries[upper]; ok {
		return &info, nil
	}

	// Try name substring match.
	for _, info := range countries {
		if strings.EqualFold(info.Name, input) {
			return &info, nil
		}
	}
	for _, info := range countries {
		if strings.Contains(strings.ToLower(info.Name), strings.ToLower(input)) {
			return &info, nil
		}
	}

	return nil, fmt.Errorf("country not found: %q", input)
}

// normalizeNumber strips non-digits from a phone number string and
// removes the country dial code prefix if present.
func normalizeNumber(raw string, info *CountryInfo) string {
	digits := stripNonDigits(raw)

	// Remove leading + (already stripped by stripNonDigits).
	// Remove country code prefix if present.
	if strings.HasPrefix(digits, info.DialCode) && len(digits) > len(info.DialCode) {
		withoutCode := digits[len(info.DialCode):]
		if len(withoutCode) >= info.MinDigits && len(withoutCode) <= info.MaxDigits {
			return withoutCode
		}
	}

	// For countries where national numbers start with 0, strip it.
	if len(digits) > 0 && digits[0] == '0' {
		withoutZero := digits[1:]
		if len(withoutZero) >= info.MinDigits-1 && len(withoutZero) <= info.MaxDigits {
			// Some patterns include the leading 0, so check min with -1.
			return withoutZero
		}
	}

	return digits
}

// applyPattern fills X placeholders in a pattern with digits.
func applyPattern(pattern string, digits string) string {
	var b strings.Builder
	idx := 0
	for _, r := range pattern {
		if r == 'X' && idx < len(digits) {
			b.WriteByte(digits[idx])
			idx++
		} else {
			b.WriteRune(r)
		}
	}
	// Append remaining digits if pattern had fewer X than digits.
	if idx < len(digits) {
		b.WriteString(digits[idx:])
	}
	return b.String()
}

// formatE164 produces the E.164 representation: +<dial-code><national-digits>.
func formatE164(digits string, info *CountryInfo) string {
	return "+" + info.DialCode + digits
}

// ── parent command ─────────────────────────────────────────────────────────

var phoneCmd = &cobra.Command{
	Use:   "phone",
	Short: "Phone number parser — format, validate, look up country codes",
	Long: `Phone number tools — format, validate, and look up country dial codes.

SUBCOMMANDS:

  format    Format a phone number (E.164, international, national)
  validate  Basic phone number validation
  country   Look up a country's dial code and format

Supports ~30 countries. All subcommands support --json/-j.

EXAMPLES:

  openGyver phone format "2025551234"
  openGyver phone format "+44 20 7946 0958" --country GB
  openGyver phone validate "2025551234" --country US
  openGyver phone country "Japan"
  openGyver phone country KR`,
}

// ── format subcommand ──────────────────────────────────────────────────────

var formatCountry string

var formatCmd = &cobra.Command{
	Use:   "format <phone-number>",
	Short: "Format a phone number",
	Long: `Format a phone number into E.164, international, and national formats.

Pass the number as a string. Non-digit characters are stripped.
Use --country to specify the country (default: US).

Examples:
  openGyver phone format "2025551234"
  openGyver phone format "+44 20 7946 0958" --country GB
  openGyver phone format "03-1234-5678" --country JP`,
	Args: cobra.ExactArgs(1),
	RunE: runFormat,
}

// FormatResult holds the three formatted representations.
type FormatResult struct {
	Input         string `json:"input"`
	Country       string `json:"country"`
	E164          string `json:"e164"`
	International string `json:"international"`
	National      string `json:"national"`
}

func formatNumber(raw, countryCode string) (*FormatResult, error) {
	info, err := lookupCountry(countryCode)
	if err != nil {
		return nil, err
	}

	digits := normalizeNumber(raw, info)
	if len(digits) == 0 {
		return nil, fmt.Errorf("no digits found in %q", raw)
	}

	return &FormatResult{
		Input:         raw,
		Country:       info.Name,
		E164:          formatE164(digits, info),
		International: applyPattern(info.IntPattern, digits),
		National:      applyPattern(info.NatPattern, digits),
	}, nil
}

func runFormat(_ *cobra.Command, args []string) error {
	result, err := formatNumber(args[0], formatCountry)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(result)
	}

	fmt.Printf("Country:       %s\n", result.Country)
	fmt.Printf("E.164:         %s\n", result.E164)
	fmt.Printf("International: %s\n", result.International)
	fmt.Printf("National:      %s\n", result.National)
	return nil
}

// ── validate subcommand ────────────────────────────────────────────────────

var validateCountry string

var validateCmd = &cobra.Command{
	Use:   "validate <phone-number>",
	Short: "Basic phone number validation",
	Long: `Validate a phone number: check digit count and prefix for the given country.

Examples:
  openGyver phone validate "2025551234" --country US
  openGyver phone validate "+81-03-1234-5678" --country JP`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

// ValidationResult holds the outcome of a phone number check.
type ValidationResult struct {
	Input       string `json:"input"`
	Country     string `json:"country"`
	Valid       bool   `json:"valid"`
	Digits      int    `json:"digits"`
	DialCode    string `json:"dialCode"`
	Reason      string `json:"reason,omitempty"`
	NatDigits   string `json:"nationalDigits"`
}

func validateNumber(raw, countryCode string) *ValidationResult {
	info, err := lookupCountry(countryCode)
	if err != nil {
		return &ValidationResult{
			Input:  raw,
			Valid:  false,
			Reason: err.Error(),
		}
	}

	digits := normalizeNumber(raw, info)
	result := &ValidationResult{
		Input:     raw,
		Country:   info.Name,
		Digits:    len(digits),
		DialCode:  info.DialCode,
		NatDigits: digits,
	}

	if len(digits) < info.MinDigits {
		result.Valid = false
		result.Reason = fmt.Sprintf("too few digits: got %d, need at least %d", len(digits), info.MinDigits)
		return result
	}
	if len(digits) > info.MaxDigits {
		result.Valid = false
		result.Reason = fmt.Sprintf("too many digits: got %d, max is %d", len(digits), info.MaxDigits)
		return result
	}

	result.Valid = true
	return result
}

func runValidate(_ *cobra.Command, args []string) error {
	result := validateNumber(args[0], validateCountry)

	if jsonOut {
		return cmd.PrintJSON(result)
	}

	if result.Valid {
		fmt.Printf("Valid (%s, %d digits)\n", result.Country, result.Digits)
	} else {
		fmt.Printf("Invalid: %s\n", result.Reason)
	}
	return nil
}

// ── country subcommand ─────────────────────────────────────────────────────

var countryCmd = &cobra.Command{
	Use:   "country <name | ISO code>",
	Short: "Look up a country's dial code and format",
	Long: `Look up a country by name or 2-letter ISO code to see its dial code,
typical digit count, and formatting patterns.

Examples:
  openGyver phone country US
  openGyver phone country "Japan"
  openGyver phone country KR
  openGyver phone country "united kingdom"`,
	Args: cobra.ExactArgs(1),
	RunE: runCountry,
}

func runCountry(_ *cobra.Command, args []string) error {
	info, err := lookupCountry(args[0])
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(info)
	}

	fmt.Printf("Country:      %s (%s)\n", info.Name, info.ISO)
	fmt.Printf("Dial Code:    +%s\n", info.DialCode)
	fmt.Printf("Digits:       %d", info.MinDigits)
	if info.MinDigits != info.MaxDigits {
		fmt.Printf("-%d", info.MaxDigits)
	}
	fmt.Println()
	fmt.Printf("National:     %s\n", info.NatPattern)
	fmt.Printf("International:%s\n", info.IntPattern)
	return nil
}

func init() {
	formatCmd.Flags().StringVar(&formatCountry, "country", "US", "country ISO code (e.g. US, GB, JP)")
	validateCmd.Flags().StringVar(&validateCountry, "country", "US", "country ISO code (e.g. US, GB, JP)")

	for _, c := range []*cobra.Command{formatCmd, validateCmd, countryCmd} {
		c.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	}

	phoneCmd.AddCommand(formatCmd, validateCmd, countryCmd)
	cmd.Register(phoneCmd)
}

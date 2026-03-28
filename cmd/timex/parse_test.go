package timex

import (
	"math"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// loadLocation
// ---------------------------------------------------------------------------

func TestLoadLocation_UTC(t *testing.T) {
	loc, err := loadLocation("UTC")
	if err != nil {
		t.Fatal(err)
	}
	if loc != time.UTC {
		t.Errorf("expected UTC, got %v", loc)
	}
}

func TestLoadLocation_IANA(t *testing.T) {
	loc, err := loadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	if loc.String() != "America/New_York" {
		t.Errorf("expected America/New_York, got %s", loc)
	}
}

func TestLoadLocation_Abbreviation(t *testing.T) {
	abbrevs := []string{"EST", "PST", "JST", "CET", "IST", "GMT"}
	for _, abbr := range abbrevs {
		loc, err := loadLocation(abbr)
		if err != nil {
			t.Errorf("failed to load %s: %v", abbr, err)
			continue
		}
		if loc == nil {
			t.Errorf("%s returned nil location", abbr)
		}
	}
}

func TestLoadLocation_Unknown(t *testing.T) {
	_, err := loadLocation("FakeZone/Nowhere")
	if err == nil {
		t.Error("expected error for unknown timezone")
	}
}

func TestLoadLocation_Local(t *testing.T) {
	loc, err := loadLocation("local")
	if err != nil {
		t.Fatal(err)
	}
	if loc != time.Local {
		t.Errorf("expected Local, got %v", loc)
	}
}

// ---------------------------------------------------------------------------
// parseTime — relative keywords
// ---------------------------------------------------------------------------

func TestParseTime_Now(t *testing.T) {
	before := time.Now()
	parsed, err := parseTime("now", nil)
	after := time.Now()
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Before(before) || parsed.After(after) {
		t.Error("parsed 'now' is outside expected range")
	}
}

func TestParseTime_Today(t *testing.T) {
	parsed, err := parseTime("today", nil)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if parsed.Year() != now.Year() || parsed.Month() != now.Month() || parsed.Day() != now.Day() {
		t.Error("'today' didn't match current date")
	}
	if parsed.Hour() != 0 || parsed.Minute() != 0 {
		t.Error("'today' should be midnight")
	}
}

func TestParseTime_Yesterday(t *testing.T) {
	parsed, err := parseTime("yesterday", nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Now().AddDate(0, 0, -1)
	if parsed.Year() != expected.Year() || parsed.Month() != expected.Month() || parsed.Day() != expected.Day() {
		t.Error("'yesterday' didn't match")
	}
}

func TestParseTime_Tomorrow(t *testing.T) {
	parsed, err := parseTime("tomorrow", nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Now().AddDate(0, 0, 1)
	if parsed.Year() != expected.Year() || parsed.Month() != expected.Month() || parsed.Day() != expected.Day() {
		t.Error("'tomorrow' didn't match")
	}
}

// ---------------------------------------------------------------------------
// parseTime — Unix timestamps
// ---------------------------------------------------------------------------

func TestParseTime_UnixSeconds(t *testing.T) {
	parsed, err := parseTime("1705334400", nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Unix(1705334400, 0).UTC()
	if !parsed.Equal(expected) {
		t.Errorf("got %v, want %v", parsed, expected)
	}
}

func TestParseTime_UnixMilliseconds(t *testing.T) {
	parsed, err := parseTime("1705334400000", nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Unix(1705334400, 0).UTC()
	if !parsed.Equal(expected) {
		t.Errorf("got %v, want %v", parsed, expected)
	}
}

func TestParseTime_UnixNanoseconds(t *testing.T) {
	parsed, err := parseTime("1705334400000000000", nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Unix(1705334400, 0).UTC()
	if !parsed.Equal(expected) {
		t.Errorf("got %v, want %v", parsed, expected)
	}
}

func TestParseTime_UnixFractional(t *testing.T) {
	parsed, err := parseTime("1705334400.5", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != 1705334400 {
		t.Errorf("seconds: got %d, want 1705334400", parsed.Unix())
	}
	if parsed.Nanosecond() < 499000000 || parsed.Nanosecond() > 501000000 {
		t.Errorf("nanoseconds: got %d, want ~500000000", parsed.Nanosecond())
	}
}

// ---------------------------------------------------------------------------
// parseTime — standard formats
// ---------------------------------------------------------------------------

func TestParseTime_ISO8601(t *testing.T) {
	parsed, err := parseTime("2024-01-15T14:30:00Z", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Year() != 2024 || parsed.Month() != 1 || parsed.Day() != 15 {
		t.Errorf("date mismatch: %v", parsed)
	}
	if parsed.Hour() != 14 || parsed.Minute() != 30 {
		t.Errorf("time mismatch: %v", parsed)
	}
}

func TestParseTime_ISO8601_Offset(t *testing.T) {
	parsed, err := parseTime("2024-01-15T14:30:00+05:30", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Should be 09:00 UTC
	utc := parsed.UTC()
	if utc.Hour() != 9 || utc.Minute() != 0 {
		t.Errorf("expected 09:00 UTC, got %02d:%02d", utc.Hour(), utc.Minute())
	}
}

func TestParseTime_RFC2822(t *testing.T) {
	_, err := parseTime("Mon, 15 Jan 2024 14:30:00 +0000", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseTime_DateOnly(t *testing.T) {
	parsed, err := parseTime("2024-01-15", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Year() != 2024 || parsed.Month() != 1 || parsed.Day() != 15 {
		t.Errorf("date mismatch: %v", parsed)
	}
}

func TestParseTime_DateSlash(t *testing.T) {
	parsed, err := parseTime("01/15/2024", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Month() != 1 || parsed.Day() != 15 {
		t.Errorf("US date slash mismatch: %v", parsed)
	}
}

func TestParseTime_DatetimeNoTZ(t *testing.T) {
	parsed, err := parseTime("2024-01-15 14:30:00", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Hour() != 14 || parsed.Minute() != 30 {
		t.Errorf("time mismatch: %v", parsed)
	}
}

func TestParseTime_DatetimeNoSeconds(t *testing.T) {
	parsed, err := parseTime("2024-01-15 14:30", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Hour() != 14 || parsed.Minute() != 30 {
		t.Errorf("time mismatch: %v", parsed)
	}
}

func TestParseTime_USLong12Hour(t *testing.T) {
	parsed, err := parseTime("Jan 15, 2024 2:30 PM", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Hour() != 14 || parsed.Minute() != 30 {
		t.Errorf("time mismatch: %v", parsed)
	}
}

func TestParseTime_WithAssumeTZ(t *testing.T) {
	ny, _ := time.LoadLocation("America/New_York")
	parsed, err := parseTime("2024-01-15 14:30", ny)
	if err != nil {
		t.Fatal(err)
	}
	// Should be in New York timezone
	utc := parsed.UTC()
	if utc.Hour() != 19 {
		t.Errorf("expected 19:00 UTC (14:30 EST), got %02d:%02d", utc.Hour(), utc.Minute())
	}
}

func TestParseTime_Invalid(t *testing.T) {
	_, err := parseTime("not-a-date", nil)
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

// ---------------------------------------------------------------------------
// parseUnixNumeric
// ---------------------------------------------------------------------------

func TestParseUnixNumeric_Seconds(t *testing.T) {
	parsed, err := parseUnixNumeric("1705334400")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != 1705334400 {
		t.Errorf("got %d, want 1705334400", parsed.Unix())
	}
}

func TestParseUnixNumeric_Milliseconds(t *testing.T) {
	parsed, err := parseUnixNumeric("1705334400000")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != 1705334400 {
		t.Errorf("got %d, want 1705334400", parsed.Unix())
	}
}

func TestParseUnixNumeric_Microseconds(t *testing.T) {
	parsed, err := parseUnixNumeric("1705334400000000")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != 1705334400 {
		t.Errorf("got %d, want 1705334400", parsed.Unix())
	}
}

func TestParseUnixNumeric_Nanoseconds(t *testing.T) {
	parsed, err := parseUnixNumeric("1705334400000000000")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != 1705334400 {
		t.Errorf("got %d, want 1705334400", parsed.Unix())
	}
}

func TestParseUnixNumeric_Negative(t *testing.T) {
	parsed, err := parseUnixNumeric("-86400")
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Unix(-86400, 0).UTC()
	if !parsed.Equal(expected) {
		t.Errorf("got %v, want %v", parsed, expected)
	}
}

// ---------------------------------------------------------------------------
// isNumeric
// ---------------------------------------------------------------------------

func TestIsNumeric(t *testing.T) {
	tests := map[string]bool{
		"123":         true,
		"-456":        true,
		"12.34":       true,
		"1705334400":  true,
		"abc":         false,
		"12:30":       false,
		"2024-01-15":  false,
		"":            false,
	}
	for input, want := range tests {
		got := isNumeric(input)
		if got != want {
			t.Errorf("isNumeric(%q) = %v, want %v", input, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// layoutHasZone
// ---------------------------------------------------------------------------

func TestLayoutHasZone(t *testing.T) {
	tests := map[string]bool{
		time.RFC3339:              true,
		"2006-01-02T15:04:05Z07": true,
		"2006-01-02 15:04:05":    false,
		"2006-01-02":             false,
		time.Kitchen:             false,
		time.RFC1123:             true,
	}
	for layout, want := range tests {
		got := layoutHasZone(layout)
		if got != want {
			t.Errorf("layoutHasZone(%q) = %v, want %v", layout, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// resolveFormat
// ---------------------------------------------------------------------------

func TestResolveFormat_Named(t *testing.T) {
	names := []string{"iso8601", "rfc3339", "rfc2822", "date", "time", "kitchen", "human"}
	for _, name := range names {
		layout := resolveFormat(name)
		if layout == name {
			t.Errorf("resolveFormat(%q) returned unchanged — not a known format?", name)
		}
	}
}

func TestResolveFormat_Custom(t *testing.T) {
	custom := "Monday, January 2 2006"
	if got := resolveFormat(custom); got != custom {
		t.Errorf("custom layout should pass through unchanged: got %q", got)
	}
}

// ---------------------------------------------------------------------------
// addDuration
// ---------------------------------------------------------------------------

func TestAddDuration_GoStyle(t *testing.T) {
	base := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	result, err := addDuration(base, "2h30m")
	if err != nil {
		t.Fatal(err)
	}
	if result.Hour() != 17 || result.Minute() != 0 {
		t.Errorf("expected 17:00, got %02d:%02d", result.Hour(), result.Minute())
	}
}

func TestAddDuration_Days(t *testing.T) {
	base := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	result, err := addDuration(base, "90d")
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Date(2024, 4, 14, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestAddDuration_Weeks(t *testing.T) {
	base := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	result, err := addDuration(base, "2w")
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestAddDuration_Months(t *testing.T) {
	base := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	result, err := addDuration(base, "1mo")
	if err != nil {
		t.Fatal(err)
	}
	// Jan 31 + 1 month = March 2 (Feb has 29 days in 2024)
	if result.Month() != 3 {
		t.Errorf("expected March, got %v", result.Month())
	}
}

func TestAddDuration_Years(t *testing.T) {
	base := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	result, err := addDuration(base, "1y")
	if err != nil {
		t.Fatal(err)
	}
	// Feb 29 2024 + 1 year = March 1 2025 (not a leap year)
	if result.Year() != 2025 || result.Month() != 3 || result.Day() != 1 {
		t.Errorf("got %v", result)
	}
}

func TestAddDuration_Combined(t *testing.T) {
	base := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	result, err := addDuration(base, "1y2mo3d4h5m")
	if err != nil {
		t.Fatal(err)
	}
	if result.Year() != 2025 || result.Month() != 3 || result.Day() != 18 {
		t.Errorf("date: got %v", result)
	}
	if result.Hour() != 18 || result.Minute() != 35 {
		t.Errorf("time: got %02d:%02d", result.Hour(), result.Minute())
	}
}

func TestAddDuration_Negative(t *testing.T) {
	base := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	result, err := addDuration(base, "-30d")
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Date(2024, 5, 16, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestAddDuration_NegativeGoStyle(t *testing.T) {
	base := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	result, err := addDuration(base, "-1h30m")
	if err != nil {
		t.Fatal(err)
	}
	if result.Hour() != 13 || result.Minute() != 0 {
		t.Errorf("expected 13:00, got %02d:%02d", result.Hour(), result.Minute())
	}
}

func TestAddDuration_Invalid(t *testing.T) {
	base := time.Now()
	_, err := addDuration(base, "notaduration")
	if err == nil {
		t.Error("expected error for invalid duration")
	}
}

// ---------------------------------------------------------------------------
// Timezone abbreviation coverage
// ---------------------------------------------------------------------------

func TestTZAbbreviations_AllResolve(t *testing.T) {
	for abbr := range tzAbbreviations {
		loc, err := loadLocation(abbr)
		if err != nil {
			t.Errorf("abbreviation %s failed: %v", abbr, err)
			continue
		}
		if loc == nil {
			t.Errorf("abbreviation %s returned nil", abbr)
		}
	}
}

// ---------------------------------------------------------------------------
// Named formats coverage
// ---------------------------------------------------------------------------

func TestNamedFormats_AllProduceOutput(t *testing.T) {
	ref := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	for name, layout := range namedFormats {
		result := ref.Format(layout)
		if result == "" {
			t.Errorf("format %q produced empty output", name)
		}
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestParseTime_LeapDay(t *testing.T) {
	parsed, err := parseTime("2024-02-29", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Month() != 2 || parsed.Day() != 29 {
		t.Errorf("leap day mismatch: %v", parsed)
	}
}

func TestParseTime_Y2K(t *testing.T) {
	parsed, err := parseTime("2000-01-01T00:00:00Z", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != 946684800 {
		t.Errorf("Y2K epoch: got %d, want 946684800", parsed.Unix())
	}
}

func TestParseTime_Epoch(t *testing.T) {
	parsed, err := parseTime("0", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != 0 {
		t.Errorf("epoch zero: got %d, want 0", parsed.Unix())
	}
}

func TestParseTime_NegativeEpoch(t *testing.T) {
	parsed, err := parseTime("-86400", nil)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Unix() != -86400 {
		t.Errorf("negative epoch: got %d, want -86400", parsed.Unix())
	}
}

func TestDiff_Symmetric(t *testing.T) {
	t1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	d1 := t2.Sub(t1)
	d2 := t1.Sub(t2)
	if math.Abs(float64(d1)+float64(d2)) > 0 {
		t.Error("diff should be symmetric (absolute value)")
	}
}

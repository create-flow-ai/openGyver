package timex

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Common timezone abbreviation mappings to IANA names.
var tzAbbreviations = map[string]string{
	"EST":  "America/New_York",
	"EDT":  "America/New_York",
	"CST":  "America/Chicago",
	"CDT":  "America/Chicago",
	"MST":  "America/Denver",
	"MDT":  "America/Denver",
	"PST":  "America/Los_Angeles",
	"PDT":  "America/Los_Angeles",
	"IST":  "Asia/Kolkata",
	"JST":  "Asia/Tokyo",
	"KST":  "Asia/Seoul",
	"CST8": "Asia/Shanghai",
	"HKT":  "Asia/Hong_Kong",
	"SGT":  "Asia/Singapore",
	"AEST": "Australia/Sydney",
	"AEDT": "Australia/Sydney",
	"NZST": "Pacific/Auckland",
	"NZDT": "Pacific/Auckland",
	"GMT":  "Europe/London",
	"BST":  "Europe/London",
	"CET":  "Europe/Berlin",
	"CEST": "Europe/Berlin",
	"EET":  "Europe/Bucharest",
	"EEST": "Europe/Bucharest",
	"WET":  "Europe/Lisbon",
	"WEST": "Europe/Lisbon",
	"BRT":  "America/Sao_Paulo",
	"ART":  "America/Argentina/Buenos_Aires",
}

// loadLocation resolves a timezone string to a *time.Location.
// Accepts IANA names (America/New_York), abbreviations (EST, PST), and UTC.
func loadLocation(tz string) (*time.Location, error) {
	if strings.EqualFold(tz, "UTC") || strings.EqualFold(tz, "utc") {
		return time.UTC, nil
	}
	if strings.EqualFold(tz, "local") {
		return time.Local, nil
	}

	// Try IANA name first
	loc, err := time.LoadLocation(tz)
	if err == nil {
		return loc, nil
	}

	// Try abbreviation lookup
	if iana, ok := tzAbbreviations[strings.ToUpper(tz)]; ok {
		return time.LoadLocation(iana)
	}

	return nil, fmt.Errorf("unknown timezone: %q (use IANA names like America/New_York or abbreviations like EST)", tz)
}

// Layouts to try when parsing time strings, ordered from most specific to least.
var parseLayouts = []string{
	time.RFC3339Nano,                 // 2006-01-02T15:04:05.999999999Z07:00
	time.RFC3339,                     // 2006-01-02T15:04:05Z07:00
	"2006-01-02T15:04:05Z0700",      // without colon in offset
	"2006-01-02T15:04:05",           // ISO without timezone
	"2006-01-02T15:04",              // ISO without seconds
	time.RFC1123Z,                    // Mon, 02 Jan 2006 15:04:05 -0700
	time.RFC1123,                     // Mon, 02 Jan 2006 15:04:05 MST
	time.RFC850,                      // Monday, 02-Jan-06 15:04:05 MST
	time.RFC822Z,                     // 02 Jan 06 15:04 -0700
	time.RFC822,                      // 02 Jan 06 15:04 MST
	time.ANSIC,                       // Mon Jan _2 15:04:05 2006
	time.UnixDate,                    // Mon Jan _2 15:04:05 MST 2006
	time.RubyDate,                    // Mon Jan 02 15:04:05 -0700 2006
	"2006-01-02 15:04:05 -0700",     // datetime with offset
	"2006-01-02 15:04:05 MST",       // datetime with tz abbr
	"2006-01-02 15:04:05",           // datetime
	"2006-01-02 15:04",              // datetime without seconds
	"2006-01-02 3:04:05 PM",         // datetime 12-hour
	"2006-01-02 3:04 PM",            // datetime 12-hour without seconds
	"Jan 2, 2006 3:04:05 PM",        // US long 12-hour
	"Jan 2, 2006 3:04 PM",           // US long 12-hour without seconds
	"Jan 2, 2006 15:04:05",          // US long 24-hour
	"Jan 2, 2006 15:04",             // US long 24-hour without seconds
	"January 2, 2006 3:04:05 PM",    // US full 12-hour
	"January 2, 2006 15:04:05",      // US full 24-hour
	"January 2, 2006",               // US full date
	"Jan 2, 2006",                   // US short date
	"Jan 2 2006",                    // US short date no comma
	"02-Jan-2006 15:04:05",          // European with time
	"02-Jan-2006",                   // European date
	"2-Jan-2006",                    // European date single digit
	"2006-01-02",                    // ISO date only
	"01/02/2006 15:04:05",           // US date slash with time
	"01/02/2006 3:04:05 PM",         // US date slash 12-hour
	"01/02/2006 15:04",              // US date slash with time no sec
	"01/02/2006",                    // US date slash
	"02/01/2006",                    // EU date slash (ambiguous, try after US)
	"2006/01/02",                    // Asian date slash
	"2006/01/02 15:04:05",           // Asian date slash with time
	time.Kitchen,                     // 3:04PM
	time.TimeOnly,                    // 15:04:05
	"15:04",                          // time only without seconds
	time.DateOnly,                    // 2006-01-02
	time.DateTime,                    // 2006-01-02 15:04:05
	time.Stamp,                       // Jan _2 15:04:05
	time.StampMilli,                  // Jan _2 15:04:05.000
	time.StampMicro,                  // Jan _2 15:04:05.000000
	time.StampNano,                   // Jan _2 15:04:05.000000000
}

// parseTime attempts to parse an input time string using multiple strategies:
// 1. Relative keywords (now, today, yesterday, tomorrow)
// 2. Unix timestamps (pure numeric input)
// 3. Standard format detection against known layouts
//
// If assumeTZ is non-nil, timezone-naive inputs are assumed to be in that location.
func parseTime(input string, assumeTZ *time.Location) (time.Time, error) {
	trimmed := strings.TrimSpace(input)

	// Relative keywords
	now := time.Now()
	switch strings.ToLower(trimmed) {
	case "now":
		return now, nil
	case "today":
		y, m, d := now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, now.Location()), nil
	case "yesterday":
		y, m, d := now.AddDate(0, 0, -1).Date()
		return time.Date(y, m, d, 0, 0, 0, 0, now.Location()), nil
	case "tomorrow":
		y, m, d := now.AddDate(0, 0, 1).Date()
		return time.Date(y, m, d, 0, 0, 0, 0, now.Location()), nil
	}

	// Unix timestamp detection — pure numeric (with optional decimal for fractional seconds)
	if isNumeric(trimmed) {
		return parseUnixNumeric(trimmed)
	}

	// Default location for naive inputs
	loc := time.UTC
	if assumeTZ != nil {
		loc = assumeTZ
	}

	// Try each layout
	for _, layout := range parseLayouts {
		t, err := time.Parse(layout, trimmed)
		if err == nil {
			// If the parsed time has no timezone info (zone name is empty or "UTC" from naive parse),
			// and the layout doesn't include a zone indicator, apply assumeTZ.
			if !layoutHasZone(layout) && assumeTZ != nil {
				t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse time: %q\nSupported formats: ISO 8601, RFC 2822, RFC 3339, Unix timestamp, and more.\nRun 'openGyver timex --help' for the full list.", input)
}

// isNumeric returns true if the string is a valid integer or float (possibly negative).
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// parseUnixNumeric parses a numeric string as a Unix timestamp.
// Auto-detects seconds, milliseconds, microseconds, or nanoseconds based on magnitude.
func parseUnixNumeric(s string) (time.Time, error) {
	// Handle fractional seconds (e.g. 1705334400.123)
	if strings.Contains(s, ".") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return time.Time{}, err
		}
		sec := int64(f)
		nsec := int64((f - float64(sec)) * 1e9)
		return time.Unix(sec, nsec).UTC(), nil
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	switch {
	case n > 1e18: // nanoseconds
		return time.Unix(0, n).UTC(), nil
	case n > 1e15: // microseconds
		return time.Unix(0, n*1000).UTC(), nil
	case n > 1e12: // milliseconds
		return time.Unix(0, n*1000000).UTC(), nil
	default: // seconds
		return time.Unix(n, 0).UTC(), nil
	}
}

// layoutHasZone checks if a Go time layout string contains timezone information.
func layoutHasZone(layout string) bool {
	// Go's reference time timezone indicators
	zoneIndicators := []string{"Z07", "Z0700", "Z07:00", "-07", "-0700", "-07:00", "MST"}
	for _, z := range zoneIndicators {
		if strings.Contains(layout, z) {
			return true
		}
	}
	return false
}

// Named output formats.
var namedFormats = map[string]string{
	"iso8601":  "2006-01-02T15:04:05Z07:00",
	"rfc3339":  time.RFC3339,
	"rfc2822":  "Mon, 02 Jan 2006 15:04:05 -0700",
	"rfc1123":  time.RFC1123Z,
	"rfc850":   time.RFC850,
	"rfc822":   time.RFC822Z,
	"ansic":    time.ANSIC,
	"unix":     time.UnixDate,
	"ruby":     time.RubyDate,
	"kitchen":  time.Kitchen,
	"date":     "2006-01-02",
	"time":     "15:04:05",
	"datetime": "2006-01-02 15:04:05",
	"us":       "01/02/2006 3:04:05 PM",
	"eu":       "02/01/2006 15:04:05",
	"short":    "Jan 2, 2006",
	"long":     "January 2, 2006 15:04:05 MST",
	"stamp":    time.Stamp,
	"human":    "Mon, Jan 2 2006 at 3:04 PM MST",
}

// resolveFormat maps a format name to a Go layout, or returns the input as a custom layout.
func resolveFormat(name string) string {
	if layout, ok := namedFormats[strings.ToLower(name)]; ok {
		return layout
	}
	return name
}

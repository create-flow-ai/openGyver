package cron

import (
	"testing"
	"time"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestCronCmd_Metadata(t *testing.T) {
	if cronCmd.Use != "cron" {
		t.Errorf("unexpected Use: %s", cronCmd.Use)
	}
	if cronCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestSubcommands_Registered(t *testing.T) {
	names := map[string]bool{}
	for _, sub := range cronCmd.Commands() {
		names[sub.Name()] = true
	}
	for _, want := range []string{"explain", "next", "validate"} {
		if !names[want] {
			t.Errorf("subcommand %q not registered", want)
		}
	}
}

// ── parseExpression ────────────────────────────────────────────────────────

func TestParseExpression_FiveFields(t *testing.T) {
	fields, hasSeconds, err := parseExpression("*/5 * * * *")
	if err != nil {
		t.Fatal(err)
	}
	if hasSeconds {
		t.Error("expected 5-field format")
	}
	if len(fields) != 5 {
		t.Errorf("expected 5 fields, got %d", len(fields))
	}
	if fields[0] != "*/5" {
		t.Errorf("first field: got %q, want */5", fields[0])
	}
}

func TestParseExpression_SixFields(t *testing.T) {
	fields, hasSeconds, err := parseExpression("0 */5 * * * *")
	if err != nil {
		t.Fatal(err)
	}
	if !hasSeconds {
		t.Error("expected 6-field format with seconds")
	}
	if len(fields) != 6 {
		t.Errorf("expected 6 fields, got %d", len(fields))
	}
}

func TestParseExpression_InvalidCount(t *testing.T) {
	for _, expr := range []string{"* * *", "* * * * * * *", ""} {
		_, _, err := parseExpression(expr)
		if err == nil {
			t.Errorf("expected error for %q", expr)
		}
	}
}

// ── expandField ────────────────────────────────────────────────────────────

func TestExpandField_Wildcard(t *testing.T) {
	vals, err := expandField("*", fieldRange{0, 59, "minute"})
	if err != nil {
		t.Fatal(err)
	}
	if vals != nil {
		t.Error("wildcard should return nil")
	}
}

func TestExpandField_Step(t *testing.T) {
	vals, err := expandField("*/15", fieldRange{0, 59, "minute"})
	if err != nil {
		t.Fatal(err)
	}
	expected := map[int]bool{0: true, 15: true, 30: true, 45: true}
	if len(vals) != len(expected) {
		t.Errorf("expected %d values, got %d", len(expected), len(vals))
	}
	for k := range expected {
		if !vals[k] {
			t.Errorf("expected value %d", k)
		}
	}
}

func TestExpandField_Range(t *testing.T) {
	vals, err := expandField("1-5", fieldRange{0, 6, "weekday"})
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 5 {
		t.Errorf("expected 5 values, got %d", len(vals))
	}
	for i := 1; i <= 5; i++ {
		if !vals[i] {
			t.Errorf("expected value %d", i)
		}
	}
}

func TestExpandField_List(t *testing.T) {
	vals, err := expandField("1,15", fieldRange{1, 31, "day"})
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Errorf("expected 2 values, got %d", len(vals))
	}
	if !vals[1] || !vals[15] {
		t.Error("expected values 1 and 15")
	}
}

func TestExpandField_Last(t *testing.T) {
	vals, err := expandField("L", fieldRange{1, 31, "day"})
	if err != nil {
		t.Fatal(err)
	}
	if !vals[31] {
		t.Error("L should map to max value")
	}
}

func TestExpandField_SingleValue(t *testing.T) {
	vals, err := expandField("30", fieldRange{0, 59, "minute"})
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 1 || !vals[30] {
		t.Error("expected single value 30")
	}
}

func TestExpandField_OutOfRange(t *testing.T) {
	_, err := expandField("60", fieldRange{0, 59, "minute"})
	if err == nil {
		t.Error("expected out-of-range error")
	}
}

func TestExpandField_InvalidStep(t *testing.T) {
	_, err := expandField("*/abc", fieldRange{0, 59, "minute"})
	if err == nil {
		t.Error("expected error for invalid step")
	}
}

// ── matches ────────────────────────────────────────────────────────────────

func TestMatches_Wildcard(t *testing.T) {
	pf := parsedField{values: nil}
	if !matches(pf, 42) {
		t.Error("wildcard should match any value")
	}
}

func TestMatches_Specific(t *testing.T) {
	pf := parsedField{values: map[int]bool{5: true, 10: true}}
	if !matches(pf, 5) {
		t.Error("should match 5")
	}
	if matches(pf, 6) {
		t.Error("should not match 6")
	}
}

// ── nextOccurrences ────────────────────────────────────────────────────────

func TestNextOccurrences_EveryMinute(t *testing.T) {
	// "* * * * *" — every minute
	fields := []string{"*", "*", "*", "*", "*"}
	parsed, err := compileFields(fields, false)
	if err != nil {
		t.Fatal(err)
	}
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	results := nextOccurrences(parsed, false, from, 3)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// First result should be 00:01 (next minute after 00:00)
	want := time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC)
	if !results[0].Equal(want) {
		t.Errorf("first result: got %v, want %v", results[0], want)
	}
}

func TestNextOccurrences_EveryFiveMinutes(t *testing.T) {
	// "*/5 * * * *"
	fields := []string{"*/5", "*", "*", "*", "*"}
	parsed, err := compileFields(fields, false)
	if err != nil {
		t.Fatal(err)
	}
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	results := nextOccurrences(parsed, false, from, 3)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	expected := []int{5, 10, 15}
	for i, r := range results {
		if r.Minute() != expected[i] {
			t.Errorf("result %d: got minute %d, want %d", i, r.Minute(), expected[i])
		}
	}
}

// ── describeField ──────────────────────────────────────────────────────────

func TestDescribeField_Wildcard(t *testing.T) {
	got := describeField("*", fieldRange{0, 59, "minute"})
	if got != "every minute" {
		t.Errorf("got %q", got)
	}
}

func TestDescribeField_Step(t *testing.T) {
	got := describeField("*/5", fieldRange{0, 59, "minute"})
	if got != "every 5 minutes" {
		t.Errorf("got %q", got)
	}
}

func TestDescribeField_Month(t *testing.T) {
	got := describeField("1", fieldRange{1, 12, "month"})
	if got != "month(s) January" {
		t.Errorf("got %q", got)
	}
}

func TestDescribeField_Weekday(t *testing.T) {
	got := describeField("1-5", fieldRange{0, 6, "weekday"})
	if got != "Monday through Friday" {
		t.Errorf("got %q", got)
	}
}

// ── compile validation ─────────────────────────────────────────────────────

func TestCompileFields_InvalidField(t *testing.T) {
	fields := []string{"abc", "*", "*", "*", "*"}
	_, err := compileFields(fields, false)
	if err == nil {
		t.Error("expected error for invalid field")
	}
}

func TestExpandField_RangeStep(t *testing.T) {
	vals, err := expandField("1-10/3", fieldRange{0, 59, "minute"})
	if err != nil {
		t.Fatal(err)
	}
	expected := map[int]bool{1: true, 4: true, 7: true, 10: true}
	if len(vals) != len(expected) {
		t.Errorf("expected %d values, got %d", len(expected), len(vals))
	}
	for k := range expected {
		if !vals[k] {
			t.Errorf("expected value %d", k)
		}
	}
}

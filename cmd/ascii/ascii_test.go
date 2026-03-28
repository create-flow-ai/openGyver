package ascii

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestAsciiCmd_Metadata(t *testing.T) {
	if asciiCmd.Use != "ascii" {
		t.Errorf("unexpected Use: %s", asciiCmd.Use)
	}
	if asciiCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestAsciiCmd_Subcommands(t *testing.T) {
	names := map[string]bool{}
	for _, sub := range asciiCmd.Commands() {
		names[sub.Name()] = true
	}
	for _, want := range []string{"banner", "table", "lookup"} {
		if !names[want] {
			t.Errorf("missing subcommand: %s", want)
		}
	}
}

// ── banner ─────────────────────────────────────────────────────────────────

func TestRenderBanner_Hello(t *testing.T) {
	banner := renderBanner("HELLO")
	lines := strings.Split(banner, "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
	// Each line should contain # characters.
	for i, line := range lines {
		if !strings.Contains(line, "#") {
			t.Errorf("line %d should contain # characters: %q", i, line)
		}
	}
}

func TestRenderBanner_SingleChar(t *testing.T) {
	banner := renderBanner("A")
	lines := strings.Split(banner, "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
	// First line of A is " ### "
	if !strings.Contains(lines[0], "###") {
		t.Errorf("unexpected first line for 'A': %q", lines[0])
	}
}

func TestRenderBanner_Space(t *testing.T) {
	banner := renderBanner("A B")
	lines := strings.Split(banner, "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
	// The middle character is a space so there should be a gap.
	if len(lines[0]) < 15 {
		t.Errorf("expected wider banner with space: %q", lines[0])
	}
}

func TestRenderBanner_Lowercase(t *testing.T) {
	// Lowercase should be auto-uppercased.
	upper := renderBanner("HI")
	lower := renderBanner("hi")
	if upper != lower {
		t.Error("lowercase input should produce same banner as uppercase")
	}
}

func TestRenderBanner_UnknownChar(t *testing.T) {
	// Unknown characters should fall back to '?'
	banner := renderBanner("@")
	qBanner := renderBanner("?")
	if banner != qBanner {
		t.Error("unknown character should render as '?'")
	}
}

// ── table ──────────────────────────────────────────────────────────────────

func TestBuildTable_Length(t *testing.T) {
	entries := buildTable()
	// Characters 32-126 = 95 entries.
	if len(entries) != 95 {
		t.Errorf("expected 95 entries, got %d", len(entries))
	}
}

func TestBuildTable_FirstEntry(t *testing.T) {
	entries := buildTable()
	first := entries[0]
	if first.Decimal != 32 {
		t.Errorf("first entry decimal: got %d, want 32", first.Decimal)
	}
	if first.Char != "SP" {
		t.Errorf("first entry char: got %q, want SP", first.Char)
	}
	if first.Hex != "0x20" {
		t.Errorf("first entry hex: got %s, want 0x20", first.Hex)
	}
}

func TestBuildTable_LastEntry(t *testing.T) {
	entries := buildTable()
	last := entries[len(entries)-1]
	if last.Decimal != 126 {
		t.Errorf("last entry decimal: got %d, want 126", last.Decimal)
	}
	if last.Char != "~" {
		t.Errorf("last entry char: got %q, want ~", last.Char)
	}
}

func TestBuildTable_LetterA(t *testing.T) {
	entries := buildTable()
	// 'A' is decimal 65, index 65-32 = 33.
	a := entries[33]
	if a.Decimal != 65 {
		t.Errorf("A decimal: got %d, want 65", a.Decimal)
	}
	if a.Char != "A" {
		t.Errorf("A char: got %q, want A", a.Char)
	}
}

// ── lookup ─────────────────────────────────────────────────────────────────

func TestCharLookup_ByDecimal(t *testing.T) {
	r, err := charLookup("65")
	if err != nil {
		t.Fatal(err)
	}
	if r.Char != "A" {
		t.Errorf("char: got %q, want A", r.Char)
	}
	if r.Decimal != 65 {
		t.Errorf("decimal: got %d, want 65", r.Decimal)
	}
	if r.Hex != "0x41" {
		t.Errorf("hex: got %s, want 0x41", r.Hex)
	}
	if r.Binary != "01000001" {
		t.Errorf("binary: got %s, want 01000001", r.Binary)
	}
}

func TestCharLookup_ByChar(t *testing.T) {
	r, err := charLookup("A")
	if err != nil {
		t.Fatal(err)
	}
	if r.Decimal != 65 {
		t.Errorf("decimal: got %d, want 65", r.Decimal)
	}
	if r.HTMLEntity != "&#65;" {
		t.Errorf("HTML entity: got %s, want &#65;", r.HTMLEntity)
	}
	if r.URLEncoding != "%41" {
		t.Errorf("URL encoding: got %s, want %%41", r.URLEncoding)
	}
}

func TestCharLookup_Space(t *testing.T) {
	r, err := charLookup("32")
	if err != nil {
		t.Fatal(err)
	}
	if r.Char != "SP" {
		t.Errorf("char: got %q, want SP", r.Char)
	}
	if r.URLEncoding != "%20" {
		t.Errorf("URL encoding: got %s, want %%20", r.URLEncoding)
	}
}

func TestCharLookup_OutOfRange(t *testing.T) {
	_, err := charLookup("200")
	if err == nil {
		t.Error("expected error for value > 127")
	}
}

func TestCharLookup_InvalidInput(t *testing.T) {
	_, err := charLookup("hello")
	if err == nil {
		t.Error("expected error for multi-char string")
	}
}

func TestCharLookup_ControlChar(t *testing.T) {
	r, err := charLookup("0")
	if err != nil {
		t.Fatal(err)
	}
	if r.Char != "NUL" {
		t.Errorf("char: got %q, want NUL", r.Char)
	}
}

// ── flags ──────────────────────────────────────────────────────────────────

func TestBannerCmd_AcceptsExactlyOneArg(t *testing.T) {
	validator := cobra.ExactArgs(1)
	if err := validator(bannerCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(bannerCmd, []string{"text"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTableCmd_NoArgs(t *testing.T) {
	validator := cobra.NoArgs
	if err := validator(tableCmd, []string{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := validator(tableCmd, []string{"extra"}); err == nil {
		t.Error("expected error with args")
	}
}

func TestLookupCmd_HasJSONFlag(t *testing.T) {
	f := lookupCmd.Flags().Lookup("json")
	if f == nil {
		t.Error("lookup should have --json flag")
	}
	if f.Shorthand != "j" {
		t.Errorf("--json shorthand: got %q, want 'j'", f.Shorthand)
	}
}

// ── isASCIIPrintable ───────────────────────────────────────────────────────

func TestIsASCIIPrintable(t *testing.T) {
	if !isASCIIPrintable('A') {
		t.Error("'A' should be printable")
	}
	if !isASCIIPrintable(' ') {
		t.Error("space should be printable")
	}
	if isASCIIPrintable(0) {
		t.Error("NUL should not be printable")
	}
	if isASCIIPrintable(127) {
		t.Error("DEL should not be printable")
	}
}

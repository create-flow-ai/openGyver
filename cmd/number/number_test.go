package number

import (
	"strconv"
	"strings"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestNumberCmd_Metadata(t *testing.T) {
	if numberCmd.Use == "" {
		t.Error("numberCmd.Use must not be empty")
	}
	if numberCmd.Short == "" {
		t.Error("numberCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"baseCmd", baseCmd.Use, baseCmd.Short},
		{"romanCmd", romanCmd.Use, romanCmd.Short},
		{"ieee754Cmd", ieee754Cmd.Use, ieee754Cmd.Short},
	}
	for _, c := range cmds {
		if c.use == "" {
			t.Errorf("%s.Use must not be empty", c.name)
		}
		if c.short == "" {
			t.Errorf("%s.Short must not be empty", c.name)
		}
	}
}

// ── flag existence and defaults ────────────────────────────────────────────

func TestNumberCmd_PersistentFlags(t *testing.T) {
	f := numberCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

func TestBaseCmd_Flags(t *testing.T) {
	f := baseCmd.Flags()
	if f.Lookup("from") == nil {
		t.Error("expected flag --from on baseCmd")
	}
	if f.Lookup("to") == nil {
		t.Error("expected flag --to on baseCmd")
	}
}

// ── base conversion ────────────────────────────────────────────────────────

func TestBaseConversion_DecToHex(t *testing.T) {
	// 255 base10 -> ff base16
	n, err := strconv.ParseInt("255", 10, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	got := strconv.FormatInt(n, 16)
	if got != "ff" {
		t.Errorf("255 base10 -> base16 = %q, want %q", got, "ff")
	}
}

func TestBaseConversion_HexToBin(t *testing.T) {
	// ff base16 -> 11111111 base2
	n, err := strconv.ParseInt("ff", 16, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	got := strconv.FormatInt(n, 2)
	if got != "11111111" {
		t.Errorf("ff base16 -> base2 = %q, want %q", got, "11111111")
	}
}

func TestBaseConversion_BinToDec(t *testing.T) {
	n, err := strconv.ParseInt("11111111", 2, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	got := strconv.FormatInt(n, 10)
	if got != "255" {
		t.Errorf("11111111 base2 -> base10 = %q, want %q", got, "255")
	}
}

func TestBaseConversion_OctalToHex(t *testing.T) {
	n, err := strconv.ParseInt("377", 8, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	got := strconv.FormatInt(n, 16)
	if got != "ff" {
		t.Errorf("377 base8 -> base16 = %q, want %q", got, "ff")
	}
}

func TestBaseConversion_DecToBase36(t *testing.T) {
	n, err := strconv.ParseInt("1000", 10, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	got := strconv.FormatInt(n, 36)
	if got != "rs" {
		t.Errorf("1000 base10 -> base36 = %q, want %q", got, "rs")
	}
}

func TestBaseConversion_Zero(t *testing.T) {
	n, err := strconv.ParseInt("0", 10, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	got := strconv.FormatInt(n, 16)
	if got != "0" {
		t.Errorf("0 base10 -> base16 = %q, want %q", got, "0")
	}
}

func TestBaseConversion_RoundTrip(t *testing.T) {
	// 255 -> hex -> bin -> back to dec
	n, _ := strconv.ParseInt("255", 10, 64)
	hexStr := strconv.FormatInt(n, 16)
	n2, _ := strconv.ParseInt(hexStr, 16, 64)
	binStr := strconv.FormatInt(n2, 2)
	n3, _ := strconv.ParseInt(binStr, 2, 64)
	if n3 != 255 {
		t.Errorf("roundtrip failed: got %d, want 255", n3)
	}
}

// ── isDecimal ──────────────────────────────────────────────────────────────

func TestIsDecimal(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"42", true},
		{"0", true},
		{"-1", true},
		{"+5", true},
		{"XLII", false},
		{"", false},
		{"-", false},
		{"+", false},
		{"12.5", false},
		{"abc", false},
	}
	for _, tt := range tests {
		got := isDecimal(tt.input)
		if got != tt.want {
			t.Errorf("isDecimal(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ── Roman numerals: intToRoman ─────────────────────────────────────────────

func TestIntToRoman_Known(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "I"},
		{4, "IV"},
		{9, "IX"},
		{14, "XIV"},
		{42, "XLII"},
		{99, "XCIX"},
		{100, "C"},
		{399, "CCCXCIX"},
		{500, "D"},
		{944, "CMXLIV"},
		{1000, "M"},
		{1994, "MCMXCIV"},
		{3999, "MMMCMXCIX"},
	}
	for _, tt := range tests {
		got, err := intToRoman(tt.n)
		if err != nil {
			t.Fatalf("intToRoman(%d): %v", tt.n, err)
		}
		if got != tt.want {
			t.Errorf("intToRoman(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestIntToRoman_OutOfRange(t *testing.T) {
	_, err := intToRoman(0)
	if err == nil {
		t.Error("expected error for intToRoman(0)")
	}

	_, err = intToRoman(-1)
	if err == nil {
		t.Error("expected error for intToRoman(-1)")
	}

	_, err = intToRoman(4000)
	if err == nil {
		t.Error("expected error for intToRoman(4000)")
	}
}

func TestIntToRoman_EdgeCases(t *testing.T) {
	// Minimum valid
	got, err := intToRoman(1)
	if err != nil {
		t.Fatalf("intToRoman(1): %v", err)
	}
	if got != "I" {
		t.Errorf("intToRoman(1) = %q, want %q", got, "I")
	}

	// Maximum valid
	got, err = intToRoman(3999)
	if err != nil {
		t.Fatalf("intToRoman(3999): %v", err)
	}
	if got != "MMMCMXCIX" {
		t.Errorf("intToRoman(3999) = %q, want %q", got, "MMMCMXCIX")
	}
}

// ── Roman numerals: romanToInt ─────────────────────────────────────────────

func TestRomanToInt_Known(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"I", 1},
		{"IV", 4},
		{"IX", 9},
		{"XLII", 42},
		{"XCIX", 99},
		{"C", 100},
		{"MCMXCIV", 1994},
		{"MMMCMXCIX", 3999},
	}
	for _, tt := range tests {
		got, err := romanToInt(tt.input)
		if err != nil {
			t.Fatalf("romanToInt(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("romanToInt(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestRomanToInt_CaseInsensitive(t *testing.T) {
	got, err := romanToInt("xlii")
	if err != nil {
		t.Fatalf("romanToInt(xlii): %v", err)
	}
	if got != 42 {
		t.Errorf("romanToInt(xlii) = %d, want 42", got)
	}

	got, err = romanToInt("mcmxciv")
	if err != nil {
		t.Fatalf("romanToInt(mcmxciv): %v", err)
	}
	if got != 1994 {
		t.Errorf("romanToInt(mcmxciv) = %d, want 1994", got)
	}
}

func TestRomanToInt_InvalidChars(t *testing.T) {
	_, err := romanToInt("ABC")
	if err == nil {
		t.Error("expected error for invalid Roman numeral chars")
	}
}

func TestRomanToInt_Empty(t *testing.T) {
	_, err := romanToInt("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestRomanToInt_NonCanonical(t *testing.T) {
	// "IIII" is not canonical (should be "IV")
	_, err := romanToInt("IIII")
	if err == nil {
		t.Error("expected error for non-canonical Roman numeral IIII")
	}
}

// ── Roman numeral roundtrip ────────────────────────────────────────────────

func TestRomanRoundTrip(t *testing.T) {
	for n := 1; n <= 3999; n++ {
		roman, err := intToRoman(n)
		if err != nil {
			t.Fatalf("intToRoman(%d): %v", n, err)
		}
		back, err := romanToInt(roman)
		if err != nil {
			t.Fatalf("romanToInt(%q): %v", roman, err)
		}
		if back != n {
			t.Errorf("roundtrip failed: %d -> %q -> %d", n, roman, back)
		}
	}
}

// ── Roman numeral: specific edge case from prompt ──────────────────────────

func TestRoman_42_XLII(t *testing.T) {
	roman, err := intToRoman(42)
	if err != nil {
		t.Fatalf("intToRoman(42): %v", err)
	}
	if roman != "XLII" {
		t.Errorf("intToRoman(42) = %q, want %q", roman, "XLII")
	}

	n, err := romanToInt("XLII")
	if err != nil {
		t.Fatalf("romanToInt(XLII): %v", err)
	}
	if n != 42 {
		t.Errorf("romanToInt(XLII) = %d, want 42", n)
	}
}

// ── signLabel ──────────────────────────────────────────────────────────────

func TestSignLabel(t *testing.T) {
	if signLabel(uint32(0)) != "+" {
		t.Error("expected '+' for sign 0")
	}
	if signLabel(uint32(1)) != "-" {
		t.Error("expected '-' for sign 1")
	}
	if signLabel(uint64(0)) != "+" {
		t.Error("expected '+' for sign 0 (uint64)")
	}
	if signLabel(uint64(1)) != "-" {
		t.Error("expected '-' for sign 1 (uint64)")
	}
}

// ── romanTable completeness ────────────────────────────────────────────────

func TestRomanTable_NotEmpty(t *testing.T) {
	if len(romanTable) == 0 {
		t.Error("romanTable must not be empty")
	}
}

func TestRomanTable_DescendingOrder(t *testing.T) {
	for i := 1; i < len(romanTable); i++ {
		if romanTable[i].value > romanTable[i-1].value {
			t.Errorf("romanTable not descending at index %d: %d > %d",
				i, romanTable[i].value, romanTable[i-1].value)
		}
	}
}

// ── romanCharValues completeness ───────────────────────────────────────────

func TestRomanCharValues_AllPresent(t *testing.T) {
	expected := []rune{'I', 'V', 'X', 'L', 'C', 'D', 'M'}
	for _, r := range expected {
		if _, ok := romanCharValues[r]; !ok {
			t.Errorf("romanCharValues missing %c", r)
		}
	}
}

// ── base conversion: edge case for case insensitivity ──────────────────────

func TestBaseConversion_CaseInsensitive(t *testing.T) {
	n1, err := strconv.ParseInt(strings.ToLower("FF"), 16, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	n2, err := strconv.ParseInt("ff", 16, 64)
	if err != nil {
		t.Fatalf("ParseInt error: %v", err)
	}
	if n1 != n2 {
		t.Errorf("case sensitivity issue: %d != %d", n1, n2)
	}
}

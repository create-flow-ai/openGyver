package color

import (
	"math"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestColorCmd_Metadata(t *testing.T) {
	if colorCmd.Use == "" {
		t.Error("colorCmd.Use must not be empty")
	}
	if colorCmd.Short == "" {
		t.Error("colorCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"convertCmd", convertCmd.Use, convertCmd.Short},
		{"contrastCmd", contrastCmd.Use, contrastCmd.Short},
		{"paletteCmd", paletteCmd.Use, paletteCmd.Short},
		{"nameCmd", nameCmd.Use, nameCmd.Short},
		{"randomCmd", randomCmd.Use, randomCmd.Short},
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

// ── flag existence ─────────────────────────────────────────────────────────

func TestColorCmd_PersistentFlags(t *testing.T) {
	f := colorCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

func TestConvertCmd_Flags(t *testing.T) {
	f := convertCmd.Flags()
	if f.Lookup("to") == nil {
		t.Error("expected flag --to on convertCmd")
	}
}

func TestPaletteCmd_Flags(t *testing.T) {
	f := paletteCmd.Flags()
	if f.Lookup("count") == nil {
		t.Error("expected flag --count on paletteCmd")
	}
	if f.Lookup("type") == nil {
		t.Error("expected flag --type on paletteCmd")
	}
}

func TestRandomCmd_Flags(t *testing.T) {
	f := randomCmd.Flags()
	if f.Lookup("format") == nil {
		t.Error("expected flag --format on randomCmd")
	}
	if f.Lookup("count") == nil {
		t.Error("expected flag --count on randomCmd")
	}
}

// ── hexToRGB ───────────────────────────────────────────────────────────────

func TestHexToRGB_SixDigit(t *testing.T) {
	tests := []struct {
		hex  string
		want RGB
	}{
		{"#ff0000", RGB{255, 0, 0}},
		{"#00ff00", RGB{0, 255, 0}},
		{"#0000ff", RGB{0, 0, 255}},
		{"#ffffff", RGB{255, 255, 255}},
		{"#000000", RGB{0, 0, 0}},
		{"ff5733", RGB{255, 87, 51}}, // without #
	}
	for _, tt := range tests {
		got, err := hexToRGB(tt.hex)
		if err != nil {
			t.Fatalf("hexToRGB(%q): %v", tt.hex, err)
		}
		if got != tt.want {
			t.Errorf("hexToRGB(%q) = %v, want %v", tt.hex, got, tt.want)
		}
	}
}

func TestHexToRGB_ThreeDigit(t *testing.T) {
	got, err := hexToRGB("#f00")
	if err != nil {
		t.Fatalf("hexToRGB(#f00): %v", err)
	}
	if got != (RGB{255, 0, 0}) {
		t.Errorf("hexToRGB(#f00) = %v, want %v", got, RGB{255, 0, 0})
	}

	got, err = hexToRGB("#fff")
	if err != nil {
		t.Fatalf("hexToRGB(#fff): %v", err)
	}
	if got != (RGB{255, 255, 255}) {
		t.Errorf("hexToRGB(#fff) = %v, want %v", got, RGB{255, 255, 255})
	}
}

func TestHexToRGB_Invalid(t *testing.T) {
	invalids := []string{"#gg0000", "#12345", "#", "zzzzzz", "#1234567"}
	for _, h := range invalids {
		_, err := hexToRGB(h)
		if err == nil {
			t.Errorf("expected error for hexToRGB(%q)", h)
		}
	}
}

// ── rgbToHex ───────────────────────────────────────────────────────────────

func TestRGBToHex(t *testing.T) {
	tests := []struct {
		rgb  RGB
		want string
	}{
		{RGB{255, 0, 0}, "#ff0000"},
		{RGB{0, 255, 0}, "#00ff00"},
		{RGB{0, 0, 255}, "#0000ff"},
		{RGB{255, 255, 255}, "#ffffff"},
		{RGB{0, 0, 0}, "#000000"},
	}
	for _, tt := range tests {
		got := rgbToHex(tt.rgb)
		if got != tt.want {
			t.Errorf("rgbToHex(%v) = %q, want %q", tt.rgb, got, tt.want)
		}
	}
}

// ── RGB to HSL conversion ──────────────────────────────────────────────────

func TestRGBToHSL_KnownValues(t *testing.T) {
	tests := []struct {
		rgb  RGB
		wantH, wantS, wantL float64
	}{
		{RGB{255, 0, 0}, 0, 100, 50},     // pure red
		{RGB{0, 255, 0}, 120, 100, 50},   // pure green
		{RGB{0, 0, 255}, 240, 100, 50},   // pure blue
		{RGB{255, 255, 255}, 0, 0, 100},  // white
		{RGB{0, 0, 0}, 0, 0, 0},          // black
		{RGB{128, 128, 128}, 0, 0, 50.2}, // gray (approximately)
	}
	for _, tt := range tests {
		hsl := rgbToHSL(tt.rgb)
		if math.Abs(hsl.H-tt.wantH) > 1 {
			t.Errorf("rgbToHSL(%v).H = %v, want ~%v", tt.rgb, hsl.H, tt.wantH)
		}
		if math.Abs(hsl.S-tt.wantS) > 1 {
			t.Errorf("rgbToHSL(%v).S = %v, want ~%v", tt.rgb, hsl.S, tt.wantS)
		}
		if math.Abs(hsl.L-tt.wantL) > 1 {
			t.Errorf("rgbToHSL(%v).L = %v, want ~%v", tt.rgb, hsl.L, tt.wantL)
		}
	}
}

// ── HSL to RGB roundtrip ───────────────────────────────────────────────────

func TestHSLToRGB_RoundTrip(t *testing.T) {
	colors := []RGB{
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
		{255, 255, 0},
		{0, 255, 255},
		{255, 0, 255},
		{128, 64, 32},
	}
	for _, c := range colors {
		hsl := rgbToHSL(c)
		back := hslToRGB(hsl)
		// Allow tolerance of 1 per channel due to rounding
		if abs(int(c.R)-int(back.R)) > 1 || abs(int(c.G)-int(back.G)) > 1 || abs(int(c.B)-int(back.B)) > 1 {
			t.Errorf("HSL roundtrip failed for %v: HSL=%v, back=%v", c, hsl, back)
		}
	}
}

func TestHSLToRGB_Achromatic(t *testing.T) {
	// Pure gray: S=0
	rgb := hslToRGB(HSL{0, 0, 50})
	if rgb.R != rgb.G || rgb.G != rgb.B {
		t.Errorf("expected achromatic gray, got %v", rgb)
	}
	expected := uint8(math.Round(0.5 * 255))
	if abs(int(rgb.R)-int(expected)) > 1 {
		t.Errorf("expected gray ~%d, got %d", expected, rgb.R)
	}
}

// ── CMYK conversion ────────────────────────────────────────────────────────

func TestRGBToCMYK_KnownValues(t *testing.T) {
	tests := []struct {
		rgb  RGB
		want CMYK
	}{
		{RGB{255, 0, 0}, CMYK{0, 100, 100, 0}},     // red
		{RGB{0, 255, 0}, CMYK{100, 0, 100, 0}},     // green
		{RGB{0, 0, 255}, CMYK{100, 100, 0, 0}},     // blue
		{RGB{255, 255, 255}, CMYK{0, 0, 0, 0}},     // white
		{RGB{0, 0, 0}, CMYK{0, 0, 0, 100}},         // black
	}
	for _, tt := range tests {
		got := rgbToCMYK(tt.rgb)
		if math.Abs(got.C-tt.want.C) > 0.5 ||
			math.Abs(got.M-tt.want.M) > 0.5 ||
			math.Abs(got.Y-tt.want.Y) > 0.5 ||
			math.Abs(got.K-tt.want.K) > 0.5 {
			t.Errorf("rgbToCMYK(%v) = %v, want %v", tt.rgb, got, tt.want)
		}
	}
}

func TestCMYKToRGB_RoundTrip(t *testing.T) {
	colors := []RGB{
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
		{128, 64, 32},
		{200, 100, 50},
	}
	for _, c := range colors {
		cmyk := rgbToCMYK(c)
		back := cmykToRGB(cmyk)
		if abs(int(c.R)-int(back.R)) > 2 || abs(int(c.G)-int(back.G)) > 2 || abs(int(c.B)-int(back.B)) > 2 {
			t.Errorf("CMYK roundtrip failed for %v: CMYK=%v, back=%v", c, cmyk, back)
		}
	}
}

// ── contrast ratio ─────────────────────────────────────────────────────────

func TestContrastRatio_BlackWhite(t *testing.T) {
	black := RGB{0, 0, 0}
	white := RGB{255, 255, 255}

	ratio := contrastRatio(black, white)
	// WCAG defines black/white contrast as 21:1
	if math.Abs(ratio-21.0) > 0.1 {
		t.Errorf("contrastRatio(black, white) = %.2f, want 21.0", ratio)
	}
}

func TestContrastRatio_SameColor(t *testing.T) {
	c := RGB{128, 128, 128}
	ratio := contrastRatio(c, c)
	if math.Abs(ratio-1.0) > 0.01 {
		t.Errorf("contrastRatio(same, same) = %.2f, want 1.0", ratio)
	}
}

func TestContrastRatio_Symmetric(t *testing.T) {
	c1 := RGB{255, 0, 0}
	c2 := RGB{0, 0, 255}
	r1 := contrastRatio(c1, c2)
	r2 := contrastRatio(c2, c1)
	if math.Abs(r1-r2) > 0.001 {
		t.Errorf("contrastRatio should be symmetric: %.4f != %.4f", r1, r2)
	}
}

func TestRelativeLuminance_Black(t *testing.T) {
	l := relativeLuminance(RGB{0, 0, 0})
	if l != 0 {
		t.Errorf("relativeLuminance(black) = %v, want 0", l)
	}
}

func TestRelativeLuminance_White(t *testing.T) {
	l := relativeLuminance(RGB{255, 255, 255})
	if math.Abs(l-1.0) > 0.001 {
		t.Errorf("relativeLuminance(white) = %v, want 1.0", l)
	}
}

// ── color name lookup ──────────────────────────────────────────────────────

func TestClosestColorName_ExactMatch(t *testing.T) {
	tests := []struct {
		rgb      RGB
		wantName string
	}{
		{RGB{255, 0, 0}, "red"},
		{RGB{0, 0, 255}, "blue"},
		{RGB{255, 255, 255}, "white"},
		{RGB{0, 0, 0}, "black"},
	}
	for _, tt := range tests {
		name, _, dist := closestColorName(tt.rgb)
		if name != tt.wantName {
			t.Errorf("closestColorName(%v) = %q, want %q", tt.rgb, name, tt.wantName)
		}
		if dist != 0 {
			t.Errorf("expected exact match (distance 0) for %q, got %.2f", tt.wantName, dist)
		}
	}
}

func TestClosestColorName_Approximate(t *testing.T) {
	// A color very close to red but not exact
	name, _, dist := closestColorName(RGB{254, 1, 1})
	if name != "red" {
		t.Errorf("expected closest name 'red' for (254,1,1), got %q", name)
	}
	if dist == 0 {
		t.Error("expected non-zero distance for approximate match")
	}
}

// ── parseColor ─────────────────────────────────────────────────────────────

func TestParseColor_Hex(t *testing.T) {
	rgb, format, err := parseColor("#ff0000")
	if err != nil {
		t.Fatalf("parseColor(#ff0000): %v", err)
	}
	if format != "hex" {
		t.Errorf("expected format 'hex', got %q", format)
	}
	if rgb != (RGB{255, 0, 0}) {
		t.Errorf("expected RGB{255,0,0}, got %v", rgb)
	}
}

func TestParseColor_RGB(t *testing.T) {
	rgb, format, err := parseColor("rgb(128,64,32)")
	if err != nil {
		t.Fatalf("parseColor: %v", err)
	}
	if format != "rgb" {
		t.Errorf("expected format 'rgb', got %q", format)
	}
	if rgb != (RGB{128, 64, 32}) {
		t.Errorf("expected RGB{128,64,32}, got %v", rgb)
	}
}

func TestParseColor_HSL(t *testing.T) {
	_, format, err := parseColor("hsl(0,100%,50%)")
	if err != nil {
		t.Fatalf("parseColor: %v", err)
	}
	if format != "hsl" {
		t.Errorf("expected format 'hsl', got %q", format)
	}
}

func TestParseColor_CMYK(t *testing.T) {
	_, format, err := parseColor("cmyk(0,100,100,0)")
	if err != nil {
		t.Fatalf("parseColor: %v", err)
	}
	if format != "cmyk" {
		t.Errorf("expected format 'cmyk', got %q", format)
	}
}

func TestParseColor_Invalid(t *testing.T) {
	invalids := []string{"notacolor", "rgb()", "hsl(abc)", "#gggggg"}
	for _, input := range invalids {
		_, _, err := parseColor(input)
		if err == nil {
			t.Errorf("expected error for parseColor(%q)", input)
		}
	}
}

// ── format helpers ─────────────────────────────────────────────────────────

func TestFormatRGB(t *testing.T) {
	got := formatRGB(RGB{255, 128, 0})
	want := "rgb(255,128,0)"
	if got != want {
		t.Errorf("formatRGB = %q, want %q", got, want)
	}
}

func TestFormatHSL(t *testing.T) {
	got := formatHSL(HSL{120, 50, 75})
	want := "hsl(120,50%,75%)"
	if got != want {
		t.Errorf("formatHSL = %q, want %q", got, want)
	}
}

func TestFormatCMYK(t *testing.T) {
	got := formatCMYK(CMYK{0, 100, 100, 0})
	want := "cmyk(0,100,100,0)"
	if got != want {
		t.Errorf("formatCMYK = %q, want %q", got, want)
	}
}

// ── palette generation ─────────────────────────────────────────────────────

func TestGeneratePalette_Count(t *testing.T) {
	base := RGB{255, 87, 51}
	types := []string{"complementary", "analogous", "triadic", "shades", "tints"}
	for _, pt := range types {
		colors, err := generatePalette(base, 5, pt)
		if err != nil {
			t.Fatalf("generatePalette(%q): %v", pt, err)
		}
		if len(colors) != 5 {
			t.Errorf("generatePalette(%q) returned %d colors, want 5", pt, len(colors))
		}
	}
}

func TestGeneratePalette_UnknownType(t *testing.T) {
	_, err := generatePalette(RGB{255, 0, 0}, 5, "invalid")
	if err == nil {
		t.Error("expected error for unknown palette type")
	}
}

// ── allFormats ─────────────────────────────────────────────────────────────

func TestAllFormats_ContainsAllKeys(t *testing.T) {
	m := allFormats(RGB{255, 0, 0})
	for _, key := range []string{"hex", "rgb", "hsl", "cmyk"} {
		if _, ok := m[key]; !ok {
			t.Errorf("allFormats missing key %q", key)
		}
	}
}

// ── passFailStr ────────────────────────────────────────────────────────────

func TestPassFailStr(t *testing.T) {
	if passFailStr(true) != "PASS" {
		t.Error("expected PASS for true")
	}
	if passFailStr(false) != "FAIL" {
		t.Error("expected FAIL for false")
	}
}

// ── helper ─────────────────────────────────────────────────────────────────

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

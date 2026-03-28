package qr

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "qr-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

// ---------------------------------------------------------------------------
// Command metadata
// ---------------------------------------------------------------------------

func TestQRCmd_Metadata(t *testing.T) {
	if qrCmd.Use != "qr <text>" {
		t.Errorf("unexpected Use: %s", qrCmd.Use)
	}
	if qrCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if qrCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestQRCmd_RequiresOneArg(t *testing.T) {
	validator := cobra.ExactArgs(1)
	if err := validator(qrCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(qrCmd, []string{"hello"}); err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}
	if err := validator(qrCmd, []string{"a", "b"}); err == nil {
		t.Error("expected error with two args")
	}
}

func TestQRCmd_Flags(t *testing.T) {
	f := qrCmd.Flags()

	if f.Lookup("output") == nil {
		t.Error("--output flag not found")
	}
	if f.Lookup("size") == nil {
		t.Error("--size flag not found")
	}
	if f.Lookup("level") == nil {
		t.Error("--level flag not found")
	}
	if f.Lookup("invert") == nil {
		t.Error("--invert flag not found")
	}
	if f.ShorthandLookup("o") == nil {
		t.Error("-o shorthand not found")
	}
}

func TestQRCmd_FlagDefaults(t *testing.T) {
	f := qrCmd.Flags()

	if v := f.Lookup("size").DefValue; v != "256" {
		t.Errorf("size default: got %q, want 256", v)
	}
	if v := f.Lookup("level").DefValue; v != "L" {
		t.Errorf("level default: got %q, want L", v)
	}
	if v := f.Lookup("invert").DefValue; v != "false" {
		t.Errorf("invert default: got %q, want false", v)
	}
}

// ---------------------------------------------------------------------------
// parseLevel
// ---------------------------------------------------------------------------

func TestParseLevel(t *testing.T) {
	tests := map[string]qrcode.RecoveryLevel{
		"L": qrcode.Low,
		"l": qrcode.Low,
		"M": qrcode.Medium,
		"m": qrcode.Medium,
		"Q": qrcode.High,
		"H": qrcode.Highest,
		"":  qrcode.Low,
		"X": qrcode.Low,
	}
	for input, want := range tests {
		got := parseLevel(input)
		if got != want {
			t.Errorf("parseLevel(%q) = %v, want %v", input, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// ASCII output
// ---------------------------------------------------------------------------

func TestPrintASCII(t *testing.T) {
	// Just verify it doesn't error
	err := printASCII("test", qrcode.Low, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintASCII_Inverted(t *testing.T) {
	err := printASCII("test", qrcode.Low, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintASCII_SingleChar(t *testing.T) {
	err := printASCII("X", qrcode.Low, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintASCII_LongContent(t *testing.T) {
	long := strings.Repeat("A", 500)
	err := printASCII(long, qrcode.Low, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// PNG output
// ---------------------------------------------------------------------------

func TestSavePNG(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.png")

	err := savePNG("Hello World", path, 256, qrcode.Low)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatal("not a valid PNG")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 256 || bounds.Dy() != 256 {
		t.Errorf("size: got %dx%d, want 256x256", bounds.Dx(), bounds.Dy())
	}
}

func TestSavePNG_CustomSize(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.png")

	err := savePNG("test", path, 512, qrcode.Low)
	if err != nil {
		t.Fatal(err)
	}

	f, _ := os.Open(path)
	defer f.Close()
	img, _ := png.Decode(f)
	if img.Bounds().Dx() != 512 {
		t.Errorf("expected 512px, got %d", img.Bounds().Dx())
	}
}

func TestSavePNG_AllLevels(t *testing.T) {
	dir := tempDir(t)
	levels := []qrcode.RecoveryLevel{qrcode.Low, qrcode.Medium, qrcode.High, qrcode.Highest}
	for _, lvl := range levels {
		path := filepath.Join(dir, "test.png")
		if err := savePNG("test", path, 128, lvl); err != nil {
			t.Errorf("level %v: %v", lvl, err)
		}
	}
}

// ---------------------------------------------------------------------------
// SVG output
// ---------------------------------------------------------------------------

func TestSaveSVG(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.svg")

	err := saveSVG("Hello World", path, qrcode.Low)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	svg := string(data)

	if !strings.Contains(svg, "<?xml") {
		t.Error("missing XML declaration")
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "fill=\"black\"") {
		t.Error("missing black modules")
	}
	if !strings.Contains(svg, "fill=\"white\"") {
		t.Error("missing white background")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing closing </svg>")
	}
}

func TestSaveSVG_URL(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "url.svg")

	err := saveSVG("https://example.com/path?q=1&r=2", path, qrcode.Medium)
	if err != nil {
		t.Fatal(err)
	}

	info, _ := os.Stat(path)
	if info.Size() < 100 {
		t.Error("SVG file too small")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestSavePNG_BadPath(t *testing.T) {
	err := savePNG("test", "/nonexistent/dir/test.png", 256, qrcode.Low)
	if err == nil {
		t.Error("expected error for bad path")
	}
}

func TestSaveSVG_BadPath(t *testing.T) {
	err := saveSVG("test", "/nonexistent/dir/test.svg", qrcode.Low)
	if err == nil {
		t.Error("expected error for bad path")
	}
}

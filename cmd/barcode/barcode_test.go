package barcode

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestBarcodeCmd_Metadata(t *testing.T) {
	if barcodeCmd.Use == "" {
		t.Error("barcodeCmd.Use must not be empty")
	}
	if barcodeCmd.Short == "" {
		t.Error("barcodeCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"code128Cmd", code128Cmd.Use, code128Cmd.Short},
		{"ean13Cmd", ean13Cmd.Use, ean13Cmd.Short},
		{"qrCmd", qrCmd.Use, qrCmd.Short},
		{"datamatrixCmd", datamatrixCmd.Use, datamatrixCmd.Short},
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

func TestBarcodeCmd_PersistentFlags(t *testing.T) {
	f := barcodeCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
	if f.Lookup("output") == nil {
		t.Error("expected persistent flag --output")
	}
	if f.Lookup("width") == nil {
		t.Error("expected persistent flag --width")
	}
	if f.Lookup("height") == nil {
		t.Error("expected persistent flag --height")
	}
}

func TestWidthDefault(t *testing.T) {
	val, err := barcodeCmd.PersistentFlags().GetInt("width")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 300 {
		t.Errorf("expected default width 300, got %d", val)
	}
}

func TestHeightDefault(t *testing.T) {
	val, err := barcodeCmd.PersistentFlags().GetInt("height")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 100 {
		t.Errorf("expected default height 100, got %d", val)
	}
}

// ── EAN-13 input validation ───────────────────────────────────────────────

func TestValidateEAN13Input_Valid12(t *testing.T) {
	if err := validateEAN13Input("590123412345"); err != nil {
		t.Errorf("12-digit input should be valid: %v", err)
	}
}

func TestValidateEAN13Input_Valid13(t *testing.T) {
	if err := validateEAN13Input("5901234123457"); err != nil {
		t.Errorf("13-digit input should be valid: %v", err)
	}
}

func TestValidateEAN13Input_TooShort(t *testing.T) {
	if err := validateEAN13Input("12345"); err == nil {
		t.Error("expected error for too-short input")
	}
}

func TestValidateEAN13Input_TooLong(t *testing.T) {
	if err := validateEAN13Input("12345678901234"); err == nil {
		t.Error("expected error for too-long input")
	}
}

func TestValidateEAN13Input_NonDigit(t *testing.T) {
	if err := validateEAN13Input("59012341234A"); err == nil {
		t.Error("expected error for non-digit character")
	}
}

// ── barcode encoding (library-level) ──────────────────────────────────────

func TestCode128_Encode(t *testing.T) {
	bc, err := code128.Encode("Hello123")
	if err != nil {
		t.Fatalf("code128 encoding failed: %v", err)
	}
	if bc.Bounds().Dx() == 0 || bc.Bounds().Dy() == 0 {
		t.Error("encoded barcode has zero dimensions")
	}
}

func TestEAN13_Encode(t *testing.T) {
	bc, err := ean.Encode("5901234123457")
	if err != nil {
		t.Fatalf("ean13 encoding failed: %v", err)
	}
	if bc.Bounds().Dx() == 0 {
		t.Error("encoded barcode has zero width")
	}
}

func TestQR_Encode(t *testing.T) {
	bc, err := qr.Encode("https://example.com", qr.M, qr.Auto)
	if err != nil {
		t.Fatalf("qr encoding failed: %v", err)
	}
	bounds := bc.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		t.Error("encoded QR has zero dimensions")
	}
}

func TestDataMatrix_Encode(t *testing.T) {
	bc, err := datamatrix.Encode("SERIAL-00123")
	if err != nil {
		t.Fatalf("datamatrix encoding failed: %v", err)
	}
	bounds := bc.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		t.Error("encoded DataMatrix has zero dimensions")
	}
}

// ── savePNG helper ────────────────────────────────────────────────────────

func TestSavePNG_WritesValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if err := savePNG(img, path); err != nil {
		t.Fatalf("savePNG error: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("cannot open saved file: %v", err)
	}
	defer f.Close()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("cannot decode saved PNG: %v", err)
	}
	if decoded.Bounds().Dx() != 10 || decoded.Bounds().Dy() != 10 {
		t.Errorf("unexpected dimensions: %v", decoded.Bounds())
	}
}

func TestSavePNG_BadPath(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	err := savePNG(img, "/nonexistent/dir/test.png")
	if err == nil {
		t.Error("expected error for bad path")
	}
}

// ── runCode128 requires output ─────────────────────────────────────────────

func TestRunCode128_NoOutput(t *testing.T) {
	output = ""
	err := runCode128(nil, []string{"test"})
	if err == nil {
		t.Error("expected error when output is empty")
	}
}

// ── runEAN13 requires output ───────────────────────────────────────────────

func TestRunEAN13_NoOutput(t *testing.T) {
	output = ""
	err := runEAN13(nil, []string{"590123412345"})
	if err == nil {
		t.Error("expected error when output is empty")
	}
}

// ── runQR requires output ──────────────────────────────────────────────────

func TestRunQR_NoOutput(t *testing.T) {
	output = ""
	err := runQR(nil, []string{"hello"})
	if err == nil {
		t.Error("expected error when output is empty")
	}
}

// ── runDataMatrix requires output ──────────────────────────────────────────

func TestRunDataMatrix_NoOutput(t *testing.T) {
	output = ""
	err := runDataMatrix(nil, []string{"hello"})
	if err == nil {
		t.Error("expected error when output is empty")
	}
}

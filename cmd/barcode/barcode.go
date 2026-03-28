package barcode

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── persistent flags (inherited by all subcommands) ─────────────────────────

var (
	jsonOut    bool
	output     string
	barcodeW   int
	barcodeH   int
)

// ── parent command ──────────────────────────────────────────────────────────

var barcodeCmd = &cobra.Command{
	Use:   "barcode",
	Short: "Generate barcodes and 2D codes",
	Long: `Generate various barcode formats as PNG images.

SUBCOMMANDS:

  code128      Generate Code 128 barcode
  ean13        Generate EAN-13 barcode
  qr           Generate QR code (via boombuler/barcode library)
  datamatrix   Generate DataMatrix 2D barcode

FLAGS (inherited by all subcommands):

  --output/-o   Output PNG file path (required)
  --width       Image width in pixels (default 300)
  --height      Image height in pixels (default 100)
  --json/-j     Output result as JSON

Examples:
  openGyver barcode code128 "Hello123" -o barcode.png
  openGyver barcode ean13 "590123412345" -o ean.png
  openGyver barcode qr "https://example.com" -o qr.png --width 256 --height 256
  openGyver barcode datamatrix "DATA" -o dm.png`,
}

// ── subcommands ─────────────────────────────────────────────────────────────

var code128Cmd = &cobra.Command{
	Use:   "code128 <data>",
	Short: "Generate Code 128 barcode",
	Long: `Generate a Code 128 barcode as a PNG image.

Code 128 is a high-density linear barcode that encodes all 128 ASCII
characters. It is widely used in shipping, packaging, and supply chain.

Examples:
  openGyver barcode code128 "Hello World" -o code128.png
  openGyver barcode code128 "ABC-12345" -o code128.png --width 400 --height 120`,
	Args: cobra.ExactArgs(1),
	RunE: runCode128,
}

var ean13Cmd = &cobra.Command{
	Use:   "ean13 <digits>",
	Short: "Generate EAN-13 barcode",
	Long: `Generate an EAN-13 barcode as a PNG image.

EAN-13 (European Article Number) is a 13-digit barcode used worldwide
for marking retail goods. Provide 12 digits (the check digit is computed
automatically) or 13 digits (check digit is validated).

Examples:
  openGyver barcode ean13 "590123412345" -o ean13.png
  openGyver barcode ean13 "5901234123457" -o ean13.png`,
	Args: cobra.ExactArgs(1),
	RunE: runEAN13,
}

var qrCmd = &cobra.Command{
	Use:   "qr <data>",
	Short: "Generate QR code via barcode library",
	Long: `Generate a QR code as a PNG image using the boombuler/barcode library.

This is an alternative QR code generator to the top-level "qr" command.
It uses a different library and always outputs to a PNG file.

For --width and --height, use equal values for a square QR code.

Examples:
  openGyver barcode qr "https://example.com" -o qr.png
  openGyver barcode qr "Hello" -o qr.png --width 512 --height 512`,
	Args: cobra.ExactArgs(1),
	RunE: runQR,
}

var datamatrixCmd = &cobra.Command{
	Use:   "datamatrix <data>",
	Short: "Generate DataMatrix 2D barcode",
	Long: `Generate a DataMatrix 2D barcode as a PNG image.

DataMatrix is a two-dimensional barcode that can store large amounts
of data in a small space. It is commonly used in electronics, healthcare,
and logistics for marking small items.

Examples:
  openGyver barcode datamatrix "SERIAL-00123" -o dm.png
  openGyver barcode datamatrix "Hello World" -o dm.png --width 200 --height 200`,
	Args: cobra.ExactArgs(1),
	RunE: runDataMatrix,
}

// ── runners ─────────────────────────────────────────────────────────────────

func runCode128(_ *cobra.Command, args []string) error {
	data := args[0]
	if output == "" {
		return fmt.Errorf("--output/-o is required")
	}

	bc, err := code128.Encode(data)
	if err != nil {
		return fmt.Errorf("encoding Code 128: %w", err)
	}

	return saveBarcode(bc, data, "code128")
}

func runEAN13(_ *cobra.Command, args []string) error {
	data := args[0]
	if output == "" {
		return fmt.Errorf("--output/-o is required")
	}

	// Validate that input is 12 or 13 digits.
	if err := validateEAN13Input(data); err != nil {
		return err
	}

	bc, err := ean.Encode(data)
	if err != nil {
		return fmt.Errorf("encoding EAN-13: %w", err)
	}

	return saveBarcode(bc, data, "ean13")
}

func runQR(_ *cobra.Command, args []string) error {
	data := args[0]
	if output == "" {
		return fmt.Errorf("--output/-o is required")
	}

	bc, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return fmt.Errorf("encoding QR code: %w", err)
	}

	return saveBarcode(bc, data, "qr")
}

func runDataMatrix(_ *cobra.Command, args []string) error {
	data := args[0]
	if output == "" {
		return fmt.Errorf("--output/-o is required")
	}

	bc, err := datamatrix.Encode(data)
	if err != nil {
		return fmt.Errorf("encoding DataMatrix: %w", err)
	}

	return saveBarcode(bc, data, "datamatrix")
}

// ── helpers ─────────────────────────────────────────────────────────────────

// validateEAN13Input checks that the input is 12 or 13 digits.
func validateEAN13Input(data string) error {
	if len(data) != 12 && len(data) != 13 {
		return fmt.Errorf("EAN-13 requires 12 or 13 digits, got %d characters", len(data))
	}
	for i, ch := range data {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("EAN-13 requires all digits, found %q at position %d", ch, i)
		}
	}
	return nil
}

// saveBarcode scales the barcode to the requested dimensions and writes a PNG file.
func saveBarcode(bc barcode.Barcode, data, format string) error {
	scaled, err := barcode.Scale(bc, barcodeW, barcodeH)
	if err != nil {
		return fmt.Errorf("scaling barcode: %w", err)
	}

	if err := savePNG(scaled, output); err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"format": format,
			"data":   data,
			"output": output,
			"width":  barcodeW,
			"height": barcodeH,
		})
	}

	fmt.Printf("Saved %s barcode to %s (%dx%d)\n", format, output, barcodeW, barcodeH)
	return nil
}

// savePNG encodes an image to a PNG file.
func savePNG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("encoding PNG: %w", err)
	}
	return nil
}

// ── registration ────────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent — inherited by all subcommands.
	barcodeCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output result as JSON")
	barcodeCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output PNG file path (required)")
	barcodeCmd.PersistentFlags().IntVar(&barcodeW, "width", 300, "image width in pixels")
	barcodeCmd.PersistentFlags().IntVar(&barcodeH, "height", 100, "image height in pixels")

	// Wire subcommands.
	barcodeCmd.AddCommand(code128Cmd)
	barcodeCmd.AddCommand(ean13Cmd)
	barcodeCmd.AddCommand(qrCmd)
	barcodeCmd.AddCommand(datamatrixCmd)

	cmd.Register(barcodeCmd)
}

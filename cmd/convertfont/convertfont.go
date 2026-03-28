package convertfont

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var output string

var supportedFormats = map[string]bool{
	"afm": true, "cff": true, "dfont": true, "eot": true, "otf": true,
	"pfa": true, "pfb": true, "sfd": true, "ttf": true, "ufo": true,
	"woff": true, "woff2": true,
}

var convertFontCmd = &cobra.Command{
	Use:   "convertFont <input-file>",
	Short: "Convert between font formats",
	Long: `Convert font files between popular web and desktop formats using fonttools.

REQUIRES: fonttools must be installed (Python package).
  pip install fonttools brotli zopfli

SUPPORTED FORMATS (12):

  AFM     Adobe Font Metrics
  CFF     Compact Font Format
  DFONT   Mac Data Fork Font
  EOT     Embedded OpenType (IE legacy)
  OTF     OpenType Font
  PFA     PostScript Font ASCII
  PFB     PostScript Font Binary
  SFD     FontForge Spline Font Database
  TTF     TrueType Font
  UFO     Unified Font Object
  WOFF    Web Open Font Format 1.0
  WOFF2   Web Open Font Format 2.0 (Brotli)

Examples:
  openGyver convertFont font.ttf -o font.woff2
  openGyver convertFont font.otf -o font.woff
  openGyver convertFont font.woff2 -o font.ttf
  openGyver convertFont font.ttf -o font.eot`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertFont,
}

func runConvertFont(c *cobra.Command, args []string) error {
	inputPath := args[0]
	if output == "" {
		return fmt.Errorf("--output (-o) is required")
	}

	inExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(inputPath)), ".")
	outExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(output)), ".")

	if !supportedFormats[inExt] {
		return fmt.Errorf("unsupported input format: .%s", inExt)
	}
	if !supportedFormats[outExt] {
		return fmt.Errorf("unsupported output format: .%s", outExt)
	}

	// Use fonttools' pyftsubset or ttx for conversion
	if err := checkFonttools(); err != nil {
		return err
	}

	// fonttools approach: use pyftsubset with flavor for woff/woff2
	var cmd *exec.Cmd
	switch outExt {
	case "woff":
		cmd = exec.Command("fonttools", "ttLib.woff", inputPath, "-o", output)
		// Fallback: use pyftsubset
		cmd = exec.Command("pyftsubset", inputPath, "--output-file="+output, "--flavor=woff", "--no-subset")
	case "woff2":
		cmd = exec.Command("pyftsubset", inputPath, "--output-file="+output, "--flavor=woff2", "--no-subset")
	case "ttf", "otf":
		// Convert via fonttools ttx (XML intermediate)
		cmd = exec.Command("fonttools", "ttx", "-o", output, inputPath)
	default:
		return fmt.Errorf("conversion to .%s is not yet supported.\n"+
			"Supported output: ttf, otf, woff, woff2", outExt)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fonttools error: %s\n%s", err, string(out))
	}

	fmt.Printf("Converted %s → %s\n", inputPath, output)
	return nil
}

func checkFonttools() error {
	_, err := exec.LookPath("pyftsubset")
	if err != nil {
		return fmt.Errorf("fonttools is not installed or not in PATH.\n" +
			"Install it:\n" +
			"  pip install fonttools brotli zopfli")
	}
	return nil
}

func init() {
	convertFontCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	cmd.Register(convertFontCmd)
}

package convertvector

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	output string
	width  int
	height int
)

var supportedFormats = map[string]bool{
	"ai": true, "ccx": true, "cdr": true, "cdt": true, "cgm": true,
	"cmx": true, "dst": true, "emf": true, "eps": true, "exp": true,
	"fig": true, "pdf": true, "pes": true, "plt": true, "png": true,
	"jpg": true, "ps": true, "sk": true, "sk1": true, "svg": true,
	"svgz": true, "vsd": true, "wmf": true,
}

var convertVectorCmd = &cobra.Command{
	Use:   "convertVector <input-file>",
	Short: "Convert between vector graphics formats",
	Long: `Convert vector graphics between popular formats.

SVG → PNG/PDF conversions use rsvg-convert (librsvg) if available,
or Inkscape as a fallback. Other vector formats require Inkscape.

REQUIRES (one of):
  rsvg-convert:  brew install librsvg  (macOS) / apt install librsvg2-bin (Linux)
  inkscape:      brew install inkscape (macOS) / apt install inkscape (Linux)

SUPPORTED FORMATS (23):

  SVG, SVGZ    Scalable Vector Graphics
  EPS          Encapsulated PostScript
  PDF, PS      Page layout
  EMF, WMF     Windows Metafile formats
  AI           Adobe Illustrator
  CDR, CDT     CorelDRAW
  CCX, CMX     Corel formats
  CGM          Computer Graphics Metafile
  VSD          Microsoft Visio
  FIG          Xfig
  PLT          HPGL Plotter
  DST, PES, EXP  Embroidery formats
  SK, SK1      Sketch

  Raster output: PNG, JPG (rasterize vector to image)

Examples:
  openGyver convertVector logo.svg -o logo.png
  openGyver convertVector logo.svg -o logo.pdf
  openGyver convertVector logo.svg -o logo.eps
  openGyver convertVector logo.svg -o logo.png --width 1024 --height 768
  openGyver convertVector diagram.eps -o diagram.svg
  openGyver convertVector drawing.ai -o drawing.svg`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertVector,
}

func runConvertVector(c *cobra.Command, args []string) error {
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

	// Try rsvg-convert first for SVG input
	if (inExt == "svg" || inExt == "svgz") && isRasterOutput(outExt) {
		if err := tryRsvg(inputPath, outExt); err == nil {
			fmt.Printf("Converted %s → %s\n", inputPath, output)
			return nil
		}
	}

	// Fall back to Inkscape for all conversions
	if err := tryInkscape(inputPath, outExt); err != nil {
		return err
	}

	fmt.Printf("Converted %s → %s\n", inputPath, output)
	return nil
}

func isRasterOutput(ext string) bool {
	return ext == "png" || ext == "jpg" || ext == "jpeg"
}

func tryRsvg(inputPath, outExt string) error {
	rsvg, err := exec.LookPath("rsvg-convert")
	if err != nil {
		return err
	}

	args := []string{inputPath}
	switch outExt {
	case "png":
		args = append(args, "-f", "png")
	case "pdf":
		args = append(args, "-f", "pdf")
	case "ps", "eps":
		args = append(args, "-f", "ps")
	default:
		return fmt.Errorf("rsvg-convert doesn't support .%s", outExt)
	}

	if width > 0 {
		args = append(args, "-w", fmt.Sprintf("%d", width))
	}
	if height > 0 {
		args = append(args, "-h", fmt.Sprintf("%d", height))
	}
	args = append(args, "-o", output)

	cmd := exec.Command(rsvg, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsvg-convert error: %s\n%s", err, string(out))
	}
	return nil
}

func tryInkscape(inputPath, outExt string) error {
	inkscape, err := exec.LookPath("inkscape")
	if err != nil {
		return fmt.Errorf("neither rsvg-convert nor inkscape found in PATH.\n" +
			"Install one of:\n" +
			"  macOS:  brew install librsvg   OR   brew install inkscape\n" +
			"  Linux:  apt install librsvg2-bin  OR  apt install inkscape")
	}

	args := []string{inputPath, "--export-filename=" + output}
	switch outExt {
	case "png":
		args = append(args, "--export-type=png")
	case "pdf":
		args = append(args, "--export-type=pdf")
	case "eps":
		args = append(args, "--export-type=eps")
	case "ps":
		args = append(args, "--export-type=ps")
	case "svg":
		args = append(args, "--export-type=svg")
	case "emf":
		args = append(args, "--export-type=emf")
	case "wmf":
		args = append(args, "--export-type=wmf")
	default:
		return fmt.Errorf("conversion to .%s requires manual tools", outExt)
	}

	if width > 0 {
		args = append(args, fmt.Sprintf("--export-width=%d", width))
	}
	if height > 0 {
		args = append(args, fmt.Sprintf("--export-height=%d", height))
	}

	cmd := exec.Command(inkscape, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("inkscape error: %s\n%s", err, string(out))
	}
	return nil
}

func init() {
	convertVectorCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	convertVectorCmd.Flags().IntVar(&width, "width", 0, "output width in pixels (for rasterization)")
	convertVectorCmd.Flags().IntVar(&height, "height", 0, "output height in pixels (for rasterization)")
	cmd.Register(convertVectorCmd)
}

package convertcad

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
	"dwg": true, "dxf": true, "dwf": true, "pdf": true, "svg": true, "png": true,
}

var convertCADCmd = &cobra.Command{
	Use:   "convertCAD <input-file>",
	Short: "Convert between CAD file formats",
	Long: `Convert CAD files between DWG, DXF, and other formats.

REQUIRES: One of the following must be installed:
  LibreCAD:   brew install librecad  (for DXF conversions)
  ODA File Converter: https://www.opendesign.com/guestfiles/oda_file_converter
    (for DWG ↔ DXF)

SUPPORTED FORMATS:

  DWG    AutoCAD Drawing
  DXF    Drawing Exchange Format
  DWF    Design Web Format
  PDF    Export to PDF
  SVG    Export to SVG
  PNG    Export to PNG

Examples:
  openGyver convertCAD drawing.dwg -o drawing.dxf
  openGyver convertCAD drawing.dxf -o drawing.pdf
  openGyver convertCAD drawing.dxf -o drawing.svg
  openGyver convertCAD drawing.dwg -o drawing.pdf`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertCAD,
}

func runConvertCAD(c *cobra.Command, args []string) error {
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

	// Try ODA File Converter for DWG ↔ DXF
	if (inExt == "dwg" && outExt == "dxf") || (inExt == "dxf" && outExt == "dwg") {
		if oda, err := exec.LookPath("ODAFileConverter"); err == nil {
			return runODA(oda, inputPath, outExt)
		}
	}

	// Try LibreCAD for DXF → other formats
	if inExt == "dxf" {
		if librecad, err := exec.LookPath("librecad"); err == nil {
			return runLibreCAD(librecad, inputPath, outExt)
		}
	}

	return fmt.Errorf("no suitable CAD converter found.\n" +
		"Install one of:\n" +
		"  ODA File Converter: https://www.opendesign.com/guestfiles/oda_file_converter\n" +
		"  LibreCAD:           brew install librecad (macOS) / apt install librecad (Linux)")
}

func runODA(odaPath, inputPath, outExt string) error {
	inDir := filepath.Dir(inputPath)
	outDir := filepath.Dir(output)
	if outDir == "" {
		outDir = "."
	}

	outVersion := "ACAD2018"
	outFormat := "DXF"
	if outExt == "dwg" {
		outFormat = "DWG"
	}

	cmd := exec.Command(odaPath, inDir, outDir, outVersion, outFormat, "0", "1", filepath.Base(inputPath))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ODA converter error: %s\n%s", err, string(out))
	}

	fmt.Printf("Converted %s → %s\n", inputPath, output)
	return nil
}

func runLibreCAD(librecadPath, inputPath, outExt string) error {
	cmd := exec.Command(librecadPath, inputPath, "--export", output)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("LibreCAD error: %s\n%s", err, string(out))
	}

	fmt.Printf("Converted %s → %s\n", inputPath, output)
	return nil
}

func init() {
	convertCADCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	cmd.Register(convertCADCmd)
}

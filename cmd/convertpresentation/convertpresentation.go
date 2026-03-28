package convertpresentation

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
	"dps": true, "key": true, "odp": true, "pot": true, "potm": true,
	"potx": true, "pps": true, "ppsm": true, "ppsx": true, "ppt": true,
	"pptm": true, "pptx": true, "pdf": true, "png": true, "jpg": true,
	"svg": true, "html": true,
}

var convertPresentationCmd = &cobra.Command{
	Use:   "convertPresentation <input-file>",
	Short: "Convert between presentation formats",
	Long: `Convert presentation files between popular formats using LibreOffice.

REQUIRES: LibreOffice must be installed.
  macOS:   brew install --cask libreoffice
  Linux:   apt install libreoffice
  Windows: https://www.libreoffice.org/download

SUPPORTED FORMATS:

  Input:   DPS, KEY, ODP, POT, POTX, PPS, PPSX, PPT, PPTM, PPTX
  Output:  ODP, PDF, PPTX, PPT, HTML, PNG, JPG, SVG

Examples:
  openGyver convertPresentation slides.pptx -o slides.pdf
  openGyver convertPresentation slides.pptx -o slides.odp
  openGyver convertPresentation slides.key -o slides.pptx
  openGyver convertPresentation slides.pptx -o slides.html
  openGyver convertPresentation old.ppt -o new.pptx`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertPresentation,
}

func runConvertPresentation(c *cobra.Command, args []string) error {
	if err := checkLibreOffice(); err != nil {
		return err
	}

	inputPath := args[0]
	if output == "" {
		return fmt.Errorf("--output (-o) is required")
	}

	outExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(output)), ".")
	outDir := filepath.Dir(output)
	if outDir == "" {
		outDir = "."
	}

	// Map extensions to LibreOffice filter names
	filterMap := map[string]string{
		"pdf":  "pdf",
		"pptx": "pptx",
		"ppt":  "ppt",
		"odp":  "odp",
		"html": "html",
		"png":  "png",
		"jpg":  "jpg",
		"svg":  "svg",
	}

	filter, ok := filterMap[outExt]
	if !ok {
		return fmt.Errorf("unsupported output format: .%s\nSupported: pdf, pptx, ppt, odp, html, png, jpg, svg", outExt)
	}

	loCmd := exec.Command("soffice",
		"--headless",
		"--convert-to", filter,
		"--outdir", outDir,
		inputPath,
	)
	out, err := loCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("libreoffice error: %s\n%s", err, string(out))
	}

	// LibreOffice outputs to outDir with the same base name
	expectedName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) + "." + outExt
	expectedPath := filepath.Join(outDir, expectedName)

	// Rename if the user specified a different output name
	if expectedPath != output {
		if err := exec.Command("mv", expectedPath, output).Run(); err != nil {
			// Non-fatal: file was created but couldn't be renamed
			fmt.Printf("Converted %s → %s (note: rename to %s failed)\n", inputPath, expectedPath, output)
			return nil
		}
	}

	fmt.Printf("Converted %s → %s\n", inputPath, output)
	return nil
}

func checkLibreOffice() error {
	_, err := exec.LookPath("soffice")
	if err != nil {
		return fmt.Errorf("LibreOffice is not installed or not in PATH.\n" +
			"Install it:\n" +
			"  macOS:   brew install --cask libreoffice\n" +
			"  Linux:   apt install libreoffice\n" +
			"  Windows: https://www.libreoffice.org/download")
	}
	return nil
}

func init() {
	convertPresentationCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	cmd.Register(convertPresentationCmd)
}

package convertebook

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
	"azw": true, "azw3": true, "azw4": true, "cbc": true, "cbr": true,
	"cbz": true, "chm": true, "epub": true, "fb2": true, "htm": true,
	"html": true, "htmlz": true, "lit": true, "lrf": true, "mobi": true,
	"pdb": true, "pdf": true, "pml": true, "prc": true, "rb": true,
	"snb": true, "tcr": true, "txt": true, "txtz": true, "docx": true,
}

var convertEbookCmd = &cobra.Command{
	Use:   "convertEbook <input-file>",
	Short: "Convert between ebook formats",
	Long: `Convert ebook files between popular formats using Calibre's ebook-convert.

REQUIRES: Calibre must be installed.
  macOS:   brew install calibre
  Linux:   apt install calibre
  Windows: https://calibre-ebook.com/download
  All:     https://calibre-ebook.com

SUPPORTED FORMATS:

  Read:   AZW, AZW3, AZW4, CBC, CBR, CBZ, CHM, DOCX, EPUB, FB2,
          HTM, HTML, HTMLZ, LIT, LRF, MOBI, PDB, PDF, PML, PRC,
          RB, SNB, TCR, TXT, TXTZ
  Write:  AZW3, DOCX, EPUB, FB2, HTM, HTMLZ, LRF, MOBI, PDB,
          PDF, PML, PRC, RB, SNB, TCR, TXT, TXTZ

Examples:
  openGyver convertEbook book.epub -o book.mobi
  openGyver convertEbook book.epub -o book.pdf
  openGyver convertEbook book.mobi -o book.epub
  openGyver convertEbook document.docx -o document.epub
  openGyver convertEbook book.epub -o book.azw3
  openGyver convertEbook book.fb2 -o book.epub`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertEbook,
}

func runConvertEbook(c *cobra.Command, args []string) error {
	if err := checkCalibre(); err != nil {
		return err
	}

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

	cmd := exec.Command("ebook-convert", inputPath, output)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ebook-convert error: %s\n%s", err, string(out))
	}

	fmt.Printf("Converted %s → %s\n", inputPath, output)
	return nil
}

func checkCalibre() error {
	_, err := exec.LookPath("ebook-convert")
	if err != nil {
		return fmt.Errorf("Calibre's ebook-convert is not installed or not in PATH.\n" +
			"Install Calibre:\n" +
			"  macOS:   brew install calibre\n" +
			"  Linux:   apt install calibre\n" +
			"  Windows: https://calibre-ebook.com/download")
	}
	return nil
}

func init() {
	convertEbookCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	cmd.Register(convertEbookCmd)
}

package convertfile

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	output    string
	sheet     string
	delimiter string
	quiet     bool
	jsonOut   bool
)

var convertFileCmd = &cobra.Command{
	Use:   "convertFile <input-file>",
	Short: "Convert between document, spreadsheet, and tabular file formats",
	Long: `Convert files between spreadsheet, document, and page-layout formats.

The input format is auto-detected from the file extension. The output format
is determined by the --output file extension, or defaults to a sensible
counterpart.

SUPPORTED FORMATS:

  Spreadsheet    .csv, .tsv, .xlsx, .numbers
  Document       .txt, .md (Markdown), .html, .docx
  Page layout    .pdf, .ps (PostScript)

SPREADSHEET CONVERSIONS:

  CSV  ↔ XLSX         Fully implemented
  CSV  → Numbers      Stub (workaround provided)
  Numbers → CSV/XLSX  Stub (workaround provided)

DOCUMENT CONVERSIONS:

  Markdown → HTML     Rendered with full CommonMark support
  HTML → Markdown     Converted back to clean Markdown
  Markdown → Text     Stripped of all formatting
  HTML → Text         Tags stripped, text extracted
  DOCX → Text         Text extracted from Word XML
  DOCX → HTML         Paragraphs converted to HTML
  DOCX → Markdown     Paragraphs converted to Markdown
  Text → HTML         Wrapped in <pre> block
  Text → Markdown     Wrapped in code fence
  Text → DOCX         Plain paragraphs in a Word document
  Markdown → DOCX     Converted via HTML to Word
  HTML → DOCX         Paragraphs extracted into Word document

PDF / POSTSCRIPT OUTPUT:

  Text → PDF          Plain text rendered to PDF
  Markdown → PDF      Rendered via HTML to a formatted PDF
  HTML → PDF          Parsed and rendered to PDF
  CSV → PDF           Tabular layout with headers
  XLSX → PDF          Tabular layout with headers
  DOCX → PDF          Text extracted and rendered to PDF
  Text → PS           PostScript text output
  Markdown → PS       Rendered to PostScript
  CSV → PS            Tabular PostScript output

CSV OPTIONS:

  Delimiter is auto-detected (comma, tab, semicolon, pipe).
  Use --delimiter to override.

XLSX OPTIONS:

  Multi-sheet files export the first sheet by default.
  Use --sheet to specify a sheet name.

Examples:
  openGyver convertFile data.csv -o report.xlsx
  openGyver convertFile report.xlsx -o report.pdf
  openGyver convertFile README.md -o README.html
  openGyver convertFile README.md -o README.pdf
  openGyver convertFile page.html -o page.md
  openGyver convertFile page.html -o page.pdf
  openGyver convertFile notes.txt -o notes.pdf
  openGyver convertFile notes.txt -o notes.docx
  openGyver convertFile report.docx -o report.pdf
  openGyver convertFile report.docx -o report.txt
  openGyver convertFile data.csv -o data.ps`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertFile,
}

func runConvertFile(c *cobra.Command, args []string) error {
	inputPath := args[0]
	inputExt := strings.ToLower(filepath.Ext(inputPath))

	outputPath := output
	if outputPath == "" {
		outputPath = defaultOutput(inputPath, inputExt)
	}
	outputExt := strings.ToLower(filepath.Ext(outputPath))

	inputFmt := extToFormat(inputExt)
	outputFmt := extToFormat(outputExt)

	if inputFmt == "" {
		return fmt.Errorf("unsupported input format: %s\nRun 'openGyver convertFile --help' for supported formats", inputExt)
	}
	if outputFmt == "" {
		return fmt.Errorf("unsupported output format: %s\nRun 'openGyver convertFile --help' for supported formats", outputExt)
	}
	if inputFmt == outputFmt {
		return fmt.Errorf("input and output formats are the same (%s). Use different extensions.", inputFmt)
	}

	key := inputFmt + "→" + outputFmt
	converter, ok := converters[key]
	if !ok {
		return fmt.Errorf("conversion %s → %s is not supported\nRun 'openGyver convertFile --help' for supported conversions", inputFmt, outputFmt)
	}

	delim := resolveDelimiter(delimiter)

	opts := ConvertOpts{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Sheet:      sheet,
		Delimiter:  delim,
	}

	if err := converter(opts); err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"success":       true,
			"input":         inputPath,
			"output":        outputPath,
			"input_format":  inputFmt,
			"output_format": outputFmt,
		})
	}
	if !quiet {
		fmt.Printf("Converted %s → %s\n", inputPath, outputPath)
	}
	return nil
}

// ConvertOpts holds options passed to converter functions.
type ConvertOpts struct {
	InputPath  string
	OutputPath string
	Sheet      string
	Delimiter  rune
}

// converters maps "inputFormat→outputFormat" to a conversion function.
var converters = map[string]func(ConvertOpts) error{
	// Spreadsheet
	"csv→xlsx":     csvToXLSX,
	"xlsx→csv":     xlsxToCSV,
	"csv→numbers":  csvToNumbers,
	"xlsx→numbers": xlsxToNumbers,
	"numbers→csv":  numbersToCSV,
	"numbers→xlsx": numbersToXLSX,

	// Document ↔ Document
	"md→html":   mdToHTML,
	"html→md":   htmlToMD,
	"md→txt":    mdToText,
	"html→txt":  htmlToText,
	"txt→html":  textToHTML,
	"txt→md":    textToMD,
	"docx→txt":  docxToText,
	"docx→html": docxToHTML,
	"docx→md":   docxToMD,
	"txt→docx":  textToDOCX,
	"md→docx":   mdToDOCX,
	"html→docx": htmlToDOCX,

	// → PDF
	"txt→pdf":  textToPDF,
	"md→pdf":   mdToPDF,
	"html→pdf": htmlToPDF,
	"csv→pdf":  csvToPDF,
	"xlsx→pdf": xlsxToPDF,
	"docx→pdf": docxToPDF,

	// → PostScript
	"txt→ps": textToPS,
	"md→ps":  mdToPS,
	"csv→ps": csvToPS,
}

func extToFormat(ext string) string {
	switch ext {
	case ".csv", ".tsv":
		return "csv"
	case ".xlsx":
		return "xlsx"
	case ".numbers":
		return "numbers"
	case ".txt", ".text":
		return "txt"
	case ".md", ".markdown":
		return "md"
	case ".html", ".htm":
		return "html"
	case ".docx":
		return "docx"
	case ".pdf":
		return "pdf"
	case ".ps":
		return "ps"
	default:
		return ""
	}
}

func defaultOutput(inputPath, inputExt string) string {
	base := strings.TrimSuffix(inputPath, inputExt)
	switch extToFormat(inputExt) {
	case "csv":
		return base + ".xlsx"
	case "xlsx":
		return base + ".csv"
	case "numbers":
		return base + ".csv"
	case "txt":
		return base + ".pdf"
	case "md":
		return base + ".html"
	case "html":
		return base + ".md"
	case "docx":
		return base + ".pdf"
	default:
		return base + ".pdf"
	}
}

func resolveDelimiter(s string) rune {
	switch s {
	case "tab", "\\t", "\t":
		return '\t'
	case "semicolon", ";":
		return ';'
	case "pipe", "|":
		return '|'
	case "", "comma", ",":
		return ','
	default:
		if len(s) > 0 {
			return rune(s[0])
		}
		return ','
	}
}

func init() {
	convertFileCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (default: input name with new extension)")
	convertFileCmd.Flags().StringVar(&sheet, "sheet", "", "sheet name for XLSX/Numbers files (default: first sheet)")
	convertFileCmd.Flags().StringVar(&delimiter, "delimiter", "", "CSV delimiter: comma, tab, semicolon, pipe, or any single character (default: comma)")
	convertFileCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output messages (for piping)")
	convertFileCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output result as JSON")
	cmd.Register(convertFileCmd)
}

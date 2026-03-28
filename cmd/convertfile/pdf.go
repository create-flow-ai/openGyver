package convertfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"
	"github.com/xuri/excelize/v2"
)

const (
	pdfMargin   = 15.0
	pdfFontSize = 10.0
	pdfLineHt   = 5.0
)

// textToPDF renders plain text to a PDF document.
func textToPDF(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading text: %w", err)
	}
	return writeTextPDF(opts.OutputPath, string(src))
}

// mdToPDF converts Markdown to PDF via HTML rendering.
func mdToPDF(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading markdown: %w", err)
	}

	html, err := renderMarkdownToHTML(src)
	if err != nil {
		return fmt.Errorf("rendering markdown: %w", err)
	}

	text := stripHTML(html)
	return writeTextPDF(opts.OutputPath, text)
}

// htmlToPDF converts HTML to PDF by extracting text.
func htmlToPDF(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading HTML: %w", err)
	}

	text := stripHTML(string(src))
	return writeTextPDF(opts.OutputPath, text)
}

// docxToPDF converts DOCX to PDF.
func docxToPDF(opts ConvertOpts) error {
	paragraphs, err := readDOCXParagraphs(opts.InputPath)
	if err != nil {
		return err
	}
	text := strings.Join(paragraphs, "\n")
	return writeTextPDF(opts.OutputPath, text)
}

// csvToPDF renders CSV data as a table in PDF.
func csvToPDF(opts ConvertOpts) error {
	records, err := readCSV(opts.InputPath, opts.Delimiter)
	if err != nil {
		return fmt.Errorf("reading CSV: %w", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}
	return writeTablePDF(opts.OutputPath, records)
}

// xlsxToPDF renders an XLSX sheet as a table in PDF.
func xlsxToPDF(opts ConvertOpts) error {
	f, err := excelize.OpenFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("opening XLSX: %w", err)
	}
	defer f.Close()

	sheetName := opts.Sheet
	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("reading sheet %q: %w", sheetName, err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("sheet %q is empty", sheetName)
	}
	return writeTablePDF(opts.OutputPath, rows)
}

// --- PDF rendering helpers ---

func writeTextPDF(path, text string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pdfMargin, pdfMargin, pdfMargin)
	pdf.SetAutoPageBreak(true, pdfMargin)
	pdf.AddPage()
	pdf.SetFont("Courier", "", pdfFontSize)

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			pdf.Ln(pdfLineHt)
			continue
		}
		pdf.MultiCell(0, pdfLineHt, line, "", "L", false)
	}

	return pdf.OutputFileAndClose(path)
}

func writeTablePDF(path string, records [][]string) error {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(pdfMargin, pdfMargin, pdfMargin)
	pdf.SetAutoPageBreak(true, pdfMargin)
	pdf.AddPage()

	pageWidth, _ := pdf.GetPageSize()
	usable := pageWidth - 2*pdfMargin

	// Determine column count and widths
	maxCols := 0
	for _, row := range records {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	if maxCols == 0 {
		return fmt.Errorf("no columns in data")
	}
	colWidth := usable / float64(maxCols)
	if colWidth < 15 {
		colWidth = 15
	}

	cellHeight := 6.0

	// Header row (first row, bold)
	if len(records) > 0 {
		pdf.SetFont("Helvetica", "B", 9)
		for _, cell := range records[0] {
			pdf.CellFormat(colWidth, cellHeight, truncate(cell, colWidth), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Data rows
	pdf.SetFont("Helvetica", "", 8)
	for i := 1; i < len(records); i++ {
		for j := 0; j < maxCols; j++ {
			cell := ""
			if j < len(records[i]) {
				cell = records[i][j]
			}
			pdf.CellFormat(colWidth, cellHeight, truncate(cell, colWidth), "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	return pdf.OutputFileAndClose(path)
}

func truncate(s string, colWidth float64) string {
	maxChars := int(colWidth / 2)
	if maxChars < 5 {
		maxChars = 5
	}
	if len(s) > maxChars {
		return s[:maxChars-1] + "~"
	}
	return s
}

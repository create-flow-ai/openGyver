package convertfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

// ---------------------------------------------------------------------------
// Text → PDF
// ---------------------------------------------------------------------------

func TestTextToPDF(t *testing.T) {
	dir := tempDir(t)
	txtPath := filepath.Join(dir, "test.txt")
	pdfPath := filepath.Join(dir, "test.pdf")

	os.WriteFile(txtPath, []byte("Hello World\nLine two\n"), 0644)

	err := textToPDF(ConvertOpts{InputPath: txtPath, OutputPath: pdfPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(pdfPath)
	if len(data) < 100 {
		t.Error("PDF file too small")
	}
	if string(data[:4]) != "%PDF" {
		t.Error("not a valid PDF file")
	}
}

// ---------------------------------------------------------------------------
// Markdown → PDF
// ---------------------------------------------------------------------------

func TestMDToPDF(t *testing.T) {
	dir := tempDir(t)
	mdPath := filepath.Join(dir, "test.md")
	pdfPath := filepath.Join(dir, "test.pdf")

	os.WriteFile(mdPath, []byte("# Title\n\nSome **text** here.\n"), 0644)

	err := mdToPDF(ConvertOpts{InputPath: mdPath, OutputPath: pdfPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(pdfPath)
	if string(data[:4]) != "%PDF" {
		t.Error("not a valid PDF")
	}
}

// ---------------------------------------------------------------------------
// HTML → PDF
// ---------------------------------------------------------------------------

func TestHTMLToPDF(t *testing.T) {
	dir := tempDir(t)
	htmlPath := filepath.Join(dir, "test.html")
	pdfPath := filepath.Join(dir, "test.pdf")

	os.WriteFile(htmlPath, []byte("<html><body><h1>Title</h1><p>Content</p></body></html>"), 0644)

	err := htmlToPDF(ConvertOpts{InputPath: htmlPath, OutputPath: pdfPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(pdfPath)
	if string(data[:4]) != "%PDF" {
		t.Error("not a valid PDF")
	}
}

// ---------------------------------------------------------------------------
// CSV → PDF
// ---------------------------------------------------------------------------

func TestCSVToPDF(t *testing.T) {
	dir := tempDir(t)
	csvPath := filepath.Join(dir, "test.csv")
	pdfPath := filepath.Join(dir, "test.pdf")

	writeCSV(t, csvPath, ',', sampleRows)

	err := csvToPDF(ConvertOpts{InputPath: csvPath, OutputPath: pdfPath, Delimiter: ','})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(pdfPath)
	if string(data[:4]) != "%PDF" {
		t.Error("not a valid PDF")
	}
}

// ---------------------------------------------------------------------------
// XLSX → PDF
// ---------------------------------------------------------------------------

func TestXLSXToPDF(t *testing.T) {
	dir := tempDir(t)
	xlsxPath := filepath.Join(dir, "test.xlsx")
	pdfPath := filepath.Join(dir, "test.pdf")

	f := excelize.NewFile()
	for i, row := range sampleRows {
		for j, cell := range row {
			col, _ := excelize.ColumnNumberToName(j + 1)
			f.SetCellValue("Sheet1", col+string(rune('1'+i)), cell)
		}
	}
	f.SaveAs(xlsxPath)
	f.Close()

	err := xlsxToPDF(ConvertOpts{InputPath: xlsxPath, OutputPath: pdfPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(pdfPath)
	if string(data[:4]) != "%PDF" {
		t.Error("not a valid PDF")
	}
}

// ---------------------------------------------------------------------------
// DOCX → PDF
// ---------------------------------------------------------------------------

func TestDOCXToPDF(t *testing.T) {
	dir := tempDir(t)
	docxPath := filepath.Join(dir, "test.docx")
	pdfPath := filepath.Join(dir, "test.pdf")

	writeDOCX(docxPath, []string{"Hello World", "Second paragraph"})

	err := docxToPDF(ConvertOpts{InputPath: docxPath, OutputPath: pdfPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(pdfPath)
	if string(data[:4]) != "%PDF" {
		t.Error("not a valid PDF")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestTextToPDF_MissingFile(t *testing.T) {
	err := textToPDF(ConvertOpts{InputPath: "/nonexistent.txt", OutputPath: "/tmp/out.pdf"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCSVToPDF_EmptyFile(t *testing.T) {
	dir := tempDir(t)
	csvPath := filepath.Join(dir, "empty.csv")
	pdfPath := filepath.Join(dir, "out.pdf")

	os.WriteFile(csvPath, []byte(""), 0644)

	err := csvToPDF(ConvertOpts{InputPath: csvPath, OutputPath: pdfPath, Delimiter: ','})
	if err == nil {
		t.Error("expected error for empty CSV")
	}
}

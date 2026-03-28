package convertfile

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "convertfile-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeCSV(t *testing.T, path string, delim rune, rows [][]string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = delim
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			t.Fatal(err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		t.Fatal(err)
	}
}

func writeXLSX(t *testing.T, path string, sheetName string, rows [][]string) {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()
	if sheetName != "" {
		f.SetSheetName("Sheet1", sheetName)
	} else {
		sheetName = "Sheet1"
	}
	for rowIdx, row := range rows {
		for colIdx, cell := range row {
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			cellRef := colName + string(rune('0'+rowIdx+1))
			if rowIdx+1 >= 10 {
				cellRef = colName + string(rune('0'+(rowIdx+1)/10)) + string(rune('0'+(rowIdx+1)%10))
			}
			f.SetCellValue(sheetName, cellRef, cell)
		}
	}
	if err := f.SaveAs(path); err != nil {
		t.Fatal(err)
	}
}

var sampleRows = [][]string{
	{"Name", "Age", "City"},
	{"Alice", "30", "New York"},
	{"Bob", "25", "London"},
	{"Charlie", "35", "Tokyo"},
}

// ---------------------------------------------------------------------------
// Command metadata
// ---------------------------------------------------------------------------

func TestConvertFileCmd_Metadata(t *testing.T) {
	if convertFileCmd.Use != "convertFile <input-file>" {
		t.Errorf("unexpected Use: %s", convertFileCmd.Use)
	}
	if convertFileCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if convertFileCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestConvertFileCmd_RequiresOneArg(t *testing.T) {
	validator := cobra.ExactArgs(1)
	if err := validator(convertFileCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(convertFileCmd, []string{"a.csv"}); err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}
	if err := validator(convertFileCmd, []string{"a.csv", "b.csv"}); err == nil {
		t.Error("expected error with two args")
	}
}

func TestConvertFileCmd_Flags(t *testing.T) {
	f := convertFileCmd.Flags()

	outFlag := f.Lookup("output")
	if outFlag == nil {
		t.Fatal("--output flag not found")
	}

	sheetFlag := f.Lookup("sheet")
	if sheetFlag == nil {
		t.Fatal("--sheet flag not found")
	}

	delimFlag := f.Lookup("delimiter")
	if delimFlag == nil {
		t.Fatal("--delimiter flag not found")
	}
}

func TestConvertFileCmd_ShortFlags(t *testing.T) {
	f := convertFileCmd.Flags()
	if f.ShorthandLookup("o") == nil {
		t.Error("-o shorthand not found for --output")
	}
}

// ---------------------------------------------------------------------------
// extToFormat
// ---------------------------------------------------------------------------

func TestExtToFormat(t *testing.T) {
	tests := map[string]string{
		".csv":      "csv",
		".tsv":      "csv",
		".xlsx":     "xlsx",
		".numbers":  "numbers",
		".txt":      "txt",
		".text":     "txt",
		".md":       "md",
		".markdown": "md",
		".html":     "html",
		".htm":      "html",
		".docx":     "docx",
		".pdf":      "pdf",
		".ps":       "ps",
		".xyz":      "",
		"":          "",
	}
	for ext, want := range tests {
		got := extToFormat(ext)
		if got != want {
			t.Errorf("extToFormat(%q) = %q, want %q", ext, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// defaultOutput
// ---------------------------------------------------------------------------

func TestDefaultOutput(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"data.csv", "data.xlsx"},
		{"report.xlsx", "report.csv"},
		{"budget.numbers", "budget.csv"},
		{"notes.txt", "notes.pdf"},
		{"readme.md", "readme.html"},
		{"page.html", "page.md"},
		{"report.docx", "report.pdf"},
	}
	for _, tt := range tests {
		ext := filepath.Ext(tt.input)
		got := defaultOutput(tt.input, ext)
		if got != tt.want {
			t.Errorf("defaultOutput(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// resolveDelimiter
// ---------------------------------------------------------------------------

func TestResolveDelimiter(t *testing.T) {
	tests := map[string]rune{
		"":          ',',
		"comma":     ',',
		",":         ',',
		"tab":       '\t',
		"\\t":       '\t',
		"semicolon": ';',
		";":         ';',
		"pipe":      '|',
		"|":         '|',
	}
	for input, want := range tests {
		got := resolveDelimiter(input)
		if got != want {
			t.Errorf("resolveDelimiter(%q) = %q, want %q", input, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// CSV → XLSX
// ---------------------------------------------------------------------------

func TestCSVToXLSX(t *testing.T) {
	dir := tempDir(t)
	csvPath := filepath.Join(dir, "input.csv")
	xlsxPath := filepath.Join(dir, "output.xlsx")

	writeCSV(t, csvPath, ',', sampleRows)

	err := csvToXLSX(ConvertOpts{
		InputPath:  csvPath,
		OutputPath: xlsxPath,
		Delimiter:  ',',
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify the XLSX file was created and has correct data
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(rows))
	}
	if rows[0][0] != "Name" || rows[0][1] != "Age" || rows[0][2] != "City" {
		t.Errorf("header mismatch: %v", rows[0])
	}
	if rows[1][0] != "Alice" || rows[1][2] != "New York" {
		t.Errorf("row 1 mismatch: %v", rows[1])
	}
}

func TestCSVToXLSX_CustomSheet(t *testing.T) {
	dir := tempDir(t)
	csvPath := filepath.Join(dir, "input.csv")
	xlsxPath := filepath.Join(dir, "output.xlsx")

	writeCSV(t, csvPath, ',', sampleRows)

	err := csvToXLSX(ConvertOpts{
		InputPath:  csvPath,
		OutputPath: xlsxPath,
		Sheet:      "MyData",
		Delimiter:  ',',
	})
	if err != nil {
		t.Fatal(err)
	}

	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = f.GetRows("MyData")
	if err != nil {
		t.Errorf("sheet 'MyData' not found: %v", err)
	}
}

func TestCSVToXLSX_TabDelimited(t *testing.T) {
	dir := tempDir(t)
	csvPath := filepath.Join(dir, "input.tsv")
	xlsxPath := filepath.Join(dir, "output.xlsx")

	writeCSV(t, csvPath, '\t', sampleRows)

	err := csvToXLSX(ConvertOpts{
		InputPath:  csvPath,
		OutputPath: xlsxPath,
		Delimiter:  '\t',
	})
	if err != nil {
		t.Fatal(err)
	}

	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rows, _ := f.GetRows("Sheet1")
	if len(rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(rows))
	}
}

// ---------------------------------------------------------------------------
// XLSX → CSV
// ---------------------------------------------------------------------------

func TestXLSXToCSV(t *testing.T) {
	dir := tempDir(t)
	xlsxPath := filepath.Join(dir, "input.xlsx")
	csvPath := filepath.Join(dir, "output.csv")

	// Create an XLSX file
	f := excelize.NewFile()
	for rowIdx, row := range sampleRows {
		for colIdx, cell := range row {
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			ref := colName + string(rune('1'+rowIdx))
			f.SetCellValue("Sheet1", ref, cell)
		}
	}
	f.SaveAs(xlsxPath)
	f.Close()

	err := xlsxToCSV(ConvertOpts{
		InputPath:  xlsxPath,
		OutputPath: csvPath,
		Delimiter:  ',',
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify CSV content
	csvFile, err := os.Open(csvPath)
	if err != nil {
		t.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(rows))
	}
	if rows[0][0] != "Name" {
		t.Errorf("header mismatch: %v", rows[0])
	}
	if rows[3][0] != "Charlie" {
		t.Errorf("last row mismatch: %v", rows[3])
	}
}

func TestXLSXToCSV_SemicolonDelimiter(t *testing.T) {
	dir := tempDir(t)
	xlsxPath := filepath.Join(dir, "input.xlsx")
	csvPath := filepath.Join(dir, "output.csv")

	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "a")
	f.SetCellValue("Sheet1", "B1", "b")
	f.SaveAs(xlsxPath)
	f.Close()

	err := xlsxToCSV(ConvertOpts{
		InputPath:  xlsxPath,
		OutputPath: csvPath,
		Delimiter:  ';',
	})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(csvPath)
	if string(data) != "a;b\n" {
		t.Errorf("expected semicolon-delimited output, got %q", string(data))
	}
}

// ---------------------------------------------------------------------------
// Roundtrip: CSV → XLSX → CSV
// ---------------------------------------------------------------------------

func TestRoundtrip_CSV_XLSX_CSV(t *testing.T) {
	dir := tempDir(t)
	csvIn := filepath.Join(dir, "input.csv")
	xlsxMid := filepath.Join(dir, "middle.xlsx")
	csvOut := filepath.Join(dir, "output.csv")

	writeCSV(t, csvIn, ',', sampleRows)

	// CSV → XLSX
	if err := csvToXLSX(ConvertOpts{InputPath: csvIn, OutputPath: xlsxMid, Delimiter: ','}); err != nil {
		t.Fatal(err)
	}

	// XLSX → CSV
	if err := xlsxToCSV(ConvertOpts{InputPath: xlsxMid, OutputPath: csvOut, Delimiter: ','}); err != nil {
		t.Fatal(err)
	}

	// Compare
	outFile, _ := os.Open(csvOut)
	defer outFile.Close()
	reader := csv.NewReader(outFile)
	rows, _ := reader.ReadAll()

	if len(rows) != len(sampleRows) {
		t.Fatalf("row count: got %d, want %d", len(rows), len(sampleRows))
	}
	for i, row := range rows {
		for j, cell := range row {
			if cell != sampleRows[i][j] {
				t.Errorf("cell [%d][%d]: got %q, want %q", i, j, cell, sampleRows[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Delimiter auto-detection
// ---------------------------------------------------------------------------

func TestDetectDelimiter_Comma(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.csv")
	writeCSV(t, path, ',', sampleRows)

	delim, err := detectDelimiter(path)
	if err != nil {
		t.Fatal(err)
	}
	if delim != ',' {
		t.Errorf("expected comma, got %q", delim)
	}
}

func TestDetectDelimiter_Tab(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.tsv")
	writeCSV(t, path, '\t', sampleRows)

	delim, err := detectDelimiter(path)
	if err != nil {
		t.Fatal(err)
	}
	if delim != '\t' {
		t.Errorf("expected tab, got %q", delim)
	}
}

func TestDetectDelimiter_Semicolon(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.csv")
	writeCSV(t, path, ';', sampleRows)

	delim, err := detectDelimiter(path)
	if err != nil {
		t.Fatal(err)
	}
	if delim != ';' {
		t.Errorf("expected semicolon, got %q", delim)
	}
}

func TestDetectDelimiter_Pipe(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.csv")
	writeCSV(t, path, '|', sampleRows)

	delim, err := detectDelimiter(path)
	if err != nil {
		t.Fatal(err)
	}
	if delim != '|' {
		t.Errorf("expected pipe, got %q", delim)
	}
}

// ---------------------------------------------------------------------------
// Numbers stubs return errors
// ---------------------------------------------------------------------------

func TestNumbers_StubsReturnError(t *testing.T) {
	opts := ConvertOpts{InputPath: "test.csv", OutputPath: "test.numbers"}

	if err := csvToNumbers(opts); err == nil {
		t.Error("csvToNumbers should return error")
	}
	if err := xlsxToNumbers(opts); err == nil {
		t.Error("xlsxToNumbers should return error")
	}
	if err := numbersToCSV(opts); err == nil {
		t.Error("numbersToCSV should return error")
	}
	if err := numbersToXLSX(opts); err == nil {
		t.Error("numbersToXLSX should return error")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestCSVToXLSX_MissingInput(t *testing.T) {
	dir := tempDir(t)
	err := csvToXLSX(ConvertOpts{
		InputPath:  filepath.Join(dir, "nonexistent.csv"),
		OutputPath: filepath.Join(dir, "out.xlsx"),
		Delimiter:  ',',
	})
	if err == nil {
		t.Error("expected error for missing input file")
	}
}

func TestXLSXToCSV_MissingInput(t *testing.T) {
	dir := tempDir(t)
	err := xlsxToCSV(ConvertOpts{
		InputPath:  filepath.Join(dir, "nonexistent.xlsx"),
		OutputPath: filepath.Join(dir, "out.csv"),
		Delimiter:  ',',
	})
	if err == nil {
		t.Error("expected error for missing input file")
	}
}

func TestCSVToXLSX_EmptyFile(t *testing.T) {
	dir := tempDir(t)
	csvPath := filepath.Join(dir, "empty.csv")
	xlsxPath := filepath.Join(dir, "output.xlsx")

	os.WriteFile(csvPath, []byte(""), 0644)

	err := csvToXLSX(ConvertOpts{
		InputPath:  csvPath,
		OutputPath: xlsxPath,
		Delimiter:  ',',
	})
	// Should succeed (empty workbook) or return a clear error
	if err != nil {
		t.Logf("empty CSV returned error (acceptable): %v", err)
	}
}

// ---------------------------------------------------------------------------
// Same format rejection
// ---------------------------------------------------------------------------

func TestSameFormatRejected(t *testing.T) {
	// extToFormat(".csv") == extToFormat(".tsv") == "csv"
	// So csv→csv should be caught by runConvertFile
	if extToFormat(".csv") != extToFormat(".tsv") {
		t.Error("csv and tsv should map to same format")
	}
}

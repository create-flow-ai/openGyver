package convertfile

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

// csvToXLSX reads a CSV file and writes it as an XLSX workbook.
func csvToXLSX(opts ConvertOpts) error {
	records, err := readCSV(opts.InputPath, opts.Delimiter)
	if err != nil {
		return fmt.Errorf("reading CSV: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	if opts.Sheet != "" {
		sheetName = opts.Sheet
		idx, _ := f.NewSheet(sheetName)
		f.SetActiveSheet(idx)
	}

	for rowIdx, row := range records {
		for colIdx, cell := range row {
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			cellRef := fmt.Sprintf("%s%d", colName, rowIdx+1)
			f.SetCellValue(sheetName, cellRef, cell)
		}
	}

	if err := f.SaveAs(opts.OutputPath); err != nil {
		return fmt.Errorf("writing XLSX: %w", err)
	}
	return nil
}

// xlsxToCSV reads an XLSX workbook and writes the specified sheet as CSV.
func xlsxToCSV(opts ConvertOpts) error {
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

	outFile, err := os.Create(opts.OutputPath)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}
	defer outFile.Close()

	w := csv.NewWriter(outFile)
	w.Comma = opts.Delimiter

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return fmt.Errorf("writing CSV row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}

// readCSV reads a CSV file with optional delimiter auto-detection.
func readCSV(path string, delim rune) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Auto-detect delimiter if default comma
	if delim == ',' {
		delim, err = detectDelimiter(path)
		if err != nil {
			delim = ','
		}
	}

	reader := csv.NewReader(file)
	reader.Comma = delim
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // allow ragged rows

	return reader.ReadAll()
}

// detectDelimiter samples the first few lines to guess the delimiter.
func detectDelimiter(path string) (rune, error) {
	file, err := os.Open(path)
	if err != nil {
		return ',', err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sample string
	for i := 0; i < 5 && scanner.Scan(); i++ {
		sample += scanner.Text() + "\n"
	}

	candidates := []rune{',', '\t', ';', '|'}
	best := ','
	bestCount := 0
	for _, c := range candidates {
		count := strings.Count(sample, string(c))
		if count > bestCount {
			bestCount = count
			best = c
		}
	}

	// Sanity: if no delimiter found at all, default to comma
	if bestCount == 0 {
		return ',', nil
	}
	return best, nil
}

// Ensure utf8 is available (used indirectly by csv package).
var _ = utf8.UTFMax

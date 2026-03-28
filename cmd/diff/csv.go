package diff

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	csvFile1 string
	csvFile2 string
)

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Diff between two CSV files",
	Long: `Compare two CSV files and show added, removed, and changed rows.

Rows are compared by their full content. The output shows which rows
were added, removed, or changed between the two files.

Examples:
  openGyver diff csv --file1 old.csv --file2 new.csv
  openGyver diff csv --file1 export-jan.csv --file2 export-feb.csv --json`,
	RunE: func(c *cobra.Command, args []string) error {
		if csvFile1 == "" || csvFile2 == "" {
			return fmt.Errorf("both --file1 and --file2 are required")
		}

		rows1, err := readCSV(csvFile1)
		if err != nil {
			return fmt.Errorf("reading file1: %w", err)
		}
		rows2, err := readCSV(csvFile2)
		if err != nil {
			return fmt.Errorf("reading file2: %w", err)
		}

		diffs := compareCSV(rows1, rows2)

		if jsonOut {
			var entries []map[string]interface{}
			for _, d := range diffs {
				entry := map[string]interface{}{
					"type": d.diffType,
					"row":  d.lineNum,
				}
				if d.diffType == "removed" || d.diffType == "changed" {
					entry["old"] = d.oldRow
				}
				if d.diffType == "added" || d.diffType == "changed" {
					entry["new"] = d.newRow
				}
				entries = append(entries, entry)
			}
			return cmd.PrintJSON(map[string]interface{}{
				"file1": csvFile1,
				"file2": csvFile2,
				"diffs": entries,
			})
		}

		if len(diffs) == 0 {
			fmt.Println("No differences found.")
			return nil
		}

		fmt.Printf("--- %s\n", csvFile1)
		fmt.Printf("+++ %s\n", csvFile2)
		fmt.Println()
		for _, d := range diffs {
			switch d.diffType {
			case "added":
				fmt.Printf("+ row %d: %s\n", d.lineNum, formatRow(d.newRow))
			case "removed":
				fmt.Printf("- row %d: %s\n", d.lineNum, formatRow(d.oldRow))
			case "changed":
				fmt.Printf("~ row %d:\n", d.lineNum)
				fmt.Printf("  - %s\n", formatRow(d.oldRow))
				fmt.Printf("  + %s\n", formatRow(d.newRow))
			}
		}
		return nil
	},
}

// csvDiff represents a single row difference.
type csvDiff struct {
	diffType string // "added", "removed", "changed"
	lineNum  int
	oldRow   []string
	newRow   []string
}

// readCSV reads and parses a CSV file.
func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // Allow variable field counts.
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV in %s: %w", path, err)
	}
	return records, nil
}

// compareCSV compares two sets of CSV rows line-by-line.
func compareCSV(a, b [][]string) []csvDiff {
	var diffs []csvDiff

	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	for i := 0; i < maxLen; i++ {
		lineNum := i + 1
		if i >= len(a) {
			diffs = append(diffs, csvDiff{diffType: "added", lineNum: lineNum, newRow: b[i]})
		} else if i >= len(b) {
			diffs = append(diffs, csvDiff{diffType: "removed", lineNum: lineNum, oldRow: a[i]})
		} else if !rowsEqual(a[i], b[i]) {
			diffs = append(diffs, csvDiff{diffType: "changed", lineNum: lineNum, oldRow: a[i], newRow: b[i]})
		}
	}

	return diffs
}

// rowsEqual checks if two CSV rows are identical.
func rowsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// formatRow formats a CSV row for human-readable output.
func formatRow(row []string) string {
	return strings.Join(row, ", ")
}

func init() {
	csvCmd.Flags().StringVar(&csvFile1, "file1", "", "first CSV file to compare (required)")
	csvCmd.Flags().StringVar(&csvFile2, "file2", "", "second CSV file to compare (required)")
	_ = csvCmd.MarkFlagRequired("file1")
	_ = csvCmd.MarkFlagRequired("file2")
	register(csvCmd)
}

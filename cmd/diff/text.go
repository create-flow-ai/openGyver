package diff

import (
	"fmt"
	"os"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	textFile1 string
	textFile2 string
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Unified text diff between two files",
	Long: `Compare two text files line-by-line and display a unified diff.

Lines only in file1 are prefixed with "-", lines only in file2 with "+",
and common lines with " ".

Examples:
  openGyver diff text --file1 a.txt --file2 b.txt
  openGyver diff text --file1 original.go --file2 modified.go --json`,
	RunE: func(c *cobra.Command, args []string) error {
		if textFile1 == "" || textFile2 == "" {
			return fmt.Errorf("both --file1 and --file2 are required")
		}

		lines1, err := readLines(textFile1)
		if err != nil {
			return fmt.Errorf("reading file1: %w", err)
		}
		lines2, err := readLines(textFile2)
		if err != nil {
			return fmt.Errorf("reading file2: %w", err)
		}

		edits := computeLCS(lines1, lines2)

		if jsonOut {
			var diffs []map[string]interface{}
			for _, e := range edits {
				diffs = append(diffs, map[string]interface{}{
					"op":   e.op,
					"line": e.line,
				})
			}
			return cmd.PrintJSON(map[string]interface{}{
				"file1": textFile1,
				"file2": textFile2,
				"diffs": diffs,
			})
		}

		fmt.Printf("--- %s\n", textFile1)
		fmt.Printf("+++ %s\n", textFile2)
		for _, e := range edits {
			switch e.op {
			case "remove":
				fmt.Printf("-%s\n", e.line)
			case "add":
				fmt.Printf("+%s\n", e.line)
			default:
				fmt.Printf(" %s\n", e.line)
			}
		}
		return nil
	},
}

// edit represents a single line in the diff output.
type edit struct {
	op   string // "keep", "add", "remove"
	line string
}

// computeLCS computes a unified diff using the Longest Common Subsequence.
func computeLCS(a, b []string) []edit {
	m, n := len(a), len(b)

	// Build LCS table.
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	// Backtrack to produce edits.
	var result []edit
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			result = append(result, edit{op: "keep", line: a[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			result = append(result, edit{op: "add", line: b[j-1]})
			j--
		} else {
			result = append(result, edit{op: "remove", line: a[i-1]})
			i--
		}
	}

	// Reverse (backtrack produces them in reverse order).
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

// readLines reads a file and splits it into lines.
func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := strings.TrimSuffix(string(data), "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}

func init() {
	textCmd.Flags().StringVar(&textFile1, "file1", "", "first file to compare (required)")
	textCmd.Flags().StringVar(&textFile2, "file2", "", "second file to compare (required)")
	_ = textCmd.MarkFlagRequired("file1")
	_ = textCmd.MarkFlagRequired("file2")
	register(textCmd)
}

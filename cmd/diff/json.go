package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	jsonFile1 string
	jsonFile2 string
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Structural diff between two JSON files",
	Long: `Compare two JSON files structurally and show added, removed, and changed keys.

Recursively walks both JSON structures and reports differences at each path.

Examples:
  openGyver diff json --file1 old.json --file2 new.json
  openGyver diff json --file1 config-a.json --file2 config-b.json --json`,
	RunE: func(c *cobra.Command, args []string) error {
		if jsonFile1 == "" || jsonFile2 == "" {
			return fmt.Errorf("both --file1 and --file2 are required")
		}

		v1, err := readJSON(jsonFile1)
		if err != nil {
			return fmt.Errorf("reading file1: %w", err)
		}
		v2, err := readJSON(jsonFile2)
		if err != nil {
			return fmt.Errorf("reading file2: %w", err)
		}

		diffs := compareJSON("", v1, v2)

		if jsonOut {
			var entries []map[string]interface{}
			for _, d := range diffs {
				entry := map[string]interface{}{
					"path": d.path,
					"type": d.diffType,
				}
				if d.diffType == "removed" || d.diffType == "changed" {
					entry["old"] = d.oldVal
				}
				if d.diffType == "added" || d.diffType == "changed" {
					entry["new"] = d.newVal
				}
				entries = append(entries, entry)
			}
			return cmd.PrintJSON(map[string]interface{}{
				"file1": jsonFile1,
				"file2": jsonFile2,
				"diffs": entries,
			})
		}

		if len(diffs) == 0 {
			fmt.Println("No differences found.")
			return nil
		}

		fmt.Printf("--- %s\n", jsonFile1)
		fmt.Printf("+++ %s\n", jsonFile2)
		fmt.Println()
		for _, d := range diffs {
			switch d.diffType {
			case "added":
				fmt.Printf("+ %s: %s\n", d.path, formatVal(d.newVal))
			case "removed":
				fmt.Printf("- %s: %s\n", d.path, formatVal(d.oldVal))
			case "changed":
				fmt.Printf("~ %s: %s -> %s\n", d.path, formatVal(d.oldVal), formatVal(d.newVal))
			}
		}
		return nil
	},
}

// jsonDiff represents a single difference in the JSON structure.
type jsonDiff struct {
	path     string
	diffType string // "added", "removed", "changed"
	oldVal   interface{}
	newVal   interface{}
}

// readJSON reads and unmarshals a JSON file.
func readJSON(path string) (interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return v, nil
}

// compareJSON recursively compares two JSON values at the given path.
func compareJSON(path string, a, b interface{}) []jsonDiff {
	var diffs []jsonDiff

	// Both nil — equal.
	if a == nil && b == nil {
		return nil
	}

	// One nil, one not.
	if a == nil {
		return []jsonDiff{{path: pathOrRoot(path), diffType: "added", newVal: b}}
	}
	if b == nil {
		return []jsonDiff{{path: pathOrRoot(path), diffType: "removed", oldVal: a}}
	}

	aMap, aIsMap := a.(map[string]interface{})
	bMap, bIsMap := b.(map[string]interface{})

	aArr, aIsArr := a.([]interface{})
	bArr, bIsArr := b.([]interface{})

	switch {
	case aIsMap && bIsMap:
		// Collect all keys.
		keys := make(map[string]bool)
		for k := range aMap {
			keys[k] = true
		}
		for k := range bMap {
			keys[k] = true
		}

		// Sort for stable output.
		sorted := make([]string, 0, len(keys))
		for k := range keys {
			sorted = append(sorted, k)
		}
		sort.Strings(sorted)

		for _, k := range sorted {
			childPath := joinPath(path, k)
			aVal, aHas := aMap[k]
			bVal, bHas := bMap[k]

			if aHas && !bHas {
				diffs = append(diffs, jsonDiff{path: childPath, diffType: "removed", oldVal: aVal})
			} else if !aHas && bHas {
				diffs = append(diffs, jsonDiff{path: childPath, diffType: "added", newVal: bVal})
			} else {
				diffs = append(diffs, compareJSON(childPath, aVal, bVal)...)
			}
		}

	case aIsArr && bIsArr:
		maxLen := len(aArr)
		if len(bArr) > maxLen {
			maxLen = len(bArr)
		}
		for i := 0; i < maxLen; i++ {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			if i >= len(aArr) {
				diffs = append(diffs, jsonDiff{path: childPath, diffType: "added", newVal: bArr[i]})
			} else if i >= len(bArr) {
				diffs = append(diffs, jsonDiff{path: childPath, diffType: "removed", oldVal: aArr[i]})
			} else {
				diffs = append(diffs, compareJSON(childPath, aArr[i], bArr[i])...)
			}
		}

	default:
		if !reflect.DeepEqual(a, b) {
			diffs = append(diffs, jsonDiff{path: pathOrRoot(path), diffType: "changed", oldVal: a, newVal: b})
		}
	}

	return diffs
}

// joinPath creates a dotted path for nested keys.
func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

// pathOrRoot returns "." for the root path.
func pathOrRoot(path string) string {
	if path == "" {
		return "."
	}
	return path
}

// formatVal produces a compact string representation of a JSON value.
func formatVal(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case float64:
		// Print integers without decimal point.
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case map[string]interface{}, []interface{}:
		b, _ := json.Marshal(val)
		s := string(b)
		if len(s) > 60 {
			return s[:57] + "..."
		}
		return s
	default:
		return fmt.Sprintf("%v", val)
	}
}

func init() {
	jsonCmd.Flags().StringVar(&jsonFile1, "file1", "", "first JSON file to compare (required)")
	jsonCmd.Flags().StringVar(&jsonFile2, "file2", "", "second JSON file to compare (required)")
	_ = jsonCmd.MarkFlagRequired("file1")
	_ = jsonCmd.MarkFlagRequired("file2")
	register(jsonCmd)
}

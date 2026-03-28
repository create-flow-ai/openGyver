package regex

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── persistent flags ───────────────────────────────────────────────────────

var jsonOut bool

// ── parent command ─────────────────────────────────────────────────────────

var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "Regular expression tools",
	Long: `Regular expression tools — test, replace, and extract with Go regexp syntax.

SUBCOMMANDS:

  test      Test a regex against input (shows match result, matches, groups)
  replace   Regex find-and-replace
  extract   Extract all matches from input

All subcommands support --json/-j for machine-readable output.

EXAMPLES:

  openGyver regex test "\d+" "order 42 has 3 items"
  openGyver regex test --global "\d+" "order 42 has 3 items"
  openGyver regex replace "\d+" "X" "order 42 has 3 items"
  openGyver regex extract "\w+@\w+\.\w+" "Contact alice@ex.com or bob@ex.com"
  openGyver regex extract "\d+" --file data.txt`,
}

// ── helpers ────────────────────────────────────────────────────────────────

// readInput returns text from positional args, --file, or stdin.
func readInput(args []string, startIdx int, filePath string) (string, error) {
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}
	if len(args) > startIdx {
		return strings.Join(args[startIdx:], " "), nil
	}
	// Try stdin if it's piped.
	info, _ := os.Stdin.Stat()
	if info.Mode()&os.ModeCharDevice == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		var sb strings.Builder
		for scanner.Scan() {
			sb.WriteString(scanner.Text())
			sb.WriteByte('\n')
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return sb.String(), nil
	}
	return "", fmt.Errorf("no input provided (pass text as an argument, use --file, or pipe via stdin)")
}

// ── test ───────────────────────────────────────────────────────────────────

var testGlobal bool

var testCmd = &cobra.Command{
	Use:   "test <pattern> <input>",
	Short: "Test a regex against input",
	Long: `Test whether a regular expression matches the given input.

Shows whether the pattern matches, all matched substrings, and any
captured groups. Use --global/-g to find all matches in the input.

Examples:
  openGyver regex test "\d+" "order 42 has 3 items"
  openGyver regex test --global "(\w+)@(\w+)" "alice@example bob@test"
  openGyver regex test "^hello" "hello world"`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		pattern := args[0]
		input, err := readInput(args, 1, "")
		if err != nil {
			return err
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		matches := re.MatchString(input)

		if testGlobal {
			allMatches := re.FindAllStringSubmatch(input, -1)

			if jsonOut {
				result := map[string]interface{}{
					"pattern": pattern,
					"input":   input,
					"matches": matches,
				}
				var matchList []map[string]interface{}
				for _, m := range allMatches {
					entry := map[string]interface{}{
						"match": m[0],
					}
					if len(m) > 1 {
						entry["groups"] = m[1:]
					}
					matchList = append(matchList, entry)
				}
				result["all_matches"] = matchList
				return cmd.PrintJSON(result)
			}

			fmt.Printf("Pattern:  %s\n", pattern)
			fmt.Printf("Input:    %s\n", input)
			fmt.Printf("Matches:  %t\n", matches)
			fmt.Println()
			if len(allMatches) == 0 {
				fmt.Println("No matches found.")
			} else {
				for i, m := range allMatches {
					fmt.Printf("Match %d:  %s\n", i+1, m[0])
					if len(m) > 1 {
						for j, g := range m[1:] {
							fmt.Printf("  Group %d: %s\n", j+1, g)
						}
					}
				}
			}
			return nil
		}

		// Single match mode
		loc := re.FindStringSubmatchIndex(input)
		var match string
		var groups []string
		if loc != nil {
			match = input[loc[0]:loc[1]]
			sub := re.FindStringSubmatch(input)
			if len(sub) > 1 {
				groups = sub[1:]
			}
		}

		if jsonOut {
			result := map[string]interface{}{
				"pattern": pattern,
				"input":   input,
				"matches": matches,
				"match":   match,
				"groups":  groups,
			}
			return cmd.PrintJSON(result)
		}

		fmt.Printf("Pattern:  %s\n", pattern)
		fmt.Printf("Input:    %s\n", input)
		fmt.Printf("Matches:  %t\n", matches)
		if match != "" {
			fmt.Printf("Match:    %s\n", match)
			if len(groups) > 0 {
				fmt.Println()
				for i, g := range groups {
					fmt.Printf("  Group %d: %s\n", i+1, g)
				}
			}
		}
		return nil
	},
}

// ── replace ────────────────────────────────────────────────────────────────

var replaceCmd = &cobra.Command{
	Use:   "replace <pattern> <replacement> <input>",
	Short: "Regex find-and-replace",
	Long: `Replace all occurrences of a regex pattern in the input string.

The replacement string may include $1, $2, etc. for captured groups,
and $0 for the entire match.

Examples:
  openGyver regex replace "\d+" "X" "order 42 has 3 items"
  openGyver regex replace "(\w+)@(\w+)" "$2/$1" "alice@example"
  openGyver regex replace "\s+" "-" "hello   world"`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(c *cobra.Command, args []string) error {
		pattern := args[0]
		replacement := args[1]
		input, err := readInput(args, 2, "")
		if err != nil {
			return err
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		result := re.ReplaceAllString(input, replacement)

		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"pattern":     pattern,
				"replacement": replacement,
				"input":       input,
				"result":      result,
			})
		}

		fmt.Println(result)
		return nil
	},
}

// ── extract ────────────────────────────────────────────────────────────────

var extractFile string

var extractCmd = &cobra.Command{
	Use:   "extract <pattern> [input]",
	Short: "Extract all matches from input",
	Long: `Extract all occurrences of a regex pattern from the input.

Outputs one match per line. Use --file/-f to read input from a file.

Examples:
  openGyver regex extract "\w+@\w+\.\w+" "Contact alice@ex.com or bob@ex.com"
  openGyver regex extract "\d+\.\d+" --file prices.txt
  echo "hello 123 world 456" | openGyver regex extract "\d+"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		pattern := args[0]
		input, err := readInput(args, 1, extractFile)
		if err != nil {
			return err
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		allMatches := re.FindAllString(input, -1)

		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"pattern": pattern,
				"count":   len(allMatches),
				"matches": allMatches,
			})
		}

		for _, m := range allMatches {
			fmt.Println(m)
		}
		return nil
	},
}

// ── init ───────────────────────────────────────────────────────────────────

func init() {
	regexCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	testCmd.Flags().BoolVarP(&testGlobal, "global", "g", false, "find all matches (not just the first)")
	regexCmd.AddCommand(testCmd)

	regexCmd.AddCommand(replaceCmd)

	extractCmd.Flags().StringVarP(&extractFile, "file", "f", "", "read input from file")
	regexCmd.AddCommand(extractCmd)

	cmd.Register(regexCmd)
}

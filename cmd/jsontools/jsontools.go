package jsontools

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	jsonOut    bool
	filePath   string
	outputPath string
	indent     int
)

// ---------------------------------------------------------------------------
// Parent command
// ---------------------------------------------------------------------------

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON tools — format, minify, validate, path query, escape/unescape",
	Long: `A collection of JSON utilities for everyday use.

SUBCOMMANDS:

  format     Format (beautify) JSON with configurable indentation
  minify     Minify JSON by removing all unnecessary whitespace
  validate   Check whether a string is valid JSON
  path       Evaluate a dot-notation path against a JSON document
  escape     Escape a string for safe embedding inside JSON
  unescape   Unescape a JSON-encoded string

INPUT:

  Most subcommands accept JSON as a positional argument or via --file/-f.
  When both are provided, --file takes precedence.

Examples:
  openGyver json format '{"a":1}'
  openGyver json minify --file data.json --output data.min.json
  openGyver json validate '{"ok":true}'
  openGyver json path --file config.json 'database.hosts[0].address'
  openGyver json escape 'line1\nline2'
  openGyver json unescape '"hello\tworld"'`,
}

// ---------------------------------------------------------------------------
// format
// ---------------------------------------------------------------------------

var formatCmd = &cobra.Command{
	Use:   "format [json-string]",
	Short: "Format (beautify) JSON",
	Long: `Format a JSON string with indentation for readability.

The input can be provided as a positional argument or read from a file
with --file/-f. The formatted output is printed to stdout unless --output/-o
is specified.

Flags:
  --indent    Number of spaces for each indentation level (default 2)
  --file/-f   Read JSON from a file instead of an argument
  --output/-o Write result to a file instead of stdout
  --json/-j   Wrap the output in a JSON envelope

Examples:
  openGyver json format '{"name":"Alice","age":30}'
  openGyver json format --indent 4 '{"a":1,"b":2}'
  openGyver json format --file input.json --output pretty.json
  openGyver json format --json '{"x":1}'`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFormat,
}

func runFormat(c *cobra.Command, args []string) error {
	raw, err := readInput(args)
	if err != nil {
		return err
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	indentStr := strings.Repeat(" ", indent)
	out, err := json.MarshalIndent(parsed, "", indentStr)
	if err != nil {
		return fmt.Errorf("formatting error: %w", err)
	}

	formatted := string(out) + "\n"

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"formatted": string(out),
			"indent":    indent,
		})
	}

	return writeOutput(formatted)
}

// ---------------------------------------------------------------------------
// minify
// ---------------------------------------------------------------------------

var minifyCmd = &cobra.Command{
	Use:   "minify [json-string]",
	Short: "Minify JSON (remove whitespace)",
	Long: `Compact a JSON string by removing all unnecessary whitespace.

The input can be provided as a positional argument or read from a file
with --file/-f. The minified output is printed to stdout unless --output/-o
is specified.

Flags:
  --file/-f   Read JSON from a file instead of an argument
  --output/-o Write result to a file instead of stdout
  --json/-j   Wrap the output in a JSON envelope

Examples:
  openGyver json minify '{  "name": "Alice",  "age": 30  }'
  openGyver json minify --file pretty.json
  openGyver json minify --file pretty.json --output compact.json
  openGyver json minify --json '{ "x": 1 }'`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMinify,
}

func runMinify(c *cobra.Command, args []string) error {
	raw, err := readInput(args)
	if err != nil {
		return err
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	out, err := json.Marshal(parsed)
	if err != nil {
		return fmt.Errorf("minify error: %w", err)
	}

	minified := string(out) + "\n"

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"minified": string(out),
		})
	}

	return writeOutput(minified)
}

// ---------------------------------------------------------------------------
// validate
// ---------------------------------------------------------------------------

var validateCmd = &cobra.Command{
	Use:   "validate [json-string]",
	Short: "Validate JSON",
	Long: `Check whether a string is valid JSON.

Prints "valid" if the input is well-formed JSON, or an error message
describing why it is not.

Flags:
  --file/-f  Read JSON from a file instead of an argument
  --json/-j  Wrap the output in a JSON envelope

Examples:
  openGyver json validate '{"ok": true}'
  openGyver json validate '{"missing": }'
  openGyver json validate --file data.json
  openGyver json validate --json '{"a":1}'`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

func runValidate(c *cobra.Command, args []string) error {
	raw, err := readInput(args)
	if err != nil {
		return err
	}

	valid := json.Valid([]byte(raw))
	var errMsg string
	if !valid {
		var js json.RawMessage
		if uerr := json.Unmarshal([]byte(raw), &js); uerr != nil {
			errMsg = uerr.Error()
		} else {
			errMsg = "unknown error"
		}
	}

	if jsonOut {
		result := map[string]interface{}{
			"valid": valid,
		}
		if errMsg != "" {
			result["error"] = errMsg
		}
		return cmd.PrintJSON(result)
	}

	if valid {
		fmt.Println("valid")
	} else {
		fmt.Printf("invalid: %s\n", errMsg)
	}
	return nil
}

// ---------------------------------------------------------------------------
// path
// ---------------------------------------------------------------------------

var pathCmd = &cobra.Command{
	Use:   "path <expression>",
	Short: "Evaluate a JSON path expression (dot notation)",
	Long: `Extract a value from a JSON document using dot-notation path expressions.

The JSON input must be provided via --file/-f. The path expression is the
positional argument.

PATH SYNTAX:

  Dot notation with optional array indexing:

    key               Access a top-level key
    key.nested        Nested object access
    key[0]            Array element by index
    data.users[2].name  Mixed object/array traversal

Flags:
  --file/-f  Read JSON from a file (required)
  --json/-j  Wrap the output in a JSON envelope

Examples:
  openGyver json path --file config.json 'database.host'
  openGyver json path --file users.json 'users[0].name'
  openGyver json path --file data.json 'results[3].tags[0]'
  openGyver json path --json --file config.json 'server.port'`,
	Args: cobra.ExactArgs(1),
	RunE: runPath,
}

func runPath(c *cobra.Command, args []string) error {
	if filePath == "" {
		return fmt.Errorf("--file/-f is required for the path command")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	expression := args[0]
	result, err := evaluatePath(root, expression)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"path":  expression,
			"value": result,
		})
	}

	// Pretty-print the result: if it's a complex type, marshal it; otherwise
	// print the scalar directly.
	switch v := result.(type) {
	case map[string]interface{}, []interface{}:
		out, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	case string:
		fmt.Println(v)
	case nil:
		fmt.Println("null")
	default:
		fmt.Println(v)
	}

	return nil
}

// evaluatePath walks a parsed JSON value using a dot-notation path with
// optional array indexing (e.g. "data.users[0].name").
func evaluatePath(root interface{}, path string) (interface{}, error) {
	segments := parsePath(path)
	current := root

	for _, seg := range segments {
		switch {
		case seg.index >= 0:
			// Array access: first resolve the key (if any), then index.
			if seg.key != "" {
				obj, ok := current.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("expected object at %q, got %T", seg.key, current)
				}
				val, exists := obj[seg.key]
				if !exists {
					return nil, fmt.Errorf("key %q not found", seg.key)
				}
				current = val
			}
			arr, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected array for index [%d], got %T", seg.index, current)
			}
			if seg.index >= len(arr) {
				return nil, fmt.Errorf("index %d out of range (length %d)", seg.index, len(arr))
			}
			current = arr[seg.index]
		default:
			// Plain key access.
			obj, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("expected object at %q, got %T", seg.key, current)
			}
			val, exists := obj[seg.key]
			if !exists {
				return nil, fmt.Errorf("key %q not found", seg.key)
			}
			current = val
		}
	}

	return current, nil
}

// pathSegment represents one step in a path expression.
// If index >= 0 an array lookup is performed (after an optional key lookup).
type pathSegment struct {
	key   string
	index int // -1 means "no index"
}

// parsePath splits "data.users[0].name" into segments:
//
//	{key:"data", index:-1}, {key:"users", index:0}, {key:"name", index:-1}
func parsePath(path string) []pathSegment {
	parts := strings.Split(path, ".")
	segments := make([]pathSegment, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}
		// Check for bracket notation: key[0], key[12], or just [0]
		if idx := strings.Index(part, "["); idx != -1 {
			key := part[:idx]
			rest := part[idx:]

			// There may be chained indices: key[0][1] — handle all of them.
			first := true
			for rest != "" {
				open := strings.Index(rest, "[")
				close := strings.Index(rest, "]")
				if open == -1 || close == -1 || close <= open+1 {
					break
				}
				numStr := rest[open+1 : close]
				n, err := strconv.Atoi(numStr)
				if err != nil {
					break
				}
				if first {
					segments = append(segments, pathSegment{key: key, index: n})
					first = false
				} else {
					segments = append(segments, pathSegment{key: "", index: n})
				}
				rest = rest[close+1:]
			}

			if first {
				// No valid index found — treat the whole thing as a key.
				segments = append(segments, pathSegment{key: part, index: -1})
			}
		} else {
			segments = append(segments, pathSegment{key: part, index: -1})
		}
	}

	return segments
}

// ---------------------------------------------------------------------------
// escape
// ---------------------------------------------------------------------------

var escapeCmd = &cobra.Command{
	Use:   "escape <string>",
	Short: "Escape a string for JSON embedding",
	Long: `Escape a raw string so it can be safely embedded inside a JSON document.

The result is a JSON string literal (with surrounding double quotes).
Special characters (\n, \t, \", \\, etc.) are properly escaped.

Flags:
  --json/-j  Wrap the output in a JSON envelope

Examples:
  openGyver json escape 'hello "world"'
  openGyver json escape 'line1\nline2'
  openGyver json escape --json 'tab\there'`,
	Args: cobra.ExactArgs(1),
	RunE: runEscape,
}

func runEscape(c *cobra.Command, args []string) error {
	raw := args[0]
	escaped, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("escape error: %w", err)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"original": raw,
			"escaped":  string(escaped),
		})
	}

	fmt.Println(string(escaped))
	return nil
}

// ---------------------------------------------------------------------------
// unescape
// ---------------------------------------------------------------------------

var unescapeCmd = &cobra.Command{
	Use:   "unescape <json-string>",
	Short: "Unescape a JSON string literal",
	Long: `Unescape a JSON-encoded string, removing surrounding quotes and
converting escape sequences back to their literal characters.

The input must be a valid JSON string (typically surrounded by double
quotes, e.g. "hello\tworld").

Flags:
  --json/-j  Wrap the output in a JSON envelope

Examples:
  openGyver json unescape '"hello\tworld"'
  openGyver json unescape '"line1\nline2"'
  openGyver json unescape --json '"escaped \"quotes\""'`,
	Args: cobra.ExactArgs(1),
	RunE: runUnescape,
}

func runUnescape(c *cobra.Command, args []string) error {
	raw := args[0]

	var unescaped string
	if err := json.Unmarshal([]byte(raw), &unescaped); err != nil {
		return fmt.Errorf("unescape error (input must be a JSON string literal): %w", err)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"original":  raw,
			"unescaped": unescaped,
		})
	}

	fmt.Println(unescaped)
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// readInput returns JSON text from --file (if set) or the first positional arg.
func readInput(args []string) (string, error) {
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}
	if len(args) > 0 {
		return args[0], nil
	}
	return "", fmt.Errorf("provide JSON as an argument or use --file/-f")
}

// writeOutput writes text to --output (if set) or prints to stdout.
func writeOutput(text string) error {
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(text), 0o644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", outputPath)
		return nil
	}
	fmt.Print(text)
	return nil
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func init() {
	// Persistent flags available to all subcommands.
	jsonCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	jsonCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "read input from a file")

	// format-specific flags.
	formatCmd.Flags().IntVar(&indent, "indent", 2, "number of spaces per indentation level")
	formatCmd.Flags().StringVarP(&outputPath, "output", "o", "", "write result to a file")

	// minify-specific flags.
	minifyCmd.Flags().StringVarP(&outputPath, "output", "o", "", "write result to a file")

	// Wire subcommands.
	jsonCmd.AddCommand(formatCmd)
	jsonCmd.AddCommand(minifyCmd)
	jsonCmd.AddCommand(validateCmd)
	jsonCmd.AddCommand(pathCmd)
	jsonCmd.AddCommand(escapeCmd)
	jsonCmd.AddCommand(unescapeCmd)

	cmd.Register(jsonCmd)
}

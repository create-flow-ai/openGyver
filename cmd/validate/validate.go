package validate

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v3"
)

// ── flags ───────────────────────────────────────────────────────────────────

var (
	jsonOut   bool
	inputFile string
)

// ── parent command ──────────────────────────────────────────────────────────

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate data files (HTML, CSV, XML, YAML, TOML)",
	Long: `Validate common data and markup formats.

SUBCOMMANDS:

  html    Validate HTML (unclosed/mismatched tags, missing doctype,
          missing alt on <img>, duplicate IDs)
  csv     Validate CSV (consistent column count, proper quoting, encoding)
  xml     Validate XML (well-formedness)
  yaml    Validate YAML syntax
  toml    Validate TOML syntax

Each subcommand accepts input as the first argument, or via --file/-f.
Output is "valid" on success, or a list of errors.
With --json/-j, output is {"valid":true/false,"errors":[...]}.

Examples:
  openGyver validate html --file index.html
  openGyver validate csv 'name,age\nAlice,30\nBob'
  openGyver validate xml '<root><item/></root>'
  openGyver validate yaml --file config.yaml --json
  openGyver validate toml --file pyproject.toml`,
}

// ── subcommands ─────────────────────────────────────────────────────────────

var htmlCmd = &cobra.Command{
	Use:   "html [input]",
	Short: "Validate HTML markup",
	Long: `Validate HTML for common issues:

  - Missing <!DOCTYPE html> declaration
  - Unclosed tags (e.g. <p> without </p>)
  - Mismatched tags (e.g. <b>...</i>)
  - Missing alt attribute on <img> elements
  - Duplicate id attributes

Uses the golang.org/x/net/html tokenizer for parsing.

Examples:
  openGyver validate html '<html><body><p>Hello</p></body></html>'
  openGyver validate html --file page.html
  openGyver validate html --file page.html --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		errs := validateHTML(input)
		return printResult(errs)
	},
}

var csvCmd = &cobra.Command{
	Use:   "csv [input]",
	Short: "Validate CSV data",
	Long: `Validate CSV for common issues:

  - Inconsistent column count across rows
  - Improper quoting / bare quotes
  - Encoding / parse errors

Uses encoding/csv in strict mode (FieldsPerRecord set from the first row).

Examples:
  openGyver validate csv 'name,age
Alice,30
Bob,25'
  openGyver validate csv --file data.csv
  openGyver validate csv --file data.csv --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		errs := validateCSV(input)
		return printResult(errs)
	},
}

var xmlCmd = &cobra.Command{
	Use:   "xml [input]",
	Short: "Validate XML well-formedness",
	Long: `Validate XML by parsing with encoding/xml.Decoder.

Reports any parse errors with byte offset position.

Examples:
  openGyver validate xml '<root><item/></root>'
  openGyver validate xml --file feed.xml
  openGyver validate xml --file feed.xml --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		errs := validateXML(input)
		return printResult(errs)
	},
}

var yamlCmd = &cobra.Command{
	Use:   "yaml [input]",
	Short: "Validate YAML syntax",
	Long: `Validate YAML by attempting to unmarshal with gopkg.in/yaml.v3.

Reports any parse errors with line/column context.

Examples:
  openGyver validate yaml 'name: hello'
  openGyver validate yaml --file config.yaml
  openGyver validate yaml --file config.yaml --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		errs := validateYAML(input)
		return printResult(errs)
	},
}

var tomlCmd = &cobra.Command{
	Use:   "toml [input]",
	Short: "Validate TOML syntax",
	Long: `Validate TOML by attempting to decode with github.com/BurntSushi/toml.

Reports any parse errors with position context.

Examples:
  openGyver validate toml 'key = "value"'
  openGyver validate toml --file pyproject.toml
  openGyver validate toml --file config.toml --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		errs := validateTOML(input)
		return printResult(errs)
	},
}

// ── input resolution ───────────────────────────────────────────────────────

func resolveInput(args []string) (string, error) {
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}
	if len(args) == 0 {
		return "", fmt.Errorf("provide input as an argument or use --file/-f")
	}
	return args[0], nil
}

// ── output formatting ──────────────────────────────────────────────────────

type jsonResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

func printResult(errs []string) error {
	if jsonOut {
		res := jsonResult{
			Valid:  len(errs) == 0,
			Errors: errs,
		}
		if res.Errors == nil {
			res.Errors = []string{}
		}
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return fmt.Errorf("JSON marshal error: %w", err)
		}
		fmt.Println(string(b))
		return nil
	}

	if len(errs) == 0 {
		fmt.Println("valid")
		return nil
	}

	for _, e := range errs {
		fmt.Println(e)
	}
	return nil
}

// ── HTML validation ────────────────────────────────────────────────────────

// voidElements are HTML elements that must not have a closing tag.
var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true,
	"embed": true, "hr": true, "img": true, "input": true,
	"link": true, "meta": true, "param": true, "source": true,
	"track": true, "wbr": true,
}

func validateHTML(input string) []string {
	var errs []string

	// Track line numbers: build a byte-offset-to-line map.
	lines := strings.Split(input, "\n")
	byteToLine := make([]int, len(input)+1)
	offset := 0
	for lineNum, line := range lines {
		for i := 0; i <= len(line); i++ {
			if offset+i <= len(input) {
				byteToLine[offset+i] = lineNum + 1
			}
		}
		offset += len(line) + 1 // +1 for the newline
	}

	// Check for DOCTYPE.
	hasDoctype := false
	lowered := strings.ToLower(input)
	if strings.Contains(lowered, "<!doctype") {
		hasDoctype = true
	}
	if !hasDoctype {
		errs = append(errs, "missing <!DOCTYPE html> declaration")
	}

	// Tokenize and check structure.
	tokenizer := html.NewTokenizer(strings.NewReader(input))
	var stack []string     // open tag stack
	var stackLines []int   // line of each open tag
	seenIDs := make(map[string]int) // id -> first line seen

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			raw := tokenizer.Raw()
			pos := strings.Index(input, string(raw))
			line := 0
			if pos >= 0 && pos < len(byteToLine) {
				line = byteToLine[pos]
			}
			errs = append(errs, fmt.Sprintf("line %d: parse error: %v", line, err))
			break
		}

		tn, hasAttr := tokenizer.TagName()
		tagName := strings.ToLower(string(tn))

		// Approximate line number from the tokenizer's raw output position.
		raw := tokenizer.Raw()
		rawStr := string(raw)
		lineNum := findLine(input, rawStr, byteToLine)

		switch tt {
		case html.StartTagToken:
			// Check <img> for alt attribute.
			if tagName == "img" {
				foundAlt := false
				if hasAttr {
					for {
						key, _, more := tokenizer.TagAttr()
						if string(key) == "alt" {
							foundAlt = true
						}
						if !more {
							break
						}
					}
				}
				if !foundAlt {
					errs = append(errs, fmt.Sprintf("line %d: <img> missing alt attribute", lineNum))
				}
			} else if hasAttr {
				// Consume attributes for non-img tags to check for duplicate IDs.
				for {
					key, val, more := tokenizer.TagAttr()
					_ = val
					_ = key
					if string(key) == "id" {
						idVal := string(val)
						if idVal != "" {
							if prevLine, exists := seenIDs[idVal]; exists {
								errs = append(errs, fmt.Sprintf(
									"line %d: duplicate id %q (first seen on line %d)",
									lineNum, idVal, prevLine))
							} else {
								seenIDs[idVal] = lineNum
							}
						}
					}
					if !more {
						break
					}
				}
			}

			// Also check <img> for duplicate IDs (already consumed attrs above for img).
			// For non-void elements, push onto stack.
			if !voidElements[tagName] {
				stack = append(stack, tagName)
				stackLines = append(stackLines, lineNum)
			}

			// For img, we already consumed attrs above, also check id there.
			if tagName == "img" && hasAttr {
				// Attrs already consumed in the img block above; re-parse
				// from raw to check for id. We do a simple string scan.
				checkIDInRaw(rawStr, lineNum, seenIDs, &errs)
			}

		case html.EndTagToken:
			if voidElements[tagName] {
				errs = append(errs, fmt.Sprintf(
					"line %d: unexpected closing tag </%s> (void element)",
					lineNum, tagName))
				continue
			}
			if len(stack) == 0 {
				errs = append(errs, fmt.Sprintf(
					"line %d: unexpected closing tag </%s> with no open tags",
					lineNum, tagName))
				continue
			}
			top := stack[len(stack)-1]
			if top != tagName {
				errs = append(errs, fmt.Sprintf(
					"line %d: mismatched tag: expected </%s> (opened on line %d), found </%s>",
					lineNum, top, stackLines[len(stackLines)-1], tagName))
				// Pop anyway to recover.
				stack = stack[:len(stack)-1]
				stackLines = stackLines[:len(stackLines)-1]
			} else {
				stack = stack[:len(stack)-1]
				stackLines = stackLines[:len(stackLines)-1]
			}

		case html.SelfClosingTagToken:
			// Self-closing tags: check for alt on img, check for duplicate IDs.
			if tagName == "img" {
				foundAlt := false
				if hasAttr {
					for {
						key, _, more := tokenizer.TagAttr()
						if string(key) == "alt" {
							foundAlt = true
						}
						if !more {
							break
						}
					}
				}
				if !foundAlt {
					errs = append(errs, fmt.Sprintf("line %d: <img> missing alt attribute", lineNum))
				}
				checkIDInRaw(rawStr, lineNum, seenIDs, &errs)
			} else if hasAttr {
				for {
					key, val, more := tokenizer.TagAttr()
					if string(key) == "id" {
						idVal := string(val)
						if idVal != "" {
							if prevLine, exists := seenIDs[idVal]; exists {
								errs = append(errs, fmt.Sprintf(
									"line %d: duplicate id %q (first seen on line %d)",
									lineNum, idVal, prevLine))
							} else {
								seenIDs[idVal] = lineNum
							}
						}
					}
					if !more {
						break
					}
				}
			}
		}
	}

	// Report unclosed tags remaining on the stack.
	for i := len(stack) - 1; i >= 0; i-- {
		errs = append(errs, fmt.Sprintf(
			"line %d: unclosed tag <%s>",
			stackLines[i], stack[i]))
	}

	return errs
}

// findLine locates raw in input and returns the approximate line number.
func findLine(input, raw string, byteToLine []int) int {
	idx := strings.Index(input, raw)
	if idx >= 0 && idx < len(byteToLine) {
		return byteToLine[idx]
	}
	return 0
}

// checkIDInRaw does a simple scan of a raw tag string for id="..." attributes.
// Used when the tokenizer has already consumed attributes (e.g. for <img>).
func checkIDInRaw(raw string, lineNum int, seenIDs map[string]int, errs *[]string) {
	lower := strings.ToLower(raw)
	idx := 0
	for {
		pos := strings.Index(lower[idx:], "id=")
		if pos < 0 {
			break
		}
		pos += idx
		valStart := pos + 3
		if valStart >= len(raw) {
			break
		}
		quote := raw[valStart]
		if quote == '"' || quote == '\'' {
			valEnd := strings.IndexByte(raw[valStart+1:], quote)
			if valEnd >= 0 {
				idVal := raw[valStart+1 : valStart+1+valEnd]
				if idVal != "" {
					if prevLine, exists := seenIDs[idVal]; exists {
						*errs = append(*errs, fmt.Sprintf(
							"line %d: duplicate id %q (first seen on line %d)",
							lineNum, idVal, prevLine))
					} else {
						seenIDs[idVal] = lineNum
					}
				}
			}
		}
		idx = valStart + 1
	}
}

// ── CSV validation ─────────────────────────────────────────────────────────

func validateCSV(input string) []string {
	var errs []string

	r := csv.NewReader(strings.NewReader(input))
	r.LazyQuotes = false
	r.FieldsPerRecord = 0 // set dynamically from first row

	rowNum := 0
	expectedCols := 0

	for {
		rowNum++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: %v", rowNum, err))
			continue
		}

		if rowNum == 1 {
			expectedCols = len(record)
			r.FieldsPerRecord = expectedCols
			continue
		}

		if len(record) != expectedCols {
			errs = append(errs, fmt.Sprintf(
				"row %d: expected %d columns, got %d",
				rowNum, expectedCols, len(record)))
		}
	}

	if rowNum <= 1 && len(errs) == 0 {
		// Only header or empty — still valid.
	}

	return errs
}

// ── XML validation ─────────────────────────────────────────────────────────

func validateXML(input string) []string {
	var errs []string

	decoder := xml.NewDecoder(strings.NewReader(input))
	decoder.Strict = true

	for {
		_, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs = append(errs, fmt.Sprintf("offset %d: %v",
				decoder.InputOffset(), err))
			break
		}
	}

	return errs
}

// ── YAML validation ────────────────────────────────────────────────────────

func validateYAML(input string) []string {
	var errs []string

	var data interface{}
	if err := yaml.Unmarshal([]byte(input), &data); err != nil {
		errs = append(errs, err.Error())
	}

	return errs
}

// ── TOML validation ────────────────────────────────────────────────────────

func validateTOML(input string) []string {
	var errs []string

	var data interface{}
	if _, err := toml.Decode(input, &data); err != nil {
		errs = append(errs, err.Error())
	}

	return errs
}

// ── init ────────────────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent (inherited by all subcommands).
	validateCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false,
		`output as {"valid":true/false,"errors":[...]}`)
	validateCmd.PersistentFlags().StringVarP(&inputFile, "file", "f", "",
		"read input from a file instead of an argument")

	// Register all subcommands.
	validateCmd.AddCommand(
		htmlCmd,
		csvCmd,
		xmlCmd,
		yamlCmd,
		tomlCmd,
	)

	cmd.Register(validateCmd)
}

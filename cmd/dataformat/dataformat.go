package dataformat

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ── flags ───────────────────────────────────────────────────────────────────

var (
	jsonOut    bool
	inputFile  string
	outputFile string
)

// ── parent command ──────────────────────────────────────────────────────────

var dataformatCmd = &cobra.Command{
	Use:   "dataformat",
	Short: "Convert between data formats (YAML, JSON, TOML, CSV, XML)",
	Long: `Convert between common data serialisation formats.

SUBCOMMANDS:

  yaml2json   Convert YAML to JSON
  json2yaml   Convert JSON to YAML
  toml2json   Convert TOML to JSON
  json2toml   Convert JSON to TOML
  csv2json    Convert CSV to JSON (first row = headers)
  json2csv    Convert JSON array of objects to CSV
  xml2json    Convert XML to JSON (simple element mapping)
  json2xml    Convert JSON to XML

Each subcommand accepts input as the first argument, or via --file/-f.
Output goes to stdout by default, or to a file via --output/-o.

FLAGS:

  --json,   -j   Wrap output in {"input_format","output_format","data"}
  --file,   -f   Read input from a file instead of an argument
  --output, -o   Write output to a file instead of stdout

Examples:
  openGyver dataformat yaml2json '{"name: hello"}'
  openGyver dataformat json2yaml --file config.json
  openGyver dataformat csv2json --file data.csv --json
  openGyver dataformat toml2json --file config.toml -o config.json`,
}

// ── subcommands ─────────────────────────────────────────────────────────────

var yaml2jsonCmd = &cobra.Command{
	Use:   "yaml2json [input]",
	Short: "Convert YAML to JSON",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "yaml", "json", yamlToJSON)
	},
}

var json2yamlCmd = &cobra.Command{
	Use:   "json2yaml [input]",
	Short: "Convert JSON to YAML",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "json", "yaml", jsonToYAML)
	},
}

var toml2jsonCmd = &cobra.Command{
	Use:   "toml2json [input]",
	Short: "Convert TOML to JSON",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "toml", "json", tomlToJSON)
	},
}

var json2tomlCmd = &cobra.Command{
	Use:   "json2toml [input]",
	Short: "Convert JSON to TOML",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "json", "toml", jsonToTOML)
	},
}

var csv2jsonCmd = &cobra.Command{
	Use:   "csv2json [input]",
	Short: "Convert CSV to JSON (first row = headers)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "csv", "json", csvToJSON)
	},
}

var json2csvCmd = &cobra.Command{
	Use:   "json2csv [input]",
	Short: "Convert JSON array of objects to CSV",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "json", "csv", jsonToCSV)
	},
}

var xml2jsonCmd = &cobra.Command{
	Use:   "xml2json [input]",
	Short: "Convert XML to JSON (simple element mapping)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "xml", "json", xmlToJSON)
	},
}

var json2xmlCmd = &cobra.Command{
	Use:   "json2xml [input]",
	Short: "Convert JSON to XML",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return convert(args, "json", "xml", jsonToXML)
	},
}

// ── shared processing ───────────────────────────────────────────────────────

// convertFn takes raw input and returns the converted output string.
type convertFn func(input string) (string, error)

// convert resolves input, calls fn, and handles --json wrapping and --output.
func convert(args []string, inFmt, outFmt string, fn convertFn) error {
	input, err := resolveInput(args)
	if err != nil {
		return err
	}

	output, err := fn(input)
	if err != nil {
		return err
	}

	var result string
	if jsonOut {
		// Parse the converted data so it nests properly in the wrapper.
		var data interface{}
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			// If output isn't valid JSON (e.g. YAML/TOML/CSV/XML text),
			// embed it as a string.
			data = output
		}
		wrapper := map[string]interface{}{
			"input_format":  inFmt,
			"output_format": outFmt,
			"data":          data,
		}
		b, err := json.MarshalIndent(wrapper, "", "  ")
		if err != nil {
			return fmt.Errorf("JSON wrapping error: %w", err)
		}
		result = string(b)
	} else {
		result = output
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(result+"\n"), 0o644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		return nil
	}

	fmt.Println(result)
	return nil
}

// resolveInput returns the input string from either the first argument or --file.
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

// ── YAML <-> JSON ──────────────────────────────────────────────────────────

func yamlToJSON(input string) (string, error) {
	var data interface{}
	if err := yaml.Unmarshal([]byte(input), &data); err != nil {
		return "", fmt.Errorf("YAML parse error: %w", err)
	}
	data = normalizeYAML(data)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %w", err)
	}
	return string(b), nil
}

// normalizeYAML converts map[string]interface{} (yaml.v3 quirk: map keys may
// not always be strings) to ensure JSON-safe output.
func normalizeYAML(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, v2 := range val {
			out[k] = normalizeYAML(v2)
		}
		return out
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, v2 := range val {
			out[fmt.Sprintf("%v", k)] = normalizeYAML(v2)
		}
		return out
	case []interface{}:
		for i, v2 := range val {
			val[i] = normalizeYAML(v2)
		}
		return val
	default:
		return v
	}
}

func jsonToYAML(input string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", fmt.Errorf("JSON parse error: %w", err)
	}
	b, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("YAML marshal error: %w", err)
	}
	return strings.TrimRight(string(b), "\n"), nil
}

// ── TOML <-> JSON ──────────────────────────────────────────────────────────

func tomlToJSON(input string) (string, error) {
	var data interface{}
	if _, err := toml.Decode(input, &data); err != nil {
		return "", fmt.Errorf("TOML parse error: %w", err)
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %w", err)
	}
	return string(b), nil
}

func jsonToTOML(input string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", fmt.Errorf("JSON parse error: %w", err)
	}
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return "", fmt.Errorf("TOML marshal error: %w", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// ── CSV <-> JSON ───────────────────────────────────────────────────────────

func csvToJSON(input string) (string, error) {
	r := csv.NewReader(strings.NewReader(input))
	records, err := r.ReadAll()
	if err != nil {
		return "", fmt.Errorf("CSV parse error: %w", err)
	}
	if len(records) < 1 {
		return "[]", nil
	}
	headers := records[0]
	var result []map[string]string
	for _, row := range records[1:] {
		obj := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(row) {
				obj[h] = row[i]
			} else {
				obj[h] = ""
			}
		}
		result = append(result, obj)
	}
	if result == nil {
		result = []map[string]string{}
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %w", err)
	}
	return string(b), nil
}

func jsonToCSV(input string) (string, error) {
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", fmt.Errorf("JSON parse error (expected array of objects): %w", err)
	}
	if len(data) == 0 {
		return "", nil
	}

	// Collect all keys for headers, sorted for deterministic output.
	keySet := make(map[string]struct{})
	for _, obj := range data {
		for k := range obj {
			keySet[k] = struct{}{}
		}
	}
	headers := make([]string, 0, len(keySet))
	for k := range keySet {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(headers); err != nil {
		return "", fmt.Errorf("CSV write error: %w", err)
	}

	for _, obj := range data {
		row := make([]string, len(headers))
		for i, h := range headers {
			if v, ok := obj[h]; ok {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("CSV write error: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("CSV flush error: %w", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// ── XML <-> JSON ───────────────────────────────────────────────────────────

// xmlNode is a simple recursive representation of an XML element.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Content  string     `xml:",chardata"`
	Children []xmlNode  `xml:",any"`
}

func xmlToJSON(input string) (string, error) {
	var root xmlNode
	if err := xml.Unmarshal([]byte(input), &root); err != nil {
		return "", fmt.Errorf("XML parse error: %w", err)
	}
	data := xmlNodeToMap(root)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %w", err)
	}
	return string(b), nil
}

// xmlNodeToMap converts an xmlNode tree into a map structure suitable for JSON.
func xmlNodeToMap(n xmlNode) map[string]interface{} {
	result := make(map[string]interface{})

	// Add attributes with a "-" prefix to distinguish them.
	for _, attr := range n.Attrs {
		result["-"+attr.Name.Local] = attr.Value
	}

	if len(n.Children) == 0 {
		// Leaf element: just the text content.
		content := strings.TrimSpace(n.Content)
		if len(n.Attrs) > 0 {
			if content != "" {
				result["#text"] = content
			}
		} else {
			// Simple leaf — return the whole node as {name: value}.
			return map[string]interface{}{n.XMLName.Local: content}
		}
		return map[string]interface{}{n.XMLName.Local: result}
	}

	// Group children by tag name to detect arrays.
	childGroups := make(map[string][]interface{})
	var childOrder []string
	for _, child := range n.Children {
		m := xmlNodeToMap(child)
		name := child.XMLName.Local
		if _, seen := childGroups[name]; !seen {
			childOrder = append(childOrder, name)
		}
		// Extract the inner value from the single-key map.
		if inner, ok := m[name]; ok {
			childGroups[name] = append(childGroups[name], inner)
		} else {
			childGroups[name] = append(childGroups[name], m)
		}
	}

	for _, name := range childOrder {
		group := childGroups[name]
		if len(group) == 1 {
			result[name] = group[0]
		} else {
			result[name] = group
		}
	}

	content := strings.TrimSpace(n.Content)
	if content != "" {
		result["#text"] = content
	}

	return map[string]interface{}{n.XMLName.Local: result}
}

func jsonToXML(input string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", fmt.Errorf("JSON parse error: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	writeXML(&buf, "", data, 0)
	return strings.TrimRight(buf.String(), "\n"), nil
}

func writeXML(w io.Writer, tagName string, data interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch val := data.(type) {
	case map[string]interface{}:
		if tagName == "" {
			// Top-level: iterate keys as root elements.
			// Sort keys for deterministic output.
			keys := sortedKeys(val)
			for _, k := range keys {
				writeXML(w, k, val[k], indent)
			}
			return
		}
		fmt.Fprintf(w, "%s<%s", prefix, tagName)

		// Separate attributes (keys starting with "-") from child elements.
		var attrs, children []string
		for k := range val {
			if strings.HasPrefix(k, "-") {
				attrs = append(attrs, k)
			} else {
				children = append(children, k)
			}
		}
		sort.Strings(attrs)
		sort.Strings(children)

		for _, a := range attrs {
			fmt.Fprintf(w, " %s=%q", a[1:], fmt.Sprintf("%v", val[a]))
		}
		fmt.Fprintf(w, ">\n")

		for _, k := range children {
			if k == "#text" {
				fmt.Fprintf(w, "%s  %v\n", prefix, val[k])
			} else {
				writeXML(w, k, val[k], indent+1)
			}
		}
		fmt.Fprintf(w, "%s</%s>\n", prefix, tagName)

	case []interface{}:
		for _, item := range val {
			writeXML(w, tagName, item, indent)
		}

	default:
		fmt.Fprintf(w, "%s<%s>%v</%s>\n", prefix, tagName, val, tagName)
	}
}

// sortedKeys returns the keys of a map sorted alphabetically.
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ── init ────────────────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent (inherited by all subcommands).
	dataformatCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false,
		`wrap output in {"input_format","output_format","data"}`)
	dataformatCmd.PersistentFlags().StringVarP(&inputFile, "file", "f", "",
		"read input from a file instead of an argument")

	// --output/-o on each subcommand.
	for _, sc := range []*cobra.Command{
		yaml2jsonCmd, json2yamlCmd,
		toml2jsonCmd, json2tomlCmd,
		csv2jsonCmd, json2csvCmd,
		xml2jsonCmd, json2xmlCmd,
	} {
		sc.Flags().StringVarP(&outputFile, "output", "o", "",
			"write output to a file instead of stdout")
	}

	// Register all subcommands.
	dataformatCmd.AddCommand(
		yaml2jsonCmd,
		json2yamlCmd,
		toml2jsonCmd,
		json2tomlCmd,
		csv2jsonCmd,
		json2csvCmd,
		xml2jsonCmd,
		json2xmlCmd,
	)

	cmd.Register(dataformatCmd)
}

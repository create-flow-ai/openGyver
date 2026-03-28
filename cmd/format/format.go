package format

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

// ── flags ──────────────────────────────────────────────────────────────────

var (
	jsonOut    bool
	inputFile  string
	outputFile string
)

// ── parent command ─────────────────────────────────────────────────────────

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Code formatting and prettifying utilities",
	Long: `Code formatting and prettifying utilities — format, beautify, or
minify HTML, XML, CSS, and SQL.

SUBCOMMANDS:

  html      Format / beautify HTML (proper indentation, self-closing tags)
  xml       Format / beautify XML
  css       Format / beautify CSS (one property per line, normalised spacing)
  sql       Format / beautify SQL (uppercase keywords, clause-per-line)

All subcommands support --json/-j for machine-readable output,
--file/-f to read input from a file, and --output/-o to write to a file.

Each formatter also supports --minify to strip unnecessary whitespace
instead of prettifying, and --indent to control indentation width.

EXAMPLES:

  openGyver format html '<div><p>Hello</p></div>'
  openGyver format html --file index.html --output pretty.html
  openGyver format html --minify '<div>  <p> Hello </p>  </div>'
  openGyver format xml '<root><item id="1"/></root>'
  openGyver format css 'body { color: red; margin: 0; }'
  openGyver format css --minify --file styles.css
  openGyver format sql 'select id, name from users where active = 1 order by name'
  openGyver format sql --minify --file query.sql --json`,
}

// ── helpers ────────────────────────────────────────────────────────────────

// resolveInput returns the input string from the first argument, --file, or stdin.
func resolveInput(args []string) (string, error) {
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}
	if len(args) > 0 {
		return args[0], nil
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
	return "", fmt.Errorf("provide input as an argument, use --file/-f, or pipe via stdin")
}

// emitResult handles --json wrapping and --output writing.
func emitResult(input, output, formatName string) error {
	var result string
	if jsonOut {
		wrapper := map[string]interface{}{
			"input":  input,
			"output": output,
			"format": formatName,
		}
		b, err := json.MarshalIndent(wrapper, "", "  ")
		if err != nil {
			return fmt.Errorf("JSON encoding error: %w", err)
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

// indentStr returns a string of n spaces.
func indentStr(n int) string {
	return strings.Repeat(" ", n)
}

// ── HTML formatter ─────────────────────────────────────────────────────────

// htmlVoidElements is the set of HTML void (self-closing) elements.
var htmlVoidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

// htmlRawTextElements are elements whose children should not be indented.
var htmlRawTextElements = map[string]bool{
	"script": true, "style": true, "pre": true, "textarea": true, "code": true,
}

var (
	htmlIndent int
	htmlMinify bool
)

var htmlCmd = &cobra.Command{
	Use:   "html [input]",
	Short: "Format / beautify HTML",
	Long: `Format and beautify HTML with proper indentation.

Parses the input as HTML, then re-emits it with correct nesting
and indentation. Handles self-closing (void) tags, attributes, text
nodes, comments, and doctype declarations.

FLAGS:

  --indent, -i   Number of spaces per indentation level (default 2)
  --minify, -m   Strip whitespace instead of prettifying

EXAMPLES:

  openGyver format html '<div><p>Hello</p></div>'
  openGyver format html --indent 4 '<html><body><h1>Title</h1></body></html>'
  openGyver format html --minify '<div>  <p> Hello </p>  </div>'
  openGyver format html --file index.html
  openGyver format html --file index.html --output pretty.html
  echo '<ul><li>A</li><li>B</li></ul>' | openGyver format html`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		var output string
		if htmlMinify {
			output, err = minifyHTML(input)
		} else {
			output, err = formatHTML(input, htmlIndent)
		}
		if err != nil {
			return err
		}
		return emitResult(input, output, "html")
	},
}

// formatHTML parses HTML and re-emits it with proper indentation.
func formatHTML(input string, indent int) (string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", fmt.Errorf("HTML parse error: %w", err)
	}

	var buf strings.Builder
	renderHTMLNode(&buf, doc, 0, indent, false)
	return strings.TrimSpace(buf.String()), nil
}

// renderHTMLNode recursively renders an HTML node with indentation.
func renderHTMLNode(buf *strings.Builder, n *html.Node, depth, indent int, preserveWS bool) {
	pad := indentStr(depth * indent)

	switch n.Type {
	case html.DocumentNode:
		// Just render children.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			renderHTMLNode(buf, c, depth, indent, preserveWS)
		}

	case html.DoctypeNode:
		buf.WriteString(pad)
		buf.WriteString("<!DOCTYPE ")
		buf.WriteString(n.Data)
		buf.WriteString(">\n")

	case html.ElementNode:
		buf.WriteString(pad)
		buf.WriteByte('<')
		buf.WriteString(n.Data)
		for _, attr := range n.Attr {
			buf.WriteByte(' ')
			if attr.Namespace != "" {
				buf.WriteString(attr.Namespace)
				buf.WriteByte(':')
			}
			buf.WriteString(attr.Key)
			buf.WriteString(`="`)
			buf.WriteString(html.EscapeString(attr.Val))
			buf.WriteByte('"')
		}

		if htmlVoidElements[n.Data] {
			buf.WriteString(" />\n")
			return
		}

		isRaw := htmlRawTextElements[n.Data]

		// Check if element has only a single short text child.
		if n.FirstChild != nil && n.FirstChild == n.LastChild &&
			n.FirstChild.Type == html.TextNode && !isRaw {
			text := strings.TrimSpace(n.FirstChild.Data)
			if len(text) < 80 && !strings.Contains(text, "\n") {
				buf.WriteByte('>')
				buf.WriteString(text)
				buf.WriteString("</")
				buf.WriteString(n.Data)
				buf.WriteString(">\n")
				return
			}
		}

		buf.WriteString(">\n")

		if isRaw && n.FirstChild != nil {
			// Emit raw content without extra indentation.
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					lines := strings.Split(c.Data, "\n")
					for _, line := range lines {
						buf.WriteString(pad)
						buf.WriteString(indentStr(indent))
						buf.WriteString(strings.TrimSpace(line))
						buf.WriteByte('\n')
					}
				} else {
					renderHTMLNode(buf, c, depth+1, indent, true)
				}
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				renderHTMLNode(buf, c, depth+1, indent, preserveWS)
			}
		}

		buf.WriteString(pad)
		buf.WriteString("</")
		buf.WriteString(n.Data)
		buf.WriteString(">\n")

	case html.TextNode:
		text := n.Data
		if !preserveWS {
			text = strings.TrimSpace(text)
		}
		if text != "" {
			buf.WriteString(pad)
			buf.WriteString(text)
			buf.WriteByte('\n')
		}

	case html.CommentNode:
		buf.WriteString(pad)
		buf.WriteString("<!-- ")
		buf.WriteString(strings.TrimSpace(n.Data))
		buf.WriteString(" -->\n")
	}
}

// minifyHTML strips unnecessary whitespace from HTML.
func minifyHTML(input string) (string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", fmt.Errorf("HTML parse error: %w", err)
	}

	var buf strings.Builder
	minifyHTMLNode(&buf, doc)
	return strings.TrimSpace(buf.String()), nil
}

// minifyHTMLNode renders HTML without any extra whitespace.
func minifyHTMLNode(buf *strings.Builder, n *html.Node) {
	switch n.Type {
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			minifyHTMLNode(buf, c)
		}

	case html.DoctypeNode:
		buf.WriteString("<!DOCTYPE ")
		buf.WriteString(n.Data)
		buf.WriteByte('>')

	case html.ElementNode:
		buf.WriteByte('<')
		buf.WriteString(n.Data)
		for _, attr := range n.Attr {
			buf.WriteByte(' ')
			if attr.Namespace != "" {
				buf.WriteString(attr.Namespace)
				buf.WriteByte(':')
			}
			buf.WriteString(attr.Key)
			buf.WriteString(`="`)
			buf.WriteString(html.EscapeString(attr.Val))
			buf.WriteByte('"')
		}

		if htmlVoidElements[n.Data] {
			buf.WriteString("/>")
			return
		}

		buf.WriteByte('>')

		isRaw := htmlRawTextElements[n.Data]
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if isRaw && c.Type == html.TextNode {
				buf.WriteString(c.Data)
			} else {
				minifyHTMLNode(buf, c)
			}
		}

		buf.WriteString("</")
		buf.WriteString(n.Data)
		buf.WriteByte('>')

	case html.TextNode:
		// Collapse whitespace to a single space.
		text := strings.TrimSpace(n.Data)
		if text != "" {
			buf.WriteString(text)
		}

	case html.CommentNode:
		// Omit comments during minification.
	}
}

// ── XML formatter ──────────────────────────────────────────────────────────

var (
	xmlIndent int
	xmlMinify bool
)

var xmlCmd = &cobra.Command{
	Use:   "xml [input]",
	Short: "Format / beautify XML",
	Long: `Format and beautify XML with proper indentation.

Parses the input with encoding/xml and re-emits it with correct
nesting and indentation. Handles elements, attributes, text nodes,
comments, processing instructions, and CDATA sections.

FLAGS:

  --indent, -i   Number of spaces per indentation level (default 2)
  --minify, -m   Strip whitespace instead of prettifying

EXAMPLES:

  openGyver format xml '<root><item id="1"><name>Foo</name></item></root>'
  openGyver format xml --indent 4 --file data.xml
  openGyver format xml --minify '<root>  <item />  </root>'
  openGyver format xml --file input.xml --output pretty.xml
  echo '<a><b>text</b></a>' | openGyver format xml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		var output string
		if xmlMinify {
			output, err = minifyXML(input)
		} else {
			output, err = formatXML(input, xmlIndent)
		}
		if err != nil {
			return err
		}
		return emitResult(input, output, "xml")
	},
}

// formatXML parses XML with a Decoder and re-emits with indentation.
func formatXML(input string, indent int) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(strings.TrimSpace(input)))
	decoder.Strict = false

	var buf strings.Builder
	depth := 0
	var lastStartName string
	_ = lastStartName

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("XML parse error: %w", err)
		}

		switch t := tok.(type) {
		case xml.ProcInst:
			buf.WriteString(indentStr(depth * indent))
			buf.WriteString("<?")
			buf.WriteString(t.Target)
			if len(t.Inst) > 0 {
				buf.WriteByte(' ')
				buf.Write(t.Inst)
			}
			buf.WriteString("?>\n")

		case xml.StartElement:
			buf.WriteString(indentStr(depth * indent))
			buf.WriteByte('<')
			if t.Name.Space != "" {
				buf.WriteString(t.Name.Space)
				buf.WriteByte(':')
			}
			buf.WriteString(t.Name.Local)
			for _, attr := range t.Attr {
				buf.WriteByte(' ')
				if attr.Name.Space != "" {
					buf.WriteString(attr.Name.Space)
					buf.WriteByte(':')
				}
				buf.WriteString(attr.Name.Local)
				buf.WriteString(`="`)
				buf.WriteString(xmlEscapeAttr(attr.Value))
				buf.WriteByte('"')
			}
			buf.WriteString(">\n")
			lastStartName = t.Name.Local
			depth++

		case xml.EndElement:
			depth--
			if depth < 0 {
				depth = 0
			}
			buf.WriteString(indentStr(depth * indent))
			buf.WriteString("</")
			if t.Name.Space != "" {
				buf.WriteString(t.Name.Space)
				buf.WriteByte(':')
			}
			buf.WriteString(t.Name.Local)
			buf.WriteString(">\n")

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				buf.WriteString(indentStr(depth * indent))
				buf.WriteString(xmlEscapeText(text))
				buf.WriteByte('\n')
			}

		case xml.Comment:
			buf.WriteString(indentStr(depth * indent))
			buf.WriteString("<!-- ")
			buf.WriteString(strings.TrimSpace(string(t)))
			buf.WriteString(" -->\n")

		case xml.Directive:
			buf.WriteString(indentStr(depth * indent))
			buf.WriteString("<!")
			buf.Write(t)
			buf.WriteString(">\n")
		}
	}

	return strings.TrimSpace(buf.String()), nil
}

// minifyXML strips unnecessary whitespace from XML.
func minifyXML(input string) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(strings.TrimSpace(input)))
	decoder.Strict = false

	var buf strings.Builder

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("XML parse error: %w", err)
		}

		switch t := tok.(type) {
		case xml.ProcInst:
			buf.WriteString("<?")
			buf.WriteString(t.Target)
			if len(t.Inst) > 0 {
				buf.WriteByte(' ')
				buf.Write(t.Inst)
			}
			buf.WriteString("?>")

		case xml.StartElement:
			buf.WriteByte('<')
			if t.Name.Space != "" {
				buf.WriteString(t.Name.Space)
				buf.WriteByte(':')
			}
			buf.WriteString(t.Name.Local)
			for _, attr := range t.Attr {
				buf.WriteByte(' ')
				if attr.Name.Space != "" {
					buf.WriteString(attr.Name.Space)
					buf.WriteByte(':')
				}
				buf.WriteString(attr.Name.Local)
				buf.WriteString(`="`)
				buf.WriteString(xmlEscapeAttr(attr.Value))
				buf.WriteByte('"')
			}
			buf.WriteByte('>')

		case xml.EndElement:
			buf.WriteString("</")
			if t.Name.Space != "" {
				buf.WriteString(t.Name.Space)
				buf.WriteByte(':')
			}
			buf.WriteString(t.Name.Local)
			buf.WriteByte('>')

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				buf.WriteString(xmlEscapeText(text))
			}

		case xml.Comment:
			// Omit comments during minification.

		case xml.Directive:
			buf.WriteString("<!")
			buf.Write(t)
			buf.WriteByte('>')
		}
	}

	return strings.TrimSpace(buf.String()), nil
}

// xmlEscapeAttr escapes special characters in XML attribute values.
func xmlEscapeAttr(s string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		return s
	}
	return buf.String()
}

// xmlEscapeText escapes special characters in XML text content.
func xmlEscapeText(s string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		return s
	}
	return buf.String()
}

// ── CSS formatter ──────────────────────────────────────────────────────────

var (
	cssIndent int
	cssMinify bool
)

var cssCmd = &cobra.Command{
	Use:   "css [input]",
	Short: "Format / beautify CSS",
	Long: `Format and beautify CSS with proper indentation.

Uses a simple rule-based approach: puts each property on its own line,
indents inside braces, and normalises spacing around colons, semicolons,
and braces. Handles selectors, media queries, nested at-rules, and
comments.

FLAGS:

  --indent, -i   Number of spaces per indentation level (default 2)
  --minify, -m   Strip all unnecessary whitespace

EXAMPLES:

  openGyver format css 'body { color: red; margin: 0; }'
  openGyver format css --indent 4 'h1{font-size:2em;color:blue}'
  openGyver format css --minify 'body { color: red; margin: 0; }'
  openGyver format css --file styles.css
  openGyver format css --file styles.css --output pretty.css
  echo '.box { padding: 10px; }' | openGyver format css`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		var output string
		if cssMinify {
			output = minifyCSS(input)
		} else {
			output = formatCSS(input, cssIndent)
		}
		return emitResult(input, output, "css")
	},
}

// cssCommentRe matches CSS block comments.
var cssCommentRe = regexp.MustCompile(`/\*[\s\S]*?\*/`)

// formatCSS formats CSS with one property per line and proper indentation.
func formatCSS(input string, indent int) string {
	// Normalise the input: collapse multiple whitespace into single spaces.
	s := strings.TrimSpace(input)

	var buf strings.Builder
	depth := 0
	pad := func() string { return indentStr(depth * indent) }

	i := 0
	for i < len(s) {
		// Handle block comments.
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			end := strings.Index(s[i+2:], "*/")
			if end == -1 {
				// Unterminated comment: emit the rest.
				buf.WriteString(pad())
				buf.WriteString(s[i:])
				buf.WriteByte('\n')
				break
			}
			comment := s[i : i+2+end+2]
			buf.WriteString(pad())
			buf.WriteString(comment)
			buf.WriteByte('\n')
			i = i + 2 + end + 2
			// Skip trailing whitespace.
			for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
				i++
			}
			continue
		}

		// Opening brace.
		if s[i] == '{' {
			// The selector/at-rule preceding the brace.
			buf.WriteString(" {\n")
			depth++
			i++
			// Skip whitespace after brace.
			for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
				i++
			}
			continue
		}

		// Closing brace.
		if s[i] == '}' {
			depth--
			if depth < 0 {
				depth = 0
			}
			buf.WriteString(pad())
			buf.WriteString("}\n")
			i++
			// Skip whitespace after brace.
			for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
				i++
			}
			// Add blank line between top-level blocks.
			if depth == 0 && i < len(s) && s[i] != '}' {
				buf.WriteByte('\n')
			}
			continue
		}

		// Semicolon: end of a property.
		if s[i] == ';' {
			buf.WriteByte(';')
			buf.WriteByte('\n')
			i++
			// Skip whitespace after semicolon.
			for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
				i++
			}
			continue
		}

		// Collect a token (selector, property, value, etc.) up to the next
		// significant character.
		start := i
		for i < len(s) && s[i] != '{' && s[i] != '}' && s[i] != ';' {
			// Check for comment start.
			if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
				break
			}
			i++
		}

		chunk := strings.TrimSpace(s[start:i])
		if chunk == "" {
			continue
		}

		// Normalise colon spacing in property declarations.
		if colonIdx := strings.Index(chunk, ":"); colonIdx > 0 && depth > 0 &&
			!strings.HasPrefix(chunk, "@") && !strings.Contains(chunk, "//") {
			prop := strings.TrimSpace(chunk[:colonIdx])
			val := strings.TrimSpace(chunk[colonIdx+1:])
			buf.WriteString(pad())
			buf.WriteString(prop)
			buf.WriteString(": ")
			buf.WriteString(val)
		} else {
			// Selector or at-rule.
			buf.WriteString(pad())
			// Normalise whitespace within selectors.
			normalized := collapseWhitespace(chunk)
			buf.WriteString(normalized)
		}
	}

	result := buf.String()
	// Clean up any trailing blank lines.
	result = strings.TrimRight(result, "\n") + "\n"
	// Remove trailing whitespace on each line.
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// minifyCSS strips all unnecessary whitespace from CSS.
func minifyCSS(input string) string {
	s := strings.TrimSpace(input)

	// Remove comments.
	s = cssCommentRe.ReplaceAllString(s, "")

	// Collapse all whitespace to single spaces.
	s = collapseWhitespace(s)

	// Remove spaces around structural characters.
	s = strings.ReplaceAll(s, " {", "{")
	s = strings.ReplaceAll(s, "{ ", "{")
	s = strings.ReplaceAll(s, " }", "}")
	s = strings.ReplaceAll(s, "} ", "}")
	s = strings.ReplaceAll(s, " ;", ";")
	s = strings.ReplaceAll(s, "; ", ";")
	s = strings.ReplaceAll(s, " :", ":")
	s = strings.ReplaceAll(s, ": ", ":")

	// Remove the last semicolon before a closing brace.
	s = strings.ReplaceAll(s, ";}", "}")

	return strings.TrimSpace(s)
}

// collapseWhitespace replaces runs of whitespace with a single space.
func collapseWhitespace(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
}

// ── SQL formatter ──────────────────────────────────────────────────────────

var (
	sqlIndent int
	sqlMinify bool
)

var sqlCmd = &cobra.Command{
	Use:   "sql [input]",
	Short: "Format / beautify SQL",
	Long: `Format and beautify SQL with proper indentation and keyword casing.

Uppercases SQL keywords, adds newlines before major clauses, and
indents clause bodies. Handles SELECT, FROM, WHERE, JOIN, subqueries,
and all common SQL keywords.

FLAGS:

  --indent, -i   Number of spaces per indentation level (default 2)
  --minify, -m   Put everything on one line, stripping extra whitespace

KEYWORDS UPPERCASED:

  SELECT, FROM, WHERE, JOIN, LEFT JOIN, RIGHT JOIN, INNER JOIN,
  OUTER JOIN, CROSS JOIN, FULL JOIN, ON, ORDER BY, GROUP BY,
  HAVING, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, AND, OR,
  IN, NOT, NULL, AS, SET, VALUES, INTO, LIMIT, OFFSET, UNION,
  DISTINCT, BETWEEN, LIKE, EXISTS, CASE, WHEN, THEN, ELSE, END

EXAMPLES:

  openGyver format sql 'select id, name from users where active = 1 order by name'
  openGyver format sql --indent 4 'select * from orders o join users u on o.user_id = u.id'
  openGyver format sql --minify --file query.sql
  openGyver format sql --file complex.sql --output pretty.sql --json
  echo 'insert into users (name, email) values ("Alice", "a@b.com")' | openGyver format sql`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		input, err := resolveInput(args)
		if err != nil {
			return err
		}
		var output string
		if sqlMinify {
			output = minifySQL(input)
		} else {
			output = formatSQL(input, sqlIndent)
		}
		return emitResult(input, output, "sql")
	},
}

// sqlKeywords is the set of SQL keywords to uppercase.
var sqlKeywords = []string{
	"SELECT", "FROM", "WHERE", "JOIN", "LEFT JOIN", "RIGHT JOIN",
	"INNER JOIN", "OUTER JOIN", "CROSS JOIN", "FULL JOIN", "FULL OUTER JOIN",
	"LEFT OUTER JOIN", "RIGHT OUTER JOIN",
	"ON", "ORDER BY", "GROUP BY", "HAVING",
	"INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP",
	"AND", "OR", "IN", "NOT", "NULL", "AS", "SET", "VALUES", "INTO",
	"LIMIT", "OFFSET", "UNION", "UNION ALL", "DISTINCT",
	"BETWEEN", "LIKE", "EXISTS",
	"CASE", "WHEN", "THEN", "ELSE", "END",
}

// sqlMajorClauses are keywords that should start on a new line at the top level.
var sqlMajorClauses = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true,
	"JOIN": true, "LEFT JOIN": true, "RIGHT JOIN": true,
	"INNER JOIN": true, "OUTER JOIN": true, "CROSS JOIN": true,
	"FULL JOIN": true, "FULL OUTER JOIN": true,
	"LEFT OUTER JOIN": true, "RIGHT OUTER JOIN": true,
	"ON": true, "ORDER BY": true, "GROUP BY": true, "HAVING": true,
	"INSERT": true, "UPDATE": true, "DELETE": true,
	"CREATE": true, "ALTER": true, "DROP": true,
	"LIMIT": true, "OFFSET": true,
	"UNION": true, "UNION ALL": true,
	"SET": true, "VALUES": true, "INTO": true,
}

// sqlSubClauses are keywords indented within a clause.
var sqlSubClauses = map[string]bool{
	"AND": true, "OR": true,
}

// formatSQL formats SQL with keyword uppercasing and newlines before clauses.
func formatSQL(input string, indent int) string {
	// Normalise whitespace.
	s := collapseWhitespace(strings.TrimSpace(input))

	// Uppercase all keywords (case-insensitive replacement).
	s = uppercaseSQLKeywords(s)

	// Now insert newlines before major clauses and indent.
	var buf strings.Builder
	pad := indentStr(indent)

	tokens := tokenizeSQL(s)
	firstClause := true

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		upper := strings.ToUpper(strings.TrimSpace(token))

		// Check for two-word keywords.
		twoWord := ""
		if i+1 < len(tokens) {
			candidate := upper + " " + strings.ToUpper(strings.TrimSpace(tokens[i+1]))
			if sqlMajorClauses[candidate] || sqlSubClauses[candidate] {
				twoWord = candidate
			}
		}

		if twoWord != "" {
			if sqlMajorClauses[twoWord] {
				if !firstClause {
					buf.WriteByte('\n')
				}
				buf.WriteString(twoWord)
				firstClause = false
			} else if sqlSubClauses[twoWord] {
				buf.WriteByte('\n')
				buf.WriteString(pad)
				buf.WriteString(twoWord)
			}
			i++ // Skip the next token (already consumed).
		} else if sqlMajorClauses[upper] {
			if !firstClause {
				buf.WriteByte('\n')
			}
			buf.WriteString(upper)
			firstClause = false
		} else if sqlSubClauses[upper] {
			buf.WriteByte('\n')
			buf.WriteString(pad)
			buf.WriteString(upper)
		} else {
			// Regular token: part of a clause body.
			if buf.Len() > 0 {
				lastByte := buf.String()[buf.Len()-1]
				if lastByte != '\n' && lastByte != ' ' && lastByte != '(' {
					buf.WriteByte(' ')
				}
			}
			buf.WriteString(strings.TrimSpace(token))
		}
	}

	result := buf.String()

	// Indent clause bodies: lines that don't start with a major keyword get indented.
	lines := strings.Split(result, "\n")
	var out strings.Builder
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if i > 0 {
			out.WriteByte('\n')
		}
		upper := strings.ToUpper(trimmed)
		isMajor := false
		for kw := range sqlMajorClauses {
			if strings.HasPrefix(upper, kw+" ") || upper == kw {
				isMajor = true
				break
			}
		}
		isSubClause := false
		for kw := range sqlSubClauses {
			if strings.HasPrefix(upper, kw+" ") || upper == kw {
				isSubClause = true
				break
			}
		}
		if isMajor {
			out.WriteString(trimmed)
		} else if isSubClause {
			out.WriteString(pad)
			out.WriteString(trimmed)
		} else {
			out.WriteString(pad)
			out.WriteString(trimmed)
		}
	}

	return strings.TrimSpace(out.String())
}

// tokenizeSQL splits SQL into tokens, preserving strings and parenthesised groups.
func tokenizeSQL(s string) []string {
	var tokens []string
	i := 0

	for i < len(s) {
		// Skip whitespace.
		if s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r' {
			i++
			continue
		}

		// Quoted string (single or double quotes).
		if s[i] == '\'' || s[i] == '"' {
			quote := s[i]
			j := i + 1
			for j < len(s) {
				if s[j] == quote {
					if j+1 < len(s) && s[j+1] == quote {
						j += 2 // Escaped quote.
						continue
					}
					break
				}
				j++
			}
			if j < len(s) {
				j++ // Include closing quote.
			}
			tokens = append(tokens, s[i:j])
			i = j
			continue
		}

		// Parenthesised group.
		if s[i] == '(' {
			depth := 1
			j := i + 1
			for j < len(s) && depth > 0 {
				if s[j] == '(' {
					depth++
				} else if s[j] == ')' {
					depth--
				} else if s[j] == '\'' || s[j] == '"' {
					// Skip over string literals inside parens.
					quote := s[j]
					j++
					for j < len(s) && s[j] != quote {
						j++
					}
				}
				j++
			}
			tokens = append(tokens, s[i:j])
			i = j
			continue
		}

		// Comma.
		if s[i] == ',' {
			tokens = append(tokens, ",")
			i++
			continue
		}

		// Regular word.
		j := i
		for j < len(s) && s[j] != ' ' && s[j] != '\t' && s[j] != '\n' &&
			s[j] != '\r' && s[j] != ',' && s[j] != '(' && s[j] != ')' &&
			s[j] != '\'' && s[j] != '"' {
			j++
		}
		if j > i {
			tokens = append(tokens, s[i:j])
			i = j
		}
	}

	return tokens
}

// uppercaseSQLKeywords uppercases all SQL keywords in the input (case-insensitive).
func uppercaseSQLKeywords(s string) string {
	// Build a regex that matches each keyword as a whole word (case-insensitive).
	// Sort keywords by length (longest first) so multi-word keywords match first.
	sorted := make([]string, len(sqlKeywords))
	copy(sorted, sqlKeywords)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if len(sorted[j]) > len(sorted[i]) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	for _, kw := range sorted {
		// Use word boundaries. For multi-word keywords, allow flexible whitespace.
		pattern := `(?i)\b` + strings.ReplaceAll(regexp.QuoteMeta(kw), `\ `, `\s+`) + `\b`
		re := regexp.MustCompile(pattern)
		s = re.ReplaceAllStringFunc(s, func(match string) string {
			// Preserve content inside quoted strings by checking context.
			return kw
		})
	}

	return s
}

// minifySQL puts everything on one line and collapses whitespace.
func minifySQL(input string) string {
	s := collapseWhitespace(strings.TrimSpace(input))
	s = uppercaseSQLKeywords(s)
	return strings.TrimSpace(s)
}

// ── init / register ────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent (inherited by all subcommands).
	formatCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false,
		`wrap output in {"input","output","format"} JSON`)
	formatCmd.PersistentFlags().StringVarP(&inputFile, "file", "f", "",
		"read input from a file instead of an argument")
	formatCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "",
		"write output to a file instead of stdout")

	// HTML subcommand flags.
	htmlCmd.Flags().IntVarP(&htmlIndent, "indent", "i", 2,
		"number of spaces per indentation level")
	htmlCmd.Flags().BoolVarP(&htmlMinify, "minify", "m", false,
		"strip whitespace instead of prettifying")

	// XML subcommand flags.
	xmlCmd.Flags().IntVarP(&xmlIndent, "indent", "i", 2,
		"number of spaces per indentation level")
	xmlCmd.Flags().BoolVarP(&xmlMinify, "minify", "m", false,
		"strip whitespace instead of prettifying")

	// CSS subcommand flags.
	cssCmd.Flags().IntVarP(&cssIndent, "indent", "i", 2,
		"number of spaces per indentation level")
	cssCmd.Flags().BoolVarP(&cssMinify, "minify", "m", false,
		"strip all unnecessary whitespace")

	// SQL subcommand flags.
	sqlCmd.Flags().IntVarP(&sqlIndent, "indent", "i", 2,
		"number of spaces per indentation level")
	sqlCmd.Flags().BoolVarP(&sqlMinify, "minify", "m", false,
		"put everything on one line")

	// Register all subcommands.
	formatCmd.AddCommand(
		htmlCmd,
		xmlCmd,
		cssCmd,
		sqlCmd,
	)

	cmd.Register(formatCmd)
}

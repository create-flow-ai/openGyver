package format

import (
	"strings"
	"testing"
)

func TestFormatSQL_KeywordUppercase(t *testing.T) {
	input := "select id, name from users where active = 1 order by name"
	result := formatSQL(input, 2)

	if !strings.Contains(result, "SELECT") {
		t.Error("expected SELECT to be uppercased")
	}
	if !strings.Contains(result, "FROM") {
		t.Error("expected FROM to be uppercased")
	}
	if !strings.Contains(result, "WHERE") {
		t.Error("expected WHERE to be uppercased")
	}
	if !strings.Contains(result, "ORDER BY") {
		t.Error("expected ORDER BY to be uppercased")
	}
}

func TestFormatSQL_AddsNewlines(t *testing.T) {
	input := "select id from users where active = 1"
	result := formatSQL(input, 2)

	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d: %s", len(lines), result)
	}
}

func TestMinifySQL(t *testing.T) {
	input := "SELECT id, name\nFROM users\nWHERE active = 1\nORDER BY name"
	result := minifySQL(input)

	if strings.Contains(result, "\n") {
		t.Error("minified SQL should not contain newlines")
	}
	// Should still have keywords uppercased.
	if !strings.Contains(result, "SELECT") {
		t.Error("minified SQL should contain uppercased keywords")
	}
}

func TestFormatCSS_AddsIndentation(t *testing.T) {
	input := "body { color: red; margin: 0; }"
	result := formatCSS(input, 2)

	if !strings.Contains(result, "  color: red") {
		t.Errorf("expected indented color property, got:\n%s", result)
	}
}

func TestFormatCSS_OnePropertyPerLine(t *testing.T) {
	input := "h1{font-size:2em;color:blue}"
	result := formatCSS(input, 2)

	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Errorf("expected multiple lines, got %d:\n%s", len(lines), result)
	}
}

func TestMinifyCSS(t *testing.T) {
	input := "body {\n  color: red;\n  margin: 0;\n}"
	result := minifyCSS(input)

	if strings.Contains(result, "\n") {
		t.Error("minified CSS should not contain newlines")
	}
	if strings.Contains(result, "  ") {
		t.Error("minified CSS should not contain double spaces")
	}
}

func TestFormatHTML(t *testing.T) {
	input := "<div><p>Hello</p></div>"
	result, err := formatHTML(input, 2)
	if err != nil {
		t.Fatalf("formatHTML error: %v", err)
	}
	if !strings.Contains(result, "<div>") {
		t.Error("formatted HTML should contain <div>")
	}
	if !strings.Contains(result, "<p>") {
		t.Error("formatted HTML should contain <p>")
	}
}

func TestMinifyHTML(t *testing.T) {
	input := "<div>  <p> Hello </p>  </div>"
	result, err := minifyHTML(input)
	if err != nil {
		t.Fatalf("minifyHTML error: %v", err)
	}
	// Minified should not have extra whitespace.
	if strings.Contains(result, "  ") {
		t.Errorf("minified HTML should not have double spaces: %s", result)
	}
}

func TestFormatXML(t *testing.T) {
	input := "<root><item id=\"1\"><name>Foo</name></item></root>"
	result, err := formatXML(input, 2)
	if err != nil {
		t.Fatalf("formatXML error: %v", err)
	}
	if !strings.Contains(result, "\n") {
		t.Error("formatted XML should contain newlines")
	}
	if !strings.Contains(result, "<root>") {
		t.Error("formatted XML should contain <root>")
	}
}

func TestMinifyXML(t *testing.T) {
	input := "<root>\n  <item>Test</item>\n</root>"
	result, err := minifyXML(input)
	if err != nil {
		t.Fatalf("minifyXML error: %v", err)
	}
	if strings.Contains(result, "\n") {
		t.Error("minified XML should not contain newlines")
	}
}

func TestCollapseWhitespace(t *testing.T) {
	input := "hello   world\n\tfoo   bar"
	result := collapseWhitespace(input)
	if result != "hello world foo bar" {
		t.Errorf("collapseWhitespace = %q, want %q", result, "hello world foo bar")
	}
}

func TestUppercaseSQLKeywords(t *testing.T) {
	input := "select id from users where active = 1"
	result := uppercaseSQLKeywords(input)
	if !strings.Contains(result, "SELECT") {
		t.Error("expected SELECT")
	}
	if !strings.Contains(result, "FROM") {
		t.Error("expected FROM")
	}
	if !strings.Contains(result, "WHERE") {
		t.Error("expected WHERE")
	}
}

func TestTokenizeSQL(t *testing.T) {
	tokens := tokenizeSQL("SELECT id, name FROM users")
	if len(tokens) == 0 {
		t.Fatal("expected some tokens")
	}
	// Should contain SELECT, id, comma, name, FROM, users.
	found := false
	for _, tok := range tokens {
		if tok == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find SELECT in tokens: %v", tokens)
	}
}

func TestFormatSQL_JoinClause(t *testing.T) {
	input := "select * from orders o join users u on o.user_id = u.id"
	result := formatSQL(input, 2)

	if !strings.Contains(result, "JOIN") {
		t.Error("expected JOIN to appear")
	}
	if !strings.Contains(result, "ON") {
		t.Error("expected ON to appear")
	}
}

func TestFormatCSS_WithComments(t *testing.T) {
	input := "/* reset */ body { margin: 0; }"
	result := formatCSS(input, 2)
	if !strings.Contains(result, "/* reset */") {
		t.Errorf("expected comment to be preserved: %s", result)
	}
}

func TestIndentStr(t *testing.T) {
	if indentStr(0) != "" {
		t.Error("indentStr(0) should be empty")
	}
	if indentStr(4) != "    " {
		t.Errorf("indentStr(4) = %q, want 4 spaces", indentStr(4))
	}
}

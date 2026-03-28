package convertfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Markdown → HTML
// ---------------------------------------------------------------------------

func TestMDToHTML(t *testing.T) {
	dir := tempDir(t)
	mdPath := filepath.Join(dir, "test.md")
	htmlPath := filepath.Join(dir, "test.html")

	os.WriteFile(mdPath, []byte("# Hello\n\nWorld **bold**\n"), 0644)

	err := mdToHTML(ConvertOpts{InputPath: mdPath, OutputPath: htmlPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(htmlPath)
	html := string(data)
	if !strings.Contains(html, "<h1>Hello</h1>") {
		t.Error("expected <h1> tag")
	}
	if !strings.Contains(html, "<strong>bold</strong>") {
		t.Error("expected <strong> tag")
	}
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("expected HTML wrapper")
	}
}

// ---------------------------------------------------------------------------
// HTML → Markdown
// ---------------------------------------------------------------------------

func TestHTMLToMD(t *testing.T) {
	dir := tempDir(t)
	htmlPath := filepath.Join(dir, "test.html")
	mdPath := filepath.Join(dir, "test.md")

	os.WriteFile(htmlPath, []byte("<h1>Title</h1><p>Hello <strong>world</strong></p>"), 0644)

	err := htmlToMD(ConvertOpts{InputPath: htmlPath, OutputPath: mdPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(mdPath)
	md := string(data)
	if !strings.Contains(md, "# Title") {
		t.Error("expected Markdown heading")
	}
	if !strings.Contains(md, "**world**") {
		t.Error("expected bold Markdown")
	}
}

// ---------------------------------------------------------------------------
// Markdown → Text
// ---------------------------------------------------------------------------

func TestMDToText(t *testing.T) {
	dir := tempDir(t)
	mdPath := filepath.Join(dir, "test.md")
	txtPath := filepath.Join(dir, "test.txt")

	os.WriteFile(mdPath, []byte("# Header\n\n**bold** and *italic* and `code`\n\n[link](http://ex.com)\n"), 0644)

	err := mdToText(ConvertOpts{InputPath: mdPath, OutputPath: txtPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(txtPath)
	text := string(data)
	if strings.Contains(text, "#") {
		t.Error("should strip headers")
	}
	if strings.Contains(text, "**") {
		t.Error("should strip bold markers")
	}
	if strings.Contains(text, "`") {
		t.Error("should strip code markers")
	}
	if !strings.Contains(text, "bold") {
		t.Error("should preserve text content")
	}
	if !strings.Contains(text, "link") {
		t.Error("should preserve link text")
	}
}

// ---------------------------------------------------------------------------
// HTML → Text
// ---------------------------------------------------------------------------

func TestHTMLToText(t *testing.T) {
	dir := tempDir(t)
	htmlPath := filepath.Join(dir, "test.html")
	txtPath := filepath.Join(dir, "test.txt")

	os.WriteFile(htmlPath, []byte("<html><body><h1>Title</h1><p>Hello &amp; world</p></body></html>"), 0644)

	err := htmlToText(ConvertOpts{InputPath: htmlPath, OutputPath: txtPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(txtPath)
	text := string(data)
	if strings.Contains(text, "<") {
		t.Error("should strip HTML tags")
	}
	if !strings.Contains(text, "Title") {
		t.Error("should preserve text")
	}
	if !strings.Contains(text, "Hello & world") {
		t.Error("should decode HTML entities")
	}
}

// ---------------------------------------------------------------------------
// Text → HTML
// ---------------------------------------------------------------------------

func TestTextToHTML(t *testing.T) {
	dir := tempDir(t)
	txtPath := filepath.Join(dir, "test.txt")
	htmlPath := filepath.Join(dir, "test.html")

	os.WriteFile(txtPath, []byte("Hello <world> & friends\n"), 0644)

	err := textToHTML(ConvertOpts{InputPath: txtPath, OutputPath: htmlPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(htmlPath)
	html := string(data)
	if !strings.Contains(html, "<pre>") {
		t.Error("expected <pre> wrapper")
	}
	if !strings.Contains(html, "&lt;world&gt;") {
		t.Error("should escape HTML entities")
	}
	if !strings.Contains(html, "&amp;") {
		t.Error("should escape ampersand")
	}
}

// ---------------------------------------------------------------------------
// Text → Markdown
// ---------------------------------------------------------------------------

func TestTextToMD(t *testing.T) {
	dir := tempDir(t)
	txtPath := filepath.Join(dir, "test.txt")
	mdPath := filepath.Join(dir, "test.md")

	os.WriteFile(txtPath, []byte("some code here\n"), 0644)

	err := textToMD(ConvertOpts{InputPath: txtPath, OutputPath: mdPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(mdPath)
	md := string(data)
	if !strings.Contains(md, "```") {
		t.Error("expected code fence")
	}
	if !strings.Contains(md, "some code here") {
		t.Error("should preserve content")
	}
}

// ---------------------------------------------------------------------------
// Roundtrip: MD → HTML → MD
// ---------------------------------------------------------------------------

func TestRoundtrip_MD_HTML_MD(t *testing.T) {
	dir := tempDir(t)
	mdIn := filepath.Join(dir, "in.md")
	htmlMid := filepath.Join(dir, "mid.html")
	mdOut := filepath.Join(dir, "out.md")

	os.WriteFile(mdIn, []byte("# Title\n\nParagraph with **bold**.\n"), 0644)

	if err := mdToHTML(ConvertOpts{InputPath: mdIn, OutputPath: htmlMid}); err != nil {
		t.Fatal(err)
	}
	if err := htmlToMD(ConvertOpts{InputPath: htmlMid, OutputPath: mdOut}); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(mdOut)
	md := string(data)
	if !strings.Contains(md, "# Title") {
		t.Error("heading lost in roundtrip")
	}
	if !strings.Contains(md, "**bold**") {
		t.Error("bold lost in roundtrip")
	}
}

// ---------------------------------------------------------------------------
// stripMarkdown
// ---------------------------------------------------------------------------

func TestStripMarkdown(t *testing.T) {
	input := "# Header\n\n**bold** *italic* `code` [link](url)\n\n- item\n1. num\n> quote\n---\n"
	result := stripMarkdown(input)

	checks := map[string]bool{
		"#":    false,
		"**":   false,
		"*":    false,
		"`":    false,
		"[":    false,
		"](":   false,
		"> ":   false,
		"---":  false,
		"bold": true,
		"code": true,
		"link": true,
	}
	for substr, shouldContain := range checks {
		if strings.Contains(result, substr) != shouldContain {
			if shouldContain {
				t.Errorf("expected %q in result", substr)
			} else {
				t.Errorf("should have stripped %q from result", substr)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// stripHTML
// ---------------------------------------------------------------------------

func TestStripHTML(t *testing.T) {
	input := "<html><body><h1>Title</h1><p>Hello &amp; world</p><br><div>More</div></body></html>"
	result := stripHTML(input)

	if strings.Contains(result, "<") {
		t.Error("tags not stripped")
	}
	if !strings.Contains(result, "Title") {
		t.Error("text lost")
	}
	if !strings.Contains(result, "Hello & world") {
		t.Error("entities not decoded")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestMDToHTML_MissingFile(t *testing.T) {
	err := mdToHTML(ConvertOpts{InputPath: "/nonexistent.md", OutputPath: "/tmp/out.html"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestHTMLToMD_MissingFile(t *testing.T) {
	err := htmlToMD(ConvertOpts{InputPath: "/nonexistent.html", OutputPath: "/tmp/out.md"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

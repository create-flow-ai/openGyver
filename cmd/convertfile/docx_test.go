package convertfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Text → DOCX → Text roundtrip
// ---------------------------------------------------------------------------

func TestTextToDOCX(t *testing.T) {
	dir := tempDir(t)
	txtPath := filepath.Join(dir, "test.txt")
	docxPath := filepath.Join(dir, "test.docx")

	os.WriteFile(txtPath, []byte("Hello World\nSecond Line\n"), 0644)

	err := textToDOCX(ConvertOpts{InputPath: txtPath, OutputPath: docxPath})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(docxPath); os.IsNotExist(err) {
		t.Fatal("DOCX file not created")
	}
}

func TestDOCXToText(t *testing.T) {
	dir := tempDir(t)
	txtIn := filepath.Join(dir, "in.txt")
	docxPath := filepath.Join(dir, "test.docx")
	txtOut := filepath.Join(dir, "out.txt")

	os.WriteFile(txtIn, []byte("Hello World\nSecond Line\n"), 0644)

	if err := textToDOCX(ConvertOpts{InputPath: txtIn, OutputPath: docxPath}); err != nil {
		t.Fatal(err)
	}
	if err := docxToText(ConvertOpts{InputPath: docxPath, OutputPath: txtOut}); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(txtOut)
	text := string(data)
	if !strings.Contains(text, "Hello World") {
		t.Error("text not preserved in roundtrip")
	}
	if !strings.Contains(text, "Second Line") {
		t.Error("second line not preserved")
	}
}

// ---------------------------------------------------------------------------
// DOCX → HTML
// ---------------------------------------------------------------------------

func TestDOCXToHTML(t *testing.T) {
	dir := tempDir(t)
	txtIn := filepath.Join(dir, "in.txt")
	docxPath := filepath.Join(dir, "test.docx")
	htmlPath := filepath.Join(dir, "test.html")

	os.WriteFile(txtIn, []byte("Paragraph one\nParagraph two\n"), 0644)

	textToDOCX(ConvertOpts{InputPath: txtIn, OutputPath: docxPath})

	err := docxToHTML(ConvertOpts{InputPath: docxPath, OutputPath: htmlPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(htmlPath)
	html := string(data)
	if !strings.Contains(html, "<p>") {
		t.Error("expected <p> tags")
	}
	if !strings.Contains(html, "Paragraph one") {
		t.Error("text not preserved")
	}
}

// ---------------------------------------------------------------------------
// DOCX → Markdown
// ---------------------------------------------------------------------------

func TestDOCXToMD(t *testing.T) {
	dir := tempDir(t)
	txtIn := filepath.Join(dir, "in.txt")
	docxPath := filepath.Join(dir, "test.docx")
	mdPath := filepath.Join(dir, "test.md")

	os.WriteFile(txtIn, []byte("First para\nSecond para\n"), 0644)

	textToDOCX(ConvertOpts{InputPath: txtIn, OutputPath: docxPath})

	err := docxToMD(ConvertOpts{InputPath: docxPath, OutputPath: mdPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(mdPath)
	if !strings.Contains(string(data), "First para") {
		t.Error("text not preserved")
	}
}

// ---------------------------------------------------------------------------
// MD → DOCX
// ---------------------------------------------------------------------------

func TestMDToDOCX(t *testing.T) {
	dir := tempDir(t)
	mdPath := filepath.Join(dir, "test.md")
	docxPath := filepath.Join(dir, "test.docx")

	os.WriteFile(mdPath, []byte("# Title\n\n**Bold** text\n"), 0644)

	err := mdToDOCX(ConvertOpts{InputPath: mdPath, OutputPath: docxPath})
	if err != nil {
		t.Fatal(err)
	}

	// Verify we can read it back
	paragraphs, err := readDOCXParagraphs(docxPath)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(paragraphs, " ")
	if !strings.Contains(joined, "Title") {
		t.Error("title not in DOCX")
	}
	if !strings.Contains(joined, "Bold") {
		t.Error("bold text not in DOCX")
	}
}

// ---------------------------------------------------------------------------
// HTML → DOCX
// ---------------------------------------------------------------------------

func TestHTMLToDOCX(t *testing.T) {
	dir := tempDir(t)
	htmlPath := filepath.Join(dir, "test.html")
	docxPath := filepath.Join(dir, "test.docx")

	os.WriteFile(htmlPath, []byte("<h1>Title</h1><p>Content</p>"), 0644)

	err := htmlToDOCX(ConvertOpts{InputPath: htmlPath, OutputPath: docxPath})
	if err != nil {
		t.Fatal(err)
	}

	paragraphs, err := readDOCXParagraphs(docxPath)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(paragraphs, " ")
	if !strings.Contains(joined, "Title") {
		t.Error("title not in DOCX")
	}
}

// ---------------------------------------------------------------------------
// DOCX writer/reader
// ---------------------------------------------------------------------------

func TestWriteAndReadDOCX(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.docx")

	input := []string{"First paragraph", "Second paragraph", "Third paragraph"}
	if err := writeDOCX(path, input); err != nil {
		t.Fatal(err)
	}

	output, err := readDOCXParagraphs(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(output) != len(input) {
		t.Fatalf("paragraph count: got %d, want %d", len(output), len(input))
	}
	for i, p := range output {
		if p != input[i] {
			t.Errorf("paragraph %d: got %q, want %q", i, p, input[i])
		}
	}
}

func TestWriteDOCX_SpecialChars(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "test.docx")

	input := []string{"Quotes \"here\" & <angle> brackets"}
	if err := writeDOCX(path, input); err != nil {
		t.Fatal(err)
	}

	output, err := readDOCXParagraphs(path)
	if err != nil {
		t.Fatal(err)
	}
	if output[0] != input[0] {
		t.Errorf("got %q, want %q", output[0], input[0])
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestDOCXToText_MissingFile(t *testing.T) {
	err := docxToText(ConvertOpts{InputPath: "/nonexistent.docx", OutputPath: "/tmp/out.txt"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestReadDOCXParagraphs_InvalidFile(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "bad.docx")
	os.WriteFile(path, []byte("not a zip file"), 0644)

	_, err := readDOCXParagraphs(path)
	if err == nil {
		t.Error("expected error for invalid DOCX")
	}
}

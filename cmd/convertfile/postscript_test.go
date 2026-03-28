package convertfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Text → PS
// ---------------------------------------------------------------------------

func TestTextToPS(t *testing.T) {
	dir := tempDir(t)
	txtPath := filepath.Join(dir, "test.txt")
	psPath := filepath.Join(dir, "test.ps")

	os.WriteFile(txtPath, []byte("Hello World\nSecond line\n"), 0644)

	err := textToPS(ConvertOpts{InputPath: txtPath, OutputPath: psPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(psPath)
	ps := string(data)
	if !strings.HasPrefix(ps, "%!PS-Adobe") {
		t.Error("not a valid PS file")
	}
	if !strings.Contains(ps, "(Hello World) show") {
		t.Error("text not in PS output")
	}
	if !strings.Contains(ps, "showpage") {
		t.Error("missing showpage")
	}
}

// ---------------------------------------------------------------------------
// Markdown → PS
// ---------------------------------------------------------------------------

func TestMDToPS(t *testing.T) {
	dir := tempDir(t)
	mdPath := filepath.Join(dir, "test.md")
	psPath := filepath.Join(dir, "test.ps")

	os.WriteFile(mdPath, []byte("# Title\n\nSome **text**.\n"), 0644)

	err := mdToPS(ConvertOpts{InputPath: mdPath, OutputPath: psPath})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(psPath)
	ps := string(data)
	if !strings.HasPrefix(ps, "%!PS-Adobe") {
		t.Error("not a valid PS file")
	}
	if !strings.Contains(ps, "Title") {
		t.Error("title not in PS output")
	}
}

// ---------------------------------------------------------------------------
// CSV → PS
// ---------------------------------------------------------------------------

func TestCSVToPS(t *testing.T) {
	dir := tempDir(t)
	csvPath := filepath.Join(dir, "test.csv")
	psPath := filepath.Join(dir, "test.ps")

	writeCSV(t, csvPath, ',', sampleRows)

	err := csvToPS(ConvertOpts{InputPath: csvPath, OutputPath: psPath, Delimiter: ','})
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(psPath)
	ps := string(data)
	if !strings.HasPrefix(ps, "%!PS-Adobe") {
		t.Error("not a valid PS file")
	}
	if !strings.Contains(ps, "Name") {
		t.Error("header not in PS output")
	}
	if !strings.Contains(ps, "Courier-Bold") {
		t.Error("header should use bold font")
	}
}

// ---------------------------------------------------------------------------
// PS escape
// ---------------------------------------------------------------------------

func TestPSEscape(t *testing.T) {
	tests := map[string]string{
		"hello":       "hello",
		"a(b)c":       "a\\(b\\)c",
		"back\\slash": "back\\\\slash",
	}
	for input, want := range tests {
		got := psEscape(input)
		if got != want {
			t.Errorf("psEscape(%q) = %q, want %q", input, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestTextToPS_MissingFile(t *testing.T) {
	err := textToPS(ConvertOpts{InputPath: "/nonexistent.txt", OutputPath: "/tmp/out.ps"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

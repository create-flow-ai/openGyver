package mimecmd

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestMimeCmd_Metadata(t *testing.T) {
	if mimeCmd.Use == "" {
		t.Error("mimeCmd.Use must not be empty")
	}
	if mimeCmd.Short == "" {
		t.Error("mimeCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"lookupCmd", lookupCmd.Use, lookupCmd.Short},
		{"extensionCmd", extensionCmd.Use, extensionCmd.Short},
		{"detectCmd", detectCmd.Use, detectCmd.Short},
	}
	for _, c := range cmds {
		if c.use == "" {
			t.Errorf("%s.Use must not be empty", c.name)
		}
		if c.short == "" {
			t.Errorf("%s.Short must not be empty", c.name)
		}
	}
}

// ── flag existence ─────────────────────────────────────────────────────────

func TestMimeCmd_PersistentFlags(t *testing.T) {
	f := mimeCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

// ── normalizeExtension ────────────────────────────────────────────────────

func TestNormalizeExtension_WithDot(t *testing.T) {
	got := normalizeExtension(".json")
	if got != ".json" {
		t.Errorf("expected .json, got %q", got)
	}
}

func TestNormalizeExtension_WithoutDot(t *testing.T) {
	got := normalizeExtension("json")
	if got != ".json" {
		t.Errorf("expected .json, got %q", got)
	}
}

func TestNormalizeExtension_Uppercase(t *testing.T) {
	got := normalizeExtension("PDF")
	if got != ".pdf" {
		t.Errorf("expected .pdf, got %q", got)
	}
}

func TestNormalizeExtension_WithSpaces(t *testing.T) {
	got := normalizeExtension("  .txt  ")
	if got != ".txt" {
		t.Errorf("expected .txt, got %q", got)
	}
}

// ── lookupMIME ────────────────────────────────────────────────────────────

func TestLookupMIME_JSON(t *testing.T) {
	got := lookupMIME(".json")
	if got != "application/json" {
		t.Errorf("expected application/json, got %q", got)
	}
}

func TestLookupMIME_PNG(t *testing.T) {
	got := lookupMIME(".png")
	if got != "image/png" {
		t.Errorf("expected image/png, got %q", got)
	}
}

func TestLookupMIME_PDF(t *testing.T) {
	got := lookupMIME(".pdf")
	if got != "application/pdf" {
		t.Errorf("expected application/pdf, got %q", got)
	}
}

func TestLookupMIME_MP3(t *testing.T) {
	got := lookupMIME(".mp3")
	if got != "audio/mpeg" {
		t.Errorf("expected audio/mpeg, got %q", got)
	}
}

func TestLookupMIME_Unknown(t *testing.T) {
	got := lookupMIME(".zzzzz_nonexistent")
	if got != "" {
		t.Errorf("expected empty string for unknown extension, got %q", got)
	}
}

func TestLookupMIME_HTML(t *testing.T) {
	got := lookupMIME(".html")
	if got != "text/html" {
		t.Errorf("expected text/html, got %q", got)
	}
}

// ── reverseLookup ─────────────────────────────────────────────────────────

func TestReverseLookup_JSON(t *testing.T) {
	buildReverse()
	got := reverseLookup("application/json")
	if got != ".json" {
		t.Errorf("expected .json, got %q", got)
	}
}

func TestReverseLookup_PNG(t *testing.T) {
	buildReverse()
	got := reverseLookup("image/png")
	if got != ".png" {
		t.Errorf("expected .png, got %q", got)
	}
}

func TestReverseLookup_Unknown(t *testing.T) {
	buildReverse()
	got := reverseLookup("application/x-totally-made-up-type-zzzzz")
	if got != "" {
		t.Errorf("expected empty string for unknown MIME type, got %q", got)
	}
}

// ── mimeTypes map coverage ────────────────────────────────────────────────

func TestMimeTypesMap_HasMinimumEntries(t *testing.T) {
	// The map should have at least 90 entries (we target ~100).
	if len(mimeTypes) < 90 {
		t.Errorf("expected at least 90 MIME types, got %d", len(mimeTypes))
	}
}

func TestMimeTypesMap_AllExtensionsHaveDot(t *testing.T) {
	for ext := range mimeTypes {
		if !strings.HasPrefix(ext, ".") {
			t.Errorf("extension %q missing dot prefix", ext)
		}
	}
}

// ── detectFile ────────────────────────────────────────────────────────────

func TestDetectFile_PNG(t *testing.T) {
	// Create a minimal PNG file (8-byte PNG signature).
	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")
	pngSig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if err := os.WriteFile(path, pngSig, 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}

	got, err := detectFile(path)
	if err != nil {
		t.Fatalf("detectFile error: %v", err)
	}
	if got != "image/png" {
		t.Errorf("expected image/png, got %q", got)
	}
}

func TestDetectFile_PlainText(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("Hello, World!"), 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}

	got, err := detectFile(path)
	if err != nil {
		t.Fatalf("detectFile error: %v", err)
	}
	// net/http.DetectContentType returns "text/plain; charset=utf-8" for plain text.
	if !strings.HasPrefix(got, "text/plain") {
		t.Errorf("expected text/plain prefix, got %q", got)
	}
}

func TestDetectFile_HTML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.html")
	if err := os.WriteFile(path, []byte("<html><body>hello</body></html>"), 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}

	got, err := detectFile(path)
	if err != nil {
		t.Fatalf("detectFile error: %v", err)
	}
	if !strings.HasPrefix(got, "text/html") {
		t.Errorf("expected text/html prefix, got %q", got)
	}
}

func TestDetectFile_NonexistentFile(t *testing.T) {
	_, err := detectFile("/nonexistent/path/to/file.bin")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestDetectFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty")
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}

	got, err := detectFile(path)
	if err != nil {
		t.Fatalf("detectFile error: %v", err)
	}
	// http.DetectContentType on empty returns "text/plain; charset=utf-8"
	expected := http.DetectContentType([]byte{})
	if got != expected {
		t.Errorf("expected %q for empty file, got %q", expected, got)
	}
}

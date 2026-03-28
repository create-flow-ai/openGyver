package pdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestPdfCmd_Metadata(t *testing.T) {
	if pdfCmd.Use != "pdf" {
		t.Errorf("unexpected Use: %s", pdfCmd.Use)
	}
	if pdfCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestPdfCmd_Subcommands(t *testing.T) {
	names := map[string]bool{}
	for _, sub := range pdfCmd.Commands() {
		names[sub.Use[:len(sub.Name())] ] = true
	}
	for _, want := range []string{"merge", "split", "pages", "info"} {
		if !names[want] {
			t.Errorf("missing subcommand: %s", want)
		}
	}
}

func TestMergeCmd_RequiresOutput(t *testing.T) {
	old := mergeOutput
	defer func() { mergeOutput = old }()
	mergeOutput = ""
	err := runMerge(mergeCmd, []string{"a.pdf", "b.pdf"})
	if err == nil {
		t.Error("expected error when --output is empty")
	}
}

func TestMergeCmd_MissingInputFile(t *testing.T) {
	old := mergeOutput
	defer func() { mergeOutput = old }()
	mergeOutput = "/tmp/out.pdf"
	err := runMerge(mergeCmd, []string{"/nonexistent/a.pdf", "/nonexistent/b.pdf"})
	if err == nil {
		t.Error("expected error for missing input files")
	}
}

// ── split ──────────────────────────────────────────────────────────────────

func TestSplitCmd_MissingInputFile(t *testing.T) {
	err := runSplit(splitCmd, []string{"/nonexistent/doc.pdf"})
	if err == nil {
		t.Error("expected error for missing input file")
	}
}

func TestSplitCmd_AcceptsExactlyOneArg(t *testing.T) {
	validator := cobra.ExactArgs(1)
	if err := validator(splitCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(splitCmd, []string{"one"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// ── pages ──────────────────────────────────────────────────────────────────

func TestCountPages_MissingFile(t *testing.T) {
	_, err := countPages("/nonexistent/missing.pdf")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCountPagesFallback_NoPagesInContent(t *testing.T) {
	// Create a temp file without PDF page markers.
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.pdf")
	os.WriteFile(path, []byte("not a real pdf"), 0644)

	n, err := countPagesFallback(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 pages, got %d", n)
	}
}

func TestCountPagesFallback_WithPageMarkers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.pdf")
	content := `/Type /Page\n/Type /Page\n/Type /Pages\n/Type /Page `
	os.WriteFile(path, []byte(content), 0644)

	n, err := countPagesFallback(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should match /Type /Page followed by non-s characters.
	// Content has 3 "/Type /Page" (not followed by "s") occurrences.
	if n < 2 {
		t.Errorf("expected at least 2 page markers, got %d", n)
	}
}

// ── info ───────────────────────────────────────────────────────────────────

func TestInfoCmd_MissingFile(t *testing.T) {
	err := runInfo(infoCmd, []string{"/nonexistent/doc.pdf"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// ── cleanDateString ────────────────────────────────────────────────────────

func TestCleanDateString(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"D:20210101120000", "20210101120000"},
		{"20210101", "20210101"},
		{"D: 20210101", "20210101"},
		{"", ""},
	}
	for _, tc := range tests {
		got := cleanDateString(tc.input)
		if got != tc.want {
			t.Errorf("cleanDateString(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// ── flags ──────────────────────────────────────────────────────────────────

func TestMergeCmd_HasOutputFlag(t *testing.T) {
	f := mergeCmd.Flags().Lookup("output")
	if f == nil {
		t.Error("merge should have --output flag")
	}
	if f.Shorthand != "o" {
		t.Errorf("merge --output shorthand: got %q, want 'o'", f.Shorthand)
	}
}

func TestSplitCmd_HasOutputDirFlag(t *testing.T) {
	f := splitCmd.Flags().Lookup("output-dir")
	if f == nil {
		t.Error("split should have --output-dir flag")
	}
}

func TestPagesCmd_HasJSONFlag(t *testing.T) {
	f := pagesCmd.Flags().Lookup("json")
	if f == nil {
		t.Error("pages should have --json flag")
	}
	if f.Shorthand != "j" {
		t.Errorf("--json shorthand: got %q, want 'j'", f.Shorthand)
	}
}

package pdf

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var jsonOut bool

// ── parent command ─────────────────────────────────────────────────────────

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "PDF tools — merge, split, page count, metadata",
	Long: `PDF tools — merge, split, count pages, and inspect metadata.

SUBCOMMANDS:

  merge   Merge multiple PDF files into one
  split   Split a PDF into individual pages
  pages   Count pages in a PDF
  info    Show PDF metadata (title, author, dates, pages)

All subcommands support --json/-j for machine-readable output.

EXAMPLES:

  openGyver pdf merge -o combined.pdf file1.pdf file2.pdf file3.pdf
  openGyver pdf split document.pdf -o ./pages/
  openGyver pdf pages document.pdf
  openGyver pdf info document.pdf`,
}

// ── merge ──────────────────────────────────────────────────────────────────

var mergeOutput string

var mergeCmd = &cobra.Command{
	Use:   "merge <file1> <file2> [file3...]",
	Short: "Merge multiple PDFs into one",
	Long: `Merge two or more PDF files into a single output file.

Examples:
  openGyver pdf merge -o combined.pdf a.pdf b.pdf
  openGyver pdf merge -o all.pdf *.pdf`,
	Args: cobra.MinimumNArgs(2),
	RunE: runMerge,
}

func runMerge(_ *cobra.Command, args []string) error {
	if mergeOutput == "" {
		return fmt.Errorf("--output (-o) is required")
	}

	// Verify all input files exist.
	for _, f := range args {
		if _, err := os.Stat(f); err != nil {
			return fmt.Errorf("input file not found: %s", f)
		}
	}

	if err := pdfcpuapi.MergeCreateFile(args, mergeOutput, false, nil); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"action": "merge",
			"inputs": args,
			"output": mergeOutput,
		})
	}
	fmt.Printf("Merged %d files → %s\n", len(args), mergeOutput)
	return nil
}

// ── split ──────────────────────────────────────────────────────────────────

var splitOutputDir string

var splitCmd = &cobra.Command{
	Use:   "split <file>",
	Short: "Split a PDF into individual pages",
	Long: `Split a PDF into individual single-page PDF files.

Each page is written as a separate file in the output directory.

Examples:
  openGyver pdf split document.pdf
  openGyver pdf split document.pdf -o ./pages/`,
	Args: cobra.ExactArgs(1),
	RunE: runSplit,
}

func runSplit(_ *cobra.Command, args []string) error {
	input := args[0]
	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input file not found: %s", input)
	}

	outDir := splitOutputDir
	if outDir == "" {
		outDir = "."
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Split with span=1 means one page per file.
	if err := pdfcpuapi.SplitFile(input, outDir, 1, nil); err != nil {
		return fmt.Errorf("split failed: %w", err)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"action":    "split",
			"input":     input,
			"outputDir": outDir,
		})
	}
	fmt.Printf("Split %s → %s/\n", input, outDir)
	return nil
}

// ── pages ──────────────────────────────────────────────────────────────────

var pagesCmd = &cobra.Command{
	Use:   "pages <file>",
	Short: "Count pages in a PDF",
	Long: `Count the number of pages in a PDF file.

Examples:
  openGyver pdf pages document.pdf
  openGyver pdf pages document.pdf --json`,
	Args: cobra.ExactArgs(1),
	RunE: runPages,
}

func runPages(_ *cobra.Command, args []string) error {
	input := args[0]
	n, err := countPages(input)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"file":  input,
			"pages": n,
		})
	}
	fmt.Printf("%s: %d pages\n", input, n)
	return nil
}

// countPages returns the page count for a PDF file.
// It tries pdfcpu first; if that fails it falls back to a simple
// regex count of "/Type /Page" (excluding "/Type /Pages") occurrences.
func countPages(path string) (int, error) {
	if _, err := os.Stat(path); err != nil {
		return 0, fmt.Errorf("file not found: %s", path)
	}

	n, err := countPagesPdfcpu(path)
	if err == nil {
		return n, nil
	}
	return countPagesFallback(path)
}

func countPagesPdfcpu(path string) (int, error) {
	ctx, err := pdfcpuapi.ReadContextFile(path)
	if err != nil {
		return 0, err
	}
	return ctx.PageCount, nil
}

// countPagesFallback reads raw bytes and counts /Type /Page tokens.
func countPagesFallback(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("reading file: %w", err)
	}
	content := string(data)

	// Match "/Type /Page" but NOT "/Type /Pages"
	re := regexp.MustCompile(`/Type\s*/Page[^s]`)
	matches := re.FindAllString(content, -1)
	if len(matches) == 0 {
		// Also try without the negative look-ahead in case file ends at /Page
		re2 := regexp.MustCompile(`/Type\s*/Page\b`)
		matches = re2.FindAllString(content, -1)
	}
	return len(matches), nil
}

// ── info ───────────────────────────────────────────────────────────────────

var infoCmd = &cobra.Command{
	Use:   "info <file>",
	Short: "Show PDF metadata",
	Long: `Show PDF metadata: title, author, creator, creation date, page count.

Examples:
  openGyver pdf info document.pdf
  openGyver pdf info document.pdf --json`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func runInfo(_ *cobra.Command, args []string) error {
	input := args[0]
	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("file not found: %s", input)
	}

	info, err := extractInfo(input)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(info)
	}

	printField := func(label, val string) {
		if val != "" {
			fmt.Printf("%-16s %s\n", label+":", val)
		}
	}

	printField("File", input)
	printField("Title", info["title"].(string))
	printField("Author", info["author"].(string))
	printField("Creator", info["creator"].(string))
	printField("Producer", info["producer"].(string))
	printField("Creation Date", info["creationDate"].(string))
	printField("Mod Date", info["modDate"].(string))
	fmt.Printf("%-16s %d\n", "Pages:", info["pages"])
	return nil
}

func extractInfo(path string) (map[string]interface{}, error) {
	ctx, err := pdfcpuapi.ReadContextFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading PDF: %w", err)
	}

	info := map[string]interface{}{
		"file":         path,
		"title":        "",
		"author":       "",
		"creator":      "",
		"producer":     "",
		"creationDate": "",
		"modDate":      "",
		"pages":        ctx.PageCount,
	}

	if ctx.XRefTable != nil {
		fillFromXRefTable(ctx.XRefTable, info)
	}

	return info, nil
}

func fillFromXRefTable(xrt *model.XRefTable, info map[string]interface{}) {
	if xrt.Title != "" {
		info["title"] = xrt.Title
	}
	if xrt.Author != "" {
		info["author"] = xrt.Author
	}
	if xrt.Creator != "" {
		info["creator"] = xrt.Creator
	}
	if xrt.Producer != "" {
		info["producer"] = xrt.Producer
	}
	if xrt.CreationDate != "" {
		info["creationDate"] = cleanDateString(xrt.CreationDate)
	}
	if xrt.ModDate != "" {
		info["modDate"] = cleanDateString(xrt.ModDate)
	}
}

// cleanDateString removes the PDF date prefix D: if present.
func cleanDateString(s string) string {
	s = strings.TrimPrefix(s, "D:")
	return strings.TrimSpace(s)
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "output PDF path (required)")
	splitCmd.Flags().StringVarP(&splitOutputDir, "output-dir", "o", ".", "output directory for split pages")

	// Add --json/-j to all leaf commands.
	for _, c := range []*cobra.Command{mergeCmd, splitCmd, pagesCmd, infoCmd} {
		c.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	}

	pdfCmd.AddCommand(mergeCmd, splitCmd, pagesCmd, infoCmd)
	cmd.Register(pdfCmd)
}

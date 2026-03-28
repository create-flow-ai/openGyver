package convertfile

import (
	"fmt"
	"os"
	"strings"
)

const (
	psPageWidth  = 612 // US Letter
	psPageHeight = 792
	psMargin     = 50
	psFontSize   = 10
	psLineHeight = 14
)

// textToPS renders plain text as a PostScript document.
func textToPS(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading text: %w", err)
	}
	return writeTextPS(opts.OutputPath, string(src))
}

// mdToPS converts Markdown to PostScript via text extraction.
func mdToPS(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading markdown: %w", err)
	}
	text := stripMarkdown(string(src))
	return writeTextPS(opts.OutputPath, text)
}

// csvToPS renders CSV data as a tabular PostScript document.
func csvToPS(opts ConvertOpts) error {
	records, err := readCSV(opts.InputPath, opts.Delimiter)
	if err != nil {
		return fmt.Errorf("reading CSV: %w", err)
	}
	return writeTablePS(opts.OutputPath, records)
}

// --- PostScript rendering helpers ---

func writeTextPS(path, text string) error {
	var ps strings.Builder

	ps.WriteString(psHeader())
	ps.WriteString(psBeginPage())

	y := psPageHeight - psMargin
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if y < psMargin+psLineHeight {
			ps.WriteString("showpage\n")
			ps.WriteString(psBeginPage())
			y = psPageHeight - psMargin
		}

		ps.WriteString(fmt.Sprintf("%d %d moveto\n", psMargin, y))
		ps.WriteString(fmt.Sprintf("(%s) show\n", psEscape(line)))
		y -= psLineHeight
	}

	ps.WriteString("showpage\n")
	ps.WriteString("%%EOF\n")

	return os.WriteFile(path, []byte(ps.String()), 0644)
}

func writeTablePS(path string, records [][]string) error {
	if len(records) == 0 {
		return writeTextPS(path, "(empty)")
	}

	maxCols := 0
	for _, row := range records {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	colWidth := (psPageWidth - 2*psMargin) / maxCols

	var ps strings.Builder
	ps.WriteString(psHeader())
	ps.WriteString(psBeginPage())

	y := psPageHeight - psMargin

	for i, row := range records {
		if y < psMargin+psLineHeight {
			ps.WriteString("showpage\n")
			ps.WriteString(psBeginPage())
			y = psPageHeight - psMargin
		}

		// Bold for header
		if i == 0 {
			ps.WriteString("/Courier-Bold findfont 10 scalefont setfont\n")
		} else if i == 1 {
			ps.WriteString("/Courier findfont 10 scalefont setfont\n")
		}

		x := psMargin
		for j := 0; j < maxCols; j++ {
			cell := ""
			if j < len(row) {
				cell = row[j]
			}
			// Truncate to fit column
			maxChars := colWidth / 6
			if len(cell) > int(maxChars) {
				cell = cell[:int(maxChars)-1] + "~"
			}
			ps.WriteString(fmt.Sprintf("%d %d moveto\n", x, y))
			ps.WriteString(fmt.Sprintf("(%s) show\n", psEscape(cell)))
			x += colWidth
		}

		// Draw horizontal line under header
		if i == 0 {
			y -= 2
			ps.WriteString(fmt.Sprintf("%d %d moveto %d %d lineto stroke\n",
				psMargin, y, psPageWidth-psMargin, y))
		}

		y -= psLineHeight
	}

	ps.WriteString("showpage\n")
	ps.WriteString("%%EOF\n")

	return os.WriteFile(path, []byte(ps.String()), 0644)
}

func psHeader() string {
	return `%!PS-Adobe-3.0
%%Creator: openGyver convertFile
%%Pages: (atend)
%%EndComments
`
}

func psBeginPage() string {
	return fmt.Sprintf("/Courier findfont %d scalefont setfont\n", psFontSize)
}

func psEscape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}

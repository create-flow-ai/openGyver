package qr

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

var (
	output string
	size   int
	level  string
	invert bool
)

var qrCmd = &cobra.Command{
	Use:   "qr <text>",
	Short: "Generate a QR code from text",
	Long: `Generate a QR code from a string and display it in the terminal as ASCII art,
or save it to a file as PNG or SVG.

OUTPUT MODES:

  Default (no -o)    Prints QR code as ASCII block characters in the terminal
  -o file.png        Saves as PNG image (use --size to set pixel dimensions)
  -o file.svg        Saves as SVG vector graphic

ERROR CORRECTION LEVELS:

  L    7%  recovery (default — smallest QR code)
  M    15% recovery
  Q    25% recovery
  H    30% recovery (largest QR code, most resilient)

Examples:
  openGyver qr "https://example.com"
  openGyver qr "Hello World" -o qr.png
  openGyver qr "Hello World" -o qr.png --size 512
  openGyver qr "Hello World" -o qr.svg
  openGyver qr "wifi:WPA;S:MyNetwork;P:secret;;" --level H
  openGyver qr "some data" --invert`,
	Args: cobra.ExactArgs(1),
	RunE: runQR,
}

func runQR(c *cobra.Command, args []string) error {
	content := args[0]
	ecLevel := parseLevel(level)

	if output == "" {
		return printASCII(content, ecLevel, invert)
	}

	ext := strings.ToLower(filepath.Ext(output))
	switch ext {
	case ".png":
		return savePNG(content, output, size, ecLevel)
	case ".svg":
		return saveSVG(content, output, ecLevel)
	default:
		return fmt.Errorf("unsupported output format: %s (use .png or .svg)", ext)
	}
}

func parseLevel(s string) qrcode.RecoveryLevel {
	switch strings.ToUpper(s) {
	case "M":
		return qrcode.Medium
	case "Q":
		return qrcode.High
	case "H":
		return qrcode.Highest
	default:
		return qrcode.Low
	}
}

// printASCII renders the QR code as Unicode block characters in the terminal.
// Uses █ for black and spaces for white (or inverted with --invert).
func printASCII(content string, level qrcode.RecoveryLevel, inverted bool) error {
	qr, err := qrcode.New(content, level)
	if err != nil {
		return fmt.Errorf("generating QR code: %w", err)
	}

	bitmap := qr.Bitmap()
	rows := len(bitmap)

	black := "██"
	white := "  "
	if inverted {
		black, white = white, black
	}

	// Process two rows at a time using half-block characters for compact output
	// ▀ = top half, ▄ = bottom half, █ = full block, " " = empty
	topBlk := "▀"
	botBlk := "▄"
	fullBlk := "█"
	emptyBlk := " "
	if inverted {
		topBlk = "▄"
		botBlk = "▀"
		fullBlk = " "
		emptyBlk = "█"
	}
	_ = black
	_ = white

	for y := 0; y < rows; y += 2 {
		for x := 0; x < len(bitmap[y]); x++ {
			top := bitmap[y][x]
			bot := false
			if y+1 < rows {
				bot = bitmap[y+1][x]
			}

			switch {
			case top && bot:
				fmt.Print(fullBlk)
			case top && !bot:
				fmt.Print(topBlk)
			case !top && bot:
				fmt.Print(botBlk)
			default:
				fmt.Print(emptyBlk)
			}
		}
		fmt.Println()
	}

	return nil
}

// savePNG saves the QR code as a PNG file.
func savePNG(content, path string, size int, level qrcode.RecoveryLevel) error {
	qr, err := qrcode.New(content, level)
	if err != nil {
		return fmt.Errorf("generating QR code: %w", err)
	}

	img := qr.Image(size)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("encoding PNG: %w", err)
	}

	fmt.Printf("Saved QR code to %s (%dx%d)\n", path, size, size)
	return nil
}

// saveSVG saves the QR code as an SVG file.
func saveSVG(content, path string, level qrcode.RecoveryLevel) error {
	qr, err := qrcode.New(content, level)
	if err != nil {
		return fmt.Errorf("generating QR code: %w", err)
	}

	bitmap := qr.Bitmap()
	rows := len(bitmap)
	cols := 0
	if rows > 0 {
		cols = len(bitmap[0])
	}

	moduleSize := 10
	svgWidth := cols * moduleSize
	svgHeight := rows * moduleSize

	var svg strings.Builder
	svg.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
<rect width="100%%" height="100%%" fill="white"/>
`, svgWidth, svgHeight, svgWidth, svgHeight))

	for y, row := range bitmap {
		for x, black := range row {
			if black {
				svg.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="black"/>
`, x*moduleSize, y*moduleSize, moduleSize, moduleSize))
			}
		}
	}

	svg.WriteString("</svg>\n")

	if err := os.WriteFile(path, []byte(svg.String()), 0644); err != nil {
		return fmt.Errorf("writing SVG: %w", err)
	}

	fmt.Printf("Saved QR code to %s (%dx%d)\n", path, svgWidth, svgHeight)
	return nil
}

func init() {
	qrCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (.png or .svg). Omit for ASCII terminal output")
	qrCmd.Flags().IntVar(&size, "size", 256, "PNG image size in pixels")
	qrCmd.Flags().StringVar(&level, "level", "L", "error correction level: L, M, Q, H")
	qrCmd.Flags().BoolVar(&invert, "invert", false, "invert colors (light-on-dark for dark terminals)")
	cmd.Register(qrCmd)
}

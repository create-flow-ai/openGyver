package ascii

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var jsonOut bool

// ── parent command ─────────────────────────────────────────────────────────

var asciiCmd = &cobra.Command{
	Use:   "ascii",
	Short: "ASCII tools — banners, table, character lookup",
	Long: `ASCII tools — generate text banners, print the ASCII table, and look up characters.

SUBCOMMANDS:

  banner   Generate a large text banner from input text
  table    Print the printable ASCII table (chars 32-126)
  lookup   Look up a character by decimal value or literal

All subcommands support --json/-j for machine-readable output.

EXAMPLES:

  openGyver ascii banner "HELLO"
  openGyver ascii table
  openGyver ascii lookup 65
  openGyver ascii lookup A`,
}

// ── block-letter font ──────────────────────────────────────────────────────
// Each letter is 5 lines high and 5 columns wide (plus 1-col spacing).

var blockFont = map[rune][5]string{
	'A': {
		" ### ",
		"#   #",
		"#####",
		"#   #",
		"#   #",
	},
	'B': {
		"#### ",
		"#   #",
		"#### ",
		"#   #",
		"#### ",
	},
	'C': {
		" ####",
		"#    ",
		"#    ",
		"#    ",
		" ####",
	},
	'D': {
		"#### ",
		"#   #",
		"#   #",
		"#   #",
		"#### ",
	},
	'E': {
		"#####",
		"#    ",
		"###  ",
		"#    ",
		"#####",
	},
	'F': {
		"#####",
		"#    ",
		"###  ",
		"#    ",
		"#    ",
	},
	'G': {
		" ####",
		"#    ",
		"# ###",
		"#   #",
		" ### ",
	},
	'H': {
		"#   #",
		"#   #",
		"#####",
		"#   #",
		"#   #",
	},
	'I': {
		"#####",
		"  #  ",
		"  #  ",
		"  #  ",
		"#####",
	},
	'J': {
		"#####",
		"    #",
		"    #",
		"#   #",
		" ### ",
	},
	'K': {
		"#   #",
		"#  # ",
		"###  ",
		"#  # ",
		"#   #",
	},
	'L': {
		"#    ",
		"#    ",
		"#    ",
		"#    ",
		"#####",
	},
	'M': {
		"#   #",
		"## ##",
		"# # #",
		"#   #",
		"#   #",
	},
	'N': {
		"#   #",
		"##  #",
		"# # #",
		"#  ##",
		"#   #",
	},
	'O': {
		" ### ",
		"#   #",
		"#   #",
		"#   #",
		" ### ",
	},
	'P': {
		"#### ",
		"#   #",
		"#### ",
		"#    ",
		"#    ",
	},
	'Q': {
		" ### ",
		"#   #",
		"# # #",
		"#  # ",
		" ## #",
	},
	'R': {
		"#### ",
		"#   #",
		"#### ",
		"#  # ",
		"#   #",
	},
	'S': {
		" ####",
		"#    ",
		" ### ",
		"    #",
		"#### ",
	},
	'T': {
		"#####",
		"  #  ",
		"  #  ",
		"  #  ",
		"  #  ",
	},
	'U': {
		"#   #",
		"#   #",
		"#   #",
		"#   #",
		" ### ",
	},
	'V': {
		"#   #",
		"#   #",
		"#   #",
		" # # ",
		"  #  ",
	},
	'W': {
		"#   #",
		"#   #",
		"# # #",
		"## ##",
		"#   #",
	},
	'X': {
		"#   #",
		" # # ",
		"  #  ",
		" # # ",
		"#   #",
	},
	'Y': {
		"#   #",
		" # # ",
		"  #  ",
		"  #  ",
		"  #  ",
	},
	'Z': {
		"#####",
		"   # ",
		"  #  ",
		" #   ",
		"#####",
	},
	'0': {
		" ### ",
		"#  ##",
		"# # #",
		"##  #",
		" ### ",
	},
	'1': {
		" #   ",
		"##   ",
		" #   ",
		" #   ",
		"#####",
	},
	'2': {
		" ### ",
		"#   #",
		"  ## ",
		" #   ",
		"#####",
	},
	'3': {
		" ### ",
		"#   #",
		"  ## ",
		"#   #",
		" ### ",
	},
	'4': {
		"#   #",
		"#   #",
		"#####",
		"    #",
		"    #",
	},
	'5': {
		"#####",
		"#    ",
		"#### ",
		"    #",
		"#### ",
	},
	'6': {
		" ### ",
		"#    ",
		"#### ",
		"#   #",
		" ### ",
	},
	'7': {
		"#####",
		"    #",
		"   # ",
		"  #  ",
		"  #  ",
	},
	'8': {
		" ### ",
		"#   #",
		" ### ",
		"#   #",
		" ### ",
	},
	'9': {
		" ### ",
		"#   #",
		" ####",
		"    #",
		" ### ",
	},
	' ': {
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
	},
	'!': {
		"  #  ",
		"  #  ",
		"  #  ",
		"     ",
		"  #  ",
	},
	'.': {
		"     ",
		"     ",
		"     ",
		"     ",
		"  #  ",
	},
	'-': {
		"     ",
		"     ",
		"#####",
		"     ",
		"     ",
	},
	'?': {
		" ### ",
		"#   #",
		"  ## ",
		"     ",
		"  #  ",
	},
}

// renderBanner builds a 5-line banner for the given text.
func renderBanner(text string) string {
	upper := strings.ToUpper(text)
	lines := [5]strings.Builder{}

	for i, ch := range upper {
		glyph, ok := blockFont[ch]
		if !ok {
			glyph = blockFont['?']
		}
		for row := 0; row < 5; row++ {
			if i > 0 {
				lines[row].WriteByte(' ')
			}
			lines[row].WriteString(glyph[row])
		}
	}

	var sb strings.Builder
	for i := 0; i < 5; i++ {
		sb.WriteString(lines[i].String())
		if i < 4 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// ── banner subcommand ──────────────────────────────────────────────────────

var bannerCmd = &cobra.Command{
	Use:   "banner <text>",
	Short: "Generate a large text banner",
	Long: `Generate a large block-letter ASCII art banner from the given text.

Supports A-Z, 0-9, space, and a few punctuation marks.

Examples:
  openGyver ascii banner "HELLO"
  openGyver ascii banner "GO 1.25"`,
	Args: cobra.ExactArgs(1),
	RunE: runBanner,
}

func runBanner(_ *cobra.Command, args []string) error {
	text := args[0]
	banner := renderBanner(text)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"text":   text,
			"banner": banner,
		})
	}
	fmt.Println(banner)
	return nil
}

// ── table subcommand ───────────────────────────────────────────────────────

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Print the ASCII table (characters 32-126)",
	Long: `Print a table of printable ASCII characters (32-126) with their
decimal, hexadecimal, octal, and character representations.

Examples:
  openGyver ascii table
  openGyver ascii table --json`,
	Args: cobra.NoArgs,
	RunE: runTable,
}

type asciiEntry struct {
	Decimal int    `json:"decimal"`
	Hex     string `json:"hex"`
	Octal   string `json:"octal"`
	Char    string `json:"char"`
}

func buildTable() []asciiEntry {
	entries := make([]asciiEntry, 0, 95)
	for i := 32; i <= 126; i++ {
		ch := string(rune(i))
		if i == 32 {
			ch = "SP"
		}
		entries = append(entries, asciiEntry{
			Decimal: i,
			Hex:     fmt.Sprintf("0x%02X", i),
			Octal:   fmt.Sprintf("0%03o", i),
			Char:    ch,
		})
	}
	return entries
}

func runTable(_ *cobra.Command, _ []string) error {
	entries := buildTable()

	if jsonOut {
		return cmd.PrintJSON(entries)
	}

	fmt.Printf("%-8s %-8s %-8s %s\n", "DEC", "HEX", "OCT", "CHAR")
	fmt.Println(strings.Repeat("-", 32))
	for _, e := range entries {
		fmt.Printf("%-8d %-8s %-8s %s\n", e.Decimal, e.Hex, e.Octal, e.Char)
	}
	return nil
}

// ── lookup subcommand ──────────────────────────────────────────────────────

var lookupCmd = &cobra.Command{
	Use:   "lookup <decimal | character>",
	Short: "Look up a character's ASCII codes",
	Long: `Look up a character by its decimal value or by the literal character.

Shows: character, decimal, hex, octal, binary, HTML entity, URL encoding.

Examples:
  openGyver ascii lookup 65
  openGyver ascii lookup A
  openGyver ascii lookup 32`,
	Args: cobra.ExactArgs(1),
	RunE: runLookup,
}

type lookupResult struct {
	Char        string `json:"char"`
	Decimal     int    `json:"decimal"`
	Hex         string `json:"hex"`
	Octal       string `json:"octal"`
	Binary      string `json:"binary"`
	HTMLEntity  string `json:"htmlEntity"`
	URLEncoding string `json:"urlEncoding"`
}

func charLookup(input string) (*lookupResult, error) {
	var code int

	// Try parsing as a decimal number first.
	if n, err := strconv.Atoi(input); err == nil {
		code = n
	} else if len([]rune(input)) == 1 {
		code = int([]rune(input)[0])
	} else {
		return nil, fmt.Errorf("invalid input %q: provide a decimal number or a single character", input)
	}

	if code < 0 || code > 127 {
		return nil, fmt.Errorf("value %d is outside the ASCII range (0-127)", code)
	}

	ch := string(rune(code))
	display := ch
	if code < 32 || code == 32 || code == 127 {
		// Non-printable or space — use a label.
		display = controlName(code)
	}

	return &lookupResult{
		Char:        display,
		Decimal:     code,
		Hex:         fmt.Sprintf("0x%02X", code),
		Octal:       fmt.Sprintf("0%03o", code),
		Binary:      fmt.Sprintf("%08b", code),
		HTMLEntity:  fmt.Sprintf("&#%d;", code),
		URLEncoding: fmt.Sprintf("%%%02X", code),
	}, nil
}

func controlName(code int) string {
	names := map[int]string{
		0: "NUL", 1: "SOH", 2: "STX", 3: "ETX", 4: "EOT", 5: "ENQ",
		6: "ACK", 7: "BEL", 8: "BS", 9: "HT", 10: "LF", 11: "VT",
		12: "FF", 13: "CR", 14: "SO", 15: "SI", 16: "DLE", 17: "DC1",
		18: "DC2", 19: "DC3", 20: "DC4", 21: "NAK", 22: "SYN", 23: "ETB",
		24: "CAN", 25: "EM", 26: "SUB", 27: "ESC", 28: "FS", 29: "GS",
		30: "RS", 31: "US", 32: "SP", 127: "DEL",
	}
	if n, ok := names[code]; ok {
		return n
	}
	return fmt.Sprintf("CTRL-%d", code)
}

func runLookup(_ *cobra.Command, args []string) error {
	result, err := charLookup(args[0])
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(result)
	}

	fmt.Printf("Character:    %s\n", result.Char)
	fmt.Printf("Decimal:      %d\n", result.Decimal)
	fmt.Printf("Hex:          %s\n", result.Hex)
	fmt.Printf("Octal:        %s\n", result.Octal)
	fmt.Printf("Binary:       %s\n", result.Binary)
	fmt.Printf("HTML Entity:  %s\n", result.HTMLEntity)
	fmt.Printf("URL Encoding: %s\n", result.URLEncoding)
	return nil
}

// ── helpers ────────────────────────────────────────────────────────────────

// isASCIIPrintable returns true for characters 32-126.
func isASCIIPrintable(r rune) bool {
	return r >= 32 && r <= 126 && unicode.IsPrint(r)
}

func init() {
	for _, c := range []*cobra.Command{bannerCmd, tableCmd, lookupCmd} {
		c.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	}

	asciiCmd.AddCommand(bannerCmd, tableCmd, lookupCmd)
	cmd.Register(asciiCmd)
}

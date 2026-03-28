package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// jsonOut controls JSON output across all subcommands.
var jsonOut bool

// ─── Parent command ─────────────────────────────────────────────────────────

var numberCmd = &cobra.Command{
	Use:   "number",
	Short: "Number conversion utilities",
	Long: `Convert numbers between bases, Roman numerals, and IEEE 754 representations.

SUBCOMMANDS:

  base      Convert a number between bases (2-36)
  roman     Convert between Roman numerals and decimal
  ieee754   Show IEEE 754 floating-point bit representation

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver number base 255 --to 16
  openGyver number base ff --from 16 --to 2
  openGyver number roman 42
  openGyver number roman XLII
  openGyver number ieee754 3.14`,
}

func init() {
	numberCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	numberCmd.AddCommand(baseCmd)
	numberCmd.AddCommand(romanCmd)
	numberCmd.AddCommand(ieee754Cmd)
	cmd.Register(numberCmd)
}

// ─── base subcommand ────────────────────────────────────────────────────────

var (
	fromBase int
	toBase   int
)

var baseCmd = &cobra.Command{
	Use:   "base <value>",
	Short: "Convert a number between bases (2-36)",
	Long: `Convert a number from one base to another.

Supports bases 2 through 36. The input value is interpreted in the base
given by --from (default 10), and printed in the base given by --to.

Letters a-z (case-insensitive) represent digits 10-35 for bases > 10.

FLAGS:

  --from   Source base (default 10)
  --to     Target base (required)

EXAMPLES:

  # Decimal to hexadecimal
  openGyver number base 255 --to 16
  ff

  # Hexadecimal to binary
  openGyver number base ff --from 16 --to 2
  11111111

  # Binary to decimal
  openGyver number base 11111111 --from 2 --to 10
  255

  # Decimal to base-36
  openGyver number base 1000 --to 36
  rs

  # Octal to hexadecimal
  openGyver number base 377 --from 8 --to 16
  ff`,
	Args: cobra.ExactArgs(1),
	RunE: runBase,
}

func init() {
	baseCmd.Flags().IntVar(&fromBase, "from", 10, "source base (2-36)")
	baseCmd.Flags().IntVar(&toBase, "to", 0, "target base (2-36, required)")
	_ = baseCmd.MarkFlagRequired("to")
}

func runBase(c *cobra.Command, args []string) error {
	value := args[0]

	if fromBase < 2 || fromBase > 36 {
		return fmt.Errorf("--from base must be between 2 and 36, got %d", fromBase)
	}
	if toBase < 2 || toBase > 36 {
		return fmt.Errorf("--to base must be between 2 and 36, got %d", toBase)
	}

	// Parse the input value in the source base.
	n, err := strconv.ParseInt(strings.ToLower(value), fromBase, 64)
	if err != nil {
		return fmt.Errorf("invalid value %q for base %d: %w", value, fromBase, err)
	}

	// Format in the target base.
	output := strconv.FormatInt(n, toBase)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":     value,
			"from_base": fromBase,
			"to_base":   toBase,
			"output":    output,
		})
	}

	fmt.Println(output)
	return nil
}

// ─── roman subcommand ───────────────────────────────────────────────────────

var romanCmd = &cobra.Command{
	Use:   "roman <value>",
	Short: "Convert between Roman numerals and decimal",
	Long: `Convert a decimal integer to Roman numerals or a Roman numeral string
to its decimal value. The direction is auto-detected from the input.

If the input is purely numeric, it is treated as a decimal integer and
converted to Roman numeral notation. Otherwise the input is parsed as
a Roman numeral string and converted to decimal.

Supports values 1 through 3999 (standard Roman numeral range).

EXAMPLES:

  # Decimal to Roman
  openGyver number roman 42
  XLII

  # Decimal to Roman
  openGyver number roman 1994
  MCMXCIV

  # Roman to decimal
  openGyver number roman XLII
  42

  # Roman to decimal (case-insensitive)
  openGyver number roman mcmxciv
  1994`,
	Args: cobra.ExactArgs(1),
	RunE: runRoman,
}

func runRoman(_ *cobra.Command, args []string) error {
	input := strings.TrimSpace(args[0])

	// Auto-detect direction: if the input is a decimal integer, convert to Roman.
	if isDecimal(input) {
		n, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("invalid decimal number: %q", input)
		}
		roman, err := intToRoman(n)
		if err != nil {
			return err
		}
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"input":   input,
				"decimal": n,
				"roman":   roman,
			})
		}
		fmt.Println(roman)
		return nil
	}

	// Otherwise, parse as Roman numeral.
	n, err := romanToInt(input)
	if err != nil {
		return err
	}
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":   input,
			"roman":   strings.ToUpper(input),
			"decimal": n,
		})
	}
	fmt.Println(n)
	return nil
}

// isDecimal reports whether s consists entirely of ASCII digits (with an
// optional leading minus sign).
func isDecimal(s string) bool {
	if len(s) == 0 {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
		start = 1
	}
	if start >= len(s) {
		return false
	}
	for _, r := range s[start:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Roman numeral value table, ordered from largest to smallest.
var romanTable = []struct {
	value  int
	symbol string
}{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

// intToRoman converts a decimal integer (1-3999) to a Roman numeral string.
func intToRoman(n int) (string, error) {
	if n < 1 || n > 3999 {
		return "", fmt.Errorf("value %d is out of Roman numeral range (1-3999)", n)
	}
	var b strings.Builder
	for _, entry := range romanTable {
		for n >= entry.value {
			b.WriteString(entry.symbol)
			n -= entry.value
		}
	}
	return b.String(), nil
}

// romanCharValues maps individual Roman numeral characters to their values.
var romanCharValues = map[rune]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
}

// romanToInt converts a Roman numeral string to a decimal integer.
func romanToInt(s string) (int, error) {
	upper := strings.ToUpper(s)
	if len(upper) == 0 {
		return 0, fmt.Errorf("empty Roman numeral string")
	}

	// Validate characters.
	for _, r := range upper {
		if _, ok := romanCharValues[r]; !ok {
			return 0, fmt.Errorf("invalid Roman numeral character: %c", r)
		}
	}

	total := 0
	runes := []rune(upper)
	for i := 0; i < len(runes); i++ {
		curr := romanCharValues[runes[i]]
		if i+1 < len(runes) {
			next := romanCharValues[runes[i+1]]
			if curr < next {
				// Subtractive notation (e.g., IV = 4, IX = 9).
				total += next - curr
				i++
				continue
			}
		}
		total += curr
	}

	// Round-trip validation: convert back and compare to catch malformed input.
	roundTrip, err := intToRoman(total)
	if err != nil {
		return 0, fmt.Errorf("parsed value %d is out of Roman numeral range", total)
	}
	if roundTrip != upper {
		return 0, fmt.Errorf("non-canonical Roman numeral %q (canonical form: %s = %d)", s, roundTrip, total)
	}

	return total, nil
}

// ─── ieee754 subcommand ─────────────────────────────────────────────────────

var ieee754Cmd = &cobra.Command{
	Use:   "ieee754 <number>",
	Short: "Show IEEE 754 floating-point bit representation",
	Long: `Display the IEEE 754 binary representation of a decimal number,
showing the sign bit, exponent bits, and mantissa (significand) bits
for both 32-bit (float32) and 64-bit (float64) formats.

EXAMPLES:

  openGyver number ieee754 3.14
  openGyver number ieee754 -1.5
  openGyver number ieee754 0
  openGyver number ieee754 inf
  openGyver number ieee754 nan`,
	Args: cobra.ExactArgs(1),
	RunE: runIEEE754,
}

func runIEEE754(_ *cobra.Command, args []string) error {
	input := strings.TrimSpace(args[0])

	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("invalid number: %q", input)
	}

	f32 := float32(val)
	bits32 := math.Float32bits(f32)
	bits64 := math.Float64bits(val)

	sign32 := (bits32 >> 31) & 1
	exp32 := (bits32 >> 23) & 0xFF
	mant32 := bits32 & 0x7FFFFF

	sign64 := (bits64 >> 63) & 1
	exp64 := (bits64 >> 52) & 0x7FF
	mant64 := bits64 & 0xFFFFFFFFFFFFF

	sign32Str := fmt.Sprintf("%01b", sign32)
	exp32Str := fmt.Sprintf("%08b", exp32)
	mant32Str := fmt.Sprintf("%023b", mant32)

	sign64Str := fmt.Sprintf("%01b", sign64)
	exp64Str := fmt.Sprintf("%011b", exp64)
	mant64Str := fmt.Sprintf("%052b", mant64)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input": input,
			"float32": map[string]interface{}{
				"value":    fmt.Sprintf("%g", f32),
				"bits":     fmt.Sprintf("%032b", bits32),
				"hex":      fmt.Sprintf("0x%08X", bits32),
				"sign":     sign32Str,
				"exponent": exp32Str,
				"mantissa": mant32Str,
			},
			"float64": map[string]interface{}{
				"value":    fmt.Sprintf("%g", val),
				"bits":     fmt.Sprintf("%064b", bits64),
				"hex":      fmt.Sprintf("0x%016X", bits64),
				"sign":     sign64Str,
				"exponent": exp64Str,
				"mantissa": mant64Str,
			},
		})
	}

	fmt.Printf("Input: %s\n\n", input)

	fmt.Println("float32:")
	fmt.Printf("  Value:    %g\n", f32)
	fmt.Printf("  Hex:      0x%08X\n", bits32)
	fmt.Printf("  Binary:   %s %s %s\n", sign32Str, exp32Str, mant32Str)
	fmt.Printf("  Sign:     %s (%s)\n", sign32Str, signLabel(sign32))
	fmt.Printf("  Exponent: %s (%d biased, %d unbiased)\n", exp32Str, exp32, int(exp32)-127)
	fmt.Printf("  Mantissa: %s\n", mant32Str)
	fmt.Println()

	fmt.Println("float64:")
	fmt.Printf("  Value:    %g\n", val)
	fmt.Printf("  Hex:      0x%016X\n", bits64)
	fmt.Printf("  Binary:   %s %s %s\n", sign64Str, exp64Str, mant64Str)
	fmt.Printf("  Sign:     %s (%s)\n", sign64Str, signLabel(sign64))
	fmt.Printf("  Exponent: %s (%d biased, %d unbiased)\n", exp64Str, exp64, int(exp64)-1023)
	fmt.Printf("  Mantissa: %s\n", mant64Str)

	return nil
}

func signLabel[T ~uint32 | ~uint64](sign T) string {
	if sign == 0 {
		return "+"
	}
	return "-"
}

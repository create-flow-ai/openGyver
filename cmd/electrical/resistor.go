package electrical

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var resistorCmd = &cobra.Command{
	Use:   "resistor <value>",
	Short: "Resistor color code calculator",
	Long: `Convert a resistance value to its corresponding color band codes for
standard 4-band and 5-band resistors.

INPUT FORMATS:

  4700       Plain number in Ohms
  4.7k       SI suffix (k = kilo, M = mega, R = Ohms)
  4k7        Mixed notation (common in schematics)
  1M         1 megaohm
  470R       470 Ohms (explicit)
  2.2M       2.2 megaohms

COLOR BANDS:

  4-band: [digit1] [digit2] [multiplier] [tolerance ±5%]
  5-band: [digit1] [digit2] [digit3] [multiplier] [tolerance ±1%]

  0 = Black    5 = Green
  1 = Brown    6 = Blue
  2 = Red      7 = Violet
  3 = Orange   8 = Grey
  4 = Yellow   9 = White

  Multipliers: ×1=Black ×10=Brown ×100=Red ×1k=Orange ×10k=Yellow
               ×100k=Green ×1M=Blue ×10M=Violet ×0.1=Gold ×0.01=Silver

EXAMPLES:

  openGyver electrical resistor 4700
  4-band: Yellow Violet Red Gold (4700 Ω ±5%)
  5-band: Yellow Violet Black Brown Brown (4700 Ω ±1%)

  openGyver electrical resistor 4.7k
  4-band: Yellow Violet Red Gold (4700 Ω ±5%)

  openGyver electrical resistor 10k
  4-band: Brown Black Orange Gold (10000 Ω ±5%)

  openGyver electrical resistor 2.2M
  4-band: Red Red Green Gold (2200000 Ω ±5%)`,
	Args: cobra.ExactArgs(1),
	RunE: runResistor,
}

func init() {
	register(resistorCmd)
}

// Color names for resistor bands (indexed by digit 0-9).
var digitColors = []string{
	"Black",  // 0
	"Brown",  // 1
	"Red",    // 2
	"Orange", // 3
	"Yellow", // 4
	"Green",  // 5
	"Blue",   // 6
	"Violet", // 7
	"Grey",   // 8
	"White",  // 9
}

// multiplierColors maps the power-of-10 exponent to a color name.
var multiplierColors = map[int]string{
	-2: "Silver",
	-1: "Gold",
	0:  "Black",
	1:  "Brown",
	2:  "Red",
	3:  "Orange",
	4:  "Yellow",
	5:  "Green",
	6:  "Blue",
	7:  "Violet",
}

// parseResistorValue parses a human-readable resistance value and returns Ohms.
func parseResistorValue(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty resistance value")
	}

	// Handle mixed notation like "4k7" -> 4700
	s = strings.ToUpper(s)
	if idx := strings.Index(s, "K"); idx >= 0 && idx < len(s)-1 {
		// e.g. "4K7" -> "4.7" * 1000
		before := s[:idx]
		after := s[idx+1:]
		if _, err := strconv.ParseFloat(before, 64); err == nil {
			if _, err2 := strconv.ParseFloat(after, 64); err2 == nil {
				val, _ := strconv.ParseFloat(before+"."+after, 64)
				return val * 1e3, nil
			}
		}
	}
	if idx := strings.Index(s, "M"); idx >= 0 && idx < len(s)-1 {
		before := s[:idx]
		after := s[idx+1:]
		if _, err := strconv.ParseFloat(before, 64); err == nil {
			if _, err2 := strconv.ParseFloat(after, 64); err2 == nil {
				val, _ := strconv.ParseFloat(before+"."+after, 64)
				return val * 1e6, nil
			}
		}
	}
	if idx := strings.Index(s, "R"); idx >= 0 && idx < len(s)-1 {
		before := s[:idx]
		after := s[idx+1:]
		if _, err := strconv.ParseFloat(before, 64); err == nil {
			if _, err2 := strconv.ParseFloat(after, 64); err2 == nil {
				val, _ := strconv.ParseFloat(before+"."+after, 64)
				return val, nil
			}
		}
	}

	// Handle suffix notation: "4.7k", "2.2M", "470R"
	if strings.HasSuffix(s, "K") {
		val, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid resistance value: %q", s)
		}
		return val * 1e3, nil
	}
	if strings.HasSuffix(s, "M") {
		val, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid resistance value: %q", s)
		}
		return val * 1e6, nil
	}
	if strings.HasSuffix(s, "R") {
		val, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid resistance value: %q", s)
		}
		return val, nil
	}

	// Plain number.
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid resistance value: %q", s)
	}
	return val, nil
}

// resistorBands holds the color band information.
type resistorBands struct {
	Bands     []string `json:"bands"`
	Value     float64  `json:"value_ohms"`
	Tolerance string   `json:"tolerance"`
	BandCount int      `json:"band_count"`
}

// colorBands4 calculates the 4-band color code for a resistance value.
func colorBands4(ohms float64) (*resistorBands, error) {
	if ohms < 0.01 || ohms > 99e6 {
		return nil, fmt.Errorf("resistance %.2f Ω is out of 4-band range (0.01 Ω to 99 MΩ)", ohms)
	}

	// Find the two significant digits and the multiplier.
	// We need to express ohms as AB × 10^n where AB is a 2-digit number (10-99).
	exp := int(math.Floor(math.Log10(ohms))) - 1
	sig := ohms / math.Pow(10, float64(exp))
	sig = math.Round(sig)

	// Handle rounding that pushes us to 3 digits.
	if sig >= 100 {
		sig /= 10
		exp++
	}

	d1 := int(sig) / 10
	d2 := int(sig) % 10

	multColor, ok := multiplierColors[exp]
	if !ok {
		return nil, fmt.Errorf("cannot represent %.2f Ω as a 4-band resistor (multiplier 10^%d not available)", ohms, exp)
	}

	return &resistorBands{
		Bands:     []string{digitColors[d1], digitColors[d2], multColor, "Gold"},
		Value:     math.Round(float64(int(sig)) * math.Pow(10, float64(exp))),
		Tolerance: "±5%",
		BandCount: 4,
	}, nil
}

// colorBands5 calculates the 5-band color code for a resistance value.
func colorBands5(ohms float64) (*resistorBands, error) {
	if ohms < 0.01 || ohms > 999e6 {
		return nil, fmt.Errorf("resistance %.2f Ω is out of 5-band range (0.01 Ω to 999 MΩ)", ohms)
	}

	// We need to express ohms as ABC × 10^n where ABC is a 3-digit number (100-999).
	exp := int(math.Floor(math.Log10(ohms))) - 2
	sig := ohms / math.Pow(10, float64(exp))
	sig = math.Round(sig)

	// Handle rounding that pushes us to 4 digits.
	if sig >= 1000 {
		sig /= 10
		exp++
	}

	d1 := int(sig) / 100
	d2 := (int(sig) / 10) % 10
	d3 := int(sig) % 10

	multColor, ok := multiplierColors[exp]
	if !ok {
		return nil, fmt.Errorf("cannot represent %.2f Ω as a 5-band resistor (multiplier 10^%d not available)", ohms, exp)
	}

	return &resistorBands{
		Bands:     []string{digitColors[d1], digitColors[d2], digitColors[d3], multColor, "Brown"},
		Value:     math.Round(float64(int(sig)) * math.Pow(10, float64(exp))),
		Tolerance: "±1%",
		BandCount: 5,
	}, nil
}

// formatOhms formats a resistance value with appropriate units.
func formatOhms(ohms float64) string {
	switch {
	case ohms >= 1e6:
		return fmt.Sprintf("%.2g MΩ", ohms/1e6)
	case ohms >= 1e3:
		return fmt.Sprintf("%.4g kΩ", ohms/1e3)
	default:
		return fmt.Sprintf("%.4g Ω", ohms)
	}
}

func runResistor(_ *cobra.Command, args []string) error {
	ohms, err := parseResistorValue(args[0])
	if err != nil {
		return err
	}

	if ohms <= 0 {
		return fmt.Errorf("resistance must be a positive value, got %.2f", ohms)
	}

	band4, err4 := colorBands4(ohms)
	band5, err5 := colorBands5(ohms)

	if err4 != nil && err5 != nil {
		return fmt.Errorf("could not calculate color bands: %v", err4)
	}

	if jsonOut {
		result := map[string]interface{}{
			"input":      args[0],
			"value_ohms": ohms,
			"formatted":  formatOhms(ohms),
		}
		if band4 != nil {
			result["4_band"] = map[string]interface{}{
				"bands":     band4.Bands,
				"tolerance": band4.Tolerance,
			}
		}
		if band5 != nil {
			result["5_band"] = map[string]interface{}{
				"bands":     band5.Bands,
				"tolerance": band5.Tolerance,
			}
		}
		return cmd.PrintJSON(result)
	}

	fmt.Printf("Value: %s (%g Ω)\n\n", formatOhms(ohms), ohms)

	if band4 != nil {
		fmt.Printf("4-band: %s  (%g Ω %s)\n",
			strings.Join(band4.Bands, "  "),
			band4.Value, band4.Tolerance)
	}
	if band5 != nil {
		fmt.Printf("5-band: %s  (%g Ω %s)\n",
			strings.Join(band5.Bands, "  "),
			band5.Value, band5.Tolerance)
	}

	return nil
}

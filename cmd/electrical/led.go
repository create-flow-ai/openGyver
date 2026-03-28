package electrical

import (
	"fmt"
	"math"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	ledSource  float64
	ledForward float64
	ledCurrent float64
)

var ledCmd = &cobra.Command{
	Use:   "led",
	Short: "LED resistor calculator",
	Long: `Calculate the required current-limiting resistor for an LED circuit.

Uses the formula: R = (Vsource - Vforward) / I

where:
  Vsource  = supply voltage (e.g. 5V, 12V)
  Vforward = LED forward voltage drop (depends on LED color)
  I        = desired current through the LED

Typical LED forward voltages:
  Red:         1.8 - 2.2 V
  Orange:      2.0 - 2.2 V
  Yellow:      2.0 - 2.2 V
  Green:       2.0 - 3.5 V
  Blue:        3.0 - 3.5 V
  White:       3.0 - 3.5 V
  Infrared:    1.1 - 1.5 V
  UV:          3.3 - 4.0 V

FLAGS:

  --source    Supply voltage in Volts (required)
  --forward   LED forward voltage in Volts (required)
  --current   Desired LED current in milliamps (default: 20 mA)

EXAMPLES:

  # Red LED on 5V supply (typical)
  openGyver electrical led --source 5 --forward 2.0
  Resistor: 150 Ω   Power: 60.0 mW

  # Blue LED on 3.3V supply
  openGyver electrical led --source 3.3 --forward 3.2 --current 10
  Resistor: 10 Ω   Power: 0.1 mW

  # White LED on 12V supply
  openGyver electrical led --source 12 --forward 3.3 --current 20
  Resistor: 435 Ω   Power: 174.0 mW

  # Multiple LEDs in series (sum forward voltages)
  openGyver electrical led --source 12 --forward 6.6 --current 20
  Resistor: 270 Ω   Power: 108.0 mW`,
	Args: cobra.NoArgs,
	RunE: runLED,
}

func init() {
	ledCmd.Flags().Float64Var(&ledSource, "source", 0, "supply voltage in Volts (V)")
	ledCmd.Flags().Float64Var(&ledForward, "forward", 0, "LED forward voltage in Volts (V)")
	ledCmd.Flags().Float64Var(&ledCurrent, "current", 20, "desired LED current in milliamps (mA)")
	_ = ledCmd.MarkFlagRequired("source")
	_ = ledCmd.MarkFlagRequired("forward")
	register(ledCmd)
}

// Standard E24 resistor values (multiplied to cover ranges).
var e24Values = []float64{
	1.0, 1.1, 1.2, 1.3, 1.5, 1.6, 1.8, 2.0, 2.2, 2.4, 2.7, 3.0,
	3.3, 3.6, 3.9, 4.3, 4.7, 5.1, 5.6, 6.2, 6.8, 7.5, 8.2, 9.1,
}

// nearestStandardResistor finds the nearest E24 series resistor value
// that is greater than or equal to the calculated value.
func nearestStandardResistor(ohms float64) float64 {
	if ohms <= 0 {
		return 1.0
	}
	exp := math.Floor(math.Log10(ohms))
	mantissa := ohms / math.Pow(10, exp)

	// Find the smallest E24 value >= mantissa.
	for _, v := range e24Values {
		if v >= mantissa-0.001 {
			return v * math.Pow(10, exp)
		}
	}
	// Wrap around to next decade.
	return e24Values[0] * math.Pow(10, exp+1)
}

func runLED(_ *cobra.Command, _ []string) error {
	if ledSource <= 0 {
		return fmt.Errorf("--source must be a positive voltage, got %.2f V", ledSource)
	}
	if ledForward < 0 {
		return fmt.Errorf("--forward must be non-negative, got %.2f V", ledForward)
	}
	if ledCurrent <= 0 {
		return fmt.Errorf("--current must be positive, got %.2f mA", ledCurrent)
	}
	if ledForward >= ledSource {
		return fmt.Errorf("forward voltage (%.2f V) must be less than source voltage (%.2f V)", ledForward, ledSource)
	}

	// Convert mA to A.
	currentA := ledCurrent / 1000.0

	// R = (Vsource - Vforward) / I
	resistance := (ledSource - ledForward) / currentA

	// Power dissipated by the resistor: P = I² × R  or  P = (Vs - Vf) × I
	power := (ledSource - ledForward) * currentA

	// Find nearest standard value (round up for safety).
	stdResistance := nearestStandardResistor(resistance)
	// Actual current with standard resistor.
	actualCurrentA := (ledSource - ledForward) / stdResistance
	actualPower := (ledSource - ledForward) * actualCurrentA

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"source_V":                    ledSource,
			"forward_V":                   ledForward,
			"desired_current_mA":          ledCurrent,
			"calculated_resistance_Ω":     math.Round(resistance*100) / 100,
			"calculated_power_mW":         math.Round(power*1e5) / 100,
			"nearest_standard_Ω":          stdResistance,
			"actual_current_mA":           math.Round(actualCurrentA*1e5) / 100,
			"actual_power_mW":             math.Round(actualPower*1e5) / 100,
		})
	}

	fmt.Printf("Source voltage:     %.2f V\n", ledSource)
	fmt.Printf("Forward voltage:    %.2f V\n", ledForward)
	fmt.Printf("Voltage across R:   %.2f V\n", ledSource-ledForward)
	fmt.Printf("Desired current:    %.2f mA\n", ledCurrent)
	fmt.Println()
	fmt.Printf("Calculated R:       %s\n", formatSI(resistance, "Ω"))
	fmt.Printf("Power dissipation:  %s\n", formatSI(power, "W"))
	fmt.Println()
	fmt.Printf("Nearest standard:   %s (E24 series)\n", formatSI(stdResistance, "Ω"))
	fmt.Printf("Actual current:     %.2f mA\n", actualCurrentA*1000)
	fmt.Printf("Actual power:       %s\n", formatSI(actualPower, "W"))

	return nil
}

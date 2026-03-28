package electrical

import (
	"fmt"
	"math"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	ohmVoltage    float64
	ohmCurrent    float64
	ohmResistance float64
	ohmVSet       bool
	ohmISet       bool
	ohmRSet       bool
)

var ohmCmd = &cobra.Command{
	Use:   "ohm",
	Short: "Ohm's law calculator",
	Long: `Calculate voltage, current, resistance, and power using Ohm's law.

Provide any two of the three fundamental values and the calculator
will determine the third, plus the dissipated power.

  V = I × R    (Voltage = Current × Resistance)
  I = V / R    (Current = Voltage / Resistance)
  R = V / I    (Resistance = Voltage / Current)
  P = V × I    (Power = Voltage × Current)

FLAGS:

  --voltage    / -v   Voltage in Volts (V)
  --current    / -i   Current in Amps (A)
  --resistance / -r   Resistance in Ohms (Ω)

EXAMPLES:

  # Given voltage and current, find resistance and power
  openGyver electrical ohm --voltage 12 --current 0.5
  Resistance: 24.00 Ω    Power: 6.00 W

  # Given voltage and resistance, find current and power
  openGyver electrical ohm -v 5 -r 1000
  Current: 0.005 A (5.00 mA)    Power: 0.025 W (25.00 mW)

  # Given current and resistance, find voltage and power
  openGyver electrical ohm -i 0.02 -r 220
  Voltage: 4.40 V    Power: 0.088 W (88.00 mW)`,
	Args: cobra.NoArgs,
	RunE: runOhm,
}

func init() {
	ohmCmd.Flags().Float64VarP(&ohmVoltage, "voltage", "v", 0, "voltage in Volts (V)")
	ohmCmd.Flags().Float64VarP(&ohmCurrent, "current", "i", 0, "current in Amps (A)")
	ohmCmd.Flags().Float64VarP(&ohmResistance, "resistance", "r", 0, "resistance in Ohms (Ω)")
	register(ohmCmd)
}

// formatSI formats a value with an appropriate SI prefix for readability.
func formatSI(value float64, unit string) string {
	abs := math.Abs(value)
	switch {
	case abs >= 1e9:
		return fmt.Sprintf("%.2f G%s", value/1e9, unit)
	case abs >= 1e6:
		return fmt.Sprintf("%.2f M%s", value/1e6, unit)
	case abs >= 1e3:
		return fmt.Sprintf("%.2f k%s", value/1e3, unit)
	case abs >= 1:
		return fmt.Sprintf("%.2f %s", value, unit)
	case abs >= 1e-3:
		return fmt.Sprintf("%.2f m%s", value*1e3, unit)
	case abs >= 1e-6:
		return fmt.Sprintf("%.2f μ%s", value*1e6, unit)
	case abs >= 1e-9:
		return fmt.Sprintf("%.2f n%s", value*1e9, unit)
	default:
		return fmt.Sprintf("%.6g %s", value, unit)
	}
}

func runOhm(c *cobra.Command, _ []string) error {
	ohmVSet = c.Flags().Changed("voltage")
	ohmISet = c.Flags().Changed("current")
	ohmRSet = c.Flags().Changed("resistance")

	count := 0
	if ohmVSet {
		count++
	}
	if ohmISet {
		count++
	}
	if ohmRSet {
		count++
	}

	if count < 2 {
		return fmt.Errorf("provide exactly 2 of --voltage, --current, --resistance (got %d)", count)
	}
	if count > 2 {
		return fmt.Errorf("provide exactly 2 of --voltage, --current, --resistance (got %d)", count)
	}

	var voltage, current, resistance, power float64
	var calculated string

	switch {
	case ohmVSet && ohmISet:
		// Calculate resistance.
		voltage = ohmVoltage
		current = ohmCurrent
		if current == 0 {
			return fmt.Errorf("current cannot be zero when calculating resistance")
		}
		resistance = voltage / current
		power = voltage * current
		calculated = "resistance"

	case ohmVSet && ohmRSet:
		// Calculate current.
		voltage = ohmVoltage
		resistance = ohmResistance
		if resistance == 0 {
			return fmt.Errorf("resistance cannot be zero when calculating current")
		}
		current = voltage / resistance
		power = voltage * current
		calculated = "current"

	case ohmISet && ohmRSet:
		// Calculate voltage.
		current = ohmCurrent
		resistance = ohmResistance
		voltage = current * resistance
		power = voltage * current
		calculated = "voltage"
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"voltage_V":    math.Round(voltage*1e6) / 1e6,
			"current_A":    math.Round(current*1e6) / 1e6,
			"resistance_Ω": math.Round(resistance*1e6) / 1e6,
			"power_W":      math.Round(power*1e6) / 1e6,
			"calculated":   calculated,
		})
	}

	fmt.Printf("Voltage:    %s\n", formatSI(voltage, "V"))
	fmt.Printf("Current:    %s\n", formatSI(current, "A"))
	fmt.Printf("Resistance: %s\n", formatSI(resistance, "Ω"))
	fmt.Printf("Power:      %s\n", formatSI(power, "W"))
	fmt.Printf("\n(Calculated: %s)\n", calculated)

	return nil
}

package electrical

import (
	"fmt"
	"math"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	divVin float64
	divR1  float64
	divR2  float64
)

var dividerCmd = &cobra.Command{
	Use:   "divider",
	Short: "Voltage divider calculator",
	Long: `Calculate the output voltage of a resistive voltage divider.

A voltage divider uses two resistors in series to produce a lower
output voltage from a higher input voltage:

         Vin
          |
         [R1]
          |──── Vout
         [R2]
          |
         GND

  Vout = Vin × R2 / (R1 + R2)

Also calculates the divider ratio, current through the divider,
and power dissipated by each resistor.

FLAGS:

  --vin   Input voltage in Volts (required)
  --r1    Upper resistor in Ohms (required)
  --r2    Lower resistor in Ohms (required)

EXAMPLES:

  # 12V down to ~3.84V with 10kΩ and 4.7kΩ
  openGyver electrical divider --vin 12 --r1 10000 --r2 4700
  Vout: 3.84 V   Ratio: 0.3197

  # 5V to 3.3V level shifting (close approximation)
  openGyver electrical divider --vin 5 --r1 1800 --r2 3300
  Vout: 3.24 V   Ratio: 0.6471

  # 24V to 5V for ADC input
  openGyver electrical divider --vin 24 --r1 38000 --r2 10000
  Vout: 5.00 V   Ratio: 0.2083

  # Check divider current draw
  openGyver electrical divider --vin 12 --r1 100000 --r2 100000
  Vout: 6.00 V   Current: 0.06 mA`,
	Args: cobra.NoArgs,
	RunE: runDivider,
}

func init() {
	dividerCmd.Flags().Float64Var(&divVin, "vin", 0, "input voltage in Volts (V)")
	dividerCmd.Flags().Float64Var(&divR1, "r1", 0, "upper resistor in Ohms (Ω)")
	dividerCmd.Flags().Float64Var(&divR2, "r2", 0, "lower resistor in Ohms (Ω)")
	_ = dividerCmd.MarkFlagRequired("vin")
	_ = dividerCmd.MarkFlagRequired("r1")
	_ = dividerCmd.MarkFlagRequired("r2")
	register(dividerCmd)
}

func runDivider(_ *cobra.Command, _ []string) error {
	if divR1 <= 0 {
		return fmt.Errorf("--r1 must be positive, got %.2f", divR1)
	}
	if divR2 <= 0 {
		return fmt.Errorf("--r2 must be positive, got %.2f", divR2)
	}

	// Vout = Vin * R2 / (R1 + R2)
	vout := divVin * divR2 / (divR1 + divR2)

	// Divider ratio.
	ratio := divR2 / (divR1 + divR2)

	// Total current through the divider: I = Vin / (R1 + R2)
	current := divVin / (divR1 + divR2)

	// Power dissipated by each resistor: P = I² × R
	powerR1 := current * current * divR1
	powerR2 := current * current * divR2
	powerTotal := powerR1 + powerR2

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"vin_V":         divVin,
			"r1_Ω":          divR1,
			"r2_Ω":          divR2,
			"vout_V":        math.Round(vout*1e4) / 1e4,
			"ratio":         math.Round(ratio*1e6) / 1e6,
			"current_A":     math.Round(current*1e6) / 1e6,
			"power_r1_W":    math.Round(powerR1*1e6) / 1e6,
			"power_r2_W":    math.Round(powerR2*1e6) / 1e6,
			"power_total_W": math.Round(powerTotal*1e6) / 1e6,
		})
	}

	fmt.Printf("Input voltage:  %.2f V\n", divVin)
	fmt.Printf("R1 (upper):     %s\n", formatSI(divR1, "Ω"))
	fmt.Printf("R2 (lower):     %s\n", formatSI(divR2, "Ω"))
	fmt.Println()
	fmt.Printf("Output voltage: %.4f V\n", vout)
	fmt.Printf("Divider ratio:  %.6f  (R2 / [R1+R2])\n", ratio)
	fmt.Printf("Current:        %s\n", formatSI(current, "A"))
	fmt.Println()
	fmt.Printf("Power (R1):     %s\n", formatSI(powerR1, "W"))
	fmt.Printf("Power (R2):     %s\n", formatSI(powerR2, "W"))
	fmt.Printf("Power (total):  %s\n", formatSI(powerTotal, "W"))

	return nil
}

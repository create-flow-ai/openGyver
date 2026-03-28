package electrical

import (
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var electricalCmd = &cobra.Command{
	Use:   "electrical",
	Short: "Electrical engineering calculators and tools",
	Long: `Electrical engineering tools for circuit design and component calculation.

SUBCOMMANDS:

  ohm       Ohm's law calculator (voltage, current, resistance, power)
  resistor  Resistor color code calculator (value to band colors)
  led       LED resistor calculator (find the right current-limiting resistor)
  divider   Voltage divider calculator (Vout from Vin, R1, R2)

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver electrical ohm --voltage 12 --current 0.5
  openGyver electrical resistor 4700
  openGyver electrical resistor 4.7k
  openGyver electrical led --source 5 --forward 2.1 --current 20
  openGyver electrical divider --vin 12 --r1 10000 --r2 4700`,
}

// jsonOut controls JSON output across all subcommands.
var jsonOut bool

func init() {
	electricalCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	cmd.Register(electricalCmd)
}

// register adds a subcommand to the electrical parent command.
func register(sub *cobra.Command) {
	electricalCmd.AddCommand(sub)
}

package geo

import (
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var geoCmd = &cobra.Command{
	Use:   "geo",
	Short: "Geolocation conversion and calculation tools",
	Long: `Geolocation tools for distance calculation, coordinate conversion, and more.

SUBCOMMANDS:

  distance  Calculate Haversine (great-circle) distance between two coordinates
  dms       Convert between decimal degrees and DMS (degrees/minutes/seconds)
  utm       Convert latitude/longitude to UTM coordinates

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver geo distance --lat1 40.7128 --lon1 -74.0060 --lat2 51.5074 --lon2 -0.1278
  openGyver geo dms 40.7128
  openGyver geo dms "40°42'46.1\"N"
  openGyver geo utm --lat 40.7128 --lon -74.0060`,
}

// jsonOut controls JSON output across all subcommands.
var jsonOut bool

func init() {
	geoCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	cmd.Register(geoCmd)
}

// register adds a subcommand to the geo parent command.
func register(sub *cobra.Command) {
	geoCmd.AddCommand(sub)
}

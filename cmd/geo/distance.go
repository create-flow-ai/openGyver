package geo

import (
	"fmt"
	"math"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// Earth's mean radius in kilometres.
const earthRadiusKm = 6371.0

// Conversion factor from kilometres to miles.
const kmToMiles = 0.621371

var (
	distLat1 float64
	distLon1 float64
	distLat2 float64
	distLon2 float64
)

var distanceCmd = &cobra.Command{
	Use:   "distance",
	Short: "Calculate Haversine (great-circle) distance between two coordinates",
	Long: `Calculate the great-circle distance between two points on Earth using
the Haversine formula.

The Haversine formula accounts for Earth's curvature and gives an
accurate distance for any two points on the globe:

  a = sin²(Δφ/2) + cos(φ1)·cos(φ2)·sin²(Δλ/2)
  c = 2·atan2(√a, √(1−a))
  d = R·c

where R is Earth's mean radius (6,371 km).

FLAGS:

  --lat1   Latitude of point 1  (decimal degrees, required)
  --lon1   Longitude of point 1 (decimal degrees, required)
  --lat2   Latitude of point 2  (decimal degrees, required)
  --lon2   Longitude of point 2 (decimal degrees, required)

EXAMPLES:

  # New York to London
  openGyver geo distance --lat1 40.7128 --lon1 -74.0060 --lat2 51.5074 --lon2 -0.1278
  5570.25 km (3461.02 mi)

  # Tokyo to Sydney
  openGyver geo distance --lat1 35.6762 --lon1 139.6503 --lat2 -33.8688 --lon2 151.2093
  7823.27 km (4861.78 mi)

  # San Francisco to Paris
  openGyver geo distance --lat1 37.7749 --lon1 -122.4194 --lat2 48.8566 --lon2 2.3522
  8964.85 km (5570.22 mi)`,
	Args: cobra.NoArgs,
	RunE: runDistance,
}

func init() {
	distanceCmd.Flags().Float64Var(&distLat1, "lat1", 0, "latitude of point 1 (decimal degrees)")
	distanceCmd.Flags().Float64Var(&distLon1, "lon1", 0, "longitude of point 1 (decimal degrees)")
	distanceCmd.Flags().Float64Var(&distLat2, "lat2", 0, "latitude of point 2 (decimal degrees)")
	distanceCmd.Flags().Float64Var(&distLon2, "lon2", 0, "longitude of point 2 (decimal degrees)")
	_ = distanceCmd.MarkFlagRequired("lat1")
	_ = distanceCmd.MarkFlagRequired("lon1")
	_ = distanceCmd.MarkFlagRequired("lat2")
	_ = distanceCmd.MarkFlagRequired("lon2")
	register(distanceCmd)
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// haversine calculates the great-circle distance between two points
// on a sphere given their latitudes and longitudes in decimal degrees.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	phi1 := degreesToRadians(lat1)
	phi2 := degreesToRadians(lat2)
	deltaPhi := degreesToRadians(lat2 - lat1)
	deltaLambda := degreesToRadians(lon2 - lon1)

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func runDistance(_ *cobra.Command, _ []string) error {
	// Validate latitude and longitude ranges.
	if distLat1 < -90 || distLat1 > 90 {
		return fmt.Errorf("--lat1 must be between -90 and 90, got %f", distLat1)
	}
	if distLat2 < -90 || distLat2 > 90 {
		return fmt.Errorf("--lat2 must be between -90 and 90, got %f", distLat2)
	}
	if distLon1 < -180 || distLon1 > 180 {
		return fmt.Errorf("--lon1 must be between -180 and 180, got %f", distLon1)
	}
	if distLon2 < -180 || distLon2 > 180 {
		return fmt.Errorf("--lon2 must be between -180 and 180, got %f", distLon2)
	}

	km := haversine(distLat1, distLon1, distLat2, distLon2)
	mi := km * kmToMiles

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"point1": map[string]float64{
				"latitude":  distLat1,
				"longitude": distLon1,
			},
			"point2": map[string]float64{
				"latitude":  distLat2,
				"longitude": distLon2,
			},
			"distance_km":    math.Round(km*100) / 100,
			"distance_miles": math.Round(mi*100) / 100,
		})
	}

	fmt.Printf("Point 1:   %.6f, %.6f\n", distLat1, distLon1)
	fmt.Printf("Point 2:   %.6f, %.6f\n", distLat2, distLon2)
	fmt.Printf("Distance:  %.2f km (%.2f mi)\n", km, mi)

	return nil
}

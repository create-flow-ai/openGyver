package geo

import (
	"fmt"
	"math"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	utmLat float64
	utmLon float64
)

var utmCmd = &cobra.Command{
	Use:   "utm",
	Short: "Convert latitude/longitude to UTM coordinates",
	Long: `Convert geographic coordinates (latitude/longitude in decimal degrees)
to Universal Transverse Mercator (UTM) coordinates.

UTM divides the Earth into 60 zones (1-60), each 6° wide. Within each
zone, locations are specified by easting (metres from the zone's central
meridian, offset by 500,000 m) and northing (metres from the equator,
with a 10,000,000 m false northing for the southern hemisphere).

The conversion uses the WGS 84 ellipsoid parameters.

FLAGS:

  --lat   Latitude  in decimal degrees (required, -80 to 84)
  --lon   Longitude in decimal degrees (required, -180 to 180)

EXAMPLES:

  # Statue of Liberty, New York
  openGyver geo utm --lat 40.6892 --lon -74.0445
  Zone: 18N   Easting: 583960.00   Northing: 4507523.00

  # Sydney Opera House
  openGyver geo utm --lat -33.8568 --lon 151.2153
  Zone: 56S   Easting: 334841.00   Northing: 6252080.00

  # Null Island (0°, 0°)
  openGyver geo utm --lat 0 --lon 0
  Zone: 31N   Easting: 166021.00   Northing: 0.00`,
	Args: cobra.NoArgs,
	RunE: runUTM,
}

func init() {
	utmCmd.Flags().Float64Var(&utmLat, "lat", 0, "latitude in decimal degrees (-80 to 84)")
	utmCmd.Flags().Float64Var(&utmLon, "lon", 0, "longitude in decimal degrees (-180 to 180)")
	_ = utmCmd.MarkFlagRequired("lat")
	_ = utmCmd.MarkFlagRequired("lon")
	register(utmCmd)
}

// WGS 84 ellipsoid constants.
const (
	wgs84A  = 6378137.0         // semi-major axis (m)
	wgs84F  = 1 / 298.257223563 // flattening
	wgs84E2 = 2*wgs84F - wgs84F*wgs84F
	utmK0   = 0.9996 // UTM scale factor
)

// utmResult holds the result of a lat/lon to UTM conversion.
type utmResult struct {
	Zone     int     `json:"zone"`
	Letter   string  `json:"letter"`
	Easting  float64 `json:"easting"`
	Northing float64 `json:"northing"`
}

// latLonToUTM converts latitude and longitude (decimal degrees) to UTM.
func latLonToUTM(lat, lon float64) utmResult {
	// UTM zone number.
	zone := int(math.Floor((lon+180)/6)) + 1

	// Handle special zones for Svalbard and Norway.
	if lat >= 56 && lat < 64 && lon >= 3 && lon < 12 {
		zone = 32
	}
	if lat >= 72 && lat < 84 {
		switch {
		case lon >= 0 && lon < 9:
			zone = 31
		case lon >= 9 && lon < 21:
			zone = 33
		case lon >= 21 && lon < 33:
			zone = 35
		case lon >= 33 && lon < 42:
			zone = 37
		}
	}

	// Central meridian of the zone.
	lonOrigin := float64((zone-1)*6-180) + 3

	// Latitude band letter.
	letter := utmLetterDesignator(lat)

	// Convert to radians.
	latRad := lat * math.Pi / 180
	lonRad := lon * math.Pi / 180
	lonOriginRad := lonOrigin * math.Pi / 180

	ePrime2 := wgs84E2 / (1 - wgs84E2)

	n := wgs84A / math.Sqrt(1-wgs84E2*math.Sin(latRad)*math.Sin(latRad))
	t := math.Tan(latRad) * math.Tan(latRad)
	c := ePrime2 * math.Cos(latRad) * math.Cos(latRad)
	a := math.Cos(latRad) * (lonRad - lonOriginRad)

	// Meridional arc length.
	m := wgs84A * ((1-wgs84E2/4-3*wgs84E2*wgs84E2/64-5*wgs84E2*wgs84E2*wgs84E2/256)*latRad -
		(3*wgs84E2/8+3*wgs84E2*wgs84E2/32+45*wgs84E2*wgs84E2*wgs84E2/1024)*math.Sin(2*latRad) +
		(15*wgs84E2*wgs84E2/256+45*wgs84E2*wgs84E2*wgs84E2/1024)*math.Sin(4*latRad) -
		(35*wgs84E2*wgs84E2*wgs84E2/3072)*math.Sin(6*latRad))

	easting := utmK0*n*(a+(1-t+c)*a*a*a/6+
		(5-18*t+t*t+72*c-58*ePrime2)*a*a*a*a*a/120) + 500000.0

	northing := utmK0 * (m + n*math.Tan(latRad)*(a*a/2+(5-t+9*c+4*c*c)*a*a*a*a/24+
		(61-58*t+t*t+600*c-330*ePrime2)*a*a*a*a*a*a/720))

	if lat < 0 {
		northing += 10000000.0
	}

	return utmResult{
		Zone:     zone,
		Letter:   letter,
		Easting:  math.Round(easting*100) / 100,
		Northing: math.Round(northing*100) / 100,
	}
}

// utmLetterDesignator returns the latitude band letter for a given latitude.
func utmLetterDesignator(lat float64) string {
	switch {
	case lat >= 72:
		return "X"
	case lat >= 64:
		return "W"
	case lat >= 56:
		return "V"
	case lat >= 48:
		return "U"
	case lat >= 40:
		return "T"
	case lat >= 32:
		return "S"
	case lat >= 24:
		return "R"
	case lat >= 16:
		return "Q"
	case lat >= 8:
		return "P"
	case lat >= 0:
		return "N"
	case lat >= -8:
		return "M"
	case lat >= -16:
		return "L"
	case lat >= -24:
		return "K"
	case lat >= -32:
		return "J"
	case lat >= -40:
		return "H"
	case lat >= -48:
		return "G"
	case lat >= -56:
		return "F"
	case lat >= -64:
		return "E"
	case lat >= -72:
		return "D"
	default:
		return "C"
	}
}

func runUTM(_ *cobra.Command, _ []string) error {
	if utmLat < -80 || utmLat > 84 {
		return fmt.Errorf("--lat must be between -80 and 84 for UTM, got %f", utmLat)
	}
	if utmLon < -180 || utmLon > 180 {
		return fmt.Errorf("--lon must be between -180 and 180, got %f", utmLon)
	}

	result := latLonToUTM(utmLat, utmLon)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input": map[string]float64{
				"latitude":  utmLat,
				"longitude": utmLon,
			},
			"zone":     result.Zone,
			"letter":   result.Letter,
			"easting":  result.Easting,
			"northing": result.Northing,
			"zone_designator": fmt.Sprintf("%d%s", result.Zone, result.Letter),
		})
	}

	fmt.Printf("Input:     %.6f, %.6f\n", utmLat, utmLon)
	fmt.Printf("Zone:      %d%s\n", result.Zone, result.Letter)
	fmt.Printf("Easting:   %.2f m\n", result.Easting)
	fmt.Printf("Northing:  %.2f m\n", result.Northing)

	return nil
}

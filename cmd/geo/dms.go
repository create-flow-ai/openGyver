package geo

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var dmsCmd = &cobra.Command{
	Use:   "dms <coordinate>",
	Short: "Convert between decimal degrees and DMS notation",
	Long: `Convert a coordinate between decimal degrees and DMS (degrees, minutes,
seconds) notation. The input format is auto-detected:

  • If the input is a plain number (e.g. "40.7128"), it is treated as
    decimal degrees and converted to DMS.
  • If the input contains degree symbols or DMS notation (e.g.
    "40°42'46.1\"N"), it is parsed as DMS and converted to decimal degrees.

DMS INPUT FORMATS (flexible parsing):

  40°42'46.1"N          Full DMS with direction
  40°42'46.1"           DMS without direction
  40 42 46.1 N          Space-separated
  -40 42 46.1           Negative sign for S/W
  40d42m46.1s           Letter notation

EXAMPLES:

  # Decimal degrees to DMS
  openGyver geo dms 40.7128
  40° 42' 46.08" N

  # Negative decimal degrees to DMS
  openGyver geo dms -- -74.0060
  74° 0' 21.60" W

  # DMS to decimal degrees
  openGyver geo dms "40°42'46.1\"N"
  40.712806

  # Space-separated DMS to decimal degrees
  openGyver geo dms "40 42 46.1 N"
  40.712806`,
	Args: cobra.ExactArgs(1),
	RunE: runDMS,
}

func init() {
	register(dmsCmd)
}

// dmsResult holds a DMS representation.
type dmsResult struct {
	Degrees   int     `json:"degrees"`
	Minutes   int     `json:"minutes"`
	Seconds   float64 `json:"seconds"`
	Direction string  `json:"direction"`
}

// dmsPattern matches various DMS formats.
// Group 1: degrees, Group 2: minutes, Group 3: seconds, Group 4: direction
var dmsPattern = regexp.MustCompile(
	`(?i)^\s*(-?)(\d+)\s*[°d]\s*(\d+)\s*[''m]\s*([\d.]+)\s*["″s]?\s*([NSEW])?\s*$`,
)

// spaceDMSPattern matches space-separated DMS: "40 42 46.1 N"
var spaceDMSPattern = regexp.MustCompile(
	`(?i)^\s*(-?)(\d+)\s+(\d+)\s+([\d.]+)\s*([NSEW])?\s*$`,
)

// isDMS returns true if the input looks like a DMS coordinate.
func isDMS(s string) bool {
	return dmsPattern.MatchString(s) || spaceDMSPattern.MatchString(s)
}

// parseDMS parses a DMS string and returns decimal degrees.
func parseDMS(s string) (float64, *dmsResult, error) {
	var (
		neg     bool
		degStr  string
		minStr  string
		secStr  string
		dirStr  string
		matches []string
	)

	if matches = dmsPattern.FindStringSubmatch(s); matches != nil {
		neg = matches[1] == "-"
		degStr = matches[2]
		minStr = matches[3]
		secStr = matches[4]
		dirStr = strings.ToUpper(matches[5])
	} else if matches = spaceDMSPattern.FindStringSubmatch(s); matches != nil {
		neg = matches[1] == "-"
		degStr = matches[2]
		minStr = matches[3]
		secStr = matches[4]
		dirStr = strings.ToUpper(matches[5])
	} else {
		return 0, nil, fmt.Errorf("unrecognised DMS format: %q", s)
	}

	deg, _ := strconv.Atoi(degStr)
	min, _ := strconv.Atoi(minStr)
	sec, _ := strconv.ParseFloat(secStr, 64)

	if min < 0 || min >= 60 {
		return 0, nil, fmt.Errorf("minutes out of range: %d", min)
	}
	if sec < 0 || sec >= 60 {
		return 0, nil, fmt.Errorf("seconds out of range: %f", sec)
	}

	decimal := float64(deg) + float64(min)/60.0 + sec/3600.0
	if neg || dirStr == "S" || dirStr == "W" {
		decimal = -decimal
	}

	dir := dirStr
	if dir == "" {
		if neg {
			dir = "-"
		} else {
			dir = "+"
		}
	}

	result := &dmsResult{
		Degrees:   deg,
		Minutes:   min,
		Seconds:   sec,
		Direction: dir,
	}

	return decimal, result, nil
}

// decimalToDMS converts decimal degrees to DMS components.
func decimalToDMS(dd float64, isLongitude bool) dmsResult {
	var direction string
	abs := math.Abs(dd)

	if isLongitude {
		if dd >= 0 {
			direction = "E"
		} else {
			direction = "W"
		}
	} else {
		if dd >= 0 {
			direction = "N"
		} else {
			direction = "S"
		}
	}

	degrees := int(abs)
	minFloat := (abs - float64(degrees)) * 60
	minutes := int(minFloat)
	seconds := (minFloat - float64(minutes)) * 60

	// Avoid 60-second rollover due to floating point.
	if seconds >= 59.995 {
		seconds = 0
		minutes++
	}
	if minutes >= 60 {
		minutes = 0
		degrees++
	}

	return dmsResult{
		Degrees:   degrees,
		Minutes:   minutes,
		Seconds:   math.Round(seconds*100) / 100,
		Direction: direction,
	}
}

func runDMS(_ *cobra.Command, args []string) error {
	input := strings.TrimSpace(args[0])

	if isDMS(input) {
		// DMS -> decimal degrees
		dd, parsed, err := parseDMS(input)
		if err != nil {
			return err
		}
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"input":           input,
				"decimal_degrees": math.Round(dd*1e6) / 1e6,
				"dms": map[string]interface{}{
					"degrees":   parsed.Degrees,
					"minutes":   parsed.Minutes,
					"seconds":   parsed.Seconds,
					"direction": parsed.Direction,
				},
			})
		}
		fmt.Printf("%f\n", math.Round(dd*1e6)/1e6)
		return nil
	}

	// Decimal degrees -> DMS
	dd, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("could not parse %q as decimal degrees or DMS: %w", input, err)
	}

	if dd < -180 || dd > 180 {
		return fmt.Errorf("value %f is out of range (-180 to 180)", dd)
	}

	// Determine if latitude or longitude based on range heuristic:
	// if |dd| <= 90, assume latitude; otherwise longitude.
	isLon := math.Abs(dd) > 90
	result := decimalToDMS(dd, isLon)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":           input,
			"decimal_degrees": dd,
			"dms": map[string]interface{}{
				"degrees":   result.Degrees,
				"minutes":   result.Minutes,
				"seconds":   result.Seconds,
				"direction": result.Direction,
			},
			"formatted": fmt.Sprintf(`%d° %d' %.2f" %s`, result.Degrees, result.Minutes, result.Seconds, result.Direction),
		})
	}

	fmt.Printf(`%d° %d' %.2f" %s`+"\n", result.Degrees, result.Minutes, result.Seconds, result.Direction)
	return nil
}

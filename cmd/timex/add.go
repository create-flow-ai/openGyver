package timex

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	addFrom string
	addTZ   string
)

var addCmd = &cobra.Command{
	Use:   "add <time> <duration>",
	Short: "Add or subtract a duration from a time",
	Long: `Parse a time and add (or subtract) a duration, then display the result.

DURATION FORMAT:

  Go-style:    1h30m, 2h, 45m, 90s, 500ms, 1h30m45s
  Extended:    30d, 2w, 3mo, 1y (days, weeks, months, years)
  Combined:    1y2mo3d4h5m6s
  Negative:    -2h, -30d (subtract)

Examples:
  openGyver timex add "2024-01-15T14:30:00Z" 2h30m
  openGyver timex add "2024-01-15" 90d
  openGyver timex add "2024-01-15" -30d
  openGyver timex add now 2w
  openGyver timex add "2024-03-01" 1y2mo
  openGyver timex add "2024-01-15" 1y --tz America/New_York
  openGyver timex add now -1h30m`,
	Args: cobra.ExactArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		var assumeTZ *time.Location
		if addFrom != "" {
			loc, err := loadLocation(addFrom)
			if err != nil {
				return err
			}
			assumeTZ = loc
		}

		t, err := parseTime(args[0], assumeTZ)
		if err != nil {
			return err
		}

		result, err := addDuration(t, args[1])
		if err != nil {
			return err
		}

		if addTZ != "" {
			loc, err := loadLocation(addTZ)
			if err != nil {
				return err
			}
			result = result.In(loc)
		}

		if brief {
			fmt.Println(result.Format("2006-01-02T15:04:05Z07:00"))
			return nil
		}
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"original": t.Format("2006-01-02T15:04:05Z07:00"),
				"duration": args[1],
				"result":   result.Format("2006-01-02T15:04:05Z07:00"),
				"unix":     result.Unix(),
			})
		}

		fmt.Printf("Original:  %s\n", t.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Printf("Duration:  %s\n", args[1])
		fmt.Printf("Result:    %s\n", result.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Printf("Unix:      %d\n", result.Unix())
		return nil
	},
}

// extDurationPattern matches extended duration components: 1y, 2mo, 30d, 2w
var extDurationPattern = regexp.MustCompile(`(-?)(?:(\d+)y)?(?:(\d+)mo)?(?:(\d+)w)?(?:(\d+)d)?(.*)`)

func addDuration(t time.Time, durStr string) (time.Time, error) {
	durStr = strings.TrimSpace(durStr)

	// Check for negative prefix
	negative := false
	cleanDur := durStr
	if strings.HasPrefix(cleanDur, "-") {
		negative = true
		cleanDur = cleanDur[1:]
	}

	// Try Go's time.ParseDuration first (handles h, m, s, ms, us, ns)
	if d, err := time.ParseDuration(durStr); err == nil {
		return t.Add(d), nil
	}

	// Parse extended components: y, mo, w, d + remaining Go duration
	m := extDurationPattern.FindStringSubmatch(cleanDur)
	if m == nil {
		return time.Time{}, fmt.Errorf("invalid duration: %q\nExamples: 2h30m, 30d, 2w, 1y2mo, -1h", durStr)
	}

	sign := 1
	if negative || m[1] == "-" {
		sign = -1
	}

	years := parseIntOr(m[2], 0) * sign
	months := parseIntOr(m[3], 0) * sign
	weeks := parseIntOr(m[4], 0) * sign
	days := parseIntOr(m[5], 0) * sign
	remainder := m[6]

	result := t.AddDate(years, months, days+weeks*7)

	// Parse remaining Go-style duration (h, m, s)
	if remainder != "" {
		if negative && !strings.HasPrefix(remainder, "-") {
			remainder = "-" + remainder
		}
		d, err := time.ParseDuration(remainder)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration component %q: %w", remainder, err)
		}
		result = result.Add(d)
	}

	return result, nil
}

func parseIntOr(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func init() {
	addCmd.Flags().StringVar(&addFrom, "from", "", "source timezone for naive time inputs")
	addCmd.Flags().StringVar(&addTZ, "tz", "", "display result in this timezone")
	register(addCmd)
}

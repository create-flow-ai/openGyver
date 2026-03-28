package timex

import (
	"fmt"
	"math"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var diffFrom string

var diffCmd = &cobra.Command{
	Use:   "diff <time1> <time2>",
	Short: "Calculate the duration between two times",
	Long: `Calculate and display the elapsed time between two date/time values.

Shows the difference in multiple units: total days, hours, minutes, seconds,
as well as a human-friendly breakdown. The result is always shown as a
positive duration regardless of order.

Use --from to specify the timezone for naive inputs.

Examples:
  openGyver timex diff "2024-01-15" "2024-06-30"
  openGyver timex diff "2024-01-15T08:00:00Z" "2024-01-15T17:30:00Z"
  openGyver timex diff "2024-01-01" now
  openGyver timex diff yesterday tomorrow
  openGyver timex diff 1705334400 1710000000`,
	Args: cobra.ExactArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		var assumeTZ *time.Location
		if diffFrom != "" {
			loc, err := loadLocation(diffFrom)
			if err != nil {
				return err
			}
			assumeTZ = loc
		}

		t1, err := parseTime(args[0], assumeTZ)
		if err != nil {
			return fmt.Errorf("first time: %w", err)
		}
		t2, err := parseTime(args[1], assumeTZ)
		if err != nil {
			return fmt.Errorf("second time: %w", err)
		}

		d := t2.Sub(t1)
		if d < 0 {
			d = -d
		}

		if brief {
			fmt.Printf("%.0f\n", d.Seconds())
			return nil
		}

		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"from":          t1.Format("2006-01-02T15:04:05Z07:00"),
				"to":            t2.Format("2006-01-02T15:04:05Z07:00"),
				"total_seconds": d.Seconds(),
				"total_hours":   d.Hours(),
				"total_days":    d.Hours() / 24,
			})
		}

		// Human breakdown
		totalSec := int64(d.Seconds())
		days := totalSec / 86400
		hours := (totalSec % 86400) / 3600
		mins := (totalSec % 3600) / 60
		secs := totalSec % 60

		fmt.Printf("From:      %s\n", t1.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Printf("To:        %s\n", t2.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Println()
		fmt.Printf("Duration:  %dd %dh %dm %ds\n", days, hours, mins, secs)
		fmt.Println()
		fmt.Printf("Total days:    %.2f\n", d.Hours()/24)
		fmt.Printf("Total hours:   %.2f\n", d.Hours())
		fmt.Printf("Total minutes: %.0f\n", d.Minutes())
		fmt.Printf("Total seconds: %.0f\n", d.Seconds())

		// Also show weeks/months/years for large spans
		if days >= 7 {
			fmt.Printf("Total weeks:   %.2f\n", float64(days)/7)
		}
		if days >= 30 {
			fmt.Printf("Total months:  %.2f\n", float64(days)/30.4375)
		}
		if days >= 365 {
			fmt.Printf("Total years:   %.2f\n", float64(days)/365.25)
		}

		// Calendar difference (years, months, days)
		if days >= 30 {
			y1, m1, d1 := t1.Date()
			y2, m2, d2 := t2.Date()
			if t1.After(t2) {
				y1, m1, d1, y2, m2, d2 = y2, m2, d2, y1, m1, d1
			}
			years := y2 - y1
			months := int(m2) - int(m1)
			ddays := d2 - d1
			if ddays < 0 {
				months--
				// Get days in the previous month
				prev := time.Date(y2, m2, 0, 0, 0, 0, 0, time.UTC)
				ddays += prev.Day()
			}
			if months < 0 {
				years--
				months += 12
			}
			fmt.Printf("\nCalendar:  %dy %dm %dd\n", years, int(math.Abs(float64(months))), int(math.Abs(float64(ddays))))
		}

		return nil
	},
}

func init() {
	diffCmd.Flags().StringVar(&diffFrom, "from", "", "source timezone for naive inputs")
	register(diffCmd)
}

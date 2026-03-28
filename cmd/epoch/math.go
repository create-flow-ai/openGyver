package epoch

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	from    int64
	hours   int
	minutes int
	days    int
	weeks   int
	months  int
	years   int
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a duration to an epoch and return the new epoch",
	Long: `Add hours, days, weeks, months, or years to an epoch timestamp.

Uses the current epoch by default. Use --from to specify a starting epoch.
Multiple duration flags can be combined.

Examples:
  openGyver epoch add --hours 2
  openGyver epoch add --days 30
  openGyver epoch add --days 7 --hours 12
  openGyver epoch add --months 3 --from 1705334400
  openGyver epoch add --years 1 --months 6
  openGyver epoch add --weeks 2 --ms`,
	RunE: func(c *cobra.Command, args []string) error {
		base := resolveBase(from)
		result := base.AddDate(years, months, days+weeks*7).
			Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)
		printEpoch(result)
		return nil
	},
}

var subtractCmd = &cobra.Command{
	Use:   "subtract",
	Short: "Subtract a duration from an epoch and return the new epoch",
	Long: `Subtract hours, days, weeks, months, or years from an epoch timestamp.

Uses the current epoch by default. Use --from to specify a starting epoch.
Multiple duration flags can be combined.

Examples:
  openGyver epoch subtract --hours 2
  openGyver epoch subtract --days 30
  openGyver epoch subtract --days 7 --hours 12
  openGyver epoch subtract --months 3 --from 1705334400
  openGyver epoch subtract --years 1
  openGyver epoch subtract --weeks 2 --ms`,
	RunE: func(c *cobra.Command, args []string) error {
		base := resolveBase(from)
		result := base.AddDate(-years, -months, -(days + weeks*7)).
			Add(-time.Duration(hours)*time.Hour - time.Duration(minutes)*time.Minute)
		printEpoch(result)
		return nil
	},
}

func resolveBase(fromEpoch int64) time.Time {
	if fromEpoch != 0 {
		return time.Unix(fromEpoch, 0).UTC()
	}
	return time.Now()
}

func addDurationFlags(c *cobra.Command) {
	c.Flags().Int64Var(&from, "from", 0, "starting epoch in seconds (default: current time)")
	c.Flags().IntVar(&hours, "hours", 0, "hours to add/subtract")
	c.Flags().IntVar(&minutes, "minutes", 0, "minutes to add/subtract")
	c.Flags().IntVar(&days, "days", 0, "days to add/subtract")
	c.Flags().IntVar(&weeks, "weeks", 0, "weeks to add/subtract")
	c.Flags().IntVar(&months, "months", 0, "months to add/subtract")
	c.Flags().IntVar(&years, "years", 0, "years to add/subtract")
}

func init() {
	addDurationFlags(addCmd)
	addDurationFlags(subtractCmd)
}

// computeAdd is exported for testing — applies add logic to a fixed base.
func computeAdd(base time.Time, y, mo, w, d, h, min int) time.Time {
	return base.AddDate(y, mo, d+w*7).
		Add(time.Duration(h)*time.Hour + time.Duration(min)*time.Minute)
}

// computeSubtract is exported for testing — applies subtract logic to a fixed base.
func computeSubtract(base time.Time, y, mo, w, d, h, min int) time.Time {
	return base.AddDate(-y, -mo, -(d+w*7)).
		Add(-time.Duration(h)*time.Hour - time.Duration(min)*time.Minute)
}

// formatEpoch returns the epoch as a string in the current precision mode (for testing).
func formatEpoch(t time.Time) string {
	switch {
	case ns:
		return fmt.Sprintf("%d", t.UnixNano())
	case us:
		return fmt.Sprintf("%d", t.UnixMicro())
	case ms:
		return fmt.Sprintf("%d", t.UnixMilli())
	default:
		return fmt.Sprintf("%d", t.Unix())
	}
}

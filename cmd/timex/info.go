package timex

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var infoFrom string

var infoCmd = &cobra.Command{
	Use:   "info <time>",
	Short: "Show detailed metadata about a date/time",
	Long: `Parse a time and display comprehensive metadata about it.

Shows: day of week, day of year, ISO week number, quarter, leap year status,
Unix timestamps, timezone offset, and more.

Examples:
  openGyver timex info "2024-01-15"
  openGyver timex info "2024-01-15T14:30:00Z"
  openGyver timex info now
  openGyver timex info 1705334400
  openGyver timex info "2024-02-29" --from America/New_York`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		var assumeTZ *time.Location
		if infoFrom != "" {
			loc, err := loadLocation(infoFrom)
			if err != nil {
				return err
			}
			assumeTZ = loc
		}

		t, err := parseTime(args[0], assumeTZ)
		if err != nil {
			return err
		}

		year, month, day := t.Date()
		hour, min, sec := t.Clock()
		isoYear, isoWeek := t.ISOWeek()
		dayOfYear := t.YearDay()
		quarter := (int(month) + 2) / 3
		isLeap := year%4 == 0 && (year%100 != 0 || year%400 == 0)
		zoneName, zoneOffset := t.Zone()

		daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, t.Location()).Day()
		daysInYear := 365
		if isLeap {
			daysInYear = 366
		}
		daysRemaining := daysInYear - dayOfYear

		// Start of day / end of day
		startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
		endOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())

		if brief {
			fmt.Println(t.Format("2006-01-02T15:04:05Z07:00"))
			return nil
		}

		fmt.Printf("Input:          %s\n", args[0])
		fmt.Printf("Parsed:         %s\n", t.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Println()

		fmt.Printf("--- Date ---\n")
		fmt.Printf("Year:           %d\n", year)
		fmt.Printf("Month:          %s (%d)\n", month, month)
		fmt.Printf("Day:            %d\n", day)
		fmt.Printf("Day of week:    %s\n", t.Weekday())
		fmt.Printf("Day of year:    %d / %d\n", dayOfYear, daysInYear)
		fmt.Printf("Days remaining: %d\n", daysRemaining)
		fmt.Printf("Days in month:  %d\n", daysInMonth)
		fmt.Printf("ISO week:       %d-W%02d\n", isoYear, isoWeek)
		fmt.Printf("Quarter:        Q%d\n", quarter)
		fmt.Printf("Leap year:      %t\n", isLeap)
		fmt.Println()

		fmt.Printf("--- Time ---\n")
		fmt.Printf("Hour:           %d\n", hour)
		fmt.Printf("Minute:         %d\n", min)
		fmt.Printf("Second:         %d\n", sec)
		fmt.Printf("Nanosecond:     %d\n", t.Nanosecond())
		fmt.Printf("12-hour:        %s\n", t.Format("3:04:05 PM"))
		fmt.Println()

		fmt.Printf("--- Timezone ---\n")
		fmt.Printf("Timezone:       %s\n", t.Location())
		fmt.Printf("Zone name:      %s\n", zoneName)
		fmt.Printf("UTC offset:     %+d seconds (%+.1f hours)\n", zoneOffset, float64(zoneOffset)/3600)
		fmt.Println()

		fmt.Printf("--- Epoch ---\n")
		fmt.Printf("Unix (s):       %d\n", t.Unix())
		fmt.Printf("Unix (ms):      %d\n", t.UnixMilli())
		fmt.Printf("Unix (us):      %d\n", t.UnixMicro())
		fmt.Printf("Unix (ns):      %d\n", t.UnixNano())
		fmt.Println()

		fmt.Printf("--- Boundaries ---\n")
		fmt.Printf("Start of day:   %s\n", startOfDay.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Printf("End of day:     %s\n", endOfDay.Format("2006-01-02T15:04:05.999999999Z07:00"))

		return nil
	},
}

func init() {
	infoCmd.Flags().StringVar(&infoFrom, "from", "", "timezone for naive inputs")
	register(infoCmd)
}

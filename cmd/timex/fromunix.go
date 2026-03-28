package timex

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	fromUnixTZ    string
	fromUnixMS    bool
	fromUnixUS    bool
	fromUnixNS    bool
	fromUnixFmt   string
)

var fromUnixCmd = &cobra.Command{
	Use:   "from-unix <timestamp>",
	Short: "Convert a Unix epoch timestamp to human-readable time",
	Long: `Convert a Unix epoch number to a human-readable date/time.

By default, the input is treated as seconds. Use --ms, --us, or --ns to
specify other precisions. Use --tz to display in a specific timezone.

Auto-detection: if no precision flag is given, large numbers are automatically
detected as milliseconds (>1e12), microseconds (>1e15), or nanoseconds (>1e18).

Examples:
  openGyver timex from-unix 1705334400
  openGyver timex from-unix 1705334400000 --ms
  openGyver timex from-unix 1705334400 --tz Asia/Tokyo
  openGyver timex from-unix 1705334400 --format kitchen
  openGyver timex from-unix 1705334400000000000 --ns`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		var t time.Time

		explicitPrecision := fromUnixMS || fromUnixUS || fromUnixNS

		if explicitPrecision {
			n, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid timestamp: %s", args[0])
			}
			switch {
			case fromUnixNS:
				t = time.Unix(0, n)
			case fromUnixUS:
				t = time.Unix(0, n*1000)
			case fromUnixMS:
				t = time.Unix(0, n*1000000)
			}
		} else {
			// Auto-detect precision
			var err error
			t, err = parseUnixNumeric(args[0])
			if err != nil {
				return fmt.Errorf("invalid timestamp: %s", args[0])
			}
		}

		loc := time.UTC
		if fromUnixTZ != "" {
			var err error
			loc, err = loadLocation(fromUnixTZ)
			if err != nil {
				return err
			}
		}
		t = t.In(loc)

		if fromUnixFmt != "" {
			layout := resolveFormat(fromUnixFmt)
			fmt.Println(t.Format(layout))
			return nil
		}

		fmt.Printf("Timezone:  %s\n", t.Location())
		fmt.Printf("ISO 8601:  %s\n", t.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Printf("RFC 2822:  %s\n", t.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
		fmt.Printf("Date:      %s\n", t.Format("2006-01-02"))
		fmt.Printf("Time:      %s\n", t.Format("15:04:05"))
		fmt.Printf("12-hour:   %s\n", t.Format("3:04:05 PM"))
		fmt.Printf("Human:     %s\n", t.Format("Mon, Jan 2 2006 at 3:04 PM MST"))
		return nil
	},
}

func init() {
	fromUnixCmd.Flags().StringVar(&fromUnixTZ, "tz", "", "display timezone (default: UTC)")
	fromUnixCmd.Flags().BoolVar(&fromUnixMS, "ms", false, "input is in milliseconds")
	fromUnixCmd.Flags().BoolVar(&fromUnixUS, "us", false, "input is in microseconds")
	fromUnixCmd.Flags().BoolVar(&fromUnixNS, "ns", false, "input is in nanoseconds")
	fromUnixCmd.Flags().StringVar(&fromUnixFmt, "format", "", "output format (iso8601, rfc2822, date, kitchen, human, etc.)")
	register(fromUnixCmd)
}

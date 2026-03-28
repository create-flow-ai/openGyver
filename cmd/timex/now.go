package timex

import (
	"fmt"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var nowTZ string
var nowFormat string

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Show current time in various formats",
	Long: `Display the current time in multiple standard formats and timezones.

By default shows: UTC, local, Unix epoch, ISO 8601, RFC 2822, and more.
Use --tz to show the time in a specific timezone.
Use --format to show only one specific format.

Examples:
  openGyver timex now
  openGyver timex now --tz Asia/Tokyo
  openGyver timex now --tz EST
  openGyver timex now --format iso8601
  openGyver timex now --tz Europe/London --format rfc2822`,
	Args: cobra.NoArgs,
	RunE: func(c *cobra.Command, args []string) error {
		now := time.Now()

		loc := now.Location()
		if nowTZ != "" {
			var err error
			loc, err = loadLocation(nowTZ)
			if err != nil {
				return err
			}
		}
		t := now.In(loc)

		if nowFormat != "" {
			layout := resolveFormat(nowFormat)
			fmt.Println(t.Format(layout))
			return nil
		}

		if brief {
			fmt.Println(t.Format("2006-01-02T15:04:05Z07:00"))
			return nil
		}

		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"timezone": t.Location().String(),
				"iso8601":  t.Format("2006-01-02T15:04:05Z07:00"),
				"rfc2822":  t.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
				"date":     t.Format("2006-01-02"),
				"time":     t.Format("15:04:05"),
				"unix":     t.Unix(),
				"unix_ms":  t.UnixMilli(),
			})
		}

		fmt.Printf("Timezone:  %s\n", t.Location())
		fmt.Printf("ISO 8601:  %s\n", t.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Printf("RFC 2822:  %s\n", t.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
		fmt.Printf("RFC 3339:  %s\n", t.Format(time.RFC3339Nano))
		fmt.Printf("Date:      %s\n", t.Format("2006-01-02"))
		fmt.Printf("Time:      %s\n", t.Format("15:04:05"))
		fmt.Printf("12-hour:   %s\n", t.Format("3:04:05 PM"))
		fmt.Printf("Human:     %s\n", t.Format("Mon, Jan 2 2006 at 3:04 PM MST"))
		fmt.Printf("Unix:      %d\n", t.Unix())
		fmt.Printf("Unix (ms): %d\n", t.UnixMilli())
		fmt.Printf("Unix (ns): %d\n", t.UnixNano())

		if nowTZ == "" {
			fmt.Printf("\nUTC:       %s\n", now.UTC().Format("2006-01-02T15:04:05Z"))
		}

		return nil
	},
}

func init() {
	nowCmd.Flags().StringVar(&nowTZ, "tz", "", "timezone to display (IANA name or abbreviation)")
	nowCmd.Flags().StringVar(&nowFormat, "format", "", "output format name (iso8601, rfc2822, rfc3339, date, time, kitchen, human, etc.)")
	register(nowCmd)
}

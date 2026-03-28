package timex

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var toUTCFrom string

var toUTCCmd = &cobra.Command{
	Use:   "to-utc <time>",
	Short: "Convert a time string to UTC",
	Long: `Parse a time string and convert it to UTC.

If the input contains timezone info (offset or zone name), it is used directly.
Otherwise, use --from to specify the source timezone (defaults to local).

Examples:
  openGyver timex to-utc "2024-01-15T14:30:00-05:00"
  openGyver timex to-utc "2024-01-15 14:30" --from America/New_York
  openGyver timex to-utc "Jan 15, 2024 2:30 PM" --from PST
  openGyver timex to-utc now
  openGyver timex to-utc 1705334400`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		var assumeTZ *time.Location
		if toUTCFrom != "" {
			loc, err := loadLocation(toUTCFrom)
			if err != nil {
				return err
			}
			assumeTZ = loc
		} else {
			assumeTZ = time.Local
		}

		t, err := parseTime(args[0], assumeTZ)
		if err != nil {
			return err
		}

		utc := t.UTC()
		fmt.Printf("UTC:       %s\n", utc.Format("2006-01-02T15:04:05Z"))
		fmt.Printf("RFC 2822:  %s\n", utc.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
		fmt.Printf("Unix:      %d\n", utc.Unix())
		return nil
	},
}

func init() {
	toUTCCmd.Flags().StringVar(&toUTCFrom, "from", "", "source timezone for naive inputs (default: local)")
	register(toUTCCmd)
}

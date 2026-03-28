package timex

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	toUnixFrom string
	toUnixMS   bool
	toUnixUS   bool
	toUnixNS   bool
)

var toUnixCmd = &cobra.Command{
	Use:   "to-unix <time>",
	Short: "Convert a time string to Unix epoch",
	Long: `Parse a time string and output the Unix epoch timestamp.

By default outputs seconds. Use flags for other precisions.

Examples:
  openGyver timex to-unix "2024-01-15T14:30:00Z"
  openGyver timex to-unix "Jan 15, 2024 2:30 PM" --from America/New_York
  openGyver timex to-unix "2024-01-15" --ms
  openGyver timex to-unix now --ns
  openGyver timex to-unix "2024-01-15 14:30" --from EST`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		var assumeTZ *time.Location
		if toUnixFrom != "" {
			loc, err := loadLocation(toUnixFrom)
			if err != nil {
				return err
			}
			assumeTZ = loc
		}

		t, err := parseTime(args[0], assumeTZ)
		if err != nil {
			return err
		}

		switch {
		case toUnixNS:
			fmt.Println(t.UnixNano())
		case toUnixUS:
			fmt.Println(t.UnixMicro())
		case toUnixMS:
			fmt.Println(t.UnixMilli())
		default:
			fmt.Println(t.Unix())
		}
		return nil
	},
}

func init() {
	toUnixCmd.Flags().StringVar(&toUnixFrom, "from", "", "source timezone for naive inputs (default: UTC)")
	toUnixCmd.Flags().BoolVar(&toUnixMS, "ms", false, "output in milliseconds")
	toUnixCmd.Flags().BoolVar(&toUnixUS, "us", false, "output in microseconds")
	toUnixCmd.Flags().BoolVar(&toUnixNS, "ns", false, "output in nanoseconds")
	register(toUnixCmd)
}

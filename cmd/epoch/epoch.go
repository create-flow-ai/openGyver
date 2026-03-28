package epoch

import (
	"fmt"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	ms bool
	us bool
	ns bool
)

var epochCmd = &cobra.Command{
	Use:   "epoch",
	Short: "Unix epoch utilities — current time, add, subtract",
	Long: `Print the current Unix epoch timestamp, or perform arithmetic on epochs.

SUBCOMMANDS:

  epoch            Print the current epoch (default)
  epoch add        Add a duration to an epoch
  epoch subtract   Subtract a duration from an epoch

By default outputs seconds. Use --ms, --us, or --ns for other precisions.

Examples:
  openGyver epoch
  openGyver epoch --ms
  openGyver epoch add --hours 2
  openGyver epoch add --days 30 --from 1705334400
  openGyver epoch subtract --years 1`,
	Args: cobra.NoArgs,
	RunE: func(c *cobra.Command, args []string) error {
		now := time.Now()
		printEpoch(now)
		return nil
	},
}

func printEpoch(t time.Time) {
	switch {
	case ns:
		fmt.Println(t.UnixNano())
	case us:
		fmt.Println(t.UnixMicro())
	case ms:
		fmt.Println(t.UnixMilli())
	default:
		fmt.Println(t.Unix())
	}
}

func init() {
	epochCmd.PersistentFlags().BoolVar(&ms, "ms", false, "output in milliseconds")
	epochCmd.PersistentFlags().BoolVar(&us, "us", false, "output in microseconds")
	epochCmd.PersistentFlags().BoolVar(&ns, "ns", false, "output in nanoseconds")

	epochCmd.AddCommand(addCmd)
	epochCmd.AddCommand(subtractCmd)

	cmd.Register(epochCmd)
}

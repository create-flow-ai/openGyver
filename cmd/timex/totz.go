package timex

import (
	"fmt"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	toTZTarget string
	toTZFrom   string
)

var toTZCmd = &cobra.Command{
	Use:   "to-tz <time>",
	Short: "Convert a time to a target timezone",
	Long: `Parse a time string and convert it to the specified timezone.

Use --tz to set the target timezone (required).
Use --from to specify the source timezone for inputs without timezone info.

Examples:
  openGyver timex to-tz "2024-01-15T14:30:00Z" --tz Asia/Tokyo
  openGyver timex to-tz "2024-01-15 09:00" --from America/New_York --tz Europe/London
  openGyver timex to-tz now --tz Australia/Sydney
  openGyver timex to-tz 1705334400 --tz America/Chicago`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		if toTZTarget == "" {
			return fmt.Errorf("--tz flag is required")
		}

		targetLoc, err := loadLocation(toTZTarget)
		if err != nil {
			return err
		}

		var assumeTZ *time.Location
		if toTZFrom != "" {
			assumeTZ, err = loadLocation(toTZFrom)
			if err != nil {
				return err
			}
		}

		t, err := parseTime(args[0], assumeTZ)
		if err != nil {
			return err
		}

		converted := t.In(targetLoc)
		if brief {
			fmt.Println(converted.Format("2006-01-02T15:04:05Z07:00"))
			return nil
		}
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"timezone": converted.Location().String(),
				"iso8601":  converted.Format("2006-01-02T15:04:05Z07:00"),
				"date":     converted.Format("2006-01-02"),
				"time":     converted.Format("15:04:05"),
			})
		}
		fmt.Printf("Timezone:  %s\n", converted.Location())
		fmt.Printf("ISO 8601:  %s\n", converted.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Printf("RFC 2822:  %s\n", converted.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
		fmt.Printf("Date:      %s\n", converted.Format("2006-01-02"))
		fmt.Printf("Time:      %s\n", converted.Format("15:04:05"))
		fmt.Printf("12-hour:   %s\n", converted.Format("3:04:05 PM"))
		return nil
	},
}

func init() {
	toTZCmd.Flags().StringVar(&toTZTarget, "tz", "", "target timezone (required)")
	toTZCmd.Flags().StringVar(&toTZFrom, "from", "", "source timezone for naive inputs (default: UTC)")
	register(toTZCmd)
}

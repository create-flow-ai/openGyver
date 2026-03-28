package timex

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	fmtTo   string
	fmtFrom string
)

var formatCmd = &cobra.Command{
	Use:   "format <time>",
	Short: "Reformat a time string into a different layout",
	Long: `Parse a time string and reformat it using a named or custom layout.

NAMED FORMATS (use with --to):

  iso8601    2006-01-02T15:04:05Z07:00
  rfc3339    2006-01-02T15:04:05Z07:00
  rfc2822    Mon, 02 Jan 2006 15:04:05 -0700
  rfc1123    Mon, 02 Jan 2006 15:04:05 -0700
  rfc850     Monday, 02-Jan-06 15:04:05 MST
  rfc822     02 Jan 06 15:04 -0700
  ansic      Mon Jan _2 15:04:05 2006
  unix       Mon Jan _2 15:04:05 MST 2006
  ruby       Mon Jan 02 15:04:05 -0700 2006
  date       2006-01-02
  time       15:04:05
  datetime   2006-01-02 15:04:05
  kitchen    3:04PM
  us         01/02/2006 3:04:05 PM
  eu         02/01/2006 15:04:05
  short      Jan 2, 2006
  long       January 2, 2006 15:04:05 MST
  stamp      Jan _2 15:04:05
  human      Mon, Jan 2 2006 at 3:04 PM MST

Or pass a custom Go time layout string as --to value.

Examples:
  openGyver timex format "2024-01-15T14:30:00Z" --to rfc2822
  openGyver timex format "Mon, 15 Jan 2024 14:30:00 +0000" --to iso8601
  openGyver timex format "2024-01-15" --to human --from America/New_York
  openGyver timex format "2024-01-15T14:30:00Z" --to kitchen
  openGyver timex format "2024-01-15T14:30:00Z" --to "Monday, January 2 2006"
  openGyver timex format 1705334400 --to short`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		if fmtTo == "" {
			// Show all named formats
			return showAllFormats(args[0])
		}

		var assumeTZ *time.Location
		if fmtFrom != "" {
			loc, err := loadLocation(fmtFrom)
			if err != nil {
				return err
			}
			assumeTZ = loc
		}

		t, err := parseTime(args[0], assumeTZ)
		if err != nil {
			return err
		}

		layout := resolveFormat(fmtTo)
		fmt.Println(t.Format(layout))
		return nil
	},
}

func showAllFormats(input string) error {
	t, err := parseTime(input, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Input:     %s\n\n", input)
	order := []string{
		"iso8601", "rfc3339", "rfc2822", "rfc1123", "rfc850", "rfc822",
		"ansic", "unix", "ruby", "date", "time", "datetime",
		"kitchen", "us", "eu", "short", "long", "stamp", "human",
	}
	for _, name := range order {
		layout := namedFormats[name]
		fmt.Printf("%-10s %s\n", name, t.Format(layout))
	}
	fmt.Printf("\n%-10s %d\n", "epoch", t.Unix())
	fmt.Printf("%-10s %d\n", "epoch_ms", t.UnixMilli())
	return nil
}

func init() {
	formatCmd.Flags().StringVar(&fmtTo, "to", "", "target format name or Go layout (omit to show all formats)")
	formatCmd.Flags().StringVar(&fmtFrom, "from", "", "source timezone for naive inputs")

	// Add format list to help
	var fmtList strings.Builder
	fmtList.WriteString("available format names: ")
	i := 0
	for name := range namedFormats {
		if i > 0 {
			fmtList.WriteString(", ")
		}
		fmtList.WriteString(name)
		i++
	}
	_ = fmtList.String() // keep for reference

	register(formatCmd)
}

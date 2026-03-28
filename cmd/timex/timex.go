package timex

import (
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var timexCmd = &cobra.Command{
	Use:   "timex",
	Short: "Time conversion and timezone utilities",
	Long: `Convert, format, and manipulate dates, times, timezones, and Unix epochs.

SUBCOMMANDS:

  now        Show current time in various formats and timezones
  to-utc     Convert a time string to UTC
  to-tz      Convert a time to a target timezone
  to-unix    Convert a time string to a Unix epoch timestamp
  from-unix  Convert a Unix epoch timestamp to human-readable time
  format     Reformat a time string into a different layout
  diff       Calculate the duration between two times
  add        Add or subtract a duration from a time
  info       Show detailed metadata about a date/time

INPUT FORMATS (auto-detected):

  ISO 8601 / RFC 3339    2024-01-15T14:30:00Z, 2024-01-15T14:30:00+05:30
  RFC 2822               Mon, 15 Jan 2024 14:30:00 +0000
  RFC 850                Monday, 15-Jan-24 14:30:00 UTC
  Date only              2024-01-15, 01/15/2024, 15-Jan-2024, Jan 15 2024
  Date + time            2024-01-15 14:30:00, 2024-01-15 14:30
  12-hour                2024-01-15 2:30 PM, Jan 15, 2024 2:30:00 PM
  Unix timestamp         1705334400 (auto-detected when input is numeric)
  Relative               now, today, yesterday, tomorrow

TIMEZONE FORMAT:

  Use IANA timezone names: America/New_York, Europe/London, Asia/Tokyo, etc.
  Also accepts: UTC, EST, PST, IST, JST, CET, and other common abbreviations.

EXAMPLES:

  openGyver timex now
  openGyver timex now --tz Asia/Tokyo
  openGyver timex to-utc "2024-01-15 14:30" --from America/New_York
  openGyver timex to-tz "2024-01-15T14:30:00Z" --tz Asia/Tokyo
  openGyver timex to-unix "2024-01-15T14:30:00Z"
  openGyver timex from-unix 1705334400
  openGyver timex format "2024-01-15T14:30:00Z" --to rfc2822
  openGyver timex diff "2024-01-15" "2024-06-30"
  openGyver timex add "2024-01-15T14:30:00Z" 2h30m
  openGyver timex info "2024-01-15"`,
}

func init() {
	cmd.Register(timexCmd)
}

// register adds a subcommand to the timex parent command.
func register(sub *cobra.Command) {
	timexCmd.AddCommand(sub)
}

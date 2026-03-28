# timex

Convert, format, and manipulate dates, times, timezones, and Unix epochs.

## Usage

```bash
openGyver timex [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--brief` | `-b` | bool | `false` | Output a single ISO 8601 line (for piping) |
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for timex |

## Input Formats (auto-detected)

| Format | Example |
|--------|---------|
| ISO 8601 / RFC 3339 | `2024-01-15T14:30:00Z`, `2024-01-15T14:30:00+05:30` |
| RFC 2822 | `Mon, 15 Jan 2024 14:30:00 +0000` |
| RFC 850 | `Monday, 15-Jan-24 14:30:00 UTC` |
| Date only | `2024-01-15`, `01/15/2024`, `15-Jan-2024`, `Jan 15 2024` |
| Date + time | `2024-01-15 14:30:00`, `2024-01-15 14:30` |
| 12-hour | `2024-01-15 2:30 PM`, `Jan 15, 2024 2:30:00 PM` |
| Unix timestamp | `1705334400` (auto-detected when input is numeric) |
| Relative | `now`, `today`, `yesterday`, `tomorrow` |

## Timezone Format

Use IANA timezone names: `America/New_York`, `Europe/London`, `Asia/Tokyo`, etc.
Also accepts: `UTC`, `EST`, `PST`, `IST`, `JST`, `CET`, and other common abbreviations.

## Subcommands

### now

Display the current time in multiple standard formats and timezones. By default shows UTC, local, Unix epoch, ISO 8601, RFC 2822, and more.

#### Arguments

No positional arguments.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--tz` | | string | | Timezone to display (IANA name or abbreviation) |
| `--format` | | string | | Output format name (iso8601, rfc2822, rfc3339, date, time, kitchen, human, etc.) |
| `--help` | `-h` | | | Help for now |

#### Examples

```bash
# Show current time in all formats
openGyver timex now

# Show current time in Tokyo
openGyver timex now --tz Asia/Tokyo

# Show current time in EST
openGyver timex now --tz EST

# Show only ISO 8601 format
openGyver timex now --format iso8601

# Show London time in RFC 2822 format
openGyver timex now --tz Europe/London --format rfc2822

# Brief output for piping
openGyver timex now --brief

# JSON output for scripting
openGyver timex now --json

# Kitchen format (e.g., 3:04PM)
openGyver timex now --format kitchen
```

---

### to-utc

Parse a time string and convert it to UTC. If the input contains timezone info (offset or zone name), it is used directly. Otherwise, use `--from` to specify the source timezone (defaults to local).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `time` | Yes | The time string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | string | local | Source timezone for naive inputs |
| `--help` | `-h` | | | Help for to-utc |

#### Examples

```bash
# Convert a time with embedded timezone
openGyver timex to-utc "2024-01-15T14:30:00-05:00"

# Specify source timezone for a naive time
openGyver timex to-utc "2024-01-15 14:30" --from America/New_York

# Convert with named month format
openGyver timex to-utc "Jan 15, 2024 2:30 PM" --from PST

# Convert "now"
openGyver timex to-utc now

# Convert a Unix timestamp
openGyver timex to-utc 1705334400

# Brief output for piping
openGyver timex to-utc "2024-01-15 14:30" --from EST --brief

# JSON output
openGyver timex to-utc "2024-01-15T14:30:00Z" --json

# Convert a date without time
openGyver timex to-utc "2024-01-15" --from America/Chicago
```

---

### to-tz

Parse a time string and convert it to the specified timezone. Use `--tz` to set the target timezone (required). Use `--from` to specify the source timezone for inputs without timezone info.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `time` | Yes | The time string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--tz` | | string | | Target timezone (required) |
| `--from` | | string | UTC | Source timezone for naive inputs |
| `--help` | `-h` | | | Help for to-tz |

#### Examples

```bash
# Convert UTC to Tokyo time
openGyver timex to-tz "2024-01-15T14:30:00Z" --tz Asia/Tokyo

# Convert New York time to London time
openGyver timex to-tz "2024-01-15 09:00" --from America/New_York --tz Europe/London

# Convert current time to Sydney
openGyver timex to-tz now --tz Australia/Sydney

# Convert a Unix timestamp to Chicago time
openGyver timex to-tz 1705334400 --tz America/Chicago

# Brief output for piping
openGyver timex to-tz "2024-01-15T14:30:00Z" --tz Asia/Seoul --brief

# JSON output
openGyver timex to-tz "2024-01-15T14:30:00Z" --tz Europe/Berlin --json

# Convert between non-UTC timezones
openGyver timex to-tz "2024-01-15 09:00" --from America/Los_Angeles --tz Asia/Tokyo

# Convert "today" to a remote timezone
openGyver timex to-tz today --tz Pacific/Auckland
```

---

### to-unix

Parse a time string and output the Unix epoch timestamp. By default outputs seconds. Use flags for other precisions.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `time` | Yes | The time string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | string | UTC | Source timezone for naive inputs |
| `--ms` | | bool | `false` | Output in milliseconds |
| `--us` | | bool | `false` | Output in microseconds |
| `--ns` | | bool | `false` | Output in nanoseconds |
| `--help` | `-h` | | | Help for to-unix |

#### Examples

```bash
# Convert ISO 8601 to Unix seconds
openGyver timex to-unix "2024-01-15T14:30:00Z"

# Specify source timezone
openGyver timex to-unix "Jan 15, 2024 2:30 PM" --from America/New_York

# Output in milliseconds
openGyver timex to-unix "2024-01-15" --ms

# Output in nanoseconds
openGyver timex to-unix now --ns

# Output in microseconds
openGyver timex to-unix "2024-01-15T14:30:00Z" --us

# From EST timezone
openGyver timex to-unix "2024-01-15 14:30" --from EST

# JSON output
openGyver timex to-unix "2024-01-15T14:30:00Z" --json

# Brief output
openGyver timex to-unix "2024-01-15" --brief
```

---

### from-unix

Convert a Unix epoch number to a human-readable date/time. By default, the input is treated as seconds. Use `--ms`, `--us`, or `--ns` to specify other precisions. Auto-detection: if no precision flag is given, large numbers are automatically detected as milliseconds (>1e12), microseconds (>1e15), or nanoseconds (>1e18).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `timestamp` | Yes | The Unix epoch timestamp |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--tz` | | string | UTC | Display timezone |
| `--format` | | string | | Output format (iso8601, rfc2822, date, kitchen, human, etc.) |
| `--ms` | | bool | `false` | Input is in milliseconds |
| `--us` | | bool | `false` | Input is in microseconds |
| `--ns` | | bool | `false` | Input is in nanoseconds |
| `--help` | `-h` | | | Help for from-unix |

#### Examples

```bash
# Convert seconds to human time
openGyver timex from-unix 1705334400

# Convert milliseconds
openGyver timex from-unix 1705334400000 --ms

# Display in a specific timezone
openGyver timex from-unix 1705334400 --tz Asia/Tokyo

# Use kitchen format
openGyver timex from-unix 1705334400 --format kitchen

# Convert nanoseconds
openGyver timex from-unix 1705334400000000000 --ns

# JSON output
openGyver timex from-unix 1705334400 --json

# Brief output
openGyver timex from-unix 1705334400 --brief

# Human-readable format in local timezone
openGyver timex from-unix 1705334400 --format human --tz America/New_York
```

---

### format

Parse a time string and reformat it using a named or custom layout.

**Named formats (use with `--to`):**

| Name | Layout |
|------|--------|
| `iso8601` | `2006-01-02T15:04:05Z07:00` |
| `rfc3339` | `2006-01-02T15:04:05Z07:00` |
| `rfc2822` | `Mon, 02 Jan 2006 15:04:05 -0700` |
| `rfc1123` | `Mon, 02 Jan 2006 15:04:05 -0700` |
| `rfc850` | `Monday, 02-Jan-06 15:04:05 MST` |
| `rfc822` | `02 Jan 06 15:04 -0700` |
| `ansic` | `Mon Jan _2 15:04:05 2006` |
| `unix` | `Mon Jan _2 15:04:05 MST 2006` |
| `ruby` | `Mon Jan 02 15:04:05 -0700 2006` |
| `date` | `2006-01-02` |
| `time` | `15:04:05` |
| `datetime` | `2006-01-02 15:04:05` |
| `kitchen` | `3:04PM` |
| `us` | `01/02/2006 3:04:05 PM` |
| `eu` | `02/01/2006 15:04:05` |
| `short` | `Jan 2, 2006` |
| `long` | `January 2, 2006 15:04:05 MST` |
| `stamp` | `Jan _2 15:04:05` |
| `human` | `Mon, Jan 2 2006 at 3:04 PM MST` |

Or pass a custom Go time layout string as `--to` value.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `time` | Yes | The time string to reformat |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--to` | | string | | Target format name or Go layout (omit to show all formats) |
| `--from` | | string | | Source timezone for naive inputs |
| `--help` | `-h` | | | Help for format |

#### Examples

```bash
# Convert to RFC 2822
openGyver timex format "2024-01-15T14:30:00Z" --to rfc2822

# Convert RFC 2822 to ISO 8601
openGyver timex format "Mon, 15 Jan 2024 14:30:00 +0000" --to iso8601

# Format with source timezone
openGyver timex format "2024-01-15" --to human --from America/New_York

# Kitchen format
openGyver timex format "2024-01-15T14:30:00Z" --to kitchen

# Custom Go layout
openGyver timex format "2024-01-15T14:30:00Z" --to "Monday, January 2 2006"

# Convert Unix timestamp to short format
openGyver timex format 1705334400 --to short

# Show all available formats
openGyver timex format "2024-01-15T14:30:00Z"

# JSON output
openGyver timex format "2024-01-15T14:30:00Z" --to rfc2822 --json
```

---

### diff

Calculate and display the elapsed time between two date/time values. Shows the difference in multiple units: total days, hours, minutes, seconds, as well as a human-friendly breakdown. The result is always shown as a positive duration regardless of order.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `time1` | Yes | First date/time value |
| `time2` | Yes | Second date/time value |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | string | | Source timezone for naive inputs |
| `--help` | `-h` | | | Help for diff |

#### Examples

```bash
# Difference between two dates
openGyver timex diff "2024-01-15" "2024-06-30"

# Difference between precise times
openGyver timex diff "2024-01-15T08:00:00Z" "2024-01-15T17:30:00Z"

# Days since a date
openGyver timex diff "2024-01-01" now

# Difference between relative dates
openGyver timex diff yesterday tomorrow

# Difference between Unix timestamps
openGyver timex diff 1705334400 1710000000

# JSON output
openGyver timex diff "2024-01-15" "2024-06-30" --json

# Brief output
openGyver timex diff "2024-01-01" "2024-12-31" --brief

# With source timezone
openGyver timex diff "2024-01-15 08:00" "2024-01-15 17:30" --from America/New_York
```

---

### add

Parse a time and add (or subtract) a duration, then display the result.

**Duration format:**

| Style | Examples |
|-------|---------|
| Go-style | `1h30m`, `2h`, `45m`, `90s`, `500ms`, `1h30m45s` |
| Extended | `30d`, `2w`, `3mo`, `1y` (days, weeks, months, years) |
| Combined | `1y2mo3d4h5m6s` |
| Negative | `-2h`, `-30d` (subtract) |

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `time` | Yes | The base time string |
| `duration` | Yes | The duration to add (or subtract with `-` prefix) |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | string | | Source timezone for naive time inputs |
| `--tz` | | string | | Display result in this timezone |
| `--help` | `-h` | | | Help for add |

#### Examples

```bash
# Add 2 hours 30 minutes
openGyver timex add "2024-01-15T14:30:00Z" 2h30m

# Add 90 days
openGyver timex add "2024-01-15" 90d

# Subtract 30 days
openGyver timex add "2024-01-15" -30d

# Add 2 weeks from now
openGyver timex add now 2w

# Add 1 year and 2 months
openGyver timex add "2024-03-01" 1y2mo

# Add with timezone display
openGyver timex add "2024-01-15" 1y --tz America/New_York

# Subtract 1.5 hours from now
openGyver timex add now -1h30m

# JSON output
openGyver timex add "2024-01-15T14:30:00Z" 2h30m --json
```

---

### info

Parse a time and display comprehensive metadata about it. Shows: day of week, day of year, ISO week number, quarter, leap year status, Unix timestamps, timezone offset, and more.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `time` | Yes | The date/time to inspect |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | string | | Timezone for naive inputs |
| `--help` | `-h` | | | Help for info |

#### Examples

```bash
# Info about a date
openGyver timex info "2024-01-15"

# Info about a full timestamp
openGyver timex info "2024-01-15T14:30:00Z"

# Info about current time
openGyver timex info now

# Info from a Unix timestamp
openGyver timex info 1705334400

# Info with source timezone
openGyver timex info "2024-02-29" --from America/New_York

# JSON output
openGyver timex info "2024-01-15" --json

# Check if a year is a leap year
openGyver timex info "2024-02-29"

# Brief output
openGyver timex info "2024-07-04" --brief
```

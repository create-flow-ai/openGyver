# epoch

Print the current Unix epoch timestamp, or perform arithmetic on epochs.

## Usage

```bash
openGyver epoch [flags]
openGyver epoch [subcommand] [flags]
```

When invoked without a subcommand, prints the current Unix epoch timestamp.

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON (includes epoch in seconds, milliseconds, nanoseconds, and ISO 8601) |
| `--ms` | | bool | `false` | Output in milliseconds |
| `--us` | | bool | `false` | Output in microseconds |
| `--ns` | | bool | `false` | Output in nanoseconds |
| `--help` | `-h` | bool | | Show help for the command |

## Default Behavior (no subcommand)

When called without a subcommand, `openGyver epoch` prints the current Unix epoch timestamp. By default the output is in seconds. Use `--ms`, `--us`, or `--ns` for other precisions.

### Examples

```bash
# Print current epoch in seconds
openGyver epoch

# Print current epoch in milliseconds
openGyver epoch --ms

# Print current epoch in microseconds
openGyver epoch --us

# Print current epoch in nanoseconds
openGyver epoch --ns

# JSON output (includes all precisions and ISO 8601)
openGyver epoch -j

# Use in a script
TIMESTAMP=$(openGyver epoch)

# Use millisecond epoch for a database timestamp
openGyver epoch --ms

# JSON output piped to jq
openGyver epoch -j | jq '.epoch'
```

### JSON Output Format

```json
{
  "epoch": 1711612800,
  "epoch_ms": 1711612800000,
  "epoch_ns": 1711612800000000000,
  "iso8601": "2024-03-28T12:00:00Z"
}
```

## Subcommands

### add

Add hours, minutes, days, weeks, months, or years to an epoch timestamp. Uses the current epoch by default. Use `--from` to specify a starting epoch in seconds. Multiple duration flags can be combined.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | int | `0` (current time) | Starting epoch in seconds (default: current time) |
| `--hours` | | int | `0` | Hours to add |
| `--minutes` | | int | `0` | Minutes to add |
| `--days` | | int | `0` | Days to add |
| `--weeks` | | int | `0` | Weeks to add |
| `--months` | | int | `0` | Months to add |
| `--years` | | int | `0` | Years to add |

#### Examples

```bash
# Add 2 hours to the current epoch
openGyver epoch add --hours 2

# Add 30 days to the current epoch
openGyver epoch add --days 30

# Combine multiple durations: 7 days and 12 hours
openGyver epoch add --days 7 --hours 12

# Add 3 months to a specific epoch
openGyver epoch add --months 3 --from 1705334400

# Add 1 year and 6 months
openGyver epoch add --years 1 --months 6

# Add 2 weeks, output in milliseconds
openGyver epoch add --weeks 2 --ms

# JSON output
openGyver epoch add --days 90 -j

# Add 45 minutes to a known epoch
openGyver epoch add --minutes 45 --from 1711612800
```

#### JSON Output Format

```json
{
  "operation": "add",
  "base_epoch": 1711612800,
  "result_epoch": 1711620000,
  "result_epoch_ms": 1711620000000,
  "result_iso8601": "2024-03-28T14:00:00Z"
}
```

---

### subtract

Subtract hours, minutes, days, weeks, months, or years from an epoch timestamp. Uses the current epoch by default. Use `--from` to specify a starting epoch in seconds. Multiple duration flags can be combined.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | int | `0` (current time) | Starting epoch in seconds (default: current time) |
| `--hours` | | int | `0` | Hours to subtract |
| `--minutes` | | int | `0` | Minutes to subtract |
| `--days` | | int | `0` | Days to subtract |
| `--weeks` | | int | `0` | Weeks to subtract |
| `--months` | | int | `0` | Months to subtract |
| `--years` | | int | `0` | Years to subtract |

#### Examples

```bash
# Subtract 2 hours from the current epoch
openGyver epoch subtract --hours 2

# Subtract 30 days from the current epoch
openGyver epoch subtract --days 30

# Combine: 7 days and 12 hours ago
openGyver epoch subtract --days 7 --hours 12

# Subtract 3 months from a specific epoch
openGyver epoch subtract --months 3 --from 1705334400

# Subtract 1 year
openGyver epoch subtract --years 1

# Subtract 2 weeks, output in milliseconds
openGyver epoch subtract --weeks 2 --ms

# JSON output for "24 hours ago"
openGyver epoch subtract --hours 24 -j

# Subtract 15 minutes from a known epoch
openGyver epoch subtract --minutes 15 --from 1711612800
```

#### JSON Output Format

```json
{
  "operation": "subtract",
  "base_epoch": 1711612800,
  "result_epoch": 1711605600,
  "result_epoch_ms": 1711605600000,
  "result_iso8601": "2024-03-28T10:00:00Z"
}
```

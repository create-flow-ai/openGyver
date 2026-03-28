# cron

Parse, explain, and validate cron expressions.

## Usage

```bash
openGyver cron [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for cron |

## Subcommands

### explain

Parse a cron expression and output a human-readable description of each field. Supports standard 5-field expressions (minute hour day month weekday) and 6-field expressions with seconds (second minute hour day month weekday).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `expression` | Yes | A 5-field or 6-field cron expression |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Explain a simple cron expression (every 5 minutes)
openGyver cron explain "*/5 * * * *"

# Explain weekday schedule (9am Mon-Fri)
openGyver cron explain "0 9 * * 1-5"

# Explain midnight on the 1st of January
openGyver cron explain "0 0 1 1 *"

# Explain a 6-field expression with seconds
openGyver cron explain "0 */5 * * * *"

# JSON output for scripting
openGyver cron explain "30 2 * * 0" --json

# Explain a complex expression
openGyver cron explain "0,30 8-17 * * 1-5"

# Explain every hour
openGyver cron explain "0 * * * *"
```

#### JSON Output Format

```json
{
  "expression": "*/5 * * * *",
  "has_seconds": false,
  "fields": [
    "every 5 minutes",
    "every hour",
    "every day",
    "every month",
    "every weekday"
  ],
  "summary": "every 5 minutes | every hour | every day | every month | every weekday"
}
```

---

### next

Show the next N run times for a cron expression. Calculates future occurrences starting from the current time. Supports both 5-field and 6-field cron expressions, as well as ranges, steps, and comma-separated lists.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `expression` | Yes | A 5-field or 6-field cron expression |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--count` | | int | `5` | Number of next run times to show |

#### Examples

```bash
# Show next 5 run times (default)
openGyver cron next "*/5 * * * *"

# Show next 10 run times
openGyver cron next "0 9 * * 1-5" --count 10

# Show next 3 runs for a midnight cron
openGyver cron next "0 0 * * *" --count 3

# JSON output for automation
openGyver cron next "0 */2 * * *" --json

# Next runs for a monthly job
openGyver cron next "0 0 1 * *" --count 12

# Next runs with JSON and custom count
openGyver cron next "30 8 * * 1-5" --count 7 --json

# Next runs for a weekend job
openGyver cron next "0 10 * * 0,6" --count 5
```

#### JSON Output Format

```json
{
  "expression": "0 9 * * 1-5",
  "count": 5,
  "next": [
    "2026-03-30T09:00:00Z",
    "2026-03-31T09:00:00Z",
    "2026-04-01T09:00:00Z",
    "2026-04-02T09:00:00Z",
    "2026-04-03T09:00:00Z"
  ]
}
```

---

### validate

Check if a cron expression is syntactically valid. Accepts 5-field and 6-field expressions. Reports whether the expression is valid and, if invalid, provides the error message.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `expression` | Yes | A cron expression to validate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Validate a correct expression
openGyver cron validate "0 0 1 1 *"

# Validate an every-5-minutes expression
openGyver cron validate "*/5 * * * *"

# Validate an invalid expression (too many fields)
openGyver cron validate "* * * * * * *"

# Validate with JSON output
openGyver cron validate "0 9 * * 1-5" --json

# Validate a 6-field expression (with seconds)
openGyver cron validate "0 0 12 * * *"

# Validate an invalid range
openGyver cron validate "0 25 * * *" --json
```

#### JSON Output Format (valid)

```json
{
  "expression": "0 0 1 1 *",
  "valid": true
}
```

#### JSON Output Format (invalid)

```json
{
  "expression": "* * * * * * *",
  "valid": false,
  "error": "expected 5 or 6 fields, got 7"
}
```

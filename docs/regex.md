# regex

Regular expression tools -- test, replace, and extract with Go regexp syntax.

## Usage

```bash
openGyver regex [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for regex |

## Subcommands

### test

Test whether a regular expression matches the given input. Shows whether the pattern matches, all matched substrings, and any captured groups. Use `--global/-g` to find all matches in the input.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `pattern` | Yes | The regular expression pattern (Go regexp syntax) |
| `input` | Yes | The input string to test against |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--global` | `-g` | bool | `false` | Find all matches (not just the first) |
| `--help` | `-h` | | | Help for test |

#### Examples

```bash
# Test for digits in a string
openGyver regex test "\d+" "order 42 has 3 items"

# Find all matches globally
openGyver regex test --global "\d+" "order 42 has 3 items"

# Test with capture groups (global)
openGyver regex test --global "(\w+)@(\w+)" "alice@example bob@test"

# Test if a string starts with a pattern
openGyver regex test "^hello" "hello world"

# Test an email-like pattern
openGyver regex test "\w+@\w+\.\w+" "contact alice@example.com"

# JSON output for scripting
openGyver regex test "\d+" "order 42" --json

# Test with a complex pattern
openGyver regex test "^[A-Z][a-z]+ \d{4}$" "January 2024"

# Test for IP address pattern
openGyver regex test "\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}" "Server at 192.168.1.1 is up"
```

#### JSON Output Format (single match)

```json
{
  "pattern": "\\d+",
  "input": "order 42 has 3 items",
  "matches": true,
  "match": "42",
  "groups": null
}
```

#### JSON Output Format (global match)

```json
{
  "pattern": "\\d+",
  "input": "order 42 has 3 items",
  "matches": true,
  "all_matches": [
    { "match": "42" },
    { "match": "3" }
  ]
}
```

---

### replace

Replace all occurrences of a regex pattern in the input string. The replacement string may include `$1`, `$2`, etc. for captured groups, and `$0` for the entire match.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `pattern` | Yes | The regular expression pattern to find |
| `replacement` | Yes | The replacement string (supports `$1`, `$2`, etc.) |
| `input` | Yes | The input string to perform replacement on |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Replace all digits with X
openGyver regex replace "\d+" "X" "order 42 has 3 items"
# Output: order X has X items

# Swap capture groups
openGyver regex replace "(\w+)@(\w+)" "$2/$1" "alice@example"
# Output: example/alice

# Collapse whitespace to single hyphens
openGyver regex replace "\s+" "-" "hello   world"
# Output: hello-world

# Remove all vowels
openGyver regex replace "[aeiouAEIOU]" "" "Hello World"

# JSON output
openGyver regex replace "\d+" "NUM" "item 42 costs 99" --json

# Wrap matches in brackets
openGyver regex replace "(\w+)" "[$1]" "hello world"

# Normalize date format
openGyver regex replace "(\d{2})/(\d{2})/(\d{4})" "$3-$1-$2" "01/15/2024"

# Strip HTML tags
openGyver regex replace "<[^>]+>" "" "<b>bold</b> and <i>italic</i>"
```

#### JSON Output Format

```json
{
  "pattern": "\\d+",
  "replacement": "X",
  "input": "order 42 has 3 items",
  "result": "order X has X items"
}
```

---

### extract

Extract all occurrences of a regex pattern from the input. Outputs one match per line. Input can come from a positional argument, `--file/-f`, or piped via stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `pattern` | Yes | The regular expression pattern to extract |
| `input` | No | The input string (or use `--file` / stdin) |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | | Read input from a file |
| `--help` | `-h` | | | Help for extract |

#### Examples

```bash
# Extract email addresses
openGyver regex extract "\w+@\w+\.\w+" "Contact alice@ex.com or bob@ex.com"

# Extract decimal numbers from a file
openGyver regex extract "\d+\.\d+" --file prices.txt

# Extract from piped input
echo "hello 123 world 456" | openGyver regex extract "\d+"

# Extract URLs
openGyver regex extract "https?://[^\s]+" "Visit https://example.com or http://test.org"

# Extract IP addresses from a log file
openGyver regex extract "\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}" --file access.log

# JSON output with match count
openGyver regex extract "\w+@\w+\.\w+" "alice@ex.com bob@ex.com" --json

# Extract hashtags
openGyver regex extract "#\w+" "Check out #golang and #devtools today"

# Extract dates from text
openGyver regex extract "\d{4}-\d{2}-\d{2}" "Created 2024-01-15, updated 2024-06-30"
```

#### JSON Output Format

```json
{
  "pattern": "\\w+@\\w+\\.\\w+",
  "count": 2,
  "matches": ["alice@ex.com", "bob@ex.com"]
}
```

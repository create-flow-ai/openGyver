# json

A collection of JSON utilities for everyday use.

## Usage

```bash
openGyver json [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | | Read input from a file |
| `--json` | `-j` | bool | `false` | Output as JSON (wrap result in a JSON envelope) |
| `--help` | `-h` | | | Help for json |

## Subcommands

### format

Format (beautify) a JSON string with configurable indentation. The input can be provided as a positional argument or read from a file with `--file/-f`. The formatted output is printed to stdout unless `--output/-o` is specified.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `json-string` | No (use `--file` instead) | The JSON string to format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--indent` | | int | `2` | Number of spaces per indentation level |
| `--output` | `-o` | string | | Write result to a file instead of stdout |
| `--help` | `-h` | | | Help for format |

#### Examples

```bash
# Format a compact JSON string
openGyver json format '{"name":"Alice","age":30}'

# Format with 4-space indentation
openGyver json format --indent 4 '{"a":1,"b":2}'

# Read from file, write formatted output to another file
openGyver json format --file input.json --output pretty.json

# Format with JSON envelope output
openGyver json format --json '{"x":1}'

# Format a deeply nested structure
openGyver json format '{"users":[{"name":"Alice","roles":["admin","user"]}]}'

# Use tab-like indentation (8 spaces)
openGyver json format --indent 8 '{"key":"value"}'

# Format and save a minified API response
openGyver json format --file api-response.json --output readable.json

# Format inline and pipe to less
openGyver json format '{"a":1,"b":{"c":2}}' | less
```

#### JSON Output Format

```json
{
  "formatted": "{\n  \"name\": \"Alice\",\n  \"age\": 30\n}",
  "indent": 2
}
```

---

### minify

Compact a JSON string by removing all unnecessary whitespace. The input can be provided as a positional argument or read from a file with `--file/-f`. The minified output is printed to stdout unless `--output/-o` is specified.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `json-string` | No (use `--file` instead) | The JSON string to minify |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write result to a file instead of stdout |
| `--help` | `-h` | | | Help for minify |

#### Examples

```bash
# Minify an indented JSON string
openGyver json minify '{  "name": "Alice",  "age": 30  }'

# Minify a file
openGyver json minify --file pretty.json

# Minify a file and save the result
openGyver json minify --file pretty.json --output compact.json

# Minify with JSON envelope
openGyver json minify --json '{ "x": 1 }'

# Minify a large config file in place (write to new file)
openGyver json minify --file config.json --output config.min.json

# Pipe minified output to clipboard (macOS)
openGyver json minify '{ "a": 1, "b": 2 }' | pbcopy

# Reduce file size of an API fixture
openGyver json minify --file test-fixtures/response.json --output test-fixtures/response.min.json
```

#### JSON Output Format

```json
{
  "minified": "{\"name\":\"Alice\",\"age\":30}"
}
```

---

### validate

Check whether a string is valid JSON. Prints "valid" if the input is well-formed JSON, or an error message describing why it is not.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `json-string` | No (use `--file` instead) | The JSON string to validate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Validate well-formed JSON
openGyver json validate '{"ok": true}'

# Detect invalid JSON
openGyver json validate '{"missing": }'

# Validate a file
openGyver json validate --file data.json

# Validate with JSON envelope output
openGyver json validate --json '{"a":1}'

# Validate a package.json
openGyver json validate --file package.json

# Check if an API response is valid JSON
openGyver json validate '{"users":[],"count":0}'

# Detect a trailing comma issue
openGyver json validate '{"a": 1, "b": 2,}'

# Validate JSON from a build artifact
openGyver json validate --file dist/manifest.json --json
```

#### JSON Output Format

```json
{
  "valid": true
}
```

```json
{
  "valid": false,
  "error": "invalid character '}' looking for beginning of value"
}
```

---

### path

Extract a value from a JSON document using dot-notation path expressions. The JSON input must be provided via `--file/-f`. The path expression is the positional argument.

**Path Syntax:**
- `key` -- Access a top-level key
- `key.nested` -- Nested object access
- `key[0]` -- Array element by index
- `data.users[2].name` -- Mixed object/array traversal

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `expression` | Yes | Dot-notation path expression to evaluate |

#### Flags

No subcommand-specific flags. Uses global flags only. Note: `--file/-f` is required for this subcommand.

#### Examples

```bash
# Extract a simple key
openGyver json path --file config.json 'database.host'

# Extract the first user's name
openGyver json path --file users.json 'users[0].name'

# Deeply nested array traversal
openGyver json path --file data.json 'results[3].tags[0]'

# Get JSON envelope output
openGyver json path --json --file config.json 'server.port'

# Extract a nested configuration value
openGyver json path --file settings.json 'app.logging.level'

# Access an array element in a complex structure
openGyver json path --file api-response.json 'data.items[0].metadata.created_at'

# Extract a top-level value
openGyver json path --file package.json 'version'

# Get chained array indices
openGyver json path --file matrix.json 'grid[0][2]'
```

#### JSON Output Format

```json
{
  "path": "database.host",
  "value": "localhost"
}
```

---

### escape

Escape a raw string so it can be safely embedded inside a JSON document. The result is a JSON string literal (with surrounding double quotes). Special characters (`\n`, `\t`, `\"`, `\\`, etc.) are properly escaped.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `string` | Yes | The raw string to escape |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Escape a string with quotes
openGyver json escape 'hello "world"'

# Escape a string with newlines
openGyver json escape 'line1\nline2'

# Escape with JSON envelope
openGyver json escape --json 'tab\there'

# Escape a file path with backslashes
openGyver json escape 'C:\Users\admin\file.txt'

# Escape multiline content
openGyver json escape 'first line
second line'

# Escape special characters
openGyver json escape 'price: $9.99 & free shipping'

# Escape for embedding in a JSON template
openGyver json escape 'She said "hello" and left.'
```

#### JSON Output Format

```json
{
  "original": "hello \"world\"",
  "escaped": "\"hello \\\"world\\\"\""
}
```

---

### unescape

Unescape a JSON-encoded string, removing surrounding quotes and converting escape sequences back to their literal characters. The input must be a valid JSON string (typically surrounded by double quotes).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `json-string` | Yes | The JSON string literal to unescape |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Unescape a tab character
openGyver json unescape '"hello\tworld"'

# Unescape newlines
openGyver json unescape '"line1\nline2"'

# Unescape with JSON envelope
openGyver json unescape --json '"escaped \"quotes\""'

# Unescape a Unicode sequence
openGyver json unescape '"caf\u00e9"'

# Unescape a complex string
openGyver json unescape '"path: C:\\\\Users\\\\admin"'

# Unescape a backslash-heavy string
openGyver json unescape '"a\\nb\\tc"'

# Restore original text from a JSON value
openGyver json unescape '"She said \"hello\""'
```

#### JSON Output Format

```json
{
  "original": "\"hello\\tworld\"",
  "unescaped": "hello\tworld"
}
```

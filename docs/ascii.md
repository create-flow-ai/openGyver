# ascii

ASCII tools -- generate text banners, print the ASCII table, and look up characters.

## Usage

```bash
openGyver ascii [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | | | Help for ascii |

## Subcommands

### banner

Generate a large block-letter ASCII art banner from the given text. Supports A-Z, 0-9, space, and a few punctuation marks (!, ., -, ?). Characters are rendered as 5-line-high block letters using `#` symbols.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | Yes | The text to render as a banner |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Generate a simple banner
openGyver ascii banner "HELLO"

# Banner with numbers
openGyver ascii banner "GO 1.25"

# Banner with punctuation
openGyver ascii banner "HI!"

# Get JSON output with the banner text
openGyver ascii banner "OK" --json

# Multi-word banner
openGyver ascii banner "HELLO WORLD"

# Banner with digits and letters
openGyver ascii banner "V2.0"
```

#### JSON Output Format

```json
{
  "text": "HELLO",
  "banner": " ###  #####  #      #      ###  \n#   # #      #      #     #   # \n##### ###    #      #     #   # \n#   # #      #      #     #   # \n#   # #####  #####  #####  ### "
}
```

---

### table

Print a table of printable ASCII characters (32-126) with their decimal, hexadecimal, octal, and character representations. Space (decimal 32) is displayed as "SP".

#### Arguments

No arguments.

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Print the full ASCII table
openGyver ascii table

# Print the ASCII table as JSON
openGyver ascii table --json

# Print the table (short flag)
openGyver ascii table -j
```

#### JSON Output Format

```json
[
  {
    "decimal": 32,
    "hex": "0x20",
    "octal": "0040",
    "char": "SP"
  },
  {
    "decimal": 65,
    "hex": "0x41",
    "octal": "0101",
    "char": "A"
  }
]
```

---

### lookup

Look up a character by its decimal value or by the literal character. Shows: character, decimal, hex, octal, binary, HTML entity, and URL encoding. Accepts input as either a decimal number (0-127) or a single ASCII character.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | Yes | A decimal value (0-127) or a single character to look up |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Look up by decimal value
openGyver ascii lookup 65

# Look up by character
openGyver ascii lookup A

# Look up space (decimal 32)
openGyver ascii lookup 32

# Look up a control character
openGyver ascii lookup 10

# JSON output for scripting
openGyver ascii lookup 65 --json

# Look up a punctuation mark
openGyver ascii lookup @

# Look up the DEL character
openGyver ascii lookup 127

# Look up a digit character
openGyver ascii lookup 0
```

#### JSON Output Format

```json
{
  "char": "A",
  "decimal": 65,
  "hex": "0x41",
  "octal": "0101",
  "binary": "01000001",
  "htmlEntity": "&#65;",
  "urlEncoding": "%41"
}
```

# number

Convert numbers between bases, Roman numerals, and IEEE 754 representations.

## Usage

```bash
openGyver number [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for number |

## Subcommands

### base

Convert a number from one base to another. Supports bases 2 through 36. The input value is interpreted in the base given by `--from` (default 10), and printed in the base given by `--to`. Letters a-z (case-insensitive) represent digits 10-35 for bases > 10.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `value` | Yes | The number to convert (in the source base) |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--from` | | int | `10` | Source base (2-36) |
| `--to` | | int | | Target base (2-36, required) |
| `--help` | `-h` | | | Help for base |

#### Examples

```bash
# Decimal to hexadecimal
openGyver number base 255 --to 16
# Output: ff

# Hexadecimal to binary
openGyver number base ff --from 16 --to 2
# Output: 11111111

# Binary to decimal
openGyver number base 11111111 --from 2 --to 10
# Output: 255

# Decimal to base-36
openGyver number base 1000 --to 36
# Output: rs

# Octal to hexadecimal
openGyver number base 377 --from 8 --to 16
# Output: ff

# Decimal to binary
openGyver number base 42 --to 2

# Hex to octal
openGyver number base ff --from 16 --to 8

# JSON output for scripting
openGyver number base 255 --to 16 --json
```

#### JSON Output Format

```json
{
  "input": "255",
  "from_base": 10,
  "to_base": 16,
  "output": "ff"
}
```

---

### roman

Convert a decimal integer to Roman numerals or a Roman numeral string to its decimal value. The direction is auto-detected from the input.

If the input is purely numeric, it is treated as a decimal integer and converted to Roman numeral notation. Otherwise the input is parsed as a Roman numeral string and converted to decimal. Supports values 1 through 3999 (standard Roman numeral range). Input is case-insensitive.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `value` | Yes | A decimal integer (1-3999) or Roman numeral string |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Decimal to Roman
openGyver number roman 42
# Output: XLII

# Decimal to Roman (large number)
openGyver number roman 1994
# Output: MCMXCIV

# Roman to decimal
openGyver number roman XLII
# Output: 42

# Roman to decimal (case-insensitive)
openGyver number roman mcmxciv
# Output: 1994

# JSON output for decimal to Roman
openGyver number roman 42 --json

# Maximum value
openGyver number roman 3999

# Minimum value
openGyver number roman 1

# JSON output for Roman to decimal
openGyver number roman XIV --json
```

#### JSON Output Format (decimal to Roman)

```json
{
  "input": "42",
  "decimal": 42,
  "roman": "XLII"
}
```

#### JSON Output Format (Roman to decimal)

```json
{
  "input": "XLII",
  "roman": "XLII",
  "decimal": 42
}
```

---

### ieee754

Display the IEEE 754 binary representation of a decimal number, showing the sign bit, exponent bits, and mantissa (significand) bits for both 32-bit (float32) and 64-bit (float64) formats.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `number` | Yes | The decimal number to represent (supports `inf`, `nan`) |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Pi representation
openGyver number ieee754 3.14

# Negative number
openGyver number ieee754 -1.5

# Zero
openGyver number ieee754 0

# Infinity
openGyver number ieee754 inf

# Not a Number
openGyver number ieee754 nan

# Small decimal
openGyver number ieee754 0.1

# Large integer
openGyver number ieee754 1000000

# JSON output
openGyver number ieee754 3.14 --json
```

#### JSON Output Format

```json
{
  "input": "3.14",
  "float32": {
    "value": "3.14",
    "bits": "01000000010010001111010111000011",
    "hex": "0x4048F5C3",
    "sign": "0",
    "exponent": "10000000",
    "mantissa": "10010001111010111000011"
  },
  "float64": {
    "value": "3.14",
    "bits": "0100000000001001000111101011100001010001111010111000010100011111",
    "hex": "0x40091EB851EB851F",
    "sign": "0",
    "exponent": "10000000000",
    "mantissa": "1001000111101011100001010001111010111000010100011111"
  }
}
```

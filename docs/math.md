# math

Math expression evaluator and utilities.

## Usage

```bash
openGyver math [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for math |

## Subcommands

### eval

Evaluate a mathematical expression using a recursive descent parser. Supported operators: `+`, `-`, `*`, `/`, `%` (modulo), `^` (power). Supported functions: sqrt, abs, ceil, floor, sin, cos, tan, log, log2, log10. Supported constants: pi, e. Parentheses are supported for grouping. Unary minus is supported.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `expression` | Yes | The math expression to evaluate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Basic arithmetic
openGyver math eval "2 + 3 * 4"

# Parenthesized grouping
openGyver math eval "(2 + 3) * 4"

# Square root
openGyver math eval "sqrt(144)"

# Power operator
openGyver math eval "2^10"

# Trigonometric function with constant
openGyver math eval "sin(pi/2)"

# Complex expression
openGyver math eval "sqrt(144) + 2^3"

# JSON output for scripting
openGyver math eval "log10(1000)" --json

# Modulo operation
openGyver math eval "17 % 5"
```

#### JSON Output Format

```json
{
  "expression": "2 + 3 * 4",
  "result": 14
}
```

---

### percent

Percentage calculator with three modes: compute X% of Y, determine what percentage X is of Y, or calculate the percentage change from X to Y. Exactly one mode flag must be specified.

#### Arguments

No positional arguments. All inputs are provided via flags.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--of` | | bool | `false` | What is X% of Y? (requires --value and --total) |
| `--is` | | bool | `false` | X is what % of Y? (requires --value and --total) |
| `--change` | | bool | `false` | % change from X to Y (requires --from and --to) |
| `--value` | | float64 | `0` | The percentage or value |
| `--total` | | float64 | `0` | The total amount |
| `--from` | | float64 | `0` | Starting value for % change |
| `--to` | | float64 | `0` | Ending value for % change |

#### Examples

```bash
# What is 15% of 200?
openGyver math percent --of --value 15 --total 200

# 30 is what % of 200?
openGyver math percent --is --value 30 --total 200

# Percentage change from 80 to 100
openGyver math percent --change --from 80 --to 100

# JSON output for "of" mode
openGyver math percent --of --value 25 --total 500 --json

# What is 7.5% of 1000?
openGyver math percent --of --value 7.5 --total 1000

# Percentage decrease
openGyver math percent --change --from 100 --to 75

# JSON output for "is" mode
openGyver math percent --is --value 45 --total 180 --json
```

#### JSON Output Format (--of mode)

```json
{
  "mode": "of",
  "value": 15,
  "total": 200,
  "result": 30
}
```

#### JSON Output Format (--is mode)

```json
{
  "mode": "is",
  "value": 30,
  "total": 200,
  "result": 15
}
```

#### JSON Output Format (--change mode)

```json
{
  "mode": "change",
  "from": 80,
  "to": 100,
  "result": 25
}
```

---

### gcd

Calculate the greatest common divisor of two integers using the Euclidean algorithm. Negative values are accepted (their absolute values are used).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `a` | Yes | First integer |
| `b` | Yes | Second integer |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# GCD of 12 and 18
openGyver math gcd 12 18

# GCD of larger numbers
openGyver math gcd 48 36

# GCD with negative numbers
openGyver math gcd -24 16

# JSON output
openGyver math gcd 100 75 --json

# Coprime numbers
openGyver math gcd 17 13

# GCD of equal numbers
openGyver math gcd 42 42
```

#### JSON Output Format

```json
{
  "a": 12,
  "b": 18,
  "result": 6
}
```

---

### lcm

Calculate the least common multiple of two integers. Returns 0 if either input is 0. Negative values are accepted (their absolute values are used).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `a` | Yes | First integer |
| `b` | Yes | Second integer |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# LCM of 4 and 6
openGyver math lcm 4 6

# LCM of larger numbers
openGyver math lcm 12 18

# JSON output
openGyver math lcm 15 20 --json

# LCM with negative numbers
openGyver math lcm -8 12

# Coprime numbers (LCM is their product)
openGyver math lcm 7 11

# LCM of equal numbers
openGyver math lcm 9 9
```

#### JSON Output Format

```json
{
  "a": 4,
  "b": 6,
  "result": 12
}
```

---

### factorial

Calculate the factorial of a non-negative integer N (N!). The maximum supported value is 170, as larger values overflow float64 representation. Returns 1 for both 0! and 1!.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `n` | Yes | Non-negative integer (0-170) |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Factorial of 10
openGyver math factorial 10

# Factorial of 0
openGyver math factorial 0

# Factorial of 20
openGyver math factorial 20

# JSON output
openGyver math factorial 5 --json

# Large factorial
openGyver math factorial 100

# Maximum supported value
openGyver math factorial 170
```

#### JSON Output Format

```json
{
  "n": 5,
  "result": 120
}
```

---

### fibonacci

Calculate the Nth Fibonacci number. F(0) = 0, F(1) = 1, F(n) = F(n-1) + F(n-2). The index must be a non-negative integer.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `n` | Yes | Non-negative Fibonacci index |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# 20th Fibonacci number
openGyver math fibonacci 20

# First Fibonacci numbers
openGyver math fibonacci 0
openGyver math fibonacci 1

# JSON output
openGyver math fibonacci 10 --json

# Larger Fibonacci number
openGyver math fibonacci 50

# Moderate index
openGyver math fibonacci 30
```

#### JSON Output Format

```json
{
  "n": 10,
  "result": 55
}
```

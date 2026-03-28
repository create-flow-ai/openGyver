# convert

Convert values between units of measurement. Automatically detects the category from the unit names.

## Usage

```bash
openGyver convert <value> <from-unit> <to-unit> [flags]
```

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--abbreviated` | `-a` | bool | `false` | Output only the converted value and unit (omit the input side) |
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | bool | `false` | Show help for convert |

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `value` | Yes | Numeric value to convert |
| `from-unit` | Yes | Source unit (case-insensitive, short or long form) |
| `to-unit` | Yes | Target unit (case-insensitive, short or long form) |

## Supported Categories and Units

### Temperature

| Unit | Aliases |
|------|---------|
| Celsius | `c`, `celsius` |
| Fahrenheit | `f`, `fahrenheit` |
| Kelvin | `k`, `kelvin` |

### Length

| Unit | Aliases |
|------|---------|
| Millimeter | `mm` |
| Centimeter | `cm` |
| Meter | `m` |
| Kilometer | `km` |
| Inch | `in` |
| Foot | `ft` |
| Yard | `yd` |
| Mile | `mi` |
| Nautical Mile | `nm` |

### Weight

| Unit | Aliases |
|------|---------|
| Milligram | `mg` |
| Gram | `g` |
| Kilogram | `kg` |
| Ounce | `oz` |
| Pound | `lb` |
| Ton (US) | `ton` |
| Tonne (metric) | `tonne` |
| Stone | `st` |

### Volume

| Unit | Aliases |
|------|---------|
| Milliliter | `ml` |
| Liter | `l` |
| Gallon | `gal` |
| Quart | `qt` |
| Pint | `pt` |
| Cup | `cup` |
| Fluid Ounce | `floz` |
| Tablespoon | `tbsp` |
| Teaspoon | `tsp` |

### Area

| Unit | Aliases |
|------|---------|
| Square Millimeter | `sqmm` |
| Square Centimeter | `sqcm` |
| Square Meter | `sqm` |
| Square Kilometer | `sqkm` |
| Square Inch | `sqin` |
| Square Foot | `sqft` |
| Square Yard | `sqyd` |
| Square Mile | `sqmi` |
| Acre | `acre` |
| Hectare | `hectare` |

### Speed

| Unit | Aliases |
|------|---------|
| Meters/second | `mps` |
| Kilometers/hour | `kph` |
| Miles/hour | `mph` |
| Knots | `knots` |
| Feet/second | `fps` |

### Data

| Unit | Aliases |
|------|---------|
| Byte | `b` |
| Kilobyte | `kb` |
| Megabyte | `mb` |
| Gigabyte | `gb` |
| Terabyte | `tb` |
| Petabyte | `pb` |
| Bit | `bit` |
| Kilobit | `kbit` |
| Megabit | `mbit` |
| Gigabit | `gbit` |

### Time

| Unit | Aliases |
|------|---------|
| Millisecond | `ms` |
| Second | `sec` |
| Minute | `min` |
| Hour | `hr` |
| Day | `day` |
| Week | `week` |
| Month | `month` |
| Year | `year` |

### Currency

Live rates via the Frankfurter API (no API key needed).

Supported currencies: `usd`, `eur`, `gbp`, `jpy`, `cad`, `aud`, `chf`, `cny`, `inr`, `mxn`, `brl`, `krw`, `sgd`, `hkd`, `nok`, `sek`, `dkk`, `nzd`, `zar`, `rub`, `try`, `pln`, `thb`, `idr`, `huf`, `czk`, `ils`, `clp`, `php`, `aed`, `cop`, `sar`, `myr`, `ron`, `bgn`, `hrk`, `isk`, `twd`.

## Examples

```bash
# Length conversion
openGyver convert 100 cm in

# Temperature conversion
openGyver convert 72 f c

# Data size conversion
openGyver convert 1.5 gb mb

# Time duration conversion
openGyver convert 365 days hours

# Speed conversion
openGyver convert 60 mph kph

# Area conversion
openGyver convert 2.5 acre sqft

# Volume conversion
openGyver convert 500 ml cup

# Weight conversion
openGyver convert 150 lb kg

# Currency conversion (live rates)
openGyver convert 100 usd eur

# Abbreviated output (value and unit only)
openGyver convert 100 cm in -a

# JSON output for scripting
openGyver convert 72 f c -j

# Pipe abbreviated output to another command
openGyver convert 100 usd eur -a | cut -d' ' -f1

# Large data conversion
openGyver convert 1 pb gb

# Nautical miles to kilometers
openGyver convert 10 nm km
```

## JSON Output Format

```json
{
  "input_value": 72,
  "input_unit": "Fahrenheit",
  "output_value": 22.222222,
  "output_unit": "Celsius",
  "category": "Temperature"
}
```

## Notes

- Unit names are case-insensitive. Both short and long forms work (e.g., `cm` or `centimeter`).
- The category is automatically detected from the unit names; you do not need to specify it.
- Cross-category conversions (e.g., `cm` to `kg`) produce an error.
- Currency conversion requires an internet connection for live exchange rates.

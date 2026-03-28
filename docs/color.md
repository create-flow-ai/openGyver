# color

Color utilities for working with hex, RGB, HSL, and CMYK colors.

## Usage

```bash
openGyver color [command] [flags]
```

## Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | `false` | Show help for color |
| `--json` | `-j` | bool | `false` | Output as JSON (available to all subcommands) |

## Subcommands

### convert

Convert a color value between hex, RGB, HSL, and CMYK formats. The input format is auto-detected. Use `--to` to specify the desired output format. When `--json` is used, all formats are returned at once regardless of `--to`.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `color` | Yes | Color value in any supported format |

#### Supported Input Formats

| Format | Syntax | Example |
|--------|--------|---------|
| hex | `#rrggbb` or `#rgb` | `"#ff5733"`, `"#f00"` |
| rgb | `rgb(r,g,b)` | `"rgb(255,87,51)"` |
| hsl | `hsl(h,s%,l%)` | `"hsl(11,100%,60%)"` |
| cmyk | `cmyk(c,m,y,k)` | `"cmyk(0,66,80,0)"` |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--to` | | string | `"hex"` | Target format: `hex`, `rgb`, `hsl`, `cmyk` |
| `--help` | `-h` | bool | `false` | Show help for convert |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Examples

```bash
# Convert hex to RGB
openGyver color convert "#ff5733" --to rgb

# Convert RGB to HSL
openGyver color convert "rgb(255,87,51)" --to hsl

# Convert HSL to hex
openGyver color convert "hsl(11,100%,60%)" --to hex

# Convert CMYK to RGB
openGyver color convert "cmyk(0,66,80,0)" --to rgb

# Get all formats at once via JSON
openGyver color convert "#ff5733" --json

# Convert shorthand hex
openGyver color convert "#f00" --to rgb

# Pipe hex output to clipboard (macOS)
openGyver color convert "rgb(100,149,237)" --to hex | pbcopy

# Use JSON output to extract a specific format
openGyver color convert "#ff5733" -j | jq -r '.rgb'
```

#### JSON Output Format

```json
{
  "cmyk": "cmyk(0,65.9,80,0)",
  "hex": "#ff5733",
  "hsl": "hsl(11.2,100%,60%)",
  "rgb": "rgb(255,87,51)"
}
```

---

### contrast

Calculate the WCAG 2.1 contrast ratio between two colors and report whether the pair passes AA and AAA accessibility levels. Both colors are auto-detected and can be in any supported format (hex, rgb, hsl, cmyk).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `color1` | Yes | First color in any supported format |
| `color2` | Yes | Second color in any supported format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | `false` | Show help for contrast |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### WCAG Levels

| Level | Context | Minimum Ratio |
|-------|---------|---------------|
| AA | Normal text | >= 4.5:1 |
| AA | Large text | >= 3.0:1 |
| AAA | Normal text | >= 7.0:1 |
| AAA | Large text | >= 4.5:1 |

#### Examples

```bash
# Check maximum contrast (black and white)
openGyver color contrast "#ffffff" "#000000"

# Mix hex and rgb formats
openGyver color contrast "#ff5733" "rgb(0,0,0)"

# Use HSL for one color
openGyver color contrast "hsl(0,0%,100%)" "#333333" --json

# Check a subtle color pair
openGyver color contrast "#666666" "#999999"

# JSON output for automation
openGyver color contrast "#ffffff" "#000000" -j

# Pipe to jq to check a specific level
openGyver color contrast "#1a1a2e" "#e0e0e0" -j | jq '.aa.normal_text'

# Check brand colors against white background
openGyver color contrast "#4285f4" "#ffffff"

# Verify dark mode contrast
openGyver color contrast "#e0e0e0" "#1a1a1a" -j
```

#### JSON Output Format

```json
{
  "color1": "#ffffff",
  "color2": "#000000",
  "ratio": 21,
  "aa": {
    "normal_text": "PASS",
    "large_text": "PASS"
  },
  "aaa": {
    "normal_text": "PASS",
    "large_text": "PASS"
  }
}
```

---

### palette

Generate a palette of colors derived from a base color.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `color` | Yes | Base color in any supported format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--type` | | string | `"complementary"` | Palette type: `complementary`, `analogous`, `triadic`, `shades`, `tints` |
| `--count` | | int | `5` | Number of colors to generate |
| `--help` | `-h` | bool | `false` | Show help for palette |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Palette Types

| Type | Description |
|------|-------------|
| `complementary` | Colors spread toward the complement (opposite on the color wheel) |
| `analogous` | Colors adjacent on the color wheel (30 degree spread) |
| `triadic` | Colors spread across a 240 degree arc (three-way split) |
| `shades` | Progressively darker versions of the base |
| `tints` | Progressively lighter versions of the base |

#### Examples

```bash
# Generate a default complementary palette (5 colors)
openGyver color palette "#ff5733"

# Generate an analogous palette with 7 colors
openGyver color palette "#ff5733" --type analogous --count 7

# Generate shades from an RGB color
openGyver color palette "rgb(100,149,237)" --type shades

# Generate tints with JSON output
openGyver color palette "#336699" --type tints --count 10 --json

# Triadic palette for design exploration
openGyver color palette "#ff5733" --type triadic --count 6

# Minimal 3-color complementary palette
openGyver color palette "#2ecc71" --type complementary --count 3

# Pipe palette to a file for later use
openGyver color palette "#ff5733" --type analogous --count 12 > palette.txt

# Get all format info for each color in JSON
openGyver color palette "#4a90d9" --type shades --count 8 -j | jq '.palette[].hex'
```

#### JSON Output Format

```json
{
  "base": "#ff5733",
  "type": "complementary",
  "count": 5,
  "palette": [
    {
      "cmyk": "cmyk(0,65.9,80,0)",
      "hex": "#ff5733",
      "hsl": "hsl(11.2,100%,60%)",
      "rgb": "rgb(255,87,51)"
    },
    {
      "cmyk": "cmyk(0,30.2,80,0)",
      "hex": "#ff9f33",
      "hsl": "hsl(32.2,100%,60%)",
      "rgb": "rgb(255,159,51)"
    }
  ]
}
```

---

### name

Find the closest CSS named color for a given color value. Compares the input color against all 148 CSS named colors using Euclidean distance in RGB space and returns the best match.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `color` | Yes | Color value in any supported format |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | `false` | Show help for name |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Examples

```bash
# Find the name for lavender
openGyver color name "#e6e6fa"

# Find nearest CSS name for a custom color
openGyver color name "#ff5734"

# Use RGB input with JSON output
openGyver color name "rgb(100,149,237)" --json

# Check if a color has an exact CSS name match
openGyver color name "#ff0000"

# Identify a dark color
openGyver color name "#2f4f4f"

# Pipe JSON to check if exact match
openGyver color name "#ffa500" -j | jq '.exact'

# Get just the name
openGyver color name "#663399" -j | jq -r '.name'

# Identify an HSL color
openGyver color name "hsl(240,100%,50%)"
```

#### JSON Output Format

```json
{
  "input": "#ff5734",
  "name": "tomato",
  "hex": "#ff6347",
  "exact": false,
  "distance": 21.63
}
```

---

### random

Generate one or more random colors in the specified format.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--format` | | string | `"hex"` | Output format: `hex`, `rgb`, `hsl` |
| `--count` | | int | `1` | Number of colors to generate |
| `--help` | `-h` | bool | `false` | Show help for random |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Examples

```bash
# Generate a single random hex color
openGyver color random

# Generate a random color in RGB format
openGyver color random --format rgb

# Generate 5 random colors in HSL format
openGyver color random --format hsl --count 5

# Generate 10 random colors with full JSON info
openGyver color random --count 10 --json

# Generate a random color and copy to clipboard (macOS)
openGyver color random | pbcopy

# Generate 3 hex colors for a quick palette
openGyver color random --count 3

# Save a batch of random colors to a file
openGyver color random --count 20 --format rgb > colors.txt

# Get all formats for random colors via JSON
openGyver color random --count 5 -j | jq '.[].hex'
```

#### JSON Output Format

```json
{
  "count": 3,
  "colors": [
    {
      "cmyk": "cmyk(58.8,0,39.6,10.2)",
      "hex": "#5fe68b",
      "hsl": "hsl(140,73.8%,63.7%)",
      "rgb": "rgb(95,230,139)"
    },
    {
      "cmyk": "cmyk(0,72.5,47.5,5.5)",
      "hex": "#f04280",
      "hsl": "hsl(339.2,85.1%,60%)",
      "rgb": "rgb(240,66,128)"
    }
  ]
}
```

# convertFont

Convert font files between popular web and desktop formats using fonttools.

## Usage

```bash
openGyver convertFont <input-file> [flags]
```

## Prerequisites

fonttools must be installed (Python package).

```bash
pip install fonttools brotli zopfli
```

The `pyftsubset` command (included with fonttools) must be available in PATH.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input-file` | Yes | Path to the input font file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--help` | `-h` | bool | `false` | Show help for convertFont |

## Supported Formats (12)

| Format | Extension | Description |
|--------|-----------|-------------|
| AFM | `.afm` | Adobe Font Metrics |
| CFF | `.cff` | Compact Font Format |
| DFONT | `.dfont` | Mac Data Fork Font |
| EOT | `.eot` | Embedded OpenType (IE legacy) |
| OTF | `.otf` | OpenType Font |
| PFA | `.pfa` | PostScript Font ASCII |
| PFB | `.pfb` | PostScript Font Binary |
| SFD | `.sfd` | FontForge Spline Font Database |
| TTF | `.ttf` | TrueType Font |
| UFO | `.ufo` | Unified Font Object |
| WOFF | `.woff` | Web Open Font Format 1.0 |
| WOFF2 | `.woff2` | Web Open Font Format 2.0 (Brotli) |

## Conversion Details

| Output Format | Tool / Method |
|---------------|---------------|
| WOFF | `pyftsubset` with `--flavor=woff --no-subset` |
| WOFF2 | `pyftsubset` with `--flavor=woff2 --no-subset` |
| TTF / OTF | `fonttools ttx` (XML intermediate) |

## Examples

```bash
# Convert TrueType to WOFF2 for modern web use
openGyver convertFont font.ttf -o font.woff2

# Convert OpenType to WOFF for broader web compatibility
openGyver convertFont font.otf -o font.woff

# Convert WOFF2 back to TrueType
openGyver convertFont font.woff2 -o font.ttf

# Convert TrueType to legacy Embedded OpenType (for old IE)
openGyver convertFont font.ttf -o font.eot

# Convert OpenType to WOFF2
openGyver convertFont font.otf -o font.woff2

# Convert WOFF to TrueType
openGyver convertFont webfont.woff -o desktop.ttf

# Convert TTF to OTF
openGyver convertFont font.ttf -o font.otf

# Convert OTF to WOFF2 for web deployment
openGyver convertFont display.otf -o display.woff2
```

## Notes

- The input and output formats are auto-detected from file extensions.
- Both the input and output format must be in the supported formats list.
- WOFF and WOFF2 conversions use `pyftsubset` with `--no-subset` to preserve all glyphs.
- WOFF2 requires the `brotli` Python package for Brotli compression.
- TTF/OTF conversions use fonttools' `ttx` command which converts via an XML intermediate representation.
- Not all format combinations may be supported by fonttools. The primary supported output formats are: TTF, OTF, WOFF, WOFF2.

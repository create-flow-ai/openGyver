# openGyver

A Swiss-army-knife CLI tool with 37 commands and 140+ subcommands for everyday conversions, encoding, hashing, generation, formatting, validation, and more. Built in Go for zero-dependency, single-binary distribution across Linux, macOS, and Windows.

Designed to be used standalone, or hooked into CI/CD pipelines, shell scripts, and AI agents.

## Installation

```bash
# Homebrew
brew tap create-flow-ai/tap
brew install opengyver

# Go
go install github.com/mj/opengyver@latest

# From source
git clone https://github.com/create-flow-ai/openGyver.git
cd openGyver && go build -o openGyver .
```

## Quick Start

```bash
openGyver convert 100 cm in                    # unit conversion
openGyver convert 100 usd eur                  # live currency rates
openGyver encode base64 "hello world"           # encoding
openGyver hash sha256 "hello"                   # hashing
openGyver generate password --length 32         # password generator
openGyver stock AAPL -f price                   # stock price
openGyver epoch                                 # current Unix epoch
openGyver timex now --tz Asia/Tokyo             # time in Tokyo
openGyver qr "https://example.com"              # QR code in terminal
openGyver uuid                                  # random UUID v4
openGyver json format '{"a":1}'                 # JSON beautify
openGyver validate html --file index.html       # HTML validation
openGyver format sql "select * from users"      # SQL formatting
openGyver color convert "#ff5733" --to rgb      # color conversion
openGyver finance loan --principal 300000 --rate 6.5 --years 30
openGyver testdata person --count 5 -j          # fake test data
openGyver weather "New York"                    # current weather
openGyver weather Tokyo --date 2024-12-25       # historical weather
```

## Output Modes

Every command supports multiple output modes:

| Mode | Flag | Description |
|------|------|-------------|
| **JSON** | `--json` / `-j` | Structured JSON output for scripting and piping |
| **Abbreviated** | `-a`, `-f`, `-b` | Single value output (command-specific) |
| **Quiet** | `--quiet` / `-q` | Suppress confirmation messages (file converters) |

```bash
openGyver convert -j 100 cm in                 # {"input_value":100,"output_value":39.37,...}
openGyver stock AAPL -f price                  # 248.80
openGyver timex now -b                         # 2026-03-28T08:44:46-04:00
openGyver convertFile data.csv -o data.xlsx -q # silent
```

## Commands

### Conversion Tools

| Command | Description | Docs |
|---------|-------------|------|
| [convert](docs/convert.md) | Unit & currency conversions (9 categories, 38 currencies) | [details](docs/convert.md) |
| [convertAudio](docs/convertAudio.md) | Audio format conversion (33 formats via ffmpeg) | [details](docs/convertAudio.md) |
| [convertCAD](docs/convertCAD.md) | CAD file conversion (DWG, DXF, DWF) | [details](docs/convertCAD.md) |
| [convertEbook](docs/convertEbook.md) | Ebook format conversion (25 formats via Calibre) | [details](docs/convertEbook.md) |
| [convertFile](docs/convertFile.md) | Document & spreadsheet conversion (CSV, XLSX, MD, HTML, DOCX, PDF, PS) | [details](docs/convertFile.md) |
| [convertFont](docs/convertFont.md) | Font format conversion (TTF, OTF, WOFF, WOFF2, EOT + 7 more) | [details](docs/convertFont.md) |
| [convertImage](docs/convertImage.md) | Image format conversion (PNG, JPEG, GIF, BMP, TIFF, WebP, RAW, SVG) | [details](docs/convertImage.md) |
| [convertPresentation](docs/convertPresentation.md) | Presentation conversion (PPTX, KEY, ODP via LibreOffice) | [details](docs/convertPresentation.md) |
| [convertVector](docs/convertVector.md) | Vector graphics conversion (SVG, EPS, AI, CDR + 20 more) | [details](docs/convertVector.md) |
| [convertVideo](docs/convertVideo.md) | Video format conversion (37 formats via ffmpeg) | [details](docs/convertVideo.md) |
| [toGif](docs/toGif.md) | Create animated GIF from images | [details](docs/toGif.md) |
| [toIco](docs/toIco.md) | Create Windows ICO from images | [details](docs/toIco.md) |

### Encoding & Hashing

| Command | Description | Docs |
|---------|-------------|------|
| [encode](docs/encode.md) | Encode/decode: Base64, Base32, Base58, URL, HTML, hex, binary, ROT13, Morse, Punycode, JWT | [details](docs/encode.md) |
| [hash](docs/hash.md) | Hashing: MD5, SHA-1/256/384/512, HMAC, bcrypt, CRC32, Adler-32 | [details](docs/hash.md) |

### Data & Format Tools

| Command | Description | Docs |
|---------|-------------|------|
| [dataformat](docs/dataformat.md) | Convert between YAML, TOML, XML, CSV, and JSON | [details](docs/dataformat.md) |
| [json](docs/json.md) | JSON format, minify, validate, path query, escape/unescape | [details](docs/json.md) |
| [format](docs/format.md) | Format/beautify/minify HTML, XML, CSS, SQL | [details](docs/format.md) |
| [validate](docs/validate.md) | Validate HTML, CSV, XML, YAML, TOML | [details](docs/validate.md) |
| [diff](docs/diff.md) | Compare files: text (unified diff), JSON (structural), CSV | [details](docs/diff.md) |
| [regex](docs/regex.md) | Regex test, replace, extract | [details](docs/regex.md) |

### Generators

| Command | Description | Docs |
|---------|-------------|------|
| [generate](docs/generate.md) | Password, passphrase, API key, secret, OTP, nanoid, snowflake, short ID | [details](docs/generate.md) |
| [uuid](docs/uuid.md) | UUID v4 (random) and v6 (time-sorted) | [details](docs/uuid.md) |
| [qr](docs/qr.md) | QR code generator (ASCII terminal, PNG, SVG) | [details](docs/qr.md) |
| [testdata](docs/testdata.md) | Fake data: people, CSV, JSON, random numbers | [details](docs/testdata.md) |

### Time & Date

| Command | Description | Docs |
|---------|-------------|------|
| [epoch](docs/epoch.md) | Unix epoch: current time, add/subtract durations | [details](docs/epoch.md) |
| [timex](docs/timex.md) | Time tools: now, to-utc, to-tz, to-unix, from-unix, format, diff, add, info | [details](docs/timex.md) |

### Lookup & Reference

| Command | Description | Docs |
|---------|-------------|------|
| [stock](docs/stock.md) | Stock ticker lookup from 35+ global markets (Yahoo Finance, no API key) | [details](docs/stock.md) |
| [weather](docs/weather.md) | Weather: current, forecast (16 days), historical (back to 1940) for any city | [details](docs/weather.md) |
| [network](docs/network.md) | DNS lookup, public IP, WHOIS, CIDR calculator, URL parser, HTTP status codes | [details](docs/network.md) |
| [color](docs/color.md) | Color convert (hex/RGB/HSL/CMYK), WCAG contrast, palette, name lookup | [details](docs/color.md) |
| [number](docs/number.md) | Number base conversion (2-36), Roman numerals, IEEE 754 | [details](docs/number.md) |

### Text Tools

| Command | Description | Docs |
|---------|-------------|------|
| [text](docs/text.md) | Count, case convert, reverse, sort, dedupe, slug, lorem, diff, wrap, trim, replace | [details](docs/text.md) |

### Crypto & Security

| Command | Description | Docs |
|---------|-------------|------|
| [crypto](docs/crypto.md) | AES encrypt/decrypt, RSA keygen, SSH keygen, self-signed cert, CSR | [details](docs/crypto.md) |

### Science & Engineering

| Command | Description | Docs |
|---------|-------------|------|
| [electrical](docs/electrical.md) | Ohm's law, resistor color codes, LED resistor, voltage divider | [details](docs/electrical.md) |
| [geo](docs/geo.md) | Haversine distance, DMS/decimal converter, lat/lon to UTM | [details](docs/geo.md) |
| [accessibility](docs/accessibility.md) | WCAG contrast checker, Flesch/Gunning Fog readability scores | [details](docs/accessibility.md) |

### Finance

| Command | Description | Docs |
|---------|-------------|------|
| [finance](docs/finance.md) | Loan/mortgage, compound interest, ROI, tip, tax, salary, discount, margin | [details](docs/finance.md) |

### Archive

| Command | Description | Docs |
|---------|-------------|------|
| [archive](docs/archive.md) | Create/extract ZIP, TAR, TAR.GZ, 7Z, RAR | [details](docs/archive.md) |

## Plugin Architecture

openGyver uses a plugin-based architecture where each command is a self-contained Go package. Adding a new command requires no changes to existing code.

### Adding a New Command

1. Create a new package under `cmd/`:

```go
// cmd/yourcommand/yourcommand.go
package yourcommand

import (
    "github.com/mj/opengyver/cmd"
    "github.com/spf13/cobra"
)

var yourCmd = &cobra.Command{
    Use:   "yourcommand <args>",
    Short: "One-line description",
    Long:  `Detailed help text with examples.`,
    RunE: func(c *cobra.Command, args []string) error {
        // implementation
        return nil
    },
}

func init() {
    yourCmd.Flags().StringVarP(&output, "output", "o", "default", "description")
    cmd.Register(yourCmd)
}
```

2. Add a blank import in `main.go`:

```go
_ "github.com/mj/opengyver/cmd/yourcommand"
```

3. Build and run. Your command appears in `openGyver --help`.

## Building from Source

```bash
git clone https://github.com/create-flow-ai/openGyver.git
cd openGyver
go build -o openGyver .
```

Run tests:

```bash
go test ./...
```

## Cross-Compilation

```bash
GOOS=linux   GOARCH=amd64 go build -o openGyver-linux-amd64 .
GOOS=linux   GOARCH=arm64 go build -o openGyver-linux-arm64 .
GOOS=windows GOARCH=amd64 go build -o openGyver-windows-amd64.exe .
GOOS=darwin  GOARCH=arm64 go build -o openGyver-darwin-arm64 .
GOOS=darwin  GOARCH=amd64 go build -o openGyver-darwin-amd64 .
GOOS=freebsd GOARCH=amd64 go build -o openGyver-freebsd-amd64 .
```

Each produces a single static binary with zero runtime dependencies.

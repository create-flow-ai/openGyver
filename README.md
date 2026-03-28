# openGyver

A CLI tool for everyday conversions — images, units, currencies, and more. Built in Go for zero-dependency, single-binary distribution across Linux, macOS, and Windows.

Designed to be used standalone, or hooked into CI/CD pipelines, shell scripts, and AI agents.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Global Usage](#global-usage)
- [Commands](#commands)
  - [toIco — Image to ICO](#toico)
  - [toGif — Images to Animated GIF](#togif)
  - [convert — Unit & Currency Conversions](#convert)
    - [Temperature](#temperature)
    - [Length](#length)
    - [Weight](#weight)
    - [Volume](#volume)
    - [Area](#area)
    - [Speed](#speed)
    - [Data (Digital Storage)](#data-digital-storage)
    - [Time](#time)
    - [Currency](#currency)
- [Plugin Architecture](#plugin-architecture)
- [Building from Source](#building-from-source)
- [Cross-Compilation](#cross-compilation)

---

## Installation

Download a prebuilt binary from the [Releases](https://github.com/mj/opengyver/releases) page, or build from source:

```bash
go install github.com/mj/opengyver@latest
```

## Quick Start

```bash
# Convert 100 centimeters to inches
openGyver convert 100 cm in

# Convert Fahrenheit to Celsius
openGyver convert 72 f c

# Convert USD to EUR (live rate)
openGyver convert 250 usd eur

# Convert a PNG to a Windows ICO file
openGyver toIco logo.png -o favicon.ico

# Create an animated GIF from a sequence of images
openGyver toGif frame1.png frame2.png frame3.png -o animation.gif
```

---

## Global Usage

```
openGyver [command] [flags]
```

| Flag       | Description              |
|------------|--------------------------|
| `--help`   | Show help for openGyver or any subcommand |

To get help for any command:

```bash
openGyver --help
openGyver <command> --help
```

### Available Commands

| Command    | Description                                |
|------------|--------------------------------------------|
| `toIco`    | Convert an image to ICO format             |
| `toGif`    | Create an animated GIF from images         |
| `convert`  | Convert between units of measurement       |

---

## Commands

---

### toIco

Convert a PNG, JPEG, or BMP image to a Windows ICO file. Supports embedding multiple icon sizes into a single `.ico` file.

#### Usage

```
openGyver toIco <image> [flags]
```

#### Arguments

| Argument  | Required | Description                          |
|-----------|----------|--------------------------------------|
| `<image>` | Yes      | Path to the source image file (PNG, JPEG, or BMP) |

#### Flags

| Flag              | Short | Type     | Default              | Description                        |
|-------------------|-------|----------|----------------------|------------------------------------|
| `--output`        | `-o`  | string   | `output.ico`         | Output file path                   |
| `--sizes`         |       | ints     | `16,32,48,256`       | Comma-separated icon sizes to embed in the ICO file |
| `--help`          | `-h`  |          |                      | Show help for toIco                |

#### Examples

```bash
# Basic conversion with default sizes (16, 32, 48, 256)
openGyver toIco logo.png

# Specify a custom output path
openGyver toIco logo.png -o favicon.ico

# Embed only specific sizes
openGyver toIco logo.png --sizes 16,32,64

# Combine flags
openGyver toIco banner.jpg -o app-icon.ico --sizes 32,48,128,256
```

#### Notes

- The source image is resized to each requested size and embedded as a separate entry in the ICO file.
- For best results, use a square source image at least as large as the largest requested size.
- Common ICO sizes for web favicons: `16,32,48`. For Windows application icons: `16,32,48,256`.

---

### toGif

Combine multiple images into a single animated GIF. Images are added as frames in the order they are provided.

#### Usage

```
openGyver toGif <image> [image...] [flags]
```

#### Arguments

| Argument            | Required | Description                             |
|---------------------|----------|-----------------------------------------|
| `<image> [image...]`| Yes (1+) | One or more image files to use as frames. Shell globbing (e.g. `frame*.png`) is supported. |

#### Flags

| Flag              | Short | Type   | Default        | Description                             |
|-------------------|-------|--------|----------------|-----------------------------------------|
| `--output`        | `-o`  | string | `output.gif`   | Output file path                        |
| `--delay`         |       | int    | `100`          | Delay between frames in milliseconds    |
| `--loop`          |       | int    | `0`            | Number of times to loop the animation. `0` = infinite loop. |
| `--help`          | `-h`  |        |                | Show help for toGif                     |

#### Examples

```bash
# Create a GIF from three images with default settings (100ms delay, infinite loop)
openGyver toGif frame1.png frame2.png frame3.png

# Use a glob pattern with a custom output name
openGyver toGif frame*.png -o animation.gif

# Speed up the animation (50ms between frames)
openGyver toGif frame*.png --delay 50

# Play the animation exactly 3 times then stop
openGyver toGif frame*.png --loop 3

# Combine all flags
openGyver toGif img_001.png img_002.png img_003.png -o demo.gif --delay 200 --loop 1
```

#### Notes

- Frames are added in the order given on the command line.
- All frames are resized to match the dimensions of the first image.
- Supported input formats: PNG, JPEG, BMP, GIF (first frame only).

---

### convert

Convert a numeric value from one unit to another. The category (length, temperature, currency, etc.) is detected automatically from the unit names.

#### Usage

```
openGyver convert <value> <from-unit> <to-unit>
```

#### Arguments

| Argument      | Required | Description                                     |
|---------------|----------|-------------------------------------------------|
| `<value>`     | Yes      | The numeric value to convert (integer or decimal) |
| `<from-unit>` | Yes      | The source unit (case-insensitive)               |
| `<to-unit>`   | Yes      | The target unit (case-insensitive)               |

#### Flags

| Flag     | Short | Description          |
|----------|-------|----------------------|
| `--help` | `-h`  | Show help for convert |

#### General Notes

- Unit names are **case-insensitive**: `CM`, `cm`, `Cm` all work.
- Both short and long forms are accepted: `cm` or `centimeter`, `f` or `fahrenheit`.
- You can only convert between units in the **same category**. Attempting to convert `cm` to `kg` will return an error.
- Output is automatically formatted: integers stay clean (`100`), decimals are trimmed of trailing zeros (`39.37`).

---

#### Temperature

Convert between Celsius, Fahrenheit, and Kelvin using exact formulas (not factor-based).

| Unit Alias       | Full Name   |
|------------------|-------------|
| `c`, `celsius`   | Celsius     |
| `f`, `fahrenheit` | Fahrenheit |
| `k`, `kelvin`    | Kelvin      |

```bash
openGyver convert 72 f c          # 72 Fahrenheit = 22.222222 Celsius
openGyver convert 100 c f         # 100 Celsius = 212 Fahrenheit
openGyver convert 0 c k           # 0 Celsius = 273.15 Kelvin
openGyver convert 300 k c         # 300 Kelvin = 26.85 Celsius
```

---

#### Length

Convert between metric and imperial length units.

| Unit Alias           | Full Name      | Base (meters) |
|----------------------|----------------|---------------|
| `mm`, `millimeter`   | millimeter     | 0.001         |
| `cm`, `centimeter`   | centimeter     | 0.01          |
| `m`, `meter`         | meter          | 1             |
| `km`, `kilometer`    | kilometer      | 1,000         |
| `in`, `inch`         | inch           | 0.0254        |
| `ft`, `foot`, `feet` | foot           | 0.3048        |
| `yd`, `yard`         | yard           | 0.9144        |
| `mi`, `mile`         | mile           | 1,609.344     |
| `nm`                 | nautical mile  | 1,852         |

```bash
openGyver convert 100 cm in       # 100 centimeter = 39.370079 inch
openGyver convert 5 km mi         # 5 kilometer = 3.106856 mile
openGyver convert 1 nm km         # 1 nautical mile = 1.852 kilometer
openGyver convert 6 ft cm         # 6 foot = 182.88 centimeter
```

---

#### Weight

Convert between metric and imperial mass/weight units.

| Unit Alias           | Full Name     | Base (grams) |
|----------------------|---------------|--------------|
| `mg`, `milligram`    | milligram     | 0.001        |
| `g`, `gram`          | gram          | 1            |
| `kg`, `kilogram`     | kilogram      | 1,000        |
| `oz`, `ounce`        | ounce         | 28.3495      |
| `lb`, `pound`        | pound         | 453.592      |
| `st`, `stone`        | stone         | 6,350.29     |
| `ton`                | short ton (US)| 907,185      |
| `tonne`              | metric tonne  | 1,000,000    |

```bash
openGyver convert 150 lb kg       # 150 pound = 68.0388 kilogram
openGyver convert 1 kg oz         # 1 kilogram = 35.27399 ounce
openGyver convert 10 st lb        # 10 stone = 140.00014 pound
openGyver convert 2.5 tonne ton   # 2.5 metric tonne = 2.755781 short ton
```

---

#### Volume

Convert between metric and US customary volume units.

| Unit Alias             | Full Name          | Base (mL) |
|------------------------|--------------------|-----------|
| `ml`, `milliliter`     | milliliter         | 1         |
| `l`, `liter`           | liter              | 1,000     |
| `gal`, `gallon`        | gallon (US)        | 3,785.41  |
| `qt`, `quart`          | quart (US)         | 946.353   |
| `pt`, `pint`           | pint (US)          | 473.176   |
| `cup`                  | cup (US)           | 236.588   |
| `floz`                 | fluid ounce (US)   | 29.5735   |
| `tbsp`, `tablespoon`   | tablespoon         | 14.7868   |
| `tsp`, `teaspoon`      | teaspoon           | 4.92892   |

```bash
openGyver convert 500 ml cup      # 500 milliliter = 2.113379 cup
openGyver convert 1 gal l         # 1 gallon = 3.78541 liter
openGyver convert 2 cup floz      # 2 cup = 16.000054 fluid ounce
openGyver convert 3 tbsp tsp      # 3 tablespoon = 9.000049 teaspoon
```

---

#### Area

Convert between metric and imperial area units.

| Unit Alias   | Full Name           | Base (sq meters) |
|--------------|---------------------|------------------|
| `sqmm`       | square millimeter   | 0.000001         |
| `sqcm`       | square centimeter   | 0.0001           |
| `sqm`        | square meter        | 1                |
| `sqkm`       | square kilometer    | 1,000,000        |
| `sqin`       | square inch         | 0.00064516       |
| `sqft`       | square foot         | 0.092903         |
| `sqyd`       | square yard         | 0.836127         |
| `sqmi`       | square mile         | 2,589,988        |
| `acre`       | acre                | 4,046.86         |
| `hectare`, `ha` | hectare          | 10,000           |

```bash
openGyver convert 2.5 acre sqft   # 2.5 acre = 108900.14316 square foot
openGyver convert 1 sqmi acre     # 1 square mile = 639.999901 acre
openGyver convert 100 sqm sqft    # 100 square meter = 1076.393035 square foot
openGyver convert 5 hectare acre  # 5 hectare = 12.35527 acre
```

---

#### Speed

Convert between common speed units.

| Unit Alias       | Full Name     | Base (m/s) |
|------------------|---------------|------------|
| `mps`, `m/s`     | meters/sec    | 1          |
| `kph`, `km/h`    | km/hour       | 0.277778   |
| `mph`            | miles/hour    | 0.44704    |
| `knot`, `knots`  | knot          | 0.514444   |
| `fps`, `ft/s`    | feet/sec      | 0.3048     |

```bash
openGyver convert 60 mph kph      # 60 miles/hour = 96.560563 km/hour
openGyver convert 100 kph mph     # 100 km/hour = 62.137168 miles/hour
openGyver convert 30 knots mph    # 30 knot = 34.523419 miles/hour
openGyver convert 10 mps kph      # 10 meters/sec = 36.000029 km/hour
```

---

#### Data (Digital Storage)

Convert between binary-based (1 KB = 1024 bytes) data storage and transfer units.

| Unit Alias | Full Name  | Base (bytes)          |
|------------|------------|-----------------------|
| `bit`      | bit        | 0.125                 |
| `b`        | byte       | 1                     |
| `kb`       | kilobyte   | 1,024                 |
| `mb`       | megabyte   | 1,048,576             |
| `gb`       | gigabyte   | 1,073,741,824         |
| `tb`       | terabyte   | 1,099,511,627,776     |
| `pb`       | petabyte   | 1,125,899,906,842,624 |
| `kbit`     | kilobit    | 128                   |
| `mbit`     | megabit    | 131,072               |
| `gbit`     | gigabit    | 134,217,728           |

```bash
openGyver convert 1.5 gb mb       # 1.5 gigabyte = 1536 megabyte
openGyver convert 1 tb gb         # 1 terabyte = 1024 gigabyte
openGyver convert 100 mbit mb     # 100 megabit = 12.5 megabyte
openGyver convert 8 bit b         # 8 bit = 1 byte
```

---

#### Time

Convert between common time duration units.

| Unit Alias             | Full Name      | Base (seconds) |
|------------------------|----------------|----------------|
| `ms`, `millisecond`    | millisecond    | 0.001          |
| `sec`, `second`        | second         | 1              |
| `min`, `minute`        | minute         | 60             |
| `hr`, `hour`, `hours`  | hour           | 3,600          |
| `day`, `days`          | day            | 86,400         |
| `week`, `weeks`        | week           | 604,800        |
| `month`, `months`      | month (30d)    | 2,592,000      |
| `year`, `years`        | year (365d)    | 31,536,000     |

```bash
openGyver convert 365 days hours  # 365 day = 8760 hour
openGyver convert 1 year days     # 1 year = 365 day
openGyver convert 90 min hr       # 90 minute = 1.5 hour
openGyver convert 5000 ms sec     # 5000 millisecond = 5 second
```

---

#### Currency

Convert between 40 world currencies using **live exchange rates** from the [Frankfurter API](https://frankfurter.app) (free, open-source, no API key required). Rates are sourced from the European Central Bank and updated daily.

**Requires an internet connection.**

| Code  | Currency                | Code  | Currency               |
|-------|-------------------------|-------|------------------------|
| `usd` | US Dollar               | `krw` | South Korean Won       |
| `eur` | Euro                    | `sgd` | Singapore Dollar       |
| `gbp` | British Pound           | `hkd` | Hong Kong Dollar       |
| `jpy` | Japanese Yen            | `nok` | Norwegian Krone        |
| `cad` | Canadian Dollar         | `sek` | Swedish Krona          |
| `aud` | Australian Dollar       | `dkk` | Danish Krone           |
| `chf` | Swiss Franc             | `nzd` | New Zealand Dollar     |
| `cny` | Chinese Yuan            | `zar` | South African Rand     |
| `inr` | Indian Rupee            | `rub` | Russian Ruble          |
| `mxn` | Mexican Peso            | `try` | Turkish Lira           |
| `brl` | Brazilian Real          | `pln` | Polish Zloty           |
| `thb` | Thai Baht               | `idr` | Indonesian Rupiah      |
| `huf` | Hungarian Forint        | `czk` | Czech Koruna           |
| `ils` | Israeli Shekel          | `clp` | Chilean Peso           |
| `php` | Philippine Peso         | `aed` | UAE Dirham             |
| `cop` | Colombian Peso          | `sar` | Saudi Riyal            |
| `myr` | Malaysian Ringgit       | `ron` | Romanian Leu           |
| `bgn` | Bulgarian Lev           | `hrk` | Croatian Kuna          |
| `isk` | Icelandic Krona         | `twd` | Taiwan Dollar          |

```bash
openGyver convert 100 usd eur     # 100 USD = 86.83 EUR (rate varies)
openGyver convert 1000 jpy usd    # 1000 JPY → USD
openGyver convert 50 gbp inr      # 50 GBP → INR
openGyver convert 500 brl eur     # 500 BRL → EUR
```

---

## Plugin Architecture

openGyver uses a plugin-based architecture where each command is a self-contained Go package. Adding a new command requires no changes to existing code.

### Project Structure

```
openGyver/
  main.go                         # Entrypoint — imports all plugins
  cmd/
    root.go                       # Root command + Register() function
    toico/
      toico.go                    # Plugin: image → ICO
    togif/
      togif.go                    # Plugin: images → animated GIF
    convert/
      convert.go                  # Plugin: unit conversion dispatcher
      units.go                    # Category registry and factor-based conversion
      temperature.go              # Temperature (custom formula)
      length.go                   # Length units
      weight.go                   # Weight/mass units
      volume.go                   # Volume units
      area.go                     # Area units
      speed.go                    # Speed units
      data.go                     # Digital storage units
      duration.go                 # Time duration units
      currency.go                 # Currency (live API)
```

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
    // Register flags
    yourCmd.Flags().StringVarP(&output, "output", "o", "default", "description")
    // Register with root
    cmd.Register(yourCmd)
}
```

2. Add a blank import in `main.go`:

```go
import (
    _ "github.com/mj/opengyver/cmd/yourcommand"
)
```

3. Build and run. Your command is now available as `openGyver yourcommand --help`.

---

## Building from Source

```bash
git clone https://github.com/mj/opengyver.git
cd opengyver
go build -o openGyver .
```

## Cross-Compilation

Go makes it trivial to build for any platform:

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o openGyver-linux-amd64 .

# Linux (arm64, e.g. Raspberry Pi 4, AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o openGyver-linux-arm64 .

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o openGyver-windows-amd64.exe .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o openGyver-darwin-arm64 .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o openGyver-darwin-amd64 .

# FreeBSD
GOOS=freebsd GOARCH=amd64 go build -o openGyver-freebsd-amd64 .
```

Each produces a **single static binary** with zero runtime dependencies.

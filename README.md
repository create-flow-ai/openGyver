# openGyver

A CLI tool for everyday conversions — images, units, currencies, documents, time, and more. Built in Go for zero-dependency, single-binary distribution across Linux, macOS, and Windows.

Designed to be used standalone, or hooked into CI/CD pipelines, shell scripts, and AI agents.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Output Modes](#output-modes)
- [Commands](#commands)
  - [accessibility — WCAG & Readability](#accessibility)
  - [archive — Create & Extract Archives](#archive)
  - [color — Color Conversion & Palette Tools](#color)
  - [convert — Unit & Currency Conversions](#convert)
  - [convertAudio — Audio Format Conversions](#convertaudio)
  - [convertCAD — CAD File Conversions](#convertcad)
  - [convertEbook — Ebook Format Conversions](#convertebook)
  - [convertFile — Document & Spreadsheet Conversions](#convertfile)
  - [convertFont — Font Format Conversions](#convertfont)
  - [convertImage — Image Format Conversions](#convertimage)
  - [convertPresentation — Presentation Conversions](#convertpresentation)
  - [convertVector — Vector Graphics Conversions](#convertvector)
  - [convertVideo — Video Format Conversions](#convertvideo)
  - [crypto — Encryption & Key Generation](#crypto)
  - [dataformat — YAML/TOML/XML/CSV Converters](#dataformat)
  - [diff — Text, JSON & CSV Diff](#diff)
  - [electrical — Circuit Calculators](#electrical)
  - [encode — Encoding & Decoding](#encode)
  - [epoch — Unix Epoch Utilities](#epoch)
  - [finance — Financial Calculators](#finance)
  - [format — HTML, XML, CSS, SQL Formatter](#format)
  - [generate — Password, Key & ID Generators](#generate)
  - [geo — Geolocation Tools](#geo)
  - [hash — Hashing & Checksums](#hash)
  - [json — JSON Tools](#json)
  - [network — DNS, IP & Web Utilities](#network)
  - [number — Number Base & Roman Numerals](#number)
  - [qr — QR Code Generator](#qr)
  - [regex — Regular Expression Tools](#regex)
  - [stock — Stock Ticker Lookup](#stock)
  - [testdata — Fake Data Generators](#testdata)
  - [text — Text Manipulation](#text)
  - [timex — Time & Timezone Utilities](#timex)
  - [toGif — Images to Animated GIF](#togif)
  - [toIco — Image to ICO](#toico)
  - [uuid — UUID Generator](#uuid)
  - [validate — HTML, CSV, XML, YAML, TOML Validator](#validate)
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
openGyver convert 100 cm in              # unit conversion
openGyver convert 100 usd eur            # live currency rates
openGyver convertFile data.csv -o data.xlsx
openGyver convertFile README.md -o README.pdf
openGyver convertImage photo.png -o photo.jpg
openGyver epoch                           # current Unix epoch
openGyver epoch add --days 30             # epoch 30 days from now
openGyver timex now --tz Asia/Tokyo       # current time in Tokyo
openGyver qr "https://example.com"        # QR code in terminal
openGyver uuid                            # random UUID v4
```

---

## Output Modes

Every command supports multiple output modes for different use cases:

### JSON Output (`--json` / `-j`)

All commands support `--json` / `-j` for structured JSON output — ideal for scripting, APIs, and piping into `jq`.

**Text commands** return the full data as JSON:

```bash
openGyver convert -j 100 cm in
```
```json
{
  "input_value": 100,
  "input_unit": "centimeter",
  "output_value": 39.37,
  "output_unit": "inch",
  "category": "Length"
}
```

```bash
openGyver stock AAPL -j
```
```json
{
  "symbol": "AAPL",
  "exchange": "NasdaqGS",
  "currency": "USD",
  "price": 248.80,
  "change": -4.09,
  "percent": -1.62,
  "previous_close": 252.89,
  "as_of": "2026-03-27T16:00:02-04:00"
}
```

```bash
openGyver epoch -j
```
```json
{
  "epoch": 1774701886,
  "epoch_ms": 1774701886035,
  "epoch_ns": 1774701886035389000,
  "iso8601": "2026-03-28T12:44:46Z"
}
```

```bash
openGyver timex now -j
```
```json
{
  "timezone": "America/New_York",
  "iso8601": "2026-03-28T08:44:46-04:00",
  "rfc2822": "Sat, 28 Mar 2026 08:44:46 -0400",
  "date": "2026-03-28",
  "time": "08:44:46",
  "unix": 1774701886,
  "unix_ms": 1774701886124
}
```

```bash
openGyver uuid --count 2 -j
```
```json
{
  "version": 4,
  "count": 2,
  "uuids": ["0e54bca5-...", "d4ba1074-..."]
}
```

**Binary/file commands** return conversion metadata as JSON:

```bash
openGyver convertImage photo.png -o photo.jpg -j
```
```json
{
  "success": true,
  "input": "photo.png",
  "output": "photo.jpg",
  "input_format": "png",
  "output_format": "jpeg",
  "width": 1024,
  "height": 768
}
```

```bash
openGyver convertFile data.csv -o data.xlsx -j
```
```json
{
  "success": true,
  "input": "data.csv",
  "output": "data.xlsx",
  "input_format": "csv",
  "output_format": "xlsx"
}
```

### Abbreviated Output

| Command | Flag | Example output |
|---|---|---|
| `convert` | `-a` | `39.37 inch` |
| `stock` | `-f price` | `248.80` |
| `stock` | `-f percent` | `-1.62` |
| `timex *` | `-b` | `2026-03-28T08:44:46-04:00` |
| `timex diff` | `-b` | `14428800` (seconds) |

### Quiet Mode (`--quiet` / `-q`)

File conversion commands support `-q` to suppress the "Converted X → Y" confirmation:

```bash
openGyver convertFile data.csv -o data.xlsx -q
openGyver convertImage photo.png -o photo.jpg -q
openGyver convertAudio song.wav -o song.mp3 -q
openGyver convertVideo input.avi -o output.mp4 -q
```

---

## Commands

---

### accessibility

WCAG contrast ratio checker and readability score calculator.

#### Subcommands

| Subcommand    | Description |
|---------------|-------------|
| `contrast`    | WCAG 2.1 contrast ratio checker — reports AA/AAA pass/fail for normal and large text |
| `readability` | Flesch Reading Ease, Flesch-Kincaid Grade Level, Gunning Fog Index |

```bash
openGyver accessibility contrast "#ffffff" "#000000"
openGyver accessibility contrast "#333" "#ccc"
openGyver accessibility readability "The quick brown fox jumps over the lazy dog."
openGyver accessibility readability --file article.txt
openGyver accessibility readability -j "Some complex text."
```

---

### archive

Create and extract archive files. Implemented in pure Go — no external tools required.

#### Usage

```
openGyver archive create -o <archive> <file> [file...]
openGyver archive extract <archive> [-o <directory>]
```

#### Subcommand: `create`

| Flag       | Short | Type   | Default | Description                |
|------------|-------|--------|---------|----------------------------|
| `--output` | `-o`  | string |         | Output archive path (required) |

#### Subcommand: `extract`

| Flag       | Short | Type   | Default | Description                        |
|------------|-------|--------|---------|------------------------------------|
| `--output` | `-o`  | string | `.`     | Extraction destination directory   |

#### Supported Formats

| Format  | Extensions          | Create | Extract |
|---------|---------------------|--------|---------|
| ZIP     | `.zip`              | yes    | yes     |
| TAR     | `.tar`              | yes    | yes     |
| TAR.GZ  | `.tar.gz`, `.tgz`   | yes    | yes     |

For other formats (7z, RAR, etc.), use system tools: `7z x file.7z`, `unrar x file.rar`.

#### Examples

```bash
openGyver archive create -o backup.zip file1.txt file2.txt dir/
openGyver archive create -o project.tar.gz src/ README.md
openGyver archive create -o files.tar doc1.txt doc2.txt
openGyver archive extract backup.zip
openGyver archive extract backup.zip -o ./extracted/
openGyver archive extract project.tar.gz -o ./project/
```

---

### color

Color conversion, contrast checking, and palette generation.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `convert`  | Convert between hex, RGB, HSL, CMYK. Auto-detects input format. `--to` flag for target. |
| `contrast` | WCAG contrast ratio between two colors |
| `palette`  | Generate palettes: `--type` complementary/analogous/triadic/shades/tints. `--count` default 5. |
| `name`     | Find closest CSS color name for a hex value (148 colors) |
| `random`   | Generate random colors. `--format` hex/rgb/hsl, `--count` |

```bash
openGyver color convert "#ff5733" --to rgb
openGyver color convert "#ff5733" --to rgb -j       # JSON with all formats
openGyver color contrast "#333" "#fff"
openGyver color palette "#ff5733" --type shades --count 5
openGyver color name "#e74c3c"
openGyver color random --count 3 --format hex
```

---

### convert

Convert a numeric value between units of measurement. The category is auto-detected from the unit names.

#### Usage

```
openGyver convert <value> <from-unit> <to-unit> [flags]
```

#### Arguments

| Argument      | Required | Description                                      |
|---------------|----------|--------------------------------------------------|
| `<value>`     | Yes      | Numeric value to convert (integer or decimal)    |
| `<from-unit>` | Yes      | Source unit (case-insensitive)                    |
| `<to-unit>`   | Yes      | Target unit (case-insensitive)                    |

#### Flags

| Flag             | Short | Type | Default | Description                             |
|------------------|-------|------|---------|-----------------------------------------|
| `--abbreviated`  | `-a`  | bool | false   | Output only the converted value and unit |

#### Categories and Units

**Temperature** — `c`, `celsius`, `f`, `fahrenheit`, `k`, `kelvin`

```bash
openGyver convert 72 f c             # 72 Fahrenheit = 22.222222 Celsius
openGyver convert 100 c f            # 100 Celsius = 212 Fahrenheit
openGyver convert 0 c k              # 0 Celsius = 273.15 Kelvin
```

**Length** — `mm`, `cm`, `m`, `km`, `in`, `inch`, `ft`, `foot`, `feet`, `yd`, `yard`, `mi`, `mile`, `nm` (nautical mile)

```bash
openGyver convert 100 cm in          # 100 centimeter = 39.370079 inch
openGyver convert 5 km mi            # 5 kilometer = 3.106856 mile
openGyver convert 6 ft cm            # 6 foot = 182.88 centimeter
```

**Weight** — `mg`, `g`, `kg`, `oz`, `ounce`, `lb`, `pound`, `st`, `stone`, `ton`, `tonne`

```bash
openGyver convert 150 lb kg          # 150 pound = 68.0388 kilogram
openGyver convert 1 kg oz            # 1 kilogram = 35.27399 ounce
```

**Volume** — `ml`, `l`, `gal`, `gallon`, `qt`, `quart`, `pt`, `pint`, `cup`, `floz`, `tbsp`, `tablespoon`, `tsp`, `teaspoon`

```bash
openGyver convert 500 ml cup         # 500 milliliter = 2.113379 cup
openGyver convert 1 gal l            # 1 gallon = 3.78541 liter
```

**Area** — `sqmm`, `sqcm`, `sqm`, `sqkm`, `sqin`, `sqft`, `sqyd`, `sqmi`, `acre`, `hectare`, `ha`

```bash
openGyver convert 2.5 acre sqft      # 2.5 acre = 108900 square foot
openGyver convert 5 hectare acre     # 5 hectare = 12.355 acre
```

**Speed** — `mps`, `m/s`, `kph`, `km/h`, `mph`, `knot`, `knots`, `fps`, `ft/s`

```bash
openGyver convert 60 mph kph         # 60 miles/hour = 96.56 km/hour
openGyver convert 100 kph mph        # 100 km/hour = 62.14 miles/hour
```

**Data** — `b`, `kb`, `mb`, `gb`, `tb`, `pb`, `bit`, `kbit`, `mbit`, `gbit`

```bash
openGyver convert 1.5 gb mb          # 1.5 gigabyte = 1536 megabyte
openGyver convert 100 mbit mb        # 100 megabit = 12.5 megabyte
```

**Time** — `ms`, `sec`, `min`, `hr`, `hour`, `hours`, `day`, `days`, `week`, `weeks`, `month`, `months`, `year`, `years`

```bash
openGyver convert 365 days hours     # 365 day = 8760 hour
openGyver convert 90 min hr          # 90 minute = 1.5 hour
```

**Currency** — 38 currencies with live rates via [Frankfurter API](https://frankfurter.app) (free, no API key):

`usd`, `eur`, `gbp`, `jpy`, `cad`, `aud`, `chf`, `cny`, `inr`, `mxn`, `brl`, `krw`, `sgd`, `hkd`, `nok`, `sek`, `dkk`, `nzd`, `zar`, `rub`, `try`, `pln`, `thb`, `idr`, `huf`, `czk`, `ils`, `clp`, `php`, `aed`, `cop`, `sar`, `myr`, `ron`, `bgn`, `hrk`, `isk`, `twd`

```bash
openGyver convert 100 usd eur        # 100 USD = 86.83 EUR (live rate)
openGyver convert -a 100 usd eur     # 86.83 EUR (abbreviated)
```

---

### convertAudio

Convert between audio formats using ffmpeg.

#### Usage

```
openGyver convertAudio <input-file> -o <output-file> [flags]
```

**Requires:** `ffmpeg` (`brew install ffmpeg` / `apt install ffmpeg`)

#### Flags

| Flag        | Short | Type   | Default | Description                         |
|-------------|-------|--------|---------|-------------------------------------|
| `--output`  | `-o`  | string |         | Output file path (required)         |
| `--bitrate` |       | string |         | Audio bitrate (e.g. `128k`, `320k`) |
| `--sample`  |       | string |         | Sample rate in Hz (e.g. `44100`)    |

#### Supported Formats

AAC, AC3, AIF, AIFC, AIFF, AMR, AU, CAF, DSS, FLAC, M4A, M4B, MP3, OGA, OGG, OPUS, VOC, WAV, WEBA, WMA

#### Examples

```bash
openGyver convertAudio song.wav -o song.mp3
openGyver convertAudio song.flac -o song.aac
openGyver convertAudio podcast.mp3 -o podcast.ogg
openGyver convertAudio song.wav -o song.mp3 --bitrate 320k
openGyver convertAudio song.mp3 -o song.wav --sample 44100
```

---

### convertCAD

Convert between CAD file formats.

#### Usage

```
openGyver convertCAD <input-file> -o <output-file>
```

**Requires:** ODA File Converter or LibreCAD

#### Supported Formats

DWG, DXF, DWF, PDF, SVG, PNG

#### Examples

```bash
openGyver convertCAD drawing.dwg -o drawing.dxf
openGyver convertCAD drawing.dxf -o drawing.pdf
openGyver convertCAD drawing.dxf -o drawing.svg
```

---

### convertEbook

Convert between ebook formats using Calibre.

#### Usage

```
openGyver convertEbook <input-file> -o <output-file>
```

**Requires:** Calibre (`brew install calibre` / `apt install calibre`)

#### Supported Formats

AZW, AZW3, AZW4, CBC, CBR, CBZ, CHM, DOCX, EPUB, FB2, HTM, HTML, HTMLZ, LIT, LRF, MOBI, PDB, PDF, PML, PRC, RB, SNB, TCR, TXT, TXTZ

#### Examples

```bash
openGyver convertEbook book.epub -o book.mobi
openGyver convertEbook book.epub -o book.pdf
openGyver convertEbook book.mobi -o book.epub
openGyver convertEbook document.docx -o document.epub
openGyver convertEbook book.epub -o book.azw3
```

---

### convertFile

Convert between document, spreadsheet, and page-layout file formats.

#### Usage

```
openGyver convertFile <input-file> [flags]
```

#### Arguments

| Argument       | Required | Description            |
|----------------|----------|------------------------|
| `<input-file>` | Yes      | Path to the input file |

#### Flags

| Flag          | Short | Type   | Default                             | Description                            |
|---------------|-------|--------|-------------------------------------|----------------------------------------|
| `--output`    | `-o`  | string | input name with new extension       | Output file path                       |
| `--sheet`     |       | string |                                     | Sheet name for XLSX/Numbers (default: first sheet) |
| `--delimiter` |       | string | auto-detected                       | CSV delimiter: `comma`, `tab`, `semicolon`, `pipe`, or any character |

#### Supported Formats

| Format     | Extensions                |
|------------|---------------------------|
| CSV        | `.csv`, `.tsv`            |
| XLSX       | `.xlsx`                   |
| Numbers    | `.numbers` (stub)         |
| Text       | `.txt`, `.text`           |
| Markdown   | `.md`, `.markdown`        |
| HTML       | `.html`, `.htm`           |
| DOCX       | `.docx`                   |
| PDF        | `.pdf` (output only)      |
| PostScript | `.ps` (output only)       |

#### Spreadsheet Conversions

```bash
openGyver convertFile data.csv                         # → data.xlsx (default)
openGyver convertFile data.csv -o report.xlsx
openGyver convertFile report.xlsx                      # → report.csv
openGyver convertFile report.xlsx --sheet "Sales Q1"
openGyver convertFile report.xlsx -o out.csv --delimiter ";"
```

#### Document Conversions

```bash
# Markdown ↔ HTML
openGyver convertFile README.md -o README.html
openGyver convertFile page.html -o page.md

# Markdown / HTML → Text (strip formatting)
openGyver convertFile README.md -o README.txt
openGyver convertFile page.html -o page.txt

# Text → HTML / Markdown
openGyver convertFile notes.txt -o notes.html
openGyver convertFile notes.txt -o notes.md

# DOCX ↔ Text / HTML / Markdown
openGyver convertFile report.docx -o report.txt
openGyver convertFile report.docx -o report.html
openGyver convertFile report.docx -o report.md
openGyver convertFile notes.txt -o notes.docx
openGyver convertFile README.md -o README.docx
openGyver convertFile page.html -o page.docx
```

#### PDF Output

```bash
openGyver convertFile notes.txt -o notes.pdf
openGyver convertFile README.md -o README.pdf
openGyver convertFile page.html -o page.pdf
openGyver convertFile data.csv -o data.pdf            # table layout
openGyver convertFile report.xlsx -o report.pdf        # table layout
openGyver convertFile report.docx -o report.pdf
```

#### PostScript Output

```bash
openGyver convertFile notes.txt -o notes.ps
openGyver convertFile README.md -o README.ps
openGyver convertFile data.csv -o data.ps              # tabular layout
```

---

### convertFont

Convert between font formats using fonttools.

#### Usage

```
openGyver convertFont <input-file> -o <output-file>
```

**Requires:** fonttools (`pip install fonttools brotli zopfli`)

#### Supported Formats

| Format | Description                       |
|--------|-----------------------------------|
| TTF    | TrueType Font                     |
| OTF    | OpenType Font                     |
| WOFF   | Web Open Font Format 1.0          |
| WOFF2  | Web Open Font Format 2.0 (Brotli) |
| EOT    | Embedded OpenType (IE legacy)     |

#### Examples

```bash
openGyver convertFont font.ttf -o font.woff2
openGyver convertFont font.otf -o font.woff
openGyver convertFont font.woff2 -o font.ttf
```

---

### convertImage

Convert between raster image formats with optional resizing.

#### Usage

```
openGyver convertImage <image> [flags]
```

#### Arguments

| Argument  | Required | Description            |
|-----------|----------|------------------------|
| `<image>` | Yes      | Path to the input image |

#### Flags

| Flag        | Short | Type   | Default | Description                                      |
|-------------|-------|--------|---------|--------------------------------------------------|
| `--output`  | `-o`  | string |         | Output file path (required)                      |
| `--quality` |       | int    | 90      | JPEG quality 1-100                               |
| `--width`   |       | int    | 0       | Resize to this width (0 = keep original)         |
| `--height`  |       | int    | 0       | Resize to this height (0 = keep original)        |
| `--format`  | `-f`  | string |         | Override input format detection                  |

If only `--width` or `--height` is given, the other scales proportionally.

#### Supported Formats

| Format | Read | Write | Extensions                |
|--------|------|-------|---------------------------|
| PNG    | yes  | yes   | `.png`                    |
| JPEG   | yes  | yes   | `.jpg`, `.jpeg`, `.jfif`  |
| GIF    | yes  | yes   | `.gif`                    |
| BMP    | yes  | yes   | `.bmp`                    |
| TIFF   | yes  | yes   | `.tiff`, `.tif`           |
| PPM    | yes  | yes   | `.ppm`, `.pgm`, `.pbm`    |
| SVG    | no   | yes*  | `.svg` (raster → vector)  |
| WebP   | yes  | no    | `.webp`                   |
| HEIC   | yes* | no    | `.heic`, `.heif`          |
| Camera RAW | yes* | no | `.cr2`, `.nef`, `.arw`, `.dng`, `.orf`, `.raf`, `.rw2`, `.pef`, and 16 more |

\* SVG output uses potrace (`brew install potrace`) or ImageMagick for raster-to-vector tracing.
\* HEIC and Camera RAW decoding requires dcraw, ImageMagick, or Apple sips.

#### Examples

```bash
# Basic format conversion
openGyver convertImage photo.png -o photo.jpg
openGyver convertImage photo.jpg -o photo.png
openGyver convertImage photo.jpg -o photo.bmp
openGyver convertImage photo.bmp -o photo.tiff
openGyver convertImage photo.webp -o photo.png

# JPEG compression quality
openGyver convertImage photo.png -o photo.jpg --quality 85
openGyver convertImage photo.png -o lo-fi.jpg --quality 30

# Resize
openGyver convertImage photo.png -o thumb.jpg --width 200
openGyver convertImage photo.png -o resized.png --width 800 --height 600

# Raster → SVG (vector tracing)
openGyver convertImage logo.png -o logo.svg
openGyver convertImage icon.bmp -o icon.svg

# Camera RAW → standard format
openGyver convertImage photo.cr2 -o photo.jpg
openGyver convertImage photo.nef -o photo.png
openGyver convertImage photo.dng -o photo.tiff

# PPM output
openGyver convertImage photo.png -o photo.ppm

# Format hint for files without standard extensions
openGyver convertImage raw.dat -o out.png --format bmp
```

---

### convertPresentation

Convert between presentation formats using LibreOffice.

#### Usage

```
openGyver convertPresentation <input-file> -o <output-file>
```

**Requires:** LibreOffice (`brew install --cask libreoffice` / `apt install libreoffice`)

#### Supported Formats

| Direction | Formats                                           |
|-----------|---------------------------------------------------|
| Input     | DPS, KEY, ODP, POT, POTX, PPS, PPSX, PPT, PPTM, PPTX |
| Output    | ODP, PDF, PPTX, PPT, HTML, PNG, JPG, SVG         |

#### Examples

```bash
openGyver convertPresentation slides.pptx -o slides.pdf
openGyver convertPresentation slides.pptx -o slides.odp
openGyver convertPresentation slides.key -o slides.pptx
openGyver convertPresentation slides.pptx -o slides.html
openGyver convertPresentation old.ppt -o new.pptx
```

---

### convertVector

Convert between vector graphics formats.

#### Usage

```
openGyver convertVector <input-file> -o <output-file> [flags]
```

**Requires:** `rsvg-convert` (librsvg) or Inkscape

| Flag       | Short | Type   | Default | Description                      |
|------------|-------|--------|---------|----------------------------------|
| `--output` | `-o`  | string |         | Output file path (required)      |
| `--width`  |       | int    | 0       | Output width (for rasterization) |
| `--height` |       | int    | 0       | Output height (for rasterization)|

#### Supported Formats

SVG, SVGZ, EPS, PDF, PS, EMF, WMF, AI, CDR, CGM, VSD + raster output (PNG, JPG)

#### Examples

```bash
openGyver convertVector logo.svg -o logo.png
openGyver convertVector logo.svg -o logo.pdf
openGyver convertVector logo.svg -o logo.eps
openGyver convertVector logo.svg -o logo.png --width 1024 --height 768
openGyver convertVector diagram.eps -o diagram.svg
```

---

### convertVideo

Convert between video formats using ffmpeg.

#### Usage

```
openGyver convertVideo <input-file> -o <output-file> [flags]
```

**Requires:** `ffmpeg` (`brew install ffmpeg` / `apt install ffmpeg`)

#### Flags

| Flag           | Short | Type   | Default | Description                              |
|----------------|-------|--------|---------|------------------------------------------|
| `--output`     | `-o`  | string |         | Output file path (required)              |
| `--resolution` |       | string |         | Output resolution (e.g. `1920x1080`)     |
| `--vbitrate`   |       | string |         | Video bitrate (e.g. `2M`, `5M`)          |
| `--abitrate`   |       | string |         | Audio bitrate (e.g. `128k`, `192k`)      |
| `--fps`        |       | string |         | Frames per second (e.g. `24`, `30`, `60`)|
| `--codec`      |       | string |         | Video codec (e.g. `libx264`, `libvpx-vp9`)|

#### Supported Formats

3G2, 3GP, 3GPP, AVI, CAVS, DV, DVR, FLV, M2TS, M4V, MKV, MOD, MOV, MP4, MPEG, MPG, MTS, MXF, OGG, OGV, RM, RMVB, SWF, TS, VOB, WEBM, WMV, WTV + GIF output

#### Examples

```bash
openGyver convertVideo input.avi -o output.mp4
openGyver convertVideo input.mkv -o output.webm
openGyver convertVideo input.mov -o output.mp4 --resolution 1920x1080
openGyver convertVideo input.mp4 -o output.mp4 --vbitrate 5M --abitrate 192k
openGyver convertVideo input.mp4 -o output.gif --fps 10
openGyver convertVideo input.mp4 -o output.webm --codec libvpx-vp9
```

---

### crypto

Cryptographic tools — AES encryption, RSA/SSH key generation, SSL certificates.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `aes`      | AES-256-GCM encrypt/decrypt. `--key` (hex or passphrase, required). `--decrypt`/`-d` to decrypt. |
| `rsa`      | RSA key pair generator. `--bits` (default 2048). `--output-dir` to write files. |
| `sshkey`   | SSH key pair generator. `--type` ed25519 (default) or rsa. `--comment`. |
| `cert`     | Self-signed certificate generator. `--cn` (required), `--days` (default 365). |
| `csr`      | Certificate Signing Request generator. `--cn` (required), `--org`, `--country`. |

```bash
openGyver crypto aes "secret message" --key "mypassword"
openGyver crypto aes "encrypted_base64" --key "mypassword" -d
openGyver crypto rsa --bits 4096 --output-dir ./keys/
openGyver crypto sshkey --type ed25519 --comment "me@host"
openGyver crypto cert --cn "localhost" --days 365 --output-dir ./certs/
openGyver crypto csr --cn "example.com" --org "My Corp" --country US
```

---

### dataformat

Convert between data serialization formats — YAML, TOML, XML, CSV, JSON.

#### Subcommands

| Subcommand  | Description |
|-------------|-------------|
| `yaml2json` | YAML to JSON |
| `json2yaml` | JSON to YAML |
| `toml2json` | TOML to JSON |
| `json2toml` | JSON to TOML |
| `csv2json`  | CSV to JSON (first row = headers, outputs array of objects) |
| `json2csv`  | JSON array of objects to CSV |
| `xml2json`  | XML to JSON |
| `json2xml`  | JSON to XML |

All accept input as argument or `--file`/`-f`. Output to stdout or `--output`/`-o`.

```bash
openGyver dataformat yaml2json --file config.yml
openGyver dataformat json2yaml '{"name":"John","age":30}'
openGyver dataformat csv2json --file data.csv
openGyver dataformat toml2json --file config.toml -o config.json
openGyver dataformat json2xml '{"user":{"name":"Alice"}}'
```

---

### diff

Compare files — unified text diff, JSON structural diff, CSV row diff.

#### Subcommands

| Subcommand | Flags | Description |
|------------|-------|-------------|
| `text`     | `--file1`, `--file2` | Unified diff with +/- markers |
| `json`     | `--file1`, `--file2` | Structural diff showing added/removed/changed keys |
| `csv`      | `--file1`, `--file2` | Row-by-row CSV comparison |

```bash
openGyver diff text --file1 old.txt --file2 new.txt
openGyver diff json --file1 v1.json --file2 v2.json
openGyver diff csv --file1 before.csv --file2 after.csv
openGyver diff json --file1 a.json --file2 b.json -j
```

---

### electrical

Circuit design calculators — Ohm's law, resistor codes, LED resistors, voltage dividers.

#### Subcommands

| Subcommand | Flags | Description |
|------------|-------|-------------|
| `ohm`      | Any 2 of: `--voltage`/`-v`, `--current`/`-i`, `--resistance`/`-r` | Calculates the third value + power |
| `resistor` | Arg: resistance value (e.g. `4700`, `4.7k`, `4k7`) | Shows 4-band and 5-band color codes |
| `led`      | `--source`, `--forward`, `--current` (default 20mA) | Calculates resistor value for LED circuit |
| `divider`  | `--vin`, `--r1`, `--r2` | Calculates Vout, ratio, current, power |

```bash
openGyver electrical ohm --voltage 12 --resistance 100
openGyver electrical resistor 4.7k
openGyver electrical led --source 5 --forward 2.1
openGyver electrical divider --vin 12 --r1 10000 --r2 4700
```

---

### encode

Encode and decode text in various formats.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `base64`   | Base64 encode/decode (`--decode`/`-d`) |
| `base32`   | Base32 encode/decode |
| `base58`   | Base58 (Bitcoin alphabet) encode/decode |
| `url`      | URL percent-encoding encode/decode |
| `html`     | HTML entity encode/decode |
| `hex`      | Text to hex / hex to text |
| `binary`   | Text to 8-bit binary / binary to text |
| `rot13`    | ROT13 (symmetric) |
| `morse`    | Morse code encode/decode |
| `punycode` | Unicode to Punycode / Punycode to Unicode |
| `jwt`      | Decode JWT token payload (no verification) |

All accept input as argument or `--file`/`-f`. All support `--json`/`-j`.

```bash
openGyver encode base64 "hello world"           # aGVsbG8gd29ybGQ=
openGyver encode base64 -d "aGVsbG8gd29ybGQ="   # hello world
openGyver encode url "hello world & foo"
openGyver encode hex "hello"
openGyver encode rot13 "hello"
openGyver encode morse "SOS"
openGyver encode jwt "eyJhbG..."
openGyver encode base64 "data" -j               # JSON output
```

---

### finance

Financial calculators for everyday money math.

#### Subcommands

| Subcommand | Key Flags | Description |
|------------|-----------|-------------|
| `loan`     | `--principal`, `--rate`, `--years` | Monthly payment, total payment, total interest |
| `compound` | `--principal`, `--rate`, `--years`, `--frequency` (default 12) | Final amount, total interest |
| `roi`      | `--initial`, `--final` | ROI percentage, profit/loss |
| `tip`      | `--amount`, `--percent` (default 18), `--split` (default 1) | Tip, total, per person |
| `tax`      | `--amount`, `--rate` | Tax amount, total |
| `salary`   | `--amount`, `--from`, `--to` (hourly/daily/weekly/monthly/yearly) | Pay period conversion |
| `discount` | `--price`, `--percent` | Discount amount, final price |
| `margin`   | `--cost`, `--revenue` | Profit, margin %, markup % |

```bash
openGyver finance loan --principal 300000 --rate 6.5 --years 30
openGyver finance compound --principal 10000 --rate 7 --years 10
openGyver finance tip --amount 85.50 --percent 20 --split 4
openGyver finance salary --amount 50 --from hourly --to yearly
openGyver finance margin --cost 40 --revenue 100 -j
```

---

### format

Format, beautify, and minify HTML, XML, CSS, and SQL.

#### Subcommands

| Subcommand | Flags | Description |
|------------|-------|-------------|
| `html`     | `--indent` (2), `--minify` | HTML beautifier/minifier. Proper indentation, handles self-closing tags. |
| `xml`      | `--indent` (2), `--minify` | XML formatter/minifier. |
| `css`      | `--indent` (2), `--minify` | CSS formatter — one property per line, proper indentation. Minify strips all whitespace. |
| `sql`      | `--indent` (2), `--minify` | SQL formatter — uppercase keywords, newlines before clauses. Minify collapses to one line. |

All accept input as argument or `--file`/`-f`. Output to stdout or `--output`/`-o`. All support `--json`/`-j`.

```bash
openGyver format html "<div><p>hello</p></div>"
openGyver format html --minify --file index.html -o index.min.html
openGyver format xml '<root><item id="1"/></root>'
openGyver format css "body{color:red;margin:0}"
openGyver format css --minify --file styles.css
openGyver format sql "select * from users where id = 1 order by name"
openGyver format sql --minify "SELECT * FROM users WHERE active = 1"
```

---

### generate

Random generators — passwords, keys, IDs, OTP secrets.

#### Subcommands

| Subcommand   | Key Flags | Description |
|--------------|-----------|-------------|
| `password`   | `--length` (16), `--no-upper/lower/digits/special`, `--count` | Random password |
| `passphrase` | `--words` (4), `--separator` ("-"), `--count` | Diceware passphrase (1296-word list) |
| `string`     | `--length` (32), `--charset` (alpha/alphanumeric/hex/base64/custom), `--count` | Random string |
| `apikey`     | `--prefix`, `--length` (32) | Base62 API key |
| `secret`     | `--length` (64 bytes) | Hex-encoded secret key |
| `otp`        | `--issuer`, `--account` | TOTP secret + otpauth:// URI |
| `nanoid`     | `--length` (21), `--alphabet` | Nano ID |
| `snowflake`  | | 64-bit timestamp-based ID |
| `shortid`    | `--length` (8) | Short base62 ID |

All use `crypto/rand`. All support `--json`/`-j`.

```bash
openGyver generate password                       # random 16-char password
openGyver generate password --length 32 --count 5
openGyver generate passphrase --words 6
openGyver generate apikey --prefix "sk_live_"
openGyver generate secret --length 32
openGyver generate otp --issuer "MyApp" --account "user@example.com"
openGyver generate nanoid
openGyver generate snowflake
```

---

### geo

Geolocation tools — distance calculation, coordinate conversion.

#### Subcommands

| Subcommand | Flags | Description |
|------------|-------|-------------|
| `distance` | `--lat1`, `--lon1`, `--lat2`, `--lon2` | Haversine great-circle distance (km and miles) |
| `dms`      | Arg: decimal or DMS | Convert between decimal degrees and DMS notation |
| `utm`      | `--lat`, `--lon` | Convert lat/lon to UTM coordinates |

```bash
openGyver geo distance --lat1 40.7128 --lon1 -74.006 --lat2 51.5074 --lon2 -0.1278
openGyver geo dms 40.7128                        # → 40°42'46.08"N
openGyver geo utm --lat 40.7128 --lon -74.006
```

---

### hash

Compute cryptographic hashes and checksums.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `md5`      | MD5 hash |
| `sha1`     | SHA-1 hash |
| `sha256`   | SHA-256 hash |
| `sha384`   | SHA-384 hash |
| `sha512`   | SHA-512 hash |
| `hmac`     | HMAC generator (`--key`, `--algorithm` sha256/sha1/md5/sha384/sha512) |
| `bcrypt`   | bcrypt hash (`--rounds` default 10) or verify (`--verify` hash) |
| `crc32`    | CRC-32 checksum |
| `adler32`  | Adler-32 checksum |

All accept input as argument or `--file`/`-f`. `--uppercase`/`-u` for hex output.

```bash
openGyver hash sha256 "hello"
openGyver hash md5 --file document.pdf
openGyver hash hmac "message" --key "secret" --algorithm sha256
openGyver hash bcrypt "password123"
openGyver hash bcrypt "password123" --verify '$2a$10$...'
openGyver hash sha256 "data" -j -u
```

---

### json

JSON tools — format, minify, validate, path query, escape/unescape.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `format`   | Beautify JSON (`--indent` default 2). `--file`/`-f`, `--output`/`-o`. |
| `minify`   | Remove whitespace from JSON. `--file`/`-f`, `--output`/`-o`. |
| `validate` | Check if JSON is valid. Outputs "valid" or error. |
| `path`     | Evaluate dot-notation path (e.g. `data.users[0].name`). `--file`/`-f`. |
| `escape`   | Escape a string for JSON embedding |
| `unescape` | Unescape a JSON string literal |

```bash
openGyver json format '{"a":1,"b":2}'
openGyver json minify --file data.json -o data.min.json
openGyver json validate '{"valid": true}'
openGyver json path "users[0].name" --file data.json
openGyver json escape 'hello "world"'
```

---

### network

Network and web utilities — DNS, WHOIS, IP, CIDR, URL parsing, HTTP status codes.

#### Subcommands

| Subcommand   | Description |
|--------------|-------------|
| `dns`        | DNS lookup. `--type` A/AAAA/MX/TXT/NS/CNAME/SOA/SRV/PTR. |
| `ip`         | Show your public IP address |
| `whois`      | WHOIS lookup for a domain |
| `cidr`       | CIDR calculator — network, broadcast, host range, subnet mask |
| `urlparse`   | Parse URL into components (scheme, host, port, path, query, fragment) |
| `httpstatus` | Look up HTTP status code meaning |
| `useragent`  | Parse User-Agent string (browser, OS, device type) |
| `jwt`        | Decode JWT token header and payload |

```bash
openGyver network dns example.com --type MX
openGyver network ip
openGyver network whois example.com
openGyver network cidr 192.168.1.0/24
openGyver network urlparse "https://example.com/path?q=1&r=2#frag"
openGyver network httpstatus 404
openGyver network useragent "Mozilla/5.0 ..."
```

---

### number

Number base conversions and Roman numerals.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `base`     | Convert between bases 2-36. `--from` (default 10), `--to` (required). |
| `roman`    | Auto-detects direction: decimal to Roman or Roman to decimal. |
| `ieee754`  | Show IEEE 754 float32/float64 representation (sign, exponent, mantissa). |

```bash
openGyver number base 255 --to 16            # ff
openGyver number base ff --from 16 --to 2    # 11111111
openGyver number roman 42                    # XLII
openGyver number roman XLII                  # 42
openGyver number ieee754 3.14
```

---

### regex

Regular expression tools — test, replace, extract.

#### Subcommands

| Subcommand | Args | Description |
|------------|------|-------------|
| `test`     | `<pattern> <input>` | Test regex, show match result and groups. `--global`/`-g` for all matches. |
| `replace`  | `<pattern> <replacement> <input>` | Regex find-and-replace. Supports `$1`, `$2` groups. |
| `extract`  | `<pattern> <input>` | Extract all matches, one per line. `--file`/`-f`. |

```bash
openGyver regex test "\d+" "order 42 has 3 items"
openGyver regex test --global "\d+" "order 42 has 3 items"
openGyver regex replace "\d+" "X" "order 42 has 3 items"
openGyver regex extract "\w+@\w+\.\w+" --file emails.txt
```

---

### testdata

Generate fake/test data for development.

#### Subcommands

| Subcommand | Key Flags | Description |
|------------|-----------|-------------|
| `person`   | `--count` (1) | Fake name, email, phone, address, age |
| `csv`      | `--rows` (10), `--columns` (name,email,age,city) | Sample CSV. Column types: name, email, number, date, bool, city, country, age, phone, uuid |
| `json`     | `--count` (5) | Sample JSON objects with id, name, email, age, active, created_at |
| `number`   | `--min` (0), `--max` (100), `--count` (1), `--float` | Random numbers |

```bash
openGyver testdata person --count 5
openGyver testdata person --count 3 -j
openGyver testdata csv --rows 100 --columns name,email,age,city,phone
openGyver testdata json --count 20
openGyver testdata number --min 1 --max 1000 --count 10
```

---

### text

Text manipulation — counting, case conversion, sorting, diffing, and more.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `count`    | Word, character, line, sentence count |
| `case`     | Case converter: `--to` upper/lower/title/sentence/camel/snake/kebab/pascal/constant/dot |
| `reverse`  | Reverse a string |
| `sort`     | Sort lines: `--by` alpha/length/numeric, `--reverse` |
| `dedupe`   | Remove duplicate lines |
| `slug`     | Generate URL slug |
| `lorem`    | Lorem Ipsum generator: `--words`, `--sentences`, `--paragraphs` |
| `diff`     | Text diff between two files: `--file1`, `--file2` |
| `wrap`     | Word wrap: `--width` (default 80) |
| `lines`    | Add line numbers |
| `trim`     | Strip whitespace, `--blank` to remove blank lines |
| `replace`  | Find and replace: `--find`, `--replace`, `--regex` |

```bash
openGyver text count "hello world foo bar"
openGyver text count --file document.txt
openGyver text case "hello world" --to snake       # hello_world
openGyver text case "helloWorld" --to kebab         # hello-world
openGyver text slug "Hello World! This is a Test"   # hello-world-this-is-a-test
openGyver text lorem --sentences 3
openGyver text sort --file list.txt --by length --reverse
openGyver text replace --find "old" --replace "new" --file doc.txt
```

---

### toIco

Convert an image to a Windows ICO file with multiple embedded sizes.

#### Usage

```
openGyver toIco <image> [flags]
```

#### Arguments

| Argument  | Required | Description                               |
|-----------|----------|-------------------------------------------|
| `<image>` | Yes      | Path to the source image (PNG, JPEG, BMP, etc.) |

#### Flags

| Flag       | Short | Type     | Default          | Description                                      |
|------------|-------|----------|------------------|--------------------------------------------------|
| `--output` | `-o`  | string   | `output.ico`     | Output file path                                 |
| `--sizes`  |       | ints     | `16,32,48,256`   | Comma-separated icon sizes to embed              |
| `--format` | `-f`  | string   |                  | Override input format detection                  |
| `--width`  |       | int      | 0                | Resize source to this width before generating ICO |
| `--height` |       | int      | 0                | Resize source to this height before generating ICO |

#### Examples

```bash
openGyver toIco logo.png
openGyver toIco logo.png -o favicon.ico
openGyver toIco logo.png --sizes 16,32,64
openGyver toIco photo.dat --format jpeg
openGyver toIco logo.png --width 512 --height 512
```

---

### toGif

Combine multiple images into a single animated GIF.

#### Usage

```
openGyver toGif <image> [image...] [flags]
```

#### Arguments

| Argument             | Required | Description                                       |
|----------------------|----------|---------------------------------------------------|
| `<image> [image...]` | Yes (1+) | Image files to use as frames (supports shell globs) |

#### Flags

| Flag       | Short | Type   | Default      | Description                                 |
|------------|-------|--------|--------------|---------------------------------------------|
| `--output` | `-o`  | string | `output.gif` | Output file path                            |
| `--delay`  |       | int    | 100          | Delay between frames in milliseconds        |
| `--loop`   |       | int    | 0            | Number of loops (0 = infinite)              |
| `--format` | `-f`  | string |              | Override input format detection             |
| `--width`  |       | int    | 0            | Output GIF width (0 = first frame's width)  |
| `--height` |       | int    | 0            | Output GIF height (0 = first frame's height)|

#### Examples

```bash
openGyver toGif frame1.png frame2.png frame3.png
openGyver toGif frame*.png -o animation.gif --delay 50
openGyver toGif frame*.png --loop 3
openGyver toGif *.dat --format bmp
openGyver toGif frame*.png --width 320 --height 240
```

---

### epoch

Print the current Unix epoch timestamp, or perform arithmetic on epoch values.

#### Usage

```
openGyver epoch [flags]
openGyver epoch add [flags]
openGyver epoch subtract [flags]
```

#### Flags (all subcommands)

| Flag   | Type | Default | Description              |
|--------|------|---------|--------------------------|
| `--ms` | bool | false   | Output in milliseconds   |
| `--us` | bool | false   | Output in microseconds   |
| `--ns` | bool | false   | Output in nanoseconds    |

#### Subcommand: `add`

Add a duration to an epoch and return the new epoch.

| Flag        | Type  | Default | Description                            |
|-------------|-------|---------|----------------------------------------|
| `--from`    | int64 | 0       | Starting epoch in seconds (0 = now)    |
| `--hours`   | int   | 0       | Hours to add                           |
| `--minutes` | int   | 0       | Minutes to add                         |
| `--days`    | int   | 0       | Days to add                            |
| `--weeks`   | int   | 0       | Weeks to add                           |
| `--months`  | int   | 0       | Months to add                          |
| `--years`   | int   | 0       | Years to add                           |

#### Subcommand: `subtract`

Subtract a duration from an epoch. Same flags as `add`.

#### Examples

```bash
openGyver epoch                                        # current epoch (seconds)
openGyver epoch --ms                                   # current epoch (milliseconds)
openGyver epoch --ns                                   # current epoch (nanoseconds)
openGyver epoch add --hours 2                          # now + 2 hours
openGyver epoch add --days 30                          # now + 30 days
openGyver epoch add --days 30 --from 1705334400        # specific epoch + 30 days
openGyver epoch add --years 1 --months 6 --days 15     # combine durations
openGyver epoch subtract --years 1                     # now - 1 year
openGyver epoch subtract --days 7 --hours 12           # now - 7 days 12 hours
openGyver epoch subtract --months 3 --from 1705334400  # specific epoch - 3 months
openGyver epoch add --weeks 2 --ms                     # output in milliseconds
```

---

### timex

Convert, format, and manipulate dates, times, timezones, and Unix epochs.

#### Usage

```
openGyver timex <subcommand> [flags]
```

#### Input Formats (auto-detected)

| Format                | Example                                   |
|-----------------------|-------------------------------------------|
| ISO 8601 / RFC 3339   | `2024-01-15T14:30:00Z`, `2024-01-15T14:30:00+05:30` |
| RFC 2822              | `Mon, 15 Jan 2024 14:30:00 +0000`        |
| RFC 850               | `Monday, 15-Jan-24 14:30:00 UTC`          |
| Date only             | `2024-01-15`, `01/15/2024`, `Jan 15, 2024` |
| Date + time           | `2024-01-15 14:30:00`, `2024-01-15 14:30` |
| 12-hour               | `2024-01-15 2:30 PM`, `Jan 15, 2024 2:30:00 PM` |
| Unix timestamp        | `1705334400` (auto-detects s/ms/us/ns)    |
| Relative              | `now`, `today`, `yesterday`, `tomorrow`   |

#### Timezone Format

IANA names (`America/New_York`, `Europe/London`, `Asia/Tokyo`) or common abbreviations (`UTC`, `EST`, `PST`, `JST`, `CET`, `IST`, and 28 others).

---

#### timex now

Show the current time in multiple formats.

| Flag       | Type   | Default | Description                              |
|------------|--------|---------|------------------------------------------|
| `--tz`     | string |         | Timezone to display                      |
| `--format` | string |         | Output a single named format             |

```bash
openGyver timex now
openGyver timex now --tz Asia/Tokyo
openGyver timex now --tz EST
openGyver timex now --format iso8601
openGyver timex now --tz Europe/London --format rfc2822
```

---

#### timex to-utc

Convert a time string to UTC.

| Flag     | Type   | Default | Description                           |
|----------|--------|---------|---------------------------------------|
| `--from` | string | local   | Source timezone for naive inputs       |

```bash
openGyver timex to-utc "2024-01-15T14:30:00-05:00"
openGyver timex to-utc "2024-01-15 14:30" --from America/New_York
openGyver timex to-utc "Jan 15, 2024 2:30 PM" --from PST
openGyver timex to-utc now
openGyver timex to-utc 1705334400
```

---

#### timex to-tz

Convert a time to a target timezone.

| Flag     | Type   | Default | Description                           |
|----------|--------|---------|---------------------------------------|
| `--tz`   | string |         | Target timezone (required)            |
| `--from` | string | UTC     | Source timezone for naive inputs       |

```bash
openGyver timex to-tz "2024-01-15T14:30:00Z" --tz Asia/Tokyo
openGyver timex to-tz "2024-01-15 09:00" --from America/New_York --tz Europe/London
openGyver timex to-tz now --tz Australia/Sydney
openGyver timex to-tz 1705334400 --tz America/Chicago
```

---

#### timex to-unix

Convert a time string to a Unix epoch timestamp.

| Flag     | Type   | Default | Description                           |
|----------|--------|---------|---------------------------------------|
| `--from` | string | UTC     | Source timezone for naive inputs       |
| `--ms`   | bool   | false   | Output in milliseconds                |
| `--us`   | bool   | false   | Output in microseconds                |
| `--ns`   | bool   | false   | Output in nanoseconds                 |

```bash
openGyver timex to-unix "2024-01-15T14:30:00Z"
openGyver timex to-unix "Jan 15, 2024 2:30 PM" --from America/New_York
openGyver timex to-unix "2024-01-15" --ms
openGyver timex to-unix now --ns
```

---

#### timex from-unix

Convert a Unix epoch timestamp to human-readable time.

| Flag       | Type   | Default | Description                           |
|------------|--------|---------|---------------------------------------|
| `--tz`     | string | UTC     | Display timezone                      |
| `--ms`     | bool   | false   | Input is in milliseconds              |
| `--us`     | bool   | false   | Input is in microseconds              |
| `--ns`     | bool   | false   | Input is in nanoseconds               |
| `--format` | string |         | Output a single named format          |

Auto-detects precision if no flag is given: seconds (<1e12), milliseconds (<1e15), microseconds (<1e18), nanoseconds.

```bash
openGyver timex from-unix 1705334400
openGyver timex from-unix 1705334400000 --ms
openGyver timex from-unix 1705334400 --tz Asia/Tokyo
openGyver timex from-unix 1705334400 --format kitchen
```

---

#### timex format

Reformat a time string into a different layout. Omit `--to` to see all formats.

| Flag     | Type   | Default | Description                           |
|----------|--------|---------|---------------------------------------|
| `--to`   | string |         | Target format name or custom Go layout |
| `--from` | string |         | Source timezone for naive inputs       |

**Named formats:** `iso8601`, `rfc3339`, `rfc2822`, `rfc1123`, `rfc850`, `rfc822`, `ansic`, `unix`, `ruby`, `date`, `time`, `datetime`, `kitchen`, `us`, `eu`, `short`, `long`, `stamp`, `human`

```bash
openGyver timex format "2024-01-15T14:30:00Z" --to rfc2822
openGyver timex format "Mon, 15 Jan 2024 14:30:00 +0000" --to iso8601
openGyver timex format "2024-01-15" --to human --from America/New_York
openGyver timex format "2024-01-15T14:30:00Z" --to kitchen
openGyver timex format "2024-01-15T14:30:00Z" --to "Monday, January 2 2006"
openGyver timex format "2024-01-15T14:30:00Z"          # show all formats
```

---

#### timex diff

Calculate the duration between two date/time values.

| Flag     | Type   | Default | Description                           |
|----------|--------|---------|---------------------------------------|
| `--from` | string |         | Source timezone for naive inputs       |

Output includes: human breakdown (days/hours/minutes/seconds), total days/hours/minutes/seconds, weeks/months/years for large spans, and a calendar difference.

```bash
openGyver timex diff "2024-01-15" "2024-06-30"
openGyver timex diff "2024-01-15T08:00:00Z" "2024-01-15T17:30:00Z"
openGyver timex diff "2024-01-01" now
openGyver timex diff yesterday tomorrow
openGyver timex diff 1705334400 1710000000
```

---

#### timex add

Add or subtract a duration from a time.

| Flag     | Type   | Default | Description                           |
|----------|--------|---------|---------------------------------------|
| `--from` | string |         | Source timezone for naive time inputs  |
| `--tz`   | string |         | Display result in this timezone       |

**Duration format:** Go-style (`2h30m`, `45m`, `90s`) or extended (`30d`, `2w`, `3mo`, `1y`, `1y2mo3d4h5m6s`). Prefix with `-` to subtract.

```bash
openGyver timex add "2024-01-15T14:30:00Z" 2h30m
openGyver timex add "2024-01-15" 90d
openGyver timex add "2024-01-15" -30d
openGyver timex add now 2w
openGyver timex add "2024-03-01" 1y2mo
openGyver timex add "2024-01-15" 1y --tz America/New_York
```

---

#### timex info

Show comprehensive metadata about a date/time.

| Flag     | Type   | Default | Description                           |
|----------|--------|---------|---------------------------------------|
| `--from` | string |         | Timezone for naive inputs             |

Output includes: year, month, day, day of week, day of year, days remaining, days in month, ISO week number, quarter, leap year status, hour/minute/second, timezone offset, Unix epoch (s/ms/us/ns), and start/end of day.

```bash
openGyver timex info "2024-01-15"
openGyver timex info "2024-01-15T14:30:00Z"
openGyver timex info now
openGyver timex info 1705334400
openGyver timex info "2024-02-29" --from America/New_York
```

---

### qr

Generate a QR code from text. Prints as ASCII art in the terminal by default, or saves to PNG/SVG.

#### Usage

```
openGyver qr <text> [flags]
```

#### Arguments

| Argument | Required | Description                 |
|----------|----------|-----------------------------|
| `<text>` | Yes      | Text or URL to encode       |

#### Flags

| Flag       | Short | Type   | Default | Description                                  |
|------------|-------|--------|---------|----------------------------------------------|
| `--output` | `-o`  | string |         | Output file path (`.png` or `.svg`). Omit for ASCII output |
| `--size`   |       | int    | 256     | PNG image size in pixels                     |
| `--level`  |       | string | `L`     | Error correction: `L` (7%), `M` (15%), `Q` (25%), `H` (30%) |
| `--invert` |       | bool   | false   | Invert colors (light-on-dark for dark terminals) |

#### Examples

```bash
openGyver qr "https://example.com"
openGyver qr "Hello World" -o qr.png
openGyver qr "Hello World" -o qr.png --size 512
openGyver qr "Hello World" -o qr.svg
openGyver qr "wifi:WPA;S:MyNetwork;P:secret;;" --level H
openGyver qr "some data" --invert
```

---

### stock

Look up current or historical stock prices from global markets. Uses Yahoo Finance data — no API key required.

#### Usage

```
openGyver stock <ticker> [flags]
```

#### Arguments

| Argument   | Required | Description                           |
|------------|----------|---------------------------------------|
| `<ticker>` | Yes      | Stock ticker symbol (e.g. AAPL, MSFT) |

#### Flags

| Flag         | Short | Type   | Default | Description                                  |
|--------------|-------|--------|---------|----------------------------------------------|
| `--date`     | `-d`  | string |         | Price on a specific date (YYYY-MM-DD)        |
| `--from`     |       | string |         | Start date for range (YYYY-MM-DD)            |
| `--to`       |       | string |         | End date for range (YYYY-MM-DD)              |
| `--market`   | `-m`  | string |         | Target exchange (see market list below)      |
| `--interval` |       | string | `1d`    | Data granularity: `1d`, `1wk`, `1mo`         |

#### Supported Markets (--market flag)

| Region   | Market names                                              |
|----------|-----------------------------------------------------------|
| US       | `nasdaq`, `nyse`, `us`                                    |
| Korea    | `kosdaq`, `kospi`, `korea`                                |
| Japan    | `tokyo`, `tse`, `japan`                                   |
| UK       | `london`, `lse`, `uk`                                     |
| China    | `shanghai`, `shenzhen`, `hongkong`, `hk`                  |
| Europe   | `frankfurt`, `xetra`, `paris`, `euronext`, `amsterdam`, `swiss`, `six` |
| Nordics  | `stockholm`, `oslo`, `copenhagen`, `helsinki`              |
| Americas | `toronto`, `tsx`, `canada`, `brazil`, `bovespa`, `mexico` |
| Asia     | `singapore`, `sgx`, `taiwan`, `twse`, `mumbai`, `nse`, `bse`, `jakarta`, `bangkok` |
| Other    | `australia`, `asx`, `johannesburg`, `newzealand`          |

Tickers can also use the Yahoo Finance suffix directly: `005930.KS`, `7203.T`, `SHEL.L`.

#### Examples

```bash
# Current price
openGyver stock AAPL
openGyver stock TSLA
openGyver stock GOOGL

# Specific date
openGyver stock AAPL --date 2024-01-15
openGyver stock MSFT -d 2024-12-20

# Historical range
openGyver stock AAPL --from 2024-01-01 --to 2024-06-30
openGyver stock AAPL --from 2024-01-01 --to 2024-12-31 --interval 1wk
openGyver stock MSFT --from 2025-03-01 --to 2025-03-07

# Global markets
openGyver stock 005930 --market kospi        # Samsung (Korea)
openGyver stock 035720 --market kosdaq       # Kakao (Korea)
openGyver stock 7203 --market tokyo          # Toyota (Japan)
openGyver stock SHEL --market london         # Shell (UK)
openGyver stock 0700 --market hk             # Tencent (Hong Kong)
openGyver stock SAP --market frankfurt       # SAP (Germany)
openGyver stock MC --market paris            # LVMH (France)
openGyver stock RY --market tsx              # Royal Bank (Canada)
openGyver stock RELIANCE --market nse        # Reliance (India)
openGyver stock 2330 --market twse           # TSMC (Taiwan)

# Or use suffix directly
openGyver stock 005930.KS
openGyver stock 7203.T
```

---

### uuid

Generate universally unique identifiers.

#### Usage

```
openGyver uuid [flags]
```

#### Flags

| Flag          | Type | Default | Description                             |
|---------------|------|---------|-----------------------------------------|
| `--version`   | int  | 4       | UUID version: `4` (random) or `6` (time-sorted) |
| `--count`     | int  | 1       | Number of UUIDs to generate             |
| `--uppercase` | bool | false   | Output in uppercase                     |

**Version 4** — 122 bits of cryptographic randomness. Best for most use cases.

**Version 6** — Reordered time-based UUID. Lexicographically sortable by creation time. Good for database primary keys.

#### Examples

```bash
openGyver uuid                            # single v4 UUID
openGyver uuid --version 4                # explicit v4
openGyver uuid --version 6                # time-sorted v6
openGyver uuid --count 5                  # generate 5 UUIDs
openGyver uuid --version 6 --count 10     # 10 time-sorted UUIDs
openGyver uuid --uppercase                # A1B2C3D4-E5F6-...
```

---

### validate

Validate HTML, CSV, XML, YAML, and TOML data for correctness.

#### Subcommands

| Subcommand | Description |
|------------|-------------|
| `html`     | Check for unclosed/mismatched tags, missing doctype, missing alt on img, duplicate IDs |
| `csv`      | Check for consistent column count, proper quoting |
| `xml`      | Check XML well-formedness |
| `yaml`     | Check YAML syntax |
| `toml`     | Check TOML syntax |

All accept input as argument or `--file`/`-f`. Output "valid" on success or list of errors. `--json`/`-j` outputs `{"valid":true/false,"errors":[...]}`.

Note: JSON validation is available via `openGyver json validate`.

```bash
openGyver validate html --file index.html
openGyver validate html '<img src="x">'           # warns: missing alt
openGyver validate csv --file data.csv
openGyver validate xml '<root><item/></root>'
openGyver validate yaml --file config.yml
openGyver validate toml --file pyproject.toml
openGyver validate yaml --file config.yml -j      # JSON output
```

---

## Plugin Architecture

openGyver uses a plugin-based architecture where each command is a self-contained Go package. Adding a new command requires no changes to existing code.

### Project Structure

```
openGyver/
  main.go                             # Entrypoint — imports all plugins
  cmd/
    root.go                           # Root command + Register() function
    accessibility/                    # WCAG contrast, readability scores
    archive/                          # ZIP, TAR, TAR.GZ, 7Z, RAR
    color/                            # Color conversion, palette, contrast
    convert/                          # Unit & currency conversions
    convertaudio/                     # Audio conversions (ffmpeg)
    convertcad/                       # CAD conversions (ODA/LibreCAD)
    convertebook/                     # Ebook conversions (Calibre)
    convertfile/                      # Document & spreadsheet conversions
    convertfont/                      # Font conversions (fonttools)
    convertimage/                     # Image format conversions
    convertpresentation/              # Presentation conversions (LibreOffice)
    convertvector/                    # Vector graphics (rsvg/Inkscape)
    convertvideo/                     # Video conversions (ffmpeg)
    crypto/                           # AES, RSA, SSH, SSL certs
    dataformat/                       # YAML/TOML/XML/CSV ↔ JSON
    diff/                             # Text, JSON, CSV diff
    electrical/                       # Ohm's law, resistor codes, LED
    encode/                           # Base64, hex, URL, Morse, JWT
    epoch/                            # Unix epoch utilities
    finance/                          # Loan, compound interest, ROI
    generate/                         # Passwords, keys, IDs, OTP
    geo/                              # Distance, DMS, UTM
    hash/                             # MD5, SHA-*, HMAC, bcrypt
    jsontools/                        # JSON format, minify, validate
    network/                          # DNS, WHOIS, IP, CIDR
    number/                           # Base conversion, Roman numerals
    qr/                               # QR code generator
    regex/                            # Regex test, replace, extract
    stock/                            # Stock ticker lookup
    testdata/                         # Fake data generators
    text/                             # Text manipulation tools
    timex/                            # Time & timezone utilities
    togif/                            # Images → animated GIF
    toico/                            # Image → ICO
    uuid/                             # UUID generator
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
    yourCmd.Flags().StringVarP(&output, "output", "o", "default", "description")
    cmd.Register(yourCmd)
}
```

2. Add a blank import in `main.go`:

```go
import (
    _ "github.com/mj/opengyver/cmd/yourcommand"
)
```

3. Build and run. Your command appears in `openGyver --help`.

---

## Building from Source

```bash
git clone https://github.com/mj/opengyver.git
cd opengyver
go build -o openGyver .
```

Run all tests:

```bash
go test ./...
```

## Cross-Compilation

Go compiles to a single static binary for any platform:

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o openGyver-linux-amd64 .

# Linux (arm64 — Raspberry Pi 4, AWS Graviton)
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

Each produces a single binary with zero runtime dependencies.

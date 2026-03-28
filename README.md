# openGyver

A CLI tool for everyday conversions — images, units, currencies, documents, time, and more. Built in Go for zero-dependency, single-binary distribution across Linux, macOS, and Windows.

Designed to be used standalone, or hooked into CI/CD pipelines, shell scripts, and AI agents.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
  - [archive — Create & Extract Archives](#archive)
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
  - [toIco — Image to ICO](#toico)
  - [toGif — Images to Animated GIF](#togif)
  - [epoch — Unix Epoch Utilities](#epoch)
  - [timex — Time & Timezone Utilities](#timex)
  - [qr — QR Code Generator](#qr)
  - [uuid — UUID Generator](#uuid)
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

## Commands

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

## Plugin Architecture

openGyver uses a plugin-based architecture where each command is a self-contained Go package. Adding a new command requires no changes to existing code.

### Project Structure

```
openGyver/
  main.go                             # Entrypoint — imports all plugins
  cmd/
    root.go                           # Root command + Register() function
    archive/                          # ZIP, TAR, TAR.GZ (pure Go)
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
    toico/                            # Image → ICO
    togif/                            # Images → animated GIF
    epoch/                            # Unix epoch utilities
    timex/                            # Time & timezone utilities
    qr/                               # QR code generator
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

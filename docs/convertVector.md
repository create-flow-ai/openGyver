# convertVector

Convert vector graphics between popular formats.

## Usage

```bash
openGyver convertVector <input-file> [flags]
```

## Prerequisites

One of the following must be installed:

| Tool | Best For | Install (macOS) | Install (Linux) |
|------|----------|-----------------|-----------------|
| rsvg-convert (librsvg) | SVG to PNG/PDF (fast, lightweight) | `brew install librsvg` | `apt install librsvg2-bin` |
| Inkscape | All vector format conversions | `brew install inkscape` | `apt install inkscape` |

SVG to PNG/PDF conversions use rsvg-convert if available, or Inkscape as a fallback. Other vector format conversions require Inkscape.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input-file` | Yes | Path to the input vector file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--width` | | int | `0` | Output width in pixels (for rasterization) |
| `--height` | | int | `0` | Output height in pixels (for rasterization) |
| `--help` | `-h` | bool | `false` | Show help for convertVector |

## Supported Formats (23)

### Vector Formats

| Format | Extension(s) | Description |
|--------|-------------|-------------|
| SVG | `.svg` | Scalable Vector Graphics |
| SVGZ | `.svgz` | Compressed SVG |
| EPS | `.eps` | Encapsulated PostScript |
| PDF | `.pdf` | Portable Document Format |
| PS | `.ps` | PostScript |
| EMF | `.emf` | Enhanced Windows Metafile |
| WMF | `.wmf` | Windows Metafile |
| AI | `.ai` | Adobe Illustrator |
| CDR | `.cdr` | CorelDRAW |
| CDT | `.cdt` | CorelDRAW Template |
| CCX | `.ccx` | Corel Compressed Exchange |
| CMX | `.cmx` | Corel Metafile Exchange |
| CGM | `.cgm` | Computer Graphics Metafile |
| VSD | `.vsd` | Microsoft Visio |
| FIG | `.fig` | Xfig |
| PLT | `.plt` | HPGL Plotter |
| DST | `.dst` | Embroidery (Tajima) |
| PES | `.pes` | Embroidery (Brother) |
| EXP | `.exp` | Embroidery (Melco) |
| SK | `.sk` | Sketch |
| SK1 | `.sk1` | Sketch v1 |

### Raster Output

| Format | Extension | Description |
|--------|-----------|-------------|
| PNG | `.png` | Rasterize vector to PNG image |
| JPG | `.jpg` | Rasterize vector to JPEG image |

## Conversion Tool Selection

| Input | Output | Preferred Tool | Fallback |
|-------|--------|---------------|----------|
| SVG/SVGZ | PNG | rsvg-convert | Inkscape |
| SVG/SVGZ | PDF | rsvg-convert | Inkscape |
| SVG/SVGZ | PS/EPS | rsvg-convert | Inkscape |
| SVG/SVGZ | SVG/EMF/WMF | Inkscape | -- |
| Any other | Any | Inkscape | -- |

## Examples

```bash
# Convert SVG to PNG
openGyver convertVector logo.svg -o logo.png

# Convert SVG to PDF
openGyver convertVector logo.svg -o logo.pdf

# Convert SVG to EPS
openGyver convertVector logo.svg -o logo.eps

# Rasterize SVG to specific dimensions
openGyver convertVector logo.svg -o logo.png --width 1024 --height 768

# Convert EPS to SVG
openGyver convertVector diagram.eps -o diagram.svg

# Convert Adobe Illustrator file to SVG
openGyver convertVector drawing.ai -o drawing.svg

# Generate a high-resolution PNG from SVG
openGyver convertVector icon.svg -o icon.png --width 2048

# Convert SVG to PostScript
openGyver convertVector chart.svg -o chart.ps

# Convert Visio diagram to SVG
openGyver convertVector diagram.vsd -o diagram.svg

# Convert SVG to EMF (Windows Metafile)
openGyver convertVector logo.svg -o logo.emf

# Convert CorelDRAW file to SVG
openGyver convertVector design.cdr -o design.svg

# Rasterize with only width (height proportional)
openGyver convertVector banner.svg -o banner.png --width 1200
```

## Notes

- The input and output formats are auto-detected from file extensions.
- When converting SVG to raster (PNG/JPG), use `--width` and `--height` to control the output resolution. If only one dimension is specified, the other is determined by the tool (rsvg-convert or Inkscape) to maintain aspect ratio.
- rsvg-convert is preferred for SVG input because it is faster and more lightweight than Inkscape.
- Inkscape is required for non-SVG input formats and for output formats like EMF and WMF.
- If neither rsvg-convert nor Inkscape is found, the command will return an error with installation instructions.
- Embroidery format support (DST, PES, EXP) requires Inkscape with appropriate extensions.

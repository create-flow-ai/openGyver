# convertImage

Convert images between popular raster formats.

## Usage

```bash
openGyver convertImage <image> [flags]
```

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `image` | Yes | Path to the input image file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--quality` | | int | `90` | JPEG quality 1-100 (ignored for other formats) |
| `--width` | | int | `0` | Resize to this width in pixels (0 = keep original) |
| `--height` | | int | `0` | Resize to this height in pixels (0 = keep original) |
| `--format` | `-f` | string | `""` | Override input format detection (`png`, `jpeg`, `gif`, `bmp`, `tiff`, `webp`, `raw`) |
| `--quiet` | `-q` | bool | `false` | Suppress output messages (for piping) |
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | bool | `false` | Show help for convertImage |

## Supported Formats

### Native Go Formats

| Format | Read | Write | Extensions |
|--------|------|-------|------------|
| PNG | Yes | Yes | `.png` |
| JPEG | Yes | Yes | `.jpg`, `.jpeg`, `.jfif`, `.jpe` |
| GIF | Yes | Yes | `.gif` (first frame) |
| BMP | Yes | Yes | `.bmp` |
| TIFF | Yes | Yes | `.tiff`, `.tif` |
| WebP | Yes | No* | `.webp` |
| PPM | Yes | Yes | `.ppm`, `.pgm`, `.pbm`, `.pnm` |
| PCX | Yes | No | `.pcx` |
| TGA | Yes | No | `.tga` |
| SVG | No | Yes* | `.svg` (raster to vector via potrace/ImageMagick) |
| ICO | No | No | `.ico` (use `toIco` command instead) |

### External Tool Formats

| Format | Extensions | Requirements |
|--------|-----------|--------------|
| HEIC | `.heic`, `.heif` | dcraw, ImageMagick, or Apple sips |
| Camera RAW | `.cr2`, `.cr3`, `.crw`, `.nef`, `.arw`, `.dng`, `.orf`, `.raf`, `.rw2`, `.pef`, `.erf`, `.mrw`, `.srf`, `.sr2`, `.3fr`, `.k25`, `.kdc`, `.mef`, `.nrw`, `.x3f`, `.dcr`, `.mos`, `.iiq`, `.raw` | dcraw, ImageMagick, or Apple sips |

\* WebP can be decoded (read) but not encoded (write) in pure Go. SVG writing uses potrace or ImageMagick for raster-to-vector tracing.

## Resize Behavior

- If only `--width` is set, height scales proportionally to maintain aspect ratio.
- If only `--height` is set, width scales proportionally to maintain aspect ratio.
- If both are set, the image is resized to the exact dimensions (may change aspect ratio).
- If neither is set (both 0), the original dimensions are kept.
- Resizing uses nearest-neighbor interpolation.

## Examples

```bash
# Convert PNG to JPEG
openGyver convertImage photo.png -o photo.jpg

# Convert JPEG to PNG
openGyver convertImage photo.jpg -o photo.png

# Convert JPEG to BMP
openGyver convertImage photo.jpg -o photo.bmp

# Convert BMP to TIFF
openGyver convertImage photo.bmp -o photo.tiff

# Convert WebP to PNG (WebP decode supported)
openGyver convertImage photo.webp -o photo.png

# Convert with custom JPEG quality
openGyver convertImage photo.png -o photo.jpg --quality 85

# Create a thumbnail (width only, height auto-calculated)
openGyver convertImage photo.png -o thumb.jpg --width 200

# Resize to exact dimensions
openGyver convertImage photo.png -o resized.png --width 800 --height 600

# Override input format detection
openGyver convertImage raw.dat -o out.png --format bmp

# Quiet mode for scripting
openGyver convertImage photo.png -o photo.jpg -q

# JSON output for automation
openGyver convertImage photo.png -o photo.jpg -j

# Convert RAW camera file (requires dcraw or ImageMagick)
openGyver convertImage photo.cr2 -o photo.png

# Convert HEIC to JPEG (requires sips on macOS)
openGyver convertImage photo.heic -o photo.jpg

# Batch convert with a shell loop
for f in *.png; do openGyver convertImage "$f" -o "${f%.png}.jpg" -q; done

# Resize and convert in one step
openGyver convertImage photo.png -o thumb.jpg --width 150 --quality 75

# Pipe JSON to check dimensions
openGyver convertImage photo.png -o photo.jpg -j | jq '{width, height}'

# Convert to PPM format
openGyver convertImage photo.png -o photo.ppm

# Trace raster to SVG (requires potrace or ImageMagick)
openGyver convertImage logo.png -o logo.svg
```

## JSON Output Format

```json
{
  "success": true,
  "input": "photo.png",
  "output": "photo.jpg",
  "input_format": "png",
  "output_format": "jpeg",
  "width": 1920,
  "height": 1080
}
```

## Notes

- The input format is auto-detected from the file header (magic bytes), not the extension.
- The output format is determined by the `--output` file extension.
- Use `--format` to override input format detection when auto-detection fails.
- JPEG quality is only relevant for JPEG output; it is ignored for all other formats.
- WebP encoding is not supported in pure Go. Use `cwebp` as a workaround: convert to PNG first, then run `cwebp temp.png -o output.webp`.
- HEIC encoding is not supported in pure Go. On macOS, use `sips -s format heic input.png --out output.heic`.
- For ICO output, use the separate `toIco` command.
- Camera RAW decoding tries dcraw first, then ImageMagick (`magick` or `convert`), then Apple `sips` on macOS.

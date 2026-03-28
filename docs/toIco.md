# toIco

Convert an image to a Windows ICO file.

## Usage

```bash
openGyver toIco <image> [flags]
```

The input format is auto-detected from the file header. Use `--format` to override detection when reading from a non-standard extension or pipe.

Supports embedding multiple sizes in a single `.ico` file. Default sizes: 16, 32, 48, 256.

Use `--width` and `--height` to resize the source image before generating the ICO entries. If only one dimension is given, the image is scaled proportionally. If neither is given, the original dimensions are used.

**Supported input formats:** png, jpeg, gif, bmp, tiff, webp.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `image` | Yes | The source image file to convert |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--format` | `-f` | string | auto-detected | Input image format (png, jpeg, gif, bmp, tiff, webp) |
| `--height` | | int | `0` | Resize source to this height before generating ICO (0 = keep original) |
| `--width` | | int | `0` | Resize source to this width before generating ICO (0 = keep original) |
| `--output` | `-o` | string | `output.ico` | Output file path |
| `--sizes` | | ints | `[16,32,48,256]` | Icon sizes to embed in the ICO file |
| `--help` | `-h` | | | Help for toIco |

## Examples

```bash
# Convert a PNG to ICO with default sizes (16, 32, 48, 256)
openGyver toIco logo.png

# Specify output file name
openGyver toIco logo.png -o favicon.ico

# Custom icon sizes
openGyver toIco logo.png --sizes 16,32,64

# Override input format for a non-standard extension
openGyver toIco photo.dat --format jpeg

# Resize source image before generating ICO
openGyver toIco logo.png --width 512 --height 512

# Create a favicon with specific sizes for web
openGyver toIco logo.png -o favicon.ico --sizes 16,32,48,64,128,256

# Create a small icon set
openGyver toIco icon.png --sizes 16,32 -o small-icon.ico

# Convert a JPEG photograph to ICO
openGyver toIco headshot.jpg -o profile.ico --sizes 32,64,128

# Create ICO from a WebP source
openGyver toIco design.webp -o app-icon.ico

# Resize and create a single-size ICO
openGyver toIco large-logo.png --width 256 --height 256 --sizes 256 -o logo.ico
```

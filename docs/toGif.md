# toGif

Combine multiple images into an animated GIF.

## Usage

```bash
openGyver toGif <image> [image...] [flags]
```

Images are added as frames in the order given. The input format is auto-detected from each file's header. Use `--format` to override detection when all inputs share the same non-standard extension.

Use `--width` and `--height` to set the output GIF dimensions. If only one dimension is given, the other is scaled proportionally from the first frame. If neither is given, the first frame's dimensions are used.

**Supported input formats:** png, jpeg, gif, bmp, tiff, webp.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `image` | Yes (at least one) | One or more image files to combine as frames |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--delay` | | int | `100` | Delay between frames in milliseconds |
| `--format` | `-f` | string | auto-detected | Input image format (png, jpeg, gif, bmp, tiff, webp) |
| `--height` | | int | `0` | Output GIF height in pixels (0 = use first frame's height) |
| `--width` | | int | `0` | Output GIF width in pixels (0 = use first frame's width) |
| `--loop` | | int | `0` | Number of loops (0 = infinite) |
| `--output` | `-o` | string | `output.gif` | Output file path |
| `--help` | `-h` | | | Help for toGif |

## Examples

```bash
# Create a GIF from three PNG frames
openGyver toGif frame1.png frame2.png frame3.png

# Use a glob pattern with custom output name and faster animation
openGyver toGif frame*.png -o animation.gif --delay 50

# Limit to 3 loops (then stop)
openGyver toGif frame*.png --loop 3

# Override format detection for non-standard file extensions
openGyver toGif *.dat --format bmp

# Set specific output dimensions
openGyver toGif frame*.png --width 320 --height 240

# Create a slow slideshow (1 second between frames)
openGyver toGif photo1.jpg photo2.jpg photo3.jpg --delay 1000 -o slideshow.gif

# Create a very fast animation (30ms between frames)
openGyver toGif sprite*.png -o sprite-animation.gif --delay 30

# Scale down a large animation
openGyver toGif screenshot*.png -o demo.gif --width 640

# Create a looping animation from TIFF frames
openGyver toGif render*.tiff -o render.gif --delay 80

# Create a GIF from mixed format images
openGyver toGif cover.png frame1.jpeg frame2.png -o mixed.gif
```

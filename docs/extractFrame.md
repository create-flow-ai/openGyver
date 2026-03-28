# extractFrame

Extract a specific frame or all frames from an animated GIF or APNG as PNG images.

## Usage

```bash
openGyver extractFrame <animated-image> [flags]
```

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `<animated-image>` | Yes | Path to an animated GIF (`.gif`) or APNG (`.png`) |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--frame` | | int | 0 | Frame number to extract (0-indexed) |
| `--all` | | bool | false | Extract all frames as individual PNG files |
| `--output` | `-o` | string | `input_frameN.png` | Output file path |
| `--quiet` | `-q` | bool | false | Suppress output messages |
| `--json` | `-j` | bool | false | Output as JSON |

## Supported Formats

| Format | Single Frame | All Frames | Engine |
|--------|-------------|------------|--------|
| GIF | yes | yes | Pure Go (image/gif) |
| APNG | frame 0 only | no | Go stdlib + ffmpeg fallback |

## Examples

```bash
# Extract first frame (default)
openGyver extractFrame animation.gif

# Extract frame 5
openGyver extractFrame animation.gif --frame 5

# Extract frame to custom path
openGyver extractFrame animation.gif --frame 5 -o keyframe.png

# Extract ALL frames as individual PNGs
openGyver extractFrame animation.gif --all
# Creates: animation_000.png, animation_001.png, animation_002.png, ...

# Extract all frames to a specific directory
openGyver extractFrame animation.gif --all -o frames/frame.png
# Creates: frames/frame_000.png, frames/frame_001.png, ...

# Extract from APNG (frame 0)
openGyver extractFrame animated.png --frame 0

# JSON output
openGyver extractFrame animation.gif --frame 3 -j

# Quiet mode
openGyver extractFrame animation.gif --frame 0 -o thumb.png -q
```

## JSON Output Format

### Single Frame

```json
{
  "success": true,
  "input": "animation.gif",
  "output": "animation_frame3.png",
  "frame": 3,
  "total_frames": 24
}
```

### All Frames

```json
{
  "success": true,
  "input": "animation.gif",
  "total_frames": 24,
  "files": [
    "animation_000.png",
    "animation_001.png",
    "animation_002.png"
  ]
}
```

## Notes

- **GIF extraction is pure Go** — no external tools needed. Works on any platform.
- **APNG frame 0** uses Go's standard `image/png` decoder (which returns the default image). For other APNG frames, ffmpeg is required.
- Frame numbers are **0-indexed**. A GIF with 10 frames has frames 0-9.
- All output frames are saved as **PNG** regardless of input format.
- When using `--all`, frames are numbered with 3-digit zero-padding (000, 001, ...).

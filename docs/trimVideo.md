# trimVideo

Trim a video file to a specific time range. By default copies streams without re-encoding (instant, lossless).

## Usage

```bash
openGyver trimVideo <input-file> [flags]
```

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `<input-file>` | Yes | Path to the video file to trim |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--start` | `-s` | string | | Start time (HH:MM:SS, HH:MM:SS.ms, or seconds) |
| `--end` | `-e` | string | | End time (alternative to `--duration`) |
| `--duration` | `-d` | string | | Duration to extract (alternative to `--end`) |
| `--output` | `-o` | string | `input_trimmed.ext` | Output file path |
| `--codec` | | string | `copy` | Video codec. `copy` = no re-encoding (instant). Use `libx264`, `libx265`, etc. to re-encode |
| `--quiet` | `-q` | bool | false | Suppress output messages |
| `--json` | `-j` | bool | false | Output as JSON |

## Requirements

Requires `ffmpeg` installed and in PATH.

```bash
# macOS
brew install ffmpeg

# Linux
apt install ffmpeg
```

## Time Format

| Format | Example | Description |
|--------|---------|-------------|
| Seconds | `90` | 1 minute 30 seconds |
| HH:MM:SS | `00:01:30` | 1 minute 30 seconds |
| HH:MM:SS.ms | `00:01:30.500` | 1 minute 30.5 seconds |

## Examples

```bash
# Trim from 1:00 to 2:30
openGyver trimVideo input.mp4 -s 00:01:00 -e 00:02:30

# Take 60 seconds starting at 30s mark
openGyver trimVideo input.mp4 -s 30 -d 60 -o clip.mp4

# Extract a scene from a movie, convert to MP4
openGyver trimVideo movie.mkv -s 00:10:00 -e 00:15:00 -o scene.mp4

# First 10 seconds with re-encoding
openGyver trimVideo input.mp4 -s 0 -d 10 --codec libx264

# First 10 seconds, keep original format (default: stream copy)
openGyver trimVideo input.mp4 -d 10

# Trim and convert AVI to WebM
openGyver trimVideo input.avi -s 00:05:00 -e 00:10:00 -o clip.webm

# Quiet mode for scripting
openGyver trimVideo input.mp4 -s 0 -d 30 -o intro.mp4 -q

# JSON output
openGyver trimVideo input.mp4 -s 60 -d 30 -j
```

## JSON Output Format

```json
{
  "success": true,
  "input": "input.mp4",
  "output": "clip.mp4",
  "start": "00:01:00",
  "end": "00:02:30",
  "duration": "",
  "codec": "copy"
}
```

## Notes

- **Default mode is stream copy** (`-c copy`) — no re-encoding, instant, lossless. The output format matches the input.
- When using stream copy, cut points may not be frame-exact (keyframe-aligned). Use `--codec libx264` for frame-exact cuts.
- The output format is determined by the `-o` file extension. Use this to convert during trimming (e.g., `.mkv` input → `.mp4` output).
- `--start` + `--end` and `--start` + `--duration` are two ways to specify the range. Don't use both `--end` and `--duration`.

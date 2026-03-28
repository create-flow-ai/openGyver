# convertVideo

Convert video files between popular formats using ffmpeg.

## Usage

```bash
openGyver convertVideo <input-file> [flags]
```

## Prerequisites

ffmpeg must be installed and available in PATH.

| Platform | Install Command |
|----------|----------------|
| macOS | `brew install ffmpeg` |
| Linux | `apt install ffmpeg` / `dnf install ffmpeg` |
| Windows | [ffmpeg.org/download.html](https://ffmpeg.org/download.html) |

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input-file` | Yes | Path to the input video file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--resolution` | | string | `""` | Output resolution (e.g., `1920x1080`, `1280x720`) |
| `--vbitrate` | | string | `""` | Video bitrate (e.g., `2M`, `5M`, `10M`) |
| `--abitrate` | | string | `""` | Audio bitrate (e.g., `128k`, `192k`, `320k`) |
| `--fps` | | string | `""` | Frames per second (e.g., `24`, `30`, `60`) |
| `--codec` | | string | `""` | Video codec (e.g., `libx264`, `libx265`, `libvpx-vp9`) |
| `--quiet` | `-q` | bool | `false` | Suppress output messages (for piping) |
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | bool | `false` | Show help for convertVideo |

## Supported Formats (37)

3G2, 3GP, 3GPP, ASF, AV1, AVI, CAVS, DIVX, DV, DVR, F4V, FLV, HEVC, M2TS, M2V, M4V, MJPEG, MKV, MOD, MOV, MP4, MPEG, MPG, MTS, MXF, OGG, OGV, RM, RMVB, SWF, TOD, TS, VOB, WEBM, WMV, WTV, XVID

Additionally, GIF is supported as an output format for creating animated GIFs from video.

## Common Video Codecs

| Codec | Description |
|-------|-------------|
| `libx264` | H.264/AVC (widely compatible, good quality) |
| `libx265` | H.265/HEVC (better compression, newer) |
| `libvpx-vp9` | VP9 (good for WebM) |
| `libaom-av1` | AV1 (best compression, slow encoding) |
| `mpeg4` | MPEG-4 Part 2 (legacy) |

## Examples

```bash
# Convert AVI to MP4
openGyver convertVideo input.avi -o output.mp4

# Convert MKV to WebM
openGyver convertVideo input.mkv -o output.webm

# Convert with specific resolution
openGyver convertVideo input.mov -o output.mp4 --resolution 1920x1080

# Set video and audio bitrate
openGyver convertVideo input.mp4 -o output.mp4 --vbitrate 5M --abitrate 192k

# Create an animated GIF from video
openGyver convertVideo input.mp4 -o output.gif --fps 10

# Use a specific video codec
openGyver convertVideo input.mp4 -o output.webm --codec libvpx-vp9

# Downscale to 720p
openGyver convertVideo input.mp4 -o output.mp4 --resolution 1280x720

# Convert with H.265 for smaller file size
openGyver convertVideo input.mov -o output.mp4 --codec libx265

# Convert old format to modern MP4
openGyver convertVideo input.wmv -o output.mp4

# High-quality conversion with all options
openGyver convertVideo input.mov -o output.mp4 --resolution 1920x1080 --vbitrate 8M --abitrate 256k --fps 30 --codec libx264

# Quiet mode for scripting
openGyver convertVideo input.avi -o output.mp4 -q

# JSON output for automation
openGyver convertVideo input.avi -o output.mp4 -j

# Convert VOB (DVD) to MP4
openGyver convertVideo input.vob -o output.mp4

# Convert MTS (camera footage) to MP4
openGyver convertVideo footage.mts -o footage.mp4

# Pipe JSON output to jq
openGyver convertVideo input.avi -o output.mp4 -j | jq '.output_format'

# Batch convert with a shell loop
for f in *.avi; do openGyver convertVideo "$f" -o "${f%.avi}.mp4" -q; done

# Create low-fps GIF for preview
openGyver convertVideo clip.mp4 -o preview.gif --fps 5 --resolution 320x240
```

## JSON Output Format

```json
{
  "success": true,
  "input": "input.avi",
  "output": "output.mp4",
  "output_format": "mp4",
  "resolution": "1920x1080",
  "video_bitrate": "5M",
  "audio_bitrate": "192k",
  "fps": "30",
  "codec": "libx264"
}
```

## Notes

- The output format is determined by the file extension of the `--output` path.
- If optional flags (`--resolution`, `--vbitrate`, `--abitrate`, `--fps`, `--codec`) are not specified, ffmpeg uses sensible defaults for the output format.
- The `-y` flag is passed to ffmpeg automatically, so existing output files will be overwritten without prompting.
- GIF is accepted as an output format even though it is not in the main supported formats list. Use `--fps` to control the frame rate of the resulting GIF.
- Resolution is specified as `WIDTHxHEIGHT` (e.g., `1920x1080`).
- Video bitrate uses ffmpeg notation: `2M` = 2 Mbps, `500k` = 500 Kbps.
- Audio bitrate uses ffmpeg notation: `128k` = 128 Kbps, `320k` = 320 Kbps.

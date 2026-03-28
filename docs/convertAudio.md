# convertAudio

Convert audio files between popular formats using ffmpeg.

## Usage

```bash
openGyver convertAudio <input-file> [flags]
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
| `input-file` | Yes | Path to the input audio file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--bitrate` | | string | `""` | Audio bitrate (e.g., `128k`, `192k`, `320k`) |
| `--sample` | | string | `""` | Sample rate in Hz (e.g., `44100`, `48000`) |
| `--quiet` | `-q` | bool | `false` | Suppress output messages (for piping) |
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | bool | `false` | Show help for convertAudio |

## Supported Formats (33)

AAC, AC3, AIF, AIFC, AIFF, AMR, APE, AU, CAF, DSS, DTS, FLAC, GSM, M4A, M4B, M4R, MP2, MP3, OGA, OGG, OPUS, RA, SHN, SND, SPX, TTA, VOC, VOX, W64, WAV, WEBA, WMA, WV

## Examples

```bash
# Convert WAV to MP3
openGyver convertAudio song.wav -o song.mp3

# Convert FLAC to AAC
openGyver convertAudio song.flac -o song.aac

# Convert MP3 to OGG
openGyver convertAudio podcast.mp3 -o podcast.ogg

# Convert with a specific bitrate
openGyver convertAudio song.wav -o song.mp3 --bitrate 320k

# Convert with a specific sample rate
openGyver convertAudio song.mp3 -o song.wav --sample 44100

# Convert with both bitrate and sample rate
openGyver convertAudio recording.flac -o recording.mp3 --bitrate 192k --sample 48000

# Quiet mode for scripting/piping
openGyver convertAudio song.wav -o song.mp3 -q

# JSON output for programmatic use
openGyver convertAudio song.wav -o song.mp3 -j

# Convert lossless to lossless
openGyver convertAudio song.wav -o song.flac

# Convert to Apple-friendly format
openGyver convertAudio podcast.mp3 -o podcast.m4a

# Convert to Opus (modern efficient codec)
openGyver convertAudio song.flac -o song.opus --bitrate 128k

# Pipe JSON output to jq
openGyver convertAudio song.wav -o song.mp3 -j | jq '.output_format'
```

## JSON Output Format

```json
{
  "success": true,
  "input": "song.wav",
  "output": "song.mp3",
  "output_format": "mp3",
  "bitrate": "320k",
  "sample_rate": "44100"
}
```

## Notes

- The output format is determined by the file extension of the `--output` path.
- If `--bitrate` is not specified, ffmpeg uses its default bitrate for the output format.
- If `--sample` is not specified, the original sample rate is preserved.
- The `-y` flag is passed to ffmpeg automatically, so existing output files will be overwritten without prompting.

# qr

Generate a QR code from text and display it in the terminal or save it as an image file.

## Usage

```bash
openGyver qr <text> [flags]
```

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Output file path (.png or .svg). Omit for ASCII terminal output |
| `--size` | | int | `256` | PNG image size in pixels |
| `--level` | | string | `L` | Error correction level: L, M, Q, H |
| `--invert` | | bool | `false` | Invert colors (light-on-dark for dark terminals) |
| `--help` | `-h` | | | Help for qr |

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | Yes | The text or data to encode in the QR code |

## Output Modes

| Mode | Description |
|------|-------------|
| Default (no `-o`) | Prints QR code as Unicode block characters in the terminal |
| `-o file.png` | Saves as PNG image (use `--size` to set pixel dimensions) |
| `-o file.svg` | Saves as SVG vector graphic |

## Error Correction Levels

| Level | Recovery | Description |
|-------|----------|-------------|
| `L` | 7% | Default -- smallest QR code |
| `M` | 15% | Medium recovery |
| `Q` | 25% | Higher recovery |
| `H` | 30% | Maximum recovery, largest QR code, most resilient |

## Examples

```bash
# Display a URL as ASCII QR code in the terminal
openGyver qr "https://example.com"

# Save as PNG image
openGyver qr "Hello World" -o qr.png

# Save as PNG with custom size (512x512 pixels)
openGyver qr "Hello World" -o qr.png --size 512

# Save as SVG vector graphic
openGyver qr "Hello World" -o qr.svg

# Encode Wi-Fi credentials with maximum error correction
openGyver qr "wifi:WPA;S:MyNetwork;P:secret;;" --level H

# Invert colors for dark terminal backgrounds
openGyver qr "some data" --invert

# Generate a large, high-quality PNG
openGyver qr "https://myapp.com/download" -o download-qr.png --size 1024 --level H

# Encode a vCard contact
openGyver qr "BEGIN:VCARD\nVERSION:3.0\nN:Doe;John\nTEL:+1234567890\nEND:VCARD" -o contact.png

# Generate a small QR code for a short text
openGyver qr "HELLO" -o small.png --size 128 --level L

# Save as SVG for print (scales perfectly)
openGyver qr "https://docs.example.com" -o print-qr.svg --level M
```

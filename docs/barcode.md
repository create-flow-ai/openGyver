# barcode

Generate barcodes and 2D codes as PNG images.

## Usage

```bash
openGyver barcode [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Output PNG file path (required) |
| `--width` | | int | `300` | Image width in pixels |
| `--height` | | int | `100` | Image height in pixels |
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | | | Help for barcode |

## Subcommands

### code128

Generate a Code 128 barcode as a PNG image. Code 128 is a high-density linear barcode that encodes all 128 ASCII characters. It is widely used in shipping, packaging, and supply chain management.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `data` | Yes | The data to encode in the barcode |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Generate a basic Code 128 barcode
openGyver barcode code128 "Hello World" -o code128.png

# Custom dimensions
openGyver barcode code128 "ABC-12345" -o code128.png --width 400 --height 120

# Encode a product code
openGyver barcode code128 "SKU-98765" -o sku.png

# JSON output with barcode generation
openGyver barcode code128 "Hello123" -o barcode.png --json

# Generate a shipping label barcode
openGyver barcode code128 "SHIP-2024-001" -o shipping.png --width 500 --height 150

# Compact barcode
openGyver barcode code128 "XYZ" -o small.png --width 200 --height 80
```

#### JSON Output Format

```json
{
  "format": "code128",
  "data": "Hello123",
  "output": "barcode.png",
  "width": 300,
  "height": 100
}
```

---

### ean13

Generate an EAN-13 barcode as a PNG image. EAN-13 (European Article Number) is a 13-digit barcode used worldwide for marking retail goods. Provide 12 digits (the check digit is computed automatically) or 13 digits (the check digit is validated).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `digits` | Yes | 12 or 13 digits for the EAN-13 barcode |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Generate with 12 digits (check digit auto-computed)
openGyver barcode ean13 "590123412345" -o ean13.png

# Generate with 13 digits (check digit validated)
openGyver barcode ean13 "5901234123457" -o ean13.png

# Custom size
openGyver barcode ean13 "590123412345" -o ean.png --width 400 --height 150

# JSON output
openGyver barcode ean13 "590123412345" -o ean.png --json

# Generate for a retail product
openGyver barcode ean13 "400399415136" -o product.png

# Wide barcode for label printing
openGyver barcode ean13 "978020137962" -o isbn.png --width 500 --height 200
```

#### JSON Output Format

```json
{
  "format": "ean13",
  "data": "590123412345",
  "output": "ean13.png",
  "width": 300,
  "height": 100
}
```

---

### qr

Generate a QR code as a PNG image using the boombuler/barcode library. This is an alternative QR code generator that always outputs to a PNG file. For square QR codes, use equal values for `--width` and `--height`. Uses medium error correction level (M).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `data` | Yes | The data to encode in the QR code |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Encode a URL
openGyver barcode qr "https://example.com" -o qr.png

# Square QR code with custom size
openGyver barcode qr "Hello" -o qr.png --width 512 --height 512

# Encode contact info
openGyver barcode qr "tel:+1234567890" -o phone-qr.png --width 256 --height 256

# JSON output
openGyver barcode qr "https://example.com" -o qr.png --json

# Encode plain text
openGyver barcode qr "Meeting at 3pm" -o note.png --width 300 --height 300

# Encode Wi-Fi credentials
openGyver barcode qr "WIFI:S:MyNetwork;T:WPA;P:password123;;" -o wifi.png --width 400 --height 400
```

#### JSON Output Format

```json
{
  "format": "qr",
  "data": "https://example.com",
  "output": "qr.png",
  "width": 512,
  "height": 512
}
```

---

### datamatrix

Generate a DataMatrix 2D barcode as a PNG image. DataMatrix is a two-dimensional barcode that can store large amounts of data in a small space. It is commonly used in electronics, healthcare, and logistics for marking small items.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `data` | Yes | The data to encode in the DataMatrix barcode |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Generate a DataMatrix barcode
openGyver barcode datamatrix "SERIAL-00123" -o dm.png

# Custom dimensions
openGyver barcode datamatrix "Hello World" -o dm.png --width 200 --height 200

# Encode a part number
openGyver barcode datamatrix "P/N:ABC-123-XYZ" -o part.png --width 150 --height 150

# JSON output
openGyver barcode datamatrix "DATA" -o dm.png --json

# Encode tracking data
openGyver barcode datamatrix "LOT:2024-A1-B2" -o lot.png --width 300 --height 300

# Small DataMatrix for component marking
openGyver barcode datamatrix "SN001" -o sn.png --width 100 --height 100
```

#### JSON Output Format

```json
{
  "format": "datamatrix",
  "data": "SERIAL-00123",
  "output": "dm.png",
  "width": 200,
  "height": 200
}
```

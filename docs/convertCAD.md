# convertCAD

Convert CAD files between DWG, DXF, and other formats.

## Usage

```bash
openGyver convertCAD <input-file> [flags]
```

## Prerequisites

One of the following must be installed:

| Tool | Use Case | Install |
|------|----------|---------|
| LibreCAD | DXF conversions (DXF to PDF, SVG, PNG) | `brew install librecad` (macOS) / `apt install librecad` (Linux) |
| ODA File Converter | DWG to DXF and DXF to DWG | [opendesign.com/guestfiles/oda_file_converter](https://www.opendesign.com/guestfiles/oda_file_converter) |

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input-file` | Yes | Path to the input CAD file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--help` | `-h` | bool | `false` | Show help for convertCAD |

## Supported Formats

| Format | Description |
|--------|-------------|
| DWG | AutoCAD Drawing |
| DXF | Drawing Exchange Format |
| DWF | Design Web Format |
| PDF | Export to PDF |
| SVG | Export to SVG |
| PNG | Export to PNG |

## Conversion Paths

| From | To | Tool Used |
|------|----|-----------|
| DWG | DXF | ODA File Converter |
| DXF | DWG | ODA File Converter |
| DXF | PDF | LibreCAD |
| DXF | SVG | LibreCAD |
| DXF | PNG | LibreCAD |
| DWG | PDF | ODA (DWG->DXF) + LibreCAD (DXF->PDF) |

## Examples

```bash
# Convert DWG to DXF (requires ODA File Converter)
openGyver convertCAD drawing.dwg -o drawing.dxf

# Convert DXF to PDF (requires LibreCAD)
openGyver convertCAD drawing.dxf -o drawing.pdf

# Convert DXF to SVG
openGyver convertCAD drawing.dxf -o drawing.svg

# Convert DXF to PNG
openGyver convertCAD drawing.dxf -o drawing.png

# Convert DWG to PDF
openGyver convertCAD drawing.dwg -o drawing.pdf

# Convert DXF back to DWG
openGyver convertCAD schematic.dxf -o schematic.dwg

# Convert floor plan to web-friendly SVG
openGyver convertCAD floorplan.dxf -o floorplan.svg

# Export mechanical drawing as PDF for review
openGyver convertCAD part.dxf -o part.pdf
```

## Notes

- The input and output formats are auto-detected from file extensions.
- DWG to DXF conversions (and vice versa) use the ODA File Converter with the ACAD2018 output version.
- DXF to PDF/SVG/PNG conversions use LibreCAD's `--export` flag.
- If neither ODA File Converter nor LibreCAD is found, the command will return an error with installation instructions.
- For DWG to PDF, you may need to first convert DWG to DXF (via ODA), then DXF to PDF (via LibreCAD).

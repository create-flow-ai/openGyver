# convertPresentation

Convert presentation files between popular formats using LibreOffice.

## Usage

```bash
openGyver convertPresentation <input-file> [flags]
```

## Prerequisites

LibreOffice must be installed (specifically the `soffice` command-line tool).

| Platform | Install Command |
|----------|----------------|
| macOS | `brew install --cask libreoffice` |
| Linux | `apt install libreoffice` |
| Windows | [libreoffice.org/download](https://www.libreoffice.org/download) |

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input-file` | Yes | Path to the input presentation file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--help` | `-h` | bool | `false` | Show help for convertPresentation |

## Supported Formats

### Input Formats (10)

| Format | Description |
|--------|-------------|
| DPS | WPS Presentation |
| KEY | Apple Keynote |
| ODP | OpenDocument Presentation |
| POT | PowerPoint Template (legacy) |
| POTX | PowerPoint Template |
| PPS | PowerPoint Show (legacy) |
| PPSX | PowerPoint Show |
| PPT | PowerPoint (legacy) |
| PPTM | PowerPoint with Macros |
| PPTX | PowerPoint |

### Output Formats (8)

| Format | Description |
|--------|-------------|
| ODP | OpenDocument Presentation |
| PDF | Portable Document Format |
| PPTX | PowerPoint |
| PPT | PowerPoint (legacy) |
| HTML | Web page |
| PNG | Image (one per slide) |
| JPG | Image (one per slide) |
| SVG | Scalable Vector Graphics |

## Examples

```bash
# Convert PowerPoint to PDF
openGyver convertPresentation slides.pptx -o slides.pdf

# Convert PowerPoint to OpenDocument
openGyver convertPresentation slides.pptx -o slides.odp

# Convert Keynote to PowerPoint
openGyver convertPresentation slides.key -o slides.pptx

# Convert PowerPoint to HTML
openGyver convertPresentation slides.pptx -o slides.html

# Upgrade legacy PPT to modern PPTX
openGyver convertPresentation old.ppt -o new.pptx

# Export slides as images
openGyver convertPresentation slides.pptx -o slides.png

# Export slides as SVG
openGyver convertPresentation slides.pptx -o slides.svg

# Convert a PowerPoint template to PDF
openGyver convertPresentation template.potx -o template.pdf

# Convert PowerPoint show to editable format
openGyver convertPresentation show.ppsx -o editable.pptx

# Convert ODP to PDF for distribution
openGyver convertPresentation presentation.odp -o presentation.pdf

# Convert macro-enabled presentation to standard format
openGyver convertPresentation macros.pptm -o clean.pptx
```

## Notes

- The conversion is performed by LibreOffice in headless mode (`soffice --headless`).
- LibreOffice outputs files to the output directory with the same base name as the input. If your specified output filename differs, the file is automatically renamed.
- The input format is also validated against the supported formats list, though LibreOffice itself handles the actual format detection.
- Image output (PNG, JPG) may produce one file per slide depending on the LibreOffice version.
- Keynote (`.key`) support depends on the LibreOffice version; newer versions have better compatibility.
- The `POTM` and `PPSM` formats are also accepted as input (present in the source code's supported formats map).

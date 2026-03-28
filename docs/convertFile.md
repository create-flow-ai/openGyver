# convertFile

Convert files between spreadsheet, document, and page-layout formats.

## Usage

```bash
openGyver convertFile <input-file> [flags]
```

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input-file` | Yes | Path to the input file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | Auto-detected | Output file path (default: input name with new extension) |
| `--sheet` | | string | `""` | Sheet name for XLSX/Numbers files (default: first sheet) |
| `--delimiter` | | string | `""` | CSV delimiter: `comma`, `tab`, `semicolon`, `pipe`, or any single character (default: auto-detected) |
| `--quiet` | `-q` | bool | `false` | Suppress output messages (for piping) |
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | bool | `false` | Show help for convertFile |

## Supported Formats

| Category | Extensions |
|----------|-----------|
| Spreadsheet | `.csv`, `.tsv`, `.xlsx`, `.numbers` |
| Document | `.txt`, `.text`, `.md`, `.markdown`, `.html`, `.htm`, `.docx` |
| Page layout | `.pdf`, `.ps` (PostScript) |

## Supported Conversions

### Spreadsheet Conversions

| From | To | Status |
|------|----|--------|
| CSV | XLSX | Fully implemented |
| XLSX | CSV | Fully implemented |
| CSV | Numbers | Stub (workaround provided) |
| XLSX | Numbers | Stub (workaround provided) |
| Numbers | CSV | Stub (workaround provided) |
| Numbers | XLSX | Stub (workaround provided) |

### Document Conversions

| From | To | Description |
|------|----|-------------|
| Markdown | HTML | Rendered with full CommonMark support |
| HTML | Markdown | Converted back to clean Markdown |
| Markdown | Text | Stripped of all formatting |
| HTML | Text | Tags stripped, text extracted |
| DOCX | Text | Text extracted from Word XML |
| DOCX | HTML | Paragraphs converted to HTML |
| DOCX | Markdown | Paragraphs converted to Markdown |
| Text | HTML | Wrapped in `<pre>` block |
| Text | Markdown | Wrapped in code fence |
| Text | DOCX | Plain paragraphs in a Word document |
| Markdown | DOCX | Converted via HTML to Word |
| HTML | DOCX | Paragraphs extracted into Word document |

### PDF / PostScript Output

| From | To | Description |
|------|----|-------------|
| Text | PDF | Plain text rendered to PDF |
| Markdown | PDF | Rendered via HTML to a formatted PDF |
| HTML | PDF | Parsed and rendered to PDF |
| CSV | PDF | Tabular layout with headers |
| XLSX | PDF | Tabular layout with headers |
| DOCX | PDF | Text extracted and rendered to PDF |
| Text | PS | PostScript text output |
| Markdown | PS | Rendered to PostScript |
| CSV | PS | Tabular PostScript output |

### Default Output Formats

When `--output` is not specified, the default output extension is chosen based on the input:

| Input Format | Default Output |
|-------------|----------------|
| CSV | `.xlsx` |
| XLSX | `.csv` |
| Numbers | `.csv` |
| Text | `.pdf` |
| Markdown | `.html` |
| HTML | `.md` |
| DOCX | `.pdf` |

## Examples

```bash
# Convert CSV to Excel
openGyver convertFile data.csv -o report.xlsx

# Convert Excel to PDF with tabular layout
openGyver convertFile report.xlsx -o report.pdf

# Render Markdown to HTML
openGyver convertFile README.md -o README.html

# Render Markdown to PDF
openGyver convertFile README.md -o README.pdf

# Convert HTML back to Markdown
openGyver convertFile page.html -o page.md

# Render HTML to PDF
openGyver convertFile page.html -o page.pdf

# Convert plain text to PDF
openGyver convertFile notes.txt -o notes.pdf

# Convert plain text to Word document
openGyver convertFile notes.txt -o notes.docx

# Convert Word document to PDF
openGyver convertFile report.docx -o report.pdf

# Extract text from Word document
openGyver convertFile report.docx -o report.txt

# Convert CSV to PostScript
openGyver convertFile data.csv -o data.ps

# Use a specific CSV delimiter
openGyver convertFile data.csv -o data.xlsx --delimiter semicolon

# Use tab as delimiter
openGyver convertFile data.tsv -o data.xlsx --delimiter tab

# Specify a sheet name for multi-sheet XLSX
openGyver convertFile report.xlsx -o report.csv --sheet "Q4 Data"

# Quiet mode for scripting
openGyver convertFile data.csv -o data.xlsx -q

# JSON output for automation
openGyver convertFile README.md -o README.html -j

# Pipe JSON to check output path
openGyver convertFile data.csv -o data.xlsx -j | jq '.output'

# Convert Word to Markdown
openGyver convertFile document.docx -o document.md

# Default output (CSV -> XLSX when -o omitted)
openGyver convertFile data.csv
```

## JSON Output Format

```json
{
  "success": true,
  "input": "data.csv",
  "output": "report.xlsx",
  "input_format": "csv",
  "output_format": "xlsx"
}
```

## Notes

- The input format is auto-detected from the file extension.
- If `--output` is not specified, a sensible default extension is chosen (e.g., CSV defaults to XLSX).
- The CSV delimiter is auto-detected from the file content (comma, tab, semicolon, pipe). Use `--delimiter` to override.
- Multi-sheet XLSX files export the first sheet by default. Use `--sheet` to specify a different sheet name.
- Converting between the same format (e.g., CSV to CSV) will produce an error.
- Numbers format support is currently a stub with workaround instructions.

# convertEbook

Convert ebook files between popular formats using Calibre's ebook-convert.

## Usage

```bash
openGyver convertEbook <input-file> [flags]
```

## Prerequisites

Calibre must be installed (specifically the `ebook-convert` command-line tool).

| Platform | Install Command |
|----------|----------------|
| macOS | `brew install calibre` |
| Linux | `apt install calibre` |
| Windows | [calibre-ebook.com/download](https://calibre-ebook.com/download) |
| All | [calibre-ebook.com](https://calibre-ebook.com) |

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input-file` | Yes | Path to the input ebook file |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output file path (required) |
| `--help` | `-h` | bool | `false` | Show help for convertEbook |

## Supported Formats

### Read Formats (25)

AZW, AZW3, AZW4, CBC, CBR, CBZ, CHM, DOCX, EPUB, FB2, HTM, HTML, HTMLZ, LIT, LRF, MOBI, PDB, PDF, PML, PRC, RB, SNB, TCR, TXT, TXTZ

### Write Formats (17)

AZW3, DOCX, EPUB, FB2, HTM, HTMLZ, LRF, MOBI, PDB, PDF, PML, PRC, RB, SNB, TCR, TXT, TXTZ

## Examples

```bash
# Convert EPUB to MOBI (for Kindle)
openGyver convertEbook book.epub -o book.mobi

# Convert EPUB to PDF
openGyver convertEbook book.epub -o book.pdf

# Convert MOBI back to EPUB
openGyver convertEbook book.mobi -o book.epub

# Convert Word document to EPUB
openGyver convertEbook document.docx -o document.epub

# Convert EPUB to AZW3 (modern Kindle format)
openGyver convertEbook book.epub -o book.azw3

# Convert FB2 (FictionBook) to EPUB
openGyver convertEbook book.fb2 -o book.epub

# Convert comic book archive to PDF
openGyver convertEbook comic.cbz -o comic.pdf

# Extract ebook to plain text
openGyver convertEbook novel.epub -o novel.txt

# Convert HTML book to EPUB
openGyver convertEbook manuscript.html -o manuscript.epub

# Convert EPUB to DOCX for editing
openGyver convertEbook book.epub -o book.docx

# Convert old Kindle format to modern one
openGyver convertEbook old.prc -o new.azw3
```

## Notes

- The input and output formats are auto-detected from file extensions.
- Both the input and output format must be in the supported formats list.
- Calibre's `ebook-convert` handles all the actual conversion work, so conversion quality and features depend on the Calibre version installed.
- Some format combinations may produce better results than others (e.g., EPUB to MOBI is well-supported; PDF output may lose some formatting).
- Read-only formats (AZW, AZW4, CBC, CBR, CBZ, CHM, LIT) can be used as input but not output.

# pdf

PDF tools -- merge, split, count pages, and inspect metadata.

## Usage

```bash
openGyver pdf [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for pdf |

## Subcommands

### merge

Merge two or more PDF files into a single output file. All input files must exist. The `--output` flag is required.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file1` | Yes | First PDF file to merge |
| `file2` | Yes | Second PDF file to merge |
| `file3...` | No | Additional PDF files to merge |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Output PDF path (required) |

#### Examples

```bash
# Merge two PDFs
openGyver pdf merge -o combined.pdf file1.pdf file2.pdf

# Merge three PDFs
openGyver pdf merge -o combined.pdf a.pdf b.pdf c.pdf

# Merge all PDFs in current directory
openGyver pdf merge -o all.pdf *.pdf

# Merge with JSON output
openGyver pdf merge -o result.pdf part1.pdf part2.pdf --json

# Merge chapter files
openGyver pdf merge -o book.pdf ch1.pdf ch2.pdf ch3.pdf ch4.pdf

# Merge scanned pages
openGyver pdf merge -o scan-complete.pdf scan001.pdf scan002.pdf scan003.pdf
```

#### JSON Output Format

```json
{
  "action": "merge",
  "inputs": ["file1.pdf", "file2.pdf"],
  "output": "combined.pdf"
}
```

---

### split

Split a PDF into individual single-page PDF files. Each page is written as a separate file in the output directory. If no output directory is specified, files are written to the current directory.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file` | Yes | PDF file to split |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output-dir` | `-o` | string | `.` | Output directory for split pages |

#### Examples

```bash
# Split into current directory
openGyver pdf split document.pdf

# Split into a specific directory
openGyver pdf split document.pdf -o ./pages/

# Split with JSON output
openGyver pdf split report.pdf -o ./split/ --json

# Split a book into pages
openGyver pdf split book.pdf -o ./book-pages/

# Split into a temporary directory
openGyver pdf split form.pdf -o /tmp/pages/

# Split and output JSON summary
openGyver pdf split presentation.pdf --json
```

#### JSON Output Format

```json
{
  "action": "split",
  "input": "document.pdf",
  "outputDir": "./pages/"
}
```

---

### pages

Count the number of pages in a PDF file. Uses pdfcpu for accurate parsing, with a regex-based fallback for files that pdfcpu cannot parse.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file` | Yes | PDF file to count pages in |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Count pages in a document
openGyver pdf pages document.pdf

# Count pages with JSON output
openGyver pdf pages document.pdf --json

# Count pages in a book
openGyver pdf pages textbook.pdf

# Quick page count for a report
openGyver pdf pages quarterly-report.pdf

# JSON output for scripting
openGyver pdf pages manual.pdf -j

# Count pages in a scanned document
openGyver pdf pages scan.pdf
```

#### JSON Output Format

```json
{
  "file": "document.pdf",
  "pages": 42
}
```

---

### info

Show PDF metadata: title, author, creator, producer, creation date, modification date, and page count. Uses pdfcpu to read the PDF cross-reference table for metadata extraction.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file` | Yes | PDF file to inspect |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Show metadata for a document
openGyver pdf info document.pdf

# JSON output for scripting
openGyver pdf info document.pdf --json

# Inspect a book's metadata
openGyver pdf info ebook.pdf

# Check author and title
openGyver pdf info report.pdf

# Machine-readable metadata
openGyver pdf info paper.pdf -j

# Inspect a scanned document
openGyver pdf info scan.pdf
```

#### JSON Output Format

```json
{
  "file": "document.pdf",
  "title": "Annual Report 2025",
  "author": "John Doe",
  "creator": "LibreOffice",
  "producer": "LibreOffice 7.6",
  "creationDate": "20250115120000",
  "modDate": "20250120093000",
  "pages": 42
}
```

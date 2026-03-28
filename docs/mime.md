# mime

MIME type lookup and detection.

## Usage

```bash
openGyver mime [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | | | Help for mime |

## Subcommands

### lookup

Look up the MIME type for a given file extension. The extension can be given with or without a leading dot. Uses Go's `mime.TypeByExtension()` supplemented with a comprehensive built-in map covering text, application, image, audio, video, and programming file types.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `extension` | Yes | File extension (with or without leading dot) |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Look up JSON MIME type
openGyver mime lookup .json

# Extension without the leading dot
openGyver mime lookup json

# Look up PDF MIME type
openGyver mime lookup pdf

# Look up with compound extension
openGyver mime lookup .tar.gz

# JSON output for scripting
openGyver mime lookup .png --json

# Look up a programming language extension
openGyver mime lookup .go

# Look up a video format
openGyver mime lookup .mp4

# Look up a font format
openGyver mime lookup .woff2
```

#### JSON Output Format

```json
{
  "extension": ".json",
  "mime_type": "application/json"
}
```

#### JSON Output Format (unknown extension)

```json
{
  "extension": ".xyz",
  "mime_type": null,
  "error": "unknown extension"
}
```

---

### extension

Reverse lookup: find the file extension for a given MIME type. Returns the most common file extension associated with the MIME type. Checks the built-in reverse map first, then falls back to Go's stdlib `mime.ExtensionsByType()`.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `mime-type` | Yes | The MIME type to look up (e.g. "application/json") |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Find extension for application/json
openGyver mime extension "application/json"

# Find extension for image/png
openGyver mime extension "image/png"

# Find extension for text/html
openGyver mime extension "text/html"

# Find extension for video/mp4
openGyver mime extension "video/mp4"

# JSON output
openGyver mime extension "application/pdf" --json

# Find extension for audio format
openGyver mime extension "audio/mpeg"

# Find extension for a compressed format
openGyver mime extension "application/zip"
```

#### JSON Output Format

```json
{
  "mime_type": "application/json",
  "extension": ".json"
}
```

#### JSON Output Format (unknown MIME type)

```json
{
  "mime_type": "application/unknown",
  "extension": null,
  "error": "unknown MIME type"
}
```

---

### detect

Detect the MIME type of a file by reading its first 512 bytes and using Go's `net/http.DetectContentType()` for content sniffing. This performs binary/magic-byte detection, not extension-based lookup, so it works even on files with incorrect or missing extensions.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file` | Yes | Path to the file to detect |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Detect a JPEG image
openGyver mime detect photo.jpg

# Detect a file with no extension
openGyver mime detect mystery-file

# Detect a binary executable
openGyver mime detect /usr/bin/ls

# JSON output
openGyver mime detect document.pdf --json

# Detect a text file
openGyver mime detect README.md

# Detect a compressed archive
openGyver mime detect archive.tar.gz

# Detect an HTML file
openGyver mime detect index.html
```

#### JSON Output Format

```json
{
  "file": "photo.jpg",
  "mime_type": "image/jpeg"
}
```

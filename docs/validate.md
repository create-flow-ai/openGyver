# validate

Validate common data and markup formats.

## Usage

```bash
openGyver validate [command] [flags]
```

Each subcommand accepts input as the first argument, or via `--file/-f`. Output is "valid" on success, or a list of errors. With `--json/-j`, output is `{"valid":true/false,"errors":[...]}`.

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | | Read input from a file instead of an argument |
| `--json` | `-j` | bool | `false` | Output as `{"valid":true/false,"errors":[...]}` |
| `--help` | `-h` | | | Help for validate |

## Subcommands

### html

Validate HTML for common issues:

- Missing `<!DOCTYPE html>` declaration
- Unclosed tags (e.g. `<p>` without `</p>`)
- Mismatched tags (e.g. `<b>...</i>`)
- Missing `alt` attribute on `<img>` elements
- Duplicate `id` attributes

Uses the `golang.org/x/net/html` tokenizer for parsing.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The HTML string to validate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Validate an HTML string
openGyver validate html '<html><body><p>Hello</p></body></html>'

# Validate a file
openGyver validate html --file page.html

# Validate with JSON output
openGyver validate html --file page.html --json

# Check a full HTML document
openGyver validate html '<!DOCTYPE html><html><head><title>Test</title></head><body><p>Hello</p></body></html>'

# Detect missing alt on images
openGyver validate html '<img src="photo.jpg">'

# Detect duplicate IDs
openGyver validate html '<div id="main"></div><div id="main"></div>'

# Validate a template file
openGyver validate html --file templates/index.html --json

# Detect unclosed tags
openGyver validate html '<div><p>unclosed paragraph</div>'
```

#### JSON Output Format (valid)

```json
{
  "valid": true,
  "errors": []
}
```

#### JSON Output Format (invalid)

```json
{
  "valid": false,
  "errors": [
    "missing <!DOCTYPE html> declaration",
    "line 1: <img> missing alt attribute"
  ]
}
```

---

### csv

Validate CSV for common issues:

- Inconsistent column count across rows
- Improper quoting / bare quotes
- Encoding / parse errors

Uses `encoding/csv` in strict mode (FieldsPerRecord set from the first row).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The CSV string to validate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Validate a CSV string
openGyver validate csv 'name,age
Alice,30
Bob,25'

# Validate a CSV file
openGyver validate csv --file data.csv

# Validate with JSON output
openGyver validate csv --file data.csv --json

# Detect inconsistent column count
openGyver validate csv 'a,b,c
1,2
3,4,5'

# Validate a large dataset
openGyver validate csv --file export.csv

# Validate and get structured error report
openGyver validate csv --file import.csv --json

# Check quoting issues
openGyver validate csv 'name,bio
Alice,"She said "hello""'

# Validate a header-only CSV
openGyver validate csv 'name,email,age'
```

#### JSON Output Format (valid)

```json
{
  "valid": true,
  "errors": []
}
```

#### JSON Output Format (invalid)

```json
{
  "valid": false,
  "errors": [
    "row 3: expected 3 columns, got 2"
  ]
}
```

---

### xml

Validate XML well-formedness by parsing with `encoding/xml.Decoder`. Reports any parse errors with byte offset position.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The XML string to validate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Validate a well-formed XML string
openGyver validate xml '<root><item/></root>'

# Validate an XML file
openGyver validate xml --file feed.xml

# Validate with JSON output
openGyver validate xml --file feed.xml --json

# Detect malformed XML
openGyver validate xml '<root><unclosed>'

# Validate an SVG file
openGyver validate xml --file icon.svg

# Validate a configuration XML
openGyver validate xml --file app-config.xml --json

# Validate a sitemap
openGyver validate xml --file sitemap.xml

# Detect mismatched tags
openGyver validate xml '<a><b></a></b>'
```

#### JSON Output Format (valid)

```json
{
  "valid": true,
  "errors": []
}
```

#### JSON Output Format (invalid)

```json
{
  "valid": false,
  "errors": [
    "offset 18: XML syntax error on line 1: unexpected end element </a>"
  ]
}
```

---

### yaml

Validate YAML syntax by attempting to unmarshal with `gopkg.in/yaml.v3`. Reports any parse errors with line/column context.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The YAML string to validate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Validate a YAML string
openGyver validate yaml 'name: hello'

# Validate a YAML file
openGyver validate yaml --file config.yaml

# Validate with JSON output
openGyver validate yaml --file config.yaml --json

# Detect indentation errors
openGyver validate yaml 'items:
- first
  - nested-wrong'

# Validate a Kubernetes manifest
openGyver validate yaml --file deployment.yaml

# Validate a CI/CD config
openGyver validate yaml --file .github/workflows/ci.yml --json

# Validate a Docker Compose file
openGyver validate yaml --file docker-compose.yml

# Validate a complex nested structure
openGyver validate yaml --file values.yaml
```

#### JSON Output Format (valid)

```json
{
  "valid": true,
  "errors": []
}
```

#### JSON Output Format (invalid)

```json
{
  "valid": false,
  "errors": [
    "yaml: line 3: did not find expected key"
  ]
}
```

---

### toml

Validate TOML syntax by attempting to decode with `github.com/BurntSushi/toml`. Reports any parse errors with position context.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The TOML string to validate |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Validate a TOML string
openGyver validate toml 'key = "value"'

# Validate a pyproject.toml
openGyver validate toml --file pyproject.toml

# Validate with JSON output
openGyver validate toml --file config.toml --json

# Validate a Cargo.toml
openGyver validate toml --file Cargo.toml

# Detect syntax errors
openGyver validate toml 'key = value without quotes'

# Validate a Hugo configuration
openGyver validate toml --file hugo.toml --json

# Validate a Rust project config
openGyver validate toml --file Cargo.toml

# Validate an application config
openGyver validate toml --file app.toml
```

#### JSON Output Format (valid)

```json
{
  "valid": true,
  "errors": []
}
```

#### JSON Output Format (invalid)

```json
{
  "valid": false,
  "errors": [
    "toml: line 1: expected value but found \"without\" instead"
  ]
}
```

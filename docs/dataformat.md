# dataformat

Convert between common data serialisation formats (YAML, JSON, TOML, CSV, XML).

## Usage

```bash
openGyver dataformat [subcommand] [input] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Wrap output in `{"input_format","output_format","data"}` JSON envelope |
| `--file` | `-f` | string | | Read input from a file instead of a positional argument |
| `--help` | `-h` | bool | | Show help for the command |

## Subcommands

### yaml2json

Convert YAML input to pretty-printed JSON output. Handles nested maps, arrays, and all standard YAML types. Keys that are non-string types in YAML are coerced to strings for JSON compatibility.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | YAML string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline YAML to JSON
openGyver dataformat yaml2json 'name: hello'

# Convert a YAML file to JSON, print to stdout
openGyver dataformat yaml2json --file config.yaml

# Convert and write to a file
openGyver dataformat yaml2json --file config.yaml -o config.json

# Wrap result in a JSON metadata envelope
openGyver dataformat yaml2json --file config.yaml --json

# Pipe YAML into the command
cat config.yaml | openGyver dataformat yaml2json

# Convert nested YAML
openGyver dataformat yaml2json 'server:
  host: localhost
  port: 8080'

# Read from file, write to file, with JSON envelope
openGyver dataformat yaml2json -f settings.yaml -o settings.json -j

# Convert a list
openGyver dataformat yaml2json '- one
- two
- three'
```

#### JSON Output Format

```json
{
  "input_format": "yaml",
  "output_format": "json",
  "data": {
    "name": "hello"
  }
}
```

---

### json2yaml

Convert JSON input to YAML output. Handles nested objects, arrays, and all standard JSON types.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | JSON string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline JSON to YAML
openGyver dataformat json2yaml '{"name": "hello", "count": 42}'

# Convert a JSON file to YAML
openGyver dataformat json2yaml --file config.json

# Convert and save to a file
openGyver dataformat json2yaml --file config.json -o config.yaml

# With JSON metadata envelope
openGyver dataformat json2yaml '{"key": "value"}' --json

# Convert a JSON array
openGyver dataformat json2yaml '[1, 2, 3]'

# Read from file, wrap in JSON envelope
openGyver dataformat json2yaml -f data.json -j

# Pipe JSON input
echo '{"host": "localhost"}' | openGyver dataformat json2yaml

# Convert nested JSON
openGyver dataformat json2yaml '{"server": {"host": "0.0.0.0", "port": 3000}}'
```

#### JSON Output Format

```json
{
  "input_format": "json",
  "output_format": "yaml",
  "data": "name: hello\ncount: 42"
}
```

---

### toml2json

Convert TOML input to pretty-printed JSON output.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | TOML string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline TOML to JSON
openGyver dataformat toml2json 'name = "hello"'

# Convert a TOML config file
openGyver dataformat toml2json --file config.toml

# Convert and write output to a file
openGyver dataformat toml2json --file config.toml -o config.json

# With JSON envelope
openGyver dataformat toml2json --file Cargo.toml --json

# Convert TOML with sections
openGyver dataformat toml2json '[database]
host = "localhost"
port = 5432'

# Pipe TOML to the command
cat pyproject.toml | openGyver dataformat toml2json

# Full pipeline: convert and save
openGyver dataformat toml2json -f config.toml -o config.json

# Read from file with JSON wrapper
openGyver dataformat toml2json -f settings.toml -j
```

#### JSON Output Format

```json
{
  "input_format": "toml",
  "output_format": "json",
  "data": {
    "name": "hello"
  }
}
```

---

### json2toml

Convert JSON input to TOML output.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | JSON string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline JSON to TOML
openGyver dataformat json2toml '{"name": "hello"}'

# Convert a JSON file to TOML
openGyver dataformat json2toml --file config.json

# Save output to file
openGyver dataformat json2toml --file config.json -o config.toml

# With JSON metadata envelope
openGyver dataformat json2toml '{"key": "value"}' --json

# Convert nested JSON
openGyver dataformat json2toml '{"database": {"host": "localhost", "port": 5432}}'

# Pipe JSON into the converter
echo '{"debug": true}' | openGyver dataformat json2toml

# Read from file with JSON wrapper
openGyver dataformat json2toml -f package.json -j

# Full pipeline: file-to-file
openGyver dataformat json2toml -f data.json -o data.toml
```

#### JSON Output Format

```json
{
  "input_format": "json",
  "output_format": "toml",
  "data": "name = \"hello\""
}
```

---

### csv2json

Convert CSV to a JSON array of objects. The first row of the CSV is used as the header/key names for the resulting JSON objects.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | CSV string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline CSV
openGyver dataformat csv2json 'name,age
Alice,30
Bob,25'

# Convert a CSV file to JSON
openGyver dataformat csv2json --file data.csv

# Save output to a JSON file
openGyver dataformat csv2json --file data.csv -o data.json

# With JSON metadata envelope
openGyver dataformat csv2json --file data.csv --json

# Pipe CSV into the command
cat export.csv | openGyver dataformat csv2json

# Convert and wrap in JSON envelope, save to file
openGyver dataformat csv2json -f users.csv -o users.json -j

# Process a large export
openGyver dataformat csv2json --file report.csv --output report.json

# Inline CSV with special characters
openGyver dataformat csv2json 'key,value
"hello, world","foo""bar"'
```

#### JSON Output Format

```json
{
  "input_format": "csv",
  "output_format": "json",
  "data": [
    {
      "name": "Alice",
      "age": "30"
    },
    {
      "name": "Bob",
      "age": "25"
    }
  ]
}
```

---

### json2csv

Convert a JSON array of objects to CSV format. All unique keys across all objects are collected and sorted alphabetically to form the CSV header row.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | JSON array of objects to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline JSON array to CSV
openGyver dataformat json2csv '[{"name":"Alice","age":30},{"name":"Bob","age":25}]'

# Convert a JSON file to CSV
openGyver dataformat json2csv --file data.json

# Save output to a CSV file
openGyver dataformat json2csv --file data.json -o data.csv

# With JSON envelope
openGyver dataformat json2csv --file users.json --json

# Pipe JSON into the converter
cat records.json | openGyver dataformat json2csv

# Full pipeline: file-to-file
openGyver dataformat json2csv -f export.json -o export.csv

# Convert API response
curl -s https://api.example.com/users | openGyver dataformat json2csv

# With envelope and output file
openGyver dataformat json2csv -f items.json -o items.csv -j
```

#### JSON Output Format

```json
{
  "input_format": "json",
  "output_format": "csv",
  "data": "age,name\n30,Alice\n25,Bob"
}
```

---

### xml2json

Convert XML to JSON using a simple element mapping. Attributes are prefixed with `-`, text content in mixed elements uses the `#text` key. Repeated child elements with the same tag name are grouped into arrays.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | XML string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline XML to JSON
openGyver dataformat xml2json '<root><name>hello</name></root>'

# Convert an XML file
openGyver dataformat xml2json --file data.xml

# Save output to a JSON file
openGyver dataformat xml2json --file feed.xml -o feed.json

# With JSON metadata envelope
openGyver dataformat xml2json --file config.xml --json

# XML with attributes
openGyver dataformat xml2json '<item id="1" active="true"><name>Widget</name></item>'

# Pipe XML into the converter
cat sitemap.xml | openGyver dataformat xml2json

# Convert and save with envelope
openGyver dataformat xml2json -f response.xml -o response.json -j

# Convert XML with repeated elements (auto-arrays)
openGyver dataformat xml2json '<list><item>A</item><item>B</item><item>C</item></list>'
```

#### JSON Output Format

```json
{
  "input_format": "xml",
  "output_format": "json",
  "data": {
    "root": {
      "name": "hello"
    }
  }
}
```

---

### json2xml

Convert JSON to XML. Objects become elements, arrays produce repeated elements, and keys prefixed with `-` become XML attributes. The `#text` key maps to text content within an element. Output includes an XML declaration header.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | JSON string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | | Write output to a file instead of stdout |

#### Examples

```bash
# Convert inline JSON to XML
openGyver dataformat json2xml '{"root": {"name": "hello"}}'

# Convert a JSON file to XML
openGyver dataformat json2xml --file data.json

# Save output to an XML file
openGyver dataformat json2xml --file data.json -o data.xml

# With JSON envelope
openGyver dataformat json2xml '{"item": {"name": "test"}}' --json

# JSON with attributes (use "-" prefix)
openGyver dataformat json2xml '{"item": {"-id": "1", "name": "Widget"}}'

# Pipe JSON into the converter
cat config.json | openGyver dataformat json2xml

# Full pipeline: file-to-file with envelope
openGyver dataformat json2xml -f records.json -o records.xml -j

# JSON array to repeated XML elements
openGyver dataformat json2xml '{"list": {"item": ["A", "B", "C"]}}'
```

#### JSON Output Format

```json
{
  "input_format": "json",
  "output_format": "xml",
  "data": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<root>\n  <name>hello</name>\n</root>"
}
```

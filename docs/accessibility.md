# accessibility

Tools for checking web accessibility compliance.

## Usage

```bash
openGyver accessibility [command] [flags]
```

## Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | `false` | Show help for accessibility |
| `--json` | `-j` | bool | `false` | Output as JSON |

## Subcommands

### contrast

Check the contrast ratio between two colors per WCAG 2.1 guidelines. Reports the ratio and pass/fail for AA and AAA at normal and large text sizes. Colors can be hex (`#fff`, `#ffffff`), `rgb(r,g,b)`, or CSS names.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `color1` | Yes | First color (hex, rgb, or CSS name) |
| `color2` | Yes | Second color (hex, rgb, or CSS name) |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | `false` | Show help for contrast |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited from parent) |

#### WCAG Levels Reference

| Level | Context | Minimum Ratio |
|-------|---------|---------------|
| AA | Normal text | >= 4.5:1 |
| AA | Large text | >= 3.0:1 |
| AAA | Normal text | >= 7.0:1 |
| AAA | Large text | >= 4.5:1 |

#### Examples

```bash
# Check contrast between white and black (maximum contrast)
openGyver accessibility contrast "#ffffff" "#000000"

# Use shorthand hex notation
openGyver accessibility contrast "#333" "#ccc"

# Use CSS color names instead of hex
openGyver accessibility contrast white black

# Check contrast for a colored background
openGyver accessibility contrast "#ff5733" "#000000"

# Get JSON output for programmatic use
openGyver accessibility contrast "#ffffff" "#000000" -j

# Use rgb() format for one or both colors
openGyver accessibility contrast "rgb(255,255,255)" "#333333"

# Pipe JSON output to jq for specific fields
openGyver accessibility contrast "#fff" "#000" -j | jq '.ratio'

# Check a low-contrast pair (will show FAIL results)
openGyver accessibility contrast "#777" "#888"
```

#### JSON Output Format

```json
{
  "color1": "#ffffff",
  "color2": "#000000",
  "ratio": 21,
  "aa_normal": true,
  "aa_large": true,
  "aaa_normal": true,
  "aaa_large": true
}
```

---

### readability

Calculate Flesch Reading Ease, Flesch-Kincaid Grade Level, and Gunning Fog Index for a given text.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No | The text to analyze (required if `--file` is not used) |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | `""` | Read text from a file instead of an argument |
| `--help` | `-h` | bool | `false` | Show help for readability |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited from parent) |

#### Flesch Reading Ease Scale

| Score | Label |
|-------|-------|
| 90-100 | Very Easy |
| 80-89 | Easy |
| 70-79 | Fairly Easy |
| 60-69 | Standard |
| 50-59 | Fairly Difficult |
| 30-49 | Difficult |
| 0-29 | Very Confusing |

#### Examples

```bash
# Analyze a simple sentence
openGyver accessibility readability "The cat sat on the mat."

# Analyze more complex text
openGyver accessibility readability "The implementation of quantum computing paradigms necessitates a fundamental reconsideration of classical algorithmic approaches."

# Read text from a file
openGyver accessibility readability --file article.txt

# Short flag for file input
openGyver accessibility readability -f essay.txt

# Get JSON output
openGyver accessibility readability -j "Complex text here."

# Read from file with JSON output
openGyver accessibility readability --file report.txt --json

# Pipe file content and capture JSON result
openGyver accessibility readability -f chapter.txt -j | jq '.flesch_reading_ease'

# Check grade level of a document
openGyver accessibility readability --file legal-document.txt -j | jq '.flesch_kincaid_grade'
```

#### JSON Output Format

```json
{
  "words": 8,
  "sentences": 1,
  "syllables": 8,
  "flesch_reading_ease": 104.83,
  "flesch_reading_ease_label": "Very Easy",
  "flesch_kincaid_grade": -1.45,
  "gunning_fog_index": 3.2
}
```

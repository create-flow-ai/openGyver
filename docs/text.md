# text

Text manipulation utilities -- count, convert case, reverse, sort, deduplicate, slugify, generate lorem ipsum, diff, wrap, number lines, trim, and find-and-replace.

## Usage

```bash
openGyver text [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for text |

## Subcommands

### count

Count the number of words, characters, lines, and sentences in text. Input can be provided as a positional argument, via `--file/-f`, or piped through stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use `--file` or stdin) | The text to analyze |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | | Read input from a file |
| `--help` | `-h` | | | Help for count |

#### Examples

```bash
# Count words in a string
openGyver text count "The quick brown fox jumps over the lazy dog."

# Count from a file
openGyver text count --file essay.txt

# Count from piped input
echo "hello world" | openGyver text count

# JSON output
openGyver text count --json "Hello, world!"

# Count a longer passage
openGyver text count "This is sentence one. This is sentence two. And a third."

# Count source code lines
openGyver text count --file main.go --json
```

#### JSON Output Format

```json
{
  "words": 9,
  "characters": 44,
  "lines": 1,
  "sentences": 1
}
```

---

### case

Convert text to a different case style.

**Supported cases (`--to`):**

| Case | Example |
|------|---------|
| `upper` | HELLO WORLD |
| `lower` | hello world |
| `title` | Hello World |
| `sentence` | Hello world |
| `camel` | helloWorld |
| `pascal` | HelloWorld |
| `snake` | hello_world |
| `kebab` | hello-world |
| `constant` | HELLO_WORLD |
| `dot` | hello.world |

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use stdin) | The text to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--to` | | string | `lower` | Target case style |
| `--help` | `-h` | | | Help for case |

#### Examples

```bash
# Convert to uppercase
openGyver text case --to upper "hello world"

# Convert to snake_case
openGyver text case --to snake "Hello World"

# Convert to camelCase
openGyver text case --to camel "some-variable-name"

# Convert to kebab-case
openGyver text case --to kebab "myVariableName"

# Convert to CONSTANT_CASE
openGyver text case --to constant "max retries"

# Convert to dot.case
openGyver text case --to dot "Hello World"

# JSON output
openGyver text case --to title "the quick brown fox" --json

# Convert to PascalCase
openGyver text case --to pascal "hello world"
```

#### JSON Output Format

```json
{
  "original": "Hello World",
  "case": "snake",
  "result": "hello_world"
}
```

---

### reverse

Reverse the characters in a string.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use stdin) | The string to reverse |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Reverse a greeting
openGyver text reverse "Hello, World!"

# Reverse a simple string
openGyver text reverse "abcdef"

# JSON output
openGyver text reverse --json "palindrome"

# Check if a word is a palindrome
openGyver text reverse "racecar"

# Reverse a number string
openGyver text reverse "1234567890"

# Reverse with unicode characters
openGyver text reverse "hello world"
```

#### JSON Output Format

```json
{
  "original": "abcdef",
  "reversed": "fedcba"
}
```

---

### sort

Sort lines of text alphabetically, by length, or numerically. Use `--reverse` to invert the order. Input via `--file/-f`, positional argument (newlines as `\n`), or stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use `--file` or stdin) | The text with newline-separated lines |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--by` | | string | `alpha` | Sort order: alpha, length, numeric |
| `--file` | `-f` | string | | Read input from a file |
| `--reverse` | | bool | `false` | Reverse the sort order |
| `--help` | `-h` | | | Help for sort |

#### Examples

```bash
# Sort a file alphabetically
openGyver text sort --file names.txt

# Sort by line length
openGyver text sort --by length --file code.txt

# Sort numerically in reverse order
openGyver text sort --by numeric --reverse --file scores.txt

# Sort piped input
echo -e "cherry\napple\nbanana" | openGyver text sort

# Sort alphabetically in reverse
openGyver text sort --reverse --file words.txt

# Sort code by line length (find longest lines)
openGyver text sort --by length --reverse --file main.go

# JSON output
openGyver text sort --file data.txt --json
```

#### JSON Output Format

```json
{
  "lines": ["apple", "banana", "cherry"],
  "order": "alpha",
  "reversed": false
}
```

---

### dedupe

Remove duplicate lines while preserving the original order. Input via `--file/-f`, positional argument, or stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use `--file` or stdin) | The text with potential duplicate lines |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | | Read input from a file |
| `--help` | `-h` | | | Help for dedupe |

#### Examples

```bash
# Deduplicate a file
openGyver text dedupe --file list.txt

# Deduplicate piped input
echo -e "a\nb\na\nc\nb" | openGyver text dedupe

# JSON output
openGyver text dedupe --json --file data.txt

# Clean up a log file
openGyver text dedupe --file error-messages.txt

# Remove duplicate imports from a list
openGyver text dedupe --file dependencies.txt

# Deduplicate with output stats
openGyver text dedupe --json --file urls.txt
```

#### JSON Output Format

```json
{
  "original_count": 5,
  "unique_count": 3,
  "lines": ["a", "b", "c"]
}
```

---

### slug

Generate a URL-safe slug from text. Converts to lowercase, replaces spaces and special characters with hyphens, and collapses multiple hyphens.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use stdin) | The text to slugify |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Slugify a blog post title
openGyver text slug "My Blog Post Title!"

# Slugify text with extra spaces and punctuation
openGyver text slug "Hello,   World!! 2024"

# JSON output
openGyver text slug --json "My Article Title"

# Slugify with special characters
openGyver text slug "Cafe & Restaurant -- The Best!"

# Slugify a product name
openGyver text slug "iPhone 15 Pro Max (256GB)"

# Slugify with international characters
openGyver text slug "Resume: My Career Journey"
```

#### JSON Output Format

```json
{
  "original": "My Blog Post Title!",
  "slug": "my-blog-post-title"
}
```

---

### lorem

Generate Lorem Ipsum placeholder text. Exactly one of `--words`, `--sentences`, or `--paragraphs` must be given.

#### Arguments

No positional arguments.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--words` | | int | | Number of words to generate |
| `--sentences` | | int | | Number of sentences to generate |
| `--paragraphs` | | int | | Number of paragraphs to generate |
| `--help` | `-h` | | | Help for lorem |

#### Examples

```bash
# Generate 50 words
openGyver text lorem --words 50

# Generate 5 sentences
openGyver text lorem --sentences 5

# Generate 3 paragraphs
openGyver text lorem --paragraphs 3

# JSON output
openGyver text lorem --paragraphs 2 --json

# Generate a short snippet
openGyver text lorem --words 10

# Generate a single paragraph
openGyver text lorem --paragraphs 1

# Generate placeholder for a form field
openGyver text lorem --sentences 2

# Save to a file for testing
openGyver text lorem --paragraphs 5 > placeholder.txt
```

#### JSON Output Format

```json
{
  "text": "Lorem ipsum dolor sit amet...",
  "words": 50
}
```

---

### diff

Show a unified diff between two text files.

#### Arguments

No positional arguments.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file1` | | string | | First file to compare |
| `--file2` | | string | | Second file to compare |
| `--help` | `-h` | | | Help for diff |

#### Examples

```bash
# Compare two text files
openGyver text diff --file1 original.txt --file2 modified.txt

# Compare source code files with JSON output
openGyver text diff --file1 a.go --file2 b.go --json

# Compare configuration files
openGyver text diff --file1 config.old.yaml --file2 config.yaml

# Compare two versions of a script
openGyver text diff --file1 deploy-v1.sh --file2 deploy-v2.sh

# Compare README versions
openGyver text diff --file1 README.old.md --file2 README.md

# Get JSON diff for programmatic processing
openGyver text diff --file1 before.txt --file2 after.txt --json
```

---

### wrap

Word-wrap text to a given column width (default 80). Input via positional argument, `--file`, or stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use `--file` or stdin) | The text to wrap |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--width` | | int | `80` | Column width to wrap at |
| `--file` | `-f` | string | | Read input from a file |
| `--help` | `-h` | | | Help for wrap |

#### Examples

```bash
# Wrap at 60 columns
openGyver text wrap --width 60 "This is a very long sentence that should be wrapped at sixty characters."

# Wrap a file at default width (80)
openGyver text wrap --file article.txt

# Wrap with JSON output
openGyver text wrap --width 40 --json "Some text to wrap."

# Narrow wrap for mobile preview
openGyver text wrap --width 30 --file long-text.txt

# Wrap for email formatting
openGyver text wrap --width 72 --file email-body.txt

# Wrap and save
openGyver text wrap --width 60 --file readme.txt > readme-wrapped.txt
```

---

### lines

Prefix each line of text with its line number. Input via `--file/-f`, positional argument, or stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use `--file` or stdin) | The text to number |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | | Read input from a file |
| `--help` | `-h` | | | Help for lines |

#### Examples

```bash
# Add line numbers to a source file
openGyver text lines --file main.go

# Number piped input
echo -e "alpha\nbeta\ngamma" | openGyver text lines

# JSON output
openGyver text lines --json --file script.sh

# Number a configuration file for review
openGyver text lines --file config.yaml

# Number a log file for reference
openGyver text lines --file app.log

# Pipe numbered output for review
openGyver text lines --file data.csv | less
```

---

### trim

Remove leading and trailing whitespace from each line. Use `--blank` to also remove entirely blank lines. Input via positional argument, `--file`, or stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use `--file` or stdin) | The text to trim |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--blank` | | bool | `false` | Also remove blank lines |
| `--file` | `-f` | string | | Read input from a file |
| `--help` | `-h` | | | Help for trim |

#### Examples

```bash
# Trim a simple string
openGyver text trim "  hello world  "

# Trim and remove blank lines from a file
openGyver text trim --blank --file messy.txt

# Trim piped input and remove blank lines
echo -e "  hi  \n\n  there  " | openGyver text trim --blank

# Trim whitespace only (keep blank lines)
openGyver text trim --file indented.txt

# JSON output
openGyver text trim --json "  padded string  "

# Clean up copy-pasted code
openGyver text trim --blank --file snippet.txt
```

---

### replace

Find and replace text using literal strings or regular expressions. Use `--find` and `--replace` to specify the search and replacement strings. Use `--regex` to interpret `--find` as a regular expression (supports Go regexp syntax including capture groups like `$1`, `$2`). Input via `--file/-f`, positional argument, or stdin.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `text` | No (use `--file` or stdin) | The text to perform replacement on |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--find` | | string | | String or pattern to find |
| `--replace` | | string | | Replacement string |
| `--regex` | | bool | `false` | Interpret `--find` as a regular expression |
| `--file` | `-f` | string | | Read input from a file |
| `--help` | `-h` | | | Help for replace |

#### Examples

```bash
# Simple literal replacement
openGyver text replace --find foo --replace bar "foo fighters"

# Replace in a file
openGyver text replace --find foo --replace bar --file input.txt

# Regex replacement with capture groups
openGyver text replace --regex --find "(\w+)@(\w+)" --replace "$1 at $2" "user@host"

# JSON output
openGyver text replace --find old --replace new --json "old is gold"

# Replace multiple occurrences
openGyver text replace --find "TODO" --replace "DONE" --file tasks.txt

# Regex: normalize whitespace
openGyver text replace --regex --find "\s+" --replace " " "too   many    spaces"

# Replace URLs
openGyver text replace --find "http://" --replace "https://" --file links.txt

# Complex regex replacement
openGyver text replace --regex --find "v(\d+)\.(\d+)" --replace "version $1.$2" "Upgraded to v3.14"
```

#### JSON Output Format

```json
{
  "original": "foo fighters",
  "result": "bar fighters",
  "find": "foo",
  "replace": "bar"
}
```

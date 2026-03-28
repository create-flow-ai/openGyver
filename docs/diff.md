# diff

Diff and compare files -- unified text diff, JSON structural diff, CSV diff.

## Usage

```bash
openGyver diff [subcommand] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output results as machine-readable JSON |
| `--help` | `-h` | bool | | Show help for the command |

## Subcommands

### text

Compare two text files line-by-line and display a unified diff. Uses the Longest Common Subsequence (LCS) algorithm to compute the minimal set of changes. Lines only in file1 are prefixed with `-`, lines only in file2 with `+`, and common lines with a space.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file1` | | string | | First file to compare (required) |
| `--file2` | | string | | Second file to compare (required) |

#### Examples

```bash
# Basic text diff between two files
openGyver diff text --file1 a.txt --file2 b.txt

# Diff two Go source files
openGyver diff text --file1 original.go --file2 modified.go

# Get JSON output for programmatic processing
openGyver diff text --file1 a.txt --file2 b.txt --json

# Diff two configuration files
openGyver diff text --file1 config.old --file2 config.new

# Diff and pipe through a pager
openGyver diff text --file1 before.txt --file2 after.txt | less

# JSON diff of shell scripts
openGyver diff text --file1 deploy-v1.sh --file2 deploy-v2.sh -j

# Compare README versions
openGyver diff text --file1 README.old.md --file2 README.md

# Diff two SQL migration files
openGyver diff text --file1 001_init.sql --file2 002_update.sql
```

#### JSON Output Format

```json
{
  "file1": "a.txt",
  "file2": "b.txt",
  "diffs": [
    {
      "op": "keep",
      "line": "line that exists in both files"
    },
    {
      "op": "remove",
      "line": "line only in file1"
    },
    {
      "op": "add",
      "line": "line only in file2"
    }
  ]
}
```

---

### json

Compare two JSON files structurally and show added, removed, and changed keys. Recursively walks both JSON structures and reports differences at each dotted path. Array elements are compared by index.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file1` | | string | | First JSON file to compare (required) |
| `--file2` | | string | | Second JSON file to compare (required) |

#### Examples

```bash
# Compare two JSON config files
openGyver diff json --file1 old.json --file2 new.json

# Structural diff with JSON output
openGyver diff json --file1 config-a.json --file2 config-b.json --json

# Compare API responses saved to files
openGyver diff json --file1 response-v1.json --file2 response-v2.json

# Compare package.json versions
openGyver diff json --file1 package.old.json --file2 package.json

# Compare and pipe to jq for filtering
openGyver diff json --file1 a.json --file2 b.json -j | jq '.diffs[] | select(.type == "changed")'

# Compare two Terraform state files
openGyver diff json --file1 state-before.json --file2 state-after.json

# Check for differences in CI config
openGyver diff json --file1 pipeline-prod.json --file2 pipeline-staging.json -j

# Compare exported database records
openGyver diff json --file1 export-jan.json --file2 export-feb.json
```

#### JSON Output Format

```json
{
  "file1": "old.json",
  "file2": "new.json",
  "diffs": [
    {
      "path": "config.timeout",
      "type": "changed",
      "old": 30,
      "new": 60
    },
    {
      "path": "config.retries",
      "type": "added",
      "new": 3
    },
    {
      "path": "config.debug",
      "type": "removed",
      "old": true
    }
  ]
}
```

---

### csv

Compare two CSV files and show added, removed, and changed rows. Rows are compared line-by-line by their full content. Supports variable field counts across files.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file1` | | string | | First CSV file to compare (required) |
| `--file2` | | string | | Second CSV file to compare (required) |

#### Examples

```bash
# Compare two CSV exports
openGyver diff csv --file1 old.csv --file2 new.csv

# JSON output for scripting
openGyver diff csv --file1 export-jan.csv --file2 export-feb.csv --json

# Compare inventory snapshots
openGyver diff csv --file1 inventory-monday.csv --file2 inventory-friday.csv

# Compare user exports
openGyver diff csv --file1 users-before.csv --file2 users-after.csv

# Pipe JSON output to jq
openGyver diff csv --file1 a.csv --file2 b.csv -j | jq '.diffs[] | select(.type == "added")'

# Compare price lists
openGyver diff csv --file1 prices-q1.csv --file2 prices-q2.csv

# Compare database dumps
openGyver diff csv --file1 dump-prod.csv --file2 dump-staging.csv -j

# Compare test results
openGyver diff csv --file1 results-baseline.csv --file2 results-current.csv
```

#### JSON Output Format

```json
{
  "file1": "old.csv",
  "file2": "new.csv",
  "diffs": [
    {
      "type": "changed",
      "row": 3,
      "old": ["Alice", "30", "NYC"],
      "new": ["Alice", "31", "NYC"]
    },
    {
      "type": "added",
      "row": 5,
      "new": ["Charlie", "28", "LA"]
    },
    {
      "type": "removed",
      "row": 4,
      "old": ["Bob", "25", "SF"]
    }
  ]
}
```

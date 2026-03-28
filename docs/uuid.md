# uuid

Generate universally unique identifiers (UUIDs).

## Usage

```bash
openGyver uuid [flags]
```

## Supported Versions

| Version | Description |
|---------|-------------|
| `4` | Random UUID (default). 122 bits of randomness. Best for most use cases. |
| `6` | Reordered time-based UUID. Sortable by creation time, includes a timestamp and random node. Good for database primary keys. |

## Arguments

No positional arguments.

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--version` | | int | `4` | UUID version: 4 (random) or 6 (time-sorted) |
| `--count` | | int | `1` | Number of UUIDs to generate |
| `--uppercase` | | bool | `false` | Output in uppercase |
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for uuid |

## Examples

```bash
# Generate a single random UUID (v4)
openGyver uuid

# Explicitly request v4
openGyver uuid --version 4

# Generate a time-sorted UUID (v6)
openGyver uuid --version 6

# Generate 5 UUIDs at once
openGyver uuid --count 5

# Generate 10 time-sorted UUIDs
openGyver uuid --version 6 --count 10

# Output in uppercase
openGyver uuid --uppercase

# JSON output for scripting
openGyver uuid --json

# Generate multiple UUIDs in uppercase JSON
openGyver uuid --count 3 --uppercase --json

# Generate a single UUID and copy to clipboard (macOS)
openGyver uuid | pbcopy

# Generate UUIDs for database seeds
openGyver uuid --version 6 --count 20

# Quick ID for a test record
openGyver uuid --version 4
```

## JSON Output Format

```json
{
  "version": 4,
  "count": 3,
  "uuids": [
    "550e8400-e29b-41d4-a716-446655440000",
    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "f47ac10b-58cc-4372-a567-0e02b2c3d479"
  ]
}
```

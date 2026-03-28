# chmod

File permission calculator -- convert between octal and symbolic file permissions and calculate umask effects.

## Usage

```bash
openGyver chmod [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for chmod |

## Subcommands

### calc

Convert file permissions between octal (e.g. 755) and symbolic (e.g. rwxr-xr-x) representations. Auto-detects the input format. Shows octal, symbolic, and a per-class breakdown of read/write/execute permissions for owner, group, and other.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `permission` | Yes | Octal (e.g. 755) or symbolic (e.g. rwxr-xr-x) permission |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Convert octal to symbolic
openGyver chmod calc 755

# Convert another common permission
openGyver chmod calc 644

# Convert symbolic to octal
openGyver chmod calc rwxr-xr-x

# Convert read-only permission
openGyver chmod calc rw-r--r--

# Full permissions
openGyver chmod calc 777

# Restrictive permissions with JSON output
openGyver chmod calc 600 --json

# Symbolic input with JSON output
openGyver chmod calc rwx------ --json
```

#### JSON Output Format

```json
{
  "input": "755",
  "octal": "755",
  "symbolic": "rwxr-xr-x",
  "owner": {
    "class": "Owner",
    "read": "yes",
    "write": "yes",
    "execute": "yes"
  },
  "group": {
    "class": "Group",
    "read": "yes",
    "write": "no",
    "execute": "yes"
  },
  "other": {
    "class": "Other",
    "read": "yes",
    "write": "no",
    "execute": "yes"
  }
}
```

---

### umask

Calculate the resulting default permissions when a umask is applied. Files start with base permission 666 and directories with 777. The umask is subtracted to produce the effective permission. Input must be an octal value (e.g. 022).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `value` | Yes | Umask value in octal (e.g. 022, 077) |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Common default umask
openGyver chmod umask 022

# Restrictive umask
openGyver chmod umask 077

# Permissive umask
openGyver chmod umask 000

# Moderate umask
openGyver chmod umask 027

# JSON output for scripting
openGyver chmod umask 022 --json

# Very restrictive umask
openGyver chmod umask 077 --json
```

#### JSON Output Format

```json
{
  "umask": "022",
  "file": {
    "octal": "644",
    "symbolic": "rw-r--r--"
  },
  "directory": {
    "octal": "755",
    "symbolic": "rwxr-xr-x"
  }
}
```

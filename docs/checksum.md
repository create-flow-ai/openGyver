# checksum

Calculate and verify file checksums using various hash algorithms.

## Usage

```bash
openGyver checksum [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--algorithm` | `-a` | string | `sha256` | Hash algorithm: md5, sha1, sha256, sha512, crc32 |
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--help` | `-h` | | | Help for checksum |

## Subcommands

### calc

Calculate the checksum of a file using the specified algorithm. The file is read in 32 KB chunks for memory efficiency, making it suitable for large files. Supported algorithms: md5, sha1, sha256 (default), sha512, crc32.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file` | Yes | Path to the file to checksum |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Calculate SHA-256 checksum (default algorithm)
openGyver checksum calc document.pdf

# Calculate MD5 checksum
openGyver checksum calc archive.tar.gz --algorithm md5

# Calculate SHA-512 checksum
openGyver checksum calc backup.iso --algorithm sha512

# Calculate CRC-32 checksum
openGyver checksum calc firmware.bin --algorithm crc32

# JSON output for scripting
openGyver checksum calc myfile.zip --json

# SHA-1 checksum with JSON output
openGyver checksum calc release.tar.gz --algorithm sha1 --json

# Quick CRC-32 check
openGyver checksum calc config.toml -a crc32
```

#### JSON Output Format

```json
{
  "file": "document.pdf",
  "algorithm": "sha256",
  "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
```

---

### verify

Verify a file's checksum against a known expected hash. Prints "MATCH" if the computed checksum equals the expected hash, or "MISMATCH" if they differ. The comparison is case-insensitive. Supported algorithms: md5, sha1, sha256 (default), sha512, crc32.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file` | Yes | Path to the file to verify |
| `expected-hash` | Yes | The expected checksum to compare against |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Verify SHA-256 checksum (default)
openGyver checksum verify document.pdf abc123def456

# Verify with a specific algorithm
openGyver checksum verify archive.tar.gz e3b0c44298fc1c14 --algorithm sha256

# Verify CRC-32 checksum
openGyver checksum verify firmware.bin 3610a686 --algorithm crc32

# Verify MD5 checksum
openGyver checksum verify download.zip d41d8cd98f00b204e9800998ecf8427e --algorithm md5

# JSON output for automation
openGyver checksum verify release.tar.gz abc123 --json

# Verify SHA-512 checksum with JSON
openGyver checksum verify backup.iso 9b71d224bd62f378 --algorithm sha512 --json

# Short flag for algorithm
openGyver checksum verify config.yaml abc123 -a sha1
```

#### JSON Output Format

```json
{
  "file": "document.pdf",
  "algorithm": "sha256",
  "expected": "abc123def456",
  "actual": "e3b0c44298fc1c149afbf4c8996fb924",
  "match": false
}
```

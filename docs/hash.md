# hash

Compute cryptographic hashes and checksums for strings or files.

## Usage

```bash
openGyver hash [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | | Hash a file's contents instead of a string argument |
| `--json` | `-j` | bool | `false` | Output result as JSON |
| `--uppercase` | `-u` | bool | `false` | Output hex digest in uppercase |
| `--help` | `-h` | | | Help for hash |

## Subcommands

### md5

Compute the MD5 message digest of a string or file. MD5 produces a 128-bit (16-byte) hash value, typically rendered as a 32-character hexadecimal string. Note: MD5 is cryptographically broken and should not be used for security purposes. It is still useful for checksums and non-security fingerprinting.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to hash |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Hash a simple string
openGyver hash md5 "hello"

# Hash the contents of a file
openGyver hash md5 --file document.pdf

# Output in uppercase hex
openGyver hash md5 "hello world" --uppercase

# Get JSON output
openGyver hash md5 "hello" --json

# Hash a file and get uppercase JSON output
openGyver hash md5 --file /etc/hosts --uppercase --json

# Pipe-friendly: hash a known config file
openGyver hash md5 --file ~/.bashrc
```

#### JSON Output Format

```json
{
  "algorithm": "md5",
  "input": "hello",
  "hash": "5d41402abc4b2a76b9719d911017c592"
}
```

---

### sha1

Compute the SHA-1 message digest of a string or file. SHA-1 produces a 160-bit (20-byte) hash value. Note: SHA-1 is considered weak against well-funded attackers and is deprecated for digital signatures. Still acceptable for non-security checksums.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to hash |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Hash a string
openGyver hash sha1 "hello"

# Hash a compressed archive
openGyver hash sha1 --file archive.tar.gz

# Uppercase digest
openGyver hash sha1 "test data" --uppercase

# JSON output for scripting
openGyver hash sha1 "hello" --json

# Verify file integrity with uppercase hex
openGyver hash sha1 --file release.zip --uppercase

# Compare two strings by hashing each
openGyver hash sha1 "version1"
```

#### JSON Output Format

```json
{
  "algorithm": "sha1",
  "input": "hello",
  "hash": "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
}
```

---

### sha256

Compute the SHA-256 message digest of a string or file. SHA-256 is part of the SHA-2 family and produces a 256-bit (32-byte) hash value. It is widely used for data integrity, digital signatures, and blockchain applications.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to hash |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Hash a string
openGyver hash sha256 "hello"

# Hash a file with JSON output
openGyver hash sha256 --file release.zip --json

# Hash the hosts file
openGyver hash sha256 --file /etc/hosts

# Uppercase output
openGyver hash sha256 "secret" --uppercase

# Quick integrity check of a binary
openGyver hash sha256 --file /usr/local/bin/openGyver

# Combine with JSON for automation
openGyver hash sha256 "deployment-v2.3.1" --json --uppercase
```

#### JSON Output Format

```json
{
  "algorithm": "sha256",
  "input": "hello",
  "hash": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
}
```

---

### sha384

Compute the SHA-384 message digest of a string or file. SHA-384 is a truncated version of SHA-512 and produces a 384-bit (48-byte) hash value. It offers a higher security margin than SHA-256 while being faster on 64-bit platforms.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to hash |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Hash a string
openGyver hash sha384 "hello"

# Hash an ISO image
openGyver hash sha384 --file image.iso

# Get uppercase JSON output
openGyver hash sha384 "sensitive data" --json --uppercase

# Hash a configuration file
openGyver hash sha384 --file config.yaml

# Pipe-friendly verification
openGyver hash sha384 --file firmware.bin

# Hash with all formatting options
openGyver hash sha384 "test" --uppercase --json
```

#### JSON Output Format

```json
{
  "algorithm": "sha384",
  "input": "hello",
  "hash": "59e1748777448c69de6b800d7a33bbfb9ff1b463e44354c3553bcdb9c666fa90125a3c79f90397bdf5f6a13de828684f"
}
```

---

### sha512

Compute the SHA-512 message digest of a string or file. SHA-512 is part of the SHA-2 family and produces a 512-bit (64-byte) hash value. It is the strongest member of SHA-2, preferred when maximum security margin is desired.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to hash |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Hash a string
openGyver hash sha512 "hello"

# Hash a backup file in uppercase
openGyver hash sha512 --file backup.tar.gz --uppercase

# JSON output for automation
openGyver hash sha512 "critical-key" --json

# Hash a secret with uppercase
openGyver hash sha512 "secret" --uppercase

# Verify download integrity
openGyver hash sha512 --file downloaded-binary

# Full output with all options
openGyver hash sha512 --file release.zip --json --uppercase
```

#### JSON Output Format

```json
{
  "algorithm": "sha512",
  "input": "hello",
  "hash": "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"
}
```

---

### hmac

Compute an HMAC (Hash-based Message Authentication Code) of a string or file. HMAC combines a cryptographic hash function with a secret key to produce a message authentication code. This is used to verify both data integrity and authenticity.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to compute HMAC for |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--key` | `-k` | string | | Secret key for HMAC (required) |
| `--algorithm` | `-a` | string | `sha256` | Hash algorithm: md5, sha1, sha256, sha384, sha512 |
| `--help` | `-h` | | | Help for hmac |

#### Examples

```bash
# Basic HMAC with default SHA-256
openGyver hash hmac "hello" --key mysecret

# HMAC with SHA-512 algorithm
openGyver hash hmac "hello" --key mysecret --algorithm sha512

# HMAC a file's contents
openGyver hash hmac --file payload.json --key apikey123

# HMAC with MD5 and uppercase output
openGyver hash hmac "message" --key secret --algorithm md5 --uppercase

# HMAC for webhook signature verification
openGyver hash hmac "webhook-body" --key webhook_secret --algorithm sha256 --json

# HMAC with SHA-1 for compatibility
openGyver hash hmac "data" --key oldkey --algorithm sha1

# Full JSON output
openGyver hash hmac "test" --key mykey --json --uppercase
```

#### JSON Output Format

```json
{
  "algorithm": "hmac-sha256",
  "input": "hello",
  "hash": "88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b"
}
```

---

### bcrypt

Hash a password using bcrypt, or verify a password against a bcrypt hash. Bcrypt is a password-hashing function designed to be computationally expensive, making brute-force attacks impractical. The `--rounds` flag controls the cost factor (higher = slower + more secure).

**Modes:**
- **Hash mode (default):** Generates a bcrypt hash of the input string.
- **Verify mode (`--verify`):** Checks whether the input string matches the given bcrypt hash. Prints "match" or "no match" and exits with code 0 or 1.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | Yes | The password string to hash or verify |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--rounds` | `-r` | int | `10` | Bcrypt cost factor (4-31) |
| `--verify` | `-v` | string | | Bcrypt hash to verify the input against |
| `--help` | `-h` | | | Help for bcrypt |

#### Examples

```bash
# Hash a password with default cost (10)
openGyver hash bcrypt "mypassword"

# Hash with higher cost factor
openGyver hash bcrypt "mypassword" --rounds 12

# Verify a password against a known hash
openGyver hash bcrypt "mypassword" --verify '$2a$10$N9qo8uLOickgx2ZMRZoMye...'

# Hash with JSON output
openGyver hash bcrypt "mypassword" --json

# Use maximum security cost factor
openGyver hash bcrypt "supersecret" --rounds 14

# Verify and get JSON result
openGyver hash bcrypt "test" --verify '$2a$10$...' --json

# Minimum cost factor for testing
openGyver hash bcrypt "test" --rounds 4
```

#### JSON Output Format (hash mode)

```json
{
  "algorithm": "bcrypt",
  "input": "mypassword",
  "hash": "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
  "rounds": 10
}
```

#### JSON Output Format (verify mode)

```json
{
  "algorithm": "bcrypt",
  "input": "mypassword",
  "hash": "$2a$10$N9qo8uLOickgx2ZMRZoMye...",
  "match": true
}
```

---

### crc32

Compute the CRC-32 checksum of a string or file using the IEEE polynomial. CRC-32 is a non-cryptographic checksum used for error detection in network transmissions and file integrity checks.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to checksum |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Checksum a string
openGyver hash crc32 "hello"

# Checksum a firmware binary
openGyver hash crc32 --file firmware.bin

# Uppercase hex output
openGyver hash crc32 "data" --uppercase

# JSON output including decimal value
openGyver hash crc32 "hello" --json

# Checksum a configuration file
openGyver hash crc32 --file config.toml

# Checksum with all options
openGyver hash crc32 --file archive.zip --json --uppercase
```

#### JSON Output Format

```json
{
  "algorithm": "crc32",
  "input": "hello",
  "hash": "3610a686",
  "decimal": 907060870
}
```

---

### adler32

Compute the Adler-32 checksum of a string or file. Adler-32 is a non-cryptographic checksum that is faster but less reliable than CRC-32. It is used by the zlib compression library.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | The string to checksum |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Checksum a string
openGyver hash adler32 "hello"

# Checksum a data file
openGyver hash adler32 --file data.bin

# Uppercase output
openGyver hash adler32 "test" --uppercase

# JSON output including decimal value
openGyver hash adler32 "hello" --json

# Quick checksum for a small file
openGyver hash adler32 --file README.md

# Full options
openGyver hash adler32 --file payload.dat --json --uppercase
```

#### JSON Output Format

```json
{
  "algorithm": "adler32",
  "input": "hello",
  "hash": "062c0215",
  "decimal": 103547413
}
```

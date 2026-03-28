# crypto

A collection of cryptographic utilities for encryption, key generation, and certificate management.

## Usage

```bash
openGyver crypto [command] [flags]
```

## Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | `false` | Show help for crypto |
| `--json` | `-j` | bool | `false` | Output as JSON (available to all subcommands) |

## Subcommands

### aes

Encrypt or decrypt data using AES-256-GCM. The `--key` flag is required and accepts either a 64-character hex string (raw 256-bit key) or any other string treated as a passphrase (key derived via PBKDF2 with 600,000 iterations and SHA-256).

Encryption output is base64-encoded and includes the nonce (first 12 bytes) and, when using a passphrase, a 16-byte salt prefix.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `plaintext` or `ciphertext` | Yes | Text to encrypt, or base64-encoded ciphertext to decrypt |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--key` | | string | `""` | Encryption key: hex string (64 chars) or passphrase (required) |
| `--decrypt` | `-d` | bool | `false` | Decrypt instead of encrypt |
| `--help` | `-h` | bool | `false` | Show help for aes |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Examples

```bash
# Encrypt with a passphrase
openGyver crypto aes "secret message" --key "my passphrase"

# Decrypt (pass the base64 output from encryption)
openGyver crypto aes "BASE64CIPHERTEXT..." --key "my passphrase" --decrypt

# Short flag for decrypt
openGyver crypto aes "BASE64..." --key "my passphrase" -d

# Encrypt with a raw 256-bit hex key (64 hex chars = 32 bytes)
openGyver crypto aes "hello" --key 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

# JSON output
openGyver crypto aes "hello" --key pass --json

# Encrypt and capture output in a variable
encrypted=$(openGyver crypto aes "sensitive data" --key "mykey")

# Round-trip: encrypt then decrypt
openGyver crypto aes "$(openGyver crypto aes 'test' --key pw)" --key pw -d

# Pipe JSON to extract just the ciphertext
openGyver crypto aes "hello world" --key mypass -j | jq -r '.output'
```

#### JSON Output Format (Encrypt)

```json
{
  "input": "secret message",
  "output": "BASE64ENCODEDCIPHERTEXT...",
  "algorithm": "aes-256-gcm"
}
```

#### JSON Output Format (Decrypt)

```json
{
  "input": "BASE64ENCODEDCIPHERTEXT...",
  "output": "secret message",
  "algorithm": "aes-256-gcm"
}
```

---

### rsa

Generate an RSA private/public key pair in PEM format (PKCS#8 private key, PKIX public key).

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--bits` | | int | `2048` | RSA key size in bits (`2048`, `3072`, `4096`) |
| `--output-dir` | | string | `""` | Directory to write key files (`private.pem` and `public.pem`) |
| `--help` | `-h` | bool | `false` | Show help for rsa |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Examples

```bash
# Generate a default 2048-bit RSA key pair (printed to stdout)
openGyver crypto rsa

# Generate a 4096-bit key pair
openGyver crypto rsa --bits 4096

# Write keys to files in a directory
openGyver crypto rsa --output-dir ./keys

# Generate 4096-bit keys and write to directory
openGyver crypto rsa --bits 4096 --output-dir ./keys

# JSON output (keys included as strings)
openGyver crypto rsa --bits 4096 --json

# Generate and save private key only
openGyver crypto rsa | head -n $(grep -c '' <<< "$(openGyver crypto rsa 2>/dev/null)") > private.pem

# Extract public key from JSON
openGyver crypto rsa -j | jq -r '.public_key'

# Generate 3072-bit key pair (good balance of security and performance)
openGyver crypto rsa --bits 3072
```

#### JSON Output Format

```json
{
  "algorithm": "RSA",
  "bits": 4096,
  "private_key": "(PEM-encoded private key)",
  "public_key": "(PEM-encoded public key)"
}
```

---

### sshkey

Generate an SSH key pair in OpenSSH format.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--type` | | string | `"ed25519"` | Key type: `ed25519` or `rsa` |
| `--comment` | | string | `""` | Key comment (e.g., `user@host`) |
| `--help` | `-h` | bool | `false` | Show help for sshkey |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Supported Key Types

| Type | Description |
|------|-------------|
| `ed25519` | Fast, small, modern (default). Recommended for most use cases. |
| `rsa` | 4096-bit RSA. Use for compatibility with older systems. |

#### Examples

```bash
# Generate a default ed25519 SSH key pair
openGyver crypto sshkey

# Generate an RSA SSH key pair
openGyver crypto sshkey --type rsa

# Generate with a comment
openGyver crypto sshkey --type ed25519 --comment "deploy@prod"

# JSON output
openGyver crypto sshkey --json

# Generate and save to files
openGyver crypto sshkey --comment "me@laptop" > id_ed25519 2>&1

# Extract just the public key from JSON
openGyver crypto sshkey -j | jq -r '.public_key'

# Generate RSA key with comment for server access
openGyver crypto sshkey --type rsa --comment "admin@server"

# Generate key and extract public portion
openGyver crypto sshkey --type ed25519 --comment "ci@github" -j | jq -r '.public_key'
```

#### JSON Output Format

```json
{
  "type": "ed25519",
  "comment": "deploy@prod",
  "private_key": "(OpenSSH private key)",
  "public_key": "(ssh-ed25519 public key string)"
}
```

---

### cert

Generate a self-signed X.509 TLS certificate and private key in PEM format. The certificate is signed with ECDSA P-256 for fast generation and small size.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--cn` | | string | `""` | Common Name (required) |
| `--days` | | int | `365` | Certificate validity in days |
| `--output-dir` | | string | `""` | Directory to write `cert.pem` and `key.pem` files |
| `--help` | `-h` | bool | `false` | Show help for cert |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Examples

```bash
# Generate a self-signed cert for a domain
openGyver crypto cert --cn example.com

# Short-lived cert for local development
openGyver crypto cert --cn localhost --days 30

# Wildcard certificate valid for 2 years, saved to disk
openGyver crypto cert --cn "*.example.com" --days 730 --output-dir ./certs

# JSON output
openGyver crypto cert --cn myapp.local --json

# Generate cert and save to specific directory
openGyver crypto cert --cn api.example.com --days 365 --output-dir /etc/ssl/custom

# Extract just the certificate from JSON
openGyver crypto cert --cn test.local -j | jq -r '.certificate'

# Generate cert for internal service
openGyver crypto cert --cn "internal.corp" --days 3650

# Check the validity dates from JSON output
openGyver crypto cert --cn example.com -j | jq '{not_before, not_after}'
```

#### JSON Output Format

```json
{
  "common_name": "example.com",
  "not_before": "2026-03-28T12:00:00Z",
  "not_after": "2027-03-28T12:00:00Z",
  "certificate": "(PEM-encoded certificate)",
  "private_key": "(PEM-encoded EC private key)"
}
```

---

### csr

Generate a PEM-encoded Certificate Signing Request and a new ECDSA P-256 private key.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--cn` | | string | `""` | Common Name (required) |
| `--org` | | string | `""` | Organization name |
| `--country` | | string | `""` | Country code (e.g., `US`) |
| `--output` | | string | `""` | File path to write the CSR |
| `--help` | `-h` | bool | `false` | Show help for csr |
| `--json` | `-j` | bool | `false` | Output as JSON (inherited) |

#### Examples

```bash
# Generate a basic CSR
openGyver crypto csr --cn example.com

# CSR with organization and country
openGyver crypto csr --cn example.com --org "Acme Inc" --country US

# Write CSR to a file (private key still printed to stdout)
openGyver crypto csr --cn example.com --output request.pem

# JSON output with all fields
openGyver crypto csr --cn example.com --json

# Full CSR with all metadata
openGyver crypto csr --cn "api.example.com" --org "Example Corp" --country DE --output api.csr

# Extract just the CSR PEM from JSON
openGyver crypto csr --cn example.com -j | jq -r '.csr'

# Extract just the private key from JSON
openGyver crypto csr --cn example.com -j | jq -r '.private_key'

# Generate CSR for wildcard domain
openGyver crypto csr --cn "*.example.com" --org "Example Inc" --country US
```

#### JSON Output Format

```json
{
  "common_name": "example.com",
  "org": "Acme Inc",
  "country": "US",
  "csr": "(PEM-encoded CSR)",
  "private_key": "(PEM-encoded EC private key)"
}
```

## Notes

- All subcommands support `--json` / `-j` for machine-readable output, inherited from the parent command.
- AES encryption uses AES-256-GCM with PBKDF2 key derivation (600,000 iterations, SHA-256) when a passphrase is provided.
- RSA key pairs are generated in PKCS#8 (private) and PKIX (public) PEM format.
- SSH keys are generated in OpenSSH format, compatible with `ssh-keygen` output.
- Certificates use ECDSA P-256 for both `cert` and `csr` subcommands, providing fast key generation and small key sizes.
- Self-signed certificates include the CN as a DNS SAN (Subject Alternative Name) and have ServerAuth extended key usage.
- Private keys written to disk use `0600` permissions for security.

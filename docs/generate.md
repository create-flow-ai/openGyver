# generate

Generate cryptographically random passwords, passphrases, API keys, secrets, OTP tokens, and various ID formats. All randomness is sourced from `crypto/rand`.

## Usage

```bash
openGyver generate [subcommand] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | bool | | Show help for the command |

## Subcommands

### password

Generate cryptographically random passwords. By default passwords include uppercase, lowercase, digits, and special characters. Disable character classes with the `--no-*` flags.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--length` | | int | `16` | Password length |
| `--count` | | int | `1` | Number of passwords to generate |
| `--no-upper` | | bool | `false` | Exclude uppercase letters (A-Z) |
| `--no-lower` | | bool | `false` | Exclude lowercase letters (a-z) |
| `--no-digits` | | bool | `false` | Exclude digits (0-9) |
| `--no-special` | | bool | `false` | Exclude special characters |

#### Examples

```bash
# Generate a default 16-character password
openGyver generate password

# Generate a 24-character password
openGyver generate password --length 24

# Generate a password without special characters
openGyver generate password --no-special

# Generate an alphanumeric-only password (no special, no upper)
openGyver generate password --no-upper --no-digits

# Generate 10 passwords at once
openGyver generate password --count 10

# Generate a long password with JSON output
openGyver generate password --length 32 --json

# Generate 5 passwords without special characters
openGyver generate password --length 24 --no-special --count 5

# Digits-only PIN
openGyver generate password --length 6 --no-upper --no-lower --no-special
```

#### JSON Output Format

```json
{
  "passwords": ["aB3$kL9!mN2@pQ5&"],
  "length": 16,
  "count": 1,
  "charset": "uppercase, lowercase, digits, special"
}
```

---

### passphrase

Generate random passphrases by picking words from the EFF short diceware word list (~1296 words). Each word is chosen independently using `crypto/rand`.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--words` | | int | `4` | Number of words in the passphrase |
| `--separator` | | string | `"-"` | Word separator |
| `--count` | | int | `1` | Number of passphrases to generate |

#### Examples

```bash
# Generate a default 4-word passphrase
openGyver generate passphrase

# Generate a 6-word passphrase
openGyver generate passphrase --words 6

# Use a dot separator
openGyver generate passphrase --separator "."

# Generate multiple passphrases with JSON output
openGyver generate passphrase --words 5 --count 3 --json

# Use an underscore separator
openGyver generate passphrase --separator "_"

# Long passphrase for high security
openGyver generate passphrase --words 8

# Space-separated passphrase
openGyver generate passphrase --separator " "

# Generate 10 passphrases for a team
openGyver generate passphrase --words 4 --count 10
```

#### JSON Output Format

```json
{
  "passphrases": ["correct-horse-battery-staple"],
  "words": 4,
  "separator": "-",
  "count": 1
}
```

---

### string

Generate cryptographically random strings from a chosen character set.

**Available charsets:**

| Charset | Characters |
|---------|-----------|
| `alpha` | A-Z a-z |
| `alphanumeric` | A-Z a-z 0-9 (default) |
| `hex` | 0-9 a-f |
| `base64` | A-Z a-z 0-9 + / = |
| `custom` | User-supplied alphabet via `--custom` |

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--length` | | int | `32` | String length |
| `--charset` | | string | `"alphanumeric"` | Character set: alpha, alphanumeric, hex, base64, custom |
| `--custom` | | string | | Custom alphabet (use with `--charset=custom`) |
| `--count` | | int | `1` | Number of strings to generate |

#### Examples

```bash
# Generate a default alphanumeric string (32 chars)
openGyver generate string

# Generate a 64-character string
openGyver generate string --length 64

# Generate a hex string (great for tokens)
openGyver generate string --charset hex --length 32

# Generate a binary string using custom alphabet
openGyver generate string --charset custom --custom "01" --length 16

# Generate multiple strings with JSON output
openGyver generate string --count 5 --json

# Alpha-only string
openGyver generate string --charset alpha --length 20

# Base64 string
openGyver generate string --charset base64 --length 44

# Custom alphabet: only vowels
openGyver generate string --charset custom --custom "aeiou" --length 10
```

#### JSON Output Format

```json
{
  "strings": ["a8Kf2Lm9Np3Qr5St7Uv0Wx4Yz1Bc6De"],
  "length": 32,
  "charset": "alphanumeric",
  "count": 1
}
```

---

### apikey

Generate an API key with a random base62 portion (A-Z a-z 0-9). Optionally add a prefix (e.g., `myapp_key_`, `myapp_pub_`). The `--length` flag controls the length of the random portion only.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--length` | | int | `32` | Length of the random portion |
| `--prefix` | | string | | Optional prefix (e.g., `myapp_key_`) |

#### Examples

```bash
# Generate a default API key
openGyver generate apikey

# JSON output
openGyver generate apikey --json

# With a prefix
openGyver generate apikey --prefix myapp_key_

# Shorter random portion
openGyver generate apikey --length 16

# Test key prefix
openGyver generate apikey --prefix myapp_pub_

# Long random portion for extra security
openGyver generate apikey --length 64

# Custom service prefix
openGyver generate apikey --prefix myapp_

# Generate and copy to clipboard (macOS)
openGyver generate apikey | pbcopy
```

#### JSON Output Format

```json
{
  "apikey": "myapp_key_a8Kf2Lm9Np3Qr5St7Uv0Wx4Yz1Bc6De",
  "length": 32
}
```

---

### secret

Generate a cryptographically random secret key and output it as a hex-encoded string. The `--length` flag specifies the number of random bytes (the hex output will be twice as long).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--length` | | int | `64` | Number of random bytes (hex output is 2x this length) |

#### Examples

```bash
# Generate a default 64-byte secret (128 hex chars)
openGyver generate secret

# Generate a 32-byte secret (64 hex chars)
openGyver generate secret --length 32

# JSON output
openGyver generate secret --json

# 16-byte secret for a HMAC key
openGyver generate secret --length 16

# 256-byte secret for maximum entropy
openGyver generate secret --length 256

# Generate and set as environment variable
export SECRET_KEY=$(openGyver generate secret --length 32)

# Generate and copy to clipboard (macOS)
openGyver generate secret | pbcopy

# Short secret for a session token
openGyver generate secret --length 24
```

#### JSON Output Format

```json
{
  "secret": "a1b2c3d4e5f6...128 hex characters...",
  "bytes": 64
}
```

---

### otp

Generate a random TOTP (Time-based One-Time Password) secret encoded in base32, along with an `otpauth://` URI suitable for QR codes and authenticator apps.

The secret is 20 bytes (160 bits) of `crypto/rand` entropy, matching the recommended length for HMAC-SHA1 TOTP.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--issuer` | | string | | Issuer name for the OTP URI |
| `--account` | | string | | Account name for the OTP URI |

#### Examples

```bash
# Generate with issuer and account
openGyver generate otp --issuer MyApp --account user@example.com

# Generate for a GitHub-like service
openGyver generate otp --issuer GitHub --account octocat

# JSON output
openGyver generate otp --json

# Generate for a production service
openGyver generate otp --issuer "Production API" --account admin@corp.com

# Generate without issuer/account (secret only)
openGyver generate otp

# Generate and pipe the URI to a QR code generator
openGyver generate otp --issuer MyApp --account test@test.com -j | jq -r '.uri'

# Generate for multiple services
openGyver generate otp --issuer Slack --account dev@company.com
openGyver generate otp --issuer AWS --account root@company.com

# JSON output for integration
openGyver generate otp --issuer MyService --account hello@world.com -j
```

#### JSON Output Format

```json
{
  "secret": "JBSWY3DPEHPK3PXP",
  "uri": "otpauth://totp/MyApp:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=MyApp",
  "issuer": "MyApp",
  "account": "user@example.com"
}
```

---

### nanoid

Generate compact, URL-friendly unique identifiers using the Nano ID algorithm. The default alphabet is `A-Za-z0-9_-` (64 characters), matching the canonical nanoid specification.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--length` | | int | `21` | ID length |
| `--alphabet` | | string | `A-Za-z0-9_-` | Custom alphabet |

#### Examples

```bash
# Generate a default 21-character Nano ID
openGyver generate nanoid

# Generate a shorter 12-character ID
openGyver generate nanoid --length 12

# Use a hex alphabet
openGyver generate nanoid --alphabet "0123456789abcdef"

# JSON output
openGyver generate nanoid --json

# Very short ID for URL slugs
openGyver generate nanoid --length 8

# Long ID for maximum collision resistance
openGyver generate nanoid --length 36

# Custom numeric-only alphabet
openGyver generate nanoid --alphabet "0123456789" --length 10

# Generate and use as a filename
openGyver generate nanoid --length 16
```

#### JSON Output Format

```json
{
  "id": "V1StGXR8_Z5jdHi6B-myT",
  "length": 21,
  "alphabet": "A-Za-z0-9_-"
}
```

---

### snowflake

Generate a Twitter-style snowflake ID -- a 64-bit integer composed of:

- **41 bits:** milliseconds since a custom epoch (2020-01-01T00:00:00Z)
- **10 bits:** node ID (randomly assigned per invocation)
- **12 bits:** sequence number

The resulting ID is sortable by creation time and unique within a single process.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | No arguments |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| (none beyond global) | | | | |

#### Examples

```bash
# Generate a snowflake ID
openGyver generate snowflake

# JSON output
openGyver generate snowflake --json

# Generate and use as a database ID
ID=$(openGyver generate snowflake)

# Generate multiple by running in a loop
for i in $(seq 5); do openGyver generate snowflake; done

# JSON output piped to jq
openGyver generate snowflake -j | jq '.id'

# Use in a script
openGyver generate snowflake | xargs -I{} echo "INSERT INTO items (id) VALUES ({});"
```

#### JSON Output Format

```json
{
  "id": 192482349234176,
  "timestamp_ms": 1711612800000,
  "node_id": 42,
  "sequence": 0
}
```

---

### shortid

Generate a short, URL-safe random identifier. Uses a base62 alphabet (A-Z a-z 0-9) for maximum density without special characters.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--length` | | int | `8` | ID length |

#### Examples

```bash
# Generate a default 8-character short ID
openGyver generate shortid

# Generate a 12-character ID
openGyver generate shortid --length 12

# JSON output
openGyver generate shortid --json

# Very short 4-character ID
openGyver generate shortid --length 4

# Longer 16-character ID for more uniqueness
openGyver generate shortid --length 16

# Use as a URL slug
openGyver generate shortid --length 6

# Generate and copy to clipboard (macOS)
openGyver generate shortid | pbcopy

# JSON output piped to jq
openGyver generate shortid -j | jq '.id'
```

#### JSON Output Format

```json
{
  "id": "a8Kf2Lm9",
  "length": 8
}
```

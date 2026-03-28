# encode

Encode and decode text using common encoding schemes.

## Usage

```bash
openGyver encode [subcommand] [input] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output structured JSON: `{"input","output","encoding"}` |
| `--file` | `-f` | string | | Read input from a file instead of a positional argument |
| `--help` | `-h` | bool | | Show help for the command |

## Subcommands

### base64

Base64 encode or decode text using the standard Base64 alphabet (RFC 4648).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or Base64 string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode a string to Base64
openGyver encode base64 "hello world"

# Decode a Base64 string
openGyver encode base64 -d "aGVsbG8gd29ybGQ="

# Encode from a file
openGyver encode base64 --file input.txt

# Decode from a file
openGyver encode base64 -d --file encoded.txt

# JSON output
openGyver encode base64 "hello" --json

# Encode binary-safe content from file
openGyver encode base64 -f image.png

# Decode and redirect to a file
openGyver encode base64 -d "aGVsbG8=" > output.txt

# Pipe into the encoder
echo "secret data" | openGyver encode base64
```

#### JSON Output Format

```json
{
  "input": "hello world",
  "output": "aGVsbG8gd29ybGQ=",
  "encoding": "base64"
}
```

---

### base32

Base32 encode or decode text using the standard Base32 alphabet (RFC 4648).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or Base32 string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode to Base32
openGyver encode base32 "hello world"

# Decode from Base32
openGyver encode base32 -d "NBSWY3DPEB3W64TMMQ======"

# Encode from a file
openGyver encode base32 --file data.txt

# JSON output
openGyver encode base32 "test" --json

# Decode from file
openGyver encode base32 -d -f encoded.txt

# Pipe input
echo "encode me" | openGyver encode base32

# Decode and save
openGyver encode base32 -d "NBSWY3DP" > decoded.txt

# Encode with JSON wrapper
openGyver encode base32 "hello" -j
```

#### JSON Output Format

```json
{
  "input": "hello world",
  "output": "NBSWY3DPEB3W64TMMQ======",
  "encoding": "base32"
}
```

---

### base58

Base58 encode or decode text using the Bitcoin alphabet (`123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz`). This encoding avoids visually ambiguous characters (0, O, I, l).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or Base58 string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode to Base58
openGyver encode base58 "hello"

# Decode from Base58
openGyver encode base58 -d "Cn8eVZg"

# JSON output
openGyver encode base58 "test" --json

# Encode from file
openGyver encode base58 -f data.txt

# Decode from file
openGyver encode base58 -d -f encoded.txt

# Use for short URL-safe identifiers
openGyver encode base58 "my-identifier"

# Pipe data to encode
echo "bitcoin" | openGyver encode base58

# Decode and verify
openGyver encode base58 -d "StV1DL6CwTryKyV"
```

#### JSON Output Format

```json
{
  "input": "hello",
  "output": "Cn8eVZg",
  "encoding": "base58"
}
```

---

### url

URL percent-encoding encode or decode text. Encodes special characters as `%XX` hex sequences as required by RFC 3986.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or URL-encoded string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode a URL query string
openGyver encode url "hello world"

# Decode a percent-encoded string
openGyver encode url -d "hello+world"

# Encode special characters
openGyver encode url "name=John Doe&city=New York"

# Decode a full URL parameter
openGyver encode url -d "name%3DJohn%20Doe%26city%3DNew%20York"

# JSON output
openGyver encode url "hello world" --json

# Encode from file
openGyver encode url -f query.txt

# Pipe URL-encoded data to decode
echo "hello%20world" | openGyver encode url -d

# Encode Unicode characters
openGyver encode url "cafe et creme"
```

#### JSON Output Format

```json
{
  "input": "hello world",
  "output": "hello+world",
  "encoding": "url"
}
```

---

### html

HTML entity encode or decode text. Encodes `<`, `>`, `&`, `"`, and `'` into their HTML entity equivalents. Decodes all standard HTML entities back to their original characters.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or HTML-encoded string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode HTML special characters
openGyver encode html '<script>alert("xss")</script>'

# Decode HTML entities
openGyver encode html -d "&lt;div&gt;Hello&lt;/div&gt;"

# Encode an ampersand-heavy string
openGyver encode html "Tom & Jerry"

# JSON output
openGyver encode html "<b>bold</b>" --json

# Encode from file
openGyver encode html -f snippet.html

# Decode from file
openGyver encode html -d -f entities.txt

# Pipe HTML to encode
echo '<p class="test">Hello</p>' | openGyver encode html

# Decode named entities
openGyver encode html -d "&copy; 2024 &mdash; All rights reserved"
```

#### JSON Output Format

```json
{
  "input": "<b>bold</b>",
  "output": "&lt;b&gt;bold&lt;/b&gt;",
  "encoding": "html"
}
```

---

### hex

Hex encode or decode text. Each byte is represented as two hexadecimal characters.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or hex string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode text to hex
openGyver encode hex "hello"

# Decode hex to text
openGyver encode hex -d "68656c6c6f"

# Decode a hex color code to see raw bytes
openGyver encode hex -d "cafe"

# JSON output
openGyver encode hex "test" --json

# Encode from file
openGyver encode hex -f binary.dat

# Decode from file
openGyver encode hex -d -f hexdump.txt

# Pipe data through hex encoding
echo "secret" | openGyver encode hex

# Decode a longer hex string
openGyver encode hex -d "48656c6c6f20576f726c64"
```

#### JSON Output Format

```json
{
  "input": "hello",
  "output": "68656c6c6f",
  "encoding": "hex"
}
```

---

### binary

Binary (0/1) encode or decode text. Each byte is represented as an 8-bit binary string, with bytes separated by spaces.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or binary string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode text to binary
openGyver encode binary "Hi"

# Decode binary back to text
openGyver encode binary -d "01001000 01101001"

# Encode a single character
openGyver encode binary "A"

# JSON output
openGyver encode binary "OK" --json

# Encode from file
openGyver encode binary -f message.txt

# Decode from file
openGyver encode binary -d -f bits.txt

# Pipe text to binary encode
echo "hello" | openGyver encode binary

# Decode a longer binary string
openGyver encode binary -d "01001000 01100101 01101100 01101100 01101111"
```

#### JSON Output Format

```json
{
  "input": "Hi",
  "output": "01001000 01101001",
  "encoding": "binary"
}
```

---

### rot13

ROT13 cipher -- a symmetric letter substitution that replaces each letter with the letter 13 positions after it in the alphabet. Since it is symmetric, applying ROT13 twice returns the original text. No `--decode` flag is needed.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to apply ROT13 to |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| (none beyond global) | | | | ROT13 is symmetric; no decode flag needed |

#### Examples

```bash
# Encode with ROT13
openGyver encode rot13 "Hello"

# Decode (just apply ROT13 again)
openGyver encode rot13 "Uryyb"

# Encode a sentence
openGyver encode rot13 "The quick brown fox"

# JSON output
openGyver encode rot13 "Secret" --json

# From file
openGyver encode rot13 -f message.txt

# Non-alpha characters pass through unchanged
openGyver encode rot13 "Hello, World! 123"

# Pipe through ROT13
echo "spoiler alert" | openGyver encode rot13

# Double-apply to verify symmetry
openGyver encode rot13 "Gur dhvpx oebja sbk"
```

#### JSON Output Format

```json
{
  "input": "Hello",
  "output": "Uryyb",
  "encoding": "rot13"
}
```

---

### morse

Morse code encode or decode. Letters are separated by spaces, words are separated by `/`. Supports A-Z, 0-9, and common punctuation.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Text to encode or Morse code string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode text to Morse code
openGyver encode morse "SOS"

# Decode Morse code to text
openGyver encode morse -d "... --- ..."

# Encode a word with spaces (words separated by /)
openGyver encode morse "HELLO WORLD"

# Decode a message with word separators
openGyver encode morse -d ".... . .-.. .-.. --- / .-- --- .-. .-.. -.."

# JSON output
openGyver encode morse "OK" --json

# Encode from file
openGyver encode morse -f message.txt

# Decode from file
openGyver encode morse -d -f morse.txt

# Encode numbers
openGyver encode morse "42"
```

#### JSON Output Format

```json
{
  "input": "SOS",
  "output": "... --- ...",
  "encoding": "morse"
}
```

---

### punycode

Punycode (IDNA) encode or decode internationalized domain names. Converts Unicode domain labels to ASCII-compatible encoding and back.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `input` | No (use `--file` instead) | Domain name to encode or Punycode string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--decode` | `-d` | bool | `false` | Decode instead of encode |

#### Examples

```bash
# Encode a Unicode domain
openGyver encode punycode "munchen.de"

# Decode a Punycode domain
openGyver encode punycode -d "xn--mnchen-3ya.de"

# Encode a Japanese domain
openGyver encode punycode "example.jp"

# JSON output
openGyver encode punycode "cafe.fr" --json

# Encode from file
openGyver encode punycode -f domains.txt

# Decode from file
openGyver encode punycode -d -f punycode.txt

# Pipe a domain name
echo "example.com" | openGyver encode punycode

# Decode and verify
openGyver encode punycode -d "xn--nxasmq6b.xn--jxalpdlp"
```

#### JSON Output Format

```json
{
  "input": "munchen.de",
  "output": "xn--mnchen-3ya.de",
  "encoding": "punycode"
}
```

---

### jwt

Decode a JWT (JSON Web Token) payload without verification. Extracts and pretty-prints the payload section of the token. Does not validate the signature -- use this for inspection only.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `token` | No (use `--file` instead) | JWT token string to decode |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| (none beyond global) | | | | No decode flag; jwt always decodes |

#### Examples

```bash
# Decode a JWT token
openGyver encode jwt "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

# JSON output
openGyver encode jwt "eyJhbGciOi..." --json

# Decode from a file containing the token
openGyver encode jwt -f token.txt

# Pipe a token from clipboard or other command
pbpaste | openGyver encode jwt

# Decode token from an environment variable
openGyver encode jwt "$AUTH_TOKEN"

# Inspect a token from a curl response
curl -s https://api.example.com/auth | jq -r '.token' | openGyver encode jwt

# JSON wrapped output
openGyver encode jwt -f jwt.txt -j

# Quick token inspection
openGyver encode jwt "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJhdXRoMCJ9.signature"
```

#### JSON Output Format

```json
{
  "input": "eyJhbGciOiJIUzI1NiIs...",
  "output": "{\n  \"sub\": \"1234567890\",\n  \"name\": \"John Doe\",\n  \"iat\": 1516239022\n}",
  "encoding": "jwt"
}
```

# network

DNS lookups, public IP, WHOIS, CIDR calculation, URL parsing, HTTP status codes, User-Agent parsing, and JWT decoding.

## Usage

```bash
openGyver network [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for network |

## Subcommands

### dns

Perform a DNS lookup for the given domain and record type.

**Supported record types:** A (default), AAAA, MX, TXT, NS, CNAME, SOA, SRV, PTR.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `domain` | Yes | The domain name (or IP for PTR) to look up |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--type` | | string | `A` | Record type: A, AAAA, MX, TXT, NS, CNAME, SOA, SRV, PTR |
| `--help` | `-h` | | | Help for dns |

#### Examples

```bash
# Basic A record lookup
openGyver network dns example.com

# IPv6 address lookup
openGyver network dns --type AAAA example.com

# Mail exchange servers
openGyver network dns --type MX gmail.com

# TXT records (SPF, DKIM, etc.)
openGyver network dns --type TXT example.com

# Name servers
openGyver network dns --type NS example.com

# Start of Authority
openGyver network dns --type SOA example.com

# SRV record lookup
openGyver network dns --type SRV _sip._tcp.example.com

# Reverse DNS lookup
openGyver network dns --type PTR 8.8.8.8

# CNAME lookup
openGyver network dns --type CNAME www.example.com

# JSON output for scripting
openGyver network dns example.com --json
```

#### JSON Output Format

```json
{
  "domain": "example.com",
  "type": "A",
  "records": ["93.184.216.34"]
}
```

---

### ip

Query an external service to determine your public-facing IP address. Tries https://api.ipify.org first, falls back to https://ifconfig.me/ip.

#### Arguments

No arguments.

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Show your public IP
openGyver network ip

# Get JSON output
openGyver network ip --json

# Use in a script
MY_IP=$(openGyver network ip)

# Pipe to clipboard (macOS)
openGyver network ip | pbcopy

# Check IP for firewall configuration
openGyver network ip --json | jq -r '.ip'
```

#### JSON Output Format

```json
{
  "ip": "203.0.113.42"
}
```

---

### whois

Query the WHOIS database for domain registration information. Connects to whois.iana.org to find the authoritative WHOIS server for the TLD, then follows the referral to get full registration details.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `domain` | Yes | The domain name to query |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Basic WHOIS lookup
openGyver network whois example.com

# WHOIS with JSON output
openGyver network whois google.com --json

# Check domain registration details
openGyver network whois github.com

# Look up a country-code TLD domain
openGyver network whois example.co.uk

# Verify domain ownership
openGyver network whois mycompany.com --json

# Check expiration dates
openGyver network whois expiring-domain.com
```

#### JSON Output Format

```json
{
  "domain": "example.com",
  "raw": "Domain Name: EXAMPLE.COM\nRegistry Domain ID: ..."
}
```

---

### cidr

Calculate network information from a CIDR notation address. Outputs the network address, broadcast address, first and last usable IPs, total number of hosts, subnet mask, and wildcard mask. Supports both IPv4 and IPv6 CIDR notation.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `CIDR` | Yes | Network address in CIDR notation (e.g. `192.168.1.0/24`) |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Calculate a /24 subnet
openGyver network cidr 192.168.1.0/24

# Large private network
openGyver network cidr 10.0.0.0/8

# JSON output for automation
openGyver network cidr 172.16.0.0/12 --json

# Smaller subnet
openGyver network cidr 192.168.1.128/25

# Single host
openGyver network cidr 10.0.0.1/32

# Common office subnet
openGyver network cidr 192.168.0.0/16

# IPv6 CIDR calculation
openGyver network cidr 2001:db8::/32 --json
```

#### JSON Output Format

```json
{
  "cidr": "192.168.1.0/24",
  "network": "192.168.1.0",
  "broadcast": "192.168.1.255",
  "first_host": "192.168.1.1",
  "last_host": "192.168.1.254",
  "total_hosts": 254,
  "subnet_mask": "255.255.255.0",
  "wildcard_mask": "0.0.0.255"
}
```

---

### urlparse

Break a URL into its constituent parts: scheme, user, host, port, path, query parameters, and fragment.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `url` | Yes | The URL to parse |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Parse a full URL with all components
openGyver network urlparse "https://example.com:8080/path?q=search&lang=en#results"

# Parse an FTP URL with credentials
openGyver network urlparse "ftp://user:pass@files.example.com/pub/file.txt"

# JSON output with duplicate query params
openGyver network urlparse "https://example.com/path?a=1&a=2" --json

# Parse a simple URL
openGyver network urlparse "https://google.com"

# Parse a URL with a fragment
openGyver network urlparse "https://docs.example.com/guide#installation"

# Parse a localhost URL
openGyver network urlparse "http://localhost:3000/api/v1/users?page=2"

# Parse a complex query string
openGyver network urlparse "https://search.example.com/q?term=hello+world&sort=date&order=desc"
```

#### JSON Output Format

```json
{
  "scheme": "https",
  "host": "example.com",
  "port": "8080",
  "path": "/path",
  "query": {
    "q": ["search"],
    "lang": ["en"]
  },
  "fragment": "results"
}
```

---

### httpstatus

Display the name and description for a given HTTP status code. Covers all standard codes from RFC 9110 and common extensions.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `code` | Yes | The HTTP status code (e.g. `200`, `404`, `503`) |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Look up 200 OK
openGyver network httpstatus 200

# Look up 404 Not Found
openGyver network httpstatus 404

# JSON output for scripting
openGyver network httpstatus 503 --json

# The teapot status
openGyver network httpstatus 418

# Server error
openGyver network httpstatus 500

# Redirect status
openGyver network httpstatus 301

# Rate limiting
openGyver network httpstatus 429

# Created
openGyver network httpstatus 201
```

#### JSON Output Format

```json
{
  "code": 404,
  "name": "Not Found",
  "description": "The server cannot find the requested resource."
}
```

---

### useragent

Extract browser, version, operating system, and device type from a User-Agent string using pattern matching. Recognizes common browsers (Chrome, Firefox, Safari, Edge, Opera, etc.), operating systems (Windows, macOS, Linux, Android, iOS), and device types.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `user-agent-string` | Yes | The User-Agent string to parse |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Parse a Chrome on macOS User-Agent
openGyver network useragent "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

# Parse an iPhone User-Agent
openGyver network useragent "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"

# Parse curl with JSON output
openGyver network useragent "curl/8.1.2" --json

# Parse a Firefox User-Agent
openGyver network useragent "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/120.0"

# Parse a bot User-Agent
openGyver network useragent "Googlebot/2.1 (+http://www.google.com/bot.html)"

# Parse an Android User-Agent
openGyver network useragent "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 Chrome/120.0.0.0 Mobile Safari/537.36"
```

#### JSON Output Format

```json
{
  "browser": "Chrome",
  "version": "120.0.0.0",
  "os": "macOS",
  "device": "Desktop"
}
```

---

### jwt

Decode a JSON Web Token (JWT) by splitting it into its three parts (header, payload, signature), base64-decoding the header and payload, and displaying them as formatted JSON.

**Note:** This does NOT verify the signature. It is an inspection/debugging tool.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `token` | Yes | The JWT token string to decode |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Decode a JWT token
openGyver network jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

# Decode with JSON output
openGyver network jwt <token> --json

# Inspect an OAuth access token
openGyver network jwt "$ACCESS_TOKEN"

# Decode a token from a cookie value
openGyver network jwt "$SESSION_JWT" --json

# Debug an API authentication token
openGyver network jwt "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

# Quick check of token claims
openGyver network jwt "$TOKEN"
```

#### JSON Output Format

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "1234567890",
    "name": "John Doe",
    "iat": 1516239022
  },
  "signature": "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

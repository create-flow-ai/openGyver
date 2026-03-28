# phone

Phone number parser -- format, validate, and look up country dial codes.

## Usage

```bash
openGyver phone [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for phone |

## Subcommands

### format

Format a phone number into E.164, international, and national formats. Non-digit characters are stripped from the input. The country dial code prefix and leading national zero are automatically handled. Supports approximately 30 countries.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `phone-number` | Yes | The phone number to format (digits, spaces, dashes, and plus sign allowed) |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--country` | | string | `US` | Country ISO code (e.g. US, GB, JP) |

#### Examples

```bash
# Format a US number (default country)
openGyver phone format "2025551234"

# Format a UK number
openGyver phone format "+44 20 7946 0958" --country GB

# Format a Japanese number
openGyver phone format "03-1234-5678" --country JP

# Format with JSON output
openGyver phone format "2025551234" --json

# Format a German number
openGyver phone format "030 12345678" --country DE

# Format a French number
openGyver phone format "01 23 45 67 89" --country FR

# Format an Australian number
openGyver phone format "0412345678" --country AU

# Format a Brazilian number
openGyver phone format "11987654321" --country BR
```

#### JSON Output Format

```json
{
  "input": "2025551234",
  "country": "United States",
  "e164": "+12025551234",
  "international": "+1 202-555-1234",
  "national": "(202) 555-1234"
}
```

---

### validate

Basic phone number validation: checks digit count and prefix for the given country. Reports whether the number is valid and, if invalid, explains the reason (too few or too many digits).

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `phone-number` | Yes | The phone number to validate |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--country` | | string | `US` | Country ISO code (e.g. US, GB, JP) |

#### Examples

```bash
# Validate a US number
openGyver phone validate "2025551234" --country US

# Validate a Japanese number
openGyver phone validate "+81-03-1234-5678" --country JP

# Validate a UK number
openGyver phone validate "+44 20 7946 0958" --country GB

# JSON output for scripting
openGyver phone validate "2025551234" --country US --json

# Validate with default country (US)
openGyver phone validate "5551234567"

# Validate a short number (will fail)
openGyver phone validate "12345" --country US

# Validate a Singapore number
openGyver phone validate "91234567" --country SG

# Validate a Korean number
openGyver phone validate "010-1234-5678" --country KR
```

#### JSON Output Format (valid)

```json
{
  "input": "2025551234",
  "country": "United States",
  "valid": true,
  "digits": 10,
  "dialCode": "1",
  "nationalDigits": "2025551234"
}
```

#### JSON Output Format (invalid)

```json
{
  "input": "12345",
  "country": "United States",
  "valid": false,
  "digits": 5,
  "dialCode": "1",
  "reason": "too few digits: got 5, need at least 10",
  "nationalDigits": "12345"
}
```

---

### country

Look up a country by name or 2-letter ISO code to see its dial code, typical digit count, and formatting patterns (national and international). Supports substring matching on country names.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `name or ISO code` | Yes | Country name (or substring) or 2-letter ISO code |

#### Flags

No subcommand-specific flags. Uses global flags only.

#### Examples

```bash
# Look up by ISO code
openGyver phone country US

# Look up by full name
openGyver phone country "Japan"

# Look up by ISO code (South Korea)
openGyver phone country KR

# Look up by name substring
openGyver phone country "united kingdom"

# JSON output
openGyver phone country US --json

# Look up a European country
openGyver phone country DE

# Look up using partial name
openGyver phone country "Sing"

# Look up with JSON for scripting
openGyver phone country "Brazil" --json
```

#### JSON Output Format

```json
{
  "name": "United States",
  "iso": "US",
  "dialCode": "1",
  "minDigits": 10,
  "maxDigits": 10,
  "natPattern": "(XXX) XXX-XXXX",
  "intPattern": "+1 XXX-XXX-XXXX"
}
```

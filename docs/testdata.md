# testdata

Generate random test data for development and testing.

## Usage

```bash
openGyver testdata [command] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for testdata |

## Subcommands

### person

Generate random person data with name, email, phone, and address. Each record includes first name, last name, email, US phone number, street address, city, state, zip code, and age.

#### Arguments

No positional arguments.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--count` | | int | `1` | Number of people to generate |
| `--help` | `-h` | | | Help for person |

#### Examples

```bash
# Generate a single person
openGyver testdata person

# Generate 5 people
openGyver testdata person --count 5

# Generate 3 people as JSON
openGyver testdata person --count 3 -j

# Generate 10 people for a test database
openGyver testdata person --count 10

# Generate a large dataset for load testing
openGyver testdata person --count 100 --json

# Pipe person data into a file
openGyver testdata person --count 50 > test-users.txt

# Generate JSON people and pipe to jq
openGyver testdata person --count 5 -j | jq '.[].email'

# Single person for quick testing
openGyver testdata person --json
```

#### Plain Output Format

```
James Smith | james.smith@gmail.com | (555) 123-4567 | 1234 Main St, New York, NY 10001 | Age: 34
```

#### JSON Output Format

```json
[
  {
    "name": "James Smith",
    "email": "james.smith@gmail.com",
    "phone": "(555) 123-4567",
    "address": "1234 Main St",
    "city": "New York",
    "state": "NY",
    "zip": "10001",
    "age": 34
  }
]
```

---

### csv

Generate CSV with configurable columns. Outputs a header row followed by the specified number of data rows.

**Available column types:** name, email, number, date, bool, city, country, age, phone, uuid

#### Arguments

No positional arguments.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--rows` | | int | `10` | Number of data rows to generate |
| `--columns` | | string | `name,email,age,city` | Comma-separated column types |
| `--help` | `-h` | | | Help for csv |

#### Examples

```bash
# Default CSV (10 rows, name/email/age/city)
openGyver testdata csv

# Custom columns with 10 rows
openGyver testdata csv --rows 10 --columns name,email,age

# Generate with city, phone, and date columns
openGyver testdata csv --rows 5 --columns name,city,phone,date

# Large CSV for testing
openGyver testdata csv --rows 1000 --columns name,email,age,city,country

# CSV with boolean and UUID columns
openGyver testdata csv --rows 20 --columns name,email,bool,uuid

# Save to a file
openGyver testdata csv --rows 50 --columns name,email,phone > contacts.csv

# Generate a dataset with all column types
openGyver testdata csv --rows 10 --columns name,email,number,date,bool,city,country,age,phone,uuid

# Minimal dataset
openGyver testdata csv --rows 3 --columns name,age
```

#### Output Format

```csv
name,email,age,city
James Smith,james42@gmail.com,34,New York
Mary Johnson,mary127@outlook.com,28,Chicago
```

---

### json

Generate random JSON objects with id, name, email, age, active, and created_at fields. Always outputs as JSON array.

#### Arguments

No positional arguments.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--count` | | int | `5` | Number of JSON records to generate |
| `--help` | `-h` | | | Help for json |

#### Examples

```bash
# Generate 5 JSON records (default)
openGyver testdata json

# Generate 10 records
openGyver testdata json --count 10

# Single record
openGyver testdata json --count 1

# Large dataset
openGyver testdata json --count 100

# Save to a file
openGyver testdata json --count 50 > test-data.json

# Pipe to jq for filtering
openGyver testdata json --count 20 | jq '[.[] | select(.active == true)]'

# Generate fixture data for tests
openGyver testdata json --count 3 > fixtures/users.json

# Quick single record for debugging
openGyver testdata json --count 1 | jq '.[0]'
```

#### JSON Output Format

```json
[
  {
    "id": 1,
    "name": "James Smith",
    "email": "james.smith@gmail.com",
    "age": 34,
    "active": true,
    "created_at": "2024-06-15T10:30:00Z"
  }
]
```

---

### number

Generate random numbers within a range. Supports both integer and floating-point generation.

#### Arguments

No positional arguments.

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--min` | | int | `0` | Minimum value |
| `--max` | | int | `100` | Maximum value |
| `--count` | | int | `1` | How many numbers to generate |
| `--float` | | bool | `false` | Generate floating-point numbers |
| `--help` | `-h` | | | Help for number |

#### Examples

```bash
# Generate a random number 1-100
openGyver testdata number --min 1 --max 100

# Generate 5 random numbers
openGyver testdata number --min 1 --max 100 --count 5

# Generate floating-point numbers
openGyver testdata number --min 0 --max 1 --float --count 3

# Generate dice rolls
openGyver testdata number --min 1 --max 6 --count 10

# JSON output
openGyver testdata number --min 1 --max 100 --count 5 --json

# Generate large random numbers
openGyver testdata number --min 1000 --max 9999 --count 10

# Generate random percentages as floats
openGyver testdata number --min 0 --max 100 --float --count 5

# Single random number for a seed
openGyver testdata number --min 0 --max 999999
```

#### Plain Output Format

```
42
87
15
```

#### JSON Output Format

```json
{
  "min": 1,
  "max": 100,
  "count": 5,
  "values": [42, 87, 15, 63, 91]
}
```

# finance

A collection of financial calculators for everyday money math.

## Usage

```bash
openGyver finance [subcommand] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output results as machine-readable JSON |
| `--help` | `-h` | bool | | Show help for the command |

## Subcommands

### loan

Calculate monthly payment, total payment, and total interest for a fixed-rate loan or mortgage.

Uses the standard amortisation formula:

```
M = P * [r(1+r)^n] / [(1+r)^n - 1]
```

where `P` = principal, `r` = monthly interest rate, `n` = total number of monthly payments.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--principal` | | float | | Loan amount in dollars (required) |
| `--rate` | | float | | Annual interest rate as a percentage (required) |
| `--years` | | int | | Loan term in years (required) |

#### Examples

```bash
# 30-year mortgage at 6.5%
openGyver finance loan --principal 250000 --rate 6.5 --years 30

# 5-year car loan at 4.9%
openGyver finance loan --principal 35000 --rate 4.9 --years 5

# 15-year mortgage at 5.75% (JSON output)
openGyver finance loan --principal 400000 --rate 5.75 --years 15 -j

# Student loan
openGyver finance loan --principal 50000 --rate 5.0 --years 10

# Small personal loan
openGyver finance loan --principal 10000 --rate 8.5 --years 3

# Compare 15-year vs 30-year mortgage
openGyver finance loan --principal 300000 --rate 6.0 --years 15
openGyver finance loan --principal 300000 --rate 6.5 --years 30

# JSON output piped to jq
openGyver finance loan --principal 200000 --rate 7 --years 30 -j | jq '.monthly_payment'
```

#### JSON Output Format

```json
{
  "principal": 250000,
  "annual_rate": 6.5,
  "years": 30,
  "monthly_payment": 1580.17,
  "total_payment": 568861.22,
  "total_interest": 318861.22
}
```

---

### compound

Calculate the future value of an investment with compound interest.

Uses the compound interest formula:

```
A = P * (1 + r/n)^(n*t)
```

where `P` = principal, `r` = annual rate, `n` = compounding frequency per year, `t` = time in years.

**Common compounding frequencies:**

| Value | Meaning |
|-------|---------|
| 1 | Annually |
| 4 | Quarterly |
| 12 | Monthly (default) |
| 52 | Weekly |
| 365 | Daily |

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--principal` | | float | | Initial investment amount (required) |
| `--rate` | | float | | Annual interest rate as a percentage (required) |
| `--years` | | int | | Investment period in years (required) |
| `--frequency` | | int | `12` | Compounding frequency per year |

#### Examples

```bash
# $10,000 at 7% for 10 years, compounded monthly
openGyver finance compound --principal 10000 --rate 7 --years 10

# $5,000 at 5% for 20 years, compounded daily
openGyver finance compound --principal 5000 --rate 5 --years 20 --frequency 365

# $25,000 at 8.5% for 5 years, compounded quarterly
openGyver finance compound --principal 25000 --rate 8.5 --years 5 --frequency 4

# Retirement savings: $100k at 7% for 30 years
openGyver finance compound --principal 100000 --rate 7 --years 30

# Annual compounding
openGyver finance compound --principal 10000 --rate 5 --years 10 --frequency 1

# JSON output
openGyver finance compound --principal 50000 --rate 6 --years 15 -j

# Weekly compounding
openGyver finance compound --principal 10000 --rate 4.5 --years 5 --frequency 52

# Compare daily vs monthly compounding
openGyver finance compound --principal 10000 --rate 5 --years 10 --frequency 12
openGyver finance compound --principal 10000 --rate 5 --years 10 --frequency 365
```

#### JSON Output Format

```json
{
  "principal": 10000,
  "annual_rate": 7,
  "years": 10,
  "frequency": 12,
  "final_amount": 20096.61,
  "total_interest": 10096.61
}
```

---

### roi

Calculate the return on investment (ROI) given an initial investment and a final value.

```
ROI = ((final - initial) / initial) * 100
```

Positive ROI indicates a profit; negative ROI indicates a loss.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--initial` | | float | | Initial investment amount (required) |
| `--final` | | float | | Final value of the investment (required) |

#### Examples

```bash
# Bought at $1,000, now worth $1,500 (50% ROI)
openGyver finance roi --initial 1000 --final 1500

# Bought at $5,000, sold for $3,750 (a loss)
openGyver finance roi --initial 5000 --final 3750

# JSON output
openGyver finance roi --initial 200 --final 350 -j

# Real estate investment
openGyver finance roi --initial 250000 --final 375000

# Stock investment loss
openGyver finance roi --initial 10000 --final 8500

# Doubled your money
openGyver finance roi --initial 5000 --final 10000

# Pipe JSON to jq
openGyver finance roi --initial 1000 --final 1250 -j | jq '.roi_percent'

# Break even
openGyver finance roi --initial 1000 --final 1000
```

#### JSON Output Format

```json
{
  "initial": 1000,
  "final": 1500,
  "profit_loss": 500,
  "roi_percent": 50
}
```

---

### tip

Calculate tip, total bill, and per-person share when splitting the bill among multiple people.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--amount` | | float | | Bill amount before tip (required) |
| `--percent` | | float | `18` | Tip percentage |
| `--split` | | int | `1` | Number of people splitting the bill |

#### Examples

```bash
# Standard 18% tip on $85.50
openGyver finance tip --amount 85.50

# 20% tip on $120, split 4 ways
openGyver finance tip --amount 120 --percent 20 --split 4

# 15% tip on $45.00
openGyver finance tip --amount 45 --percent 15

# Generous 25% tip
openGyver finance tip --amount 200 --percent 25

# Split between 6 people
openGyver finance tip --amount 350 --percent 20 --split 6

# JSON output
openGyver finance tip --amount 85.50 --percent 20 -j

# Minimal tip
openGyver finance tip --amount 30 --percent 10

# Large party dinner split
openGyver finance tip --amount 800 --percent 18 --split 10
```

#### JSON Output Format

```json
{
  "amount": 85.5,
  "percent": 20,
  "tip": 17.1,
  "total": 102.6,
  "split": 4,
  "per_person": 25.65
}
```

---

### tax

Calculate the tax amount and total price given a pre-tax amount and a tax rate.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--amount` | | float | | Pre-tax amount (required) |
| `--rate` | | float | | Tax rate as a percentage (required) |

#### Examples

```bash
# 8.25% sales tax on $99.99
openGyver finance tax --amount 99.99 --rate 8.25

# 10% tax on $250
openGyver finance tax --amount 250 --rate 10

# JSON output
openGyver finance tax --amount 49.95 --rate 7.5 -j

# High tax rate
openGyver finance tax --amount 500 --rate 20

# Vehicle purchase tax
openGyver finance tax --amount 35000 --rate 6.5

# Low tax rate
openGyver finance tax --amount 15.99 --rate 4

# Pipe JSON to jq
openGyver finance tax --amount 100 --rate 8.875 -j | jq '.total'

# Multiple items total
openGyver finance tax --amount 1234.56 --rate 9.5
```

#### JSON Output Format

```json
{
  "amount": 99.99,
  "rate": 8.25,
  "tax": 8.25,
  "total": 108.24
}
```

---

### salary

Convert salary or wages between different pay periods.

Assumes a standard full-time work schedule:
- 8 hours per day
- 5 days per week
- 4.33 weeks per month (52/12)
- 52 weeks per year

**Valid period values:** `hourly`, `daily`, `weekly`, `monthly`, `yearly`

The text output shows all pay period conversions at once, regardless of the `--to` flag.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--amount` | | float | | Salary/wage amount (required) |
| `--from` | | string | | Source pay period: hourly, daily, weekly, monthly, yearly (required) |
| `--to` | | string | | Target pay period: hourly, daily, weekly, monthly, yearly (required) |

#### Examples

```bash
# Hourly to yearly
openGyver finance salary --amount 50 --from hourly --to yearly

# Yearly to hourly
openGyver finance salary --amount 100000 --from yearly --to hourly

# Monthly to weekly
openGyver finance salary --amount 8000 --from monthly --to weekly

# Daily to monthly
openGyver finance salary --amount 400 --from daily --to monthly

# JSON output
openGyver finance salary --amount 75 --from hourly --to yearly -j

# Weekly to yearly
openGyver finance salary --amount 2000 --from weekly --to yearly

# Compare hourly rates
openGyver finance salary --amount 120000 --from yearly --to hourly

# Monthly to daily
openGyver finance salary --amount 6500 --from monthly --to daily
```

#### JSON Output Format

```json
{
  "amount": 50,
  "from": "hourly",
  "to": "yearly",
  "result": 104000
}
```

---

### discount

Calculate the discount amount and final price after applying a percentage discount.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--price` | | float | | Original price (required) |
| `--percent` | | float | | Discount percentage (required, 0-100) |

#### Examples

```bash
# 25% off $199.99
openGyver finance discount --price 199.99 --percent 25

# 50% off $80
openGyver finance discount --price 80 --percent 50

# 10% off $1,250 (JSON output)
openGyver finance discount --price 1250 --percent 10 -j

# Black Friday sale: 40% off
openGyver finance discount --price 599.99 --percent 40

# Small discount
openGyver finance discount --price 29.99 --percent 5

# Full price (0% discount)
openGyver finance discount --price 100 --percent 0

# Pipe JSON to jq for final price only
openGyver finance discount --price 200 --percent 30 -j | jq '.final_price'

# Compare discounts
openGyver finance discount --price 500 --percent 15
openGyver finance discount --price 500 --percent 25
```

#### JSON Output Format

```json
{
  "original_price": 199.99,
  "discount_percent": 25,
  "discount_amount": 50,
  "final_price": 149.99
}
```

---

### margin

Calculate profit, profit margin percentage, and markup percentage given cost and revenue (selling price).

```
Profit   = Revenue - Cost
Margin % = (Profit / Revenue) * 100
Markup % = (Profit / Cost) * 100
```

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--cost` | | float | | Cost of goods/services (required) |
| `--revenue` | | float | | Selling price / revenue (required) |

#### Examples

```bash
# Cost $40, selling for $100
openGyver finance margin --cost 40 --revenue 100

# Cost $15.50, selling for $29.99
openGyver finance margin --cost 15.50 --revenue 29.99

# JSON output
openGyver finance margin --cost 250 --revenue 400 -j

# High margin product
openGyver finance margin --cost 5 --revenue 49.99

# Low margin product
openGyver finance margin --cost 80 --revenue 90

# Selling at a loss
openGyver finance margin --cost 100 --revenue 75

# Pipe to jq for margin percent only
openGyver finance margin --cost 30 --revenue 60 -j | jq '.margin_percent'

# Wholesale vs retail
openGyver finance margin --cost 12 --revenue 24.99
```

#### JSON Output Format

```json
{
  "cost": 40,
  "revenue": 100,
  "profit": 60,
  "margin_percent": 60,
  "markup_percent": 150
}
```

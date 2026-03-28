# stock

Look up current or historical stock prices from global markets.

## Usage

```bash
openGyver stock <ticker> [flags]
```

Uses Yahoo Finance data (no API key required). Tickers are searched universally -- just type the symbol and it auto-resolves. Use `--market` to target a specific exchange.

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--date` | `-d` | string | | Look up price on a specific date (YYYY-MM-DD) |
| `--from` | | string | | Start date for historical range (YYYY-MM-DD) |
| `--to` | | string | | End date for historical range (YYYY-MM-DD) |
| `--interval` | | string | `1d` | Data interval: 1d, 1wk, 1mo |
| `--market` | `-m` | string | | Target market/exchange (e.g. kosdaq, tokyo, london) |
| `--field` | `-f` | string | | Output a single field value (ideal for piping) |
| `--json` | `-j` | bool | `false` | Output as JSON |
| `--help` | `-h` | | | Help for stock |

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `ticker` | Yes | Stock ticker symbol (e.g. AAPL, MSFT, 005930) |

## Ticker Format

| Format | Example | Description |
|--------|---------|-------------|
| US stocks | `AAPL`, `MSFT`, `GOOGL`, `TSLA` | No suffix needed |
| With suffix | `005930.KS`, `7203.T` | Explicit exchange suffix |
| With `--market` | `005930 --market kospi` | Auto-appends exchange suffix |

## Supported Markets (for `--market` flag)

| Region | Market names |
|--------|-------------|
| US | nasdaq, nyse, us |
| Korea | kosdaq, kospi, korea |
| Japan | tokyo, tse, japan |
| UK | london, lse, uk |
| China | shanghai, shenzhen, hongkong, hk |
| Europe | frankfurt, xetra, paris, euronext, amsterdam, swiss |
| Americas | toronto, tsx, brazil, bovespa, mexico |
| Asia | singapore, sgx, taiwan, mumbai, nse, jakarta, bangkok |
| Other | australia, asx, johannesburg, oslo, stockholm, copenhagen, helsinki, newzealand |

## Abbreviated Output Fields (`--field` / `-f`)

| Field | Description |
|-------|-------------|
| `price` | Current / closing price |
| `change` | Price change (absolute) |
| `percent` | Price change (percentage) |
| `open` | Opening price |
| `high` | High price |
| `low` | Low price |
| `close` | Closing price |
| `volume` | Trading volume |
| `currency` | Currency code |
| `exchange` | Exchange name |

## Examples

```bash
# Current price of Apple
openGyver stock AAPL

# Historical price on a specific date
openGyver stock MSFT --date 2024-01-15

# Date range with daily data
openGyver stock AAPL --from 2024-01-01 --to 2024-06-30

# Weekly historical data
openGyver stock AAPL --from 2024-01-01 --to 2024-06-30 --interval 1wk

# Monthly historical data
openGyver stock AAPL --from 2023-01-01 --to 2024-01-01 --interval 1mo

# Korean stock by market name
openGyver stock 005930 --market kospi

# Japanese stock
openGyver stock 7203 --market tokyo

# UK stock
openGyver stock SHEL --market london

# Hong Kong stock
openGyver stock 0700 --market hk

# Extract just the price (for piping)
openGyver stock AAPL -f price

# Extract just the percent change
openGyver stock AAPL -f percent

# Extract absolute change
openGyver stock AAPL -f change

# JSON output for automation
openGyver stock AAPL --json

# Get volume as a single value
openGyver stock TSLA -f volume

# Full JSON for a specific date
openGyver stock GOOGL --date 2024-06-15 --json
```

## JSON Output Format (current price)

```json
{
  "symbol": "AAPL",
  "exchange": "NMS",
  "currency": "USD",
  "price": 248.80,
  "change": 2.30,
  "percent": 0.93,
  "previous_close": 246.50,
  "as_of": "2024-12-20T16:00:00-05:00"
}
```

## JSON Output Format (specific date)

```json
{
  "symbol": "AAPL",
  "exchange": "NMS",
  "currency": "USD",
  "date": "2024-01-15",
  "open": 182.16,
  "high": 185.10,
  "low": 180.17,
  "close": 185.10,
  "price": 185.10,
  "volume": 65076600,
  "change": 0.0,
  "percent": 0.0
}
```

## JSON Output Format (date range)

```json
{
  "symbol": "AAPL",
  "exchange": "NMS",
  "currency": "USD",
  "interval": "1d",
  "data": [
    {
      "date": "2024-01-02",
      "open": 187.15,
      "high": 188.44,
      "low": 183.89,
      "close": 185.64,
      "volume": 82488700
    }
  ]
}
```

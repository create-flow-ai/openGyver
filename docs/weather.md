# weather

Look up current weather, forecasts (up to 16 days), and historical data (back to 1940) for any city worldwide. Uses Open-Meteo API — free, no API key required.

## Usage

```bash
openGyver weather <city> [flags]
openGyver weather --lat <latitude> --lon <longitude> [flags]
```

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `<city>` | Yes (or use `--lat`/`--lon`) | City name. Geocoded automatically. Examples: `"New York"`, `Tokyo`, `"San Francisco"` |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--date` | `-d` | string | | Weather on a specific date (YYYY-MM-DD). Past dates use historical API, future up to 16 days uses forecast. |
| `--from` | | string | | Start date for date range (YYYY-MM-DD) |
| `--to` | | string | | End date for date range (YYYY-MM-DD) |
| `--lat` | | float | 0 | Latitude (use with `--lon` instead of city name) |
| `--lon` | | float | 0 | Longitude (use with `--lat` instead of city name) |
| `--units` | | string | `celsius` | Temperature units: `celsius` or `fahrenheit` |
| `--field` | `-f` | string | | Output a single value for piping. See [available fields](#available-fields). |
| `--json` | `-j` | bool | false | Output as JSON |

## Available Fields

For `--field` / `-f`:

| Field | Description |
|-------|-------------|
| `temperature` | Current temperature |
| `feels_like` | Apparent temperature (wind chill / heat index) |
| `humidity` | Relative humidity (%) |
| `wind_speed` | Wind speed |
| `wind_direction` | Wind direction (degrees) |
| `precipitation` | Precipitation (mm) |
| `cloud_cover` | Cloud cover (%) |
| `pressure` | Surface pressure (hPa) |
| `uv_index` | UV index |
| `description` | Weather condition text (e.g., "Clear sky", "Moderate rain") |
| `temp_max` | Daily maximum temperature |
| `temp_min` | Daily minimum temperature |
| `sunrise` | Sunrise time |
| `sunset` | Sunset time |

## Data Returned

Current weather includes: temperature, feels-like, daily high/low, humidity, wind speed and direction, precipitation, cloud cover, surface pressure, UV index, sunrise/sunset, and weather condition description.

Date range queries return a table with: date, condition, high, low, precipitation, max wind speed, and UV index.

## Examples

### Current Weather

```bash
# Current weather for a city
openGyver weather "New York"

# Fahrenheit
openGyver weather Tokyo --units fahrenheit

# Using coordinates
openGyver weather --lat 40.7128 --lon -74.006

# JSON output
openGyver weather Paris -j

# Single field for scripting
openGyver weather London -f temperature
openGyver weather Berlin -f humidity
openGyver weather "San Francisco" -f description
```

### Specific Date

```bash
# Weather on a past date (back to 1940)
openGyver weather London --date 2024-12-25

# Weather forecast for a future date (up to 16 days ahead)
openGyver weather Seoul --date 2025-04-05

# Historical date with JSON
openGyver weather "New York" --date 2020-03-15 -j

# Historical date with single field
openGyver weather Tokyo --date 2023-07-04 -f temp_max
```

### Date Range

```bash
# Recent week
openGyver weather Paris --from 2025-03-20 --to 2025-03-27

# Historical range
openGyver weather Berlin --from 2024-01-01 --to 2024-01-31

# Forecast range
openGyver weather Tokyo --from 2025-03-28 --to 2025-04-10

# Range with Fahrenheit
openGyver weather Chicago --from 2025-03-01 --to 2025-03-07 --units fahrenheit

# Range as JSON (returns data[] array)
openGyver weather London --from 2025-03-01 --to 2025-03-07 -j

# Extract just temperatures from a range
openGyver weather Seoul --from 2025-03-24 --to 2025-03-27 -f temp_max
# 18.5
# 22.5
# 24.0
# 21.6
```

### Piping and Scripting

```bash
# Get temperature as a number
TEMP=$(openGyver weather "New York" -f temperature)
echo "Current temp in NYC: $TEMP"

# Check if it's raining
openGyver weather London -f description | grep -i rain

# Get weekly forecast as JSON and process with jq
openGyver weather Tokyo --from 2025-03-28 --to 2025-04-03 -j | jq '.data[].temp_max'

# Compare two cities
echo "NYC: $(openGyver weather 'New York' -f temperature)C"
echo "LA:  $(openGyver weather 'Los Angeles' -f temperature)C"
```

## JSON Output Format

### Current Weather

```json
{
  "city": "New York, United States",
  "temperature": 1.0,
  "feels_like": -5.5,
  "temp_max": 5.6,
  "temp_min": -0.6,
  "humidity": 56,
  "wind_speed": 21.3,
  "wind_direction": 330,
  "precipitation": 0.0,
  "cloud_cover": 0,
  "pressure": 1020,
  "uv_index": 0.0,
  "description": "Clear sky",
  "sunrise": "2026-03-28T06:45",
  "sunset": "2026-03-28T19:16",
  "time": "2026-03-28T12:15"
}
```

### Date Range

```json
{
  "city": "Tokyo, Japan",
  "data": [
    {
      "date": "2025-03-24",
      "temp_max": 18.5,
      "temp_min": 8.9,
      "precipitation_mm": 0.2,
      "wind_max": 11.4,
      "wind_direction": 180,
      "uv_max": 4.2,
      "description": "Light drizzle",
      "sunrise": "2025-03-24T05:37",
      "sunset": "2025-03-24T17:54"
    }
  ]
}
```

## Notes

- **Data source**: Open-Meteo (https://open-meteo.com) — free for non-commercial use, no API key.
- **Historical data**: Available back to 1940 for most locations.
- **Forecast data**: Up to 16 days ahead.
- **Geocoding**: City names are automatically resolved to coordinates. For ambiguous names, be specific (e.g., "Portland, Oregon" vs "Portland, Maine").
- **Weather codes**: Uses WMO standard weather interpretation codes for condition descriptions.
- **Timezone**: All times are in the location's local timezone.

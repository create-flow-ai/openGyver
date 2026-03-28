# geo

Geolocation tools for distance calculation, coordinate conversion, and more.

## Usage

```bash
openGyver geo [subcommand] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output results as machine-readable JSON |
| `--help` | `-h` | bool | | Show help for the command |

## Subcommands

### distance

Calculate the great-circle distance between two points on Earth using the Haversine formula. The Haversine formula accounts for Earth's curvature and gives an accurate distance for any two points on the globe.

**Formula:**

```
a = sin^2(delta_phi/2) + cos(phi1) * cos(phi2) * sin^2(delta_lambda/2)
c = 2 * atan2(sqrt(a), sqrt(1-a))
d = R * c
```

where `R` is Earth's mean radius (6,371 km).

Output includes the distance in both kilometres and miles.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--lat1` | | float | | Latitude of point 1 in decimal degrees (required, -90 to 90) |
| `--lon1` | | float | | Longitude of point 1 in decimal degrees (required, -180 to 180) |
| `--lat2` | | float | | Latitude of point 2 in decimal degrees (required, -90 to 90) |
| `--lon2` | | float | | Longitude of point 2 in decimal degrees (required, -180 to 180) |

#### Examples

```bash
# New York to London
openGyver geo distance --lat1 40.7128 --lon1 -74.0060 --lat2 51.5074 --lon2 -0.1278

# Tokyo to Sydney
openGyver geo distance --lat1 35.6762 --lon1 139.6503 --lat2 -33.8688 --lon2 151.2093

# San Francisco to Paris
openGyver geo distance --lat1 37.7749 --lon1 -122.4194 --lat2 48.8566 --lon2 2.3522

# JSON output
openGyver geo distance --lat1 40.7128 --lon1 -74.0060 --lat2 51.5074 --lon2 -0.1278 -j

# Los Angeles to Tokyo
openGyver geo distance --lat1 34.0522 --lon1 -118.2437 --lat2 35.6762 --lon2 139.6503

# North Pole to South Pole
openGyver geo distance --lat1 90 --lon1 0 --lat2 -90 --lon2 0

# Same city (should be ~0)
openGyver geo distance --lat1 40.7128 --lon1 -74.0060 --lat2 40.7580 --lon2 -73.9855

# Pipe JSON to jq for kilometres only
openGyver geo distance --lat1 40.7128 --lon1 -74.0060 --lat2 51.5074 --lon2 -0.1278 -j | jq '.distance_km'
```

#### JSON Output Format

```json
{
  "point1": {
    "latitude": 40.7128,
    "longitude": -74.006
  },
  "point2": {
    "latitude": 51.5074,
    "longitude": -0.1278
  },
  "distance_km": 5570.25,
  "distance_miles": 3461.02
}
```

---

### dms

Convert a coordinate between decimal degrees and DMS (degrees, minutes, seconds) notation. The input format is auto-detected:

- If the input is a plain number (e.g., `40.7128`), it is treated as decimal degrees and converted to DMS.
- If the input contains degree symbols or DMS notation (e.g., `40d42m46.1sN`), it is parsed as DMS and converted to decimal degrees.

For decimal-degree input, the direction is inferred heuristically: if the absolute value is <= 90, it is assumed to be latitude (N/S); otherwise it is assumed to be longitude (E/W).

**Supported DMS input formats:**

| Format | Example | Description |
|--------|---------|-------------|
| Full DMS with direction | `40d42'46.1"N` | Degree symbol, minute, second, direction |
| DMS without direction | `40d42'46.1"` | Without cardinal direction |
| Space-separated | `40 42 46.1 N` | Spaces between components |
| Negative sign | `-40 42 46.1` | Negative for S/W |
| Letter notation | `40d42m46.1s` | d/m/s letters |

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `coordinate` | Yes | Decimal degrees value or DMS string to convert |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| (none beyond global) | | | | |

#### Examples

```bash
# Decimal degrees to DMS
openGyver geo dms 40.7128

# Negative decimal degrees to DMS (use -- to prevent flag parsing)
openGyver geo dms -- -74.0060

# DMS to decimal degrees (with degree/minute/second symbols)
openGyver geo dms "40°42'46.1\"N"

# Space-separated DMS to decimal degrees
openGyver geo dms "40 42 46.1 N"

# Longitude value (auto-detected as E/W when > 90)
openGyver geo dms 151.2093

# JSON output for decimal-to-DMS
openGyver geo dms 40.7128 -j

# JSON output for DMS-to-decimal
openGyver geo dms "40 42 46.1 N" -j

# Letter notation DMS
openGyver geo dms "40d42m46.1s"
```

#### JSON Output Format

When converting decimal degrees to DMS:

```json
{
  "input": "40.7128",
  "decimal_degrees": 40.7128,
  "dms": {
    "degrees": 40,
    "minutes": 42,
    "seconds": 46.08,
    "direction": "N"
  },
  "formatted": "40° 42' 46.08\" N"
}
```

When converting DMS to decimal degrees:

```json
{
  "input": "40 42 46.1 N",
  "decimal_degrees": 40.712806,
  "dms": {
    "degrees": 40,
    "minutes": 42,
    "seconds": 46.1,
    "direction": "N"
  }
}
```

---

### utm

Convert geographic coordinates (latitude/longitude in decimal degrees) to Universal Transverse Mercator (UTM) coordinates.

UTM divides the Earth into 60 zones (1-60), each 6 degrees wide. Within each zone, locations are specified by easting (metres from the zone's central meridian, offset by 500,000 m) and northing (metres from the equator, with a 10,000,000 m false northing for the southern hemisphere).

The conversion uses the WGS 84 ellipsoid parameters. Special zone handling is applied for Norway and Svalbard.

**Latitude band letters** are assigned from C (at -80 degrees) through X (at +84 degrees), covering the valid UTM range.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--lat` | | float | | Latitude in decimal degrees (required, -80 to 84) |
| `--lon` | | float | | Longitude in decimal degrees (required, -180 to 180) |

#### Examples

```bash
# Statue of Liberty, New York
openGyver geo utm --lat 40.6892 --lon -74.0445

# Sydney Opera House
openGyver geo utm --lat -33.8568 --lon 151.2153

# Null Island (0, 0)
openGyver geo utm --lat 0 --lon 0

# JSON output
openGyver geo utm --lat 40.7128 --lon -74.0060 -j

# Eiffel Tower, Paris
openGyver geo utm --lat 48.8584 --lon 2.2945

# Tokyo Tower
openGyver geo utm --lat 35.6586 --lon 139.7454

# Southern hemisphere: Cape Town
openGyver geo utm --lat -33.9249 --lon 18.4241

# Pipe JSON to jq for zone only
openGyver geo utm --lat 51.5074 --lon -0.1278 -j | jq '.zone_designator'
```

#### JSON Output Format

```json
{
  "input": {
    "latitude": 40.6892,
    "longitude": -74.0445
  },
  "zone": 18,
  "letter": "T",
  "easting": 583960.00,
  "northing": 4507523.00,
  "zone_designator": "18T"
}
```

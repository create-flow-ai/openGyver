# electrical

Electrical engineering tools for circuit design and component calculation.

## Usage

```bash
openGyver electrical [subcommand] [flags]
```

## Global Flags (inherited by all subcommands)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--json` | `-j` | bool | `false` | Output results as machine-readable JSON |
| `--help` | `-h` | bool | | Show help for the command |

## Subcommands

### ohm

Calculate voltage, current, resistance, and power using Ohm's law. Provide any two of the three fundamental values (voltage, current, resistance) and the calculator determines the third, plus the dissipated power.

Formulas used:

- `V = I x R` (Voltage = Current x Resistance)
- `I = V / R` (Current = Voltage / Resistance)
- `R = V / I` (Resistance = Voltage / Current)
- `P = V x I` (Power = Voltage x Current)

Values are displayed with appropriate SI prefixes (mA, kOhm, mW, etc.) for readability.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--voltage` | `-v` | float | `0` | Voltage in Volts (V) |
| `--current` | `-i` | float | `0` | Current in Amps (A) |
| `--resistance` | `-r` | float | `0` | Resistance in Ohms |

Exactly 2 of the 3 flags must be provided. The third value is calculated.

#### Examples

```bash
# Given voltage and current, find resistance and power
openGyver electrical ohm --voltage 12 --current 0.5

# Given voltage and resistance, find current and power
openGyver electrical ohm -v 5 -r 1000

# Given current and resistance, find voltage and power
openGyver electrical ohm -i 0.02 -r 220

# JSON output for scripting
openGyver electrical ohm --voltage 12 --current 0.5 -j

# High-voltage circuit
openGyver electrical ohm --voltage 240 --resistance 1000

# Low-power microcontroller pin
openGyver electrical ohm -v 3.3 -r 10000

# Calculate resistance for a known power scenario
openGyver electrical ohm -v 9 -i 0.001

# Short form flags for quick calculations
openGyver electrical ohm -i 0.1 -r 47
```

#### JSON Output Format

```json
{
  "voltage_V": 12,
  "current_A": 0.5,
  "resistance_Ohm": 24,
  "power_W": 6,
  "calculated": "resistance"
}
```

---

### resistor

Convert a resistance value to its corresponding color band codes for standard 4-band and 5-band resistors.

**Input formats:**

| Format | Example | Description |
|--------|---------|-------------|
| Plain number | `4700` | Value in Ohms |
| SI suffix | `4.7k` | `k` = kilo, `M` = mega, `R` = Ohms |
| Mixed notation | `4k7` | Common in schematics (4.7k) |
| Megaohm | `2.2M` | 2.2 megaohms |
| Explicit Ohms | `470R` | 470 Ohms |

**Color band mapping:**

| Digit | Color | Digit | Color |
|-------|-------|-------|-------|
| 0 | Black | 5 | Green |
| 1 | Brown | 6 | Blue |
| 2 | Red | 7 | Violet |
| 3 | Orange | 8 | Grey |
| 4 | Yellow | 9 | White |

**Multiplier colors:** x1=Black, x10=Brown, x100=Red, x1k=Orange, x10k=Yellow, x100k=Green, x1M=Blue, x10M=Violet, x0.1=Gold, x0.01=Silver

**Band types:**
- 4-band: `[digit1] [digit2] [multiplier] [tolerance +/-5% Gold]`
- 5-band: `[digit1] [digit2] [digit3] [multiplier] [tolerance +/-1% Brown]`

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `value` | Yes | Resistance value (plain number, SI suffix, or mixed notation) |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| (none beyond global) | | | | |

#### Examples

```bash
# Standard 4.7k resistor
openGyver electrical resistor 4700

# Using SI suffix notation
openGyver electrical resistor 4.7k

# Mixed schematic notation
openGyver electrical resistor 4k7

# 10k resistor
openGyver electrical resistor 10k

# Megaohm range
openGyver electrical resistor 2.2M

# Explicit Ohms notation
openGyver electrical resistor 470R

# JSON output
openGyver electrical resistor 4700 -j

# Small resistance
openGyver electrical resistor 10
```

#### JSON Output Format

```json
{
  "input": "4700",
  "value_ohms": 4700,
  "formatted": "4.7 kOhm",
  "4_band": {
    "bands": ["Yellow", "Violet", "Red", "Gold"],
    "tolerance": "+/-5%"
  },
  "5_band": {
    "bands": ["Yellow", "Violet", "Black", "Brown", "Brown"],
    "tolerance": "+/-1%"
  }
}
```

---

### led

Calculate the required current-limiting resistor for an LED circuit. Uses the formula `R = (Vsource - Vforward) / I`. Also finds the nearest standard E24 series resistor value (rounded up for safety) and computes the actual current and power with that standard resistor.

**Typical LED forward voltages:**

| Color | Forward Voltage |
|-------|----------------|
| Red | 1.8 - 2.2 V |
| Orange | 2.0 - 2.2 V |
| Yellow | 2.0 - 2.2 V |
| Green | 2.0 - 3.5 V |
| Blue | 3.0 - 3.5 V |
| White | 3.0 - 3.5 V |
| Infrared | 1.1 - 1.5 V |
| UV | 3.3 - 4.0 V |

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--source` | | float | | Supply voltage in Volts (V) (required) |
| `--forward` | | float | | LED forward voltage in Volts (V) (required) |
| `--current` | | float | `20` | Desired LED current in milliamps (mA) |

#### Examples

```bash
# Red LED on 5V supply (typical 20mA)
openGyver electrical led --source 5 --forward 2.0

# Blue LED on 3.3V supply at 10mA
openGyver electrical led --source 3.3 --forward 3.2 --current 10

# White LED on 12V supply
openGyver electrical led --source 12 --forward 3.3 --current 20

# Multiple LEDs in series (sum forward voltages: 3 red LEDs = 6.0V)
openGyver electrical led --source 12 --forward 6.0 --current 20

# JSON output for integration
openGyver electrical led --source 5 --forward 2.0 -j

# Infrared LED on 3.3V
openGyver electrical led --source 3.3 --forward 1.2 --current 50

# UV LED on 5V at low current
openGyver electrical led --source 5 --forward 3.5 --current 10

# High-current LED on 24V supply
openGyver electrical led --source 24 --forward 3.3 --current 350
```

#### JSON Output Format

```json
{
  "source_V": 5,
  "forward_V": 2,
  "desired_current_mA": 20,
  "calculated_resistance_Ohm": 150,
  "calculated_power_mW": 60,
  "nearest_standard_Ohm": 150,
  "actual_current_mA": 20,
  "actual_power_mW": 60
}
```

---

### divider

Calculate the output voltage of a resistive voltage divider. Given an input voltage and two resistor values, computes the output voltage, divider ratio, current through the divider, and power dissipated by each resistor.

Circuit diagram:

```
     Vin
      |
     [R1]
      |---- Vout
     [R2]
      |
     GND
```

Formula: `Vout = Vin x R2 / (R1 + R2)`

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| (none) | | All input is provided via flags |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--vin` | | float | | Input voltage in Volts (V) (required) |
| `--r1` | | float | | Upper resistor in Ohms (required) |
| `--r2` | | float | | Lower resistor in Ohms (required) |

#### Examples

```bash
# 12V down to ~3.84V with 10k and 4.7k
openGyver electrical divider --vin 12 --r1 10000 --r2 4700

# 5V to 3.3V level shifting (close approximation)
openGyver electrical divider --vin 5 --r1 1800 --r2 3300

# 24V to 5V for ADC input
openGyver electrical divider --vin 24 --r1 38000 --r2 10000

# Check divider current draw (high impedance)
openGyver electrical divider --vin 12 --r1 100000 --r2 100000

# JSON output for scripting
openGyver electrical divider --vin 12 --r1 10000 --r2 4700 -j

# 3.3V to 1.8V for a voltage reference
openGyver electrical divider --vin 3.3 --r1 8200 --r2 10000

# Battery voltage monitoring (12V to 3.3V)
openGyver electrical divider --vin 12.6 --r1 27000 --r2 10000

# Equal resistors for half voltage
openGyver electrical divider --vin 9 --r1 10000 --r2 10000
```

#### JSON Output Format

```json
{
  "vin_V": 12,
  "r1_Ohm": 10000,
  "r2_Ohm": 4700,
  "vout_V": 3.8367,
  "ratio": 0.319728,
  "current_A": 0.000816,
  "power_r1_W": 0.006662,
  "power_r2_W": 0.003131,
  "power_total_W": 0.009796
}
```

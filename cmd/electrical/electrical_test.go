package electrical

import (
	"math"
	"testing"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestFormatSI(t *testing.T) {
	tests := []struct {
		value float64
		unit  string
		want  string
	}{
		{1200, "Ω", "1.20 kΩ"},
		{0.005, "A", "5.00 mA"},
		{1.44, "W", "1.44 W"},
		{4700000, "Ω", "4.70 MΩ"},
	}
	for _, tt := range tests {
		got := formatSI(tt.value, tt.unit)
		if got != tt.want {
			t.Errorf("formatSI(%g, %q) = %q, want %q", tt.value, tt.unit, got, tt.want)
		}
	}
}

func TestOhmsLaw_VoltageAndResistance(t *testing.T) {
	// Given V=12, R=100 -> I=0.12A, P=1.44W
	voltage := 12.0
	resistance := 100.0
	current := voltage / resistance
	power := voltage * current

	if !almostEqual(current, 0.12, 1e-9) {
		t.Errorf("current = %f, want 0.12", current)
	}
	if !almostEqual(power, 1.44, 1e-9) {
		t.Errorf("power = %f, want 1.44", power)
	}
}

func TestOhmsLaw_VoltageAndCurrent(t *testing.T) {
	// Given V=5, I=0.5 -> R=10, P=2.5
	voltage := 5.0
	current := 0.5
	resistance := voltage / current
	power := voltage * current

	if !almostEqual(resistance, 10.0, 1e-9) {
		t.Errorf("resistance = %f, want 10.0", resistance)
	}
	if !almostEqual(power, 2.5, 1e-9) {
		t.Errorf("power = %f, want 2.5", power)
	}
}

func TestOhmsLaw_CurrentAndResistance(t *testing.T) {
	// Given I=0.02, R=220 -> V=4.4, P=0.088
	current := 0.02
	resistance := 220.0
	voltage := current * resistance
	power := voltage * current

	if !almostEqual(voltage, 4.4, 1e-9) {
		t.Errorf("voltage = %f, want 4.4", voltage)
	}
	if !almostEqual(power, 0.088, 1e-9) {
		t.Errorf("power = %f, want 0.088", power)
	}
}

func TestLEDResistorCalc(t *testing.T) {
	// Red LED on 5V supply, forward 2.0V, desired current 20mA.
	source := 5.0
	forward := 2.0
	currentMA := 20.0
	currentA := currentMA / 1000.0

	resistance := (source - forward) / currentA
	power := (source - forward) * currentA

	if !almostEqual(resistance, 150.0, 0.01) {
		t.Errorf("LED resistance = %f, want 150.0", resistance)
	}
	if !almostEqual(power, 0.06, 0.001) {
		t.Errorf("LED power = %f, want 0.06", power)
	}
}

func TestNearestStandardResistor(t *testing.T) {
	tests := []struct {
		ohms float64
		want float64
	}{
		{150, 150},
		{148, 150},
		{100, 100},
		{4700, 4700},
		{435, 470},
	}
	for _, tt := range tests {
		got := nearestStandardResistor(tt.ohms)
		if !almostEqual(got, tt.want, 0.1) {
			t.Errorf("nearestStandardResistor(%g) = %g, want %g", tt.ohms, got, tt.want)
		}
	}
}

func TestVoltageDivider(t *testing.T) {
	// Vin=12, R1=10k, R2=4.7k -> Vout = 12 * 4700 / (10000 + 4700) ~= 3.836V
	vin := 12.0
	r1 := 10000.0
	r2 := 4700.0

	vout := vin * r2 / (r1 + r2)
	ratio := r2 / (r1 + r2)
	current := vin / (r1 + r2)

	expectedVout := 3.836734693877551
	if !almostEqual(vout, expectedVout, 0.001) {
		t.Errorf("Vout = %f, want ~%f", vout, expectedVout)
	}
	if !almostEqual(ratio, 0.319727, 0.001) {
		t.Errorf("ratio = %f, want ~0.3197", ratio)
	}
	if current <= 0 {
		t.Error("current should be positive")
	}
}

func TestParseResistorValue(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"4700", 4700},
		{"4.7k", 4700},
		{"4k7", 4700},
		{"1M", 1e6},
		{"470R", 470},
		{"2.2M", 2.2e6},
		{"10k", 10000},
		{"100", 100},
	}
	for _, tt := range tests {
		got, err := parseResistorValue(tt.input)
		if err != nil {
			t.Errorf("parseResistorValue(%q) error: %v", tt.input, err)
			continue
		}
		if !almostEqual(got, tt.want, 0.01) {
			t.Errorf("parseResistorValue(%q) = %g, want %g", tt.input, got, tt.want)
		}
	}
}

func TestParseResistorValueErrors(t *testing.T) {
	badInputs := []string{"", "abc", "k"}
	for _, input := range badInputs {
		_, err := parseResistorValue(input)
		if err == nil {
			t.Errorf("parseResistorValue(%q) expected error, got nil", input)
		}
	}
}

func TestColorBands4(t *testing.T) {
	// 4700 Ohm: Yellow Violet Red Gold
	bands, err := colorBands4(4700)
	if err != nil {
		t.Fatalf("colorBands4(4700) error: %v", err)
	}
	expected := []string{"Yellow", "Violet", "Red", "Gold"}
	if len(bands.Bands) != 4 {
		t.Fatalf("expected 4 bands, got %d", len(bands.Bands))
	}
	for i, want := range expected {
		if bands.Bands[i] != want {
			t.Errorf("band %d = %q, want %q", i, bands.Bands[i], want)
		}
	}
}

func TestColorBands5(t *testing.T) {
	// 4700 Ohm: Yellow Violet Black Brown Brown
	bands, err := colorBands5(4700)
	if err != nil {
		t.Fatalf("colorBands5(4700) error: %v", err)
	}
	expected := []string{"Yellow", "Violet", "Black", "Brown", "Brown"}
	if len(bands.Bands) != 5 {
		t.Fatalf("expected 5 bands, got %d", len(bands.Bands))
	}
	for i, want := range expected {
		if bands.Bands[i] != want {
			t.Errorf("band %d = %q, want %q", i, bands.Bands[i], want)
		}
	}
}

func TestFormatOhms(t *testing.T) {
	tests := []struct {
		ohms float64
		want string
	}{
		{4700, "4.7 kΩ"},
		{1e6, "1 MΩ"},
		{100, "100 Ω"},
	}
	for _, tt := range tests {
		got := formatOhms(tt.ohms)
		if got != tt.want {
			t.Errorf("formatOhms(%g) = %q, want %q", tt.ohms, got, tt.want)
		}
	}
}

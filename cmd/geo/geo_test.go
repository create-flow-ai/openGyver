package geo

import (
	"math"
	"testing"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestHaversine_NYCtoLondon(t *testing.T) {
	// NYC: 40.7128, -74.0060
	// London: 51.5074, -0.1278
	// Expected: ~5570 km
	km := haversine(40.7128, -74.0060, 51.5074, -0.1278)
	if !almostEqual(km, 5570.0, 10.0) {
		t.Errorf("NYC to London = %.2f km, want ~5570 km", km)
	}
}

func TestHaversine_SamePoint(t *testing.T) {
	km := haversine(40.7128, -74.0060, 40.7128, -74.0060)
	if km != 0 {
		t.Errorf("same point distance = %f, want 0", km)
	}
}

func TestHaversine_Antipodal(t *testing.T) {
	// Points roughly opposite: ~20000 km apart.
	km := haversine(0, 0, 0, 180)
	if !almostEqual(km, math.Pi*earthRadiusKm, 1.0) {
		t.Errorf("antipodal distance = %.2f km, want ~%.2f km", km, math.Pi*earthRadiusKm)
	}
}

func TestDegreesToRadians(t *testing.T) {
	if !almostEqual(degreesToRadians(180), math.Pi, 1e-10) {
		t.Errorf("180 degrees = %f rad, want Pi", degreesToRadians(180))
	}
	if !almostEqual(degreesToRadians(90), math.Pi/2, 1e-10) {
		t.Errorf("90 degrees = %f rad, want Pi/2", degreesToRadians(90))
	}
	if !almostEqual(degreesToRadians(0), 0, 1e-10) {
		t.Errorf("0 degrees = %f rad, want 0", degreesToRadians(0))
	}
}

func TestDecimalToDMS(t *testing.T) {
	// 40.7128 -> 40 degrees 42' 46.08" N
	result := decimalToDMS(40.7128, false)
	if result.Degrees != 40 {
		t.Errorf("degrees = %d, want 40", result.Degrees)
	}
	if result.Minutes != 42 {
		t.Errorf("minutes = %d, want 42", result.Minutes)
	}
	if !almostEqual(result.Seconds, 46.08, 0.1) {
		t.Errorf("seconds = %f, want ~46.08", result.Seconds)
	}
	if result.Direction != "N" {
		t.Errorf("direction = %q, want N", result.Direction)
	}
}

func TestDecimalToDMS_Negative(t *testing.T) {
	// -74.0060 is longitude (> 90 in abs), so isLongitude should be true.
	result := decimalToDMS(-74.0060, true)
	if result.Direction != "W" {
		t.Errorf("direction = %q, want W", result.Direction)
	}
	if result.Degrees != 74 {
		t.Errorf("degrees = %d, want 74", result.Degrees)
	}
}

func TestParseDMS(t *testing.T) {
	tests := []struct {
		input   string
		wantDD  float64
		wantErr bool
	}{
		{`40°42'46.1"N`, 40.712806, false},
		{`40 42 46.1 N`, 40.712806, false},
		{`-40°42'46.1"`, -40.712806, false},
		{"not dms", 0, true},
	}
	for _, tt := range tests {
		if !isDMS(tt.input) && !tt.wantErr {
			t.Errorf("isDMS(%q) = false, want true", tt.input)
			continue
		}
		if tt.wantErr {
			if isDMS(tt.input) {
				t.Errorf("isDMS(%q) = true, want false", tt.input)
			}
			continue
		}
		dd, _, err := parseDMS(tt.input)
		if err != nil {
			t.Errorf("parseDMS(%q) error: %v", tt.input, err)
			continue
		}
		if !almostEqual(dd, tt.wantDD, 0.001) {
			t.Errorf("parseDMS(%q) = %f, want ~%f", tt.input, dd, tt.wantDD)
		}
	}
}

func TestLatLonToUTM(t *testing.T) {
	// NYC: 40.7128, -74.0060 -> Zone 18T
	result := latLonToUTM(40.7128, -74.0060)
	if result.Zone != 18 {
		t.Errorf("zone = %d, want 18", result.Zone)
	}
	if result.Letter != "T" {
		t.Errorf("letter = %q, want T", result.Letter)
	}
	// Easting should be around 583960.
	if !almostEqual(result.Easting, 583960, 100) {
		t.Errorf("easting = %f, want ~583960", result.Easting)
	}
	// Northing should be around 4507000 (exact value depends on implementation).
	if !almostEqual(result.Northing, 4507351, 500) {
		t.Errorf("northing = %f, want ~4507351", result.Northing)
	}
}

func TestLatLonToUTM_Equator(t *testing.T) {
	// 0, 0 -> Zone 31N
	result := latLonToUTM(0, 0)
	if result.Zone != 31 {
		t.Errorf("zone = %d, want 31", result.Zone)
	}
	if result.Letter != "N" {
		t.Errorf("letter = %q, want N", result.Letter)
	}
}

func TestLatLonToUTM_SouthernHemisphere(t *testing.T) {
	// Sydney: -33.8568, 151.2153
	result := latLonToUTM(-33.8568, 151.2153)
	if result.Zone != 56 {
		t.Errorf("zone = %d, want 56", result.Zone)
	}
	// Should have "H" letter for southern hemisphere in the -40 to -32 band.
	// -33.8568 is >= -40, so letter "H".
	if result.Letter != "H" {
		t.Errorf("letter = %q, want H", result.Letter)
	}
}

func TestUTMLetterDesignator(t *testing.T) {
	tests := []struct {
		lat  float64
		want string
	}{
		{80, "X"},
		{50, "U"},
		{40, "T"},
		{0, "N"},
		{-10, "L"},
		{-80, "C"},
	}
	for _, tt := range tests {
		got := utmLetterDesignator(tt.lat)
		if got != tt.want {
			t.Errorf("utmLetterDesignator(%f) = %q, want %q", tt.lat, got, tt.want)
		}
	}
}

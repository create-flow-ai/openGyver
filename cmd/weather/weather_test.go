package weather

import (
	"encoding/json"
	"testing"
)

func TestGeocode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping geocode test in short mode (requires network)")
	}

	result, err := geocode("London")
	if err != nil {
		t.Fatalf("geocode(London) error: %v", err)
	}
	if result.Name == "" {
		t.Error("geocode result name is empty")
	}
	if result.Lat == 0 && result.Lon == 0 {
		t.Error("geocode result has zero coordinates")
	}
}

func TestGeocode_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping geocode test in short mode (requires network)")
	}

	_, err := geocode("zzzxyznonexistentcity12345")
	if err == nil {
		t.Error("expected error for nonexistent city")
	}
}

func TestWeatherDescription(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{0, "Clear sky"},
		{1, "Mainly clear"},
		{2, "Partly cloudy"},
		{3, "Overcast"},
		{45, "Fog"},
		{51, "Light drizzle"},
		{61, "Slight rain"},
		{63, "Moderate rain"},
		{65, "Heavy rain"},
		{71, "Slight snow"},
		{95, "Thunderstorm"},
		{999, "Unknown"},
	}
	for _, tt := range tests {
		got := weatherDescription(tt.code)
		if got != tt.want {
			t.Errorf("weatherDescription(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2024-12-25T08:30", "8:30 AM"},
		{"2024-12-25T14:00", "2:00 PM"},
		{"2024-12-25T00:00", "12:00 AM"},
		{"2024-12-25T12:00", "12:00 PM"},
	}
	for _, tt := range tests {
		got := formatTime(tt.input)
		if got != tt.want {
			t.Errorf("formatTime(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatTime_InvalidInput(t *testing.T) {
	// Invalid input should return the original string.
	input := "not-a-time"
	got := formatTime(input)
	if got != input {
		t.Errorf("formatTime(%q) = %q, want original string returned", input, got)
	}
}

func TestJsonFloat_Null(t *testing.T) {
	var f jsonFloat
	err := json.Unmarshal([]byte("null"), &f)
	if err != nil {
		t.Fatalf("unmarshal null error: %v", err)
	}
	if float64(f) != 0 {
		t.Errorf("jsonFloat(null) = %f, want 0", float64(f))
	}
}

func TestJsonFloat_Number(t *testing.T) {
	var f jsonFloat
	err := json.Unmarshal([]byte("42.5"), &f)
	if err != nil {
		t.Fatalf("unmarshal number error: %v", err)
	}
	if float64(f) != 42.5 {
		t.Errorf("jsonFloat(42.5) = %f, want 42.5", float64(f))
	}
}

func TestJsonFloat_Zero(t *testing.T) {
	var f jsonFloat
	err := json.Unmarshal([]byte("0"), &f)
	if err != nil {
		t.Fatalf("unmarshal 0 error: %v", err)
	}
	if float64(f) != 0 {
		t.Errorf("jsonFloat(0) = %f, want 0", float64(f))
	}
}

func TestJsonFloat_Invalid(t *testing.T) {
	var f jsonFloat
	err := json.Unmarshal([]byte(`"not a number"`), &f)
	if err == nil {
		t.Error("expected error for invalid jsonFloat input")
	}
}

func TestSafeIndex(t *testing.T) {
	arr := []jsonFloat{10, 20, 30}

	if safeIndex(arr, 0) != 10 {
		t.Errorf("safeIndex(arr, 0) = %f, want 10", safeIndex(arr, 0))
	}
	if safeIndex(arr, 2) != 30 {
		t.Errorf("safeIndex(arr, 2) = %f, want 30", safeIndex(arr, 2))
	}
	// Out of range: returns last element.
	if safeIndex(arr, 10) != 30 {
		t.Errorf("safeIndex(arr, 10) = %f, want 30 (last)", safeIndex(arr, 10))
	}
	// Empty array.
	if safeIndex(nil, 0) != 0 {
		t.Errorf("safeIndex(nil, 0) = %f, want 0", safeIndex(nil, 0))
	}
}

func TestSafeIdxFloat(t *testing.T) {
	arr := []jsonFloat{1.5, 2.5}
	if safeIdxFloat(arr, 0) != 1.5 {
		t.Errorf("safeIdxFloat(arr, 0) = %f", safeIdxFloat(arr, 0))
	}
	// Out of range.
	if safeIdxFloat(arr, 5) != 0 {
		t.Errorf("safeIdxFloat(arr, 5) = %f, want 0", safeIdxFloat(arr, 5))
	}
}

func TestSafeIdxStr(t *testing.T) {
	arr := []string{"a", "b"}
	if safeIdxStr(arr, 0) != "a" {
		t.Errorf("safeIdxStr(arr, 0) = %q", safeIdxStr(arr, 0))
	}
	if safeIdxStr(arr, 5) != "" {
		t.Errorf("safeIdxStr(arr, 5) = %q, want empty", safeIdxStr(arr, 5))
	}
}

package epoch

import (
	"testing"
	"time"
)

// base is 2024-01-15T14:00:00Z (epoch 1705327200)
var base = time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)

// ---------------------------------------------------------------------------
// computeAdd
// ---------------------------------------------------------------------------

func TestComputeAdd_Hours(t *testing.T) {
	result := computeAdd(base, 0, 0, 0, 0, 2, 0)
	want := base.Add(2 * time.Hour)
	if !result.Equal(want) {
		t.Errorf("got %d, want %d", result.Unix(), want.Unix())
	}
}

func TestComputeAdd_Minutes(t *testing.T) {
	result := computeAdd(base, 0, 0, 0, 0, 0, 45)
	want := base.Add(45 * time.Minute)
	if !result.Equal(want) {
		t.Errorf("got %d, want %d", result.Unix(), want.Unix())
	}
}

func TestComputeAdd_Days(t *testing.T) {
	result := computeAdd(base, 0, 0, 0, 30, 0, 0)
	want := time.Date(2024, 2, 14, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeAdd_Weeks(t *testing.T) {
	result := computeAdd(base, 0, 0, 2, 0, 0, 0)
	want := time.Date(2024, 1, 29, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeAdd_Months(t *testing.T) {
	result := computeAdd(base, 0, 3, 0, 0, 0, 0)
	want := time.Date(2024, 4, 15, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeAdd_Years(t *testing.T) {
	result := computeAdd(base, 1, 0, 0, 0, 0, 0)
	want := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeAdd_Combined(t *testing.T) {
	result := computeAdd(base, 1, 6, 0, 15, 3, 30)
	want := time.Date(2025, 7, 30, 17, 30, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeAdd_Zero(t *testing.T) {
	result := computeAdd(base, 0, 0, 0, 0, 0, 0)
	if !result.Equal(base) {
		t.Errorf("adding zero should return same time")
	}
}

// ---------------------------------------------------------------------------
// computeSubtract
// ---------------------------------------------------------------------------

func TestComputeSubtract_Hours(t *testing.T) {
	result := computeSubtract(base, 0, 0, 0, 0, 2, 0)
	want := base.Add(-2 * time.Hour)
	if !result.Equal(want) {
		t.Errorf("got %d, want %d", result.Unix(), want.Unix())
	}
}

func TestComputeSubtract_Days(t *testing.T) {
	result := computeSubtract(base, 0, 0, 0, 30, 0, 0)
	want := time.Date(2023, 12, 16, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeSubtract_Weeks(t *testing.T) {
	result := computeSubtract(base, 0, 0, 1, 0, 0, 0)
	want := time.Date(2024, 1, 8, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeSubtract_Months(t *testing.T) {
	result := computeSubtract(base, 0, 3, 0, 0, 0, 0)
	want := time.Date(2023, 10, 15, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeSubtract_Years(t *testing.T) {
	result := computeSubtract(base, 1, 0, 0, 0, 0, 0)
	want := time.Date(2023, 1, 15, 14, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeSubtract_Combined(t *testing.T) {
	result := computeSubtract(base, 1, 0, 0, 7, 12, 0)
	want := time.Date(2023, 1, 8, 2, 0, 0, 0, time.UTC)
	if !result.Equal(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestComputeSubtract_Zero(t *testing.T) {
	result := computeSubtract(base, 0, 0, 0, 0, 0, 0)
	if !result.Equal(base) {
		t.Errorf("subtracting zero should return same time")
	}
}

// ---------------------------------------------------------------------------
// resolveBase
// ---------------------------------------------------------------------------

func TestResolveBase_FromEpoch(t *testing.T) {
	result := resolveBase(1705334400)
	if result.Unix() != 1705334400 {
		t.Errorf("got %d, want 1705334400", result.Unix())
	}
}

func TestResolveBase_CurrentTime(t *testing.T) {
	before := time.Now().Unix()
	result := resolveBase(0)
	after := time.Now().Unix()
	if result.Unix() < before || result.Unix() > after {
		t.Error("resolveBase(0) should return current time")
	}
}

// ---------------------------------------------------------------------------
// Add subcommand metadata
// ---------------------------------------------------------------------------

func TestAddCmd_Metadata(t *testing.T) {
	if addCmd.Use != "add" {
		t.Errorf("unexpected Use: %s", addCmd.Use)
	}
	if addCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestSubtractCmd_Metadata(t *testing.T) {
	if subtractCmd.Use != "subtract" {
		t.Errorf("unexpected Use: %s", subtractCmd.Use)
	}
	if subtractCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestAddCmd_Flags(t *testing.T) {
	f := addCmd.Flags()
	for _, name := range []string{"from", "hours", "minutes", "days", "weeks", "months", "years"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found on add", name)
		}
	}
}

func TestSubtractCmd_Flags(t *testing.T) {
	f := subtractCmd.Flags()
	for _, name := range []string{"from", "hours", "minutes", "days", "weeks", "months", "years"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found on subtract", name)
		}
	}
}

// ---------------------------------------------------------------------------
// Add/Subtract are inverse operations
// ---------------------------------------------------------------------------

func TestAddSubtract_Inverse(t *testing.T) {
	added := computeAdd(base, 1, 2, 3, 4, 5, 6)
	restored := computeSubtract(added, 1, 2, 3, 4, 5, 6)
	if !restored.Equal(base) {
		t.Errorf("add then subtract should return original: got %v, want %v", restored, base)
	}
}

// ---------------------------------------------------------------------------
// Edge: month overflow
// ---------------------------------------------------------------------------

func TestComputeAdd_MonthOverflow(t *testing.T) {
	// Jan 31 + 1 month = Mar 2 or 3 (Go normalizes)
	jan31 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	result := computeAdd(jan31, 0, 1, 0, 0, 0, 0)
	// Go's AddDate normalizes: Jan 31 + 1mo = Feb 31 → Mar 2
	if result.Month() != 3 || result.Day() != 2 {
		t.Errorf("expected March 2, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// formatEpoch
// ---------------------------------------------------------------------------

func TestFormatEpoch_Seconds(t *testing.T) {
	ms, us, ns = false, false, false
	s := formatEpoch(base)
	if s != "1705327200" {
		t.Errorf("got %s, want 1705327200", s)
	}
}

func TestFormatEpoch_Milliseconds(t *testing.T) {
	ms, us, ns = true, false, false
	defer func() { ms = false }()
	s := formatEpoch(base)
	if s != "1705327200000" {
		t.Errorf("got %s, want 1705327200000", s)
	}
}

package accessibility

import (
	"math"
	"testing"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestContrastRatio_BlackWhite(t *testing.T) {
	l1 := relativeLuminance(255, 255, 255) // white
	l2 := relativeLuminance(0, 0, 0)       // black
	ratio := contrastRatio(l1, l2)

	if !almostEqual(ratio, 21.0, 0.1) {
		t.Errorf("black/white contrast ratio = %.2f, want 21.0", ratio)
	}
}

func TestContrastRatio_SameColor(t *testing.T) {
	l := relativeLuminance(128, 128, 128)
	ratio := contrastRatio(l, l)

	if !almostEqual(ratio, 1.0, 0.01) {
		t.Errorf("same color contrast ratio = %.2f, want 1.0", ratio)
	}
}

func TestContrastRatio_OrderIndependent(t *testing.T) {
	l1 := relativeLuminance(255, 0, 0)
	l2 := relativeLuminance(0, 0, 255)

	r1 := contrastRatio(l1, l2)
	r2 := contrastRatio(l2, l1)

	if !almostEqual(r1, r2, 0.001) {
		t.Errorf("contrast ratio is order-dependent: %.4f vs %.4f", r1, r2)
	}
}

func TestRelativeLuminance(t *testing.T) {
	// Black should be 0.
	lBlack := relativeLuminance(0, 0, 0)
	if lBlack != 0 {
		t.Errorf("luminance of black = %f, want 0", lBlack)
	}

	// White should be 1.
	lWhite := relativeLuminance(255, 255, 255)
	if !almostEqual(lWhite, 1.0, 0.001) {
		t.Errorf("luminance of white = %f, want 1.0", lWhite)
	}

	// Mid-gray should be between 0 and 1.
	lGray := relativeLuminance(128, 128, 128)
	if lGray <= 0 || lGray >= 1 {
		t.Errorf("luminance of gray = %f, expected between 0 and 1", lGray)
	}
}

func TestLinearize(t *testing.T) {
	// Low value (below threshold 0.03928).
	low := linearize(0.01)
	if !almostEqual(low, 0.01/12.92, 1e-6) {
		t.Errorf("linearize(0.01) = %f, want %f", low, 0.01/12.92)
	}

	// High value.
	high := linearize(0.5)
	expected := math.Pow((0.5+0.055)/1.055, 2.4)
	if !almostEqual(high, expected, 1e-6) {
		t.Errorf("linearize(0.5) = %f, want %f", high, expected)
	}
}

func TestParseColor_Hex(t *testing.T) {
	tests := []struct {
		input      string
		r, g, b    uint8
		shouldFail bool
	}{
		{"#ffffff", 255, 255, 255, false},
		{"#000000", 0, 0, 0, false},
		{"#fff", 255, 255, 255, false},
		{"#000", 0, 0, 0, false},
		{"ff0000", 255, 0, 0, false},
		{"#abc", 170, 187, 204, false},
	}
	for _, tt := range tests {
		r, g, b, err := parseColor(tt.input)
		if tt.shouldFail {
			if err == nil {
				t.Errorf("parseColor(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseColor(%q) error: %v", tt.input, err)
			continue
		}
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("parseColor(%q) = (%d,%d,%d), want (%d,%d,%d)", tt.input, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestParseColor_RGB(t *testing.T) {
	r, g, b, err := parseColor("rgb(255,0,128)")
	if err != nil {
		t.Fatalf("parseColor(rgb) error: %v", err)
	}
	if r != 255 || g != 0 || b != 128 {
		t.Errorf("parseColor(rgb) = (%d,%d,%d), want (255,0,128)", r, g, b)
	}
}

func TestParseColor_Named(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
	}{
		{"black", 0, 0, 0},
		{"white", 255, 255, 255},
		{"red", 255, 0, 0},
		{"green", 0, 128, 0},
		{"blue", 0, 0, 255},
	}
	for _, tt := range tests {
		r, g, b, err := parseColor(tt.name)
		if err != nil {
			t.Errorf("parseColor(%q) error: %v", tt.name, err)
			continue
		}
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("parseColor(%q) = (%d,%d,%d), want (%d,%d,%d)", tt.name, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestParseColor_Invalid(t *testing.T) {
	_, _, _, err := parseColor("not-a-color")
	if err == nil {
		t.Error("expected error for invalid color")
	}
}

func TestCountSyllablesWord(t *testing.T) {
	tests := []struct {
		word string
		want int
	}{
		{"the", 1},
		{"hello", 2},
		{"beautiful", 3},
		{"a", 1},
	}
	for _, tt := range tests {
		got := countSyllablesWord(tt.word)
		if got != tt.want {
			t.Errorf("countSyllablesWord(%q) = %d, want %d", tt.word, got, tt.want)
		}
	}
}

func TestCountWords(t *testing.T) {
	if countWords("hello world") != 2 {
		t.Errorf("countWords(\"hello world\") = %d, want 2", countWords("hello world"))
	}
	if countWords("one") != 1 {
		t.Errorf("countWords(\"one\") = %d, want 1", countWords("one"))
	}
	if countWords("  ") != 0 {
		t.Errorf("countWords(\"  \") = %d, want 0", countWords("  "))
	}
}

func TestCountSentences(t *testing.T) {
	if countSentences("Hello. World!") != 2 {
		t.Errorf("countSentences = %d, want 2", countSentences("Hello. World!"))
	}
	if countSentences("No punctuation") != 1 {
		t.Errorf("countSentences(no punct) = %d, want 1", countSentences("No punctuation"))
	}
}

func TestFleschReadingEase_SimpleText(t *testing.T) {
	// Simple text should score high (easy to read).
	text := "The cat sat on the mat. The dog ran in the park."
	words := countWords(text)
	sentences := countSentences(text)
	syllables := countSyllablesText(text)

	if words == 0 || sentences == 0 {
		t.Fatal("words or sentences is 0")
	}

	wps := float64(words) / float64(sentences)
	spw := float64(syllables) / float64(words)
	fre := 206.835 - 1.015*wps - 84.6*spw

	// Simple text should be "Easy" or higher (FRE >= 70).
	if fre < 70 {
		t.Errorf("Flesch Reading Ease = %.1f, expected >= 70 for simple text", fre)
	}
}

func TestFleschLabel(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{95, "Very Easy"},
		{85, "Easy"},
		{75, "Fairly Easy"},
		{65, "Standard"},
		{55, "Fairly Difficult"},
		{35, "Difficult"},
		{10, "Very Confusing"},
	}
	for _, tt := range tests {
		got := fleschLabel(tt.score)
		if got != tt.want {
			t.Errorf("fleschLabel(%f) = %q, want %q", tt.score, got, tt.want)
		}
	}
}

func TestPassFail(t *testing.T) {
	if passFail(true) != "PASS" {
		t.Error("passFail(true) != PASS")
	}
	if passFail(false) != "FAIL" {
		t.Error("passFail(false) != FAIL")
	}
}

func TestCountComplexWords(t *testing.T) {
	// "beautiful" has 3 syllables, "the" has 1, "cat" has 1.
	text := "The beautiful cat sat on a beautiful mat."
	count := countComplexWords(text)
	if count != 2 {
		t.Errorf("countComplexWords = %d, want 2", count)
	}
}

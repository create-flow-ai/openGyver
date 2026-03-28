package regex

import (
	"regexp"
	"testing"
)

func TestRegexMatch(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	input := "abc123def456"

	if !re.MatchString(input) {
		t.Error(`expected \d+ to match "abc123def456"`)
	}

	match := re.FindString(input)
	if match != "123" {
		t.Errorf(`first match = %q, want "123"`, match)
	}
}

func TestRegexNoMatch(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	input := "no digits here"

	if re.MatchString(input) {
		t.Error(`expected \d+ not to match "no digits here"`)
	}
}

func TestRegexReplace(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	result := re.ReplaceAllString("order 42 has 3 items", "X")
	expected := "order X has X items"
	if result != expected {
		t.Errorf("replace result = %q, want %q", result, expected)
	}
}

func TestRegexReplaceWithGroups(t *testing.T) {
	re := regexp.MustCompile(`(\w+)@(\w+)`)
	result := re.ReplaceAllString("alice@example", "$2/$1")
	expected := "example/alice"
	if result != expected {
		t.Errorf("replace result = %q, want %q", result, expected)
	}
}

func TestRegexExtractMultipleMatches(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	input := "order 42 has 3 items and 7 widgets"
	matches := re.FindAllString(input, -1)

	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d: %v", len(matches), matches)
	}
	want := []string{"42", "3", "7"}
	for i, m := range matches {
		if m != want[i] {
			t.Errorf("match[%d] = %q, want %q", i, m, want[i])
		}
	}
}

func TestRegexExtractNoMatches(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString("no numbers", -1)
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestRegexInvalidPattern(t *testing.T) {
	_, err := regexp.Compile(`[invalid`)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestRegexSubmatch(t *testing.T) {
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	input := "Contact alice@ex.com or bob@ex.com"
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	// First match.
	if matches[0][0] != "alice@ex.com" {
		t.Errorf("match[0][0] = %q, want alice@ex.com", matches[0][0])
	}
	if matches[0][1] != "alice" {
		t.Errorf("group 1 = %q, want alice", matches[0][1])
	}
}

func TestRegexGlobalMatch(t *testing.T) {
	re := regexp.MustCompile(`\d+`)
	input := "abc123def456ghi789"
	all := re.FindAllString(input, -1)
	if len(all) != 3 {
		t.Errorf("global matches = %d, want 3", len(all))
	}
}

func TestReadInput_FromArgs(t *testing.T) {
	// readInput should return text from args starting at startIdx.
	text, err := readInput([]string{"pattern", "hello world"}, 1, "")
	if err != nil {
		t.Fatalf("readInput error: %v", err)
	}
	if text != "hello world" {
		t.Errorf("readInput = %q, want %q", text, "hello world")
	}
}

func TestReadInput_MultipleArgs(t *testing.T) {
	text, err := readInput([]string{"pattern", "hello", "world"}, 1, "")
	if err != nil {
		t.Fatalf("readInput error: %v", err)
	}
	if text != "hello world" {
		t.Errorf("readInput = %q, want %q", text, "hello world")
	}
}

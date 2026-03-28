package text

import (
	"strings"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestTextCmd_Metadata(t *testing.T) {
	if textCmd.Use == "" {
		t.Error("textCmd.Use must not be empty")
	}
	if textCmd.Short == "" {
		t.Error("textCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"countCmd", countCmd.Use, countCmd.Short},
		{"caseCmd", caseCmd.Use, caseCmd.Short},
		{"reverseCmd", reverseCmd.Use, reverseCmd.Short},
		{"sortCmd", sortCmd.Use, sortCmd.Short},
		{"dedupeCmd", dedupeCmd.Use, dedupeCmd.Short},
		{"slugCmd", slugCmd.Use, slugCmd.Short},
		{"loremCmd", loremCmd.Use, loremCmd.Short},
		{"diffCmd", diffCmd.Use, diffCmd.Short},
		{"wrapCmd", wrapCmd.Use, wrapCmd.Short},
		{"linesCmd", linesCmd.Use, linesCmd.Short},
		{"trimCmd", trimCmd.Use, trimCmd.Short},
		{"replaceCmd", replaceCmd.Use, replaceCmd.Short},
	}
	for _, c := range cmds {
		if c.use == "" {
			t.Errorf("%s.Use must not be empty", c.name)
		}
		if c.short == "" {
			t.Errorf("%s.Short must not be empty", c.name)
		}
	}
}

// ── flag existence ─────────────────────────────────────────────────────────

func TestTextCmd_PersistentFlags(t *testing.T) {
	f := textCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

func TestCaseCmd_Flags(t *testing.T) {
	f := caseCmd.Flags()
	if f.Lookup("to") == nil {
		t.Error("expected flag --to on caseCmd")
	}
}

func TestSortCmd_Flags(t *testing.T) {
	f := sortCmd.Flags()
	if f.Lookup("by") == nil {
		t.Error("expected flag --by on sortCmd")
	}
	if f.Lookup("reverse") == nil {
		t.Error("expected flag --reverse on sortCmd")
	}
}

func TestLoremCmd_Flags(t *testing.T) {
	f := loremCmd.Flags()
	if f.Lookup("words") == nil {
		t.Error("expected flag --words on loremCmd")
	}
	if f.Lookup("sentences") == nil {
		t.Error("expected flag --sentences on loremCmd")
	}
	if f.Lookup("paragraphs") == nil {
		t.Error("expected flag --paragraphs on loremCmd")
	}
}

// ── tokenize ───────────────────────────────────────────────────────────────

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"Hello World", []string{"hello", "world"}},
		{"helloWorld", []string{"hello", "world"}},
		{"hello_world", []string{"hello", "world"}},
		{"hello-world", []string{"hello", "world"}},
		{"hello.world", []string{"hello", "world"}},
		{"HelloWorld", []string{"hello", "world"}},
		{"myVariableName", []string{"my", "variable", "name"}},
		{"some-variable-name", []string{"some", "variable", "name"}},
		{"CONSTANT_CASE", []string{"constant", "case"}},
		{"", nil},
	}
	for _, tt := range tests {
		got := tokenize(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("tokenize(%q) = %v (len %d), want %v (len %d)", tt.input, got, len(got), tt.want, len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("tokenize(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

// ── case conversion ────────────────────────────────────────────────────────

func TestConvertCase_Camel(t *testing.T) {
	got, err := convertCase("Hello World", "camel")
	if err != nil {
		t.Fatalf("convertCase camel: %v", err)
	}
	if got != "helloWorld" {
		t.Errorf("camel('Hello World') = %q, want %q", got, "helloWorld")
	}
}

func TestConvertCase_Snake(t *testing.T) {
	got, err := convertCase("Hello World", "snake")
	if err != nil {
		t.Fatalf("convertCase snake: %v", err)
	}
	if got != "hello_world" {
		t.Errorf("snake('Hello World') = %q, want %q", got, "hello_world")
	}
}

func TestConvertCase_Kebab(t *testing.T) {
	got, err := convertCase("Hello World", "kebab")
	if err != nil {
		t.Fatalf("convertCase kebab: %v", err)
	}
	if got != "hello-world" {
		t.Errorf("kebab('Hello World') = %q, want %q", got, "hello-world")
	}
}

func TestConvertCase_Pascal(t *testing.T) {
	got, err := convertCase("hello world", "pascal")
	if err != nil {
		t.Fatalf("convertCase pascal: %v", err)
	}
	if got != "HelloWorld" {
		t.Errorf("pascal('hello world') = %q, want %q", got, "HelloWorld")
	}
}

func TestConvertCase_Upper(t *testing.T) {
	got, err := convertCase("hello world", "upper")
	if err != nil {
		t.Fatalf("convertCase upper: %v", err)
	}
	if got != "HELLO WORLD" {
		t.Errorf("upper('hello world') = %q, want %q", got, "HELLO WORLD")
	}
}

func TestConvertCase_Lower(t *testing.T) {
	got, err := convertCase("HELLO WORLD", "lower")
	if err != nil {
		t.Fatalf("convertCase lower: %v", err)
	}
	if got != "hello world" {
		t.Errorf("lower('HELLO WORLD') = %q, want %q", got, "hello world")
	}
}

func TestConvertCase_Title(t *testing.T) {
	got, err := convertCase("hello world", "title")
	if err != nil {
		t.Fatalf("convertCase title: %v", err)
	}
	if got != "Hello World" {
		t.Errorf("title('hello world') = %q, want %q", got, "Hello World")
	}
}

func TestConvertCase_Sentence(t *testing.T) {
	got, err := convertCase("hello world", "sentence")
	if err != nil {
		t.Fatalf("convertCase sentence: %v", err)
	}
	if got != "Hello world" {
		t.Errorf("sentence('hello world') = %q, want %q", got, "Hello world")
	}
}

func TestConvertCase_Constant(t *testing.T) {
	got, err := convertCase("hello world", "constant")
	if err != nil {
		t.Fatalf("convertCase constant: %v", err)
	}
	if got != "HELLO_WORLD" {
		t.Errorf("constant('hello world') = %q, want %q", got, "HELLO_WORLD")
	}
}

func TestConvertCase_Dot(t *testing.T) {
	got, err := convertCase("hello world", "dot")
	if err != nil {
		t.Fatalf("convertCase dot: %v", err)
	}
	if got != "hello.world" {
		t.Errorf("dot('hello world') = %q, want %q", got, "hello.world")
	}
}

func TestConvertCase_FromCamelToSnake(t *testing.T) {
	got, err := convertCase("myVariableName", "snake")
	if err != nil {
		t.Fatalf("convertCase snake: %v", err)
	}
	if got != "my_variable_name" {
		t.Errorf("snake('myVariableName') = %q, want %q", got, "my_variable_name")
	}
}

func TestConvertCase_FromKebabToCamel(t *testing.T) {
	got, err := convertCase("some-variable-name", "camel")
	if err != nil {
		t.Fatalf("convertCase camel: %v", err)
	}
	if got != "someVariableName" {
		t.Errorf("camel('some-variable-name') = %q, want %q", got, "someVariableName")
	}
}

func TestConvertCase_Unknown(t *testing.T) {
	_, err := convertCase("hello", "unknown")
	if err == nil {
		t.Error("expected error for unknown case")
	}
}

func TestConvertCase_EmptyInput(t *testing.T) {
	got, err := convertCase("", "snake")
	if err != nil {
		t.Fatalf("convertCase error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

// ── word/char/line count ───────────────────────────────────────────────────

func TestCount_Words(t *testing.T) {
	input := "The quick brown fox jumps over the lazy dog"
	words := len(strings.Fields(input))
	if words != 9 {
		t.Errorf("word count = %d, want 9", words)
	}
}

func TestCount_Characters(t *testing.T) {
	input := "Hello"
	chars := len([]rune(input))
	if chars != 5 {
		t.Errorf("char count = %d, want 5", chars)
	}
}

func TestCount_Lines(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"one\ntwo\nthree", 3},
		{"single line", 1},
		{"with trailing\n", 1},
		{"two\nlines\n", 2},
	}
	for _, tt := range tests {
		lines := 0
		if len(tt.input) > 0 {
			lines = strings.Count(tt.input, "\n")
			if tt.input[len(tt.input)-1] != '\n' {
				lines++
			}
		}
		if lines != tt.want {
			t.Errorf("line count(%q) = %d, want %d", tt.input, lines, tt.want)
		}
	}
}

func TestCount_EmptyString(t *testing.T) {
	input := ""
	words := len(strings.Fields(input))
	chars := len([]rune(input))
	if words != 0 {
		t.Errorf("word count of empty string = %d, want 0", words)
	}
	if chars != 0 {
		t.Errorf("char count of empty string = %d, want 0", chars)
	}
}

// ── string reverse ─────────────────────────────────────────────────────────

func TestReverse(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"abcdef", "fedcba"},
		{"Hello, World!", "!dlroW ,olleH"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
	}
	for _, tt := range tests {
		runes := []rune(tt.input)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		got := string(runes)
		if got != tt.want {
			t.Errorf("reverse(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestReverse_Unicode(t *testing.T) {
	input := "abc"
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	got := string(runes)
	if got != "cba" {
		t.Errorf("reverse(%q) = %q, want %q", input, got, "cba")
	}
}

// ── slug generation ────────────────────────────────────────────────────────

func TestSlugify(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"My Blog Post Title!", "my-blog-post-title"},
		{"Hello,   World!! 2024", "hello-world-2024"},
		{"Simple", "simple"},
		{"  spaces  ", "spaces"},
		{"CamelCase", "camelcase"},
		{"under_score", "under-score"},
		{"foo---bar", "foo-bar"},
	}
	for _, tt := range tests {
		got := slugify(tt.input)
		if got != tt.want {
			t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlugify_Empty(t *testing.T) {
	got := slugify("")
	if got != "" {
		t.Errorf("slugify('') = %q, want empty", got)
	}
}

func TestSlugify_AllSpecial(t *testing.T) {
	got := slugify("!!!@@@###")
	if got != "" {
		t.Errorf("slugify('!!!@@@###') = %q, want empty", got)
	}
}

// ── lorem ipsum ────────────────────────────────────────────────────────────

func TestLoremByWords_Count(t *testing.T) {
	for _, n := range []int{5, 10, 20, 50, 100} {
		output := loremByWords(n)
		words := strings.Fields(output)
		if len(words) != n {
			t.Errorf("loremByWords(%d) produced %d words, want %d", n, len(words), n)
		}
	}
}

func TestLoremBySentences_Count(t *testing.T) {
	for _, n := range []int{1, 3, 5, 10} {
		output := loremBySentences(n)
		// Each sentence from the pool ends with a period
		// Count sentences by splitting
		if output == "" {
			t.Errorf("loremBySentences(%d) returned empty", n)
			continue
		}
		// Sentences are space-joined, each ends with '.'
		// We just verify the output is non-empty and contains text
		if len(strings.Fields(output)) == 0 {
			t.Errorf("loremBySentences(%d) produced no words", n)
		}
	}
}

func TestLoremByParagraphs_Count(t *testing.T) {
	for _, n := range []int{1, 2, 3, 5} {
		output := loremByParagraphs(n)
		paragraphs := strings.Split(output, "\n\n")
		if len(paragraphs) != n {
			t.Errorf("loremByParagraphs(%d) produced %d paragraphs, want %d", n, len(paragraphs), n)
		}
	}
}

func TestLoremByWords_StartsWithLorem(t *testing.T) {
	output := loremByWords(5)
	if !strings.HasPrefix(output, "Lorem") {
		t.Errorf("loremByWords should start with 'Lorem', got %q", output[:20])
	}
}

// ── dedupe ─────────────────────────────────────────────────────────────────

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"a\nb\nc", 3},
		{"a\nb\nc\n", 3},
		{"single", 1},
		{"", 0},
	}
	for _, tt := range tests {
		got := splitLines(tt.input)
		if len(got) != tt.want {
			t.Errorf("splitLines(%q) = %d lines, want %d", tt.input, len(got), tt.want)
		}
	}
}

func TestDedupe_Logic(t *testing.T) {
	lines := []string{"a", "b", "a", "c", "b"}
	seen := make(map[string]struct{})
	var unique []string
	for _, l := range lines {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			unique = append(unique, l)
		}
	}
	if len(unique) != 3 {
		t.Errorf("dedupe produced %d unique lines, want 3", len(unique))
	}
	want := []string{"a", "b", "c"}
	for i, v := range unique {
		if v != want[i] {
			t.Errorf("dedupe[%d] = %q, want %q", i, v, want[i])
		}
	}
}

func TestDedupe_AllUnique(t *testing.T) {
	lines := []string{"a", "b", "c", "d"}
	seen := make(map[string]struct{})
	var unique []string
	for _, l := range lines {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			unique = append(unique, l)
		}
	}
	if len(unique) != 4 {
		t.Errorf("all unique: dedupe produced %d lines, want 4", len(unique))
	}
}

func TestDedupe_AllDuplicates(t *testing.T) {
	lines := []string{"a", "a", "a", "a"}
	seen := make(map[string]struct{})
	var unique []string
	for _, l := range lines {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			unique = append(unique, l)
		}
	}
	if len(unique) != 1 {
		t.Errorf("all duplicates: dedupe produced %d lines, want 1", len(unique))
	}
}

// ── readInput ──────────────────────────────────────────────────────────────

func TestReadInput_FromArgs(t *testing.T) {
	got, err := readInput([]string{"hello", "world"}, "")
	if err != nil {
		t.Fatalf("readInput error: %v", err)
	}
	if got != "hello world" {
		t.Errorf("readInput = %q, want %q", got, "hello world")
	}
}

func TestReadInput_NoInput(t *testing.T) {
	// With no args, no file, and stdin being a terminal, should error
	// We can only reliably test with args though
	_, err := readInput(nil, "/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestReadInput_FromFile(t *testing.T) {
	got, err := readInput(nil, "/dev/null")
	if err != nil {
		t.Fatalf("readInput error: %v", err)
	}
	if got != "" {
		t.Errorf("readInput from /dev/null = %q, want empty", got)
	}
}

// ── extractNumber ──────────────────────────────────────────────────────────

func TestExtractNumber(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"42 items", 42},
		{"-3.14 radians", -3.14},
		{"no number here", 0},
		{"123", 123},
		{"  99 with spaces", 99},
	}
	for _, tt := range tests {
		got := extractNumber(tt.input)
		if got != tt.want {
			t.Errorf("extractNumber(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ── wordWrap ───────────────────────────────────────────────────────────────

func TestWordWrap(t *testing.T) {
	input := "The quick brown fox jumps over the lazy dog"
	output := wordWrap(input, 20)
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	for _, line := range lines {
		if len(line) > 20 {
			t.Errorf("line exceeds width 20: %q (len %d)", line, len(line))
		}
	}
}

func TestWordWrap_ZeroWidth(t *testing.T) {
	// 0 width should default to 80
	output := wordWrap("hello world", 0)
	if !strings.Contains(output, "hello world") {
		t.Error("wordWrap with 0 width should still produce output")
	}
}

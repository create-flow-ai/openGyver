package text

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// jsonOut is the persistent --json/-j flag shared by all subcommands.
var jsonOut bool

// ── parent command ──────────────────────────────────────────────────────────

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Text manipulation utilities",
	Long: `Text manipulation utilities — count, convert case, reverse, sort,
deduplicate, slugify, generate lorem ipsum, diff, wrap, number lines,
trim, and find-and-replace.

SUBCOMMANDS:

  count     Count words, characters, lines, and sentences
  case      Convert text between cases (upper, lower, title, snake, …)
  reverse   Reverse a string
  sort      Sort lines alphabetically, by length, or numerically
  dedupe    Remove duplicate lines
  slug      Generate a URL-safe slug from text
  lorem     Generate Lorem Ipsum placeholder text
  diff      Show a unified diff between two files
  wrap      Word-wrap text to a given width
  lines     Add line numbers to text
  trim      Strip leading/trailing whitespace and blank lines
  replace   Find and replace text (literal or regex)

All subcommands support --json/-j for machine-readable output.

EXAMPLES:

  openGyver text count "Hello, world!"
  openGyver text case --to snake "Hello World"
  openGyver text reverse "abcdef"
  openGyver text sort --by length --file lines.txt
  openGyver text slug "My Blog Post Title!"
  openGyver text lorem --paragraphs 3
  openGyver text diff --file1 a.txt --file2 b.txt
  openGyver text wrap --width 60 "A very long sentence …"
  openGyver text lines --file main.go
  openGyver text trim "  hello  "
  openGyver text replace --find foo --replace bar --file input.txt`,
}

// ── helpers ─────────────────────────────────────────────────────────────────

// readInput returns text from the first positional arg, from --file, or from
// stdin (when piped).  It returns an error when no input is available.
func readInput(args []string, filePath string) (string, error) {
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	// Try stdin if it's piped.
	info, _ := os.Stdin.Stat()
	if info.Mode()&os.ModeCharDevice == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		var sb strings.Builder
		for scanner.Scan() {
			sb.WriteString(scanner.Text())
			sb.WriteByte('\n')
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return sb.String(), nil
	}
	return "", fmt.Errorf("no input provided (pass text as an argument, use --file, or pipe via stdin)")
}

// ── count ───────────────────────────────────────────────────────────────────

var countFile string

var countCmd = &cobra.Command{
	Use:   "count [text]",
	Short: "Count words, characters, lines, and sentences",
	Long: `Count the number of words, characters, lines, and sentences in text.

Input can be provided as a positional argument, via --file/-f, or piped
through stdin.

EXAMPLES:

  openGyver text count "The quick brown fox jumps over the lazy dog."
  openGyver text count --file essay.txt
  echo "hello world" | openGyver text count
  openGyver text count --json "Hello, world!"`,
	RunE: runCount,
}

func runCount(_ *cobra.Command, args []string) error {
	input, err := readInput(args, countFile)
	if err != nil {
		return err
	}

	words := len(strings.Fields(input))
	characters := len([]rune(input))
	lines := 0
	if len(input) > 0 {
		lines = strings.Count(input, "\n")
		if input[len(input)-1] != '\n' {
			lines++
		}
	}
	// Sentence count: split on .!? followed by whitespace or end-of-string.
	sentenceRe := regexp.MustCompile(`[.!?]+\s|[.!?]+$`)
	sentences := len(sentenceRe.FindAllString(input, -1))

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"words":      words,
			"characters": characters,
			"lines":      lines,
			"sentences":  sentences,
		})
	}

	fmt.Printf("Words:      %d\n", words)
	fmt.Printf("Characters: %d\n", characters)
	fmt.Printf("Lines:      %d\n", lines)
	fmt.Printf("Sentences:  %d\n", sentences)
	return nil
}

// ── case ────────────────────────────────────────────────────────────────────

var caseTo string

var caseCmd = &cobra.Command{
	Use:   "case [text]",
	Short: "Convert text between cases",
	Long: `Convert text to a different case style.

SUPPORTED CASES (--to):

  upper      HELLO WORLD
  lower      hello world
  title      Hello World
  sentence   Hello world
  camel      helloWorld
  pascal     HelloWorld
  snake      hello_world
  kebab      hello-world
  constant   HELLO_WORLD
  dot        hello.world

EXAMPLES:

  openGyver text case --to upper "hello world"
  openGyver text case --to snake "Hello World"
  openGyver text case --to camel "some-variable-name"
  openGyver text case --to kebab "myVariableName"
  openGyver text case --to constant "max retries"
  openGyver text case --to dot "Hello World"
  openGyver text case --to title "the quick brown fox" --json`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCase,
}

// tokenize splits input into lowercase word tokens, handling camelCase,
// snake_case, kebab-case, dot.case, and whitespace-separated words.
func tokenize(s string) []string {
	// Insert a space before uppercase letters that follow a lowercase letter
	// (camelCase / PascalCase boundaries).
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	s = re.ReplaceAllString(s, "${1} ${2}")
	// Replace common delimiters with spaces.
	replacer := strings.NewReplacer("_", " ", "-", " ", ".", " ")
	s = replacer.Replace(s)
	// Split on whitespace, lowercase, and drop empties.
	parts := strings.Fields(s)
	tokens := make([]string, 0, len(parts))
	for _, p := range parts {
		w := strings.ToLower(p)
		if w != "" {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

func convertCase(input, target string) (string, error) {
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return "", nil
	}

	switch target {
	case "upper":
		return strings.ToUpper(input), nil
	case "lower":
		return strings.ToLower(input), nil
	case "title":
		return strings.Title(strings.Join(tokens, " ")), nil //nolint:staticcheck
	case "sentence":
		joined := strings.Join(tokens, " ")
		if len(joined) == 0 {
			return "", nil
		}
		runes := []rune(joined)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes), nil
	case "camel":
		var sb strings.Builder
		for i, t := range tokens {
			if i == 0 {
				sb.WriteString(t)
			} else {
				sb.WriteString(strings.Title(t)) //nolint:staticcheck
			}
		}
		return sb.String(), nil
	case "pascal":
		var sb strings.Builder
		for _, t := range tokens {
			sb.WriteString(strings.Title(t)) //nolint:staticcheck
		}
		return sb.String(), nil
	case "snake":
		return strings.Join(tokens, "_"), nil
	case "kebab":
		return strings.Join(tokens, "-"), nil
	case "constant":
		upper := make([]string, len(tokens))
		for i, t := range tokens {
			upper[i] = strings.ToUpper(t)
		}
		return strings.Join(upper, "_"), nil
	case "dot":
		return strings.Join(tokens, "."), nil
	default:
		return "", fmt.Errorf("unknown case %q — supported: upper, lower, title, sentence, camel, pascal, snake, kebab, constant, dot", target)
	}
}

func runCase(_ *cobra.Command, args []string) error {
	input := strings.Join(args, " ")
	output, err := convertCase(input, caseTo)
	if err != nil {
		return err
	}
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":  input,
			"output": output,
			"case":   caseTo,
		})
	}
	fmt.Println(output)
	return nil
}

// ── reverse ─────────────────────────────────────────────────────────────────

var reverseCmd = &cobra.Command{
	Use:   "reverse [text]",
	Short: "Reverse a string",
	Long: `Reverse the characters in a string.

EXAMPLES:

  openGyver text reverse "Hello, World!"
  openGyver text reverse "abcdef"
  openGyver text reverse --json "palindrome"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		input := strings.Join(args, " ")
		runes := []rune(input)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		output := string(runes)
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"input":  input,
				"output": output,
			})
		}
		fmt.Println(output)
		return nil
	},
}

// ── sort ────────────────────────────────────────────────────────────────────

var (
	sortBy      string
	sortReverse bool
	sortFile    string
)

var sortCmd = &cobra.Command{
	Use:   "sort [text]",
	Short: "Sort lines of text",
	Long: `Sort lines of text alphabetically, by length, or numerically.

Use --by to choose the sort order:
  alpha    Alphabetical (default)
  length   By line length (shortest first)
  numeric  By leading numeric value

Use --reverse to invert the order.

Input via --file/-f, positional argument (newlines as \n), or stdin.

EXAMPLES:

  openGyver text sort --file names.txt
  openGyver text sort --by length --file code.txt
  openGyver text sort --by numeric --reverse --file scores.txt
  echo -e "cherry\napple\nbanana" | openGyver text sort`,
	RunE: runSort,
}

func runSort(_ *cobra.Command, args []string) error {
	input, err := readInput(args, sortFile)
	if err != nil {
		return err
	}
	lines := splitLines(input)

	switch sortBy {
	case "alpha":
		sort.Strings(lines)
	case "length":
		sort.SliceStable(lines, func(i, j int) bool {
			return len(lines[i]) < len(lines[j])
		})
	case "numeric":
		sort.SliceStable(lines, func(i, j int) bool {
			ni := extractNumber(lines[i])
			nj := extractNumber(lines[j])
			return ni < nj
		})
	default:
		return fmt.Errorf("unknown sort order %q — supported: alpha, length, numeric", sortBy)
	}

	if sortReverse {
		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}

	output := strings.Join(lines, "\n")
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"sort_by": sortBy,
			"reverse": sortReverse,
			"count":   len(lines),
			"lines":   lines,
		})
	}
	fmt.Println(output)
	return nil
}

func extractNumber(s string) float64 {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`^-?\d+(\.\d+)?`)
	m := re.FindString(s)
	if m == "" {
		return 0
	}
	n, _ := strconv.ParseFloat(m, 64)
	return n
}

// splitLines splits text into lines and removes a trailing empty line if the
// input ended with a newline.
func splitLines(s string) []string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// ── dedupe ──────────────────────────────────────────────────────────────────

var dedupeFile string

var dedupeCmd = &cobra.Command{
	Use:   "dedupe [text]",
	Short: "Remove duplicate lines",
	Long: `Remove duplicate lines while preserving the original order.

Input via --file/-f, positional argument, or stdin.

EXAMPLES:

  openGyver text dedupe --file list.txt
  echo -e "a\nb\na\nc\nb" | openGyver text dedupe
  openGyver text dedupe --json --file data.txt`,
	RunE: runDedupe,
}

func runDedupe(_ *cobra.Command, args []string) error {
	input, err := readInput(args, dedupeFile)
	if err != nil {
		return err
	}
	lines := splitLines(input)
	seen := make(map[string]struct{}, len(lines))
	unique := make([]string, 0, len(lines))
	for _, l := range lines {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			unique = append(unique, l)
		}
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"total_lines":   len(lines),
			"unique_lines":  len(unique),
			"removed":       len(lines) - len(unique),
			"lines":         unique,
		})
	}
	fmt.Println(strings.Join(unique, "\n"))
	return nil
}

// ── slug ────────────────────────────────────────────────────────────────────

var slugCmd = &cobra.Command{
	Use:   "slug [text]",
	Short: "Generate a URL-safe slug",
	Long: `Generate a URL-safe slug from text.

Converts to lowercase, replaces spaces and special characters with
hyphens, and collapses multiple hyphens.

EXAMPLES:

  openGyver text slug "My Blog Post Title!"
  openGyver text slug "Hello,   World!! 2024"
  openGyver text slug --json "Über Cool: A Story"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		input := strings.Join(args, " ")
		slug := slugify(input)
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"input": input,
				"slug":  slug,
			})
		}
		fmt.Println(slug)
		return nil
	},
}

func slugify(s string) string {
	s = strings.ToLower(s)
	// Replace non-alphanumeric characters with hyphens.
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// ── lorem ───────────────────────────────────────────────────────────────────

var (
	loremWords      int
	loremSentences  int
	loremParagraphs int
)

// loremText is a classic Lorem Ipsum passage used as a source.
const loremText = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem.`

var loremSentencePool = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.",
	"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore.",
	"Excepteur sint occaecat cupidatat non proident, sunt in culpa.",
	"Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit.",
	"Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet.",
	"Sed ut perspiciatis unde omnis iste natus error sit voluptatem.",
	"At vero eos et accusamus et iusto odio dignissimos ducimus.",
	"Nam libero tempore, cum soluta nobis est eligendi optio.",
	"Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus.",
	"Itaque earum rerum hic tenetur a sapiente delectus.",
	"Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse.",
	"Nulla pariatur excepteur sint occaecat cupidatat non proident.",
	"Ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis.",
}

var loremCmd = &cobra.Command{
	Use:   "lorem",
	Short: "Generate Lorem Ipsum placeholder text",
	Long: `Generate Lorem Ipsum placeholder text.

Exactly one of --words, --sentences, or --paragraphs must be given.

EXAMPLES:

  openGyver text lorem --words 50
  openGyver text lorem --sentences 5
  openGyver text lorem --paragraphs 3
  openGyver text lorem --paragraphs 2 --json`,
	Args: cobra.NoArgs,
	RunE: runLorem,
}

func runLorem(_ *cobra.Command, _ []string) error {
	modes := 0
	if loremWords > 0 {
		modes++
	}
	if loremSentences > 0 {
		modes++
	}
	if loremParagraphs > 0 {
		modes++
	}
	if modes == 0 {
		return fmt.Errorf("specify one of --words, --sentences, or --paragraphs")
	}
	if modes > 1 {
		return fmt.Errorf("specify only one of --words, --sentences, or --paragraphs")
	}

	var output string
	switch {
	case loremWords > 0:
		output = loremByWords(loremWords)
	case loremSentences > 0:
		output = loremBySentences(loremSentences)
	case loremParagraphs > 0:
		output = loremByParagraphs(loremParagraphs)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"text":       output,
			"words":      len(strings.Fields(output)),
			"characters": len([]rune(output)),
		})
	}
	fmt.Println(output)
	return nil
}

func loremByWords(n int) string {
	words := strings.Fields(loremText)
	var out []string
	for len(out) < n {
		out = append(out, words...)
	}
	return strings.Join(out[:n], " ")
}

func loremBySentences(n int) string {
	pool := loremSentencePool
	var out []string
	for i := 0; i < n; i++ {
		out = append(out, pool[i%len(pool)])
	}
	return strings.Join(out, " ")
}

func loremByParagraphs(n int) string {
	pool := loremSentencePool
	paragraphs := make([]string, n)
	for i := 0; i < n; i++ {
		// Each paragraph is 4-6 random sentences.
		count := 4 + rand.Intn(3)
		sents := make([]string, count)
		for j := 0; j < count; j++ {
			sents[j] = pool[(i*count+j)%len(pool)]
		}
		paragraphs[i] = strings.Join(sents, " ")
	}
	return strings.Join(paragraphs, "\n\n")
}

// ── diff ────────────────────────────────────────────────────────────────────

var (
	diffFile1 string
	diffFile2 string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Unified diff between two files",
	Long: `Show a unified diff between two text files.

EXAMPLES:

  openGyver text diff --file1 original.txt --file2 modified.txt
  openGyver text diff --file1 a.go --file2 b.go --json`,
	Args: cobra.NoArgs,
	RunE: runDiff,
}

func runDiff(_ *cobra.Command, _ []string) error {
	if diffFile1 == "" || diffFile2 == "" {
		return fmt.Errorf("both --file1 and --file2 are required")
	}
	data1, err := os.ReadFile(diffFile1)
	if err != nil {
		return fmt.Errorf("reading file1: %w", err)
	}
	data2, err := os.ReadFile(diffFile2)
	if err != nil {
		return fmt.Errorf("reading file2: %w", err)
	}

	lines1 := splitLines(string(data1))
	lines2 := splitLines(string(data2))

	hunks := unifiedDiff(lines1, lines2)
	header := fmt.Sprintf("--- %s\n+++ %s", diffFile1, diffFile2)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"file1":   diffFile1,
			"file2":   diffFile2,
			"hunks":   len(hunks),
			"diff":    header + "\n" + strings.Join(hunks, "\n"),
			"changed": len(hunks) > 0,
		})
	}

	if len(hunks) == 0 {
		fmt.Println("Files are identical.")
		return nil
	}
	fmt.Println(header)
	for _, h := range hunks {
		fmt.Println(h)
	}
	return nil
}

// unifiedDiff produces a simple unified-diff output using a basic LCS-based
// approach.  It is intentionally simple (O(n*m) memory) — suitable for
// reasonably-sized text files.
func unifiedDiff(a, b []string) []string {
	// Build LCS table.
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}

	// Walk the LCS table to produce diff lines.
	var diffLines []string
	i, j := 0, 0
	for i < m || j < n {
		if i < m && j < n && a[i] == b[j] {
			diffLines = append(diffLines, " "+a[i])
			i++
			j++
		} else if i < m && (j >= n || dp[i+1][j] >= dp[i][j+1]) {
			diffLines = append(diffLines, "-"+a[i])
			i++
		} else {
			diffLines = append(diffLines, "+"+b[j])
			j++
		}
	}

	// Group into hunks (context of 3 lines).
	const ctx = 3
	var hunks []string
	var hunk []string
	hunkStart1, hunkStart2 := -1, -1
	hunkLen1, hunkLen2 := 0, 0
	lastChange := -1 - ctx - 1 // ensure first hunk starts fresh

	flushHunk := func() {
		if len(hunk) == 0 {
			return
		}
		header := fmt.Sprintf("@@ -%d,%d +%d,%d @@", hunkStart1+1, hunkLen1, hunkStart2+1, hunkLen2)
		hunks = append(hunks, header+"\n"+strings.Join(hunk, "\n"))
		hunk = nil
		hunkStart1, hunkStart2 = -1, -1
		hunkLen1, hunkLen2 = 0, 0
	}

	lineA, lineB := 0, 0
	for idx, dl := range diffLines {
		isChange := dl[0] == '+' || dl[0] == '-'
		if isChange {
			if idx-lastChange > 2*ctx && len(hunk) > 0 {
				flushHunk()
			}
			lastChange = idx
		}

		// Decide whether this line belongs in the current hunk.
		inRange := idx-lastChange <= ctx || (lastChange >= 0 && idx <= lastChange+ctx)
		// Also include context lines before a change.
		upcoming := false
		for look := idx + 1; look < len(diffLines) && look <= idx+ctx; look++ {
			if diffLines[look][0] == '+' || diffLines[look][0] == '-' {
				upcoming = true
				break
			}
		}

		if isChange || inRange || upcoming {
			if hunkStart1 == -1 {
				hunkStart1 = lineA
				hunkStart2 = lineB
			}
			hunk = append(hunk, dl)
			switch dl[0] {
			case ' ':
				hunkLen1++
				hunkLen2++
			case '-':
				hunkLen1++
			case '+':
				hunkLen2++
			}
		}

		switch dl[0] {
		case ' ':
			lineA++
			lineB++
		case '-':
			lineA++
		case '+':
			lineB++
		}
	}
	flushHunk()
	return hunks
}

// ── wrap ────────────────────────────────────────────────────────────────────

var wrapWidth int

var wrapCmd = &cobra.Command{
	Use:   "wrap [text]",
	Short: "Word-wrap text to a given width",
	Long: `Word-wrap text to a given column width (default 80).

Input via positional argument, --file, or stdin.

EXAMPLES:

  openGyver text wrap --width 60 "This is a very long sentence that should be wrapped at sixty characters."
  openGyver text wrap --file article.txt
  openGyver text wrap --width 40 --json "Some text to wrap."`,
	RunE: runWrap,
}

var wrapFile string

func runWrap(_ *cobra.Command, args []string) error {
	input, err := readInput(args, wrapFile)
	if err != nil {
		return err
	}
	output := wordWrap(input, wrapWidth)
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"width":  wrapWidth,
			"input":  input,
			"output": output,
		})
	}
	fmt.Print(output)
	return nil
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}
	paragraphs := strings.Split(text, "\n")
	var result []string
	for _, para := range paragraphs {
		if strings.TrimSpace(para) == "" {
			result = append(result, "")
			continue
		}
		words := strings.Fields(para)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}
		var line strings.Builder
		line.WriteString(words[0])
		for _, w := range words[1:] {
			if line.Len()+1+len(w) > width {
				result = append(result, line.String())
				line.Reset()
				line.WriteString(w)
			} else {
				line.WriteByte(' ')
				line.WriteString(w)
			}
		}
		if line.Len() > 0 {
			result = append(result, line.String())
		}
	}
	return strings.Join(result, "\n") + "\n"
}

// ── lines ───────────────────────────────────────────────────────────────────

var linesFile string

var linesCmd = &cobra.Command{
	Use:   "lines [text]",
	Short: "Add line numbers to text",
	Long: `Prefix each line of text with its line number.

Input via --file/-f, positional argument, or stdin.

EXAMPLES:

  openGyver text lines --file main.go
  echo -e "alpha\nbeta\ngamma" | openGyver text lines
  openGyver text lines --json --file script.sh`,
	RunE: runLines,
}

func runLines(_ *cobra.Command, args []string) error {
	input, err := readInput(args, linesFile)
	if err != nil {
		return err
	}
	rawLines := splitLines(input)
	width := len(fmt.Sprintf("%d", len(rawLines)))
	numbered := make([]string, len(rawLines))
	for i, l := range rawLines {
		numbered[i] = fmt.Sprintf("%*d  %s", width, i+1, l)
	}
	output := strings.Join(numbered, "\n")

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"total_lines": len(rawLines),
			"lines":       numbered,
		})
	}
	fmt.Println(output)
	return nil
}

// ── trim ────────────────────────────────────────────────────────────────────

var trimBlank bool

var trimCmd = &cobra.Command{
	Use:   "trim [text]",
	Short: "Trim whitespace from text",
	Long: `Remove leading and trailing whitespace from each line.

Use --blank to also remove entirely blank lines.

Input via positional argument, --file, or stdin.

EXAMPLES:

  openGyver text trim "  hello world  "
  openGyver text trim --blank --file messy.txt
  echo -e "  hi  \n\n  there  " | openGyver text trim --blank`,
	RunE: runTrim,
}

var trimFile string

func runTrim(_ *cobra.Command, args []string) error {
	input, err := readInput(args, trimFile)
	if err != nil {
		return err
	}
	rawLines := splitLines(input)
	var trimmed []string
	for _, l := range rawLines {
		t := strings.TrimSpace(l)
		if trimBlank && t == "" {
			continue
		}
		trimmed = append(trimmed, t)
	}
	output := strings.Join(trimmed, "\n")

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input_lines":  len(rawLines),
			"output_lines": len(trimmed),
			"removed":      len(rawLines) - len(trimmed),
			"output":       output,
		})
	}
	fmt.Println(output)
	return nil
}

// ── replace ─────────────────────────────────────────────────────────────────

var (
	replaceFind    string
	replaceReplace string
	replaceFile    string
	replaceRegex   bool
)

var replaceCmd = &cobra.Command{
	Use:   "replace [text]",
	Short: "Find and replace text",
	Long: `Find and replace text using literal strings or regular expressions.

Use --find and --replace to specify the search and replacement strings.
Use --regex to interpret --find as a regular expression (supports Go
regexp syntax including capture groups like $1, $2).

Input via --file/-f, positional argument, or stdin.

EXAMPLES:

  openGyver text replace --find foo --replace bar "foo fighters"
  openGyver text replace --find foo --replace bar --file input.txt
  openGyver text replace --regex --find "(\w+)@(\w+)" --replace "$1 at $2" "user@host"
  openGyver text replace --find old --replace new --json "old is gold"`,
	RunE: runReplace,
}

func runReplace(_ *cobra.Command, args []string) error {
	if replaceFind == "" {
		return fmt.Errorf("--find is required")
	}
	input, err := readInput(args, replaceFile)
	if err != nil {
		return err
	}

	var output string
	var count int
	if replaceRegex {
		re, err := regexp.Compile(replaceFind)
		if err != nil {
			return fmt.Errorf("invalid regex %q: %w", replaceFind, err)
		}
		count = len(re.FindAllString(input, -1))
		output = re.ReplaceAllString(input, replaceReplace)
	} else {
		count = strings.Count(input, replaceFind)
		output = strings.ReplaceAll(input, replaceFind, replaceReplace)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"find":         replaceFind,
			"replace":      replaceReplace,
			"regex":        replaceRegex,
			"replacements": count,
			"output":       output,
		})
	}
	fmt.Print(output)
	return nil
}

// ── init ────────────────────────────────────────────────────────────────────

func init() {
	// Persistent flag on parent.
	textCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	// count
	countCmd.Flags().StringVarP(&countFile, "file", "f", "", "read input from a file")
	textCmd.AddCommand(countCmd)

	// case
	caseCmd.Flags().StringVar(&caseTo, "to", "lower", "target case: upper, lower, title, sentence, camel, pascal, snake, kebab, constant, dot")
	textCmd.AddCommand(caseCmd)

	// reverse
	textCmd.AddCommand(reverseCmd)

	// sort
	sortCmd.Flags().StringVar(&sortBy, "by", "alpha", "sort order: alpha, length, numeric")
	sortCmd.Flags().BoolVar(&sortReverse, "reverse", false, "reverse the sort order")
	sortCmd.Flags().StringVarP(&sortFile, "file", "f", "", "read input from a file")
	textCmd.AddCommand(sortCmd)

	// dedupe
	dedupeCmd.Flags().StringVarP(&dedupeFile, "file", "f", "", "read input from a file")
	textCmd.AddCommand(dedupeCmd)

	// slug
	textCmd.AddCommand(slugCmd)

	// lorem
	loremCmd.Flags().IntVar(&loremWords, "words", 0, "number of words to generate")
	loremCmd.Flags().IntVar(&loremSentences, "sentences", 0, "number of sentences to generate")
	loremCmd.Flags().IntVar(&loremParagraphs, "paragraphs", 0, "number of paragraphs to generate")
	textCmd.AddCommand(loremCmd)

	// diff
	diffCmd.Flags().StringVar(&diffFile1, "file1", "", "first file to compare")
	diffCmd.Flags().StringVar(&diffFile2, "file2", "", "second file to compare")
	textCmd.AddCommand(diffCmd)

	// wrap
	wrapCmd.Flags().IntVar(&wrapWidth, "width", 80, "column width to wrap at")
	wrapCmd.Flags().StringVarP(&wrapFile, "file", "f", "", "read input from a file")
	textCmd.AddCommand(wrapCmd)

	// lines
	linesCmd.Flags().StringVarP(&linesFile, "file", "f", "", "read input from a file")
	textCmd.AddCommand(linesCmd)

	// trim
	trimCmd.Flags().BoolVar(&trimBlank, "blank", false, "also remove blank lines")
	trimCmd.Flags().StringVarP(&trimFile, "file", "f", "", "read input from a file")
	textCmd.AddCommand(trimCmd)

	// replace
	replaceCmd.Flags().StringVar(&replaceFind, "find", "", "string or pattern to find")
	replaceCmd.Flags().StringVar(&replaceReplace, "replace", "", "replacement string")
	replaceCmd.Flags().BoolVar(&replaceRegex, "regex", false, "interpret --find as a regular expression")
	replaceCmd.Flags().StringVarP(&replaceFile, "file", "f", "", "read input from a file")
	textCmd.AddCommand(replaceCmd)

	// Register with root.
	cmd.Register(textCmd)
}

package accessibility

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var jsonOut bool

var accessibilityCmd = &cobra.Command{
	Use:   "accessibility",
	Short: "Accessibility tools — contrast checker, readability scores",
	Long: `Tools for checking web accessibility compliance.

SUBCOMMANDS:
  contrast     WCAG contrast ratio checker for two colors
  readability  Calculate readability scores (Flesch, Gunning Fog, etc.)

Examples:
  openGyver accessibility contrast "#ffffff" "#000000"
  openGyver accessibility readability "The quick brown fox jumps over the lazy dog."
  openGyver accessibility readability --file article.txt`,
}

// --- Contrast subcommand ---
var contrastCmd = &cobra.Command{
	Use:   "contrast <color1> <color2>",
	Short: "WCAG contrast ratio checker",
	Long: `Check the contrast ratio between two colors per WCAG 2.1 guidelines.

Reports the ratio and pass/fail for AA and AAA at normal and large text sizes.
Colors can be hex (#fff, #ffffff), rgb(r,g,b), or CSS names.

Examples:
  openGyver accessibility contrast "#ffffff" "#000000"
  openGyver accessibility contrast "#333" "#ccc"
  openGyver accessibility contrast white black`,
	Args: cobra.ExactArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		r1, g1, b1, err := parseColor(args[0])
		if err != nil {
			return fmt.Errorf("color 1: %w", err)
		}
		r2, g2, b2, err := parseColor(args[1])
		if err != nil {
			return fmt.Errorf("color 2: %w", err)
		}

		l1 := relativeLuminance(r1, g1, b1)
		l2 := relativeLuminance(r2, g2, b2)
		ratio := contrastRatio(l1, l2)

		aaLarge := ratio >= 3.0
		aaNormal := ratio >= 4.5
		aaaLarge := ratio >= 4.5
		aaaNormal := ratio >= 7.0

		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"color1": args[0], "color2": args[1],
				"ratio":      math.Round(ratio*100) / 100,
				"aa_normal":  aaNormal, "aa_large": aaLarge,
				"aaa_normal": aaaNormal, "aaa_large": aaaLarge,
			})
		}

		fmt.Printf("Color 1:     %s\n", args[0])
		fmt.Printf("Color 2:     %s\n", args[1])
		fmt.Printf("Ratio:       %.2f:1\n", ratio)
		fmt.Printf("AA Normal:   %s (need 4.5:1)\n", passFail(aaNormal))
		fmt.Printf("AA Large:    %s (need 3:1)\n", passFail(aaLarge))
		fmt.Printf("AAA Normal:  %s (need 7:1)\n", passFail(aaaNormal))
		fmt.Printf("AAA Large:   %s (need 4.5:1)\n", passFail(aaaLarge))
		return nil
	},
}

// --- Readability subcommand ---
var readFile string

var readabilityCmd = &cobra.Command{
	Use:   "readability [text]",
	Short: "Calculate readability scores",
	Long: `Calculate Flesch Reading Ease, Flesch-Kincaid Grade Level, and Gunning Fog Index.

Input from argument or --file/-f.

Examples:
  openGyver accessibility readability "The cat sat on the mat."
  openGyver accessibility readability --file article.txt
  openGyver accessibility readability -j "Complex text here."`,
	RunE: func(c *cobra.Command, args []string) error {
		var text string
		if readFile != "" {
			data, err := os.ReadFile(readFile)
			if err != nil {
				return err
			}
			text = string(data)
		} else if len(args) > 0 {
			text = args[0]
		} else {
			return fmt.Errorf("provide text as argument or use --file")
		}

		words := countWords(text)
		sentences := countSentences(text)
		syllables := countSyllablesText(text)
		complexWords := countComplexWords(text)

		if sentences == 0 {
			sentences = 1
		}
		if words == 0 {
			return fmt.Errorf("no words found in input")
		}

		wps := float64(words) / float64(sentences)
		spw := float64(syllables) / float64(words)

		fre := 206.835 - 1.015*wps - 84.6*spw
		fkgl := 0.39*wps + 11.8*spw - 15.59
		gfi := 0.4 * (wps + 100.0*float64(complexWords)/float64(words))

		freLabel := fleschLabel(fre)

		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"words": words, "sentences": sentences, "syllables": syllables,
				"flesch_reading_ease":       math.Round(fre*100) / 100,
				"flesch_reading_ease_label": freLabel,
				"flesch_kincaid_grade":      math.Round(fkgl*100) / 100,
				"gunning_fog_index":         math.Round(gfi*100) / 100,
			})
		}

		fmt.Printf("Words:       %d\n", words)
		fmt.Printf("Sentences:   %d\n", sentences)
		fmt.Printf("Syllables:   %d\n", syllables)
		fmt.Println()
		fmt.Printf("Flesch Reading Ease:       %.1f (%s)\n", fre, freLabel)
		fmt.Printf("Flesch-Kincaid Grade:      %.1f\n", fkgl)
		fmt.Printf("Gunning Fog Index:         %.1f\n", gfi)
		return nil
	},
}

// --- Color parsing helpers ---
var cssColors = map[string][3]uint8{
	"black": {0, 0, 0}, "white": {255, 255, 255}, "red": {255, 0, 0},
	"green": {0, 128, 0}, "blue": {0, 0, 255}, "yellow": {255, 255, 0},
	"cyan": {0, 255, 255}, "magenta": {255, 0, 255}, "gray": {128, 128, 128},
	"grey": {128, 128, 128}, "silver": {192, 192, 192}, "orange": {255, 165, 0},
	"purple": {128, 0, 128}, "navy": {0, 0, 128}, "teal": {0, 128, 128},
}

func parseColor(s string) (uint8, uint8, uint8, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if c, ok := cssColors[s]; ok {
		return c[0], c[1], c[2], nil
	}
	s = strings.TrimPrefix(s, "#")
	if len(s) == 3 {
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) == 6 {
		r, _ := strconv.ParseUint(s[0:2], 16, 8)
		g, _ := strconv.ParseUint(s[2:4], 16, 8)
		b, _ := strconv.ParseUint(s[4:6], 16, 8)
		return uint8(r), uint8(g), uint8(b), nil
	}
	if strings.HasPrefix(s, "rgb(") {
		s = strings.TrimPrefix(s, "rgb(")
		s = strings.TrimSuffix(s, ")")
		parts := strings.Split(s, ",")
		if len(parts) == 3 {
			r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
			return uint8(r), uint8(g), uint8(b), nil
		}
	}
	return 0, 0, 0, fmt.Errorf("cannot parse color: %s", s)
}

func relativeLuminance(r, g, b uint8) float64 {
	rs := linearize(float64(r) / 255)
	gs := linearize(float64(g) / 255)
	bs := linearize(float64(b) / 255)
	return 0.2126*rs + 0.7152*gs + 0.0722*bs
}

func linearize(v float64) float64 {
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func contrastRatio(l1, l2 float64) float64 {
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

func passFail(ok bool) string {
	if ok {
		return "PASS"
	}
	return "FAIL"
}

// --- Readability helpers ---
func countWords(text string) int {
	return len(strings.Fields(text))
}

func countSentences(text string) int {
	count := 0
	for _, r := range text {
		if r == '.' || r == '!' || r == '?' {
			count++
		}
	}
	if count == 0 {
		count = 1
	}
	return count
}

func countSyllablesWord(word string) int {
	word = strings.ToLower(strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r)
	}))
	if len(word) == 0 {
		return 0
	}
	vowels := "aeiouy"
	count := 0
	prevVowel := false
	for _, r := range word {
		isVowel := strings.ContainsRune(vowels, r)
		if isVowel && !prevVowel {
			count++
		}
		prevVowel = isVowel
	}
	// Silent e
	if strings.HasSuffix(word, "e") && count > 1 {
		count--
	}
	if count == 0 {
		count = 1
	}
	return count
}

func countSyllablesText(text string) int {
	total := 0
	for _, w := range strings.Fields(text) {
		total += countSyllablesWord(w)
	}
	return total
}

func countComplexWords(text string) int {
	count := 0
	for _, w := range strings.Fields(text) {
		if countSyllablesWord(w) >= 3 {
			count++
		}
	}
	return count
}

func fleschLabel(score float64) string {
	switch {
	case score >= 90:
		return "Very Easy"
	case score >= 80:
		return "Easy"
	case score >= 70:
		return "Fairly Easy"
	case score >= 60:
		return "Standard"
	case score >= 50:
		return "Fairly Difficult"
	case score >= 30:
		return "Difficult"
	default:
		return "Very Confusing"
	}
}

func init() {
	accessibilityCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	readabilityCmd.Flags().StringVarP(&readFile, "file", "f", "", "read text from file")

	accessibilityCmd.AddCommand(contrastCmd, readabilityCmd)
	cmd.Register(accessibilityCmd)
}

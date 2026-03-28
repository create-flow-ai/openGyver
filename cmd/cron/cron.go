package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── flags ──────────────────────────────────────────────────────────────────

var (
	jsonOut   bool
	nextCount int
)

// ── parent command ─────────────────────────────────────────────────────────

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Cron expression tools",
	Long: `Parse, explain, and validate cron expressions.

SUBCOMMANDS:

  explain    Parse a cron expression and describe each field
  next       Show the next N run times for a cron expression
  validate   Check if a cron expression is valid

Supports standard 5-field expressions (minute hour day month weekday)
and 6-field expressions with seconds (second minute hour day month weekday).

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver cron explain "*/5 * * * *"
  openGyver cron next "0 9 * * 1-5" --count 10
  openGyver cron validate "0 0 1 1 *"`,
}

// ── explain subcommand ─────────────────────────────────────────────────────

var explainCmd = &cobra.Command{
	Use:   "explain <expression>",
	Short: "Parse a cron expression and output human-readable description",
	Args:  cobra.ExactArgs(1),
	RunE:  runExplain,
}

func runExplain(_ *cobra.Command, args []string) error {
	expr := args[0]
	fields, hasSeconds, err := parseExpression(expr)
	if err != nil {
		return err
	}

	desc := describeFields(fields, hasSeconds)

	if jsonOut {
		result := map[string]interface{}{
			"expression":  expr,
			"has_seconds": hasSeconds,
			"fields":      desc,
			"summary":     buildSummary(desc, hasSeconds),
		}
		return cmd.PrintJSON(result)
	}

	fmt.Printf("Expression: %s\n", expr)
	if hasSeconds {
		fmt.Println("Format:     6-field (with seconds)")
	} else {
		fmt.Println("Format:     5-field (standard)")
	}
	fmt.Println()

	labels := fieldLabels(hasSeconds)
	for i, d := range desc {
		fmt.Printf("  %-10s %s\n", labels[i]+":", d)
	}
	return nil
}

// ── next subcommand ────────────────────────────────────────────────────────

var nextCmd = &cobra.Command{
	Use:   "next <expression>",
	Short: "Show next N run times for a cron expression",
	Args:  cobra.ExactArgs(1),
	RunE:  runNext,
}

func runNext(_ *cobra.Command, args []string) error {
	expr := args[0]
	fields, hasSeconds, err := parseExpression(expr)
	if err != nil {
		return err
	}

	parsed, err := compileFields(fields, hasSeconds)
	if err != nil {
		return err
	}

	times := nextOccurrences(parsed, hasSeconds, time.Now(), nextCount)

	if jsonOut {
		strs := make([]string, len(times))
		for i, t := range times {
			strs[i] = t.Format(time.RFC3339)
		}
		return cmd.PrintJSON(map[string]interface{}{
			"expression": expr,
			"count":      len(times),
			"next":       strs,
		})
	}

	for i, t := range times {
		fmt.Printf("  %d. %s\n", i+1, t.Format("Mon 2006-01-02 15:04:05"))
	}
	return nil
}

// ── validate subcommand ────────────────────────────────────────────────────

var validateCmd = &cobra.Command{
	Use:   "validate <expression>",
	Short: "Check if a cron expression is valid",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func runValidate(_ *cobra.Command, args []string) error {
	expr := args[0]
	_, _, err := parseExpression(expr)
	valid := err == nil

	if jsonOut {
		result := map[string]interface{}{
			"expression": expr,
			"valid":      valid,
		}
		if err != nil {
			result["error"] = err.Error()
		}
		return cmd.PrintJSON(result)
	}

	if valid {
		fmt.Printf("Valid cron expression: %s\n", expr)
	} else {
		fmt.Printf("Invalid cron expression: %s\n  Error: %s\n", expr, err)
	}
	return nil
}

// ── init ───────────────────────────────────────────────────────────────────

func init() {
	cronCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	nextCmd.Flags().IntVar(&nextCount, "count", 5, "number of next run times to show")
	cronCmd.AddCommand(explainCmd)
	cronCmd.AddCommand(nextCmd)
	cronCmd.AddCommand(validateCmd)
	cmd.Register(cronCmd)
}

// ── cron parsing ───────────────────────────────────────────────────────────

// parseExpression splits the expression into fields and detects 5 vs 6 field format.
func parseExpression(expr string) ([]string, bool, error) {
	parts := strings.Fields(strings.TrimSpace(expr))
	switch len(parts) {
	case 5:
		return parts, false, nil
	case 6:
		return parts, true, nil
	default:
		return nil, false, fmt.Errorf("expected 5 or 6 fields, got %d", len(parts))
	}
}

// fieldRange represents the valid range for each cron field.
type fieldRange struct {
	min, max int
	name     string
}

func standardRanges() []fieldRange {
	return []fieldRange{
		{0, 59, "minute"},
		{0, 23, "hour"},
		{1, 31, "day"},
		{1, 12, "month"},
		{0, 6, "weekday"},
	}
}

func sixFieldRanges() []fieldRange {
	return []fieldRange{
		{0, 59, "second"},
		{0, 59, "minute"},
		{0, 23, "hour"},
		{1, 31, "day"},
		{1, 12, "month"},
		{0, 6, "weekday"},
	}
}

func fieldLabels(hasSeconds bool) []string {
	if hasSeconds {
		return []string{"Second", "Minute", "Hour", "Day", "Month", "Weekday"}
	}
	return []string{"Minute", "Hour", "Day", "Month", "Weekday"}
}

// parsedField holds the expanded set of allowed values for one field.
type parsedField struct {
	values map[int]bool // nil means wildcard (all values)
}

// compileFields turns raw field strings into parsedField slices for scheduling.
func compileFields(fields []string, hasSeconds bool) ([]parsedField, error) {
	ranges := standardRanges()
	if hasSeconds {
		ranges = sixFieldRanges()
	}

	result := make([]parsedField, len(fields))
	for i, raw := range fields {
		vals, err := expandField(raw, ranges[i])
		if err != nil {
			return nil, fmt.Errorf("field %d (%s): %w", i+1, ranges[i].name, err)
		}
		result[i] = parsedField{values: vals}
	}
	return result, nil
}

// expandField parses a single cron field token (e.g. "*/5", "1-3", "1,15", "*", "L")
// and returns the set of matching values, or nil for wildcard.
func expandField(raw string, fr fieldRange) (map[int]bool, error) {
	if raw == "*" {
		return nil, nil // wildcard
	}

	values := map[int]bool{}

	// Handle comma-separated list.
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)

		// Handle L (last) — for day-of-month or weekday.
		if strings.ToUpper(part) == "L" {
			values[fr.max] = true
			continue
		}

		// Handle step: */N or N-M/N
		if strings.Contains(part, "/") {
			stepParts := strings.SplitN(part, "/", 2)
			step, err := strconv.Atoi(stepParts[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step %q", part)
			}
			start, end := fr.min, fr.max
			if stepParts[0] != "*" {
				rangeParts := strings.SplitN(stepParts[0], "-", 2)
				start, err = strconv.Atoi(rangeParts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid value %q", part)
				}
				if len(rangeParts) == 2 {
					end, err = strconv.Atoi(rangeParts[1])
					if err != nil {
						return nil, fmt.Errorf("invalid range %q", part)
					}
				}
			}
			for v := start; v <= end; v += step {
				values[v] = true
			}
			continue
		}

		// Handle range: N-M
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			lo, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range start %q", part)
			}
			hi, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end %q", part)
			}
			if lo > hi {
				return nil, fmt.Errorf("invalid range %d-%d", lo, hi)
			}
			for v := lo; v <= hi; v++ {
				values[v] = true
			}
			continue
		}

		// Plain number.
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q", part)
		}
		if n < fr.min || n > fr.max {
			return nil, fmt.Errorf("value %d out of range [%d-%d]", n, fr.min, fr.max)
		}
		values[n] = true
	}
	return values, nil
}

// matches checks if a value is accepted by a parsed field.
func matches(pf parsedField, val int) bool {
	if pf.values == nil {
		return true // wildcard
	}
	return pf.values[val]
}

// nextOccurrences calculates the next N times the cron expression fires.
func nextOccurrences(parsed []parsedField, hasSeconds bool, from time.Time, count int) []time.Time {
	var results []time.Time
	var t time.Time
	if hasSeconds {
		t = from.Truncate(time.Second).Add(time.Second)
	} else {
		t = from.Truncate(time.Minute).Add(time.Minute)
	}

	limit := 525960 // safety limit: one year of minutes
	for len(results) < count && limit > 0 {
		limit--
		if matchesTime(parsed, hasSeconds, t) {
			results = append(results, t)
		}
		if hasSeconds {
			t = t.Add(time.Second)
		} else {
			t = t.Add(time.Minute)
		}
	}
	return results
}

func matchesTime(parsed []parsedField, hasSeconds bool, t time.Time) bool {
	if hasSeconds {
		return matches(parsed[0], t.Second()) &&
			matches(parsed[1], t.Minute()) &&
			matches(parsed[2], t.Hour()) &&
			matches(parsed[3], t.Day()) &&
			matches(parsed[4], int(t.Month())) &&
			matches(parsed[5], int(t.Weekday()))
	}
	return matches(parsed[0], t.Minute()) &&
		matches(parsed[1], t.Hour()) &&
		matches(parsed[2], t.Day()) &&
		matches(parsed[3], int(t.Month())) &&
		matches(parsed[4], int(t.Weekday()))
}

// ── description helpers ────────────────────────────────────────────────────

var monthNames = []string{"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December"}

var weekdayNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func describeFields(fields []string, hasSeconds bool) []string {
	ranges := standardRanges()
	if hasSeconds {
		ranges = sixFieldRanges()
	}
	descs := make([]string, len(fields))
	for i, f := range fields {
		descs[i] = describeField(f, ranges[i])
	}
	return descs
}

func describeField(raw string, fr fieldRange) string {
	if raw == "*" {
		return fmt.Sprintf("every %s", fr.name)
	}

	if strings.HasPrefix(raw, "*/") {
		step := raw[2:]
		return fmt.Sprintf("every %s %ss", step, fr.name)
	}

	if strings.ToUpper(raw) == "L" {
		return fmt.Sprintf("last %s", fr.name)
	}

	// Translate specific values for month/weekday.
	if fr.name == "month" {
		return "month(s) " + translateList(raw, func(n int) string {
			if n >= 1 && n <= 12 {
				return monthNames[n]
			}
			return strconv.Itoa(n)
		})
	}
	if fr.name == "weekday" {
		return translateList(raw, func(n int) string {
			if n >= 0 && n <= 6 {
				return weekdayNames[n]
			}
			return strconv.Itoa(n)
		})
	}

	return raw
}

func translateList(raw string, nameFn func(int) string) string {
	parts := strings.Split(raw, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.Contains(p, "-") {
			rangeParts := strings.SplitN(p, "-", 2)
			lo, err1 := strconv.Atoi(rangeParts[0])
			hi, err2 := strconv.Atoi(rangeParts[1])
			if err1 == nil && err2 == nil {
				out = append(out, nameFn(lo)+" through "+nameFn(hi))
			} else {
				out = append(out, p)
			}
		} else {
			n, err := strconv.Atoi(p)
			if err == nil {
				out = append(out, nameFn(n))
			} else {
				out = append(out, p)
			}
		}
	}
	return strings.Join(out, ", ")
}

func buildSummary(descs []string, hasSeconds bool) string {
	return strings.Join(descs, " | ")
}

package mathcmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── flags ──────────────────────────────────────────────────────────────────

var (
	jsonOut bool

	// percent flags
	percentOf     bool
	percentIs     bool
	percentChange bool
	percentValue  float64
	percentTotal  float64
	percentFrom   float64
	percentTo     float64
)

// ── parent command ─────────────────────────────────────────────────────────

var mathCmd = &cobra.Command{
	Use:   "math",
	Short: "Math expression evaluator and utilities",
	Long: `Evaluate math expressions and perform common calculations.

SUBCOMMANDS:

  eval        Evaluate a math expression
  percent     Percentage calculator
  gcd         Greatest common divisor of two numbers
  lcm         Least common multiple of two numbers
  factorial   Factorial of N
  fibonacci   Nth Fibonacci number

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver math eval "2 + 3 * 4"
  openGyver math eval "sqrt(144) + 2^3"
  openGyver math percent --of --value 15 --total 200
  openGyver math gcd 12 18
  openGyver math lcm 4 6
  openGyver math factorial 10
  openGyver math fibonacci 20`,
}

// ── eval subcommand ────────────────────────────────────────────────────────

var evalCmd = &cobra.Command{
	Use:   "eval <expression>",
	Short: "Evaluate a math expression",
	Long: `Evaluate a mathematical expression.

Supported operators: + - * / % ^ (power)
Supported functions: sqrt, abs, ceil, floor, sin, cos, tan, log, log2, log10
Supported constants: pi, e
Parentheses are supported for grouping.

EXAMPLES:

  openGyver math eval "2 + 3 * 4"
  openGyver math eval "(2 + 3) * 4"
  openGyver math eval "sqrt(144)"
  openGyver math eval "2^10"
  openGyver math eval "sin(pi/2)"`,
	Args: cobra.ExactArgs(1),
	RunE: runEval,
}

func runEval(_ *cobra.Command, args []string) error {
	expr := args[0]
	result, err := evaluate(expr)
	if err != nil {
		return fmt.Errorf("evaluation error: %w", err)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"expression": expr,
			"result":     result,
		})
	}

	fmt.Println(formatNumber(result))
	return nil
}

// ── percent subcommand ─────────────────────────────────────────────────────

var percentCmd = &cobra.Command{
	Use:   "percent",
	Short: "Percentage calculator",
	Long: `Calculate percentages in three modes:

  --of      What is X% of Y?          (requires --value and --total)
  --is      X is what % of Y?         (requires --value and --total)
  --change  % change from X to Y      (requires --from and --to)

EXAMPLES:

  openGyver math percent --of --value 15 --total 200
  openGyver math percent --is --value 30 --total 200
  openGyver math percent --change --from 80 --to 100`,
	RunE: runPercent,
}

func runPercent(_ *cobra.Command, args []string) error {
	switch {
	case percentOf:
		result := (percentValue / 100) * percentTotal
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"mode":   "of",
				"value":  percentValue,
				"total":  percentTotal,
				"result": result,
			})
		}
		fmt.Printf("%.4g%% of %.4g = %.4g\n", percentValue, percentTotal, result)

	case percentIs:
		if percentTotal == 0 {
			return fmt.Errorf("total cannot be zero")
		}
		result := (percentValue / percentTotal) * 100
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"mode":   "is",
				"value":  percentValue,
				"total":  percentTotal,
				"result": result,
			})
		}
		fmt.Printf("%.4g is %.4g%% of %.4g\n", percentValue, result, percentTotal)

	case percentChange:
		if percentFrom == 0 {
			return fmt.Errorf("from value cannot be zero")
		}
		result := ((percentTo - percentFrom) / percentFrom) * 100
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"mode":   "change",
				"from":   percentFrom,
				"to":     percentTo,
				"result": result,
			})
		}
		fmt.Printf("Change from %.4g to %.4g = %.4g%%\n", percentFrom, percentTo, result)

	default:
		return fmt.Errorf("specify one of --of, --is, or --change")
	}
	return nil
}

// ── gcd subcommand ─────────────────────────────────────────────────────────

var gcdCmd = &cobra.Command{
	Use:   "gcd <a> <b>",
	Short: "Greatest common divisor of two numbers",
	Args:  cobra.ExactArgs(2),
	RunE:  runGCD,
}

func runGCD(_ *cobra.Command, args []string) error {
	a, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", args[0], err)
	}
	b, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", args[1], err)
	}

	result := gcd(abs64(a), abs64(b))

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"a":      a,
			"b":      b,
			"result": result,
		})
	}
	fmt.Println(result)
	return nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// ── lcm subcommand ─────────────────────────────────────────────────────────

var lcmCmd = &cobra.Command{
	Use:   "lcm <a> <b>",
	Short: "Least common multiple of two numbers",
	Args:  cobra.ExactArgs(2),
	RunE:  runLCM,
}

func runLCM(_ *cobra.Command, args []string) error {
	a, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", args[0], err)
	}
	b, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", args[1], err)
	}

	result := lcm(abs64(a), abs64(b))

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"a":      a,
			"b":      b,
			"result": result,
		})
	}
	fmt.Println(result)
	return nil
}

func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return (a / gcd(a, b)) * b
}

// ── factorial subcommand ───────────────────────────────────────────────────

var factorialCmd = &cobra.Command{
	Use:   "factorial <n>",
	Short: "Factorial of N",
	Args:  cobra.ExactArgs(1),
	RunE:  runFactorial,
}

func runFactorial(_ *cobra.Command, args []string) error {
	n, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", args[0], err)
	}
	if n < 0 {
		return fmt.Errorf("factorial is not defined for negative numbers")
	}
	if n > 170 {
		return fmt.Errorf("factorial overflow: n must be <= 170 for float64")
	}

	result := factorial(n)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"n":      n,
			"result": result,
		})
	}
	fmt.Println(formatNumber(result))
	return nil
}

func factorial(n int) float64 {
	if n <= 1 {
		return 1
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result
}

// ── fibonacci subcommand ───────────────────────────────────────────────────

var fibonacciCmd = &cobra.Command{
	Use:   "fibonacci <n>",
	Short: "Nth Fibonacci number",
	Args:  cobra.ExactArgs(1),
	RunE:  runFibonacci,
}

func runFibonacci(_ *cobra.Command, args []string) error {
	n, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", args[0], err)
	}
	if n < 0 {
		return fmt.Errorf("fibonacci index must be non-negative")
	}

	result := fibonacci(n)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"n":      n,
			"result": result,
		})
	}
	fmt.Println(result)
	return nil
}

func fibonacci(n int) int64 {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	var a, b int64 = 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// ── init ───────────────────────────────────────────────────────────────────

func init() {
	mathCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	percentCmd.Flags().BoolVar(&percentOf, "of", false, "what is X% of Y")
	percentCmd.Flags().BoolVar(&percentIs, "is", false, "X is what % of Y")
	percentCmd.Flags().BoolVar(&percentChange, "change", false, "% change from X to Y")
	percentCmd.Flags().Float64Var(&percentValue, "value", 0, "the percentage or value")
	percentCmd.Flags().Float64Var(&percentTotal, "total", 0, "the total amount")
	percentCmd.Flags().Float64Var(&percentFrom, "from", 0, "starting value for % change")
	percentCmd.Flags().Float64Var(&percentTo, "to", 0, "ending value for % change")

	mathCmd.AddCommand(evalCmd)
	mathCmd.AddCommand(percentCmd)
	mathCmd.AddCommand(gcdCmd)
	mathCmd.AddCommand(lcmCmd)
	mathCmd.AddCommand(factorialCmd)
	mathCmd.AddCommand(fibonacciCmd)
	cmd.Register(mathCmd)
}

// ── expression evaluator (recursive descent parser) ────────────────────────

// evaluate parses and evaluates a mathematical expression.
func evaluate(expr string) (float64, error) {
	p := &parser{input: expr}
	result, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	p.skipSpaces()
	if p.pos < len(p.input) {
		return 0, fmt.Errorf("unexpected character at position %d: %q", p.pos, string(p.input[p.pos]))
	}
	return result, nil
}

type parser struct {
	input string
	pos   int
}

func (p *parser) skipSpaces() {
	for p.pos < len(p.input) && p.input[p.pos] == ' ' {
		p.pos++
	}
}

func (p *parser) peek() byte {
	p.skipSpaces()
	if p.pos < len(p.input) {
		return p.input[p.pos]
	}
	return 0
}

func (p *parser) consume(ch byte) bool {
	p.skipSpaces()
	if p.pos < len(p.input) && p.input[p.pos] == ch {
		p.pos++
		return true
	}
	return false
}

// parseExpr handles + and - (lowest precedence).
func (p *parser) parseExpr() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpaces()
		if p.pos >= len(p.input) {
			break
		}
		op := p.input[p.pos]
		if op != '+' && op != '-' {
			break
		}
		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return 0, err
		}
		if op == '+' {
			left += right
		} else {
			left -= right
		}
	}
	return left, nil
}

// parseTerm handles *, /, % (medium precedence).
func (p *parser) parseTerm() (float64, error) {
	left, err := p.parsePower()
	if err != nil {
		return 0, err
	}
	for {
		p.skipSpaces()
		if p.pos >= len(p.input) {
			break
		}
		op := p.input[p.pos]
		if op != '*' && op != '/' && op != '%' {
			break
		}
		p.pos++
		right, err := p.parsePower()
		if err != nil {
			return 0, err
		}
		switch op {
		case '*':
			left *= right
		case '/':
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		case '%':
			if right == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			left = math.Mod(left, right)
		}
	}
	return left, nil
}

// parsePower handles ^ (right-associative, high precedence).
func (p *parser) parsePower() (float64, error) {
	base, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	p.skipSpaces()
	if p.pos < len(p.input) && p.input[p.pos] == '^' {
		p.pos++
		exp, err := p.parsePower() // right-associative
		if err != nil {
			return 0, err
		}
		return math.Pow(base, exp), nil
	}
	return base, nil
}

// parseUnary handles unary minus.
func (p *parser) parseUnary() (float64, error) {
	p.skipSpaces()
	if p.pos < len(p.input) && p.input[p.pos] == '-' {
		p.pos++
		val, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}
	return p.parsePrimary()
}

// parsePrimary handles numbers, constants, functions, and parentheses.
func (p *parser) parsePrimary() (float64, error) {
	p.skipSpaces()

	// Parenthesized expression.
	if p.consume('(') {
		val, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		if !p.consume(')') {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		return val, nil
	}

	// Number.
	if p.pos < len(p.input) && (p.input[p.pos] >= '0' && p.input[p.pos] <= '9' || p.input[p.pos] == '.') {
		return p.parseNumber()
	}

	// Identifier: constant or function.
	if p.pos < len(p.input) && (unicode.IsLetter(rune(p.input[p.pos])) || p.input[p.pos] == '_') {
		return p.parseIdentifier()
	}

	return 0, fmt.Errorf("unexpected token at position %d", p.pos)
}

func (p *parser) parseNumber() (float64, error) {
	start := p.pos
	for p.pos < len(p.input) && (p.input[p.pos] >= '0' && p.input[p.pos] <= '9' || p.input[p.pos] == '.') {
		p.pos++
	}
	val, err := strconv.ParseFloat(p.input[start:p.pos], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", p.input[start:p.pos])
	}
	return val, nil
}

func (p *parser) parseIdentifier() (float64, error) {
	start := p.pos
	for p.pos < len(p.input) && (unicode.IsLetter(rune(p.input[p.pos])) || unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '_') {
		p.pos++
	}
	name := strings.ToLower(p.input[start:p.pos])

	// Constants.
	switch name {
	case "pi":
		return math.Pi, nil
	case "e":
		return math.E, nil
	}

	// Functions — must be followed by '('.
	p.skipSpaces()
	if p.pos >= len(p.input) || p.input[p.pos] != '(' {
		return 0, fmt.Errorf("unknown identifier %q", name)
	}
	p.pos++ // consume '('

	arg, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	if !p.consume(')') {
		return 0, fmt.Errorf("missing closing parenthesis for %s()", name)
	}

	switch name {
	case "sqrt":
		return math.Sqrt(arg), nil
	case "abs":
		return math.Abs(arg), nil
	case "ceil":
		return math.Ceil(arg), nil
	case "floor":
		return math.Floor(arg), nil
	case "sin":
		return math.Sin(arg), nil
	case "cos":
		return math.Cos(arg), nil
	case "tan":
		return math.Tan(arg), nil
	case "log":
		return math.Log(arg), nil
	case "log2":
		return math.Log2(arg), nil
	case "log10":
		return math.Log10(arg), nil
	default:
		return 0, fmt.Errorf("unknown function %q", name)
	}
}

// formatNumber prints a number nicely (integer if whole, otherwise float).
func formatNumber(f float64) string {
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatFloat(f, 'f', 0, 64)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

package mathcmd

import (
	"math"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestMathCmd_Metadata(t *testing.T) {
	if mathCmd.Use != "math" {
		t.Errorf("unexpected Use: %s", mathCmd.Use)
	}
	if mathCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestSubcommands_Registered(t *testing.T) {
	names := map[string]bool{}
	for _, sub := range mathCmd.Commands() {
		names[sub.Name()] = true
	}
	for _, want := range []string{"eval", "percent", "gcd", "lcm", "factorial", "fibonacci"} {
		if !names[want] {
			t.Errorf("subcommand %q not registered", want)
		}
	}
}

// ── evaluate (expression parser) ───────────────────────────────────────────

func TestEvaluate_BasicArithmetic(t *testing.T) {
	tests := []struct {
		expr string
		want float64
	}{
		{"2 + 3", 5},
		{"10 - 4", 6},
		{"3 * 7", 21},
		{"20 / 4", 5},
		{"7 % 3", 1},
	}
	for _, tt := range tests {
		got, err := evaluate(tt.expr)
		if err != nil {
			t.Errorf("evaluate(%q): %v", tt.expr, err)
			continue
		}
		if got != tt.want {
			t.Errorf("evaluate(%q) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestEvaluate_Precedence(t *testing.T) {
	got, err := evaluate("2 + 3 * 4")
	if err != nil {
		t.Fatal(err)
	}
	if got != 14 {
		t.Errorf("got %v, want 14", got)
	}
}

func TestEvaluate_Parentheses(t *testing.T) {
	got, err := evaluate("(2 + 3) * 4")
	if err != nil {
		t.Fatal(err)
	}
	if got != 20 {
		t.Errorf("got %v, want 20", got)
	}
}

func TestEvaluate_Power(t *testing.T) {
	got, err := evaluate("2^10")
	if err != nil {
		t.Fatal(err)
	}
	if got != 1024 {
		t.Errorf("got %v, want 1024", got)
	}
}

func TestEvaluate_Functions(t *testing.T) {
	tests := []struct {
		expr string
		want float64
	}{
		{"sqrt(144)", 12},
		{"abs(-5)", 5},
		{"ceil(4.2)", 5},
		{"floor(4.8)", 4},
		{"log10(100)", 2},
		{"log2(8)", 3},
	}
	for _, tt := range tests {
		got, err := evaluate(tt.expr)
		if err != nil {
			t.Errorf("evaluate(%q): %v", tt.expr, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("evaluate(%q) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestEvaluate_Constants(t *testing.T) {
	got, err := evaluate("pi")
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-math.Pi) > 1e-9 {
		t.Errorf("got %v, want pi", got)
	}

	got, err = evaluate("e")
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-math.E) > 1e-9 {
		t.Errorf("got %v, want e", got)
	}
}

func TestEvaluate_SinPiOverTwo(t *testing.T) {
	got, err := evaluate("sin(pi / 2)")
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-1.0) > 1e-9 {
		t.Errorf("got %v, want 1.0", got)
	}
}

func TestEvaluate_UnaryMinus(t *testing.T) {
	got, err := evaluate("-5 + 3")
	if err != nil {
		t.Fatal(err)
	}
	if got != -2 {
		t.Errorf("got %v, want -2", got)
	}
}

func TestEvaluate_DivisionByZero(t *testing.T) {
	_, err := evaluate("1 / 0")
	if err == nil {
		t.Error("expected division by zero error")
	}
}

func TestEvaluate_InvalidExpression(t *testing.T) {
	for _, expr := range []string{"", "2 +", "foo(1)", "(2 + 3"} {
		_, err := evaluate(expr)
		if err == nil {
			t.Errorf("expected error for %q", expr)
		}
	}
}

// ── gcd ────────────────────────────────────────────────────────────────────

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b, want int64
	}{
		{12, 18, 6},
		{100, 75, 25},
		{7, 13, 1},
		{0, 5, 5},
		{48, 36, 12},
	}
	for _, tt := range tests {
		got := gcd(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("gcd(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

// ── lcm ────────────────────────────────────────────────────────────────────

func TestLCM(t *testing.T) {
	tests := []struct {
		a, b, want int64
	}{
		{4, 6, 12},
		{3, 5, 15},
		{12, 18, 36},
		{0, 5, 0},
		{7, 7, 7},
	}
	for _, tt := range tests {
		got := lcm(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("lcm(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

// ── factorial ──────────────────────────────────────────────────────────────

func TestFactorial(t *testing.T) {
	tests := []struct {
		n    int
		want float64
	}{
		{0, 1},
		{1, 1},
		{5, 120},
		{10, 3628800},
		{20, 2432902008176640000},
	}
	for _, tt := range tests {
		got := factorial(tt.n)
		if got != tt.want {
			t.Errorf("factorial(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

// ── fibonacci ──────────────────────────────────────────────────────────────

func TestFibonacci(t *testing.T) {
	tests := []struct {
		n    int
		want int64
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{10, 55},
		{20, 6765},
	}
	for _, tt := range tests {
		got := fibonacci(tt.n)
		if got != tt.want {
			t.Errorf("fibonacci(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

// ── formatNumber ───────────────────────────────────────────────────────────

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		val  float64
		want string
	}{
		{42, "42"},
		{3.14, "3.14"},
		{0, "0"},
		{-7, "-7"},
		{100.0, "100"},
	}
	for _, tt := range tests {
		got := formatNumber(tt.val)
		if got != tt.want {
			t.Errorf("formatNumber(%v) = %q, want %q", tt.val, got, tt.want)
		}
	}
}

// ── flags ──────────────────────────────────────────────────────────────────

func TestMathCmd_PersistentFlags(t *testing.T) {
	f := mathCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

func TestPercentCmd_Flags(t *testing.T) {
	f := percentCmd.Flags()
	for _, name := range []string{"of", "is", "change", "value", "total", "from", "to"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found on percent command", name)
		}
	}
}

// ── nested expression ──────────────────────────────────────────────────────

func TestEvaluate_Nested(t *testing.T) {
	got, err := evaluate("sqrt(abs(-16))")
	if err != nil {
		t.Fatal(err)
	}
	if got != 4 {
		t.Errorf("got %v, want 4", got)
	}
}

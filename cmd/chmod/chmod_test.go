package chmod

import (
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestChmodCmd_Metadata(t *testing.T) {
	if chmodCmd.Use != "chmod" {
		t.Errorf("unexpected Use: %s", chmodCmd.Use)
	}
	if chmodCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestSubcommands_Registered(t *testing.T) {
	names := map[string]bool{}
	for _, sub := range chmodCmd.Commands() {
		names[sub.Name()] = true
	}
	for _, want := range []string{"calc", "umask"} {
		if !names[want] {
			t.Errorf("subcommand %q not registered", want)
		}
	}
}

// ── parsePerm (octal) ──────────────────────────────────────────────────────

func TestParsePerm_Octal755(t *testing.T) {
	perm, err := parsePerm("755")
	if err != nil {
		t.Fatal(err)
	}
	if perm != 0755 {
		t.Errorf("got %03o, want 755", perm)
	}
}

func TestParsePerm_Octal644(t *testing.T) {
	perm, err := parsePerm("644")
	if err != nil {
		t.Fatal(err)
	}
	if perm != 0644 {
		t.Errorf("got %03o, want 644", perm)
	}
}

func TestParsePerm_Octal777(t *testing.T) {
	perm, err := parsePerm("777")
	if err != nil {
		t.Fatal(err)
	}
	if perm != 0777 {
		t.Errorf("got %03o, want 777", perm)
	}
}

func TestParsePerm_Octal000(t *testing.T) {
	perm, err := parsePerm("000")
	if err != nil {
		t.Fatal(err)
	}
	if perm != 0 {
		t.Errorf("got %03o, want 000", perm)
	}
}

// ── parsePerm (symbolic) ───────────────────────────────────────────────────

func TestParsePerm_Symbolic_rwxrxrx(t *testing.T) {
	perm, err := parsePerm("rwxr-xr-x")
	if err != nil {
		t.Fatal(err)
	}
	if perm != 0755 {
		t.Errorf("got %03o, want 755", perm)
	}
}

func TestParsePerm_Symbolic_rwrr(t *testing.T) {
	perm, err := parsePerm("rw-r--r--")
	if err != nil {
		t.Fatal(err)
	}
	if perm != 0644 {
		t.Errorf("got %03o, want 644", perm)
	}
}

func TestParsePerm_Symbolic_none(t *testing.T) {
	perm, err := parsePerm("---------")
	if err != nil {
		t.Fatal(err)
	}
	if perm != 0 {
		t.Errorf("got %03o, want 000", perm)
	}
}

// ── parsePerm (invalid) ────────────────────────────────────────────────────

func TestParsePerm_Invalid(t *testing.T) {
	invalids := []string{"abc", "999", "rwxrwx", "12", ""}
	for _, input := range invalids {
		_, err := parsePerm(input)
		if err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

// ── toSymbolic ─────────────────────────────────────────────────────────────

func TestToSymbolic(t *testing.T) {
	tests := []struct {
		perm uint16
		want string
	}{
		{0755, "rwxr-xr-x"},
		{0644, "rw-r--r--"},
		{0777, "rwxrwxrwx"},
		{0000, "---------"},
		{0700, "rwx------"},
		{0100, "--x------"},
	}
	for _, tt := range tests {
		got := toSymbolic(tt.perm)
		if got != tt.want {
			t.Errorf("toSymbolic(%03o) = %q, want %q", tt.perm, got, tt.want)
		}
	}
}

// ── fromSymbolic ───────────────────────────────────────────────────────────

func TestFromSymbolic(t *testing.T) {
	tests := []struct {
		sym  string
		want uint16
	}{
		{"rwxr-xr-x", 0755},
		{"rw-r--r--", 0644},
		{"rwxrwxrwx", 0777},
		{"---------", 0000},
	}
	for _, tt := range tests {
		got, err := fromSymbolic(tt.sym)
		if err != nil {
			t.Errorf("fromSymbolic(%q): %v", tt.sym, err)
			continue
		}
		if got != tt.want {
			t.Errorf("fromSymbolic(%q) = %03o, want %03o", tt.sym, got, tt.want)
		}
	}
}

func TestFromSymbolic_InvalidLength(t *testing.T) {
	_, err := fromSymbolic("rwx")
	if err == nil {
		t.Error("expected error for short symbolic string")
	}
}

// ── permBreakdown ──────────────────────────────────────────────────────────

func TestPermBreakdown_755(t *testing.T) {
	bd := permBreakdown(0755)
	if len(bd) != 3 {
		t.Fatalf("expected 3 classes, got %d", len(bd))
	}
	// Owner: rwx
	if bd[0]["read"] != "yes" || bd[0]["write"] != "yes" || bd[0]["execute"] != "yes" {
		t.Errorf("owner breakdown wrong: %v", bd[0])
	}
	// Group: r-x
	if bd[1]["read"] != "yes" || bd[1]["write"] != "no" || bd[1]["execute"] != "yes" {
		t.Errorf("group breakdown wrong: %v", bd[1])
	}
	// Other: r-x
	if bd[2]["read"] != "yes" || bd[2]["write"] != "no" || bd[2]["execute"] != "yes" {
		t.Errorf("other breakdown wrong: %v", bd[2])
	}
}

// ── roundtrip ──────────────────────────────────────────────────────────────

func TestRoundtrip_OctalToSymbolicAndBack(t *testing.T) {
	perms := []uint16{0755, 0644, 0777, 0700, 0600, 0100, 0000}
	for _, p := range perms {
		sym := toSymbolic(p)
		back, err := fromSymbolic(sym)
		if err != nil {
			t.Errorf("roundtrip failed for %03o: %v", p, err)
			continue
		}
		if back != p {
			t.Errorf("roundtrip: %03o -> %q -> %03o", p, sym, back)
		}
	}
}

// ── flags ──────────────────────────────────────────────────────────────────

func TestChmodCmd_PersistentFlags(t *testing.T) {
	f := chmodCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

// ── isOctalString / isSymbolicString ───────────────────────────────────────

func TestIsOctalString(t *testing.T) {
	if !isOctalString("755") {
		t.Error("755 should be octal")
	}
	if isOctalString("999") {
		t.Error("999 should not be octal")
	}
	if isOctalString("abc") {
		t.Error("abc should not be octal")
	}
}

func TestIsSymbolicString(t *testing.T) {
	if !isSymbolicString("rwxr-xr-x") {
		t.Error("rwxr-xr-x should be symbolic")
	}
	if isSymbolicString("rwxabc123") {
		t.Error("rwxabc123 should not be symbolic")
	}
}

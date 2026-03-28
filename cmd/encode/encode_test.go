package encode

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestEncodeCmd_Metadata(t *testing.T) {
	if encodeCmd.Use == "" {
		t.Error("encodeCmd.Use must not be empty")
	}
	if encodeCmd.Short == "" {
		t.Error("encodeCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"base64Cmd", base64Cmd.Use, base64Cmd.Short},
		{"hexCmd", hexCmd.Use, hexCmd.Short},
		{"urlCmd", urlCmd.Use, urlCmd.Short},
		{"binaryCmd", binaryCmd.Use, binaryCmd.Short},
		{"rot13Cmd", rot13Cmd.Use, rot13Cmd.Short},
		{"morseCmd", morseCmd.Use, morseCmd.Short},
		{"base32Cmd", base32Cmd.Use, base32Cmd.Short},
		{"base58Cmd", base58Cmd.Use, base58Cmd.Short},
		{"htmlCmd", htmlCmd.Use, htmlCmd.Short},
		{"punycodeCmd", punycodeCmd.Use, punycodeCmd.Short},
		{"jwtCmd", jwtCmd.Use, jwtCmd.Short},
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

// ── flag existence and defaults ────────────────────────────────────────────

func TestEncodeCmd_PersistentFlags(t *testing.T) {
	f := encodeCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
	if f.Lookup("file") == nil {
		t.Error("expected persistent flag --file")
	}
}

func TestDecodeFlag_ExistsOnSubcommands(t *testing.T) {
	for _, sc := range []struct {
		name string
		hasD bool
	}{
		{"base64", true},
		{"hex", true},
		{"url", true},
		{"binary", true},
		{"morse", true},
	} {
		found := false
		for _, sub := range encodeCmd.Commands() {
			if sub.Name() == sc.name {
				if sub.Flags().Lookup("decode") != nil {
					found = true
				}
				break
			}
		}
		if sc.hasD && !found {
			t.Errorf("expected --decode flag on %s subcommand", sc.name)
		}
	}
}

// ── resolveInput ───────────────────────────────────────────────────────────

func TestResolveInput_FromArgs(t *testing.T) {
	file = ""
	got, err := resolveInput([]string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestResolveInput_NoArgs_NoFile(t *testing.T) {
	file = ""
	_, err := resolveInput(nil)
	if err == nil {
		t.Error("expected error when no args and no file")
	}
}

func TestResolveInput_FromFile(t *testing.T) {
	file = "/dev/null"
	got, err := resolveInput(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string from /dev/null, got %q", got)
	}
	file = ""
}

func TestResolveInput_BadFile(t *testing.T) {
	file = "/nonexistent/path/to/file.txt"
	_, err := resolveInput(nil)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	file = ""
}

// ── processEncode ──────────────────────────────────────────────────────────

func TestProcessEncode_CallsFn(t *testing.T) {
	file = ""
	decode = false
	jsonOut = false

	called := false
	fn := func(input string, dec bool) (string, error) {
		called = true
		if input != "test" {
			t.Errorf("expected input %q, got %q", "test", input)
		}
		if dec {
			t.Error("expected decode=false")
		}
		return "result", nil
	}

	err := processEncode([]string{"test"}, "test-encoding", fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected fn to be called")
	}
}

func TestProcessEncode_NoInput(t *testing.T) {
	file = ""
	decode = false
	jsonOut = false

	fn := func(input string, dec bool) (string, error) {
		return "", nil
	}

	err := processEncode(nil, "test", fn)
	if err == nil {
		t.Error("expected error when no input")
	}
}

// ── base64 encode/decode ───────────────────────────────────────────────────

func TestBase64_EncodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"hello world", "aGVsbG8gd29ybGQ="},
		{"", ""},
		{"a", "YQ=="},
		{"openGyver", "b3Blbkd5dmVy"},
	}
	for _, tt := range tests {
		got := base64.StdEncoding.EncodeToString([]byte(tt.input))
		if got != tt.want {
			t.Errorf("base64 encode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBase64_DecodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"aGVsbG8gd29ybGQ=", "hello world"},
		{"", ""},
		{"YQ==", "a"},
	}
	for _, tt := range tests {
		b, err := base64.StdEncoding.DecodeString(tt.input)
		if err != nil {
			t.Fatalf("base64 decode(%q): %v", tt.input, err)
		}
		if string(b) != tt.want {
			t.Errorf("base64 decode(%q) = %q, want %q", tt.input, string(b), tt.want)
		}
	}
}

func TestBase64_RoundTrip(t *testing.T) {
	inputs := []string{"hello world", "openGyver", "special chars: !@#$%"}
	for _, input := range inputs {
		encoded := base64.StdEncoding.EncodeToString([]byte(input))
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			t.Fatalf("decode error for %q: %v", input, err)
		}
		if string(decoded) != input {
			t.Errorf("roundtrip failed: got %q, want %q", string(decoded), input)
		}
	}
}

// ── hex encode/decode ──────────────────────────────────────────────────────

func TestHex_EncodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"hello", "68656c6c6f"},
		{"", ""},
		{"AB", "4142"},
	}
	for _, tt := range tests {
		got := hex.EncodeToString([]byte(tt.input))
		if got != tt.want {
			t.Errorf("hex encode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestHex_DecodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"68656c6c6f", "hello"},
		{"4142", "AB"},
	}
	for _, tt := range tests {
		b, err := hex.DecodeString(tt.input)
		if err != nil {
			t.Fatalf("hex decode(%q): %v", tt.input, err)
		}
		if string(b) != tt.want {
			t.Errorf("hex decode(%q) = %q, want %q", tt.input, string(b), tt.want)
		}
	}
}

func TestHex_DecodeInvalid(t *testing.T) {
	_, err := hex.DecodeString("zzzz")
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

// ── url encode/decode ──────────────────────────────────────────────────────

func TestURL_EncodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"hello world", "hello+world"},
		{"foo=bar&baz=qux", "foo%3Dbar%26baz%3Dqux"},
		{"", ""},
	}
	for _, tt := range tests {
		got := url.QueryEscape(tt.input)
		if got != tt.want {
			t.Errorf("url encode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestURL_DecodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"hello+world", "hello world"},
		{"foo%3Dbar", "foo=bar"},
	}
	for _, tt := range tests {
		got, err := url.QueryUnescape(tt.input)
		if err != nil {
			t.Fatalf("url decode(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("url decode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ── rot13 ──────────────────────────────────────────────────────────────────

func TestRot13_Known(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Hello", "Uryyb"},
		{"hello", "uryyb"},
		{"ABC", "NOP"},
		{"abc", "nop"},
		{"Hello, World!", "Uryyb, Jbeyq!"},
		{"123", "123"},   // digits unchanged
		{"", ""},          // empty string
	}
	for _, tt := range tests {
		got := rot13(tt.input)
		if got != tt.want {
			t.Errorf("rot13(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRot13_Symmetric(t *testing.T) {
	inputs := []string{"Hello, World!", "ABCxyz", "Testing 123!"}
	for _, input := range inputs {
		if rot13(rot13(input)) != input {
			t.Errorf("rot13 is not symmetric for %q", input)
		}
	}
}

// ── binary encode/decode ───────────────────────────────────────────────────

func TestBinary_EncodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"A", "01000001"},
		{"AB", "01000001 01000010"},
		{"Hi", "01001000 01101001"},
	}
	for _, tt := range tests {
		got := binaryEncode(tt.input)
		if got != tt.want {
			t.Errorf("binaryEncode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBinary_DecodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"01000001", "A"},
		{"01000001 01000010", "AB"},
		{"01001000 01101001", "Hi"},
	}
	for _, tt := range tests {
		got, err := binaryDecode(tt.input)
		if err != nil {
			t.Fatalf("binaryDecode(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("binaryDecode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBinary_RoundTrip(t *testing.T) {
	inputs := []string{"Hello", "openGyver", "test 123"}
	for _, input := range inputs {
		encoded := binaryEncode(input)
		decoded, err := binaryDecode(encoded)
		if err != nil {
			t.Fatalf("binaryDecode error: %v", err)
		}
		if decoded != input {
			t.Errorf("binary roundtrip failed: got %q, want %q", decoded, input)
		}
	}
}

func TestBinary_DecodeInvalid(t *testing.T) {
	_, err := binaryDecode("0102")
	if err == nil {
		t.Error("expected error for invalid binary character")
	}
}

// ── morse encode/decode ────────────────────────────────────────────────────

func TestMorse_EncodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"SOS", "... --- ..."},
		{"HELLO", ".... . .-.. .-.. ---"},
		{"A B", ".- / -..."},
	}
	for _, tt := range tests {
		got, err := morseEncode(tt.input)
		if err != nil {
			t.Fatalf("morseEncode(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("morseEncode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMorse_DecodeKnown(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"... --- ...", "SOS"},
		{".... . .-.. .-.. ---", "HELLO"},
		{".- / -...", "A B"},
	}
	for _, tt := range tests {
		got, err := morseDecode(tt.input)
		if err != nil {
			t.Fatalf("morseDecode(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("morseDecode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMorse_RoundTrip(t *testing.T) {
	inputs := []string{"SOS", "HELLO WORLD", "TEST 123"}
	for _, input := range inputs {
		encoded, err := morseEncode(input)
		if err != nil {
			t.Fatalf("morseEncode(%q): %v", input, err)
		}
		decoded, err := morseDecode(encoded)
		if err != nil {
			t.Fatalf("morseDecode(%q): %v", encoded, err)
		}
		if decoded != strings.ToUpper(input) {
			t.Errorf("morse roundtrip failed: got %q, want %q", decoded, strings.ToUpper(input))
		}
	}
}

func TestMorse_EncodeUnsupportedChar(t *testing.T) {
	_, err := morseEncode("\x00")
	if err == nil {
		t.Error("expected error for unsupported character")
	}
}

func TestMorse_DecodeUnknownCode(t *testing.T) {
	_, err := morseDecode("........")
	if err == nil {
		t.Error("expected error for unknown morse code")
	}
}

// ── base58 encode/decode ───────────────────────────────────────────────────

func TestBase58_RoundTrip(t *testing.T) {
	inputs := []string{"hello", "Hello World", "test"}
	for _, input := range inputs {
		encoded := base58Encode([]byte(input))
		decoded, err := base58Decode(encoded)
		if err != nil {
			t.Fatalf("base58Decode error: %v", err)
		}
		if string(decoded) != input {
			t.Errorf("base58 roundtrip failed: got %q, want %q", string(decoded), input)
		}
	}
}

func TestBase58_DecodeInvalidChar(t *testing.T) {
	// '0', 'O', 'I', 'l' are not in base58 alphabet
	_, err := base58Decode("0OIl")
	if err == nil {
		t.Error("expected error for invalid base58 characters")
	}
}

// ── jwt decode ─────────────────────────────────────────────────────────────

func TestJWTDecode_Valid(t *testing.T) {
	// This is a well-known test JWT from jwt.io
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	result, err := jwtDecode(token)
	if err != nil {
		t.Fatalf("jwtDecode error: %v", err)
	}
	if !strings.Contains(result, "John Doe") {
		t.Errorf("expected payload to contain 'John Doe', got %q", result)
	}
	if !strings.Contains(result, "1234567890") {
		t.Errorf("expected payload to contain '1234567890', got %q", result)
	}
}

func TestJWTDecode_InvalidParts(t *testing.T) {
	_, err := jwtDecode("not.a.valid.jwt.token")
	if err == nil {
		t.Error("expected error for token with wrong number of parts")
	}

	_, err = jwtDecode("onlytwoparts.here")
	if err == nil {
		t.Error("expected error for token with only two parts")
	}
}

func TestJWTDecode_InvalidBase64(t *testing.T) {
	_, err := jwtDecode("header.!!!invalid!!!.signature")
	if err == nil {
		t.Error("expected error for invalid base64 payload")
	}
}

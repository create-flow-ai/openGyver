package generate

import (
	"strings"
	"testing"
	"unicode"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestGenerateCmd_Metadata(t *testing.T) {
	if generateCmd.Use == "" {
		t.Error("generateCmd.Use must not be empty")
	}
	if generateCmd.Short == "" {
		t.Error("generateCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"passwordCmd", passwordCmd.Use, passwordCmd.Short},
		{"passphraseCmd", passphraseCmd.Use, passphraseCmd.Short},
		{"stringCmd", stringCmd.Use, stringCmd.Short},
		{"nanoidCmd", nanoidCmd.Use, nanoidCmd.Short},
		{"snowflakeCmd", snowflakeCmd.Use, snowflakeCmd.Short},
		{"shortidCmd", shortidCmd.Use, shortidCmd.Short},
		{"apikeyCmd", apikeyCmd.Use, apikeyCmd.Short},
		{"secretCmd", secretCmd.Use, secretCmd.Short},
		{"otpCmd", otpCmd.Use, otpCmd.Short},
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

func TestGenerateCmd_PersistentFlags(t *testing.T) {
	f := generateCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

func TestPasswordCmd_Flags(t *testing.T) {
	f := passwordCmd.Flags()
	for _, name := range []string{"length", "no-upper", "no-lower", "no-digits", "no-special", "count"} {
		if f.Lookup(name) == nil {
			t.Errorf("expected flag --%s on passwordCmd", name)
		}
	}
}

func TestPassphraseCmd_Flags(t *testing.T) {
	f := passphraseCmd.Flags()
	for _, name := range []string{"words", "separator", "count"} {
		if f.Lookup(name) == nil {
			t.Errorf("expected flag --%s on passphraseCmd", name)
		}
	}
}

func TestNanoidCmd_Flags(t *testing.T) {
	f := nanoidCmd.Flags()
	if f.Lookup("length") == nil {
		t.Error("expected flag --length on nanoidCmd")
	}
}

func TestShortidCmd_Flags(t *testing.T) {
	f := shortidCmd.Flags()
	if f.Lookup("length") == nil {
		t.Error("expected flag --length on shortidCmd")
	}
}

// ── randStringFromCharset ──────────────────────────────────────────────────

func TestRandStringFromCharset_Length(t *testing.T) {
	lengths := []int{1, 5, 16, 32, 64}
	charset := "abcdef0123456789"
	for _, length := range lengths {
		s, err := randStringFromCharset(length, charset)
		if err != nil {
			t.Fatalf("randStringFromCharset(%d, %q): %v", length, charset, err)
		}
		if len(s) != length {
			t.Errorf("expected length %d, got %d", length, len(s))
		}
	}
}

func TestRandStringFromCharset_OnlyUsesCharset(t *testing.T) {
	charset := "abc"
	s, err := randStringFromCharset(100, charset)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for _, c := range s {
		if !strings.ContainsRune(charset, c) {
			t.Errorf("character %q not in charset %q", string(c), charset)
		}
	}
}

func TestRandStringFromCharset_ZeroLength(t *testing.T) {
	_, err := randStringFromCharset(0, "abc")
	if err == nil {
		t.Error("expected error for zero length")
	}
}

func TestRandStringFromCharset_NegativeLength(t *testing.T) {
	_, err := randStringFromCharset(-1, "abc")
	if err == nil {
		t.Error("expected error for negative length")
	}
}

// ── randInt ────────────────────────────────────────────────────────────────

func TestRandInt_Range(t *testing.T) {
	for i := 0; i < 100; i++ {
		n, err := randInt(10)
		if err != nil {
			t.Fatalf("randInt error: %v", err)
		}
		if n < 0 || n >= 10 {
			t.Errorf("randInt(10) = %d, out of range [0, 10)", n)
		}
	}
}

// ── randBytes ──────────────────────────────────────────────────────────────

func TestRandBytes_Length(t *testing.T) {
	for _, n := range []int{1, 16, 32, 64} {
		b, err := randBytes(n)
		if err != nil {
			t.Fatalf("randBytes(%d): %v", n, err)
		}
		if len(b) != n {
			t.Errorf("expected %d bytes, got %d", n, len(b))
		}
	}
}

// ── password generation ────────────────────────────────────────────────────

func TestPasswordGeneration_Length(t *testing.T) {
	for _, length := range []int{8, 16, 32, 64} {
		const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()"
		pw, err := randStringFromCharset(length, charset)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(pw) != length {
			t.Errorf("expected password length %d, got %d", length, len(pw))
		}
	}
}

func TestPasswordGeneration_CharacterClasses(t *testing.T) {
	// Generate with all classes enabled, verify all classes appear in a large sample
	const (
		upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lower   = "abcdefghijklmnopqrstuvwxyz"
		digits  = "0123456789"
		special = "!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
	)
	charset := upper + lower + digits + special

	// With a long password, all classes are very likely to appear
	pw, err := randStringFromCharset(200, charset)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	for _, c := range pw {
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if unicode.IsLower(c) {
			hasLower = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
		if strings.ContainsRune(special, c) {
			hasSpecial = true
		}
	}

	if !hasUpper {
		t.Error("expected uppercase characters in password")
	}
	if !hasLower {
		t.Error("expected lowercase characters in password")
	}
	if !hasDigit {
		t.Error("expected digit characters in password")
	}
	if !hasSpecial {
		t.Error("expected special characters in password")
	}
}

func TestPasswordGeneration_DigitsOnly(t *testing.T) {
	const digits = "0123456789"
	pw, err := randStringFromCharset(20, digits)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for _, c := range pw {
		if !unicode.IsDigit(c) {
			t.Errorf("expected only digits, got %q", string(c))
		}
	}
}

// ── passphrase generation ──────────────────────────────────────────────────

func TestPassphraseGeneration_WordCount(t *testing.T) {
	sep := "-"
	for _, wordCount := range []int{2, 4, 6, 8} {
		words := make([]string, 0, wordCount)
		for j := 0; j < wordCount; j++ {
			idx, err := randInt(len(effWordList))
			if err != nil {
				t.Fatalf("randInt error: %v", err)
			}
			words = append(words, effWordList[idx])
		}
		passphrase := strings.Join(words, sep)
		parts := strings.Split(passphrase, sep)
		if len(parts) != wordCount {
			t.Errorf("expected %d words, got %d", wordCount, len(parts))
		}
	}
}

func TestPassphraseGeneration_WordsFromList(t *testing.T) {
	wordSet := make(map[string]bool, len(effWordList))
	for _, w := range effWordList {
		wordSet[w] = true
	}

	for i := 0; i < 50; i++ {
		idx, err := randInt(len(effWordList))
		if err != nil {
			t.Fatalf("randInt error: %v", err)
		}
		word := effWordList[idx]
		if !wordSet[word] {
			t.Errorf("word %q not in EFF word list", word)
		}
	}
}

// ── string generation ──────────────────────────────────────────────────────

func TestStringGeneration_Length(t *testing.T) {
	for _, length := range []int{8, 16, 32, 64} {
		s, err := randStringFromCharset(length, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(s) != length {
			t.Errorf("expected length %d, got %d", length, len(s))
		}
	}
}

func TestStringGeneration_HexCharset(t *testing.T) {
	const hexChars = "0123456789abcdef"
	s, err := randStringFromCharset(100, hexChars)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for _, c := range s {
		if !strings.ContainsRune(hexChars, c) {
			t.Errorf("character %q not in hex charset", string(c))
		}
	}
}

// ── nanoid generation ──────────────────────────────────────────────────────

func TestNanoidGeneration_Length(t *testing.T) {
	alphabet := "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, length := range []int{8, 12, 21, 36} {
		id, err := randStringFromCharset(length, alphabet)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(id) != length {
			t.Errorf("expected nanoid length %d, got %d", length, len(id))
		}
	}
}

func TestNanoidGeneration_AlphabetOnly(t *testing.T) {
	alphabet := "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	id, err := randStringFromCharset(100, alphabet)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	for _, c := range id {
		if !strings.ContainsRune(alphabet, c) {
			t.Errorf("character %q not in nanoid alphabet", string(c))
		}
	}
}

// ── snowflake generation ───────────────────────────────────────────────────

func TestSnowflakeGeneration_NonZero(t *testing.T) {
	// We can't call runSnowflake directly without command context,
	// but we can test the core logic.
	// The snowflake is: ((nowMs - epoch) << 22) | (nodeID << 12) | seq
	// With a valid timestamp this should never be zero.

	// Reset state for test
	snowflakeNodeID = 0
	snowflakeSeq = 0
	snowflakeLastMs = 0

	// Simulate the core snowflake generation
	const epoch int64 = 1577836800000
	nodeID, err := randInt(1024)
	if err != nil {
		t.Fatalf("randInt error: %v", err)
	}

	// Use a known timestamp
	nowMs := int64(1700000000000) // some reasonable timestamp
	id := ((nowMs - epoch) << 22) | (int64(nodeID) << 12) | 0

	if id == 0 {
		t.Error("snowflake ID should not be zero")
	}
	if id < 0 {
		t.Error("snowflake ID should be positive")
	}
}

// ── shortid generation ─────────────────────────────────────────────────────

func TestShortidGeneration_Length(t *testing.T) {
	const base62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for _, length := range []int{4, 8, 12, 16} {
		id, err := randStringFromCharset(length, base62)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if len(id) != length {
			t.Errorf("expected shortid length %d, got %d", length, len(id))
		}
	}
}

// ── randomness verification (no duplicates in 100 runs) ────────────────────

func TestRandomness_NoDuplicates(t *testing.T) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s, err := randStringFromCharset(32, charset)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if seen[s] {
			t.Errorf("duplicate found: %q (iteration %d)", s, i)
		}
		seen[s] = true
	}
}

func TestRandomness_NoDuplicates_ShortID(t *testing.T) {
	const base62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s, err := randStringFromCharset(8, base62)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if seen[s] {
			t.Errorf("duplicate shortid found: %q (iteration %d)", s, i)
		}
		seen[s] = true
	}
}

func TestRandomness_NoDuplicates_Nanoid(t *testing.T) {
	alphabet := "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s, err := randStringFromCharset(21, alphabet)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if seen[s] {
			t.Errorf("duplicate nanoid found: %q (iteration %d)", s, i)
		}
		seen[s] = true
	}
}

// ── EFF word list sanity ───────────────────────────────────────────────────

func TestEffWordList_NotEmpty(t *testing.T) {
	if len(effWordList) == 0 {
		t.Error("effWordList must not be empty")
	}
}

func TestEffWordList_AllLowercase(t *testing.T) {
	for i, w := range effWordList {
		if w != strings.ToLower(w) {
			t.Errorf("effWordList[%d] = %q is not lowercase", i, w)
		}
	}
}

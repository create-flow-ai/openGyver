package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash/crc32"
	"os"
	"path/filepath"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestChecksumCmd_Metadata(t *testing.T) {
	if checksumCmd.Use == "" {
		t.Error("checksumCmd.Use must not be empty")
	}
	if checksumCmd.Short == "" {
		t.Error("checksumCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"calcCmd", calcCmd.Use, calcCmd.Short},
		{"verifyCmd", verifyCmd.Use, verifyCmd.Short},
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

func TestChecksumCmd_PersistentFlags(t *testing.T) {
	f := checksumCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
	if f.Lookup("algorithm") == nil {
		t.Error("expected persistent flag --algorithm")
	}
}

func TestAlgorithmDefault(t *testing.T) {
	val, err := checksumCmd.PersistentFlags().GetString("algorithm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "sha256" {
		t.Errorf("expected default algorithm sha256, got %q", val)
	}
}

// ── newHash ───────────────────────────────────────────────────────────────

func TestNewHash_AllAlgorithms(t *testing.T) {
	algos := []string{"md5", "sha1", "sha256", "sha512", "crc32"}
	for _, a := range algos {
		h, err := newHash(a)
		if err != nil {
			t.Errorf("newHash(%q) returned error: %v", a, err)
		}
		if h == nil {
			t.Errorf("newHash(%q) returned nil hash", a)
		}
	}
}

func TestNewHash_Unsupported(t *testing.T) {
	_, err := newHash("blake2b")
	if err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}

// ── computeChecksum known values ──────────────────────────────────────────

func writeTestFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}
	return path
}

func TestComputeChecksum_SHA256(t *testing.T) {
	path := writeTestFile(t, "hello")

	got, err := computeChecksum(path, "sha256")
	if err != nil {
		t.Fatalf("computeChecksum error: %v", err)
	}

	h := sha256.Sum256([]byte("hello"))
	want := hex.EncodeToString(h[:])
	if got != want {
		t.Errorf("sha256 mismatch: got %q, want %q", got, want)
	}
}

func TestComputeChecksum_MD5(t *testing.T) {
	path := writeTestFile(t, "hello")

	got, err := computeChecksum(path, "md5")
	if err != nil {
		t.Fatalf("computeChecksum error: %v", err)
	}

	h := md5.Sum([]byte("hello"))
	want := hex.EncodeToString(h[:])
	if got != want {
		t.Errorf("md5 mismatch: got %q, want %q", got, want)
	}
}

func TestComputeChecksum_SHA1(t *testing.T) {
	path := writeTestFile(t, "hello")

	got, err := computeChecksum(path, "sha1")
	if err != nil {
		t.Fatalf("computeChecksum error: %v", err)
	}

	h := sha1.Sum([]byte("hello"))
	want := hex.EncodeToString(h[:])
	if got != want {
		t.Errorf("sha1 mismatch: got %q, want %q", got, want)
	}
}

func TestComputeChecksum_SHA512(t *testing.T) {
	path := writeTestFile(t, "hello")

	got, err := computeChecksum(path, "sha512")
	if err != nil {
		t.Fatalf("computeChecksum error: %v", err)
	}

	h := sha512.Sum512([]byte("hello"))
	want := hex.EncodeToString(h[:])
	if got != want {
		t.Errorf("sha512 mismatch: got %q, want %q", got, want)
	}
}

func TestComputeChecksum_CRC32(t *testing.T) {
	path := writeTestFile(t, "hello")

	got, err := computeChecksum(path, "crc32")
	if err != nil {
		t.Fatalf("computeChecksum error: %v", err)
	}

	h := crc32.NewIEEE()
	h.Write([]byte("hello"))
	want := "3610a686"
	if got != want {
		t.Errorf("crc32 mismatch: got %q, want %q", got, want)
	}
}

func TestComputeChecksum_EmptyFile(t *testing.T) {
	path := writeTestFile(t, "")

	got, err := computeChecksum(path, "sha256")
	if err != nil {
		t.Fatalf("computeChecksum error: %v", err)
	}

	h := sha256.Sum256([]byte(""))
	want := hex.EncodeToString(h[:])
	if got != want {
		t.Errorf("sha256 of empty file: got %q, want %q", got, want)
	}
}

func TestComputeChecksum_NonexistentFile(t *testing.T) {
	_, err := computeChecksum("/nonexistent/path/to/file.bin", "sha256")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestComputeChecksum_UnsupportedAlgorithm(t *testing.T) {
	path := writeTestFile(t, "test")
	_, err := computeChecksum(path, "blake2b")
	if err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}

// ── runVerify logic (unit-test the match comparison) ─────────────────────

func TestVerify_Match(t *testing.T) {
	path := writeTestFile(t, "hello")
	h := sha256.Sum256([]byte("hello"))
	expected := hex.EncodeToString(h[:])

	// Reset package-level state for test.
	algorithm = "sha256"
	jsonOut = false

	err := runVerify(nil, []string{path, expected})
	if err != nil {
		t.Fatalf("runVerify error: %v", err)
	}
}

func TestVerify_Mismatch(t *testing.T) {
	path := writeTestFile(t, "hello")

	algorithm = "sha256"
	jsonOut = false

	// Pass a wrong hash — should not error, just print MISMATCH.
	err := runVerify(nil, []string{path, "0000000000000000000000000000000000000000000000000000000000000000"})
	if err != nil {
		t.Fatalf("runVerify error: %v", err)
	}
}

func TestVerify_CaseInsensitive(t *testing.T) {
	path := writeTestFile(t, "hello")
	h := sha256.Sum256([]byte("hello"))
	// Use uppercase hex — should still match.
	expected := hex.EncodeToString(h[:])
	uppercased := ""
	for _, c := range expected {
		if c >= 'a' && c <= 'f' {
			uppercased += string(c - 32)
		} else {
			uppercased += string(c)
		}
	}

	algorithm = "sha256"
	jsonOut = false

	err := runVerify(nil, []string{path, uppercased})
	if err != nil {
		t.Fatalf("runVerify error: %v", err)
	}
}

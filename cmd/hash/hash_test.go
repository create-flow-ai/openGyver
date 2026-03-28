package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash/adler32"
	"hash/crc32"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestHashCmd_Metadata(t *testing.T) {
	if hashCmd.Use == "" {
		t.Error("hashCmd.Use must not be empty")
	}
	if hashCmd.Short == "" {
		t.Error("hashCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"md5Cmd", md5Cmd.Use, md5Cmd.Short},
		{"sha1Cmd", sha1Cmd.Use, sha1Cmd.Short},
		{"sha256Cmd", sha256Cmd.Use, sha256Cmd.Short},
		{"sha384Cmd", sha384Cmd.Use, sha384Cmd.Short},
		{"sha512Cmd", sha512Cmd.Use, sha512Cmd.Short},
		{"hmacCmd", hmacCmd.Use, hmacCmd.Short},
		{"bcryptCmd", bcryptCmd.Use, bcryptCmd.Short},
		{"crc32Cmd", crc32Cmd.Use, crc32Cmd.Short},
		{"adler32Cmd", adler32Cmd.Use, adler32Cmd.Short},
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

func TestHashCmd_PersistentFlags(t *testing.T) {
	f := hashCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
	if f.Lookup("uppercase") == nil {
		t.Error("expected persistent flag --uppercase")
	}
	if f.Lookup("file") == nil {
		t.Error("expected persistent flag --file")
	}
}

func TestHmacCmd_Flags(t *testing.T) {
	f := hmacCmd.Flags()
	if f.Lookup("key") == nil {
		t.Error("expected flag --key on hmacCmd")
	}
	if f.Lookup("algorithm") == nil {
		t.Error("expected flag --algorithm on hmacCmd")
	}
}

func TestBcryptCmd_Flags(t *testing.T) {
	f := bcryptCmd.Flags()
	if f.Lookup("rounds") == nil {
		t.Error("expected flag --rounds on bcryptCmd")
	}
	if f.Lookup("verify") == nil {
		t.Error("expected flag --verify on bcryptCmd")
	}
}

// ── readInput ──────────────────────────────────────────────────────────────

func TestReadInput_FromArgs(t *testing.T) {
	filePath = ""
	input, data, err := readInput([]string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input != "hello" {
		t.Errorf("expected input %q, got %q", "hello", input)
	}
	if string(data) != "hello" {
		t.Errorf("expected data %q, got %q", "hello", string(data))
	}
}

func TestReadInput_NoArgs_NoFile(t *testing.T) {
	filePath = ""
	_, _, err := readInput(nil)
	if err == nil {
		t.Error("expected error when no args and no file")
	}
}

func TestReadInput_FromFile(t *testing.T) {
	filePath = "/dev/null"
	input, data, err := readInput(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// /dev/null yields the file path as input and empty data
	if input != "/dev/null" {
		t.Errorf("expected input %q, got %q", "/dev/null", input)
	}
	if len(data) != 0 {
		t.Errorf("expected empty data from /dev/null, got %d bytes", len(data))
	}
	filePath = ""
}

func TestReadInput_BadFile(t *testing.T) {
	filePath = "/nonexistent/path/to/file.txt"
	_, _, err := readInput(nil)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	filePath = ""
}

// ── MD5 known hash ─────────────────────────────────────────────────────────

func TestMD5_KnownHash(t *testing.T) {
	h := md5.New()
	h.Write([]byte("hello"))
	got := hex.EncodeToString(h.Sum(nil))
	want := "5d41402abc4b2a76b9719d911017c592"
	if got != want {
		t.Errorf("md5(hello) = %q, want %q", got, want)
	}
}

func TestMD5_EmptyString(t *testing.T) {
	h := md5.New()
	h.Write([]byte(""))
	got := hex.EncodeToString(h.Sum(nil))
	want := "d41d8cd98f00b204e9800998ecf8427e"
	if got != want {
		t.Errorf("md5('') = %q, want %q", got, want)
	}
}

// ── SHA256 known hash ──────────────────────────────────────────────────────

func TestSHA256_KnownHash(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("hello"))
	got := hex.EncodeToString(h.Sum(nil))
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got != want {
		t.Errorf("sha256(hello) = %q, want %q", got, want)
	}
}

func TestSHA256_EmptyString(t *testing.T) {
	h := sha256.New()
	h.Write([]byte(""))
	got := hex.EncodeToString(h.Sum(nil))
	want := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got != want {
		t.Errorf("sha256('') = %q, want %q", got, want)
	}
}

// ── CRC32 known checksum ───────────────────────────────────────────────────

func TestCRC32_KnownChecksum(t *testing.T) {
	h := crc32.NewIEEE()
	h.Write([]byte("hello"))
	got := h.Sum32()
	want := uint32(0x3610a686)
	if got != want {
		t.Errorf("crc32(hello) = 0x%08x, want 0x%08x", got, want)
	}
}

func TestCRC32_EmptyString(t *testing.T) {
	h := crc32.NewIEEE()
	h.Write([]byte(""))
	got := h.Sum32()
	want := uint32(0x00000000)
	if got != want {
		t.Errorf("crc32('') = 0x%08x, want 0x%08x", got, want)
	}
}

// ── Adler32 known checksum ─────────────────────────────────────────────────

func TestAdler32_KnownChecksum(t *testing.T) {
	h := adler32.New()
	h.Write([]byte("hello"))
	got := h.Sum32()
	want := uint32(0x062c0215)
	if got != want {
		t.Errorf("adler32(hello) = 0x%08x, want 0x%08x", got, want)
	}
}

func TestAdler32_EmptyString(t *testing.T) {
	h := adler32.New()
	h.Write([]byte(""))
	got := h.Sum32()
	want := uint32(0x00000001) // Adler32 of empty string is 1
	if got != want {
		t.Errorf("adler32('') = 0x%08x, want 0x%08x", got, want)
	}
}

// ── bcrypt hash then verify ────────────────────────────────────────────────

func TestBcrypt_HashAndVerify(t *testing.T) {
	password := "mypassword123"
	cost := 4 // low cost for speed in tests

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		t.Fatalf("bcrypt hash error: %v", err)
	}

	// Verify correct password
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Errorf("bcrypt verify failed for correct password: %v", err)
	}

	// Verify wrong password
	err = bcrypt.CompareHashAndPassword(hashed, []byte("wrongpassword"))
	if err == nil {
		t.Error("bcrypt verify should fail for wrong password")
	}
}

func TestBcrypt_DifferentHashesForSamePassword(t *testing.T) {
	password := "testpassword"
	cost := 4

	h1, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		t.Fatalf("bcrypt hash error: %v", err)
	}
	h2, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		t.Fatalf("bcrypt hash error: %v", err)
	}

	if string(h1) == string(h2) {
		t.Error("bcrypt should produce different hashes for the same password (different salts)")
	}
}

package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── persistent flags ────────────────────────────────────────────────────────

var (
	jsonOut   bool
	algorithm string
)

// ── parent command ──────────────────────────────────────────────────────────

var checksumCmd = &cobra.Command{
	Use:   "checksum",
	Short: "Calculate and verify file checksums",
	Long: `Calculate and verify file checksums using various hash algorithms.

SUBCOMMANDS:

  calc     Calculate the checksum of a file
  verify   Verify a file against a known checksum

ALGORITHMS:

  md5      MD5 (128-bit) — fast, not cryptographically secure
  sha1     SHA-1 (160-bit) — deprecated for security use
  sha256   SHA-256 (256-bit) — recommended default
  sha512   SHA-512 (512-bit) — maximum security margin
  crc32    CRC-32 (32-bit) — non-cryptographic error detection

FLAGS (inherited by all subcommands):

  --algorithm/-a   Hash algorithm (default: sha256)
  --json/-j        Output result as JSON

Examples:
  openGyver checksum calc myfile.zip
  openGyver checksum calc myfile.zip --algorithm md5
  openGyver checksum verify myfile.zip abc123def456 --algorithm sha256
  openGyver checksum calc largefile.iso --algorithm sha512 --json`,
}

// ── subcommands ─────────────────────────────────────────────────────────────

var calcCmd = &cobra.Command{
	Use:   "calc <file>",
	Short: "Calculate checksum of a file",
	Long: `Calculate the checksum of a file using the specified algorithm.

The file is read in chunks for memory efficiency, making it suitable
for large files.

ALGORITHMS:
  md5, sha1, sha256 (default), sha512, crc32

Examples:
  openGyver checksum calc document.pdf
  openGyver checksum calc archive.tar.gz --algorithm md5
  openGyver checksum calc backup.iso --algorithm sha512 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runCalc,
}

var verifyCmd = &cobra.Command{
	Use:   "verify <file> <expected-hash>",
	Short: "Verify file against a known checksum",
	Long: `Verify a file's checksum against a known expected hash.

Prints "MATCH" if the computed checksum equals the expected hash,
or "MISMATCH" if they differ.

ALGORITHMS:
  md5, sha1, sha256 (default), sha512, crc32

Examples:
  openGyver checksum verify document.pdf abc123def456
  openGyver checksum verify archive.tar.gz e3b0c44298fc1c14 --algorithm sha256
  openGyver checksum verify firmware.bin 3610a686 --algorithm crc32`,
	Args: cobra.ExactArgs(2),
	RunE: runVerify,
}

// ── runners ─────────────────────────────────────────────────────────────────

func runCalc(_ *cobra.Command, args []string) error {
	filePath := args[0]

	digest, err := computeChecksum(filePath, algorithm)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"file":      filePath,
			"algorithm": algorithm,
			"checksum":  digest,
		})
	}

	fmt.Println(digest)
	return nil
}

func runVerify(_ *cobra.Command, args []string) error {
	filePath := args[0]
	expected := strings.ToLower(strings.TrimSpace(args[1]))

	digest, err := computeChecksum(filePath, algorithm)
	if err != nil {
		return err
	}

	match := digest == expected

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"file":      filePath,
			"algorithm": algorithm,
			"expected":  expected,
			"actual":    digest,
			"match":     match,
		})
	}

	if match {
		fmt.Println("MATCH")
	} else {
		fmt.Println("MISMATCH")
	}
	return nil
}

// ── helpers ─────────────────────────────────────────────────────────────────

// newHash returns a hash.Hash for the given algorithm name.
// For crc32, the returned hash implements hash.Hash (not hash.Hash32).
func newHash(algo string) (hash.Hash, error) {
	switch strings.ToLower(algo) {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	case "crc32":
		return crc32.NewIEEE(), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm %q (supported: md5, sha1, sha256, sha512, crc32)", algo)
	}
}

// computeChecksum reads a file in chunks and returns its hex digest.
func computeChecksum(filePath, algo string) (string, error) {
	h, err := newHash(algo)
	if err != nil {
		return "", err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 32*1024) // 32 KB chunks
	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	// For CRC32, format as zero-padded 8-char hex.
	if strings.ToLower(algo) == "crc32" {
		if h32, ok := h.(hash.Hash32); ok {
			return fmt.Sprintf("%08x", h32.Sum32()), nil
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// ── registration ────────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent — inherited by all subcommands.
	checksumCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output result as JSON")
	checksumCmd.PersistentFlags().StringVarP(&algorithm, "algorithm", "a", "sha256", "hash algorithm: md5, sha1, sha256, sha512, crc32")

	// Wire subcommands.
	checksumCmd.AddCommand(calcCmd)
	checksumCmd.AddCommand(verifyCmd)

	cmd.Register(checksumCmd)
}

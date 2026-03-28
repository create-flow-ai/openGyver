package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

// ── persistent flags (inherited by all subcommands) ─────────────────────────

var (
	jsonOut   bool
	uppercase bool
	filePath  string
)

// ── hmac-specific flags ─────────────────────────────────────────────────────

var (
	hmacKey       string
	hmacAlgorithm string
)

// ── bcrypt-specific flags ───────────────────────────────────────────────────

var (
	bcryptRounds int
	bcryptVerify string
)

// ── parent command ──────────────────────────────────────────────────────────

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Compute cryptographic hashes and checksums",
	Long: `Compute cryptographic hashes and checksums for strings or files.

SUBCOMMANDS:

  md5       Compute MD5 hash (128-bit)
  sha1      Compute SHA-1 hash (160-bit)
  sha256    Compute SHA-256 hash (256-bit)
  sha384    Compute SHA-384 hash (384-bit)
  sha512    Compute SHA-512 hash (512-bit)
  hmac      Compute HMAC with a secret key
  bcrypt    Hash or verify a password with bcrypt
  crc32     Compute CRC-32 checksum
  adler32   Compute Adler-32 checksum

INPUT:

  Pass the value to hash as a positional argument, or use --file/-f to hash
  the contents of a file instead.

FLAGS (inherited by all subcommands):

  --file/-f       Hash a file's contents instead of a string argument
  --json/-j       Output result as JSON
  --uppercase/-u  Output hex digest in uppercase

Examples:
  openGyver hash md5 "hello world"
  openGyver hash sha256 --file /etc/hosts
  openGyver hash sha512 "secret" --uppercase
  openGyver hash hmac "message" --key mysecret --algorithm sha512
  openGyver hash bcrypt "mypassword"
  openGyver hash bcrypt "mypassword" --verify '$2a$10$...'
  openGyver hash crc32 "hello"
  openGyver hash sha256 "hello" --json`,
}

// ── subcommands ─────────────────────────────────────────────────────────────

var md5Cmd = &cobra.Command{
	Use:   "md5 [input]",
	Short: "Compute MD5 hash (128-bit)",
	Long: `Compute the MD5 message digest of a string or file.

MD5 produces a 128-bit (16-byte) hash value, typically rendered as a
32-character hexadecimal string. Note: MD5 is cryptographically broken
and should not be used for security purposes. It is still useful for
checksums and non-security fingerprinting.

Examples:
  openGyver hash md5 "hello"
  openGyver hash md5 --file document.pdf`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processHash("md5", md5.New(), args)
	},
}

var sha1Cmd = &cobra.Command{
	Use:   "sha1 [input]",
	Short: "Compute SHA-1 hash (160-bit)",
	Long: `Compute the SHA-1 message digest of a string or file.

SHA-1 produces a 160-bit (20-byte) hash value. Note: SHA-1 is
considered weak against well-funded attackers and is deprecated for
digital signatures. Still acceptable for non-security checksums.

Examples:
  openGyver hash sha1 "hello"
  openGyver hash sha1 --file archive.tar.gz`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processHash("sha1", sha1.New(), args)
	},
}

var sha256Cmd = &cobra.Command{
	Use:   "sha256 [input]",
	Short: "Compute SHA-256 hash (256-bit)",
	Long: `Compute the SHA-256 message digest of a string or file.

SHA-256 is part of the SHA-2 family and produces a 256-bit (32-byte)
hash value. It is widely used for data integrity, digital signatures,
and blockchain applications.

Examples:
  openGyver hash sha256 "hello"
  openGyver hash sha256 --file release.zip --json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processHash("sha256", sha256.New(), args)
	},
}

var sha384Cmd = &cobra.Command{
	Use:   "sha384 [input]",
	Short: "Compute SHA-384 hash (384-bit)",
	Long: `Compute the SHA-384 message digest of a string or file.

SHA-384 is a truncated version of SHA-512 and produces a 384-bit
(48-byte) hash value. It offers a higher security margin than SHA-256
while being faster on 64-bit platforms.

Examples:
  openGyver hash sha384 "hello"
  openGyver hash sha384 --file image.iso`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processHash("sha384", sha512.New384(), args)
	},
}

var sha512Cmd = &cobra.Command{
	Use:   "sha512 [input]",
	Short: "Compute SHA-512 hash (512-bit)",
	Long: `Compute the SHA-512 message digest of a string or file.

SHA-512 is part of the SHA-2 family and produces a 512-bit (64-byte)
hash value. It is the strongest member of SHA-2, preferred when maximum
security margin is desired.

Examples:
  openGyver hash sha512 "hello"
  openGyver hash sha512 --file backup.tar.gz --uppercase`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processHash("sha512", sha512.New(), args)
	},
}

var hmacCmd = &cobra.Command{
	Use:   "hmac [input]",
	Short: "Compute HMAC with a secret key",
	Long: `Compute an HMAC (Hash-based Message Authentication Code) of a string or file.

HMAC combines a cryptographic hash function with a secret key to produce
a message authentication code. This is used to verify both data integrity
and authenticity.

FLAGS:
  --key/-k         Secret key (required)
  --algorithm/-a   Hash algorithm: md5, sha1, sha256, sha384, sha512
                   (default: sha256)

Examples:
  openGyver hash hmac "hello" --key mysecret
  openGyver hash hmac "hello" --key mysecret --algorithm sha512
  openGyver hash hmac --file payload.json --key apikey123`,
	Args: cobra.MaximumNArgs(1),
	RunE: runHMAC,
}

var bcryptCmd = &cobra.Command{
	Use:   "bcrypt [input]",
	Short: "Hash or verify a password with bcrypt",
	Long: `Hash a password using bcrypt, or verify a password against a bcrypt hash.

Bcrypt is a password-hashing function designed to be computationally
expensive, making brute-force attacks impractical. The --rounds flag
controls the cost factor (higher = slower + more secure).

MODES:
  Hash mode (default):
    Generates a bcrypt hash of the input string.

  Verify mode (--verify):
    Checks whether the input string matches the given bcrypt hash.
    Prints "match" or "no match" and exits with code 0 or 1.

FLAGS:
  --rounds/-r   Cost factor, range 4-31 (default: 10)
  --verify/-v   Bcrypt hash to verify against

Examples:
  openGyver hash bcrypt "mypassword"
  openGyver hash bcrypt "mypassword" --rounds 12
  openGyver hash bcrypt "mypassword" --verify '$2a$10$N9qo8uLOickgx2ZMRZoMye...'`,
	Args: cobra.ExactArgs(1),
	RunE: runBcrypt,
}

var crc32Cmd = &cobra.Command{
	Use:   "crc32 [input]",
	Short: "Compute CRC-32 checksum",
	Long: `Compute the CRC-32 checksum of a string or file using the IEEE polynomial.

CRC-32 is a non-cryptographic checksum used for error detection in
network transmissions and file integrity checks.

Examples:
  openGyver hash crc32 "hello"
  openGyver hash crc32 --file firmware.bin`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processChecksum("crc32", crc32.NewIEEE(), args)
	},
}

var adler32Cmd = &cobra.Command{
	Use:   "adler32 [input]",
	Short: "Compute Adler-32 checksum",
	Long: `Compute the Adler-32 checksum of a string or file.

Adler-32 is a non-cryptographic checksum that is faster but less
reliable than CRC-32. It is used by the zlib compression library.

Examples:
  openGyver hash adler32 "hello"
  openGyver hash adler32 --file data.bin`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processChecksum("adler32", adler32.New(), args)
	},
}

// ── shared helpers ──────────────────────────────────────────────────────────

// readInput returns the bytes to hash. When --file is set it reads the file;
// otherwise it uses the first positional argument.
func readInput(args []string) (input string, data []byte, err error) {
	if filePath != "" {
		f, ferr := os.Open(filePath)
		if ferr != nil {
			return "", nil, fmt.Errorf("cannot open file: %w", ferr)
		}
		defer f.Close()
		data, err = io.ReadAll(f)
		if err != nil {
			return "", nil, fmt.Errorf("cannot read file: %w", err)
		}
		return filePath, data, nil
	}
	if len(args) == 0 {
		return "", nil, fmt.Errorf("provide an input string as argument or use --file/-f")
	}
	return args[0], []byte(args[0]), nil
}

// processHash computes a cryptographic hash and prints the result.
func processHash(algorithm string, h hash.Hash, args []string) error {
	input, data, err := readInput(args)
	if err != nil {
		return err
	}
	h.Write(data)
	digest := hex.EncodeToString(h.Sum(nil))
	if uppercase {
		digest = strings.ToUpper(digest)
	}
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"algorithm": algorithm,
			"input":     input,
			"hash":      digest,
		})
	}
	fmt.Println(digest)
	return nil
}

// processChecksum computes a non-cryptographic checksum and prints the result.
// It outputs the checksum as an unsigned 32-bit hex value.
func processChecksum(algorithm string, h hash.Hash32, args []string) error {
	input, data, err := readInput(args)
	if err != nil {
		return err
	}
	h.Write(data)
	sum := h.Sum32()
	digest := fmt.Sprintf("%08x", sum)
	if uppercase {
		digest = strings.ToUpper(digest)
	}
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"algorithm": algorithm,
			"input":     input,
			"hash":      digest,
			"decimal":   sum,
		})
	}
	fmt.Println(digest)
	return nil
}

// ── hmac runner ─────────────────────────────────────────────────────────────

func runHMAC(c *cobra.Command, args []string) error {
	if hmacKey == "" {
		return fmt.Errorf("--key/-k is required for HMAC")
	}
	var hf func() hash.Hash
	algo := strings.ToLower(hmacAlgorithm)
	switch algo {
	case "md5":
		hf = md5.New
	case "sha1":
		hf = sha1.New
	case "sha256":
		hf = sha256.New
	case "sha384":
		hf = sha512.New384
	case "sha512":
		hf = sha512.New
	default:
		return fmt.Errorf("unsupported HMAC algorithm %q (supported: md5, sha1, sha256, sha384, sha512)", hmacAlgorithm)
	}

	input, data, err := readInput(args)
	if err != nil {
		return err
	}

	mac := hmac.New(hf, []byte(hmacKey))
	mac.Write(data)
	digest := hex.EncodeToString(mac.Sum(nil))
	if uppercase {
		digest = strings.ToUpper(digest)
	}

	label := "hmac-" + algo
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"algorithm": label,
			"input":     input,
			"hash":      digest,
		})
	}
	fmt.Println(digest)
	return nil
}

// ── bcrypt runner ───────────────────────────────────────────────────────────

func runBcrypt(c *cobra.Command, args []string) error {
	password := args[0]

	// Verify mode.
	if bcryptVerify != "" {
		err := bcrypt.CompareHashAndPassword([]byte(bcryptVerify), []byte(password))
		if err != nil {
			if jsonOut {
				return cmd.PrintJSON(map[string]interface{}{
					"algorithm": "bcrypt",
					"input":     password,
					"hash":      bcryptVerify,
					"match":     false,
				})
			}
			fmt.Println("no match")
			os.Exit(1)
		}
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"algorithm": "bcrypt",
				"input":     password,
				"hash":      bcryptVerify,
				"match":     true,
			})
		}
		fmt.Println("match")
		return nil
	}

	// Hash mode.
	if bcryptRounds < 4 || bcryptRounds > 31 {
		return fmt.Errorf("--rounds must be between 4 and 31, got %d", bcryptRounds)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptRounds)
	if err != nil {
		return fmt.Errorf("bcrypt error: %w", err)
	}

	digest := string(hashed)
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"algorithm": "bcrypt",
			"input":     password,
			"hash":      digest,
			"rounds":    bcryptRounds,
		})
	}
	fmt.Println(digest)
	return nil
}

// ── registration ────────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent — inherited by all subcommands.
	hashCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output result as JSON")
	hashCmd.PersistentFlags().BoolVarP(&uppercase, "uppercase", "u", false, "output hex digest in uppercase")
	hashCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "hash a file's contents instead of a string argument")

	// HMAC-specific flags.
	hmacCmd.Flags().StringVarP(&hmacKey, "key", "k", "", "secret key for HMAC (required)")
	hmacCmd.Flags().StringVarP(&hmacAlgorithm, "algorithm", "a", "sha256", "hash algorithm: md5, sha1, sha256, sha384, sha512")

	// Bcrypt-specific flags.
	bcryptCmd.Flags().IntVarP(&bcryptRounds, "rounds", "r", 10, "bcrypt cost factor (4-31)")
	bcryptCmd.Flags().StringVarP(&bcryptVerify, "verify", "v", "", "bcrypt hash to verify the input against")

	// Wire subcommands.
	hashCmd.AddCommand(md5Cmd)
	hashCmd.AddCommand(sha1Cmd)
	hashCmd.AddCommand(sha256Cmd)
	hashCmd.AddCommand(sha384Cmd)
	hashCmd.AddCommand(sha512Cmd)
	hashCmd.AddCommand(hmacCmd)
	hashCmd.AddCommand(bcryptCmd)
	hashCmd.AddCommand(crc32Cmd)
	hashCmd.AddCommand(adler32Cmd)

	cmd.Register(hashCmd)
}

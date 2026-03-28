package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/elliptic"
	gocrypto "crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh"
)

// ──────────────────────────────────────────────────────────────
// Shared flags
// ──────────────────────────────────────────────────────────────

var jsonOut bool

// ──────────────────────────────────────────────────────────────
// Parent command
// ──────────────────────────────────────────────────────────────

var cryptoCmd = &cobra.Command{
	Use:   "crypto",
	Short: "Cryptographic tools — encrypt, generate keys, certificates",
	Long: `A collection of cryptographic utilities.

SUBCOMMANDS:

  aes       AES-256-GCM encrypt / decrypt
  rsa       Generate RSA key pairs (PEM)
  sshkey    Generate SSH key pairs (OpenSSH format)
  cert      Generate self-signed TLS certificates
  csr       Generate Certificate Signing Requests

All subcommands support --json / -j for machine-readable output.

Examples:
  openGyver crypto aes "hello world" --key mypassphrase
  openGyver crypto aes "BASE64CIPHER" --key mypassphrase --decrypt
  openGyver crypto rsa --bits 4096
  openGyver crypto sshkey --type ed25519 --comment "me@host"
  openGyver crypto cert --cn example.com --days 730
  openGyver crypto csr --cn example.com --org "Acme Inc"`,
}

// ──────────────────────────────────────────────────────────────
// aes subcommand
// ──────────────────────────────────────────────────────────────

var (
	aesKey     string
	aesDecrypt bool
)

var aesCmd = &cobra.Command{
	Use:   "aes <plaintext | ciphertext>",
	Short: "AES-256-GCM encrypt or decrypt",
	Long: `Encrypt or decrypt data using AES-256-GCM.

The --key flag is required and accepts either:
  - A 64-character hex string (raw 256-bit key)
  - Any other string, treated as a passphrase (key derived via PBKDF2)

Encryption output is base64-encoded and includes the nonce (first 12 bytes)
and, when using a passphrase, a 16-byte salt prefix.

Use --decrypt / -d to reverse the operation.

Examples:
  # Encrypt with a passphrase
  openGyver crypto aes "secret message" --key "my passphrase"

  # Decrypt
  openGyver crypto aes "BASE64..." --key "my passphrase" --decrypt

  # Encrypt with a raw hex key (64 hex chars = 32 bytes)
  openGyver crypto aes "hello" --key 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

  # JSON output
  openGyver crypto aes "hello" --key pass --json`,
	Args: cobra.ExactArgs(1),
	RunE: runAES,
}

func isHexKey(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func deriveKey(passphrase, salt []byte) []byte {
	return pbkdf2.Key(passphrase, salt, 600_000, 32, sha256.New)
}

func runAES(c *cobra.Command, args []string) error {
	input := args[0]

	if aesDecrypt {
		return aesDecryptFlow(input)
	}
	return aesEncryptFlow(input)
}

func aesEncryptFlow(plaintext string) error {
	var keyBytes []byte
	var prefix []byte // prepended to ciphertext (salt for passphrase mode)

	if isHexKey(aesKey) {
		var err error
		keyBytes, err = hex.DecodeString(aesKey)
		if err != nil {
			return fmt.Errorf("invalid hex key: %w", err)
		}
	} else {
		// Passphrase mode — generate random salt.
		salt := make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			return fmt.Errorf("generating salt: %w", err)
		}
		keyBytes = deriveKey([]byte(aesKey), salt)
		prefix = salt
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("creating GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generating nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	blob := append(prefix, ciphertext...)
	encoded := base64.StdEncoding.EncodeToString(blob)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":     plaintext,
			"output":    encoded,
			"algorithm": "aes-256-gcm",
		})
	}
	fmt.Println(encoded)
	return nil
}

func aesDecryptFlow(encoded string) error {
	blob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("invalid base64 input: %w", err)
	}

	var keyBytes []byte
	var ciphertext []byte

	if isHexKey(aesKey) {
		keyBytes, err = hex.DecodeString(aesKey)
		if err != nil {
			return fmt.Errorf("invalid hex key: %w", err)
		}
		ciphertext = blob
	} else {
		// Passphrase mode — first 16 bytes are salt.
		if len(blob) < 16 {
			return fmt.Errorf("ciphertext too short (missing salt)")
		}
		salt := blob[:16]
		ciphertext = blob[16:]
		keyBytes = deriveKey([]byte(aesKey), salt)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("creating GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short (missing nonce)")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed (wrong key or corrupted data): %w", err)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":     encoded,
			"output":    string(plaintext),
			"algorithm": "aes-256-gcm",
		})
	}
	fmt.Println(string(plaintext))
	return nil
}

// ──────────────────────────────────────────────────────────────
// rsa subcommand
// ──────────────────────────────────────────────────────────────

var (
	rsaBits      int
	rsaOutputDir string
)

var rsaCmd = &cobra.Command{
	Use:   "rsa",
	Short: "Generate an RSA key pair (PEM)",
	Long: `Generate an RSA private/public key pair in PEM format.

The --bits flag controls key size (default 2048). Common values: 2048, 3072, 4096.

Use --output-dir to write the keys to files (private.pem and public.pem)
instead of printing to stdout.

Examples:
  openGyver crypto rsa
  openGyver crypto rsa --bits 4096
  openGyver crypto rsa --output-dir ./keys
  openGyver crypto rsa --bits 4096 --json`,
	Args: cobra.NoArgs,
	RunE: runRSA,
}

func runRSA(c *cobra.Command, args []string) error {
	privKey, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return fmt.Errorf("generating RSA key: %w", err)
	}

	privDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("marshalling private key: %w", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privDER,
	})

	pubDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshalling public key: %w", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	})

	if rsaOutputDir != "" {
		if err := writeKeyFiles(rsaOutputDir, "private.pem", privPEM, "public.pem", pubPEM); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Wrote %s and %s\n",
			filepath.Join(rsaOutputDir, "private.pem"),
			filepath.Join(rsaOutputDir, "public.pem"))
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"algorithm":   "RSA",
			"bits":        rsaBits,
			"private_key": string(privPEM),
			"public_key":  string(pubPEM),
		})
	}

	if rsaOutputDir == "" {
		fmt.Print(string(privPEM))
		fmt.Print(string(pubPEM))
	}
	return nil
}

// ──────────────────────────────────────────────────────────────
// sshkey subcommand
// ──────────────────────────────────────────────────────────────

var (
	sshKeyType    string
	sshKeyComment string
)

var sshKeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "Generate an SSH key pair (OpenSSH format)",
	Long: `Generate an SSH key pair in OpenSSH format.

Supported key types:
  ed25519   (default) — fast, small, modern
  rsa       — 4096-bit RSA

The --comment flag sets the key comment (e.g. user@host).

Examples:
  openGyver crypto sshkey
  openGyver crypto sshkey --type rsa
  openGyver crypto sshkey --type ed25519 --comment "deploy@prod"
  openGyver crypto sshkey --json`,
	Args: cobra.NoArgs,
	RunE: runSSHKey,
}

func runSSHKey(c *cobra.Command, args []string) error {
	var (
		privPEM  []byte
		pubBytes []byte
	)

	switch strings.ToLower(sshKeyType) {
	case "ed25519":
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("generating ed25519 key: %w", err)
		}

		privBlock, err := ssh.MarshalPrivateKey(gocrypto.PrivateKey(privKey), sshKeyComment)
		if err != nil {
			return fmt.Errorf("marshalling private key: %w", err)
		}
		privPEM = pem.EncodeToMemory(privBlock)

		sshPub, err := ssh.NewPublicKey(pubKey)
		if err != nil {
			return fmt.Errorf("converting public key: %w", err)
		}
		pubBytes = ssh.MarshalAuthorizedKey(sshPub)
		if sshKeyComment != "" {
			pubBytes = []byte(strings.TrimSpace(string(pubBytes)) + " " + sshKeyComment + "\n")
		}

	case "rsa":
		privKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return fmt.Errorf("generating RSA key: %w", err)
		}

		privBlock, err := ssh.MarshalPrivateKey(gocrypto.PrivateKey(privKey), sshKeyComment)
		if err != nil {
			return fmt.Errorf("marshalling private key: %w", err)
		}
		privPEM = pem.EncodeToMemory(privBlock)

		sshPub, err := ssh.NewPublicKey(&privKey.PublicKey)
		if err != nil {
			return fmt.Errorf("converting public key: %w", err)
		}
		pubBytes = ssh.MarshalAuthorizedKey(sshPub)
		if sshKeyComment != "" {
			pubBytes = []byte(strings.TrimSpace(string(pubBytes)) + " " + sshKeyComment + "\n")
		}

	default:
		return fmt.Errorf("unsupported key type %q (supported: ed25519, rsa)", sshKeyType)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":        sshKeyType,
			"comment":     sshKeyComment,
			"private_key": string(privPEM),
			"public_key":  strings.TrimSpace(string(pubBytes)),
		})
	}

	fmt.Print(string(privPEM))
	fmt.Print(string(pubBytes))
	return nil
}

// ──────────────────────────────────────────────────────────────
// cert subcommand
// ──────────────────────────────────────────────────────────────

var (
	certCN        string
	certDays      int
	certOutputDir string
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Generate a self-signed TLS certificate",
	Long: `Generate a self-signed X.509 certificate and private key in PEM format.

The certificate is signed with ECDSA P-256 for fast generation and small size.

Use --output-dir to write cert.pem and key.pem to disk.

Examples:
  openGyver crypto cert --cn example.com
  openGyver crypto cert --cn localhost --days 30
  openGyver crypto cert --cn "*.example.com" --days 730 --output-dir ./certs
  openGyver crypto cert --cn myapp.local --json`,
	Args: cobra.NoArgs,
	RunE: runCert,
}

func runCert(c *cobra.Command, args []string) error {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return fmt.Errorf("generating serial: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(certDays) * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: certCN,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{certCN},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return fmt.Errorf("creating certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyDER, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("marshalling private key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyDER,
	})

	if certOutputDir != "" {
		if err := writeKeyFiles(certOutputDir, "cert.pem", certPEM, "key.pem", keyPEM); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Wrote %s and %s\n",
			filepath.Join(certOutputDir, "cert.pem"),
			filepath.Join(certOutputDir, "key.pem"))
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"common_name": certCN,
			"not_before":  notBefore.UTC().Format(time.RFC3339),
			"not_after":   notAfter.UTC().Format(time.RFC3339),
			"certificate": string(certPEM),
			"private_key": string(keyPEM),
		})
	}

	if certOutputDir == "" {
		fmt.Print(string(certPEM))
		fmt.Print(string(keyPEM))
	}
	return nil
}

// ──────────────────────────────────────────────────────────────
// csr subcommand
// ──────────────────────────────────────────────────────────────

var (
	csrCN      string
	csrOrg     string
	csrCountry string
	csrOutput  string
)

var csrCmd = &cobra.Command{
	Use:   "csr",
	Short: "Generate a Certificate Signing Request (CSR)",
	Long: `Generate a PEM-encoded Certificate Signing Request.

A new ECDSA P-256 private key is generated alongside the CSR.

Use --output to write the CSR to a file. The private key is always
printed to stdout (or included in JSON output).

Examples:
  openGyver crypto csr --cn example.com
  openGyver crypto csr --cn example.com --org "Acme Inc" --country US
  openGyver crypto csr --cn example.com --output request.pem
  openGyver crypto csr --cn example.com --json`,
	Args: cobra.NoArgs,
	RunE: runCSR,
}

func runCSR(c *cobra.Command, args []string) error {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	subject := pkix.Name{
		CommonName: csrCN,
	}
	if csrOrg != "" {
		subject.Organization = []string{csrOrg}
	}
	if csrCountry != "" {
		subject.Country = []string{csrCountry}
	}

	template := x509.CertificateRequest{
		Subject:            subject,
		DNSNames:           []string{csrCN},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	if err != nil {
		return fmt.Errorf("creating CSR: %w", err)
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	keyDER, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("marshalling private key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyDER,
	})

	if csrOutput != "" {
		if err := os.WriteFile(csrOutput, csrPEM, 0644); err != nil {
			return fmt.Errorf("writing CSR file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote %s\n", csrOutput)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"common_name": csrCN,
			"org":         csrOrg,
			"country":     csrCountry,
			"csr":         string(csrPEM),
			"private_key": string(keyPEM),
		})
	}

	if csrOutput == "" {
		fmt.Print(string(csrPEM))
	}
	fmt.Print(string(keyPEM))
	return nil
}

// ──────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────

func writeKeyFiles(dir, name1 string, data1 []byte, name2 string, data2 []byte) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name1), data1, 0600); err != nil {
		return fmt.Errorf("writing %s: %w", name1, err)
	}
	if err := os.WriteFile(filepath.Join(dir, name2), data2, 0600); err != nil {
		return fmt.Errorf("writing %s: %w", name2, err)
	}
	return nil
}

// ──────────────────────────────────────────────────────────────
// Registration
// ──────────────────────────────────────────────────────────────

func init() {
	// Parent persistent flags.
	cryptoCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	// aes flags
	aesCmd.Flags().StringVar(&aesKey, "key", "", "encryption key (hex string or passphrase)")
	_ = aesCmd.MarkFlagRequired("key")
	aesCmd.Flags().BoolVarP(&aesDecrypt, "decrypt", "d", false, "decrypt instead of encrypt")

	// rsa flags
	rsaCmd.Flags().IntVar(&rsaBits, "bits", 2048, "RSA key size in bits (2048, 3072, 4096)")
	rsaCmd.Flags().StringVar(&rsaOutputDir, "output-dir", "", "directory to write key files")

	// sshkey flags
	sshKeyCmd.Flags().StringVar(&sshKeyType, "type", "ed25519", "key type: ed25519 or rsa")
	sshKeyCmd.Flags().StringVar(&sshKeyComment, "comment", "", "key comment (e.g. user@host)")

	// cert flags
	certCmd.Flags().StringVar(&certCN, "cn", "", "Common Name (required)")
	_ = certCmd.MarkFlagRequired("cn")
	certCmd.Flags().IntVar(&certDays, "days", 365, "certificate validity in days")
	certCmd.Flags().StringVar(&certOutputDir, "output-dir", "", "directory to write cert and key files")

	// csr flags
	csrCmd.Flags().StringVar(&csrCN, "cn", "", "Common Name (required)")
	_ = csrCmd.MarkFlagRequired("cn")
	csrCmd.Flags().StringVar(&csrOrg, "org", "", "Organization name")
	csrCmd.Flags().StringVar(&csrCountry, "country", "", "Country code (e.g. US)")
	csrCmd.Flags().StringVar(&csrOutput, "output", "", "file path to write the CSR")

	// Wire subcommands.
	cryptoCmd.AddCommand(aesCmd)
	cryptoCmd.AddCommand(rsaCmd)
	cryptoCmd.AddCommand(sshKeyCmd)
	cryptoCmd.AddCommand(certCmd)
	cryptoCmd.AddCommand(csrCmd)

	cmd.Register(cryptoCmd)
}

package encode

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"math/big"
	"net/url"
	"os"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/net/idna"
)

// ── flags ───────────────────────────────────────────────────────────────────

var (
	decode  bool
	jsonOut bool
	file    string
)

// ── parent command ──────────────────────────────────────────────────────────

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode and decode text in various formats",
	Long: `Encode and decode text using common encoding schemes.

SUBCOMMANDS:

  base64    Base64 encode/decode
  base32    Base32 encode/decode
  base58    Base58 (Bitcoin alphabet) encode/decode
  url       URL percent-encoding encode/decode
  html      HTML entity encode/decode
  hex       Hex encode/decode
  binary    Binary (0/1) encode/decode
  rot13     ROT13 cipher (symmetric — no --decode needed)
  morse     Morse code encode/decode
  punycode  Punycode (IDNA) encode/decode
  jwt       Decode a JWT token payload (no verification)

Each subcommand accepts input as the first argument, or via --file/-f.
By default each subcommand encodes; pass --decode/-d to reverse.

FLAGS:

  --decode, -d   Decode instead of encode (on each subcommand)
  --json,   -j   Output structured JSON: {"input","output","encoding"}
  --file,   -f   Read input from a file instead of an argument

Examples:
  openGyver encode base64 "hello world"
  openGyver encode base64 -d "aGVsbG8gd29ybGQ="
  openGyver encode url "hello world"
  openGyver encode hex "cafe" -d
  openGyver encode rot13 "Hello"
  openGyver encode morse "SOS"
  openGyver encode jwt "eyJhbGciOi..."
  openGyver encode base64 --file input.txt
  openGyver encode base64 "hello" --json`,
}

// ── subcommands ─────────────────────────────────────────────────────────────

var base64Cmd = &cobra.Command{
	Use:   "base64 [input]",
	Short: "Base64 encode or decode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "base64", func(in string, dec bool) (string, error) {
			if dec {
				b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(in))
				if err != nil {
					return "", fmt.Errorf("base64 decode: %w", err)
				}
				return string(b), nil
			}
			return base64.StdEncoding.EncodeToString([]byte(in)), nil
		})
	},
}

var base32Cmd = &cobra.Command{
	Use:   "base32 [input]",
	Short: "Base32 encode or decode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "base32", func(in string, dec bool) (string, error) {
			if dec {
				b, err := base32.StdEncoding.DecodeString(strings.TrimSpace(in))
				if err != nil {
					return "", fmt.Errorf("base32 decode: %w", err)
				}
				return string(b), nil
			}
			return base32.StdEncoding.EncodeToString([]byte(in)), nil
		})
	},
}

var base58Cmd = &cobra.Command{
	Use:   "base58 [input]",
	Short: "Base58 (Bitcoin alphabet) encode or decode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "base58", func(in string, dec bool) (string, error) {
			if dec {
				b, err := base58Decode(strings.TrimSpace(in))
				if err != nil {
					return "", err
				}
				return string(b), nil
			}
			return base58Encode([]byte(in)), nil
		})
	},
}

var urlCmd = &cobra.Command{
	Use:   "url [input]",
	Short: "URL percent-encoding encode or decode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "url", func(in string, dec bool) (string, error) {
			if dec {
				s, err := url.QueryUnescape(strings.TrimSpace(in))
				if err != nil {
					return "", fmt.Errorf("url decode: %w", err)
				}
				return s, nil
			}
			return url.QueryEscape(in), nil
		})
	},
}

var htmlCmd = &cobra.Command{
	Use:   "html [input]",
	Short: "HTML entity encode or decode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "html", func(in string, dec bool) (string, error) {
			if dec {
				return html.UnescapeString(in), nil
			}
			return html.EscapeString(in), nil
		})
	},
}

var hexCmd = &cobra.Command{
	Use:   "hex [input]",
	Short: "Hex encode or decode text",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "hex", func(in string, dec bool) (string, error) {
			if dec {
				b, err := hex.DecodeString(strings.TrimSpace(in))
				if err != nil {
					return "", fmt.Errorf("hex decode: %w", err)
				}
				return string(b), nil
			}
			return hex.EncodeToString([]byte(in)), nil
		})
	},
}

var binaryCmd = &cobra.Command{
	Use:   "binary [input]",
	Short: "Binary (0/1) encode or decode text",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "binary", func(in string, dec bool) (string, error) {
			if dec {
				return binaryDecode(strings.TrimSpace(in))
			}
			return binaryEncode(in), nil
		})
	},
}

var rot13Cmd = &cobra.Command{
	Use:   "rot13 [input]",
	Short: "ROT13 cipher (symmetric)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "rot13", func(in string, _ bool) (string, error) {
			return rot13(in), nil
		})
	},
}

var morseCmd = &cobra.Command{
	Use:   "morse [input]",
	Short: "Morse code encode or decode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "morse", func(in string, dec bool) (string, error) {
			if dec {
				return morseDecode(in)
			}
			return morseEncode(in)
		})
	},
}

var punycodeCmd = &cobra.Command{
	Use:   "punycode [input]",
	Short: "Punycode (IDNA) encode or decode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "punycode", func(in string, dec bool) (string, error) {
			if dec {
				s, err := idna.ToUnicode(strings.TrimSpace(in))
				if err != nil {
					return "", fmt.Errorf("punycode decode: %w", err)
				}
				return s, nil
			}
			s, err := idna.ToASCII(in)
			if err != nil {
				return "", fmt.Errorf("punycode encode: %w", err)
			}
			return s, nil
		})
	},
}

var jwtCmd = &cobra.Command{
	Use:   "jwt [token]",
	Short: "Decode a JWT token payload (no verification)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return processEncode(args, "jwt", func(in string, _ bool) (string, error) {
			return jwtDecode(strings.TrimSpace(in))
		})
	},
}

// ── shared processing ───────────────────────────────────────────────────────

// encodeFn performs the actual encode/decode. The bool is true when decoding.
type encodeFn func(input string, decode bool) (string, error)

// processEncode resolves the input (argument or --file), calls fn, and handles
// --json output.
func processEncode(args []string, encoding string, fn encodeFn) error {
	input, err := resolveInput(args)
	if err != nil {
		return err
	}

	output, err := fn(input, decode)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":    input,
			"output":   output,
			"encoding": encoding,
		})
	}

	fmt.Println(output)
	return nil
}

// resolveInput returns the input string from either the first argument or the
// --file flag.
func resolveInput(args []string) (string, error) {
	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}
	if len(args) == 0 {
		return "", fmt.Errorf("provide input as an argument or use --file/-f")
	}
	return args[0], nil
}

// ── base58 (Bitcoin alphabet) ───────────────────────────────────────────────

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func base58Encode(data []byte) string {
	n := new(big.Int).SetBytes(data)
	zero := big.NewInt(0)
	base := big.NewInt(58)
	mod := new(big.Int)

	var result []byte
	for n.Cmp(zero) > 0 {
		n.DivMod(n, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}
	// preserve leading zero bytes
	for _, b := range data {
		if b != 0 {
			break
		}
		result = append(result, base58Alphabet[0])
	}
	// reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

func base58Decode(s string) ([]byte, error) {
	n := big.NewInt(0)
	base := big.NewInt(58)
	for _, c := range []byte(s) {
		idx := strings.IndexByte(base58Alphabet, c)
		if idx < 0 {
			return nil, fmt.Errorf("base58 decode: invalid character %q", c)
		}
		n.Mul(n, base)
		n.Add(n, big.NewInt(int64(idx)))
	}
	result := n.Bytes()
	// restore leading zero bytes
	for _, c := range []byte(s) {
		if c != base58Alphabet[0] {
			break
		}
		result = append([]byte{0}, result...)
	}
	return result, nil
}

// ── binary ──────────────────────────────────────────────────────────────────

func binaryEncode(s string) string {
	parts := make([]string, len(s))
	for i, b := range []byte(s) {
		parts[i] = fmt.Sprintf("%08b", b)
	}
	return strings.Join(parts, " ")
}

func binaryDecode(s string) (string, error) {
	fields := strings.Fields(s)
	var out []byte
	for _, f := range fields {
		var b byte
		for _, ch := range f {
			if ch != '0' && ch != '1' {
				return "", fmt.Errorf("binary decode: invalid character %q", ch)
			}
			b = b<<1 | byte(ch-'0')
		}
		out = append(out, b)
	}
	return string(out), nil
}

// ── rot13 ───────────────────────────────────────────────────────────────────

func rot13(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		default:
			return r
		}
	}, s)
}

// ── morse ───────────────────────────────────────────────────────────────────

var charToMorse = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".",
	'F': "..-.", 'G': "--.", 'H': "....", 'I': "..", 'J': ".---",
	'K': "-.-", 'L': ".-..", 'M': "--", 'N': "-.", 'O': "---",
	'P': ".--.", 'Q': "--.-", 'R': ".-.", 'S': "...", 'T': "-",
	'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-", 'Y': "-.--",
	'Z': "--..",
	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-",
	'5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",
	'.': ".-.-.-", ',': "--..--", '?': "..--..", '!': "-.-.--",
	'/': "-..-.", '(': "-.--.", ')': "-.--.-", '&': ".-...",
	':': "---...", ';': "-.-.-.", '=': "-...-", '+': ".-.-.",
	'-': "-....-", '_': "..--.-", '"': ".-..-.", '\'': ".----.",
	'@': ".--.-.",
}

var morseToChar map[string]rune

func initMorseLookup() {
	morseToChar = make(map[string]rune, len(charToMorse))
	for k, v := range charToMorse {
		morseToChar[v] = k
	}
}

func morseEncode(s string) (string, error) {
	var parts []string
	for _, r := range strings.ToUpper(s) {
		if r == ' ' {
			parts = append(parts, "/")
			continue
		}
		code, ok := charToMorse[r]
		if !ok {
			return "", fmt.Errorf("morse encode: unsupported character %q", r)
		}
		parts = append(parts, code)
	}
	return strings.Join(parts, " "), nil
}

func morseDecode(s string) (string, error) {
	if morseToChar == nil {
		initMorseLookup()
	}
	words := strings.Split(s, "/")
	var out strings.Builder
	for i, word := range words {
		if i > 0 {
			out.WriteRune(' ')
		}
		codes := strings.Fields(strings.TrimSpace(word))
		for _, code := range codes {
			ch, ok := morseToChar[code]
			if !ok {
				return "", fmt.Errorf("morse decode: unknown code %q", code)
			}
			out.WriteRune(ch)
		}
	}
	return out.String(), nil
}

// ── jwt ─────────────────────────────────────────────────────────────────────

func jwtDecode(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("jwt: expected 3 dot-separated parts, got %d", len(parts))
	}
	payload := parts[1]
	// JWT uses base64url without padding
	if m := len(payload) % 4; m != 0 {
		payload += strings.Repeat("=", 4-m)
	}
	b, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("jwt payload decode: %w", err)
	}
	// Pretty-print JSON
	var obj interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		// Not valid JSON — return raw
		return string(b), nil
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return string(b), nil
	}
	return string(pretty), nil
}

// ── init ────────────────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent (inherited by all subcommands).
	encodeCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON {input, output, encoding}")
	encodeCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "read input from a file instead of an argument")

	// --decode/-d on each subcommand that uses it.
	for _, sc := range []*cobra.Command{
		base64Cmd, base32Cmd, base58Cmd, urlCmd, htmlCmd, hexCmd, binaryCmd, morseCmd, punycodeCmd,
	} {
		sc.Flags().BoolVarP(&decode, "decode", "d", false, "decode instead of encode")
	}

	// Register all subcommands.
	encodeCmd.AddCommand(
		base64Cmd,
		base32Cmd,
		base58Cmd,
		urlCmd,
		htmlCmd,
		hexCmd,
		binaryCmd,
		rot13Cmd,
		morseCmd,
		punycodeCmd,
		jwtCmd,
	)

	cmd.Register(encodeCmd)
}

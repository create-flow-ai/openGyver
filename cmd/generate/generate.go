package generate

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// Shared state
// ---------------------------------------------------------------------------

var jsonOut bool

// ---------------------------------------------------------------------------
// Parent command
// ---------------------------------------------------------------------------

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Random generators — passwords, keys, IDs, OTP secrets, and more",
	Long: `Generate cryptographically random passwords, passphrases, API keys,
secrets, OTP tokens, and various ID formats.

All randomness is sourced from crypto/rand.

SUBCOMMANDS:

  password     Random password with configurable character classes
  passphrase   Random passphrase from a word list
  string       Random string in a chosen character set
  apikey       API key with optional prefix
  secret       Hex-encoded secret key
  otp          TOTP secret (base32) with otpauth:// URI
  nanoid       Nano ID (URL-friendly unique string)
  snowflake    Snowflake-style 64-bit timestamp-based ID
  shortid      Short random ID

Examples:
  openGyver generate password
  openGyver generate password --length 24 --no-special --count 5
  openGyver generate passphrase --words 6 --separator "."
  openGyver generate string --charset hex --length 64
  openGyver generate apikey --prefix sk_live_
  openGyver generate secret --length 32
  openGyver generate otp --issuer MyApp --account user@example.com
  openGyver generate nanoid --length 12
  openGyver generate snowflake
  openGyver generate shortid --length 10
  openGyver generate password --json`,
}

// ---------------------------------------------------------------------------
// password
// ---------------------------------------------------------------------------

var (
	pwLength    int
	pwNoUpper   bool
	pwNoLower   bool
	pwNoDigits  bool
	pwNoSpecial bool
	pwCount     int
)

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Generate random passwords",
	Long: `Generate cryptographically random passwords.

By default passwords include uppercase, lowercase, digits, and special
characters. Disable character classes with --no-upper, --no-lower,
--no-digits, --no-special.

Examples:
  openGyver generate password
  openGyver generate password --length 24
  openGyver generate password --no-special
  openGyver generate password --no-upper --no-digits
  openGyver generate password --count 10
  openGyver generate password --length 32 --json`,
	Args: cobra.NoArgs,
	RunE: runPassword,
}

func runPassword(_ *cobra.Command, _ []string) error {
	const (
		upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lower   = "abcdefghijklmnopqrstuvwxyz"
		digits  = "0123456789"
		special = "!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
	)

	var charset string
	if !pwNoUpper {
		charset += upper
	}
	if !pwNoLower {
		charset += lower
	}
	if !pwNoDigits {
		charset += digits
	}
	if !pwNoSpecial {
		charset += special
	}
	if charset == "" {
		return fmt.Errorf("all character classes disabled; at least one must be enabled")
	}

	results := make([]string, 0, pwCount)
	for i := 0; i < pwCount; i++ {
		pw, err := randStringFromCharset(pwLength, charset)
		if err != nil {
			return err
		}
		results = append(results, pw)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":   "password",
			"length": pwLength,
			"count":  pwCount,
			"values": results,
		})
	}
	for _, v := range results {
		fmt.Println(v)
	}
	return nil
}

// ---------------------------------------------------------------------------
// passphrase
// ---------------------------------------------------------------------------

var (
	ppWords     int
	ppSeparator string
	ppCount     int
)

var passphraseCmd = &cobra.Command{
	Use:   "passphrase",
	Short: "Generate random passphrases from a word list",
	Long: `Generate random passphrases by picking words from the EFF short
diceware word list (~1296 words). Each word is chosen independently
using crypto/rand.

Examples:
  openGyver generate passphrase
  openGyver generate passphrase --words 6
  openGyver generate passphrase --separator "."
  openGyver generate passphrase --words 5 --count 3 --json`,
	Args: cobra.NoArgs,
	RunE: runPassphrase,
}

func runPassphrase(_ *cobra.Command, _ []string) error {
	results := make([]string, 0, ppCount)
	for i := 0; i < ppCount; i++ {
		words := make([]string, 0, ppWords)
		for j := 0; j < ppWords; j++ {
			idx, err := randInt(len(effWordList))
			if err != nil {
				return err
			}
			words = append(words, effWordList[idx])
		}
		results = append(results, strings.Join(words, ppSeparator))
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":      "passphrase",
			"words":     ppWords,
			"separator": ppSeparator,
			"count":     ppCount,
			"values":    results,
		})
	}
	for _, v := range results {
		fmt.Println(v)
	}
	return nil
}

// ---------------------------------------------------------------------------
// string
// ---------------------------------------------------------------------------

var (
	strLength  int
	strCharset string
	strCustom  string
	strCount   int
)

var stringCmd = &cobra.Command{
	Use:   "string",
	Short: "Generate random strings in various character sets",
	Long: `Generate cryptographically random strings.

CHARSETS:

  alpha         A-Z a-z
  alphanumeric  A-Z a-z 0-9 (default)
  hex           0-9 a-f
  base64        A-Z a-z 0-9 + / =
  custom        Supply your own alphabet via --custom

Examples:
  openGyver generate string
  openGyver generate string --length 64
  openGyver generate string --charset hex --length 32
  openGyver generate string --charset custom --custom "01" --length 16
  openGyver generate string --count 5 --json`,
	Args: cobra.NoArgs,
	RunE: runString,
}

func runString(_ *cobra.Command, _ []string) error {
	var charset string
	switch strCharset {
	case "alpha":
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "alphanumeric":
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	case "hex":
		charset = "0123456789abcdef"
	case "base64":
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	case "custom":
		if strCustom == "" {
			return fmt.Errorf("--custom alphabet must not be empty when --charset=custom")
		}
		charset = strCustom
	default:
		return fmt.Errorf("unknown charset %q; choose alpha, alphanumeric, hex, base64, or custom", strCharset)
	}

	results := make([]string, 0, strCount)
	for i := 0; i < strCount; i++ {
		s, err := randStringFromCharset(strLength, charset)
		if err != nil {
			return err
		}
		results = append(results, s)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":    "string",
			"charset": strCharset,
			"length":  strLength,
			"count":   strCount,
			"values":  results,
		})
	}
	for _, v := range results {
		fmt.Println(v)
	}
	return nil
}

// ---------------------------------------------------------------------------
// apikey
// ---------------------------------------------------------------------------

var (
	akPrefix string
	akLength int
)

var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Generate API keys with optional prefix",
	Long: `Generate a URL-safe random API key. Optionally prepend a prefix
(e.g. "sk_live_", "pk_test_").

The random portion is base62 (A-Z a-z 0-9). The --length flag controls
the length of the random part, not including the prefix.

Examples:
  openGyver generate apikey
  openGyver generate apikey --prefix sk_live_
  openGyver generate apikey --length 48 --prefix pk_test_
  openGyver generate apikey --json`,
	Args: cobra.NoArgs,
	RunE: runAPIKey,
}

func runAPIKey(_ *cobra.Command, _ []string) error {
	const base62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	body, err := randStringFromCharset(akLength, base62)
	if err != nil {
		return err
	}
	key := akPrefix + body

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":   "apikey",
			"prefix": akPrefix,
			"length": akLength,
			"key":    key,
		})
	}
	fmt.Println(key)
	return nil
}

// ---------------------------------------------------------------------------
// secret
// ---------------------------------------------------------------------------

var secretLength int

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Generate hex-encoded secret keys",
	Long: `Generate a cryptographically random secret key and output it as a
hex-encoded string. The --length flag specifies the number of random
bytes (the hex output will be twice as long).

Examples:
  openGyver generate secret
  openGyver generate secret --length 32
  openGyver generate secret --json`,
	Args: cobra.NoArgs,
	RunE: runSecret,
}

func runSecret(_ *cobra.Command, _ []string) error {
	buf := make([]byte, secretLength)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Errorf("crypto/rand: %w", err)
	}
	secret := hex.EncodeToString(buf)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":       "secret",
			"bytes":      secretLength,
			"hex_length": secretLength * 2,
			"secret":     secret,
		})
	}
	fmt.Println(secret)
	return nil
}

// ---------------------------------------------------------------------------
// otp
// ---------------------------------------------------------------------------

var (
	otpIssuer  string
	otpAccount string
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Generate TOTP secrets with otpauth:// URIs",
	Long: `Generate a random TOTP secret encoded in base32, along with an
otpauth:// URI suitable for QR codes and authenticator apps.

The secret is 20 bytes (160 bits) of crypto/rand entropy, matching
the recommended length for HMAC-SHA1 TOTP.

Examples:
  openGyver generate otp --issuer MyApp --account user@example.com
  openGyver generate otp --issuer GitHub --account octocat
  openGyver generate otp --json`,
	Args: cobra.NoArgs,
	RunE: runOTP,
}

func runOTP(_ *cobra.Command, _ []string) error {
	// 20 bytes = 160 bits, standard for TOTP HMAC-SHA1
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Errorf("crypto/rand: %w", err)
	}
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf)

	// Build otpauth:// URI
	label := otpAccount
	if otpIssuer != "" && otpAccount != "" {
		label = otpIssuer + ":" + otpAccount
	}
	uri := fmt.Sprintf("otpauth://totp/%s?secret=%s&algorithm=SHA1&digits=6&period=30",
		url.PathEscape(label),
		secret,
	)
	if otpIssuer != "" {
		uri += "&issuer=" + url.QueryEscape(otpIssuer)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":    "otp",
			"secret":  secret,
			"uri":     uri,
			"issuer":  otpIssuer,
			"account": otpAccount,
		})
	}
	fmt.Printf("Secret:  %s\n", secret)
	fmt.Printf("URI:     %s\n", uri)
	return nil
}

// ---------------------------------------------------------------------------
// nanoid
// ---------------------------------------------------------------------------

var (
	nanoidLength   int
	nanoidAlphabet string
)

var nanoidCmd = &cobra.Command{
	Use:   "nanoid",
	Short: "Generate Nano IDs (URL-friendly unique strings)",
	Long: `Generate compact, URL-friendly unique identifiers using the Nano ID
algorithm. The default alphabet is A-Za-z0-9_- (64 characters),
matching the canonical nanoid specification.

Examples:
  openGyver generate nanoid
  openGyver generate nanoid --length 12
  openGyver generate nanoid --alphabet "0123456789abcdef"
  openGyver generate nanoid --json`,
	Args: cobra.NoArgs,
	RunE: runNanoid,
}

func runNanoid(_ *cobra.Command, _ []string) error {
	alphabet := nanoidAlphabet
	if alphabet == "" {
		alphabet = "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if len(alphabet) == 0 || len(alphabet) > 255 {
		return fmt.Errorf("alphabet must be between 1 and 255 characters")
	}

	id, err := randStringFromCharset(nanoidLength, alphabet)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":   "nanoid",
			"length": nanoidLength,
			"id":     id,
		})
	}
	fmt.Println(id)
	return nil
}

// ---------------------------------------------------------------------------
// snowflake
// ---------------------------------------------------------------------------

var (
	snowflakeOnce    sync.Mutex
	snowflakeSeq     int64
	snowflakeLastMs  int64
	snowflakeNodeID  int64
)

var snowflakeCmd = &cobra.Command{
	Use:   "snowflake",
	Short: "Generate snowflake-style 64-bit IDs",
	Long: `Generate a Twitter-style snowflake ID — a 64-bit integer composed of:

  - 41 bits: milliseconds since a custom epoch (2020-01-01T00:00:00Z)
  - 10 bits: node ID (randomly assigned per invocation)
  - 12 bits: sequence number

The resulting ID is sortable by creation time and unique within a
single process.

Examples:
  openGyver generate snowflake
  openGyver generate snowflake --json`,
	Args: cobra.NoArgs,
	RunE: runSnowflake,
}

func runSnowflake(_ *cobra.Command, _ []string) error {
	// Custom epoch: 2020-01-01T00:00:00Z
	const epoch int64 = 1577836800000

	snowflakeOnce.Lock()
	defer snowflakeOnce.Unlock()

	// Random 10-bit node ID
	if snowflakeNodeID == 0 {
		n, err := randInt(1024)
		if err != nil {
			return err
		}
		snowflakeNodeID = int64(n)
	}

	nowMs := time.Now().UnixMilli()
	if nowMs == snowflakeLastMs {
		snowflakeSeq++
		if snowflakeSeq >= 4096 { // 12-bit overflow
			// Spin until next millisecond
			for nowMs <= snowflakeLastMs {
				nowMs = time.Now().UnixMilli()
			}
			snowflakeSeq = 0
		}
	} else {
		snowflakeSeq = 0
	}
	snowflakeLastMs = nowMs

	id := ((nowMs - epoch) << 22) | (snowflakeNodeID << 12) | snowflakeSeq

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":      "snowflake",
			"id":        id,
			"timestamp": nowMs,
			"node_id":   snowflakeNodeID,
			"sequence":  snowflakeSeq,
		})
	}
	fmt.Println(id)
	return nil
}

// ---------------------------------------------------------------------------
// shortid
// ---------------------------------------------------------------------------

var shortidLength int

var shortidCmd = &cobra.Command{
	Use:   "shortid",
	Short: "Generate short random IDs",
	Long: `Generate a short, URL-safe random identifier. Uses a base62 alphabet
(A-Z a-z 0-9) for maximum density without special characters.

Examples:
  openGyver generate shortid
  openGyver generate shortid --length 12
  openGyver generate shortid --json`,
	Args: cobra.NoArgs,
	RunE: runShortID,
}

func runShortID(_ *cobra.Command, _ []string) error {
	const base62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	id, err := randStringFromCharset(shortidLength, base62)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"type":   "shortid",
			"length": shortidLength,
			"id":     id,
		})
	}
	fmt.Println(id)
	return nil
}

// ---------------------------------------------------------------------------
// Crypto helpers
// ---------------------------------------------------------------------------

// randInt returns a cryptographically random int in [0, max).
func randInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, fmt.Errorf("crypto/rand: %w", err)
	}
	return int(n.Int64()), nil
}

// randStringFromCharset builds a random string of the given length by
// picking characters uniformly from charset using crypto/rand.
func randStringFromCharset(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}
	// For charsets that are a power of 2, we could use masking, but
	// big.Int rejection sampling is correct for any alphabet size.
	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("crypto/rand: %w", err)
		}
		b[i] = charset[idx.Int64()]
	}
	return string(b), nil
}

// randBytes returns n cryptographically random bytes.
func randBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("crypto/rand: %w", err)
	}
	return buf, nil
}

// ---------------------------------------------------------------------------
// init — register everything
// ---------------------------------------------------------------------------

func init() {
	// Persistent --json flag on the parent
	generateCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	// password
	passwordCmd.Flags().IntVar(&pwLength, "length", 16, "password length")
	passwordCmd.Flags().BoolVar(&pwNoUpper, "no-upper", false, "exclude uppercase letters")
	passwordCmd.Flags().BoolVar(&pwNoLower, "no-lower", false, "exclude lowercase letters")
	passwordCmd.Flags().BoolVar(&pwNoDigits, "no-digits", false, "exclude digits")
	passwordCmd.Flags().BoolVar(&pwNoSpecial, "no-special", false, "exclude special characters")
	passwordCmd.Flags().IntVar(&pwCount, "count", 1, "number of passwords to generate")

	// passphrase
	passphraseCmd.Flags().IntVar(&ppWords, "words", 4, "number of words in the passphrase")
	passphraseCmd.Flags().StringVar(&ppSeparator, "separator", "-", "word separator")
	passphraseCmd.Flags().IntVar(&ppCount, "count", 1, "number of passphrases to generate")

	// string
	stringCmd.Flags().IntVar(&strLength, "length", 32, "string length")
	stringCmd.Flags().StringVar(&strCharset, "charset", "alphanumeric", "character set: alpha, alphanumeric, hex, base64, custom")
	stringCmd.Flags().StringVar(&strCustom, "custom", "", "custom alphabet (use with --charset=custom)")
	stringCmd.Flags().IntVar(&strCount, "count", 1, "number of strings to generate")

	// apikey
	apikeyCmd.Flags().StringVar(&akPrefix, "prefix", "", "key prefix (e.g. sk_live_)")
	apikeyCmd.Flags().IntVar(&akLength, "length", 32, "length of the random portion")

	// secret
	secretCmd.Flags().IntVar(&secretLength, "length", 64, "number of random bytes (hex output is 2x)")

	// otp
	otpCmd.Flags().StringVar(&otpIssuer, "issuer", "", "issuer name for the OTP URI")
	otpCmd.Flags().StringVar(&otpAccount, "account", "", "account name for the OTP URI")

	// nanoid
	nanoidCmd.Flags().IntVar(&nanoidLength, "length", 21, "ID length")
	nanoidCmd.Flags().StringVar(&nanoidAlphabet, "alphabet", "", "custom alphabet (default: A-Za-z0-9_-)")

	// snowflake — no extra flags

	// shortid
	shortidCmd.Flags().IntVar(&shortidLength, "length", 8, "ID length")

	// Wire subcommands
	generateCmd.AddCommand(
		passwordCmd,
		passphraseCmd,
		stringCmd,
		apikeyCmd,
		secretCmd,
		otpCmd,
		nanoidCmd,
		snowflakeCmd,
		shortidCmd,
	)

	cmd.Register(generateCmd)
}

// ---------------------------------------------------------------------------
// EFF short diceware word list (~1296 words)
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_2_0.txt
// ---------------------------------------------------------------------------

// The unused import guards are satisfied by the helpers above; add
// compile-time assertions for the imports used only in certain code paths.
var (
	_ = base64.StdEncoding
	_ = binary.BigEndian
)

var effWordList = []string{
	"acid", "acme", "acre", "acts", "aged", "agent", "aging", "agony",
	"ahead", "aide", "aims", "ajar", "alarm", "alias", "alibi", "alien",
	"align", "alike", "alive", "alley", "allot", "allow", "aloft", "alone",
	"alpha", "amaze", "amen", "amid", "amiss", "ammo", "amp", "angel",
	"anger", "angle", "ankle", "annex", "antic", "anvil", "apart", "apex",
	"apple", "apply", "apron", "arena", "argue", "arise", "armor", "army",
	"aroma", "array", "arrow", "arson", "ash", "asset", "atlas", "atom",
	"attic", "audio", "audit", "aunts", "avid", "avoid", "awake", "award",
	"awash", "awful", "axle", "azure", "bacon", "badge", "badly", "bagel",
	"baggy", "baked", "baker", "balmy", "ban", "banks", "barn", "baron",
	"bash", "basic", "basis", "batch", "bath", "baton", "bay", "beach",
	"beads", "beam", "beans", "beard", "beast", "begin", "being", "bell",
	"belly", "below", "bench", "bent", "berry", "best", "bias", "bible",
	"bike", "bills", "bind", "bingo", "birth", "black", "blade", "blame",
	"bland", "blank", "blast", "blaze", "bleak", "bleed", "blend", "bless",
	"blimp", "blind", "bliss", "blitz", "block", "blog", "blond", "blood",
	"bloom", "blown", "blues", "bluff", "blunt", "blur", "blurt", "board",
	"boast", "boat", "body", "bogus", "bolt", "bomb", "bond", "bones",
	"bonus", "book", "boost", "booth", "bore", "born", "boss", "botch",
	"both", "bound", "bow", "boxer", "brain", "brand", "brass", "brave",
	"bread", "break", "breed", "brick", "bride", "brief", "bring", "brink",
	"brisk", "broad", "broil", "broke", "brook", "broom", "broth", "brown",
	"brush", "buddy", "budge", "buggy", "build", "built", "bulge", "bulk",
	"bully", "bunch", "bunny", "burn", "burst", "buyer", "cabin", "cable",
	"cadet", "cage", "cake", "calm", "camel", "candy", "cape", "cargo",
	"carol", "carry", "carve", "case", "cash", "cast", "catch", "cause",
	"cave", "cedar", "cell", "chain", "chair", "chalk", "champ", "chant",
	"chaos", "charm", "chart", "chase", "cheap", "check", "cheek", "cheer",
	"chess", "chest", "chick", "chief", "child", "chill", "china", "chip",
	"choir", "chose", "chunk", "cinch", "civic", "civil", "claim", "clamp",
	"clap", "clash", "clasp", "class", "clean", "clear", "clerk", "click",
	"cliff", "climb", "cling", "clip", "cloak", "clock", "clone", "close",
	"cloth", "cloud", "clown", "club", "clue", "clump", "clung", "coach",
	"coast", "coat", "cobra", "cocoa", "coil", "coin", "cold", "comet",
	"comic", "comma", "cone", "coral", "cord", "core", "cork", "corn",
	"couch", "count", "court", "cover", "cozy", "crack", "craft", "cramp",
	"crane", "crash", "crate", "crave", "crawl", "crazy", "creak", "cream",
	"crest", "crew", "crime", "crisp", "cross", "crowd", "crown", "crude",
	"crush", "crust", "cubic", "curve", "cycle", "daily", "dairy", "daisy",
	"dance", "dare", "dawn", "deals", "death", "debug", "decay", "decks",
	"decor", "decoy", "deity", "delay", "delta", "delve", "demon", "denim",
	"dense", "depot", "depth", "derby", "desk", "detox", "deuce", "devil",
	"diary", "diced", "dig", "digit", "dimly", "diner", "disco", "dish",
	"ditty", "diver", "dizzy", "dock", "dodge", "doing", "doll", "donor",
	"donut", "doom", "doors", "dose", "doubt", "dough", "dove", "down",
	"dozen", "draft", "drain", "drama", "drank", "drape", "drawl", "drawn",
	"dread", "dream", "dress", "dried", "drift", "drill", "drink", "drive",
	"droit", "drone", "drool", "drops", "drove", "drums", "drunk", "dryer",
	"dual", "dub", "dude", "dug", "duke", "dummy", "dumps", "dune",
	"dunce", "duo", "dusk", "dust", "duty", "dwarf", "dying", "each",
	"eagle", "early", "earth", "easel", "east", "eaten", "eater", "ebony",
	"edges", "edict", "eight", "eject", "elbow", "elder", "elect", "elite",
	"elm", "ember", "emit", "empty", "ended", "enemy", "enjoy", "entry",
	"envoy", "equal", "equip", "erase", "error", "essay", "ethic", "evade",
	"even", "event", "every", "evil", "evoke", "exact", "exalt", "exam",
	"excel", "exert", "exile", "exist", "extra", "eyed", "fable", "faced",
	"facet", "facts", "faded", "faith", "false", "fancy", "fang", "fatal",
	"favor", "feast", "feats", "fence", "fend", "ferry", "fetch", "fever",
	"fiber", "field", "fifth", "fifty", "fight", "filth", "final", "finch",
	"fined", "fire", "firm", "first", "fish", "fist", "five", "fixed",
	"fizzy", "flag", "flame", "flask", "flat", "flaw", "fleet", "flesh",
	"flick", "fling", "flint", "flip", "float", "flock", "flood", "floor",
	"flora", "floss", "flour", "flown", "fluid", "flush", "flute", "foam",
	"focal", "focus", "foggy", "foil", "folly", "fonts", "force", "forge",
	"forgo", "forms", "forth", "forum", "fossil", "found", "fox", "foyer",
	"frail", "frame", "frank", "fraud", "freak", "freed", "fresh", "friar",
	"fried", "from", "front", "frost", "froze", "fruit", "fuels", "fully",
	"fumes", "funds", "funny", "fused", "fussy", "fuzzy", "gaffe", "gains",
	"gala", "gale", "games", "gamma", "gang", "gaps", "gases", "gauge",
	"gave", "gaze", "gears", "gem", "genre", "ghost", "giant", "gifts",
	"given", "glad", "glare", "glass", "gleam", "glide", "glint", "globe",
	"gloom", "glory", "gloss", "glove", "glow", "glue", "goat", "going",
	"gold", "golf", "gone", "goofy", "goose", "gorge", "grace", "grade",
	"grain", "grand", "grant", "grape", "graph", "grasp", "grass", "grave",
	"gravy", "great", "greed", "green", "greet", "grief", "grill", "grime",
	"grind", "gripe", "groan", "groom", "gross", "group", "grove", "growl",
	"grown", "grub", "grunt", "guard", "guess", "guide", "guild", "guilt",
	"guise", "gulch", "gulf", "gummy", "guru", "gusts", "habit", "half",
	"halt", "hands", "handy", "hang", "happy", "hardy", "harem", "harm",
	"haste", "hasty", "hatch", "haunt", "haven", "havoc", "hawk", "hazel",
	"heads", "heart", "heath", "heavy", "hedge", "hefty", "hello", "hence",
	"herbs", "heron", "hiker", "hills", "hilly", "hinge", "hippo", "hired",
	"hobby", "hoist", "holly", "homer", "honey", "honor", "hooks", "hoped",
	"horns", "horse", "host", "hotel", "hover", "howl", "hub", "human",
	"humid", "humor", "hung", "hunts", "hurry", "hutch", "hymns", "hyper",
	"icing", "icon", "ideal", "idiot", "igloo", "image", "imply", "inbox",
	"index", "indie", "inept", "infer", "ingot", "ink", "inner", "input",
	"intel", "inter", "intro", "ionic", "iron", "issue", "ivory", "jab",
	"jacks", "jaded", "jaunt", "jaws", "jazzy", "jeans", "jelly", "jewel",
	"jiffy", "jig", "jobs", "jog", "join", "joker", "jolly", "jolt",
	"joust", "judge", "juice", "juicy", "jumbo", "jumps", "jumpy", "junco",
	"juice", "jury", "karma", "kayak", "keen", "keeps", "kept", "ketch",
	"kicks", "kilt", "kind", "king", "kiosk", "kite", "knack", "knead",
	"knee", "kneel", "knelt", "knife", "knit", "knobs", "knock", "knoll",
	"knot", "known", "label", "laced", "lacks", "lance", "lanes", "lapse",
	"large", "laser", "latch", "later", "latte", "laugh", "lava", "lawn",
	"layer", "leads", "leaky", "lean", "leapt", "lease", "least", "leave",
	"ledge", "legal", "lemon", "lend", "level", "lever", "light", "lilac",
	"limbs", "limit", "linen", "liner", "links", "lions", "lists", "liter",
	"lived", "liver", "llama", "loads", "lobby", "local", "locus", "lodge",
	"lofty", "logic", "logo", "loose", "lorry", "lost", "lotus", "lousy",
	"loved", "lover", "loyal", "lucid", "lucky", "lumen", "lumps", "lunar",
	"lunch", "lunge", "lurch", "lying", "lynch", "lyric", "macro", "magic",
	"major", "maker", "mango", "manor", "maple", "march", "marsh", "mason",
	"match", "mayor", "mealy", "meant", "medal", "media", "medic", "melee",
	"melon", "memos", "mend", "mercy", "merge", "merit", "merry", "mesh",
	"metal", "meter", "midst", "might", "milky", "mills", "mimic", "mince",
	"minds", "minor", "minus", "mirth", "miser", "misty", "mixed", "mixer",
	"moat", "model", "moist", "molar", "mold", "money", "month", "moods",
	"moose", "moral", "motel", "moth", "motor", "motto", "mound", "mouse",
	"mouth", "moved", "mover", "movie", "much", "mulch", "mural", "murky",
	"music", "musty", "muted", "nacho", "nails", "naive", "named", "nanny",
	"nap", "naval", "needy", "nerve", "never", "newly", "nexus", "niche",
	"night", "noble", "noise", "nomad", "none", "north", "notch", "noted",
	"novel", "nudge", "nurse", "nutty", "nylon", "oasis", "obese", "occur",
	"ocean", "oddly", "offer", "often", "olive", "omega", "onset", "opera",
	"opted", "orbit", "order", "organ", "other", "otter", "ought", "ounce",
	"outer", "owned", "owner", "oxide", "ozone", "paced", "paddy", "pagan",
	"paint", "pairs", "palms", "panel", "panic", "paper", "parch", "parks",
	"party", "pasta", "patch", "patio", "pause", "peace", "peach", "pearl",
	"pedal", "peels", "penal", "penny", "perch", "peril", "perks", "petal",
	"petty", "phase", "phone", "photo", "piano", "picks", "piece", "piggy",
	"pilot", "pinch", "pint", "pious", "pipes", "pitch", "pivot", "pixel",
	"pizza", "place", "plaid", "plain", "plane", "plank", "plant", "plate",
	"plaza", "plead", "pleat", "plied", "plop", "plots", "pluck", "plugs",
	"plumb", "plump", "plums", "plush", "poach", "pods", "poem", "poets",
	"point", "poise", "polar", "polls", "polyp", "pond", "pooch", "pools",
	"poppy", "porch", "pork", "posed", "poser", "posse", "post", "pouch",
	"pound", "power", "prank", "prawn", "press", "price", "pride", "prime",
	"print", "prior", "prism", "privy", "prize", "probe", "prone", "proof",
	"props", "prose", "proud", "prove", "prowl", "prude", "prune", "psalm",
	"puck", "pull", "pulp", "pulse", "pumps", "punch", "pupil", "puppy",
	"purge", "purse", "pushy", "qualm", "quart", "queen", "query", "quest",
	"queue", "quick", "quiet", "quill", "quirk", "quota", "quote", "raced",
	"radar", "radio", "raft", "rage", "rainy", "raise", "rally", "ramp",
	"ranch", "range", "rapid", "rash", "rated", "ratio", "raven", "ray",
	"razor", "reach", "react", "reads", "ready", "realm", "rebel", "recap",
	"reign", "relax", "relay", "renal", "renew", "repay", "reply", "rerun",
	"reset", "resin", "retro", "rider", "ridge", "rifle", "rigid", "rings",
	"rinse", "riots", "risen", "risky", "ritzy", "rival", "river", "roads",
	"roast", "robin", "robot", "rocks", "rocky", "rogue", "roots", "ropes",
	"roses", "rotor", "rouge", "rough", "round", "route", "rover", "royal",
	"rugby", "ruins", "ruled", "ruler", "rumor", "rural", "rusty", "sadly",
	"safer", "saint", "salad", "salon", "salsa", "salty", "salve", "sands",
	"sandy", "satin", "sauce", "sauna", "saved", "savor", "scale", "scalp",
	"scam", "scare", "scarf", "scene", "scent", "scope", "score", "scout",
	"scowl", "scrap", "seals", "sedan", "seeds", "seize", "sense", "serum",
	"serve", "setup", "seven", "shack", "shade", "shaft", "shake", "shall",
	"shame", "shape", "share", "shark", "sharp", "shave", "shawl", "sheds",
	"sheen", "sheep", "sheer", "sheet", "shelf", "shell", "shift", "shiny",
	"ships", "shirt", "shock", "shoes", "shook", "shoot", "shore", "short",
	"shout", "shove", "shown", "showy", "shrub", "shrug", "sided", "siege",
	"sight", "sigma", "signs", "silky", "silly", "since", "siren", "sixth",
	"sixty", "sized", "skate", "skies", "skill", "skull", "skunk", "slate",
	"slave", "sleek", "sleep", "slept", "slice", "slide", "slope", "slots",
	"slows", "slug", "slump", "slush", "smart", "smell", "smile", "smirk",
	"smith", "smock", "smoke", "snack", "snake", "snare", "sneak", "snore",
	"snout", "snowy", "soapy", "sober", "solar", "solid", "solve", "sonic",
	"sorry", "souls", "south", "space", "spade", "spare", "spark", "spawn",
	"speak", "speed", "spend", "spent", "spice", "spied", "spike", "spine",
	"split", "spoke", "spoon", "sport", "spots", "spray", "spree", "sprig",
	"squad", "stab", "stack", "staff", "stage", "stain", "stair", "stake",
	"stale", "stall", "stamp", "stand", "stank", "stark", "start", "stash",
	"state", "stays", "steak", "steal", "steam", "steel", "steep", "steer",
	"stems", "steps", "stern", "stew", "stick", "stiff", "still", "sting",
	"stint", "stock", "stoic", "stoke", "stole", "stomp", "stone", "stood",
	"stool", "stoop", "stops", "store", "storm", "story", "stout", "stove",
	"straw", "stray", "strip", "strut", "stuck", "studs", "stuff", "stump",
	"stung", "stunk", "stunt", "style", "suave", "sugar", "suite", "sulky",
	"sunny", "super", "surge", "sushi", "swamp", "swans", "swarm", "swear",
	"sweat", "sweep", "sweet", "swept", "swift", "swim", "swine", "swing",
	"swirl", "sword", "swore", "sworn", "swung", "syrup", "tabby", "table",
	"tacky", "taint", "taken", "tales", "talks", "tally", "talon", "tamed",
	"tangy", "tanks", "tapes", "tardy", "tasks", "taste", "tasty", "taxes",
	"teach", "teams", "tease", "tempo", "tends", "tenor", "tense", "tenth",
	"terms", "tests", "theft", "theme", "thick", "thief", "thigh", "thing",
	"think", "third", "thorn", "those", "three", "threw", "throw", "thud",
	"thumb", "tidal", "tidy", "tiger", "tight", "tiles", "tilts", "timer",
	"timid", "tipsy", "tired", "titan", "title", "toast", "token", "tonal",
	"torch", "total", "touch", "tough", "towel", "tower", "towns", "toxic",
	"trace", "track", "trade", "trail", "train", "trait", "trash", "trawl",
	"treat", "trees", "trend", "trial", "tribe", "trick", "tried", "trims",
	"trio", "trips", "troll", "troop", "trout", "truce", "truck", "truly",
	"trump", "trunk", "trust", "truth", "tubes", "tucks", "tulip", "tumor",
	"tuned", "tuner", "tunes", "turbo", "turns", "tutor", "twang", "tweed",
	"twice", "twigs", "twist", "tying", "udder", "ultra", "uncle", "under",
	"unfit", "union", "unite", "unity", "until", "upper", "upset", "urban",
	"usage", "usher", "using", "usual", "utter", "vague", "valid", "valor",
	"valve", "vapid", "vault", "veins", "venue", "verge", "verse", "vibes",
	"video", "vigil", "vigor", "vinyl", "viola", "virus", "visit", "visor",
	"vista", "vital", "vivid", "vocal", "vodka", "voice", "voter", "vouch",
	"vowed", "vowel", "wacky", "wade", "wages", "wagon", "waist", "walks",
	"walls", "waltz", "wants", "wards", "waste", "watch", "water", "watts",
	"waves", "waved", "wavy", "wax", "weak", "weary", "weave", "wedge",
	"weeds", "weigh", "weird", "wells", "whale", "wheat", "wheel", "where",
	"which", "while", "whine", "whirl", "whole", "widen", "width", "wield",
	"wills", "winds", "windy", "wines", "wings", "wiped", "wired", "witch",
	"wives", "woke", "woman", "women", "woods", "woody", "words", "wordy",
	"works", "world", "worms", "worry", "worse", "worst", "worth", "would",
	"wound", "wrath", "wreck", "wrist", "wrote", "yacht", "yards", "yearn",
	"yeast", "yield", "young", "youth", "zebra", "zones",
}

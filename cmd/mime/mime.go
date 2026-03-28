package mimecmd

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── persistent flags ────────────────────────────────────────────────────────

var jsonOut bool

// ── comprehensive MIME type map ─────────────────────────────────────────────

// mimeTypes maps file extensions (with dot prefix) to MIME types.
// This supplements Go's stdlib mime package with commonly needed types.
var mimeTypes = map[string]string{
	// Text
	".txt":        "text/plain",
	".html":       "text/html",
	".htm":        "text/html",
	".css":        "text/css",
	".csv":        "text/csv",
	".tsv":        "text/tab-separated-values",
	".xml":        "text/xml",
	".rtf":        "text/rtf",
	".markdown":   "text/markdown",
	".md":         "text/markdown",
	".yaml":       "text/yaml",
	".yml":        "text/yaml",
	".ics":        "text/calendar",
	".vcf":        "text/vcard",

	// Application
	".json":       "application/json",
	".jsonld":     "application/ld+json",
	".js":         "application/javascript",
	".mjs":        "application/javascript",
	".pdf":        "application/pdf",
	".zip":        "application/zip",
	".gz":         "application/gzip",
	".gzip":       "application/gzip",
	".bz2":        "application/x-bzip2",
	".xz":         "application/x-xz",
	".zst":        "application/zstd",
	".tar":        "application/x-tar",
	".rar":        "application/vnd.rar",
	".7z":         "application/x-7z-compressed",
	".jar":        "application/java-archive",
	".war":        "application/java-archive",
	".doc":        "application/msword",
	".docx":       "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":        "application/vnd.ms-excel",
	".xlsx":       "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":        "application/vnd.ms-powerpoint",
	".pptx":       "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".odt":        "application/vnd.oasis.opendocument.text",
	".ods":        "application/vnd.oasis.opendocument.spreadsheet",
	".odp":        "application/vnd.oasis.opendocument.presentation",
	".epub":       "application/epub+zip",
	".wasm":       "application/wasm",
	".sql":        "application/sql",
	".graphql":    "application/graphql",
	".woff":       "application/font-woff",
	".woff2":      "font/woff2",
	".ttf":        "font/ttf",
	".otf":        "font/otf",
	".eot":        "application/vnd.ms-fontobject",
	".swf":        "application/x-shockwave-flash",
	".bin":        "application/octet-stream",
	".exe":        "application/vnd.microsoft.portable-executable",
	".dll":        "application/vnd.microsoft.portable-executable",
	".deb":        "application/vnd.debian.binary-package",
	".rpm":        "application/x-rpm",
	".dmg":        "application/x-apple-diskimage",
	".iso":        "application/x-iso9660-image",
	".apk":        "application/vnd.android.package-archive",
	".ipa":        "application/octet-stream",
	".toml":       "application/toml",
	".protobuf":   "application/protobuf",
	".msgpack":    "application/msgpack",
	".cbor":       "application/cbor",
	".atom":       "application/atom+xml",
	".rss":        "application/rss+xml",
	".xhtml":      "application/xhtml+xml",

	// Image
	".png":        "image/png",
	".jpg":        "image/jpeg",
	".jpeg":       "image/jpeg",
	".gif":        "image/gif",
	".bmp":        "image/bmp",
	".ico":        "image/x-icon",
	".svg":        "image/svg+xml",
	".webp":       "image/webp",
	".tiff":       "image/tiff",
	".tif":        "image/tiff",
	".avif":       "image/avif",
	".heic":       "image/heic",
	".heif":       "image/heif",
	".jxl":        "image/jxl",
	".psd":        "image/vnd.adobe.photoshop",

	// Audio
	".mp3":        "audio/mpeg",
	".wav":        "audio/wav",
	".ogg":        "audio/ogg",
	".flac":       "audio/flac",
	".aac":        "audio/aac",
	".m4a":        "audio/mp4",
	".wma":        "audio/x-ms-wma",
	".opus":       "audio/opus",
	".mid":        "audio/midi",
	".midi":       "audio/midi",
	".aiff":       "audio/aiff",

	// Video
	".mp4":        "video/mp4",
	".avi":        "video/x-msvideo",
	".mkv":        "video/x-matroska",
	".mov":        "video/quicktime",
	".wmv":        "video/x-ms-wmv",
	".flv":        "video/x-flv",
	".webm":       "video/webm",
	".m4v":        "video/mp4",
	".mpeg":       "video/mpeg",
	".mpg":        "video/mpeg",
	".3gp":        "video/3gpp",
	".m2ts":       "video/mp2t",

	// Programming / config
	".go":         "text/x-go",
	".py":         "text/x-python",
	".rb":         "text/x-ruby",
	".rs":         "text/x-rust",
	".java":       "text/x-java",
	".c":          "text/x-c",
	".cpp":        "text/x-c++",
	".h":          "text/x-c",
	".hpp":        "text/x-c++",
	".sh":         "application/x-sh",
	".bat":        "application/x-msdos-program",
	".ps1":        "application/x-powershell",
	".php":        "application/x-httpd-php",
	".ts":         "text/typescript",
	".tsx":        "text/typescript",
	".jsx":        "text/jsx",
	".swift":      "text/x-swift",
	".kt":         "text/x-kotlin",
	".scala":      "text/x-scala",
	".r":          "text/x-r",
	".lua":        "text/x-lua",
	".pl":         "text/x-perl",
	".ini":        "text/plain",
	".cfg":        "text/plain",
	".conf":       "text/plain",
	".env":        "text/plain",
	".log":        "text/plain",
}

// reverseMap builds extension-to-MIME reverse lookup.
// Built lazily on first use.
var reverseMIME map[string]string

func buildReverse() {
	if reverseMIME != nil {
		return
	}
	reverseMIME = make(map[string]string, len(mimeTypes))
	for ext, mt := range mimeTypes {
		// First extension wins (some MIME types have multiple extensions).
		if _, exists := reverseMIME[mt]; !exists {
			reverseMIME[mt] = ext
		}
	}
}

// ── parent command ──────────────────────────────────────────────────────────

var mimeCmd = &cobra.Command{
	Use:   "mime",
	Short: "MIME type lookup and detection",
	Long: `Look up MIME types by extension, find extensions for MIME types,
or detect the MIME type of a file by reading its magic bytes.

SUBCOMMANDS:

  lookup      Look up MIME type from a file extension
  extension   Find file extension for a MIME type (reverse lookup)
  detect      Detect MIME type of a file by reading magic bytes

Examples:
  openGyver mime lookup .json
  openGyver mime lookup pdf
  openGyver mime extension "application/json"
  openGyver mime detect photo.jpg`,
}

// ── subcommands ─────────────────────────────────────────────────────────────

var lookupCmd = &cobra.Command{
	Use:   "lookup <extension>",
	Short: "Look up MIME type from file extension",
	Long: `Look up the MIME type for a given file extension.

The extension can be given with or without a leading dot.
Uses Go's mime.TypeByExtension() supplemented with a comprehensive
built-in map of ~100 common types.

Examples:
  openGyver mime lookup .json
  openGyver mime lookup json
  openGyver mime lookup pdf
  openGyver mime lookup .tar.gz`,
	Args: cobra.ExactArgs(1),
	RunE: runLookup,
}

var extensionCmd = &cobra.Command{
	Use:   "extension <mime-type>",
	Short: "Find file extension for a MIME type",
	Long: `Reverse lookup: find the file extension for a given MIME type.

Returns the most common file extension associated with the MIME type.

Examples:
  openGyver mime extension "application/json"
  openGyver mime extension "image/png"
  openGyver mime extension "text/html"`,
	Args: cobra.ExactArgs(1),
	RunE: runExtension,
}

var detectCmd = &cobra.Command{
	Use:   "detect <file>",
	Short: "Detect MIME type of a file by reading magic bytes",
	Long: `Detect the MIME type of a file by reading its first 512 bytes
and using Go's net/http.DetectContentType() for content sniffing.

This performs binary/magic-byte detection, not extension-based lookup.

Examples:
  openGyver mime detect photo.jpg
  openGyver mime detect mystery-file
  openGyver mime detect /usr/bin/ls`,
	Args: cobra.ExactArgs(1),
	RunE: runDetect,
}

// ── runners ─────────────────────────────────────────────────────────────────

func runLookup(_ *cobra.Command, args []string) error {
	ext := normalizeExtension(args[0])
	mimeType := lookupMIME(ext)

	if mimeType == "" {
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"extension": ext,
				"mime_type": nil,
				"error":     "unknown extension",
			})
		}
		return fmt.Errorf("unknown extension: %s", ext)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"extension": ext,
			"mime_type": mimeType,
		})
	}

	fmt.Println(mimeType)
	return nil
}

func runExtension(_ *cobra.Command, args []string) error {
	mimeType := strings.TrimSpace(args[0])
	buildReverse()

	ext := reverseLookup(mimeType)

	if ext == "" {
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"mime_type": mimeType,
				"extension": nil,
				"error":     "unknown MIME type",
			})
		}
		return fmt.Errorf("unknown MIME type: %s", mimeType)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"mime_type": mimeType,
			"extension": ext,
		})
	}

	fmt.Println(ext)
	return nil
}

func runDetect(_ *cobra.Command, args []string) error {
	filePath := args[0]

	mimeType, err := detectFile(filePath)
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"file":      filePath,
			"mime_type": mimeType,
		})
	}

	fmt.Println(mimeType)
	return nil
}

// ── helpers ─────────────────────────────────────────────────────────────────

// normalizeExtension ensures the extension starts with a dot.
func normalizeExtension(ext string) string {
	ext = strings.TrimSpace(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return strings.ToLower(ext)
}

// lookupMIME checks the built-in map first, then falls back to Go's stdlib.
func lookupMIME(ext string) string {
	if mt, ok := mimeTypes[ext]; ok {
		return mt
	}
	return mime.TypeByExtension(ext)
}

// reverseLookup finds an extension for a MIME type.
// Checks our built-in reverse map first, then falls back to stdlib.
func reverseLookup(mimeType string) string {
	if ext, ok := reverseMIME[mimeType]; ok {
		return ext
	}

	// Try stdlib.
	exts, err := mime.ExtensionsByType(mimeType)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ""
}

// detectFile reads the first 512 bytes of a file and uses
// net/http.DetectContentType() for MIME type sniffing.
func detectFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("cannot read file: %w", err)
	}

	return http.DetectContentType(buf[:n]), nil
}

// ── registration ────────────────────────────────────────────────────────────

func init() {
	// Persistent flags on the parent — inherited by all subcommands.
	mimeCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output result as JSON")

	// Wire subcommands.
	mimeCmd.AddCommand(lookupCmd)
	mimeCmd.AddCommand(extensionCmd)
	mimeCmd.AddCommand(detectCmd)

	cmd.Register(mimeCmd)
}

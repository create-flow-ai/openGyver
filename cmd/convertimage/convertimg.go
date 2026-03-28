package convertimage

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/webp" // register WebP decoder
)

var (
	output  string
	quality int
	width   int
	height  int
	format  string
)

var convertImageCmd = &cobra.Command{
	Use:   "convertImage <image>",
	Short: "Convert between image formats",
	Long: `Convert images between popular raster formats.

The input format is auto-detected from the file header. The output format
is determined by the --output file extension.

SUPPORTED FORMATS:

  Format   Read   Write   Extensions
  ───────  ────   ─────   ──────────
  PNG      yes    yes     .png
  JPEG     yes    yes     .jpg, .jpeg
  GIF      yes    yes     .gif (first frame)
  BMP      yes    yes     .bmp
  TIFF     yes    yes     .tiff, .tif
  WebP     yes    no*     .webp
  HEIC     no*    no*     .heic, .heif

  * WebP can be decoded (read) but not encoded (write) in pure Go.
    HEIC requires platform-specific codecs — use Apple's sips or
    ImageMagick as a workaround.

OPTIONS:

  --quality    JPEG quality (1-100, default 90). Ignored for other formats.
  --width      Resize output to this width (0 = keep original).
  --height     Resize output to this height (0 = keep original).
               If only one dimension is set, the other scales proportionally.
  --format     Override input format detection (png, jpeg, gif, bmp, tiff, webp).

Examples:
  openGyver convertImage photo.png -o photo.jpg
  openGyver convertImage photo.jpg -o photo.png
  openGyver convertImage photo.jpg -o photo.bmp
  openGyver convertImage photo.bmp -o photo.tiff
  openGyver convertImage photo.webp -o photo.png
  openGyver convertImage photo.png -o photo.jpg --quality 85
  openGyver convertImage photo.png -o thumb.jpg --width 200
  openGyver convertImage photo.png -o resized.png --width 800 --height 600
  openGyver convertImage raw.dat -o out.png --format bmp`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertImg,
}

func runConvertImg(c *cobra.Command, args []string) error {
	inputPath := args[0]

	outputPath := output
	if outputPath == "" {
		return fmt.Errorf("--output (-o) is required")
	}

	outputExt := strings.ToLower(filepath.Ext(outputPath))
	outFmt := extToImgFormat(outputExt)
	if outFmt == "" {
		return fmt.Errorf("unsupported output format: %s\nSupported: .png, .jpg, .jpeg, .gif, .bmp, .tiff, .tif", outputExt)
	}

	// Decode input
	img, inputFmt, err := decodeImage(inputPath, format)
	if err != nil {
		return err
	}

	// Resize if requested
	if width > 0 || height > 0 {
		img = resize(img, width, height)
	}

	// Encode output
	if err := encodeImage(outputPath, outFmt, img); err != nil {
		return err
	}

	bounds := img.Bounds()
	fmt.Printf("Converted %s (%s) → %s (%s) [%dx%d]\n",
		inputPath, inputFmt, outputPath, outFmt, bounds.Dx(), bounds.Dy())
	return nil
}

// extToImgFormat maps file extensions to canonical format names.
func extToImgFormat(ext string) string {
	switch ext {
	case ".png":
		return "png"
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".gif":
		return "gif"
	case ".bmp":
		return "bmp"
	case ".tiff", ".tif":
		return "tiff"
	case ".webp":
		return "webp"
	case ".heic", ".heif":
		return "heic"
	default:
		return ""
	}
}

// decodeImage opens and decodes an image file. If formatHint is non-empty,
// it's used to select the decoder; otherwise the format is auto-detected.
func decodeImage(path, formatHint string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("opening image: %w", err)
	}
	defer f.Close()

	// If format hint is provided, use format-specific decoder
	if formatHint != "" {
		img, err := decodeWithHint(f, formatHint)
		if err != nil {
			return nil, "", err
		}
		return img, formatHint, nil
	}

	// Auto-detect using image.Decode (uses registered decoders)
	img, detected, err := image.Decode(f)
	if err != nil {
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".heic" || ext == ".heif" {
			return nil, "", fmt.Errorf("HEIC/HEIF decoding is not supported in pure Go.\n" +
				"Workaround: convert with Apple's sips first:\n" +
				"  sips -s format png input.heic --out input.png\n" +
				"Or use ImageMagick:\n" +
				"  magick input.heic input.png")
		}
		return nil, "", fmt.Errorf("decoding image: %w", err)
	}
	return img, detected, nil
}

func decodeWithHint(f *os.File, hint string) (image.Image, error) {
	switch strings.ToLower(hint) {
	case "png":
		return png.Decode(f)
	case "jpeg", "jpg":
		return jpeg.Decode(f)
	case "gif":
		return gif.Decode(f)
	case "bmp":
		return bmp.Decode(f)
	case "tiff", "tif":
		return tiff.Decode(f)
	default:
		// Fall back to auto-detection
		img, _, err := image.Decode(f)
		return img, err
	}
}

// encodeImage writes an image to a file in the specified format.
func encodeImage(path, format string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}
	defer f.Close()

	switch format {
	case "png":
		return png.Encode(f, img)
	case "jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
	case "gif":
		return gif.Encode(f, img, nil)
	case "bmp":
		return bmp.Encode(f, img)
	case "tiff":
		return tiff.Encode(f, img, nil)
	case "webp":
		return fmt.Errorf("WebP encoding is not supported in pure Go.\n" +
			"Workaround: convert to PNG first, then use cwebp:\n" +
			"  openGyver convertImg input.xxx -o temp.png\n" +
			"  cwebp temp.png -o output.webp")
	case "heic":
		return fmt.Errorf("HEIC encoding is not supported in pure Go.\n" +
			"Workaround: use Apple's sips:\n" +
			"  sips -s format heic input.png --out output.heic")
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// resize scales an image to the target dimensions using nearest-neighbor interpolation.
// If only one dimension is set (the other is 0), the missing dimension is calculated
// to maintain the aspect ratio.
func resize(img image.Image, targetW, targetH int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if targetW <= 0 && targetH <= 0 {
		return img
	}
	if targetW <= 0 {
		targetW = srcW * targetH / srcH
	}
	if targetH <= 0 {
		targetH = srcH * targetW / srcW
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

	for y := 0; y < targetH; y++ {
		for x := 0; x < targetW; x++ {
			srcX := x * srcW / targetW + bounds.Min.X
			srcY := y * srcH / targetH + bounds.Min.Y
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}

func init() {
	convertImageCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	convertImageCmd.Flags().IntVar(&quality, "quality", 90, "JPEG quality 1-100 (default 90)")
	convertImageCmd.Flags().IntVar(&width, "width", 0, "resize to this width (0 = keep original)")
	convertImageCmd.Flags().IntVar(&height, "height", 0, "resize to this height (0 = keep original)")
	convertImageCmd.Flags().StringVarP(&format, "format", "f", "", "override input format detection (png, jpeg, gif, bmp, tiff, webp)")
	cmd.Register(convertImageCmd)
}

package convertimage

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
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
	quiet   bool
	jsonOut bool
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
  JPEG     yes    yes     .jpg, .jpeg, .jfif
  GIF      yes    yes     .gif (first frame)
  BMP      yes    yes     .bmp
  TIFF     yes    yes     .tiff, .tif
  WebP     yes    no*     .webp
  HEIC     no**   no**    .heic, .heif
  PPM      yes    yes     .ppm, .pgm, .pbm, .pnm
  PCX      yes    no      .pcx
  TGA      yes    no      .tga
  SVG      no     yes*    .svg (raster → vector via potrace/ImageMagick)
  ICO      no     no      .ico (use toIco command instead)

  Camera RAW (read via dcraw/ImageMagick):
  CR2, CR3, CRW, NEF, ARW, DNG, ORF, RAF, RW2, PEF, ERF, MRW,
  SRF, SR2,3FR, K25, KDC, MEF, NRW, X3F, DCR, MOS, IIQ, RAW

  * WebP can be decoded (read) but not encoded (write) in pure Go.
  ** HEIC/Camera RAW require dcraw, ImageMagick, or Apple sips.

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
	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"success":       true,
			"input":         inputPath,
			"output":        outputPath,
			"input_format":  inputFmt,
			"output_format": outFmt,
			"width":         bounds.Dx(),
			"height":        bounds.Dy(),
		})
	}
	if !quiet {
		fmt.Printf("Converted %s (%s) → %s (%s) [%dx%d]\n",
			inputPath, inputFmt, outputPath, outFmt, bounds.Dx(), bounds.Dy())
	}
	return nil
}

// extToImgFormat maps file extensions to canonical format names.
func extToImgFormat(ext string) string {
	switch ext {
	case ".png":
		return "png"
	case ".jpg", ".jpeg", ".jfif", ".jpe":
		return "jpeg"
	case ".gif":
		return "gif"
	case ".bmp":
		return "bmp"
	case ".tiff", ".tif":
		return "tiff"
	case ".webp":
		return "webp"
	case ".svg":
		return "svg"
	case ".heic", ".heif":
		return "heic"
	case ".ppm", ".pgm", ".pbm", ".pnm":
		return "ppm"
	case ".tga":
		return "tga"
	case ".pcx":
		return "pcx"
	// Camera RAW formats — decoded via dcraw/ImageMagick
	case ".cr2", ".cr3", ".crw", ".nef", ".arw", ".dng", ".orf",
		".raf", ".rw2", ".pef", ".erf", ".mrw", ".srf", ".sr2",
		".3fr", ".k25", ".kdc", ".mef", ".nrw", ".x3f", ".dcr",
		".mos", ".iiq", ".raw":
		return "raw"
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

	// Check if it's a camera RAW or HEIC that needs external tools
	ext := strings.ToLower(filepath.Ext(path))
	detectedFmt := extToImgFormat(ext)
	if detectedFmt == "raw" || detectedFmt == "heic" {
		f.Close()
		return decodeWithExternalTool(path, detectedFmt)
	}

	// Auto-detect using image.Decode (uses registered decoders)
	img, detected, err := image.Decode(f)
	if err != nil {
		// Try external tool as fallback
		f.Close()
		if img2, fmt2, err2 := decodeWithExternalTool(path, detectedFmt); err2 == nil {
			return img2, fmt2, nil
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
	case "ppm":
		return encodePPM(f, img)
	case "svg":
		f.Close()
		return rasterToSVG(path, img)
	case "webp":
		return fmt.Errorf("WebP encoding is not supported in pure Go.\n" +
			"Workaround: convert to PNG first, then use cwebp:\n" +
			"  openGyver convertImage input.xxx -o temp.png\n" +
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

// decodeWithExternalTool uses dcraw or ImageMagick to decode RAW/HEIC images
// by converting to PNG in a temp file and reading that.
func decodeWithExternalTool(path, format string) (image.Image, string, error) {
	tmpFile, err := os.CreateTemp("", "opengyver-*.png")
	if err != nil {
		return nil, "", err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Try dcraw first for RAW files
	if format == "raw" {
		if dcraw, err := exec.LookPath("dcraw"); err == nil {
			cmd := exec.Command(dcraw, "-c", "-T", path)
			tiffData, err := cmd.Output()
			if err == nil {
				// dcraw -c -T outputs TIFF to stdout
				os.WriteFile(tmpPath, tiffData, 0644)
				f, err := os.Open(tmpPath)
				if err == nil {
					defer f.Close()
					img, _, err := image.Decode(f)
					if err == nil {
						return img, "raw", nil
					}
				}
			}
		}
	}

	// Try ImageMagick (magick or convert)
	magick := "magick"
	if _, err := exec.LookPath(magick); err != nil {
		magick = "convert" // older ImageMagick
		if _, err := exec.LookPath(magick); err != nil {
			// Try sips on macOS
			if sips, err := exec.LookPath("sips"); err == nil {
				cmd := exec.Command(sips, "-s", "format", "png", path, "--out", tmpPath)
				if err := cmd.Run(); err == nil {
					f, err := os.Open(tmpPath)
					if err == nil {
						defer f.Close()
						img, _, err := image.Decode(f)
						if err == nil {
							return img, format, nil
						}
					}
				}
			}
			return nil, "", fmt.Errorf("no RAW/HEIC decoder found. Install one of:\n"+
				"  dcraw:        brew install dcraw\n"+
				"  ImageMagick:  brew install imagemagick\n"+
				"  Apple sips:   built-in on macOS")
		}
	}

	cmd := exec.Command(magick, path, tmpPath)
	if err := cmd.Run(); err != nil {
		return nil, "", fmt.Errorf("%s failed: %w", magick, err)
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

// encodePPM writes an image in PPM (P6 binary) format.
func encodePPM(w *os.File, img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
	if _, err := w.WriteString(header); err != nil {
		return err
	}

	buf := make([]byte, width*3)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			idx := (x - bounds.Min.X) * 3
			buf[idx] = uint8(r >> 8)
			buf[idx+1] = uint8(g >> 8)
			buf[idx+2] = uint8(b >> 8)
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

// rasterToSVG traces a raster image to SVG using potrace or ImageMagick.
// Potrace requires BMP input, so we convert to BMP first if needed.
func rasterToSVG(outputPath string, img image.Image) error {
	// Try potrace first (best quality tracing)
	if potrace, err := exec.LookPath("potrace"); err == nil {
		return traceWithPotrace(potrace, outputPath, img)
	}

	// Try ImageMagick
	for _, magick := range []string{"magick", "convert"} {
		if m, err := exec.LookPath(magick); err == nil {
			return traceWithMagick(m, outputPath, img)
		}
	}

	return fmt.Errorf("no raster-to-SVG tracer found. Install one of:\n" +
		"  potrace:      brew install potrace (best quality)\n" +
		"  ImageMagick:  brew install imagemagick")
}

func traceWithPotrace(potracePath, outputPath string, img image.Image) error {
	// potrace requires BMP input
	tmpBMP, err := os.CreateTemp("", "opengyver-*.bmp")
	if err != nil {
		return err
	}
	tmpBMPPath := tmpBMP.Name()
	defer os.Remove(tmpBMPPath)

	if err := bmp.Encode(tmpBMP, img); err != nil {
		tmpBMP.Close()
		return fmt.Errorf("encoding temp BMP: %w", err)
	}
	tmpBMP.Close()

	cmd := exec.Command(potracePath, tmpBMPPath, "-s", "-o", outputPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("potrace error: %s\n%s", err, string(out))
	}
	return nil
}

func traceWithMagick(magickPath, outputPath string, img image.Image) error {
	// Write to temp PNG, then convert to SVG via ImageMagick
	tmpPNG, err := os.CreateTemp("", "opengyver-*.png")
	if err != nil {
		return err
	}
	tmpPNGPath := tmpPNG.Name()
	defer os.Remove(tmpPNGPath)

	if err := png.Encode(tmpPNG, img); err != nil {
		tmpPNG.Close()
		return err
	}
	tmpPNG.Close()

	cmd := exec.Command(magickPath, tmpPNGPath, outputPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ImageMagick error: %s\n%s", err, string(out))
	}
	return nil
}

// Ensure color is available for PPM encoding.
var _ = color.RGBA{}

func init() {
	convertImageCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	convertImageCmd.Flags().IntVar(&quality, "quality", 90, "JPEG quality 1-100 (default 90)")
	convertImageCmd.Flags().IntVar(&width, "width", 0, "resize to this width (0 = keep original)")
	convertImageCmd.Flags().IntVar(&height, "height", 0, "resize to this height (0 = keep original)")
	convertImageCmd.Flags().StringVarP(&format, "format", "f", "", "override input format detection (png, jpeg, gif, bmp, tiff, webp, raw)")
	convertImageCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output messages (for piping)")
	convertImageCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output result as JSON")
	cmd.Register(convertImageCmd)
}

package convertimage

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "convertimg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

// createTestImage generates a simple 64x64 RGBA test image with a red/blue gradient.
func createTestImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{uint8(x * 4), 0, uint8(y * 4), 255})
		}
	}
	return img
}

// savePNG creates a PNG test file.
func saveTestPNG(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, createTestImage()); err != nil {
		t.Fatal(err)
	}
}

// saveTestJPEG creates a JPEG test file.
func saveTestJPEG(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := jpeg.Encode(f, createTestImage(), &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
}

// saveTestBMP creates a BMP test file.
func saveTestBMP(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := bmp.Encode(f, createTestImage()); err != nil {
		t.Fatal(err)
	}
}

// saveTestGIF creates a GIF test file.
func saveTestGIF(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := gif.Encode(f, createTestImage(), nil); err != nil {
		t.Fatal(err)
	}
}

// saveTestTIFF creates a TIFF test file.
func saveTestTIFF(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := tiff.Encode(f, createTestImage(), nil); err != nil {
		t.Fatal(err)
	}
}

// verifyImage decodes an image file and checks it has valid dimensions.
func verifyImage(t *testing.T, path string) image.Image {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decoding %s: %v", path, err)
	}
	if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
		t.Fatalf("%s has zero dimensions", path)
	}
	return img
}

// ---------------------------------------------------------------------------
// Command metadata
// ---------------------------------------------------------------------------

func TestConvertImgCmd_Metadata(t *testing.T) {
	if convertImageCmd.Use != "convertImage <image>" {
		t.Errorf("unexpected Use: %s", convertImageCmd.Use)
	}
	if convertImageCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertImgCmd_RequiresOneArg(t *testing.T) {
	validator := cobra.ExactArgs(1)
	if err := validator(convertImageCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(convertImageCmd, []string{"img.png"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConvertImgCmd_Flags(t *testing.T) {
	f := convertImageCmd.Flags()
	for _, name := range []string{"output", "quality", "width", "height", "format"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found", name)
		}
	}
	if f.ShorthandLookup("o") == nil {
		t.Error("-o shorthand not found")
	}
	if f.ShorthandLookup("f") == nil {
		t.Error("-f shorthand not found")
	}
}

func TestConvertImgCmd_FlagDefaults(t *testing.T) {
	f := convertImageCmd.Flags()
	if v := f.Lookup("quality").DefValue; v != "90" {
		t.Errorf("quality default: got %s, want 90", v)
	}
	if v := f.Lookup("width").DefValue; v != "0" {
		t.Errorf("width default: got %s, want 0", v)
	}
}

// ---------------------------------------------------------------------------
// extToImgFormat
// ---------------------------------------------------------------------------

func TestExtToImgFormat(t *testing.T) {
	tests := map[string]string{
		".png":  "png",
		".jpg":  "jpeg",
		".jpeg": "jpeg",
		".jfif": "jpeg",
		".jpe":  "jpeg",
		".gif":  "gif",
		".bmp":  "bmp",
		".tiff": "tiff",
		".tif":  "tiff",
		".webp": "webp",
		".heic": "heic",
		".heif": "heic",
		".ppm":  "ppm",
		".pgm":  "ppm",
		".pbm":  "ppm",
		".pnm":  "ppm",
		".tga":  "tga",
		".pcx":  "pcx",
		".svg":  "svg",
		".cr2":  "raw",
		".nef":  "raw",
		".arw":  "raw",
		".dng":  "raw",
		".orf":  "raw",
		".raf":  "raw",
		".rw2":  "raw",
		".pef":  "raw",
		".raw":  "raw",
		".xyz":  "",
		"":      "",
	}
	for ext, want := range tests {
		got := extToImgFormat(ext)
		if got != want {
			t.Errorf("extToImgFormat(%q) = %q, want %q", ext, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// PNG → all formats
// ---------------------------------------------------------------------------

func TestPNGToJPEG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "test.jpg")
	saveTestPNG(t, src)

	quality = 90
	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestPNGToGIF(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "test.gif")
	saveTestPNG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestPNGToBMP(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "test.bmp")
	saveTestPNG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestPNGToTIFF(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "test.tiff")
	saveTestPNG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

// ---------------------------------------------------------------------------
// JPEG → all formats
// ---------------------------------------------------------------------------

func TestJPEGToPNG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.jpg")
	dst := filepath.Join(dir, "test.png")
	saveTestJPEG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestJPEGToBMP(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.jpg")
	dst := filepath.Join(dir, "test.bmp")
	saveTestJPEG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestJPEGToTIFF(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.jpg")
	dst := filepath.Join(dir, "test.tiff")
	saveTestJPEG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestJPEGToGIF(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.jpg")
	dst := filepath.Join(dir, "test.gif")
	saveTestJPEG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

// ---------------------------------------------------------------------------
// BMP → other formats
// ---------------------------------------------------------------------------

func TestBMPToPNG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.bmp")
	dst := filepath.Join(dir, "test.png")
	saveTestBMP(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestBMPToJPEG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.bmp")
	dst := filepath.Join(dir, "test.jpg")
	saveTestBMP(t, src)

	quality = 85
	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

// ---------------------------------------------------------------------------
// GIF → other formats
// ---------------------------------------------------------------------------

func TestGIFToPNG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.gif")
	dst := filepath.Join(dir, "test.png")
	saveTestGIF(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestGIFToJPEG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.gif")
	dst := filepath.Join(dir, "test.jpg")
	saveTestGIF(t, src)

	quality = 90
	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

// ---------------------------------------------------------------------------
// TIFF → other formats
// ---------------------------------------------------------------------------

func TestTIFFToPNG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.tiff")
	dst := filepath.Join(dir, "test.png")
	saveTestTIFF(t, src)

	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

func TestTIFFToJPEG(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.tiff")
	dst := filepath.Join(dir, "test.jpg")
	saveTestTIFF(t, src)

	quality = 90
	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

// ---------------------------------------------------------------------------
// Resize
// ---------------------------------------------------------------------------

func TestResize_BothDimensions(t *testing.T) {
	img := createTestImage() // 64x64
	resized := resize(img, 32, 32)
	if resized.Bounds().Dx() != 32 || resized.Bounds().Dy() != 32 {
		t.Errorf("expected 32x32, got %dx%d", resized.Bounds().Dx(), resized.Bounds().Dy())
	}
}

func TestResize_WidthOnly(t *testing.T) {
	img := createTestImage() // 64x64
	resized := resize(img, 32, 0)
	if resized.Bounds().Dx() != 32 || resized.Bounds().Dy() != 32 {
		t.Errorf("expected 32x32, got %dx%d", resized.Bounds().Dx(), resized.Bounds().Dy())
	}
}

func TestResize_HeightOnly(t *testing.T) {
	img := createTestImage() // 64x64
	resized := resize(img, 0, 16)
	if resized.Bounds().Dx() != 16 || resized.Bounds().Dy() != 16 {
		t.Errorf("expected 16x16, got %dx%d", resized.Bounds().Dx(), resized.Bounds().Dy())
	}
}

func TestResize_NoDimensions(t *testing.T) {
	img := createTestImage()
	resized := resize(img, 0, 0)
	if resized.Bounds().Dx() != 64 {
		t.Error("should return original when no dimensions set")
	}
}

func TestResize_NonSquare(t *testing.T) {
	// Create a 100x50 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	resized := resize(img, 50, 0)
	if resized.Bounds().Dx() != 50 || resized.Bounds().Dy() != 25 {
		t.Errorf("expected 50x25, got %dx%d", resized.Bounds().Dx(), resized.Bounds().Dy())
	}
}

func TestConvertWithResize(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "thumb.jpg")
	saveTestPNG(t, src)

	quality = 90
	err := runConversion(src, dst, "", 32, 0)
	if err != nil {
		t.Fatal(err)
	}

	img := verifyImage(t, dst)
	if img.Bounds().Dx() != 32 {
		t.Errorf("expected width 32, got %d", img.Bounds().Dx())
	}
}

// ---------------------------------------------------------------------------
// Roundtrip: PNG → JPEG → BMP → TIFF → GIF → PNG
// ---------------------------------------------------------------------------

func TestRoundtrip_AllFormats(t *testing.T) {
	dir := tempDir(t)
	quality = 95

	pngPath := filepath.Join(dir, "1.png")
	saveTestPNG(t, pngPath)

	steps := []struct{ from, to string }{
		{"1.png", "2.jpg"},
		{"2.jpg", "3.bmp"},
		{"3.bmp", "4.tiff"},
		{"4.tiff", "5.gif"},
		{"5.gif", "6.png"},
	}

	for _, s := range steps {
		src := filepath.Join(dir, s.from)
		dst := filepath.Join(dir, s.to)
		if err := runConversion(src, dst, "", 0, 0); err != nil {
			t.Fatalf("%s → %s: %v", s.from, s.to, err)
		}
	}

	// Verify final image is valid
	verifyImage(t, filepath.Join(dir, "6.png"))
}

// ---------------------------------------------------------------------------
// Format hint override
// ---------------------------------------------------------------------------

func TestFormatHint(t *testing.T) {
	dir := tempDir(t)
	// Save a PNG but with a .dat extension
	datPath := filepath.Join(dir, "test.dat")
	saveTestPNG(t, datPath)

	dst := filepath.Join(dir, "out.jpg")
	quality = 90
	err := runConversion(datPath, dst, "png", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	verifyImage(t, dst)
}

// ---------------------------------------------------------------------------
// WebP/HEIC stubs
// ---------------------------------------------------------------------------

func TestPNGToPPM(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "test.ppm")
	saveTestPNG(t, src)

	quality = 90
	err := runConversion(src, dst, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(dst)
	if !strings.HasPrefix(string(data), "P6") {
		t.Error("PPM file should start with P6 header")
	}
}

func TestPPM_ContentValid(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "test.ppm")
	saveTestPNG(t, src)

	runConversion(src, dst, "", 0, 0)

	data, _ := os.ReadFile(dst)
	content := string(data)
	if !strings.HasPrefix(content, "P6\n64 64\n255\n") {
		t.Errorf("PPM header wrong, got prefix: %q", content[:min(30, len(content))])
	}
	// P6 header + 64*64*3 bytes of pixel data
	expectedSize := len("P6\n64 64\n255\n") + 64*64*3
	if len(data) != expectedSize {
		t.Errorf("PPM size: got %d, want %d", len(data), expectedSize)
	}
}

func TestPNGToSVG_RequiresExternalTool(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	dst := filepath.Join(dir, "test.svg")
	saveTestPNG(t, src)

	err := runConversion(src, dst, "", 0, 0)
	// Will succeed if potrace or ImageMagick is installed, error otherwise
	if err != nil {
		if !strings.Contains(err.Error(), "no raster-to-SVG tracer found") &&
			!strings.Contains(err.Error(), "potrace") &&
			!strings.Contains(err.Error(), "ImageMagick") {
			t.Errorf("unexpected error: %v", err)
		}
		t.Logf("SVG conversion not available (expected in CI): %v", err)
	} else {
		// Verify the SVG was created
		data, _ := os.ReadFile(dst)
		if !strings.Contains(string(data), "<svg") && !strings.Contains(string(data), "<?xml") {
			t.Error("output doesn't look like SVG")
		}
	}
}

func TestEncodeWebP_ReturnsError(t *testing.T) {
	err := encodeImage("/tmp/test.webp", "webp", createTestImage())
	if err == nil {
		t.Error("expected error for WebP encoding")
	}
}

func TestEncodeHEIC_ReturnsError(t *testing.T) {
	err := encodeImage("/tmp/test.heic", "heic", createTestImage())
	if err == nil {
		t.Error("expected error for HEIC encoding")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestConvert_MissingInput(t *testing.T) {
	err := runConversion("/nonexistent.png", "/tmp/out.jpg", "", 0, 0)
	if err == nil {
		t.Error("expected error for missing input")
	}
}

func TestConvert_BadOutputDir(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "test.png")
	saveTestPNG(t, src)

	err := runConversion(src, "/nonexistent/dir/out.jpg", "", 0, 0)
	if err == nil {
		t.Error("expected error for bad output directory")
	}
}

func TestConvert_UnsupportedOutputFormat(t *testing.T) {
	err := encodeImage("/tmp/test.xyz", "xyz", createTestImage())
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

// ---------------------------------------------------------------------------
// Helper: direct conversion without going through cobra
// ---------------------------------------------------------------------------

func runConversion(inputPath, outputPath, fmtHint string, w, h int) error {
	img, inputFmt, err := decodeImage(inputPath, fmtHint)
	if err != nil {
		return err
	}
	_ = inputFmt

	if w > 0 || h > 0 {
		img = resize(img, w, h)
	}

	outFmt := extToImgFormat(filepath.Ext(outputPath))
	return encodeImage(outputPath, outFmt, img)
}

package extractframe

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "extractframe-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func createTestGIF(t *testing.T, path string, frames int) {
	t.Helper()
	g := &gif.GIF{}
	for i := 0; i < frames; i++ {
		img := image.NewPaletted(image.Rect(0, 0, 10, 10), color.Palette{
			color.RGBA{uint8(i * 50), 0, 0, 255},
			color.White,
		})
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, 10)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := gif.EncodeAll(f, g); err != nil {
		t.Fatal(err)
	}
}

func TestExtractFrameCmd_Metadata(t *testing.T) {
	if extractFrameCmd.Use != "extractFrame <animated-image>" {
		t.Errorf("unexpected Use: %s", extractFrameCmd.Use)
	}
	if extractFrameCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestExtractFrameCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(extractFrameCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := v(extractFrameCmd, []string{"anim.gif"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExtractFrameCmd_Flags(t *testing.T) {
	f := extractFrameCmd.Flags()
	for _, name := range []string{"output", "frame", "all", "quiet", "json"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found", name)
		}
	}
}

func TestExtractGIFFrame_First(t *testing.T) {
	dir := tempDir(t)
	gifPath := filepath.Join(dir, "test.gif")
	outPath := filepath.Join(dir, "frame0.png")
	createTestGIF(t, gifPath, 5)

	output = outPath
	frame = 0
	defer func() { output = ""; frame = 0 }()

	err := extractGIFFrame(gifPath, 0)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = png.Decode(f)
	if err != nil {
		t.Fatal("output is not valid PNG")
	}
}

func TestExtractGIFFrame_Last(t *testing.T) {
	dir := tempDir(t)
	gifPath := filepath.Join(dir, "test.gif")
	outPath := filepath.Join(dir, "frame4.png")
	createTestGIF(t, gifPath, 5)

	output = outPath
	defer func() { output = "" }()

	err := extractGIFFrame(gifPath, 4)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}

func TestExtractGIFFrame_OutOfRange(t *testing.T) {
	dir := tempDir(t)
	gifPath := filepath.Join(dir, "test.gif")
	createTestGIF(t, gifPath, 3)

	output = filepath.Join(dir, "out.png")
	defer func() { output = "" }()

	err := extractGIFFrame(gifPath, 5)
	if err == nil {
		t.Error("expected error for out-of-range frame")
	}
}

func TestExtractGIFFrame_NegativeIndex(t *testing.T) {
	dir := tempDir(t)
	gifPath := filepath.Join(dir, "test.gif")
	createTestGIF(t, gifPath, 3)

	output = filepath.Join(dir, "out.png")
	defer func() { output = "" }()

	err := extractGIFFrame(gifPath, -1)
	if err == nil {
		t.Error("expected error for negative frame index")
	}
}

func TestExtractAllGIFFrames(t *testing.T) {
	dir := tempDir(t)
	gifPath := filepath.Join(dir, "test.gif")
	createTestGIF(t, gifPath, 3)

	output = ""
	err := extractAllGIFFrames(gifPath)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		base := filepath.Join(dir, "test")
		path := base + "_" + padZero(i) + ".png"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("frame %d not extracted: %s", i, path)
		}
	}
}

func padZero(n int) string {
	if n < 10 {
		return "00" + string(rune('0'+n))
	}
	return ""
}

func TestExtractGIFFrame_MissingFile(t *testing.T) {
	err := extractGIFFrame("/nonexistent.gif", 0)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestExtractGIFFrame_InvalidFile(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "bad.gif")
	os.WriteFile(path, []byte("not a gif"), 0644)

	output = filepath.Join(dir, "out.png")
	defer func() { output = "" }()

	err := extractGIFFrame(path, 0)
	if err == nil {
		t.Error("expected error for invalid GIF")
	}
}

func TestExtractAPNGFrame_UnsupportedIndex(t *testing.T) {
	dir := tempDir(t)
	// Create a simple PNG (not animated)
	path := filepath.Join(dir, "test.png")
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	f, _ := os.Create(path)
	png.Encode(f, img)
	f.Close()

	output = filepath.Join(dir, "out.png")
	defer func() { output = "" }()

	// Frame 0 should work
	err := extractAPNGFrame(path, 0)
	if err != nil {
		t.Errorf("frame 0 should work: %v", err)
	}

	// Frame > 0 should error (needs ffmpeg)
	err = extractAPNGFrame(path, 5)
	if err == nil {
		t.Error("expected error for APNG frame > 0 without ffmpeg")
	}
}

func TestUnsupportedFormat(t *testing.T) {
	// Test via the runExtractFrame function
	err := runExtractFrame(extractFrameCmd, []string{"video.mp4"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

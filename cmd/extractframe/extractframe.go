package extractframe

import (
	"fmt"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	output  string
	frame   int
	all     bool
	quiet   bool
	jsonOut bool
)

var extractFrameCmd = &cobra.Command{
	Use:   "extractFrame <animated-image>",
	Short: "Extract frames from animated GIF or APNG",
	Long: `Extract a specific frame or all frames from an animated GIF or animated PNG.

By default extracts frame 0 (the first frame). Use --frame to pick a
specific frame, or --all to extract every frame as individual PNG files.

SUPPORTED INPUT:

  .gif    Animated GIF
  .png    Animated PNG (APNG) — extracts using ffmpeg if available

OUTPUT:

  Single frame: saved as PNG (default: input_frameN.png)
  All frames:   saved as input_frame000.png, input_frame001.png, ...

EXAMPLES:

  openGyver extractFrame animation.gif
  openGyver extractFrame animation.gif --frame 5
  openGyver extractFrame animation.gif --frame 5 -o keyframe.png
  openGyver extractFrame animation.gif --all
  openGyver extractFrame animation.gif --all -o frames/frame.png
  openGyver extractFrame animated.png --frame 0`,
	Args: cobra.ExactArgs(1),
	RunE: runExtractFrame,
}

func runExtractFrame(c *cobra.Command, args []string) error {
	inputPath := args[0]
	ext := strings.ToLower(filepath.Ext(inputPath))

	switch ext {
	case ".gif":
		if all {
			return extractAllGIFFrames(inputPath)
		}
		return extractGIFFrame(inputPath, frame)
	case ".png":
		return extractAPNGFrame(inputPath, frame)
	default:
		return fmt.Errorf("unsupported format: %s (use .gif or .png)", ext)
	}
}

func extractGIFFrame(inputPath string, idx int) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return fmt.Errorf("decoding GIF: %w", err)
	}

	if idx < 0 || idx >= len(g.Image) {
		return fmt.Errorf("frame %d out of range (GIF has %d frames, 0-%d)", idx, len(g.Image), len(g.Image)-1)
	}

	outPath := output
	if outPath == "" {
		base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
		outPath = fmt.Sprintf("%s_frame%d.png", base, idx)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}
	defer out.Close()

	if err := png.Encode(out, g.Image[idx]); err != nil {
		return fmt.Errorf("encoding PNG: %w", err)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"success":      true,
			"input":        inputPath,
			"output":       outPath,
			"frame":        idx,
			"total_frames": len(g.Image),
		})
	}
	if !quiet {
		fmt.Printf("Extracted frame %d/%d → %s\n", idx, len(g.Image)-1, outPath)
	}
	return nil
}

func extractAllGIFFrames(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return fmt.Errorf("decoding GIF: %w", err)
	}

	base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
	outDir := filepath.Dir(base)
	outBase := filepath.Base(base)

	if output != "" {
		outDir = filepath.Dir(output)
		outBase = strings.TrimSuffix(filepath.Base(output), filepath.Ext(output))
		os.MkdirAll(outDir, 0755)
	}

	var paths []string
	for i, img := range g.Image {
		outPath := filepath.Join(outDir, fmt.Sprintf("%s_%03d.png", outBase, i))
		out, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("creating %s: %w", outPath, err)
		}
		if err := png.Encode(out, img); err != nil {
			out.Close()
			return fmt.Errorf("encoding frame %d: %w", i, err)
		}
		out.Close()
		paths = append(paths, outPath)
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"success":      true,
			"input":        inputPath,
			"total_frames": len(g.Image),
			"files":        paths,
		})
	}
	if !quiet {
		fmt.Printf("Extracted %d frames from %s\n", len(g.Image), inputPath)
	}
	return nil
}

func extractAPNGFrame(inputPath string, idx int) error {
	// APNG is complex to parse in pure Go. For frame 0, we can just
	// decode the default image (which Go's png.Decode returns).
	if idx == 0 {
		f, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		img, err := png.Decode(f)
		if err != nil {
			return fmt.Errorf("decoding PNG: %w", err)
		}

		outPath := output
		if outPath == "" {
			base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
			outPath = fmt.Sprintf("%s_frame0.png", base)
		}

		out, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("creating output: %w", err)
		}
		defer out.Close()

		if err := png.Encode(out, img); err != nil {
			return fmt.Errorf("encoding PNG: %w", err)
		}

		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{
				"success": true, "input": inputPath,
				"output": outPath, "frame": 0,
			})
		}
		if !quiet {
			fmt.Printf("Extracted frame 0 → %s\n", outPath)
		}
		return nil
	}

	// For other frames, APNG requires ffmpeg
	return fmt.Errorf("extracting frame %d from APNG requires ffmpeg.\n"+
		"Use: ffmpeg -i %s -vf \"select=eq(n\\,%d)\" -vframes 1 frame.png", idx, inputPath, idx)
}

func init() {
	extractFrameCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (default: input_frameN.png)")
	extractFrameCmd.Flags().IntVar(&frame, "frame", 0, "frame number to extract (0-indexed)")
	extractFrameCmd.Flags().BoolVar(&all, "all", false, "extract all frames as individual PNGs")
	extractFrameCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output messages")
	extractFrameCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	cmd.Register(extractFrameCmd)
}

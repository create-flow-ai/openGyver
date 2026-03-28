package trimvideo

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	output  string
	start   string
	end     string
	dur     string
	codec   string
	quiet   bool
	jsonOut bool
)

var trimVideoCmd = &cobra.Command{
	Use:   "trimVideo <input-file>",
	Short: "Trim a video file to a specific time range",
	Long: `Extract a portion of a video file by specifying start/end times or duration.

By default, copies the streams without re-encoding (instant, lossless).
Use --codec to re-encode to a different format.

REQUIRES: ffmpeg must be installed and available in PATH.

TIME FORMAT:

  Seconds:         90 (= 1m30s into the video)
  HH:MM:SS:        00:01:30
  HH:MM:SS.ms:     00:01:30.500

OPTIONS:

  --start / -s     Start time (default: beginning of file)
  --end / -e       End time (alternative to --duration)
  --duration / -d  Duration to extract (alternative to --end)
  --output / -o    Output file (default: input_trimmed.ext)
  --codec          Re-encode with codec (e.g. libx264, libx265, copy)

EXAMPLES:

  openGyver trimVideo input.mp4 -s 00:01:00 -e 00:02:30
  openGyver trimVideo input.mp4 -s 30 -d 60 -o clip.mp4
  openGyver trimVideo movie.mkv -s 00:10:00 -e 00:15:00 -o scene.mp4
  openGyver trimVideo input.mp4 -s 0 -d 10 --codec libx264
  openGyver trimVideo input.avi -s 00:05:00 -e 00:10:00 -o clip.webm`,
	Args: cobra.ExactArgs(1),
	RunE: runTrimVideo,
}

func runTrimVideo(c *cobra.Command, args []string) error {
	if err := checkFFmpeg(); err != nil {
		return err
	}

	inputPath := args[0]
	outputPath := output
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(inputPath, ext)
		outputPath = base + "_trimmed" + ext
	}

	ffArgs := []string{"-y"}

	// Start time
	if start != "" {
		ffArgs = append(ffArgs, "-ss", start)
	}

	ffArgs = append(ffArgs, "-i", inputPath)

	// End time or duration
	if end != "" {
		ffArgs = append(ffArgs, "-to", end)
	} else if dur != "" {
		ffArgs = append(ffArgs, "-t", dur)
	}

	// Codec
	if codec != "" {
		ffArgs = append(ffArgs, "-c:v", codec)
		if codec == "copy" {
			ffArgs = append(ffArgs, "-c:a", "copy")
		}
	} else {
		// Default: stream copy (no re-encoding, instant)
		ffArgs = append(ffArgs, "-c", "copy")
	}

	ffArgs = append(ffArgs, outputPath)

	ffCmd := exec.Command("ffmpeg", ffArgs...)
	out, err := ffCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %s\n%s", err, string(out))
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"success":  true,
			"input":    inputPath,
			"output":   outputPath,
			"start":    start,
			"end":      end,
			"duration": dur,
			"codec":    codecOrCopy(),
		})
	}
	if !quiet {
		fmt.Printf("Trimmed %s → %s\n", inputPath, outputPath)
	}
	return nil
}

func codecOrCopy() string {
	if codec != "" {
		return codec
	}
	return "copy"
}

func checkFFmpeg() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg is not installed or not in PATH.\n" +
			"Install it:\n" +
			"  macOS:   brew install ffmpeg\n" +
			"  Linux:   apt install ffmpeg\n" +
			"  Windows: https://ffmpeg.org/download.html")
	}
	return nil
}

func init() {
	trimVideoCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (default: input_trimmed.ext)")
	trimVideoCmd.Flags().StringVarP(&start, "start", "s", "", "start time (HH:MM:SS or seconds)")
	trimVideoCmd.Flags().StringVarP(&end, "end", "e", "", "end time (HH:MM:SS or seconds)")
	trimVideoCmd.Flags().StringVarP(&dur, "duration", "d", "", "duration to extract (HH:MM:SS or seconds)")
	trimVideoCmd.Flags().StringVar(&codec, "codec", "", "video codec (default: copy = no re-encoding)")
	trimVideoCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output messages")
	trimVideoCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	cmd.Register(trimVideoCmd)
}

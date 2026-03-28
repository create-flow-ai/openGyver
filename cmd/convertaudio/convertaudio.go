package convertaudio

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
	bitrate string
	sample  string
)

var supportedFormats = map[string]bool{
	"aac": true, "ac3": true, "aif": true, "aifc": true, "aiff": true,
	"amr": true, "au": true, "caf": true, "dss": true, "flac": true,
	"m4a": true, "m4b": true, "mp3": true, "oga": true, "ogg": true,
	"opus": true, "voc": true, "wav": true, "weba": true, "wma": true,
}

var convertAudioCmd = &cobra.Command{
	Use:   "convertAudio <input-file>",
	Short: "Convert between audio formats",
	Long: `Convert audio files between popular formats using ffmpeg.

REQUIRES: ffmpeg must be installed and available in PATH.
  macOS:   brew install ffmpeg
  Linux:   apt install ffmpeg / dnf install ffmpeg
  Windows: https://ffmpeg.org/download.html

SUPPORTED FORMATS:

  AAC, AC3, AIF, AIFC, AIFF, AMR, AU, CAF, DSS, FLAC,
  M4A, M4B, MP3, OGA, OGG, OPUS, VOC, WAV, WEBA, WMA

Examples:
  openGyver convertAudio song.wav -o song.mp3
  openGyver convertAudio song.flac -o song.aac
  openGyver convertAudio podcast.mp3 -o podcast.ogg
  openGyver convertAudio song.wav -o song.mp3 --bitrate 320k
  openGyver convertAudio song.mp3 -o song.wav --sample 44100`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertAudio,
}

func runConvertAudio(c *cobra.Command, args []string) error {
	if err := checkFFmpeg(); err != nil {
		return err
	}

	inputPath := args[0]
	if output == "" {
		return fmt.Errorf("--output (-o) is required")
	}

	outExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(output)), ".")
	if !supportedFormats[outExt] {
		return fmt.Errorf("unsupported output format: .%s", outExt)
	}

	ffArgs := []string{"-i", inputPath, "-y"}
	if bitrate != "" {
		ffArgs = append(ffArgs, "-b:a", bitrate)
	}
	if sample != "" {
		ffArgs = append(ffArgs, "-ar", sample)
	}
	ffArgs = append(ffArgs, output)

	cmd := exec.Command("ffmpeg", ffArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %s\n%s", err, string(out))
	}

	fmt.Printf("Converted %s → %s\n", inputPath, output)
	return nil
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
	convertAudioCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	convertAudioCmd.Flags().StringVar(&bitrate, "bitrate", "", "audio bitrate (e.g. 128k, 192k, 320k)")
	convertAudioCmd.Flags().StringVar(&sample, "sample", "", "sample rate in Hz (e.g. 44100, 48000)")
	cmd.Register(convertAudioCmd)
}

package convertvideo

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	output     string
	resolution string
	vbitrate   string
	abitrate   string
	fps        string
	codec      string
)

var supportedFormats = map[string]bool{
	"3g2": true, "3gp": true, "3gpp": true, "asf": true, "avi": true,
	"av1": true, "cavs": true, "divx": true, "dv": true, "dvr": true,
	"f4v": true, "flv": true, "hevc": true, "m2ts": true, "m2v": true,
	"m4v": true, "mjpeg": true, "mkv": true, "mod": true, "mov": true,
	"mp4": true, "mpeg": true, "mpg": true, "mts": true, "mxf": true,
	"ogg": true, "ogv": true, "rm": true, "rmvb": true, "swf": true,
	"tod": true, "ts": true, "vob": true, "webm": true, "wmv": true,
	"wtv": true, "xvid": true,
}

var convertVideoCmd = &cobra.Command{
	Use:   "convertVideo <input-file>",
	Short: "Convert between video formats",
	Long: `Convert video files between popular formats using ffmpeg.

REQUIRES: ffmpeg must be installed and available in PATH.
  macOS:   brew install ffmpeg
  Linux:   apt install ffmpeg / dnf install ffmpeg
  Windows: https://ffmpeg.org/download.html

SUPPORTED FORMATS (37):

  3G2, 3GP, 3GPP, ASF, AV1, AVI, CAVS, DIVX, DV, DVR, F4V, FLV,
  HEVC, M2TS, M2V, M4V, MJPEG, MKV, MOD, MOV, MP4, MPEG, MPG,
  MTS, MXF, OGG, OGV, RM, RMVB, SWF, TOD, TS, VOB, WEBM, WMV,
  WTV, XVID

Examples:
  openGyver convertVideo input.avi -o output.mp4
  openGyver convertVideo input.mkv -o output.webm
  openGyver convertVideo input.mov -o output.mp4 --resolution 1920x1080
  openGyver convertVideo input.mp4 -o output.mp4 --vbitrate 5M --abitrate 192k
  openGyver convertVideo input.mp4 -o output.gif --fps 10
  openGyver convertVideo input.mp4 -o output.webm --codec libvpx-vp9`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertVideo,
}

func runConvertVideo(c *cobra.Command, args []string) error {
	if err := checkFFmpeg(); err != nil {
		return err
	}

	inputPath := args[0]
	if output == "" {
		return fmt.Errorf("--output (-o) is required")
	}

	outExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(output)), ".")
	if !supportedFormats[outExt] && outExt != "gif" {
		return fmt.Errorf("unsupported output format: .%s", outExt)
	}

	ffArgs := []string{"-i", inputPath, "-y"}
	if resolution != "" {
		ffArgs = append(ffArgs, "-s", resolution)
	}
	if vbitrate != "" {
		ffArgs = append(ffArgs, "-b:v", vbitrate)
	}
	if abitrate != "" {
		ffArgs = append(ffArgs, "-b:a", abitrate)
	}
	if fps != "" {
		ffArgs = append(ffArgs, "-r", fps)
	}
	if codec != "" {
		ffArgs = append(ffArgs, "-c:v", codec)
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
	convertVideoCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	convertVideoCmd.Flags().StringVar(&resolution, "resolution", "", "output resolution (e.g. 1920x1080, 1280x720)")
	convertVideoCmd.Flags().StringVar(&vbitrate, "vbitrate", "", "video bitrate (e.g. 2M, 5M, 10M)")
	convertVideoCmd.Flags().StringVar(&abitrate, "abitrate", "", "audio bitrate (e.g. 128k, 192k, 320k)")
	convertVideoCmd.Flags().StringVar(&fps, "fps", "", "frames per second (e.g. 24, 30, 60)")
	convertVideoCmd.Flags().StringVar(&codec, "codec", "", "video codec (e.g. libx264, libx265, libvpx-vp9)")
	cmd.Register(convertVideoCmd)
}

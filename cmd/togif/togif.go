package togif

import (
	"fmt"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	delay  int
	output string
	loop   int
	format string
	width  int
	height int
)

var togifCmd = &cobra.Command{
	Use:   "toGif <image> [image...]",
	Short: "Create an animated GIF from images",
	Long: `Combine multiple images into an animated GIF.

Images are added as frames in the order given. The input format is
auto-detected from each file's header. Use --format to override
detection when all inputs share the same non-standard extension.

Use --width and --height to set the output GIF dimensions. If only
one dimension is given, the other is scaled proportionally from the
first frame. If neither is given, the first frame's dimensions are used.

Supported input formats: png, jpeg, gif, bmp, tiff, webp.

Examples:
  openGyver toGif frame1.png frame2.png frame3.png
  openGyver toGif frame*.png -o animation.gif --delay 50
  openGyver toGif frame*.png --loop 3
  openGyver toGif *.dat --format bmp
  openGyver toGif frame*.png --width 320 --height 240`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		detected := format
		if detected == "" {
			detected = "auto"
		}

		dim := "from first frame"
		if width > 0 || height > 0 {
			dim = fmt.Sprintf("%dx%d", width, height)
		}

		fmt.Printf("Creating animated GIF → %s (%d frames, format=%s, size=%s, delay=%dms, loop=%d)\n",
			output, len(args), detected, dim, delay, loop)
		// TODO: implement GIF assembly
		return nil
	},
}

func init() {
	togifCmd.Flags().StringVarP(&output, "output", "o", "output.gif", "output file path")
	togifCmd.Flags().IntVar(&delay, "delay", 100, "delay between frames in ms")
	togifCmd.Flags().IntVar(&loop, "loop", 0, "number of loops (0 = infinite)")
	togifCmd.Flags().StringVarP(&format, "format", "f", "", "input image format (png, jpeg, gif, bmp, tiff, webp). Auto-detected if omitted")
	togifCmd.Flags().IntVar(&width, "width", 0, "output GIF width in pixels (0 = use first frame's width)")
	togifCmd.Flags().IntVar(&height, "height", 0, "output GIF height in pixels (0 = use first frame's height)")
	cmd.Register(togifCmd)
}

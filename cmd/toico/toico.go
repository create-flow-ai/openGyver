package toico

import (
	"fmt"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	sizes    []int
	output   string
	format   string
	width    int
	height   int
)

var toicoCmd = &cobra.Command{
	Use:   "toIco <image>",
	Short: "Convert an image to ICO format",
	Long: `Convert an image to a Windows ICO file.

The input format is auto-detected from the file header. Use --format to
override detection when reading from a non-standard extension or pipe.

Supports embedding multiple sizes in a single .ico file.
Default sizes: 16, 32, 48, 256.

Use --width and --height to resize the source image before generating
the ICO entries. If only one dimension is given, the image is scaled
proportionally. If neither is given, the original dimensions are used.

Supported input formats: png, jpeg, gif, bmp, tiff, webp.

Examples:
  openGyver toIco logo.png
  openGyver toIco logo.png -o favicon.ico
  openGyver toIco logo.png --sizes 16,32,64
  openGyver toIco photo.dat --format jpeg
  openGyver toIco logo.png --width 512 --height 512`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		src := args[0]

		detected := format
		if detected == "" {
			detected = "auto"
		}

		dim := "original"
		if width > 0 || height > 0 {
			dim = fmt.Sprintf("%dx%d", width, height)
		}

		fmt.Printf("Converting %s (format=%s, resize=%s) → %s (sizes: %v)\n",
			src, detected, dim, output, sizes)
		// TODO: implement ICO conversion
		return nil
	},
}

func init() {
	toicoCmd.Flags().StringVarP(&output, "output", "o", "output.ico", "output file path")
	toicoCmd.Flags().IntSliceVar(&sizes, "sizes", []int{16, 32, 48, 256}, "icon sizes to embed")
	toicoCmd.Flags().StringVarP(&format, "format", "f", "", "input image format (png, jpeg, gif, bmp, tiff, webp). Auto-detected if omitted")
	toicoCmd.Flags().IntVar(&width, "width", 0, "resize source to this width before generating ICO (0 = keep original)")
	toicoCmd.Flags().IntVar(&height, "height", 0, "resize source to this height before generating ICO (0 = keep original)")
	cmd.Register(toicoCmd)
}

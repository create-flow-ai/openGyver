package color

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// Flag variables
// ---------------------------------------------------------------------------

var jsonOut bool

// convert flags
var toFormat string

// palette flags
var (
	paletteCount int
	paletteType  string
)

// random flags
var (
	randomFormat string
	randomCount  int
)

// ---------------------------------------------------------------------------
// Color types & conversions
// ---------------------------------------------------------------------------

type RGB struct{ R, G, B uint8 }
type HSL struct{ H, S, L float64 } // H 0-360, S 0-100, L 0-100
type CMYK struct{ C, M, Y, K float64 } // each 0-100

// --- Hex <-> RGB ---

func hexToRGB(hex string) (RGB, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color: %q", hex)
	}
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color: %q", hex)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color: %q", hex)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color: %q", hex)
	}
	return RGB{uint8(r), uint8(g), uint8(b)}, nil
}

func rgbToHex(c RGB) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// --- RGB <-> HSL ---

func rgbToHSL(c RGB) HSL {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	l := (max + min) / 2.0

	if delta == 0 {
		return HSL{0, 0, l * 100}
	}

	var s float64
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2.0 - max - min)
	}

	var h float64
	switch max {
	case r:
		h = (g - b) / delta
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/delta + 2
	case b:
		h = (r-g)/delta + 4
	}
	h *= 60

	return HSL{math.Round(h*10) / 10, math.Round(s*1000) / 10, math.Round(l*1000) / 10}
}

func hslToRGB(c HSL) RGB {
	h := c.H
	s := c.S / 100.0
	l := c.L / 100.0

	if s == 0 {
		v := uint8(math.Round(l * 255))
		return RGB{v, v, v}
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	hueToRGB := func(p, q, t float64) float64 {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		switch {
		case t < 1.0/6.0:
			return p + (q-p)*6*t
		case t < 1.0/2.0:
			return q
		case t < 2.0/3.0:
			return p + (q-p)*(2.0/3.0-t)*6
		default:
			return p
		}
	}

	hNorm := h / 360.0
	r := hueToRGB(p, q, hNorm+1.0/3.0)
	g := hueToRGB(p, q, hNorm)
	b := hueToRGB(p, q, hNorm-1.0/3.0)

	return RGB{
		uint8(math.Round(r * 255)),
		uint8(math.Round(g * 255)),
		uint8(math.Round(b * 255)),
	}
}

// --- RGB <-> CMYK ---

func rgbToCMYK(c RGB) CMYK {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	k := 1.0 - math.Max(r, math.Max(g, b))
	if k == 1.0 {
		return CMYK{0, 0, 0, 100}
	}

	cm := (1 - r - k) / (1 - k) * 100
	m := (1 - g - k) / (1 - k) * 100
	y := (1 - b - k) / (1 - k) * 100

	return CMYK{
		math.Round(cm*10) / 10,
		math.Round(m*10) / 10,
		math.Round(y*10) / 10,
		math.Round(k*1000) / 10,
	}
}

func cmykToRGB(c CMYK) RGB {
	cm := c.C / 100.0
	m := c.M / 100.0
	y := c.Y / 100.0
	k := c.K / 100.0

	r := 255 * (1 - cm) * (1 - k)
	g := 255 * (1 - m) * (1 - k)
	b := 255 * (1 - y) * (1 - k)

	return RGB{
		uint8(math.Round(r)),
		uint8(math.Round(g)),
		uint8(math.Round(b)),
	}
}

// --- Format helpers ---

func formatRGB(c RGB) string   { return fmt.Sprintf("rgb(%d,%d,%d)", c.R, c.G, c.B) }
func formatHSL(h HSL) string   { return fmt.Sprintf("hsl(%g,%g%%,%g%%)", h.H, h.S, h.L) }
func formatCMYK(c CMYK) string { return fmt.Sprintf("cmyk(%g,%g,%g,%g)", c.C, c.M, c.Y, c.K) }

func allFormats(c RGB) map[string]string {
	hsl := rgbToHSL(c)
	cmyk := rgbToCMYK(c)
	return map[string]string{
		"hex":  rgbToHex(c),
		"rgb":  formatRGB(c),
		"hsl":  formatHSL(hsl),
		"cmyk": formatCMYK(cmyk),
	}
}

// ---------------------------------------------------------------------------
// Input parsing (auto-detect)
// ---------------------------------------------------------------------------

func parseColor(input string) (RGB, string, error) {
	input = strings.TrimSpace(input)

	// Hex: #rgb or #rrggbb
	if strings.HasPrefix(input, "#") {
		c, err := hexToRGB(input)
		if err != nil {
			return RGB{}, "", err
		}
		return c, "hex", nil
	}

	lower := strings.ToLower(input)

	// rgb(r,g,b)
	if strings.HasPrefix(lower, "rgb(") && strings.HasSuffix(lower, ")") {
		inner := input[4 : len(input)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 3 {
			return RGB{}, "", fmt.Errorf("invalid rgb() format: %q", input)
		}
		r, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil || r < 0 || r > 255 {
			return RGB{}, "", fmt.Errorf("invalid rgb() red value: %q", parts[0])
		}
		g, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || g < 0 || g > 255 {
			return RGB{}, "", fmt.Errorf("invalid rgb() green value: %q", parts[1])
		}
		b, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil || b < 0 || b > 255 {
			return RGB{}, "", fmt.Errorf("invalid rgb() blue value: %q", parts[2])
		}
		return RGB{uint8(r), uint8(g), uint8(b)}, "rgb", nil
	}

	// hsl(h,s%,l%)
	if strings.HasPrefix(lower, "hsl(") && strings.HasSuffix(lower, ")") {
		inner := input[4 : len(input)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 3 {
			return RGB{}, "", fmt.Errorf("invalid hsl() format: %q", input)
		}
		h, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil || h < 0 || h > 360 {
			return RGB{}, "", fmt.Errorf("invalid hsl() hue value: %q", parts[0])
		}
		sStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(parts[1]), "%"))
		s, err := strconv.ParseFloat(sStr, 64)
		if err != nil || s < 0 || s > 100 {
			return RGB{}, "", fmt.Errorf("invalid hsl() saturation value: %q", parts[1])
		}
		lStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(parts[2]), "%"))
		l, err := strconv.ParseFloat(lStr, 64)
		if err != nil || l < 0 || l > 100 {
			return RGB{}, "", fmt.Errorf("invalid hsl() lightness value: %q", parts[2])
		}
		return hslToRGB(HSL{h, s, l}), "hsl", nil
	}

	// cmyk(c,m,y,k)
	if strings.HasPrefix(lower, "cmyk(") && strings.HasSuffix(lower, ")") {
		inner := input[5 : len(input)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 4 {
			return RGB{}, "", fmt.Errorf("invalid cmyk() format: %q", input)
		}
		vals := make([]float64, 4)
		for i, p := range parts {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err != nil || v < 0 || v > 100 {
				return RGB{}, "", fmt.Errorf("invalid cmyk() value: %q", p)
			}
			vals[i] = v
		}
		return cmykToRGB(CMYK{vals[0], vals[1], vals[2], vals[3]}), "cmyk", nil
	}

	// Try bare hex (no #)
	if len(input) == 6 || len(input) == 3 {
		c, err := hexToRGB(input)
		if err == nil {
			return c, "hex", nil
		}
	}

	return RGB{}, "", fmt.Errorf("unrecognized color format: %q (supported: hex, rgb(), hsl(), cmyk())", input)
}

// ---------------------------------------------------------------------------
// WCAG contrast ratio
// ---------------------------------------------------------------------------

func relativeLuminance(c RGB) float64 {
	linearize := func(v uint8) float64 {
		s := float64(v) / 255.0
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	r := linearize(c.R)
	g := linearize(c.G)
	b := linearize(c.B)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func contrastRatio(c1, c2 RGB) float64 {
	l1 := relativeLuminance(c1)
	l2 := relativeLuminance(c2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// ---------------------------------------------------------------------------
// CSS named colors (148 entries, hex -> name)
// ---------------------------------------------------------------------------

var cssColors = map[string]string{
	"#f0f8ff": "aliceblue",
	"#faebd7": "antiquewhite",
	"#00ffff": "cyan",
	"#7fffd4": "aquamarine",
	"#f0ffff": "azure",
	"#f5f5dc": "beige",
	"#ffe4c4": "bisque",
	"#000000": "black",
	"#ffebcd": "blanchedalmond",
	"#0000ff": "blue",
	"#8a2be2": "blueviolet",
	"#a52a2a": "brown",
	"#deb887": "burlywood",
	"#5f9ea0": "cadetblue",
	"#7fff00": "chartreuse",
	"#d2691e": "chocolate",
	"#ff7f50": "coral",
	"#6495ed": "cornflowerblue",
	"#fff8dc": "cornsilk",
	"#dc143c": "crimson",
	"#00008b": "darkblue",
	"#008b8b": "darkcyan",
	"#b8860b": "darkgoldenrod",
	"#a9a9a9": "darkgray",
	"#006400": "darkgreen",
	"#bdb76b": "darkkhaki",
	"#8b008b": "darkmagenta",
	"#556b2f": "darkolivegreen",
	"#ff8c00": "darkorange",
	"#9932cc": "darkorchid",
	"#8b0000": "darkred",
	"#e9967a": "darksalmon",
	"#8fbc8f": "darkseagreen",
	"#483d8b": "darkslateblue",
	"#2f4f4f": "darkslategray",
	"#00ced1": "darkturquoise",
	"#9400d3": "darkviolet",
	"#ff1493": "deeppink",
	"#00bfff": "deepskyblue",
	"#696969": "dimgray",
	"#1e90ff": "dodgerblue",
	"#b22222": "firebrick",
	"#fffaf0": "floralwhite",
	"#228b22": "forestgreen",
	"#ff00ff": "magenta",
	"#dcdcdc": "gainsboro",
	"#f8f8ff": "ghostwhite",
	"#ffd700": "gold",
	"#daa520": "goldenrod",
	"#808080": "gray",
	"#008000": "green",
	"#adff2f": "greenyellow",
	"#f0fff0": "honeydew",
	"#ff69b4": "hotpink",
	"#cd5c5c": "indianred",
	"#4b0082": "indigo",
	"#fffff0": "ivory",
	"#f0e68c": "khaki",
	"#e6e6fa": "lavender",
	"#fff0f5": "lavenderblush",
	"#7cfc00": "lawngreen",
	"#fffacd": "lemonchiffon",
	"#add8e6": "lightblue",
	"#f08080": "lightcoral",
	"#e0ffff": "lightcyan",
	"#fafad2": "lightgoldenrodyellow",
	"#d3d3d3": "lightgray",
	"#90ee90": "lightgreen",
	"#ffb6c1": "lightpink",
	"#ffa07a": "lightsalmon",
	"#20b2aa": "lightseagreen",
	"#87cefa": "lightskyblue",
	"#778899": "lightslategray",
	"#b0c4de": "lightsteelblue",
	"#ffffe0": "lightyellow",
	"#00ff00": "lime",
	"#32cd32": "limegreen",
	"#faf0e6": "linen",
	"#800000": "maroon",
	"#66cdaa": "mediumaquamarine",
	"#0000cd": "mediumblue",
	"#ba55d3": "mediumorchid",
	"#9370db": "mediumpurple",
	"#3cb371": "mediumseagreen",
	"#7b68ee": "mediumslateblue",
	"#00fa9a": "mediumspringgreen",
	"#48d1cc": "mediumturquoise",
	"#c71585": "mediumvioletred",
	"#191970": "midnightblue",
	"#f5fffa": "mintcream",
	"#ffe4e1": "mistyrose",
	"#ffe4b5": "moccasin",
	"#ffdead": "navajowhite",
	"#000080": "navy",
	"#fdf5e6": "oldlace",
	"#808000": "olive",
	"#6b8e23": "olivedrab",
	"#ffa500": "orange",
	"#ff4500": "orangered",
	"#da70d6": "orchid",
	"#eee8aa": "palegoldenrod",
	"#98fb98": "palegreen",
	"#afeeee": "paleturquoise",
	"#db7093": "palevioletred",
	"#ffefd5": "papayawhip",
	"#ffdab9": "peachpuff",
	"#cd853f": "peru",
	"#ffc0cb": "pink",
	"#dda0dd": "plum",
	"#b0e0e6": "powderblue",
	"#800080": "purple",
	"#663399": "rebeccapurple",
	"#ff0000": "red",
	"#bc8f8f": "rosybrown",
	"#4169e1": "royalblue",
	"#8b4513": "saddlebrown",
	"#fa8072": "salmon",
	"#f4a460": "sandybrown",
	"#2e8b57": "seagreen",
	"#fff5ee": "seashell",
	"#a0522d": "sienna",
	"#c0c0c0": "silver",
	"#87ceeb": "skyblue",
	"#6a5acd": "slateblue",
	"#708090": "slategray",
	"#fffafa": "snow",
	"#00ff7f": "springgreen",
	"#4682b4": "steelblue",
	"#d2b48c": "tan",
	"#008080": "teal",
	"#d8bfd8": "thistle",
	"#ff6347": "tomato",
	"#40e0d0": "turquoise",
	"#ee82ee": "violet",
	"#f5deb3": "wheat",
	"#ffffff": "white",
	"#f5f5f5": "whitesmoke",
	"#ffff00": "yellow",
	"#9acd32": "yellowgreen",
}

// Build a reverse lookup: name -> hex (for the name-to-hex direction if needed).
// For closest-name lookup we need name -> RGB.
type namedColor struct {
	hex  string
	name string
	rgb  RGB
}

var namedColors []namedColor

func init() {
	// deduplicate (aqua/cyan and fuchsia/magenta share hex values)
	seen := map[string]bool{}
	for hex, name := range cssColors {
		if seen[hex] {
			continue
		}
		seen[hex] = true
		rgb, _ := hexToRGB(hex)
		namedColors = append(namedColors, namedColor{hex: hex, name: name, rgb: rgb})
	}
}

func closestColorName(c RGB) (string, string, float64) {
	bestDist := math.MaxFloat64
	bestName := ""
	bestHex := ""
	for _, nc := range namedColors {
		dr := float64(int(c.R) - int(nc.rgb.R))
		dg := float64(int(c.G) - int(nc.rgb.G))
		db := float64(int(c.B) - int(nc.rgb.B))
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestName = nc.name
			bestHex = nc.hex
		}
	}
	return bestName, bestHex, math.Sqrt(bestDist)
}

// ---------------------------------------------------------------------------
// Palette generation
// ---------------------------------------------------------------------------

func generatePalette(base RGB, count int, ptype string) ([]RGB, error) {
	hsl := rgbToHSL(base)
	colors := make([]RGB, 0, count)

	switch ptype {
	case "complementary":
		colors = append(colors, base)
		for i := 1; i < count; i++ {
			// Distribute evenly from base toward complement
			offset := 180.0 * float64(i) / float64(count-1)
			if count == 1 {
				offset = 180
			}
			h := math.Mod(hsl.H+offset, 360)
			colors = append(colors, hslToRGB(HSL{h, hsl.S, hsl.L}))
		}
	case "analogous":
		spread := 30.0
		start := hsl.H - spread*float64(count-1)/2
		for i := 0; i < count; i++ {
			h := math.Mod(start+spread*float64(i)+360, 360)
			colors = append(colors, hslToRGB(HSL{h, hsl.S, hsl.L}))
		}
	case "triadic":
		colors = append(colors, base)
		for i := 1; i < count; i++ {
			offset := 120.0 * float64(i) / float64(count-1) * 2
			if count == 1 {
				offset = 0
			}
			h := math.Mod(hsl.H+offset, 360)
			colors = append(colors, hslToRGB(HSL{h, hsl.S, hsl.L}))
		}
	case "shades":
		for i := 0; i < count; i++ {
			l := hsl.L * (1.0 - float64(i)/float64(count))
			colors = append(colors, hslToRGB(HSL{hsl.H, hsl.S, l}))
		}
	case "tints":
		for i := 0; i < count; i++ {
			l := hsl.L + (100-hsl.L)*float64(i)/float64(count)
			colors = append(colors, hslToRGB(HSL{hsl.H, hsl.S, l}))
		}
	default:
		return nil, fmt.Errorf("unknown palette type: %q (supported: complementary, analogous, triadic, shades, tints)", ptype)
	}
	return colors, nil
}

// ---------------------------------------------------------------------------
// Commands
// ---------------------------------------------------------------------------

var colorCmd = &cobra.Command{
	Use:   "color",
	Short: "Color conversion, contrast checking, and palette generation",
	Long: `Color utilities for working with hex, RGB, HSL, and CMYK colors.

SUBCOMMANDS:

  convert    Convert a color between formats (hex, rgb, hsl, cmyk)
  contrast   Check WCAG contrast ratio between two colors
  palette    Generate a color palette from a base color
  name       Look up the closest CSS color name for a color
  random     Generate random colors

Use --json/-j on any subcommand for structured JSON output.

Examples:
  openGyver color convert "#ff5733" --to rgb
  openGyver color convert "rgb(255,87,51)" --to hsl
  openGyver color contrast "#ffffff" "#000000"
  openGyver color palette "#ff5733" --type analogous --count 7
  openGyver color name "#e6e6fa"
  openGyver color random --format rgb --count 3`,
}

// --- convert ------------------------------------------------------------------

var convertCmd = &cobra.Command{
	Use:   "convert <color>",
	Short: "Convert between color formats",
	Long: `Convert a color value between hex, RGB, HSL, and CMYK formats.

The input format is auto-detected. Use --to to specify the desired output
format. When --json is used, all formats are returned at once regardless
of --to.

SUPPORTED FORMATS:
  hex    #rrggbb or #rgb (e.g. "#ff5733", "#f00")
  rgb    rgb(r,g,b)      (e.g. "rgb(255,87,51)")
  hsl    hsl(h,s%,l%)    (e.g. "hsl(11,100%,60%)")
  cmyk   cmyk(c,m,y,k)   (e.g. "cmyk(0,66,80,0)")

Examples:
  openGyver color convert "#ff5733" --to rgb
  openGyver color convert "rgb(255,87,51)" --to hsl
  openGyver color convert "hsl(11,100%,60%)" --to hex
  openGyver color convert "cmyk(0,66,80,0)" --to rgb
  openGyver color convert "#ff5733" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func runConvert(_ *cobra.Command, args []string) error {
	rgb, _, err := parseColor(args[0])
	if err != nil {
		return err
	}

	if jsonOut {
		return cmd.PrintJSON(allFormats(rgb))
	}

	switch strings.ToLower(toFormat) {
	case "hex":
		fmt.Println(rgbToHex(rgb))
	case "rgb":
		fmt.Println(formatRGB(rgb))
	case "hsl":
		fmt.Println(formatHSL(rgbToHSL(rgb)))
	case "cmyk":
		fmt.Println(formatCMYK(rgbToCMYK(rgb)))
	default:
		return fmt.Errorf("unknown target format: %q (supported: hex, rgb, hsl, cmyk)", toFormat)
	}
	return nil
}

// --- contrast -----------------------------------------------------------------

var contrastCmd = &cobra.Command{
	Use:   "contrast <color1> <color2>",
	Short: "WCAG contrast ratio checker",
	Long: `Calculate the WCAG 2.1 contrast ratio between two colors and report
whether the pair passes AA and AAA accessibility levels.

Both colors are auto-detected and can be in any supported format (hex,
rgb, hsl, cmyk).

WCAG LEVELS:
  AA  normal text   >= 4.5:1
  AA  large text    >= 3.0:1
  AAA normal text   >= 7.0:1
  AAA large text    >= 4.5:1

Examples:
  openGyver color contrast "#ffffff" "#000000"
  openGyver color contrast "#ff5733" "rgb(0,0,0)"
  openGyver color contrast "hsl(0,0%,100%)" "#333333" --json`,
	Args: cobra.ExactArgs(2),
	RunE: runContrast,
}

func runContrast(_ *cobra.Command, args []string) error {
	c1, _, err := parseColor(args[0])
	if err != nil {
		return fmt.Errorf("color 1: %w", err)
	}
	c2, _, err := parseColor(args[1])
	if err != nil {
		return fmt.Errorf("color 2: %w", err)
	}

	ratio := contrastRatio(c1, c2)
	ratio = math.Round(ratio*100) / 100

	aaNormal := ratio >= 4.5
	aaLarge := ratio >= 3.0
	aaaNormal := ratio >= 7.0
	aaaLarge := ratio >= 4.5

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"color1": rgbToHex(c1),
			"color2": rgbToHex(c2),
			"ratio":  ratio,
			"aa": map[string]interface{}{
				"normal_text": passFailStr(aaNormal),
				"large_text":  passFailStr(aaLarge),
			},
			"aaa": map[string]interface{}{
				"normal_text": passFailStr(aaaNormal),
				"large_text":  passFailStr(aaaLarge),
			},
		})
	}

	fmt.Printf("Contrast ratio: %.2f:1\n", ratio)
	fmt.Printf("  AA  normal text (>= 4.5:1): %s\n", passFailStr(aaNormal))
	fmt.Printf("  AA  large text  (>= 3.0:1): %s\n", passFailStr(aaLarge))
	fmt.Printf("  AAA normal text (>= 7.0:1): %s\n", passFailStr(aaaNormal))
	fmt.Printf("  AAA large text  (>= 4.5:1): %s\n", passFailStr(aaaLarge))
	return nil
}

func passFailStr(pass bool) string {
	if pass {
		return "PASS"
	}
	return "FAIL"
}

// --- palette ------------------------------------------------------------------

var paletteCmd = &cobra.Command{
	Use:   "palette <color>",
	Short: "Generate a color palette from a base color",
	Long: `Generate a palette of colors derived from a base color.

PALETTE TYPES:
  complementary   Colors spread toward the complement (opposite on the wheel)
  analogous       Colors adjacent on the color wheel (30 degree spread)
  triadic         Colors spread across a 240 degree arc (three-way split)
  shades          Progressively darker versions of the base
  tints           Progressively lighter versions of the base

Use --count to control the number of colors (default 5).

Examples:
  openGyver color palette "#ff5733"
  openGyver color palette "#ff5733" --type analogous --count 7
  openGyver color palette "rgb(100,149,237)" --type shades
  openGyver color palette "#336699" --type tints --count 10 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runPalette,
}

func runPalette(_ *cobra.Command, args []string) error {
	base, _, err := parseColor(args[0])
	if err != nil {
		return err
	}

	colors, err := generatePalette(base, paletteCount, paletteType)
	if err != nil {
		return err
	}

	if jsonOut {
		out := make([]map[string]string, len(colors))
		for i, c := range colors {
			out[i] = allFormats(c)
		}
		return cmd.PrintJSON(map[string]interface{}{
			"base":    rgbToHex(base),
			"type":    paletteType,
			"count":   len(colors),
			"palette": out,
		})
	}

	for _, c := range colors {
		fmt.Println(rgbToHex(c))
	}
	return nil
}

// --- name ---------------------------------------------------------------------

var nameCmd = &cobra.Command{
	Use:   "name <color>",
	Short: "Look up the closest CSS color name",
	Long: `Find the closest CSS named color for a given color value.

Compares the input color against all 148 CSS named colors using
Euclidean distance in RGB space and returns the best match.

Examples:
  openGyver color name "#e6e6fa"
  openGyver color name "#ff5734"
  openGyver color name "rgb(100,149,237)" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runName,
}

func runName(_ *cobra.Command, args []string) error {
	rgb, _, err := parseColor(args[0])
	if err != nil {
		return err
	}

	name, hex, dist := closestColorName(rgb)

	if jsonOut {
		exact := dist == 0
		return cmd.PrintJSON(map[string]interface{}{
			"input":    rgbToHex(rgb),
			"name":     name,
			"hex":      hex,
			"exact":    exact,
			"distance": math.Round(dist*100) / 100,
		})
	}

	if dist == 0 {
		fmt.Printf("%s (%s) — exact match\n", name, hex)
	} else {
		fmt.Printf("%s (%s) — closest match (distance: %.2f)\n", name, hex, dist)
	}
	return nil
}

// --- random -------------------------------------------------------------------

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Generate random colors",
	Long: `Generate one or more random colors in the specified format.

FORMATS:
  hex   #rrggbb (default)
  rgb   rgb(r,g,b)
  hsl   hsl(h,s%,l%)

Use --count to generate multiple colors at once.

Examples:
  openGyver color random
  openGyver color random --format rgb
  openGyver color random --format hsl --count 5
  openGyver color random --count 10 --json`,
	Args: cobra.NoArgs,
	RunE: runRandom,
}

func runRandom(_ *cobra.Command, args []string) error {
	colors := make([]RGB, randomCount)
	for i := range colors {
		colors[i] = RGB{
			R: uint8(rand.IntN(256)),
			G: uint8(rand.IntN(256)),
			B: uint8(rand.IntN(256)),
		}
	}

	if jsonOut {
		out := make([]map[string]string, len(colors))
		for i, c := range colors {
			out[i] = allFormats(c)
		}
		return cmd.PrintJSON(map[string]interface{}{
			"count":  len(colors),
			"colors": out,
		})
	}

	for _, c := range colors {
		switch strings.ToLower(randomFormat) {
		case "hex":
			fmt.Println(rgbToHex(c))
		case "rgb":
			fmt.Println(formatRGB(c))
		case "hsl":
			fmt.Println(formatHSL(rgbToHSL(c)))
		default:
			return fmt.Errorf("unknown format: %q (supported: hex, rgb, hsl)", randomFormat)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func init() {
	// Persistent flag on the parent: --json/-j available to all subcommands.
	colorCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	// convert flags
	convertCmd.Flags().StringVar(&toFormat, "to", "hex", "target format: hex, rgb, hsl, cmyk")

	// palette flags
	paletteCmd.Flags().IntVar(&paletteCount, "count", 5, "number of colors to generate")
	paletteCmd.Flags().StringVar(&paletteType, "type", "complementary", "palette type: complementary, analogous, triadic, shades, tints")

	// random flags
	randomCmd.Flags().StringVar(&randomFormat, "format", "hex", "output format: hex, rgb, hsl")
	randomCmd.Flags().IntVar(&randomCount, "count", 1, "number of colors to generate")

	// Wire subcommands
	colorCmd.AddCommand(convertCmd)
	colorCmd.AddCommand(contrastCmd)
	colorCmd.AddCommand(paletteCmd)
	colorCmd.AddCommand(nameCmd)
	colorCmd.AddCommand(randomCmd)

	cmd.Register(colorCmd)
}

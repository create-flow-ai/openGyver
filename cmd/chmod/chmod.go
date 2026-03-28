package chmod

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── flags ──────────────────────────────────────────────────────────────────

var jsonOut bool

// ── parent command ─────────────────────────────────────────────────────────

var chmodCmd = &cobra.Command{
	Use:   "chmod",
	Short: "File permission calculator",
	Long: `Convert between octal and symbolic file permissions and calculate umask effects.

SUBCOMMANDS:

  calc    Convert between octal and symbolic permissions
  umask   Show resulting default permissions for a given umask

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver chmod calc 755
  openGyver chmod calc rwxr-xr-x
  openGyver chmod umask 022`,
}

// ── calc subcommand ────────────────────────────────────────────────────────

var calcCmd = &cobra.Command{
	Use:   "calc <permission>",
	Short: "Convert between octal and symbolic permissions",
	Long: `Convert file permissions between octal (e.g. 755) and symbolic (e.g. rwxr-xr-x).

Auto-detects the input format. Shows octal, symbolic, and a per-class
breakdown of read/write/execute for owner, group, and other.

EXAMPLES:

  openGyver chmod calc 755
  openGyver chmod calc 644
  openGyver chmod calc rwxr-xr-x
  openGyver chmod calc rw-r--r--`,
	Args: cobra.ExactArgs(1),
	RunE: runCalc,
}

func runCalc(_ *cobra.Command, args []string) error {
	input := args[0]
	perm, err := parsePerm(input)
	if err != nil {
		return err
	}

	octal := fmt.Sprintf("%03o", perm)
	symbolic := toSymbolic(perm)
	breakdown := permBreakdown(perm)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"input":    input,
			"octal":    octal,
			"symbolic": symbolic,
			"owner":    breakdown[0],
			"group":    breakdown[1],
			"other":    breakdown[2],
		})
	}

	fmt.Printf("Octal:    %s\n", octal)
	fmt.Printf("Symbolic: %s\n", symbolic)
	fmt.Println()
	fmt.Println("  Class   Read  Write  Execute")
	fmt.Println("  ─────   ────  ─────  ───────")
	for _, b := range breakdown {
		fmt.Printf("  %-7s %-5s %-6s %s\n", b["class"], b["read"], b["write"], b["execute"])
	}
	return nil
}

// ── umask subcommand ───────────────────────────────────────────────────────

var umaskCmd = &cobra.Command{
	Use:   "umask <value>",
	Short: "Show resulting default permissions for a given umask",
	Long: `Calculate the resulting default permissions when a umask is applied.

Files start with base permission 666 and directories with 777.
The umask is subtracted to produce the effective permission.

EXAMPLES:

  openGyver chmod umask 022
  openGyver chmod umask 077`,
	Args: cobra.ExactArgs(1),
	RunE: runUmask,
}

func runUmask(_ *cobra.Command, args []string) error {
	input := args[0]
	mask, err := strconv.ParseUint(input, 8, 16)
	if err != nil || mask > 0777 {
		return fmt.Errorf("invalid umask value %q (expected octal, e.g. 022)", input)
	}

	filePerm := uint16(0666) &^ uint16(mask)
	dirPerm := uint16(0777) &^ uint16(mask)

	fileOctal := fmt.Sprintf("%03o", filePerm)
	dirOctal := fmt.Sprintf("%03o", dirPerm)
	fileSymbolic := toSymbolic(filePerm)
	dirSymbolic := toSymbolic(dirPerm)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"umask": input,
			"file": map[string]interface{}{
				"octal":    fileOctal,
				"symbolic": fileSymbolic,
			},
			"directory": map[string]interface{}{
				"octal":    dirOctal,
				"symbolic": dirSymbolic,
			},
		})
	}

	fmt.Printf("Umask: %s\n\n", input)
	fmt.Printf("  Files:       %s  (%s)\n", fileOctal, fileSymbolic)
	fmt.Printf("  Directories: %s  (%s)\n", dirOctal, dirSymbolic)
	return nil
}

// ── init ───────────────────────────────────────────────────────────────────

func init() {
	chmodCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	chmodCmd.AddCommand(calcCmd)
	chmodCmd.AddCommand(umaskCmd)
	cmd.Register(chmodCmd)
}

// ── permission helpers ─────────────────────────────────────────────────────

// parsePerm detects whether input is octal or symbolic and returns the numeric permission.
func parsePerm(input string) (uint16, error) {
	// Try octal first: 3 or 4 digit octal number.
	if len(input) >= 3 && len(input) <= 4 && isOctalString(input) {
		n, err := strconv.ParseUint(input, 8, 16)
		if err != nil {
			return 0, fmt.Errorf("invalid octal permission %q", input)
		}
		return uint16(n), nil
	}

	// Try symbolic: exactly 9 characters like rwxr-xr-x.
	if len(input) == 9 && isSymbolicString(input) {
		return fromSymbolic(input)
	}

	return 0, fmt.Errorf("unrecognized permission format %q (expected octal like 755 or symbolic like rwxr-xr-x)", input)
}

func isOctalString(s string) bool {
	for _, c := range s {
		if c < '0' || c > '7' {
			return false
		}
	}
	return true
}

func isSymbolicString(s string) bool {
	for i, c := range s {
		pos := i % 3
		switch pos {
		case 0:
			if c != 'r' && c != '-' {
				return false
			}
		case 1:
			if c != 'w' && c != '-' {
				return false
			}
		case 2:
			if c != 'x' && c != '-' {
				return false
			}
		}
	}
	return true
}

// fromSymbolic converts "rwxr-xr-x" to the numeric permission.
func fromSymbolic(s string) (uint16, error) {
	if len(s) != 9 {
		return 0, fmt.Errorf("symbolic permission must be 9 characters")
	}
	var perm uint16
	for i, c := range s {
		bit := uint16(8 - i)
		switch {
		case i%3 == 0 && c == 'r':
			perm |= 1 << bit
		case i%3 == 1 && c == 'w':
			perm |= 1 << bit
		case i%3 == 2 && c == 'x':
			perm |= 1 << bit
		case c == '-':
			// no permission
		default:
			return 0, fmt.Errorf("unexpected character %q at position %d", string(c), i)
		}
	}
	return perm, nil
}

// toSymbolic converts a numeric permission to "rwxr-xr-x" form.
func toSymbolic(perm uint16) string {
	var b strings.Builder
	chars := [3]byte{'r', 'w', 'x'}
	for shift := 6; shift >= 0; shift -= 3 {
		triplet := (perm >> uint(shift)) & 7
		for i, c := range chars {
			if triplet&(1<<uint(2-i)) != 0 {
				b.WriteByte(c)
			} else {
				b.WriteByte('-')
			}
		}
	}
	return b.String()
}

// permBreakdown returns per-class permission info.
func permBreakdown(perm uint16) []map[string]string {
	classes := []string{"Owner", "Group", "Other"}
	result := make([]map[string]string, 3)
	for i, class := range classes {
		shift := uint(6 - i*3)
		triplet := (perm >> shift) & 7
		result[i] = map[string]string{
			"class":   class,
			"read":    yesNo(triplet&4 != 0),
			"write":   yesNo(triplet&2 != 0),
			"execute": yesNo(triplet&1 != 0),
		}
	}
	return result
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

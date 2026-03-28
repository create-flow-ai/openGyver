package diff

import (
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// ── persistent flags ───────────────────────────────────────────────────────

var jsonOut bool

// ── parent command ─────────────────────────────────────────────────────────

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff and compare tools",
	Long: `Diff and compare files — unified text diff, JSON structural diff, CSV diff.

SUBCOMMANDS:

  text   Unified text diff between two files
  json   Structural diff between two JSON files
  csv    Diff between two CSV files

All subcommands support --json/-j for machine-readable output.

EXAMPLES:

  openGyver diff text --file1 a.txt --file2 b.txt
  openGyver diff json --file1 old.json --file2 new.json
  openGyver diff csv  --file1 old.csv  --file2 new.csv`,
}

// register adds a subcommand to the diff parent command.
func register(sub *cobra.Command) {
	diffCmd.AddCommand(sub)
}

func init() {
	diffCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	cmd.Register(diffCmd)
}

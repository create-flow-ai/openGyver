package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openGyver",
	Short: "openGyver — a Swiss-army-knife CLI for everyday conversions",
	Long: `openGyver is a plugin-based CLI tool for image conversions,
unit conversions, and other handy utilities.

Each command is a self-contained plugin with its own help and flags.
Run "openGyver <command> --help" for details on any command.`,
}

// Register adds a subcommand to the root. Plugins call this from their init().
func Register(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func Execute() error {
	return rootCmd.Execute()
}

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd_Metadata(t *testing.T) {
	if rootCmd.Use != "openGyver" {
		t.Errorf("unexpected Use: %s", rootCmd.Use)
	}
	if rootCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if rootCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestRegister(t *testing.T) {
	before := len(rootCmd.Commands())

	Register(&cobra.Command{
		Use:   "testplugin",
		Short: "test plugin",
	})

	after := len(rootCmd.Commands())
	if after != before+1 {
		t.Errorf("expected %d commands after Register, got %d", before+1, after)
	}

	// Clean up: remove the test command
	rootCmd.RemoveCommand(rootCmd.Commands()[after-1])
}

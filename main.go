package main

import (
	"os"

	"github.com/mj/opengyver/cmd"

	// Plugins — each init() registers itself with the root command.
	_ "github.com/mj/opengyver/cmd/convert"
	_ "github.com/mj/opengyver/cmd/togif"
	_ "github.com/mj/opengyver/cmd/toico"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

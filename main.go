package main

import (
	"os"

	"github.com/mj/opengyver/cmd"

	// Plugins — each init() registers itself with the root command.
	_ "github.com/mj/opengyver/cmd/archive"
	_ "github.com/mj/opengyver/cmd/color"
	_ "github.com/mj/opengyver/cmd/convert"
	_ "github.com/mj/opengyver/cmd/convertaudio"
	_ "github.com/mj/opengyver/cmd/convertcad"
	_ "github.com/mj/opengyver/cmd/convertebook"
	_ "github.com/mj/opengyver/cmd/convertfile"
	_ "github.com/mj/opengyver/cmd/convertfont"
	_ "github.com/mj/opengyver/cmd/convertimage"
	_ "github.com/mj/opengyver/cmd/convertpresentation"
	_ "github.com/mj/opengyver/cmd/convertvector"
	_ "github.com/mj/opengyver/cmd/convertvideo"
	_ "github.com/mj/opengyver/cmd/crypto"
	_ "github.com/mj/opengyver/cmd/encode"
	_ "github.com/mj/opengyver/cmd/epoch"
	_ "github.com/mj/opengyver/cmd/generate"
	_ "github.com/mj/opengyver/cmd/hash"
	_ "github.com/mj/opengyver/cmd/jsontools"
	_ "github.com/mj/opengyver/cmd/network"
	_ "github.com/mj/opengyver/cmd/number"
	_ "github.com/mj/opengyver/cmd/qr"
	_ "github.com/mj/opengyver/cmd/stock"
	_ "github.com/mj/opengyver/cmd/text"
	_ "github.com/mj/opengyver/cmd/timex"
	_ "github.com/mj/opengyver/cmd/togif"
	_ "github.com/mj/opengyver/cmd/toico"
	_ "github.com/mj/opengyver/cmd/uuid"
)

// Set by GoReleaser via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version, commit, date)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

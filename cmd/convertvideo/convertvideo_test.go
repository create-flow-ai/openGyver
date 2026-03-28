package convertvideo

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestConvertVideoCmd_Metadata(t *testing.T) {
	if convertVideoCmd.Use != "convertVideo <input-file>" {
		t.Errorf("unexpected Use: %s", convertVideoCmd.Use)
	}
	if convertVideoCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertVideoCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(convertVideoCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
}

func TestConvertVideoCmd_Flags(t *testing.T) {
	f := convertVideoCmd.Flags()
	for _, name := range []string{"output", "resolution", "vbitrate", "abitrate", "fps", "codec"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found", name)
		}
	}
}

func TestVideoSupportedFormats(t *testing.T) {
	expected := []string{
		"mp4", "mkv", "avi", "mov", "webm", "flv", "wmv", "mpeg", "ts", "m4v",
		"av1", "hevc", "divx", "xvid", "f4v", "asf", "m2v", "mjpeg", "tod",
	}
	for _, fmt := range expected {
		if !supportedFormats[fmt] {
			t.Errorf("format %s should be supported", fmt)
		}
	}
}

package trimvideo

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestTrimVideoCmd_Metadata(t *testing.T) {
	if trimVideoCmd.Use != "trimVideo <input-file>" {
		t.Errorf("unexpected Use: %s", trimVideoCmd.Use)
	}
	if trimVideoCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestTrimVideoCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(trimVideoCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := v(trimVideoCmd, []string{"input.mp4"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTrimVideoCmd_Flags(t *testing.T) {
	f := trimVideoCmd.Flags()
	for _, name := range []string{"output", "start", "end", "duration", "codec", "quiet", "json"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found", name)
		}
	}
}

func TestTrimVideoCmd_ShortFlags(t *testing.T) {
	f := trimVideoCmd.Flags()
	for _, short := range []string{"o", "s", "e", "d", "q", "j"} {
		if f.ShorthandLookup(short) == nil {
			t.Errorf("-%s shorthand not found", short)
		}
	}
}

func TestCodecOrCopy_Default(t *testing.T) {
	codec = ""
	if got := codecOrCopy(); got != "copy" {
		t.Errorf("got %q, want copy", got)
	}
}

func TestCodecOrCopy_Custom(t *testing.T) {
	codec = "libx264"
	defer func() { codec = "" }()
	if got := codecOrCopy(); got != "libx264" {
		t.Errorf("got %q, want libx264", got)
	}
}

func TestDefaultOutputName(t *testing.T) {
	// Test that when output is empty, the command would generate input_trimmed.ext
	// We can't test the full run without ffmpeg, but we can verify the flag default
	f := trimVideoCmd.Flags()
	outFlag := f.Lookup("output")
	if outFlag.DefValue != "" {
		t.Errorf("output default: got %q, want empty", outFlag.DefValue)
	}
}

func TestStartEndFlags(t *testing.T) {
	f := trimVideoCmd.Flags()
	if f.Lookup("start").DefValue != "" {
		t.Error("start should default to empty")
	}
	if f.Lookup("end").DefValue != "" {
		t.Error("end should default to empty")
	}
	if f.Lookup("duration").DefValue != "" {
		t.Error("duration should default to empty")
	}
}

package convertaudio

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestConvertAudioCmd_Metadata(t *testing.T) {
	if convertAudioCmd.Use != "convertAudio <input-file>" {
		t.Errorf("unexpected Use: %s", convertAudioCmd.Use)
	}
	if convertAudioCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertAudioCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(convertAudioCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := v(convertAudioCmd, []string{"in.wav"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConvertAudioCmd_Flags(t *testing.T) {
	f := convertAudioCmd.Flags()
	for _, name := range []string{"output", "bitrate", "sample"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found", name)
		}
	}
	if f.ShorthandLookup("o") == nil {
		t.Error("-o shorthand not found")
	}
}

func TestSupportedFormats(t *testing.T) {
	expected := []string{"mp3", "wav", "flac", "aac", "ogg", "m4a", "wma", "opus"}
	for _, fmt := range expected {
		if !supportedFormats[fmt] {
			t.Errorf("format %s should be supported", fmt)
		}
	}
}

func TestUnsupportedFormat(t *testing.T) {
	if supportedFormats["xyz"] {
		t.Error("xyz should not be supported")
	}
}

package convertfont

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestConvertFontCmd_Metadata(t *testing.T) {
	if convertFontCmd.Use != "convertFont <input-file>" {
		t.Errorf("unexpected Use: %s", convertFontCmd.Use)
	}
	if convertFontCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertFontCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(convertFontCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
}

func TestConvertFontCmd_Flags(t *testing.T) {
	if convertFontCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
}

func TestFontSupportedFormats(t *testing.T) {
	expected := []string{"ttf", "otf", "woff", "woff2", "eot"}
	for _, fmt := range expected {
		if !supportedFormats[fmt] {
			t.Errorf("format %s should be supported", fmt)
		}
	}
}

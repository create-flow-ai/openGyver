package convertpresentation

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestConvertPresentationCmd_Metadata(t *testing.T) {
	if convertPresentationCmd.Use != "convertPresentation <input-file>" {
		t.Errorf("unexpected Use: %s", convertPresentationCmd.Use)
	}
	if convertPresentationCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertPresentationCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(convertPresentationCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
}

func TestConvertPresentationCmd_Flags(t *testing.T) {
	if convertPresentationCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
}

func TestPresentationSupportedFormats(t *testing.T) {
	expected := []string{"pptx", "ppt", "odp", "pdf", "key", "ppsx", "potx", "html", "png"}
	for _, fmt := range expected {
		if !supportedFormats[fmt] {
			t.Errorf("format %s should be supported", fmt)
		}
	}
}

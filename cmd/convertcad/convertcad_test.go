package convertcad

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestConvertCADCmd_Metadata(t *testing.T) {
	if convertCADCmd.Use != "convertCAD <input-file>" {
		t.Errorf("unexpected Use: %s", convertCADCmd.Use)
	}
	if convertCADCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertCADCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(convertCADCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
}

func TestConvertCADCmd_Flags(t *testing.T) {
	if convertCADCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
}

func TestCADSupportedFormats(t *testing.T) {
	expected := []string{"dwg", "dxf", "dwf", "pdf", "svg", "png"}
	for _, fmt := range expected {
		if !supportedFormats[fmt] {
			t.Errorf("format %s should be supported", fmt)
		}
	}
}

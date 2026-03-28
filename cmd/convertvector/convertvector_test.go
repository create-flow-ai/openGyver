package convertvector

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestConvertVectorCmd_Metadata(t *testing.T) {
	if convertVectorCmd.Use != "convertVector <input-file>" {
		t.Errorf("unexpected Use: %s", convertVectorCmd.Use)
	}
	if convertVectorCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertVectorCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(convertVectorCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
}

func TestConvertVectorCmd_Flags(t *testing.T) {
	f := convertVectorCmd.Flags()
	for _, name := range []string{"output", "width", "height"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found", name)
		}
	}
}

func TestVectorSupportedFormats(t *testing.T) {
	expected := []string{
		"svg", "svgz", "eps", "pdf", "emf", "wmf", "ai", "cdr", "png", "jpg",
		"ccx", "cdt", "cmx", "dst", "exp", "fig", "pes", "plt", "sk", "sk1",
	}
	for _, fmt := range expected {
		if !supportedFormats[fmt] {
			t.Errorf("format %s should be supported", fmt)
		}
	}
}

func TestIsRasterOutput(t *testing.T) {
	if !isRasterOutput("png") {
		t.Error("png should be raster")
	}
	if !isRasterOutput("jpg") {
		t.Error("jpg should be raster")
	}
	if isRasterOutput("svg") {
		t.Error("svg should not be raster")
	}
	if isRasterOutput("pdf") {
		t.Error("pdf should not be raster")
	}
}

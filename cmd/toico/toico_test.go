package toico

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestToIcoCmd_Metadata(t *testing.T) {
	if toicoCmd.Use != "toIco <image>" {
		t.Errorf("unexpected Use: %s", toicoCmd.Use)
	}
	if toicoCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if toicoCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestToIcoCmd_RequiresExactlyOneArg(t *testing.T) {
	validator := cobra.ExactArgs(1)

	if err := validator(toicoCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(toicoCmd, []string{"a.png", "b.png"}); err == nil {
		t.Error("expected error with two args")
	}
	if err := validator(toicoCmd, []string{"a.png"}); err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}
}

func TestToIcoCmd_DefaultFlags(t *testing.T) {
	f := toicoCmd.Flags()

	outFlag := f.Lookup("output")
	if outFlag == nil {
		t.Fatal("--output flag not found")
	}
	if outFlag.DefValue != "output.ico" {
		t.Errorf("output default: got %q, want %q", outFlag.DefValue, "output.ico")
	}

	sizesFlag := f.Lookup("sizes")
	if sizesFlag == nil {
		t.Fatal("--sizes flag not found")
	}
	if sizesFlag.DefValue != "[16,32,48,256]" {
		t.Errorf("sizes default: got %q, want %q", sizesFlag.DefValue, "[16,32,48,256]")
	}

	formatFlag := f.Lookup("format")
	if formatFlag == nil {
		t.Fatal("--format flag not found")
	}
	if formatFlag.DefValue != "" {
		t.Errorf("format default: got %q, want empty", formatFlag.DefValue)
	}

	widthFlag := f.Lookup("width")
	if widthFlag == nil {
		t.Fatal("--width flag not found")
	}
	if widthFlag.DefValue != "0" {
		t.Errorf("width default: got %q, want %q", widthFlag.DefValue, "0")
	}

	heightFlag := f.Lookup("height")
	if heightFlag == nil {
		t.Fatal("--height flag not found")
	}
	if heightFlag.DefValue != "0" {
		t.Errorf("height default: got %q, want %q", heightFlag.DefValue, "0")
	}
}

func TestToIcoCmd_ShortFlags(t *testing.T) {
	f := toicoCmd.Flags()

	outShort := f.ShorthandLookup("o")
	if outShort == nil {
		t.Error("-o shorthand not found for --output")
	}

	fmtShort := f.ShorthandLookup("f")
	if fmtShort == nil {
		t.Error("-f shorthand not found for --format")
	}
}

func TestToIcoCmd_RunE(t *testing.T) {
	err := toicoCmd.RunE(toicoCmd, []string{"test.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToIcoCmd_RunE_WithDimensions(t *testing.T) {
	// Set dimension flags
	width = 512
	height = 512
	defer func() { width = 0; height = 0 }()

	err := toicoCmd.RunE(toicoCmd, []string{"test.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToIcoCmd_RunE_WithFormat(t *testing.T) {
	format = "jpeg"
	defer func() { format = "" }()

	err := toicoCmd.RunE(toicoCmd, []string{"test.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

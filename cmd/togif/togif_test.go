package togif

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestToGifCmd_Metadata(t *testing.T) {
	if togifCmd.Use != "toGif <image> [image...]" {
		t.Errorf("unexpected Use: %s", togifCmd.Use)
	}
	if togifCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if togifCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestToGifCmd_RequiresAtLeastOneArg(t *testing.T) {
	validator := cobra.MinimumNArgs(1)

	if err := validator(togifCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := validator(togifCmd, []string{"a.png"}); err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}
	if err := validator(togifCmd, []string{"a.png", "b.png", "c.png"}); err != nil {
		t.Errorf("unexpected error with three args: %v", err)
	}
}

func TestToGifCmd_DefaultFlags(t *testing.T) {
	f := togifCmd.Flags()

	outFlag := f.Lookup("output")
	if outFlag == nil {
		t.Fatal("--output flag not found")
	}
	if outFlag.DefValue != "output.gif" {
		t.Errorf("output default: got %q, want %q", outFlag.DefValue, "output.gif")
	}

	delayFlag := f.Lookup("delay")
	if delayFlag == nil {
		t.Fatal("--delay flag not found")
	}
	if delayFlag.DefValue != "100" {
		t.Errorf("delay default: got %q, want %q", delayFlag.DefValue, "100")
	}

	loopFlag := f.Lookup("loop")
	if loopFlag == nil {
		t.Fatal("--loop flag not found")
	}
	if loopFlag.DefValue != "0" {
		t.Errorf("loop default: got %q, want %q", loopFlag.DefValue, "0")
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

func TestToGifCmd_ShortFlags(t *testing.T) {
	f := togifCmd.Flags()

	outShort := f.ShorthandLookup("o")
	if outShort == nil {
		t.Error("-o shorthand not found for --output")
	}

	fmtShort := f.ShorthandLookup("f")
	if fmtShort == nil {
		t.Error("-f shorthand not found for --format")
	}
}

func TestToGifCmd_RunE_OneFrame(t *testing.T) {
	err := togifCmd.RunE(togifCmd, []string{"frame1.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToGifCmd_RunE_MultipleFrames(t *testing.T) {
	err := togifCmd.RunE(togifCmd, []string{"frame1.png", "frame2.png", "frame3.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToGifCmd_RunE_WithDimensions(t *testing.T) {
	width = 320
	height = 240
	defer func() { width = 0; height = 0 }()

	err := togifCmd.RunE(togifCmd, []string{"frame1.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToGifCmd_RunE_WithFormat(t *testing.T) {
	format = "bmp"
	defer func() { format = "" }()

	err := togifCmd.RunE(togifCmd, []string{"frame1.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestToGifCmd_RunE_WithDelayAndLoop(t *testing.T) {
	delay = 50
	loop = 3
	defer func() { delay = 100; loop = 0 }()

	err := togifCmd.RunE(togifCmd, []string{"frame1.png", "frame2.png"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

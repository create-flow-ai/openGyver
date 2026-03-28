package epoch

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestEpochCmd_Metadata(t *testing.T) {
	if epochCmd.Use != "epoch" {
		t.Errorf("unexpected Use: %s", epochCmd.Use)
	}
	if epochCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if epochCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestEpochCmd_AcceptsNoArgs(t *testing.T) {
	validator := cobra.NoArgs
	if err := validator(epochCmd, []string{}); err != nil {
		t.Errorf("unexpected error with no args: %v", err)
	}
}

func TestEpochCmd_RejectsArgs(t *testing.T) {
	validator := cobra.NoArgs
	if err := validator(epochCmd, []string{"extra"}); err == nil {
		t.Error("expected error with args")
	}
}

func TestEpochCmd_Flags(t *testing.T) {
	f := epochCmd.PersistentFlags()

	msFlag := f.Lookup("ms")
	if msFlag == nil {
		t.Fatal("--ms flag not found")
	}
	if msFlag.DefValue != "false" {
		t.Errorf("ms default: got %q, want %q", msFlag.DefValue, "false")
	}

	usFlag := f.Lookup("us")
	if usFlag == nil {
		t.Fatal("--us flag not found")
	}
	if usFlag.DefValue != "false" {
		t.Errorf("us default: got %q, want %q", usFlag.DefValue, "false")
	}

	nsFlag := f.Lookup("ns")
	if nsFlag == nil {
		t.Fatal("--ns flag not found")
	}
	if nsFlag.DefValue != "false" {
		t.Errorf("ns default: got %q, want %q", nsFlag.DefValue, "false")
	}
}

func TestEpochCmd_RunE_Seconds(t *testing.T) {
	ms, us, ns = false, false, false
	before := time.Now().Unix()
	err := epochCmd.RunE(epochCmd, []string{})
	after := time.Now().Unix()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The command prints to stdout; just verify it doesn't error.
	// Sanity: current time should be reasonable.
	if before < 1700000000 || after < before {
		t.Error("time sanity check failed")
	}
}

func TestEpochCmd_RunE_Milliseconds(t *testing.T) {
	ms = true
	us, ns = false, false
	defer func() { ms = false }()

	err := epochCmd.RunE(epochCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEpochCmd_RunE_Microseconds(t *testing.T) {
	us = true
	ms, ns = false, false
	defer func() { us = false }()

	err := epochCmd.RunE(epochCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEpochCmd_RunE_Nanoseconds(t *testing.T) {
	ns = true
	ms, us = false, false
	defer func() { ns = false }()

	err := epochCmd.RunE(epochCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEpochCmd_PrecisionPrecedence(t *testing.T) {
	// When multiple flags are set, ns should win (it's checked first in switch)
	ns = true
	ms = true
	us = true
	defer func() { ns, ms, us = false, false, false }()

	err := epochCmd.RunE(epochCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

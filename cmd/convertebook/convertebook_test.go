package convertebook

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestConvertEbookCmd_Metadata(t *testing.T) {
	if convertEbookCmd.Use != "convertEbook <input-file>" {
		t.Errorf("unexpected Use: %s", convertEbookCmd.Use)
	}
	if convertEbookCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestConvertEbookCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(convertEbookCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
}

func TestConvertEbookCmd_Flags(t *testing.T) {
	if convertEbookCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
	if convertEbookCmd.Flags().ShorthandLookup("o") == nil {
		t.Error("-o shorthand not found")
	}
}

func TestEbookSupportedFormats(t *testing.T) {
	expected := []string{"epub", "mobi", "azw3", "pdf", "fb2", "txt", "html", "docx", "cbz", "cbr"}
	for _, fmt := range expected {
		if !supportedFormats[fmt] {
			t.Errorf("format %s should be supported", fmt)
		}
	}
}

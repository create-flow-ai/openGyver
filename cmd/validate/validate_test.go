package validate

import (
	"strings"
	"testing"
)

func TestValidateHTML_Valid(t *testing.T) {
	input := `<!DOCTYPE html><html><head><title>Test</title></head><body><p>Hello</p></body></html>`
	errs := validateHTML(input)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid HTML, got: %v", errs)
	}
}

func TestValidateHTML_MissingDoctype(t *testing.T) {
	input := `<html><head><title>Test</title></head><body><p>Hello</p></body></html>`
	errs := validateHTML(input)

	found := false
	for _, e := range errs {
		if strings.Contains(strings.ToLower(e), "doctype") {
			found = true
		}
	}
	if !found {
		t.Error("expected missing DOCTYPE error")
	}
}

func TestValidateHTML_MissingAlt(t *testing.T) {
	input := `<!DOCTYPE html><html><body><img src="test.jpg"></body></html>`
	errs := validateHTML(input)

	found := false
	for _, e := range errs {
		if strings.Contains(e, "alt") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected missing alt attribute error, got: %v", errs)
	}
}

func TestValidateHTML_ImgWithAlt(t *testing.T) {
	input := `<!DOCTYPE html><html><body><img src="test.jpg" alt="A test image"></body></html>`
	errs := validateHTML(input)

	for _, e := range errs {
		if strings.Contains(e, "alt") {
			t.Errorf("unexpected alt error when alt is present: %s", e)
		}
	}
}

func TestValidateHTML_UnclosedTag(t *testing.T) {
	input := `<!DOCTYPE html><html><body><p>Hello</body></html>`
	errs := validateHTML(input)

	found := false
	for _, e := range errs {
		if strings.Contains(e, "unclosed") || strings.Contains(e, "mismatched") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected unclosed/mismatched tag error, got: %v", errs)
	}
}

func TestValidateCSV_Valid(t *testing.T) {
	input := "name,age,city\nAlice,30,NYC\nBob,25,LA"
	errs := validateCSV(input)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid CSV, got: %v", errs)
	}
}

func TestValidateCSV_InconsistentColumns(t *testing.T) {
	input := "name,age,city\nAlice,30\nBob,25,LA"
	errs := validateCSV(input)
	if len(errs) == 0 {
		t.Error("expected error for inconsistent column count")
	}
}

func TestValidateCSV_SingleRow(t *testing.T) {
	input := "name,age"
	errs := validateCSV(input)
	if len(errs) != 0 {
		t.Errorf("expected no errors for header-only CSV, got: %v", errs)
	}
}

func TestValidateXML_WellFormed(t *testing.T) {
	input := "<root><item id=\"1\">Hello</item></root>"
	errs := validateXML(input)
	if len(errs) != 0 {
		t.Errorf("expected no errors for well-formed XML, got: %v", errs)
	}
}

func TestValidateXML_Malformed(t *testing.T) {
	input := "<root><item>Hello</root>"
	errs := validateXML(input)
	if len(errs) == 0 {
		t.Error("expected error for malformed XML")
	}
}

func TestValidateXML_Unclosed(t *testing.T) {
	input := "<root><item>"
	errs := validateXML(input)
	if len(errs) == 0 {
		t.Error("expected error for unclosed XML elements")
	}
}

func TestValidateYAML_Valid(t *testing.T) {
	input := "name: hello\nage: 30\nitems:\n  - one\n  - two"
	errs := validateYAML(input)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid YAML, got: %v", errs)
	}
}

func TestValidateYAML_Invalid(t *testing.T) {
	input := "name: hello\n  bad indent: oops\n  : missing key"
	errs := validateYAML(input)
	if len(errs) == 0 {
		t.Error("expected error for invalid YAML")
	}
}

func TestValidateTOML_Valid(t *testing.T) {
	input := `title = "Test"
[database]
server = "localhost"
port = 5432`
	errs := validateTOML(input)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid TOML, got: %v", errs)
	}
}

func TestValidateTOML_Invalid(t *testing.T) {
	input := `[bad
key = `
	errs := validateTOML(input)
	if len(errs) == 0 {
		t.Error("expected error for invalid TOML")
	}
}

func TestResolveInput_FromArgs(t *testing.T) {
	text, err := resolveInput([]string{"hello world"})
	if err != nil {
		t.Fatalf("resolveInput error: %v", err)
	}
	if text != "hello world" {
		t.Errorf("resolveInput = %q", text)
	}
}

func TestResolveInput_NoInput(t *testing.T) {
	origFile := inputFile
	inputFile = ""
	defer func() { inputFile = origFile }()

	_, err := resolveInput(nil)
	if err == nil {
		t.Error("expected error when no input provided")
	}
}

func TestPrintResult_NoErrors(t *testing.T) {
	// Just verify it doesn't panic with empty error list.
	err := printResult([]string{})
	if err != nil {
		t.Errorf("printResult error: %v", err)
	}
}

func TestVoidElements(t *testing.T) {
	// Verify a few known void elements.
	voids := []string{"img", "br", "hr", "input", "meta", "link"}
	for _, v := range voids {
		if !voidElements[v] {
			t.Errorf("expected %q to be a void element", v)
		}
	}
	// Non-void.
	if voidElements["div"] {
		t.Error("div should not be a void element")
	}
}

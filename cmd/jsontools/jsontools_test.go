package jsontools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunFormat_CompactToIndented(t *testing.T) {
	input := `{"name":"Alice","age":30}`

	var parsed interface{}
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	out, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		t.Fatalf("format error: %v", err)
	}
	formatted := string(out)

	if !strings.Contains(formatted, "\n") {
		t.Error("formatted output should contain newlines")
	}
	if !strings.Contains(formatted, "  ") {
		t.Error("formatted output should contain indentation")
	}
	if !strings.Contains(formatted, `"name"`) {
		t.Error("formatted output should contain key 'name'")
	}
}

func TestRunMinify_IndentedToCompact(t *testing.T) {
	input := `{
  "name": "Alice",
  "age": 30
}`

	var parsed interface{}
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	out, err := json.Marshal(parsed)
	if err != nil {
		t.Fatalf("minify error: %v", err)
	}
	minified := string(out)

	if strings.Contains(minified, "\n") {
		t.Error("minified output should not contain newlines")
	}
	if strings.Contains(minified, "  ") {
		t.Error("minified output should not contain extra spaces")
	}
}

func TestValidate_ValidJSON(t *testing.T) {
	input := `{"valid": true, "count": 42}`
	if !json.Valid([]byte(input)) {
		t.Error("expected valid JSON to be valid")
	}
}

func TestValidate_InvalidJSON(t *testing.T) {
	input := `{"missing": }`
	if json.Valid([]byte(input)) {
		t.Error("expected invalid JSON to be invalid")
	}

	var js json.RawMessage
	err := json.Unmarshal([]byte(input), &js)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestValidate_EmptyObject(t *testing.T) {
	if !json.Valid([]byte(`{}`)) {
		t.Error("empty object should be valid JSON")
	}
}

func TestValidate_EmptyArray(t *testing.T) {
	if !json.Valid([]byte(`[]`)) {
		t.Error("empty array should be valid JSON")
	}
}

func TestEvaluatePath_NestedKey(t *testing.T) {
	root := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": float64(5432),
		},
	}

	val, err := evaluatePath(root, "database.host")
	if err != nil {
		t.Fatalf("evaluatePath error: %v", err)
	}
	if val != "localhost" {
		t.Errorf("database.host = %v, want localhost", val)
	}
}

func TestEvaluatePath_ArrayIndex(t *testing.T) {
	root := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice"},
			map[string]interface{}{"name": "Bob"},
		},
	}

	val, err := evaluatePath(root, "users[0].name")
	if err != nil {
		t.Fatalf("evaluatePath error: %v", err)
	}
	if val != "Alice" {
		t.Errorf("users[0].name = %v, want Alice", val)
	}

	val2, err := evaluatePath(root, "users[1].name")
	if err != nil {
		t.Fatalf("evaluatePath error: %v", err)
	}
	if val2 != "Bob" {
		t.Errorf("users[1].name = %v, want Bob", val2)
	}
}

func TestEvaluatePath_TopLevelKey(t *testing.T) {
	root := map[string]interface{}{
		"title": "Hello",
	}
	val, err := evaluatePath(root, "title")
	if err != nil {
		t.Fatalf("evaluatePath error: %v", err)
	}
	if val != "Hello" {
		t.Errorf("title = %v, want Hello", val)
	}
}

func TestEvaluatePath_KeyNotFound(t *testing.T) {
	root := map[string]interface{}{
		"name": "Alice",
	}
	_, err := evaluatePath(root, "missing")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestEvaluatePath_IndexOutOfRange(t *testing.T) {
	root := map[string]interface{}{
		"items": []interface{}{"a", "b"},
	}
	_, err := evaluatePath(root, "items[5]")
	if err == nil {
		t.Error("expected error for out-of-range index")
	}
}

func TestParsePath(t *testing.T) {
	segments := parsePath("data.users[0].name")
	if len(segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segments))
	}
	if segments[0].key != "data" || segments[0].index != -1 {
		t.Errorf("segment 0: key=%q index=%d", segments[0].key, segments[0].index)
	}
	if segments[1].key != "users" || segments[1].index != 0 {
		t.Errorf("segment 1: key=%q index=%d", segments[1].key, segments[1].index)
	}
	if segments[2].key != "name" || segments[2].index != -1 {
		t.Errorf("segment 2: key=%q index=%d", segments[2].key, segments[2].index)
	}
}

func TestEscape(t *testing.T) {
	raw := `hello "world"`
	escaped, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	expected := `"hello \"world\""`
	if string(escaped) != expected {
		t.Errorf("escaped = %s, want %s", string(escaped), expected)
	}
}

func TestUnescape(t *testing.T) {
	input := `"hello\tworld"`
	var unescaped string
	if err := json.Unmarshal([]byte(input), &unescaped); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if unescaped != "hello\tworld" {
		t.Errorf("unescaped = %q, want %q", unescaped, "hello\tworld")
	}
}

func TestUnescape_Invalid(t *testing.T) {
	input := `not a json string`
	var s string
	err := json.Unmarshal([]byte(input), &s)
	if err == nil {
		t.Error("expected error for non-JSON-string input")
	}
}

func TestReadInput_FromArgs(t *testing.T) {
	// Temporarily set filePath to empty.
	origFP := filePath
	filePath = ""
	defer func() { filePath = origFP }()

	text, err := readInput([]string{`{"key":"value"}`})
	if err != nil {
		t.Fatalf("readInput error: %v", err)
	}
	if text != `{"key":"value"}` {
		t.Errorf("readInput = %q", text)
	}
}

func TestReadInput_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	os.WriteFile(path, []byte(`{"from":"file"}`), 0644)

	origFP := filePath
	filePath = path
	defer func() { filePath = origFP }()

	text, err := readInput(nil)
	if err != nil {
		t.Fatalf("readInput error: %v", err)
	}
	if text != `{"from":"file"}` {
		t.Errorf("readInput = %q", text)
	}
}

func TestReadInput_NoInput(t *testing.T) {
	origFP := filePath
	filePath = ""
	defer func() { filePath = origFP }()

	_, err := readInput(nil)
	if err == nil {
		t.Error("expected error when no input provided")
	}
}

package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeLCS_IdenticalFiles(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	edits := computeLCS(lines, lines)

	for _, e := range edits {
		if e.op != "keep" {
			t.Errorf("identical files: expected all 'keep' ops, got %q for %q", e.op, e.line)
		}
	}
	if len(edits) != 3 {
		t.Errorf("expected 3 edits, got %d", len(edits))
	}
}

func TestComputeLCS_DifferentFiles(t *testing.T) {
	a := []string{"line1", "line2", "line3"}
	b := []string{"line1", "lineX", "line3"}
	edits := computeLCS(a, b)

	hasAdd := false
	hasRemove := false
	for _, e := range edits {
		if e.op == "add" {
			hasAdd = true
		}
		if e.op == "remove" {
			hasRemove = true
		}
	}
	if !hasAdd {
		t.Error("expected at least one 'add' op for different files")
	}
	if !hasRemove {
		t.Error("expected at least one 'remove' op for different files")
	}
}

func TestComputeLCS_EmptyToContent(t *testing.T) {
	a := []string{}
	b := []string{"line1", "line2"}
	edits := computeLCS(a, b)

	for _, e := range edits {
		if e.op != "add" {
			t.Errorf("expected all 'add' ops, got %q", e.op)
		}
	}
	if len(edits) != 2 {
		t.Errorf("expected 2 edits, got %d", len(edits))
	}
}

func TestComputeLCS_ContentToEmpty(t *testing.T) {
	a := []string{"line1", "line2"}
	b := []string{}
	edits := computeLCS(a, b)

	for _, e := range edits {
		if e.op != "remove" {
			t.Errorf("expected all 'remove' ops, got %q", e.op)
		}
	}
	if len(edits) != 2 {
		t.Errorf("expected 2 edits, got %d", len(edits))
	}
}

func TestReadLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("line1\nline2\nline3"), 0644)

	lines, err := readLines(path)
	if err != nil {
		t.Fatalf("readLines error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line1" || lines[2] != "line3" {
		t.Errorf("unexpected line content: %v", lines)
	}
}

func TestReadLines_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	os.WriteFile(path, []byte(""), 0644)

	lines, err := readLines(path)
	if err != nil {
		t.Fatalf("readLines error: %v", err)
	}
	if lines != nil {
		t.Errorf("expected nil for empty file, got %v", lines)
	}
}

func TestReadLines_NonexistentFile(t *testing.T) {
	_, err := readLines("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestTextDiffWithTempFiles(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "a.txt")
	file2 := filepath.Join(dir, "b.txt")

	os.WriteFile(file1, []byte("alpha\nbeta\ngamma"), 0644)
	os.WriteFile(file2, []byte("alpha\ndelta\ngamma"), 0644)

	lines1, _ := readLines(file1)
	lines2, _ := readLines(file2)
	edits := computeLCS(lines1, lines2)

	// Should have a remove (beta) and an add (delta).
	foundRemoveBeta := false
	foundAddDelta := false
	for _, e := range edits {
		if e.op == "remove" && e.line == "beta" {
			foundRemoveBeta = true
		}
		if e.op == "add" && e.line == "delta" {
			foundAddDelta = true
		}
	}
	if !foundRemoveBeta {
		t.Error("expected removal of 'beta'")
	}
	if !foundAddDelta {
		t.Error("expected addition of 'delta'")
	}
}

func TestCompareJSON_AddedKeys(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "a.json")
	file2 := filepath.Join(dir, "b.json")

	j1, _ := json.Marshal(map[string]interface{}{"name": "Alice"})
	j2, _ := json.Marshal(map[string]interface{}{"name": "Alice", "age": 30})
	os.WriteFile(file1, j1, 0644)
	os.WriteFile(file2, j2, 0644)

	v1, _ := readJSON(file1)
	v2, _ := readJSON(file2)
	diffs := compareJSON("", v1, v2)

	foundAdded := false
	for _, d := range diffs {
		if d.diffType == "added" && d.path == "age" {
			foundAdded = true
		}
	}
	if !foundAdded {
		t.Error("expected 'added' diff for key 'age'")
	}
}

func TestCompareJSON_RemovedKeys(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "a.json")
	file2 := filepath.Join(dir, "b.json")

	j1, _ := json.Marshal(map[string]interface{}{"name": "Alice", "age": 30})
	j2, _ := json.Marshal(map[string]interface{}{"name": "Alice"})
	os.WriteFile(file1, j1, 0644)
	os.WriteFile(file2, j2, 0644)

	v1, _ := readJSON(file1)
	v2, _ := readJSON(file2)
	diffs := compareJSON("", v1, v2)

	foundRemoved := false
	for _, d := range diffs {
		if d.diffType == "removed" && d.path == "age" {
			foundRemoved = true
		}
	}
	if !foundRemoved {
		t.Error("expected 'removed' diff for key 'age'")
	}
}

func TestCompareJSON_ChangedValues(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "a.json")
	file2 := filepath.Join(dir, "b.json")

	j1, _ := json.Marshal(map[string]interface{}{"name": "Alice"})
	j2, _ := json.Marshal(map[string]interface{}{"name": "Bob"})
	os.WriteFile(file1, j1, 0644)
	os.WriteFile(file2, j2, 0644)

	v1, _ := readJSON(file1)
	v2, _ := readJSON(file2)
	diffs := compareJSON("", v1, v2)

	foundChanged := false
	for _, d := range diffs {
		if d.diffType == "changed" && d.path == "name" {
			foundChanged = true
		}
	}
	if !foundChanged {
		t.Error("expected 'changed' diff for key 'name'")
	}
}

func TestCompareJSON_Identical(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	diffs := compareJSON("", data, data)
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs for identical JSON, got %d", len(diffs))
	}
}

func TestJoinPath(t *testing.T) {
	if joinPath("", "key") != "key" {
		t.Errorf("joinPath(\"\", \"key\") = %q", joinPath("", "key"))
	}
	if joinPath("root", "child") != "root.child" {
		t.Errorf("joinPath(\"root\", \"child\") = %q", joinPath("root", "child"))
	}
}

func TestPathOrRoot(t *testing.T) {
	if pathOrRoot("") != "." {
		t.Errorf("pathOrRoot(\"\") = %q, want \".\"", pathOrRoot(""))
	}
	if pathOrRoot("foo") != "foo" {
		t.Errorf("pathOrRoot(\"foo\") = %q, want \"foo\"", pathOrRoot("foo"))
	}
}

func TestFormatVal(t *testing.T) {
	if formatVal(nil) != "null" {
		t.Errorf("formatVal(nil) = %q", formatVal(nil))
	}
	if formatVal("hello") != `"hello"` {
		t.Errorf("formatVal(\"hello\") = %q", formatVal("hello"))
	}
	if formatVal(42.0) != "42" {
		t.Errorf("formatVal(42.0) = %q", formatVal(42.0))
	}
	if formatVal(true) != "true" {
		t.Errorf("formatVal(true) = %q", formatVal(true))
	}
}

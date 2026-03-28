package dataformat

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestYAMLToJSON(t *testing.T) {
	input := "name: hello\nage: 30"
	result, err := yamlToJSON(input)
	if err != nil {
		t.Fatalf("yamlToJSON error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if data["name"] != "hello" {
		t.Errorf("name = %v, want hello", data["name"])
	}
	// YAML numbers are parsed as int by yaml.v3, then normalized.
	if data["age"] != float64(30) {
		t.Errorf("age = %v (%T), want 30", data["age"], data["age"])
	}
}

func TestJSONToYAML(t *testing.T) {
	input := `{"name":"Alice","age":25}`
	result, err := jsonToYAML(input)
	if err != nil {
		t.Fatalf("jsonToYAML error: %v", err)
	}
	if !strings.Contains(result, "name: Alice") {
		t.Errorf("YAML output missing 'name: Alice': %s", result)
	}
	if !strings.Contains(result, "age: 25") {
		t.Errorf("YAML output missing 'age: 25': %s", result)
	}
}

func TestJSONToYAMLRoundtrip(t *testing.T) {
	original := `{"color":"blue","count":42}`
	yamlOut, err := jsonToYAML(original)
	if err != nil {
		t.Fatalf("jsonToYAML error: %v", err)
	}
	jsonOut, err := yamlToJSON(yamlOut)
	if err != nil {
		t.Fatalf("yamlToJSON error: %v", err)
	}

	var orig, result map[string]interface{}
	json.Unmarshal([]byte(original), &orig)
	json.Unmarshal([]byte(jsonOut), &result)

	if orig["color"] != result["color"] {
		t.Errorf("roundtrip color mismatch: %v vs %v", orig["color"], result["color"])
	}
	if orig["count"] != result["count"] {
		t.Errorf("roundtrip count mismatch: %v vs %v", orig["count"], result["count"])
	}
}

func TestCSVToJSON(t *testing.T) {
	input := "name,age,city\nAlice,30,NYC\nBob,25,LA"
	result, err := csvToJSON(input)
	if err != nil {
		t.Fatalf("csvToJSON error: %v", err)
	}

	var data []map[string]string
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid JSON array: %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(data))
	}
	if data[0]["name"] != "Alice" {
		t.Errorf("row 0 name = %q, want Alice", data[0]["name"])
	}
	if data[0]["age"] != "30" {
		t.Errorf("row 0 age = %q, want 30", data[0]["age"])
	}
	if data[1]["city"] != "LA" {
		t.Errorf("row 1 city = %q, want LA", data[1]["city"])
	}
}

func TestCSVToJSON_EmptyBody(t *testing.T) {
	input := "name,age"
	result, err := csvToJSON(input)
	if err != nil {
		t.Fatalf("csvToJSON error: %v", err)
	}
	if result != "[]" {
		t.Errorf("expected empty array, got %s", result)
	}
}

func TestTOMLToJSON(t *testing.T) {
	input := `title = "Test"
[database]
server = "localhost"
port = 5432`
	result, err := tomlToJSON(input)
	if err != nil {
		t.Fatalf("tomlToJSON error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if data["title"] != "Test" {
		t.Errorf("title = %v, want Test", data["title"])
	}
	db, ok := data["database"].(map[string]interface{})
	if !ok {
		t.Fatal("database is not a map")
	}
	if db["server"] != "localhost" {
		t.Errorf("database.server = %v, want localhost", db["server"])
	}
}

func TestXMLToJSON(t *testing.T) {
	input := `<root><name>Alice</name><age>30</age></root>`
	result, err := xmlToJSON(input)
	if err != nil {
		t.Fatalf("xmlToJSON error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	root, ok := data["root"].(map[string]interface{})
	if !ok {
		t.Fatal("root is not a map")
	}
	if root["name"] != "Alice" {
		t.Errorf("root.name = %v, want Alice", root["name"])
	}
}

func TestJSONToCSV(t *testing.T) {
	input := `[{"name":"Alice","age":"30"},{"name":"Bob","age":"25"}]`
	result, err := jsonToCSV(input)
	if err != nil {
		t.Fatalf("jsonToCSV error: %v", err)
	}
	if !strings.Contains(result, "Alice") || !strings.Contains(result, "Bob") {
		t.Errorf("CSV output missing data: %s", result)
	}
	// Should have header row.
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines (header + 2 rows), got %d", len(lines))
	}
}

func TestJSONToXML(t *testing.T) {
	input := `{"item":"hello"}`
	result, err := jsonToXML(input)
	if err != nil {
		t.Fatalf("jsonToXML error: %v", err)
	}
	if !strings.Contains(result, "<item>hello</item>") {
		t.Errorf("XML output missing expected element: %s", result)
	}
}

func TestJSONToTOML(t *testing.T) {
	input := `{"name":"test","value":42}`
	result, err := jsonToTOML(input)
	if err != nil {
		t.Fatalf("jsonToTOML error: %v", err)
	}
	if !strings.Contains(result, "name") || !strings.Contains(result, "test") {
		t.Errorf("TOML output missing expected content: %s", result)
	}
}

func TestNormalizeYAML(t *testing.T) {
	// Test with map[string]interface{}.
	input := map[string]interface{}{
		"key": "value",
		"nested": map[string]interface{}{
			"inner": "data",
		},
	}
	result := normalizeYAML(input)
	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map[string]interface{}")
	}
	if m["key"] != "value" {
		t.Errorf("key = %v, want value", m["key"])
	}
}

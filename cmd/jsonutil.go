package cmd

import (
	"encoding/json"
	"fmt"
)

// PrintJSON marshals v to indented JSON and prints it to stdout.
func PrintJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON encoding error: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

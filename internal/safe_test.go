/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"os"
	"testing"
)

func TestSaveToFile(t *testing.T) {
	// Define test file and content
	filePath := "test_output.txt"
	content := "Hello, test!"

	// Clean up after test
	defer os.Remove(filePath)

	// Call the function
	err := SaveToFile(content, filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read file to verify content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != content {
		t.Errorf("file content mismatch. expected %q, got %q", content, string(data))
	}
}

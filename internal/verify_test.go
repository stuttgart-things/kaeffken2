package internal

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Test: Existing file
	exists, err := FileExists(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error checking existing file: %v", err)
	}
	if !exists {
		t.Errorf("Expected file %s to exist", tmpFile.Name())
	}

	// Test: Non-existent file
	exists, err = FileExists("/non/existent/file")
	if err != nil {
		t.Fatalf("Unexpected error checking non-existent file: %v", err)
	}
	if exists {
		t.Error("Expected non-existent file check to return false")
	}
}

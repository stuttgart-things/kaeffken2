package internal

import (
	"fmt"
	"os"
)

// SaveToFile saves the provided content to the specified file path.
func SaveToFile(content string, filePath string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

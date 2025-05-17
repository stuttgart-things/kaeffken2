/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"fmt"
	"os"
	"regexp"
)

// SaveToFile saves the provided content to the specified file path.
func SaveToFile(content string, filePath string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

// CleanUpLines applies regex to fix the quoted config lines
func CleanUpLines(input string) string {
	re := regexp.MustCompile(`- '?"([^"]+)",?'`)
	return re.ReplaceAllString(input, `- "$1",`)
}

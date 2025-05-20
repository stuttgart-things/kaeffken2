/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var quoteRegex = regexp.MustCompile(`^["'](.*)["']$`)

// SaveToFile saves the provided content to the specified file path.
func SaveToFile(content string, filePath string) error {

	err := os.WriteFile(filePath, []byte(cleanUpLines(content)), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func cleanUpLines(input string) string {
	lines := strings.Split(input, "\n")
	re := regexp.MustCompile(`^(\s*)-\s*'?"([^"]+?)["']?,?'?\s*$`)

	for i, line := range lines {
		if match := re.FindStringSubmatch(line); match != nil {
			indent := match[1]
			content := match[2]
			lines[i] = indent + `- "` + content + `"`
		}
	}

	return strings.Join(lines, "\n")
}

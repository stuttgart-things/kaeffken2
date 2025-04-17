package internal

import (
	"regexp"
	"strings"
)

func CleanString(input string) string {
	// Remove triple-quoted key+-value patterns
	re := regexp.MustCompile(`'''[\w\-+]+'''`)
	cleaned := re.ReplaceAllString(input, "")

	// Remove any remaining single quotes
	cleaned = strings.ReplaceAll(cleaned, "'", "")

	// Remove lines that are just "-" or "- " or only whitespace
	lines := strings.Split(cleaned, "\n")
	var nonEmptyLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "-" && trimmed != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	return strings.Join(nonEmptyLines, "\n")
}

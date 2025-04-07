package internal

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/stuttgart-things/survey"

	"strings"
)

// ExtractQuestionsFromKCLFile parses a KCL file and extracts survey questions
// based on the special comment syntax
func ExtractQuestionsFromKCLFile(filename string) ([]*survey.Question, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read KCL file: %w", err)
	}

	return parseKCLQuestions(string(data))
}

// PARSE KCL-QUESTIONS EXTRACTS QUESTIONS FROM KCL FILE CONTENT
func parseKCLQuestions(content string) ([]*survey.Question, error) {
	var questions []*survey.Question

	// Regex to match KCL option declarations with special comments
	re := regexp.MustCompile(`_(\w+)\s*=\s*option\("([^"]+)"\)\s*or\s*"([^"]*)"\s*#\s*([^;]+)(?:;([^-+]*))?(?:-min(\d+))?(?:\+max(\d+))?`)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) < 5 { // Not a match
			continue
		}

		// matches[0] = full match
		// matches[1] = variable name
		// matches[2] = option name
		// matches[3] = default value
		// matches[4] = kind (ask/select)
		// matches[5] = options (for select)
		// matches[6] = minLength
		// matches[7] = maxLength

		question := &survey.Question{
			Name:    matches[2], // option name
			Default: matches[3], // default value
			Kind:    matches[4], // ask or select
		}

		// Set prompt (capitalize name)
		if len(question.Name) > 0 {
			question.Prompt = strings.ToUpper(question.Name) + "?"
		}

		// Handle select options if present
		if question.Kind == "select" && matches[5] != "" {
			question.Options = strings.Split(matches[5], ",")
			// Trim whitespace from each option
			for i, opt := range question.Options {
				question.Options[i] = strings.TrimSpace(opt)
			}
		}

		// Handle length constraints
		if matches[6] != "" {
			if min, err := strconv.Atoi(matches[6]); err == nil {
				question.MinLength = min
			}
		}
		if matches[7] != "" {
			if max, err := strconv.Atoi(matches[7]); err == nil {
				question.MaxLength = max
			}
		}

		// Set type based on kind
		if question.Kind == "ask" {
			question.Type = "string" // default type
		}

		questions = append(questions, question)
	}

	return questions, nil
}

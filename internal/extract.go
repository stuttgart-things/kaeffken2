package internal

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/stuttgart-things/survey"
	"gopkg.in/yaml.v3"

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

func ParseTemplateValues(values []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, item := range values {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			continue // or return an error if you want to be strict
		}

		key := parts[0]
		valStr := parts[1]

		// Try to convert string to int or bool if applicable
		if intVal, err := strconv.Atoi(valStr); err == nil {
			result[key] = intVal
		} else if boolVal, err := strconv.ParseBool(valStr); err == nil {
			result[key] = boolVal
		} else {
			result[key] = valStr
		}
	}

	return result
}

func ReadDictEntry(filePath, dictName, key string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var root map[string]interface{}
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("unmarshaling YAML: %w", err)
	}

	dicts, ok := root["dicts"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'dicts' not found or not a map")
	}

	dict, ok := dicts[dictName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' not found or not a map", dictName)
	}

	entry, ok := dict[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("entry '%s' not found or not a map", key)
	}

	return entry, nil
}

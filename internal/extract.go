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

func ExtractListsFromKCLFile(filename string) map[string]interface{} {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	content := parseKCLLists(string(data))

	return content
}

func parseKCLLists(content string) (allOptions map[string]interface{}) {

	allOptions = make(map[string]interface{})

	// Regex to match multiline list blocks
	listRe := regexp.MustCompile(`_(\w+)\s*=\s*option\("([^"]+)"\)\s*or\s*"""([\s\S]*?)"""\s*#\s*(list)`)

	// First, parse all multiline #list questions
	listMatches := listRe.FindAllStringSubmatch(content, -1)

	for _, match := range listMatches {

		// Extract list options
		rawBlock := match[3]

		fmt.Println(match[1])
		// fmt.Println("KEYYYYY", match[3])
		// fmt.Println("KEYYYYY", match[2])

		blockLines := strings.Split(rawBlock, "\n")
		var trimmedOptions []string

		for _, l := range blockLines {
			trimmed := strings.TrimSpace(l)
			if trimmed != "" {
				trimmedOptions = append(trimmedOptions, trimmed)
			}
		}

		allOptions[match[1]] = trimmedOptions
	}

	return allOptions
}

// PARSE KCL-QUESTIONS EXTRACTS QUESTIONS FROM KCL FILE CONTENT
func parseKCLQuestions(content string) ([]*survey.Question, error) {
	var questions []*survey.Question

	// Regex to match single-line ask/select options
	lineRe := regexp.MustCompile(`_(\w+)\s*=\s*option\("([^"]+)"\)\s*or\s*"([^"]*)"\s*#\s*([^;]+)(?:;([^-+]*))?(?:-min(\d+))?(?:\+max(\d+))?`)

	// Now parse all inline ask/select questions line-by-line
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if matches := lineRe.FindStringSubmatch(line); len(matches) >= 5 {
			question := &survey.Question{
				Name:    matches[2],
				Default: matches[3],
				Kind:    matches[4],
			}

			if len(question.Name) > 0 {
				question.Prompt = strings.ToUpper(question.Name) + "?"
			}

			if question.Kind == "select" && matches[5] != "" {
				question.Options = strings.Split(matches[5], ",")
				for i, opt := range question.Options {
					question.Options[i] = strings.TrimSpace(opt)
				}
			}

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

			if question.Kind == "ask" {
				question.Type = "string"
			}

			questions = append(questions, question)
		}
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

// ReadDicts reads the full "dicts" block into a map[string]interface{}
func ReadDicts(filePath, key string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var fullYAML map[string]interface{}
	if err := yaml.Unmarshal(data, &fullYAML); err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml: %w", err)
	}

	dicts, ok := fullYAML[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'dicts' not found or invalid format")
	}

	return dicts, nil
}

func GetValueFromDicts(dicts map[string]interface{}, DictKey, key string) map[string]interface{} {

	kinds := dicts[DictKey].(map[string]interface{})
	return ConvertToInterfaceMap(kinds[key])
}

func ConvertToInterfaceMap(val interface{}) map[string]interface{} {
	if m, ok := val.(map[string]interface{}); ok {
		return m
	}
	return nil
}

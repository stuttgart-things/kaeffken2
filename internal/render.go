/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"fmt"
	"log"
	"os"
	"regexp"

	kcl "kcl-lang.io/kcl-go"
)

func RenderKCL(kclFile string, allAnswers map[string]interface{}) string {

	// READ MAIN KCL FILE
	content, err := os.ReadFile(kclFile)
	if err != nil {
		log.Fatalf("Error reading KCL file: %v", err)
	}

	// OUTPUT ALL ANSWERS + MODIFY
	fmt.Println("RENDERING KCL")
	for key, value := range allAnswers {
		fmt.Printf("%s=%v\n", key, value)
	}

	values := convertToOptionStrings(allAnswers)

	fmt.Println("OPTS", values)
	//fmt.Println("OPTS", values)

	// options := []string{"name=kcl", fmt.Sprintf("cpu='1'")}

	// fmt.Println("OPTS2", options)

	// // Prepare KCL options with explicit key-value pairs
	opts := []kcl.Option{
		kcl.WithCode(string(content)),
		kcl.WithOptions(values...),
	}

	// Execute KCL
	result, err := kcl.Run(kclFile, opts...)
	if err != nil {
		log.Fatalf("KCL execution failed: %v", err)
	}

	// Output generated YAML
	return replaceTripleQuotes(result.GetRawYamlResult())
}

func convertToOptionStrings(answers map[string]interface{}) []string {
	var options []string

	for key, value := range answers {
		// Convert the value to a quoted string
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = fmt.Sprintf("'%s'", v)
		case int, int32, int64, float32, float64, bool:
			strValue = fmt.Sprintf("'%v'", v)
			// default:
			// 	strValue = fmt.Sprintf("'%v'", v) // Fallback for other types
		}

		options = append(options, fmt.Sprintf("%s=%s", key, strValue))
	}

	return options
}

// replaceTripleQuotes replaces ”'value”' with 'value' in a string
func replaceTripleQuotes(input string) string {
	// Updated regex to handle empty values
	re := regexp.MustCompile(`'''([^']*)'''`)
	return re.ReplaceAllString(input, `'$1'`)
}

// fixQuotesInMap processes a map replacing ”'value”' with 'value' in all values
func fixQuotesInMap(data map[string]string) map[string]string {
	re := regexp.MustCompile(`'''([^']*)'''`)
	result := make(map[string]string, len(data))

	for k, v := range data {
		result[k] = re.ReplaceAllString(v, `'$1'`)
	}
	return result
}

/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package modules

import (
	"fmt"
	"log"

	"github.com/stuttgart-things/kaeffken2/internal"
	"github.com/stuttgart-things/survey"
)

func ReadKCLQuestions(kclTemplatePath string) ([]*survey.Question, error) {
	questions, err := internal.ExtractQuestionsFromKCLFile(kclTemplatePath)
	if err != nil {
		// HANDLE ERROR
		log.Fatalf("Error extracting questions from KCL file: %v", err)
	}

	if len(questions) == 0 {
		fmt.Println("No questions found.")
	}

	return questions, nil
}

func ReadKCLList(kclTemplatePath string) map[string]interface{} {

	return internal.ExtractListsFromKCLFile(kclTemplatePath)
}

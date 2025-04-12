package modules

import (
	"fmt"

	"github.com/stuttgart-things/survey"
)

func SetAnswers(questions []*survey.Question) map[string]interface{} {

	allAnswers := make(map[string]interface{})

	// SET ANWERS TO ALL VALUES
	for _, question := range questions {
		allAnswers[question.Name] = question.Default
	}

	// OUTPUT ALL ANSWERS + MODIFY
	for key, value := range allAnswers {
		fmt.Printf("%s: %v\n", key, value)
	}

	return allAnswers
}

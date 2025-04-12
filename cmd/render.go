/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stuttgart-things/kaeffken2/internal"
	"github.com/stuttgart-things/kaeffken2/modules"

	"github.com/stuttgart-things/survey"
)

var (
	allAnswers   = make(map[string]interface{})
	kclTemplate  = "tests/proxmoxvm-template.k"
	renderedYAML string
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render templates",
	Long:  `Render templates based on profiles.`,
	Run: func(cmd *cobra.Command, args []string) {

		// LOAD THE QUESTIONS FROM A KCL FILE
		questions, err := modules.ReadKCLQuestions(kclTemplate)
		internal.CheckErr(err, "Error reading KCL questions")

		// BUILD THE SURVEY FORM AND GET A MAP FOR ANSWERS
		surveyForm, _, err := survey.BuildSurvey(questions)
		internal.CheckErr(err, "Error building survey")

		// RUN THE INTERACTIVE SURVEY
		err = surveyForm.Run()
		internal.CheckErr(err, "Error running survey")

		// SET ANWERS TO ALL VALUES
		allAnswers = modules.SetAnswers(questions)

		// RENDER KCL FILE TO YAML
		renderedYaml := internal.RenderKCL(kclTemplate, allAnswers)

		// INITIALIZE AND RUN THE TERMINAL EDITOR PROGRAM.
		renderedYaml = modules.RunEditor(renderedYaml)

		// SAVE DIALOG
		modules.SaveDialog(renderedYaml)

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

}

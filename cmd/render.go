/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stuttgart-things/kaeffken2/internal"
	"github.com/stuttgart-things/kaeffken2/modules"

	"github.com/stuttgart-things/survey"
)

var (
	allAnswers   = make(map[string]interface{})
	templatePath = "tests/proxmoxvm-template.k"
	renderedYAML string
	values       map[string]interface{}
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render templates",
	Long:  `Render templates based on profiles.`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})

		// log.Info().Str("version", "v1.2.3").Msg("starting CLI")

		// log.Info().Str("component", "cli").Msg("Application started")
		// log.Warn().Msg("This is a warning")

		// GET TEMPLATE PATH + CHECK EXISTENCE
		templatePath, _ := cmd.LocalFlags().GetString("template")
		exists, err := internal.FileExists(templatePath)
		fmt.Println("EXISTS", exists)
		fmt.Println("ERROR", err)

		internal.CheckErr(err, "ERROR READING KCL QUESTIONS")

		// GET/PARSE VALUES
		templateValues, _ := cmd.Flags().GetStringSlice("values")
		values = internal.ParseTemplateValues(templateValues)

		// IF TEMPLATE IS GIVEN
		// READ VALUES
		// READ VALUES (IF DEFINED)

		// FLAGS: REQUEST
		// FLAGS: REQUEST-CONFIG

		// IF REQUEST IS GIVEN, READ REQUEST - SKIP IF TEMPLATE IS GIVEN
		// IF REQUEST IS GIVEN, READ REQUEST-CONFIG e.g vmRequestConfig.yaml

		// READ GIVEN FIELDS IN REQUEST-CONFIG
		// IF MANDORY FILEDS NOT GIVEN GET RANDOM VALUES
		// GET DICT VALUES FROM ALL FIELDS

		// LOAD THE QUESTIONS FROM A KCL FILE
		questions, err := modules.ReadKCLQuestions(templatePath)
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
		renderedYaml := internal.RenderKCL(templatePath, allAnswers)

		// INITIALIZE AND RUN THE TERMINAL EDITOR PROGRAM.
		renderedYaml = modules.RunEditor(renderedYaml)

		// SAVE DIALOG
		modules.SaveDialog(renderedYaml)

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringSlice("values", []string{}, "templating values")
	renderCmd.Flags().String("template", "", "path to to be rendered template")
	renderCmd.Flags().String("destination", "", "path to output (if output file)")
}

// FLAGS:
// TEMPLATE
// REQUEST
// REQUEST-CONFIG
// OUTPUT
// VALUES

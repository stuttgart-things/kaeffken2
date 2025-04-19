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
	allAnswers         = make(map[string]interface{})
	templatePath       = "tests/proxmoxvm-template.k"
	renderedYAML       string
	values             map[string]interface{}
	templateFileExists bool
	configFileExists   bool
	requestFileExists  bool
	inputFiles         []inputFile
	dicts              = make(map[string]interface{})
	err                error
	listAnswers        map[string]interface{}
)

type inputFile struct {
	Name       string
	Path       string
	FileExists bool
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render templates",
	Long:  `Render templates based on profiles.`,
	Run: func(cmd *cobra.Command, args []string) {

		runSurvey, _ := cmd.LocalFlags().GetBool("survey")

		// INIT LOGGER
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})

		// GET/PARSE VALUES
		valueFLags, _ := cmd.Flags().GetStringSlice("values")
		destinationPath, _ := cmd.Flags().GetString("destination")

		if len(valueFLags) > 0 {
			values = internal.ParseTemplateValues(valueFLags)
			log.Info().Fields(values).Msg("values")
		} else {
			log.Info().Msg("NO VALUES GIVEN")
		}

		// VERIFY INPUT FILES
		requestFile, _ := cmd.Flags().GetString("request")
		configFile, _ := cmd.Flags().GetString("config")
		templatePath, _ := cmd.LocalFlags().GetString("template")

		inputFiles = append(inputFiles, inputFile{Name: "template", Path: templatePath, FileExists: false})
		inputFiles = append(inputFiles, inputFile{Name: "config", Path: configFile, FileExists: false})
		inputFiles = append(inputFiles, inputFile{Name: "request", Path: requestFile, FileExists: false})

		for i := range inputFiles {
			exists, err := internal.FileExists(inputFiles[i].Path)
			if err != nil {
				log.Error().Err(err).Str("path", inputFiles[i].Path).Msg("Error checking file")
			}
			inputFiles[i].FileExists = exists
		}

		for _, f := range inputFiles {
			switch f.Name + fmt.Sprintf(":%t", f.FileExists) {
			case "template:true":
				log.Info().Str("path", f.Path).Msg("Template exists ✅")

			case "template:false":
				log.Warn().Str("path", f.Path).Msg("Template missing ❌")

			case "request:true":
				log.Info().Str("path", f.Path).Msg("Request exists ✅")
				requestSpec, _ := internal.ReadSpecSection(f.Path)
				fmt.Println("SPEC:", requestSpec)

			case "request:false":
				log.Warn().Str("path", f.Path).Msg("Request missing ❌")

			case "config:true":
				log.Info().Str("path", f.Path).Msg("Config exists ✅")
				configSpec, _ := internal.ReadSpecSection(f.Path)
				fmt.Println("SPEC CONFIG:", configSpec)

				dicts, err = internal.ReadDicts(f.Path, "dicts")
				internal.CheckErr(err, "ERROR READING CONFIG DICTS")

			case "config:false":
				log.Warn().Str("path", f.Path).Msg("Config missing ❌")

			default:
				log.Warn().Str("name", f.Name).Str("path", f.Path).Msg("Unknown input file type or state ❌")
			}
		}

		// IF TEMPLATE IS DEFINED AND NO OTHER CONFIG
		// GET DEFAULT ANSWERS FROM SURVEY
		// MERGE WITH VALUES (VALUES ARE MOST IMPORTNAT)
		// TEST FOR ATTENDED AND UNATTENDED MODE

		fmt.Println("DICTS", dicts)

		// HOW TO GET DICT VALUES
		bla := internal.GetValueFromDicts(dicts, "kinds", "labul_proxmoxvm")
		fmt.Println("BLA", bla)

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
		internal.CheckErr(err, "ERROR READING KCL QUESTIONS")

		if !runSurvey {
			allAnswers = survey.GetRandomAnswers(questions)
			fmt.Println("RANDOM ANSWERS", allAnswers)
		}

		// BUILD THE SURVEY FORM AND GET A MAP FOR ANSWERS
		surveyForm, _, err := survey.BuildSurvey(questions)
		internal.CheckErr(err, "ERROR BUILDING SURVEY")

		// RUN THE INTERACTIVE SURVEY
		if runSurvey {
			err = surveyForm.Run()
			internal.CheckErr(err, "ERROR RUNNING SURVEY")
			// SET ANWERS TO ALL VALUES
			allAnswers = modules.SetAnswers(questions)
		}

		// LIST VALUES
		listDefaults := modules.ReadKCLList(templatePath)
		fmt.Println("LIST DEFAULTS", listDefaults)

		if runSurvey {
			listAnswers = modules.RunListEditor(listDefaults)
		} else {
			listAnswers = listDefaults
		}

		// MERGE ALL ANSWERS WITH VALUES
		allAnswers = internal.MergeMaps(allAnswers, internal.CleanMap(listAnswers))

		//reg := []string{"whateve", "vvdfvfdf", "patrick", "klaus", "test"}

		fmt.Println("ALL ANSWERS", allAnswers)

		// RENDER KCL FILE TO YAML
		renderedYaml := internal.RenderKCL(templatePath, allAnswers)

		fmt.Println("RENDERED YAML", renderedYaml)

		if runSurvey {

			// INITIALIZE AND RUN THE TERMINAL EDITOR PROGRAM.
			renderedYaml = modules.RunEditor(internal.CleanString(renderedYaml))

			// SAVE DIALOG
			modules.SaveDialog(renderedYaml)
		} else {
			internal.SaveToFile(renderedYaml, destinationPath)
			log.Info().Str("path", destinationPath).Msg("Outputfile written ✅")
		}
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringSlice("values", []string{}, "templating values")
	renderCmd.Flags().String("template", "", "path to to be rendered template")
	renderCmd.Flags().String("config", "", "path to config file")
	renderCmd.Flags().String("request", "", "path to request file")
	renderCmd.Flags().String("destination", "", "path to output (if output=file)")
	renderCmd.Flags().Bool("survey", true, "run survey")

}

// FLAGS:
// TEMPLATE
// REQUEST
// REQUEST-CONFIG
// OUTPUT
// VALUES

// go run main.go render --template tests/ansiblerun.k --values name=bla --request tests/vmRequest.yaml --config /home/sthings/projects/golang/kaeffken2/tests/vmRequestConfig.yaml

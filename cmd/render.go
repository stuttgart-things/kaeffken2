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
	inputFiles         []modules.InputFile
	configValues       = make(map[string]interface{})
	configSpec         = make(map[string]interface{})
	requestSpec        = make(map[string]interface{})
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

		// GET FLAGS
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

		// CHECK INPUT FILES
		inputFiles = append(inputFiles, modules.InputFile{Name: "template", Path: templatePath, FileExists: false})
		inputFiles = append(inputFiles, modules.InputFile{Name: "config", Path: configFile, FileExists: false})
		inputFiles = append(inputFiles, modules.InputFile{Name: "request", Path: requestFile, FileExists: false})
		configSpec, configValues, requestSpec = modules.CheckInputFiles(inputFiles)

		// IF TEMPLATE IS DEFINED AND NO OTHER CONFIG
		// GET DEFAULT ANSWERS FROM SURVEY
		// MERGE WITH VALUES (VALUES ARE MOST IMPORTANT)
		// TEST FOR ATTENDED AND UNATTENDED MODE

		// READ REQUEST SPEC
		if len(requestSpec) > 0 {
			fmt.Println("REQUEST", requestSpec)
		} else {
			log.Warn().Msg("NO REQUEST GIVEN")
		}

		// READ CONFIG (IF DEFINED)
		if len(configSpec) > 0 {

			// LOOP OVER ALL CONFIG KEYS
			for key := range configSpec {

				log.Info().Str("key", key).Msg("KEY SELECTED ✅")

				// GET RANDOM FROM KEY
				randomConfigKey, err := internal.GetRandomStringFromMap(configSpec, key)
				internal.CheckErr(err, "ERROR GETTING RANDOM VALUE FOR CONFIG")
				log.Info().Str("random", randomConfigKey).Msg("RANDOM CONFIG KEY SELECTED ✅")

				// GET VALUES AND SET TO ALL CONFIG VALUES
				allConfigValues := internal.GetValueFromDicts(configValues, key+"s", randomConfigKey)
				log.Info().Interface("config", allConfigValues).Msg("LOADED CONFIG VALUES")

				// MERGE ALL CONFIG VALUES WITH VALUES
				allAnswers = internal.MergeMaps(allAnswers, allConfigValues)
				fmt.Println("ALL CONFIG VALUES", allAnswers)
			}

		} else {
			log.Warn().Msg("NO CONFIG GIVEN")
		}

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
			// MERGE ALL RANDOM ANSWERS WITH FLAG VALUES
			allAnswers = internal.MergeMaps(allAnswers, values)
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

		// MERGE ALL ANSWERS WITH LIST ANSWERS
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
			log.Info().Str("path", destinationPath).Msg("OUTPUTFILE WRITTEN ✅")
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

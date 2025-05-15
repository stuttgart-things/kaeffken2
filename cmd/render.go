/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package cmd

import (
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
		// destinationPath, _ := cmd.Flags().GetString("destination")

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

		// READ CONFIG (IF DEFINED)
		if len(configSpec) > 0 {

			// LOOP OVER ALL CONFIG KEYS
			for key := range configSpec {

				var randomConfigKey string
				log.Info().Str("key", key).Msg("KEY SELECTED ✅")

				// CHECK IF KEY IS SET IN REQUEST
				if value, ok := requestSpec[key]; ok {
					log.Info().Str("key", key).Msg("CONFIG KEY IS SET IN REQUEST ✅")
					randomConfigKey = value.(string)
				} else {
					log.Info().Str("key", key).Msg("CONFIG KEY NOT SET IN REQUEST ✅")
					randomConfigKey, err = internal.GetRandomStringFromMap(configSpec, key)
					internal.CheckErr(err, "ERROR GETTING RANDOM VALUE FOR CONFIG")
					log.Info().Str("random", randomConfigKey).Msg("RANDOM CONFIG KEY SELECTED ✅")
				}

				// GET VALUES AND SET TO ALL CONFIG VALUES
				allConfigValues := internal.GetValueFromDicts(configValues, key+"s", randomConfigKey)
				log.Info().Interface("config", allConfigValues).Msg("LOADED CONFIG VALUES")

				// MERGE ALL CONFIG VALUES WITH VALUES
				allAnswers = internal.MergeMaps(allAnswers, allConfigValues)
				log.Info().Interface("config", allAnswers).Msg("MERGED CONFIG VALUES")
			}

		} else {
			log.Error().Msg("NO CONFIG FILE GIVEN - EXITING..")
			// EXIT PROGRAM
			os.Exit(1)
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
		// Retrieve template paths from allAnswers
		templates, err := internal.GetTemplatesPaths(allAnswers, values, "templates")
		if err != nil {
			log.Fatal().Err(err).Msg("FAILED TO GET TEMPLATE PATHS")
		}

		// Collect questions from all templates
		var allQuestions []*survey.Question
		for _, template := range templates {
			questions, err := modules.ReadKCLQuestions(template.Source)
			if err != nil {
				log.Fatal().Err(err).Str("template", template.Source).Msg("ERROR READING KCL QUESTIONS")
			}
			allQuestions = append(allQuestions, questions...)
		}

		// Get answers either randomly or via survey
		if !runSurvey {
			allAnswers = survey.GetRandomAnswers(allQuestions)
			allAnswers = internal.MergeMaps(allAnswers, values)
		}

		// Build and run survey if enabled
		surveyForm, _, err := survey.BuildSurvey(allQuestions)
		internal.CheckErr(err, "ERROR BUILDING SURVEY")

		if runSurvey {
			err = surveyForm.Run()
			internal.CheckErr(err, "ERROR RUNNING SURVEY")
			allAnswers = modules.SetAnswers(allQuestions)
		}

		// Process list inputs from all templates
		listDefaults := make(map[string]interface{})
		for _, template := range templates {
			defaults := modules.ReadKCLList(template.Source)
			listDefaults = internal.MergeMaps(listDefaults, defaults)
		}

		if runSurvey {
			listAnswers = modules.RunListEditor(listDefaults)
		} else {
			listAnswers = listDefaults
		}

		// Merge all answers and render each template
		allAnswers = internal.MergeMaps(allAnswers, internal.CleanMap(listAnswers))
		log.Info().Fields(allAnswers).Msg("COMBINED ANSWERS")

		// PROCESS EACH TEMPLATE
		for _, template := range templates {

			// RENDER TEMPLATE
			renderedYaml := internal.RenderKCL(template.Source, allAnswers)

			// HANDLE USER INTERACTION
			if runSurvey {
				renderedYaml = modules.RunEditor(internal.CleanString(renderedYaml))
				modules.SaveDialog(renderedYaml)
			} else {
				internal.SaveToFile(renderedYaml, template.Destination)
				log.Info().Str("path", template.Destination).Msg("OUTPUT FILE WRITTEN ✅")
			}
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

// go run main.go render --template tests/ansiblerun.k --values name=bla --request tests/vmRequest.yaml --config /home/sthings/projects/golang/kaeffken2/tests/vmRequestConfig.yaml

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/stuttgart-things/kaeffken2/internal"

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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// LOAD THE QUESTIONS FROM A KCL FILE
		questions, err := internal.ExtractQuestionsFromKCLFile(kclTemplate)
		if err != nil {
			// HANDLE ERROR
			log.Fatalf("Error extracting questions from KCL file: %v", err)
		}

		if len(questions) == 0 {
			fmt.Println("No questions found.")
			return
		}

		// BUILD THE SURVEY FORM AND GET A MAP FOR ANSWERS
		surveyForm, _, err := survey.BuildSurvey(questions)
		if err != nil {
			log.Fatalf("Error building survey: %v", err)
		}

		// RUN THE INTERACTIVE SURVEY
		err = surveyForm.Run()
		if err != nil {
			log.Fatalf("Error running survey: %v", err)
		}

		// SET ANWERS TO ALL VALUES
		for _, question := range questions {
			allAnswers[question.Name] = question.Default
		}

		// OUTPUT ALL ANSWERS + MODIFY
		for key, value := range allAnswers {
			fmt.Printf("%s: %v\n", key, value)
		}

		renderedYaml := internal.RenderKCL(kclTemplate, allAnswers)
		fmt.Println(renderedYaml)

		// INITIALIZE AND RUN THE TERMINAL EDITOR PROGRAM.
		p := tea.NewProgram(survey.InitialModel(renderedYaml), tea.WithAltScreen())
		m, err := p.Run()
		if err != nil {
			fmt.Println("Error running editor:", err)
			os.Exit(1)
		}

		// PRINT THE FINAL YAML CONTENT
		if result, ok := m.(survey.Text); ok && result.ErrMsg == "" {
			fmt.Println("\n" + lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1, 2).
				Render("Final YAML") + "\n")
			fmt.Println(result.Textarea.Value())
			renderedYAML = result.Textarea.Value()
		}

		// SAVE DIALOG
		p = tea.NewProgram(survey.InitialSaveModel(renderedYaml))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

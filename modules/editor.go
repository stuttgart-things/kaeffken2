/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package modules

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/charmbracelet/lipgloss"
	"github.com/stuttgart-things/survey"
)

func RunEditor(renderedYaml string) string {

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
		renderedYaml = result.Textarea.Value()
	}

	return renderedYaml

}

func RunListEditor(listDefaults map[string]interface{}) map[string]interface{} {

	fmt.Println("WHATEVER", listDefaults)

	// Initialize the listDefaults map if it's nil
	output := make(map[string]interface{})

	for key, value := range listDefaults {

		val, ok := value.([]string)
		if ok {
			fmt.Println("Converted:", val)
		} else {
			fmt.Println("Not a []string")
		}

		fmt.Println("KEY", key)
		fmt.Println("VALUE", val)

		p := tea.NewProgram(survey.InitListModel(val))
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error running prompt for %s: %v\n", key, err)
			continue
		}

		if model, ok := m.(survey.ListModel); ok && model.FinalOutput != "" {
			output[key] = model.FinalOutput
		}
	}

	return output
}

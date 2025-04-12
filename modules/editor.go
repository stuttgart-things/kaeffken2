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

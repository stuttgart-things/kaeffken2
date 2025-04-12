package modules

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/stuttgart-things/survey"
)

func SaveDialog(renderedYaml string) {

	// SAVE DIALOG
	p := tea.NewProgram(survey.InitialSaveModel(renderedYaml))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

}

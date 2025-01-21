package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	command "konbini/cli/commands"

	"konbini/cli/models"
)

func main() {
	if len(os.Args) == 1 {
		m := models.NewApp()
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err := p.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := command.Execute()
		if err != nil {
			log.Fatal(err)
		}
	}
}

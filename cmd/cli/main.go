package main

import (
	"log"
	"os"

	command "konbini/cli/commands"
	"konbini/cli/config"

	tea "github.com/charmbracelet/bubbletea"

	"konbini/cli/tui"
)

func main() {
	config.Init()
	if len(os.Args) == 1 {
		m := tui.New()
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

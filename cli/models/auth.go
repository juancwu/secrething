package models

import (
	"konbini/cli/router"

	tea "github.com/charmbracelet/bubbletea"
)

type authModel struct {
}

func newLoginModel() authModel {
	return authModel{}
}

func (m authModel) Init() tea.Cmd {
	return nil
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, router.NewNavigationMsg(menuPageID, nil)
		}
	}
	return m, nil
}

func (m authModel) View() string {
	return "auth"
}

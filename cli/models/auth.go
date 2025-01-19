package models

import (
	"konbini/cli/router"

	tea "github.com/charmbracelet/bubbletea"
)

type loginModel struct {
}

func newLoginModel() loginModel {
	return loginModel{}
}

func (m loginModel) Init() tea.Cmd {
	return nil
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, router.NewNavigationMsg(menuPageID, nil)
		}
	}
	return m, nil
}

func (m loginModel) View() string {
	return "auth"
}

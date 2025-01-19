package models

import tea "github.com/charmbracelet/bubbletea/v2"

type authModel struct {
}

func newLoginModel() authModel {
	return authModel{}
}

func (m authModel) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m authModel) View() string {
	return "auth"
}

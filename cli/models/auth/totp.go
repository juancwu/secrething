package auth

import tea "github.com/charmbracelet/bubbletea"

type totpModel struct{}

func NewTOTPModel() totpModel {
	return totpModel{}
}

func (m totpModel) Init() tea.Cmd {
	return nil
}

func (m totpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m totpModel) View() string {
	return "totp"
}

package tui

import tea "github.com/charmbracelet/bubbletea"

type setupTOTPModel struct {
}

func newSetupTOTP(params map[string]interface{}) setupTOTPModel {
	return setupTOTPModel{}
}

func (m setupTOTPModel) Init() tea.Cmd {
	return nil
}

func (m setupTOTPModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m setupTOTPModel) View() string {
	return "totp"
}

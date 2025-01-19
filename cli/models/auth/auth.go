package auth

import (
	tea "github.com/charmbracelet/bubbletea"
)

type loginModel struct {
	showModal bool
}

func NewLogin() loginModel {
	return loginModel{showModal: false}
}

func (m loginModel) Init() tea.Cmd {
	return nil
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			m.showModal = !m.showModal
		}
	}
	return m, nil
}

func (m loginModel) View() string {
	return "login"
}

type registerModel struct {
}

func NewRegister() registerModel {
	return registerModel{}
}

func (m registerModel) Init() tea.Cmd {
	return nil
}

func (m registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m registerModel) View() string {
	return "register"
}

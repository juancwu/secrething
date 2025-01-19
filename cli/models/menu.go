package models

import (
	"konbini/cli/router"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// menuModel represents the main menu of the konbini cli which has options to
// register/login, manage bentos, manage groups, and permissions
type menuModel struct {
	list list.Model
}

// newMenuModel creates a new menu model
func newMenuModel() menuModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	return menuModel{
		list: l,
	}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			return m, router.NewNavigationMsg(loginPageID, nil)
		case "r":
			return m, router.NewNavigationMsg(registerPageID, nil)
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	return "hello"
}

package tui

import (
	"github.com/juancwu/konbini/cli/config"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// menuItem represents a single menu option
type menuItem struct {
	title       string
	description string
	route       string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

// menuModel represents the main menu of the konbini cli
type menuModel struct {
	list   list.Model
	width  int
	height int
}

type Menu struct {
	list list.Model

	width  int
	height int
}

// menu creates a new Menu model
func newMenu(params map[string]interface{}) Menu {
	var l list.Model

	width := params[paramAppWidth].(int)
	height := params[paramAppHeight].(int)

	auth := config.GetAuth()
	if auth == nil {
		l = list.New([]list.Item{
			menuItem{
				title:       "Register",
				description: "Create a new account",
				route:       "register",
			},
			menuItem{
				title:       "Login",
				description: "Sign in to your account",
				route:       "login",
			},
		}, list.NewDefaultDelegate(), width, height)
	} else if auth.Token != "" {
		l = list.New([]list.Item{
			menuItem{
				title:       "Manage Bentos",
				description: "View and manage your bento configurations",
				route:       "bentos", // Add this route when ready
			},
			menuItem{
				title:       "Manage Groups",
				description: "Configure and manage user groups",
				route:       "groups", // Add this route when ready
			},
			menuItem{
				title:       "Permissions",
				description: "Manage access controls and permissions",
				route:       "permissions", // Add this route when ready
			},
		}, list.NewDefaultDelegate(), width, height)
	}

	return Menu{
		list:   l,
		width:  width,
		height: height,
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(m.width, m.height)
		return m, nil
	case tea.KeyMsg:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Menu) View() string {
	return m.list.View()
}

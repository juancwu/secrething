package menu

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"konbini/cli/router"
)

var (
	itemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2)
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

// New creates a new menu model with pre-populated items
func New(width int, height int) menuModel {
	// Define menu items
	items := []list.Item{
		menuItem{
			title:       "Login",
			description: "Sign in to your account",
			route:       "/login",
		},
		menuItem{
			title:       "Register",
			description: "Create a new account",
			route:       "/register",
		},
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
	}

	// Initialize list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62"))

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62"))

	l := list.New(items, delegate, 50, 4*len(items))
	l.Title = "Konbini CLI"
	l.SetShowHelp(true)

	return menuModel{
		list:   l,
		width:  width,
		height: height,
	}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if i, ok := m.list.SelectedItem().(menuItem); ok {
				return m, router.NewNavigationMsg(i.route, nil)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m menuModel) View() string {
	return fmt.Sprintf("\n%s", m.list.View())
}

// Implement LifecycleHooks interface
func (m menuModel) OnEnter(params map[string]interface{}) tea.Cmd {
	return nil
}

func (m menuModel) OnExit() tea.Cmd {
	return nil
}

func (m menuModel) BeforeNavigateAway() bool {
	return true
}

func (m menuModel) AfterNavigateAway() tea.Cmd {
	return nil
}

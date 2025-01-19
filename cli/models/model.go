package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

// model is the main model that is used to render the TUI
type model struct {
	availableModels   map[page]tea.Model
	activeModel       page
	initializedModels map[page]bool
}

func NewModel() model {
	return model{
		activeModel: menuPageID,
		availableModels: map[page]tea.Model{
			menuPageID: newMenuModel(),
		},
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	// initialize all the available models
	for _, value := range m.availableModels {
		cmd := value.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	activeModel, cmd := m.availableModels[m.activeModel].Update(msg)
	// update the active model store
	m.availableModels[m.activeModel] = activeModel
	return m, cmd
}

// View only renders the activeModel view
func (m model) View() string {
	return m.availableModels[m.activeModel].View()
}

package models

import (
	"konbini/cli/router"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea/v2"
)

// app is the main app that is used to render the TUI
type app struct {
	// pointer to a router because the router maps should not be copied
	// but just in-place updated.
	router *router.Router

	keys keyMap
	help help.Model

	err error

	debugProfile debugProfile
	debugMode    bool
}

type debugProfile struct {
	Input string
}

func NewApp() app {
	r := router.NewRouter()

	r.RegisterPage(
		menuPageID,
		func(params map[string]interface{}) tea.Model {
			return newMenuModel()
		},
	)
	r.RegisterPage(
		loginPageID,
		func(params map[string]interface{}) tea.Model {
			return newLoginModel()
		},
	)

	return app{
		router:       r,
		keys:         defaultKeyMap(),
		help:         help.New(),
		debugProfile: debugProfile{},
		debugMode:    true,
	}
}

func (a app) Init() (tea.Model, tea.Cmd) {
	cmd, err := a.router.SetInitialPage(menuPageID, nil)
	if err != nil {
		panic(err) // TODO: do something better
	}
	return a, cmd
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		a.debugProfile.Input = msg.String()
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "ctrl+b": // navigate back
			model, cmd := a.router.Back()
			a.router.UpdateCurrentModel(model)
			return a, cmd
		case "ctrl+f": // navigate forward
			model, cmd := a.router.Forward()
			a.router.UpdateCurrentModel(model)
			return a, cmd
		}
	case router.NavigationMsg:
		cmd, err := a.router.Navigate(msg.To, msg.Params)
		if err != nil {
			return a, newErrMsg(err)
		}
		return a, cmd
	}

	// Update current page
	currentModel, cmd := a.router.CurrentModel().Update(msg)
	a.router.UpdateCurrentModel(currentModel)
	return a, cmd
}

// View only renders the activeModel view
func (a app) View() string {
	var builder strings.Builder

	if a.debugMode {
		builder.WriteString(a.debugPrint())
	}

	builder.WriteString(a.router.HistoryString() + "\n")

	modelView := a.router.CurrentModel().View()
	builder.WriteString(modelView + "\n")

	builder.WriteString(a.help.View(a.keys))

	return builder.String()
}

func (a *app) debug(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		a.debugProfile.Input = msg.String()
	}
}

func (a *app) debugPrint() string {
	return "input: " + a.debugProfile.Input + "\n"
}

package models

import (
	"konbini/cli/models/menu"
	"konbini/cli/router"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// app is the main app that is used to render the TUI
type app struct {
	// pointer to a router because the router maps should not be copied
	// but just in-place updated.
	router *router.Router

	keys keyMap
	help help.Model

	err error

	width  int
	height int

	ready bool

	debugProfile debugProfile
	debugMode    bool
	debugOverlay debugOverlay
	showDebug    bool
}

type debugProfile struct {
	Input string
}

func NewApp() app {
	r := router.NewRouter()

	r.RegisterPage(
		menuPageID,
		func(params map[string]interface{}) tea.Model {
			return menu.New()
		},
	)
	r.RegisterPage(
		loginPageID,
		func(params map[string]interface{}) tea.Model {
			return newLoginModel()
		},
	)
	r.RegisterPage(
		registerPageID,
		func(params map[string]interface{}) tea.Model {
			return newRegisterModel()
		},
	)

	return app{
		router:       r,
		keys:         defaultKeyMap(),
		help:         help.New(),
		debugProfile: debugProfile{},
		debugMode:    true,
		debugOverlay: newDebugOverlay(0, 0),
		showDebug:    false,
	}
}

func (a app) Init() tea.Cmd {
	cmd, err := a.router.SetInitialPage(menuPageID, nil)
	if err != nil {
		panic(err)
	}
	return cmd
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
		case "alt+ctrl+d":
			a.showDebug = !a.showDebug
			return a, nil
		}
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		a.debugOverlay = newDebugOverlay(a.width, a.height)
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
	if !a.ready {
		return "not ready"
	}

	if a.showDebug && a.debugMode {
		return a.debugOverlay.View(a.debugProfile, a.router.HistoryString())
	}

	view := a.router.CurrentModel().View()
	return view
}

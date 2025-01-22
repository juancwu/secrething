package models

import (
	"konbini/cli/models/auth"
	"konbini/cli/models/menu"
	"konbini/cli/router"
	"konbini/cli/secrets"

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

	windowSizeDone bool
	authCheckDone  bool

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
			var (
				width  int
				height int
			)

			if i, ok := params["app_width"]; ok {
				width = i.(int)
			}
			if i, ok := params["app_height"]; ok {
				height = i.(int)
			}

			return menu.New(width, height)
		},
	)
	r.RegisterPage(
		loginPageID,
		func(params map[string]interface{}) tea.Model {
			return auth.NewLogin()
		},
	)
	r.RegisterPage(
		registerPageID,
		func(params map[string]interface{}) tea.Model {
			return auth.NewRegister()
		},
	)
	r.RegisterPage(
		setupTOTPPageID,
		func(params map[string]interface{}) tea.Model {
			return auth.NewTOTPModel()
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
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var err error
	cmd, err = a.router.SetInitialPage(menuPageID, nil)
	if err != nil {
		panic(err)
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	cmds = append(cmds, a.persistAuth())
	return tea.Batch(cmds...)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		a.debugProfile.Input = msg.String()
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "alt+ctrl+d":
			a.showDebug = !a.showDebug
			return a, nil
		}
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.windowSizeDone = true
		a.debugOverlay = newDebugOverlay(a.width, a.height)
	case router.NavigationMsg:
		if msg.Params == nil {
			msg.Params = make(map[string]interface{})
		}
		msg.Params["app_width"] = a.width
		msg.Params["app_height"] = a.height
		cmd, err := a.router.Navigate(msg.To, msg.Params)
		if err != nil {
			return a, newErrMsg(err)
		}
		return a, cmd
	case initMsg:
		a.authCheckDone = true
	}

	// Update current page
	currentModel, cmd := a.router.CurrentModel().Update(msg)
	a.router.UpdateCurrentModel(currentModel)
	return a, cmd
}

// View only renders the activeModel view
func (a app) View() string {
	if a.showDebug && a.debugMode {
		return a.debugOverlay.View(a.debugProfile, a.router.HistoryString())
	}

	view := a.router.CurrentModel().View()
	return view
}

type initMsg struct{}

// persistAuth checks if the current
func (a app) persistAuth() tea.Cmd {
	return func() tea.Msg {
		secrets.CheckAuth()
		return initMsg{}
	}
}

func (a app) ready() bool {
	return a.authCheckDone && a.windowSizeDone
}

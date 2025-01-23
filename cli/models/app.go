package models

import (
	"konbini/cli/config"
	"konbini/cli/models/auth"
	"konbini/cli/models/menu"
	"konbini/cli/router"
	"konbini/cli/secrets"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// app is the main app that is used to render the TUI
type app struct {
	// pointer to a router because the router maps should not be copied
	// but just in-place updated.
	router *router.Router

	spinner spinner.Model

	keys keyMap
	help help.Model

	err error

	width  int
	height int

	// Steps that needs to be done before hand
	windowSizeDone bool
	authCheckDone  bool

	// If it has been all initialized or not
	initialized bool

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

	s := spinner.New()
	s.Spinner = spinner.Dot

	return app{
		router:         r,
		spinner:        s,
		keys:           defaultKeyMap(),
		help:           help.New(),
		debugProfile:   debugProfile{},
		debugMode:      true,
		debugOverlay:   newDebugOverlay(0, 0),
		showDebug:      false,
		authCheckDone:  false,
		windowSizeDone: false,
		initialized:    false,
	}
}

func (a app) Init() tea.Cmd {
	return a.spinner.Tick
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
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
		config.UpdateTermSize(a.width, a.height)
	case router.NavigationMsg:
		cmd, err := a.router.Navigate(msg.To, msg.Params)
		if err != nil {
			return a, newErrMsg(err)
		}
		return a, cmd
	case authCheckMsg:
		a.authCheckDone = true
	}

	if !a.ready() {
		a.spinner, cmd = a.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else if !a.initialized {
		// ready to set initial page
		cmd, err := a.router.SetInitialPage(menuPageID, nil)
		if err != nil {
			panic(err)
		}
		a.initialized = true
		return a, cmd
	} else {
		// Update current page
		currentModel, cmd := a.router.CurrentModel().Update(msg)
		a.router.UpdateCurrentModel(currentModel)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

// View only renders the activeModel view
func (a app) View() string {
	if !a.ready() {
		return a.spinner.View() + " Loading..."
	}

	if a.showDebug && a.debugMode {
		return a.debugOverlay.View(a.debugProfile, a.router.HistoryString())
	}

	view := a.router.CurrentModel().View()
	return view
}

type authCheckMsg struct {
	redirectTo string
}

// persistAuth checks if the current
func (a app) persistAuth() tea.Msg {
	err := secrets.CheckAuth()
	msg := authCheckMsg{redirectTo: menuPageID}
	// check for partial token
	if err == nil && !secrets.TOTPSet() {
		// redirect to totp setup
		msg.redirectTo = setupTOTPPageID
	}
	return msg
}

func (a app) ready() bool {
	return a.authCheckDone && a.windowSizeDone
}

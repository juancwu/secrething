package tui

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"konbini/cli/config"
	"konbini/cli/router"
	"konbini/common/api"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	keyring "github.com/zalando/go-keyring"
)

type GlobalKeyMap struct {
	ForceQuit key.Binding
}

func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ForceQuit,
	}
}

func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ForceQuit},
	}
}

type App struct {
	s spinner.Model
	h help.Model
	r *router.Router

	keys GlobalKeyMap

	width  int
	height int

	// checkpoints
	windowSizeCheck bool
	authCheckDone   bool

	routerInitialized bool
}

// New creates a new app model
func New() App {
	r := router.NewRouter()

	r.RegisterPage("menu", func(params map[string]interface{}) tea.Model {
		return newMenu(params)
	})

	keys := GlobalKeyMap{
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Force quit"),
		),
	}

	return App{
		s: newSpinner(),
		r: r,
		h: help.New(),

		keys: keys,

		width:  0,
		height: 0,

		windowSizeCheck: false,
		authCheckDone:   false,

		routerInitialized: false,
	}
}

func (m App) Init() tea.Cmd {
	return tea.Batch(
		m.s.Tick,
		m.authCheck,
	)
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
		err  error
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ForceQuit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.windowSizeCheck {
			m.width = msg.Width
			m.height = msg.Height
			m.windowSizeCheck = true
		}
	case router.NavigationMsg:
		params := m.globalParams(msg.Params)
		cmd, err = m.r.Navigate(msg.To, params)
		if err != nil {
			panic(err)
		}
		cmds = append(cmds, cmd)
	case spinner.TickMsg:
		if !m.ready() {
			m.s, cmd = m.s.Update(msg)
			return m, cmd
		}
	case authCheckMsg:
		m.authCheckDone = true
	}

	if m.ready() {
		if !m.routerInitialized {
			m.r.SetInitialPage("menu", m.globalParams(nil))
			m.routerInitialized = true
		}

		model := m.r.CurrentModel()
		model, cmd = model.Update(msg)
		m.r.UpdateCurrentModel(model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m App) View() string {
	if !m.ready() {
		return m.s.View() + " Loading..."
	}

	return m.r.CurrentModel().View()
}

type authCheckMsg struct {
	Err error
}

func (m App) authCheck() tea.Msg {
	token, err := keyring.Get("konbini", "user")
	if err != nil {
		return authCheckMsg{Err: err}
	}

	reqBody := api.CheckAuthTokenRequest{AuthToken: token}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return authCheckMsg{Err: err}
	}
	reader := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(
		http.MethodPost,
		config.BackendUrl(api.UriCheckToken),
		reader,
	)
	if err != nil {
		return authCheckMsg{Err: err}
	}
	req.Header.Add(api.HeaderContentType, api.MimeApplicationJson)

	c := http.Client{Timeout: time.Second * 30}
	res, err := c.Do(req)
	if err != nil {
		return authCheckMsg{Err: err}
	}

	data, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var resBody api.ErrorResponse
		err := json.Unmarshal(data, &resBody)
		if err == nil {
			err = errors.New(resBody.Message)
		}
		return authCheckMsg{Err: err}
	}

	var resBody api.CheckAuthResponse
	err = json.Unmarshal(data, &resBody)
	if err != nil {
		return authCheckMsg{Err: err}
	}

	if resBody.AuthToken != "" {
		keyring.Set("konbini", "user", resBody.AuthToken)
	}

	config.SetAuth(config.Auth{
		Token:         resBody.AuthToken,
		TokenType:     resBody.TokenType,
		TOTP:          resBody.TOTP,
		EmailVerified: resBody.EmailVerified,
	})

	return authCheckMsg{}
}

// Global parameters
const (
	paramAppWidth  = "_app_width"
	paramAppHeight = "_app_height"
)

// globalParams extends the given params with global params, pass nil is no params
func (m App) globalParams(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		params = make(map[string]interface{})
	}
	params[paramAppWidth] = m.width
	params[paramAppHeight] = m.height
	return params
}

func (m App) ready() bool {
	return m.windowSizeCheck && m.authCheckDone
}

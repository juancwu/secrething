package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"konbini/cli/config"
	"konbini/cli/router"
	"konbini/cli/secrets"
	"konbini/common/api"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	wrapperStyles = lipgloss.NewStyle().Margin(1)

	redTextStyles = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	spinnerStyles = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type loginKeyMap struct {
	Tab   key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func (k loginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Quit}
}

func (k loginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.Enter},
		{k.Quit},
	}
}

type loginModel struct {
	spinner                    spinner.Model
	emailInput                 textinput.Model
	passwordInput              textinput.Model
	totpInput                  textinput.Model
	showTOTPInput              bool
	loading                    bool
	keys                       loginKeyMap
	help                       help.Model
	needConfirmationBeforeExit bool
	err                        error
}

func NewLogin() loginModel {
	eti := textinput.New()
	eti.Placeholder = "Email"
	eti.Focus()
	eti.Validate = validateEmail

	pti := textinput.New()
	pti.Placeholder = "Password"
	pti.EchoMode = textinput.EchoPassword
	pti.Validate = func(s string) error {
		return validatePasswords(s)
	}

	totp := textinput.New()
	totp.Placeholder = "Enter 2FA Code (6 digits)"
	totp.Validate = func(s string) error {
		match, err := regexp.MatchString("^[0-9]{6}$", s)
		if err != nil {
			return err
		}

		if !match {
			return errors.New("Invalid code")
		}

		return nil
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyles

	keys := loginKeyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "Switch inputs"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Login"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Return to menu"),
		),
	}

	return loginModel{
		spinner:       s,
		emailInput:    eti,
		passwordInput: pti,
		totpInput:     totp,
		showTOTPInput: false,
		keys:          keys,
		help:          help.New(),
		loading:       false,
	}
}

func (m loginModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if !m.loading {
			// if set then just exit if the input is "y"
			if m.needConfirmationBeforeExit && msg.String() == "y" {
				return m, router.NewNavigationMsg("menu", nil)
			} else if m.needConfirmationBeforeExit {
				m.needConfirmationBeforeExit = false
				return m, nil
			}

			switch msg.Type {
			case tea.KeyEsc:
				m.needConfirmationBeforeExit = m.emailInput.Value() != "" || m.passwordInput.Value() != ""
				if m.needConfirmationBeforeExit {
					return m, nil
				}

				return m, router.NewNavigationMsg("menu", nil)
			case tea.KeyEnter:
				if m.showTOTPInput {
					if err := m.totpInput.Validate(m.totpInput.Value()); err != nil {
						// TODO: do something here
						return m, nil
					}
					m.loading = true
					return m, m.login
				}
				if m.emailInput.Focused() {
					m.emailInput.Blur()
					m.passwordInput.Focus()
				} else {
					// make request here
					m.loading = true
					return m, m.login
				}
			case tea.KeyTab:
				if m.emailInput.Focused() {
					m.emailInput.Blur()
					m.passwordInput.Focus()
				} else {
					m.passwordInput.Blur()
					m.emailInput.Focus()
				}
			}

			if m.showTOTPInput {
				m.totpInput, cmd = m.totpInput.Update(msg)
				return m, cmd
			}

			if m.emailInput.Focused() {
				m.emailInput, cmd = m.emailInput.Update(msg)
				return m, cmd
			}
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			return m, cmd
		}

	case totpMsg:
		m.loading = false
		m.showTOTPInput = true
		m.totpInput.Focus()
		m.emailInput.Blur()
		m.passwordInput.Blur()

	case api.LoginResponse:
		m.loading = false
		secrets.SaveCredentials(msg.Token)
		return m, router.NewNavigationMsg("menu", nil)

	case error:
		m.loading = false
		m.err = msg
	}

	return m, nil
}

func (m loginModel) View() string {
	var parts []string

	if m.showTOTPInput {
		parts = append(parts, m.totpInput.View())
	} else {
		parts = append(parts, m.emailInput.View(), m.passwordInput.View())

	}

	if m.err != nil {
		parts = append(parts, redTextStyles.Render(m.err.Error()))
	}

	if m.loading {
		parts = append(parts, m.spinner.View()+" Waiting...")
	}

	if m.needConfirmationBeforeExit {
		parts = append(parts, redTextStyles.Render("Are you sure you want to return to the menu? Press 'y' to confirm (any other key to cancel)"))
	}

	parts = append(parts, m.help.View(m.keys))

	return wrapperStyles.Render(strings.Join(parts, "\n\n"))
}

// totpMsg signals that totp code is needed to login
type totpMsg struct{}

// login makes a login request and passes the response in a loginResponseMsg
func (m loginModel) login() tea.Msg {
	reqBody := api.LoginRequest{
		Email:    m.emailInput.Value(),
		Password: m.passwordInput.Value(),
		TOTPCode: nil,
	}

	if m.totpInput.Value() != "" {
		code := m.totpInput.Value()
		reqBody.TOTPCode = &code
	}

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(reqBodyData)

	req, err := http.NewRequest(
		http.MethodPost,
		config.BackendUrl()+"/auth/login",
		reader,
	)
	req.Header.Add("Content-Type", "application/json")

	c := http.Client{Timeout: time.Second * 10}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		var body api.LoginResponse
		err = json.Unmarshal(data, &body)
		if err != nil {
			return err
		}
		return body
	default:
		var body api.ErrorResponse
		err = json.Unmarshal(data, &body)
		if err != nil {
			return err
		}
		if strings.Contains(body.Message, "code is required") {
			return totpMsg{}
		}
		return errors.New(body.Message)
	}
}

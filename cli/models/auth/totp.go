package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"konbini/cli/config"
	"konbini/cli/router"
	"konbini/cli/secrets"
	commonAPI "konbini/common/api"
	"konbini/common/qrcode"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Retry key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Retry, k.Enter},
		{k.Quit},
	}
}

type totpModel struct {
	spinner   spinner.Model
	codeInput textinput.Model
	qr        string
	otpURL    string
	keys      keyMap
	loading   bool
	err       error
	help      help.Model

	showCodes bool
	codes     []string
}

func NewTOTPModel() totpModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	i := textinput.New()
	i.Placeholder = "6 digits code"
	i.Focus()
	i.Validate = func(s string) error {
		match, err := regexp.MatchString("^[0-9]{6}$", s)
		if err != nil {
			return err
		}

		if !match {
			return errors.New("Invalid code")
		}

		return nil
	}

	keys := keyMap{
		Retry: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Retry")),
		Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Confirm")),
		Quit:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Quit")),
	}

	return totpModel{
		spinner:   s,
		codeInput: i,
		keys:      keys,
		qr:        "",
		loading:   true, // start at loading because the model makes a setup totp request on load
		help:      help.New(),
	}
}

func (m totpModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.requestSetup)
}

func (m totpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(m.codes) > 0 {
				return m, router.NewNavigationMsg("menu", nil)
			}
			if err := m.codeInput.Validate(m.codeInput.Value()); err != nil {
				m.err = err
			} else {
				// try to lock down the totp
				m.loading = true
				return m, m.lockTOTP(m.codeInput.Value())
			}
		case "esc":
			return m, tea.Quit
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "backspace":
			if msg.String() == "backspace" || len(m.codeInput.Value()) < 6 {
				m.codeInput, cmd = m.codeInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case requestMsg:
		m.loading = false
		m.qr = qrcode.New(msg.otpURL)
		m.otpURL = msg.otpURL
	case failRequestMsg:
		m.loading = false
		m.err = msg.Err
	case commonAPI.LockTOTPResponse:
		m.loading = false
		m.showCodes = true
		m.codes = msg.RecoveryCodes
		m.qr = ""
		m.otpURL = ""

		// save new token
		secrets.SaveCredentials(msg.Token)

	case error:
		m.loading = false
		m.err = msg
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m totpModel) View() string {
	parts := []string{"Setup TOTP/2FA (This step is mandatory)"}

	if m.err != nil {
		parts = append(parts, redTextStyles.Render(m.err.Error()))
	}

	if len(m.codes) > 0 {
		parts = append(parts, "Recovery Codes", strings.Join(m.codes, "\n"), "Make sure to not lose these codes and they won't be shown again.", "Press Enter to continue")
	} else {
		if m.loading {
			parts = append(parts, m.spinner.View()+" Loading...")
		} else {
			parts = append(parts, m.codeInput.View())
		}
		parts = append(parts, m.qr, m.otpURL)
	}

	parts = append(parts, m.help.View(m.keys))

	return strings.Join(parts, "\n\n")
}

type requestMsg struct {
	otpURL string
}

type failRequestMsg struct {
	Err error
}

func (m totpModel) requestSetup() tea.Msg {
	token := secrets.AuthToken()
	req, err := http.NewRequest(
		http.MethodPost,
		config.BackendUrl()+"/auth/totp/setup",
		nil,
	)
	if err != nil {
		return failRequestMsg{Err: err}
	}
	req.Header.Add("Authorization", "Bearer "+token)

	c := http.Client{Timeout: time.Second * 10}
	res, err := c.Do(req)
	if err != nil {
		return failRequestMsg{Err: err}
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return failRequestMsg{Err: err}
	}

	if res.StatusCode == http.StatusOK {
		var body commonAPI.SetupTOTPResponse
		err = json.Unmarshal(data, &body)
		if err != nil {
			return failRequestMsg{Err: err}
		}

		return requestMsg{otpURL: body.URL}
	}

	var body commonAPI.ErrorResponse
	err = json.Unmarshal(data, &body)
	if err != nil {
		return failRequestMsg{Err: err}
	}

	return failRequestMsg{
		Err: errors.New(body.Message),
	}
}

func (m totpModel) lockTOTP(code string) tea.Cmd {
	return func() tea.Msg {
		reqBodyData, err := json.Marshal(commonAPI.SetupTOTPLockRequest{Code: code})
		if err != nil {
			return failRequestMsg{Err: err}
		}

		reader := bytes.NewReader(reqBodyData)

		req, err := http.NewRequest(
			http.MethodPost,
			config.BackendUrl()+"/auth/totp/lock",
			reader,
		)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+secrets.AuthToken())

		c := http.Client{Timeout: time.Second * 10}
		res, err := c.Do(req)
		if err != nil {
			return err
		}

		data, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			var body commonAPI.ErrorResponse
			err = json.Unmarshal(data, &body)
			if err != nil {
				return err
			}
			return errors.New(body.Message)
		}

		var body commonAPI.LockTOTPResponse
		err = json.Unmarshal(data, &body)
		if err != nil {
			return err
		}

		return body
	}
}

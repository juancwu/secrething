package auth

import (
	"konbini/cli/router"
	"konbini/cli/secrets"
	"konbini/cli/services"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type registerKeyMap struct {
	Tab   key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func (k registerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Quit}
}

func (k registerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.Enter},
		{k.Quit},
	}
}

type registerModel struct {
	spinner                    spinner.Model
	inputs                     []textinput.Model
	currentInputIdx            int
	loading                    bool
	keys                       registerKeyMap
	help                       help.Model
	needConfirmationBeforeExit bool
	err                        error
}

func NewRegister() registerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyles

	ei := textinput.New()
	ei.Placeholder = "Enter email"
	ei.Validate = validateEmail
	ei.Focus()

	ni := textinput.New()
	ni.Placeholder = "Enter nickname"
	ni.Validate = validateNickname

	pi := textinput.New()
	pi.Placeholder = "Enter password"
	pi.Validate = func(s string) error {
		return validatePasswords(s)
	}

	cpi := textinput.New()
	cpi.Placeholder = "Confirm password"
	cpi.Validate = func(s string) error {
		return nil
	}

	keys := registerKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Confirm"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab", "up", "down"),
			key.WithHelp("tab/up/down", "Switch inputs"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Return to menu"),
		),
	}

	return registerModel{
		spinner:                    s,
		inputs:                     []textinput.Model{ei, ni, pi, cpi},
		currentInputIdx:            0,
		loading:                    false,
		keys:                       keys,
		help:                       help.New(),
		needConfirmationBeforeExit: false,
	}
}

func (m registerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.needConfirmationBeforeExit = m.isDirty()
				if m.needConfirmationBeforeExit {
					return m, nil
				}
				return m, router.NewNavigationMsg("menu", nil)
			case tea.KeyEnter:
				if m.currentInputIdx == len(m.inputs)-1 {
					for _, i := range m.inputs {
						err := i.Validate(i.Value())
						if err != nil {
							m.err = err
							return m, nil
						}
					}

					// match passwords
					if err := validatePasswords(m.inputs[2].Value(), m.inputs[3].Value()); err != nil {
						m.err = err
						return m, nil
					}

					m.loading = true
					return m, m.register
				}
				m.switchInput(1)
				return m, nil
			case tea.KeyTab:
				m.switchInput(1)
				return m, nil
			case tea.KeyUp:
				m.switchInput(-1)
				return m, nil
			case tea.KeyDown:
				m.switchInput(1)
				return m, nil
			}

			m.inputs[m.currentInputIdx], cmd = m.inputs[m.currentInputIdx].Update(msg)

			return m, cmd
		}

	case registerMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}

		// save the partial token
		secrets.SaveCredentials(msg.Token, m.inputs[0].Value())

		// move onto setting up 2FA
		return m, router.NewNavigationMsg("setup-totp", nil)
	}

	return m, nil
}

func (m registerModel) View() string {
	parts := []string{}

	if m.err != nil {
		parts = append(parts, redTextStyles.Render(m.err.Error()))
	}

	for _, i := range m.inputs {
		parts = append(parts, i.View())
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

type registerMsg struct {
	Token string
	Err   error
}

func (m registerModel) register() tea.Msg {
	email := m.inputs[0].Value()
	nickname := m.inputs[1].Value()
	password := m.inputs[2].Value()
	res, err := services.Register(email, nickname, password)
	return registerMsg{
		Token: res.Token,
		Err:   err,
	}
}

func (m registerModel) isDirty() bool {
	for _, i := range m.inputs {
		if i.Value() != "" {
			return true
		}
	}
	return false
}

func (m *registerModel) switchInput(inc int) {
	m.inputs[m.currentInputIdx].Blur()
	m.currentInputIdx = (m.currentInputIdx + inc) % len(m.inputs)
	if m.currentInputIdx < 0 {
		m.currentInputIdx = len(m.inputs) - 1
	}
	m.inputs[m.currentInputIdx].Focus()
}

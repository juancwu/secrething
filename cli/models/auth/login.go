package auth

import (
	"konbini/cli/router"
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
	loading                    bool
	keys                       loginKeyMap
	help                       help.Model
	needConfirmationBeforeExit bool
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
			key.WithHelp("enter", "Confirm input/submit"),
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

			if m.emailInput.Focused() {
				m.emailInput, cmd = m.emailInput.Update(msg)
				return m, cmd
			}
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			return m, cmd
		}

	case loginResponseMsg:
		m.loading = false
	}

	return m, nil
}

func (m loginModel) View() string {

	parts := []string{m.emailInput.View(), m.passwordInput.View()}

	if m.loading {
		parts = append(parts, m.spinner.View()+" Waiting...")
	}

	if m.needConfirmationBeforeExit {
		parts = append(parts, redTextStyles.Render("Are you sure you want to return to the menu? Press 'y' to confirm (any other key to cancel)"))
	}

	parts = append(parts, m.help.View(m.keys))

	return wrapperStyles.Render(strings.Join(parts, "\n\n"))
}

type loginResponseMsg struct {
	Status    int
	AuthToken string
}

// login makes a login request and passes the response in a loginResponseMsg
func (m loginModel) login() tea.Msg {
	time.Sleep(time.Second * 3)
	return loginResponseMsg{Status: 200, AuthToken: ""}
}

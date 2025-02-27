package tui

import (
	"fmt"
	"github.com/juancwu/konbini/cli/services"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type verifyEmailKeyMap struct {
	Resend key.Binding
	Quit   key.Binding
}

func (k verifyEmailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Resend, k.Quit}
}

func (k verifyEmailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Resend},
		{k.Quit},
	}
}

type verifyEmail struct {
	countDown time.Duration
	done      bool
	keys      verifyEmailKeyMap
	help      help.Model

	err error
}

func newVerifyEmail(_ map[string]interface{}) verifyEmail {

	keys := verifyEmailKeyMap{
		Resend: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "Resend verification email"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "Quit"),
		),
	}

	return verifyEmail{
		keys:      keys,
		help:      help.New(),
		done:      true,
		countDown: time.Second * 60,
	}
}

func (m verifyEmail) Init() tea.Cmd {
	return m.tick()
}

func (m verifyEmail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Resend):
			if m.done {
				m.done = false
				return m, m.resend
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	case tickMsg:
		m.countDown -= time.Second
		if m.countDown <= 0 && !m.done {
			m.done = true
			m.countDown = time.Second * 60
			return m, nil
		}
		return m, m.tick()
	case resendMsg:
		if msg != nil {
			m.err = msg
		}
	}
	return m, nil
}

func (m verifyEmail) View() string {
	parts := []string{
		"Please verify your email to continue. Check your email inbox and spam.",
	}

	if m.err != nil {
		parts = append(parts, errTextStyle.Render(m.err.Error()))
	}

	if m.done {
		parts = append(parts, "Press 'R' to resend the verification email.")
	} else {
		parts = append(parts, fmt.Sprintf("Wait %ss before resending a new verification email.", m.countDown))
	}

	return strings.Join(parts, "\n\n")
}

type resendMsg error

func (m verifyEmail) resend() tea.Msg {
	err := services.ResendVerificationEmail()
	return resendMsg(err)
}

type tickMsg time.Time

func (m verifyEmail) tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

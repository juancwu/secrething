package models

import tea "github.com/charmbracelet/bubbletea/v2"

type errMsg struct {
	err error
}

func newErrMsg(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg{err: err}
	}
}

func (e errMsg) Err() error {
	return e.err
}

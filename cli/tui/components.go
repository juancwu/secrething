package tui

import "github.com/charmbracelet/bubbles/spinner"

// newSpinner returns a new spinner.Model
func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return s
}

package tui

import "github.com/charmbracelet/lipgloss"

var (
	menuItemStyle = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2)

	errTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)

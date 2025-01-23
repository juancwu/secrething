package app

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type debugOverlay struct {
	width  int
	height int
}

var (
	overlayStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("87")).
			Bold(true)
)

func newDebugOverlay(width, height int) debugOverlay {
	return debugOverlay{
		width:  width,
		height: height,
	}
}

func (d debugOverlay) View(profile debugProfile, routerHistory string) string {
	// Calculate dimensions for the overlay
	overlayWidth := d.width / 2
	overlayHeight := d.height / 2

	var content strings.Builder
	content.WriteString(titleStyle.Render("Debug Information\n"))
	content.WriteString("Last Input: " + profile.Input + "\n")
	content.WriteString("Route History:\n" + routerHistory + "\n")

	// Center the overlay
	return lipgloss.Place(
		d.width,
		d.height,
		lipgloss.Center,
		lipgloss.Center,
		overlayStyle.Width(overlayWidth).Height(overlayHeight).Render(content.String()),
	)
}

package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// Catppuccin Mocha palette
const (
	colorBase    = lipgloss.Color("#1e1e2e")
	colorSurface = lipgloss.Color("#313244")
	colorOverlay = lipgloss.Color("#45475a")
	colorText    = lipgloss.Color("#cdd6f4")
	colorSubtext = lipgloss.Color("#a6adc8")
	colorBlue    = lipgloss.Color("#89b4fa")
	colorGreen   = lipgloss.Color("#a6e3a1")
	colorYellow  = lipgloss.Color("#f9e2af")
	colorRed     = lipgloss.Color("#f38ba8")
	colorMauve   = lipgloss.Color("#cba6f7")
	colorPeach   = lipgloss.Color("#fab387")
	colorTeal    = lipgloss.Color("#94e2d5")
)

var (
	StylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorOverlay)

	StyleTitle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	StyleSelected = lipgloss.NewStyle().
			Foreground(colorBase).
			Background(colorBlue).
			Bold(true)

	StyleDimmed = lipgloss.NewStyle().
			Foreground(colorSubtext)

	StyleTimer = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(colorRed)

	StyleProject = lipgloss.NewStyle().
			Foreground(colorMauve).
			Bold(true)

	StyleTask = lipgloss.NewStyle().
			Foreground(colorText)

	StyleDuration = lipgloss.NewStyle().
			Foreground(colorPeach)

	StyleInputLabel = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)
)

// catppuccinHelpStyles returns help.Styles using the Mocha palette.
func catppuccinHelpStyles() help.Styles {
	keyStyle  := lipgloss.NewStyle().Foreground(colorBlue)
	descStyle := lipgloss.NewStyle().Foreground(colorSubtext)
	sepStyle  := lipgloss.NewStyle().Foreground(colorOverlay)
	return help.Styles{
		ShortKey:       keyStyle,
		ShortDesc:      descStyle,
		ShortSeparator: sepStyle,
		Ellipsis:       descStyle,
		FullKey:        keyStyle,
		FullDesc:       descStyle,
		FullSeparator:  sepStyle,
	}
}

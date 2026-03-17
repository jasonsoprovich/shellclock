package ui

import "github.com/charmbracelet/lipgloss"

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
)

var (
	StylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorOverlay).
			Padding(0, 1)

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

	StyleHelp = lipgloss.NewStyle().
			Foreground(colorSubtext)

	StyleProject = lipgloss.NewStyle().
			Foreground(colorMauve).
			Bold(true)

	StyleTask = lipgloss.NewStyle().
			Foreground(colorText)

	StyleDuration = lipgloss.NewStyle().
			Foreground(colorPeach)
)

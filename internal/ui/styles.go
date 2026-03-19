package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// Current palette — mutable so ApplyTheme can update them at runtime.
// Initialised to Catppuccin Mocha values.
var (
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

// Named styles — rebuilt by ApplyTheme when the theme changes.
// Initialised to Catppuccin Mocha at package load.
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

	StyleLogo = lipgloss.NewStyle().
			Foreground(colorMauve)
)

// helpStyles returns help.Styles built from the current palette variables.
// Called inside every View() so theme changes take effect on the next render
// without needing to store styles in the model.
func helpStyles() help.Styles {
	return help.Styles{
		ShortKey:       lipgloss.NewStyle().Foreground(colorBlue),
		ShortDesc:      lipgloss.NewStyle().Foreground(colorSubtext),
		ShortSeparator: lipgloss.NewStyle().Foreground(colorOverlay),
		Ellipsis:       lipgloss.NewStyle().Foreground(colorSubtext),
		FullKey:        lipgloss.NewStyle().Foreground(colorBlue),
		FullDesc:       lipgloss.NewStyle().Foreground(colorSubtext),
		FullSeparator:  lipgloss.NewStyle().Foreground(colorOverlay),
	}
}

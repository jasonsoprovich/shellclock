package ui

import "github.com/charmbracelet/lipgloss"

// Theme holds the complete color palette for one visual theme.
// Each field maps to a semantic role used across all views.
type Theme struct {
	Name    string         // identifier stored in JSON (e.g. "catppuccin-mocha")
	Label   string         // human-readable display name
	Dark    bool           // true for dark backgrounds
	Base    lipgloss.Color // panel / background
	Surface lipgloss.Color // selection highlight background
	Overlay lipgloss.Color // border, dividers, dimmed rules
	Text    lipgloss.Color // primary text (StyleTask)
	Subtext lipgloss.Color // secondary / dimmed text (StyleDimmed)
	Blue    lipgloss.Color // titles, key-hint keys, selection highlight bg
	Green   lipgloss.Color // live timer, progress-bar fill (StyleTimer)
	Yellow  lipgloss.Color // input labels (StyleInputLabel)
	Red     lipgloss.Color // validation errors (StyleError)
	Mauve   lipgloss.Color // project names (StyleProject)
	Peach   lipgloss.Color // durations (StyleDuration)
	Teal    lipgloss.Color // palette completeness; used in swatch display
}

// allThemes is the ordered list shown in the theme picker.
var allThemes = []Theme{
	{
		Name: "catppuccin-mocha", Label: "Catppuccin Mocha", Dark: true,
		Base: "#1e1e2e", Surface: "#313244", Overlay: "#45475a",
		Text: "#cdd6f4", Subtext: "#a6adc8",
		Blue: "#89b4fa", Green: "#a6e3a1", Yellow: "#f9e2af",
		Red: "#f38ba8", Mauve: "#cba6f7", Peach: "#fab387", Teal: "#94e2d5",
	},
	{
		Name: "catppuccin-latte", Label: "Catppuccin Latte", Dark: false,
		Base: "#eff1f5", Surface: "#ccd0da", Overlay: "#9ca0b0",
		Text: "#4c4f69", Subtext: "#6c6f85",
		Blue: "#1e66f5", Green: "#40a02b", Yellow: "#df8e1d",
		Red: "#d20f39", Mauve: "#8839ef", Peach: "#fe640b", Teal: "#179299",
	},
	{
		Name: "dracula", Label: "Dracula", Dark: true,
		Base: "#282a36", Surface: "#44475a", Overlay: "#6272a4",
		Text: "#f8f8f2", Subtext: "#6272a4",
		Blue: "#bd93f9", Green: "#50fa7b", Yellow: "#f1fa8c",
		Red: "#ff5555", Mauve: "#ff79c6", Peach: "#ffb86c", Teal: "#8be9fd",
	},
	{
		Name: "nord", Label: "Nord", Dark: true,
		Base: "#2e3440", Surface: "#3b4252", Overlay: "#4c566a",
		Text: "#eceff4", Subtext: "#d8dee9",
		Blue: "#88c0d0", Green: "#a3be8c", Yellow: "#ebcb8b",
		Red: "#bf616a", Mauve: "#b48ead", Peach: "#d08770", Teal: "#8fbcbb",
	},
	{
		Name: "gruvbox", Label: "Gruvbox", Dark: true,
		Base: "#282828", Surface: "#3c3836", Overlay: "#504945",
		Text: "#ebdbb2", Subtext: "#a89984",
		Blue: "#83a598", Green: "#b8bb26", Yellow: "#fabd2f",
		Red: "#fb4934", Mauve: "#d3869b", Peach: "#fe8019", Teal: "#8ec07c",
	},
	{
		Name: "tokyo-night", Label: "Tokyo Night", Dark: true,
		Base: "#1a1b26", Surface: "#24283b", Overlay: "#414868",
		Text: "#c0caf5", Subtext: "#565f89",
		Blue: "#7aa2f7", Green: "#9ece6a", Yellow: "#e0af68",
		Red: "#f7768e", Mauve: "#bb9af7", Peach: "#ff9e64", Teal: "#73daca",
	},
}

// activeTheme is the currently applied theme.
var activeTheme = allThemes[0] // Catppuccin Mocha

// AllThemes returns the ordered slice of all available themes.
func AllThemes() []Theme { return allThemes }

// ThemeByName returns the theme with the given name, falling back to
// Catppuccin Mocha when name is empty or not found.
func ThemeByName(name string) Theme {
	for _, t := range allThemes {
		if t.Name == name {
			return t
		}
	}
	return allThemes[0]
}

// ThemeIndex returns the index of the theme with the given name in allThemes,
// or 0 (Catppuccin Mocha) if not found.
func ThemeIndex(name string) int {
	for i, t := range allThemes {
		if t.Name == name {
			return i
		}
	}
	return 0
}

// ApplyTheme updates all global palette variables and named style vars to
// match t, then records t as the active theme.  Must be called from the main
// goroutine (Bubble Tea update loop).
func ApplyTheme(t Theme) {
	activeTheme = t

	// Update palette vars so helpStyles() and any direct color references
	// (e.g. timer.go's paused-clock style) pick up the new values.
	colorBase = t.Base
	colorSurface = t.Surface
	colorOverlay = t.Overlay
	colorText = t.Text
	colorSubtext = t.Subtext
	colorBlue = t.Blue
	colorGreen = t.Green
	colorYellow = t.Yellow
	colorRed = t.Red
	colorMauve = t.Mauve
	colorPeach = t.Peach
	colorTeal = t.Teal

	// Rebuild every named style.
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
}

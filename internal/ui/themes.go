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
	Teal    lipgloss.Color // palette completeness; shown in swatch row
}

// allThemes is the ordered list shown in the theme picker.
// Catppuccin variants appear first, remaining themes are alphabetical.
var allThemes = []Theme{

	// ── Catppuccin ────────────────────────────────────────────────────────

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

	// ── A ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/dempfi/ayu — light variant
		Name: "ayu-light", Label: "Ayu Light", Dark: false,
		Base: "#fafafa", Surface: "#f3f3f3", Overlay: "#c5c5c5",
		Text: "#5c6166", Subtext: "#8a9199",
		Blue: "#399ee6", Green: "#86b300", Yellow: "#f2ae49",
		Red: "#f07171", Mauve: "#a37acc", Peach: "#fa8d3e", Teal: "#4cbf99",
	},
	{
		// https://github.com/dempfi/ayu — mirage (dark) variant
		Name: "ayu-mirage", Label: "Ayu Mirage", Dark: true,
		Base: "#1f2430", Surface: "#242b38", Overlay: "#3d4752",
		Text: "#cccac2", Subtext: "#607080",
		Blue: "#73d0ff", Green: "#bae67e", Yellow: "#ffd173",
		Red: "#f28779", Mauve: "#dfbfff", Peach: "#ffad66", Teal: "#5ccfe6",
	},

	// ── C ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/wesbos/cobalt2 — Cobalt2 by Wes Bos
		Name: "cobalt2", Label: "Cobalt2", Dark: true,
		Base: "#193549", Surface: "#1f4662", Overlay: "#2e6890",
		Text: "#ffffff", Subtext: "#7ea8c9",
		Blue: "#0088ff", Green: "#3ad900", Yellow: "#ffc600",
		Red: "#ff628c", Mauve: "#9effff", Peach: "#ff9d00", Teal: "#80fcff",
	},

	// ── D ─────────────────────────────────────────────────────────────────

	{
		// https://draculatheme.com/
		Name: "dracula", Label: "Dracula", Dark: true,
		Base: "#282a36", Surface: "#44475a", Overlay: "#6272a4",
		Text: "#f8f8f2", Subtext: "#6272a4",
		Blue: "#bd93f9", Green: "#50fa7b", Yellow: "#f1fa8c",
		Red: "#ff5555", Mauve: "#ff79c6", Peach: "#ffb86c", Teal: "#8be9fd",
	},

	// ── E ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/sainnhe/everforest — dark hard
		Name: "everforest-dark", Label: "Everforest Dark", Dark: true,
		Base: "#1e2326", Surface: "#272e33", Overlay: "#374145",
		Text: "#d3c6aa", Subtext: "#859289",
		Blue: "#7fbbb3", Green: "#a7c080", Yellow: "#dbbc7f",
		Red: "#e67e80", Mauve: "#d699b6", Peach: "#e69875", Teal: "#83c092",
	},
	{
		// https://github.com/sainnhe/everforest — light soft
		Name: "everforest-light", Label: "Everforest Light", Dark: false,
		Base: "#f2efdf", Surface: "#eae4ca", Overlay: "#a6b0a0",
		Text: "#5c6a72", Subtext: "#829181",
		Blue: "#3a94c5", Green: "#8da101", Yellow: "#dfa000",
		Red: "#f85552", Mauve: "#df69ba", Peach: "#f57d26", Teal: "#35a77c",
	},

	// ── F ─────────────────────────────────────────────────────────────────

	{
		// https://stephango.com/flexoki — dark
		Name: "flexoki-dark", Label: "Flexoki Dark", Dark: true,
		Base: "#1c1b1a", Surface: "#282726", Overlay: "#575653",
		Text: "#cecdc3", Subtext: "#878580",
		Blue: "#4385be", Green: "#879a39", Yellow: "#d0a215",
		Red: "#d14d41", Mauve: "#8b7ec8", Peach: "#da702c", Teal: "#3aa99f",
	},
	{
		// https://stephango.com/flexoki — light
		Name: "flexoki-light", Label: "Flexoki Light", Dark: false,
		Base: "#fffcf0", Surface: "#f2f0e5", Overlay: "#b7b5ac",
		Text: "#100f0f", Subtext: "#6f6e69",
		Blue: "#4385be", Green: "#879a39", Yellow: "#d0a215",
		Red: "#d14d41", Mauve: "#8b7ec8", Peach: "#da702c", Teal: "#3aa99f",
	},

	// ── G ─────────────────────────────────────────────────────────────────

	{
		// https://primer.style/primitives — GitHub Dark
		Name: "github-dark", Label: "GitHub Dark", Dark: true,
		Base: "#0d1117", Surface: "#161b22", Overlay: "#30363d",
		Text: "#c9d1d9", Subtext: "#8b949e",
		Blue: "#58a6ff", Green: "#3fb950", Yellow: "#d29922",
		Red: "#f85149", Mauve: "#a371f7", Peach: "#e3b341", Teal: "#39c5cf",
	},
	{
		// https://primer.style/primitives — GitHub Light
		Name: "github-light", Label: "GitHub Light", Dark: false,
		Base: "#ffffff", Surface: "#f6f8fa", Overlay: "#d0d7de",
		Text: "#24292f", Subtext: "#57606a",
		Blue: "#0969da", Green: "#1a7f37", Yellow: "#9a6700",
		Red: "#cf222e", Mauve: "#8250df", Peach: "#953800", Teal: "#0550ae",
	},
	{
		// https://github.com/morhetz/gruvbox — dark
		Name: "gruvbox", Label: "Gruvbox", Dark: true,
		Base: "#282828", Surface: "#3c3836", Overlay: "#504945",
		Text: "#ebdbb2", Subtext: "#a89984",
		Blue: "#83a598", Green: "#b8bb26", Yellow: "#fabd2f",
		Red: "#fb4934", Mauve: "#d3869b", Peach: "#fe8019", Teal: "#8ec07c",
	},
	{
		// https://github.com/morhetz/gruvbox — light
		Name: "gruvbox-light", Label: "Gruvbox Light", Dark: false,
		Base: "#fbf1c7", Surface: "#ebdbb2", Overlay: "#928374",
		Text: "#3c3836", Subtext: "#7c6f64",
		Blue: "#458588", Green: "#98971a", Yellow: "#d79921",
		Red: "#cc241d", Mauve: "#b16286", Peach: "#d65d0e", Teal: "#689d6a",
	},

	// ── K ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/rebelot/kanagawa.nvim — wave
		Name: "kanagawa-wave", Label: "Kanagawa Wave", Dark: true,
		Base: "#1f1f28", Surface: "#2a2a37", Overlay: "#54546d",
		Text: "#dcd7ba", Subtext: "#727169",
		Blue: "#7e9cd8", Green: "#76946a", Yellow: "#c0a36e",
		Red: "#c34043", Mauve: "#957fb8", Peach: "#ffa066", Teal: "#6a9589",
	},
	{
		// https://github.com/rebelot/kanagawa.nvim — dragon
		Name: "kanagawa-dragon", Label: "Kanagawa Dragon", Dark: true,
		Base: "#181616", Surface: "#282727", Overlay: "#625a5a",
		Text: "#c5c9c5", Subtext: "#71716f",
		Blue: "#8ba4b0", Green: "#87a987", Yellow: "#c4b28a",
		Red: "#c4746e", Mauve: "#8992a7", Peach: "#b6927b", Teal: "#8ea4a2",
	},

	// ── M ─────────────────────────────────────────────────────────────────

	{
		// https://material-theme.com — palenight
		Name: "material-palenight", Label: "Material Palenight", Dark: true,
		Base: "#292d3e", Surface: "#32374d", Overlay: "#676e95",
		Text: "#a6accd", Subtext: "#676e95",
		Blue: "#82aaff", Green: "#c3e88d", Yellow: "#ffcb6b",
		Red: "#f07178", Mauve: "#c792ea", Peach: "#f78c6c", Teal: "#89ddff",
	},
	{
		// https://monokai.pro/
		Name: "monokai-pro", Label: "Monokai Pro", Dark: true,
		Base: "#2d2a2e", Surface: "#403e41", Overlay: "#5b595c",
		Text: "#fcfcfa", Subtext: "#939293",
		Blue: "#78dce8", Green: "#a9dc76", Yellow: "#ffd866",
		Red: "#ff6188", Mauve: "#ab9df2", Peach: "#fc9867", Teal: "#a1efe4",
	},

	// ── N ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/EdenEast/nightfox.nvim — nightfox
		Name: "nightfox", Label: "Nightfox", Dark: true,
		Base: "#192330", Surface: "#212e3f", Overlay: "#29394f",
		Text: "#cdcecf", Subtext: "#738091",
		Blue: "#719cd6", Green: "#81b29a", Yellow: "#dbc074",
		Red: "#c94f6d", Mauve: "#9d79d6", Peach: "#f4a261", Teal: "#63cdcf",
	},
	{
		// https://www.nordtheme.com/
		Name: "nord", Label: "Nord", Dark: true,
		Base: "#2e3440", Surface: "#3b4252", Overlay: "#4c566a",
		Text: "#eceff4", Subtext: "#d8dee9",
		Blue: "#88c0d0", Green: "#a3be8c", Yellow: "#ebcb8b",
		Red: "#bf616a", Mauve: "#b48ead", Peach: "#d08770", Teal: "#8fbcbb",
	},

	// ── O ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/Binaryify/OneDark-Pro
		Name: "one-dark-pro", Label: "One Dark Pro", Dark: true,
		Base: "#282c34", Surface: "#3e4451", Overlay: "#4b5263",
		Text: "#abb2bf", Subtext: "#5c6370",
		Blue: "#61afef", Green: "#98c379", Yellow: "#e5c07b",
		Red: "#e06c75", Mauve: "#c678dd", Peach: "#d19a66", Teal: "#56b6c2",
	},
	{
		// https://github.com/atom/one-light-syntax
		Name: "one-light", Label: "One Light", Dark: false,
		Base: "#fafafa", Surface: "#e5e5e6", Overlay: "#a0a1a7",
		Text: "#383a42", Subtext: "#696c77",
		Blue: "#4078f2", Green: "#50a14f", Yellow: "#c18401",
		Red: "#e45649", Mauve: "#a626a4", Peach: "#986801", Teal: "#0184bc",
	},
	{
		// https://github.com/nyoom-engineering/oxocarbon.nvim (IBM Carbon)
		Name: "oxocarbon", Label: "Oxocarbon", Dark: true,
		Base: "#161616", Surface: "#262626", Overlay: "#393939",
		Text: "#f2f4f8", Subtext: "#525252",
		Blue: "#78a9ff", Green: "#42be65", Yellow: "#ffe97b",
		Red: "#fa4d56", Mauve: "#be95ff", Peach: "#ff832b", Teal: "#3ddbd9",
	},

	// ── P ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/olivercederborg/poimandres.nvim
		Name: "poimandres", Label: "Poimandres", Dark: true,
		Base: "#1b1e28", Surface: "#252734", Overlay: "#3d3f4e",
		Text: "#a6accd", Subtext: "#767c9d",
		Blue: "#89ddff", Green: "#5de4c7", Yellow: "#fffac2",
		Red: "#d0679d", Mauve: "#a277ff", Peach: "#e9a16b", Teal: "#add7ff",
	},

	// ── R ─────────────────────────────────────────────────────────────────

	{
		// https://rosepinetheme.com/ — main (dark)
		Name: "rose-pine", Label: "Rosé Pine", Dark: true,
		Base: "#191724", Surface: "#1f1d2e", Overlay: "#6e6a86",
		Text: "#e0def4", Subtext: "#908caa",
		Blue: "#c4a7e7", Green: "#9ccfd8", Yellow: "#f6c177",
		Red: "#eb6f92", Mauve: "#ebbcba", Peach: "#31748f", Teal: "#9ccfd8",
	},
	{
		// https://rosepinetheme.com/ — dawn (light)
		Name: "rose-pine-dawn", Label: "Rosé Pine Dawn", Dark: false,
		Base: "#faf4ed", Surface: "#f2e9e1", Overlay: "#9893a5",
		Text: "#575279", Subtext: "#797593",
		Blue: "#286983", Green: "#56949f", Yellow: "#ea9d34",
		Red: "#b4637a", Mauve: "#907aa9", Peach: "#d7827e", Teal: "#56949f",
	},

	// ── S ─────────────────────────────────────────────────────────────────

	{
		// https://ethanschoonover.com/solarized/ — dark
		Name: "solarized-dark", Label: "Solarized Dark", Dark: true,
		Base: "#002b36", Surface: "#073642", Overlay: "#586e75",
		Text: "#839496", Subtext: "#657b83",
		Blue: "#268bd2", Green: "#859900", Yellow: "#b58900",
		Red: "#dc322f", Mauve: "#6c71c4", Peach: "#cb4b16", Teal: "#2aa198",
	},
	{
		// https://ethanschoonover.com/solarized/ — light
		Name: "solarized-light", Label: "Solarized Light", Dark: false,
		Base: "#fdf6e3", Surface: "#eee8d5", Overlay: "#93a1a1",
		Text: "#657b83", Subtext: "#839496",
		Blue: "#268bd2", Green: "#859900", Yellow: "#b58900",
		Red: "#dc322f", Mauve: "#6c71c4", Peach: "#cb4b16", Teal: "#2aa198",
	},

	// ── T ─────────────────────────────────────────────────────────────────

	{
		// https://github.com/tokyo-night/tokyo-night-vscode-theme
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

	StyleLogo = lipgloss.NewStyle().
		Foreground(colorMauve)
}

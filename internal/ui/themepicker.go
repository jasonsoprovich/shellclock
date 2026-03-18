package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jasonsoprovich/shellclock/internal/model"
)

// ThemePickerModel is the theme-selection view.
//
// Navigation (↑/↓ or k/j) live-previews each theme by calling ApplyTheme
// immediately.  Pressing Esc reverts to the theme that was active when the
// picker opened; pressing Enter confirms and persists to disk.
type ThemePickerModel struct {
	store         *model.Store
	keys          KeyMap
	cursor        int
	width         int
	height        int
	help          help.Model
	previousTheme Theme // restored on Esc

	SwitchToTree bool
}

func NewThemePickerModel(store *model.Store, keys KeyMap) ThemePickerModel {
	h := help.New()
	h.Styles = helpStyles()
	return ThemePickerModel{
		store:         store,
		keys:          keys,
		cursor:        ThemeIndex(activeTheme.Name),
		help:          h,
		previousTheme: activeTheme,
	}
}

// ── Update ──────────────────────────────────────────────────────────────────

func (m ThemePickerModel) Update(msg tea.Msg) (ThemePickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				ApplyTheme(allThemes[m.cursor])
			}

		case "down", "j":
			if m.cursor < len(allThemes)-1 {
				m.cursor++
				ApplyTheme(allThemes[m.cursor])
			}

		case "enter":
			// Confirm: persist name to store.
			m.store.Theme = activeTheme.Name
			_ = m.store.Save()
			m.SwitchToTree = true

		case "esc", "q":
			// Cancel: revert to the theme that was active on entry.
			ApplyTheme(m.previousTheme)
			m.SwitchToTree = true
		}
	}
	return m, nil
}

// ── View ─────────────────────────────────────────────────────────────────────

// labelWidth is the display width of the longest theme label, used for column
// alignment.  "Catppuccin Mocha" / "Catppuccin Latte" = 16 chars.
const labelWidth = 16

// swatchColors returns the seven representative colors shown as ■ squares.
func swatchColors(t Theme) []lipgloss.Color {
	return []lipgloss.Color{t.Blue, t.Green, t.Yellow, t.Red, t.Mauve, t.Peach, t.Teal}
}

func (m ThemePickerModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	innerW := w - 4
	if innerW < 36 {
		innerW = 36
	}

	var sb strings.Builder

	// ── Header ─────────────────────────────────────────────────────────────
	sb.WriteString(StyleTitle.Render("Themes"))
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render(strings.Repeat("─", innerW)))
	sb.WriteString("\n\n")

	// ── Theme list ─────────────────────────────────────────────────────────
	for i, t := range allThemes {
		selected := i == m.cursor

		// Indicator — 2 visible chars.
		var indicator string
		if selected {
			indicator = StyleTimer.Render("▸ ")
		} else {
			indicator = "  "
		}

		// Label — padded to labelWidth.
		var nameStyle lipgloss.Style
		if selected {
			nameStyle = StyleTitle
		} else {
			nameStyle = StyleDimmed
		}
		// Render the styled name then pad the *visible* width manually so ANSI
		// codes don't upset the column alignment.
		styledName := nameStyle.Render(t.Label)
		namePad := labelWidth - lipgloss.Width(t.Label)
		if namePad < 0 {
			namePad = 0
		}

		// Swatches — each ■ rendered in its own theme's color, independent of
		// the currently applied (previewed) theme.
		var swatchSB strings.Builder
		for _, c := range swatchColors(t) {
			swatchSB.WriteString(lipgloss.NewStyle().Foreground(c).Render("■"))
		}
		swatches := swatchSB.String()

		sb.WriteString(indicator)
		sb.WriteString(styledName)
		sb.WriteString(strings.Repeat(" ", namePad+2)) // gap after label
		sb.WriteString(swatches)
		sb.WriteString("\n")
	}

	// ── Instruction line ────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render("navigate to preview · enter to apply · esc to cancel"))
	sb.WriteString("\n")

	// ── Help bar ─────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	m.help.Styles = helpStyles()
	m.help.Width = innerW
	sb.WriteString(m.help.View(themePickerKeyMap{m.keys}))

	return StylePanel.Width(innerW).Padding(0, 1).Render(sb.String())
}

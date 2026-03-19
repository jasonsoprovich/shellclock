package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderTagPickerOverlay composites a centered tag-picker modal over a dimmed
// copy of bg. tags is the sorted list of all unique tags; selectedIdx is the
// currently highlighted row (0 = "show all", 1+ = tags[i-1]); activeFilter is
// the tag currently applied (empty = none).
func renderTagPickerOverlay(bg string, tags []string, selectedIdx int, activeFilter string, termW, termH int) string {
	if termW <= 0 {
		termW = 80
	}
	if termH <= 0 {
		termH = 24
	}

	titleStyle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true)

	var lines []string
	lines = append(lines, titleStyle.Render("Filter by Tag"))
	lines = append(lines, "")

	// "(show all)" row — index 0.
	if selectedIdx == 0 {
		lines = append(lines, StyleSelected.Render(" (show all) "))
	} else if activeFilter == "" {
		lines = append(lines, StyleProject.Render("  (show all)") + " " + StyleTimer.Render("✓"))
	} else {
		lines = append(lines, StyleDimmed.Render("  (show all)"))
	}

	for i, tag := range tags {
		idx := i + 1
		if idx == selectedIdx {
			lines = append(lines, StyleSelected.Render(" "+tag+" "))
		} else {
			entry := "  " + renderTagPill(tag)
			if tag == activeFilter {
				entry += " " + StyleTimer.Render("✓")
			}
			lines = append(lines, entry)
		}
	}

	lines = append(lines, "")
	lines = append(lines, StyleDimmed.Render("↑/↓ navigate   enter select   esc cancel"))

	content := strings.Join(lines, "\n")

	modal := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue).
		Background(colorBase).
		Padding(1, 4).
		Render(content)

	// Dim the background.
	const dimColor = "\033[2m\033[38;2;80;80;100m"
	const rst = "\033[0m"

	bgLines := strings.Split(bg, "\n")
	for len(bgLines) < termH {
		bgLines = append(bgLines, "")
	}
	dimmed := make([]string, len(bgLines))
	for i, ln := range bgLines {
		dimmed[i] = dimColor + stripAnsi(ln) + rst
	}

	modalLines := strings.Split(modal, "\n")
	mH := len(modalLines)
	mW := 0
	for _, l := range modalLines {
		if w := lipgloss.Width(l); w > mW {
			mW = w
		}
	}

	startRow := (termH - mH) / 2
	startCol := (termW - mW) / 2
	if startRow < 0 {
		startRow = 0
	}
	if startCol < 0 {
		startCol = 0
	}

	for i, mLine := range modalLines {
		idx := startRow + i
		if idx >= 0 && idx < len(dimmed) {
			dimmed[idx] = overlayLine(dimmed[idx], mLine, startCol)
		}
	}

	return strings.Join(dimmed, "\n")
}

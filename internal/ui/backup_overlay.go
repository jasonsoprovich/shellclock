package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderBackupOverlay composites a centered backup-info modal over a dimmed
// copy of bg. backups is the list of backup file names (newest first).
func renderBackupOverlay(bg string, backups []string, termW, termH int) string {
	if termW <= 0 {
		termW = 80
	}
	if termH <= 0 {
		termH = 24
	}

	titleStyle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true)
	pathStyle := lipgloss.NewStyle().Foreground(colorTeal)
	dimStyle := lipgloss.NewStyle().Foreground(colorSubtext)

	title := titleStyle.Render("Backups")
	pathLine := pathStyle.Render("~/.config/shellclock/backups/")

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")
	lines = append(lines, dimStyle.Render("Backups are stored at"))
	lines = append(lines, pathLine)
	lines = append(lines, "")

	if len(backups) == 0 {
		lines = append(lines, dimStyle.Render("No backups yet."))
	} else {
		for _, name := range backups {
			lines = append(lines, StyleTask.Render("  "+name))
		}
	}

	lines = append(lines, "")
	lines = append(lines, dimStyle.Render("[esc] / any key  close"))

	content := strings.Join(lines, "\n")

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue).
		Background(colorBase).
		Padding(1, 4)

	modal := modalStyle.Render(content)

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
		if idx < 0 || idx >= len(dimmed) {
			continue
		}
		dimmed[idx] = overlayLine(dimmed[idx], mLine, startCol)
	}

	return strings.Join(dimmed, "\n")
}

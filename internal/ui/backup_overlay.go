package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderBackupOverlay composites an interactive backup picker over a dimmed
// copy of bg. cursor is the index of the selected backup; confirmActive
// switches the footer to the restore-confirmation prompt.
func renderBackupOverlay(bg string, backups []string, cursor int, confirmActive bool, termW, termH int) string {
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
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("[esc] close"))
	} else {
		for i, name := range backups {
			if i == cursor {
				lines = append(lines, StyleSelected.Render("  "+name+"  "))
			} else {
				lines = append(lines, StyleTask.Render("  "+name))
			}
		}
		lines = append(lines, "")
		if confirmActive {
			lines = append(lines, lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("Restore this backup? Current data will be overwritten."))
			lines = append(lines, StyleTimer.Render("[y]")+StyleTask.Render(" restore")+"   "+dimStyle.Render("[n] / [esc] cancel"))
		} else {
			lines = append(lines, dimStyle.Render("↑/↓ select   enter restore   esc close"))
		}
	}

	content := strings.Join(lines, "\n")

	borderColor := colorBlue
	if confirmActive {
		borderColor = colorYellow
	}
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(colorBase).
		Padding(1, 4)

	modal := modalStyle.Render(content)

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

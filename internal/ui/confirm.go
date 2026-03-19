package ui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ansiRe matches ANSI SGR escape sequences (colours, bold, dim, reset, etc.)
var ansiRe = regexp.MustCompile(`\033\[[0-9;]*[A-Za-z]`)

// stripAnsi returns s with all ANSI escape sequences removed.
func stripAnsi(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

// renderConfirmOverlay composites a centered delete-confirmation modal over a
// dimmed copy of bg. msg is the first line of the modal (e.g. `Delete project "Foo"?`).
// termW / termH are the full terminal dimensions.
func renderConfirmOverlay(bg, msg string, termW, termH int) string {
	if termW <= 0 {
		termW = 80
	}
	if termH <= 0 {
		termH = 24
	}

	// ── Build modal content ────────────────────────────────────────────────
	titleStyle := lipgloss.NewStyle().Foreground(colorRed).Bold(true)
	title := titleStyle.Render(msg)
	sub := StyleDimmed.Render("This cannot be undone.")
	keys := StyleTimer.Render("[y]") + StyleTask.Render(" confirm") +
		"   " + StyleDimmed.Render("[n] / [esc] cancel")

	content := title + "\n" + sub + "\n\n" + keys

	// ── Modal box ──────────────────────────────────────────────────────────
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorRed).
		Background(colorBase).
		Padding(1, 4)

	modal := modalStyle.Render(content)

	// ── Dim the background: strip ANSI, re-render with uniform dim style ───
	const dimColor = "\033[2m\033[38;2;80;80;100m"
	const rst = "\033[0m"

	bgLines := strings.Split(bg, "\n")
	// Ensure we have a full-height canvas.
	for len(bgLines) < termH {
		bgLines = append(bgLines, "")
	}
	dimmed := make([]string, len(bgLines))
	for i, ln := range bgLines {
		dimmed[i] = dimColor + stripAnsi(ln) + rst
	}

	// ── Center and overlay the modal ───────────────────────────────────────
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

// renderResetOverlay composites a centered system-reset confirmation modal over
// a dimmed copy of bg. inputView is the rendered text input the user types into.
func renderResetOverlay(bg, inputView string, termW, termH int) string {
	if termW <= 0 {
		termW = 80
	}
	if termH <= 0 {
		termH = 24
	}

	titleStyle := lipgloss.NewStyle().Foreground(colorRed).Bold(true)
	title := titleStyle.Render("⚠  System Reset")
	warn1 := StyleError.Render("This will permanently erase ALL projects, tasks, and sessions.")
	warn2 := StyleDimmed.Render("Backups are not affected. This action cannot be undone.")
	prompt := StyleInputLabel.Render("Type CONFIRM to proceed:")
	cancel := StyleDimmed.Render("[esc] cancel   [enter] execute if CONFIRM is typed")

	content := title + "\n\n" +
		warn1 + "\n" +
		warn2 + "\n\n" +
		prompt + "\n" +
		inputView + "\n\n" +
		cancel

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorRed).
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

// overlayLine inserts fg starting at visible column col within the already-dimmed
// bg line.  bg is plain text wrapped in dim ANSI codes; we strip those to find
// column positions, then re-wrap the left/right portions around fg.
func overlayLine(bg, fg string, col int) string {
	const dimColor = "\033[2m\033[38;2;80;80;100m"
	const rst = "\033[0m"

	plain := []rune(stripAnsi(bg))

	leftEnd := col
	if leftEnd > len(plain) {
		leftEnd = len(plain)
	}
	left := dimColor + string(plain[:leftEnd]) + rst

	// Pad with spaces if col extends past the end of the line.
	if col > len(plain) {
		left += strings.Repeat(" ", col-len(plain))
	}

	rightStart := col + lipgloss.Width(fg)
	var right string
	if rightStart < len(plain) {
		right = dimColor + string(plain[rightStart:]) + rst
	}

	return left + fg + right
}

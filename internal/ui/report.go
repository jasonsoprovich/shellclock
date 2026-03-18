package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

// ── Row types ──────────────────────────────────────────────────────────────

type reportRowKind int

const (
	reportRowProject reportRowKind = iota
	reportRowTask
	reportRowBlank
)

type reportRow struct {
	kind      reportRowKind
	projectID string
	taskID    string // empty for project rows
}

// ── Model ──────────────────────────────────────────────────────────────────

// ReportModel renders the time-tracking summary.
type ReportModel struct {
	store    *model.Store
	keys     KeyMap
	width    int
	height   int
	help     help.Model
	showFull bool

	rows   []reportRow
	offset int // first visible row index

	SwitchToTree bool
}

func NewReportModel(store *model.Store, keys KeyMap) ReportModel {
	h := help.New()
	h.Styles = helpStyles()
	m := ReportModel{store: store, keys: keys, help: h}
	m.buildRows()
	return m
}

func (m *ReportModel) buildRows() {
	m.rows = nil
	for _, p := range m.store.Projects {
		m.rows = append(m.rows, reportRow{kind: reportRowProject, projectID: p.ID})
		for _, t := range p.Tasks {
			m.rows = append(m.rows, reportRow{kind: reportRowTask, projectID: p.ID, taskID: t.ID})
		}
		m.rows = append(m.rows, reportRow{kind: reportRowBlank})
	}
}

// listHeight returns the number of data rows the view can display at once.
//
// Fixed content lines (always present):
//
//	header: title+total (1) + rule (1)     = 2
//	scroll hint (always reserved)          = 1
//	blank + help                           = 2
//	border                                 = 2
//	total                                  = 7
//
// Optional lines added to fixed:
//
//	active timer notice                    +1
//	full help vs short help                +3
func (m *ReportModel) listHeight() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	fixed := 7
	if m.store.ActiveTimer != nil {
		fixed++ // timer notice line is always written when timer is active
	}
	if m.showFull {
		fixed += 3
	}
	lh := h - fixed
	if lh < 1 {
		lh = 1
	}
	return lh
}

func (m *ReportModel) scrollUp() {
	if m.offset > 0 {
		m.offset--
	}
}

func (m *ReportModel) scrollDown() {
	lh := m.listHeight()
	if m.offset+lh < len(m.rows) {
		m.offset++
	}
}

// ── Update ─────────────────────────────────────────────────────────────────

func (m ReportModel) Update(msg tea.Msg) (ReportModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width - 4
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "R":
			m.SwitchToTree = true
		case "up", "k":
			m.scrollUp()
		case "down", "j":
			m.scrollDown()
		case "?":
			m.showFull = !m.showFull
			m.help.ShowAll = m.showFull
		}
	}
	return m, nil
}

// ── View ───────────────────────────────────────────────────────────────────

func (m ReportModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	innerW := w - 4 // border(2) + padding(2)
	if innerW < 30 {
		innerW = 30
	}

	// Column widths.
	//   durW  = 11  (right-aligned duration, e.g. "1h 23m 45s")
	//   barW  = 20  (progress bar)
	//   gap   = 2   (space between bar and duration)
	//   nameW = remainder
	const durW, barW, gap = 11, 20, 2
	nameW := innerW - barW - gap - durW
	if nameW < 10 {
		nameW = 10
	}

	// Compute grand total across all committed sessions.
	var grandTotal int64
	for i := range m.store.Projects {
		grandTotal += m.store.Projects[i].TotalSeconds()
	}

	// ── Header ────────────────────────────────────────────────────────────
	var header strings.Builder
	header.WriteString(StyleTitle.Render("Report"))

	totalStr := util.FormatDuration(grandTotal)
	if grandTotal == 0 {
		totalStr = "no time tracked"
	}
	headerRight := StyleDimmed.Render("total  ") + StyleDuration.Render(totalStr)
	// Right-align the total in the header line.
	titleW := lipgloss.Width("Report")
	rightW := lipgloss.Width(headerRight)
	headerPad := innerW - titleW - rightW
	if headerPad > 0 {
		header.WriteString(strings.Repeat(" ", headerPad))
	} else {
		header.WriteString("  ")
	}
	header.WriteString(headerRight)
	header.WriteString("\n")
	header.WriteString(StyleDimmed.Render(strings.Repeat("─", innerW)))
	header.WriteString("\n")

	// Active-timer notice — always written as a line when timer is active so
	// listHeight() can safely add 1 to fixed when store.ActiveTimer != nil.
	if m.store.ActiveTimer != nil {
		at := m.store.ActiveTimer
		p := m.store.FindProject(at.ProjectID)
		t := m.store.FindTask(at.ProjectID, at.TaskID)
		pName, tName := "?", "?"
		if p != nil {
			pName = truncate(p.Name, 20)
		}
		if t != nil {
			tName = truncate(t.Name, 20)
		}
		// Only add time.Since when the timer is not paused; AccumulatedSeconds
		// already banks the interval up to the most recent pause.
		elapsed := at.AccumulatedSeconds
		if !at.Paused {
			elapsed += int64(time.Since(at.Start).Seconds())
		}
		header.WriteString(
			StyleTimer.Render("⚡ ") +
				StyleDimmed.Render(fmt.Sprintf("%s › %s  %s  (not yet saved)",
					pName, tName, util.FormatDuration(elapsed))) +
				"\n",
		)
	}

	// ── Row list ──────────────────────────────────────────────────────────
	lh := m.listHeight()
	end := m.offset + lh
	if end > len(m.rows) {
		end = len(m.rows)
	}

	var body strings.Builder

	// Empty state: render the message as the first body row, then pad.
	if len(m.rows) == 0 {
		body.WriteString(StyleDimmed.Render("No data yet — create a project and track some time."))
		body.WriteString("\n")
		for i := 1; i < lh; i++ {
			body.WriteString("\n")
		}
	} else {
		for i := m.offset; i < end; i++ {
			row := m.rows[i]
			switch row.kind {
			case reportRowBlank:
				body.WriteString("\n")

			case reportRowProject:
				p := m.store.FindProject(row.projectID)
				if p == nil {
					continue
				}
				secs := p.TotalSeconds()
				name := truncate(p.Name, nameW-2) // "▸ " prefix = 2
				nameCol := lipgloss.NewStyle().Width(nameW).Render(
					StyleProject.Render("▸ " + name),
				)
				barCol := renderBar(secs, grandTotal, barW)
				durCol := lipgloss.NewStyle().Width(durW).Align(lipgloss.Right).
					Render(StyleDuration.Render(util.FormatDuration(secs)))
				body.WriteString(nameCol + strings.Repeat(" ", gap) + barCol + " " + durCol)
				body.WriteString("\n")

			case reportRowTask:
				p := m.store.FindProject(row.projectID)
				t := m.store.FindTask(row.projectID, row.taskID)
				if p == nil || t == nil {
					continue
				}
				secs := t.TotalSeconds()
				name := truncate(t.Name, nameW-4) // "  · " prefix = 4
				nameCol := lipgloss.NewStyle().Width(nameW).Render(
					StyleDimmed.Render("  · ") + StyleTask.Render(name),
				)
				barCol := renderBar(secs, p.TotalSeconds(), barW)
				durCol := lipgloss.NewStyle().Width(durW).Align(lipgloss.Right).
					Render(StyleDuration.Render(util.FormatDuration(secs)))
				body.WriteString(nameCol + strings.Repeat(" ", gap) + barCol + " " + durCol)
				body.WriteString("\n")
			}
		}

		// Pad unused rows to stabilise layout.
		rendered := end - m.offset
		for i := rendered; i < lh; i++ {
			body.WriteString("\n")
		}
	}

	// Scroll indicator — always written as exactly one line so layout height
	// stays stable regardless of whether the list is scrollable.
	scrollHint := ""
	canScrollUp := m.offset > 0
	canScrollDown := m.offset+lh < len(m.rows)
	if canScrollUp || canScrollDown {
		parts := []string{}
		if canScrollUp {
			parts = append(parts, "↑ more above")
		}
		if canScrollDown {
			parts = append(parts, "↓ more below")
		}
		scrollHint = StyleDimmed.Render(strings.Join(parts, "   "))
	}

	// ── Help bar ──────────────────────────────────────────────────────────
	m.help.Styles = helpStyles()
	m.help.Width = innerW
	helpStr := m.help.View(reportKeyMap{m.keys})

	// ── Assemble ──────────────────────────────────────────────────────────
	var sb strings.Builder
	sb.WriteString(header.String())
	sb.WriteString(body.String())
	sb.WriteString(scrollHint + "\n") // always 1 line (blank when not scrollable)
	sb.WriteString("\n")
	sb.WriteString(helpStr)

	return StylePanel.Width(innerW).Padding(0, 1).Render(sb.String())
}

// ── Helpers ────────────────────────────────────────────────────────────────

// renderBar builds a fixed-width progress bar using block characters.
// The filled portion is proportional to secs/total.
func renderBar(secs, total int64, width int) string {
	filled := 0
	if total > 0 && secs > 0 {
		filled = int(int64(width) * secs / total)
		if filled > width {
			filled = width
		}
	}
	empty := width - filled
	bar := StyleTimer.Render(strings.Repeat("█", filled)) +
		StyleDimmed.Render(strings.Repeat("░", empty))
	return bar
}

// truncate clips s to at most max visible characters, appending "…" if cut.
// It counts Unicode code points (runes) as 1 column wide each — close enough
// for project/task names which are almost always Latin text.
func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}

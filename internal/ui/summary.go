package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

type summaryScope int

const (
	scopeToday summaryScope = iota
	scopeWeek
	scopeMonth
)

type summaryRowKind int

const (
	summaryKindProject summaryRowKind = iota
	summaryKindTask
	summaryKindSession
	summaryKindBlank
)

type summaryRow struct {
	kind     summaryRowKind
	label    string // project name, task name, or time range string
	duration int64  // seconds; 0 for blank rows
}

// SummaryModel shows sessions logged today (default) or this week, grouped by
// project and task. Opened with 's' from the tree view; 'w' toggles the scope.
type SummaryModel struct {
	store      *model.Store
	keys       KeyMap
	scope      summaryScope
	rows       []summaryRow
	grandTotal int64
	offset     int
	width      int
	height     int
	help       help.Model

	SwitchToTree bool
}

func NewSummaryModel(store *model.Store, keys KeyMap) SummaryModel {
	h := help.New()
	h.Styles = helpStyles()
	m := SummaryModel{
		store: store,
		keys:  keys,
		scope: scopeToday,
		help:  h,
	}
	m.buildRows()
	return m
}

// buildRows filters sessions by scope and groups them by project → task.
func (m *SummaryModel) buildRows() {
	now := time.Now()
	var windowStart time.Time
	switch m.scope {
	case scopeToday:
		windowStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case scopeWeek:
		// Start of ISO week (Monday).
		wd := int(now.Weekday())
		if wd == 0 {
			wd = 7
		}
		monday := now.AddDate(0, 0, -(wd - 1))
		windowStart = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
	case scopeMonth:
		windowStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	m.rows = nil
	m.grandTotal = 0

	for _, p := range m.store.Projects {
		var projSecs int64
		var taskRows []summaryRow

		for _, t := range p.Tasks {
			var taskSecs int64
			var sessRows []summaryRow

			for _, sess := range t.Sessions {
				if sess.Start.Before(windowStart) {
					continue
				}
				taskSecs += sess.DurationSeconds
				timeRange := sess.Start.Format("3:04 PM") + " — " + sess.End.Format("3:04 PM")
				sessRows = append(sessRows, summaryRow{
					kind:     summaryKindSession,
					label:    timeRange,
					duration: sess.DurationSeconds,
				})
			}

			if len(sessRows) == 0 {
				continue
			}

			projSecs += taskSecs
			taskRows = append(taskRows, summaryRow{
				kind:     summaryKindTask,
				label:    t.Name,
				duration: taskSecs,
			})
			taskRows = append(taskRows, sessRows...)
		}

		if len(taskRows) == 0 {
			continue
		}

		m.grandTotal += projSecs
		m.rows = append(m.rows, summaryRow{
			kind:     summaryKindProject,
			label:    p.Name,
			duration: projSecs,
		})
		m.rows = append(m.rows, taskRows...)
		m.rows = append(m.rows, summaryRow{kind: summaryKindBlank})
	}
}

func (m SummaryModel) Init() tea.Cmd { return nil }

func (m SummaryModel) Update(msg tea.Msg) (SummaryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width - 4

	case tickMsg:
		// Rebuild on tick so a just-saved timer session appears without
		// requiring the user to reopen the view.
		m.buildRows()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.SwitchToTree = true
		case "w":
			switch m.scope {
			case scopeToday:
				m.scope = scopeWeek
			case scopeWeek:
				m.scope = scopeMonth
			default:
				m.scope = scopeToday
			}
			m.offset = 0
			m.buildRows()
		case "up", "k":
			if m.offset > 0 {
				m.offset--
			}
		case "down", "j":
			max := len(m.rows) - m.visibleLines()
			if max < 0 {
				max = 0
			}
			if m.offset < max {
				m.offset++
			}
		}
	}
	return m, nil
}

func (m *SummaryModel) visibleLines() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	// 2 border + 1 title + 1 rule + 1 rule (footer sep) + 1 total + 1 spacer + 1 help = 8
	v := h - 10
	if v < 1 {
		v = 1
	}
	return v
}

func (m SummaryModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	innerW := w - 4
	if innerW < 20 {
		innerW = 20
	}

	var sb strings.Builder

	// ── Title ─────────────────────────────────────────────────────────────
	scopeLabel := "Today"
	switch m.scope {
	case scopeWeek:
		scopeLabel = "This Week"
	case scopeMonth:
		scopeLabel = "This Month"
	}
	title := StyleTitle.Render("◷  Summary — " + scopeLabel)
	toggleHint := StyleDimmed.Render("  [w] toggle")
	titleLine := title + toggleHint
	sb.WriteString(titleLine)
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render(strings.Repeat("─", innerW)))
	sb.WriteString("\n")

	// ── Session rows ──────────────────────────────────────────────────────
	vl := m.visibleLines()
	end := m.offset + vl
	if end > len(m.rows) {
		end = len(m.rows)
	}

	if len(m.rows) == 0 {
		emptyMsg := "No sessions logged "
		switch m.scope {
		case scopeToday:
			emptyMsg += "today."
		case scopeWeek:
			emptyMsg += "this week."
		case scopeMonth:
			emptyMsg += "this month."
		}
		sb.WriteString(StyleDimmed.Render(emptyMsg))
		sb.WriteString("\n")
		for i := 1; i < vl; i++ {
			sb.WriteString("\n")
		}
	} else {
		// Scroll hint arrows
		if m.offset > 0 {
			sb.WriteString(StyleDimmed.Render("  ↑ more"))
			sb.WriteString("\n")
		}

		for _, row := range m.rows[m.offset:end] {
			sb.WriteString(m.renderRow(row, innerW))
			sb.WriteString("\n")
		}

		shown := end - m.offset
		extra := 0
		if m.offset > 0 {
			extra = 1 // the ↑ more line
		}
		for i := shown + extra; i < vl; i++ {
			sb.WriteString("\n")
		}
	}

	// ── Footer: grand total ───────────────────────────────────────────────
	sb.WriteString(StyleDimmed.Render(strings.Repeat("─", innerW)))
	sb.WriteString("\n")
	totalLabel := StyleDimmed.Render("Total")
	totalDur := StyleDuration.Render(util.FormatDuration(m.grandTotal))
	labelW := lipgloss.Width(totalLabel)
	durW := lipgloss.Width(totalDur)
	space := innerW - labelW - durW
	if space < 1 {
		space = 1
	}
	sb.WriteString(totalLabel + strings.Repeat(" ", space) + totalDur)
	sb.WriteString("\n")

	// ── Help bar ──────────────────────────────────────────────────────────
	sb.WriteString("\n")
	m.help.Styles = helpStyles()
	m.help.Width = innerW
	sb.WriteString(m.help.View(summaryKeyMap{m.keys}))

	return StylePanel.
		Width(innerW + 2).
		Padding(0, 1).
		Render(sb.String())
}

func (m *SummaryModel) renderRow(row summaryRow, innerW int) string {
	switch row.kind {
	case summaryKindBlank:
		return ""

	case summaryKindProject:
		name := StyleProject.Render("▶ " + row.label)
		dur := StyleDuration.Render(util.FormatDuration(row.duration))
		nameW := lipgloss.Width(name)
		durW := lipgloss.Width(dur)
		space := innerW - nameW - durW
		if space < 1 {
			space = 1
		}
		return name + strings.Repeat(" ", space) + dur

	case summaryKindTask:
		name := StyleTask.Render("  · " + row.label)
		dur := StyleDimmed.Render(util.FormatDuration(row.duration))
		nameW := lipgloss.Width(name)
		durW := lipgloss.Width(dur)
		space := innerW - nameW - durW
		if space < 1 {
			space = 1
		}
		return name + strings.Repeat(" ", space) + dur

	case summaryKindSession:
		timeRange := lipgloss.NewStyle().Foreground(colorSubtext).Render("      " + row.label)
		dur := StyleDimmed.Render(util.FormatDuration(row.duration))
		rangeW := lipgloss.Width(timeRange)
		durW := lipgloss.Width(dur)
		space := innerW - rangeW - durW
		if space < 1 {
			space = 1
		}
		return timeRange + strings.Repeat(" ", space) + dur
	}
	return ""
}

// ── Key map ───────────────────────────────────────────────────────────────────

type summaryKeyMap struct{ km KeyMap }

func (k summaryKeyMap) ShortHelp() []key.Binding {
	toggleScope := key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "cycle today/week/month"))
	return []key.Binding{k.km.Up, k.km.Down, toggleScope, k.km.Esc}
}

func (k summaryKeyMap) FullHelp() [][]key.Binding {
	toggleScope := key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "cycle today/week/month"))
	return [][]key.Binding{
		{k.km.Up, k.km.Down},
		{toggleScope, k.km.Esc},
	}
}

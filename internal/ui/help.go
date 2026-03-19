package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpModel is a full-screen scrollable reference covering all keybindings
// and CLI commands. Opened with H from any view; dismissed with esc or q.
type HelpModel struct {
	keys   KeyMap
	lines  []string // pre-rendered content lines
	offset int      // first visible line
	width  int
	height int

	SwitchToTree bool
}

func NewHelpModel(keys KeyMap) HelpModel {
	m := HelpModel{keys: keys}
	m.buildLines()
	return m
}

func (m HelpModel) Init() tea.Cmd { return nil }

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.SwitchToTree = true
		case "up", "k":
			if m.offset > 0 {
				m.offset--
			}
		case "down", "j":
			max := len(m.lines) - m.visibleLines()
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

func (m *HelpModel) visibleLines() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	// 2 border + 2 padding rows + 1 footer line + 1 spacer = 6
	v := h - 6
	if v < 1 {
		v = 1
	}
	return v
}

func (m HelpModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	innerW := w - 4 // 2 border + 2 padding (Padding(0,1) * 2 sides)
	if innerW < 20 {
		innerW = 20
	}

	vl := m.visibleLines()
	end := m.offset + vl
	if end > len(m.lines) {
		end = len(m.lines)
	}

	var sb strings.Builder
	for _, ln := range m.lines[m.offset:end] {
		sb.WriteString(ln)
		sb.WriteString("\n")
	}
	// Pad remaining space so the footer stays at the bottom.
	shown := end - m.offset
	for i := shown; i < vl; i++ {
		sb.WriteString("\n")
	}

	// Footer: scroll hint + dismiss key.
	var footerParts []string
	if m.offset > 0 || end < len(m.lines) {
		hint := ""
		if m.offset > 0 {
			hint += "↑ "
		}
		if end < len(m.lines) {
			hint += "↓ "
		}
		footerParts = append(footerParts, StyleDimmed.Render(strings.TrimSpace(hint)+" scroll"))
	}
	footerParts = append(footerParts,
		StyleDimmed.Render("[q] / [esc] close"),
	)
	sb.WriteString(strings.Join(footerParts, "   "))

	return StylePanel.
		Width(innerW + 2).
		Padding(0, 1).
		Render(sb.String())
}

// buildLines pre-renders all help content into styled lines.
func (m *HelpModel) buildLines() {
	h1 := func(s string) string {
		return StyleTitle.Render(s)
	}
	h2 := func(s string) string {
		return StyleProject.Render(s)
	}
	dim := func(s string) string {
		return StyleDimmed.Render(s)
	}
	key := func(k string) string {
		return lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render(k)
	}
	row := func(keys, desc string) string {
		k := key(keys)
		kw := lipgloss.Width(k)
		pad := 18 - kw
		if pad < 1 {
			pad = 1
		}
		return k + strings.Repeat(" ", pad) + StyleTask.Render(desc)
	}
	sep := func() string {
		return dim(strings.Repeat("─", 44))
	}
	note := func(s string) string {
		return dim("  " + s)
	}

	var ls []string
	add := func(s string) { ls = append(ls, s) }
	blank := func() { ls = append(ls, "") }

	// ── Header ──────────────────────────────────────────────────────────────
	add(h1("◷  shellclock — Help Reference"))
	blank()
	add(dim("All views share the same data file. Changes are saved automatically."))
	blank()

	// ── Tree View ────────────────────────────────────────────────────────────
	add(sep())
	add(h2("Tree View  (main screen)"))
	add(sep())
	blank()
	add(h2("Navigation"))
	add(row("↑ / k", "Move cursor up"))
	add(row("↓ / j", "Move cursor down"))
	add(row("→ / l", "Expand selected project"))
	add(row("← / h", "Collapse selected project"))
	add(row("enter", "Expand/collapse project; open task detail"))
	blank()
	add(h2("Projects & Tasks"))
	add(row("N", "Create a new project"))
	add(row("n", "Create a new task under the focused project"))
	add(row("E", "Rename the focused project or task"))
	add(row("d", "Delete focused project or task (confirmation required)"))
	add(row("#", "Edit tags for the focused project"))
	add(note("After creating or renaming a project you'll be prompted for tags."))
	add(note("Tags are comma-separated: work, client, learning"))
	blank()
	add(h2("Timer"))
	add(row("p", "Start timer on focused task / pause / resume"))
	add(row("S", "Stop timer and save the session"))
	add(row("r", "Reset timer (discards accumulated time)"))
	add(note("Only one timer can run at a time."))
	add(note("The running timer survives app restarts."))
	blank()
	add(h2("Navigation & Other"))
	add(row("e", "Open session editor for focused task"))
	add(row("s", "Open summary view (today's sessions)"))
	add(row("R", "Open report view"))
	add(row("T", "Open theme picker"))
	add(row("$", "Set hourly rate for focused project (blank to clear)"))
	add(row("B", "Open backup picker — navigate and restore a backup"))
	add(row("W", "Idle warn settings — set threshold in minutes (0 = disable)"))
	add(row("X", "System reset — erase all data (requires typing CONFIRM)"))
	add(row("H", "Open this help screen"))
	add(row("?", "Toggle compact key reference bar"))
	add(row("q / ctrl+c", "Quit"))
	blank()

	// ── Summary View ─────────────────────────────────────────────────────────
	add(sep())
	add(h2("Summary View"))
	add(sep())
	add(note("Open with s from the tree view. Shows sessions logged today or this week."))
	blank()
	add(row("↑ / k", "Scroll up"))
	add(row("↓ / j", "Scroll down"))
	add(row("w", "Toggle between today and this week"))
	add(row("esc / q", "Return to tree"))
	blank()

	// ── Task Detail View ─────────────────────────────────────────────────────
	add(sep())
	add(h2("Task Detail View"))
	add(sep())
	add(note("Open with enter on any task in the tree."))
	blank()
	add(row("s / enter", "Start / pause / resume timer for this task"))
	add(row("S", "Stop timer and save session"))
	add(row("r", "Reset timer"))
	add(row("↑ / k", "Scroll session list up"))
	add(row("↓ / j", "Scroll session list down"))
	add(row("n", "Add a new session manually"))
	add(row("e", "Edit selected session"))
	add(row("d", "Delete selected session (confirmation required)"))
	add(row("esc / q", "Return to tree view"))
	blank()

	// ── Session Editor ───────────────────────────────────────────────────────
	add(sep())
	add(h2("Session Editor"))
	add(sep())
	add(note("Open with e on any task in the tree or task detail view."))
	blank()
	add(h2("List mode"))
	add(row("↑ / k", "Move up"))
	add(row("↓ / j", "Move down"))
	add(row("n", "Add a new session"))
	add(row("e", "Edit selected session"))
	add(row("d", "Delete selected session (confirmation required)"))
	add(row("esc / q", "Return to tree"))
	blank()
	add(h2("Form mode  (add / edit session)"))
	add(row("tab", "Move to next field  (start → end → notes)"))
	add(row("shift+tab", "Move to previous field"))
	add(row("enter", "Confirm field / save on the last field"))
	add(row("esc", "Cancel and discard changes"))
	add(note("Time format: YYYY-MM-DD HH:MM:SS"))
	add(note("Notes field is optional — max 120 characters."))
	blank()

	// ── Report View ──────────────────────────────────────────────────────────
	add(sep())
	add(h2("Report View"))
	add(sep())
	add(note("Open with R from the tree view."))
	blank()
	add(row("↑ / k", "Scroll up"))
	add(row("↓ / j", "Scroll down"))
	add(row("f", "Open tag filter (show only projects with a selected tag)"))
	add(row("$", "Toggle earnings column (shows $rate×hours per project)"))
	add(row("x", "Export report — choose CSV or plain text"))
	add(row("R / esc / q", "Return to tree"))
	blank()
	add(note("CSV export always includes Hourly Rate and Earnings columns when rates are set."))
	add(note("Reports are saved to ~/.config/shellclock/reports/"))
	blank()

	// ── Theme Picker ─────────────────────────────────────────────────────────
	add(sep())
	add(h2("Theme Picker"))
	add(sep())
	add(note("Open with T from the tree view. Themes preview live."))
	blank()
	add(row("↑ / k", "Previous theme (live preview)"))
	add(row("↓ / j", "Next theme (live preview)"))
	add(row("enter", "Apply and save selected theme"))
	add(row("esc / q", "Cancel — revert to previous theme"))
	add(note("31 built-in themes: Catppuccin, Dracula, Nord, Tokyo Night, Gruvbox,"))
	add(note("Rosé Pine, Kanagawa, Everforest, Ayu, Monokai Pro, and more."))
	blank()

	// ── Backups ──────────────────────────────────────────────────────────────
	add(sep())
	add(h2("Backups"))
	add(sep())
	blank()
	add(dim("shellclock automatically backs up your data on every launch."))
	blank()
	add(row("B", "Open backup picker (from tree view)"))
	blank()
	add(note("One backup per calendar day; last 7 are retained."))
	add(note("Stored at: ~/.config/shellclock/backups/"))
	blank()
	add(h2("Restore from backup"))
	add(row("↑ / k", "Select previous backup"))
	add(row("↓ / j", "Select next backup"))
	add(row("enter", "Prompt to restore selected backup"))
	add(row("y", "Confirm restore (overwrites current data)"))
	add(row("n / esc", "Cancel"))
	blank()

	// ── Idle Warn ────────────────────────────────────────────────────────────
	add(sep())
	add(h2("Idle Timer Warning"))
	add(sep())
	blank()
	add(dim("When a timer runs past the configured threshold, a flashing ⚠ Nh+"))
	add(dim("indicator appears next to the timer badge in the tree and task detail views."))
	add(dim("The warning only fires when the timer is running (not paused)."))
	blank()
	add(row("W", "Open idle-warn settings (from tree view)"))
	add(note("Enter a threshold in minutes.  Enter 0 to disable the warning."))
	add(note("Default: 120 minutes (2 hours).  Setting is saved to the data file."))
	blank()

	// ── CLI Commands ─────────────────────────────────────────────────────────
	add(sep())
	add(h2("CLI Commands  (run in your terminal, not inside the app)"))
	add(sep())
	blank()
	add(h2("Launch"))
	add(row("shellclock", "Open the TUI"))
	blank()
	add(h2("Import"))
	add(StyleTask.Render("shellclock import toggl /path/to/export.csv"))
	blank()
	add(dim("  Imports a Toggl Detailed CSV export into your shellclock data."))
	add(dim("  The TUI is not launched — the import runs and exits."))
	blank()
	add(dim("  How to export from Toggl:"))
	add(dim("    1. Toggl Track → Reports → Detailed"))
	add(dim("    2. Set your date range"))
	add(dim("    3. Export → Download as CSV"))
	blank()
	add(dim("  Column mapping:"))
	add(dim("    Toggl Project     → shellclock Project"))
	add(dim("    Toggl Task        → shellclock Task  (falls back to Description)"))
	add(dim("    Toggl Description → session notes  (when it differs from task name)"))
	add(dim("    Start date/time + End date/time → session start/end"))
	blank()
	add(dim("  Projects and tasks matched by name — existing ones are reused,"))
	add(dim("  sessions are always appended and never duplicated."))
	blank()
	add(dim("  Output on success:"))
	add(StyleTask.Render("    3 projects imported, 4 tasks imported, 20 sessions imported"))
	blank()

	// ── Data & Files ─────────────────────────────────────────────────────────
	add(sep())
	add(h2("Data & Files"))
	add(sep())
	blank()
	add(dim("  Data file:    ~/.config/shellclock/shellclock.json"))
	add(dim("  Backups:      ~/.config/shellclock/backups/"))
	add(dim("  Reports:      ~/.config/shellclock/reports/"))
	blank()
	add(dim("  All writes are atomic (tmp file + rename). The data file is"))
	add(dim("  saved on every change — no manual save needed."))
	blank()

	m.lines = ls
}


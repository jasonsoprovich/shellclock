package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

// TimerModel manages the live timer view.
type TimerModel struct {
	store  *model.Store
	keys   KeyMap
	width  int
	height int
	help   help.Model

	SwitchToTree bool
}

func NewTimerModel(store *model.Store, keys KeyMap) TimerModel {
	h := help.New()
	h.Styles = catppuccinHelpStyles()
	return TimerModel{store: store, keys: keys, help: h}
}

// Init starts the tick chain when a timer is already running — covers both the
// "user navigated here from tree" path and the "app restarted with a persisted
// timer" path.
func (m TimerModel) Init() tea.Cmd {
	if m.store.ActiveTimer != nil && !m.store.ActiveTimer.Paused {
		return tick()
	}
	return nil
}

// elapsed returns the total seconds tracked so far: time banked before the
// last pause plus the current running interval measured from Start.
func (m TimerModel) elapsed() int64 {
	at := m.store.ActiveTimer
	if at == nil {
		return 0
	}
	acc := at.AccumulatedSeconds
	if !at.Paused {
		acc += int64(time.Since(at.Start).Seconds())
	}
	return acc
}

func (m TimerModel) Update(msg tea.Msg) (TimerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width - 4
		return m, nil

	case tickMsg:
		// Keep the chain alive while the timer is running.
		if m.store.ActiveTimer != nil && !m.store.ActiveTimer.Paused {
			return m, tick()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			// Go back to the tree without touching the active timer — it keeps
			// running in the background and is persisted in the JSON file.
			m.SwitchToTree = true

		case "p":
			at := m.store.ActiveTimer
			if at == nil {
				break
			}
			if at.Paused {
				// Resume: bank is already correct; just restart the interval.
				at.Start = time.Now()
				at.Paused = false
				_ = m.store.Save()
				return m, tick()
			}
			// Pause: bank the seconds accumulated during this interval.
			at.AccumulatedSeconds += int64(time.Since(at.Start).Seconds())
			at.Paused = true
			_ = m.store.Save()

		case "S":
			// Stop: commit the session and clear the active timer.
			at := m.store.ActiveTimer
			if at == nil {
				break
			}
			now := time.Now()
			secs := m.elapsed()
			if secs > 0 {
				m.store.AddSession(at.ProjectID, at.TaskID, model.Session{
					ID:              uuid.NewString(),
					Start:           at.OriginalStart, // wall-clock start, not last resume
					End:             now,
					DurationSeconds: secs,
				})
			}
			m.store.ActiveTimer = nil
			_ = m.store.Save()
			m.SwitchToTree = true

		case "r":
			// Reset: discard all accumulated time and restart from now.
			at := m.store.ActiveTimer
			if at == nil {
				break
			}
			now := time.Now()
			at.AccumulatedSeconds = 0
			at.OriginalStart = now
			at.Start = now
			at.Paused = false
			_ = m.store.Save()
			return m, tick()
		}
	}

	return m, nil
}

func (m TimerModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	h := m.height
	if h == 0 {
		h = 24
	}
	innerW := w - 4
	if innerW < 24 {
		innerW = 24
	}

	at := m.store.ActiveTimer

	// ── No active timer ────────────────────────────────────────────────────
	if at == nil {
		body := StyleDimmed.Render("No active timer.") +
			"\n\n" +
			StyleDimmed.Render("Select a task in the tree and press enter.")
		return StylePanel.Width(innerW).Padding(0, 1).Render(body)
	}

	// ── Look up names ──────────────────────────────────────────────────────
	p := m.store.FindProject(at.ProjectID)
	t := m.store.FindTask(at.ProjectID, at.TaskID)

	projectName, taskName := "Unknown Project", "Unknown Task"
	if p != nil {
		projectName = p.Name
	}
	if t != nil {
		taskName = t.Name
	}
	// Clamp names so the breadcrumb never overflows innerW.
	half := (innerW - 7) / 2 // 7 = "  ›  " + margin
	if half < 6 {
		half = 6
	}
	projectName = truncate(projectName, half)
	taskName = truncate(taskName, half)

	elapsed := m.elapsed()

	// ── Compose centred blocks ─────────────────────────────────────────────
	centre := lipgloss.NewStyle().Width(innerW).Align(lipgloss.Center)

	sep := StyleDimmed.Render(strings.Repeat("─", innerW))

	var clockStyle lipgloss.Style
	if at.Paused {
		clockStyle = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
	} else {
		clockStyle = StyleTimer
	}
	clock := centre.Render(clockStyle.Render(util.FormatDurationShort(elapsed)))

	var badge string
	if at.Paused {
		badge = centre.Render(lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("⏸  PAUSED"))
	} else {
		badge = centre.Render(StyleTimer.Render("●  RUNNING"))
	}

	crumb := centre.Render(
		StyleProject.Render(projectName) +
			StyleDimmed.Render("  ›  ") +
			StyleTask.Render(taskName),
	)

	started := centre.Render(
		StyleDimmed.Render("started  ") +
			StyleTask.Render(at.OriginalStart.Format("2006-01-02  15:04:05")),
	)

	prevTotal := ""
	if t != nil && t.TotalSeconds() > 0 {
		prevTotal = centre.Render(
			StyleDimmed.Render("prev total  ") +
				StyleDuration.Render(util.FormatDuration(t.TotalSeconds())),
		)
	}

	// Help bar.
	m.help.Width = innerW
	helpStr := m.help.View(timerKeyMap{m.keys})

	// ── Vertical centering ─────────────────────────────────────────────────
	// Panel outer height = content lines + 2 (border; Padding(0,1) is horizontal only).
	// Fixed content lines:
	//   header : title(1) + blank(1) + crumb(1) + blank(1)  = 4
	//   block  : sep(1)+clock(1)+sep(1)+blank(1)+badge(1)+blank(1)+started(1)
	//            [+prevTotal(1)]                             = 7 or 8
	//   footer : blank(1) + help(1)                         = 2
	blockLines := 7
	if prevTotal != "" {
		blockLines++
	}
	fixed := 2 + 4 + blockLines + 2
	topPad := (h - fixed) / 2
	if topPad < 0 {
		topPad = 0
	}

	// ── Assemble ───────────────────────────────────────────────────────────
	var sb strings.Builder

	sb.WriteString(StyleTitle.Render("Timer"))
	sb.WriteString("\n\n")
	sb.WriteString(crumb)
	sb.WriteString("\n\n")

	for range topPad {
		sb.WriteString("\n")
	}

	sb.WriteString(sep + "\n")
	sb.WriteString(clock + "\n")
	sb.WriteString(sep + "\n\n")
	sb.WriteString(badge + "\n\n")
	sb.WriteString(started)
	if prevTotal != "" {
		sb.WriteString("\n" + prevTotal)
	}

	for range topPad {
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(helpStr)

	return StylePanel.Width(innerW).Padding(0, 1).Render(sb.String())
}

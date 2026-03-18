package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

// TaskDetailModel shows everything for one task: live timer controls and the
// full session list with inline add/edit/delete.
type TaskDetailModel struct {
	store     *model.Store
	keys      KeyMap
	ProjectID string
	TaskID    string

	// list navigation
	cursor int
	offset int
	width  int
	height int

	// help
	help     help.Model
	showFull bool

	// session form state (shared constants from edit.go)
	inputMode   editInputMode
	activeField editField
	startInput  textinput.Model
	endInput    textinput.Model
	errMsg      string
	editingID   string

	SwitchToTree bool
}

func NewTaskDetailModel(store *model.Store, keys KeyMap) TaskDetailModel {
	si := textinput.New()
	si.CharLimit = 19
	si.Placeholder = "2006-01-02 15:04:05"
	si.PlaceholderStyle = StyleDimmed
	si.PromptStyle = StyleInputLabel
	si.TextStyle = StyleTask

	ei := textinput.New()
	ei.CharLimit = 19
	ei.Placeholder = "2006-01-02 15:04:05"
	ei.PlaceholderStyle = StyleDimmed
	ei.PromptStyle = StyleInputLabel
	ei.TextStyle = StyleTask

	h := help.New()
	h.Styles = helpStyles()

	return TaskDetailModel{
		store:      store,
		keys:       keys,
		startInput: si,
		endInput:   ei,
		help:       h,
	}
}

func (m TaskDetailModel) sessions() []model.Session {
	t := m.store.FindTask(m.ProjectID, m.TaskID)
	if t == nil {
		return nil
	}
	return t.Sessions
}

func (m TaskDetailModel) isActiveTask() bool {
	return m.store.ActiveTimer != nil && m.store.ActiveTimer.TaskID == m.TaskID
}

func (m TaskDetailModel) elapsed() int64 {
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

// listHeight computes the number of session rows that fit in the panel.
//
// Fixed lines (always present):
//
//	title + rule              = 2
//	timer section (2 lines)   = 2
//	blank after timer         = 1
//	sessions header + rule    = 2
//	col header                = 1
//	scroll hint               = 1
//	blank + help              = 2
//	border                    = 2
//	total                     = 13
//
// Optional:
//
//	form: blank+label+start+end+err = +5
//	full help                       = +3
func (m *TaskDetailModel) listHeight() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	fixed := 13
	if m.inputMode != editModeNone {
		fixed += 5
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

func (m *TaskDetailModel) scrollToCursor() {
	lh := m.listHeight()
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+lh {
		m.offset = m.cursor - lh + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *TaskDetailModel) clampCursor() {
	n := len(m.sessions())
	if n == 0 {
		m.cursor = 0
		return
	}
	if m.cursor >= n {
		m.cursor = n - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *TaskDetailModel) commitForm() {
	startStr := strings.TrimSpace(m.startInput.Value())
	endStr := strings.TrimSpace(m.endInput.Value())

	start, err := time.ParseInLocation(timeFmt, startStr, time.Local)
	if err != nil {
		m.errMsg = "invalid start — use 2006-01-02 15:04:05"
		return
	}
	end, err := time.ParseInLocation(timeFmt, endStr, time.Local)
	if err != nil {
		m.errMsg = "invalid end — use 2006-01-02 15:04:05"
		return
	}
	if !end.After(start) {
		m.errMsg = "end must be after start"
		return
	}
	secs := int64(end.Sub(start).Seconds())
	switch m.inputMode {
	case editModeAdd:
		m.store.AddSession(m.ProjectID, m.TaskID, model.Session{
			ID:              uuid.NewString(),
			Start:           start,
			End:             end,
			DurationSeconds: secs,
		})
	case editModeEdit:
		m.store.UpdateSession(m.ProjectID, m.TaskID, model.Session{
			ID:              m.editingID,
			Start:           start,
			End:             end,
			DurationSeconds: secs,
		})
	}
	_ = m.store.Save()
	m.closeForm()
}

func (m *TaskDetailModel) closeForm() {
	m.inputMode = editModeNone
	m.errMsg = ""
	m.editingID = ""
	m.startInput.Blur()
	m.endInput.Blur()
	m.startInput.SetValue("")
	m.endInput.SetValue("")
	m.clampCursor()
	m.scrollToCursor()
}

func (m *TaskDetailModel) openAdd() tea.Cmd {
	now := time.Now().Truncate(time.Second)
	m.inputMode = editModeAdd
	m.errMsg = ""
	m.activeField = fieldStart
	m.startInput.SetValue(now.Format(timeFmt))
	m.endInput.SetValue(now.Format(timeFmt))
	m.startInput.CursorEnd()
	return m.startInput.Focus()
}

func (m *TaskDetailModel) openEdit(sess model.Session) tea.Cmd {
	m.inputMode = editModeEdit
	m.errMsg = ""
	m.editingID = sess.ID
	m.activeField = fieldStart
	m.startInput.SetValue(sess.Start.Format(timeFmt))
	endVal := ""
	if !sess.End.IsZero() {
		endVal = sess.End.Format(timeFmt)
	}
	m.endInput.SetValue(endVal)
	m.startInput.CursorEnd()
	return m.startInput.Focus()
}

// startTimer creates an active timer for this task and returns a tick cmd to
// start the one-second chain.
func (m *TaskDetailModel) startTimer() tea.Cmd {
	if m.store.ActiveTimer != nil {
		return nil
	}
	now := time.Now()
	m.store.ActiveTimer = &model.ActiveTimer{
		ProjectID:     m.ProjectID,
		TaskID:        m.TaskID,
		OriginalStart: now,
		Start:         now,
	}
	_ = m.store.Save()
	return tick()
}

// ── Update ───────────────────────────────────────────────────────────────────

func (m TaskDetailModel) Update(msg tea.Msg) (TaskDetailModel, tea.Cmd) {
	var cmd tea.Cmd

	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		m.help.Width = ws.Width - 4
		return m, nil
	}

	key, isKey := msg.(tea.KeyMsg)

	// ── Form mode ────────────────────────────────────────────────────────────
	if m.inputMode != editModeNone && isKey {
		switch key.Type {
		case tea.KeyEscape:
			m.closeForm()
			return m, nil

		case tea.KeyTab, tea.KeyShiftTab:
			if m.activeField == fieldStart {
				m.activeField = fieldEnd
				m.startInput.Blur()
				cmd = m.endInput.Focus()
			} else {
				m.activeField = fieldStart
				m.endInput.Blur()
				cmd = m.startInput.Focus()
			}
			return m, cmd

		case tea.KeyEnter:
			if m.activeField == fieldStart {
				m.activeField = fieldEnd
				m.startInput.Blur()
				cmd = m.endInput.Focus()
				return m, cmd
			}
			m.commitForm()
			return m, nil
		}

		if m.activeField == fieldStart {
			m.startInput, cmd = m.startInput.Update(msg)
		} else {
			m.endInput, cmd = m.endInput.Update(msg)
		}
		m.errMsg = ""
		return m, cmd
	}

	// ── Normal mode ──────────────────────────────────────────────────────────
	if isKey {
		sessions := m.sessions()
		switch key.String() {
		case "esc", "q":
			m.SwitchToTree = true

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.scrollToCursor()
			}

		case "down", "j":
			if m.cursor < len(sessions)-1 {
				m.cursor++
				m.scrollToCursor()
			}

		case "s", "enter":
			at := m.store.ActiveTimer
			if at == nil {
				// Start.
				cmd = m.startTimer()
			} else if at.TaskID == m.TaskID {
				if at.Paused {
					// Resume.
					at.Start = time.Now()
					at.Paused = false
					_ = m.store.Save()
					cmd = tick()
				} else {
					// Pause.
					at.AccumulatedSeconds += int64(time.Since(at.Start).Seconds())
					at.Paused = true
					_ = m.store.Save()
				}
			}

		case "S":
			at := m.store.ActiveTimer
			if at != nil && at.TaskID == m.TaskID {
				now := time.Now()
				secs := m.elapsed()
				if secs > 0 {
					m.store.AddSession(at.ProjectID, at.TaskID, model.Session{
						ID:              uuid.NewString(),
						Start:           at.OriginalStart,
						End:             now,
						DurationSeconds: secs,
					})
				}
				m.store.ActiveTimer = nil
				_ = m.store.Save()
			}

		case "r":
			at := m.store.ActiveTimer
			if at != nil && at.TaskID == m.TaskID {
				now := time.Now()
				at.AccumulatedSeconds = 0
				at.OriginalStart = now
				at.Start = now
				at.Paused = false
				_ = m.store.Save()
				cmd = tick()
			}

		case "n":
			cmd = m.openAdd()

		case "e":
			if len(sessions) > 0 && m.cursor < len(sessions) {
				cmd = m.openEdit(sessions[m.cursor])
			}

		case "d":
			if len(sessions) > 0 && m.cursor < len(sessions) {
				m.store.DeleteSession(m.ProjectID, m.TaskID, sessions[m.cursor].ID)
				_ = m.store.Save()
				m.clampCursor()
				m.scrollToCursor()
			}

		case "?":
			m.showFull = !m.showFull
			m.help.ShowAll = m.showFull
			m.scrollToCursor()
		}
	}

	return m, cmd
}

// ── View ─────────────────────────────────────────────────────────────────────

func (m TaskDetailModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	innerW := w - 4
	if innerW < 58 {
		innerW = 58
	}

	p := m.store.FindProject(m.ProjectID)
	t := m.store.FindTask(m.ProjectID, m.TaskID)

	projectName, taskName := "Unknown", "Unknown"
	if p != nil {
		projectName = truncate(p.Name, 30)
	}
	if t != nil {
		taskName = truncate(t.Name, 40)
	}

	var sb strings.Builder
	centre := lipgloss.NewStyle().Width(innerW).Align(lipgloss.Center)

	// ── Header ───────────────────────────────────────────────────────────────
	sb.WriteString(
		StyleProject.Render(projectName) +
			StyleDimmed.Render("  ›  ") +
			StyleTask.Render(taskName),
	)
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render(strings.Repeat("─", innerW)))
	sb.WriteString("\n")

	// ── Timer section (always exactly 2 lines) ────────────────────────────────
	at := m.store.ActiveTimer
	isActive := at != nil && at.TaskID == m.TaskID

	if isActive {
		elapsed := m.elapsed()
		if at.Paused {
			pauseStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
			sb.WriteString(centre.Render(pauseStyle.Render("⏸  PAUSED  " + util.FormatDurationShort(elapsed))))
		} else {
			sb.WriteString(centre.Render(StyleTimer.Render("●  RUNNING  " + util.FormatDurationShort(elapsed))))
		}
		sb.WriteString("\n")
		sb.WriteString(centre.Render(StyleDimmed.Render("started  ") + StyleTask.Render(at.OriginalStart.Format("2006-01-02  15:04:05"))))
		sb.WriteString("\n")
	} else if at != nil {
		// Different task has an active timer.
		otherP := m.store.FindProject(at.ProjectID)
		otherT := m.store.FindTask(at.ProjectID, at.TaskID)
		pn, tn := "?", "?"
		if otherP != nil {
			pn = truncate(otherP.Name, 25)
		}
		if otherT != nil {
			tn = truncate(otherT.Name, 25)
		}
		sb.WriteString(StyleDimmed.Render("timer running on: ") +
			StyleProject.Render(pn) + StyleDimmed.Render(" › ") + StyleTask.Render(tn))
		sb.WriteString("\n")
		sb.WriteString("\n") // second line blank for layout stability
	} else {
		sb.WriteString(StyleDimmed.Render("no timer running  —  press s or enter to start"))
		sb.WriteString("\n")
		sb.WriteString("\n") // second line blank for layout stability
	}

	// blank after timer section
	sb.WriteString("\n")

	// ── Sessions section ──────────────────────────────────────────────────────
	sessions := m.sessions()
	var total int64
	for _, s := range sessions {
		total += s.DurationSeconds
	}

	sectionLeft := StyleTitle.Render("Sessions")
	sectionRight := StyleDimmed.Render("total  ") + StyleDuration.Render(util.FormatDuration(total))
	pad := innerW - lipgloss.Width(sectionLeft) - lipgloss.Width(sectionRight)
	if pad < 1 {
		pad = 1
	}
	sb.WriteString(sectionLeft + strings.Repeat(" ", pad) + sectionRight)
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render(strings.Repeat("─", innerW)))
	sb.WriteString("\n")

	// Column header.
	hdrIdx := lipgloss.NewStyle().Width(colIdx).Render(StyleDimmed.Render("#"))
	hdrStart := lipgloss.NewStyle().Width(colStart).Render(StyleDimmed.Render("Start"))
	hdrEnd := lipgloss.NewStyle().Width(colEnd).Render(StyleDimmed.Render("End"))
	hdrDur := lipgloss.NewStyle().Width(colDur).Align(lipgloss.Right).Render(StyleDimmed.Render("Duration"))
	sb.WriteString(hdrIdx + hdrStart + hdrEnd + hdrDur + "\n")

	// Session rows.
	lh := m.listHeight()
	end := m.offset + lh
	if end > len(sessions) {
		end = len(sessions)
	}

	if len(sessions) == 0 {
		sb.WriteString(StyleDimmed.Render("No sessions yet."))
		sb.WriteString("\n")
		for i := 1; i < lh; i++ {
			sb.WriteString("\n")
		}
	} else {
		for i := m.offset; i < end; i++ {
			sess := sessions[i]
			selected := i == m.cursor

			numStr := fmt.Sprintf("%d", i+1)
			startStr := sess.Start.Format(timeFmt)
			endStr := ""
			if !sess.End.IsZero() {
				endStr = sess.End.Format(timeFmt)
			}
			durStr := util.FormatDuration(sess.DurationSeconds)

			if selected {
				row := fmt.Sprintf("%-*s%-*s%-*s%*s",
					colIdx, numStr,
					colStart, startStr,
					colEnd, endStr,
					colDur, durStr,
				)
				sb.WriteString(highlightRow(row, innerW))
			} else {
				idxC := lipgloss.NewStyle().Width(colIdx).Render(StyleDimmed.Render(numStr))
				startC := lipgloss.NewStyle().Width(colStart).Render(StyleTask.Render(startStr))
				endC := lipgloss.NewStyle().Width(colEnd).Render(StyleTask.Render(endStr))
				durC := lipgloss.NewStyle().Width(colDur).Align(lipgloss.Right).Render(StyleDuration.Render(durStr))
				sb.WriteString(idxC + startC + endC + durC)
			}
			sb.WriteString("\n")
		}
		rendered := end - m.offset
		for i := rendered; i < lh; i++ {
			sb.WriteString("\n")
		}
	}

	// Scroll hint — always 1 line.
	scrollHint := ""
	canUp := m.offset > 0
	canDown := m.offset+lh < len(sessions)
	if canUp || canDown {
		parts := []string{}
		if canUp {
			parts = append(parts, "↑ more above")
		}
		if canDown {
			parts = append(parts, "↓ more below")
		}
		scrollHint = StyleDimmed.Render(strings.Join(parts, "   "))
	}
	sb.WriteString(scrollHint + "\n")

	// ── Input form (5 lines when active) ─────────────────────────────────────
	if m.inputMode != editModeNone {
		sb.WriteString("\n")
		label := "Add session"
		if m.inputMode == editModeEdit {
			label = "Edit session"
		}
		sb.WriteString(StyleInputLabel.Render(label))
		sb.WriteString("\n")

		var startPrompt, endPrompt string
		if m.activeField == fieldStart {
			startPrompt = StyleTimer.Render("▸") + StyleDimmed.Render(" start: ")
			endPrompt = StyleDimmed.Render("  end:   ")
		} else {
			startPrompt = StyleDimmed.Render("  start: ")
			endPrompt = StyleTimer.Render("▸") + StyleDimmed.Render(" end:   ")
		}
		sb.WriteString(startPrompt + m.startInput.View() + "\n")
		sb.WriteString(endPrompt + m.endInput.View() + "\n")

		errLine := ""
		if m.errMsg != "" {
			errLine = StyleError.Render("  " + m.errMsg)
		}
		sb.WriteString(errLine + "\n")
	}

	// ── Help bar ──────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	m.help.Styles = helpStyles()
	m.help.Width = innerW
	var km help.KeyMap
	if m.inputMode != editModeNone {
		km = editFormKeyMap{m.keys}
	} else if isActive {
		km = taskDetailActiveKeyMap{m.keys}
	} else {
		km = taskDetailKeyMap{m.keys}
	}
	sb.WriteString(m.help.View(km))

	return StylePanel.Width(innerW + 2).Padding(0, 1).Render(sb.String())
}

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

const timeFmt = "2006-01-02 15:04"

// editInputMode controls what the two text inputs are collecting.
type editInputMode int

const (
	editModeNone   editInputMode = iota
	editModeAdd                  // entering times for a brand-new session
	editModeEdit                 // editing an existing session's times
)

// editField tracks which of the two inputs is focused.
type editField int

const (
	fieldStart editField = iota
	fieldEnd
)

// EditModel allows viewing, adding, editing, and deleting sessions for a task.
type EditModel struct {
	store     *model.Store
	keys      KeyMap
	ProjectID string
	TaskID    string
	cursor    int
	offset    int
	width     int
	height    int
	help      help.Model
	showFull  bool

	// input state
	inputMode    editInputMode
	activeField  editField
	startInput   textinput.Model
	endInput     textinput.Model
	errMsg       string
	editingID    string // session ID being edited (editModeEdit only)

	SwitchToTree bool
}

func NewEditModel(store *model.Store, keys KeyMap) EditModel {
	si := textinput.New()
	si.CharLimit = 16
	si.Placeholder = "2006-01-02 15:04"
	si.PlaceholderStyle = StyleDimmed
	si.PromptStyle = StyleInputLabel
	si.TextStyle = StyleTask

	ei := textinput.New()
	ei.CharLimit = 16
	ei.Placeholder = "2006-01-02 15:04"
	ei.PlaceholderStyle = StyleDimmed
	ei.PromptStyle = StyleInputLabel
	ei.TextStyle = StyleTask

	h := help.New()
	h.Styles = catppuccinHelpStyles()

	return EditModel{
		store:      store,
		keys:       keys,
		startInput: si,
		endInput:   ei,
		help:       h,
	}
}

func (m EditModel) sessions() []model.Session {
	t := m.store.FindTask(m.ProjectID, m.TaskID)
	if t == nil {
		return nil
	}
	return t.Sessions
}

// listHeight returns how many session rows fit in the view.
// Fixed overhead: border(2) + header(4) + total(2) + footer(2) = 10.
// In input mode add 4 more rows for the form.
func (m *EditModel) listHeight() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	fixed := 10
	if m.inputMode != editModeNone {
		fixed += 4
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

func (m *EditModel) scrollToCursor() {
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

func (m *EditModel) clampCursor() {
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

// commitForm validates the two inputs and writes the session to the store.
func (m *EditModel) commitForm() {
	startStr := strings.TrimSpace(m.startInput.Value())
	endStr := strings.TrimSpace(m.endInput.Value())

	start, err := time.ParseInLocation(timeFmt, startStr, time.Local)
	if err != nil {
		m.errMsg = "invalid start time — use 2006-01-02 15:04"
		return
	}
	end, err := time.ParseInLocation(timeFmt, endStr, time.Local)
	if err != nil {
		m.errMsg = "invalid end time — use 2006-01-02 15:04"
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

func (m *EditModel) closeForm() {
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

func (m *EditModel) openAdd() {
	m.inputMode = editModeAdd
	m.errMsg = ""
	m.activeField = fieldStart
	m.startInput.SetValue("")
	m.endInput.SetValue("")
	_ = m.startInput.Focus()
}

func (m *EditModel) openEdit(sess model.Session) {
	m.inputMode = editModeEdit
	m.errMsg = ""
	m.editingID = sess.ID
	m.activeField = fieldStart
	m.startInput.SetValue(sess.Start.Format(timeFmt))
	m.endInput.SetValue(sess.End.Format(timeFmt))
	_ = m.startInput.Focus()
}

// ── Update ──────────────────────────────────────────────────────────────────

func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	var cmd tea.Cmd

	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		m.help.Width = ws.Width - 4
		return m, nil
	}

	key, isKey := msg.(tea.KeyMsg)

	// ── Form input mode ───────────────────────────────────────────────────
	if m.inputMode != editModeNone && isKey {
		switch key.Type {
		case tea.KeyEscape:
			m.closeForm()
			return m, nil

		case tea.KeyTab, tea.KeyShiftTab:
			// Toggle between start and end fields.
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
				// Advance to the end field on Enter in start field.
				m.activeField = fieldEnd
				m.startInput.Blur()
				cmd = m.endInput.Focus()
				return m, cmd
			}
			// Enter in end field commits the form.
			m.commitForm()
			return m, nil
		}

		// Route keystrokes to the focused input.
		if m.activeField == fieldStart {
			m.startInput, cmd = m.startInput.Update(msg)
		} else {
			m.endInput, cmd = m.endInput.Update(msg)
		}
		m.errMsg = "" // clear error on any edit
		return m, cmd
	}

	// ── Normal navigation mode ────────────────────────────────────────────
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

		case "n":
			m.openAdd()

		case "e":
			if len(sessions) > 0 && m.cursor < len(sessions) {
				m.openEdit(sessions[m.cursor])
			}

		case "d":
			if len(sessions) > 0 && m.cursor < len(sessions) {
				sess := sessions[m.cursor]
				m.store.DeleteSession(m.ProjectID, m.TaskID, sess.ID)
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

	return m, nil
}

// ── View ─────────────────────────────────────────────────────────────────────

func (m EditModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	innerW := w - 4
	if innerW < 30 {
		innerW = 30
	}

	// ── Header ─────────────────────────────────────────────────────────────
	t := m.store.FindTask(m.ProjectID, m.TaskID)
	taskName := "?"
	if t != nil {
		taskName = truncate(t.Name, 40)
	}

	var sb strings.Builder
	sb.WriteString(StyleTitle.Render("Edit Sessions"))
	sb.WriteString(StyleDimmed.Render("  —  "))
	sb.WriteString(StyleTask.Render(taskName))
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render(strings.Repeat("─", innerW)))
	sb.WriteString("\n\n")

	// ── Session list ────────────────────────────────────────────────────────
	sessions := m.sessions()
	lh := m.listHeight()
	end := m.offset + lh
	if end > len(sessions) {
		end = len(sessions)
	}

	// Column widths: idx(3) + gap(2) + start(16) + gap(2) + end(16) + gap(2) + dur(11)
	const startW, endW, durW = 16, 16, 11

	if len(sessions) == 0 {
		sb.WriteString(StyleDimmed.Render("No sessions — press n to add one."))
		sb.WriteString("\n")
		for i := 1; i < lh; i++ {
			sb.WriteString("\n")
		}
	} else {
		// Column header.
		idxCol := lipgloss.NewStyle().Width(3).Render(StyleDimmed.Render("#"))
		startCol := lipgloss.NewStyle().Width(startW + 2).Render(StyleDimmed.Render("Start"))
		endCol := lipgloss.NewStyle().Width(endW + 2).Render(StyleDimmed.Render("End"))
		durCol := lipgloss.NewStyle().Width(durW).Align(lipgloss.Right).Render(StyleDimmed.Render("Duration"))
		sb.WriteString(idxCol + startCol + endCol + durCol + "\n")

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
				row := fmt.Sprintf("%-3s  %-16s  %-16s  %11s", numStr, startStr, endStr, durStr)
				sb.WriteString(highlightRow(row, innerW))
			} else {
				idxC := lipgloss.NewStyle().Width(3).Render(StyleDimmed.Render(numStr))
				startC := lipgloss.NewStyle().Width(startW + 2).Render(StyleTask.Render(startStr))
				endC := lipgloss.NewStyle().Width(endW + 2).Render(StyleTask.Render(endStr))
				durC := lipgloss.NewStyle().Width(durW).Align(lipgloss.Right).Render(StyleDuration.Render(durStr))
				sb.WriteString(idxC + startC + endC + durC)
			}
			sb.WriteString("\n")
		}

		// Pad unused rows.
		rendered := end - m.offset
		for i := rendered; i < lh; i++ {
			sb.WriteString("\n")
		}
	}

	// Scroll hint.
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
		sb.WriteString(StyleDimmed.Render(strings.Join(parts, "   ")))
		sb.WriteString("\n")
	}

	// ── Total ───────────────────────────────────────────────────────────────
	var total int64
	for _, s := range sessions {
		total += s.DurationSeconds
	}
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render("total  ") + StyleDuration.Render(util.FormatDuration(total)))
	sb.WriteString("\n")

	// ── Input form ──────────────────────────────────────────────────────────
	if m.inputMode != editModeNone {
		sb.WriteString("\n")
		label := "Add session"
		if m.inputMode == editModeEdit {
			label = "Edit session"
		}
		sb.WriteString(StyleInputLabel.Render(label))
		sb.WriteString("\n")

		startPrompt := "  start: "
		endPrompt := "    end: "
		if m.activeField == fieldStart {
			startPrompt = StyleTimer.Render("▸ start: ")
			endPrompt = StyleDimmed.Render("  end:   ")
		} else {
			startPrompt = StyleDimmed.Render("  start: ")
			endPrompt = StyleTimer.Render("▸ end:   ")
		}
		sb.WriteString(startPrompt + m.startInput.View() + "\n")
		sb.WriteString(endPrompt + m.endInput.View() + "\n")

		if m.errMsg != "" {
			sb.WriteString(StyleError.Render("  " + m.errMsg) + "\n")
		}
	}

	// ── Help bar ─────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	m.help.Width = innerW
	var km help.KeyMap
	if m.inputMode != editModeNone {
		km = editFormKeyMap{m.keys}
	} else {
		km = editKeyMap{m.keys}
	}
	sb.WriteString(m.help.View(km))

	return StylePanel.Width(innerW).Padding(0, 1).Render(sb.String())
}

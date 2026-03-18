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

const timeFmt = "2006-01-02 15:04:05"

// editInputMode controls what the two text inputs are collecting.
type editInputMode int

const (
	editModeNone editInputMode = iota
	editModeAdd                // entering times for a brand-new session
	editModeEdit               // editing an existing session's times
)

// editField tracks which of the two inputs is focused.
type editField int

const (
	fieldStart editField = iota
	fieldEnd
)

// Session list column visible widths.
// Layout: [colIdx][colStart][colEnd][colDur] = 3+21+21+11 = 56
// The start and end columns include a 2-char trailing gap so adjacent
// columns read as separated without an explicit separator character.
const (
	colIdx   = 3
	colStart = 21 // 19 content + 2 trailing gap
	colEnd   = 21 // 19 content + 2 trailing gap
	colDur   = 11 // right-aligned duration
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
	inputMode   editInputMode
	activeField editField
	startInput  textinput.Model
	endInput    textinput.Model
	errMsg      string
	editingID   string // session ID being edited (editModeEdit only)

	// confirm-delete modal
	confirmActive    bool
	confirmMsg       string
	confirmSessionID string

	SwitchToTree bool
}

func NewEditModel(store *model.Store, keys KeyMap) EditModel {
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
//
// Fixed content lines (always present, no form):
//
//	title + rule + blank                   = 3
//	column header                          = 1
//	scroll hint (always reserved)          = 1
//	blank + total                          = 2
//	blank + help                           = 2
//	border                                 = 2
//	total                                  = 11
//
// Optional lines added to fixed:
//
//	form: blank + label + start + end + err = +5
//	full help vs short help                 = +3
func (m *EditModel) listHeight() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	fixed := 11
	if m.inputMode != editModeNone {
		fixed += 5 // blank + label + start + end + error
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

// openAdd starts the add-session form pre-filled with the current time so the
// user can navigate with arrow keys and adjust rather than typing from scratch.
func (m *EditModel) openAdd() tea.Cmd {
	now := time.Now().Truncate(time.Second)
	m.inputMode = editModeAdd
	m.errMsg = ""
	m.activeField = fieldStart
	m.startInput.SetValue(now.Format(timeFmt))
	m.endInput.SetValue(now.Format(timeFmt))
	m.startInput.CursorEnd()
	return m.startInput.Focus()
}

// openEdit pre-fills the form with sess's existing times and returns the
// Focus tea.Cmd.
func (m *EditModel) openEdit(sess model.Session) tea.Cmd {
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
				// Advance to end field on Enter in start.
				m.activeField = fieldEnd
				m.startInput.Blur()
				cmd = m.endInput.Focus()
				return m, cmd
			}
			// Enter in end field commits.
			m.commitForm()
			return m, nil
		}

		// Route keystrokes to the focused input; clear any stale error.
		if m.activeField == fieldStart {
			m.startInput, cmd = m.startInput.Update(msg)
		} else {
			m.endInput, cmd = m.endInput.Update(msg)
		}
		m.errMsg = ""
		return m, cmd
	}

	// ── Confirm-delete modal ──────────────────────────────────────────────
	if m.confirmActive && isKey {
		switch key.String() {
		case "y":
			m.store.DeleteSession(m.ProjectID, m.TaskID, m.confirmSessionID)
			_ = m.store.Save()
			m.clampCursor()
			m.scrollToCursor()
			m.confirmActive = false
			m.confirmMsg = ""
			m.confirmSessionID = ""
		case "n", "esc":
			m.confirmActive = false
			m.confirmMsg = ""
			m.confirmSessionID = ""
		}
		return m, nil
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
			cmd = m.openAdd()

		case "e":
			if len(sessions) > 0 && m.cursor < len(sessions) {
				cmd = m.openEdit(sessions[m.cursor])
			}

		case "d":
			if len(sessions) > 0 && m.cursor < len(sessions) {
				sess := sessions[m.cursor]
				m.confirmMsg = fmt.Sprintf("Delete session %d (%s)?", m.cursor+1, sess.Start.Format("2006-01-02 15:04"))
				m.confirmSessionID = sess.ID
				m.confirmActive = true
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

func (m EditModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	innerW := w - 4
	if innerW < 58 {
		innerW = 58
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

	// ── Column header — always shown for layout stability ──────────────────
	hdrIdx := lipgloss.NewStyle().Width(colIdx).Render(StyleDimmed.Render("#"))
	hdrStart := lipgloss.NewStyle().Width(colStart).Render(StyleDimmed.Render("Start"))
	hdrEnd := lipgloss.NewStyle().Width(colEnd).Render(StyleDimmed.Render("End"))
	hdrDur := lipgloss.NewStyle().Width(colDur).Align(lipgloss.Right).Render(StyleDimmed.Render("Duration"))
	sb.WriteString(hdrIdx + hdrStart + hdrEnd + hdrDur + "\n")

	// ── Session list ────────────────────────────────────────────────────────
	sessions := m.sessions()
	lh := m.listHeight()
	end := m.offset + lh
	if end > len(sessions) {
		end = len(sessions)
	}

	if len(sessions) == 0 {
		sb.WriteString(StyleDimmed.Render("No sessions — press n to add one."))
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
				// Use the same column widths as non-selected rows so the
				// highlighted row never shifts the layout.
				// colIdx=3  colStart=18  colEnd=18  colDur=11 → total=50
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
				durC := lipgloss.NewStyle().Width(colDur).Align(lipgloss.Right).
					Render(StyleDuration.Render(durStr))
				sb.WriteString(idxC + startC + endC + durC)
			}
			sb.WriteString("\n")
		}

		// Pad unused rows to stabilise layout height.
		rendered := end - m.offset
		for i := rendered; i < lh; i++ {
			sb.WriteString("\n")
		}
	}

	// Scroll hint — always written as exactly one line so layout height stays
	// stable regardless of whether the list is scrollable.
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

	// ── Total ───────────────────────────────────────────────────────────────
	var total int64
	for _, s := range sessions {
		total += s.DurationSeconds
	}
	sb.WriteString("\n")
	sb.WriteString(StyleDimmed.Render("total  ") + StyleDuration.Render(util.FormatDuration(total)))
	sb.WriteString("\n")

	// ── Input form ──────────────────────────────────────────────────────────
	// The form occupies exactly 5 lines (blank + label + start + end + error)
	// so listHeight() can add a fixed +5 when inputMode != None.
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

		// Error line is always written (blank when no error) to keep layout stable.
		errLine := ""
		if m.errMsg != "" {
			errLine = StyleError.Render("  " + m.errMsg)
		}
		sb.WriteString(errLine + "\n")
	}

	// ── Help bar ─────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	m.help.Styles = helpStyles()
	m.help.Width = innerW
	var km help.KeyMap
	if m.inputMode != editModeNone {
		km = editFormKeyMap{m.keys}
	} else {
		km = editKeyMap{m.keys}
	}
	sb.WriteString(m.help.View(km))

	panel := StylePanel.Width(innerW + 2).Padding(0, 1).Render(sb.String())
	if m.confirmActive {
		return renderConfirmOverlay(panel, m.confirmMsg, m.width, m.height)
	}
	return panel
}

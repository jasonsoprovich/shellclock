package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

type inputMode int

const (
	inputNone    inputMode = iota
	inputProject           // naming a new project
	inputTask              // naming a new task
)

// treeItem is one visible row in the flattened tree.
type treeItem struct {
	isProject bool
	projectID string
	taskID    string
	name      string
	expanded  bool // only meaningful for projects
}

// TreeModel manages the project/task tree view.
type TreeModel struct {
	store    *model.Store
	keys     KeyMap
	items    []treeItem
	cursor   int
	offset   int // index of the first visible row
	width    int
	height   int
	expanded map[string]bool // set of expanded project IDs

	// text-input state
	mode      inputMode
	textInput textinput.Model
	parentID  string // project ID when mode == inputTask

	// help component
	help     help.Model
	showFull bool

	// signals consumed by App
	WantsQuit           bool
	SwitchToTimer       bool
	SwitchToEdit        bool
	SwitchToReport      bool
	SwitchToThemePicker bool
	SelectedProjectID   string
	SelectedTaskID      string
}

func NewTreeModel(store *model.Store, keys KeyMap) TreeModel {
	ti := textinput.New()
	ti.CharLimit = 64
	ti.Prompt = "> "
	ti.PromptStyle = StyleInputLabel
	ti.TextStyle = StyleTask
	ti.PlaceholderStyle = StyleDimmed

	h := help.New()
	h.Styles = helpStyles()

	m := TreeModel{
		store:     store,
		keys:      keys,
		expanded:  make(map[string]bool),
		textInput: ti,
		help:      h,
	}
	for _, p := range store.Projects {
		m.expanded[p.ID] = true
	}
	m.buildItems()
	return m
}

// buildItems flattens the project/task tree into a navigable list.
func (m *TreeModel) buildItems() {
	m.items = nil
	for _, p := range m.store.Projects {
		exp := m.expanded[p.ID]
		m.items = append(m.items, treeItem{
			isProject: true,
			projectID: p.ID,
			name:      p.Name,
			expanded:  exp,
		})
		if exp {
			for _, t := range p.Tasks {
				m.items = append(m.items, treeItem{
					isProject: false,
					projectID: p.ID,
					taskID:    t.ID,
					name:      t.Name,
				})
			}
		}
	}
}

// listHeight returns how many rows the list area should occupy.
// Panel outer height = content height + 2 (border, no vertical padding).
// Fixed content lines: title(1) + blank(1) + spacer(1) + help(1) = 4.
func (m *TreeModel) listHeight() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	fixed := 6 // 2 border + 4 fixed content lines
	if m.mode != inputNone {
		fixed += 2 // blank line + input line
	}
	if m.showFull {
		fixed += 3 // full help is 4 rows instead of 1
	}
	lh := h - fixed
	if lh < 1 {
		lh = 1
	}
	return lh
}

func (m *TreeModel) scrollToCursor() {
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

func (m *TreeModel) clampCursor() {
	if len(m.items) == 0 {
		m.cursor = 0
		return
	}
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// cursorProjectID returns the project ID relevant to the current cursor
// position (for adding a task under the right project).
func (m *TreeModel) cursorProjectID() string {
	if len(m.items) == 0 {
		return ""
	}
	return m.items[m.cursor].projectID
}

func (m TreeModel) Init() tea.Cmd { return nil }

func (m TreeModel) Update(msg tea.Msg) (TreeModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width - 4
		m.scrollToCursor()
		return m, nil

	case tickMsg:
		m.buildItems()
		return m, nil

	case tea.KeyMsg:
		// Always quit on ctrl+c; handled globally in App but guard here too.
		if msg.String() == "ctrl+c" {
			m.WantsQuit = true
			return m, nil
		}

		// ── Input mode: route keys to text input ──────────────────────────
		if m.mode != inputNone {
			switch msg.Type {
			case tea.KeyEnter:
				name := strings.TrimSpace(m.textInput.Value())
				if name != "" {
					if m.mode == inputProject {
						p := m.store.AddProject(name)
						_ = m.store.Save()
						m.expanded[p.ID] = true
						m.buildItems()
						// Move cursor to the new project row.
						for i, item := range m.items {
							if item.isProject && item.projectID == p.ID {
								m.cursor = i
								break
							}
						}
					} else { // inputTask
						t := m.store.AddTask(m.parentID, name)
						_ = m.store.Save()
						m.expanded[m.parentID] = true
						m.buildItems()
						if t != nil {
							for i, item := range m.items {
								if !item.isProject && item.taskID == t.ID {
									m.cursor = i
									break
								}
							}
						}
					}
				}
				m.mode = inputNone
				m.textInput.Blur()
				m.textInput.SetValue("")
				m.scrollToCursor()
				return m, nil

			case tea.KeyEscape:
				m.mode = inputNone
				m.textInput.Blur()
				m.textInput.SetValue("")
				return m, nil
			}

			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		// ── Normal navigation mode ─────────────────────────────────────────
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.scrollToCursor()
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.scrollToCursor()
			}

		case "left", "h":
			if len(m.items) > 0 {
				item := m.items[m.cursor]
				if item.isProject && m.expanded[item.projectID] {
					m.expanded[item.projectID] = false
					m.buildItems()
					m.clampCursor()
					m.scrollToCursor()
				}
			}

		case "right", "l":
			if len(m.items) > 0 {
				item := m.items[m.cursor]
				if item.isProject && !m.expanded[item.projectID] {
					m.expanded[item.projectID] = true
					m.buildItems()
					m.scrollToCursor()
				}
			}

		case "N":
			m.mode = inputProject
			m.textInput.Placeholder = "Project name…"
			m.textInput.SetValue("")
			cmd = m.textInput.Focus()
			return m, cmd

		case "n":
			pid := m.cursorProjectID()
			if pid == "" {
				break
			}
			m.mode = inputTask
			m.parentID = pid
			m.textInput.Placeholder = "Task name…"
			m.textInput.SetValue("")
			cmd = m.textInput.Focus()
			return m, cmd

		case "d":
			if len(m.items) == 0 {
				break
			}
			item := m.items[m.cursor]
			if item.isProject {
				if m.store.ActiveTimer != nil &&
					m.store.ActiveTimer.ProjectID == item.projectID {
					m.store.ActiveTimer = nil
				}
				delete(m.expanded, item.projectID)
				m.store.DeleteProject(item.projectID)
			} else {
				if m.store.ActiveTimer != nil &&
					m.store.ActiveTimer.TaskID == item.taskID {
					m.store.ActiveTimer = nil
				}
				m.store.DeleteTask(item.projectID, item.taskID)
			}
			_ = m.store.Save()
			m.buildItems()
			m.clampCursor()
			m.scrollToCursor()

		case "enter":
			if len(m.items) == 0 {
				break
			}
			item := m.items[m.cursor]
			if item.isProject {
				// Toggle expand/collapse on projects.
				m.expanded[item.projectID] = !m.expanded[item.projectID]
				m.buildItems()
				m.scrollToCursor()
			} else {
				m.startTimerForCursor()
			}

		case "s":
			if len(m.items) > 0 && !m.items[m.cursor].isProject {
				m.startTimerForCursor()
			}

		case "e":
			if len(m.items) == 0 {
				break
			}
			item := m.items[m.cursor]
			if !item.isProject {
				m.SelectedProjectID = item.projectID
				m.SelectedTaskID = item.taskID
				m.SwitchToEdit = true
			}

		case "R":
			m.SwitchToReport = true

		case "T":
			m.SwitchToThemePicker = true

		case "?":
			m.showFull = !m.showFull
			m.help.ShowAll = m.showFull
			m.scrollToCursor()

		case "q":
			m.WantsQuit = true
		}
	}

	return m, nil
}

func (m *TreeModel) startTimerForCursor() {
	if len(m.items) == 0 {
		return
	}
	item := m.items[m.cursor]
	if item.isProject {
		return
	}
	// If the running timer belongs to this exact task, navigate back to it
	// without starting a new one.  If it belongs to a different task, do
	// nothing — only one timer can run at a time.
	if m.store.ActiveTimer != nil {
		if m.store.ActiveTimer.TaskID == item.taskID {
			m.SelectedProjectID = item.projectID
			m.SelectedTaskID = item.taskID
			m.SwitchToTimer = true
		}
		return
	}
	now := time.Now()
	m.store.ActiveTimer = &model.ActiveTimer{
		ProjectID:     item.projectID,
		TaskID:        item.taskID,
		OriginalStart: now,
		Start:         now,
	}
	_ = m.store.Save()
	m.SelectedProjectID = item.projectID
	m.SelectedTaskID = item.taskID
	m.SwitchToTimer = true
}

func (m TreeModel) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	// innerW is the content width inside border (no vertical padding on panel).
	// StylePanel uses Padding(0,1): 1 char left + 1 char right = 2, plus 2 border = 4 total.
	innerW := w - 4
	if innerW < 20 {
		innerW = 20
	}

	var sb strings.Builder

	// ── Title ──────────────────────────────────────────────────────────────
	sb.WriteString(StyleTitle.Render("shellclock"))
	sb.WriteString("\n\n")

	// ── Tree list ──────────────────────────────────────────────────────────
	lh := m.listHeight()
	end := m.offset + lh
	if end > len(m.items) {
		end = len(m.items)
	}

	if len(m.items) == 0 {
		hint := StyleDimmed.Render("No projects yet — press N to create one.")
		sb.WriteString(hint)
		sb.WriteString("\n")
		// Fill remaining list rows.
		for i := 1; i < lh; i++ {
			sb.WriteString("\n")
		}
	} else {
		for i := m.offset; i < end; i++ {
			sb.WriteString(m.renderItem(i, innerW))
			sb.WriteString("\n")
		}
		// Pad to keep the layout stable when list is shorter than lh.
		for i := end - m.offset; i < lh; i++ {
			sb.WriteString("\n")
		}
	}

	// ── Input prompt (only in input mode) ─────────────────────────────────
	if m.mode != inputNone {
		label := "New project"
		if m.mode == inputTask {
			label = "New task"
		}
		sb.WriteString("\n")
		sb.WriteString(StyleInputLabel.Render(label+":") + " " + m.textInput.View())
		sb.WriteString("\n")
	}

	// ── Help bar ───────────────────────────────────────────────────────────
	sb.WriteString("\n")
	m.help.Styles = helpStyles()
	m.help.Width = innerW
	var km help.KeyMap
	if m.mode != inputNone {
		km = inputKeyMap{m.keys}
	} else {
		km = treeKeyMap{m.keys}
	}
	sb.WriteString(m.help.View(km))

	return StylePanel.
		Width(innerW).
		Padding(0, 1).
		Render(sb.String())
}

func (m TreeModel) renderItem(i, innerW int) string {
	item := m.items[i]
	selected := i == m.cursor

	if item.isProject {
		toggle := "▶"
		if item.expanded {
			toggle = "▼"
		}
		p := m.store.FindProject(item.projectID)
		nameText := toggle + " " + item.name
		durText := ""
		if p != nil && p.TotalSeconds() > 0 {
			durText = "  " + util.FormatDuration(p.TotalSeconds())
		}
		if selected {
			return highlightRow(nameText+durText, innerW)
		}
		return StyleProject.Render(nameText) + StyleDuration.Render(durText)
	}

	// Task row
	t := m.store.FindTask(item.projectID, item.taskID)
	nameText := "  · " + item.name
	activeStr := ""
	if m.store.ActiveTimer != nil && m.store.ActiveTimer.TaskID == item.taskID {
		activeStr = " ●"
	}
	durText := ""
	if t != nil && t.TotalSeconds() > 0 {
		durText = "  " + util.FormatDuration(t.TotalSeconds())
	}

	if selected {
		return highlightRow(nameText+activeStr+durText, innerW)
	}
	line := StyleTask.Render(nameText)
	if activeStr != "" {
		line += StyleTimer.Render(activeStr)
	}
	if durText != "" {
		line += StyleDuration.Render(durText)
	}
	return line
}

// highlightRow renders text using the selection style pinned to exactly w
// visible columns.  Width(w) pads short text; MaxWidth(w) clips long text.
// Together they guarantee a single-line highlight that never bleeds onto the
// next row regardless of terminal width.
func highlightRow(text string, w int) string {
	// Clip first so lipgloss never has to word-wrap.
	if lipgloss.Width(text) > w {
		runes := []rune(text)
		for len(runes) > 0 && lipgloss.Width(string(runes)) > w-1 {
			runes = runes[:len(runes)-1]
		}
		text = string(runes) + "…"
	}
	return StyleSelected.Width(w).MaxWidth(w).Render(text)
}

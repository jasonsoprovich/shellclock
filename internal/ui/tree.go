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

// renderLogo returns a 3-line lipgloss bordered wordmark for "shellclock".
func renderLogo() string {
	inner := StyleLogo.Render("◷") + "  " + StyleLogo.Bold(true).Render("shellclock")
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorOverlay).
		Padding(0, 1).
		Render(inner)
}

type inputMode int

const (
	inputNone    inputMode = iota
	inputProject           // naming a new project
	inputTask              // naming a new task
	inputRename            // renaming an existing project or task
	inputTags              // editing comma-separated tags for a project
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
	mode         inputMode
	textInput    textinput.Model
	parentID     string // project ID when mode == inputTask
	tagProjectID string // project ID when mode == inputTags

	// help component
	help     help.Model
	showFull bool

	// confirm-delete modal
	confirmActive bool
	confirmMsg    string
	confirmProjID string // project to delete (or parent of task)
	confirmTaskID string // task to delete; empty → delete the project

	// backup info overlay
	backupInfoActive bool
	backupList       []string

	// signals consumed by App
	WantsQuit           bool
	SwitchToTaskDetail  bool
	SwitchToEdit        bool
	SwitchToReport      bool
	SwitchToSummary     bool
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
// Fixed content lines: logo(3) + subline(1) + timer(1) + spacer(1) + help(1) = 7.
func (m *TreeModel) listHeight() int {
	h := m.height
	if h == 0 {
		h = 24
	}
	fixed := 9 // 2 border + 7 fixed content lines
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

		// ── Backup info overlay — any key closes it ───────────────────────
		if m.backupInfoActive {
			m.backupInfoActive = false
			return m, nil
		}

		// ── Confirm-delete modal ───────────────────────────────────────────
		if m.confirmActive {
			switch msg.String() {
			case "y":
				if m.confirmTaskID != "" {
					if m.store.ActiveTimer != nil && m.store.ActiveTimer.TaskID == m.confirmTaskID {
						m.store.ActiveTimer = nil
					}
					m.store.DeleteTask(m.confirmProjID, m.confirmTaskID)
				} else {
					if m.store.ActiveTimer != nil && m.store.ActiveTimer.ProjectID == m.confirmProjID {
						m.store.ActiveTimer = nil
					}
					delete(m.expanded, m.confirmProjID)
					m.store.DeleteProject(m.confirmProjID)
				}
				_ = m.store.Save()
				m.buildItems()
				m.clampCursor()
				m.scrollToCursor()
				m.confirmActive = false
				m.confirmMsg = ""
				m.confirmProjID = ""
				m.confirmTaskID = ""
			case "n", "esc":
				m.confirmActive = false
				m.confirmMsg = ""
				m.confirmProjID = ""
				m.confirmTaskID = ""
			}
			return m, nil
		}

		// ── Input mode: route keys to text input ──────────────────────────
		if m.mode != inputNone {
			switch msg.Type {
			case tea.KeyEnter:
				val := strings.TrimSpace(m.textInput.Value())
				switch m.mode {
				case inputProject:
					if val != "" {
						p := m.store.AddProject(val)
						_ = m.store.Save()
						m.expanded[p.ID] = true
						m.buildItems()
						for i, item := range m.items {
							if item.isProject && item.projectID == p.ID {
								m.cursor = i
								break
							}
						}
						// Transition to tag input for the new project.
						m.tagProjectID = p.ID
						m.mode = inputTags
						m.textInput.Placeholder = "Tags (comma-separated)…"
						m.textInput.SetValue("")
						cmd = m.textInput.Focus()
						return m, cmd
					}
				case inputTask:
					if val != "" {
						t := m.store.AddTask(m.parentID, val)
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
				case inputRename:
					if val != "" {
						item := m.items[m.cursor]
						if item.isProject {
							m.store.RenameProject(item.projectID, val)
							_ = m.store.Save()
							m.buildItems()
							// Transition to tag input for the renamed project.
							p := m.store.FindProject(item.projectID)
							currentTags := ""
							if p != nil && len(p.Tags) > 0 {
								currentTags = strings.Join(p.Tags, ", ")
							}
							m.tagProjectID = item.projectID
							m.mode = inputTags
							m.textInput.Placeholder = "Tags (comma-separated)…"
							m.textInput.SetValue(currentTags)
							m.textInput.CursorEnd()
							cmd = m.textInput.Focus()
							return m, cmd
						}
						m.store.RenameTask(item.projectID, item.taskID, val)
						_ = m.store.Save()
						m.buildItems()
					}
				case inputTags:
					// Allow empty value to clear all tags.
					var tags []string
					for _, t := range strings.Split(val, ",") {
						t = strings.TrimSpace(t)
						if t != "" {
							tags = append(tags, t)
						}
					}
					m.store.UpdateProjectTags(m.tagProjectID, tags)
					_ = m.store.Save()
					m.buildItems()
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
				p := m.store.FindProject(item.projectID)
				name := item.name
				if p != nil {
					name = p.Name
				}
				m.confirmMsg = "Delete project \"" + name + "\"?"
				m.confirmProjID = item.projectID
				m.confirmTaskID = ""
			} else {
				t := m.store.FindTask(item.projectID, item.taskID)
				name := item.name
				if t != nil {
					name = t.Name
				}
				m.confirmMsg = "Delete task \"" + name + "\"?"
				m.confirmProjID = item.projectID
				m.confirmTaskID = item.taskID
			}
			m.confirmActive = true

		case "enter":
			if len(m.items) == 0 {
				break
			}
			item := m.items[m.cursor]
			if item.isProject {
				m.expanded[item.projectID] = !m.expanded[item.projectID]
				m.buildItems()
				m.scrollToCursor()
			} else {
				m.SelectedProjectID = item.projectID
				m.SelectedTaskID = item.taskID
				m.SwitchToTaskDetail = true
			}

		case "s":
			m.SwitchToSummary = true

		case "p":
			at := m.store.ActiveTimer
			if at == nil {
				// Start timer for the focused task.
				if len(m.items) > 0 && !m.items[m.cursor].isProject {
					cmd = m.startTimerForCursor()
				}
			} else if at.Paused {
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

		case "S":
			at := m.store.ActiveTimer
			if at != nil {
				now := time.Now()
				elapsed := at.AccumulatedSeconds
				if !at.Paused {
					elapsed += int64(time.Since(at.Start).Seconds())
				}
				if elapsed > 0 {
					m.store.AddSession(at.ProjectID, at.TaskID, model.Session{
						ID:              uuid.NewString(),
						Start:           at.OriginalStart,
						End:             now,
						DurationSeconds: elapsed,
					})
				}
				m.store.ActiveTimer = nil
				_ = m.store.Save()
				m.buildItems()
			}

		case "r":
			at := m.store.ActiveTimer
			if at != nil {
				now := time.Now()
				at.AccumulatedSeconds = 0
				at.OriginalStart = now
				at.Start = now
				at.Paused = false
				_ = m.store.Save()
				cmd = tick()
			}

		case "E":
			if len(m.items) == 0 {
				break
			}
			item := m.items[m.cursor]
			m.mode = inputRename
			m.textInput.Placeholder = "New name…"
			m.textInput.SetValue(item.name)
			m.textInput.CursorEnd()
			cmd = m.textInput.Focus()
			return m, cmd

		case "#":
			if len(m.items) > 0 && m.items[m.cursor].isProject {
				item := m.items[m.cursor]
				p := m.store.FindProject(item.projectID)
				currentTags := ""
				if p != nil && len(p.Tags) > 0 {
					currentTags = strings.Join(p.Tags, ", ")
				}
				m.tagProjectID = item.projectID
				m.mode = inputTags
				m.textInput.Placeholder = "Tags (comma-separated)…"
				m.textInput.SetValue(currentTags)
				m.textInput.CursorEnd()
				cmd = m.textInput.Focus()
				return m, cmd
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

		case "B":
			backs, _ := model.ListBackups()
			m.backupList = backs
			m.backupInfoActive = true

		case "?":
			m.showFull = !m.showFull
			m.help.ShowAll = m.showFull
			m.scrollToCursor()

		case "q":
			m.WantsQuit = true
		}
	}

	return m, cmd
}

// startTimerForCursor starts a timer on the focused task and returns a tick
// cmd to begin the one-second update chain. Does nothing if a timer is already
// running on any task (only one timer at a time).
func (m *TreeModel) startTimerForCursor() tea.Cmd {
	if len(m.items) == 0 {
		return nil
	}
	item := m.items[m.cursor]
	if item.isProject || m.store.ActiveTimer != nil {
		return nil
	}
	now := time.Now()
	m.store.ActiveTimer = &model.ActiveTimer{
		ProjectID:     item.projectID,
		TaskID:        item.taskID,
		OriginalStart: now,
		Start:         now,
	}
	_ = m.store.Save()
	return tick()
}

// renderHeader renders the lipgloss logo wordmark and a one-line stats summary.
func (m *TreeModel) renderHeader() string {
	var sb strings.Builder
	sb.WriteString(renderLogo())
	sb.WriteString("\n")

	nProjects := len(m.store.Projects)
	nTasks := 0
	var totalSecs int64
	for _, p := range m.store.Projects {
		nTasks += len(p.Tasks)
		totalSecs += p.TotalSeconds()
	}

	pWord := "projects"
	if nProjects == 1 {
		pWord = "project"
	}
	tWord := "tasks"
	if nTasks == 1 {
		tWord = "task"
	}
	totalStr := util.FormatDuration(totalSecs) + " total"
	if totalSecs == 0 {
		totalStr = "no time tracked"
	}

	dateStr := time.Now().Format("Mon Jan 2, 2006")
	stats := fmt.Sprintf("%s  ·  %d %s  ·  %d %s  ·  %s",
		dateStr, nProjects, pWord, nTasks, tWord, totalStr)
	sb.WriteString(StyleDimmed.Render(stats))
	return sb.String()
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

	// ── Logo + stats subline (4 lines) ────────────────────────────────────
	sb.WriteString(m.renderHeader())
	sb.WriteString("\n")

	// ── Active timer status (always 1 line) ───────────────────────────────
	if m.store.ActiveTimer != nil {
		at := m.store.ActiveTimer
		elapsed := at.AccumulatedSeconds
		if !at.Paused {
			elapsed += int64(time.Since(at.Start).Seconds())
		}
		p := m.store.FindProject(at.ProjectID)
		t := m.store.FindTask(at.ProjectID, at.TaskID)
		pn, tn := "?", "?"
		if p != nil {
			pn = truncate(p.Name, 20)
		}
		if t != nil {
			tn = truncate(t.Name, 20)
		}
		pauseStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		var badge string
		if at.Paused {
			badge = pauseStyle.Render("⏸  PAUSED  " + util.FormatDurationShort(elapsed))
		} else {
			badge = StyleTimer.Render("●  RUNNING  " + util.FormatDurationShort(elapsed))
		}
		sb.WriteString(
			badge +
				StyleDimmed.Render("  ") +
				StyleProject.Render(pn) +
				StyleDimmed.Render(" › ") +
				StyleTask.Render(tn),
		)
	}
	sb.WriteString("\n")

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
		var label string
		switch m.mode {
		case inputProject:
			label = "New project"
		case inputTask:
			label = "New task"
		case inputRename:
			label = "Rename"
		case inputTags:
			label = "Tags"
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

	panel := StylePanel.
		Width(innerW + 2).
		Padding(0, 1).
		Render(sb.String())
	if m.confirmActive {
		return renderConfirmOverlay(panel, m.confirmMsg, m.width, m.height)
	}
	if m.backupInfoActive {
		return renderBackupOverlay(panel, m.backupList, m.width, m.height)
	}
	return panel
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
			tagText := ""
			if p != nil && len(p.Tags) > 0 {
				tagText = "  [" + strings.Join(p.Tags, ", ") + "]"
			}
			return highlightRow(nameText+tagText+durText, innerW)
		}
		line := StyleProject.Render(nameText)
		if p != nil {
			for _, tag := range p.Tags {
				line += " " + renderTagPill(tag)
			}
		}
		line += StyleDuration.Render(durText)
		return line
	}

	// Task row
	t := m.store.FindTask(item.projectID, item.taskID)
	nameText := "  · " + item.name
	activeStr := ""
	if m.store.ActiveTimer != nil && m.store.ActiveTimer.TaskID == item.taskID {
		at := m.store.ActiveTimer
		elapsed := at.AccumulatedSeconds
		if !at.Paused {
			elapsed += int64(time.Since(at.Start).Seconds())
		}
		if at.Paused {
			activeStr = "  ⏸  " + util.FormatDurationShort(elapsed)
		} else {
			activeStr = "  ●  " + util.FormatDurationShort(elapsed)
		}
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
// visible columns. Text is manually padded/clipped to w chars first so
// lipgloss never word-wraps at spaces inside the styled element.
func highlightRow(text string, w int) string {
	vis := lipgloss.Width(text)
	if vis > w {
		runes := []rune(text)
		for len(runes) > 0 && lipgloss.Width(string(runes)) > w-1 {
			runes = runes[:len(runes)-1]
		}
		text = string(runes) + "…"
		vis = w
	}
	if vis < w {
		text += strings.Repeat(" ", w-vis)
	}
	return StyleSelected.Render(text)
}

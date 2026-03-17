package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

// treeItem is a flattened row in the tree view.
type treeItem struct {
	isProject bool
	projectID string
	taskID    string
	name      string
	depth     int
	expanded  bool // only meaningful for projects
}

// TreeModel manages the project/task tree.
type TreeModel struct {
	store    *model.Store
	keys     KeyMap
	items    []treeItem
	cursor   int
	expanded map[string]bool // project IDs

	// navigation signals to App
	SwitchToTimer    bool
	SwitchToEdit     bool
	SwitchToReport   bool
	SelectedProjectID string
	SelectedTaskID    string
}

func NewTreeModel(store *model.Store, keys KeyMap) TreeModel {
	m := TreeModel{
		store:    store,
		keys:     keys,
		expanded: make(map[string]bool),
	}
	// Expand all projects by default.
	for _, p := range store.Projects {
		m.expanded[p.ID] = true
	}
	m.buildItems()
	return m
}

func (m *TreeModel) buildItems() {
	m.items = nil
	for _, p := range m.store.Projects {
		exp := m.expanded[p.ID]
		m.items = append(m.items, treeItem{
			isProject: true,
			projectID: p.ID,
			name:      p.Name,
			depth:     0,
			expanded:  exp,
		})
		if exp {
			for _, t := range p.Tasks {
				m.items = append(m.items, treeItem{
					isProject: false,
					projectID: p.ID,
					taskID:    t.ID,
					name:      t.Name,
					depth:     1,
				})
			}
		}
	}
}

func (m TreeModel) Init() tea.Cmd { return nil }

func (m TreeModel) Update(msg tea.Msg) (TreeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "up" || msg.String() == "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case msg.String() == "down" || msg.String() == "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case msg.String() == "left" || msg.String() == "h":
			if m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if item.isProject && m.expanded[item.projectID] {
					m.expanded[item.projectID] = false
					m.buildItems()
				}
			}
		case msg.String() == "right" || msg.String() == "l":
			if m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if item.isProject && !m.expanded[item.projectID] {
					m.expanded[item.projectID] = true
					m.buildItems()
				}
			}
		case msg.String() == "enter", msg.String() == "s":
			if m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if !item.isProject {
					m.SelectedProjectID = item.projectID
					m.SelectedTaskID = item.taskID
					m.SwitchToTimer = true
				}
			}
		case msg.String() == "e":
			if m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if !item.isProject {
					m.SelectedProjectID = item.projectID
					m.SelectedTaskID = item.taskID
					m.SwitchToEdit = true
				}
			}
		case msg.String() == "R":
			m.SwitchToReport = true
		case msg.String() == "n":
			// Add a placeholder project; a real implementation would prompt.
			m.store.AddProject(fmt.Sprintf("Project %d", len(m.store.Projects)+1))
			_ = m.store.Save()
			m.buildItems()
		}
	case tickMsg:
		// Rebuild to refresh totals if a timer is running.
		m.buildItems()
	}
	return m, nil
}

func (m TreeModel) View() string {
	var sb strings.Builder

	sb.WriteString(StyleTitle.Render("shellclock") + "\n\n")

	for i, item := range m.items {
		indent := strings.Repeat("  ", item.depth)
		var line string

		if item.isProject {
			toggle := "▶"
			if item.expanded {
				toggle = "▼"
			}
			p := m.store.FindProject(item.projectID)
			dur := ""
			if p != nil {
				dur = " " + StyleDuration.Render(util.FormatDuration(p.TotalSeconds()))
			}
			text := fmt.Sprintf("%s%s %s%s", indent, toggle, item.name, dur)
			if i == m.cursor {
				line = StyleSelected.Render(text)
			} else {
				line = StyleProject.Render(text)
			}
		} else {
			t := m.store.FindTask(item.projectID, item.taskID)
			dur := ""
			if t != nil {
				dur = " " + StyleDuration.Render(util.FormatDuration(t.TotalSeconds()))
			}
			active := ""
			if m.store.ActiveTimer != nil && m.store.ActiveTimer.TaskID == item.taskID {
				active = StyleTimer.Render(" ●")
			}
			text := fmt.Sprintf("%s  %s%s%s", indent, item.name, active, dur)
			if i == m.cursor {
				line = StyleSelected.Render(text)
			} else {
				line = StyleTask.Render(text)
			}
		}
		sb.WriteString(line + "\n")
	}

	if len(m.items) == 0 {
		sb.WriteString(StyleDimmed.Render("No projects yet. Press n to create one.") + "\n")
	}

	sb.WriteString("\n" + StyleHelp.Render("↑/↓ navigate  →/← expand/collapse  s start timer  e edit sessions  R report  n new project  q quit"))
	return StylePanel.Render(sb.String())
}

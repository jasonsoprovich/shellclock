package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

// EditModel allows viewing and deleting sessions for a task.
type EditModel struct {
	store     *model.Store
	keys      KeyMap
	ProjectID string
	TaskID    string
	cursor    int
	width     int
	height    int

	SwitchToTree bool
}

func NewEditModel(store *model.Store, keys KeyMap) EditModel {
	return EditModel{store: store, keys: keys}
}

func (m EditModel) sessions() []model.Session {
	t := m.store.FindTask(m.ProjectID, m.TaskID)
	if t == nil {
		return nil
	}
	return t.Sessions
}

func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		return m, nil
	}
	sessions := m.sessions()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.SwitchToTree = true
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(sessions)-1 {
				m.cursor++
			}
		case "d":
			if m.cursor < len(sessions) {
				sess := sessions[m.cursor]
				m.store.DeleteSession(m.ProjectID, m.TaskID, sess.ID)
				_ = m.store.Save()
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "n":
			// Add a 1-minute placeholder session for manual entry.
			sess := model.Session{
				ID:              uuid.NewString(),
				Start:           time.Now().Add(-time.Minute),
				End:             time.Now(),
				DurationSeconds: 60,
			}
			m.store.AddSession(m.ProjectID, m.TaskID, sess)
			_ = m.store.Save()
		}
	}
	return m, nil
}

func (m EditModel) View() string {
	var sb strings.Builder

	t := m.store.FindTask(m.ProjectID, m.TaskID)
	taskName := "?"
	if t != nil {
		taskName = t.Name
	}

	sb.WriteString(StyleTitle.Render("Edit Sessions") + " — " + StyleTask.Render(taskName) + "\n\n")

	sessions := m.sessions()
	if len(sessions) == 0 {
		sb.WriteString(StyleDimmed.Render("No sessions. Press n to add one.") + "\n")
	}

	for i, sess := range sessions {
		line := fmt.Sprintf("%s  start: %s  dur: %s",
			sess.ID[:8],
			sess.Start.Format("2006-01-02 15:04"),
			util.FormatDuration(sess.DurationSeconds),
		)
		if i == m.cursor {
			sb.WriteString(StyleSelected.Render(line) + "\n")
		} else {
			sb.WriteString(StyleTask.Render(line) + "\n")
		}
	}

	total := int64(0)
	for _, s := range sessions {
		total += s.DurationSeconds
	}
	sb.WriteString("\n" + StyleDuration.Render(fmt.Sprintf("Total: %s", util.FormatDuration(total))) + "\n")
	sb.WriteString("\n" + StyleDimmed.Render("↑/↓ navigate  d delete  n add placeholder  esc back"))

	return strings.Join([]string{StylePanel.Render(sb.String())}, "")
}

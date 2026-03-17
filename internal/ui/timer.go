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

// TimerModel manages the live timer view.
type TimerModel struct {
	store *model.Store
	keys  KeyMap

	SwitchToTree bool
}

func NewTimerModel(store *model.Store, keys KeyMap) TimerModel {
	return TimerModel{store: store, keys: keys}
}

func (m TimerModel) Init() tea.Cmd {
	if m.store.ActiveTimer != nil && !m.store.ActiveTimer.Paused {
		return tick()
	}
	return nil
}

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
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.SwitchToTree = true
		case "s":
			// Start timer if none running.
			if m.store.ActiveTimer == nil {
				return m, nil // need project/task context — handled via tree
			}
		case "p":
			at := m.store.ActiveTimer
			if at == nil {
				return m, nil
			}
			if at.Paused {
				// Resume
				at.Paused = false
				at.Start = time.Now()
				_ = m.store.Save()
				return m, tick()
			}
			// Pause
			at.AccumulatedSeconds += int64(time.Since(at.Start).Seconds())
			at.Paused = true
			_ = m.store.Save()
		case "S":
			// Stop and commit session.
			at := m.store.ActiveTimer
			if at == nil {
				return m, nil
			}
			elapsed := m.elapsed()
			sess := model.Session{
				ID:              uuid.NewString(),
				Start:           at.Start,
				End:             time.Now(),
				DurationSeconds: elapsed,
			}
			m.store.AddSession(at.ProjectID, at.TaskID, sess)
			m.store.ActiveTimer = nil
			_ = m.store.Save()
			m.SwitchToTree = true
		case "r":
			// Reset timer.
			at := m.store.ActiveTimer
			if at == nil {
				return m, nil
			}
			at.AccumulatedSeconds = 0
			at.Start = time.Now()
			at.Paused = false
			_ = m.store.Save()
			return m, tick()
		}

	case tickMsg:
		if m.store.ActiveTimer != nil && !m.store.ActiveTimer.Paused {
			return m, tick()
		}
	}
	return m, nil
}

func (m TimerModel) View() string {
	var sb strings.Builder

	at := m.store.ActiveTimer
	if at == nil {
		sb.WriteString(StyleDimmed.Render("No active timer.") + "\n")
		sb.WriteString("\n" + StyleHelp.Render("esc back"))
		return StylePanel.Render(sb.String())
	}

	p := m.store.FindProject(at.ProjectID)
	t := m.store.FindTask(at.ProjectID, at.TaskID)

	projectName := "?"
	taskName := "?"
	if p != nil {
		projectName = p.Name
	}
	if t != nil {
		taskName = t.Name
	}

	status := "RUNNING"
	statusStyle := StyleTimer
	if at.Paused {
		status = "PAUSED"
		statusStyle = StyleDimmed
	}

	elapsed := m.elapsed()

	sb.WriteString(StyleTitle.Render("Timer") + "\n\n")
	sb.WriteString(fmt.Sprintf("%s  /  %s\n\n", StyleProject.Render(projectName), StyleTask.Render(taskName)))
	sb.WriteString(StyleTimer.Render(util.FormatDurationShort(elapsed)) + "\n\n")
	sb.WriteString(statusStyle.Render(status) + "\n")
	sb.WriteString("\n" + StyleHelp.Render("p pause/resume  S stop & save  r reset  esc back"))

	return StylePanel.Render(sb.String())
}

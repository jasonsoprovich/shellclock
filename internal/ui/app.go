package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jasonsoprovich/shellclock/internal/model"
)

// view identifies which screen is active.
type view int

const (
	viewTree view = iota
	viewTimer
	viewReport
	viewEdit
)

// tickMsg is sent every second to update a running timer.
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// App is the root Bubble Tea model.
type App struct {
	store   *model.Store
	keys    KeyMap
	current view

	tree   TreeModel
	timer  TimerModel
	report ReportModel
	edit   EditModel

	width  int
	height int
}

// New creates an initialised App.
func New(store *model.Store) App {
	keys := DefaultKeyMap()
	return App{
		store:   store,
		keys:    keys,
		current: viewTree,
		tree:    NewTreeModel(store, keys),
		timer:   NewTimerModel(store, keys),
		report:  NewReportModel(store, keys),
		edit:    NewEditModel(store, keys),
	}
}

func (a App) Init() tea.Cmd {
	cmds := []tea.Cmd{a.tree.Init()}
	if a.store.ActiveTimer != nil && !a.store.ActiveTimer.Paused {
		cmds = append(cmds, tick())
	}
	return tea.Batch(cmds...)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return a, tea.Quit
		case msg.String() == "q" && a.current == viewTree:
			return a, tea.Quit
		}

	case tickMsg:
		if a.store.ActiveTimer != nil && !a.store.ActiveTimer.Paused {
			return a, tick()
		}
	}

	var cmd tea.Cmd
	switch a.current {
	case viewTree:
		a.tree, cmd = a.tree.Update(msg)
		if a.tree.SwitchToTimer {
			a.tree.SwitchToTimer = false
			a.timer = NewTimerModel(a.store, a.keys)
			a.current = viewTimer
			return a, tea.Batch(cmd, a.timer.Init())
		}
		if a.tree.SwitchToEdit {
			a.tree.SwitchToEdit = false
			a.edit = NewEditModel(a.store, a.keys)
			a.edit.ProjectID = a.tree.SelectedProjectID
			a.edit.TaskID = a.tree.SelectedTaskID
			a.current = viewEdit
			return a, cmd
		}
		if a.tree.SwitchToReport {
			a.tree.SwitchToReport = false
			a.report = NewReportModel(a.store, a.keys)
			a.current = viewReport
			return a, cmd
		}

	case viewTimer:
		a.timer, cmd = a.timer.Update(msg)
		if a.timer.SwitchToTree {
			a.timer.SwitchToTree = false
			a.tree = NewTreeModel(a.store, a.keys)
			a.current = viewTree
		}

	case viewReport:
		a.report, cmd = a.report.Update(msg)
		if a.report.SwitchToTree {
			a.report.SwitchToTree = false
			a.current = viewTree
		}

	case viewEdit:
		a.edit, cmd = a.edit.Update(msg)
		if a.edit.SwitchToTree {
			a.edit.SwitchToTree = false
			a.tree = NewTreeModel(a.store, a.keys)
			a.current = viewTree
		}
	}

	return a, cmd
}

func (a App) View() string {
	switch a.current {
	case viewTimer:
		return a.timer.View()
	case viewReport:
		return a.report.View()
	case viewEdit:
		return a.edit.View()
	default:
		return a.tree.View()
	}
}

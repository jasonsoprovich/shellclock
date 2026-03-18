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
	viewThemePicker
)

// tickMsg is sent every second to drive the live timer display.
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
	picker ThemePickerModel

	width  int
	height int
}

// New creates an initialised App.  The stored theme (if any) is applied before
// any sub-models are built so every style var is correct from the first render.
// If a timer was running when the process last exited we open straight to the
// timer view so the user sees the live elapsed time immediately.
func New(store *model.Store) App {
	// Apply the persisted theme first so all Style vars are correct before any
	// sub-model is constructed.
	ApplyTheme(ThemeByName(store.Theme))

	keys := DefaultKeyMap()
	app := App{
		store:  store,
		keys:   keys,
		tree:   NewTreeModel(store, keys),
		timer:  NewTimerModel(store, keys),
		report: NewReportModel(store, keys),
		edit:   NewEditModel(store, keys),
		picker: NewThemePickerModel(store, keys),
	}
	if store.ActiveTimer != nil {
		app.current = viewTimer
	} else {
		app.current = viewTree
	}
	return app
}

// Init starts the appropriate sub-model.  When a timer is already running it
// delegates to TimerModel.Init() which starts the one-second tick chain.
func (a App) Init() tea.Cmd {
	if a.current == viewTimer {
		return a.timer.Init()
	}
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ctrl+c always quits regardless of active view or input mode.
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return a, tea.Quit
	}

	// Forward window size to every sub-model so they size themselves correctly.
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		a.width = ws.Width
		a.height = ws.Height
		a.tree, _ = a.tree.Update(ws)
		a.timer, _ = a.timer.Update(ws)
		a.report, _ = a.report.Update(ws)
		a.edit, _ = a.edit.Update(ws)
		a.picker, _ = a.picker.Update(ws)
		return a, nil
	}

	var cmd tea.Cmd

	switch a.current {
	// ── Tree ───────────────────────────────────────────────────────────────
	case viewTree:
		a.tree, cmd = a.tree.Update(msg)

		if a.tree.WantsQuit {
			return a, tea.Quit
		}
		if a.tree.SwitchToTimer {
			a.tree.SwitchToTimer = false
			a.timer = NewTimerModel(a.store, a.keys)
			a.timer.width = a.width
			a.timer.height = a.height
			a.current = viewTimer
			return a, tea.Batch(cmd, a.timer.Init())
		}
		if a.tree.SwitchToEdit {
			a.tree.SwitchToEdit = false
			a.edit = NewEditModel(a.store, a.keys)
			a.edit.ProjectID = a.tree.SelectedProjectID
			a.edit.TaskID = a.tree.SelectedTaskID
			a.edit.width = a.width
			a.edit.height = a.height
			a.current = viewEdit
			return a, cmd
		}
		if a.tree.SwitchToReport {
			a.tree.SwitchToReport = false
			a.report = NewReportModel(a.store, a.keys)
			a.report.width = a.width
			a.report.height = a.height
			a.current = viewReport
			return a, cmd
		}
		if a.tree.SwitchToThemePicker {
			a.tree.SwitchToThemePicker = false
			a.picker = NewThemePickerModel(a.store, a.keys)
			a.picker.width = a.width
			a.picker.height = a.height
			a.current = viewThemePicker
			return a, cmd
		}

	// ── Timer ──────────────────────────────────────────────────────────────
	case viewTimer:
		a.timer, cmd = a.timer.Update(msg)
		if a.timer.SwitchToTree {
			a.timer.SwitchToTree = false
			a.tree.buildItems() // refresh totals/indicator
			a.current = viewTree
		}

	// ── Report ─────────────────────────────────────────────────────────────
	case viewReport:
		a.report, cmd = a.report.Update(msg)
		if a.report.SwitchToTree {
			a.report.SwitchToTree = false
			a.current = viewTree
		}

	// ── Edit ───────────────────────────────────────────────────────────────
	case viewEdit:
		a.edit, cmd = a.edit.Update(msg)
		if a.edit.SwitchToTree {
			a.edit.SwitchToTree = false
			a.tree.buildItems() // refresh session totals
			a.current = viewTree
		}

	// ── Theme picker ────────────────────────────────────────────────────────
	case viewThemePicker:
		a.picker, cmd = a.picker.Update(msg)
		if a.picker.SwitchToTree {
			a.picker.SwitchToTree = false
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
	case viewThemePicker:
		return a.picker.View()
	default:
		return a.tree.View()
	}
}

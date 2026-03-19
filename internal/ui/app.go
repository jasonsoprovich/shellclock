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
	viewTaskDetail
	viewReport
	viewEdit
	viewThemePicker
	viewHelp
	viewSummary
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
	version string
	current view

	tree    TreeModel
	detail  TaskDetailModel
	report  ReportModel
	edit    EditModel
	picker  ThemePickerModel
	help    HelpModel
	summary SummaryModel

	width  int
	height int
}

// New creates an initialised App. The stored theme is applied before any
// sub-model is built so every style var is correct from the first render.
// The app always opens on the tree view; the active timer (if any) is shown
// inline there.
func New(store *model.Store, version string) App {
	ApplyTheme(ThemeByName(store.Theme))

	keys := DefaultKeyMap()
	return App{
		store:   store,
		keys:    keys,
		version: version,
		current: viewTree,
		tree:    NewTreeModel(store, keys),
		detail:  NewTaskDetailModel(store, keys),
		report:  NewReportModel(store, keys),
		edit:    NewEditModel(store, keys),
		picker:  NewThemePickerModel(store, keys),
		help:    NewHelpModel(keys, version),
		summary: NewSummaryModel(store, keys),
	}
}

// Init starts the tick chain if a timer is already running when the app
// launches. The App owns the tick chain for all views.
func (a App) Init() tea.Cmd {
	if a.store.ActiveTimer != nil && !a.store.ActiveTimer.Paused {
		return tick()
	}
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return a, tea.Quit
	}

	// Forward window size to every sub-model.
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		a.width = ws.Width
		a.height = ws.Height
		a.tree, _ = a.tree.Update(ws)
		a.detail, _ = a.detail.Update(ws)
		a.report, _ = a.report.Update(ws)
		a.edit, _ = a.edit.Update(ws)
		a.picker, _ = a.picker.Update(ws)
		a.help, _ = a.help.Update(ws)
		a.summary, _ = a.summary.Update(ws)
		return a, nil
	}

	// H opens the help screen from any view.
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "H" {
		a.help = NewHelpModel(a.keys, a.version)
		a.help.width = a.width
		a.help.height = a.height
		a.current = viewHelp
		return a, nil
	}

	var cmds []tea.Cmd

	// Global tick chain: App continues it for every view whenever the timer
	// is running. Sub-models only return tick() when they *start* or *resume*
	// a timer — they never self-perpetuate the chain.
	if _, isTick := msg.(tickMsg); isTick {
		at := a.store.ActiveTimer
		if at != nil && !at.Paused {
			cmds = append(cmds, tick())
		}
	}

	switch a.current {
	// ── Tree ─────────────────────────────────────────────────────────────────
	case viewTree:
		var c tea.Cmd
		a.tree, c = a.tree.Update(msg)
		cmds = append(cmds, c)

		if a.tree.WantsQuit {
			return a, tea.Quit
		}
		if a.tree.SwitchToTaskDetail {
			a.tree.SwitchToTaskDetail = false
			a.detail = NewTaskDetailModel(a.store, a.keys)
			a.detail.ProjectID = a.tree.SelectedProjectID
			a.detail.TaskID = a.tree.SelectedTaskID
			a.detail.width = a.width
			a.detail.height = a.height
			a.current = viewTaskDetail
		}
		if a.tree.SwitchToEdit {
			a.tree.SwitchToEdit = false
			a.edit = NewEditModel(a.store, a.keys)
			a.edit.ProjectID = a.tree.SelectedProjectID
			a.edit.TaskID = a.tree.SelectedTaskID
			a.edit.width = a.width
			a.edit.height = a.height
			a.current = viewEdit
		}
		if a.tree.SwitchToReport {
			a.tree.SwitchToReport = false
			a.report = NewReportModel(a.store, a.keys)
			a.report.width = a.width
			a.report.height = a.height
			a.current = viewReport
		}
		if a.tree.SwitchToSummary {
			a.tree.SwitchToSummary = false
			a.summary = NewSummaryModel(a.store, a.keys)
			a.summary.width = a.width
			a.summary.height = a.height
			a.current = viewSummary
		}
		if a.tree.SwitchToThemePicker {
			a.tree.SwitchToThemePicker = false
			a.picker = NewThemePickerModel(a.store, a.keys)
			a.picker.width = a.width
			a.picker.height = a.height
			a.current = viewThemePicker
		}

	// ── Task detail ───────────────────────────────────────────────────────────
	case viewTaskDetail:
		var c tea.Cmd
		a.detail, c = a.detail.Update(msg)
		cmds = append(cmds, c)

		if a.detail.SwitchToTree {
			a.detail.SwitchToTree = false
			a.tree.buildItems()
			a.current = viewTree
		}

	// ── Report ────────────────────────────────────────────────────────────────
	case viewReport:
		var c tea.Cmd
		a.report, c = a.report.Update(msg)
		cmds = append(cmds, c)

		if a.report.SwitchToTree {
			a.report.SwitchToTree = false
			a.current = viewTree
		}

	// ── Edit ──────────────────────────────────────────────────────────────────
	case viewEdit:
		var c tea.Cmd
		a.edit, c = a.edit.Update(msg)
		cmds = append(cmds, c)

		if a.edit.SwitchToTree {
			a.edit.SwitchToTree = false
			a.tree.buildItems()
			a.current = viewTree
		}

	// ── Theme picker ──────────────────────────────────────────────────────────
	case viewThemePicker:
		var c tea.Cmd
		a.picker, c = a.picker.Update(msg)
		cmds = append(cmds, c)

		if a.picker.SwitchToTree {
			a.picker.SwitchToTree = false
			a.current = viewTree
		}

	// ── Help screen ───────────────────────────────────────────────────────────
	case viewHelp:
		var c tea.Cmd
		a.help, c = a.help.Update(msg)
		cmds = append(cmds, c)

		if a.help.SwitchToTree {
			a.help.SwitchToTree = false
			a.current = viewTree
		}

	// ── Summary ───────────────────────────────────────────────────────────────
	case viewSummary:
		var c tea.Cmd
		a.summary, c = a.summary.Update(msg)
		cmds = append(cmds, c)

		if a.summary.SwitchToTree {
			a.summary.SwitchToTree = false
			a.current = viewTree
		}
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	switch a.current {
	case viewTaskDetail:
		return a.detail.View()
	case viewReport:
		return a.report.View()
	case viewEdit:
		return a.edit.View()
	case viewThemePicker:
		return a.picker.View()
	case viewHelp:
		return a.help.View()
	case viewSummary:
		return a.summary.View()
	default:
		return a.tree.View()
	}
}

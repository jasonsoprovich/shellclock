package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all application-wide keybindings.
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Enter        key.Binding
	Esc          key.Binding
	NewProject   key.Binding
	NewTask      key.Binding
	Delete       key.Binding
	Edit         key.Binding
	Rename       key.Binding
	EditTags     key.Binding
	Start        key.Binding
	TreeTimer    key.Binding // start/pause/resume timer from tree view (p)
	Stop         key.Binding
	Reset        key.Binding
	Report       key.Binding
	Summary      key.Binding
	ThemePicker  key.Binding
	Export       key.Binding
	Filter       key.Binding
	BackupInfo   key.Binding
	MasterReset  key.Binding
	Quit         key.Binding
	Help         key.Binding
	HelpScreen   key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:        key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "collapse")),
		Right:       key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "expand")),
		Enter:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select/expand")),
		Esc:         key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		NewProject:  key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "new project")),
		NewTask:     key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new task")),
		Delete:      key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		Edit:        key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit sessions")),
		Rename:      key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "rename")),
		EditTags:    key.NewBinding(key.WithKeys("#"), key.WithHelp("#", "edit tags")),
		Start:       key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start/pause/resume")),
		TreeTimer:   key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play/pause timer")),
		Stop:        key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "stop & save")),
		Reset:       key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset")),
		Report:      key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "report")),
		Summary:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "summary")),
		ThemePicker: key.NewBinding(key.WithKeys("T"), key.WithHelp("T", "theme")),
		Export:      key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "export")),
		Filter:      key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "filter by tag")),
		BackupInfo:  key.NewBinding(key.WithKeys("B"), key.WithHelp("B", "backups")),
		MasterReset: key.NewBinding(key.WithKeys("X"), key.WithHelp("X", "system reset")),
		Quit:        key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
		Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more keys")),
		HelpScreen:  key.NewBinding(key.WithKeys("H"), key.WithHelp("H", "help")),
	}
}

// ── Tree ────────────────────────────────────────────────────────────────────

type treeKeyMap struct{ km KeyMap }

func (k treeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.km.Up, k.km.Down,
		k.km.NewProject, k.km.NewTask,
		k.km.Enter, k.km.Delete,
		k.km.Summary, k.km.Report, k.km.Quit, k.km.HelpScreen, k.km.Help,
	}
}

func (k treeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.km.Up, k.km.Down, k.km.Left, k.km.Right},
		{k.km.NewProject, k.km.NewTask, k.km.Rename, k.km.Delete},
		{k.km.Enter, k.km.Edit, k.km.Report, k.km.ThemePicker},
		{k.km.TreeTimer, k.km.Stop, k.km.Reset, k.km.EditTags},
		{k.km.Summary, k.km.BackupInfo, k.km.MasterReset, k.km.HelpScreen, k.km.Quit, k.km.Help},
	}
}

// ── Report ───────────────────────────────────────────────────────────────────

type reportKeyMap struct{ km KeyMap }

func (k reportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.km.Up, k.km.Down, k.km.Filter, k.km.Export, k.km.Esc, k.km.Help}
}

func (k reportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.km.Up, k.km.Down},
		{k.km.Filter, k.km.Export, k.km.Esc, k.km.Help},
	}
}

// ── Tree text-input prompt ────────────────────────────────────────────────────

type inputKeyMap struct{ km KeyMap }

func (k inputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.km.Enter, k.km.Esc}
}

func (k inputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.km.Enter, k.km.Esc}}
}

// ── Session edit (normal mode) ────────────────────────────────────────────────

// editKeyMap uses inline bindings so the help bar shows session-specific
// labels ("add session", "edit session") instead of the tree-view labels
// ("new task", "edit sessions") that share the same keys.
type editKeyMap struct{ km KeyMap }

func (k editKeyMap) ShortHelp() []key.Binding {
	addSess := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "add session"))
	editSess := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit session"))
	delSess := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete session"))
	return []key.Binding{
		k.km.Up, k.km.Down,
		addSess, editSess, delSess,
		k.km.Esc, k.km.Help,
	}
}

func (k editKeyMap) FullHelp() [][]key.Binding {
	addSess := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "add session"))
	editSess := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit session"))
	delSess := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete session"))
	return [][]key.Binding{
		{k.km.Up, k.km.Down},
		{addSess, editSess, delSess},
		{k.km.Esc, k.km.Help},
	}
}

// ── Session edit (form mode) ──────────────────────────────────────────────────

type editFormKeyMap struct{ km KeyMap }

func (k editFormKeyMap) ShortHelp() []key.Binding {
	tab := key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field"))
	return []key.Binding{tab, k.km.Enter, k.km.Esc}
}

func (k editFormKeyMap) FullHelp() [][]key.Binding {
	tab := key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field"))
	shiftTab := key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev field"))
	return [][]key.Binding{
		{tab, shiftTab},
		{k.km.Enter, k.km.Esc},
	}
}

// ── Task detail (no active timer on this task) ────────────────────────────────

type taskDetailKeyMap struct{ km KeyMap }

func (k taskDetailKeyMap) ShortHelp() []key.Binding {
	start := key.NewBinding(key.WithKeys("s", "enter"), key.WithHelp("s/enter", "start timer"))
	add := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "add session"))
	edit := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit session"))
	del := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))
	return []key.Binding{k.km.Up, k.km.Down, start, add, edit, del, k.km.Esc, k.km.Help}
}

func (k taskDetailKeyMap) FullHelp() [][]key.Binding {
	start := key.NewBinding(key.WithKeys("s", "enter"), key.WithHelp("s/enter", "start timer"))
	add := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "add session"))
	edit := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit session"))
	del := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))
	return [][]key.Binding{
		{k.km.Up, k.km.Down, start},
		{add, edit, del},
		{k.km.Esc, k.km.Help},
	}
}

// ── Task detail (timer active on this task) ───────────────────────────────────

type taskDetailActiveKeyMap struct{ km KeyMap }

func (k taskDetailActiveKeyMap) ShortHelp() []key.Binding {
	add := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "add session"))
	edit := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit session"))
	del := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))
	return []key.Binding{k.km.Up, k.km.Down, k.km.Start, k.km.Stop, k.km.Reset, add, edit, del, k.km.Esc}
}

func (k taskDetailActiveKeyMap) FullHelp() [][]key.Binding {
	add := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "add session"))
	edit := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit session"))
	del := key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))
	return [][]key.Binding{
		{k.km.Start, k.km.Stop, k.km.Reset},
		{k.km.Up, k.km.Down, add, edit, del},
		{k.km.Esc, k.km.Help},
	}
}

// ── Theme picker ──────────────────────────────────────────────────────────────

type themePickerKeyMap struct{ km KeyMap }

func (k themePickerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.km.Up, k.km.Down, k.km.Enter, k.km.Esc}
}

func (k themePickerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.km.Up, k.km.Down},
		{k.km.Enter, k.km.Esc},
	}
}

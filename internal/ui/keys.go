package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all application-wide keybindings.
type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Enter      key.Binding
	Esc        key.Binding
	NewProject key.Binding
	NewTask    key.Binding
	Delete     key.Binding
	Edit       key.Binding
	Start      key.Binding
	Pause      key.Binding
	Stop       key.Binding
	Reset      key.Binding
	Report     key.Binding
	Quit       key.Binding
	Help       key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:       key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "collapse")),
		Right:      key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "expand")),
		Enter:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select/expand")),
		Esc:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		NewProject: key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "new project")),
		NewTask:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new task")),
		Delete:     key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		Edit:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit sessions")),
		Start:      key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start timer")),
		Pause:      key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "pause/resume")),
		Stop:       key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "stop & save")),
		Reset:      key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset")),
		Report:     key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "report")),
		Quit:       key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more keys")),
	}
}

// treeKeyMap implements help.KeyMap for the tree view.
type treeKeyMap struct{ km KeyMap }

func (k treeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.km.Up, k.km.Down,
		k.km.NewProject, k.km.NewTask,
		k.km.Delete,
		k.km.Enter,
		k.km.Report,
		k.km.Quit,
		k.km.Help,
	}
}

func (k treeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.km.Up, k.km.Down, k.km.Left, k.km.Right},
		{k.km.NewProject, k.km.NewTask, k.km.Delete},
		{k.km.Enter, k.km.Edit, k.km.Report},
		{k.km.Quit, k.km.Help},
	}
}

// timerKeyMap implements help.KeyMap for the timer view.
type timerKeyMap struct{ km KeyMap }

func (k timerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.km.Pause, k.km.Stop, k.km.Reset, k.km.Esc}
}

func (k timerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.km.Pause, k.km.Stop},
		{k.km.Reset, k.km.Esc},
	}
}

// reportKeyMap implements help.KeyMap for the report view.
type reportKeyMap struct{ km KeyMap }

func (k reportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.km.Up, k.km.Down, k.km.Esc, k.km.Help}
}

func (k reportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.km.Up, k.km.Down},
		{k.km.Esc, k.km.Help},
	}
}

// inputKeyMap is shown while the tree's text prompt is active.
type inputKeyMap struct{ km KeyMap }

func (k inputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.km.Enter, k.km.Esc}
}

func (k inputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.km.Enter, k.km.Esc}}
}

// editKeyMap implements help.KeyMap for the session edit view (normal mode).
// Inline bindings override the global KeyMap labels to be session-specific.
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

// editFormKeyMap is shown while the add/edit session form is active.
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

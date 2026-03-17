package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all application-wide keybindings.
type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Enter   key.Binding
	Esc     key.Binding
	New     key.Binding
	Delete  key.Binding
	Edit    key.Binding
	Start   key.Binding
	Pause   key.Binding
	Stop    key.Binding
	Reset   key.Binding
	Report  key.Binding
	Quit    key.Binding
	Help    key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Left:   key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "collapse")),
		Right:  key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "expand")),
		Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Esc:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Start:  key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start")),
		Pause:  key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "pause")),
		Stop:   key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "stop")),
		Reset:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset")),
		Report: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "report")),
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

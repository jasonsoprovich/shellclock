package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/ui"
)

func main() {
	model.RunBackup()

	store, err := model.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "shellclock: failed to load data: %v\n", err)
		os.Exit(1)
	}

	app := ui.New(store)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "shellclock: %v\n", err)
		os.Exit(1)
	}
}

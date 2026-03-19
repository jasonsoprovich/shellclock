package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jasonsoprovich/shellclock/internal/importer"
	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/ui"
)

func main() {
	// Handle CLI subcommands before launching the TUI.
	if len(os.Args) >= 2 && os.Args[1] == "import" {
		runImport(os.Args[2:])
		return
	}

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

func runImport(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: shellclock import toggl /path/to/export.csv")
		os.Exit(1)
	}
	source := args[0]
	csvPath := args[1]

	if source != "toggl" {
		fmt.Fprintf(os.Stderr, "shellclock: unknown import source %q (supported: toggl)\n", source)
		os.Exit(1)
	}

	store, err := model.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "shellclock: failed to load data: %v\n", err)
		os.Exit(1)
	}

	projects, tasks, sessions, err := importer.ImportToggl(store, csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shellclock: import failed: %v\n", err)
		os.Exit(1)
	}

	if err := store.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "shellclock: failed to save data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d projects imported, %d tasks imported, %d sessions imported\n", projects, tasks, sessions)
}

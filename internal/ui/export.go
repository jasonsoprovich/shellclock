package ui

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

// reportsDir returns ~/.config/shellclock/reports, creating it if needed.
func reportsDir() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cfg, "shellclock", "reports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// exportCSV writes a CSV report to ~/.config/shellclock/reports/ and returns the full path.
// Columns: Project, Task, Total Duration (h:mm:ss), Total Sessions
func exportCSV(store *model.Store) (string, error) {
	dir, err := reportsDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("shellclock-report-%s.csv", time.Now().Format("2006-01-02")))
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	_ = w.Write([]string{"Project", "Task", "Total Duration (h:mm:ss)", "Total Sessions"})
	for _, p := range store.Projects {
		for _, t := range p.Tasks {
			_ = w.Write([]string{
				p.Name,
				t.Name,
				util.FormatDurationShort(t.TotalSeconds()),
				strconv.Itoa(len(t.Sessions)),
			})
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return path, nil
}

// exportText writes a plain-text summary to ~/.config/shellclock/reports/ and returns the full path.
func exportText(store *model.Store) (string, error) {
	dir, err := reportsDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("shellclock-report-%s.txt", time.Now().Format("2006-01-02")))

	var grandTotal int64
	for i := range store.Projects {
		grandTotal += store.Projects[i].TotalSeconds()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("shellclock Report — %s\n", time.Now().Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("Total: %s\n", util.FormatDuration(grandTotal)))
	sb.WriteString(strings.Repeat("─", 50) + "\n\n")

	for _, p := range store.Projects {
		sb.WriteString(fmt.Sprintf("▸ %s  %s\n", p.Name, util.FormatDuration(p.TotalSeconds())))
		for _, t := range p.Tasks {
			noun := "sessions"
			if len(t.Sessions) == 1 {
				noun = "session"
			}
			sb.WriteString(fmt.Sprintf("  · %s  %s  (%d %s)\n",
				t.Name, util.FormatDuration(t.TotalSeconds()), len(t.Sessions), noun))
		}
		sb.WriteString("\n")
	}

	return path, os.WriteFile(path, []byte(sb.String()), 0o644)
}

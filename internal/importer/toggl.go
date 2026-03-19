package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jasonsoprovich/shellclock/internal/model"
)

// ImportToggl reads a Toggl CSV export at csvPath and merges it into store.
// Toggl columns used: Project, Task, Description, Start date, Start time,
// End date, End time. Returns counts of new projects, tasks, and sessions added.
func ImportToggl(store *model.Store, csvPath string) (newProjects, newTasks, newSessions int, err error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("open %s: %w", csvPath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("read CSV header: %w", err)
	}

	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.TrimSpace(h)] = i
	}

	for _, col := range []string{"Project", "Start date", "Start time", "End date", "End time"} {
		if _, ok := idx[col]; !ok {
			return 0, 0, 0, fmt.Errorf("CSV is missing required column %q — is this a Toggl export?", col)
		}
	}

	get := func(row []string, name string) string {
		i, ok := idx[name]
		if !ok || i >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[i])
	}

	// Build lookup maps seeded from existing data so we never duplicate.
	projectIDs := make(map[string]string) // project name → project ID
	taskIDs := make(map[string]string)    // projectID+":"+taskName → task ID
	for _, p := range store.Projects {
		projectIDs[p.Name] = p.ID
		for _, t := range p.Tasks {
			taskIDs[p.ID+":"+t.Name] = t.ID
		}
	}

	for {
		row, rerr := r.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return newProjects, newTasks, newSessions, fmt.Errorf("read CSV row: %w", rerr)
		}

		projectName := get(row, "Project")
		if projectName == "" {
			continue
		}

		// Prefer Toggl "Task" as the shellclock task name; fall back to Description.
		taskName := get(row, "Task")
		if taskName == "" {
			taskName = get(row, "Description")
		}
		if taskName == "" {
			taskName = "(untitled)"
		}

		// Parse start/end timestamps. Toggl exports "YYYY-MM-DD HH:MM:SS".
		startStr := get(row, "Start date") + " " + get(row, "Start time")
		endStr := get(row, "End date") + " " + get(row, "End time")

		start, perr := time.ParseInLocation("2006-01-02 15:04:05", startStr, time.Local)
		if perr != nil {
			continue // skip rows with unparseable times
		}
		end, perr := time.ParseInLocation("2006-01-02 15:04:05", endStr, time.Local)
		if perr != nil {
			continue
		}
		if !end.After(start) {
			continue // skip zero-length or inverted sessions
		}

		// Find or create project.
		projectID, exists := projectIDs[projectName]
		if !exists {
			p := store.AddProject(projectName)
			projectID = p.ID
			projectIDs[projectName] = projectID
			newProjects++
		}

		// Find or create task.
		taskKey := projectID + ":" + taskName
		taskID, exists := taskIDs[taskKey]
		if !exists {
			t := store.AddTask(projectID, taskName)
			taskID = t.ID
			taskIDs[taskKey] = taskID
			newTasks++
		}

		// Store Description as session notes when it differs from the task name.
		notes := get(row, "Description")
		if notes == taskName {
			notes = ""
		}
		if len(notes) > 120 {
			notes = notes[:120]
		}

		store.AddSession(projectID, taskID, model.Session{
			Start:           start,
			End:             end,
			DurationSeconds: int64(end.Sub(start).Seconds()),
			Notes:           notes,
		})
		newSessions++
	}

	return newProjects, newTasks, newSessions, nil
}

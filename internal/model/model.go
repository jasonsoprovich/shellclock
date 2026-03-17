package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Session represents a single timed work interval.
type Session struct {
	ID              string    `json:"id"`
	Start           time.Time `json:"start"`
	End             time.Time `json:"end,omitempty"`
	DurationSeconds int64     `json:"duration_seconds"`
}

// Task belongs to a Project and holds zero or more Sessions.
type Task struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Sessions []Session `json:"sessions"`
}

// TotalSeconds returns the sum of all session durations for this task.
func (t *Task) TotalSeconds() int64 {
	var total int64
	for _, s := range t.Sessions {
		total += s.DurationSeconds
	}
	return total
}

// Project is the top-level grouping of tasks.
type Project struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

// TotalSeconds returns the sum of all task durations for this project.
func (p *Project) TotalSeconds() int64 {
	var total int64
	for i := range p.Tasks {
		total += p.Tasks[i].TotalSeconds()
	}
	return total
}

// ActiveTimer persists an in-progress timer across restarts.
type ActiveTimer struct {
	ProjectID string    `json:"project_id"`
	TaskID    string    `json:"task_id"`
	SessionID string    `json:"session_id"`
	Start     time.Time `json:"start"`
	Paused    bool      `json:"paused"`
	// AccumulatedSeconds holds time banked before a pause.
	AccumulatedSeconds int64 `json:"accumulated_seconds"`
}

// Store is the root data structure serialised to disk.
type Store struct {
	Projects    []Project    `json:"projects"`
	ActiveTimer *ActiveTimer `json:"active_timer,omitempty"`

	path string
}

func dataPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("config dir: %w", err)
	}
	return filepath.Join(dir, "shellclock", "shellclock.json"), nil
}

// Load reads the JSON store from disk, creating it if absent.
func Load() (*Store, error) {
	p, err := dataPath()
	if err != nil {
		return nil, err
	}

	s := &Store{path: p}

	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read store: %w", err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("parse store: %w", err)
	}
	return s, nil
}

// Save writes the store to disk atomically.
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

// --- Project helpers ---

func (s *Store) AddProject(name string) *Project {
	p := Project{ID: uuid.NewString(), Name: name}
	s.Projects = append(s.Projects, p)
	return &s.Projects[len(s.Projects)-1]
}

func (s *Store) FindProject(id string) *Project {
	for i := range s.Projects {
		if s.Projects[i].ID == id {
			return &s.Projects[i]
		}
	}
	return nil
}

func (s *Store) DeleteProject(id string) {
	for i, p := range s.Projects {
		if p.ID == id {
			s.Projects = append(s.Projects[:i], s.Projects[i+1:]...)
			return
		}
	}
}

// --- Task helpers ---

func (s *Store) AddTask(projectID, name string) *Task {
	p := s.FindProject(projectID)
	if p == nil {
		return nil
	}
	t := Task{ID: uuid.NewString(), Name: name}
	p.Tasks = append(p.Tasks, t)
	return &p.Tasks[len(p.Tasks)-1]
}

func (s *Store) FindTask(projectID, taskID string) *Task {
	p := s.FindProject(projectID)
	if p == nil {
		return nil
	}
	for i := range p.Tasks {
		if p.Tasks[i].ID == taskID {
			return &p.Tasks[i]
		}
	}
	return nil
}

func (s *Store) DeleteTask(projectID, taskID string) {
	p := s.FindProject(projectID)
	if p == nil {
		return
	}
	for i, t := range p.Tasks {
		if t.ID == taskID {
			p.Tasks = append(p.Tasks[:i], p.Tasks[i+1:]...)
			return
		}
	}
}

// --- Session helpers ---

func (s *Store) AddSession(projectID, taskID string, sess Session) {
	t := s.FindTask(projectID, taskID)
	if t == nil {
		return
	}
	if sess.ID == "" {
		sess.ID = uuid.NewString()
	}
	t.Sessions = append(t.Sessions, sess)
}

func (s *Store) DeleteSession(projectID, taskID, sessionID string) {
	t := s.FindTask(projectID, taskID)
	if t == nil {
		return
	}
	for i, sess := range t.Sessions {
		if sess.ID == sessionID {
			t.Sessions = append(t.Sessions[:i], t.Sessions[i+1:]...)
			return
		}
	}
}

func (s *Store) UpdateSession(projectID, taskID string, updated Session) {
	t := s.FindTask(projectID, taskID)
	if t == nil {
		return
	}
	for i, sess := range t.Sessions {
		if sess.ID == updated.ID {
			t.Sessions[i] = updated
			return
		}
	}
}

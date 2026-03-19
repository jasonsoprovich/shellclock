package model

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const maxBackups = 7

// BackupDir returns the path to the backup directory.
func BackupDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "shellclock", "backups"), nil
}

// RunBackup copies the current data file into the backup directory under
// today's date (shellclock-YYYY-MM-DD.json). Only one backup is created per
// calendar day. Older backups are pruned so at most maxBackups are kept.
// All errors are silently ignored — backup is best-effort.
func RunBackup() {
	src, err := dataPath()
	if err != nil {
		return
	}
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return // nothing to back up on first launch
	}

	dir, err := BackupDir()
	if err != nil {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}

	today := time.Now().Format("2006-01-02")
	dest := filepath.Join(dir, "shellclock-"+today+".json")
	if _, err := os.Stat(dest); err == nil {
		return // today's backup already exists
	}

	if err := copyFile(src, dest); err != nil {
		return
	}

	pruneBackups(dir)
}

// ListBackups returns backup file names, newest first.
func ListBackups() ([]string, error) {
	dir, err := BackupDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && isBackupFile(e.Name()) {
			names = append(names, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names, nil
}

// RestoreFromBackup reads the named backup file (filename only, not a full
// path), unmarshals it, applies its contents to the live store, and saves the
// result back to the data file. The store's path and Theme are preserved when
// the backup contains no theme entry.
func RestoreFromBackup(filename string, store *Store) error {
	dir, err := BackupDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	var tmp Store
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	store.Projects = tmp.Projects
	store.ActiveTimer = tmp.ActiveTimer
	if tmp.Theme != "" {
		store.Theme = tmp.Theme
	}
	return store.Save()
}

func isBackupFile(name string) bool {
	return strings.HasPrefix(name, "shellclock-") && strings.HasSuffix(name, ".json")
}

func pruneBackups(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && isBackupFile(e.Name()) {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names) // oldest first (YYYY-MM-DD lexicographic order)
	for len(names) > maxBackups {
		_ = os.Remove(filepath.Join(dir, names[0]))
		names = names[1:]
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

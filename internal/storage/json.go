// Package storage – JSON file-backed implementation of the Storage interface.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iamminhquan/gotodo/internal/task"
)

// JSONStorage persists tasks as a JSON array in a single file.
// Writes are atomic: data is first written to a temp file, then renamed
// over the real file so a crash during write never corrupts the data.
type JSONStorage struct {
	filePath string
	tasks    []*task.Task
}

// NewJSONStorage opens (or creates) the JSON store at filePath.
func NewJSONStorage(filePath string) (*JSONStorage, error) {
	// Ensure the parent directory exists.
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating data directory: %w", err)
	}

	s := &JSONStorage{filePath: filePath}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// load reads all tasks from disk. If the file doesn't exist yet, we start with
// an empty slice – that is not an error.
func (s *JSONStorage) load() error {
	data, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) {
		s.tasks = []*task.Task{}
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading tasks file: %w", err)
	}

	// Handle an empty file gracefully.
	if len(data) == 0 {
		s.tasks = []*task.Task{}
		return nil
	}

	if err := json.Unmarshal(data, &s.tasks); err != nil {
		return fmt.Errorf("parsing tasks file: %w", err)
	}
	return nil
}

// save writes all tasks to disk atomically.
func (s *JSONStorage) save() error {
	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("serialising tasks: %w", err)
	}

	// Write to a sibling temp file first.
	dir := filepath.Dir(s.filePath)
	tmp, err := os.CreateTemp(dir, ".gotodo-*.tmp")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpName := tmp.Name()

	// Ensure temp file is cleaned up on any failure path.
	defer func() {
		if _, statErr := os.Stat(tmpName); statErr == nil {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	// Atomic rename – on POSIX this is guaranteed; on Windows it is best-effort.
	if err := os.Rename(tmpName, s.filePath); err != nil {
		return fmt.Errorf("replacing tasks file: %w", err)
	}
	return nil
}

// findIndex returns the slice index of the task whose ID starts with the given
// string (allows short-ID prefix matching).  Returns -1 if not found.
func (s *JSONStorage) findIndex(id string) int {
	for i, t := range s.tasks {
		if t.ID == id || strings.HasPrefix(t.ID, id) {
			return i
		}
	}
	return -1
}

// ── Storage interface implementation ──────────────────────────────────────────

// Add appends a new task and saves.
func (s *JSONStorage) Add(t *task.Task) error {
	s.tasks = append(s.tasks, t)
	return s.save()
}

// List returns tasks that satisfy the filter.
func (s *JSONStorage) List(filter TaskFilter) ([]*task.Task, error) {
	var out []*task.Task
	for _, t := range s.tasks {
		if !matchesFilter(t, filter) {
			continue
		}
		out = append(out, t)
	}
	return out, nil
}

// matchesFilter checks whether a single task should be included.
func matchesFilter(t *task.Task, f TaskFilter) bool {
	// done / pending gate
	if t.Done && !f.ShowDone {
		return false
	}
	if !t.Done && !f.ShowPending {
		return false
	}

	// Priority filter
	if f.Priority != "" && t.Priority != f.Priority {
		return false
	}

	// Tag filter
	if f.Tag != "" {
		found := false
		for _, tag := range t.Tags {
			if strings.EqualFold(tag, f.Tag) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Due-today filter
	if f.DueToday && !t.IsDueToday() {
		return false
	}

	return true
}

// GetByID retrieves a task by exact ID or ID prefix.
func (s *JSONStorage) GetByID(id string) (*task.Task, error) {
	idx := s.findIndex(id)
	if idx < 0 {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	// Return a copy so callers cannot inadvertently mutate the in-memory slice.
	copy := *s.tasks[idx]
	return &copy, nil
}

// MarkDone toggles the done flag and saves.
func (s *JSONStorage) MarkDone(id string, done bool) error {
	idx := s.findIndex(id)
	if idx < 0 {
		return fmt.Errorf("task not found: %s", id)
	}
	if done {
		s.tasks[idx].MarkDone()
	} else {
		s.tasks[idx].MarkUndone()
	}
	return s.save()
}

// Delete removes a task by ID / prefix and saves.
func (s *JSONStorage) Delete(id string) error {
	idx := s.findIndex(id)
	if idx < 0 {
		return fmt.Errorf("task not found: %s", id)
	}
	s.tasks = append(s.tasks[:idx], s.tasks[idx+1:]...)
	return s.save()
}

// Update replaces the stored task that shares the same ID and saves.
func (s *JSONStorage) Update(t *task.Task) error {
	idx := s.findIndex(t.ID)
	if idx < 0 {
		return fmt.Errorf("task not found: %s", t.ID)
	}
	// Replace the record with the caller-supplied copy.
	updated := *t
	updated.ID = s.tasks[idx].ID // guard against accidental ID change
	s.tasks[idx] = &updated
	return s.save()
}

// Close is a no-op for the JSON backend (data is saved after every mutation).
func (s *JSONStorage) Close() error { return nil }

// ── Diagnostics ───────────────────────────────────────────────────────────────

// FilePath returns the underlying file path (useful for config display).
func (s *JSONStorage) FilePath() string { return s.filePath }

// Count returns the total number of stored tasks.
func (s *JSONStorage) Count() int { return len(s.tasks) }

// ── Test helpers ──────────────────────────────────────────────────────────────

// NewInMemoryJSONStorage creates a JSONStorage backed by a temp file – useful
// in tests where you don't want to pollute the real data directory.
func NewInMemoryJSONStorage() (*JSONStorage, error) {
	tmp, err := os.CreateTemp("", "gotodo-test-*.json")
	if err != nil {
		return nil, err
	}
	// Write an empty JSON array so load() doesn't fail.
	if _, err := tmp.WriteString("[]"); err != nil {
		_ = tmp.Close()
		return nil, err
	}
	if err := tmp.Close(); err != nil {
		return nil, err
	}

	s := &JSONStorage{filePath: tmp.Name(), tasks: []*task.Task{}}
	return s, nil
}

// ensure compile-time satisfaction of the interface.
var _ Storage = (*JSONStorage)(nil)

// parseDate is a helper exposed from this package so the CLI layer can reuse
// it for parsing --due flags.
func ParseDateString(s string) (time.Time, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	switch s {
	case "today":
		return today, nil
	case "tomorrow":
		return today.AddDate(0, 0, 1), nil
	case "next week", "nextweek":
		return today.AddDate(0, 0, 7), nil
	}

	// Try common date formats.
	formats := []string{
		"2006-01-02",
		"02-01-2006",
		"01/02/2006",
		"2006/01/02",
		"Jan 2 2006",
		"2 Jan 2006",
	}
	for _, f := range formats {
		if t, err := time.ParseInLocation(f, s, now.Location()); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date %q; try YYYY-MM-DD, 'today', or 'tomorrow'", s)
}

// Package storage defines the Storage interface used by all CLI commands.
// Keeping the interface here makes it straightforward to swap the JSON backend
// for SQLite, PostgreSQL, or any other data store without touching the CLI layer.
package storage

import "github.com/iamminhquan/gotodo/internal/task"

// TaskFilter controls which tasks are returned by List.
type TaskFilter struct {
	// ShowDone includes completed tasks when true.
	ShowDone bool
	// ShowPending includes pending tasks when true.
	ShowPending bool
	// Priority filters by this value when non-empty.
	Priority string
	// Tag filters by this tag when non-empty.
	Tag string
	// DueToday filters tasks that are due today when true.
	DueToday bool
}

// AllFilter is a convenience filter that includes every task.
var AllFilter = TaskFilter{ShowDone: true, ShowPending: true}

// PendingFilter is a convenience filter that includes only pending tasks.
var PendingFilter = TaskFilter{ShowPending: true}

// DoneFilter is a convenience filter that includes only completed tasks.
var DoneFilter = TaskFilter{ShowDone: true}

// Storage is the repository interface.  All CLI commands depend only on this
// interface, not on any concrete implementation.
type Storage interface {
	// Add persists a new task.
	Add(t *task.Task) error

	// List returns tasks that match the given filter.
	List(filter TaskFilter) ([]*task.Task, error)

	// GetByID retrieves a single task by its full or prefix ID.
	// Returns an error if the task is not found.
	GetByID(id string) (*task.Task, error)

	// MarkDone toggles the done state for the task identified by id.
	MarkDone(id string, done bool) error

	// Delete permanently removes a task.
	Delete(id string) error

	// Update replaces an existing task record (matched by ID).
	Update(t *task.Task) error

	// Close flushes and releases any held resources.
	Close() error
}

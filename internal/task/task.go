// Package task defines the core Task data model and related helpers.
package task

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Priority levels supported by gotodo.
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// validPriorities is the set of accepted priority strings.
var validPriorities = map[string]bool{
	PriorityHigh:   true,
	PriorityMedium: true,
	PriorityLow:    true,
}

// defaultTagPool is the pool of tags to randomly assign when the user provides
// none.  The slice order is deterministic so seeded rand gives stable results.
var defaultTagPool = []string{
	"work", "todo", "personal", "urgent",
	"home", "learning", "health", "finance",
}

// tagRng is a package-level random source seeded at init time.
// Using a dedicated source avoids mutating the global math/rand default and
// keeps behaviour reproducible in tests (call SeedTagRng).
var tagRng = rand.New(rand.NewSource(time.Now().UnixNano()))

// SeedTagRng replaces the random source with a deterministic seed – useful in
// tests that need repeatable default-tag selection.
func SeedTagRng(seed int64) {
	tagRng = rand.New(rand.NewSource(seed))
}

// randomDefaultTags picks 1–2 distinct tags from the default pool.
func randomDefaultTags() []string {
	count := 1 + tagRng.Intn(2) // 1 or 2 tags
	perm := tagRng.Perm(len(defaultTagPool))
	out := make([]string, count)
	for i := 0; i < count; i++ {
		out[i] = defaultTagPool[perm[i]]
	}
	return out
}

// Task represents a single to-do item.
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Priority    string     `json:"priority"`
	Due         *time.Time `json:"due,omitempty"`
	Done        bool       `json:"done"`
	Tags        []string   `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// New creates a new Task with a generated UUID and sane defaults.
// If tags is nil or empty, 1–2 random default tags are assigned.
func New(title, priority string, due *time.Time, tags []string) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	priority = strings.ToLower(strings.TrimSpace(priority))
	if priority == "" {
		priority = PriorityMedium
	}
	if !validPriorities[priority] {
		return nil, errors.New("priority must be one of: high, medium, low")
	}

	// Normalise tags: trim spaces, lowercase, deduplicate.
	tags = normaliseTags(tags)

	// Assign random default tags when the user provides none.
	if len(tags) == 0 {
		tags = randomDefaultTags()
	}

	return &Task{
		ID:        uuid.New().String(),
		Title:     title,
		Priority:  priority,
		Due:       due,
		Done:      false,
		Tags:      tags,
		CreatedAt: time.Now(),
	}, nil
}

// Validate checks the task for consistency errors.
func (t *Task) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if !validPriorities[t.Priority] {
		return errors.New("priority must be one of: high, medium, low")
	}
	return nil
}

// PriorityEmoji returns an emoji that visually represents the priority level.
func (t *Task) PriorityEmoji() string {
	switch t.Priority {
	case PriorityHigh:
		return "🔴"
	case PriorityMedium:
		return "🟡"
	case PriorityLow:
		return "🔵"
	default:
		return "⚪"
	}
}

// StatusEmoji returns an emoji reflecting the done/pending state.
func (t *Task) StatusEmoji() string {
	if t.Done {
		return "✅"
	}
	return "⬜"
}

// IsOverdue reports whether the task has a due date that is in the past and is
// not yet completed.
func (t *Task) IsOverdue() bool {
	if t.Done || t.Due == nil {
		return false
	}
	return time.Now().After(*t.Due)
}

// IsDueToday reports whether the task is due on the current calendar day.
func (t *Task) IsDueToday() bool {
	if t.Due == nil {
		return false
	}
	now := time.Now()
	y1, m1, d1 := now.Date()
	y2, m2, d2 := t.Due.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// MarkDone sets the task as done and records the completion timestamp.
func (t *Task) MarkDone() {
	t.Done = true
	now := time.Now()
	t.CompletedAt = &now
}

// MarkUndone clears the done state and the completion timestamp.
func (t *Task) MarkUndone() {
	t.Done = false
	t.CompletedAt = nil
}

// normaliseTags deduplicates and lowercases a slice of tag strings.
func normaliseTags(tags []string) []string {
	seen := make(map[string]bool, len(tags))
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		out = append(out, tag)
	}
	return out
}

// ParsePriority is a convenience wrapper that validates and normalises a
// priority string.  Returns an error if the value is not recognised.
func ParsePriority(s string) (string, error) {
	p := strings.ToLower(strings.TrimSpace(s))
	if !validPriorities[p] {
		return "", errors.New("priority must be one of: high, medium, low")
	}
	return p, nil
}

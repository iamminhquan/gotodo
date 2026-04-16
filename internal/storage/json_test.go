package storage_test

import (
	"testing"
	"time"

	"github.com/iamminhquan/gotodo/internal/storage"
	"github.com/iamminhquan/gotodo/internal/task"
)

// newTestStorage creates an in-memory (temp file) JSON storage for tests.
func newTestStorage(t *testing.T) *storage.JSONStorage {
	t.Helper()
	s, err := storage.NewInMemoryJSONStorage()
	if err != nil {
		t.Fatalf("creating test storage: %v", err)
	}
	return s
}

func TestJSONStorage_AddAndList(t *testing.T) {
	s := newTestStorage(t)

	tk, _ := task.New("Test task", "high", nil, nil)
	if err := s.Add(tk); err != nil {
		t.Fatalf("Add: %v", err)
	}

	tasks, err := s.List(storage.AllFilter)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Title != "Test task" {
		t.Errorf("title mismatch: got %q", tasks[0].Title)
	}
}

func TestJSONStorage_GetByID(t *testing.T) {
	s := newTestStorage(t)
	tk, _ := task.New("Find me", "low", nil, nil)
	_ = s.Add(tk)

	got, err := s.GetByID(tk.ID[:6]) // prefix match
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != "Find me" {
		t.Errorf("expected 'Find me', got %q", got.Title)
	}
}

func TestJSONStorage_MarkDone(t *testing.T) {
	s := newTestStorage(t)
	tk, _ := task.New("Complete me", "medium", nil, nil)
	_ = s.Add(tk)

	if err := s.MarkDone(tk.ID, true); err != nil {
		t.Fatalf("MarkDone: %v", err)
	}

	got, _ := s.GetByID(tk.ID)
	if !got.Done {
		t.Error("expected task to be done")
	}

	// Undo.
	_ = s.MarkDone(tk.ID, false)
	got, _ = s.GetByID(tk.ID)
	if got.Done {
		t.Error("expected task to be pending after undo")
	}
}

func TestJSONStorage_Delete(t *testing.T) {
	s := newTestStorage(t)
	tk, _ := task.New("Delete me", "low", nil, nil)
	_ = s.Add(tk)

	if err := s.Delete(tk.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	tasks, _ := s.List(storage.AllFilter)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks after delete, got %d", len(tasks))
	}
}

func TestJSONStorage_Update(t *testing.T) {
	s := newTestStorage(t)
	tk, _ := task.New("Original", "low", nil, nil)
	_ = s.Add(tk)

	tk.Title = "Updated"
	tk.Priority = task.PriorityHigh
	if err := s.Update(tk); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := s.GetByID(tk.ID)
	if got.Title != "Updated" {
		t.Errorf("expected 'Updated', got %q", got.Title)
	}
	if got.Priority != task.PriorityHigh {
		t.Errorf("expected high priority, got %q", got.Priority)
	}
}

func TestJSONStorage_FilterByPriority(t *testing.T) {
	s := newTestStorage(t)
	t1, _ := task.New("High task", "high", nil, nil)
	t2, _ := task.New("Low task", "low", nil, nil)
	_ = s.Add(t1)
	_ = s.Add(t2)

	f := storage.TaskFilter{ShowPending: true, Priority: "high"}
	tasks, _ := s.List(f)
	if len(tasks) != 1 || tasks[0].Priority != "high" {
		t.Errorf("expected 1 high-priority task, got %d", len(tasks))
	}
}

func TestJSONStorage_FilterByTag(t *testing.T) {
	s := newTestStorage(t)
	t1, _ := task.New("Tagged", "medium", nil, []string{"work"})
	t2, _ := task.New("Untagged", "medium", nil, nil)
	_ = s.Add(t1)
	_ = s.Add(t2)

	f := storage.TaskFilter{ShowPending: true, Tag: "work"}
	tasks, _ := s.List(f)
	if len(tasks) != 1 {
		t.Errorf("expected 1 tagged task, got %d", len(tasks))
	}
}

func TestParseDateString(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"2026-05-01", true},
		{"today", true},
		{"tomorrow", true},
		{"next week", true},
		{"not-a-date", false},
		{"32-13-2026", false},
	}
	for _, c := range cases {
		_, err := storage.ParseDateString(c.input)
		if c.valid && err != nil {
			t.Errorf("expected %q to parse, got error: %v", c.input, err)
		}
		if !c.valid && err == nil {
			t.Errorf("expected %q to fail, but it parsed OK", c.input)
		}
	}
}

// Ensure JSONStorage implements the Storage interface at compile time.
var _ storage.Storage = (*storage.JSONStorage)(nil)

// Silence the "time imported and not used" error when the test file is compiled
// without the time-related tests.
var _ = time.Now

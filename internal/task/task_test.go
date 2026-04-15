package task_test

import (
	"testing"
	"time"

	"github.com/iamminhquan/gotodo/internal/task"
)

func TestNew_ValidTask(t *testing.T) {
	tk, err := task.New("Buy milk", "high", nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tk.Title != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %q", tk.Title)
	}
	if tk.Priority != task.PriorityHigh {
		t.Errorf("expected priority high, got %q", tk.Priority)
	}
	if tk.ID == "" {
		t.Error("expected a non-empty UUID")
	}
	if tk.Done {
		t.Error("new task should not be done")
	}
}

func TestNew_EmptyTitle(t *testing.T) {
	_, err := task.New("   ", "medium", nil, nil)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestNew_InvalidPriority(t *testing.T) {
	_, err := task.New("Some task", "urgent", nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
}

func TestNew_DefaultPriorityMedium(t *testing.T) {
	tk, err := task.New("Task", "", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tk.Priority != task.PriorityMedium {
		t.Errorf("expected medium priority as default, got %q", tk.Priority)
	}
}

func TestMarkDone_SetsCompletedAt(t *testing.T) {
	tk, _ := task.New("Task", "low", nil, nil)
	tk.MarkDone()
	if !tk.Done {
		t.Error("expected task to be done")
	}
	if tk.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestMarkUndone_ClearsState(t *testing.T) {
	tk, _ := task.New("Task", "low", nil, nil)
	tk.MarkDone()
	tk.MarkUndone()
	if tk.Done {
		t.Error("expected task to be pending")
	}
	if tk.CompletedAt != nil {
		t.Error("expected CompletedAt to be nil after undone")
	}
}

func TestIsOverdue(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	tk, _ := task.New("Overdue task", "high", &past, nil)
	if !tk.IsOverdue() {
		t.Error("expected task to be overdue")
	}
}

func TestPriorityEmoji(t *testing.T) {
	cases := []struct {
		priority string
		emoji    string
	}{
		{task.PriorityHigh, "🔴"},
		{task.PriorityMedium, "🟡"},
		{task.PriorityLow, "🔵"},
	}
	for _, c := range cases {
		tk, _ := task.New("t", c.priority, nil, nil)
		if got := tk.PriorityEmoji(); got != c.emoji {
			t.Errorf("priority %s: expected emoji %s, got %s", c.priority, c.emoji, got)
		}
	}
}

func TestParsePriority(t *testing.T) {
	for _, valid := range []string{"high", "medium", "low", "HIGH", "Medium"} {
		if _, err := task.ParsePriority(valid); err != nil {
			t.Errorf("expected %q to be valid priority, got error: %v", valid, err)
		}
	}
	if _, err := task.ParsePriority("critical"); err == nil {
		t.Error("expected error for invalid priority 'critical'")
	}
}

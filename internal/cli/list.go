package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/iamminhquan/gotodo/internal/storage"
)

var (
	listPending  bool
	listDone     bool
	listAll      bool
	listToday    bool
	listPriority string
	listTag      string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `Display your tasks in a rich table.

Filter flags (mutually exclusive if combining done/pending):
  --pending   show only pending tasks (default when no flag is given)
  --done      show only completed tasks
  --all       show every task
  --today     show tasks due today (pending only)
  --priority  filter by priority (high/medium/low)
  --tag       filter by tag

Examples:
  gotodo list
  gotodo list --all
  gotodo list --done
  gotodo list --priority high --tag work
  gotodo list --today`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := resolveFilter()

		tasks, err := state.repo.List(filter)
		if err != nil {
			return fmt.Errorf("listing tasks: %w", err)
		}

		if len(tasks) == 0 {
			emptyMsg := emptyMessage(filter)
			printNoTasks(emptyMsg)
			return nil
		}

		// Build and render the table.
		t := buildTable(tasks, state.cfg.UseColor)
		t.SetOutputMirror(os.Stdout)
		t.Render()

		// Summary line.
		fmt.Println()
		total := len(tasks)
		doneCount := 0
		for _, tk := range tasks {
			if tk.Done {
				doneCount++
			}
		}
		pending := total - doneCount
		color.New(color.Faint).Printf("  %d task(s) shown", total)
		if pending > 0 && doneCount > 0 {
			color.New(color.Faint).Printf("  ·  %d pending  ·  %d done", pending, doneCount)
		}
		fmt.Println()
		fmt.Println()

		return nil
	},
}

// resolveFilter converts CLI flags into a storage.TaskFilter.
func resolveFilter() storage.TaskFilter {
	f := storage.TaskFilter{}

	switch {
	case listAll:
		f.ShowDone = true
		f.ShowPending = true
	case listDone:
		f.ShowDone = true
	default: // --pending is the default
		f.ShowPending = true
	}

	f.DueToday = listToday
	f.Priority = listPriority
	f.Tag = listTag
	return f
}

// emptyMessage returns an appropriate no-results message based on the filter.
func emptyMessage(f storage.TaskFilter) string {
	switch {
	case f.ShowDone && !f.ShowPending:
		return "🎉 No completed tasks yet."
	case f.DueToday:
		return "📅 No tasks due today – enjoy your day!"
	case f.Priority != "":
		return fmt.Sprintf("No %s-priority tasks found.", f.Priority)
	case f.Tag != "":
		return fmt.Sprintf("No tasks tagged #%s.", f.Tag)
	default:
		return "✨ No pending tasks – all clear!"
	}
}

func init() {
	listCmd.Flags().BoolVar(&listPending, "pending", false, "Show only pending tasks (default)")
	listCmd.Flags().BoolVar(&listDone, "done", false, "Show only completed tasks")
	listCmd.Flags().BoolVar(&listAll, "all", false, "Show all tasks")
	listCmd.Flags().BoolVar(&listToday, "today", false, "Show tasks due today")
	listCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter by priority (high/medium/low)")
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter by tag")
}

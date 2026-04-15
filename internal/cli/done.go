package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var doneUndo bool

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark a task as done (or undo with --undo)",
	Long: `Mark a task as completed.  Pass --undo to revert it back to pending.

The <id> can be the full UUID or just the first few characters (prefix match).

Examples:
  gotodo done a1b2c3d4
  gotodo done a1b2c3d4 --undo`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		// Verify the task exists before changing state.
		t, err := state.repo.GetByID(id)
		if err != nil {
			return fmt.Errorf("task not found (id: %s)", id)
		}

		markDone := !doneUndo

		// Idempotency guard with a friendly message.
		if t.Done == markDone {
			if markDone {
				color.Yellow("  Task [%s] is already marked as done.\n", shortID(t.ID))
			} else {
				color.Yellow("  Task [%s] is already pending.\n", shortID(t.ID))
			}
			return nil
		}

		if err := state.repo.MarkDone(t.ID, markDone); err != nil {
			return fmt.Errorf("updating task: %w", err)
		}

		if markDone {
			color.New(color.FgGreen, color.Bold).Printf("✅ Done! ")
			fmt.Printf("[%s] %s\n", shortID(t.ID), t.Title)
		} else {
			color.New(color.FgYellow, color.Bold).Printf("↩  Undone! ")
			fmt.Printf("[%s] %s is now pending again.\n", shortID(t.ID), t.Title)
		}
		return nil
	},
}

func init() {
	doneCmd.Flags().BoolVar(&doneUndo, "undo", false, "Revert a completed task back to pending")
}

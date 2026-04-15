package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Permanently delete a task",
	Long: `Remove a task from the list forever.

You will be asked to confirm unless --force is passed.

Examples:
  gotodo delete a1b2c3d4
  gotodo delete a1b2c3d4 --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		// Verify the task exists and show its details before deletion.
		t, err := state.repo.GetByID(id)
		if err != nil {
			return fmt.Errorf("task not found (id: %s)", id)
		}

		if !deleteForce {
			fmt.Printf("\n  %s %s [%s]\n\n", t.PriorityEmoji(), t.Title, shortID(t.ID))
			if !confirmPrompt("  Delete this task permanently?") {
				color.Yellow("  Cancelled.\n")
				return nil
			}
		}

		if err := state.repo.Delete(t.ID); err != nil {
			return fmt.Errorf("deleting task: %w", err)
		}

		color.New(color.FgRed, color.Bold).Printf("🗑  Deleted ")
		fmt.Printf("[%s] %s\n", shortID(t.ID), t.Title)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation prompt")
}

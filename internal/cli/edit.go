package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/iamminhquan/gotodo/internal/storage"
	"github.com/iamminhquan/gotodo/internal/task"
)

var (
	editTitle    string
	editPriority string
	editDue      string
	editTags     string
	editClearDue bool
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing task",
	Long: `Modify one or more fields of an existing task.
Only the flags you provide will be updated; everything else stays the same.

Examples:
  gotodo edit a1b2c3d4 --title "Buy organic groceries"
  gotodo edit a1b2c3d4 --priority low --due 2026-05-15
  gotodo edit a1b2c3d4 --tags work,urgent
  gotodo edit a1b2c3d4 --clear-due`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		existing, err := state.repo.GetByID(id)
		if err != nil {
			return fmt.Errorf("task not found (id: %s)", id)
		}

		changed := false

		// ── title ────────────────────────────────────────────────────────────
		if editTitle != "" {
			trimmed := strings.TrimSpace(editTitle)
			if trimmed == "" {
				return fmt.Errorf("title cannot be empty")
			}
			existing.Title = trimmed
			changed = true
		}

		// ── priority ─────────────────────────────────────────────────────────
		if editPriority != "" {
			p, err := task.ParsePriority(editPriority)
			if err != nil {
				return err
			}
			existing.Priority = p
			changed = true
		}

		// ── due date ─────────────────────────────────────────────────────────
		if editClearDue {
			existing.Due = nil
			changed = true
		} else if editDue != "" {
			t, err := storage.ParseDateString(editDue)
			if err != nil {
				return err
			}
			existing.Due = &t
			changed = true
		}

		// ── tags ─────────────────────────────────────────────────────────────
		if cmd.Flags().Changed("tags") {
			var tags []string
			for _, tag := range strings.Split(editTags, ",") {
				if t := strings.TrimSpace(tag); t != "" {
					tags = append(tags, t)
				}
			}
			existing.Tags = tags
			changed = true
		}

		if !changed {
			color.Yellow("  Nothing to update – please provide at least one flag.\n")
			return nil
		}

		if err := existing.Validate(); err != nil {
			return err
		}

		if err := state.repo.Update(existing); err != nil {
			return fmt.Errorf("saving task: %w", err)
		}

		color.New(color.FgGreen, color.Bold).Printf("✏️  Updated! ")
		fmt.Printf("[%s] %s %s\n", shortID(existing.ID), existing.PriorityEmoji(), existing.Title)
		if existing.Due != nil {
			fmt.Printf("   📅 Due: %s\n", humanDate(*existing.Due))
		}
		if len(existing.Tags) > 0 {
			fmt.Printf("   🏷  Tags: %s\n", strings.Join(existing.Tags, ", "))
		}
		return nil
	},
}

func init() {
	editCmd.Flags().StringVarP(&editTitle, "title", "T", "", "New title")
	editCmd.Flags().StringVarP(&editPriority, "priority", "p", "", "New priority: high, medium, low")
	editCmd.Flags().StringVarP(&editDue, "due", "d", "", "New due date: YYYY-MM-DD, 'today', 'tomorrow'")
	editCmd.Flags().StringVarP(&editTags, "tags", "t", "", "New comma-separated tags (replaces existing tags)")
	editCmd.Flags().BoolVar(&editClearDue, "clear-due", false, "Remove the due date from this task")
}

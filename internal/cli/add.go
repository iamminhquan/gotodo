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
	addPriority string
	addDue      string
	addTags     string
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new task",
	Long: `Create a new task with an optional priority, due date, and tags.

Examples:
  gotodo add "Read Clean Code" --priority high --due 2026-05-01 --tags reading,learning
  gotodo add "Buy coffee" --due tomorrow
  gotodo add "Weekly review" --priority medium --tags work`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")

		// Resolve priority (fall back to config default).
		priority := addPriority
		if priority == "" {
			priority = state.cfg.DefaultPriority
		}
		if _, err := task.ParsePriority(priority); err != nil {
			return fmt.Errorf("invalid priority %q: %w", priority, err)
		}

		// Parse optional due date.
		var due *timePtr
		if addDue != "" {
			t, err := storage.ParseDateString(addDue)
			if err != nil {
				return err
			}
			due = &t
		}

		// Parse tags.
		var tags []string
		if addTags != "" {
			for _, tag := range strings.Split(addTags, ",") {
				if t := strings.TrimSpace(tag); t != "" {
					tags = append(tags, t)
				}
			}
		}

		t, err := task.New(title, priority, due, tags)
		if err != nil {
			return err
		}

		if err := state.repo.Add(t); err != nil {
			return fmt.Errorf("saving task: %w", err)
		}

		// ── friendly confirmation ────────────────────────────────────────────
		green := color.New(color.FgGreen, color.Bold)
		green.Printf("✅ Task added! ")
		fmt.Printf("[%s] %s %s\n", shortID(t.ID), t.PriorityEmoji(), t.Title)
		if t.Due != nil {
			fmt.Printf("   📅 Due: %s\n", humanDate(*t.Due))
		}
		if len(t.Tags) > 0 {
			fmt.Printf("   🏷  Tags: %s\n", strings.Join(t.Tags, ", "))
		}
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addPriority, "priority", "p", "", "Priority: high, medium, low (default: config default_priority)")
	addCmd.Flags().StringVarP(&addDue, "due", "d", "", "Due date: YYYY-MM-DD, 'today', 'tomorrow', 'next week'")
	addCmd.Flags().StringVarP(&addTags, "tags", "t", "", "Comma-separated list of tags")
}

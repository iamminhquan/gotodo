// Package cli wires all Cobra commands together and owns the shared application
// state (config + storage) that is initialised once in PersistentPreRunE.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/iamminhquan/gotodo/internal/config"
	"github.com/iamminhquan/gotodo/internal/storage"
	"github.com/iamminhquan/gotodo/internal/version"
)

// appState holds the shared objects that are initialised once and injected into
// every sub-command.  Using a struct instead of globals makes testing easier.
type appState struct {
	cfg  *config.Config
	repo storage.Storage
}

// state is the process-level singleton.  CLI commands access it after
// PersistentPreRunE has populated it.
var state = &appState{}

// NewRootCmd constructs and returns the configured root Cobra command.
// Calling this from main.go keeps the wiring testable.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gotodo",
		Short: "A fast, beautiful CLI to-do manager",
		Long: `gotodo - manage your tasks without leaving the terminal.

Examples:
  gotodo add "Buy groceries" --priority high --due tomorrow --tags personal,errands
  gotodo list --pending
  gotodo done <id>
  gotodo edit <id> --priority low
  gotodo delete <id>`,
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		// Run without a sub-command → show list (same as `gotodo list`).
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCmd.RunE(cmd, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initState()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if state.repo != nil {
				return state.repo.Close()
			}
			return nil
		},
	}

	// Version flag is already handled by cobra when Version is set.

	// Register sub-commands.
	root.AddCommand(
		addCmd,
		listCmd,
		doneCmd,
		deleteCmd,
		editCmd,
	)

	return root
}

// initState loads config and opens the storage backend.  It is idempotent –
// subsequent calls within the same process are no-ops.
func initState() error {
	if state.cfg != nil {
		return nil // already initialised
	}

	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	state.cfg = cfg

	tasksFile := config.GetTasksFilePath(cfg)
	repo, err := storage.NewJSONStorage(tasksFile)
	if err != nil {
		return fmt.Errorf("opening storage: %w", err)
	}
	state.repo = repo
	return nil
}

// Execute is the single entry point called from main.go.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

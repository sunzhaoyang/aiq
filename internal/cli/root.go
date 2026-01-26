package cli

import (
	"fmt"
	"os"

	"github.com/aiq/aiq/internal/config"
	"github.com/aiq/aiq/internal/sql"
	"github.com/aiq/aiq/internal/ui"
)

// Run starts the main CLI application
func Run() error {
	// Check for first-run and run wizard if needed
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("failed to check config: %w", err)
	}

	if !exists {
		cfg, err := config.RunWizard()
		if err != nil {
			return fmt.Errorf("configuration wizard failed: %w", err)
		}
		_ = cfg // Config is saved by wizard
	}

	// Load config to verify it's valid
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.IsEmpty() {
		ui.ShowWarning("Configuration is incomplete. Please run the wizard again.")
		cfg, err = config.RunWizard()
		if err != nil {
			return fmt.Errorf("configuration wizard failed: %w", err)
		}
		_ = cfg
	}

	// Main menu loop
	for {
		items := []ui.MenuItem{
			{Label: "config   - Manage LLM configuration", Value: "config"},
			{Label: "source   - Manage database connections", Value: "source"},
			{Label: "sql      - Query database with natural language", Value: "sql"},
			{Label: "exit     - Exit application", Value: "exit"},
		}

		choice, err := ui.ShowMenu("AIQ - Main Menu", items)
		if err != nil {
			// User cancelled (Ctrl+C)
			fmt.Println()
			ui.ShowInfo("Goodbye!")
			return nil
		}

		switch choice {
		case "config":
			if err := RunConfigMenu(); err != nil {
				ui.ShowError(err.Error())
			}
		case "source":
			if err := RunSourceMenu(); err != nil {
				ui.ShowError(err.Error())
			}
		case "sql":
			if err := sql.RunSQLMode(); err != nil {
				ui.ShowError(err.Error())
			}
		case "exit":
			ui.ShowInfo("Goodbye!")
			os.Exit(0)
		}

		fmt.Println()
	}
}

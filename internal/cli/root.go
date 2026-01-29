package cli

import (
	"fmt"
	"os"

	"github.com/aiq/aiq/internal/config"
	"github.com/aiq/aiq/internal/prompt"
	"github.com/aiq/aiq/internal/sql"
	"github.com/aiq/aiq/internal/ui"
)

// Run starts the main CLI application
// sessionFile is optional path to a session file to restore
func Run(sessionFile string) error {
	// Ensure directory structure exists
	if err := config.EnsureDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create config directory structure: %w", err)
	}

	// Initialize prompt loader to ensure default prompt files are created
	// This happens at startup so prompts are available even if user doesn't enter chat mode
	_, err := prompt.NewLoader()
	if err != nil {
		// Log warning but don't fail - prompts will use fallback defaults
		fmt.Printf("Warning: Failed to initialize prompts: %v. Using default prompts.\n", err)
	}

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
			{Label: "chat     - Query database with natural language", Value: "chat"},
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
		case "chat":
			if err := sql.RunSQLMode(sessionFile); err != nil {
				// Check if error is ErrReturnToMenu - this is expected and means return to main menu
				if err == sql.ErrReturnToMenu {
					// Normal return to menu, don't show error
				} else {
					ui.ShowError(err.Error())
				}
			}
			// Clear sessionFile after first use (only restore once)
			sessionFile = ""
		case "exit":
			ui.ShowInfo("Goodbye!")
			os.Exit(0)
		}

		fmt.Println()
	}
}

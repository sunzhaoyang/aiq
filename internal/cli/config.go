package cli

import (
	"fmt"
	"strings"

	"github.com/aiq/aiq/internal/config"
	"github.com/aiq/aiq/internal/ui"
)

// RunConfigMenu runs the configuration management menu
func RunConfigMenu() error {
	for {
		items := []ui.MenuItem{
			{Label: "view    - View current configuration", Value: "view"},
			{Label: "url     - Update LLM API URL", Value: "update_url"},
			{Label: "model   - Update model name", Value: "update_model"},
			{Label: "key     - Update LLM API key", Value: "update_key"},
			{Label: "back    - Back to main menu", Value: "back"},
		}

		choice, err := ui.ShowMenu("Configuration", items)
		if err != nil {
			return err
		}

		switch choice {
		case "view":
			if err := viewConfig(); err != nil {
				ui.ShowError(err.Error())
			}
		case "update_url":
			if err := updateURL(); err != nil {
				ui.ShowError(err.Error())
			} else {
				ui.ShowSuccess("LLM URL updated successfully!")
			}
		case "update_model":
			if err := updateModel(); err != nil {
				ui.ShowError(err.Error())
			} else {
				ui.ShowSuccess("Model name updated successfully!")
			}
		case "update_key":
			if err := updateAPIKey(); err != nil {
				ui.ShowError(err.Error())
			} else {
				ui.ShowSuccess("API Key updated successfully!")
			}
		case "back":
			return nil
		}

		fmt.Println()
	}
}

func viewConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println()
	ui.ShowInfo("Current Configuration:")
	fmt.Printf("  LLM URL: %s\n", cfg.LLM.URL)
	fmt.Printf("  Model: %s\n", cfg.LLM.Model)
	fmt.Printf("  API Key: %s\n", maskAPIKey(cfg.LLM.APIKey))
	fmt.Println()

	return nil
}

func updateURL() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println()
	fmt.Println("LLM API URL Format:")
	fmt.Println("  Enter the base URL of your LLM API endpoint.")
	fmt.Println("  Examples:")
	fmt.Println("    - https://api.openai.com/v1")
	fmt.Println("    - https://api.anthropic.com/v1")
	fmt.Println("    - https://api.example.com/v1")
	fmt.Println()
	fmt.Println("  Note: The '/chat/completions' path will be added automatically.")
	fmt.Println()

	newURL, err := ui.ShowInput("Enter new LLM URL", cfg.LLM.URL)
	if err != nil {
		return fmt.Errorf("failed to get URL: %w", err)
	}

	cfg.LLM.URL = newURL

	if err := config.ValidatePartialLLMConfig(&cfg.LLM); err != nil {
		return err
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

func updateModel() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println()
	fmt.Println("Model Name:")
	fmt.Println("  Enter the model name to use for SQL translation.")
	fmt.Println("  Examples:")
	fmt.Println("    - gpt-3.5-turbo")
	fmt.Println("    - gpt-4")
	fmt.Println("    - claude-3-opus")
	fmt.Println("    - deepseek-chat")
	fmt.Println()

	newModel, err := ui.ShowInput("Enter Model Name", cfg.LLM.Model)
	if err != nil {
		return fmt.Errorf("failed to get model name: %w", err)
	}

	cfg.LLM.Model = strings.TrimSpace(newModel)
	if cfg.LLM.Model == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

func updateAPIKey() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	newKey, err := ui.ShowPassword("Enter new API Key")
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	cfg.LLM.APIKey = newKey

	if err := config.ValidatePartialLLMConfig(&cfg.LLM); err != nil {
		return err
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

func maskAPIKey(key string) string {
	if len(key) == 0 {
		return ""
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

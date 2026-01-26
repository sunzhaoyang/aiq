package config

import (
	"fmt"

	"github.com/aiq/aiq/internal/ui"
)

// RunWizard runs the first-run configuration wizard
func RunWizard() (*Config, error) {
	ui.ShowInfo("Welcome to AIQ! Let's set up your configuration.")
	fmt.Println()

	config := NewConfig()

	// Get LLM URL with format hint
	fmt.Println("LLM API URL Format:")
	fmt.Println("  Enter the base URL of your LLM API endpoint.")
	fmt.Println("  Examples:")
	fmt.Println("    - https://api.openai.com/v1")
	fmt.Println("    - https://api.anthropic.com/v1")
	fmt.Println("    - https://api.example.com/v1")
	fmt.Println()
	fmt.Println("  Note: The '/chat/completions' path will be added automatically.")
	fmt.Println()
	
	url, err := ui.ShowInput("Enter LLM API URL", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM URL: %w", err)
	}
	config.LLM.URL = url

	// Get Model Name
	fmt.Println()
	fmt.Println("Model Name:")
	fmt.Println("  Enter the model name to use for SQL translation.")
	fmt.Println("  Examples:")
	fmt.Println("    - gpt-3.5-turbo")
	fmt.Println("    - gpt-4")
	fmt.Println("    - claude-3-opus")
	fmt.Println("    - deepseek-chat")
	fmt.Println()
	
	model, err := ui.ShowInput("Enter Model Name", "gpt-3.5-turbo")
	if err != nil {
		return nil, fmt.Errorf("failed to get model name: %w", err)
	}
	config.LLM.Model = model

	// Get API Key
	fmt.Println()
	apiKey, err := ui.ShowPassword("Enter LLM API Key")
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}
	config.LLM.APIKey = apiKey

	// Validate configuration
	if err := Validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Save configuration
	if err := Save(config); err != nil {
		return nil, fmt.Errorf("failed to save configuration: %w", err)
	}

	ui.ShowSuccess("Configuration saved successfully!")
	fmt.Println()

	return config, nil
}

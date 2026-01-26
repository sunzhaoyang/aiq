package config

import (
	"fmt"
	"net/url"
	"strings"
)

// Validate validates the configuration
func Validate(config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	if err := ValidateLLMConfig(&config.LLM); err != nil {
		return fmt.Errorf("LLM config validation failed: %w", err)
	}

	return nil
}

// ValidateLLMConfig validates LLM configuration
func ValidateLLMConfig(llm *LLMConfig) error {
	if llm == nil {
		return fmt.Errorf("LLM config is nil")
	}

	// Validate URL
	if llm.URL == "" {
		return fmt.Errorf("LLM URL is required")
	}

	parsedURL, err := url.Parse(llm.URL)
	if err != nil {
		return fmt.Errorf("invalid LLM URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("LLM URL must use http or https scheme")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("LLM URL must have a host")
	}

	// Validate API Key
	if llm.APIKey == "" {
		return fmt.Errorf("LLM API key is required")
	}

	if strings.TrimSpace(llm.APIKey) == "" {
		return fmt.Errorf("LLM API key cannot be empty")
	}

	// Validate Model
	if llm.Model == "" {
		return fmt.Errorf("LLM model name is required")
	}

	if strings.TrimSpace(llm.Model) == "" {
		return fmt.Errorf("LLM model name cannot be empty")
	}

	return nil
}

// ValidatePartialLLMConfig validates LLM configuration allowing empty values
func ValidatePartialLLMConfig(llm *LLMConfig) error {
	if llm == nil {
		return fmt.Errorf("LLM config is nil")
	}

	// If URL is provided, validate it
	if llm.URL != "" {
		parsedURL, err := url.Parse(llm.URL)
		if err != nil {
			return fmt.Errorf("invalid LLM URL format: %w", err)
		}

		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return fmt.Errorf("LLM URL must use http or https scheme")
		}

		if parsedURL.Host == "" {
			return fmt.Errorf("LLM URL must have a host")
		}
	}

	return nil
}

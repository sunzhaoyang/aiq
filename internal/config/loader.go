package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// GetConfigPath returns the full path to the configuration file
// Deprecated: Use GetConfigFilePath() instead
func GetConfigPath() (string, error) {
	return GetConfigFilePath()
}

// Load loads configuration from file
func Load() (*Config, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return NewConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves configuration to file
func Save(config *Config) error {
	// Ensure directory structure exists
	if err := EnsureDirectoryStructure(); err != nil {
		return err
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists checks if configuration file exists
func Exists() (bool, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

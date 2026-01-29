package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// ConfigDir is the base configuration directory name
	ConfigDir = ".aiq"

	// Subdirectories within config directory
	ConfigSubdir   = "config"
	SessionsSubdir = "sessions"
	SkillsSubdir   = "skills"
	ToolsSubdir    = "tools"
	PromptsSubdir  = "prompts"
	BinSubdir      = "bin"

	// Config files
	ConfigFile  = "config.yaml"
	SourcesFile = "sources.yaml"
)

// GetBaseConfigDir returns the base configuration directory path (~/.aiq)
func GetBaseConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ConfigDir), nil
}

// GetConfigDir returns the config subdirectory path (~/.aiq/config)
func GetConfigDir() (string, error) {
	baseDir, err := GetBaseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, ConfigSubdir), nil
}

// GetSessionsDir returns the sessions subdirectory path (~/.aiq/sessions)
func GetSessionsDir() (string, error) {
	baseDir, err := GetBaseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, SessionsSubdir), nil
}

// GetSkillsDir returns the skills subdirectory path (~/.aiq/skills)
func GetSkillsDir() (string, error) {
	baseDir, err := GetBaseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, SkillsSubdir), nil
}

// GetToolsDir returns the tools subdirectory path (~/.aiq/tools)
func GetToolsDir() (string, error) {
	baseDir, err := GetBaseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, ToolsSubdir), nil
}

// GetPromptsDir returns the prompts subdirectory path (~/.aiq/prompts)
func GetPromptsDir() (string, error) {
	baseDir, err := GetBaseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, PromptsSubdir), nil
}

// GetBinDir returns the bin subdirectory path (~/.aiq/bin)
func GetBinDir() (string, error) {
	baseDir, err := GetBaseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, BinSubdir), nil
}

// GetConfigFilePath returns the full path to the configuration file (~/.aiq/config/config.yaml)
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, ConfigFile), nil
}

// GetSourcesFilePath returns the full path to the sources file (~/.aiq/config/sources.yaml)
func GetSourcesFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, SourcesFile), nil
}

// EnsureDirectoryStructure creates all required subdirectories if they don't exist
func EnsureDirectoryStructure() error {
	dirs := []struct {
		name string
		get  func() (string, error)
	}{
		{"config", GetConfigDir},
		{"sessions", GetSessionsDir},
		{"skills", GetSkillsDir},
		{"tools", GetToolsDir},
		{"prompts", GetPromptsDir},
		{"bin", GetBinDir},
	}

	for _, dir := range dirs {
		path, err := dir.get()
		if err != nil {
			return fmt.Errorf("failed to get %s directory path: %w", dir.name, err)
		}

		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dir.name, err)
		}
	}

	return nil
}

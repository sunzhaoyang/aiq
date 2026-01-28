package source

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/aiq/aiq/internal/config"
)

// GetSourcesPath returns the full path to the sources file
func GetSourcesPath() (string, error) {
	return config.GetSourcesFilePath()
}

// LoadSources loads all sources from file
func LoadSources() ([]*Source, error) {
	sourcesPath, err := GetSourcesPath()
	if err != nil {
		return nil, err
	}

	// Check if sources file exists
	if _, err := os.Stat(sourcesPath); os.IsNotExist(err) {
		return []*Source{}, nil
	}

	data, err := os.ReadFile(sourcesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sources file: %w", err)
	}

	var sources []*Source
	if err := yaml.Unmarshal(data, &sources); err != nil {
		return nil, fmt.Errorf("failed to parse sources file: %w", err)
	}

	return sources, nil
}

// SaveSources saves sources to file
func SaveSources(sources []*Source) error {
	// Ensure directory structure exists
	if err := config.EnsureDirectoryStructure(); err != nil {
		return err
	}

	sourcesPath, err := GetSourcesPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(sources)
	if err != nil {
		return fmt.Errorf("failed to marshal sources: %w", err)
	}

	if err := os.WriteFile(sourcesPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write sources file: %w", err)
	}

	return nil
}

// AddSource adds a new source
func AddSource(source *Source) error {
	sources, err := LoadSources()
	if err != nil {
		return err
	}

	// Check if source with same name already exists
	for _, s := range sources {
		if s.Name == source.Name {
			return fmt.Errorf("source with name '%s' already exists", source.Name)
		}
	}

	sources = append(sources, source)
	return SaveSources(sources)
}

// RemoveSource removes a source by name
func RemoveSource(name string) error {
	sources, err := LoadSources()
	if err != nil {
		return err
	}

	found := false
	newSources := make([]*Source, 0, len(sources))
	for _, s := range sources {
		if s.Name != name {
			newSources = append(newSources, s)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("source with name '%s' not found", name)
	}

	return SaveSources(newSources)
}

// GetSource returns a source by name
func GetSource(name string) (*Source, error) {
	sources, err := LoadSources()
	if err != nil {
		return nil, err
	}

	for _, s := range sources {
		if s.Name == name {
			return s, nil
		}
	}

	return nil, fmt.Errorf("source with name '%s' not found", name)
}

// GenerateUniqueSourceName generates a unique source name based on host, port, and user
// Format: {host}-{port}-{user}, with numeric suffix if collision occurs
func GenerateUniqueSourceName(host string, port int, user string) (string, error) {
	baseName := fmt.Sprintf("%s-%d-%s", host, port, user)
	
	sources, err := LoadSources()
	if err != nil {
		return "", fmt.Errorf("failed to load sources: %w", err)
	}

	// Check if base name exists
	nameExists := false
	for _, s := range sources {
		if s.Name == baseName {
			nameExists = true
			break
		}
	}

	if !nameExists {
		return baseName, nil
	}

	// Find unique name by appending numeric suffix
	for i := 2; i < 1000; i++ {
		candidateName := fmt.Sprintf("%s-%d", baseName, i)
		exists := false
		for _, s := range sources {
			if s.Name == candidateName {
				exists = true
				break
			}
		}
		if !exists {
			return candidateName, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique source name after 998 attempts")
}

// FindExistingSourceByConnection finds an existing source with the same host, port, and username
// Returns the source name if found, empty string if not found
func FindExistingSourceByConnection(host string, port int, username string) (string, error) {
	sources, err := LoadSources()
	if err != nil {
		return "", fmt.Errorf("failed to load sources: %w", err)
	}

	for _, s := range sources {
		if s.Host == host && s.Port == port && s.Username == username {
			return s.Name, nil
		}
	}

	return "", nil // Not found, but no error
}

// AddSourceWithAutoName adds a source with an auto-generated unique name
// If a source with the same host, port, and username already exists, returns the existing source name
func AddSourceWithAutoName(source *Source) (string, error) {
	// First check if a source with the same connection parameters already exists
	existingName, err := FindExistingSourceByConnection(source.Host, source.Port, source.Username)
	if err != nil {
		return "", err
	}
	if existingName != "" {
		// Source already exists, return existing name
		return existingName, nil
	}

	// No existing source found, create a new one with auto-generated name
	name, err := GenerateUniqueSourceName(source.Host, source.Port, source.Username)
	if err != nil {
		return "", err
	}
	source.Name = name
	return name, AddSource(source)
}

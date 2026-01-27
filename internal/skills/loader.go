package skills

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aiq/aiq/internal/config"
)

// LoadSkillsMetadata loads metadata from all Skills in the skills directory
// This is called on startup for progressive loading
func LoadSkillsMetadata() ([]*Metadata, error) {
	skillsDir, err := config.GetSkillsDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get skills directory: %w", err)
	}

	// Check if skills directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		// Skills directory doesn't exist yet, return empty list
		return []*Metadata{}, nil
	}

	var metadataList []*Metadata

	// Walk through skill directories
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()
		skillPath := filepath.Join(skillsDir, skillName, "SKILL.md")

		// Check if SKILL.md exists
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			// Skip directories without SKILL.md
			continue
		}

		// Read only frontmatter for progressive loading (efficient)
		metadata, err := LoadSkillMetadataOnly(skillPath)
		if err != nil {
			// Log error but continue with other skills
			log.Printf("Warning: Failed to load metadata for skill '%s': %v", skillName, err)
			continue
		}

		metadata.Path = skillPath
		metadataList = append(metadataList, metadata)
	}

	return metadataList, nil
}

// LoadSkillMetadataOnly reads only the frontmatter section of a Skill file
// This is more efficient for progressive loading - only reads until closing ---
func LoadSkillMetadataOnly(skillPath string) (*Metadata, error) {
	file, err := os.Open(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open skill file: %w", err)
	}
	defer file.Close()

	// Read file line by line until we find the closing ---
	var frontmatterLines []string
	foundFirstDelimiter := false
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if !foundFirstDelimiter {
			// Looking for first ---
			if strings.TrimSpace(line) == "---" {
				foundFirstDelimiter = true
				continue
			}
			// If file doesn't start with ---, it's invalid
			return nil, fmt.Errorf("skill file must start with '---' delimiter")
		}
		
		// Check if this is the closing ---
		if strings.TrimSpace(line) == "---" {
			// Found closing delimiter, stop reading
			break
		}
		
		// Collect frontmatter lines
		frontmatterLines = append(frontmatterLines, line)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read skill file: %w", err)
	}
	
	if !foundFirstDelimiter {
		return nil, fmt.Errorf("missing frontmatter delimiter")
	}
	
	if len(frontmatterLines) == 0 {
		return nil, fmt.Errorf("empty frontmatter")
	}
	
	// Parse metadata from frontmatter
	frontmatter := strings.Join(frontmatterLines, "\n")
	metadata, err := ParseSkillMetadata(frontmatter)
	if err != nil {
		return nil, err
	}
	
	return metadata, nil
}

// LoadSkillContent loads the full content of a Skill from its file path
func LoadSkillContent(metadata *Metadata) (*Skill, error) {
	content, err := os.ReadFile(metadata.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill file: %w", err)
	}

	skill, err := ParseSkillFile(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse skill file: %w", err)
	}

	return skill, nil
}

// LoadSkillContentFromPath loads a Skill from a file path
func LoadSkillContentFromPath(path string) (*Skill, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill file: %w", err)
	}

	skill, err := ParseSkillFile(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse skill file: %w", err)
	}

	return skill, nil
}

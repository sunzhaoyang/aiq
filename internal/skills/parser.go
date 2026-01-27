package skills

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseSkillFile parses a SKILL.md file with YAML frontmatter and Markdown content
func ParseSkillFile(content string) (*Skill, error) {
	// Split frontmatter and markdown content
	frontmatter, markdownContent, err := splitFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to split frontmatter: %w", err)
	}

	// Parse YAML frontmatter
	var metadata struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	}

	if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	// Validate required fields
	if metadata.Name == "" {
		return nil, fmt.Errorf("missing required field: name")
	}
	if metadata.Description == "" {
		return nil, fmt.Errorf("missing required field: description")
	}

	skill := &Skill{
		Name:        metadata.Name,
		Description: metadata.Description,
		Content:     strings.TrimSpace(markdownContent),
		Priority:    PriorityInactive,
		Loaded:      true, // Content is loaded when parsed
	}

	return skill, nil
}

// ParseSkillMetadata parses only the YAML frontmatter to extract metadata
// This is used for progressive loading (load metadata first, content later)
// Input can be either full file content or just frontmatter string
func ParseSkillMetadata(content string) (*Metadata, error) {
	// Check if content contains frontmatter delimiters (full file)
	// or is just frontmatter (from LoadSkillMetadataOnly)
	var frontmatter string
	if strings.Contains(content, "---") {
		// Full file content, need to split
		var err error
		frontmatter, _, err = splitFrontmatter(content)
		if err != nil {
			return nil, fmt.Errorf("failed to split frontmatter: %w", err)
		}
	} else {
		// Already just frontmatter
		frontmatter = content
	}

	var metadata struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	}

	if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	if metadata.Name == "" {
		return nil, fmt.Errorf("missing required field: name")
	}
	if metadata.Description == "" {
		return nil, fmt.Errorf("missing required field: description")
	}

	return &Metadata{
		Name:        metadata.Name,
		Description: metadata.Description,
	}, nil
}

// splitFrontmatter splits content into YAML frontmatter and Markdown body
// Frontmatter is between --- delimiters at the start of the file
func splitFrontmatter(content string) (string, string, error) {
	content = strings.TrimSpace(content)

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---") {
		return "", content, fmt.Errorf("missing YAML frontmatter delimiter")
	}

	// Find the end of frontmatter (second ---)
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return "", "", fmt.Errorf("invalid frontmatter format")
	}

	// First line should be "---"
	if strings.TrimSpace(lines[0]) != "---" {
		return "", "", fmt.Errorf("invalid frontmatter format: first line must be '---'")
	}

	// Find the closing "---"
	var frontmatterLines []string
	var markdownStart int
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			markdownStart = i + 1
			break
		}
		frontmatterLines = append(frontmatterLines, lines[i])
	}

	if markdownStart == 0 {
		return "", "", fmt.Errorf("missing closing frontmatter delimiter")
	}

	frontmatter := strings.Join(frontmatterLines, "\n")
	markdownContent := strings.Join(lines[markdownStart:], "\n")

	return frontmatter, markdownContent, nil
}

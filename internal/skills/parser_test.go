package skills

import (
	"testing"
)

func TestParseSkillFile_ValidSkill(t *testing.T) {
	content := `---
name: test-skill
description: A test skill for testing
---

# Test Skill

This is test content.
`

	skill, err := ParseSkillFile(content)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if skill.Name != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%s'", skill.Name)
	}

	if skill.Description != "A test skill for testing" {
		t.Errorf("Expected description 'A test skill for testing', got '%s'", skill.Description)
	}

	expectedContent := "# Test Skill\n\nThis is test content."
	if skill.Content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, skill.Content)
	}

	if skill.Priority != PriorityInactive {
		t.Errorf("Expected priority PriorityInactive, got %v", skill.Priority)
	}

	if !skill.Loaded {
		t.Error("Expected Loaded to be true")
	}
}

func TestParseSkillFile_InvalidYAML(t *testing.T) {
	content := `---
name: test-skill
description: [invalid yaml
---

# Test Skill
`

	_, err := ParseSkillFile(content)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	if !contains(err.Error(), "failed to parse YAML frontmatter") {
		t.Errorf("Expected YAML parsing error, got: %v", err)
	}
}

func TestParseSkillFile_MissingName(t *testing.T) {
	content := `---
description: A test skill
---

# Test Skill
`

	_, err := ParseSkillFile(content)
	if err == nil {
		t.Error("Expected error for missing name, got nil")
	}

	if !contains(err.Error(), "missing required field: name") {
		t.Errorf("Expected missing name error, got: %v", err)
	}
}

func TestParseSkillFile_MissingDescription(t *testing.T) {
	content := `---
name: test-skill
---

# Test Skill
`

	_, err := ParseSkillFile(content)
	if err == nil {
		t.Error("Expected error for missing description, got nil")
	}

	if !contains(err.Error(), "missing required field: description") {
		t.Errorf("Expected missing description error, got: %v", err)
	}
}

func TestParseSkillFile_MissingFrontmatterDelimiter(t *testing.T) {
	content := `name: test-skill
description: A test skill

# Test Skill
`

	_, err := ParseSkillFile(content)
	if err == nil {
		t.Error("Expected error for missing frontmatter delimiter, got nil")
	}

	if !contains(err.Error(), "missing YAML frontmatter delimiter") {
		t.Errorf("Expected missing delimiter error, got: %v", err)
	}
}

func TestParseSkillFile_MissingClosingDelimiter(t *testing.T) {
	content := `---
name: test-skill
description: A test skill

# Test Skill
`

	_, err := ParseSkillFile(content)
	if err == nil {
		t.Error("Expected error for missing closing delimiter, got nil")
	}

	if !contains(err.Error(), "missing closing frontmatter delimiter") {
		t.Errorf("Expected missing closing delimiter error, got: %v", err)
	}
}

func TestParseSkillMetadata_ValidMetadata(t *testing.T) {
	frontmatter := `name: test-skill
description: A test skill for testing
`

	metadata, err := ParseSkillMetadata(frontmatter)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if metadata.Name != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%s'", metadata.Name)
	}

	if metadata.Description != "A test skill for testing" {
		t.Errorf("Expected description 'A test skill for testing', got '%s'", metadata.Description)
	}
}

func TestParseSkillMetadata_FromFullFile(t *testing.T) {
	content := `---
name: test-skill
description: A test skill
---

# Test Skill

This is content.
`

	metadata, err := ParseSkillMetadata(content)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if metadata.Name != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%s'", metadata.Name)
	}

	if metadata.Description != "A test skill" {
		t.Errorf("Expected description 'A test skill', got '%s'", metadata.Description)
	}
}

func TestParseSkillMetadata_InvalidYAML(t *testing.T) {
	frontmatter := `name: test-skill
description: [invalid yaml
`

	_, err := ParseSkillMetadata(frontmatter)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	if !contains(err.Error(), "failed to parse YAML frontmatter") {
		t.Errorf("Expected YAML parsing error, got: %v", err)
	}
}

func TestParseSkillMetadata_MissingName(t *testing.T) {
	frontmatter := `description: A test skill
`

	_, err := ParseSkillMetadata(frontmatter)
	if err == nil {
		t.Error("Expected error for missing name, got nil")
	}

	if !contains(err.Error(), "missing required field: name") {
		t.Errorf("Expected missing name error, got: %v", err)
	}
}

func TestParseSkillMetadata_MissingDescription(t *testing.T) {
	frontmatter := `name: test-skill
`

	_, err := ParseSkillMetadata(frontmatter)
	if err == nil {
		t.Error("Expected error for missing description, got nil")
	}

	if !contains(err.Error(), "missing required field: description") {
		t.Errorf("Expected missing description error, got: %v", err)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

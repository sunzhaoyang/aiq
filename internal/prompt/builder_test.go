package prompt

import (
	"strings"
	"testing"

	"github.com/aiq/aiq/internal/skills"
)

func TestBuilder_BuildSystemPrompt_NoSkills(t *testing.T) {
	basePrompt := "You are a helpful assistant."
	builder := NewBuilder(basePrompt)
	
	result := builder.BuildSystemPrompt([]*skills.Skill{})
	
	if result != basePrompt {
		t.Errorf("Expected base prompt only, got: %s", result)
	}
}

func TestBuilder_BuildSystemPrompt_WithSkills(t *testing.T) {
	basePrompt := "You are a helpful assistant."
	builder := NewBuilder(basePrompt)
	
	skillsList := []*skills.Skill{
		{
			Name:        "test-skill",
			Description: "A test skill",
			Content:     "# Test Content\n\nThis is test content.",
			Priority:    skills.PriorityActive,
			Loaded:      true,
		},
	}
	
	result := builder.BuildSystemPrompt(skillsList)
	
	// Check that base prompt is included
	if !strings.Contains(result, basePrompt) {
		t.Error("Base prompt should be included")
	}
	
	// Check that Skills section is included
	if !strings.Contains(result, "## Available Skills (Context Information)") {
		t.Error("Skills section should be included")
	}
	
	// Check that skill name is included
	if !strings.Contains(result, "test-skill") {
		t.Error("Skill name should be included")
	}
	
	// Check that skill description is included
	if !strings.Contains(result, "A test skill") {
		t.Error("Skill description should be included")
	}
	
	// Check that skill content is included
	if !strings.Contains(result, "Test Content") {
		t.Error("Skill content should be included")
	}
}

func TestBuilder_BuildSystemPrompt_MultipleSkills(t *testing.T) {
	basePrompt := "You are a helpful assistant."
	builder := NewBuilder(basePrompt)
	
	skillsList := []*skills.Skill{
		{
			Name:        "skill1",
			Description: "First skill",
			Content:     "Content 1",
			Priority:    skills.PriorityActive,
			Loaded:      true,
		},
		{
			Name:        "skill2",
			Description: "Second skill",
			Content:     "Content 2",
			Priority:    skills.PriorityRelevant,
			Loaded:      true,
		},
	}
	
	result := builder.BuildSystemPrompt(skillsList)
	
	// Check that both skills are included
	if !strings.Contains(result, "skill1") {
		t.Error("First skill should be included")
	}
	
	if !strings.Contains(result, "skill2") {
		t.Error("Second skill should be included")
	}
	
	// Count occurrences of skill section header
	skillHeaderCount := strings.Count(result, "### Skill:")
	if skillHeaderCount != 2 {
		t.Errorf("Expected 2 skill headers, got %d", skillHeaderCount)
	}
}

func TestBuilder_FormatSkill(t *testing.T) {
	builder := NewBuilder("Base prompt")
	
	skill := &skills.Skill{
		Name:        "test-skill",
		Description: "A test skill",
		Content:     "# Test Content\n\nThis is test content.",
		Priority:    skills.PriorityActive,
		Loaded:      true,
	}
	
	result := builder.formatSkill(skill)
	
	// Check format components
	if !strings.Contains(result, "### Skill: test-skill") {
		t.Error("Skill name should be in header")
	}
	
	if !strings.Contains(result, "(Context Only - NOT a Tool)") {
		t.Error("Tool warning should be included")
	}
	
	if !strings.Contains(result, "**Description**: A test skill") {
		t.Error("Description should be formatted correctly")
	}
	
	if !strings.Contains(result, "# Test Content") {
		t.Error("Content should be included")
	}
}

func TestBuilder_FormatSkill_EmptyContent(t *testing.T) {
	builder := NewBuilder("Base prompt")
	
	skill := &skills.Skill{
		Name:        "test-skill",
		Description: "A test skill",
		Content:     "",
		Priority:    skills.PriorityActive,
		Loaded:      true,
	}
	
	result := builder.formatSkill(skill)
	
	// Should still include name and description
	if !strings.Contains(result, "test-skill") {
		t.Error("Skill name should be included even with empty content")
	}
	
	if !strings.Contains(result, "A test skill") {
		t.Error("Skill description should be included even with empty content")
	}
}

func TestBuilder_BuildFullPrompt_NoHistory(t *testing.T) {
	builder := NewBuilder("Base prompt")
	
	systemPrompt := "System prompt"
	currentQuery := "What is the weather?"
	
	result := builder.BuildFullPrompt(systemPrompt, []string{}, currentQuery)
	
	if !strings.Contains(result, systemPrompt) {
		t.Error("System prompt should be included")
	}
	
	if !strings.Contains(result, "## Current Query") {
		t.Error("Current query section should be included")
	}
	
	if !strings.Contains(result, currentQuery) {
		t.Error("Current query should be included")
	}
	
	// Should not contain history section
	if strings.Contains(result, "## Conversation History") {
		t.Error("History section should not be included when empty")
	}
}

func TestBuilder_BuildFullPrompt_WithHistory(t *testing.T) {
	builder := NewBuilder("Base prompt")
	
	systemPrompt := "System prompt"
	history := []string{
		"user: Hello",
		"assistant: Hi there!",
	}
	currentQuery := "What is the weather?"
	
	result := builder.BuildFullPrompt(systemPrompt, history, currentQuery)
	
	// Check all components
	if !strings.Contains(result, systemPrompt) {
		t.Error("System prompt should be included")
	}
	
	if !strings.Contains(result, "## Conversation History") {
		t.Error("History section should be included")
	}
	
	if !strings.Contains(result, "Hello") {
		t.Error("History content should be included")
	}
	
	if !strings.Contains(result, "## Current Query") {
		t.Error("Current query section should be included")
	}
	
	if !strings.Contains(result, currentQuery) {
		t.Error("Current query should be included")
	}
}

func TestBuilder_BuildFullPrompt_MultipleHistoryMessages(t *testing.T) {
	builder := NewBuilder("Base prompt")
	
	systemPrompt := "System prompt"
	history := []string{
		"user: First message",
		"assistant: First response",
		"user: Second message",
		"assistant: Second response",
	}
	currentQuery := "Third query"
	
	result := builder.BuildFullPrompt(systemPrompt, history, currentQuery)
	
	// Check that all history messages are included
	if !strings.Contains(result, "First message") {
		t.Error("First history message should be included")
	}
	
	if !strings.Contains(result, "Second message") {
		t.Error("Second history message should be included")
	}
	
	// Check that history messages are separated
	historySection := strings.Split(result, "## Conversation History\n")[1]
	historySection = strings.Split(historySection, "\n\n## Current Query")[0]
	
	// Should contain separators between messages
	if strings.Count(historySection, "\n\n") < len(history)-1 {
		t.Error("History messages should be properly separated")
	}
}

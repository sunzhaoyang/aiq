package prompt

import (
	"fmt"
	"strings"

	"github.com/aiq/aiq/internal/skills"
)

// Builder builds prompts with Skills content integration
type Builder struct {
	systemPromptBase string
}

// NewBuilder creates a new prompt builder
func NewBuilder(systemPromptBase string) *Builder {
	return &Builder{
		systemPromptBase: systemPromptBase,
	}
}

// BuildSystemPrompt builds the system prompt with Skills content
func (b *Builder) BuildSystemPrompt(loadedSkills []*skills.Skill) string {
	var parts []string

	// Base system prompt
	parts = append(parts, b.systemPromptBase)

	// Add Skills content
	if len(loadedSkills) > 0 {
		parts = append(parts, "\n\n## Available Skills (Context Information)\n")
		parts = append(parts, "IMPORTANT: Skills are NOT tools. They provide context and guidance on how to use the available tools.\n")
		parts = append(parts, "Do NOT call Skills as tools. Use the built-in tools (execute_sql, execute_command, http_request, file_operations, etc.) to accomplish tasks described in Skills.\n\n")
		parts = append(parts, "CRITICAL: When Skills contain bash commands or code blocks:\n")
		parts = append(parts, "- If the user requests an ACTION (like 'install', 'setup', 'configure', 'run'), you MUST EXECUTE the commands using execute_command tool, NOT just show them to the user.\n")
		parts = append(parts, "- Execute commands step by step as described in the Skill, checking results before proceeding.\n")
		parts = append(parts, "- Only show commands to the user if they explicitly ask for instructions or if execution fails and you need to explain.\n")
		parts = append(parts, "- For installation/setup tasks: Execute the commands automatically, don't just provide instructions.\n\n")
		for _, skill := range loadedSkills {
			skillSection := b.formatSkill(skill)
			parts = append(parts, skillSection)
		}
	}

	return strings.Join(parts, "")
}

// formatSkill formats a Skill for inclusion in the prompt
// Format: ## Skill: <name>\n<description>\n\n<content>
func (b *Builder) formatSkill(skill *skills.Skill) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("### Skill: %s (Context Only - NOT a Tool)", skill.Name))
	parts = append(parts, fmt.Sprintf("**Description**: %s", skill.Description))
	parts = append(parts, "**Note**: This is context information. Use the available tools to accomplish tasks described below.\n")
	if skill.Content != "" {
		parts = append(parts, skill.Content)
	}

	return strings.Join(parts, "\n") + "\n\n"
}

// BuildFullPrompt builds a complete prompt with system prompt, conversation history, and current query
func (b *Builder) BuildFullPrompt(systemPrompt string, conversationHistory []string, currentQuery string) string {
	var parts []string

	parts = append(parts, systemPrompt)

	if len(conversationHistory) > 0 {
		parts = append(parts, "\n\n## Conversation History\n")
		parts = append(parts, strings.Join(conversationHistory, "\n\n"))
	}

	parts = append(parts, "\n\n## Current Query\n")
	parts = append(parts, currentQuery)

	return strings.Join(parts, "")
}

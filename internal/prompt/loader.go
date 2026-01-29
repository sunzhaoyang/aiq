package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aiq/aiq/internal/config"
)

const (
	// Base prompt file names (using .md format with YAML frontmatter)
	FreeModeBasePromptFile = "free-mode-base.md"
	DatabaseBasePromptFile = "database-base.md"
	CommonPromptFile       = "common.md"
	
	// Database-specific prompt patch files (optional, appended to database-base.md)
	MySQLPatchFile     = "mysql.md"
	PostgreSQLPatchFile = "postgresql.md"
	SeekDBPatchFile    = "seekdb.md"
)

// Loader manages loading and initialization of prompt templates
type Loader struct {
	promptsDir string
	prompts    map[string]string
}

// NewLoader creates a new prompt loader
func NewLoader() (*Loader, error) {
	promptsDir, err := config.GetPromptsDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get prompts directory: %w", err)
	}

	loader := &Loader{
		promptsDir: promptsDir,
		prompts:    make(map[string]string),
	}

	// Ensure prompts directory exists
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create prompts directory: %w", err)
	}

	// Initialize default prompts if they don't exist
	if err := loader.initializeDefaults(); err != nil {
		return nil, fmt.Errorf("failed to initialize default prompts: %w", err)
	}

	// Load prompts from files
	if err := loader.loadPrompts(); err != nil {
		return nil, fmt.Errorf("failed to load prompts: %w", err)
	}

	return loader, nil
}

// initializeDefaults creates default prompt files if they don't exist
func (l *Loader) initializeDefaults() error {
	// Default free mode base prompt with YAML frontmatter
	freeModePrompt := `---
description: "Free mode base prompt - used when AIQ runs without database connection"
usage: "Applied when user enters chat mode without selecting a database source"
version: "1.1"
---

<MODE>
FREE MODE - No database connection available. SQL execution is not available.
</MODE>

<ROLE>
You are a helpful AI assistant. You can have natural conversations and help with system operations using available tools.
</ROLE>

<TOOLS>
- execute_command: System operations (install, setup, configuration). Not for database queries.
- http_request: Make HTTP requests.
- file_operations: Read/write files.
</TOOLS>

<POLICY>
- If the user asks for database operations, explain that no database is connected and ask whether they want to select a source.
- Do not guess database commands or run mysql/psql in free mode.
- If the request is ambiguous for the current mode, ask a clarifying question before acting.
</POLICY>
`

	// Default database mode base prompt (generic, no database-specific syntax)
	databaseBasePrompt := `---
description: "Database mode base prompt - generic database query assistant"
usage: "Applied when user enters chat mode with a selected database source. Database-specific syntax patches are loaded separately."
version: "2.0"
placeholders:
  - name: "{{DATABASE_TYPE}}"
    description: "Will be replaced with the database engine type (MySQL, PostgreSQL, etc.)"
  - name: "{{SCHEMA_CONTEXT}}"
    description: "Will be replaced with database schema information"
---

<MODE>
DATABASE MODE - Connected to a database.
</MODE>

<ROLE>
You are a helpful AI assistant for database queries and related tasks.
</ROLE>

<CONTEXT>
- Database engine type: {{DATABASE_TYPE}}
- Database connection and schema information:
{{SCHEMA_CONTEXT}}
</CONTEXT>

<POLICY>
- Use execute_sql for database queries. Do not use execute_command to run mysql/psql.
- Respect engine-specific syntax differences. Database-specific syntax guidance is provided in separate sections.
- If a request is not a database query, use the appropriate non-SQL tools.
- When unsure about syntax, rely on schema context or ask a clarifying question.
</POLICY>

<TOOLS>
- execute_sql: Execute SQL queries against the database.
- render_table: Format query results as a table.
- render_chart: Format query results as a chart when user explicitly asks for visualization.
- execute_command: System operations (install, setup, configuration). Not for database queries.
- http_request: Make HTTP requests.
- file_operations: Read/write files.
</TOOLS>
`

	// MySQL-specific syntax patch
	mysqlPatch := `---
description: "MySQL-specific syntax guidance patch"
usage: "Appended to database-base.md when database type is MySQL or seekdb"
version: "2.0"
---

<MYSQL_SYNTAX>
- Use SHOW TABLES; or SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE();
- Use DATABASE() function to get current database name.
- Schema name in WHERE table_schema should be the actual database name, not the engine type.
</MYSQL_SYNTAX>
`

	// PostgreSQL-specific syntax patch
	postgresqlPatch := `---
description: "PostgreSQL-specific syntax guidance patch"
usage: "Appended to database-base.md when database type is PostgreSQL"
version: "2.0"
---

<POSTGRESQL_SYNTAX>
- Use SELECT tablename FROM pg_tables WHERE schemaname = 'public'; or SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
- Default schema is 'public' unless otherwise specified.
- Use current_database() function to get current database name.
</POSTGRESQL_SYNTAX>
`

	// SeekDB-specific syntax patch (MySQL-compatible, but may have differences)
	seekdbPatch := `---
description: "SeekDB-specific syntax guidance patch"
usage: "Appended to database-base.md when database type is seekdb"
version: "2.0"
---

<SEEKDB_SYNTAX>
- SeekDB is MySQL-compatible, so use MySQL syntax patterns.
- Use SHOW TABLES; or SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE();
- Use DATABASE() function to get current database name.
</SEEKDB_SYNTAX>
`

	// Default common prompt with YAML frontmatter (used by both modes)
	commonPrompt := `---
description: "Common prompt section - appended to both free mode and database mode prompts"
usage: "Contains instructions that apply to both modes"
version: "1.1"
---

<EXECUTION>
- For system operations, use execute_command with explicit commands.
- If a command requires elevated privileges or interactive input, ask the user to run it manually and explain why.
- Do not fabricate command outputs. Use tool results to decide the next step.
</EXECUTION>
`

	// Create default files if they don't exist
	files := map[string]string{
		FreeModeBasePromptFile: freeModePrompt,
		DatabaseBasePromptFile: databaseBasePrompt,
		CommonPromptFile:       commonPrompt,
		MySQLPatchFile:         mysqlPatch,
		PostgreSQLPatchFile:    postgresqlPatch,
		SeekDBPatchFile:        seekdbPatch,
	}

	for filename, content := range files {
		filePath := filepath.Join(l.promptsDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create default prompt file %s: %w", filename, err)
			}
		}
	}

	return nil
}

// loadPrompts loads prompt templates from files
// Files use Markdown format with YAML frontmatter:
// - Frontmatter (between ---) contains metadata (description, usage, placeholders, etc.)
// - Body contains the actual prompt content
// - HTML comments (<!-- -->) in body are removed from final prompt
func (l *Loader) loadPrompts() error {
	// Load base prompts
	baseFiles := []string{
		FreeModeBasePromptFile,
		DatabaseBasePromptFile,
		CommonPromptFile,
	}

	for _, filename := range baseFiles {
		filePath := filepath.Join(l.promptsDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read prompt file %s: %w", filename, err)
		}

		// Parse markdown file: extract body (after frontmatter) and remove HTML comments
		promptContent, err := parsePromptFile(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse prompt file %s: %w", filename, err)
		}

		l.prompts[filename] = promptContent
	}

	// Load database-specific patches (optional, may not exist)
	patchFiles := []string{
		MySQLPatchFile,
		PostgreSQLPatchFile,
		SeekDBPatchFile,
	}

	for _, filename := range patchFiles {
		filePath := filepath.Join(l.promptsDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			// Patch files are optional, skip if not found
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to read prompt patch file %s: %w", filename, err)
		}

		// Parse markdown file: extract body (after frontmatter) and remove HTML comments
		patchContent, err := parsePromptFile(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse prompt patch file %s: %w", filename, err)
		}

		l.prompts[filename] = patchContent
	}

	return nil
}

// parsePromptFile parses a markdown prompt file with YAML frontmatter
// Returns the prompt body content with HTML comments removed
func parsePromptFile(content string) (string, error) {
	content = strings.TrimSpace(content)

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---") {
		// No frontmatter, return content as-is (backward compatibility)
		return removeHTMLComments(content), nil
	}

	// Find the end of frontmatter (second ---)
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid frontmatter format")
	}

	// First line should be "---"
	if strings.TrimSpace(lines[0]) != "---" {
		return "", fmt.Errorf("invalid frontmatter format: first line must be '---'")
	}

	// Find the closing "---"
	var markdownStart int
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			markdownStart = i + 1
			break
		}
	}

	if markdownStart == 0 {
		return "", fmt.Errorf("missing closing frontmatter delimiter")
	}

	// Extract markdown body (after frontmatter)
	markdownContent := strings.Join(lines[markdownStart:], "\n")
	
	// Remove HTML comments and return
	return removeHTMLComments(markdownContent), nil
}

// removeHTMLComments removes HTML-style comments (<!-- ... -->) from prompt content
// This allows users to add comments in prompt files without affecting the actual prompt
// HTML comments are markdown-compatible and won't conflict with markdown syntax
func removeHTMLComments(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Remove HTML comments (single-line: <!-- ... -->)
		if strings.HasPrefix(trimmed, "<!--") && strings.HasSuffix(trimmed, "-->") {
			continue // Skip comment lines
		}
		// Keep all other lines (including lines that might contain <!-- --> in the middle)
		// For simplicity, we only remove full-line comments
		result.WriteString(line)
		result.WriteString("\n")
	}
	
	return strings.TrimSpace(result.String())
}

// GetFreeModeBasePrompt returns the free mode base prompt
func (l *Loader) GetFreeModeBasePrompt() string {
	return l.prompts[FreeModeBasePromptFile]
}

// GetDatabaseModeBasePrompt returns the database mode base prompt with placeholders replaced
// and database-specific syntax patch appended if available
func (l *Loader) GetDatabaseModeBasePrompt(databaseType, schemaContext string) string {
	// Load base prompt
	prompt := l.prompts[DatabaseBasePromptFile]
	
	// Replace placeholders
	prompt = strings.ReplaceAll(prompt, "{{DATABASE_TYPE}}", databaseType)
	prompt = strings.ReplaceAll(prompt, "{{SCHEMA_CONTEXT}}", schemaContext)
	
	// Append database-specific syntax patch based on database type
	// Note: databaseType comes from Source.GetDatabaseType() which returns "MySQL", "PostgreSQL", "seekdb"
	var patchFile string
	switch strings.ToLower(databaseType) {
	case "mysql":
		patchFile = MySQLPatchFile
	case "postgresql":
		patchFile = PostgreSQLPatchFile
	case "seekdb":
		patchFile = SeekDBPatchFile
	default:
		// Default to MySQL patch for unknown types (backward compatibility)
		patchFile = MySQLPatchFile
	}
	
	// Append patch if available
	if patchFile != "" {
		if patch, exists := l.prompts[patchFile]; exists && patch != "" {
			prompt = prompt + "\n\n" + patch
		}
	}
	
	return prompt
}

// GetCommonPrompt returns the common prompt used by both modes
func (l *Loader) GetCommonPrompt() string {
	return l.prompts[CommonPromptFile]
}

// Reload reloads prompts from files (useful for testing or hot-reload scenarios)
func (l *Loader) Reload() error {
	return l.loadPrompts()
}

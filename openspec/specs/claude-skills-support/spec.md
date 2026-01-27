## ADDED Requirements

### Requirement: Skills directory structure
The system SHALL support Skills stored in `~/.aiqconfig/skills/<skill-name>/SKILL.md` format.

#### Scenario: Skills directory organization
- **WHEN** user adds a Skill
- **THEN** Skill is stored in its own subdirectory: `~/.aiqconfig/skills/<skill-name>/SKILL.md`

#### Scenario: Multiple Skills support
- **WHEN** user has multiple Skills installed
- **THEN** each Skill is stored in a separate subdirectory under `~/.aiqconfig/skills/`

### Requirement: Skills file format parsing
The system SHALL parse SKILL.md files with YAML frontmatter and Markdown content.

#### Scenario: Parse Skill metadata
- **WHEN** system loads a Skill file
- **THEN** system extracts YAML frontmatter containing `name` and `description` fields

#### Scenario: Parse Skill content
- **WHEN** system loads a Skill file
- **THEN** system extracts Markdown content after YAML frontmatter

#### Scenario: Handle invalid Skill format
- **WHEN** Skill file has invalid YAML frontmatter or missing required fields
- **THEN** system logs error and skips the Skill without crashing

### Requirement: Progressive Skills loading
The system SHALL load Skills progressively based on query relevance, not all at once.

#### Scenario: Load Skills metadata on startup
- **WHEN** application starts
- **THEN** system loads only Skills metadata (name, description) from all Skill files

#### Scenario: Match Skills to user query
- **WHEN** user submits a query
- **THEN** system matches query keywords against Skills metadata to find relevant Skills

#### Scenario: Load matched Skills content
- **WHEN** Skills are matched to user query
- **THEN** system loads full content of matched Skills (default: top 3 most relevant)

#### Scenario: Cache loaded Skills
- **WHEN** Skill content is loaded
- **THEN** system caches loaded content to avoid re-reading from disk

### Requirement: Skills matching algorithm
The system SHALL match user queries to Skills using keyword-based relevance scoring.

#### Scenario: Extract query keywords
- **WHEN** user submits a query
- **THEN** system extracts keywords from the query

#### Scenario: Score Skills by relevance
- **WHEN** system matches Skills to query
- **THEN** system scores Skills based on: exact name match > partial name match > description keyword match

#### Scenario: Select top N Skills
- **WHEN** multiple Skills match a query
- **THEN** system selects top N most relevant Skills (default: 3) for loading

### Requirement: Prompt management and compression
The system SHALL manage prompt length with token monitoring and compression strategies.

#### Scenario: Estimate token count
- **WHEN** system builds prompt for LLM
- **THEN** system estimates token count for all prompt components (system prompt, conversation history, Skills content, tools)

#### Scenario: Compress conversation history at 80% threshold
- **WHEN** estimated token count exceeds 80% of context window
- **THEN** system compresses oldest conversation messages into summary, keeping recent N messages (default: 10) in full detail

#### Scenario: Evict low-priority Skills at 90% threshold
- **WHEN** estimated token count exceeds 90% of context window after compression
- **THEN** system evicts inactive Skills (not referenced in recent queries), keeping only active and relevant Skills

#### Scenario: Aggressive compression at 95% threshold
- **WHEN** estimated token count exceeds 95% of context window
- **THEN** system summarizes all conversation history except last 5 messages and keeps only top 1 most relevant Skill

#### Scenario: Skills priority management
- **WHEN** system manages Skills in prompt
- **THEN** system prioritizes: Active (referenced in current/recent queries) > Relevant (matched but not used) > Inactive (not matched)

### Requirement: Skills integration into prompt
The system SHALL integrate Skills content into LLM system prompt dynamically.

#### Scenario: Build prompt with Skills
- **WHEN** system prepares prompt for LLM
- **THEN** system includes matched Skills content in system prompt section

#### Scenario: Format Skills content
- **WHEN** Skills content is added to prompt
- **THEN** Skills are formatted as: `## Skill: <name>\n<description>\n\n<content>`

#### Scenario: Order Skills by priority
- **WHEN** multiple Skills are included in prompt
- **THEN** Skills are ordered by priority (active first, then relevant)

### Requirement: Built-in tools for Skills
The system SHALL provide built-in tools that Skills can reference.

#### Scenario: HTTP request tool
- **WHEN** Skill needs to make HTTP requests
- **THEN** system provides `http_request` tool supporting GET, POST, PUT, DELETE with headers, query parameters, and body

#### Scenario: Command execution tool
- **WHEN** Skill needs to execute shell commands
- **THEN** system provides `execute_command` tool with timeout and working directory support, restricted to safe commands (configurable allowlist)

#### Scenario: File operations tool
- **WHEN** Skill needs file operations
- **THEN** system provides `file_operations` tool for reading, writing, listing files, restricted to user config directory and current working directory

#### Scenario: Database query tool
- **WHEN** Skill needs to query database
- **THEN** system provides `execute_sql` tool that reuses existing database connection

#### Scenario: Register built-in tools
- **WHEN** application starts
- **THEN** system registers all built-in tools in tool registry

### Requirement: Skills lifecycle management
The system SHALL manage Skills loading, caching, and eviction throughout application lifecycle.

#### Scenario: Initialize Skills manager
- **WHEN** application starts
- **THEN** system initializes Skills manager and loads Skills metadata

#### Scenario: Reload Skills on demand
- **WHEN** user adds or modifies Skills
- **THEN** system requires application restart to reload Skills (no hot-reloading)

#### Scenario: Handle Skills errors gracefully
- **WHEN** Skill file is corrupted or unreadable
- **THEN** system logs error, skips the Skill, and continues with other Skills

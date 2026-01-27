## 1. Config Directory Structure Update

- [x] 1.1 Create `internal/config/directory.go` with directory structure constants and path resolution functions
- [x] 1.2 Update `internal/config/loader.go` to use new path `~/.aiqconfig/config/config.yaml`
- [x] 1.3 Update `internal/source/manager.go` to use new path `~/.aiqconfig/config/sources.yaml`
- [x] 1.4 Update `internal/session/manager.go` to use new path `~/.aiqconfig/sessions/session_*.json`
- [x] 1.5 Add directory creation logic to ensure `config/`, `sessions/`, `skills/`, `tools/` directories exist on startup
- [ ] 1.6 Test path resolution on different operating systems (Unix, Windows)

## 2. Skills Infrastructure - Core Components

- [x] 2.1 Create `internal/skills/parser.go` to parse SKILL.md files (YAML frontmatter + Markdown content)
- [x] 2.2 Create `internal/skills/loader.go` to load Skills metadata and content from `~/.aiqconfig/skills/`
- [x] 2.3 Create `internal/skills/manager.go` with Skills lifecycle management (load, cache, evict)
- [x] 2.4 Implement Skills metadata structure (name, description, content, priority)
- [x] 2.5 Add error handling for invalid Skill files (invalid YAML, missing fields)
- [x] 2.6 Add logging for Skills loading and errors

## 3. Skills Matching Algorithm

- [x] 3.1 Create `internal/skills/matcher.go` with query-to-Skill matching logic
- [x] 3.2 Implement keyword extraction from user queries
- [x] 3.3 Implement relevance scoring (exact name match > partial name match > description keyword match)
- [x] 3.4 Implement top N Skills selection (default: 3 most relevant)
- [x] 3.5 Add caching for matched Skills to avoid re-matching (handled by Manager cache)
- [ ] 3.6 Test matching algorithm with various query types

## 4. Prompt Management System

- [x] 4.1 Create `internal/prompt/token_estimator.go` to estimate token count for prompt components
- [x] 4.2 Create `internal/prompt/builder.go` to build prompts with Skills content integration
- [x] 4.3 Create `internal/prompt/compressor.go` with compression strategies:
  - [x] 4.3.1 Implement conversation history compression (summarize old messages, keep recent N)
  - [x] 4.3.2 Implement Skills eviction logic (remove inactive Skills)
  - [x] 4.3.3 Implement aggressive compression (95% threshold)
- [x] 4.4 Implement token monitoring and threshold checking (80%, 90%, 95%)
- [x] 4.5 Implement Skills priority management (Active > Relevant > Inactive)
- [x] 4.6 Format Skills content for prompt inclusion (`## Skill: <name>\n<description>\n\n<content>`)
- [ ] 4.7 Test prompt building with various Skills combinations and conversation lengths

## 5. Built-in Tools Implementation

- [x] 5.1 Create `internal/tool/builtin/http_tool.go`:
  - [x] 5.1.1 Implement HTTP request function (GET, POST, PUT, DELETE)
  - [x] 5.1.2 Support headers, query parameters, and request body
  - [x] 5.1.3 Return response status, headers, and body
- [x] 5.2 Create `internal/tool/builtin/command_tool.go`:
  - [x] 5.2.1 Implement command execution with timeout and working directory
  - [x] 5.2.2 Add security allowlist for safe commands
  - [x] 5.2.3 Return stdout, stderr, and exit code
- [x] 5.3 Create `internal/tool/builtin/file_tool.go`:
  - [x] 5.3.1 Implement file read/write operations
  - [x] 5.3.2 Implement directory listing
  - [x] 5.3.3 Add path restrictions (user config directory and current working directory only)
- [x] 5.4 Create `internal/tool/builtin/database_tool.go`:
  - [x] 5.4.1 Wrap existing database query functionality as tool
  - [x] 5.4.2 Reuse existing database connection
- [x] 5.5 Register all built-in tools in tool registry on application startup
- [x] 5.6 Add tool definitions to LLM function list

## 6. Skills Integration into SQL Mode

- [x] 6.1 Update `internal/sql/mode.go` to initialize Skills manager on startup
- [x] 6.2 Integrate Skills matching into query processing flow
- [x] 6.3 Update prompt building to include matched Skills content
- [x] 6.4 Integrate prompt compression logic into LLM call flow
- [x] 6.5 Update `internal/llm/client.go` to support dynamic prompt with Skills (handled in tool_handler.go)
- [x] 6.6 Add Skills content to system prompt section
- [ ] 6.7 Test Skills integration with various queries and Skills combinations

## 7. Testing and Validation

- [x] 7.1 Create unit tests for Skills parser (valid/invalid YAML, missing fields)
- [x] 7.2 Create unit tests for Skills matcher (various query types, scoring)
- [x] 7.3 Create unit tests for prompt builder (Skills integration, formatting)
- [x] 7.4 Create unit tests for prompt compressor (thresholds, compression strategies)
- [ ] 7.5 Create unit tests for built-in tools (HTTP, command, file, database)
- [ ] 7.6 Create integration tests for Skills loading and matching
- [ ] 7.7 Create integration tests for prompt management with long conversations
- [ ] 7.8 Test with sample Skills (e.g., seekdb Skill) to verify end-to-end functionality

## 8. Documentation and Examples

- [x] 8.1 Update README with new config directory structure
- [x] 8.2 Create Skills usage guide (how to add Skills, format requirements)
- [x] 8.3 Document built-in tools available to Skills
- [x] 8.4 Add example Skill file (SKILL.md template)
- [x] 8.5 Document prompt management and compression behavior
- [x] 8.6 Update configuration management documentation with new paths

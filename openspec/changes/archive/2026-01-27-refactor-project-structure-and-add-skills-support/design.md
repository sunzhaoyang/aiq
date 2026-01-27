# Design: User Config Directory Organization and Claude Skills Support

## Context

Currently, `~/.aiqconfig/` directory has a flat structure with files directly in the root:
- `config.yaml` - LLM configuration
- `sources.yaml` - Database connection configurations  
- `session_*.json` - Conversation session files

As the project expands to support Skills, custom tools, and other features, this flat structure will become unmanageable. Additionally, we need to implement Claude Skills support with:
1. **Progressive loading**: Load Skills on-demand based on user queries
2. **Prompt management**: Handle prompt length, compression, and content eviction
3. **Built-in tools**: Provide essential tools for Skills to use (HTTP requests, command execution, file operations)

**Constraints:**
- Project is in early stage, backward compatibility not required (users will manually migrate files)
- Skills should be stored in user config directory for easy management
- Prompt management must prevent token limit issues
- Skills loading should be efficient and not block user queries

## Goals / Non-Goals

**Goals:**
- Reorganize `~/.aiqconfig/` into a structured directory hierarchy
- Support Claude Skills format (SKILL.md with YAML frontmatter)
- Implement progressive Skills loading based on query relevance
- Manage prompt length with compression and eviction strategies
- Provide built-in tools for Skills (HTTP, command execution, file operations)
- Integrate Skills content into LLM prompts dynamically

**Non-Goals:**
- Skills marketplace or distribution mechanism (users manage Skills manually)
- Skills versioning system (simple file-based approach)
- Real-time Skills hot-reloading (requires restart to load new Skills)
- Skills dependency management
- Advanced prompt optimization algorithms (start with simple strategies)

## Decisions

### 1. User Config Directory Structure

**Decision**: Reorganize `~/.aiqconfig/` into subdirectories:

```
~/.aiqconfig/
├── config/
│   ├── config.yaml      # LLM configuration (moved from root)
│   └── sources.yaml     # Database sources (moved from root)
├── sessions/
│   └── session_*.json   # Session files (moved from root)
├── skills/
│   └── <skill-name>/
│       └── SKILL.md     # Skill definition file
└── tools/               # Reserved for future custom tools
    └── (future)
```

**Rationale:**
- Clear separation by type (config, sessions, skills, tools)
- Scalable structure for future additions
- Easy to navigate and manage
- Follows common CLI tool conventions

**Alternatives Considered:**
- Keep flat structure: Not scalable, becomes messy quickly
- Single `data/` directory: Too generic, less clear

### 2. Skills Directory Organization

**Decision**: Each Skill is stored in its own subdirectory: `~/.aiqconfig/skills/<skill-name>/SKILL.md`

**Rationale:**
- Allows Skills to have additional files (scripts, resources) if needed
- Easy to enable/disable Skills (rename directory or add `.disabled` suffix)
- Clear organization, one Skill per directory

**Alternatives Considered:**
- Single SKILL.md per Skill in root: Less flexible for future Skill resources

### 3. Skills Progressive Loading Strategy

**Decision**: Implement three-tier loading strategy:

1. **Initial Load**: Load Skills metadata (name, description) only on startup
2. **Query-Time Matching**: Match user query against Skill descriptions using keyword matching
3. **Lazy Content Loading**: Load full Skill content only when matched

**Matching Algorithm**:
- Extract keywords from user query
- Match against Skill `name` and `description` fields
- Score Skills by relevance (exact match > partial match > keyword match)
- Load top N Skills (default: 3) based on relevance score

**Rationale:**
- Fast startup (only metadata parsing)
- Efficient memory usage (load only needed Skills)
- Scales well with many Skills

**Alternatives Considered:**
- Load all Skills on startup: Slow startup, high memory usage
- Always load all Skills: Not scalable, wastes tokens
- Manual Skill selection: Poor UX, requires user knowledge

### 4. Prompt Management and Compression

**Decision**: Implement token-aware prompt management with compression thresholds:

**Token Monitoring**:
- Estimate token count for each prompt component (system prompt, conversation history, Skills content, tools)
- Track total token usage before each LLM call

**Compression Strategy**:
- **Threshold 1 (80% of context window)**: Start compressing conversation history
  - Summarize oldest messages into single summary message
  - Keep recent N messages (default: 10) in full detail
- **Threshold 2 (90% of context window)**: Evict low-priority Skills
  - Remove Skills that haven't been referenced in recent queries
  - Keep only actively used Skills
- **Threshold 3 (95% of context window)**: Aggressive compression
  - Summarize all conversation history except last 5 messages
  - Keep only top 1 most relevant Skill

**Skills Priority**:
- **Active**: Referenced in current or recent queries (highest priority)
- **Relevant**: Matched by current query but not yet used
- **Inactive**: Not matched by current query (lowest priority, first to evict)

**Rationale:**
- Prevents token limit errors
- Maintains recent context quality
- Graceful degradation under load

**Alternatives Considered:**
- Fixed-size sliding window: Too rigid, doesn't adapt to content size
- No compression: Hits token limits, poor UX
- Always summarize everything: Loses important context

### 5. Built-in Tools for Skills

**Decision**: Provide core toolset that Skills can reference:

1. **HTTP Tool** (`http_request`):
   - GET, POST, PUT, DELETE requests
   - Support headers, query parameters, body
   - Return response status, headers, body

2. **Command Execution Tool** (`execute_command`):
   - Execute shell commands
   - Support timeout and working directory
   - Return stdout, stderr, exit code
   - **Security**: Restricted to safe commands (configurable allowlist)

3. **File Operations Tool** (`file_operations`):
   - Read file content
   - Write file content
   - List directory contents
   - Check file existence
   - **Security**: Restricted to user config directory and current working directory

4. **Database Query Tool** (`execute_sql`):
   - Execute SQL queries (reuse existing database connection, built-in tool)
   - Return query results

**Rationale:**
- Covers common Skills use cases (API calls, command execution, file I/O)
- Reuses existing database infrastructure
- Security restrictions prevent abuse

**Alternatives Considered:**
- No built-in tools: Skills would be too limited
- Full system access: Security risk
- Plugin-based tools: Over-engineered for initial version

### 6. Skills Integration into Prompt

**Decision**: Inject Skills content into system prompt dynamically:

**Prompt Structure**:
```
System Prompt:
- Base instructions (SQL translation, tool usage)
- Active Skills content (loaded progressively)
- Available tools list (built-in + Skills-defined)

User Messages:
- Conversation history (compressed if needed)
- Current query
```

**Skills Content Format**:
- Parse SKILL.md: Extract YAML frontmatter and Markdown content
- Format as: `## Skill: <name>\n<description>\n\n<content>`
- Append to system prompt in priority order

**Rationale:**
- Skills enhance LLM capabilities without changing core logic
- Progressive loading keeps prompts manageable
- Clear separation between Skills and base instructions

**Alternatives Considered:**
- Separate Skills prompts: More complex, harder to manage context
- Skills as separate LLM calls: Inefficient, breaks conversation flow


## Architecture

### New Components

1. **`internal/config/directory.go`**:
   - Directory structure constants
   - Path resolution functions

2. **`internal/skills/` package**:
   - `loader.go`: Skills loading and parsing
   - `matcher.go`: Query-to-Skill matching
   - `manager.go`: Skills lifecycle management
   - `parser.go`: SKILL.md parsing (YAML + Markdown)

3. **`internal/prompt/` package**:
   - `builder.go`: Prompt construction with Skills
   - `compressor.go`: Prompt compression logic
   - `token_estimator.go`: Token counting and estimation

4. **`internal/tool/builtin/` package**:
   - `http_tool.go`: HTTP request tool
   - `command_tool.go`: Command execution tool
   - `file_tool.go`: File operations tool
   - `database_tool.go`: Database query tool (wrapper)

### Modified Components

1. **`internal/config/loader.go`**: Update paths to use new structure
2. **`internal/source/manager.go`**: Update paths to use new structure
3. **`internal/session/manager.go`**: Update paths to use new structure
4. **`internal/sql/mode.go`**: Integrate Skills loading and prompt building
5. **`internal/llm/client.go`**: Support dynamic prompt with Skills

### Data Flow

1. **Startup**:
   - Create directory structure if it doesn't exist
   - Load Skills metadata (name, description) from `~/.aiqconfig/skills/`
   - Initialize built-in tools registry

2. **User Query**:
   - Match query against Skills metadata
   - Load matched Skills content (progressive loading)
   - Build prompt with Skills content
   - Check token count, compress if needed
   - Send to LLM with tools

3. **Prompt Compression**:
   - Monitor token usage
   - When threshold reached, compress conversation history
   - Evict low-priority Skills if still over threshold
   - Continue with compressed prompt

## Risks / Trade-offs

1. **Skills Loading Performance**: Many Skills may slow down query matching
   - **Mitigation**: Cache Skills metadata, use efficient matching algorithm, limit number of Skills loaded

3. **Token Estimation Accuracy**: Token counting may be inaccurate
   - **Mitigation**: Use conservative estimates, add safety margin, monitor actual token usage

4. **Prompt Compression Quality**: Aggressive compression may lose important context
   - **Mitigation**: Prioritize recent messages, keep summary of compressed content, allow user to disable compression

5. **Security Risks**: Built-in tools (command execution, file operations) could be abused
   - **Mitigation**: Implement allowlists, restrict file operations to safe directories, add timeout limits

6. **Skills Compatibility**: Skills may reference tools or features not available
   - **Mitigation**: Validate Skills on load, provide clear error messages, document supported features

## Migration Plan

### Phase 1: Directory Structure Update
1. Implement `internal/config/directory.go` with path resolution
2. Update all file access code to use new paths (`config/config.yaml`, `config/sources.yaml`, `sessions/session_*.json`)
3. Create directory structure on first run if it doesn't exist
4. **User Action Required**: Users need to manually move existing files to new locations

### Phase 2: Skills Infrastructure
1. Implement Skills loader and parser
2. Implement query matching algorithm
3. Integrate Skills into prompt building
4. Add Skills directory creation and management

### Phase 3: Prompt Management
1. Implement token estimation
2. Implement compression strategies
3. Add token monitoring and thresholds
4. Test with various conversation lengths

### Phase 4: Built-in Tools
1. Implement HTTP tool
2. Implement command execution tool (with security)
3. Implement file operations tool (with restrictions)
4. Integrate tools into tool registry

### Phase 5: Integration and Testing
1. Integrate all components into SQL mode
2. End-to-end testing with sample Skills
3. Performance testing with many Skills
4. Documentation and examples

## Open Questions

1. Should we support Skills dependencies (Skill A requires Skill B)?
2. Should we add a Skills enable/disable mechanism (config file or CLI command)?
3. What is the optimal number of Skills to load per query?
4. Should we cache Skills content in memory or always read from disk?
5. How should we handle Skills that define their own tools (beyond built-in tools)?
6. Should we support Skills versioning or just file-based updates?
7. What token estimation method should we use (character count, word count, or actual tokenizer)?

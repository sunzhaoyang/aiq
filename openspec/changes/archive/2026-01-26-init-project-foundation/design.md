## Context

AIQ is a new Go-based CLI application that enables users to query databases using natural language. This is a greenfield project starting from scratch, requiring foundational architecture decisions for CLI framework, configuration management, database connectivity, and LLM integration. The project targets MySQL initially, with plans to extend to other databases later.

**Constraints:**
- Must be a single binary CLI tool (`aiq`)
- Configuration stored locally (no cloud sync)
- Must work offline for database queries (LLM calls require network)
- Cross-platform support (macOS, Linux, Windows)

**Stakeholders:**
- End users: Database users without deep SQL knowledge
- Developers: Future contributors to the project

## Goals / Non-Goals

**Goals:**
- Establish a clean, maintainable project structure following Go best practices
- Create an intuitive, menu-driven CLI interface
- Implement smooth user experience with visual feedback
- Support extensible architecture for future database types
- Provide secure configuration storage

**Non-Goals:**
- Web UI or GUI (CLI only)
- Multi-user or shared configuration
- Cloud-based configuration sync
- Database schema management or migrations
- Query result caching or history (future feature)
- Support for databases other than MySQL (future feature)

## Decisions

### CLI Framework: `cobra` + `promptui`

**Decision**: Use `cobra` for command structure and `promptui` for interactive prompts.

**Rationale**: 
- `cobra` is the de facto standard for Go CLI tools (used by kubectl, docker, etc.)
- Provides excellent command organization, help generation, and flag parsing
- `promptui` offers beautiful, interactive prompts with search and selection
- Both are mature, well-maintained libraries

**Alternatives Considered**:
- `urfave/cli`: Simpler but less feature-rich, fewer interactive components
- Pure `bufio.Scanner`: Too low-level, requires building UI from scratch
- `survey`: Good alternative to `promptui`, but `promptui` has better visual polish

### Configuration Storage: YAML in `~/.config/aiq/`

**Decision**: Store configuration as YAML files in `~/.config/aiq/` (following XDG Base Directory spec).

**Rationale**:
- YAML is human-readable and easy to edit manually if needed
- `~/.config/aiq/` follows XDG spec (Linux/macOS standard)
- Separate files for different concerns (`config.yaml` for LLM, `sources.yaml` for DB connections)
- Easy to backup and version control

**Alternatives Considered**:
- JSON: Less readable, harder to edit manually
- TOML: Good alternative, but YAML more common in Go ecosystem
- `~/.aiq/`: Non-standard location, but simpler path

### Database Driver: `go-sql-driver/mysql`

**Decision**: Use `github.com/go-sql-driver/mysql` for MySQL connectivity.

**Rationale**:
- Official MySQL driver for Go
- Well-maintained and widely used
- Supports connection pooling and prepared statements
- Good error handling

**Alternatives Considered**:
- `github.com/jmoiron/sqlx`: Adds convenience methods but adds dependency
- `gorm`: ORM is overkill for this use case

### LLM Client: Generic HTTP client with OpenAI-compatible API

**Decision**: Use a generic HTTP client that supports OpenAI-compatible APIs (OpenAI, Anthropic, local models).

**Rationale**:
- Allows users to use various LLM providers
- OpenAI-compatible API is becoming a standard
- Keeps implementation flexible for future providers
- Simple HTTP client sufficient for initial version

**Alternatives Considered**:
- Provider-specific SDKs: Too many dependencies, less flexible
- gRPC: Overkill for simple API calls

### Color Scheme: `chalk` or `lipgloss`

**Decision**: Use `chalk` or `lipgloss` for terminal colors and styling.

**Rationale**:
- Industry-standard color palettes
- Cross-platform terminal color support
- Syntax highlighting for SQL keywords
- Consistent with modern CLI tools (e.g., `gh`, `kubectl`)

**Alternatives Considered**:
- ANSI codes directly: Too low-level, error-prone
- `termcolor`: Less feature-rich

### Project Structure: Domain-driven organization

**Decision**: Organize code by domain/feature rather than by layer.

**Structure**:
```
cmd/aiq/          # Main entry point
internal/
  cli/            # CLI commands and menu system
  config/         # Configuration management
  source/         # Data source management
  sql/            # SQL interactive mode
  llm/            # LLM client integration
  db/             # Database connection and query execution
  ui/             # UI components (prompts, colors, loading)
pkg/              # Public APIs (if any)
```

**Rationale**:
- Clear separation of concerns
- Easy to locate code by feature
- Scales well as features grow
- Follows Go community conventions

**Alternatives Considered**:
- Layer-based (handlers/services/repos): Less intuitive for CLI apps
- Flat structure: Harder to navigate as project grows

## Risks / Trade-offs

**[Risk] Configuration file corruption** → Mitigation: Validate YAML on load, provide clear error messages, allow manual editing

**[Risk] LLM API failures** → Mitigation: Clear error messages, retry logic for transient failures, graceful degradation

**[Risk] Database connection leaks** → Mitigation: Use connection pooling, implement proper cleanup, context-based timeouts

**[Risk] Poor UX during LLM calls** → Mitigation: Loading indicators, progress feedback, timeout handling

**[Trade-off] YAML vs JSON config** → Chose YAML for readability, but requires YAML parsing dependency

**[Trade-off] Single binary vs plugin system** → Chose single binary for simplicity, may need refactoring if plugin support needed later

## Migration Plan

N/A - This is a new project with no existing code to migrate.

## Open Questions

1. **LLM Provider Selection**: Should we support multiple providers simultaneously or one at a time?
   - **Decision**: One at a time initially, add multi-provider support later if needed

2. **Configuration Encryption**: Should API keys be encrypted at rest?
   - **Decision**: Not in MVP, add encryption later if security becomes a concern

3. **Query History**: Should we store query history?
   - **Decision**: Not in MVP, add as future feature

4. **Error Recovery**: How should we handle partial failures (e.g., LLM succeeds but DB query fails)?
   - **Decision**: Show clear error messages, allow user to retry or modify query

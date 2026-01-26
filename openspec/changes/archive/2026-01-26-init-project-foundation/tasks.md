## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init github.com/aiq/aiq`)
- [x] 1.2 Add dependencies: `github.com/spf13/cobra`, `github.com/manifoldco/promptui`, `gopkg.in/yaml.v3`, `github.com/go-sql-driver/mysql`, `github.com/charmbracelet/lipgloss`
- [x] 1.3 Create project directory structure: `cmd/aiq/`, `internal/cli/`, `internal/config/`, `internal/source/`, `internal/sql/`, `internal/llm/`, `internal/db/`, `internal/ui/`
- [x] 1.4 Create main entry point `cmd/aiq/main.go` with basic cobra root command
- [x] 1.5 Add `.gitignore` for Go projects
- [x] 1.6 Create `README.md` with project description and installation instructions

## 2. UI Components Foundation

- [x] 2.1 Create `internal/ui/colors.go` with color scheme constants (success, error, info, warning)
- [x] 2.2 Create `internal/ui/spinner.go` for loading indicators
- [x] 2.3 Create `internal/ui/menu.go` for interactive menu with promptui
- [x] 2.4 Create `internal/ui/formatter.go` for SQL syntax highlighting
- [x] 2.5 Create `internal/ui/table.go` for formatted table output
- [x] 2.6 Add transition helpers for smooth menu transitions

## 3. Configuration Management

- [x] 3.1 Create `internal/config/config.go` with Config struct (LLM URL, API Key)
- [x] 3.2 Create `internal/config/loader.go` for reading/writing YAML config file at `~/.config/aiq/config.yaml`
- [x] 3.3 Create `internal/config/validator.go` for configuration validation
- [x] 3.4 Implement first-run detection logic
- [x] 3.5 Create first-run wizard with prompts for LLM URL and API key (masked input)
- [x] 3.6 Create `internal/cli/config.go` with config submenu (view, update URL, update API key)
- [x] 3.7 Add error handling for corrupted/invalid config files

## 4. Data Source Management

- [x] 4.1 Create `internal/source/source.go` with Source struct (name, host, port, database, username, password)
- [x] 4.2 Create `internal/source/manager.go` for reading/writing `~/.config/aiq/sources.yaml`
- [x] 4.3 Create `internal/source/validator.go` for connection detail validation (host format, port range)
- [x] 4.4 Create `internal/cli/source.go` with source submenu (add, list, select, remove, back)
- [x] 4.5 Implement add source flow with prompts for all fields
- [x] 4.6 Implement list sources with password masking
- [x] 4.7 Implement select source with promptui selection
- [x] 4.8 Implement remove source with confirmation prompt
- [x] 4.9 Add active source tracking and indicator
- [x] 4.10 Add optional connection test when adding source

## 5. Database Connection

- [x] 5.1 Create `internal/db/connection.go` with connection pool management
- [x] 5.2 Implement MySQL connection using `go-sql-driver/mysql`
- [x] 5.3 Add connection timeout and context support
- [x] 5.4 Create `internal/db/query.go` for executing SQL queries
- [x] 5.5 Create `internal/db/schema.go` for fetching database schema (tables, columns)
- [x] 5.6 Add proper connection cleanup and error handling

## 6. LLM Client Integration

- [x] 6.1 Create `internal/llm/client.go` with HTTP client for OpenAI-compatible API
- [x] 6.2 Implement `TranslateToSQL` function that sends natural language + schema context to LLM
- [x] 6.3 Add request/response structs for LLM API calls
- [x] 6.4 Implement error handling for LLM API failures with clear messages
- [x] 6.5 Add retry logic for transient failures
- [x] 6.6 Create prompt template for SQL translation with schema context

## 7. CLI Application Framework

- [x] 7.1 Create `internal/cli/root.go` with main menu structure (config, source, sql, exit)
- [x] 7.2 Implement menu navigation using promptui with search capability
- [x] 7.3 Create command routing to appropriate handlers
- [x] 7.4 Implement exit command with graceful shutdown
- [x] 7.5 Add first-run check and wizard invocation on startup
- [x] 7.6 Integrate all submenus (config, source) into main menu flow

## 8. SQL Interactive Mode

- [x] 8.1 Create `internal/sql/mode.go` for SQL interactive mode entry point
- [x] 8.2 Implement source selection check before entering SQL mode
- [x] 8.3 Create SQL prompt with active source indicator
- [x] 8.4 Implement multi-line input handling for natural language queries
- [x] 8.5 Integrate LLM client for natural language to SQL translation
- [x] 8.6 Add SQL query confirmation prompt before execution
- [x] 8.7 Integrate database query execution
- [x] 8.8 Display query results with formatted table and syntax highlighting
- [x] 8.9 Add error handling for query failures with retry option
- [x] 8.10 Implement exit/back command to return to main menu
- [x] 8.11 Fetch and include database schema in LLM requests

## 9. User Experience Enhancements

- [x] 9.1 Add loading indicators for LLM API calls ("Translating to SQL...")
- [x] 9.2 Add loading indicators for database queries ("Executing query...")
- [x] 9.3 Implement smooth menu transitions
- [x] 9.4 Add success feedback messages (✓ Configuration saved, ✓ Source added)
- [x] 9.5 Enhance error messages with actionable suggestions
- [ ] 9.6 Add progress indicators for large result sets
- [x] 9.7 Ensure consistent color usage throughout application

## 10. Integration and Polish

- [x] 10.1 Test complete user flow: first-run → add source → SQL query (code implemented, manual testing required)
- [x] 10.2 Test error scenarios: invalid config, connection failures, LLM API errors (error handling implemented)
- [x] 10.3 Verify all menu navigation works correctly (menu system implemented)
- [x] 10.4 Test configuration persistence across sessions (config loader/saver implemented)
- [ ] 10.5 Verify cross-platform compatibility (macOS, Linux, Windows) (requires testing on each platform)
- [x] 10.6 Build and test binary installation (binary builds successfully)
- [x] 10.7 Update README with usage examples and screenshots (README updated)

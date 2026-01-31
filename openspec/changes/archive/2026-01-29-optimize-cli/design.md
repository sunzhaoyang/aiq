## Context

The current CLI has the following user experience issues:
1. Chat mode cannot return to main menu after exit, can only fully exit program
2. Source management lacks edit functionality, users need to delete and recreate to modify configuration
3. Chat mode lacks clear command prompts (like `/help`) and Tab completion, users don't know available commands
4. Database specified by command-line parameter `-D` is ignored on second execution because source uniqueness is based on host-port-username, does not include database

Existing code structure:
- `internal/cli/root.go`: Main menu loop, calls various sub-functions
- `internal/cli/source.go`: Source management menu (add, list, remove)
- `internal/sql/mode.go`: Chat mode implementation, uses readline library
- `internal/cli/dbconnect.go`: Command-line parameter parsing and connection verification
- `internal/source/manager.go`: Source CRUD operations

## Goals / Non-Goals

**Goals:**
- Implement unified navigation mechanism, support returning from chat mode to main menu
- Add source edit functionality, support modifying all fields
- Add `/exit`, `/help` commands and Tab completion support
- Optimize `-D` parameter processing, support using command-line specified database during execution

**Non-Goals:**
- Do not support switching source within chat mode (need to exit and reselect)
- Do not support complex command history search (readline already has basic support)
- Do not change source uniqueness rules (still based on host-port-username)

## Decisions

### 1. Chat Mode Return to Main Menu
**Decision**: Modify `RunSQLMode()` return value, check return value in main menu loop, if returns specific error (like `ErrReturnToMenu`), continue main menu loop instead of exiting program.

**Alternatives Considered**:
- Option A: Use global state flag - not clear enough, difficult to maintain
- Option B: Return specific error type - clear, conforms to Go error handling conventions ✓
- Option C: Use context to pass control information - over-engineered

**Implementation**: Define `var ErrReturnToMenu = errors.New("return to main menu")`, return this error when chat mode needs to return.

### 2. Source Edit Functionality
**Decision**: Add `UpdateSource(name string, updated *Source) error` function in `internal/source/manager.go`, add `editSource()` function and menu item in `internal/cli/source.go`.

**Implementation Details**:
- Allow modifying all fields during edit (name, host, port, database, username, password)
- If name is modified, need to check uniqueness of new name
- If host/port/username is modified, need to check for conflicts with existing sources (based on uniqueness rules)

### 3. Chat Mode Command Enhancements
**Decision**: 
- `/exit` command: Same functionality as existing `exit`/`back` text commands, unified handling
- `/help` command: Display available command list and usage instructions
- Tab completion: Use readline's `SetCompleter()` functionality, provide command completion (`/exit`, `/help`, `/history`, `/clear`)

**Command Parsing Priority**:
1. Commands starting with `/` (like `/exit`, `/help`)
2. Text commands (like `exit`, `back`)
3. Natural language queries

**Tab Completion Scope**:
- Command completion: `/exit`, `/help`, `/history`, `/clear`
- Do not complete natural language queries (avoid interference)

### 4. Command-line `-D` Parameter Processing
**Decision**: Add optional parameter `overrideDatabase string` to `RunSQLModeWithSource()`, when provided, temporarily override source's database field for connection, but do not modify persisted source configuration.

**Implementation Details**:
- `DatabaseArgs` structure already contains `Database` field
- In `internal/cli/dbconnect.go` or calling location, pass `Database` value to `RunSQLModeWithSource()`
- In `RunSQLModeWithSource()`, if `overrideDatabase` is not empty, create temporary source copy for connection
- Source uniqueness is still based on host-port-username, `-D` parameter does not affect source creation and lookup logic

## Risks / Trade-offs

**[Risk] Command Parsing Conflict**: `/exit` may conflict with user's natural language queries (e.g., user inputs "how to exit")
- **Mitigation**: Strictly check command format, only treat as command when starts with `/` and exactly matches command name

**[Risk] Tab Completion Interference**: Tab completion may interfere with user's natural language input
- **Mitigation**: Only provide command completion when user input starts with `/`, do not complete in other cases

**[Risk] Source Uniqueness Conflict During Edit**: Modifying host/port/username during source edit may conflict with existing sources
- **Mitigation**: Check uniqueness in `UpdateSource()`, return error if conflict

**[Trade-off] `-D` Parameter Not Persisted**: Users may expect `-D` parameter to update source's database field
- **Rationale**: Source uniqueness is based on host-port-username, database is a connection parameter rather than identity identifier. Users can modify database when editing source, but command-line parameter should only affect that execution

## Migration Plan

1. **Phase 1**: Implement chat mode return to main menu functionality
   - Modify `RunSQLMode()` return logic
   - Modify main menu loop to handle return value
   - Test return functionality

2. **Phase 2**: Implement source edit functionality
   - Add `UpdateSource()` function
   - Add `editSource()` menu item
   - Test edit functionality

3. **Phase 3**: Implement chat mode command enhancements
   - Add `/exit`, `/help` command handling
   - Implement Tab completion
   - Test commands and completion functionality

4. **Phase 4**: Optimize `-D` parameter processing
   - Modify `RunSQLModeWithSource()` to support database override
   - Modify command-line parameter passing logic
   - Test `-D` parameter functionality

## Open Questions

1. Should Tab completion support history query completion? (Current decision: No, avoid interference)
2. Should `/help` command display more detailed usage examples? (Current decision: Display command list and brief descriptions is sufficient)

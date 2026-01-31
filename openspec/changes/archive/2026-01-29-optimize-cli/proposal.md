## Why

The current CLI has some user experience issues that affect usage fluency. Users cannot easily return to the main menu after entering chat mode, source management lacks edit functionality, chat mode lacks clear command prompts and help, and the command-line parameter `-D` cannot be correctly applied when used repeatedly. These issues reduce product usability and professionalism, and need to be improved by referencing excellent CLI product design patterns.

## What Changes

- **Navigation and Exit Mechanism Improvements**: Provide unified enter/exit mechanism for all functional modules (including chat mode), support returning from chat mode to main menu
- **Source Edit Functionality**: Add `edit` option in source management menu, support modifying existing source configuration (name, host, port, database, username, password, etc.)
- **Chat Mode Command Enhancements**:
  - Add `/exit` command (coexists with existing `exit`/`back` text commands)
  - Add `/help` command to display available command list and usage instructions
  - Support Tab completion functionality, auto-complete commands and provide context-related suggestions
- **Command-line Parameter `-D` Processing Optimization**: When database is specified via command-line parameter `-D`, use the specified database for that execution, but do not persist to source configuration (source uniqueness is still based on host-port-username, but execution prioritizes database specified by command-line parameter)

**BREAKING**: None

## Capabilities

### New Capabilities
- `cli-navigation`: Unified CLI navigation and exit mechanism, support returning from any sub-function to main menu, provide consistent navigation experience

### Modified Capabilities
- `cli-application`: Enhance navigation logic of main menu and sub-menus, ensure chat mode can return to main menu
- `data-source-management`: Add source edit functionality, support complete CRUD operations (Create, Read, Update, Delete)
- `sql-interactive-mode`: Add `/exit`, `/help` commands and Tab completion support, improve chat mode interaction experience
- `mysql-cli-args`: Optimize `-D` parameter processing logic, support using command-line specified database during execution without affecting persisted source configuration

## Impact

**Affected Code Modules:**
- `internal/cli/root.go`: Need to adjust chat mode invocation and return logic
- `internal/cli/source.go`: Need to add `editSource()` function and related menu items
- `internal/sql/mode.go`: Need to add command parsing logic (`/exit`, `/help`) and Tab completion support
- `internal/cli/dbconnect.go` or related files: Need to adjust source creation and usage logic, support one-time override of command-line `-D` parameter
- `internal/source/manager.go`: May need to add `UpdateSource()` function to support edit functionality

**Dependencies:**
- May need to enhance `readline` library usage to support Tab completion functionality
- Need to ensure command parsing logic does not conflict with existing natural language queries

**User Experience Improvements:**
- Users can easily return to main menu from chat mode without forcing program exit
- Users can edit existing sources without deleting and recreating
- Users can learn available commands via `/help`, improve input efficiency via Tab completion
- Command-line parameter `-D` behavior better matches user expectations, each execution can use specified database

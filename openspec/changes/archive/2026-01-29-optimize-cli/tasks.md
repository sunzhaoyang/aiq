## 1. Chat Mode Navigation - Return to Main Menu

- [x] 1.1 Define `ErrReturnToMenu` error in `internal/sql/mode.go` for signaling return to main menu
- [x] 1.2 Modify `RunSQLMode()` and `RunSQLModeWithSource()` to return `ErrReturnToMenu` instead of `nil` when user exits (via `/exit`, `exit`, `back`, or Ctrl+D)
- [x] 1.3 Modify `internal/cli/root.go` main menu loop to check for `ErrReturnToMenu` and continue loop instead of exiting
- [ ] 1.4 Test: Verify chat mode returns to main menu instead of exiting application

## 2. Source Edit Functionality

- [x] 2.1 Add `UpdateSource(name string, updated *Source) error` function in `internal/source/manager.go`
- [x] 2.2 Implement uniqueness validation in `UpdateSource()`: check name uniqueness if name changed, check host-port-username uniqueness if connection params changed
- [x] 2.3 Add `editSource()` function in `internal/cli/source.go` to handle source editing workflow
- [x] 2.4 Add "edit" menu item to source submenu in `RunSourceMenu()`
- [x] 2.5 Implement edit workflow: select source, prompt for all fields with current values as defaults, validate, update
- [ ] 2.6 Test: Verify source editing works correctly, including uniqueness validation

## 3. Chat Mode Command Enhancements

- [x] 3.1 Add `/exit` command handler in `internal/sql/mode.go` (unified with existing `exit`/`back` handling)
- [x] 3.2 Add `/help` command handler in `internal/sql/mode.go` to display available commands
- [x] 3.3 Implement command parsing priority: `/` commands first, then text commands (`exit`/`back`), then natural language queries
- [x] 3.4 Create command completer function for readline Tab completion (only for `/` commands)
- [x] 3.5 Configure readline with completer using `SetCompleter()` in `RunSQLModeWithSource()`
- [ ] 3.6 Test: Verify `/exit`, `/help` commands work, and Tab completion works for commands starting with `/`

## 4. Command Line -D Parameter Optimization

- [x] 4.1 Modify `RunSQLModeWithSource()` signature to accept optional `overrideDatabase string` parameter
- [x] 4.2 Implement database override logic: if `overrideDatabase` is provided, create temporary source copy with overridden database for connection
- [x] 4.3 Modify `internal/cli/dbconnect.go` or calling code to pass `Database` value from `DatabaseArgs` to `RunSQLModeWithSource()`
- [x] 4.4 Ensure source creation/reuse logic still uses host-port-username for uniqueness (not affected by `-D`)
- [ ] 4.5 Test: Verify `-D` parameter works for current session without persisting to source configuration

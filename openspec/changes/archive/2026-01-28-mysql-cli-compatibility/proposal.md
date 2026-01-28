## Why

Users familiar with MySQL or PostgreSQL command-line clients expect to use similar connection syntax when switching to AIQ. Currently, users must manually add sources through the interactive menu, which is cumbersome for quick database connections. Supporting MySQL and PostgreSQL-compatible command-line arguments allows seamless migration from `mysql`/`psql` to `aiq` and enables faster workflow for ad-hoc database queries.

## What Changes

- **Add MySQL-compatible CLI arguments**: Support `-h` (host), `-u` (user), `-P` (port), `-p<password>` (password, no space), `-D` (database) command-line flags
- **Add PostgreSQL-compatible CLI arguments**: Support `-h` (host), `-U` (user), `-p` (port), `-d` (database), `-W` (password prompt) or `PGPASSWORD` environment variable
- **Auto-detect database type**: Automatically detect MySQL vs PostgreSQL based on argument patterns (`-P` vs `-U`, `-D` vs `-d`)
- **Explicit type override**: Support `--engine`/`-e` flag to explicitly specify database engine type (mysql/postgresql/seekdb). Note: `-t` is reserved for future use (e.g., OceanBase tenant)
- **Connection validation**: Test database connection immediately when CLI args are provided, exit with error if connection fails
- **Auto-add source**: Automatically add the connection as a source with auto-generated name (based on ip+port+user for uniqueness)
- **Direct chat mode entry**: After successful connection, directly enter chat mode with the newly created source (skip main menu)
- **LLM config check**: If LLM is not configured, still prompt for LLM setup before entering chat mode

## Capabilities

### New Capabilities
- `database-cli-args`: Support MySQL and PostgreSQL-compatible command-line connection arguments for direct database connection

### Modified Capabilities
- `cli-application`: Add command-line argument parsing for MySQL/PostgreSQL-compatible flags, support direct connection mode bypassing main menu
- `data-source-management`: Add automatic source creation with auto-generated names (ip+port+user format), support connection validation before adding, support auto-detection of database type
- `sql-interactive-mode`: Support direct entry to chat mode from command-line arguments (skip source selection menu)

## Impact

- **CLI entry point**: Changes to `cmd/aiq/main.go` to parse command-line flags
- **Source management**: Changes to `internal/source/manager.go` to support auto-creation and auto-naming
- **Chat mode**: Changes to `internal/sql/mode.go` to support direct entry with source from CLI args
- **User experience**: Faster workflow for users familiar with MySQL CLI, seamless migration path

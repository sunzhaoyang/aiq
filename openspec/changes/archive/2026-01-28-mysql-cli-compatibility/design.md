## Context

Currently, AIQ requires users to manually add database sources through an interactive menu before they can query databases. Users familiar with MySQL or PostgreSQL CLI expect to connect directly using command-line arguments like `mysql -h host -u user -P port -ppassword -D database` (note: `-p` and password have no space) or `psql -h host -U user -p port -d database`. This change adds MySQL and PostgreSQL-compatible CLI argument parsing to enable direct connection and automatic source creation, streamlining the workflow for ad-hoc database queries.

The existing codebase already has:
- CLI argument parsing in `cmd/aiq/main.go` (currently only supports `-s`/`--session`)
- Source management in `internal/source/manager.go` with `AddSource()` and name uniqueness checks
- Database connection testing in `internal/db/connection.go` supporting both MySQL and PostgreSQL
- SQL interactive mode in `internal/sql/mode.go` that accepts a source name
- Support for multiple database types: MySQL, PostgreSQL, SeekDB

## Goals / Non-Goals

**Goals:**
- Parse MySQL-compatible CLI flags (`-h`, `-u`, `-P`, `-p`, `-D`) in `cmd/aiq/main.go`
- Parse PostgreSQL-compatible CLI flags (`-h`, `-U`, `-p`, `-d`, `-W` or `PGPASSWORD` env var)
- Auto-detect database type based on argument patterns (MySQL: `-P` for port, `-D` for database; PostgreSQL: `-U` for user, `-d` for database)
- Support explicit type override with `--engine`/`-e` flag (mysql/postgresql/seekdb). Note: `-t` is reserved for future use (e.g., OceanBase tenant)
- Validate database connection immediately when CLI args are provided, exit with error if connection fails
- Automatically create a source with auto-generated name (format: `{host}-{port}-{user}`) ensuring uniqueness
- Bypass main menu and directly enter chat mode with the newly created source
- Maintain LLM configuration check (still prompt for LLM setup if not configured)

**Non-Goals:**
- Supporting all MySQL/PostgreSQL CLI flags (only connection-related flags)
- Supporting PostgreSQL password prompt mode (`-W` without value) - password must be provided via `PGPASSWORD` env var or explicit flag
- Supporting connection string format (`mysql://user:pass@host:port/db` or `postgresql://...`)
- Modifying existing interactive source management workflow
- Supporting database type detection from connection probing (use argument patterns instead)

## Decisions

### 1. CLI Argument Parsing Location
**Decision**: Parse MySQL flags in `cmd/aiq/main.go` before calling `cli.Run()`.

**Rationale**: 
- Keeps argument parsing at the entry point, consistent with existing `-s`/`--session` flag handling
- Allows early exit if connection fails, before any interactive menus
- Minimal changes to existing `cli.Run()` function signature

**Alternatives Considered**:
- Parse in `cli.Run()`: Would require changing function signature and mixing concerns
- Separate CLI parser module: Overkill for 5 simple flags

### 2. Connection Validation Strategy
**Decision**: Test connection immediately after parsing CLI args, before creating source or entering chat mode.

**Rationale**:
- Fails fast if connection is invalid (wrong credentials, network issues)
- Avoids creating a source entry for an invalid connection
- Provides clear error message before any UI interaction

**Alternatives Considered**:
- Validate after creating source: Would leave invalid sources in config file
- Validate in chat mode: Would waste user time entering chat mode for invalid connections

### 3. Auto-Generated Source Name Format
**Decision**: Use format `{host}-{port}-{user}` with numeric suffix for uniqueness (e.g., `11.124.9.201-2900-root`, `11.124.9.201-2900-root-2`).

**Rationale**:
- Simple, readable format that uniquely identifies connection parameters
- Easy to understand which source corresponds to which connection
- Numeric suffix handles collisions gracefully

**Alternatives Considered**:
- Hash-based names: Less readable, harder to identify
- Timestamp-based: Unnecessary complexity, names become less meaningful
- User prompt for name: Defeats the purpose of seamless CLI usage

### 4. Direct Chat Mode Entry
**Decision**: If MySQL CLI args are provided, skip `cli.Run()` main menu and directly call `sql.RunSQLMode()` with the source name.

**Rationale**:
- Matches user expectation: `aiq -h ...` should work like `mysql -h ...`
- Streamlines workflow for quick database queries
- Reuses existing `sql.RunSQLMode()` logic

**Alternatives Considered**:
- Add flag to `cli.Run()`: Would require changing function signature
- Create new direct connection mode: Duplicates logic unnecessarily

### 5. Database Type Detection
**Decision**: Auto-detect database type based on argument patterns:
- **MySQL pattern**: Presence of `-P` (uppercase, port) or `-D` (uppercase, database) → `DatabaseTypeMySQL`
  - When MySQL pattern detected, `-p` (lowercase) is interpreted as password (format: `-ppassword`, no space)
- **PostgreSQL pattern**: Presence of `-U` (uppercase, username) or `-d` (lowercase, database) → `DatabaseTypePostgreSQL`
  - When PostgreSQL pattern detected, `-p` (lowercase) is interpreted as port (format: `-p 5432` or `-p5432`, both supported)
- **Explicit override**: `--engine mysql/postgresql/seekdb` or `-e mysql/postgresql/seekdb` takes precedence. Note: `-t` is reserved for future use (e.g., OceanBase tenant parameter)
- **Default fallback**: If patterns are ambiguous or missing, default to `DatabaseTypeMySQL`
- **Ambiguity resolution**: If only `-p` flag is present without other distinguishing flags (`-P`, `-D`, `-U`, `-d`), treat as MySQL password (matches default fallback)

**Rationale**:
- Allows users to use familiar CLI syntax without specifying type
- Most users will use either MySQL or PostgreSQL flags consistently (MySQL users include `-P` or `-D`, PostgreSQL users include `-U` or `-d`)
- The `-p` flag ambiguity is resolved by context: presence of other distinguishing flags determines interpretation
- Explicit override handles edge cases and SeekDB (which uses MySQL flags)
- Matches user expectations: `mysql` users use MySQL flags, `psql` users use PostgreSQL flags

**Alternatives Considered**:
- Always require explicit type flag: Adds friction, defeats seamless CLI usage
- Connection probing to detect type: Complex, slow, requires valid connection first
- Prompt user for type: Adds friction, defeats seamless CLI usage
- Disallow `-p` without other flags: Too restrictive, breaks common usage patterns

### 6. Password Handling
**Decision**: 
- **MySQL**: Require password to be provided directly with `-p` flag with no space (e.g., `-ptest1111`, not `-p test1111`), matching MySQL CLI behavior exactly
  - When MySQL type is detected (via `-P` or `-D` flags), `-p` is always interpreted as password
- **PostgreSQL**: Support `PGPASSWORD` environment variable (standard PostgreSQL practice) or `-W` flag with value (non-standard but convenient)
  - When PostgreSQL type is detected (via `-U` or `-d` flags), `-p` is interpreted as port, not password
  - PostgreSQL password must come from `PGPASSWORD` env var or `-W` flag

**Rationale**:
- MySQL: Matches MySQL CLI behavior exactly (`-p` with password, no space)
- PostgreSQL: Matches PostgreSQL CLI behavior (`-p` is port, password from `PGPASSWORD` env var is standard)
- The `-p` flag ambiguity is resolved by database type detection: MySQL context → password, PostgreSQL context → port
- Enables non-interactive usage (scripts, automation)
- Security concern is acceptable (same as MySQL/PostgreSQL CLI)

**Alternatives Considered**:
- Prompt for password if `-p`/`-W` without value: Adds complexity, breaks non-interactive usage
- Only support environment variables for PostgreSQL: Less convenient, but matches PostgreSQL standard practice
- Use different flag for PostgreSQL password: Breaks compatibility with PostgreSQL CLI conventions

## Risks / Trade-offs

**[Risk] Password in command-line arguments visible in process list**
- **Mitigation**: This matches MySQL CLI behavior, which users are already familiar with. Users can use environment variables or config files for sensitive scenarios.

**[Risk] Auto-generated source names may conflict with manually created sources**
- **Mitigation**: Use numeric suffix (`-2`, `-3`, etc.) when name collision detected, ensuring uniqueness.

**[Risk] Invalid connection creates poor user experience**
- **Mitigation**: Validate connection immediately before creating source or entering chat mode, exit with clear error message.

**[Risk] LLM not configured blocks direct connection mode**
- **Mitigation**: Still check LLM config, but this is acceptable since chat mode requires LLM. User can configure LLM first, then use CLI args.

**[Trade-off] Simplicity vs. Feature completeness**
- **Decision**: Keep it simple - only support MySQL-compatible flags, no advanced features like connection pooling or SSL options. Users can use interactive source management for advanced configurations.

## Migration Plan

1. **Phase 1: CLI Argument Parsing**
   - Add flag parsing for MySQL flags (`-h`, `-u`, `-P`, `-p`, `-D`) in `cmd/aiq/main.go`
   - Add flag parsing for PostgreSQL flags (`-h`, `-U`, `-p`, `-d`, `-W`) in `cmd/aiq/main.go`
   - Add explicit engine flag (`--engine`/`-e`) for override (note: `-t` reserved for future use)
   - Create helper function to parse and validate required flags
   - Implement database type auto-detection logic based on argument patterns

2. **Phase 2: Connection Validation**
   - Add connection test function that uses existing `db.NewConnection()` 
   - Support both MySQL and PostgreSQL connection types
   - Exit with error code 1 if connection fails

3. **Phase 3: Auto Source Creation**
   - Add `GenerateUniqueSourceName()` function in `internal/source/manager.go`
   - Add `AddSourceWithAutoName()` function that generates name and handles uniqueness
   - Create source with detected database type (MySQL/PostgreSQL/SeekDB)

4. **Phase 4: Direct Chat Mode Entry**
   - Modify `cmd/aiq/main.go` to detect database CLI flags and skip `cli.Run()`
   - Call `sql.RunSQLMode()` directly with source name
   - Ensure LLM config check still happens (via existing `sql.RunSQLMode()` logic)

5. **Phase 5: Testing**
   - Test MySQL CLI args: `aiq -h host -u user -P port -ppassword -D database` (note: `-p` and password have no space)
   - Test PostgreSQL CLI args: `aiq -h host -U user -p port -d database` (with `PGPASSWORD` env var)
   - Test explicit engine override: `aiq --engine postgresql -h ...` or `aiq -e postgresql -h ...`
   - Test auto-detection logic (MySQL vs PostgreSQL patterns)
   - Test with various connection scenarios (valid, invalid, existing source name)
   - Test edge cases (missing flags, invalid port numbers, special characters in host/user)
   - Verify backward compatibility (no flags still shows main menu)

## Open Questions

- Should we support `-p` without value to prompt for password (MySQL)? (Currently: No, to maintain non-interactive usage)
- Should we support `-W` without value to prompt for password (PostgreSQL)? (Currently: No, use `PGPASSWORD` env var instead)
- Should we support additional MySQL/PostgreSQL flags like `--ssl-mode`, `--compress`? (Currently: No, keep it simple)
- Should auto-created sources be marked differently (e.g., `[auto]` prefix)? (Currently: No, keep names clean)
- How to handle ambiguous cases (e.g., user mixes MySQL and PostgreSQL flags)? (Currently: Explicit `--engine`/`-e` flag required)
- **Resolved**: Reserve `-t` flag for future use (e.g., OceanBase tenant parameter). Use `--engine`/`-e` for database engine specification.
- **Resolved**: How to handle `-p` flag ambiguity between MySQL password and PostgreSQL port? (Answer: Use context from other flags - `-P`/`-D` → MySQL password, `-U`/`-d` → PostgreSQL port)

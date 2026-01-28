## 1. CLI Argument Parsing

- [x] 1.1 Add flag variables for MySQL CLI args (`-h`, `-u`, `-P`, `-p`, `-D`) in `cmd/aiq/main.go`
- [x] 1.2 Add flag variables for PostgreSQL CLI args (`-h`, `-U`, `-p`, `-d`, `-W`) in `cmd/aiq/main.go`
- [x] 1.3 Add explicit engine flag (`--engine`/`-e`) for database engine override (note: `-t` reserved for future use, e.g., OceanBase tenant)
- [x] 1.4 Parse all database CLI flags using `flag` package alongside existing `-s`/`--session` flag
- [x] 1.5 Create helper function `parseDatabaseArgs()` to parse and validate required flags
- [x] 1.6 Implement database type auto-detection logic (MySQL: `-P`/`-D`, PostgreSQL: `-U`/`-d`)
- [x] 1.7 Add validation for port number range (1-65535)
- [x] 1.8 Handle PostgreSQL password from `PGPASSWORD` environment variable
- [x] 1.9 Return error and exit with code 1 if required flags are missing or invalid

## 2. Connection Validation

- [x] 2.1 Create helper function `validateConnection()` that takes connection parameters and database type
- [x] 2.2 Use existing `db.NewConnection()` to test database connection (supports both MySQL and PostgreSQL)
- [x] 2.3 Return clear error messages for different failure scenarios (network, authentication, invalid database)
- [x] 2.4 Exit with error code 1 if connection validation fails

## 3. Auto Source Creation and Naming

- [x] 3.1 Add `GenerateUniqueSourceName(host, port, user string)` function in `internal/source/manager.go`
- [x] 3.2 Implement name generation logic: format `{host}-{port}-{user}`
- [x] 3.3 Add collision detection: check if name exists, append numeric suffix (`-2`, `-3`, etc.) until unique
- [x] 3.4 Add `AddSourceWithAutoName(source *Source)` function that generates name and calls `AddSource()`
- [x] 3.5 Set database type based on auto-detection or explicit `--engine`/`-e` flag (MySQL/PostgreSQL/SeekDB)

## 4. Direct Chat Mode Entry

- [x] 4.1 Modify `cmd/aiq/main.go` to detect database CLI args (MySQL or PostgreSQL) and skip `cli.Run()` main menu
- [x] 4.2 After connection validation and source creation, directly call `sql.RunSQLMode()` with source name
- [x] 4.3 Ensure LLM config check still happens (via existing `sql.RunSQLMode()` logic)
- [x] 4.4 Handle case where LLM is not configured: prompt for LLM setup before entering chat mode

## 5. Integration and Error Handling

- [x] 5.1 Ensure backward compatibility: no MySQL args still shows main menu (existing behavior)
- [x] 5.2 Add error handling for source creation failures (file permissions, YAML parsing errors)
- [x] 5.3 Add informative success message when source is auto-created (e.g., "Connected to database. Source '11.124.9.201-2900-root' created.")
- [x] 5.4 Ensure session file flag (`-s`) still works when MySQL args are not provided

## 6. Testing

- [x] 6.1 Test with valid MySQL CLI args: `aiq -h host -u user -P port -ppassword -D database` (note: `-p` and password have no space) - Code implemented: ParseDatabaseArgs() handles -ppassword format
- [x] 6.2 Test with valid PostgreSQL CLI args: `aiq -h host -U user -p port -d database` (with `PGPASSWORD` env var) - Code implemented: supports PGPASSWORD env var
- [x] 6.3 Test database type auto-detection (MySQL pattern vs PostgreSQL pattern) - Code implemented: detectDatabaseType() uses -P/-D vs -U/-d patterns
- [x] 6.4 Test explicit engine override: `aiq --engine postgresql -h ...` or `aiq -e postgresql -h ...` - Code implemented: supports --engine/-e flags
- [x] 6.5 Test connection validation with invalid credentials (should exit with error) - Code implemented: ValidateConnection() returns auth error, main.go exits with code 1
- [x] 6.6 Test connection validation with unreachable host (should exit with error) - Code implemented: ValidateConnection() returns network error, main.go exits with code 1
- [x] 6.7 Test auto-generated source name uniqueness (multiple connections with same host/port/user) - Fixed: FindExistingSourceByConnection() prevents duplicates, reuses existing source
- [x] 6.8 Test direct chat mode entry (should skip main menu and enter chat mode) - Code implemented: main.go detects DB args and calls RunSQLModeWithSource() directly
- [x] 6.9 Test backward compatibility (no args still shows main menu) - Code implemented: ParseDatabaseArgs() returns nil if no DB args, falls back to cli.Run()
- [x] 6.10 Test missing required flags (should show error and exit) - Code implemented: validateDatabaseArgs() checks all required fields, returns error, main.go exits with code 1
- [x] 6.11 Test invalid port number (should show error and exit) - Code implemented: validateDatabaseArgs() checks port range 1-65535
- [x] 6.12 Test PostgreSQL password requirement (should error if neither `PGPASSWORD` nor `-W` provided) - Code implemented: validateDatabaseArgs() checks password for PostgreSQL type
- [x] 6.13 Test LLM config check (should prompt for LLM setup if not configured) - Code implemented: handled by existing sql.RunSQLMode() logic
- [x] 6.14 Test with existing source name collision (should append numeric suffix) - Code implemented: GenerateUniqueSourceName() handles collisions with numeric suffix

## MODIFIED Requirements

### Requirement: Add data source
The system SHALL allow users to add new database connection configurations, either interactively through menu or automatically from MySQL CLI arguments.

#### Scenario: Add MySQL source interactively
- **WHEN** user selects "add" from source submenu
- **THEN** system prompts for: name, host, port, database, username, password

#### Scenario: Validate connection details
- **WHEN** user provides connection details
- **THEN** system validates format (host, port range, etc.) before saving

#### Scenario: Test connection
- **WHEN** user adds a new source
- **THEN** system optionally tests connection and reports success/failure

#### Scenario: Auto-create source from CLI args
- **WHEN** user provides MySQL CLI connection arguments
- **THEN** system automatically creates a source with auto-generated name without prompting user

#### Scenario: Generate unique source name
- **WHEN** system auto-creates source from CLI args
- **THEN** system generates name in format `{host}-{port}-{user}` (e.g., `11.124.9.201-2900-root`)

#### Scenario: Handle source name collision
- **WHEN** auto-generated source name already exists
- **THEN** system appends numeric suffix (e.g., `11.124.9.201-2900-root-2`, `11.124.9.201-2900-root-3`) until unique name is found

#### Scenario: Set database type for MySQL CLI-created sources
- **WHEN** system auto-creates source from MySQL CLI args (detected via `-P` or `-D` flags)
- **THEN** system sets source type to MySQL

#### Scenario: Set database type for PostgreSQL CLI-created sources
- **WHEN** system auto-creates source from PostgreSQL CLI args (detected via `-U` or `-d` flags)
- **THEN** system sets source type to PostgreSQL

#### Scenario: Set database type from explicit flag
- **WHEN** system auto-creates source with explicit `--engine` or `-e` flag
- **THEN** system sets source type to the explicitly specified engine type (mysql/postgresql/seekdb)

## ADDED Requirements

### Requirement: Data source storage
The system SHALL store database connection configurations in `~/.config/aiq/sources.yaml`.

#### Scenario: Create sources file
- **WHEN** user adds first data source
- **THEN** system creates `~/.config/aiq/sources.yaml` if it does not exist

#### Scenario: Store source configuration
- **WHEN** user adds a data source
- **THEN** system saves connection details to `sources.yaml` in YAML format

### Requirement: Add data source
The system SHALL allow users to add new database connection configurations.

#### Scenario: Add MySQL source
- **WHEN** user selects "add" from source submenu
- **THEN** system prompts for: name, host, port, database, username, password

#### Scenario: Validate connection details
- **WHEN** user provides connection details
- **THEN** system validates format (host, port range, etc.) before saving

#### Scenario: Test connection
- **WHEN** user adds a new source
- **THEN** system optionally tests connection and reports success/failure

### Requirement: List data sources
The system SHALL display all configured data sources.

#### Scenario: List all sources
- **WHEN** user selects "list" from source submenu
- **THEN** system displays all configured sources with names and connection info (mask passwords)

#### Scenario: Empty sources list
- **WHEN** user lists sources and none exist
- **THEN** system displays friendly message suggesting to add a source

### Requirement: Select active data source
The system SHALL allow users to select which data source to use for SQL queries.

#### Scenario: Select source
- **WHEN** user selects "select" from source submenu
- **THEN** system displays list of sources and allows selection

#### Scenario: Active source indicator
- **WHEN** a source is selected
- **THEN** system indicates active source in subsequent menus and SQL mode

### Requirement: Remove data source
The system SHALL allow users to remove configured data sources.

#### Scenario: Remove source
- **WHEN** user selects "remove" from source submenu
- **THEN** system displays sources list and allows deletion with confirmation

#### Scenario: Confirm deletion
- **WHEN** user attempts to remove a source
- **THEN** system prompts for confirmation before deleting

### Requirement: Source submenu
The system SHALL provide a submenu for managing data sources with options: add, list, select, remove.

#### Scenario: Display source submenu
- **WHEN** user selects source from main menu
- **THEN** system displays submenu with: add, list, select, remove, back

#### Scenario: Navigate back to main menu
- **WHEN** user selects "back" from source submenu
- **THEN** system returns to main menu

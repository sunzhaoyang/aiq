## MODIFIED Requirements

### Requirement: Source submenu
The system SHALL provide a submenu for managing data sources with options: add, list, edit, remove, back.

#### Scenario: Display source submenu
- **WHEN** user selects source from main menu
- **THEN** system displays submenu with: add, list, edit, remove, back

#### Scenario: Navigate back to main menu
- **WHEN** user selects "back" from source submenu
- **THEN** system returns to main menu

### Requirement: Edit data source
The system SHALL allow users to edit existing data source configurations.

#### Scenario: Edit source from submenu
- **WHEN** user selects "edit" from source submenu
- **THEN** system displays list of sources and allows selection

#### Scenario: Edit source fields
- **WHEN** user selects a source to edit
- **THEN** system prompts for all source fields (name, host, port, database, username, password) with current values as defaults

#### Scenario: Update source with modified fields
- **WHEN** user provides modified source configuration
- **THEN** system validates the configuration and updates the source

#### Scenario: Validate unique name on edit
- **WHEN** user edits source and changes the name to an existing name
- **THEN** system displays error message and prevents update

#### Scenario: Validate unique connection on edit
- **WHEN** user edits source and changes host/port/username to match an existing source
- **THEN** system displays error message and prevents update (source uniqueness is based on host-port-username)

#### Scenario: Test connection after edit
- **WHEN** user edits a source
- **THEN** system optionally tests connection and reports success/failure before saving

#### Scenario: Confirm update
- **WHEN** user completes editing a source
- **THEN** system saves the updated configuration and displays success message

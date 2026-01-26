## ADDED Requirements

### Requirement: CLI application entry point
The system SHALL provide a command-line executable named `aiq` that launches an interactive CLI interface.

#### Scenario: Launch interactive CLI
- **WHEN** user runs `aiq` command
- **THEN** system displays the main menu with available options

### Requirement: Main menu system
The system SHALL display a main menu with four core functions: `config`, `source`, `sql`, and `exit`.

#### Scenario: Display main menu
- **WHEN** CLI application starts
- **THEN** system shows menu with options: config, source, sql, exit

#### Scenario: Navigate to config menu
- **WHEN** user selects `config` from main menu
- **THEN** system displays configuration management submenu

#### Scenario: Navigate to source menu
- **WHEN** user selects `source` from main menu
- **THEN** system displays data source management submenu

#### Scenario: Navigate to SQL mode
- **WHEN** user selects `sql` from main menu
- **THEN** system prompts user to select a data source (if none selected) and enters SQL interactive mode

#### Scenario: Exit application
- **WHEN** user selects `exit` from main menu
- **THEN** system gracefully exits with exit code 0

### Requirement: Command routing
The system SHALL route user menu selections to appropriate command handlers.

#### Scenario: Route config command
- **WHEN** user selects config option
- **THEN** system invokes configuration management handler

#### Scenario: Route source command
- **WHEN** user selects source option
- **THEN** system invokes data source management handler

#### Scenario: Route sql command
- **WHEN** user selects sql option
- **THEN** system invokes SQL interactive mode handler

### Requirement: Interactive prompt system
The system SHALL use interactive prompts for menu selection with search and navigation capabilities.

#### Scenario: Select menu option
- **WHEN** menu is displayed
- **THEN** user can navigate options using arrow keys and select with Enter

#### Scenario: Search menu options
- **WHEN** menu is displayed
- **THEN** user can type to filter/search menu options

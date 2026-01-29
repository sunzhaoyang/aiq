## MODIFIED Requirements

### Requirement: Main menu system
The system SHALL display a main menu with core functions and support navigation back from all sub-functions.

#### Scenario: Display main menu
- **WHEN** CLI application starts
- **THEN** system shows menu with options: config, source, chat, exit

#### Scenario: Navigate to config menu
- **WHEN** user selects `config` from main menu
- **THEN** system displays configuration management submenu

#### Scenario: Navigate to source menu
- **WHEN** user selects `source` from main menu
- **THEN** system displays data source management submenu

#### Scenario: Navigate to chat mode
- **WHEN** user selects `chat` from main menu
- **THEN** system prompts user to select a data source (if none selected) and enters chat mode

#### Scenario: Return from chat mode to main menu
- **WHEN** user exits chat mode (via `/exit`, `exit`, `back`, or Ctrl+D)
- **THEN** system returns to main menu instead of exiting the application

#### Scenario: Exit application
- **WHEN** user selects `exit` from main menu
- **THEN** system gracefully exits with exit code 0

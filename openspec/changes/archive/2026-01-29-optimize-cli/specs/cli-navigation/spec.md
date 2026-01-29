## ADDED Requirements

### Requirement: Unified CLI navigation and exit mechanism
The system SHALL provide a unified navigation mechanism that allows users to return to the main menu from any sub-function, including chat mode.

#### Scenario: Return from chat mode to main menu
- **WHEN** user exits chat mode (via `/exit`, `exit`, `back`, or Ctrl+D)
- **THEN** system returns to main menu instead of exiting the application

#### Scenario: Return from source menu to main menu
- **WHEN** user selects "back" from source submenu
- **THEN** system returns to main menu

#### Scenario: Return from config menu to main menu
- **WHEN** user selects "back" from config submenu
- **THEN** system returns to main menu

#### Scenario: Consistent navigation experience
- **WHEN** user navigates between different functions
- **THEN** all functions provide consistent "back" or "exit" options to return to main menu

#### Scenario: Exit application from main menu
- **WHEN** user selects "exit" from main menu
- **THEN** system gracefully exits the application with exit code 0

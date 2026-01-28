## MODIFIED Requirements

### Requirement: CLI application entry point
The system SHALL provide a command-line executable named `aiq` that launches an interactive CLI interface, or directly enters chat mode when MySQL-compatible connection arguments are provided.

#### Scenario: Launch interactive CLI
- **WHEN** user runs `aiq` command without arguments
- **THEN** system displays the main menu with available options

#### Scenario: Launch with MySQL CLI args
- **WHEN** user runs `aiq` with MySQL-compatible connection arguments (e.g., `-h host -u user -P port -ppassword -D database`, note: `-p` and password have no space)
- **THEN** system bypasses main menu and directly enters chat mode after validating connection and creating source

#### Scenario: Launch with PostgreSQL CLI args
- **WHEN** user runs `aiq` with PostgreSQL-compatible connection arguments (e.g., `-h host -U user -p port -d database` with `PGPASSWORD` env var)
- **THEN** system bypasses main menu and directly enters chat mode after validating connection and creating source

#### Scenario: Launch with session file flag
- **WHEN** user runs `aiq -s session.json` or `aiq --session session.json`
- **THEN** system restores session and enters chat mode (existing behavior)

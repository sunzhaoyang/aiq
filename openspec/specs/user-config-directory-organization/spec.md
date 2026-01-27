## ADDED Requirements

### Requirement: Structured config directory hierarchy
The system SHALL organize user configuration files into a structured directory hierarchy under `~/.aiqconfig/`.

#### Scenario: Directory structure creation
- **WHEN** application starts and config directory does not exist
- **THEN** system creates `~/.aiqconfig/` with subdirectories: `config/`, `sessions/`, `skills/`, `tools/`

#### Scenario: Config files location
- **WHEN** system saves configuration files
- **THEN** LLM config is saved to `~/.aiqconfig/config/config.yaml` and data sources config is saved to `~/.aiqconfig/config/sources.yaml`

#### Scenario: Session files location
- **WHEN** system saves session files
- **THEN** session files are saved to `~/.aiqconfig/sessions/session_<timestamp>.json`

#### Scenario: Skills directory location
- **WHEN** system loads Skills
- **THEN** Skills are loaded from `~/.aiqconfig/skills/<skill-name>/SKILL.md`

### Requirement: Config file path resolution
The system SHALL provide path resolution functions for all config directory paths.

#### Scenario: Get config file path
- **WHEN** system needs to access LLM configuration
- **THEN** system resolves path to `~/.aiqconfig/config/config.yaml`

#### Scenario: Get sources file path
- **WHEN** system needs to access data sources configuration
- **THEN** system resolves path to `~/.aiqconfig/config/sources.yaml`

#### Scenario: Get sessions directory path
- **WHEN** system needs to save or load session files
- **THEN** system resolves path to `~/.aiqconfig/sessions/`

#### Scenario: Get skills directory path
- **WHEN** system needs to load Skills
- **THEN** system resolves path to `~/.aiqconfig/skills/`

### Requirement: Directory structure initialization
The system SHALL create required directory structure on first run if it does not exist.

#### Scenario: Create directories on startup
- **WHEN** application starts and `~/.aiqconfig/` does not exist
- **THEN** system creates all required subdirectories (`config/`, `sessions/`, `skills/`, `tools/`)

#### Scenario: Handle existing directories
- **WHEN** application starts and some directories already exist
- **THEN** system creates only missing directories without error

#### Scenario: Directory permissions
- **WHEN** system creates directories
- **THEN** directories are created with permissions 0755

### Requirement: Path constants and utilities
The system SHALL provide constants and utility functions for config directory paths.

#### Scenario: Access directory constants
- **WHEN** code needs config directory paths
- **THEN** system provides constants for directory names and path resolution functions

#### Scenario: Cross-platform path handling
- **WHEN** system runs on different operating systems
- **THEN** path resolution functions handle platform-specific path separators correctly

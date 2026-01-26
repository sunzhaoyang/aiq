## ADDED Requirements

### Requirement: First-run configuration wizard
The system SHALL detect first launch and guide user through LLM configuration setup.

#### Scenario: Detect first launch
- **WHEN** application starts and configuration file does not exist
- **THEN** system automatically launches first-run configuration wizard

#### Scenario: Collect LLM URL
- **WHEN** first-run wizard starts
- **THEN** system prompts user to enter LLM API URL

#### Scenario: Collect LLM API key
- **WHEN** user provides LLM URL
- **THEN** system prompts user to enter LLM API key (with masked input)

#### Scenario: Save configuration
- **WHEN** user completes first-run wizard
- **THEN** system saves configuration to `~/.config/aiq/config.yaml`

### Requirement: Configuration file management
The system SHALL store LLM configuration in YAML format at `~/.config/aiq/config.yaml`.

#### Scenario: Read configuration
- **WHEN** application starts and configuration file exists
- **THEN** system loads configuration from `~/.config/aiq/config.yaml`

#### Scenario: Configuration file format
- **WHEN** configuration is saved
- **THEN** file contains LLM URL and API key in YAML format

### Requirement: Configuration submenu
The system SHALL provide a submenu for managing LLM and tool configuration.

#### Scenario: View current configuration
- **WHEN** user selects "view" from config submenu
- **THEN** system displays current LLM URL (mask API key) and other settings

#### Scenario: Update LLM URL
- **WHEN** user selects "update URL" from config submenu
- **THEN** system prompts for new URL and updates configuration

#### Scenario: Update LLM API key
- **WHEN** user selects "update API key" from config submenu
- **THEN** system prompts for new API key (masked) and updates configuration

#### Scenario: Validate configuration
- **WHEN** user updates configuration
- **THEN** system validates format and provides error feedback if invalid

### Requirement: Configuration persistence
The system SHALL persist configuration changes immediately.

#### Scenario: Save on update
- **WHEN** user updates any configuration value
- **THEN** system immediately writes changes to configuration file

#### Scenario: Handle configuration errors
- **WHEN** configuration file is corrupted or invalid
- **THEN** system displays clear error message and offers to reset configuration

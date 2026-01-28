## MODIFIED Requirements

### Requirement: SQL interactive mode entry
The system SHALL enter SQL interactive mode when user selects `chat` from main menu, with optional source selection, or directly when MySQL CLI connection arguments are provided.

#### Scenario: Enter chat mode without sources configured
- **WHEN** user selects chat and no data sources are configured
- **THEN** system enters free chat mode without requiring source selection

#### Scenario: Enter chat mode with source selection prompt
- **WHEN** user selects chat and sources are available
- **THEN** system prompts user to select a source or skip to enter free mode

#### Scenario: Enter chat mode with selected source
- **WHEN** user selects a source from available sources
- **THEN** system enters chat mode with selected source and database connection

#### Scenario: Enter chat mode in free mode
- **WHEN** user chooses to skip source selection
- **THEN** system enters free chat mode without database connection

#### Scenario: Direct entry from MySQL CLI args
- **WHEN** user provides MySQL CLI connection arguments
- **THEN** system validates connection, creates source automatically, and directly enters chat mode with the newly created source (skipping source selection menu)

#### Scenario: Direct entry from PostgreSQL CLI args
- **WHEN** user provides PostgreSQL CLI connection arguments
- **THEN** system validates connection, creates source automatically, and directly enters chat mode with the newly created source (skipping source selection menu)

#### Scenario: LLM config check before direct entry
- **WHEN** user provides MySQL CLI args but LLM is not configured
- **THEN** system prompts user to configure LLM first before entering chat mode

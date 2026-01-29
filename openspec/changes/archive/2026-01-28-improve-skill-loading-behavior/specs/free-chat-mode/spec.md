## MODIFIED Requirements

### Requirement: Free chat mode without database source
The system SHALL support chat mode operation without requiring a database source to be selected, with enhanced input validation and mode awareness.

#### Scenario: Enter chat mode without sources
- **WHEN** user selects chat mode and no data sources are configured
- **THEN** system enters free chat mode without requiring source selection

#### Scenario: Enter chat mode with source selection skipped
- **WHEN** user selects chat mode and chooses to skip source selection
- **THEN** system enters free chat mode

#### Scenario: Free mode capabilities
- **WHEN** system is in free chat mode
- **THEN** system supports general conversation, Skills-based operations (execute_command, http_request, file_operations), but NOT SQL execution (execute_sql tool unavailable)

#### Scenario: Free mode prompt indication
- **WHEN** system is in free chat mode
- **THEN** system displays prompt as `aiq> ` (without source name)

#### Scenario: Validate database queries in free mode
- **WHEN** user asks database-related query (show tables, query data, list tables) in free chat mode
- **THEN** system recognizes query as invalid for free mode and asks clarifying question: "I notice you're asking about database operations, but no database is currently connected. Would you like to select a database source to enable SQL queries?"

#### Scenario: SQL execution attempt in free mode
- **WHEN** user attempts to execute SQL in free chat mode
- **THEN** system displays error message indicating that a database source is required for SQL operations

#### Scenario: Prevent command guessing for invalid queries
- **WHEN** user asks database query in free mode
- **THEN** system does NOT attempt to execute shell commands (mysql, psql) as workaround, but instead asks clarifying question

#### Scenario: Switch from free mode to database mode
- **WHEN** user is in free chat mode and wants to use database features
- **THEN** user must exit chat mode and re-enter with source selection (no mid-session switching)

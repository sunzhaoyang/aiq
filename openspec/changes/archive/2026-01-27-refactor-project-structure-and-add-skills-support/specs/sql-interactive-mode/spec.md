## ADDED Requirements

### Requirement: SQL interactive mode entry
The system SHALL enter SQL interactive mode when user selects `sql` from main menu.

#### Scenario: Enter SQL mode without selected source
- **WHEN** user selects sql and no source is selected
- **THEN** system prompts user to select a data source first

#### Scenario: Enter SQL mode with selected source
- **WHEN** user selects sql and a source is already selected
- **THEN** system enters SQL interactive mode with selected source

### Requirement: Natural language to SQL translation
The system SHALL translate user's natural language questions into SQL queries using LLM.

#### Scenario: Submit natural language query
- **WHEN** user enters a natural language question in SQL mode
- **THEN** system sends question to LLM API with database schema context

#### Scenario: Receive SQL translation
- **WHEN** LLM returns SQL query
- **THEN** system displays translated SQL query to user

#### Scenario: Confirm SQL execution
- **WHEN** SQL is translated
- **THEN** system prompts user to confirm execution or modify query

### Requirement: SQL query execution
The system SHALL execute SQL queries against the selected database connection.

#### Scenario: Execute confirmed query
- **WHEN** user confirms SQL execution
- **THEN** system executes query against selected database

#### Scenario: Display query results
- **WHEN** query executes successfully
- **THEN** system displays results in formatted table with syntax highlighting

#### Scenario: Handle query errors
- **WHEN** query execution fails
- **THEN** system displays clear error message and allows user to retry or modify

### Requirement: SQL mode interface
The system SHALL provide an interactive interface for SQL queries with prompt and command handling.

#### Scenario: Display SQL prompt
- **WHEN** SQL mode is active
- **THEN** system displays prompt indicating SQL mode and active source

#### Scenario: Accept multi-line input
- **WHEN** user enters SQL query
- **THEN** system accepts multi-line input until user submits (e.g., Ctrl+D or special command)

#### Scenario: Exit SQL mode
- **WHEN** user types `exit` or `back` in SQL mode
- **THEN** system returns to main menu

### Requirement: Database schema context
The system SHALL provide database schema information to LLM for accurate SQL generation.

#### Scenario: Fetch schema on source selection
- **WHEN** user selects a data source
- **THEN** system optionally fetches schema information (tables, columns) for context

#### Scenario: Include schema in LLM request
- **WHEN** translating natural language to SQL
- **THEN** system includes relevant schema information in LLM API request

## MODIFIED Requirements

### Requirement: Natural language to SQL translation
The system SHALL translate user's natural language questions into SQL queries using LLM, with Skills-enhanced prompts.

#### Scenario: Submit natural language query
- **WHEN** user enters a natural language question in SQL mode
- **THEN** system sends question to LLM API with database schema context, matched Skills content, and available tools

#### Scenario: Match Skills to query
- **WHEN** user submits a natural language query
- **THEN** system matches query against Skills metadata and loads relevant Skills content

#### Scenario: Build prompt with Skills
- **WHEN** system prepares prompt for LLM translation
- **THEN** system includes matched Skills content in system prompt section, ordered by priority

#### Scenario: Manage prompt length
- **WHEN** prompt token count exceeds thresholds
- **THEN** system compresses conversation history and evicts low-priority Skills to stay within token limits

#### Scenario: Receive SQL translation
- **WHEN** LLM returns SQL query
- **THEN** system displays translated SQL query to user

#### Scenario: Confirm SQL execution
- **WHEN** SQL is translated
- **THEN** system prompts user to confirm execution or modify query

## MODIFIED Requirements

### Requirement: SQL mode interface
The system SHALL provide an interactive interface for SQL queries with prompt, command handling, and tab completion support.

#### Scenario: Display SQL prompt with source
- **WHEN** SQL mode is active with a database source
- **THEN** system displays prompt as `aiq[source-name]> ` indicating active source

#### Scenario: Display SQL prompt with source and database
- **WHEN** SQL mode is active with a database source and database name is available (from source config or `-D` override)
- **THEN** system displays prompt as `aiq[source-name database-name]> ` indicating active source and database being used

#### Scenario: Display prompt in free mode
- **WHEN** SQL mode is active without database source (free mode)
- **THEN** system displays prompt as `aiq> ` (without source name)

#### Scenario: Accept multi-line input
- **WHEN** user enters SQL query or general query
- **THEN** system accepts multi-line input until user submits (e.g., Ctrl+D or special command)

#### Scenario: Exit SQL mode with /exit command
- **WHEN** user types `/exit` in chat mode
- **THEN** system saves session and returns to main menu


#### Scenario: Display help with /help command
- **WHEN** user types `/help` in chat mode
- **THEN** system displays list of available commands and their descriptions

#### Scenario: Help command shows available commands
- **WHEN** user types `/help`
- **THEN** system displays:
  - `/exit` - Exit chat mode and return to main menu
  - `/help` - Show this help message
  - `/history` - View conversation history
  - `/clear` - Clear conversation history

#### Scenario: Tab completion for commands
- **WHEN** user types `/` and presses Tab
- **THEN** system displays available commands: `/exit`, `/help`, `/history`, `/clear`

#### Scenario: Tab completion completes command name
- **WHEN** user types `/ex` and presses Tab
- **THEN** system completes to `/exit`

#### Scenario: Tab completion only for commands
- **WHEN** user types natural language query (not starting with `/`) and presses Tab
- **THEN** system does not provide completion (to avoid interfering with natural language input)

#### Scenario: Return to main menu on exit
- **WHEN** user exits chat mode (via `/exit`, `exit`, `back`, or Ctrl+D)
- **THEN** system returns to main menu instead of exiting the application

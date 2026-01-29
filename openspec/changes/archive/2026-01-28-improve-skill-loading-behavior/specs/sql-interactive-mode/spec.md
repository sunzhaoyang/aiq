## MODIFIED Requirements

### Requirement: Natural language to SQL translation
The system SHALL translate user's natural language questions into SQL queries using LLM, with Skills-enhanced prompts, supporting both database mode and free mode, with enhanced mode awareness and tool usage understanding.

#### Scenario: Submit natural language query in database mode
- **WHEN** user enters a natural language question in chat mode with database source
- **THEN** system sends question to LLM API with database schema context, matched Skills content, and available tools (including execute_sql), with explicit guidance that execute_sql tool should be used for database queries

#### Scenario: Submit natural language query in free mode
- **WHEN** user enters a natural language question in free chat mode
- **THEN** system sends question to LLM API with matched Skills content and available tools (excluding execute_sql), with explicit guidance that database queries are invalid in free mode

#### Scenario: Match Skills to query
- **WHEN** user submits a natural language query
- **THEN** system matches query against Skills metadata using LLM semantic matching with precision filtering and loads relevant Skills content

#### Scenario: Build prompt with Skills
- **WHEN** system prepares prompt for LLM translation
- **THEN** system includes matched Skills content in system prompt section, ordered by priority

#### Scenario: Manage prompt length
- **WHEN** prompt token count exceeds thresholds
- **THEN** system uses LLM semantic compression for conversation history and evicts low-priority Skills to stay within token limits

#### Scenario: Receive SQL translation
- **WHEN** LLM returns SQL query
- **THEN** system displays translated SQL query to user

#### Scenario: Confirm SQL execution
- **WHEN** SQL is translated
- **THEN** system prompts user to confirm execution or modify query

#### Scenario: LLM understands tool availability
- **WHEN** system is in database mode and user asks database query
- **THEN** LLM uses execute_sql tool, NOT execute_command with mysql/psql commands

#### Scenario: LLM validates query appropriateness
- **WHEN** user submits query in free mode that requires database
- **THEN** LLM recognizes query as invalid for current mode and asks clarifying question instead of attempting tool execution

#### Scenario: LLM distinguishes SQL tool from shell commands
- **WHEN** user asks database query in database mode
- **THEN** LLM uses execute_sql tool for SQL queries, and only uses execute_command for system operations (not database queries)

## ADDED Requirements

### Requirement: Command execution display optimization
The system SHALL display command execution output in a user-friendly manner that avoids screen flooding, provides real-time feedback, and intelligently truncates output for LLM processing.

#### Scenario: Display tool call with gray text
- **WHEN** LLM calls `execute_command` tool
- **THEN** system displays tool call information (command name and arguments) using gray/dimmed text color

#### Scenario: Show real-time output in rolling window
- **WHEN** command is executing and producing output
- **THEN** system displays output in a rolling window of 2-3 lines near the current cursor position, updating in real-time without scrolling the entire screen

#### Scenario: Avoid screen flooding
- **WHEN** command produces large amounts of output
- **THEN** system does not flood the screen with all output lines, but maintains the rolling window display

#### Scenario: Truncate successful command output for LLM
- **WHEN** command executes successfully (exit code 0) and produces output
- **THEN** system truncates output to last 10-20 lines before sending to LLM, reducing token consumption

#### Scenario: Include more output for failed commands
- **WHEN** command execution fails (exit code non-zero) and produces output
- **THEN** system truncates output to last 50-100 lines before sending to LLM, providing sufficient context for error diagnosis

#### Scenario: Preserve full output metadata
- **WHEN** command execution completes
- **THEN** system preserves full output (stdout, stderr, exit_code) internally, even if truncated output is sent to LLM

#### Scenario: Display execution status
- **WHEN** command execution completes
- **THEN** system displays clear success or failure status message to user

#### Scenario: Handle long-running commands
- **WHEN** command takes significant time to execute
- **THEN** system shows loading indicator or progress feedback while maintaining rolling output window

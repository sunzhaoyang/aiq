## Why

The current default timeout for command execution is 30 seconds, which is too short for some commands that require longer execution time (like compilation, downloads, database operations, etc.). When commands timeout, the system directly returns an error, and users cannot choose to continue waiting, causing need to re-execute commands and affecting user experience. Increasing default timeout to 60 seconds and asking users whether to continue waiting after timeout can significantly improve user experience for long-running commands.

## What Changes

- **Increase Default Timeout**: Increase default timeout of `execute_command` tool from 30 seconds to 60 seconds
- **User Prompt After Timeout**: When command execution times out, no longer directly return error, but ask user whether to continue waiting
- **User Choice Handling**: If user chooses to continue waiting, reset timeout timer and continue execution; if chooses to cancel, return timeout error

## Capabilities

### New Capabilities
- `command-timeout-user-prompt`: Ability to ask user whether to continue waiting after command execution timeout

### Modified Capabilities
- `command-execution-display`: Modify default timeout from 30 seconds to 60 seconds, and add user prompt scenario after timeout

## Impact

**Affected Code:**
- `internal/tool/builtin/command_tool.go`: Modify default timeout, add user interaction logic after timeout
- `internal/sql/tool_handler.go`: May need to handle user interaction for timeout prompts
- `internal/ui/`: May need to add UI function for user confirmation prompts

**User Experience:**
- Reduce command execution failures due to timeout
- Provide better support for long-running commands
- Users can choose whether to continue waiting instead of being forced to re-execute

**Backward Compatibility:**
- Change in default timeout won't affect existing code, only default value changes
- Users can still customize timeout via `timeout` parameter

## Why

When LLM calls the `execute_command` tool, it sometimes returns double-encoded JSON strings (e.g., `"{\"command\":\"brew services list | grep mysql\"}"`) instead of direct JSON objects. The current `ParseArguments()` function attempts to handle this situation but still fails in some cases, causing tool call failures and displaying error message "json: cannot unmarshal string into Go value of type map[string]interface {}". This affects user experience, especially when Skills use the `execute_command` tool.

## What Changes

- **Fix JSON Parameter Parsing Logic**: Improve `ParseArguments()` function to more robustly handle double-encoded JSON strings
- **Enhance Error Handling**: Improve error handling logic in `tool_handler.go`, provide clearer error messages
- **Add Test Cases**: Add test cases for double-encoded JSON strings to ensure the fix is effective

## Capabilities

### New Capabilities
<!-- No new capabilities -->

### Modified Capabilities
- `tool-execution`: Improve robustness of tool parameter parsing, ensure correct handling of various JSON formats returned by LLM

## Impact

**Affected Code:**
- `internal/llm/client.go`: `ParseArguments()` function needs to improve JSON parsing logic
- `internal/sql/tool_handler.go`: Error handling logic in `ExecuteTool()` function may need adjustment

**Testing:**
- Need to add test cases to verify parsing of double-encoded JSON strings
- Need to verify various edge cases (empty strings, invalid JSON, nested quotes, etc.)

**User Experience:**
- After fix, JSON parameter formats returned by LLM can be correctly parsed even if not standard enough
- Reduce tool call failures, improve system stability

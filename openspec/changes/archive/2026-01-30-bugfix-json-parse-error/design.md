## Context

When the current system processes tool call parameters returned by LLM, the `ParseArguments()` function attempts to parse JSON strings. Although handling for double-encoded JSON strings has been implemented (detect and unquote outer quotes), it still fails in some cases.

**Current Implementation Issues:**
1. `ParseArguments()` only handles one layer of quote wrapping, if there are multiple nested JSON strings, may not handle correctly
2. Error handling logic is scattered in two places: `ParseArguments()` and `tool_handler.go`, causing code duplication
3. Error messages are not clear enough, difficult to diagnose issues

**Error Scenario Example:**
- LLM returns: `"{\"command\":\"brew services list | grep mysql\"}"`
- `ParseArguments()` attempts to unquote and gets: `{\"command\":\"brew services list | grep mysql\"}`
- But parsing may fail due to improper escape character handling

## Goals / Non-Goals

**Goals:**
- Improve `ParseArguments()` function to correctly handle various JSON parameter formats (including multi-layer encoding, escape characters, etc.)
- Unify error handling logic, reduce code duplication
- Provide clearer error messages for easier debugging
- Add comprehensive test cases covering various edge cases

**Non-Goals:**
- Do not change tool call API interface
- Do not modify other LLM client functionality
- Do not change tool definition structure

## Decisions

### 1. Improve JSON Parsing Strategy

**Decision**: Use recursive approach to handle multi-layer encoded JSON strings until unable to unquote further.

**Rationale**:
- Current implementation only handles one layer of quotes, cannot handle multi-layer nested cases
- Recursive processing ensures correct parsing regardless of nesting depth

**Implementation**:
```go
func (tc *ToolCall) ParseArguments() (map[string]interface{}, error) {
    argsStr := tc.Function.Arguments
    
    // Recursively unquote until unable to unquote further
    for {
        trimmed := strings.TrimSpace(argsStr)
        if len(trimmed) < 2 || trimmed[0] != '"' || trimmed[len(trimmed)-1] != '"' {
            break
        }
        var unquoted string
        if err := json.Unmarshal([]byte(trimmed), &unquoted); err != nil {
            break
        }
        argsStr = unquoted
    }
    
    // Attempt to parse as JSON object
    var args map[string]interface{}
    if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
        return nil, fmt.Errorf("failed to parse arguments: %w", err)
    }
    return args, nil
}
```

**Alternatives Considered**:
- **Option A**: Only handle one layer of quotes (current approach) - Cannot handle multi-layer nesting
- **Option B**: Use regex matching - Not robust enough, may mis-match
- **Option C**: Recursive unquote (chosen) - Most robust, can handle any number of layers

### 2. Simplify Error Handling Logic

**Decision**: Centralize error handling logic in `ParseArguments()`, remove duplicate handling in `tool_handler.go`.

**Rationale**:
- Current error handling logic is scattered in two places, causing code duplication and maintenance difficulties
- Unified handling logic improves code maintainability

**Implementation**:
- Handle all JSON parsing-related errors in `ParseArguments()`
- Only handle tool execution-related errors in `tool_handler.go`

**Alternatives Considered**:
- **Option A**: Keep current dual handling - Code duplication, difficult to maintain
- **Option B**: Unify to `ParseArguments()` (chosen) - Clearer, easier to maintain

### 3. Improve Error Messages

**Decision**: Include original parameters and parsing steps in error messages for easier debugging.

**Rationale**:
- Current error messages are not detailed enough, difficult to diagnose issues
- Including more context information helps quickly locate problems

**Implementation**:
```go
return nil, fmt.Errorf("failed to parse arguments after unquoting: %w (original: %s)", err, truncateString(tc.Function.Arguments, 100))
```

## Risks / Trade-offs

**Risk 1**: Recursive unquote may cause infinite loop
- **Mitigation**: Add maximum recursion depth limit (e.g., 10 layers)

**Risk 2**: Excessive processing may cause performance issues
- **Mitigation**: Recursion depth limit prevents performance issues, and most cases only need to handle 1-2 layers

**Risk 3**: May mis-parse certain special format strings
- **Mitigation**: After unquoting, verify if it's a valid JSON object, stop processing if not

**Trade-offs**:
- **Robustness vs Performance**: Choose robustness, because tool call failure impact far exceeds tiny parsing performance overhead
- **Simplicity vs Completeness**: Choose completeness, ensure handling of various edge cases

## Migration Plan

**Deployment Steps:**
1. Implement improved `ParseArguments()` function
2. Add test cases covering various edge cases
3. Run existing tests to ensure no regression
4. Submit code and merge to main branch

**Rollback Strategy:**
- If issues found, can quickly rollback to previous version
- Improvements are backward compatible, won't affect existing functionality

**Testing Strategy:**
- Unit tests: Test parsing of various JSON formats
- Integration tests: Test actual tool call scenarios
- Edge tests: Test empty strings, invalid JSON, multi-layer nesting, etc.

## Open Questions

1. **Should we support other parameter formats?** For example YAML, TOML, etc.
   - **Decision**: Not supported for now, only support JSON format

2. **Should we log parsing failures?** For monitoring and debugging
   - **Decision**: Not logged for now, error messages are detailed enough

3. **Should we add configuration options?** For example maximum recursion depth, whether to allow certain formats, etc.
   - **Decision**: Not needed for now, use reasonable defaults

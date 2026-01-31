## Context

The current `execute_command` tool uses idle timeout mechanism: when command has no output within 30 seconds, it times out and directly terminates the command. This mechanism is too strict for commands that run for a long time but occasionally have output (like compilation, downloads, etc.), causing commands to be prematurely terminated.

**Current Implementation:**
- Default idle timeout: 30 seconds
- Timeout handling: Directly kill process and return error
- Users cannot intervene in timeout decisions

**Constraints:**
- Command execution is asynchronous, needs to handle concurrent scenarios
- UI interaction needs to occur during command execution, cannot block output stream
- Need to maintain backward compatibility

## Goals / Non-Goals

**Goals:**
- Increase default idle timeout from 30 seconds to 60 seconds
- Ask user whether to continue waiting after timeout
- Support users extending timeout multiple times (each extension resets timer)
- Maintain existing idle timeout mechanism (reset timer when there's output)

**Non-Goals:**
- Do not change basic idle timeout mechanism (still resets when there's output)
- Do not add configuration options to disable timeout prompts (keep it simple)
- Do not change timeout mechanisms of other tools (like http_request)

## Decisions

### 1. Timeout Prompt Interaction Method

**Decision**: Use existing `ui.ShowConfirm()` function for user confirmation, prompt message clearly explains current situation.

**Rationale**:
- `ui.ShowConfirm()` already exists and is proven
- Simple yes/no choice is sufficient for needs
- No additional UI components needed

**Implementation**:
```go
// When timeout occurs
continueWaiting, err := ui.ShowConfirm(
    fmt.Sprintf("Command has been idle for %v. Continue waiting?", idleTimeout))
if err != nil {
    // User interruption (Ctrl+C), treat as cancel
    cmd.Process.Kill()
    return nil, fmt.Errorf("command execution cancelled by user")
}
if !continueWaiting {
    // User chooses not to continue
    cmd.Process.Kill()
    return nil, fmt.Errorf("command execution timeout: no output for %v", idleTimeout)
}
// User chooses to continue, reset timer
idleTimer.Reset(idleTimeout)
```

**Alternatives Considered**:
- **Option A**: Automatically extend fixed time (e.g., wait another 30 seconds) - Not flexible enough, user cannot control
- **Option B**: Ask user and allow inputting new timeout time - Too complex, increases user burden
- **Option C**: Ask user whether to continue (chosen) - Simple and direct, meets needs

### 2. Timeout Prompt Display Timing

**Decision**: Immediately display prompt when idle timeout is detected, pause and wait for user response.

**Rationale**:
- Timeout means command may be stuck, timely asking user can avoid unnecessary waiting
- Pausing wait doesn't affect normal command execution (command still runs in background)

**Implementation**:
- In `case <-idleTimer.C:` branch, first stop timer, then display confirmation prompt
- Command continues running in background while waiting for user response
- If user chooses to continue, reset timer; if chooses to cancel, terminate command

**Alternatives Considered**:
- **Option A**: Warn a few seconds in advance - Increases complexity, may false alarm
- **Option B**: Immediately prompt after timeout (chosen) - Simple and direct

### 3. Handling Multiple Timeout Extensions

**Decision**: Allow users to extend timeout multiple times, each extension resets timer.

**Rationale**:
- Some commands may need very long time, allowing multiple extensions provides better flexibility
- Simple implementation, only need to reset timer when user chooses to continue

**Implementation**:
- Ask user every time timeout occurs
- After user chooses to continue, reset timer and continue waiting
- No maximum extension count limit (user can terminate by choosing "no")

**Alternatives Considered**:
- **Option A**: Limit maximum extension count - May not be flexible enough
- **Option B**: Allow unlimited extensions (chosen) - More flexible, user controllable

### 4. Concurrency Safety Handling

**Decision**: When displaying confirmation prompt, need to ensure command output doesn't interfere with prompt display.

**Rationale**:
- Command may produce output during user response
- Need to avoid output and prompt mixing together

**Implementation**:
- Before displaying prompt, pause output callback (if any)
- Or use independent output channel to ensure prompt displays on separate line
- Using `ui.ShowConfirm()` already handles input/output isolation

**Alternatives Considered**:
- **Option A**: Pause output stream - May lose output
- **Option B**: Use independent line to display prompt (chosen) - Doesn't interfere with output, clearer

### 5. Default Timeout Time Modification

**Decision**: Change default idle timeout from 30 seconds to 60 seconds.

**Rationale**:
- 60 seconds is more reasonable for most commands
- Users can still customize via `timeout` parameter

**Implementation**:
- Modify default value in `command_tool.go`: `idleTimeout := 60 * time.Second`
- Update default value description in tool definition
- Update related documentation

## Risks / Trade-offs

**Risk 1**: User interaction may block command execution
- **Mitigation**: Use asynchronous UI interaction, command continues running in background

**Risk 2**: Users may extend timeout indefinitely, causing command to never end
- **Mitigation**: This is user's choice, if command is really stuck, user can choose to cancel

**Risk 3**: Timeout prompt may mix with command output
- **Mitigation**: Use `ui.ShowConfirm()` to ensure prompt displays on independent line

**Risk 4**: Default timeout increase may cause excessive waiting time in some scenarios
- **Mitigation**: Users can customize via `timeout` parameter, or choose to cancel

**Trade-offs**:
- **Flexibility vs Simplicity**: Choose flexibility, allow users to control timeout
- **Automation vs User Control**: Choose user control, provide better experience

## Migration Plan

**Deployment Steps:**
1. Modify default timeout time in `command_tool.go`
2. Add user prompt logic after timeout
3. Update timeout description in tool definition
4. Add test cases
5. Run existing tests to ensure no regression

**Rollback Strategy:**
- If issues found, can quickly rollback to previous version
- Default timeout time change is backward compatible

**Testing Strategy:**
- Unit tests: Test timeout prompt logic
- Integration tests: Test actual command execution scenarios
- Manual tests: Test user interaction flow

## Open Questions

1. **Should we display command information in prompt?** For example display command name being executed
   - **Decision**: Not needed for now, prompt message is clear enough

2. **Should we display elapsed waiting time in prompt?** For example "Command has been idle for 60s. Continue waiting?"
   - **Decision**: Can add, provides more context information

3. **Should we support automatic extension in non-interactive mode?** For example via environment variable configuration
   - **Decision**: Not supported for now, keep it simple

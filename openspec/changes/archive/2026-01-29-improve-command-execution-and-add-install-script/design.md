## Context

The current system uses simple loading spinner and full output display when executing `execute_command` tool, with the following issues:
1. Command output floods the screen, affecting users' ability to view context
2. Full output is returned to LLM, wasting tokens
3. Lack of installation script requires users to manually download and configure PATH

Referencing excellent designs from Cursor, zsh, brew, etc., we need to improve command execution display and add installation scripts.

**Current Architecture**:
- `internal/sql/tool_handler.go`: Handles tool call loop, displays tool call information and uses `ui.ShowLoading()` to show waiting state
- `internal/tool/builtin/command_tool.go`: Executes commands and returns full output
- `internal/ui/`: Provides color and formatting functions, already has `HintText()` for gray text
- Uses `github.com/charmbracelet/lipgloss` for style control

## Goals / Non-Goals

**Goals:**
- Implement scrolling window display (2-3 lines) during command execution to avoid screen flooding
- Use gray font to display tool call information
- Intelligent output truncation: truncate last 10-20 lines on success, last 50-100 lines on failure
- Provide cross-platform installation scripts (Unix/Linux/macOS and Windows)
- Installation scripts automatically detect latest version and system architecture
- Use CDN acceleration for downloads, support mainland China access

**Non-Goals:**
- Do not implement full terminal emulator (like tmux or screen)
- Do not implement complete command output history functionality
- Do not implement GUI interface for installation scripts
- Do not implement auto-update functionality (only install latest version)

## Decisions

### Decision 1: Scrolling Window Display Implementation

**Choice**: Use ANSI escape sequences for in-place updates instead of third-party libraries

**Rationale**:
- Currently have `lipgloss` but mainly for static styling
- Scrolling window needs dynamic cursor position updates and line clearing, ANSI escape sequences are more direct
- Avoid introducing new dependencies

**Alternatives Considered**:
- Use `github.com/charmbracelet/bubbletea`: Too complex, introduces unnecessary dependencies
- Use `github.com/gdamore/tcell`: Requires full TUI framework, over-engineered

**Implementation Details**:
- Use `\r` to return to line start, `\033[K` to clear to end of line
- Use `\033[A` to move cursor up, `\033[B` to move down
- Save current cursor position, update within 2-3 line range

### Decision 2: Output Truncation Strategy

**Choice**: Truncate output in `command_tool.go`'s `Execute()` method while retaining full output for user display

**Rationale**:
- Truncation logic is closely related to command execution, more reasonable to place in command tool
- Need to distinguish success/failure scenarios, command tool already has exit code information
- Retain full output for scrolling window display, only truncate portion returned to LLM

**Alternatives Considered**:
- Truncate in `tool_handler.go`: Need to parse JSON, increases complexity
- Truncate before LLM call: Cannot utilize exit code information

**Implementation Details**:
- `CommandResult` structure remains unchanged, add `TruncatedStdout` and `TruncatedStderr` fields
- On success (exit code 0): Truncate last 20 lines
- On failure (exit code != 0): Truncate last 100 lines
- If output is less than truncation threshold, return full output

### Decision 3: Real-time Output Streaming Display

**Choice**: For long-running commands, use goroutine to read output and update scrolling window in real-time

**Rationale**:
- Provides better user experience, users can see command execution progress
- Avoid displaying all output only after command completion

**Implementation Details**:
- Use `cmd.StdoutPipe()` and `cmd.StderrPipe()` to get output streams
- Start goroutine to read output and update scrolling window
- Main goroutine waits for command completion

### Decision 4: Installation Script Version Detection

**Choice**: Use GitHub Releases API (`https://api.github.com/repos/sunzhaoyang/aiq/releases/latest`) to get latest version

**Rationale**:
- GitHub API is stable and reliable, can get public repo releases without authentication
- Returns JSON format, easy to parse
- Supports fallback to direct download

**Alternatives Considered**:
- Use GitHub Tags API: Requires additional parsing, releases API is more direct
- Hardcode version number: Does not meet "automatically install latest version" requirement

**Implementation Details**:
- Use `curl` (Unix) or `powershell` (Windows) to call API
- Parse JSON to get `tag_name` field (e.g., `v0.0.1`)
- If API call fails, fallback to hardcoded latest known version

### Decision 5: CDN Acceleration Solution

**Choice**: Use jsdelivr CDN, format: `https://cdn.jsdelivr.net/gh/sunzhaoyang/aiq@<tag>/releases/download/<tag>/<binary>`

**Rationale**:
- jsdelivr has good access speed in mainland China
- Supports CDN acceleration for GitHub Releases
- No additional configuration needed, directly use GitHub repo path

**Alternatives Considered**:
- GitHub Releases direct link: May be slow in mainland China
- Self-hosted CDN: Increases operational costs

**Implementation Details**:
- Prioritize jsdelivr CDN
- If download fails (timeout or 404), fallback to GitHub Releases direct link
- Use `curl`'s `--fail` and `--max-time` options to detect failures

### Decision 6: PATH Configuration Method

**Unix/Linux/macOS**:
- Detect shell type (via `$SHELL` environment variable)
- Prioritize updating `~/.zshrc` (zsh) or `~/.bashrc` (bash)
- If neither exists, update `~/.profile`
- Check if PATH already contains installation directory to avoid duplicate addition

**Windows**:
- Use `setx PATH "%PATH%;<install_dir>"` to update user environment variables
- If `setx` fails (insufficient permissions), prompt user to run as administrator
- Note: `setx` requires new terminal to take effect, script should prompt user

### Decision 7: Installation Directory Selection

**Unix/Linux/macOS**:
- Default install to `~/.local/bin/aiq` (follows XDG specification)
- If `~/.local/bin` doesn't exist, create it
- Add `~/.local/bin` to PATH

**Windows**:
- Default install to `%LOCALAPPDATA%\aiq\aiq.exe`
- Add `%LOCALAPPDATA%\aiq` to PATH

## Risks / Trade-offs

**[Risk] Scrolling window implementation may be incompatible with some terminals**
- **Mitigation**: Detect terminal capabilities, if ANSI escape sequences are not supported, fallback to simple output mode

**[Risk] Real-time output streaming display may affect command execution performance**
- **Mitigation**: Use buffered reading, avoid frequent updates, set minimum update interval (e.g., 100ms)

**[Risk] Output truncation may lose important information**
- **Mitigation**: 
  - Truncate more lines on failure (100 lines vs 20 lines)
  - Retain full output for user viewing, only truncate portion returned to LLM
  - If output is less than truncation threshold, return full output

**[Risk] Installation script PATH update may fail (permissions, shell configuration, etc.)**
- **Mitigation**: 
  - Provide clear error messages
  - Provide instructions for manual PATH configuration
  - Verify installation success (try executing `aiq --version`)

**[Risk] CDN may be unavailable or return old version**
- **Mitigation**: 
  - Implement fallback to GitHub Releases direct link
  - Verify downloaded binary version (if binary supports `--version`)

**[Risk] Windows `setx` requires new terminal to take effect**
- **Mitigation**: 
  - Clearly prompt user to open new terminal after installation completes
  - Provide verification command for user testing

**[Trade-off] Real-time output vs Performance**
- Choice: Prioritize user experience, accept slight performance overhead
- Implement buffering and throttling mechanisms to reduce overhead

**[Trade-off] Full output retention vs Memory usage**
- Choice: Retain full output for display, but limit maximum output size (e.g., 10MB)
- When exceeding limit, only retain last N lines

## Migration Plan

1. **Phase 1: Implement Command Execution Display Optimization**
   - Modify `internal/ui/` to add scrolling window display function
   - Modify `internal/tool/builtin/command_tool.go` to implement output truncation
   - Modify `internal/sql/tool_handler.go` to use new display method
   - Test various command scenarios (success, failure, long-running)

2. **Phase 2: Implement Installation Scripts**
   - Create `scripts/install.sh` (Unix/Linux/macOS)
   - Create `scripts/install.bat` (Windows)
   - Test different platforms and shell environments
   - Update README to add installation instructions

3. **Phase 3: Documentation and Verification**
   - Update documentation to explain new command execution display method
   - Provide usage instructions for installation scripts
   - Verify CDN acceleration effectiveness

**Rollback Strategy**:
- If scrolling window implementation has issues, can quickly rollback to current simple display method
- Installation scripts are new features, don't affect existing functionality, no rollback needed

## Open Questions

1. **Should output truncation line count threshold be configurable?**
   - Current design: Hardcoded (20 lines on success, 100 lines on failure)
   - Consider: Support customization via environment variables or config files
   - Decision: Implement hardcoded version first, consider making configurable based on feedback

2. **Should we support output to file functionality?**
   - Current design: Only display in terminal
   - Consider: Some scenarios may need to save full output
   - Decision: Not in this implementation, consider based on future requirements

3. **Should installation scripts support specifying version for installation?**
   - Current design: Only support installing latest version
   - Consider: Some scenarios may need to install specific version
   - Decision: Not in this implementation, keep it simple, can extend later

4. **Should Windows installation script support PowerShell?**
   - Current design: Use batch script (.bat)
   - Consider: PowerShell is more powerful but requires PowerShell 5.0+
   - Decision: Implement .bat first, can consider providing PowerShell version later

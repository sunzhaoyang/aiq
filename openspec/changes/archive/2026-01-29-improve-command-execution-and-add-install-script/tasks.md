## 1. Command Execution Display - UI Components

- [x] 1.1 Add `ShowRollingOutput()` function in `internal/ui/` package for displaying rolling window output (2-3 lines)
- [x] 1.2 Add ANSI escape sequence utilities for cursor control (move up/down, clear line, save/restore position)
- [x] 1.3 Add `DimText()` or update `HintText()` to use dimmed gray color for tool call display (using existing `HintText()`)
- [x] 1.4 Add output buffer management utilities for maintaining rolling window state

## 2. Command Execution Display - Output Truncation

- [x] 2.1 Modify `CommandResult` struct in `internal/tool/builtin/command_tool.go` to add `TruncatedStdout` and `TruncatedStderr` fields
- [x] 2.2 Implement truncation logic in `CommandTool.Execute()`: success (20 lines), failure (100 lines)
- [x] 2.3 Preserve full output in original `Stdout`/`Stderr` fields for user display
- [x] 2.4 Add maximum output size limit (10MB) to prevent memory issues

## 3. Command Execution Display - Real-time Streaming

- [x] 3.1 Modify `CommandTool.Execute()` to use `cmd.StdoutPipe()` and `cmd.StderrPipe()` for streaming output
- [ ] 3.2 Implement goroutine to read output streams and update rolling window in real-time (basic implementation done, full real-time display requires callback mechanism)
- [x] 3.3 Add buffering and throttling (100ms minimum update interval) to reduce update frequency
- [x] 3.4 Handle command completion and cleanup goroutines properly

## 4. Command Execution Display - Integration

- [x] 4.1 Modify `tool_handler.go` to use gray/dimmed text for tool call display (`ui.HintText()`)
- [x] 4.2 Replace `ui.ShowLoading()` with rolling window display for `execute_command` tool
- [x] 4.3 Update `ExecuteTool()` to use truncated output when sending results to LLM
- [x] 4.4 Ensure full output is still available for user viewing in rolling window
- [x] 4.5 Add fallback to simple display mode if ANSI escape sequences are not supported

## 5. Installation Script - Unix/Linux/macOS (install.sh)

- [x] 5.1 Create `scripts/install.sh` with shebang and error handling (`set -e`)
- [x] 5.2 Implement version detection using GitHub Releases API (`curl` to `https://api.github.com/repos/sunzhaoyang/aiq/releases/latest`)
- [x] 5.3 Implement architecture detection (`uname -m` and `uname -s` for darwin-amd64, darwin-arm64, linux-amd64, linux-arm64)
- [x] 5.4 Implement binary download with CDN fallback (jsdelivr → GitHub Releases)
- [x] 5.5 Create `~/.local/bin` directory if it doesn't exist
- [x] 5.6 Implement PATH detection and update logic (detect shell, update `.zshrc`, `.bashrc`, or `.profile`)
- [x] 5.7 Add check to avoid duplicate PATH entries
- [x] 5.8 Implement installation verification (`aiq --version` check)
- [x] 5.9 Add error handling and user-friendly error messages
- [ ] 5.10 Test on macOS (darwin-amd64 and darwin-arm64)
- [ ] 5.11 Test on Linux (linux-amd64 and linux-arm64)
- [ ] 5.12 Test with different shells (bash, zsh)

## 6. Installation Script - Windows (install.bat)

- [x] 6.1 Create `scripts/install.bat` with error handling (`@echo off`, `setlocal`)
- [x] 6.2 Implement version detection using GitHub Releases API (`powershell` or `curl` if available)
- [x] 6.3 Implement architecture detection (assume windows-amd64, check `PROCESSOR_ARCHITECTURE`)
- [x] 6.4 Implement binary download with CDN fallback (jsdelivr → GitHub Releases)
- [x] 6.5 Create `%LOCALAPPDATA%\aiq` directory if it doesn't exist
- [x] 6.6 Implement PATH update using `setx PATH` command
- [x] 6.7 Add check to avoid duplicate PATH entries
- [x] 6.8 Handle permission errors (prompt user to run as administrator)
- [x] 6.9 Implement installation verification (`aiq.exe --version` check)
- [x] 6.10 Add error handling and user-friendly error messages
- [x] 6.11 Add note about needing new terminal for PATH to take effect
- [ ] 6.12 Test on Windows 10/11

## 7. Documentation and Testing

- [x] 7.1 Update README.md with installation instructions (curl command for install.sh)
- [x] 7.2 Add installation script usage examples for both platforms
- [ ] 7.3 Document new command execution display behavior
- [ ] 7.4 Test command execution display with various command types (short output, long output, errors)
- [ ] 7.5 Test output truncation with different exit codes
- [ ] 7.6 Test rolling window display with commands that produce continuous output
- [ ] 7.7 Verify CDN download speed and fallback behavior
- [ ] 7.8 Test installation scripts on clean systems (no existing aiq installation)

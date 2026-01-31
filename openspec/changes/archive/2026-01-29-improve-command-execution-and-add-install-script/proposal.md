## Why

The current command execution flow has poor user experience: command output floods the screen during execution, affecting users' ability to view context; returning all execution results to LLM wastes tokens; lack of convenient one-click installation script requires users to manually download and configure PATH, reducing product usability and professionalism. Referencing excellent designs from Cursor, zsh, brew, etc., we need to improve command execution display and add installation scripts.

## What Changes

- **Command Execution Display Optimization**:
  - Use gray font to display tool call information, scroll to show latest output within 2-3 lines at current position (similar to Cursor's design)
  - Truncate small amount of output when execution succeeds and return to LLM (e.g., last 10-20 lines) to reduce token consumption
  - Truncate more output when execution fails and return to LLM (e.g., last 50-100 lines) to help LLM determine error cause
  - Avoid screen flooding, keep interface clean
- **One-Click Installation Script**:
  - **Unix/Linux/macOS** (`install.sh`):
    - Automatically detect latest version (get latest tag via GitHub Releases API)
    - Automatically detect system architecture (darwin-amd64, darwin-arm64, linux-amd64, linux-arm64)
    - Automatically download corresponding platform binary package
    - Automatically add `aiq` to `$PATH` (support bash/zsh, detect and update `.bashrc`, `.zshrc` or `.profile`)
    - Use CDN acceleration accessible in mainland China (jsdelivr) to download Release packages
    - Provide installation verification and error handling
  - **Windows** (`install.bat`):
    - Automatically detect latest version (get latest tag via GitHub Releases API)
    - Automatically detect system architecture (windows-amd64)
    - Automatically download corresponding platform binary package (`.exe`)
    - Automatically add `aiq.exe` to `%PATH%` (update user environment variables)
    - Use CDN acceleration accessible in mainland China (jsdelivr) to download Release packages
    - Provide installation verification and error handling

**BREAKING**: None

## Capabilities

### New Capabilities

- `command-execution-display`: Real-time display optimization during command execution, including gray font, scrolling output, output truncation strategy
- `installation-script`: One-click installation script, supports Unix/Linux/macOS (`install.sh`) and Windows (`install.bat`), automatically detects latest version, system architecture, downloads binary, configures PATH, CDN acceleration

### Modified Capabilities

- `cli-application`: May need to add installation instructions or installation script references (if need to prompt users how to install in CLI)

## Impact

**Affected Code Modules:**

- `internal/sql/tool_handler.go`: Need to modify display logic of `execute_command` tool to implement scrolling output and output truncation
- `internal/ui/`: May need to add new UI components for scrolling command output display (gray font, position control)
- `internal/tool/builtin/command_tool.go`: May need to modify result format returned to LLM to implement output truncation logic
- New `scripts/install.sh`: Unix/Linux/macOS installation script (supports version detection)
- New `scripts/install.bat`: Windows installation script (supports version detection)

**Dependencies:**

- May need terminal control library (e.g., `github.com/charmbracelet/lipgloss` already in use) for formatting output
- Need to support ANSI escape sequences for scrolling display and color control

**User Experience Improvements:**

- Interface is cleaner during command execution, no screen flooding
- Users can see command execution progress in real-time
- Installation process is simpler, one-click completion, supports Unix/Linux/macOS and Windows
- Automatically install latest version, no need to manually specify version number
- Faster download speed for users in mainland China (CDN acceleration)


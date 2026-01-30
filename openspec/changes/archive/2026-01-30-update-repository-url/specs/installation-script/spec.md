## MODIFIED Requirements

### Requirement: Unix/Linux/macOS installation script
The system SHALL provide a shell script (`scripts/install.sh`) that automatically installs the latest version of `aiq` binary for Unix/Linux/macOS systems.

#### Scenario: Detect latest version
- **WHEN** user runs `install.sh`
- **THEN** system queries GitHub Releases API at `https://api.github.com/repos/sunetic/aiq/releases/latest` to get the latest release tag (e.g., `v0.0.1`)

#### Scenario: Detect system architecture
- **WHEN** installation script runs
- **THEN** system automatically detects architecture: darwin-amd64, darwin-arm64, linux-amd64, or linux-arm64

#### Scenario: Download binary from GitHub
- **WHEN** system has determined version and architecture
- **THEN** system downloads binary from GitHub Releases at `https://github.com/sunetic/aiq/releases/download` (no timeout, user can Ctrl+C if slow)

#### Scenario: Install to aiq home directory
- **WHEN** binary is downloaded
- **THEN** system installs to `~/.aiq/bin/aiq` and makes it executable

#### Scenario: Print PATH instruction
- **WHEN** installation completes
- **THEN** system prints the PATH export command for user to add manually
- **AND** system does NOT automatically modify shell config files (.zshrc, .bashrc, etc.)

#### Scenario: Show PATH command for zsh
- **WHEN** user's default shell is zsh
- **THEN** system prints: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.zshrc`

#### Scenario: Show PATH command for bash
- **WHEN** user's default shell is bash
- **THEN** system prints: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.bashrc`

#### Scenario: Verify installation
- **WHEN** installation completes
- **THEN** system verifies that binary exists and is executable, displays success message

#### Scenario: Handle installation errors
- **WHEN** any step of installation fails (download, permission)
- **THEN** system displays clear error message and exits with non-zero exit code

### Requirement: Windows installation script
The system SHALL provide a batch script (`scripts/install.bat`) that automatically installs the latest version of `aiq.exe` for Windows systems.

#### Scenario: Detect latest version (Windows)
- **WHEN** user runs `install.bat`
- **THEN** system queries GitHub Releases API at `https://api.github.com/repos/sunetic/aiq/releases/latest` using PowerShell to get the latest release tag

#### Scenario: Detect Windows architecture
- **WHEN** installation script runs on Windows
- **THEN** system assumes architecture as windows-amd64

#### Scenario: Download Windows binary
- **WHEN** system has determined version and architecture
- **THEN** system downloads `aiq-windows-amd64.exe` from GitHub Releases at `https://github.com/sunetic/aiq/releases/download` using PowerShell

#### Scenario: Install to aiq home directory (Windows)
- **WHEN** binary is downloaded
- **THEN** system installs to `%USERPROFILE%\.aiq\bin\aiq.exe`

#### Scenario: Print PATH instruction (Windows)
- **WHEN** installation completes on Windows
- **THEN** system prints the setx command for user to add PATH manually
- **AND** system does NOT automatically modify PATH environment variable

#### Scenario: Verify Windows installation
- **WHEN** installation completes on Windows
- **THEN** system verifies that `aiq.exe` exists and displays success message

#### Scenario: Handle Windows installation errors
- **WHEN** any step of Windows installation fails
- **THEN** system displays clear error message and pauses for user to read

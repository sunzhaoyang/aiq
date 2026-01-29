## ADDED Requirements

### Requirement: Unix/Linux/macOS installation script
The system SHALL provide a shell script (`install.sh`) that automatically installs the latest version of `aiq` binary for Unix/Linux/macOS systems.

#### Scenario: Detect latest version
- **WHEN** user runs `install.sh`
- **THEN** system queries GitHub Releases API to get the latest release tag (e.g., `v0.0.1`)

#### Scenario: Detect system architecture
- **WHEN** installation script runs
- **THEN** system automatically detects architecture: darwin-amd64, darwin-arm64, linux-amd64, or linux-arm64

#### Scenario: Download binary from CDN
- **WHEN** system has determined version and architecture
- **THEN** system downloads binary from jsdelivr CDN (e.g., `https://cdn.jsdelivr.net/gh/sunzhaoyang/aiq@latest/releases/download/v0.0.1/aiq-darwin-amd64`) for faster access in mainland China

#### Scenario: Fallback to GitHub direct download
- **WHEN** CDN download fails
- **THEN** system falls back to direct GitHub Releases download URL

#### Scenario: Add binary to PATH (bash)
- **WHEN** user's shell is bash and `~/.bashrc` exists
- **THEN** system adds `export PATH="$PATH:/path/to/aiq/bin"` to `~/.bashrc`

#### Scenario: Add binary to PATH (zsh)
- **WHEN** user's shell is zsh and `~/.zshrc` exists
- **THEN** system adds `export PATH="$PATH:/path/to/aiq/bin"` to `~/.zshrc`

#### Scenario: Add binary to PATH (fallback)
- **WHEN** neither `~/.bashrc` nor `~/.zshrc` exists
- **THEN** system adds PATH export to `~/.profile`

#### Scenario: Verify installation
- **WHEN** installation completes
- **THEN** system verifies that `aiq` command is available in PATH and displays success message

#### Scenario: Handle installation errors
- **WHEN** any step of installation fails (download, permission, PATH update)
- **THEN** system displays clear error message and exits with non-zero exit code

#### Scenario: Handle permission errors
- **WHEN** script lacks permissions to write to installation directory or shell config files
- **THEN** system displays error message suggesting to run with appropriate permissions

### Requirement: Windows installation script
The system SHALL provide a batch script (`install.bat`) that automatically installs the latest version of `aiq.exe` for Windows systems.

#### Scenario: Detect latest version (Windows)
- **WHEN** user runs `install.bat`
- **THEN** system queries GitHub Releases API to get the latest release tag (e.g., `v0.0.1`)

#### Scenario: Detect Windows architecture
- **WHEN** installation script runs on Windows
- **THEN** system automatically detects architecture as windows-amd64

#### Scenario: Download Windows binary from CDN
- **WHEN** system has determined version and architecture
- **THEN** system downloads `aiq-windows-amd64.exe` from jsdelivr CDN for faster access in mainland China

#### Scenario: Fallback to GitHub direct download (Windows)
- **WHEN** CDN download fails on Windows
- **THEN** system falls back to direct GitHub Releases download URL

#### Scenario: Add binary to PATH (Windows)
- **WHEN** installation completes on Windows
- **THEN** system adds installation directory to user's PATH environment variable using `setx` command or registry update

#### Scenario: Verify Windows installation
- **WHEN** installation completes on Windows
- **THEN** system verifies that `aiq.exe` command is available in PATH and displays success message

#### Scenario: Handle Windows installation errors
- **WHEN** any step of Windows installation fails (download, permission, PATH update)
- **THEN** system displays clear error message and exits with non-zero exit code

#### Scenario: Handle Windows permission errors
- **WHEN** script lacks permissions to update PATH environment variable
- **THEN** system displays error message suggesting to run as administrator

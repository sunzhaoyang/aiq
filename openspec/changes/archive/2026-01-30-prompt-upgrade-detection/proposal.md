## Why

To facilitate user modification of prompt files, the system outputs prompt files to the `~/.aiq/prompts` directory. However, during upgrades, the system is uncertain whether to use user-modified prompt files or overwrite with new versions. This leads to potentially using outdated prompts after upgrades, or accidentally overwriting user customizations.

We need to add a detection mechanism that checks at startup whether prompt file versions match built-in versions, and if inconsistent, ask the user whether to overwrite, and record the user's choice to avoid repeated prompts.

## What Changes

- **Add Application Version Management**:
  - Dynamically get application version at runtime (prioritize build info, if not available use git describe, finally use default value)
  - Get commit id (via git rev-parse or build-time injection)
  - Used to record user choices, avoid repeated prompts for same version
  - Support `-v` / `--version` command-line parameters, print version number and commit id then exit
- **Prompt Content Detection Mechanism**:
  - Calculate content hash (SHA256) of built-in prompt strings in code
  - Calculate content hash (SHA256) of prompt files in `~/.aiq/prompts/` directory
  - Compare built-in content hash with user file content hash at startup to detect if user has modified files
- **User Interaction Prompt**: When user-modified prompt files are detected, ask user whether to overwrite existing prompt files
- **Version Choice Recording**: Record user's choice for each application version (overwrite/keep), stored in `~/.aiq/config/prompt-version-choices.yaml`, ask only once per application version
- **Prompt Information**: If user chooses to keep, display prompt information telling how to manually delete files to trigger rebuild

## Capabilities

### New Capabilities
- `prompt-version-detection`: Detect consistency between prompt file versions and built-in versions, prompt user when inconsistent
- `application-version-management`: Manage application version numbers for version detection and upgrade prompts

### Modified Capabilities
- `configuration-management`: Need to add storage mechanism for version choice records (may be stored in new file under `~/.aiq/config/` directory)

## Impact

- **Code Changes**:
  - `internal/prompt/loader.go`: 
    - Add content hash calculation function (SHA256)
    - Add content detection logic, compare built-in content hash with user file content hash
    - Modify `initializeDefaults()` and `NewLoader()` methods
  - `internal/config/`: Add storage and reading functionality for version choice records
  - `internal/cli/root.go`: Call version detection logic at startup
  - `cmd/aiq/main.go`: Add `-v` / `--version` parameter handling, print version number and commit id then exit
  - New file: `internal/version/version.go` provides runtime functions to get version information:
    - `GetVersion()`: Prioritize getting version number from build-time injection (via `-ldflags`), if not available try using `git describe`, finally return default value "dev"
    - `GetCommitID()`: Prioritize getting commit id from build-time injection, if not available try using `git rev-parse HEAD`, finally return default value "unknown"
    - `GetVersionInfo()`: Return formatted version information string (e.g., "aiq v1.0.0 (commit: abc1234)")

- **Configuration Changes**:
  - New version choice record file (e.g., `~/.aiq/config/prompt-version-choices.yaml`), records user's choice for each version

- **Build Changes**:
  - CI/CD builds inject version number and commit id via `-ldflags`:
    - `-X github.com/aiq/aiq/internal/version.Version=${{ github.ref_name }}` (get from git tag)
    - `-X github.com/aiq/aiq/internal/version.CommitID=${{ github.sha }}` (get commit SHA from GitHub Actions)
  - Local development builds try to get from git describe and git rev-parse if version number not injected
  - If none available, use default values "dev" and "unknown"

- **User Experience**:
  - Do not print version number at startup (avoid interfering with normal use)
  - Support `aiq -v` or `aiq --version` to view version number and commit id
  - May display version detection prompt at startup (only when version inconsistent and user hasn't chosen before)
  - First startup after upgrade will ask whether to overwrite prompt files

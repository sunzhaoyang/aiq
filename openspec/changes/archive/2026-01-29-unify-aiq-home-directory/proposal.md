## Why

The current configuration directory `~/.aiqconfig` naming is not concise enough, and the installation script places binaries in `~/.local/bin` (non-standard macOS path). Referencing designs from Rust (`~/.cargo`), Go (`~/go`) and other tools, unify to use `~/.aiq` as AIQ's home directory, containing configuration and binary files, more concise and consistent.

## What Changes

- **Directory Rename**: `~/.aiqconfig` → `~/.aiq`, subdirectory structure remains unchanged
- **New bin Subdirectory**: `~/.aiq/bin` for storing aiq binary files
- **Installation Script Improvements**:
  - Installation location changed from `~/.local/bin` to `~/.aiq/bin`
  - No longer automatically modify user's shell configuration files
  - After installation completes, print PATH command for user to add to `.zshrc`/`.bashrc` themselves

## Capabilities

### New Capabilities

None

### Modified Capabilities

- `user-config-directory-organization`: Directory changed from `~/.aiqconfig` to `~/.aiq`, added `bin` subdirectory
- `installation-script`: Installation location changed to `~/.aiq/bin`, no longer automatically modify shell configuration, changed to print PATH command

## Impact

**Affected Code:**

- `internal/config/directory.go`: Modify `ConfigDir` constant and related paths
- `scripts/install.sh`: Modify installation directory and PATH handling logic
- `scripts/install.bat`: Modify Windows installation directory and PATH handling logic
- `README.md`: Update path references in documentation

**User Impact:**

- Users need to manually add `~/.aiq/bin` to PATH

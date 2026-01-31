## Context

Currently AIQ uses `~/.aiqconfig` as configuration directory, and installation script places binaries in `~/.local/bin`. This design has the following issues:

1. **Inconsistent Directory Naming**: Configuration directory uses `.aiqconfig` suffix, not concise enough
2. **Unified Installation Location**: Binary files separated from configuration, doesn't conform to "one tool one directory" design principle
3. **PATH Management Intrusiveness**: Installation script automatically modifies user's shell configuration files, may cause user dissatisfaction

Referencing designs from Rust (`~/.cargo`), Go (`~/go`), Node (`~/.nvm`) and other tools, unify to use `~/.aiq` as AIQ's home directory, containing all related files (configuration, sessions, skills, binaries).

## Goals / Non-Goals

**Goals:**
- Unify directory structure: `~/.aiq` contains all AIQ-related files
- Simplify naming: change from `.aiqconfig` to `.aiq`
- Unify installation location: binary files placed in `~/.aiq/bin`
- Reduce intrusiveness: no longer automatically modify shell configuration, changed to print commands for user to add themselves

**Non-Goals:**
- Do not modify existing subdirectory structure (config/, sessions/, skills/, etc. remain unchanged)

## Decisions

### 1. Directory Rename: `.aiqconfig` → `.aiq`

**Decision**: Change configuration directory from `~/.aiqconfig` to `~/.aiq`

**Rationale**:
- More concise, conforms to common tool naming conventions (`.cargo`, `.nvm`, `.go`)
- Remove `config` suffix because directory contains not only configuration but also sessions, skills, etc.

**Alternatives**:
- Keep `.aiqconfig`: Doesn't conform to conciseness principle
- Use `~/.config/aiq`: XDG standard, but macOS users unfamiliar

### 2. New `bin/` Subdirectory

**Decision**: Add new `bin/` subdirectory under `~/.aiq/` to store binary files

**Rationale**:
- Unified management: All AIQ-related files under one directory
- Conforms to "one tool one directory" design principle
- Easier for users to understand and manage

**Alternatives**:
- Continue using `~/.local/bin`: Separated from configuration, doesn't conform to unified management principle
- Use `/usr/local/bin`: Requires sudo, increases installation complexity

### 3. Do Not Automatically Modify Shell Configuration

**Decision**: Installation script no longer automatically modifies `.zshrc`/`.bashrc`, changed to print PATH command

**Rationale**:
- Reduce intrusiveness: Don't modify user configuration files
- User control: Let users decide whether to add and how to add
- Reference best practices: Similar to `rustup`, `nvm` approaches

**Alternatives**:
- Continue automatic modification: May overwrite user custom configurations, cause dissatisfaction
- Provide option for user to choose: Increases script complexity

## Risks / Trade-offs

### [Risk] User Forgets to Add PATH

**Impact**: Users may not be able to directly use `aiq` command after installation

**Mitigation**:
- Clearly print PATH command after installation completes
- Check PATH during installation verification, prompt if not in PATH
- Clearly state in README that PATH needs to be added

### [Trade-off] Installation Location vs System Standard Location

**Trade-off**: Use `~/.aiq/bin` instead of `/usr/local/bin` or `~/.local/bin`

**Choice**: `~/.aiq/bin` - Unified management, no sudo needed, conforms to tool design principles

**Cost**: Users need to manually add PATH (but this is an acceptable trade-off)

## Migration Plan

### Phase 1: Code Modifications
1. Modify `internal/config/directory.go`:
   - `ConfigDir` constant: `.aiqconfig` → `.aiq`
   - Add new `BinSubdir = "bin"`
   - Add new `GetBinDir()` function
   - Update `EnsureDirectoryStructure()` to include `bin/`

2. Modify `scripts/install.sh`:
   - Installation directory: `~/.local/bin` → `~/.aiq/bin`
   - Remove automatic shell configuration modification logic
   - Add logic to print PATH command

3. Modify `scripts/install.bat`:
   - Installation directory: `%LOCALAPPDATA%\aiq` → `%USERPROFILE%\.aiq\bin`
   - Remove automatic PATH modification logic
   - Add logic to print setx command

### Phase 2: Documentation Updates
1. Update `README.md`:
   - Change all `~/.aiqconfig` references to `~/.aiq`
   - Update installation instructions to state PATH needs to be manually added

## Open Questions

None

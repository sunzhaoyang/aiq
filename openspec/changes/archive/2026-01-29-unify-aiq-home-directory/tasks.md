## 1. Modify Configuration Directory Path

- [x] 1.1 Modify `ConfigDir` constant in `internal/config/directory.go`: `.aiqconfig` → `.aiq`
- [x] 1.2 Add new `BinSubdir = "bin"` constant in `internal/config/directory.go`
- [x] 1.3 Add new `GetBinDir()` function in `internal/config/directory.go` to return `~/.aiq/bin` path
- [x] 1.4 Update `EnsureDirectoryStructure()` function to add `bin/` subdirectory creation logic
- [x] 1.5 Update all path references in comments (from `~/.aiqconfig` to `~/.aiq`)

## 2. Modify Unix/Linux/macOS Installation Script

- [x] 2.1 Modify `INSTALL_DIR` in `scripts/install.sh`: from `~/.local/bin` to `~/.aiq/bin`
- [x] 2.2 Remove automatic shell detection and modification of `.zshrc`/`.bashrc`/`.profile` logic in `scripts/install.sh`
- [x] 2.3 After installation completes in `scripts/install.sh`, print corresponding PATH command based on user shell:
  - zsh: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.zshrc`
  - bash: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.bashrc`
  - other: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.profile`
- [x] 2.4 Update installation directory display information in `scripts/install.sh`
- [x] 2.5 When verifying installation in `scripts/install.sh`, check if PATH contains `~/.aiq/bin`, if not, prompt user to add

## 3. Modify Windows Installation Script

- [x] 3.1 Modify `INSTALL_DIR` in `scripts/install.bat`: from `%LOCALAPPDATA%\aiq` to `%USERPROFILE%\.aiq\bin`
- [x] 3.2 Remove automatic `setx` PATH modification logic in `scripts/install.bat`
- [x] 3.3 After installation completes in `scripts/install.bat`, print setx command for user to manually execute:
  - `setx PATH "%PATH%;%USERPROFILE%\.aiq\bin"`
- [x] 3.4 Update installation directory display information in `scripts/install.bat`
- [x] 3.5 When verifying installation in `scripts/install.bat`, check if PATH contains installation directory, if not, prompt user to add

## 4. Update Documentation

- [x] 4.1 Update all `~/.aiqconfig` references to `~/.aiq` in `README.md`
- [x] 4.2 Update installation instructions in `README.md` to specify installation location as `~/.aiq/bin`
- [x] 4.3 Add PATH configuration instructions in `README.md` (users need to manually add)

## 5. Testing and Verification

- [x] 5.1 Test new installation: Run `install.sh`, verify binary installs to `~/.aiq/bin` (verified: installation script works correctly)
- [x] 5.2 Test PATH command printing: Verify installation script correctly prints PATH command (verified: zsh/bash/other shells all correctly print corresponding commands)
- [x] 5.3 Test directory creation: First run of program, verify `~/.aiq` and all subdirectories (including `bin/`) are correctly created (verified: bin/, config/, sessions/, skills/, tools/, prompts/ all created)
- [ ] 5.4 Test Windows installation script: Test `install.bat` on Windows, verify installation location and PATH command printing (requires Windows environment, code already implemented)
- [x] 5.5 Verify all path resolution functions work correctly (config, sessions, skills, tools, prompts, bin) (verified: GetBinDir() implemented and used in EnsureDirectoryStructure(), all path functions work correctly)

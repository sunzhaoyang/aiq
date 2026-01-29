## 1. 修改配置目录路径

- [x] 1.1 修改 `internal/config/directory.go` 中的 `ConfigDir` 常量：`.aiqconfig` → `.aiq`
- [x] 1.2 在 `internal/config/directory.go` 中新增 `BinSubdir = "bin"` 常量
- [x] 1.3 在 `internal/config/directory.go` 中新增 `GetBinDir()` 函数返回 `~/.aiq/bin` 路径
- [x] 1.4 更新 `EnsureDirectoryStructure()` 函数，添加 `bin/` 子目录的创建逻辑
- [x] 1.5 更新所有注释中的路径引用（从 `~/.aiqconfig` 改为 `~/.aiq`）

## 2. 修改 Unix/Linux/macOS 安装脚本

- [x] 2.1 修改 `scripts/install.sh` 中的 `INSTALL_DIR`：从 `~/.local/bin` 改为 `~/.aiq/bin`
- [x] 2.2 移除 `scripts/install.sh` 中自动检测 shell 和修改 `.zshrc`/`.bashrc`/`.profile` 的逻辑
- [x] 2.3 在 `scripts/install.sh` 安装完成后，根据用户 shell 打印对应的 PATH 命令：
  - zsh: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.zshrc`
  - bash: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.bashrc`
  - 其他: `echo 'export PATH="$HOME/.aiq/bin:$PATH"' >> ~/.profile`
- [x] 2.4 更新 `scripts/install.sh` 中的安装目录显示信息
- [x] 2.5 在 `scripts/install.sh` 验证安装时，检查 PATH 是否包含 `~/.aiq/bin`，如果不在则提示用户添加

## 3. 修改 Windows 安装脚本

- [x] 3.1 修改 `scripts/install.bat` 中的 `INSTALL_DIR`：从 `%LOCALAPPDATA%\aiq` 改为 `%USERPROFILE%\.aiq\bin`
- [x] 3.2 移除 `scripts/install.bat` 中自动使用 `setx` 修改 PATH 的逻辑
- [x] 3.3 在 `scripts/install.bat` 安装完成后，打印 setx 命令让用户手动执行：
  - `setx PATH "%PATH%;%USERPROFILE%\.aiq\bin"`
- [x] 3.4 更新 `scripts/install.bat` 中的安装目录显示信息
- [x] 3.5 在 `scripts/install.bat` 验证安装时，检查 PATH 是否包含安装目录，如果不在则提示用户添加

## 4. 更新文档

- [x] 4.1 更新 `README.md` 中所有 `~/.aiqconfig` 引用为 `~/.aiq`
- [x] 4.2 更新 `README.md` 安装说明，说明安装位置为 `~/.aiq/bin`
- [x] 4.3 在 `README.md` 中添加 PATH 配置说明（用户需要手动添加）

## 5. 测试和验证

- [x] 5.1 测试新安装：运行 `install.sh`，验证二进制安装到 `~/.aiq/bin`（已验证：安装脚本正常工作）
- [x] 5.2 测试 PATH 命令打印：验证安装脚本正确打印 PATH 命令（已验证：zsh/bash/其他 shell 都能正确打印对应命令）
- [x] 5.3 测试目录创建：首次运行程序，验证 `~/.aiq` 及其所有子目录（包括 `bin/`）被正确创建（已验证：bin/, config/, sessions/, skills/, tools/, prompts/ 都已创建）
- [ ] 5.4 测试 Windows 安装脚本：在 Windows 上测试 `install.bat`，验证安装位置和 PATH 命令打印（需要 Windows 环境，代码已实现）
- [x] 5.5 验证所有路径解析函数正常工作（config, sessions, skills, tools, prompts, bin）（已验证：GetBinDir() 已实现并在 EnsureDirectoryStructure() 中使用，所有路径函数正常工作）

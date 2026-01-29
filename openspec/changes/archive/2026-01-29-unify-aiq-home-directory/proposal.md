## Why

当前配置目录 `~/.aiqconfig` 命名不够简洁，且安装脚本将二进制放在 `~/.local/bin`（非 macOS 标准路径）。参考 Rust (`~/.cargo`)、Go (`~/go`) 等工具的设计，统一使用 `~/.aiq` 作为 AIQ 的 home 目录，包含配置和二进制文件，更加简洁一致。

## What Changes

- **目录重命名**：`~/.aiqconfig` → `~/.aiq`，子目录结构保持不变
- **新增 bin 子目录**：`~/.aiq/bin` 用于存放 aiq 二进制文件
- **安装脚本改进**：
  - 安装位置从 `~/.local/bin` 改为 `~/.aiq/bin`
  - 不再自动修改用户的 shell 配置文件
  - 安装完成后打印 PATH 命令，让用户自行添加到 `.zshrc`/`.bashrc`

## Capabilities

### New Capabilities

无

### Modified Capabilities

- `user-config-directory-organization`: 目录从 `~/.aiqconfig` 改为 `~/.aiq`，新增 `bin` 子目录
- `installation-script`: 安装位置改为 `~/.aiq/bin`，不再自动修改 shell 配置，改为打印 PATH 命令

## Impact

**受影响的代码：**

- `internal/config/directory.go`: 修改 `ConfigDir` 常量和相关路径
- `scripts/install.sh`: 修改安装目录和 PATH 处理逻辑
- `scripts/install.bat`: 修改 Windows 安装目录和 PATH 处理逻辑
- `README.md`: 更新文档中的路径引用

**用户影响：**

- 用户需要手动将 `~/.aiq/bin` 添加到 PATH

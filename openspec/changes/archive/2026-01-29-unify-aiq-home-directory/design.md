## Context

当前 AIQ 使用 `~/.aiqconfig` 作为配置目录，安装脚本将二进制放在 `~/.local/bin`。这种设计存在以下问题：

1. **目录命名不一致**：配置目录使用 `.aiqconfig` 后缀，不够简洁
2. **安装位置不统一**：二进制文件与配置分离，不符合"一个工具一个目录"的设计原则
3. **PATH 管理侵入性**：安装脚本自动修改用户的 shell 配置文件，可能引起用户反感

参考 Rust (`~/.cargo`)、Go (`~/go`)、Node (`~/.nvm`) 等工具的设计，统一使用 `~/.aiq` 作为 AIQ 的 home 目录，包含所有相关文件（配置、会话、技能、二进制）。

## Goals / Non-Goals

**Goals:**
- 统一目录结构：`~/.aiq` 包含所有 AIQ 相关文件
- 简化命名：从 `.aiqconfig` 改为 `.aiq`
- 统一安装位置：二进制文件放在 `~/.aiq/bin`
- 减少侵入性：不再自动修改 shell 配置，改为打印命令让用户自行添加

**Non-Goals:**
- 不修改现有的子目录结构（config/, sessions/, skills/ 等保持不变）

## Decisions

### 1. 目录重命名：`.aiqconfig` → `.aiq`

**决策**：将配置目录从 `~/.aiqconfig` 改为 `~/.aiq`

**理由**：
- 更简洁，符合常见工具命名习惯（`.cargo`, `.nvm`, `.go`）
- 去掉 `config` 后缀，因为目录不仅包含配置，还包含会话、技能等

**替代方案**：
- 保持 `.aiqconfig`：不符合简洁性原则
- 使用 `~/.config/aiq`：XDG 标准，但 macOS 用户不熟悉

### 2. 新增 `bin/` 子目录

**决策**：在 `~/.aiq/` 下新增 `bin/` 子目录存放二进制文件

**理由**：
- 统一管理：所有 AIQ 相关文件在一个目录下
- 符合"一个工具一个目录"的设计原则
- 便于用户理解和管理

**替代方案**：
- 继续使用 `~/.local/bin`：与配置分离，不符合统一管理原则
- 使用 `/usr/local/bin`：需要 sudo，增加安装复杂度

### 3. 不自动修改 shell 配置

**决策**：安装脚本不再自动修改 `.zshrc`/`.bashrc`，改为打印 PATH 命令

**理由**：
- 减少侵入性：不修改用户配置文件
- 用户控制：让用户决定是否添加和如何添加
- 参考最佳实践：类似 `rustup`、`nvm` 的做法

**替代方案**：
- 继续自动修改：可能覆盖用户自定义配置，引起反感
- 提供选项让用户选择：增加脚本复杂度

## Risks / Trade-offs

### [风险] 用户忘记添加 PATH

**影响**：安装后用户可能无法直接使用 `aiq` 命令

**缓解措施**：
- 安装完成后清晰打印 PATH 命令
- 验证安装时检查 PATH，如果不在 PATH 中给出提示
- 在 README 中明确说明需要添加 PATH

### [权衡] 安装位置 vs 系统标准位置

**权衡**：使用 `~/.aiq/bin` 而非 `/usr/local/bin` 或 `~/.local/bin`

**选择**：`~/.aiq/bin` - 统一管理，无需 sudo，符合工具设计原则

**代价**：用户需要手动添加 PATH（但这是可接受的权衡）

## Migration Plan

### Phase 1: 代码修改
1. 修改 `internal/config/directory.go`：
   - `ConfigDir` 常量：`.aiqconfig` → `.aiq`
   - 新增 `BinSubdir = "bin"`
   - 新增 `GetBinDir()` 函数
   - 更新 `EnsureDirectoryStructure()` 包含 `bin/`

2. 修改 `scripts/install.sh`：
   - 安装目录：`~/.local/bin` → `~/.aiq/bin`
   - 移除自动修改 shell 配置的逻辑
   - 添加打印 PATH 命令的逻辑

3. 修改 `scripts/install.bat`：
   - 安装目录：`%LOCALAPPDATA%\aiq` → `%USERPROFILE%\.aiq\bin`
   - 移除自动修改 PATH 的逻辑
   - 添加打印 setx 命令的逻辑

### Phase 2: 文档更新
1. 更新 `README.md`：
   - 所有 `~/.aiqconfig` 引用改为 `~/.aiq`
   - 更新安装说明，说明 PATH 需要手动添加

## Open Questions

无

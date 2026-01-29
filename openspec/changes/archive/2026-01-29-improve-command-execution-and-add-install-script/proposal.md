## Why

当前命令执行流程的用户体验不够友好：执行命令时输出会刷屏，影响用户查看上下文；执行结果全部返回给 LLM 导致 token 浪费；缺少便捷的一键安装脚本，用户需要手动下载和配置 PATH，降低了产品的易用性和专业度。参考 Cursor、zsh、brew 等优秀产品的设计，需要改进命令执行显示方式并增加安装脚本。

## What Changes

- **命令执行显示优化**：
  - 使用灰色字体显示 tool call 信息，在当前位置的 2-3 行范围内滚动展示最新输出（类似 Cursor 的设计）
  - 执行成功后截取少量输出返回给 LLM（例如最后 10-20 行），减少 token 消耗
  - 执行失败时多截取输出返回给 LLM（例如最后 50-100 行），帮助 LLM 判断错误原因
  - 避免刷屏，保持界面整洁
- **一键安装脚本**：
  - **Unix/Linux/macOS** (`install.sh`)：
    - 自动检测最新版本（通过 GitHub Releases API 获取 latest tag）
    - 自动检测系统架构（darwin-amd64, darwin-arm64, linux-amd64, linux-arm64）
    - 自动下载对应平台的 binary 包
    - 自动将 `aiq` 添加到 `$PATH`（支持 bash/zsh，检测并更新 `.bashrc`、`.zshrc` 或 `.profile`）
    - 使用中国大陆可访问的 CDN 加速（jsdelivr）下载 Release 包
    - 提供安装验证和错误处理
  - **Windows** (`install.bat`)：
    - 自动检测最新版本（通过 GitHub Releases API 获取 latest tag）
    - 自动检测系统架构（windows-amd64）
    - 自动下载对应平台的 binary 包（`.exe`）
    - 自动将 `aiq.exe` 添加到 `%PATH%`（更新用户环境变量）
    - 使用中国大陆可访问的 CDN 加速（jsdelivr）下载 Release 包
    - 提供安装验证和错误处理

**BREAKING**: 无

## Capabilities

### New Capabilities

- `command-execution-display`: 命令执行时的实时显示优化，包括灰色字体、滚动输出、输出截断策略
- `installation-script`: 一键安装脚本，支持 Unix/Linux/macOS (`install.sh`) 和 Windows (`install.bat`)，自动检测最新版本、系统架构、下载 binary、配置 PATH、CDN 加速

### Modified Capabilities

- `cli-application`: 可能需要添加安装说明或安装脚本的引用（如果需要在 CLI 中提示用户如何安装）

## Impact

**受影响的代码模块：**

- `internal/sql/tool_handler.go`: 需要修改 `execute_command` tool 的显示逻辑，实现滚动输出和输出截断
- `internal/ui/`: 可能需要添加新的 UI 组件用于滚动显示命令输出（灰色字体、位置控制）
- `internal/tool/builtin/command_tool.go`: 可能需要修改返回给 LLM 的结果格式，实现输出截断逻辑
- 新增 `scripts/install.sh`: Unix/Linux/macOS 安装脚本（支持版本检测）
- 新增 `scripts/install.bat`: Windows 安装脚本（支持版本检测）

**依赖项：**

- 可能需要终端控制库（如 `github.com/charmbracelet/lipgloss` 已在使用）用于格式化输出
- 需要支持 ANSI 转义序列用于滚动显示和颜色控制

**用户体验改进：**

- 命令执行时界面更整洁，不会刷屏
- 用户可以实时看到命令执行进度
- 安装过程更简单，一键完成，支持 Unix/Linux/macOS 和 Windows
- 自动安装最新版本，无需手动指定版本号
- 中国大陆用户下载速度更快（CDN 加速）


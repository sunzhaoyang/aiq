## Context

当前系统在执行 `execute_command` tool 时，使用简单的 loading spinner 和完整输出显示，存在以下问题：
1. 命令输出会刷屏，影响用户查看上下文
2. 完整输出全部返回给 LLM，浪费 token
3. 缺少安装脚本，用户需要手动下载和配置 PATH

参考 Cursor、zsh、brew 等优秀产品的设计，需要改进命令执行显示方式并增加安装脚本。

**当前架构**：
- `internal/sql/tool_handler.go`: 处理 tool call 循环，显示 tool call 信息和使用 `ui.ShowLoading()` 显示等待状态
- `internal/tool/builtin/command_tool.go`: 执行命令并返回完整输出
- `internal/ui/`: 提供颜色和格式化函数，已有 `HintText()` 用于灰色文本
- 使用 `github.com/charmbracelet/lipgloss` 进行样式控制

## Goals / Non-Goals

**Goals:**
- 实现命令执行时的滚动窗口显示（2-3 行），避免刷屏
- 使用灰色字体显示 tool call 信息
- 智能截断输出：成功时截取最后 10-20 行，失败时截取最后 50-100 行
- 提供跨平台安装脚本（Unix/Linux/macOS 和 Windows）
- 安装脚本自动检测最新版本和系统架构
- 使用 CDN 加速下载，支持中国大陆访问

**Non-Goals:**
- 不实现完整的终端模拟器（如 tmux 或 screen）
- 不实现命令输出的完整历史记录功能
- 不实现安装脚本的 GUI 界面
- 不实现自动更新功能（仅安装最新版本）

## Decisions

### Decision 1: 滚动窗口显示实现方式

**选择**: 使用 ANSI 转义序列实现原地更新，而不是使用第三方库

**理由**:
- 当前已有 `lipgloss`，但主要用于静态样式
- 滚动窗口需要动态更新光标位置和清除行，ANSI 转义序列更直接
- 避免引入新的依赖

**替代方案考虑**:
- 使用 `github.com/charmbracelet/bubbletea`: 功能过于复杂，引入不必要的依赖
- 使用 `github.com/gdamore/tcell`: 需要完整的 TUI 框架，过度设计

**实现细节**:
- 使用 `\r` 回到行首，`\033[K` 清除到行尾
- 使用 `\033[A` 向上移动光标，`\033[B` 向下移动
- 保存当前光标位置，在 2-3 行范围内更新

### Decision 2: 输出截断策略

**选择**: 在 `command_tool.go` 的 `Execute()` 方法中截断输出，同时保留完整输出用于用户显示

**理由**:
- 截断逻辑与命令执行紧密相关，放在 command tool 中更合理
- 需要区分成功/失败场景，command tool 已有 exit code 信息
- 保留完整输出用于滚动窗口显示，只截断返回给 LLM 的部分

**替代方案考虑**:
- 在 `tool_handler.go` 中截断: 需要解析 JSON，增加复杂度
- 在 LLM 调用前截断: 无法利用 exit code 信息

**实现细节**:
- `CommandResult` 结构保持不变，添加 `TruncatedStdout` 和 `TruncatedStderr` 字段
- 成功时（exit code 0）: 截取最后 20 行
- 失败时（exit code != 0）: 截取最后 100 行
- 如果输出少于截断阈值，返回完整输出

### Decision 3: 实时输出流式显示

**选择**: 对于长时间运行的命令，使用 goroutine 读取输出并实时更新滚动窗口

**理由**:
- 提供更好的用户体验，用户可以看到命令执行进度
- 避免命令执行完成后才显示所有输出

**实现细节**:
- 使用 `cmd.StdoutPipe()` 和 `cmd.StderrPipe()` 获取输出流
- 启动 goroutine 读取输出并更新滚动窗口
- 主 goroutine 等待命令完成

### Decision 4: 安装脚本版本检测

**选择**: 使用 GitHub Releases API (`https://api.github.com/repos/sunzhaoyang/aiq/releases/latest`) 获取最新版本

**理由**:
- GitHub API 稳定可靠，无需认证即可获取 public repo 的 releases
- 返回 JSON 格式，易于解析
- 支持 fallback 到直接下载

**替代方案考虑**:
- 使用 GitHub Tags API: 需要额外解析，releases API 更直接
- 硬编码版本号: 不符合"自动安装最新版本"的需求

**实现细节**:
- 使用 `curl` (Unix) 或 `powershell` (Windows) 调用 API
- 解析 JSON 获取 `tag_name` 字段（如 `v0.0.1`）
- 如果 API 调用失败，fallback 到硬编码的最新已知版本

### Decision 5: CDN 加速方案

**选择**: 使用 jsdelivr CDN，格式为 `https://cdn.jsdelivr.net/gh/sunzhaoyang/aiq@<tag>/releases/download/<tag>/<binary>`

**理由**:
- jsdelivr 在中国大陆有良好的访问速度
- 支持 GitHub Releases 的 CDN 加速
- 无需额外配置，直接使用 GitHub repo 路径

**替代方案考虑**:
- GitHub Releases 直链: 在中国大陆可能较慢
- 自建 CDN: 增加运维成本

**实现细节**:
- 优先尝试 jsdelivr CDN
- 如果下载失败（超时或 404），fallback 到 GitHub Releases 直链
- 使用 `curl` 的 `--fail` 和 `--max-time` 选项检测失败

### Decision 6: PATH 配置方式

**Unix/Linux/macOS**:
- 检测 shell 类型（通过 `$SHELL` 环境变量）
- 优先更新 `~/.zshrc` (zsh) 或 `~/.bashrc` (bash)
- 如果都不存在，更新 `~/.profile`
- 检查 PATH 中是否已包含安装目录，避免重复添加

**Windows**:
- 使用 `setx PATH "%PATH%;<install_dir>"` 更新用户环境变量
- 如果 `setx` 失败（权限不足），提示用户以管理员身份运行
- 注意：`setx` 需要新开终端才能生效，脚本中提示用户

### Decision 7: 安装目录选择

**Unix/Linux/macOS**:
- 默认安装到 `~/.local/bin/aiq`（符合 XDG 规范）
- 如果 `~/.local/bin` 不存在，创建它
- 将 `~/.local/bin` 添加到 PATH

**Windows**:
- 默认安装到 `%LOCALAPPDATA%\aiq\aiq.exe`
- 将 `%LOCALAPPDATA%\aiq` 添加到 PATH

## Risks / Trade-offs

**[Risk] 滚动窗口实现可能在某些终端不兼容**
- **Mitigation**: 检测终端能力，如果不支持 ANSI 转义序列，fallback 到简单输出模式

**[Risk] 实时输出流式显示可能影响命令执行性能**
- **Mitigation**: 使用缓冲读取，避免频繁更新，设置最小更新间隔（如 100ms）

**[Risk] 输出截断可能丢失重要信息**
- **Mitigation**: 
  - 失败时截取更多行（100 行 vs 20 行）
  - 保留完整输出用于用户查看，只截断返回给 LLM 的部分
  - 如果输出少于截断阈值，返回完整输出

**[Risk] 安装脚本的 PATH 更新可能失败（权限、shell 配置等）**
- **Mitigation**: 
  - 提供清晰的错误消息
  - 提供手动配置 PATH 的说明
  - 验证安装是否成功（尝试执行 `aiq --version`）

**[Risk] CDN 可能不可用或返回旧版本**
- **Mitigation**: 
  - 实现 fallback 到 GitHub Releases 直链
  - 验证下载的 binary 版本（如果 binary 支持 `--version`）

**[Risk] Windows 的 `setx` 需要新终端才能生效**
- **Mitigation**: 
  - 安装完成后明确提示用户需要新开终端
  - 提供验证命令让用户测试

**[Trade-off] 实时输出 vs 性能**
- 选择：优先用户体验，接受轻微性能开销
- 实现缓冲和节流机制减少开销

**[Trade-off] 完整输出保留 vs 内存占用**
- 选择：保留完整输出用于显示，但限制最大输出大小（如 10MB）
- 超过限制时只保留最后 N 行

## Migration Plan

1. **Phase 1: 实现命令执行显示优化**
   - 修改 `internal/ui/` 添加滚动窗口显示函数
   - 修改 `internal/tool/builtin/command_tool.go` 实现输出截断
   - 修改 `internal/sql/tool_handler.go` 使用新的显示方式
   - 测试各种命令场景（成功、失败、长时间运行）

2. **Phase 2: 实现安装脚本**
   - 创建 `scripts/install.sh`（Unix/Linux/macOS）
   - 创建 `scripts/install.bat`（Windows）
   - 测试不同平台和 shell 环境
   - 更新 README 添加安装说明

3. **Phase 3: 文档和验证**
   - 更新文档说明新的命令执行显示方式
   - 提供安装脚本的使用说明
   - 验证 CDN 加速效果

**Rollback Strategy**:
- 如果滚动窗口实现有问题，可以快速回退到当前的简单显示方式
- 安装脚本是新增功能，不影响现有功能，无需 rollback

## Open Questions

1. **输出截断的行数阈值是否需要可配置？**
   - 当前设计：硬编码（成功 20 行，失败 100 行）
   - 考虑：通过环境变量或配置文件支持自定义
   - 决定：先实现硬编码版本，后续根据反馈考虑可配置

2. **是否需要支持输出到文件的功能？**
   - 当前设计：只显示在终端
   - 考虑：某些场景可能需要保存完整输出
   - 决定：不在本次实现中，后续根据需求考虑

3. **安装脚本是否需要支持指定版本安装？**
   - 当前设计：只支持安装最新版本
   - 考虑：某些场景可能需要安装特定版本
   - 决定：不在本次实现中，保持简单，后续可扩展

4. **Windows 安装脚本是否需要支持 PowerShell？**
   - 当前设计：使用批处理脚本（.bat）
   - 考虑：PowerShell 功能更强大，但需要 PowerShell 5.0+
   - 决定：先实现 .bat，后续可考虑提供 PowerShell 版本

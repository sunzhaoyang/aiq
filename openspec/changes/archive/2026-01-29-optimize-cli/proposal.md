## Why

当前 CLI 的用户体验存在一些不够完善的地方，影响了使用流畅度。用户在进入 chat 模式后无法方便地返回主菜单，source 管理缺少编辑功能，chat 模式缺少明确的命令提示和帮助，命令行参数 `-D` 在重复使用时无法正确应用。这些问题降低了产品的易用性和专业性，需要参考优秀 CLI 产品的设计模式进行改进。

## What Changes

- **导航和退出机制改进**：为所有功能模块（包括 chat 模式）提供统一的进入/退出机制，支持从 chat 模式返回到主菜单
- **Source 编辑功能**：在 source 管理菜单中增加 `edit` 选项，支持修改已存在的 source 配置（名称、主机、端口、数据库、用户名、密码等）
- **Chat 模式命令增强**：
  - 增加 `/exit` 命令（与现有的 `exit`/`back` 文本命令并存）
  - 增加 `/help` 命令显示可用命令列表和使用说明
  - 支持 Tab 补全功能，自动补全命令和提供上下文相关的建议
- **命令行参数 `-D` 处理优化**：当通过命令行参数 `-D` 指定数据库时，在当次执行中使用指定的数据库，但不持久化到 source 配置中（source 的唯一性仍基于 host-port-username，但执行时优先使用命令行参数指定的 database）

**BREAKING**: 无

## Capabilities

### New Capabilities
- `cli-navigation`: 统一的 CLI 导航和退出机制，支持从任何子功能返回到主菜单，提供一致的导航体验

### Modified Capabilities
- `cli-application`: 增强主菜单和子菜单的导航逻辑，确保 chat 模式可以返回到主菜单
- `data-source-management`: 增加 source 编辑功能，支持完整的 CRUD 操作（Create, Read, Update, Delete）
- `sql-interactive-mode`: 增加 `/exit`、`/help` 命令和 Tab 补全支持，提升 chat 模式的交互体验
- `mysql-cli-args`: 优化 `-D` 参数处理逻辑，支持当次执行时使用命令行指定的 database，而不影响持久化的 source 配置

## Impact

**受影响的代码模块：**
- `internal/cli/root.go`: 需要调整 chat 模式的调用和返回逻辑
- `internal/cli/source.go`: 需要增加 `editSource()` 函数和相关菜单项
- `internal/sql/mode.go`: 需要增加命令解析逻辑（`/exit`, `/help`）和 Tab 补全支持
- `internal/cli/dbconnect.go` 或相关文件: 需要调整 source 创建和使用逻辑，支持命令行 `-D` 参数的当次覆盖
- `internal/source/manager.go`: 可能需要增加 `UpdateSource()` 函数支持编辑功能

**依赖项：**
- 可能需要增强 `readline` 库的使用，支持 Tab 补全功能
- 需要确保命令解析逻辑不与现有的自然语言查询冲突

**用户体验改进：**
- 用户可以从 chat 模式方便地返回到主菜单，无需强制退出程序
- 用户可以编辑已存在的 source，无需删除重建
- 用户可以通过 `/help` 了解可用命令，通过 Tab 补全提高输入效率
- 命令行参数 `-D` 的行为更加符合用户预期，每次执行都能使用指定的数据库

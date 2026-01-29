## Context

当前 CLI 存在以下用户体验问题：
1. Chat 模式退出后无法返回主菜单，只能完全退出程序
2. Source 管理缺少编辑功能，用户需要删除重建才能修改配置
3. Chat 模式缺少明确的命令提示（如 `/help`）和 Tab 补全，用户不知道可用命令
4. 命令行参数 `-D` 指定的 database 在第二次执行时被忽略，因为 source 的唯一性基于 host-port-username，不包含 database

现有代码结构：
- `internal/cli/root.go`: 主菜单循环，调用各个子功能
- `internal/cli/source.go`: Source 管理菜单（add, list, remove）
- `internal/sql/mode.go`: Chat 模式实现，使用 readline 库
- `internal/cli/dbconnect.go`: 命令行参数解析和连接验证
- `internal/source/manager.go`: Source 的 CRUD 操作

## Goals / Non-Goals

**Goals:**
- 实现统一的导航机制，支持从 chat 模式返回到主菜单
- 增加 source 编辑功能，支持修改所有字段
- 增加 `/exit`、`/help` 命令和 Tab 补全支持
- 优化 `-D` 参数处理，支持当次执行时使用命令行指定的 database

**Non-Goals:**
- 不支持在 chat 模式中切换 source（需要退出重新选择）
- 不支持复杂的命令历史搜索（readline 已有基础支持）
- 不改变 source 的唯一性规则（仍基于 host-port-username）

## Decisions

### 1. Chat 模式返回主菜单
**决策**: 修改 `RunSQLMode()` 返回值，在主菜单循环中检查返回值，如果返回特定错误（如 `ErrReturnToMenu`），则继续主菜单循环而不是退出程序。

**替代方案考虑**:
- 方案 A: 使用全局状态标志 - 不够清晰，难以维护
- 方案 B: 返回特定错误类型 - 清晰，符合 Go 错误处理习惯 ✓
- 方案 C: 使用 context 传递控制信息 - 过度设计

**实现**: 定义 `var ErrReturnToMenu = errors.New("return to main menu")`，在 chat 模式中需要返回时返回此错误。

### 2. Source 编辑功能
**决策**: 在 `internal/source/manager.go` 中增加 `UpdateSource(name string, updated *Source) error` 函数，在 `internal/cli/source.go` 中增加 `editSource()` 函数和菜单项。

**实现细节**:
- 编辑时允许修改所有字段（name, host, port, database, username, password）
- 如果修改了 name，需要检查新 name 的唯一性
- 如果修改了 host/port/username，需要检查是否与现有 source 冲突（基于唯一性规则）

### 3. Chat 模式命令增强
**决策**: 
- `/exit` 命令：与现有的 `exit`/`back` 文本命令功能相同，统一处理
- `/help` 命令：显示可用命令列表和使用说明
- Tab 补全：使用 readline 的 `SetCompleter()` 功能，提供命令补全（`/exit`, `/help`, `/history`, `/clear`）

**命令解析优先级**:
1. 以 `/` 开头的命令（如 `/exit`, `/help`）
2. 文本命令（如 `exit`, `back`）
3. 自然语言查询

**Tab 补全范围**:
- 命令补全：`/exit`, `/help`, `/history`, `/clear`
- 不补全自然语言查询（避免干扰）

### 4. 命令行 `-D` 参数处理
**决策**: 在 `RunSQLModeWithSource()` 中增加可选参数 `overrideDatabase string`，当提供时，临时覆盖 source 的 database 字段用于连接，但不修改持久化的 source 配置。

**实现细节**:
- `DatabaseArgs` 结构已包含 `Database` 字段
- 在 `internal/cli/dbconnect.go` 或调用处，传递 `Database` 值到 `RunSQLModeWithSource()`
- 在 `RunSQLModeWithSource()` 中，如果 `overrideDatabase` 不为空，创建临时的 source 副本用于连接
- Source 的唯一性仍基于 host-port-username，`-D` 参数不影响 source 的创建和查找逻辑

## Risks / Trade-offs

**[Risk] 命令解析冲突**: `/exit` 可能与用户输入的自然语言查询冲突（如用户输入 "how to exit"）
- **Mitigation**: 严格检查命令格式，只有以 `/` 开头且完全匹配命令名时才视为命令

**[Risk] Tab 补全干扰**: Tab 补全可能干扰用户的自然语言输入
- **Mitigation**: 只在用户输入以 `/` 开头时提供命令补全，其他情况不补全

**[Risk] Source 编辑时的唯一性冲突**: 编辑 source 时修改 host/port/username 可能与现有 source 冲突
- **Mitigation**: 在 `UpdateSource()` 中检查唯一性，如果冲突则返回错误

**[Trade-off] `-D` 参数不持久化**: 用户可能期望 `-D` 参数能更新 source 的 database 字段
- **Rationale**: Source 的唯一性基于 host-port-username，database 是连接参数而非身份标识。用户可以在编辑 source 时修改 database，但命令行参数应该只影响当次执行

## Migration Plan

1. **Phase 1**: 实现 chat 模式返回主菜单功能
   - 修改 `RunSQLMode()` 返回逻辑
   - 修改主菜单循环处理返回值
   - 测试返回功能

2. **Phase 2**: 实现 source 编辑功能
   - 增加 `UpdateSource()` 函数
   - 增加 `editSource()` 菜单项
   - 测试编辑功能

3. **Phase 3**: 实现 chat 模式命令增强
   - 增加 `/exit`、`/help` 命令处理
   - 实现 Tab 补全
   - 测试命令和补全功能

4. **Phase 4**: 优化 `-D` 参数处理
   - 修改 `RunSQLModeWithSource()` 支持 database 覆盖
   - 修改命令行参数传递逻辑
   - 测试 `-D` 参数功能

## Open Questions

1. Tab 补全是否应该支持历史查询补全？（当前决策：不支持，避免干扰）
2. `/help` 命令是否应该显示更详细的使用示例？（当前决策：显示命令列表和简要说明即可）

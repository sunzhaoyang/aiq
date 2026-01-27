## Why

当前 `~/.aiqconfig/` 目录中文件直接平铺在根目录（`config.yaml`、`sources.yaml`、`session_*.json`），随着功能扩展（即将添加 Skills、tools 等），目录会变得混乱。需要在架构层面规划好用户配置目录的组织结构，建立清晰的层次和命名规范，避免后续重构成本和用户数据迁移问题。

同时，项目需要支持 Claude Skills 以增强 AI Agent 的能力，允许用户通过自定义 Skills 扩展功能。Skills 需要存储在用户配置目录中，也需要纳入目录组织规划。

## What Changes

1. **用户配置目录结构重构**：重新组织 `~/.aiqconfig/` 目录结构，建立清晰的层次：
   - 将配置文件按类型分类到子目录（如 `config/`、`sessions/`、`skills/`、`tools/` 等）
   - 建立统一的命名规范和文件组织规则
   - 定义未来扩展的目录组织规范
   - **注意**：项目处于早期阶段，不考虑向后兼容，用户需手动迁移现有文件

2. **Claude Skills 支持**：
   - 实现 Skills 加载机制，支持从 `~/.aiqconfig/skills/` 目录加载用户自定义 Skills
   - 解析 SKILL.md 格式文件（YAML frontmatter + Markdown 内容）
   - **渐进式加载**：根据用户查询和上下文需要，动态加载相关的 Skills，而不是一次性加载所有 Skills
   - **Prompt 管理机制**：
     - 监控 prompt 长度和 token 使用量
     - 实现 prompt 压缩策略（当接近 token 限制时，压缩或淘汰低优先级内容）
     - 建立 Skills 优先级和淘汰机制，避免 prompt 过长影响性能
   - **内置 Tools**：提供基础工具集以支持 Skills 中的常见操作：
     - HTTP 请求工具（支持 Skills 中的 URL 操作）
     - 命令执行工具（支持 Skills 中的命令调用）
     - 文件操作工具（读取、写入、列出文件）
     - 其他基础工具（根据 Skills 需求扩展）
   - 实现类似 Claude Agent SDK 的功能（Go 版本）

3. **配置目录组织规范**：建立 `~/.aiqconfig/` 目录的组织规则和最佳实践文档，确保未来扩展时保持结构清晰

## Capabilities

### New Capabilities

- `user-config-directory-organization`: 重构用户配置目录（`~/.aiqconfig/`）的组织结构，建立清晰的层次和命名规范
- `claude-skills-support`: 支持 Claude Skills 格式，实现 Skills 加载、解析和集成机制

### Modified Capabilities

- `sql-interactive-mode`: 需要集成 Skills 内容到 prompt 中，增强 AI Agent 能力
- `configuration-management`: 可能需要添加 Skills 目录配置项

## Impact

- **用户配置目录结构**：重构 `~/.aiqconfig/` 目录结构
- **配置文件路径**：所有访问配置文件的代码需要更新路径（`config.yaml` → `config/config.yaml`，`sources.yaml` → `config/sources.yaml`，`session_*.json` → `sessions/session_*.json`）
- **配置系统**：需要添加 Skills 相关配置（Skills 目录路径等）
- **LLM 集成**：需要将 Skills 内容整合到 prompt 构建逻辑中，实现渐进式加载和 prompt 管理
- **Prompt 管理**：需要实现 token 监控、prompt 压缩和内容淘汰机制
- **工具系统**：需要扩展 tools 注册机制，提供内置工具集以支持 Skills 中的操作
- **Skills 匹配**：需要实现 Skills 与用户查询的匹配算法，决定何时加载哪些 Skills
- **依赖管理**：可能需要添加 YAML 解析库（如果尚未使用）用于解析 Skills frontmatter
- **文档**：需要更新配置目录结构文档和 Skills 使用指南

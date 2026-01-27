<div align="center">

# AIQ

**一个将自然语言转换为 SQL 查询的智能 SQL 客户端**

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg?style=flat-square)](LICENSE)

*用自然语言提问，获得精确的 SQL 查询，将结果可视化为精美的图表*

</div>

---

## 📖 简介

AIQ (AI Query) 是一个智能 SQL 客户端，通过自然语言与数据库交互。无需编写 SQL，只需用自然语言提问，AIQ 会自动生成 SQL 查询并执行，还能将结果可视化为精美的图表。

### ✨ 核心特性

- 🗣️ **自然语言查询** - 用中文或英文提问，自动生成 SQL
- 💬 **多轮对话** - 保持对话上下文，支持查询优化和后续问题
- 📊 **图表可视化** - 自动检测并渲染图表（柱状图、折线图、饼图、散点图）
- 🔌 **多数据库支持** - [seekdb](https://www.oceanbase.ai/)、MySQL、PostgreSQL
- 🎯 **Skills 系统** - 通过自定义领域知识扩展 AI 能力
- 🎨 **美观的 CLI 界面** - 流畅的交互体验和彩色输出
- 💾 **会话持久化** - 保存和恢复对话会话

## 🚀 快速开始

### 安装

```bash
# 克隆并构建
git clone https://github.com/aiq/aiq.git
cd aiq
go build -o aiq cmd/aiq/main.go

# 安装（可选）
sudo mv aiq /usr/local/bin/
```

### 首次使用

1. **启动 AIQ**: `aiq`
2. **配置 LLM**: 输入 API URL、API Key 和模型名称（首次运行会启动配置向导）
3. **添加数据源**: 选择 `source` → `add` → 输入数据库连接信息
4. **开始查询**: 选择 `chat` → 选择数据源 → 用自然语言提问

**示例查询:**
```
aiq> 显示最近一周的销售额
aiq> 统计每个类别的商品数量
aiq> 查看用户注册趋势
```

## 📚 使用指南

### 主菜单

```
AIQ - Main Menu
? config   - Manage LLM configuration
  source   - Manage database connections
  chat     - Query database with natural language
  exit     - Exit application
```

### 聊天模式

**多轮对话:**
```
aiq> 显示上周的总销售额
[生成 SQL 和结果...]

aiq> 修改为只显示最近 3 天
[AIQ 理解上下文并生成更新的 SQL...]
```

**命令:**
- `/history` - 查看对话历史
- `/clear` - 清除对话历史
- `exit` 或 `back` - 退出聊天模式（会话自动保存）

**恢复会话:**
```bash
aiq -s ~/.aiqconfig/sessions/session_20260126100000.json
```

### 图表可视化

AIQ 会根据查询结果自动检测合适的图表类型：
- **分类 + 数值** → 柱状图或饼图
- **时间 + 数值** → 折线图
- **数值 + 数值** → 散点图

## 🎯 Skills - 扩展 AI 能力

Skills 允许你通过提供自定义指令和上下文来扩展 AIQ 的能力。Skills 会根据你的查询自动匹配和加载。

### 快速开始

1. **创建 Skill 目录:**
```bash
mkdir -p ~/.aiqconfig/skills/my-skill
```

2. **创建 SKILL.md 文件:**
```markdown
---
name: my-skill
description: 针对指标、仪表板和 SQL 模式的领域特定指导
---

# My Custom Skill

此 Skill 提供分析工作流和常见 SQL 模式的指导。

## 核心概念

- 指标和维度的命名规范
- KPI 计算模式和注意事项
- 基于时间的聚合和队列分析

## 使用示例

### 周度 KPI 汇总
```sql
SELECT DATE_TRUNC('week', created_at) AS week,
       COUNT(*) AS orders,
       SUM(amount) AS revenue
FROM orders
GROUP BY week
ORDER BY week;
```
```

3. **重启 AIQ** - Skills 会在启动时自动加载

4. **使用** - 当你查询与 Skill 描述匹配的主题时，它会自动加载

### Skill 文件格式

每个 Skill 必须包含：

- **YAML Frontmatter**（必需）：
  - `name`: Skill 名称（小写，使用连字符，如 `my-skill`）
  - `description`: Skill 描述（最多 200 字符，用于查询匹配）

- **Markdown 内容**: 指令、示例和指导

### 工作原理

1. **启动时**: AIQ 从 `~/.aiqconfig/skills/` 加载所有 Skills 的元数据（name, description）
2. **查询时**: 系统提取关键词并与 Skills 元数据匹配
3. **自动加载**: 最相关的 Top 3 Skills 会被加载到 prompt 中
4. **智能压缩**: 系统自动管理 prompt 长度（压缩历史记录，淘汰低优先级 Skills）

### 匹配规则

Skills 基于相关性评分进行匹配：
- **精确名称匹配**（最高优先级）：查询完全匹配 Skill 名称
- **部分名称匹配**：查询包含 Skill 名称或反之
- **描述关键词匹配**：查询关键词出现在 Skill 描述中

### 推荐 Skills

- **[seekdb Skill](https://github.com/oceanbase/seekdb-ecology-plugins/blob/main/claudecode-plugin/skills/seekdb/SKILL.md)** - SeekDB 文档目录和使用指导

### 内置工具

Skills 可以在其指令中使用以下内置工具：

- **`execute_sql`** - 执行数据库 SQL 查询
- **`http_request`** - 发起 HTTP 请求（GET, POST, PUT, DELETE）
- **`execute_command`** - 执行 shell 命令（有安全白名单限制）
- **`file_operations`** - 读写文件（限制在安全目录内）

**注意**: Skills 是上下文信息，不是工具本身。它们指导 AI 如何使用内置工具。

### Prompt 管理

系统自动管理 prompt 长度：
- **80% 阈值**: 压缩对话历史（保留最近 10 条消息）
- **90% 阈值**: 淘汰低优先级 Skills（保留活跃和相关的）
- **95% 阈值**: 激进压缩（只保留最近 5 条消息和 Top Skills）

### 目录结构

Skills 存储在 `~/.aiqconfig/skills/<skill-name>/SKILL.md`:

```
~/.aiqconfig/
└── skills/
    ├── my-skill/
    │   └── SKILL.md
    └── data-analysis/
        └── SKILL.md
```

**注意**: 每个 Skill 目录只包含一个 `SKILL.md` 文件。如果需要多个文件，将内容合并到一个文件中，或拆分为多个更小的 Skills。

### 故障排除

**Skills 未加载:**
- 检查目录结构: `~/.aiqconfig/skills/<skill-name>/SKILL.md`
- 验证 YAML frontmatter 格式（必须以 `---` 开始和结束）
- 确保 `name` 和 `description` 字段存在
- 查看启动日志中的错误

**Skills 未匹配:**
- 在 Skill `description` 中包含相关关键词
- 尝试在查询中使用 Skill 名称
- 检查是否有多个 Skills 竞争（只选择 Top 3）

## ⚙️ 配置

配置文件存储在 `~/.aiqconfig/`:

- **config/config.yaml** - LLM 配置（URL、API Key、模型）
- **config/sources.yaml** - 数据库连接配置
- **sessions/** - 对话会话文件（自动生成）
- **skills/** - 自定义 Skills（见上方 Skills 部分）

**示例 config.yaml:**
```yaml
llm:
  url: https://api.openai.com/v1
  apiKey: sk-...
  model: gpt-4
```

**示例 sources.yaml:**
```yaml
sources:
  - name: local-mysql
    type: MySQL
    host: localhost
    port: 3306
    database: testdb
    username: root
    password: password
```

## 🛠️ 开发

### 项目结构

```
aiq/
├── cmd/aiq/          # 主程序入口
├── internal/
│   ├── cli/          # CLI 命令和菜单系统
│   ├── config/       # 配置管理
│   ├── source/       # 数据源管理
│   ├── sql/          # SQL 交互模式（chat 模式）
│   ├── skills/       # Skills 系统（匹配、加载、管理）
│   ├── prompt/       # Prompt 构建和压缩
│   ├── llm/          # LLM 客户端集成
│   ├── db/           # 数据库连接和查询执行
│   ├── chart/        # 图表可视化
│   ├── tool/         # 工具系统（内置工具）
│   └── ui/           # UI 组件
└── openspec/         # OpenSpec 变更管理
```

### 构建

```bash
go build -o aiq cmd/aiq/main.go
```

### 运行测试

```bash
go test ./...
```

## 📝 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

---

<div align="center">

**Made with ❤️ using Go**

[报告问题](https://github.com/aiq/aiq/issues) · [提交功能请求](https://github.com/aiq/aiq/issues) · [查看文档](https://github.com/aiq/aiq)

</div>

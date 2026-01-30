## Why

仓库路径已从 `sunzhaoyang/aiq` 迁移到 `sunetic/aiq`。需要更新项目中所有引用旧仓库路径的地方，确保安装脚本、文档和代码中的仓库引用指向新的正确路径。

## What Changes

- 更新安装脚本中的仓库路径引用（`scripts/install.sh` 和 `scripts/install.bat`）
- 更新 README 文档中的仓库 URL（`README.md` 和 `README_CN.md`）
- 检查并更新代码中可能存在的仓库路径引用
- **注意**: `go.mod` 中的模块路径通常应保持不变，除非确实需要迁移 Go 模块路径

## Capabilities

### New Capabilities
<!-- 无新能力 -->

### Modified Capabilities
- `installation-script`: 安装脚本需要更新 GitHub 仓库路径引用，从 `sunzhaoyang/aiq` 改为 `sunetic/aiq`

## Impact

**受影响的文件**:
- `scripts/install.sh` - Unix/Linux/macOS 安装脚本
- `scripts/install.bat` - Windows 安装脚本
- `README.md` - 英文文档
- `README_CN.md` - 中文文档
- 可能存在的其他文档或注释中的仓库引用

**影响范围**:
- 用户安装体验：确保安装脚本能从正确的仓库下载二进制文件
- 文档准确性：确保文档中的链接和示例指向正确的仓库
- 开发体验：确保开发者能找到正确的仓库地址

**注意事项**:
- `go.mod` 中的模块路径 `github.com/aiq/aiq` 可能需要评估是否需要更新
- 需要检查是否有其他配置文件或文档引用了旧路径

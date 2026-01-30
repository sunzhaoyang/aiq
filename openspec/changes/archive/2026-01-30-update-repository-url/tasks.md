## 1. 更新安装脚本

- [x] 1.1 更新 `scripts/install.sh` 中的 `REPO` 变量从 `sunzhaoyang/aiq` 改为 `sunetic/aiq`
- [x] 1.2 更新 `scripts/install.bat` 中的 `REPO` 变量从 `sunzhaoyang/aiq` 改为 `sunetic/aiq`
- [x] 1.3 验证安装脚本语法正确性（bash 和 batch 脚本）

## 2. 更新文档

- [x] 2.1 更新 `README.md` 第 40 行的安装命令 URL：`https://raw.githubusercontent.com/sunzhaoyang/aiq/main/scripts/install.sh` → `https://raw.githubusercontent.com/sunetic/aiq/main/scripts/install.sh`
- [x] 2.2 更新 `README.md` 第 47 行的 Windows 安装命令 URL：`https://raw.githubusercontent.com/sunzhaoyang/aiq/main/scripts/install.bat` → `https://raw.githubusercontent.com/sunetic/aiq/main/scripts/install.bat`
- [x] 2.3 更新 `README.md` 第 82 行的 git clone 命令：`https://github.com/sunzhaoyang/aiq.git` → `https://github.com/sunetic/aiq.git`
- [x] 2.4 更新 `README.md` 第 186 行的 issue 链接：`https://github.com/aiq/aiq/issues` → `https://github.com/sunetic/aiq/issues`（如果存在）
- [x] 2.5 检查并更新 `README_CN.md` 中的相关 URL（如果存在 `sunzhaoyang/aiq` 或需要更新的链接）
- [x] 2.6 验证所有链接格式正确，确保 URL 可访问

## 3. 全面检查

- [x] 3.1 使用 grep 搜索所有 `sunzhaoyang` 引用，确认没有遗漏的文件
- [x] 3.2 使用 grep 搜索所有 GitHub URL 引用（`raw.githubusercontent.com`, `github.com/.*/releases`, `github.com/.*/blob`, `github.com/.*/tree`），确认都已更新
- [x] 3.3 确认归档文档（`openspec/changes/archive/`）保持不变
- [x] 3.4 确认 `go.mod` 中的模块路径 `github.com/aiq/aiq` 保持不变（符合设计决策）

## 4. 测试验证

- [x] 4.1 验证 `scripts/install.sh` 语法正确（使用 `bash -n scripts/install.sh`）
- [x] 4.2 验证 `scripts/install.bat` 语法正确（在 Windows 环境中测试，或使用语法检查工具）
- [x] 4.3 验证文档中的链接格式正确（检查 Markdown 链接语法）
- [x] 4.4 确认代码构建正常（运行 `go build` 验证模块路径未改变不影响构建）
- [x] 4.5 如果新仓库已有 releases，验证安装脚本能从新仓库正确下载（可选，取决于新仓库状态）

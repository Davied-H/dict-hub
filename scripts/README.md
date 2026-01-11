# Dict-Hub 快捷脚本

本目录包含 Raycast 和 Alfred 的快捷脚本，可以让你快速查词而无需打开浏览器。

## 前置要求

1. **Dict-Hub 服务运行中**
   ```bash
   # 在项目根目录运行
   ./start.sh start
   ```

2. **安装 jq**（JSON 处理工具）
   ```bash
   brew install jq
   ```

## Raycast 脚本

### 安装方法

1. 打开 Raycast 设置 → Extensions → Script Commands
2. 点击 "Add Script Directory"
3. 选择 `scripts/raycast` 目录
4. 脚本会自动出现在 Raycast 中

### 可用命令

| 命令 | 功能 | 快捷键建议 |
|------|------|-----------|
| **Dict Search** | 查询单词释义 | `⌘ + Shift + D` |
| **Dict Suggest** | 获取单词建议/自动补全 | `⌘ + Shift + S` |
| **Dict History** | 查看搜索历史 | `⌘ + Shift + H` |

### 使用示例

1. 打开 Raycast（默认 `⌥ + Space`）
2. 输入 "Dict Search" 或设置的快捷键
3. 输入要查询的单词，按回车查询
4. 查看释义预览，点击链接在浏览器中打开完整页面

### 环境变量配置

可以通过环境变量自定义 API 地址：

```bash
export DICT_HUB_API="http://localhost:8080"
export DICT_HUB_FRONTEND="http://localhost:3000"
```

---

## Alfred Workflow

### 安装方法

1. 双击 `alfred/Dict-Hub.alfredworkflow` 文件夹（或将其打包为 `.alfredworkflow` 文件）
2. Alfred 会自动导入 Workflow
3. 在 Alfred Preferences → Workflows 中配置环境变量（可选）

**手动打包方法：**
```bash
cd scripts/alfred
zip -r Dict-Hub.alfredworkflow.zip Dict-Hub.alfredworkflow
mv Dict-Hub.alfredworkflow.zip Dict-Hub.alfredworkflow
```

### 可用关键字

| 关键字 | 功能 | 说明 |
|--------|------|------|
| `dict [word]` | 查词 + 自动补全 | 输入时显示建议，回车打开浏览器 |
| `dicth` | 搜索历史 | 显示最近搜索的单词 |
| `dictf` | 词频统计 | 显示最常搜索的单词 |

### 使用示例

1. 打开 Alfred（默认 `⌘ + Space`）
2. 输入 `dict hello`
3. 在下拉列表中选择建议的单词
4. 按 Enter 在浏览器中打开
5. 按 `⌘ + Enter` 复制单词到剪贴板

### Workflow 配置

在 Alfred Preferences → Workflows → Dict-Hub → Configure 中设置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DICT_HUB_API` | `http://localhost:8080` | API 服务地址 |
| `DICT_HUB_FRONTEND` | `http://localhost:3000` | 前端页面地址 |

---

## 目录结构

```
scripts/
├── README.md                              # 本文档
├── raycast/                               # Raycast 脚本
│   ├── dict-search.sh                     # 查词脚本
│   ├── dict-suggest.sh                    # 自动补全脚本
│   └── dict-history.sh                    # 历史记录脚本
└── alfred/
    └── Dict-Hub.alfredworkflow/           # Alfred Workflow
        ├── info.plist                     # Workflow 配置
        ├── suggest.sh                     # 查词 + 建议脚本
        ├── history.sh                     # 历史记录脚本
        └── wordfreq.sh                    # 词频统计脚本
```

## 常见问题

### Q: 脚本无法连接到 API？
A: 确保 Dict-Hub 服务正在运行：
```bash
./start.sh status
# 或
curl http://localhost:8080/api/v1/health
```

### Q: 提示 "jq is not installed"？
A: 安装 jq：
```bash
brew install jq
```

### Q: 如何修改默认端口？
A: 设置环境变量：
- Raycast：在系统环境变量或 `.zshrc` 中设置
- Alfred：在 Workflow 配置中设置

### Q: Alfred Workflow 不显示结果？
A: 检查脚本权限：
```bash
chmod +x scripts/alfred/Dict-Hub.alfredworkflow/*.sh
```

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这些脚本！

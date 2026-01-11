# Dict-Hub

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react" alt="React Version" />
  <img src="https://img.shields.io/badge/TypeScript-5.9-3178C6?style=flat-square&logo=typescript" alt="TypeScript Version" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License" />
</p>

<p align="center">
  一个现代化的 MDX 词典查询平台，支持多词典管理、快速查词与历史记录追踪。
</p>

---

## ✨ 特性

- 🔍 **快速查词** - 高性能词条检索，毫秒级响应
- 📚 **多词典支持** - 支持加载多个 MDX 格式词典，灵活切换
- 📊 **词频统计** - 自动记录查词频率，助力词汇学习
- 📝 **历史记录** - 完整的查询历史追踪与管理
- 🎨 **现代 UI** - 简洁优雅的用户界面，支持深色模式
- 🚀 **高性能** - Go 后端 + React 前端，极致性能体验

## 🏗️ 技术栈

### 后端
- **Go 1.24** - 高性能后端语言
- **Gin** - 轻量级 Web 框架
- **GORM** - ORM 框架
- **SQLite** - 轻量级数据库
- **MDX Parser** - MDX 词典格式解析

### 前端
- **React 19** - 前端框架
- **TypeScript** - 类型安全
- **Vite** - 下一代构建工具
- **Tailwind CSS** - 原子化 CSS 框架
- **HeroUI** - 现代 UI 组件库
- **Zustand** - 轻量级状态管理
- **React Query** - 数据获取与缓存

## 📦 快速开始

### 环境要求

- Go 1.24+
- Node.js 18+
- npm 或 pnpm

### 安装与运行

1. **克隆项目**
```bash
git clone https://github.com/your-username/dict-hub.git
cd dict-hub
```

2. **一键启动**
```bash
chmod +x start.sh
./start.sh start
```

3. **访问应用**
- 前端页面: http://localhost:3000
- 后端 API: http://localhost:8080

### 手动启动

**启动后端**
```bash
cd backend
go run cmd/server/main.go
```

**启动前端**
```bash
cd frontend
npm install
npm run dev
```

## 📁 项目结构

```
dict-hub/
├── backend/                 # Go 后端
│   ├── cmd/server/         # 应用入口
│   ├── configs/            # 配置文件
│   ├── dicts/              # MDX 词典文件目录
│   ├── internal/           # 内部模块
│   │   ├── cache/          # 缓存层
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库初始化
│   │   ├── handler/        # HTTP 处理器
│   │   ├── middleware/     # 中间件
│   │   ├── model/          # 数据模型
│   │   ├── router/         # 路由定义
│   │   └── service/        # 业务逻辑
│   ├── pkg/                # 公共包
│   │   ├── mdict/          # MDX 解析器
│   │   └── response/       # 响应封装
│   └── thirdparty/         # 第三方库
│
├── frontend/               # React 前端
│   ├── src/
│   │   ├── api/            # API 客户端
│   │   ├── components/     # UI 组件
│   │   ├── hooks/          # 自定义 Hooks
│   │   ├── pages/          # 页面组件
│   │   ├── stores/         # 状态管理
│   │   └── types/          # 类型定义
│   └── ...
│
├── start.sh                # 一键启动脚本
└── README.md
```

## 🔧 配置

后端配置文件位于 `backend/configs/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug

database:
  path: ./data/dict-hub.db

mdx:
  dict_dir: ./dicts
```

## 📖 API 文档

### 字典查询

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/v1/dicts` | GET | 获取已加载的词典列表 |
| `/api/v1/dicts/:name/lookup` | GET | 查询单词释义 |
| `/api/v1/dicts/:name/suggest` | GET | 获取词条建议 |

### 词典管理

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/v1/dictionaries` | GET | 获取所有词典源 |
| `/api/v1/dictionaries/:id/enable` | POST | 启用词典 |
| `/api/v1/dictionaries/:id/disable` | POST | 禁用词典 |

### 历史记录

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/v1/history` | GET | 获取查询历史 |
| `/api/v1/history` | DELETE | 清空历史记录 |

### 词频统计

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/v1/wordfreq` | GET | 获取词频统计 |
| `/api/v1/wordfreq/top` | GET | 获取高频词汇 |

## 🎯 使用指南

### 添加词典

1. 将 `.mdx` 格式的词典文件放入 `backend/dicts/` 目录
2. 重启服务或通过 API 重新加载词典
3. 在前端界面中启用词典

### 支持的词典格式

- MDX (MDict Dictionary)
- 配套的 MDD 资源文件（可选）

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [MDict](https://www.mdict.cn/) - MDX 词典格式
- [Gin](https://gin-gonic.com/) - Go Web 框架
- [React](https://react.dev/) - 前端框架
- [HeroUI](https://heroui.com/) - UI 组件库

---

<p align="center">
  Made with ❤️ by Dict-Hub Team
</p>
